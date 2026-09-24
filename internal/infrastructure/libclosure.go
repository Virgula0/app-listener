package ebpf

import (
	"bufio"
	"bytes"
	"debug/elf"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

// rootOwnedSafe reports why path is not safe to auto-trust as a system library, or nil when it is:
// the file AND its parent directory must be owned by root (uid 0) and not writable by group (unless
// the group is root) or other. A user-writable parent lets its owner replace even a root-owned file
// by unlink+create, so both levels are checked — the same rule as guard_trust.bpf.c's
// is_system_trusted. Symlinks are resolved by the caller (resolveSoname), so path is a real file.
func rootOwnedSafe(path string) error {
	if err := rootOwnedInode(path); err != nil {
		return err
	}
	if err := rootOwnedInode(filepath.Dir(path)); err != nil {
		return fmt.Errorf("parent %s", err)
	}
	return nil
}

// SystemTrusted reports whether path passes rootOwnedSafe: a non-root user can't replace it.
func SystemTrusted(path string) bool {
	return rootOwnedSafe(path) == nil
}

func rootOwnedInode(path string) error {
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return fmt.Errorf("cannot stat %s: %w", path, err)
	}
	if st.Uid != 0 {
		return fmt.Errorf("%s is not root-owned (uid %d)", path, st.Uid)
	}
	if st.Mode&syscall.S_IWOTH != 0 {
		return fmt.Errorf("%s is world-writable", path)
	}
	if st.Mode&syscall.S_IWGRP != 0 && st.Gid != 0 {
		return fmt.Errorf("%s is group-writable by non-root gid %d", path, st.Gid)
	}
	return nil
}

// ResolveLibraryClosure returns the absolute paths of every shared object a dynamically linked
// binary loads at startup, by STATIC analysis only: the program interpreter (ld-linux*.so) plus the
// transitive DT_NEEDED closure. It never runs the target (unlike ldd), so it's safe on any
// whitelisted binary.
//
// Deliberately NOT complete: dlopen'd libraries (NSS, gconv, GL drivers, plugins) are invisible
// here and must come from allow_lib. The result is the automatic base the daemon trusts for every
// whitelisted binary (libc, the loader) so operators needn't list them.
//
// Empty closure, no error, for: a statically linked binary (no INTERP/dynamic section); a script
// (#!), which runs as its interpreter, a separate binary resolved on its own; an ELF for an
// architecture this machine can't run natively (Steam ships pressure-vessel-arm64 for FEX), whose
// interpreter isn't on the host.
func ResolveLibraryClosure(binaryPath string) ([]string, error) {
	trusted, _, err := LibraryClosure(binaryPath)
	return trusted, err
}

// LibraryClosure is ResolveLibraryClosure plus the closure members it refused to auto-trust, with
// why. The kernel may still accept those (guarded tree, reserved name); the caller decides.
func LibraryClosure(binaryPath string) (trusted []string, rejected map[string]error, err error) {
	if isScript(binaryPath) {
		return nil, nil, nil
	}
	f, err := elf.Open(binaryPath)
	if err != nil {
		return nil, nil, fmt.Errorf("opening ELF %s: %w", binaryPath, err)
	}
	defer f.Close()
	if !runsNatively(f.Machine) {
		return nil, nil, nil
	}

	out := make(map[string]struct{})

	if interp, ierr := elfInterp(f); ierr == nil && interp != "" {
		if abs, aerr := filepath.Abs(interp); aerr == nil {
			out[abs] = struct{}{}
		}
	}

	visited := make(map[string]bool)
	walkNeeded(binaryPath, f, out, visited)

	// Trust only root-owned, non-user-writable libraries. The resolved set is fed to the daemon-wide
	// TRUSTED_LIB map, an unconditional pass in trust_mmap, so a library reachable through the
	// binary's own RPATH/$ORIGIN under a user-writable directory (many catalog apps live in $HOME and
	// ship RUNPATH=$ORIGIN) would otherwise let same-user malware plant an LD_PRELOAD payload that
	// gets trusted into every whitelisted process. This mirrors the kernel's is_system_trusted
	// auto-trust rule; a legitimate non-root library must be listed via allow_lib after review.
	trusted = make([]string, 0, len(out))
	for p := range out {
		if rerr := rootOwnedSafe(p); rerr != nil {
			if rejected == nil {
				rejected = make(map[string]error)
			}
			rejected[p] = rerr
			continue
		}
		trusted = append(trusted, p)
	}
	return trusted, rejected, nil
}

// walkNeeded resolves this object's DT_NEEDED entries to absolute paths and recurses, guarding
// cycles with visited (keyed by resolved path).
func walkNeeded(objPath string, f *elf.File, out map[string]struct{}, visited map[string]bool) {
	needed, err := f.ImportedLibraries()
	if err != nil {
		// No dynamic section (static binary) — not an error, just no deps.
		return
	}
	search := searchDirs(objPath, f)
	for _, soname := range needed {
		resolved := resolveSoname(soname, search)
		if resolved == "" {
			// A soname we can't place statically (odd layout, or itself dlopen-resolved): leave it
			// out rather than guess; the operator can add it via allow_lib if a load is later
			// denied.
			continue
		}
		if visited[resolved] {
			continue
		}
		visited[resolved] = true
		out[resolved] = struct{}{}

		child, cerr := elf.Open(resolved)
		if cerr != nil {
			continue
		}
		walkNeeded(resolved, child, out, visited)
		child.Close()
	}
}

// elfInterp returns the program interpreter path from the PT_INTERP segment,
// or "" when there is none (a static binary).
func elfInterp(f *elf.File) (string, error) {
	for _, p := range f.Progs {
		if p.Type != elf.PT_INTERP {
			continue
		}
		buf := make([]byte, p.Filesz)
		if _, err := p.ReadAt(buf, 0); err != nil {
			return "", err
		}
		if i := bytes.IndexByte(buf, 0); i >= 0 {
			buf = buf[:i]
		}
		return string(buf), nil
	}
	return "", nil
}

// searchDirs is objPath's ordered library search path: its DT_RPATH/DT_RUNPATH ($ORIGIN expanded),
// then the /etc/ld.so.conf directories (the set ldconfig uses), then conventional defaults.
// RPATH/RUNPATH first approximates the loader's precedence well enough for the static closure.
func searchDirs(objPath string, f *elf.File) []string {
	var dirs []string
	origin := filepath.Dir(objPath)
	for _, tag := range []elf.DynTag{elf.DT_RUNPATH, elf.DT_RPATH} {
		vals, err := f.DynString(tag)
		if err != nil {
			continue
		}
		for _, v := range vals {
			for _, d := range strings.Split(v, ":") {
				d = strings.ReplaceAll(d, "$ORIGIN", origin)
				d = strings.ReplaceAll(d, "${ORIGIN}", origin)
				if d != "" {
					dirs = append(dirs, d)
				}
			}
		}
	}
	dirs = append(dirs, ldSoConfDirs()...)
	dirs = append(dirs, defaultLibDirs()...)
	return dirs
}

// resolveSoname places a DT_NEEDED soname on disk: an absolute one is used as-is, else each search
// dir is tried in order (a soname is itself the filename, e.g. libc.so.6).
func resolveSoname(soname string, search []string) string {
	if filepath.IsAbs(soname) {
		if fileExists(soname) {
			return soname
		}
		return ""
	}
	for _, d := range search {
		cand := filepath.Join(d, soname)
		if fileExists(cand) {
			if abs, err := filepath.Abs(cand); err == nil {
				return abs
			}
			return cand
		}
	}
	return ""
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.Mode().IsRegular()
}

func defaultLibDirs() []string {
	// Conventional GNU/Linux defaults, covering the common multiarch triplet.
	return []string{
		"/lib", "/usr/lib",
		"/lib64", "/usr/lib64",
		"/lib/x86_64-linux-gnu", "/usr/lib/x86_64-linux-gnu",
		"/usr/local/lib",
	}
}

var (
	ldSoConfOnce sync.Once
	ldSoConfList []string
)

// ldSoConfDirs returns the library dirs from /etc/ld.so.conf (including `include` globs), parsed in
// pure Go once, no external process. On any error it yields nil and resolution uses the defaults.
func ldSoConfDirs() []string {
	ldSoConfOnce.Do(func() {
		p := &ldConfParser{seen: make(map[string]bool), visited: make(map[string]bool)}
		p.parseFile("/etc/ld.so.conf", 0)
		ldSoConfList = p.dirs
	})
	return ldSoConfList
}

// ldConfParser accumulates the ld.so.conf directory set across include globs,
// deduplicating directories and guarding against include cycles.
type ldConfParser struct {
	dirs    []string
	seen    map[string]bool
	visited map[string]bool
}

func (p *ldConfParser) parseFile(path string, depth int) {
	if depth > 16 || p.visited[path] {
		return // guard against include cycles / runaway recursion
	}
	p.visited[path] = true
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		p.parseLine(strings.TrimSpace(sc.Text()), path, depth)
	}
	_ = sc.Err() // best-effort: a partial dir list still resolves most sonames
}

func (p *ldConfParser) parseLine(line, path string, depth int) {
	if line == "" || strings.HasPrefix(line, "#") {
		return
	}
	if glob, ok := ldConfInclude(line, path); ok {
		matches, _ := filepath.Glob(glob)
		for _, m := range matches {
			p.parseFile(m, depth+1)
		}
		return
	}
	if !p.seen[line] { // a bare directory line
		p.seen[line] = true
		p.dirs = append(p.dirs, line)
	}
}

// ldConfInclude returns the (glob-expanded, absolute) target of an `include`
// directive on line, or ("", false) when line is not an include.
func ldConfInclude(line, fromPath string) (string, bool) {
	rest, ok := strings.CutPrefix(line, "include")
	if !ok || rest == "" || (rest[0] != ' ' && rest[0] != '\t') {
		return "", false
	}
	glob := strings.TrimSpace(rest)
	if !filepath.IsAbs(glob) {
		glob = filepath.Join(filepath.Dir(fromPath), glob)
	}
	return glob, true
}

// LdSoPreloadPaths returns the paths in /etc/ld.so.preload, which the loader force-loads into EVERY
// dynamically linked process: they must be trusted or nothing runs, so the daemon adds them to the
// trusted set unconditionally. nil if the file is missing.
func LdSoPreloadPaths() ([]string, error) {
	data, err := os.ReadFile("/etc/ld.so.preload")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var paths []string
	for _, field := range strings.Fields(string(data)) {
		if strings.HasPrefix(field, "#") {
			break
		}
		if fileExists(field) {
			paths = append(paths, field)
		}
	}
	return paths, nil
}

// isScript reports whether path starts with a "#!" interpreter line.
func isScript(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	magic := make([]byte, 2)
	n, _ := f.Read(magic)
	return n == 2 && bytes.Equal(magic, []byte("#!"))
}

// runsNatively reports whether an ELF for machine runs here without emulation: the native
// architecture plus its 32-bit compat one (i386 on amd64; Steam's client is 32-bit). An unknown
// architecture is assumed native, so the closure is still resolved rather than silently skipped.
func runsNatively(machine elf.Machine) bool {
	switch runtime.GOARCH {
	case "amd64":
		return machine == elf.EM_X86_64 || machine == elf.EM_386
	case "arm64":
		return machine == elf.EM_AARCH64 || machine == elf.EM_ARM
	default:
		return true
	}
}
