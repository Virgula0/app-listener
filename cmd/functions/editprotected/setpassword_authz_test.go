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

func fixed(pw string) func() (string, error) { return func() (string, error) { return pw, nil } }

func setDaemonRunning(t *testing.T, up bool) {
	t.Helper()
	old := daemonRunning
	daemonRunning = func() bool { return up }
	t.Cleanup(func() { daemonRunning = old })
}

func TestSetPasswordOfflineChangeNeedsCurrent(t *testing.T) {
	useTempStore(t)
	setDaemonRunning(t, false)
	const cur, next = "Correct-Horse-9", "Brand-New-Horse-7"
	if err := StoreHash(mustHash(t, cur, OriginInstall)); err != nil {
		t.Fatal(err)
	}
	if err := setPassword(fixed("Wrong-Horse-99"), fixed(next)); err == nil {
		t.Fatal("changed the password without the current one")
	}
	if err := CheckPassword(cur); err != nil {
		t.Fatalf("a refused change touched the password: %v", err)
	}
	if err := setPassword(fixed(cur), fixed(next)); err != nil {
		t.Fatalf("change with the current password refused: %v", err)
	}
	if err := CheckPassword(next); err != nil {
		t.Fatalf("new password not stored: %v", err)
	}
}

// An installer-chosen password is changed like any other: the current one is the proof.
func TestSetPasswordChangesInstallOrigin(t *testing.T) {
	useTempStore(t)
	setDaemonRunning(t, false)
	if err := StoreHash(mustHash(t, "Correct-Horse-9", OriginInstall)); err != nil {
		t.Fatal(err)
	}
	if err := setPassword(fixed("Correct-Horse-9"), fixed("Brand-New-Horse-7")); err != nil {
		t.Fatal(err)
	}
}

func TestSetPasswordFirstTimeRefusedWhileRunning(t *testing.T) {
	useTempStore(t)
	setDaemonRunning(t, true)
	asked := false
	next := func() (string, error) { asked = true; return "Brand-New-Horse-7", nil }
	if err := setPassword(fixed(""), next); err == nil || asked {
		t.Fatalf("first password set on a running daemon (err=%v, prompted=%v)", err, asked)
	}
	if ok, _ := HashConfigured(); ok {
		t.Fatal("a refused first-time set stored a password")
	}
}

func TestSetPasswordWeakNewRefused(t *testing.T) {
	useTempStore(t)
	setDaemonRunning(t, false)
	weak := func() (string, error) { return "short", ValidatePassword("short") }
	if err := setPassword(fixed(""), weak); err == nil {
		t.Fatal("a weak first password was stored")
	}
	if err := StoreHash(mustHash(t, "Correct-Horse-9", OriginInstall)); err != nil {
		t.Fatal(err)
	}
	if err := setPassword(fixed("Correct-Horse-9"), weak); err == nil {
		t.Fatal("a weak new password was stored")
	}
}
