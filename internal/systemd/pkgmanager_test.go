package systemd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageManagerString(t *testing.T) {
	cases := map[PackageManager]string{
		PkgPacman:          "pacman",
		PkgApt:             "apt",
		PkgDnf:             "dnf",
		PkgZypper:          "zypper",
		PackageManager(99): "unknown",
	}
	for pm, want := range cases {
		if got := pm.String(); got != want {
			t.Errorf("PackageManager(%d).String() = %q, want %q", pm, got, want)
		}
	}
}

func TestPackageManagerHasReloadHook(t *testing.T) {
	for _, pm := range []PackageManager{PkgPacman, PkgApt} {
		if !pm.HasReloadHook() {
			t.Errorf("%s should ship a reload hook", pm)
		}
	}
	for _, pm := range []PackageManager{PkgDnf, PkgZypper} {
		if pm.HasReloadHook() {
			t.Errorf("%s should not claim a bundled reload hook", pm)
		}
	}
}

func TestDetectPackageManagers(t *testing.T) {
	tmp := t.TempDir()
	pacmanDir := filepath.Join(tmp, "pacman.d")
	aptDir := filepath.Join(tmp, "apt")
	if err := os.MkdirAll(pacmanDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(aptDir, 0o755); err != nil {
		t.Fatal(err)
	}

	probes := []packageManagerProbe{
		{PkgPacman, []string{pacmanDir}, []string{"\x00no-such-bin"}},
		{PkgApt, []string{aptDir}, []string{"\x00no-such-bin"}},
		{PkgDnf, []string{filepath.Join(tmp, "absent")}, []string{"\x00no-such-bin"}},
	}

	got := detectPackageManagers(probes)
	want := []PackageManager{PkgPacman, PkgApt}
	if len(got) != len(want) {
		t.Fatalf("detectPackageManagers = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("detectPackageManagers = %v, want %v", got, want)
		}
	}
}

func TestDetectPackageManagersByBinary(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "fake-pm")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	probes := []packageManagerProbe{
		{PkgApt, []string{filepath.Join(dir, "absent")}, []string{"fake-pm"}},
	}
	got := detectPackageManagers(probes)
	if len(got) != 1 || got[0] != PkgApt {
		t.Fatalf("detectPackageManagers by binary = %v, want [apt]", got)
	}
}
