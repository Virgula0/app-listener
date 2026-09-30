package editprotected

import "testing"

// Creating the FIRST edit-protected password arms live mode. On a RUNNING daemon (vaults unlocked)
// this must be refused: otherwise any root process runs `edit-protected --set-password`, and because
// the writer is the daemon binary the self-guard permits the write, the daemon SIGHUPs open the
// control socket, and that process then reads every protected resource — with no proof of ownership.
// The first password belongs to the audited `install` flow. Offline (daemon stopped) it is allowed.
func TestAuthorizeFirstTimeSetPasswordRefusedWhileDaemonRunning(t *testing.T) {
	if err := authorizeFirstTimeSetPassword(true); err == nil {
		t.Fatal("first-time --set-password was allowed while the daemon is running: any root " +
			"process can arm live edit-protected on unlocked vaults without proving ownership")
	}
	if err := authorizeFirstTimeSetPassword(false); err != nil {
		t.Fatalf("first-time --set-password must be allowed offline (daemon stopped): %v", err)
	}
}
