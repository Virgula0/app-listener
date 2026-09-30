package daemon

import (
	"path/filepath"
	"slices"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// generalTools run attacker-chosen code or copy attacker-chosen bytes from their arguments, so
// whitelisting one never makes it an updater: `git`, `node -e`, `cp` would otherwise rewrite the
// app binaries sharing its resource. Matched on the basename; a trailing version (python3.12) too.
var generalTools = []string{
	"sh", "bash", "dash", "zsh", "fish", "ksh", "busybox", "toybox", "coreutils", "env",
	"git", "node", "nodejs", "npm", "npx", "bun", "deno",
	"python", "perl", "ruby", "php", "lua", "java",
	"cp", "mv", "install", "rsync", "tar", "curl", "wget",
}

// mayUpdate judges the binary the path resolves to, since updater rights attach to its inode: a
// multi-call binary (uutils coreutils hard-linked under every applet name, busybox behind symlinks)
// runs whichever applet argv[0] names, and the caller picks argv[0] (`exec -a cp`), so whitelisting
// its harmless `date` would make it `cp`. More than one hard link, or a symlink target named like a
// general tool, therefore counts as one. A path that doesn't resolve yet (a deferred binary) is
// judged by its name. The hard-link rule alone is waived for a user-writable lib_binary:
// pressure-vessel hard-links its runtime into every var/tmp-XXXXXX and must delete those links.
// Never for a root-owned binary (distro coreutils, busybox).
func mayUpdate(path string, libBinary bool) bool {
	if resolvesToNamedTool(path) {
		return false
	}
	return !hardLinked(path) || (libBinary && !ebpf.SystemTrusted(path))
}

func hardLinked(path string) bool {
	var st unix.Stat_t
	return unix.Stat(path, &st) == nil && st.Nlink > 1
}

func resolvesToNamedTool(path string) bool {
	if namedGeneralTool(path) {
		return true
	}
	resolved, err := filepath.EvalSymlinks(path)
	return err == nil && namedGeneralTool(resolved)
}

func namedGeneralTool(path string) bool {
	return slices.Contains(generalTools, strings.TrimRight(filepath.Base(path), "0123456789."))
}

// planUpdaters assigns one bit per distinct resource binary set: a resource's binaries and
// allow_libs own its bit, and its binaries, general tools excepted (see mayUpdate), are its
// updaters. Past 64 sets the rest own no bit, so nobody may modify their binaries (fail closed).
// [libraries] allow_libs belong to no resource and get no owner either.
func planUpdaters(cfg *daemonconfig.Config) guard.UpdaterPlan {
	p := guard.UpdaterPlan{Owners: map[string]uint64{}, Updaters: map[string]uint64{}}
	bitOf := map[string]uint64{}
	for i := range cfg.Resources {
		res := &cfg.Resources[i]
		set := resourceBinaries(res)
		if len(set) == 0 {
			continue
		}
		key := strings.Join(set, "\x00")
		bit, ok := bitOf[key]
		if !ok {
			if len(bitOf) >= 64 {
				log.Warnf("trust guard: too many distinct whitelists, binaries of %s cannot be updated "+
					"in place by any process until the whitelists are merged", res.Path)
				continue
			}
			bit = uint64(1) << len(bitOf)
			bitOf[key] = bit
		}
		libBins := libBinaries(res)
		for _, b := range set {
			p.Owners[b] |= bit
			if mayUpdate(b, libBins[b]) {
				p.Updaters[b] |= bit
			}
		}
		for _, l := range append(slices.Clone(res.AllowLibs), res.PendingLibs...) {
			p.Owners[l] |= bit
		}
		if tools := generalToolsToWarn(set); len(tools) > 0 {
			log.Warnf("trust guard: %s whitelists general tool(s) %s beside user-writable binaries — "+
				"they are not allowed to update those binaries, but any arguments given to them reach %s",
				res.Path, strings.Join(tools, ", "), res.Path)
		}
	}
	return p
}

func resourceBinaries(res *daemonconfig.Resource) []string {
	var set []string
	for _, list := range [][]daemonconfig.BinaryRule{res.Binaries, res.PendingBinaries} {
		for _, b := range list {
			set = append(set, b.Path)
		}
	}
	sort.Strings(set)
	return slices.Compact(set)
}

func libBinaries(res *daemonconfig.Resource) map[string]bool {
	out := map[string]bool{}
	for _, list := range [][]daemonconfig.BinaryRule{res.Binaries, res.PendingBinaries} {
		for _, b := range list {
			if b.LibBinary {
				out[b.Path] = true
			}
		}
	}
	return out
}

// generalToolsToWarn lists the general tools whitelisted beside a user-writable binary: no updater,
// but they still read the resource's secrets on attacker-chosen arguments. A hard link alone isn't
// warned about: runtimes (Steam's pressure-vessel) hard-link ordinary binaries; mayUpdate judges
// those.
func generalToolsToWarn(set []string) []string {
	var tools []string
	userWritable := false
	for _, b := range set {
		switch {
		case resolvesToNamedTool(b):
			tools = append(tools, b)
		case !ebpf.SystemTrusted(b):
			userWritable = true
		}
	}
	if !userWritable {
		return nil
	}
	return tools
}
