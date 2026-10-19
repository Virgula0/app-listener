package gui

import "github.com/Virgula0/app-listener/internal/wizard"

// Options configures the window.
type Options struct {
	Version string
	// Release is set when the release pipeline stamped Version: the Install page downloads that
	// tag's headless asset. A local build has no release to fetch and asks for a local file.
	Release   bool
	Install   InstallHooks
	Update    UpdateHooks
	Uninstall UninstallHooks
}

// InstallMode is an installer entry point (the wizard or one `install` subcommand).
type InstallMode int

const (
	InstallWizard InstallMode = iota
	InstallRefreshCatalog
	InstallDiffCatalog
	InstallRestoreBackups
	InstallDeleteBackups
)

// InstallRequest is one installer run; Binary is the verified headless build to deploy, GUIBinary
// the same release's desktop build ("" = the release has none).
type InstallRequest struct {
	Mode                InstallMode
	Binary, GUIBinary   string
	Live                bool
	AllowMetadataOutput bool
}

// StagedRelease is a downloaded, verified release: Binary is its headless asset, GUI its desktop
// build ("" when the release ships none for this architecture).
type StagedRelease struct {
	Tag, Binary, GUI string
}

// InstallHooks are the installer and release functions the Install page drives. They live in cmd
// packages, so the caller injects them.
type InstallHooks struct {
	// Run runs req with every prompt on ui.
	Run func(ui wizard.UI, req InstallRequest) error
	// Fetch downloads and verifies the headless and desktop builds of tag into dir.
	Fetch func(tag, dir string) (StagedRelease, error)
	// FetchLatest is Fetch for the newest stable (or pre-) release.
	FetchLatest func(preRelease bool, dir string) (StagedRelease, error)
	// InstalledVersion is the version the installed binary reports.
	InstalledVersion func() string
	// InstalledGUIVersion is the installed desktop build's version, "" when none is installed.
	InstalledGUIVersion func() string
}

// UninstallRequest is one uninstall: `uninstall --delete-key`, `--keep-gui`.
type UninstallRequest struct {
	DeleteKey, KeepGUI bool
}

// UninstallHooks is the uninstaller the Install page drives (cmd/functions/uninstall).
type UninstallHooks struct {
	// Run uninstalls with every prompt on ui.
	Run func(ui wizard.UI, req UninstallRequest) error
}

// UpdateHooks are the updater functions the Update page drives (cmd/functions/update).
type UpdateHooks struct {
	// Run updates to the latest stable (or pre-) release with every prompt on ui; withGUI also
	// installs or updates the desktop build.
	Run func(ui wizard.UI, preRelease, withGUI bool) error
	// Latest returns the newest stable (or pre-) release's tag without downloading it, and whether
	// Run would install it over installed.
	Latest func(preRelease bool, installed string) (tag string, newer bool, err error)
}
