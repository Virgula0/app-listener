package daemon

import "testing"

func TestPendingGroupUnlockStopRetiresDrainOnce(t *testing.T) {
	done := make(chan struct{})
	p := &pendingGroupUnlock{drained: done}
	p.stop()
	select {
	case <-done:
	default:
		t.Fatal("stop must close drained so drainEphemeral readers exit")
	}
	p.stop() // deferred second call on every startup path: must not panic
}
