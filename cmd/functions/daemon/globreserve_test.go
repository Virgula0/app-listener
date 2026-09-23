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
