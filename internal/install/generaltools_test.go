package install

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	logtest "github.com/sirupsen/logrus/hooks/test"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func TestNamedGeneralTool(t *testing.T) {
	for _, p := range []string{
		"/usr/bin/python3.12", "/usr/bin/gawk-5.4.0", "/bin/nc.openbsd", "/usr/bin/vim.basic", "/usr/bin/vim.gtk3",
		"/usr/bin/tclsh8.6", "/usr/bin/sort", "/usr/bin/sqlite3", "/usr/bin/batcat", "/usr/bin/zcat",
	} {
		if !namedGeneralTool(p) {
			t.Errorf("%s must be a general tool", p)
		}
	}
	for _, p := range []string{
		"/usr/bin/ssh", "/usr/lib/ssh/ssh-keysign", "/usr/bin/gpg-agent", "/usr/bin/git-lfs-helper-x",
		"/usr/bin/lua-language-server", "/usr/bin/code-cli", "/opt/discord/Discord", "/usr/bin/claude",
	} {
		if namedGeneralTool(p) {
			t.Errorf("%s is not a general tool", p)
		}
	}
}

func TestScriptsOf_NamesTheConfiguredLine(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "versions", "pass")
	if err := os.MkdirAll(filepath.Dir(script), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(script, []byte("#!/bin/bash\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "pass")
	res := &daemonconfig.Resource{Binaries: []daemonconfig.BinaryRule{
		{Path: script, Link: link}, {Path: script}, {Path: self},
	}}
	if got := ScriptsOf(res); !slices.Equal(got, []string{link, script}) {
		t.Errorf("ScriptsOf = %v, want [%s %s]", got, link, script)
	}
}

func TestWarnWhitelistLines_OneWarningPerScriptSet(t *testing.T) {
	script := filepath.Join(t.TempDir(), "pv-locale-gen")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: "/home/u/.steam", Binaries: []daemonconfig.BinaryRule{{Path: script}}},
		{Path: "/home/u/.local/share/Steam/config", Binaries: []daemonconfig.BinaryRule{{Path: script}}},
	}}
	hook := logtest.NewGlobal()
	defer hook.Reset()
	WarnWhitelistLines(cfg)
	var hits []string
	for _, e := range hook.AllEntries() {
		if strings.Contains(e.Message, "never match") {
			hits = append(hits, e.Message)
		}
	}
	if len(hits) != 1 || !strings.Contains(hits[0], "/home/u/.steam, /home/u/.local/share/Steam/config") {
		t.Errorf("want one warning naming both resources, got %q", hits)
	}
}
