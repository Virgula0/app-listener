package install

import (
	"strings"
	"testing"
)

// TestSampleFilesPresent asserts every daemon-samples file the installer
// depends on by name is embedded and non-empty.
func TestSampleFilesPresent(t *testing.T) {
	want := []string{
		"app-listener-daemon.service",
		"app-listener-catalog-refresh.service",
		"ssh-agent.service",
		"daemon.conf",
		"50-app-listener-reload.hook", // pacman hook
		"apt-app-listener-reload",     // apt hook sample
	}

	got, err := SampleFiles()
	if err != nil {
		t.Fatalf("SampleFiles: %v", err)
	}
	have := make(map[string]bool, len(got))
	for _, n := range got {
		have[n] = true
	}

	for _, n := range want {
		if !have[n] {
			t.Errorf("daemon-samples is missing %q (embedded set: %v)", n, got)
			continue
		}
		content, err := SampleContent(n)
		if err != nil || len(content) == 0 {
			t.Errorf("SampleContent(%q): %d bytes, err=%v", n, len(content), err)
		}
	}
}

// TestCatalogRefreshUnitShape guards the two properties the boot-time refresh
// relies on: it orders after the daemon and runs the --live refresh.
func TestCatalogRefreshUnitShape(t *testing.T) {
	b, err := SampleContent("app-listener-catalog-refresh.service")
	if err != nil {
		t.Fatal(err)
	}
	unit := string(b)
	for _, needle := range []string{
		"After=app-listener-daemon.service",
		"Requires=app-listener-daemon.service",
		"--update-catalog-only --live --yes",
		"Type=oneshot",
	} {
		if !strings.Contains(unit, needle) {
			t.Errorf("catalog-refresh unit missing %q\n%s", needle, unit)
		}
	}
}

// TestAptHookShape: the apt hook must call --update-catalog-only and must not
// fail an apt run on a refresh error.
func TestAptHookShape(t *testing.T) {
	b, err := SampleContent("apt-app-listener-reload")
	if err != nil {
		t.Fatal(err)
	}
	hook := string(b)
	for _, needle := range []string{
		"DPkg::Post-Invoke",
		"install --update-catalog-only --yes",
		"|| true",
	} {
		if !strings.Contains(hook, needle) {
			t.Errorf("apt hook missing %q\n%s", needle, hook)
		}
	}
}
