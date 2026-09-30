package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// setupDumpLog runs in PersistentPreRunE, i.e. as root before any subcommand starts and before the
// daemon's self-guards attach. It os.Stat()s the operator-supplied path and then opens it with
// O_CREATE|O_WRONLY|O_TRUNC — both follow symlinks.
//
// The "refuse to overwrite under --headless" branch only triggers when Stat SUCCEEDS, so a symlink
// aimed at a path that does not exist yet skips the refusal entirely and root creates the link
// target. Aimed at an existing file it truncates it instead: pointing it at
// /etc/app-listener/fscrypt.key destroys the vault master key.
func TestSetupDumpLogRefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	victim := filepath.Join(dir, "victim-does-not-exist-yet")
	link := filepath.Join(dir, "trace.log")
	if err := os.Symlink(victim, link); err != nil {
		t.Fatal(err)
	}

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Bool("headless", true, "")

	prev := dumpLogFile
	dumpLogFile = link
	t.Cleanup(func() { dumpLogFile = prev })

	err := setupDumpLog(cmd)

	if _, serr := os.Lstat(victim); serr == nil {
		t.Fatalf("root created %s through a planted --dump-log symlink", victim)
	}
	if err == nil {
		t.Fatal("setupDumpLog accepted a symlinked path; expected a refusal")
	}
}
