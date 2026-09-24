package guard

import "testing"

func TestPlanTrustedDirs_LibDirsOnlyForTheirWriters(t *testing.T) {
	roots, loaders, refused := planTrustedDirs([]TrustedDir{
		{Path: "/h/.ssh"},
		{Path: "/h/steam/linux64", Loaders: []string{"/h/steam/steam", "/h/steam/pv"}},
		{Path: "/h/steam/steamrt64", Loaders: []string{"/h/steam/pv", "/h/steam/steam"}},
		{Path: "/h/other/lib", Loaders: []string{"/h/other/app"}},
		{Path: "/h/empty", Loaders: []string{}},
	})
	if len(refused) != 0 {
		t.Fatalf("refused %v", refused)
	}
	if roots["/h/.ssh"] != trustedDirAny {
		t.Errorf("a whitelist-mode root is trusted for every whitelisted binary: %#x", roots["/h/.ssh"])
	}
	steam := roots["/h/steam/linux64"]
	if steam == 0 || steam&trustedDirAny != 0 || roots["/h/steam/steamrt64"] != steam {
		t.Errorf("lib_dirs with the same writers share one bit: %#x %#x", steam, roots["/h/steam/steamrt64"])
	}
	if loaders["/h/steam/steam"]&steam == 0 || loaders["/h/other/app"]&steam != 0 {
		t.Errorf("only a dir's own writers carry its bit: %v", loaders)
	}
	if roots["/h/empty"] == 0 {
		// An empty loader set still gets its own bit, carried by no one.
		for _, b := range loaders {
			if b&roots["/h/empty"] != 0 {
				t.Errorf("a lib_dir without writers must be loadable by no one")
			}
		}
	}
}

func TestPlanTrustedDirs_OverflowFailsClosed(t *testing.T) {
	var dirs []TrustedDir
	for i := range 70 {
		dirs = append(dirs, TrustedDir{Path: string(rune('a'+i%26)) + string(rune('A'+i/26)),
			Loaders: []string{string(rune('0' + i))}})
	}
	roots, _, refused := planTrustedDirs(dirs)
	if len(refused) != 70-63 {
		t.Fatalf("refused %d, want %d", len(refused), 70-63)
	}
	for _, p := range refused {
		if m, ok := roots[p]; !ok || m != 0 {
			t.Errorf("an overflowing lib_dir must be recorded with no loaders, got %#x (present %v)", m, ok)
		}
	}
}
