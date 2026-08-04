package constants

type VerboseLevel int

const (
	NeededInfo VerboseLevel = iota
	InternalWarningsLevelTwo
	PrintAdditionalInfoLevelOne
	PrintAdditionalInfoLevelTwo
)

const (
	AppName = "app-listener"
	Version = "v0.1.0"
)
