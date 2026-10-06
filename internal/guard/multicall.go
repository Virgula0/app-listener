package guard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
)

// Applet identity (guard.bpf.c guard_exec_applet, exe_supersede.h MC_TAG_*): the key of a uutils
// multicall's process carries its attested applet in dev bits 32-47, so each applet is its own
// whitelist identity. Single binaries keep their inode key untouched.
const (
	mcTagShift = 32
	mcTagMask  = uint64(0xffff) << mcTagShift
	mcTagNone  = uint32(0xffff)
)

// Real is k without its applet tag: the file's own key, as supersede marks and trust rows hold it.
func (k GuardInodeKey) Real() GuardInodeKey {
	k.Dev &^= mcTagMask
	return k
}

func (k GuardInodeKey) tagged() bool { return k.Dev&mcTagMask != 0 }

func (k GuardInodeKey) withTag(tag uint32) GuardInodeKey {
	k.Dev = k.Dev&^mcTagMask | uint64(tag)<<mcTagShift
	return k
}

// mcReg: one tag per applet name for the process's lifetime, so an applet keeps its key across a
// multicall's replacement (a package upgrade) and in every build; and each multicall inode's
// applets, from its classification.
var mcReg struct {
	mu     sync.Mutex
	tags   map[string]uint32
	builds map[GuardInodeKey][]string
}

func appletTag(name string) (uint32, error) {
	mcReg.mu.Lock()
	defer mcReg.mu.Unlock()
	if t, ok := mcReg.tags[name]; ok {
		return t, nil
	}
	if mcReg.tags == nil {
		mcReg.tags = make(map[string]uint32)
	}
	t := uint32(len(mcReg.tags)) + 1 //nolint:gosec // bounded below
	if t >= mcTagNone {
		return 0, errors.New("too many distinct multicall applet names")
	}
	mcReg.tags[name] = t
	return t, nil
}

func registerBuild(file GuardInodeKey, applets []string) {
	mcReg.mu.Lock()
	defer mcReg.mu.Unlock()
	if mcReg.builds == nil {
		mcReg.builds = make(map[GuardInodeKey][]string)
	}
	mcReg.builds[file] = applets
}

func buildOf(file GuardInodeKey) ([]string, bool) {
	mcReg.mu.Lock()
	defer mcReg.mu.Unlock()
	a, ok := mcReg.builds[file]
	return a, ok
}

// MulticallError refuses a multicall path whose applet identity can't be established.
type MulticallError struct {
	Path     string
	Family   string
	Reason   string
	Siblings []string // other names of the same inode beside it
}

func (e *MulticallError) Error() string {
	msg := fmt.Sprintf("%s is a %s multicall binary: %s", e.Path, e.Family, e.Reason)
	if len(e.Siblings) > 0 {
		msg += fmt.Sprintf(" (the same file also runs as: %s)", strings.Join(e.Siblings, ", "))
	}
	return msg + "; whitelist a dedicated binary instead (on Ubuntu, the GNU gnu* tools or coreutils-from-gnu)"
}

// IsMulticallRefusal reports whether err is a MulticallError.
func IsMulticallRefusal(err error) bool {
	var me *MulticallError
	return errors.As(err, &me)
}

// appletKeyOrRefuse is ExeKey for a live (re-)admission of path: a refusal is logged once per inode
// and the caller admits nothing.
func (g *Guard) appletKeyOrRefuse(path string, f *os.File, file GuardInodeKey) (GuardInodeKey, bool) {
	exe, err := ExeKey(path, f, file)
	if err != nil {
		if g.refused.firstTime(path, file) {
			log.Errorf("guard %s: %s not admitted: %v", g.path, logging.SanitizeText(path), err)
		}
		return file, false
	}
	return exe, true
}

// ExeKey is the key the whitelist rows for the binary at path use, file being the key of the inode f
// holds (path is opened confined when f is nil): file itself for a single binary; file tagged with
// the applet path names for a uutils multicall, where path's basename and the file it reaches must both
// be that applet of the build's own list. Any other multicall is refused (MulticallError): one
// identity would admit every applet.
func ExeKey(path string, f *os.File, file GuardInodeKey) (GuardInodeKey, error) {
	if f == nil {
		var err error
		if f, err = ebpf.OpenConfined(path); err != nil {
			return file, err
		}
		defer f.Close()
		dev, ino, err := ebpf.StatFile(f)
		if err != nil {
			return file, err
		}
		if (GuardInodeKey{Dev: dev, Ino: ino}) != file {
			return file, fmt.Errorf("%s changed while being resolved", path)
		}
	}
	m, err := ebpf.ClassifyMulticall(f)
	if err != nil {
		return file, fmt.Errorf("classifying %s: %w", path, err)
	}
	switch m.Kind {
	case ebpf.SingleBinary:
		return file, nil
	case ebpf.OpaqueMulticall:
		return file, &MulticallError{Path: path, Family: m.Family, Siblings: siblingNames(f),
			Reason: "its applets cannot be told apart, so whitelisting one would admit them all"}
	}
	resolved, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return file, err
	}
	name := filepath.Base(path)
	if filepath.Base(resolved) != name || !slices.Contains(m.Applets, name) {
		return file, &MulticallError{Path: path, Family: m.Family, Siblings: siblingNames(f),
			Reason: fmt.Sprintf("%q is not one of its applets (the name must be an applet's own link)", name)}
	}
	tag, err := appletTag(name)
	if err != nil {
		return file, err
	}
	registerBuild(file, m.Applets)
	return file.withTag(tag), nil
}

// siblingNames lists up to 8 other hard links of f's inode in its directory, for a refusal message.
func siblingNames(f *os.File) []string {
	resolved, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return nil
	}
	dev, ino, err := ebpf.StatFile(f)
	if err != nil {
		return nil
	}
	entries, err := os.ReadDir(filepath.Dir(resolved))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		p := filepath.Join(filepath.Dir(resolved), e.Name())
		if p == resolved || !e.Type().IsRegular() {
			continue
		}
		if d, i, serr := ebpf.StatInode(p); serr == nil && d == dev && i == ino {
			if len(out) == 8 {
				return append(out, "…")
			}
			out = append(out, e.Name())
		}
	}
	return out
}

// ensureMulticall writes the applet names of file (a multicall's own key), then its guard_multicall
// row, before any tagged row of it: a process of a multicall missing from guard_multicall would be
// keyed by its inode.
func (e *engine) ensureMulticall(file GuardInodeKey) error {
	applets, ok := buildOf(file)
	if !ok {
		return fmt.Errorf("multicall inode %d was never classified", file.Ino)
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.mcNames == nil {
		e.mcNames = make(map[GuardInodeKey][]string)
	}
	if prev, done := e.mcNames[file]; !done || !slices.Equal(prev, applets) {
		if err := e.syncMulticallNamesLocked(file, prev, applets); err != nil {
			return err
		}
		e.mcNames[file] = applets
	}
	return e.objs.GuardMulticall.Put(file, uint32(1))
}

// syncMulticallNamesLocked writes applets and drops the names of prev (an earlier classification of
// the same number) that applets lacks.
func (e *engine) syncMulticallNamesLocked(file GuardInodeKey, prev, applets []string) error {
	for _, a := range applets {
		tag, err := appletTag(a)
		if err != nil {
			return err
		}
		if err := e.objs.GuardMulticallNames.Put(mcNameKey(file, a), tag); err != nil {
			return fmt.Errorf("writing multicall applet %s: %w", a, err)
		}
	}
	for _, a := range prev {
		if slices.Contains(applets, a) {
			continue
		}
		err := e.objs.GuardMulticallNames.Delete(mcNameKey(file, a))
		if err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return fmt.Errorf("dropping multicall applet %s: %w", a, err)
		}
	}
	return nil
}

func mcNameKey(file GuardInodeKey, name string) GuardMcNameKey {
	k := GuardMcNameKey{Exe: file}
	for i := 0; i < len(name) && i < len(k.Name)-1; i++ {
		k.Name[i] = int8(name[i]) //nolint:gosec // applet names are ASCII (validAppletName)
	}
	return k
}

// withRealKeys adds, to each resource's allows, the real key of every applet key it holds: the
// exec-time hooks (bprm_check_security, bprm_committed_creds, the supersede hooks) see the file
// before its applet is attested.
func withRealKeys(allows map[uint32]map[GuardInodeKey]uint8) map[uint32]map[GuardInodeKey]uint8 {
	out := make(map[uint32]map[GuardInodeKey]uint8, len(allows))
	for res, set := range allows {
		cp := make(map[GuardInodeKey]uint8, len(set))
		for k, a := range set {
			cp[k] = a
		}
		for k, a := range set {
			if k.tagged() {
				if _, ok := cp[k.Real()]; !ok {
					cp[k.Real()] = a
				}
			}
		}
		out[res] = cp
	}
	return out
}

// singleBinaryPaths drops the paths reaching a multicall from process-side trust rights (updater,
// writer, loader bits): keyed by the file, they would extend to every applet. A path that won't open
// stays, for its resolver to skip.
func singleBinaryPaths(paths map[string]uint64, what string) map[string]uint64 {
	out := make(map[string]uint64, len(paths))
	for p, bits := range paths {
		f, err := ebpf.OpenConfined(p)
		if err != nil {
			out[p] = bits
			continue
		}
		m, cerr := ebpf.ClassifyMulticall(f)
		f.Close()
		if cerr != nil || m.Kind != ebpf.SingleBinary {
			log.Warnf("trust guard: %s %s is not granted %s rights: a multicall binary's applets share its file",
				logging.SanitizeText(p), multicallFamily(m, cerr), what)
			continue
		}
		out[p] = bits
	}
	return out
}

func multicallFamily(m ebpf.Multicall, err error) string {
	if err != nil {
		return "(unclassifiable: " + err.Error() + ")"
	}
	return "(" + m.Family + " multicall)"
}
