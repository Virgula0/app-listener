package daemon

import (
	"slices"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/install"
)

// mayUpdate judges the binary the path resolves to, since updater rights attach to its inode. A
// general tool (install.ResolvesToGeneralTool) never updates. A multi-call binary (uutils
// coreutils hard-linked under every applet name, busybox behind symlinks) runs whichever applet
// argv[0] names, and the caller picks argv[0] (`exec -a cp`), so whitelisting its harmless `date`
// would make it `cp`: more than one hard link counts as a general tool too. A path that doesn't
// resolve yet (a deferred binary) is judged by its name. The hard-link rule alone is waived for a
// user-writable lib_binary: pressure-vessel hard-links its runtime into every var/tmp-XXXXXX and
// must delete those links. Never for a root-owned binary (distro coreutils, busybox).
func mayUpdate(path string, libBinary bool) bool {
	if install.ResolvesToGeneralTool(path) {
		return false
	}
	return !hardLinked(path) || (libBinary && !ebpf.SystemTrusted(path))
}

func hardLinked(path string) bool {
	var st unix.Stat_t
	return unix.Stat(path, &st) == nil && st.Nlink > 1
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
		// A link line owns the bit too: a re-pointed link's new target is judged by the link
		// (guard.resyncLink -> AllowReplacement). It resolves to a target above, so no new row.
		for _, l := range resourceLinks(res) {
			p.Owners[l] |= bit
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

func resourceLinks(res *daemonconfig.Resource) []string {
	var out []string
	for _, b := range res.Binaries {
		if b.Link != "" {
			out = append(out, b.Link)
		}
	}
	return out
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
