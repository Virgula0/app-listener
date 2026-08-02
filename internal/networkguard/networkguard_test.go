package networkguard

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestInfraFromRunning(t *testing.T) {
	cases := []struct {
		name     string
		running  map[string]bool
		expected []string
	}{
		{
			name:     "empty running set",
			running:  map[string]bool{},
			expected: nil,
		},
		{
			name: "only resolvers running",
			running: map[string]bool{
				"/usr/lib/systemd/systemd-resolved": true,
			},
			expected: []string{"/usr/lib/systemd/systemd-resolved"},
		},
		{
			name:     "nothing infra running",
			running:  map[string]bool{"/usr/bin/firefox": true, "/usr/bin/bash": true},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := infraFromRunning(tc.running)
			sort.Strings(got)
			expected := append([]string(nil), tc.expected...)
			sort.Strings(expected)
			if len(got) != len(expected) {
				t.Fatalf("got %v, want %v", got, expected)
			}
			for i := range got {
				if got[i] != expected[i] {
					t.Fatalf("got %v, want %v", got, expected)
				}
			}
		})
	}
}

func TestDiscoverInfraBinaries(t *testing.T) {
	got, err := DiscoverInfraBinaries()
	if err != nil {
		t.Fatalf("DiscoverInfraBinaries() error: %v", err)
	}
	for _, p := range got {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("discovered path %q does not exist on disk: %v", p, err)
		}
		resolved, err := filepath.EvalSymlinks(p)
		if err != nil {
			t.Errorf("resolving %q: %v", p, err)
			continue
		}
		if !pathIsRunning(resolved) {
			t.Errorf("discovered path %q (%s) is not currently running", p, resolved)
		}
	}
}

func pathIsRunning(resolved string) bool {
	procs, err := os.ReadDir("/proc")
	if err != nil {
		return false
	}
	for _, dir := range procs {
		if !dir.IsDir() {
			continue
		}
		exe, err := os.Readlink(filepath.Join("/proc", dir.Name(), "exe"))
		if err != nil {
			continue
		}
		r, err := filepath.EvalSymlinks(exe)
		if err != nil {
			continue
		}
		if r == resolved {
			return true
		}
	}
	return false
}