package maintainer

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	maxEntries = 500_000
	// maxRecord bounds a path record, the one file kept whose content the lib_dir need not hold.
	maxRecord = 4096
)

type opKind uint8

const (
	opDelete opKind = iota + 1
	opDir
	opLink
	opFile
)

// op is one change below a lib_dir, in walk order: a directory before its entries.
type op struct {
	kind opKind
	rel  string // below the lib_dir; "." is the lib_dir itself (an opaque opDir only)
	mode uint32
	// opaque: a directory the run replaced. Tree entries not in names are removed.
	opaque bool
	names  map[string]bool
	target string // opLink
	src    string // opFile: its copy in the stage
	size   int64
	sum    [32]byte
	data   []byte // opFile: a path record, committed from memory
}

type libPlan struct {
	dir string
	ops []op
}

// Plan is a run's validated changes, per lib_dir. Close it to remove the stage.
type Plan struct {
	st       *stage
	uid, gid uint32
	libs     []libPlan
}

// Close removes the run's stage.
func (p *Plan) Close() { p.st.close() }

// Changes counts the plan's operations.
func (p *Plan) Changes() int {
	n := 0
	for i := range p.libs {
		n += len(p.libs[i].ops)
	}
	return n
}

// Touches reports whether the plan writes or removes path, or removes or replaces a directory
// above it.
func (p *Plan) Touches(target string) bool {
	for i := range p.libs {
		lp := &p.libs[i]
		rel, err := filepath.Rel(lp.dir, target)
		if err != nil || rel == ".." || strings.HasPrefix(rel, "../") {
			continue
		}
		for j := range lp.ops {
			o := &lp.ops[j]
			// Anything but a kept directory replaces what was below it: a symlink there too.
			whole := o.kind != opDir || o.opaque
			if o.rel == rel || (whole && (o.rel == "." || strings.HasPrefix(rel, o.rel+"/"))) {
				return true
			}
		}
	}
	return false
}

func planFrom(spec *Spec, root string, st *stage) (*Plan, error) {
	p := &Plan{st: st, uid: spec.UID, gid: spec.GID}
	for _, lib := range spec.LibDirs {
		rel, err := filepath.Rel(root, lib)
		if err != nil {
			return nil, err
		}
		ops, err := scanLib(st.upper, rel)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", lib, err)
		}
		if len(ops) == 0 {
			continue
		}
		if err := validate(lib, ops); err != nil {
			return nil, fmt.Errorf("refusing the run's changes to %s: %w", lib, err)
		}
		p.libs = append(p.libs, libPlan{dir: lib, ops: ops})
	}
	// Links last: one may lead into another lib_dir, judged as the whole plan leaves it.
	v := newLinkView(spec.LibDirs, p.libs)
	for i := range p.libs {
		if err := v.checkLinks(&p.libs[i]); err != nil {
			return nil, fmt.Errorf("refusing the run's changes to %s: %w", p.libs[i].dir, err)
		}
	}
	return p, nil
}

// scanLib lists what the run changed below rel, the lib_dir's path in the overlay's upper layer.
func scanLib(upper, rel string) ([]op, error) {
	cur := upper
	replaced := false
	for c := range strings.SplitSeq(rel, "/") {
		cur = filepath.Join(cur, c)
		fi, err := os.Lstat(cur)
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		if isWhiteout(fi) || !fi.IsDir() {
			return nil, errors.New("the run removed or replaced the lib_dir itself")
		}
		// A replaced ancestor hides the whole tree: the lib_dir's upper copy is then all of it.
		replaced = replaced || opaque(cur)
	}
	w := walker{}
	if replaced {
		w.ops = append(w.ops, op{kind: opDir, rel: ".", opaque: true})
	}
	names, err := w.walk(cur, "", 0)
	if err != nil {
		return nil, err
	}
	if replaced {
		w.ops[0].names = names
	}
	return w.ops, nil
}

type walker struct {
	ops []op
	n   int
}

// walk appends dir's changes and returns the names it holds (whiteouts excluded).
func (w *walker) walk(dir, rel string, depth int) (map[string]bool, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("%s: deeper than %d", rel, maxDepth)
	}
	des, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make(map[string]bool, len(des))
	for _, de := range des {
		w.n++
		if w.n > maxEntries {
			return nil, fmt.Errorf("more than %d changed entries", maxEntries)
		}
		p, r := filepath.Join(dir, de.Name()), path.Join(rel, de.Name())
		fi, err := os.Lstat(p)
		if err != nil {
			return nil, err
		}
		if isWhiteout(fi) {
			w.ops = append(w.ops, op{kind: opDelete, rel: r})
			continue
		}
		names[de.Name()] = true
		switch mode := fi.Mode(); {
		case mode.IsDir():
			i := len(w.ops)
			w.ops = append(w.ops, op{kind: opDir, rel: r, mode: uint32(mode.Perm()), opaque: opaque(p)})
			sub, err := w.walk(p, r, depth+1)
			if err != nil {
				return nil, err
			}
			w.ops[i].names = sub
		case mode&os.ModeSymlink != 0:
			target, err := os.Readlink(p)
			if err != nil {
				return nil, err
			}
			w.ops = append(w.ops, op{kind: opLink, rel: r, target: target})
		case mode.IsRegular():
			// Perm drops setuid, setgid and sticky.
			w.ops = append(w.ops, op{kind: opFile, rel: r, mode: uint32(mode.Perm()), src: p, size: fi.Size()})
		default:
			return nil, fmt.Errorf("%s: a %v, never committed", r, mode.Type())
		}
	}
	return names, nil
}

// validate accepts each file only when its content is already in the lib_dir (a file, or a member
// of an archive there) or it is a path record. Symlinks are linkView's; directories and deletions
// add no content.
func validate(lib string, ops []op) error {
	src := newSources(lib)
	for i := range ops {
		o := &ops[i]
		if o.kind != opFile {
			continue
		}
		if err := src.vouch(o); err != nil {
			return fmt.Errorf("%s: %w", o.rel, err)
		}
	}
	return nil
}

// vouch hashes o's staged copy and accepts it.
func (s *sources) vouch(o *op) error {
	f, err := os.OpenFile(o.src, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	var head bytes.Buffer
	n, err := io.Copy(io.MultiWriter(h, &limitedWriter{w: &head, n: maxRecord + 1}), f)
	if err != nil {
		return err
	}
	o.size = n
	copy(o.sum[:], h.Sum(nil))
	if n <= maxRecord && pathRecord(head.Bytes()) {
		o.data = bytes.Clone(head.Bytes())
		o.mode &^= 0o111
		return nil
	}
	ok, err := s.has(n, o.sum)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("content (sha256 %x) is neither in the lib_dir nor a path record", o.sum)
	}
	return nil
}

// pathRecord: empty, or lines that are each empty or an absolute path of [A-Za-z0-9/._+@~:,-]
// (setup.sh's pin records). No quote, bracket, operator or space: not script, markup or a command
// with arguments in any language, and committed without exec bits.
func pathRecord(b []byte) bool {
	for line := range bytes.SplitSeq(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		if line[0] != '/' {
			return false
		}
		for _, c := range line {
			if !pathByte(c) {
				return false
			}
		}
	}
	return true
}

func pathByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || strings.IndexByte("/._+@~:,-", c) >= 0
}

// isWhiteout: overlayfs records a deletion as a 0:0 character device.
func isWhiteout(fi os.FileInfo) bool {
	st, ok := fi.Sys().(*syscall.Stat_t)
	return ok && fi.Mode()&os.ModeCharDevice != 0 && st.Rdev == 0
}

func opaque(dir string) bool {
	buf := make([]byte, 2)
	n, err := unix.Lgetxattr(dir, "trusted.overlay.opaque", buf)
	return err == nil && string(buf[:n]) == "y"
}
