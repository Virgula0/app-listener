package daemon

import (
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
)

// launchRulesByPath maps each whitelisted binary's resolved path to its launch rules: the built-ins
// of its config line's and its target's basenames, and the [launch_rules] naming either path.
func launchRulesByPath(cfg *daemonconfig.Config) map[string][]daemonconfig.LaunchRule {
	out := make(map[string][]daemonconfig.LaunchRule)
	used := make([]bool, len(cfg.LaunchRules))
	for i := range cfg.Resources {
		for _, list := range [][]daemonconfig.BinaryRule{cfg.Resources[i].Binaries, cfg.Resources[i].PendingBinaries} {
			for _, b := range list {
				if _, done := out[b.Path]; !done && !b.LibBinary {
					out[b.Path] = binaryLaunchRules(cfg, b, used)
				}
			}
		}
	}
	for i, r := range cfg.LaunchRules {
		if !used[i] {
			log.Warnf("daemon: [launch_rules] %s is not a whitelisted binary: its rule does nothing",
				logging.SanitizeText(r.Path))
		}
	}
	for p, rules := range out {
		if len(rules) == 0 {
			delete(out, p)
		}
	}
	return out
}

// binaryLaunchRules lists b's rules, marking in used the [launch_rules] lines it takes.
func binaryLaunchRules(cfg *daemonconfig.Config, b daemonconfig.BinaryRule, used []bool) []daemonconfig.LaunchRule {
	rules := guard.BuiltinLaunchRules(filepath.Base(b.Path))
	if b.Link != "" && filepath.Base(b.Link) != filepath.Base(b.Path) {
		rules = append(rules, guard.BuiltinLaunchRules(filepath.Base(b.Link))...)
	}
	for i, r := range cfg.LaunchRules {
		if r.Path == b.Path || (b.Link != "" && r.Path == b.Link) {
			rules = append(rules, r)
			used[i] = true
		}
	}
	return rules
}

// allowedEnv returns the built-in and allow_env NAME=/path strings whose path only root decides
// (ebpf.RootPlacedName): a value a user can fill would let a caller pick the CA bundle again.
func allowedEnv(cfg *daemonconfig.Config) []string {
	var out []string
	for _, a := range guard.BuiltinAllowedEnv() {
		if _, p, _ := strings.Cut(a, "="); ebpf.RootPlacedName(p) == nil {
			out = append(out, a)
		}
	}
	for _, a := range cfg.AllowEnv {
		_, p, _ := strings.Cut(a, "=")
		if err := ebpf.RootPlacedName(p); err != nil {
			log.Warnf("daemon: [launch_rules] allow_env %s refused, launches setting it stay marked: %v",
				logging.SanitizeText(a), logging.SanitizeText(err.Error()))
			continue
		}
		out = append(out, a)
	}
	return out
}
