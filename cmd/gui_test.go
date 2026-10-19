package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestCommandTree pins the subcommand tree and the old spellings a binary of another version
// runs (installed units, install.sh, older `update`s), which stay listed in --help.
func TestCommandTree(t *testing.T) {
	for _, path := range []string{
		"monitor", "guard", "network monitor", "network guard", "net guard",
		"daemon check", "daemon lockdown", "daemon genkey", "daemon-unit",
		"edit-protected forward", "edit-protected config record", "edit-protected put", "edit-protected password",
		"install binary", "install catalog refresh", "install catalog diff",
		"install backups restore", "install backups delete", "trust list", "trust add",
		"uninstall", "update",
	} {
		args := strings.Fields(path)
		c, _, err := rootCmd.Find(args)
		if err != nil || c.Name() != args[len(args)-1] && !c.HasAlias(args[len(args)-1]) {
			t.Errorf("%q: found %v, %v", path, c.Name(), err)
		}
	}
	daemonCmd, _, _ := rootCmd.Find([]string{"daemon"})
	for _, name := range []string{"check", "lockdown", "genkey", "verifier-only"} {
		if f := daemonCmd.Flags().Lookup(name); f == nil || f.Hidden {
			t.Errorf("daemon --%s must stay registered and listed", name)
		}
	}
	ep, _, _ := rootCmd.Find([]string{"edit-protected"})
	if f := ep.Flags().Lookup("backend"); f == nil || f.Hidden {
		t.Error("edit-protected --backend must stay registered and listed (the GUI runs it)")
	}
	inst, _, _ := rootCmd.Find([]string{"install"})
	if f := inst.Flags().Lookup("mark-complete"); f == nil || f.Hidden {
		t.Error("install --mark-complete must stay registered and listed (the GUI installer runs it)")
	}
	if unit, _, _ := rootCmd.Find([]string{"daemon-unit"}); unit.Hidden {
		t.Error("daemon-unit must stay listed")
	}
	if guiCmd.HasSubCommands() {
		t.Error("the gui build's root must only open the window")
	}
}

func TestMovedFlagError(t *testing.T) {
	find := func(path string) *cobra.Command {
		c, _, err := rootCmd.Find(strings.Fields(path))
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
	cases := []struct{ cmd, err, want string }{
		{"edit-protected", "unknown flag: --forward", "use `app-listener edit-protected forward`"},
		{"edit-protected", "unknown shorthand flag: 'w' in -w", "use `app-listener edit-protected forward`"},
		{"install", "unknown flag: --binary-only", "use `app-listener install binary`"},
		{"monitor", "unknown shorthand flag: 'w' in -w", "monitor PATH..."},
		{"trust", "unknown flag: --list", "use `app-listener trust list`"},
		{"install", "unknown flag: --nope", ""},
	}
	for _, c := range cases {
		got := movedFlagError(find(c.cmd), errors.New(c.err)).Error()
		if c.want == "" && got != c.err || !strings.Contains(got, c.want) {
			t.Errorf("%s %s: %q, want %q", c.cmd, c.err, got, c.want)
		}
	}
}
