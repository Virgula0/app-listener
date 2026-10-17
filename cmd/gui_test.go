package cmd

import "testing"

func TestSubcommandOutranksGUI(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"--gui"}, ""},
		{[]string{"--dump-log", "/tmp/x", "--gui"}, ""},
		{[]string{"-v", "2", "--gui"}, ""},
		{[]string{"monitor", "--gui"}, "monitor"},
		{[]string{"--gui", "network-guard", "-w", "/bin/curl"}, "network-guard"},
		{[]string{"--dump-log", "/tmp/x", "guard", "--gui"}, "guard"},
		{[]string{"help", "--gui"}, "help"},
		{[]string{"--gui", "nosuchcommand"}, ""},
	}
	for _, c := range cases {
		got := ""
		if sub := subcommand(c.args); sub != nil {
			got = sub.Name()
		}
		if got != c.want {
			t.Errorf("subcommand(%q) = %q, want %q", c.args, got, c.want)
		}
	}
}
