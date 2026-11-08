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
	if ssh := got["/usr/bin/ssh"]; len(ssh) != 2 || ssh[0].Short == nil || ssh[1].Env != "SSH_ASKPASS" {
		t.Errorf("ssh = %+v", ssh)
	}
	if len(got["/home/u/.local/bin/bitcoin-cli"]) != 8 {
		t.Errorf("a pending bitcoin-cli gets its catalog rules: %+v", got["/home/u/.local/bin/bitcoin-cli"])
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

// The inode's rules gather every spelling that names it, whichever section lists it first.
func TestLaunchRulesByPathMergesSpellings(t *testing.T) {
	cfg := &daemonconfig.Config{
		Resources: []daemonconfig.Resource{
			{Binaries: []daemonconfig.BinaryRule{{Path: "/opt/t/tool"}}},
			{Binaries: []daemonconfig.BinaryRule{{Path: "/opt/t/tool", Link: "/opt/t/gh"}}},
		},
		LaunchRules: []daemonconfig.LaunchRule{{Path: "/opt/t/gh", Words: []string{"--leak"}}},
	}
	tool := launchRulesByPath(cfg)["/opt/t/tool"]
	if !slices.ContainsFunc(tool, func(r daemonconfig.LaunchRule) bool { return slices.Equal(r.Words, []string{"--leak"}) }) {
		t.Errorf("the link's [launch_rules] line was dropped: %+v", tool)
	}
	if !slices.ContainsFunc(tool, func(r daemonconfig.LaunchRule) bool { return slices.Equal(r.Words, []string{"auth", "token"}) }) {
		t.Errorf("the link name's catalog rules were dropped: %+v", tool)
	}
}

func TestRefuseUnfitLaunchRules(t *testing.T) {
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{{Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/gpg"}}}}}
	if err := refuseUnfitLaunchRules(cfg); err != nil {
		t.Fatalf("the catalog's gpg rules fit: %v", err)
	}
	for i := range 11 {
		cfg.LaunchRules = append(cfg.LaunchRules, daemonconfig.LaunchRule{Path: "/usr/bin/gpg",
			Words: []string{"--x" + strings.Repeat("y", i)}})
	}
	if err := refuseUnfitLaunchRules(cfg); err == nil || !strings.Contains(err.Error(), "/usr/bin/gpg") {
		t.Fatalf("17 argv rules must be refused, naming the binary: %v", err)
	}
}
