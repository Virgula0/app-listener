package guard

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func touchPins(t *testing.T, dir string, names ...string) {
	t.Helper()
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func pinNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range entries {
		out = append(out, e.Name())
	}
	return out
}

func TestHandoverFindsOnlyTheMarkedGeneration(t *testing.T) {
	dir := t.TempDir()
	// g1 handed over one runtime section; g2 is a SIGKILLed run (no marker).
	touchPins(t, dir, "al-g1-shared-handover", "al-g1-shared-file-open", "al-g1-shared-exe-actions",
		"al-g1-aaaaaaaaaaaa-res-id", "al-g2-shared-file-open", "al-g2-bbbbbbbbbbbb-res-id")

	gen, err := handoverGen(dir)
	if err != nil || gen != "g1" {
		t.Fatalf("handoverGen = %q, %v; want g1", gen, err)
	}
	res, err := handoverResources(dir, gen)
	if err != nil || !slices.Equal(res, []string{filepath.Join(dir, "al-g1-aaaaaaaaaaaa-")}) {
		t.Fatalf("handoverResources = %v, %v", res, err)
	}
	if err := RetireHandover(dir); err != nil {
		t.Fatal(err)
	}
	if got := pinNames(t, dir); !slices.Equal(got, []string{"al-g2-bbbbbbbbbbbb-res-id", "al-g2-shared-file-open"}) {
		t.Errorf("after RetireHandover: %v; want only g2's pins (a killed run's go at the start's sweep)", got)
	}
}

func TestHandoverNoneOrAmbiguous(t *testing.T) {
	dir := t.TempDir()
	if gen, err := handoverGen(dir); gen != "" || err != nil {
		t.Errorf("empty base: %q, %v", gen, err)
	}
	if gen, err := handoverGen(""); gen != "" || err != nil {
		t.Errorf("no pinning: %q, %v", gen, err)
	}
	touchPins(t, dir, "al-g1-shared-handover", "al-g1-aaaaaaaaaaaa-res-id", "al-g3-shared-handover", "al-g4-shared-file-open")
	if gen, err := handoverGen(dir); err == nil {
		t.Fatalf("two handovers: got %q, want an error", gen)
	}
	if gen, err := AdmitSelfInHandover(dir, nil); err == nil {
		t.Fatalf("two handovers admitted %q", gen)
	}
	if err := RetireHandover(dir); err != nil {
		t.Fatal(err)
	}
	if got := pinNames(t, dir); !slices.Equal(got, []string{"al-g4-shared-file-open"}) {
		t.Errorf("after RetireHandover: %v; want every handover gone, the other run's pins kept", got)
	}
}
