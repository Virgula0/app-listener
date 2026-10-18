//go:build gui

package gui

import (
	"os"
	"strings"

	"github.com/Virgula0/app-listener/internal/systemd"
)

// isPreRelease reports whether version is a pre-YYYYMMDD build: its channel is pre-release, as
// `update --channel` would track it.
func isPreRelease(version string) bool { return strings.HasPrefix(version, "pre-") }

func boolIndex(b bool) int {
	if b {
		return 1
	}
	return 0
}

// checkForUpdate asks GitHub, off the UI goroutine, whether the installed daemon's channel has a
// newer release, and badges the Update entry. A failed lookup (offline, a network whitelist in this
// window) leaves the badge as it was.
func (s *shell) checkForUpdate() {
	versionOf, latest := s.opts.Install.InstalledVersion, s.opts.Update.Latest
	if versionOf == nil || latest == nil || s.running() {
		return
	}
	var tag string
	var newer, ok bool
	s.background(func() {
		if _, err := os.Stat(systemd.InstallBinaryPath); err != nil {
			ok = true // nothing installed: nothing to update
			return
		}
		installed := versionOf()
		var err error
		tag, newer, err = latest(isPreRelease(installed), installed)
		ok = err == nil
	}, func() {
		if ok {
			s.setUpdateAvailable(tag, newer)
		}
	})
}

// setUpdateAvailable badges the Update entry while newer holds for tag.
func (s *shell) setUpdateAvailable(tag string, newer bool) {
	if n := s.nav[pageUpdate]; n != nil {
		n.setBadge(newer && tag != "")
	}
}
