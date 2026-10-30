package daemon

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
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
