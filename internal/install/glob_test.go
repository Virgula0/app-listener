package install

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBoundedGlobMatchesFilepathGlob(t *testing.T) {
	root := t.TempDir()
	for _, p := range []string{"0.0.1/Discord", "0.0.2/Discord", "0.0.3/other", ".hidden/Discord"} {
		mkdirAll(t, filepath.Join(root, filepath.Dir(p)))
		if err := os.WriteFile(filepath.Join(root, p), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, pat := range []string{root + "/*/Discord", root + "/0.0.?/*", root + "/0.0.1/Discord", root + "/missing/*"} {
		want, _ := filepath.Glob(pat)
		got, err := boundedGlob(pat)
		if err != nil {
			t.Fatal(err)
		}
		if len(want) == 0 {
			want = nil
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("boundedGlob(%s) = %v, filepath.Glob = %v", pat, got, want)
		}
	}
}

func TestBoundedGlobStopsAtTheBudget(t *testing.T) {
	old := maxGlobEntries
	maxGlobEntries = 50
	defer func() { maxGlobEntries = old }()
	root := t.TempDir()
	for i := range 200 {
		mkdirAll(t, filepath.Join(root, "junk", string(rune('a'+i%26))+string(rune('a'+i/26))))
	}
	got, err := boundedGlob(root + "/junk/*/Discord")
	if err != nil || len(got) != 0 {
		t.Fatalf("a flooded glob root: %v, %v", got, err)
	}
}

func mkdirAll(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
}
