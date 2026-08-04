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
	Version = "v0.1.0"
)
