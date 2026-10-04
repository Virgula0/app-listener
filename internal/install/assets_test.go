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
		"ssh-agent.service",
		"daemon.conf",
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

// The daemon refreshes its catalog itself: no package-manager hook or boot unit is shipped, so an
// install never deploys one again.
func TestNoLegacyCatalogRefreshSamples(t *testing.T) {
	got, err := SampleFiles()
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range got {
		switch n {
		case "app-listener-catalog-refresh.service", "50-app-listener-reload.hook", "apt-app-listener-reload":
			t.Errorf("daemon-samples still ships the legacy refresh trigger %s", n)
		}
	}
}

// No user code may run before the guards attach (issue #80): the unit must signal readiness only
// once enforcing and order logins, user managers and cron after it.
func TestDaemonUnitOrdersUserCodeAfterReadiness(t *testing.T) {
	unit, err := SampleContent("app-listener-daemon.service")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"\nType=notify\n",
		"Before=systemd-user-sessions.service cron.service crond.service cronie.service atd.service\n",
	} {
		if !strings.Contains(string(unit), want) {
			t.Errorf("daemon unit lacks %q", strings.TrimSpace(want))
		}
	}
}
