package install

import (
	"errors"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// TestVerifyResourcesLocked exercises the post-stop lockdown verification:
// every encrypted resource must be keyless, unencrypted ones are skipped,
// and a still-unlocked directory is a fatal error.
func TestVerifyResourcesLocked(t *testing.T) {
	locked := func(path string) (bool, error) { return false, nil }
	unlocked := func(path string) (bool, error) { return true, nil }
	checkErr := func(path string) (bool, error) { return false, errors.New("boom") }

	// Everything locked: pass.
	if err := verifyResourcesLocked([]daemonconfig.Resource{
		{Path: "/enc/a", NeedEncryption: true},
		{Path: "/enc/b", NeedEncryption: true},
	}, locked); err != nil {
		t.Errorf("all locked: unexpected error: %v", err)
	}

	// One still unlocked: fatal error naming the directory.
	err := verifyResourcesLocked([]daemonconfig.Resource{
		{Path: "/enc/a", NeedEncryption: true},
		{Path: "/enc/b", NeedEncryption: true},
	}, unlocked)
	if err == nil {
		t.Fatal("expected error when a resource is still unlocked")
	}
	if !strings.Contains(err.Error(), "/enc/b") {
		t.Errorf("error must name the unlocked resource, got: %v", err)
	}

	// Unencrypted resources are never checked, even when "unlocked".
	if err := verifyResourcesLocked([]daemonconfig.Resource{
		{Path: "/plain", NeedEncryption: false},
	}, unlocked); err != nil {
		t.Errorf("unencrypted resource must be skipped: %v", err)
	}

	// Empty configuration: nothing to verify.
	if err := verifyResourcesLocked(nil, unlocked); err != nil {
		t.Errorf("empty configuration: unexpected error: %v", err)
	}

	// A failing lock-state check is propagated, not swallowed.
	if err := verifyResourcesLocked([]daemonconfig.Resource{
		{Path: "/enc/a", NeedEncryption: true},
	}, checkErr); err == nil || !strings.Contains(err.Error(), "/enc/a") {
		t.Errorf("check failure must be wrapped with the path, got: %v", err)
	}
}
