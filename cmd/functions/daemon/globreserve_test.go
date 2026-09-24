package daemon

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/install"
)

func mkdirs(t *testing.T, paths ...string) {
	t.Helper()
	for _, p := range paths {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func bitOf(t *testing.T, r guard.GlobReservations, name string) uint64 {
	t.Helper()
	i := slices.Index(r.Patterns, name)
	if i < 0 {
		t.Fatalf("pattern %q not reserved; have %v", name, r.Patterns)
	}
	return 1 << uint(i)
}

func TestBuildGlobReservations_Steam(t *testing.T) {
	home := t.TempDir()
	steam := filepath.Join(home, ".local/share/Steam")
	common := filepath.Join(steam, "steamapps/common")
	// compatibilitytools.d is missing: its globs must fall back to the deepest existing dir.
	mkdirs(t, filepath.Join(steam, "config"), common)
	client := filepath.Join(steam, "ubuntu12_32/steam")
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(steam, "config"), Binaries: []daemonconfig.BinaryRule{{Path: client}}},
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}
	r := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	ws := bitOf(t, r, "wineserver")
	if r.Roots[common]&ws == 0 {
		t.Errorf("wineserver not reserved below %s: %v", common, r.Roots)
	}
	if r.Roots[steam]&ws == 0 {
		t.Errorf("compatibilitytools.d/*/files/bin/wineserver must fall back to %s: %v", steam, r.Roots)
	}
	if r.Roots[steam]&bitOf(t, r, "steam") == 0 {
		t.Errorf("*/steam not reserved below %s", steam)
	}
	if r.Writers[client]&ws == 0 {
		t.Errorf("the Steam client must be a writer of wineserver: %v", r.Writers)
	}
	if _, ok := r.Writers["/usr/bin/ssh"]; ok {
		t.Errorf("another entry's binary must not write Steam's reserved names")
	}
	wantChild := guard.GlobChild{Parent: filepath.Join(steam, "steamapps"), Name: "common"}
	found := false
	for _, c := range r.Children {
		if c.Parent == wantChild.Parent && c.Name == wantChild.Name && c.Bits&ws != 0 {
			found = true
		}
	}
	if !found {
		t.Errorf("existing root %s must be reserved in its parent: %+v", common, r.Children)
	}
}

func TestBuildGlobReservations_UnconfiguredEntryReservesNothing(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, ".local/share/Steam/steamapps/common"))
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}
	r := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	if len(r.Patterns) != 0 || len(r.Roots) != 0 || len(r.Writers) != 0 {
		t.Fatalf("Steam is not configured, so nothing may be reserved: %+v", r)
	}
}

// A glob whose root lies inside a guarded resource is already write-gated by that resource.
func TestBuildGlobReservations_InTreeRootSkipped(t *testing.T) {
	home := t.TempDir()
	foundry := filepath.Join(home, ".foundry")
	mkdirs(t, filepath.Join(foundry, "bin"))
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: foundry, Binaries: []daemonconfig.BinaryRule{{Path: filepath.Join(foundry, "bin/forge")}}},
	}}
	r := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	if len(r.Patterns) != 0 {
		t.Fatalf("foundry's bin/* is inside its guarded tree, nothing to reserve: %v", r.Patterns)
	}
}

func discordLibConfig(home string) (*daemonconfig.Config, string) {
	bin := filepath.Join(home, ".config/discord/0.0.1/Discord")
	return &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".config/discord/sentry"), Binaries: []daemonconfig.BinaryRule{{Path: bin}}},
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}, bin
}

func TestBuildGlobReservations_DiscordLibs(t *testing.T) {
	home := t.TempDir()
	discord := filepath.Join(home, ".config/discord")
	mkdirs(t, discord)
	cfg, bin := discordLibConfig(home)
	r := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	for _, name := range []string{"*.so", "lib*", "*.node"} {
		bit := bitOf(t, r, name)
		if r.Roots[discord]&bit == 0 {
			t.Errorf("%s not reserved below %s: %v", name, discord, r.Roots)
		}
		if r.Writers[bin]&bit == 0 {
			t.Errorf("Discord must be a writer of %s: %v", name, r.Writers)
		}
		if r.Writers["/usr/bin/ssh"]&bit != 0 {
			t.Errorf("another entry's binary must not write (or load) Discord's %s", name)
		}
	}
	lib := bitOf(t, r, "lib*")
	found := false
	for _, c := range r.Children {
		if c.Parent == filepath.Join(home, ".config") && c.Name == "discord" && c.Bits&lib != 0 {
			found = true
		}
	}
	if !found {
		t.Errorf("the library root must be reserved in its parent: %+v", r.Children)
	}
}

// Falling back to ~/.config would reserve lib* for every app there.
func TestBuildGlobReservations_MissingLibRootNotWidened(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, ".config"))
	cfg, _ := discordLibConfig(home)
	r := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	lib := bitOf(t, r, "lib*")
	for root, bits := range r.Roots {
		if bits&lib != 0 {
			t.Errorf("lib* reserved below %s although %s/.config/discord is missing", root, home)
		}
	}
}

func TestGlobBuilder_LibBitsPerEntry(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, "a"), filepath.Join(home, "b"))
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".a"), Binaries: []daemonconfig.BinaryRule{{Path: "/bin/a"}}},
		{Path: filepath.Join(home, ".b"), Binaries: []daemonconfig.BinaryRule{{Path: "/bin/b"}}},
	}}
	b := globBuilder{
		cfg:      cfg,
		r:        guard.GlobReservations{Roots: map[string]uint64{}, Writers: map[string]uint64{}},
		bitOf:    map[string]int{},
		children: map[[2]string]uint64{},
	}
	u := install.User{Name: "u", Home: home}
	b.addEntry(&install.CandidateDir{Name: "A", RelPaths: []string{".a"}, ReservedLibs: []string{"%HOME%/a/*.so"}}, u)
	b.addEntry(&install.CandidateDir{Name: "B", RelPaths: []string{".b"}, ReservedLibs: []string{"%HOME%/b/*.so"}}, u)

	if len(b.r.Patterns) != 2 {
		t.Fatalf("each entry's *.so needs its own bit: %v", b.r.Patterns)
	}
	if b.r.Writers["/bin/a"]&b.r.Roots[filepath.Join(home, "b")] != 0 {
		t.Errorf("A's writer may plant or load *.so below B's root: writers %v roots %v", b.r.Writers, b.r.Roots)
	}
}

// A lib_dir glob match is guarded by the next catalog refresh with whatever it already holds, so
// its wildcard component must be reserved for the entry's writers like a whitelist glob.
func TestBuildGlobReservations_SteamLibDirGlobs(t *testing.T) {
	home := t.TempDir()
	steam := filepath.Join(home, ".local/share/Steam")
	common := filepath.Join(steam, "steamapps/common")
	mkdirs(t, filepath.Join(steam, "config"), common)
	client := filepath.Join(steam, "ubuntu12_32/steam")
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(steam, "config"), Binaries: []daemonconfig.BinaryRule{{Path: client}}},
	}}
	r := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	for _, name := range []string{"Proton*", "SteamLinuxRuntime_*"} {
		bit := bitOf(t, r, name)
		if r.Roots[common]&bit == 0 {
			t.Errorf("%s not reserved below %s: %v", name, common, r.Roots)
		}
		if r.Writers[client]&bit == 0 {
			t.Errorf("the Steam client must be a writer of %s: %v", name, r.Writers)
		}
	}
}
