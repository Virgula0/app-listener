package guard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPinPrefixShape(t *testing.T) {
	p := PinPrefix("/sys/fs/bpf/app-listener", "dl956onp5giy", "/home/angelo/.config/discord")
	// Flat: one path segment under base, ending with "-" so the guard appends
	// the hook name. No dots, no underscores, no extra nesting.
	rel := strings.TrimPrefix(p, "/sys/fs/bpf/app-listener/")
	if strings.Contains(rel, "/") {
		t.Errorf("pin prefix nests below base: %q", p)
	}
	if !strings.HasPrefix(rel, "al-") || !strings.HasSuffix(rel, "-") {
		t.Errorf("unexpected pin prefix shape: %q", rel)
	}
	for _, r := range rel {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-'
		if !ok {
			t.Errorf("pin prefix %q contains disallowed rune %q", rel, r)
		}
	}
}

func TestPinPrefixDistinctAndStable(t *testing.T) {
	a := PinPrefix("/b", "g", "/home/angelo/.ssh")
	if a != PinPrefix("/b", "g", "/home/angelo/.ssh") {
		t.Error("PinPrefix is not deterministic")
	}
	if a == PinPrefix("/b", "g", "/home/angelo/.gnupg") {
		t.Error("different resources share a pin prefix")
	}
	if PinPrefix("/b", "g1", "/x") == PinPrefix("/b", "g2", "/x") {
		t.Error("different generations share a pin prefix")
	}
	// Sanitize-collision paths must still differ (hash disambiguates).
	if PinPrefix("/b", "g", "/a/b") == PinPrefix("/b", "g", "/a_b") {
		t.Error("sanitize collision not disambiguated")
	}
}

func TestGenOfPinFile(t *testing.T) {
	cases := map[string]string{
		"al-dl956onp5giy-266d672f9abc-file-open": "dl956onp5giy", // hook underscores are normalised to dashes
		"al-g1-abcdef012345-path-rename":         "g1",
		"al-g2-abcdef012345-inode-removexattr":   "g2",
		"cilium_something":                       "",
		"al-onlytwo":                             "",
		"random":                                 "",
	}
	for name, want := range cases {
		if got := genOfPinFile(name); got != want {
			t.Errorf("genOfPinFile(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestCleanupStalePins(t *testing.T) {
	base := t.TempDir()
	write := func(name string) {
		if err := os.WriteFile(filepath.Join(base, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Two stale generations, one live, plus a foreign tool's pin.
	write("al-oldgenaa-111111111111-file_open")
	write("al-oldgenaa-111111111111-file_permission")
	write("al-oldgenbb-222222222222-path_unlink")
	write("al-livegen0-333333333333-file_open")
	write("cilium-tc-globals") // not ours: must survive
	if err := os.Mkdir(filepath.Join(base, "somedir"), 0o755); err != nil {
		t.Fatal(err)
	}

	n, err := CleanupStalePins(base, map[string]bool{"livegen0": true})
	if err != nil {
		t.Fatalf("CleanupStalePins: %v", err)
	}
	if n != 3 {
		t.Errorf("removed %d, want 3", n)
	}
	for _, gone := range []string{"al-oldgenaa-111111111111-file_open", "al-oldgenbb-222222222222-path_unlink"} {
		if _, err := os.Stat(filepath.Join(base, gone)); !os.IsNotExist(err) {
			t.Errorf("%s should have been removed", gone)
		}
	}
	for _, kept := range []string{"al-livegen0-333333333333-file_open", "cilium-tc-globals"} {
		if _, err := os.Stat(filepath.Join(base, kept)); err != nil {
			t.Errorf("%s must be kept: %v", kept, err)
		}
	}
}

func TestCleanupStalePinsMissingBase(t *testing.T) {
	n, err := CleanupStalePins(filepath.Join(t.TempDir(), "nope"), nil)
	if err != nil || n != 0 {
		t.Errorf("missing base: got (%d, %v), want (0, nil)", n, err)
	}
}
