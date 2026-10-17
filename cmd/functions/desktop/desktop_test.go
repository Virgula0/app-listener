package desktop

import (
	"slices"
	"testing"
)

func TestRequested(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		on      bool
		ignored []string
	}{
		{"absent", []string{"monitor", "-w", "/tmp"}, false, []string{"monitor", "-w", "/tmp"}},
		{"bare", []string{"--gui"}, true, nil},
		{"with subcommand and flags", []string{"guard", "/x", "--gui", "-w", "/bin/cat"}, true,
			[]string{"guard", "/x", "-w", "/bin/cat"}},
		{"explicit true", []string{"--gui=true", "daemon"}, true, []string{"daemon"}},
		{"explicit false", []string{"--gui=false", "monitor"}, false, []string{"monitor"}},
		{"bad value left to cobra", []string{"--gui=maybe"}, false, []string{"--gui=maybe"}},
		{"after terminator is an argument", []string{"monitor", "--", "--gui"}, false,
			[]string{"monitor", "--", "--gui"}},
		{"single dash is not the flag", []string{"-gui"}, false, []string{"-gui"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			on, ignored := Requested(c.args)
			if on != c.on || !slices.Equal(ignored, c.ignored) {
				t.Fatalf("Requested(%q) = %v, %q; want %v, %q", c.args, on, ignored, c.on, c.ignored)
			}
		})
	}
}
