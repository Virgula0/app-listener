package guard

import "testing"

func TestParseGlobName(t *testing.T) {
	cases := []struct {
		in   string
		kind uint8
		text string
		ok   bool
	}{
		{"wineserver", globExact, "wineserver", true},
		{"pv-*", globPrefix, "pv-", true},
		{"*-capsule-capture-libs", globSuffix, "-capsule-capture-libs", true},
		{"*", 0, "", false},
		{"a*b", 0, "", false},
		{"*mid*", 0, "", false},
		{"wine?", 0, "", false},
		{"[ab]", 0, "", false},
		{"0123456789012345678901234567890123", 0, "", false}, // >= 32 bytes
	}
	for _, c := range cases {
		kind, text, ok := ParseGlobName(c.in)
		if kind != c.kind || text != c.text || ok != c.ok {
			t.Errorf("ParseGlobName(%q) = (%d, %q, %v), want (%d, %q, %v)", c.in, kind, text, ok, c.kind, c.text, c.ok)
		}
	}
}

func TestCompileGlobNames_RepeatedPatternOrsBits(t *testing.T) {
	names, affixes, err := compileGlobNames([]string{"*.so", "wineserver", "*.so"})
	if err != nil {
		t.Fatal(err)
	}
	k := GuardTrustGlobNameKey{Kind: globSuffix, Len: 3, S: globText(".so")}
	if names[k] != 0b101 {
		t.Errorf("*.so bits = %b, want 101", names[k])
	}
	if len(affixes) != 1 {
		t.Errorf("one *.so shape, got %v", affixes)
	}
}

func TestGlobReservationsLibTrusted(t *testing.T) {
	r := GlobReservations{
		Patterns: []string{"*.so", "lib*", "wineserver"},
		Roots:    map[string]uint64{"/h/.config/discord": 0b011, "/h/steam": 0b100},
		Writers:  map[string]uint64{"/h/.config/discord/app/Discord": 0b011},
	}
	const discord = "/h/.config/discord/app/Discord"
	cases := []struct {
		bin, lib string
		want     bool
	}{
		{discord, "/h/.config/discord/app/libffmpeg.so", true},
		{discord, "/h/.config/discord/app/modules/x/libvulkan.so.1", true}, // lib* at depth
		{discord, "/h/.config/discord/app/probe.dat", false},                // unreserved name
		{discord, "/h/.config/other/libffmpeg.so", false},                   // outside the root
		{"/usr/bin/ssh", "/h/.config/discord/app/libffmpeg.so", false},      // not a writer
		{discord, "/h/steam/wineserver", false},                             // another entry's bit
	}
	for _, c := range cases {
		if got := r.LibTrusted(c.bin, c.lib); got != c.want {
			t.Errorf("LibTrusted(%s, %s) = %v, want %v", c.bin, c.lib, got, c.want)
		}
	}
}
