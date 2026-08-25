package constants

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
var Version = "v0.1.0"
