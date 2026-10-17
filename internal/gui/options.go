package gui

import "github.com/Virgula0/app-listener/internal/wizard"

// Options configures the window.
type Options struct {
	Version string
	// Release is set when the release pipeline stamped Version: the Install page downloads that
	// tag's headless asset. A local build has no release to fetch and asks for a local file.
	Release bool
	Install InstallHooks
}

// InstallMode is an installer entry point (the wizard or one maintenance flag of `install`).
type InstallMode int

const (
	InstallWizard InstallMode = iota
	InstallUpdateCatalog
	InstallDiffCatalog
	InstallRestoreBackups
	InstallDeleteBackups
)

// InstallRequest is one installer run; Binary is the verified headless build to deploy.
type InstallRequest struct {
	Mode                InstallMode
	Binary              string
	Live                bool
	AllowMetadataOutput bool
}

// InstallHooks are the installer and release functions the Install page drives. They live in cmd
// packages, so the caller injects them.
type InstallHooks struct {
	// Run runs req with every prompt on ui.
	Run func(ui wizard.UI, req InstallRequest) error
	// Fetch downloads and verifies the headless release binary of tag into dir.
	Fetch func(tag, dir string) (string, error)
	// FetchLatest is Fetch for the newest stable (or pre-) release; it returns the release's tag.
	FetchLatest func(preRelease bool, dir string) (path, tag string, err error)
	// InstalledVersion is the version the installed binary reports.
	InstalledVersion func() string
}
