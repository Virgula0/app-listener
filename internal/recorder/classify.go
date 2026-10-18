package recorder

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/install"
)

// Recording is a classified recording: what to review before anything is written.
type Recording struct {
	Targets  []Target
	Started  time.Time
	Stopped  time.Time
	Events   uint64
	Dropped  uint64
	Binaries []Binary
}

// Binary is one executable that touched a target, judged by the inode the kernel reported.
type Binary struct {
	// Line is the whitelist line (the executable's path, or a multicall applet's link); "" when
	// the executable can't be whitelisted as recorded (Problem says why).
	Line         string
	Exe          string
	Comms        []string
	Ops          []string
	Mask         []string
	Targets      []string
	UIDs         []uint32
	Count        uint64
	First        time.Time
	Last         time.Time
	SHA256       string
	Owner        uint32
	Problem      string
	Notes        []string
	Suggested    bool
	SystemPlaced bool
	GeneralTool  bool
	ShellHelper  bool
	// Libs are the code mappings an allow_lib would have to cover (Need), or that can't be
	// proposed (Problem); system and guarded-tree libraries are only counted (TrustedLibs).
	Libs        []Lib
	TrustedLibs int
}

// Admissible reports whether b can be whitelisted.
func (b *Binary) Admissible() bool { return b.Line != "" && b.Problem == "" }

// Lib is a library a recorded executable mapped as code.
type Lib struct {
	Path    string
	Why     string
	Problem string
}

// shellHelpers run from the shell's prompt or cd hooks in whatever directory it is in, so they show
// up whenever a recorded path is the cwd. Matched on the basename (or a `name-` prefix, e.g.
// gitstatusd-linux-x86_64); it only clears a default selection, never admits anything.
var shellHelpers = []string{
	"starship", "oh-my-posh", "gitstatusd", "direnv", "zoxide", "atuin", "mcfly", "mise", "rtx",
	"fzf", "carapace",
}

func isShellHelper(paths ...string) bool {
	for _, p := range paths {
		base := strings.TrimRight(filepath.Base(p), "0123456789.")
		for _, h := range shellHelpers {
			if base == h || strings.HasPrefix(base, h+"-") {
				return true
			}
		}
	}
	return false
}

// worldWritableRoots never hold a library or executable worth proposing: anyone can plant there.
var worldWritableRoots = []string{"/tmp", "/var/tmp", "/dev/shm", "/run/user", "/run/shm"}

// Classify verifies every recorded executable and sorts what may be whitelisted from what may not.
// base is the configuration the proposal will be added to (its guarded trees and allow_lib entries
// already cover some libraries).
func Classify(raw *Raw, base *daemonconfig.Config) *Recording {
	rec := &Recording{Targets: raw.Targets, Started: raw.Started, Stopped: raw.Stopped, Events: raw.Events,
		Dropped: raw.Dropped}
	groups := map[exeKey][]RawAccessor{}
	var order []exeKey
	for i := range raw.Accessors {
		a := &raw.Accessors[i]
		k := exeKey{a.Dev, a.Ino}
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}
		groups[k] = append(groups[k], *a)
	}
	for _, k := range order {
		rec.Binaries = append(rec.Binaries, classifyGroup(k, groups[k], raw.Targets, base)...)
	}
	slices.SortStableFunc(rec.Binaries, func(a, b Binary) int {
		if a.Admissible() != b.Admissible() {
			if a.Admissible() {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Line+a.Exe, b.Line+b.Exe)
	})
	return rec
}

// classifyGroup judges one executable inode. A uutils multicall yields one Binary per applet it
// ran as (its whitelist line is the applet's link); anything else is merged into one.
func classifyGroup(k exeKey, group []RawAccessor, targets []Target, base *daemonconfig.Config) []Binary {
	exe := ""
	for i := range group {
		if group[i].Exe != "" {
			exe = group[i].Exe
			break
		}
	}
	v := verifyExe(exe, k)
	if v.problem == "" {
		if refusal := guard.MulticallRefusal(exe); refusal != nil {
			out := make([]Binary, 0, len(group))
			for i := range group {
				b := newBinary(exe, group[i:i+1], v, targets, base)
				b.Line = ""
				if line := appletLine(group[i].ExecName, k); line != "" {
					b.Line = line
				} else {
					b.Problem = refusal.Error()
				}
				finishBinary(&b, v)
				out = append(out, b)
			}
			return out
		}
	}
	b := newBinary(exe, group, v, targets, base)
	finishBinary(&b, v)
	return []Binary{b}
}

// appletLine is name when it is a uutils applet link of the multicall k that the daemon admits as a
// whitelist line, else "".
func appletLine(name string, k exeKey) string {
	if name == "" || !filepath.IsAbs(name) {
		return ""
	}
	f, err := ebpf.OpenConfined(name)
	if err != nil {
		return ""
	}
	key, err := statKey(f)
	f.Close()
	if err != nil || key != k || guard.MulticallRefusal(name) != nil {
		return ""
	}
	return name
}

type exeVerdict struct {
	problem string
	hash    [32]byte
	owner   uint32
	placed  bool
}

// verifyExe checks that path is still the inode k the kernel reported and hashes that inode. A
// path the process was named by is never trusted on its own: only the inode match makes the line
// mean what ran.
func verifyExe(path string, k exeKey) exeVerdict {
	switch {
	case path == "":
		return exeVerdict{problem: "the executable could not be named (it ran from a mount this process " +
			"can't see, or exited before it was looked up)"}
	case strings.HasSuffix(path, " (deleted)"):
		return exeVerdict{problem: "the executable was deleted or replaced since it ran"}
	}
	f, err := ebpf.OpenConfined(path)
	if err != nil {
		return exeVerdict{problem: err.Error()}
	}
	defer f.Close()
	key, err := statKey(f)
	if err != nil {
		return exeVerdict{problem: err.Error()}
	}
	if key != k {
		return exeVerdict{problem: path + " is no longer the file that ran (replaced during or after the recording)"}
	}
	var st unix.Stat_t
	if ferr := unix.Fstat(int(f.Fd()), &st); ferr != nil {
		return exeVerdict{problem: ferr.Error()}
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG {
		return exeVerdict{problem: path + " is not a regular file"}
	}
	if st.Mode&unix.S_IWOTH != 0 {
		return exeVerdict{problem: path + " is world-writable"}
	}
	if dir := worldWritableDir(filepath.Dir(path)); dir != "" {
		return exeVerdict{problem: "it runs from " + dir + ", which anyone can write"}
	}
	entry, err := ebpf.ComputeBinaryEntryFile(f, path)
	if err != nil {
		return exeVerdict{problem: err.Error()}
	}
	return exeVerdict{hash: entry.Hash, owner: st.Uid, placed: ebpf.SystemPlacedInode(path, k.dev, k.ino)}
}

// worldWritableDir is dir (or the world-writable root it lies in) when anyone may add files there.
func worldWritableDir(dir string) string {
	for _, r := range worldWritableRoots {
		if dir == r || strings.HasPrefix(dir, r+"/") {
			return r
		}
	}
	var st unix.Stat_t
	if unix.Stat(dir, &st) == nil && st.Mode&unix.S_IWOTH != 0 {
		return dir
	}
	return ""
}

func newBinary(exe string, group []RawAccessor, v exeVerdict, targets []Target, base *daemonconfig.Config) Binary {
	b := Binary{Line: exe, Exe: exe, Problem: v.problem, Owner: v.owner, SystemPlaced: v.placed}
	if v.problem == "" {
		b.SHA256 = hex.EncodeToString(v.hash[:])
	}
	comms, tgts, libs := map[string]struct{}{}, map[string]struct{}{}, map[string]struct{}{}
	ops := map[ebpf.EventType]struct{}{}
	uids := map[uint32]struct{}{}
	for i := range group {
		a := &group[i]
		for _, c := range a.Comms {
			comms[c] = struct{}{}
		}
		for _, t := range a.Targets {
			tgts[t] = struct{}{}
		}
		for _, l := range a.Libs {
			libs[l] = struct{}{}
		}
		for _, o := range a.Ops {
			ops[o] = struct{}{}
		}
		for _, u := range a.UIDs {
			uids[u] = struct{}{}
		}
		b.Count += a.Count
		if b.First.IsZero() || a.First.Before(b.First) {
			b.First = a.First
		}
		if a.Last.After(b.Last) {
			b.Last = a.Last
		}
	}
	b.Comms, b.Targets, b.UIDs = sortedKeys(comms), sortedKeys(tgts), sortedUIDs(uids)
	b.Ops, b.Mask = opNames(ops), maskFor(ops)
	if v.problem == "" {
		b.Libs, b.TrustedLibs = classifyLibs(sortedKeys(libs), exe, targets, base)
	}
	return b
}

// finishBinary adds the review notes and the default selection: only a root-placed executable that
// isn't a general tool or a shell helper starts selected.
func finishBinary(b *Binary, v exeVerdict) {
	if !b.Admissible() {
		return
	}
	b.GeneralTool = install.ResolvesToGeneralTool(b.Line)
	b.ShellHelper = !b.GeneralTool && isShellHelper(b.Line, b.Exe)
	if !v.placed {
		b.Notes = append(b.Notes, fmt.Sprintf("not root-placed (owner uid %d): its owner, or malware running "+
			"as them, could have put it there; once whitelisted the daemon write-protects it", v.owner))
	}
	if b.GeneralTool {
		b.Notes = append(b.Notes, "a general tool (interpreter, shell, copier): whitelisting it lets every "+
			"script or command run through it reach the path")
	}
	if b.ShellHelper {
		b.Notes = append(b.Notes, "a shell prompt or hook: it runs in whatever directory the shell is in, "+
			"so it likely touched the path only because the shell was there")
	}
	if len(b.UIDs) > 1 {
		b.Notes = append(b.Notes, fmt.Sprintf("ran as %d different users", len(b.UIDs)))
	}
	b.Suggested = v.placed && !b.GeneralTool && !b.ShellHelper
}

// maskFor is the least-privilege event list for what was observed: listing a directory needs READ
// after its OPEN, and STAT accompanies nearly every access (READ, WRITE and MMAP imply OPEN in the
// guard). A hint for review, not a recording of every hook the guard checks.
func maskFor(ops map[ebpf.EventType]struct{}) []string {
	want := map[ebpf.EventType]struct{}{ebpf.EventStat: {}}
	for o := range ops {
		want[o] = struct{}{}
	}
	if _, ok := ops[ebpf.EventOpen]; ok {
		want[ebpf.EventRead] = struct{}{}
	}
	var out []string
	for _, t := range ebpf.EventTypes() {
		if _, ok := want[t]; ok {
			out = append(out, t.String())
		}
	}
	return out
}

// classifyLibs splits exe's code mappings: those the trust guard already accepts (counted), those
// needing allow_lib, and those never proposed.
func classifyLibs(libs []string, exe string, targets []Target, base *daemonconfig.Config) (out []Lib, trusted int) {
	for _, l := range libs {
		if l == exe {
			continue
		}
		if libCovered(l, exe, targets, base) || ebpf.SystemTrusted(l) {
			trusted++
			continue
		}
		lib := Lib{Path: l}
		switch info, err := os.Lstat(l); {
		case err != nil:
			lib.Problem = "gone since it was mapped (generated at run time): not proposed"
		case info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular():
			lib.Problem = "not a regular file: not proposed"
		case worldWritableDir(filepath.Dir(l)) != "":
			lib.Problem = "in a directory anyone can write: never proposed"
		default:
			lib.Why = "not root-owned, and outside every guarded tree: the trust guard refuses it to a " +
				"whitelisted binary unless allow_lib names it"
		}
		out = append(out, lib)
	}
	return out, trusted
}

// libCovered mirrors the trust guard's accepted locations: a guarded tree (a lib_dir only for its
// writers), an existing allow_lib, or one of the targets about to be guarded.
func libCovered(lib, exe string, targets []Target, base *daemonconfig.Config) bool {
	for _, t := range targets {
		if lib == t.Path || (t.IsDir && strings.HasPrefix(lib, t.Path+"/")) {
			return true
		}
	}
	if base == nil {
		return false
	}
	if slices.Contains(base.SharedAllowLibs, lib) {
		return true
	}
	for i := range base.Resources {
		if slices.Contains(base.Resources[i].AllowLibs, lib) {
			return true
		}
	}
	inner := innermostResource(lib, base)
	if inner == nil {
		return false
	}
	return !inner.ReadOnly || slices.ContainsFunc(slices.Concat(inner.Binaries, inner.PendingBinaries),
		func(b daemonconfig.BinaryRule) bool { return b.Path == exe || b.Link == exe })
}

func innermostResource(path string, cfg *daemonconfig.Config) *daemonconfig.Resource {
	var inner *daemonconfig.Resource
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if (path == r.Path || strings.HasPrefix(path, r.Path+"/")) && (inner == nil || len(r.Path) > len(inner.Path)) {
			inner = r
		}
	}
	return inner
}
