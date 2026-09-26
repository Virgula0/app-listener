package common

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// --serve mirrors the root daemon's live event stream (guarded paths, comms, PIDs, peer addresses)
// over a loopback TCP socket. HTTP Basic Auth is wired only when --user/--password are given, and
// they are optional: the browser-facing Origin/Host checks constrain only browsers, not a local
// non-browser WebSocket client, and loopback TCP carries no uid restriction. Enabling --serve must
// therefore require credentials (or an equivalent peer-uid check), never default to open.
func TestParseServeFlagsRequiresCredentials(t *testing.T) {
	cmd := &cobra.Command{Use: "guard"}
	AddServeFlags(cmd)
	cmd.Flags().Bool("headless", false, "")
	if err := cmd.Flags().Set("serve", "127.0.0.1:9999"); err != nil {
		t.Fatal(err)
	}

	_, err := ParseServeFlags(cmd)
	if err == nil {
		t.Fatal("ParseServeFlags accepted --serve with no --user/--password: the event stream is " +
			"exposed to any local user")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "user") && !strings.Contains(msg, "password") &&
		!strings.Contains(msg, "auth") && !strings.Contains(msg, "credential") {
		t.Fatalf("--serve without credentials must be rejected for missing authentication, "+
			"got an unrelated error: %v", err)
	}
}
