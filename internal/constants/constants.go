package constants

import "errors"

// ErrCriticalStartup marks a daemon startup failure that will reproduce
// identically on every restart — daemon.conf failing to load or parse, or a
// configured resource whose on-disk encryption state no longer matches
// need_encryption (see issue #53: a backup restore, or any out-of-band
// change, leaving daemon.conf and reality out of sync) — as opposed to a
// transient environment condition (BPF LSM not yet ready, stale BPF pins
// from a just-crashed instance, a momentary I/O hiccup) that IS worth
// systemd retrying. Wrap it with %w at the point the permanent cause is
// identified; main.go checks for it with errors.Is and exits with
// CriticalExitCode instead of the default retryable status, paired with
// RestartPreventExitStatus in the systemd unit (see
// internal/install/daemon-samples/app-listener-daemon.service) so systemd
// stops trying instead of crash-looping every RestartSec forever for a
// condition no restart can ever fix.
var ErrCriticalStartup = errors.New("critical startup failure: administrator action required")

// CriticalExitCode is the process exit status for ErrCriticalStartup — 78 is
// EX_CONFIG from sysexits.h, the conventional "configuration problem, not a
// runtime one" status.
const CriticalExitCode = 78

// VerboseLevel is the CLI-wide verbosity ladder driven by --verbose; logging.VerboseToLogrus
// maps it onto logrus for both the headless console stream and the --dump-log file:
//
//	0 ErrorsOnly                  -> error (quiet mode: denials/errors only)
//	1 NeededInfo                  -> warn   (essential lines)
//	2 InternalWarningsLevelTwo    -> info   (= historical default display)
//	3 PrintAdditionalInfoLevelOne -> debug  (+ internal details)
//	4 PrintAdditionalInfoLevelTwo -> trace  (everything)
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

// Version is overridden at build time by the release workflow via
// -ldflags "-X .../constants.Version=pre-<date>-<sha>" (a var, not a const,
// enables ldflags injection). Local builds keep the dev marker, which the
// updater treats as older than any pre-release.
var Version = "v0.2.0"
