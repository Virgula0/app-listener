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
)

// ResolveLibraryClosure returns the absolute paths of every shared object a
// dynamically linked binary loads at startup by STATIC analysis only — the
// program interpreter (ld-linux*.so) plus the transitive DT_NEEDED closure —
// without executing the binary. It never runs the target (unlike ldd, which
// invokes the loader), so it is safe to call on any whitelisted binary.
//
// It is deliberately NOT complete: libraries a program brings in later via
// dlopen (NSS, gconv, GL drivers, plugins) are invisible to static analysis
// and must be supplied out of band (the daemon's allow_lib list, populated
// with help from the observe pass). The result is the automatic base the
// daemon trusts for every whitelisted binary so an operator never has to list
// libc/the loader by hand.
//
// A statically linked binary (no INTERP, no dynamic section) yields an empty
// closure and no error.
//
// Two kinds of whitelisted file have no closure of their own and yield an
// empty one with no error: a script (#!), which the kernel runs as its
// interpreter — a separate binary resolved on its own — and an ELF for an
// architecture this machine cannot execute natively (Steam ships
// pressure-vessel-arm64 for FEX), whose interpreter is not on the host.
func ResolveLibraryClosure(binaryPath string) ([]string, error) {
	if isScript(binaryPath) {
		return nil, nil
	}
	f, err := elf.Open(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("opening ELF %s: %w", binaryPath, err)
	}
	defer f.Close()
	if !runsNatively(f.Machine) {
		return nil, nil
	}

	out := make(map[string]struct{})

	if interp, ierr := elfInterp(f); ierr == nil && interp != "" {
		if abs, aerr := filepath.Abs(interp); aerr == nil {
			out[abs] = struct{}{}
		}
	}

	visited := make(map[string]bool)
	walkNeeded(binaryPath, f, out, visited)

	paths := make([]string, 0, len(out))
	for p := range out {
		paths = append(paths, p)
	}
	return paths, nil
}

// walkNeeded resolves this object's DT_NEEDED entries to absolute paths and
// recurses into each, guarding against cycles with visited (keyed by resolved
// path).
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
			// A soname we cannot place statically (unusual layout, or it is
			// itself dlopen-resolved). Leave it out rather than guess; the
			// operator can add it via allow_lib if a load is later denied.
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

// searchDirs is the ordered library search path for objPath: its DT_RPATH /
// DT_RUNPATH (with $ORIGIN expanded), then the directories configured in
// /etc/ld.so.conf (the same set ldconfig builds its cache from), then the
// conventional defaults. RPATH/RUNPATH first matches the loader's own
// precedence closely enough for the static closure.
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

// resolveSoname places a DT_NEEDED soname on disk: an absolute soname is used
// as-is; otherwise each search dir is tried in order (RPATH/RUNPATH, the
// ld.so.conf directories, then the defaults). A soname is itself the on-disk
// filename (libc.so.6), so a plain directory search resolves it.
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

// ldSoConfDirs returns the library directories configured in /etc/ld.so.conf
// (and everything its `include` globs pull in) — the same directory set
// ldconfig builds its cache from — parsed in pure Go, once, with no external
// process. On any error it yields nil and resolution falls back to the
// conventional defaults.
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

// LdSoPreloadPaths returns the absolute library paths listed in
// /etc/ld.so.preload, which the loader force-loads into EVERY dynamically
// linked process. They must be trusted or nothing runs; the daemon adds them
// to the trusted-library set unconditionally. A missing file yields nil.
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

// runsNatively reports whether an ELF for machine executes on this host
// without emulation: the native architecture plus its 32-bit compat one
// (i386 on amd64 — Steam's own client is a 32-bit binary). An architecture
// this list does not know is assumed native, so the closure is still
// resolved rather than silently skipped.
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
