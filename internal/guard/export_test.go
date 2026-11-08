package guard

import "github.com/Virgula0/app-listener/internal/daemonconfig"

const ArgvRootOnly = argvRootOnly

// EncodeLaunchRules encodes rules as a whitelisted binary gets them (caEnvRules, BuiltinAllowedEnv),
// returning the argv word entries by word.
func EncodeLaunchRules(rules []daemonconfig.LaunchRule) (map[string]uint64, GuardTrustArgvSet, error) {
	canon, allowed, _ := canonLaunchRules(append(rules, caEnvRules...), BuiltinAllowedEnv())
	entries, set, err := encodeLaunchSet(canon, allowed)
	words := map[string]uint64{}
	for e, v := range entries {
		if e.kind == argvWord {
			words[e.name] = v
		}
	}
	return words, set, err
}

func ShortRefused(set GuardTrustArgvSet, word string) bool { return shortHit(set, word) }
