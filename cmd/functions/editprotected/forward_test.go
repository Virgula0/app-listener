package editprotected

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// TestForwardFlagSurface checks the rules that hold before root or the daemon are needed.
func TestForwardFlagSurface(t *testing.T) {
	reset := func() {
		whitelistFlags, blacklistFlags, eventsFlags, forwardResources = nil, nil, nil, nil
		EditProtectedCmd.SetArgs(nil)
		for _, c := range EditProtectedCmd.Commands() {
			c.Flags().VisitAll(func(f *pflag.Flag) { f.Changed = false })
		}
	}
	t.Cleanup(reset)
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"no binaries", []string{"forward"}, "at least one of the flags"},
		{"-w and -b", []string{"forward", "-w", "/bin/x", "-b", "/bin/y"}, "none of the others"},
		{"bad -e", []string{"forward", "-w", "/bin/x", "-e", "EXEC"}, "unknown event type"},
		{"-w outside forward", []string{"-w", "/bin/x"}, "unknown shorthand flag"},
		{"put without --resource", []string{"put", "f"}, `"resource" not set`},
		{"record without paths", []string{"config", "record"}, "requires at least 1 arg"},
	}
	for _, c := range cases {
		reset()
		EditProtectedCmd.SetArgs(c.args)
		EditProtectedCmd.SilenceUsage, EditProtectedCmd.SilenceErrors = true, true
		err := EditProtectedCmd.Execute()
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v, want %q", c.name, err, c.want)
		}
	}
}
