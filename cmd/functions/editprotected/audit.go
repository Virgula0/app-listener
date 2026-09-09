package editprotected

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// auditAfterEdit loads the installed daemon.conf (if any) and audits the
// edited tree. Used by the offline flow, which has no Config in hand.
func auditAfterEdit(resource string) {
	cfg, err := daemonconfig.Load(systemd.SystemConfigPath)
	if err != nil {
		auditAfterEditWithConfig(nil, resource)
		return
	}
	auditAfterEditWithConfig(cfg, resource)
}

// auditAfterEditWithConfig inspects the edited resource tree for things that
// would sit outside the daemon's protection or weaken it, and — if it finds
// any — prints them and asks the operator to acknowledge before returning.
// Advisory only: nothing is changed or blocked.
func auditAfterEditWithConfig(cfg *daemonconfig.Config, resource string) {
	findings := make([]string, 0, 8)
	findings = append(findings, auditTree(resource)...)
	findings = append(findings, auditGroupCoverage(cfg, resource)...)

	if len(findings) == 0 {
		log.Info("post-edit audit: no issues found")
		return
	}

	log.Warn("post-edit audit found potential issues:")
	for _, f := range findings {
		log.Warnf("  - %s", f)
	}

	// A one-way acknowledgement, not a decision: the edit is already written
	// and nothing here changes or reverts it — the prompt only makes sure the
	// operator saw the warnings.
	_ = huh.NewForm(huh.NewGroup(
		huh.NewNote().
			Title("Post-edit audit — review the warnings above").
			Description("Advisory only: the edit was saved and nothing here is changed or blocked.\n\n" +
				strings.Join(bullet(findings), "\n")).
			Next(true).
			NextLabel("OK"),
	)).Run()
}

// bullet prefixes each finding with "  • " for the acknowledgement note.
func bullet(findings []string) []string {
	out := make([]string, len(findings))
	for i, f := range findings {
		out[i] = "  • " + f
	}
	return out
}

// auditTree walks resource and flags symlinks, group/other-accessible files,
// and executable files.
func auditTree(resource string) []string {
	root, err := filepath.Abs(resource)
	if err != nil {
		root = resource
	}
	findings := make([]string, 0, 8)

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil //nolint:nilerr // skip unreadable entries, keep walking
		}
		info, infoErr := d.Info()
		if infoErr != nil {
			return nil //nolint:nilerr // entry vanished mid-walk, keep walking
		}
		if f := auditEntry(root, path, info); f != "" {
			findings = append(findings, f)
		}
		return nil
	})
	return findings
}

// auditEntry classifies one filesystem entry; "" means nothing noteworthy.
func auditEntry(root, path string, info os.FileInfo) string {
	rel := relOrPath(root, path)

	if info.Mode()&os.ModeSymlink != 0 {
		target, _ := os.Readlink(path)
		resolved := target
		if !filepath.IsAbs(resolved) {
			resolved = filepath.Join(filepath.Dir(path), target)
		}
		if !within(root, resolved) {
			return "symlink " + rel + " points OUTSIDE the guarded tree -> " + target +
				" (the daemon guards the link, not its target)"
		}
		return "symlink " + rel + " -> " + target + " (symlinks in a protected tree are unusual)"
	}
	if !info.Mode().IsRegular() {
		return ""
	}
	if info.Mode().Perm()&0o077 != 0 {
		return "file " + rel + " is readable/writable by group or others (mode " +
			info.Mode().Perm().String() + ") — protected content should be 0600"
	}
	if info.Mode().Perm()&0o111 != 0 {
		return "file " + rel + " is executable — it will not be on any whitelist and could be a planted helper"
	}
	return ""
}

// auditGroupCoverage flags entries under a grouped encryption root that fall
// outside every guarded watch path of that group: they are encrypted at rest
// but nothing guards access to them.
func auditGroupCoverage(cfg *daemonconfig.Config, resource string) []string {
	if cfg == nil {
		return nil
	}
	self := findResource(cfg, resource)
	if self == nil || self.EncryptionRoot == "" {
		return nil // ungrouped: the watch path is its own root, fully covered
	}
	watchPaths := groupWatchPaths(cfg, self.EncryptionRoot)

	findings := make([]string, 0, 4)
	seen := map[string]bool{}
	_ = filepath.WalkDir(self.EncryptionRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || path == self.EncryptionRoot {
			return nil //nolint:nilerr // skip unreadable entries, keep walking
		}
		if coveredByWatch(path, watchPaths) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		head := strings.SplitN(relOrPath(self.EncryptionRoot, path), string(os.PathSeparator), 2)[0]
		if !seen[head] {
			seen[head] = true
			findings = append(findings, "entry "+head+" under the encryption root "+self.EncryptionRoot+
				" is outside every guarded watch path — encrypted at rest but access is NOT guarded")
		}
		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	return findings
}

func findResource(cfg *daemonconfig.Config, path string) *daemonconfig.Resource {
	for i := range cfg.Resources {
		if cfg.Resources[i].Path == path {
			return &cfg.Resources[i]
		}
	}
	return nil
}

func groupWatchPaths(cfg *daemonconfig.Config, encryptionRoot string) []string {
	var out []string
	for i := range cfg.Resources {
		if cfg.Resources[i].EncryptionRoot == encryptionRoot {
			out = append(out, cfg.Resources[i].Path)
		}
	}
	return out
}

func coveredByWatch(path string, watchPaths []string) bool {
	for _, wp := range watchPaths {
		if path == wp || within(wp, path) {
			return true
		}
	}
	return false
}

func relOrPath(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil {
		return rel
	}
	return path
}

// within reports whether path is at or below dir.
func within(dir, path string) bool {
	dirClean := filepath.Clean(dir)
	p := filepath.Clean(path)
	return p == dirClean || strings.HasPrefix(p, dirClean+string(os.PathSeparator))
}
