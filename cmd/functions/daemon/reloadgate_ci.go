//go:build ci

package daemon

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// reloadGateEnv names a file the reload waits for between config validation and guard build, so
// the integration suite can change the filesystem inside that window deterministically. Compiled
// only with -tags ci, which release builds never set.
const reloadGateEnv = "APPLISTENER_TEST_RELOAD_GATE"

func awaitReloadGate() {
	gate := os.Getenv(reloadGateEnv)
	if gate == "" {
		return
	}
	log.Warnf("test: reload paused before guard build")
	for deadline := time.Now().Add(time.Minute); time.Now().Before(deadline); {
		if _, err := os.Stat(gate); err == nil {
			_ = os.Remove(gate)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
