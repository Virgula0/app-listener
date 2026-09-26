package constants

import "errors"

// ErrCriticalStartup marks a startup failure that reproduces identically on every restart
// (daemon.conf fails to load/parse, or a resource's on-disk encryption state no longer matches
// need_encryption; issue #53: a backup restore or other out-of-band change) as opposed to a
// transient condition (BPF LSM not ready, stale pins from a just-crashed instance, an I/O hiccup)
// worth retrying. Wrap with %w where the permanent cause is identified; main.go checks errors.Is
// and exits CriticalExitCode, which the systemd unit's RestartPreventExitStatus
// (internal/install/daemon-samples/app-listener-daemon.service) uses to stop instead of
// crash-looping.
var ErrCriticalStartup = errors.New("critical startup failure: administrator action required")

// CriticalExitCode is ErrCriticalStartup's exit status: 78 = EX_CONFIG (sysexits.h), "configuration
// problem, not runtime".
const CriticalExitCode = 78

// VerboseLevel is the CLI verbosity ladder (--verbose); logging.VerboseToLogrus maps it onto logrus
// for the headless console and --dump-log:
//
//	0 ErrorsOnly                  -> error (quiet: denials/errors only)
//	1 NeededInfo                  -> warn  (essential lines)
//	2 InternalWarningsLevelTwo    -> info  (historical default)
//	3 PrintAdditionalInfoLevelOne -> debug (+ internal details)
//	4 PrintAdditionalInfoLevelTwo -> trace (everything)
type VerboseLevel int

const (
	ErrorsOnly VerboseLevel = iota
	NeededInfo
	InternalWarningsLevelTwo
	PrintAdditionalInfoLevelOne
	PrintAdditionalInfoLevelTwo
)

const (
	AppName = "App-Listener"
)

// Version is set at release via -ldflags "-X .../constants.Version=pre-<date>-<sha>" (a var so
// ldflags can inject). Local builds keep "dev", which the updater treats as older than any
// pre-release.
var Version = "v0.2.0"
