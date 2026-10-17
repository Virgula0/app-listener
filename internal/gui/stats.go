//go:build gui

package gui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Virgula0/app-listener/internal/workbench"
)

// procStats samples the whole process (every workspace's engines run in it): CPU share of one core,
// resident memory and the event rate across workspaces.
type procStats struct {
	at     time.Time
	cpu    time.Duration
	events uint64
}

type statsSample struct {
	cpuPercent float64
	rss        uint64
	perSecond  float64
}

func cpuTime() time.Duration {
	var ru syscall.Rusage
	if syscall.Getrusage(syscall.RUSAGE_SELF, &ru) != nil {
		return 0
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

// rssBytes is the resident set from /proc/self/statm (pages), 0 when unreadable.
func rssBytes() uint64 {
	data, err := os.ReadFile("/proc/self/statm")
	if err != nil {
		return 0
	}
	f := strings.Fields(string(data))
	if len(f) < 2 {
		return 0
	}
	pages, err := strconv.ParseUint(f[1], 10, 64)
	if err != nil {
		return 0
	}
	return pages * uint64(os.Getpagesize()) //nolint:gosec // page size is positive
}

// sample returns the rates since the previous call (zero rates on the first).
func (p *procStats) sample() statsSample {
	now, cpu, events := time.Now(), cpuTime(), workbench.Ingested()
	out := statsSample{rss: rssBytes()}
	if !p.at.IsZero() {
		if dt := now.Sub(p.at).Seconds(); dt > 0 {
			out.cpuPercent = (cpu - p.cpu).Seconds() / dt * 100
			out.perSecond = float64(events-p.events) / dt
		}
	}
	p.at, p.cpu, p.events = now, cpu, events
	return out
}

func humanBytes(n uint64) string {
	const mib = 1 << 20
	if n >= 1<<30 {
		return fmt.Sprintf("%.2f GiB", float64(n)/(1<<30))
	}
	return fmt.Sprintf("%.0f MiB", float64(n)/mib)
}

func (s statsSample) lines() (cpu, mem, rate string) {
	return fmt.Sprintf("CPU %.1f%%", s.cpuPercent), "Memory " + humanBytes(s.rss),
		fmt.Sprintf("%.0f events/s", s.perSecond)
}
