package constants

type VerboseLevel int

const (
	NeededInfo VerboseLevel = iota
	InternalWarningsLevelTwo
	PrintAdditionalInfoLevelOne
	PrintAdditionalInfoLevelTwo
)

const (
	AppName = "App-Listener"
)

// Version is overridden at build time by the release workflow via
// -ldflags "-X .../constants.Version=pre-<date>-<sha>"; a var (not a
// const) is what makes ldflags injection possible. Local builds keep the
// default dev marker, which the updater treats as older than any
// pre-release.
var Version = "v0.1.0"
