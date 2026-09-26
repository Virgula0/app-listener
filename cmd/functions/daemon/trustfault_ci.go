//go:build ci

package daemon

import (
	"fmt"
	"os"
)

// trustFaultEnv fails the trust guard at the named stage ("start", "after-unlock") so the
// integration suite can exercise the refuse-to-start path. Compiled only with -tags ci, which
// release builds never set.
const trustFaultEnv = "APPLISTENER_TEST_TRUST_FAIL"

func injectedTrustFault(stage string) error {
	if os.Getenv(trustFaultEnv) == stage {
		return fmt.Errorf("injected trust guard failure at %s", stage)
	}
	return nil
}
