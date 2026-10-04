package daemon

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func rules(paths ...string) []daemonconfig.BinaryRule {
	out := make([]daemonconfig.BinaryRule, len(paths))
	for i, p := range paths {
		out[i] = daemonconfig.BinaryRule{Path: p}
	}
	return out
}

func TestPlanUpdaters_ScopedToOwnResources(t *testing.T) {
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: "/protected", Binaries: rules("/home/u/.local/bin/app", "/usr/bin/git"), AllowLibs: []string{"/home/u/lib.so"}},
		{Path: "/other", Binaries: rules("/opt/other/cp", "/opt/other/tool")},
		{Path: "/same", Binaries: rules("/usr/bin/git", "/home/u/.local/bin/app")},
	}}
	p := planUpdaters(cfg)

	app := p.Owners["/home/u/.local/bin/app"]
	if app == 0 || p.Owners["/home/u/lib.so"] != app {
		t.Fatalf("the resource's binaries and allow_libs must share its bit: %v", p.Owners)
	}
	if p.Updaters["/home/u/.local/bin/app"]&app == 0 {
		t.Errorf("a binary must update its own resource: %v", p.Updaters)
	}
	if p.Updaters["/opt/other/tool"]&app != 0 {
		t.Errorf("another resource's binary must not update app: %v", p.Updaters)
	}
	for _, tool := range []string{"/usr/bin/git", "/opt/other/cp"} {
		if p.Updaters[tool] != 0 {
			t.Errorf("general tool %s must never be an updater: %v", tool, p.Updaters)
		}
	}
	if p.Owners["/usr/bin/git"] != app {
		t.Errorf("identical binary sets must share one bit: %v", p.Owners)
	}
}

// AllowReplacement judges a re-pointed link's new target by the link: it must own the bit.
func TestPlanUpdaters_LinkOwnsItsResourceBit(t *testing.T) {
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{{Path: "/protected", Binaries: []daemonconfig.BinaryRule{
		{Path: "/home/u/.local/share/app/versions/2", Link: "/home/u/.local/bin/app"},
	}}}}
	p := planUpdaters(cfg)
	bit := p.Owners["/home/u/.local/share/app/versions/2"]
	if bit == 0 || p.Owners["/home/u/.local/bin/app"] != bit {
		t.Fatalf("the link must own its target's bit: %v", p.Owners)
	}
	if _, ok := p.Updaters["/home/u/.local/bin/app"]; ok {
		t.Errorf("the link adds no updater row (it resolves to the target): %v", p.Updaters)
	}
}

func TestIsGeneralTool(t *testing.T) {
	for path, want := range map[string]bool{
		"/usr/bin/python3.12": true, "/usr/bin/node": true, "/bin/sh": true,
		"/home/u/.local/bin/claude": false, "/opt/steam/steam": false,
	} {
		if got := !mayUpdate(path, false); got != want {
			t.Errorf("general tool %s = %v, want %v", path, got, want)
		}
	}
}

// Updater rights attach to the inode: a hard-linked multi-call binary or a symlink to busybox is a
// general tool whatever applet name is whitelisted.
func TestIsGeneralToolJudgesTheInode(t *testing.T) {
	dir := t.TempDir()
	multi := filepath.Join(dir, "date")
	single := filepath.Join(dir, "app")
	for _, p := range []string{multi, single, filepath.Join(dir, "busybox")} {
		if err := os.WriteFile(p, []byte("ELF"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Link(multi, filepath.Join(dir, "cp")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "busybox"), filepath.Join(dir, "true")); err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]bool{
		multi: true, filepath.Join(dir, "true"): true, single: false,
		filepath.Join(dir, "missing"): false,
	} {
		if got := !mayUpdate(path, false); got != want {
			t.Errorf("general tool %s = %v, want %v", path, got, want)
		}
	}
}

// A hard link keeps a binary out of the updaters but only a tool's name warns: Steam's runtime
// hard-links every pressure-vessel binary, which flooded the journal.
func TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(t *testing.T) {
	dir := t.TempDir()
	app, linked, git := filepath.Join(dir, "app"), filepath.Join(dir, "pv-adverb"), filepath.Join(dir, "git")
	for _, p := range []string{app, linked, git} {
		if err := os.WriteFile(p, []byte("ELF"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Link(linked, filepath.Join(dir, "copy")); err != nil {
		t.Fatal(err)
	}

	if got := generalToolsToWarn([]string{app, linked}); len(got) != 0 {
		t.Errorf("hard-linked non-tool must not warn: %v", got)
	}
	if got := generalToolsToWarn([]string{app, linked, git}); len(got) != 1 || got[0] != git {
		t.Errorf("named tool beside a user-writable binary must warn alone: %v", got)
	}
	if got := generalToolsToWarn([]string{git}); len(got) != 0 {
		t.Errorf("no user-writable sibling, no warning: %v", got)
	}

	p := planUpdaters(&daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: "/protected", Binaries: rules(app, linked)},
	}})
	if p.Updaters[linked] != 0 || p.Updaters[app] == 0 {
		t.Errorf("hard-linked binary must stay out of the updaters: %v", p.Updaters)
	}
}

// The hard-link waiver is for a user-writable lib_binary only; a tool's name is never waived.
func TestMayUpdate_HardLinkWaivedForUserWritableLibBinary(t *testing.T) {
	dir := t.TempDir()
	pv, cp := filepath.Join(dir, "pv-adverb"), filepath.Join(dir, "cp")
	for _, p := range []string{pv, cp} {
		if err := os.WriteFile(p, []byte("ELF"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Link(p, p+".link"); err != nil {
			t.Fatal(err)
		}
	}
	for _, c := range []struct {
		path      string
		libBinary bool
		want      bool
	}{
		{pv, true, true},
		{pv, false, false},
		{cp, true, false},
	} {
		if got := mayUpdate(c.path, c.libBinary); got != c.want {
			t.Errorf("mayUpdate(%s, libBinary=%v) = %v, want %v", c.path, c.libBinary, got, c.want)
		}
	}

	p := planUpdaters(&daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: "/lib", ReadOnly: true, Binaries: []daemonconfig.BinaryRule{{Path: pv, LibBinary: true}}},
		{Path: "/secret", Binaries: rules(pv, filepath.Join(dir, "app"))},
	}})
	lib := p.Owners[pv] &^ p.Owners[filepath.Join(dir, "app")]
	if p.Updaters[pv] != lib {
		t.Errorf("hard-linked lib_binary must update only its lib_dir set: updaters=%b lib=%b", p.Updaters[pv], lib)
	}
}
