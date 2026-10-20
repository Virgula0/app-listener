package common

import (
	"os"
	"os/signal"
	"syscall"
)

// DrainUntilSignal hands each event to handle until events closes or SIGINT/SIGTERM arrives.
func DrainUntilSignal[E any](events <-chan E, handle func(*E)) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}
			handle(&ev)
		case <-sig:
			return
		}
	}
}
