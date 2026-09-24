package daemon

import (
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func rules(paths ...string) []daemonconfig.BinaryRule {
	out := make([]daemonconfig.BinaryRule, len(paths))
	for i, p := range paths {
		out[i] = daemonconfig.BinaryRule{Path: p}
	}
	return out
}

func TestPlanUpdaters_ScopedToOwnResources(t *testing.T) {
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: "/protected", Binaries: rules("/home/u/.local/bin/app", "/usr/bin/git"), AllowLibs: []string{"/home/u/lib.so"}},
		{Path: "/other", Binaries: rules("/opt/other/cp", "/opt/other/tool")},
		{Path: "/same", Binaries: rules("/usr/bin/git", "/home/u/.local/bin/app")},
	}}
	p := planUpdaters(cfg)

	app := p.Owners["/home/u/.local/bin/app"]
	if app == 0 || p.Owners["/home/u/lib.so"] != app {
		t.Fatalf("the resource's binaries and allow_libs must share its bit: %v", p.Owners)
	}
	if p.Updaters["/home/u/.local/bin/app"]&app == 0 {
		t.Errorf("a binary must update its own resource: %v", p.Updaters)
	}
	if p.Updaters["/opt/other/tool"]&app != 0 {
		t.Errorf("another resource's binary must not update app: %v", p.Updaters)
	}
	for _, tool := range []string{"/usr/bin/git", "/opt/other/cp"} {
		if p.Updaters[tool] != 0 {
			t.Errorf("general tool %s must never be an updater: %v", tool, p.Updaters)
		}
	}
	if p.Owners["/usr/bin/git"] != app {
		t.Errorf("identical binary sets must share one bit: %v", p.Owners)
	}
}

func TestIsGeneralTool(t *testing.T) {
	for path, want := range map[string]bool{
		"/usr/bin/python3.12": true, "/usr/bin/node": true, "/bin/sh": true,
		"/home/u/.local/bin/claude": false, "/opt/steam/steam": false,
	} {
		if got := isGeneralTool(path); got != want {
			t.Errorf("isGeneralTool(%s) = %v, want %v", path, got, want)
		}
	}
}
