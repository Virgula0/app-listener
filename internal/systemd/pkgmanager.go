package systemd

import (
	"os"
	"os/exec"
)

// PackageManager identifies a host package manager, so the installer can drop
// the matching post-transaction catalog-refresh hook.
type PackageManager int

const (
	// PkgPacman is Arch's pacman (hook: /etc/pacman.d/hooks/*.hook).
	PkgPacman PackageManager = iota
	// PkgApt is Debian/Ubuntu apt+dpkg (hook: /etc/apt/apt.conf.d/*).
	PkgApt
	// PkgDnf is Fedora/RHEL dnf — detected but not hooked automatically
	// (its post-transaction-actions plugin is optional).
	PkgDnf
	// PkgZypper is openSUSE zypper — detected but not hooked automatically.
	PkgZypper
)

func (p PackageManager) String() string {
	switch p {
	case PkgPacman:
		return "pacman"
	case PkgApt:
		return "apt"
	case PkgDnf:
		return "dnf"
	case PkgZypper:
		return "zypper"
	default:
		return "unknown"
	}
}

// HasReloadHook reports whether the installer ships an automatic
// catalog-refresh hook for this package manager.
func (p PackageManager) HasReloadHook() bool {
	return p == PkgPacman || p == PkgApt
}

// packageManagerProbe pairs a manager with the signals that mark it present:
// the hook drop-in directory (the thing that actually has to exist for a hook
// to work) or the tool on PATH.
type packageManagerProbe struct {
	pm   PackageManager
	dirs []string
	bins []string
}

var packageManagerProbes = []packageManagerProbe{
	{PkgPacman, []string{"/etc/pacman.d"}, []string{"pacman"}},
	{PkgApt, []string{"/etc/apt/apt.conf.d"}, []string{"apt-get", "apt"}},
	{PkgDnf, []string{"/etc/dnf"}, []string{"dnf", "dnf5"}},
	{PkgZypper, []string{"/etc/zypp"}, []string{"zypper"}},
}

// DetectPackageManagers returns every package manager present on the host, in
// a stable order. More than one can match (containers, mixed systems); the
// caller installs a hook for each that has one.
func DetectPackageManagers() []PackageManager {
	return detectPackageManagers(packageManagerProbes)
}

func detectPackageManagers(probes []packageManagerProbe) []PackageManager {
	var out []PackageManager
	for _, probe := range probes {
		if anyDirExists(probe.dirs) || anyBinExists(probe.bins) {
			out = append(out, probe.pm)
		}
	}
	return out
}

func anyDirExists(paths []string) bool {
	for _, p := range paths {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func anyBinExists(names []string) bool {
	for _, n := range names {
		if _, err := exec.LookPath(n); err == nil {
			return true
		}
	}
	return false
}
