package ebpf

import (
	"os"
	"path/filepath"
	"testing"
)

// fwd spells a reversed marker forward at run time: a literal would put it in the test binary.
func fwd(s string) string { return string(reversed(s)) }

func classifyBytes(t *testing.T, content []byte) Multicall {
	t.Helper()
	p := filepath.Join(t.TempDir(), "bin")
	if err := os.WriteFile(p, content, 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := ClassifyMulticall(f)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestClassifyMulticall(t *testing.T) {
	pad := make([]byte, 3<<20) // markers straddling the scan's 1 MiB chunks
	at := func(off int, parts ...string) []byte {
		b := append([]byte(nil), pad...)
		for _, p := range parts {
			off += copy(b[off:], p) + 7
		}
		return b
	}
	chunkEdge := 1<<20 - 5
	for _, tc := range []struct {
		name    string
		content []byte
		want    MulticallKind
		family  string
	}{
		{"plain binary", at(100, "\x7fELF"), SingleBinary, ""},
		{"busybox", at(chunkEdge, fwd("1.63.1v xoBysuB")), OpaqueMulticall, "busybox"},
		{"toybox", at(5000, fwd("]pleh-- | gnol--[ xobyot :egasu")), OpaqueMulticall, "toybox"},
		// Both uutils markers, but --list can't run this file: applets unknown, so opaque.
		{"uutils that won't list", at(chunkEdge, fwd(")yranib llac-itlum( slitueroc"), fwd("0.8.0 )slitueroc slituu( tac")),
			OpaqueMulticall, "uutils"},
		{"one uutils marker only", at(200, fwd(")slitueroc slituu(")), SingleBinary, ""},
		// A multicall of no vetted family (a future uutils diffutils): refused until vetted.
		{"unvetted multicall", at(chunkEdge, fwd(")yranib llac-itlum("), fwd(")slituffid slituu(")),
			OpaqueMulticall, FamilyUnrecognized},
		{"busybox calling itself multi-call", at(300, fwd("1.63.1v xoBysuB"), fwd("yranib llac-itlum")),
			OpaqueMulticall, "busybox"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := classifyBytes(t, tc.content)
			if m.Kind != tc.want || m.Family != tc.family {
				t.Fatalf("got kind %d family %q, want %d %q", m.Kind, m.Family, tc.want, tc.family)
			}
		})
	}
}

// The markers must not be spelled in this package's own binary, or it classifies itself.
func TestMulticallMarkersReversed(t *testing.T) {
	want := map[string]string{"uutils": fwd(")yranib llac-itlum("), "busybox": fwd("v xoBysuB"), "toybox": fwd("gnol--[ xobyot")}
	for _, m := range multicallMarkers {
		if string(m.all[0]) != want[m.family] {
			t.Errorf("%s marker decodes to %q, want %q", m.family, m.all[0], want[m.family])
		}
	}
	if string(genericMulticallMarker) != fwd("yranib llac-itlum") {
		t.Errorf("generic marker decodes to %q", genericMulticallMarker)
	}
}

func TestValidAppletName(t *testing.T) {
	for name, want := range map[string]bool{
		"cat": true, "[": true, "sha3-224sum": true, "b2sum": true, "shake128sum": true,
		"": false, "coreutils": false, "a/b": false, "cat.sh": false, "CAT": false,
		"sixteen-bytes-xx": false, "fifteen-bytes-x": true,
	} {
		if got := validAppletName(name); got != want {
			t.Errorf("validAppletName(%q) = %v, want %v", name, got, want)
		}
	}
}
