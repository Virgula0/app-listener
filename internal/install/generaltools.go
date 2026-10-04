package install

import (
	"path/filepath"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// generalTools run attacker-chosen code or copy attacker-chosen bytes from their arguments (`git`,
// `node -e`, `java -cp`, `cp`). Identity is the exe inode, so a resource whitelisting one is
// readable by any script run through it (issue #79). Matched on the basename; a trailing version
// (python3.12) too.
var generalTools = []string{
	"sh", "bash", "dash", "zsh", "fish", "ksh", "busybox", "toybox", "coreutils", "env",
	"git", "node", "nodejs", "npm", "npx", "bun", "deno",
	"python", "perl", "ruby", "php", "lua", "java",
	"cp", "mv", "install", "rsync", "tar", "curl", "wget",
}

// ResolvesToGeneralTool reports whether path, or the file it resolves to, is named like a general
// tool.
func ResolvesToGeneralTool(path string) bool {
	if namedGeneralTool(path) {
		return true
	}
	resolved, err := filepath.EvalSymlinks(path)
	return err == nil && namedGeneralTool(resolved)
}

func namedGeneralTool(path string) bool {
	return slices.Contains(generalTools, strings.TrimRight(filepath.Base(path), "0123456789."))
}

// GeneralToolsOf lists the general tools res whitelists.
func GeneralToolsOf(res *daemonconfig.Resource) []string {
	var out []string
	for _, list := range [][]daemonconfig.BinaryRule{res.Binaries, res.PendingBinaries} {
		for _, b := range list {
			if ResolvesToGeneralTool(b.Path) && !slices.Contains(out, b.Path) {
				out = append(out, b.Path)
			}
		}
	}
	return out
}

// WarnGeneralTools logs, per resource, the general tools that make it readable by any script.
func WarnGeneralTools(cfg *daemonconfig.Config) {
	for i := range cfg.Resources {
		res := &cfg.Resources[i]
		if res.ReadOnly {
			continue // a lib_dir holds no secret: a tool there is a writer, judged by the updater plan
		}
		if tools := GeneralToolsOf(res); len(tools) > 0 {
			log.Warnf("%s whitelists general tool(s) %s: identity is the executable, so any script or "+
				"arguments run through them can read it (not an updater either)", res.Path, strings.Join(tools, ", "))
		}
	}
}
