package daemon

import "testing"

func TestParseReservedPorts(t *testing.T) {
	reserved := parseReservedPorts("8332,40000-40010,junk,9-x\n")
	for p, want := range map[uint16]bool{8332: true, 40000: true, 40005: true, 40010: true, 40011: false, 9: false, 8333: false} {
		if reserved(p) != want {
			t.Errorf("port %d reserved=%v, want %v", p, !want, want)
		}
	}
	if parseReservedPorts("\n")(0) {
		t.Error("an empty list reserves port 0")
	}
}
