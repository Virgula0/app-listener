package procstats

import (
	"sync"
	"testing"
)

func TestReadConcurrentRace(t *testing.T) {
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, err := Read()
			if err != nil {
				t.Logf("Read error: %v", err)
				return
			}
			_ = s.RSS
			_ = s.VMSize
			_ = s.CPUUser
			_ = s.CPUSys
			_ = s.NumCPU
		}()
	}
	wg.Wait()
}

func TestStatsFieldsNonZero(t *testing.T) {
	s, err := Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if s.NumCPU == 0 {
		t.Error("NumCPU should be > 0")
	}
	if s.RSS == 0 {
		t.Error("RSS should be > 0 (process always uses memory)")
	}
	if s.VMSize < s.RSS {
		t.Error("VMSize should be >= RSS")
	}
}

func TestStatsValuesMonotonic(t *testing.T) {
	s1, err := Read()
	if err != nil {
		t.Fatalf("first Read() error: %v", err)
	}

	s2, err := Read()
	if err != nil {
		t.Fatalf("second Read() error: %v", err)
	}

	if s2.CPUUser < s1.CPUUser {
		t.Errorf("CPUUser decreased: %v -> %v", s1.CPUUser, s2.CPUUser)
	}
	if s2.CPUSys < s1.CPUSys {
		t.Errorf("CPUSys decreased: %v -> %v", s1.CPUSys, s2.CPUSys)
	}
}
