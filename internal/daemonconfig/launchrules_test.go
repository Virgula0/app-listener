package daemonconfig

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestLaunchRulesBlockParsed(t *testing.T) {
	cfg, err := Load(writeConfig(t, "[launch_rules]\ndeny_args /usr/bin/gh repo   delete\n"+
		"warn_args \"/opt/my tool/cli\" export --all\ndeny_env /usr/bin/kubectl HTTPS_PROXY\n"+
		"warn_env /usr/bin/x/../kubectl KUBECONFIG\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := []LaunchRule{
		{Path: "/usr/bin/gh", Words: []string{"repo", "delete"}},
		{Path: "/opt/my tool/cli", Words: []string{"export", "--all"}, Warn: true},
		{Path: "/usr/bin/kubectl", Env: "HTTPS_PROXY"},
		{Path: "/usr/bin/kubectl", Env: "KUBECONFIG", Warn: true},
	}
	if !slices.EqualFunc(cfg.LaunchRules, want, LaunchRule.Same) {
		t.Fatalf("LaunchRules = %+v", cfg.LaunchRules)
	}
}

func TestLaunchRulesBlockRejectsMalformedLines(t *testing.T) {
	for name, line := range map[string]string{
		"unknown directive": "allow_args /usr/bin/gh auth",
		"relative":          "deny_args bin/gh auth token",
		"no words":          "deny_args /usr/bin/gh",
		"five words":        "deny_args /usr/bin/gh a b c d e",
		"flag with value":   "deny_args /usr/bin/gh --hostname=x",
		"long word":         "deny_args /usr/bin/gh " + strings.Repeat("w", MaxLaunchWordLen+1),
		"two env names":     "deny_env /usr/bin/gh A B",
		"bad env name":      "deny_env /usr/bin/gh 1PROXY",
		"env with value":    "deny_env /usr/bin/gh A=b",
		"duplicate":         "deny_args /usr/bin/gh auth token\ndeny_args /usr/bin/gh auth token",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Load(writeConfig(t, "[launch_rules]\n"+line+"\n")); err == nil {
				t.Fatalf("%q must be refused", line)
			}
		})
	}
	var b strings.Builder
	b.WriteString("[launch_rules]\ndeny_env /usr/bin/gh A\n")
	for i := 0; i <= MaxLaunchRules; i++ {
		fmt.Fprintf(&b, "deny_args /usr/bin/gh w%d\n", i)
	}
	if _, err := Load(writeConfig(t, b.String())); err == nil {
		t.Fatal("more than MaxLaunchRules argv rules for one binary must be refused")
	}
}
