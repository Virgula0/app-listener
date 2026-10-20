package pulse

import (
	"testing"
	"time"
)

func recv(s *Sub, within time.Duration) bool {
	select {
	case <-s.C():
		s.Ack()
		return true
	case <-time.After(within):
		return false
	}
}

func TestSubFiresEveryPeriod(t *testing.T) {
	p := New(10 * time.Millisecond)
	a, b := p.Subscribe(), p.Subscribe()
	for i := range 3 {
		if !recv(a, time.Second) || !recv(b, time.Second) {
			t.Fatalf("tick %d not delivered to both subscribers", i)
		}
	}
}

// A loop busy across a tick still gets it, once.
func TestSubDeliversMissedTickLate(t *testing.T) {
	p := New(10 * time.Millisecond)
	s := p.Subscribe()
	ch := s.C() // grabbed before the tick, as a loop blocked in its select would
	time.Sleep(35 * time.Millisecond)
	select {
	case <-ch:
	default:
		t.Fatal("a channel grabbed before the tick must be closed by it")
	}
	select {
	case <-s.C():
	default:
		t.Fatal("ticks during a busy stretch must be pending, not lost")
	}
	s.Ack()
	select {
	case <-s.C():
		t.Fatal("several missed ticks must collapse into one")
	default:
	}
}

func TestSubscribeDoesNotFireImmediately(t *testing.T) {
	p := New(time.Hour)
	p.Subscribe()
	select {
	case <-p.Subscribe().C():
		t.Fatal("a new subscriber fired before any tick")
	case <-time.After(20 * time.Millisecond):
	}
}
