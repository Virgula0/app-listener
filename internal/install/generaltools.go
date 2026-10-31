package install

import (
	"path/filepath"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// generalTools run attacker-chosen code or copy attacker-chosen bytes from their arguments (`git`,
// `node -e`, `java -cp`, `cp`). Identity is the exe inode, so a resource whitelisting one is
// readable by any script run through it (issue #79). Matched on the basename, also without a
// trailing version (python3.12, gawk-5.4.0) or a Debian alternative's suffix (nc.openbsd, vim.basic).
var generalTools = []string{
	"sh", "bash", "dash", "zsh", "fish", "ksh", "mksh", "csh", "tcsh", "ash", "yash",
	"busybox", "toybox", "coreutils", "env",
	"git", "node", "nodejs", "npm", "npx", "bun", "deno",
	"python", "perl", "ruby", "irb", "php", "lua", "luajit", "java", "tclsh", "wish", "expect", "guile",
	"Rscript", "R", "julia", "pwsh", "gdb", "m4", "make", "sqlite3",
	"cp", "mv", "install", "rsync", "tar", "curl", "wget",
	"cat", "tac", "head", "tail", "dd", "od", "xxd", "hexdump", "hd", "strings", "base64", "base32",
	"basenc", "sed", "awk", "gawk", "mawk", "grep", "egrep", "fgrep", "rg", "tee", "less", "more",
	"sort", "uniq", "nl", "cut", "paste", "fold", "fmt", "pr", "rev", "expand", "unexpand", "join",
	"comm", "diff", "split", "csplit", "shuf", "column", "iconv", "look", "tsort", "ptx", "bat", "batcat",
	"vi", "vim", "view", "ex", "ed", "red", "nvim", "nano", "emacs", "openssl", "nc", "ncat", "netcat",
	"socat", "jq", "gzip", "gunzip", "zcat", "xz", "xzcat", "zstd", "zstdcat", "bzip2", "bzcat", "zip",
	"unzip", "7z", "cpio", "pv",
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
	base := filepath.Base(path)
	unversioned := func(n string) string { return strings.TrimRight(strings.TrimRight(n, "0123456789."), "-_") }
	alt, _, _ := strings.Cut(base, ".")
	for _, n := range []string{base, unversioned(base), alt, unversioned(alt)} {
		if slices.Contains(generalTools, n) {
			return true
		}
	}
	return false
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

// ScriptsOf lists res's whitelist lines that are "#!" scripts: the exe is their interpreter, so
// such a line never matches.
func ScriptsOf(res *daemonconfig.Resource) []string {
	var out []string
	for _, b := range res.Binaries {
		line := b.Path
		if b.Link != "" {
			line = b.Link
		}
		if ebpf.IsScript(b.Path) && !slices.Contains(out, line) {
			out = append(out, line)
		}
	}
	return out
}

// WarnWhitelistLines logs, per resource, the general tools that make it readable by any script and
// the script lines that never match.
func WarnWhitelistLines(cfg *daemonconfig.Config) {
	for i := range cfg.Resources {
		res := &cfg.Resources[i]
		if scripts := ScriptsOf(res); len(scripts) > 0 {
			log.Warnf("%s whitelists script(s) %s, which never match: a #! script runs as its interpreter "+
				"(never whitelist that), so the app is denied too; whitelist the ELF it ends up running, if any",
				res.Path, strings.Join(scripts, ", "))
		}
		if res.ReadOnly {
			continue // a lib_dir holds no secret: a tool there is a writer, judged by the updater plan
		}
		if tools := GeneralToolsOf(res); len(tools) > 0 {
			log.Warnf("%s whitelists general tool(s) %s: identity is the executable, so any script or "+
				"arguments run through them can read it (not an updater either)", res.Path, strings.Join(tools, ", "))
		}
	}
}
