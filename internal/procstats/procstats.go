package procstats

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Stats struct {
	RSS     uint64
	VMSize  uint64
	CPUUser time.Duration
	CPUSys  time.Duration
	NumCPU  int
}

func Read() (*Stats, error) {
	s := &Stats{NumCPU: runtime.NumCPU()}

	status, err := os.ReadFile("/proc/self/status")
	if err == nil {
		parseStatus(string(status), s)
	}

	stat, err := os.ReadFile("/proc/self/stat")
	if err == nil {
		parseStat(string(stat), s)
	}

	return s, nil
}

func parseStatus(data string, s *Stats) {
	for _, line := range strings.Split(data, "\n") {
		switch {
		case strings.HasPrefix(line, "VmRSS:"):
			s.RSS = parseKBValue(line)
		case strings.HasPrefix(line, "VmSize:"):
			s.VMSize = parseKBValue(line)
		}
	}
}

func parseKBValue(line string) uint64 {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0
	}
	v, _ := strconv.ParseUint(parts[1], 10, 64)
	return v * 1024
}

func parseStat(data string, s *Stats) {
	r := strings.LastIndex(data, ")")
	if r < 0 || r+1 >= len(data) {
		return
	}
	fields := strings.Fields(data[r+2:])
	if len(fields) < 13 {
		return
	}

	utime, _ := strconv.ParseInt(fields[11], 10, 64)
	stime, _ := strconv.ParseInt(fields[12], 10, 64)

	s.CPUUser = time.Duration(utime) * time.Second / 100
	s.CPUSys = time.Duration(stime) * time.Second / 100
}
