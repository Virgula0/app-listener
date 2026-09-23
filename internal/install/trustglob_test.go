package install

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/guard"
)

func TestTrustGlobs_SplitsFixedRootAndName(t *testing.T) {
	e := CandidateDir{
		Whitelist: map[string][]string{
			"%HOME%/.local/share/Steam/steamapps/common/*/files/bin/wineserver": nil,
			"%HOME%/.local/bin/claude": nil, // fixed path: trusted on first use, not reserved
			"/opt/*/bin/fsnotifier":    nil, // root-owned tree
		},
		LibDirWriters: []string{
			"%HOME%/.local/share/Steam/steamapps/common/*/files/bin/wineserver", // duplicate
			"%HOME%/.local/share/Steam/steamrt64/pv-runtime/*/pressure-vessel/bin/pressure-vessel-*",
		},
	}
	got := e.TrustGlobs("u", "/home/u")
	want := []TrustGlob{
		{Home: "/home/u", Fixed: []string{".local", "share", "Steam", "steamapps", "common"}, Name: "wineserver"},
		{Home: "/home/u", Fixed: []string{".local", "share", "Steam", "steamrt64", "pv-runtime"}, Name: "pressure-vessel-*"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("TrustGlobs:\n got %+v\nwant %+v", got, want)
	}
	if r := got[0].Root(); r != "/home/u/.local/share/Steam/steamapps/common" {
		t.Fatalf("Root = %s", r)
	}
}

// Every catalog glob that grants trust from outside its own guarded tree must be enforceable by
// the trust object, and the whole set must fit its fixed-size maps.
func TestCatalogTrustGlobsAreReservable(t *testing.T) {
	const home = "/home/u"
	names := make(map[string]bool)
	affixes := make(map[string]bool)
	for i := range Catalog {
		e := &Catalog[i]
		for _, g := range e.TrustGlobs("u", home) {
			if inOwnTree(e, g.Root()) {
				continue
			}
			if _, _, ok := guard.ParseGlobName(g.Name); !ok {
				t.Errorf("%s: %s/%s is neither exact, prefix* nor *suffix", e.Name, g.Root(), g.Name)
			}
			for _, c := range g.Fixed {
				if _, _, ok := guard.ParseGlobName(c); !ok {
					t.Errorf("%s: fixed component %q of %s is not reservable", e.Name, c, g.Root())
				}
			}
			names[g.Name] = true
			if strings.HasPrefix(g.Name, "*") {
				affixes["s"+string(rune(len(g.Name)-1))] = true
			} else if strings.HasSuffix(g.Name, "*") {
				affixes["p"+string(rune(len(g.Name)-1))] = true
			}
		}
	}
	if len(names) > 64 {
		t.Errorf("%d reserved names; guard_trust.bpf.c has 64 bits", len(names))
	}
	if len(affixes) > 8 {
		t.Errorf("%d prefix/suffix shapes; guard_trust.bpf.c probes at most 8", len(affixes))
	}
	if !names["wineserver"] || !names["Discord"] {
		t.Errorf("expected the Steam and Discord globs to be reserved, got %v", names)
	}
}

// inOwnTree: root lies in a tree the entry guards. With WatchRelPaths only those subtrees are
// guarded; the RelPath is just their encryption root.
func inOwnTree(e *CandidateDir, root string) bool {
	guarded := e.PathsFor("/home/u", "u")
	if len(e.WatchRelPaths) > 0 {
		guarded = nil
		for _, rel := range e.WatchRelPaths {
			guarded = append(guarded, filepath.Join("/home/u", rel))
		}
	}
	for _, p := range guarded {
		if ok, _ := filepath.Match(p, root); ok || strings.HasPrefix(root, p+"/") {
			return true
		}
	}
	return false
}

func TestOwnsPath(t *testing.T) {
	var steam, discord *CandidateDir
	for i := range Catalog {
		switch Catalog[i].Name {
		case "Steam":
			steam = &Catalog[i]
		case "Discord":
			discord = &Catalog[i]
		}
	}
	cases := []struct {
		e    *CandidateDir
		path string
		want bool
	}{
		{steam, "/home/u/.local/share/Steam/config", true},
		{steam, "/home/u/.local/share/Steam/userdata/42/config/localconfig.vdf", true},
		{steam, "/home/u/.local/share/Steam/steamapps/common/SteamLinuxRuntime_sniper", true}, // lib_dir
		{steam, "/home/u/.local/share/Steam/steamapps/common/evil", false},
		{steam, "/home/other/.local/share/Steam/config", false},
		{discord, "/home/u/.config/discord", true},
		{discord, "/home/u/.config/discord/Local State", true},
	}
	for _, c := range cases {
		if got := c.e.OwnsPath("u", "/home/u", c.path); got != c.want {
			t.Errorf("%s OwnsPath(%s) = %v, want %v", c.e.Name, c.path, got, c.want)
		}
	}
}
