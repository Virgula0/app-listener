package daemon

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func TestLaunchRulesByPath(t *testing.T) {
	cfg := &daemonconfig.Config{
		Resources: []daemonconfig.Resource{
			{Binaries: []daemonconfig.BinaryRule{
				{Path: "/opt/gh/2.80/gh-real", Link: "/home/u/.local/bin/gh"},
				{Path: "/usr/bin/ssh"},
				{Path: "/usr/bin/gh-helper", LibBinary: true},
			}},
			{PendingBinaries: []daemonconfig.BinaryRule{{Path: "/home/u/.local/bin/bitcoin-cli"}}},
		},
		LaunchRules: []daemonconfig.LaunchRule{
			{Path: "/home/u/.local/bin/gh", Words: []string{"repo", "delete"}},
			{Path: "/usr/bin/ssh", Env: "SSH_ASKPASS"},
			{Path: "/usr/bin/nowhere", Words: []string{"x"}},
		},
	}
	got := launchRulesByPath(cfg)
	gh := got["/opt/gh/2.80/gh-real"]
	if len(gh) != 7 || gh[6].Words[0] != "repo" {
		t.Errorf("gh (via its link's name, plus the rule naming the link) = %+v", gh)
	}
	if ssh := got["/usr/bin/ssh"]; len(ssh) != 1 || ssh[0].Env != "SSH_ASKPASS" {
		t.Errorf("ssh = %+v", ssh)
	}
	if len(got["/home/u/.local/bin/bitcoin-cli"]) != 8 {
		t.Errorf("a pending bitcoin-cli gets its built-ins: %+v", got["/home/u/.local/bin/bitcoin-cli"])
	}
	if _, ok := got["/usr/bin/gh-helper"]; ok {
		t.Error("a lib_binary line is no whitelisted exe")
	}
	if len(got) != 3 {
		t.Errorf("binaries without rules are left out: %v", got)
	}
}

func TestAllowedEnvKeepsOnlyRootPlaced(t *testing.T) {
	cfg := &daemonconfig.Config{AllowEnv: []string{
		"SSL_CERT_FILE=/etc/passwd", "SSL_CERT_FILE=" + t.TempDir() + "/ca.pem", "SSL_CERT_DIR=/tmp",
	}}
	got := allowedEnv(cfg)
	if !slices.Contains(got, "SSL_CERT_FILE=/etc/passwd") {
		t.Errorf("a root-placed allow_env is dropped: %v", got)
	}
	for _, a := range got {
		if strings.HasPrefix(a, "SSL_CERT_DIR=/tmp") || strings.Contains(a, os.TempDir()) {
			t.Errorf("a user-fillable path is allowed: %s", a)
		}
	}
}
