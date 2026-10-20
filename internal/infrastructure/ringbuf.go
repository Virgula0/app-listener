package ebpf

import "github.com/cilium/ebpf/ringbuf"

// PumpRingbuf forwards read's events to out until rd is closed or done fires, then closes rd. read
// returns ok=false once rd is closed, and a nil event for a record to skip.
func PumpRingbuf[E any](rd *ringbuf.Reader, read func(*ringbuf.Reader) (*E, bool), out chan<- E, done <-chan struct{}) {
	defer rd.Close()

	for {
		ev, ok := read(rd)
		if !ok {
			return
		}
		if ev == nil {
			continue
		}

		select {
		case out <- *ev:
		case <-done:
			return
		}
	}
}
