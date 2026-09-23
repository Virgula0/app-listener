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
