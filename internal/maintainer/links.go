package maintainer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"syscall"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// maxLinkHops bounds the symlinks a target is followed through, as the kernel's MAXSYMLINKS.
const maxLinkHops = 40

type entryKind uint8

const (
	entryMissing entryKind = iota
	entryDir
	entryFile
	entryLink
)

// linkView is the lib_dirs as the plan leaves them: its ops over the live tree.
type linkView struct {
	libs []string
	ops  map[string]*op // by absolute path; a replaced lib_dir root under the lib_dir's own
	// placed caches ebpf.RootPlacedName by directory: a runtime's links share a few targets.
	placed map[string]error
}

func newLinkView(libDirs []string, libs []libPlan) *linkView {
	v := &linkView{ops: map[string]*op{}, placed: map[string]error{}}
	for _, l := range libDirs {
		v.libs = append(v.libs, path.Clean(l))
	}
	for i := range libs {
		for j := range libs[i].ops {
			o := &libs[i].ops[j]
			v.ops[path.Join(libs[i].dir, o.rel)] = o
		}
	}
	return v
}

// checkLinks accepts each symlink only when it leads, through every hop, into a lib_dir or to a name
// only root decides. Whatever writes the stage may be hostile (same-uid malware reaches the overlay
// through /proc/<pid>/cwd), and a symlink is a name the app's whitelisted binaries trust as their
// own: one out to the user's files would hand them content no lib_dir vouches for.
func (v *linkView) checkLinks(lp *libPlan) error {
	for i := range lp.ops {
		o := &lp.ops[i]
		if o.kind != opLink {
			continue
		}
		abs := path.Join(lp.dir, o.rel)
		next, err := linkTarget(abs, o.target)
		if err == nil {
			err = v.resolve(next, 1)
		}
		if err != nil {
			return fmt.Errorf("%s: the symlink to %q: %w", o.rel, o.target, err)
		}
	}
	return nil
}

// linkTarget is where target leads from the symlink at abs. Only a leading run of ".." is
// accepted: it climbs the link's own real directories (Apply opens them RESOLVE_NO_SYMLINKS), so
// the lexical path is the one the kernel walks; a ".." after a name could climb out of a
// symlinked directory instead.
func linkTarget(abs, target string) (string, error) {
	named := false
	for c := range strings.SplitSeq(target, "/") {
		switch c {
		case "", ".":
		case "..":
			if named {
				return "", errors.New(`a ".." after a name`)
			}
		default:
			named = true
		}
	}
	if target == "" {
		return "", errors.New("an empty target")
	}
	if path.IsAbs(target) {
		return path.Clean(target), nil
	}
	return path.Join(path.Dir(abs), target), nil
}

// resolve walks the clean absolute p as the kernel will once the plan is applied. Below a lib_dir
// every name is the app's (only its writers and the maintainer write there), so a missing one is
// fine; a symlink is followed and judged again. Anywhere else, root must decide what p names.
func (v *linkView) resolve(p string, hops int) error {
	if hops > maxLinkHops {
		return errors.New("too many levels of symbolic links")
	}
	lib := v.libOf(p)
	if lib == "" {
		return v.rootPlaced(p)
	}
	rest := strings.Split(strings.TrimPrefix(p[len(lib):], "/"), "/")
	cur := lib
	for i, c := range rest {
		if c == "" {
			continue
		}
		cur += "/" + c
		kind, target, err := v.kind(cur)
		if err != nil {
			return err
		}
		switch kind {
		case entryMissing, entryFile:
			return nil
		case entryLink:
			next, err := linkTarget(cur, target)
			if err != nil {
				return fmt.Errorf("through %s: %w", cur, err)
			}
			return v.resolve(path.Join(append([]string{next}, rest[i+1:]...)...), hops+1)
		}
	}
	return nil
}

func (v *linkView) libOf(p string) string {
	best := ""
	for _, l := range v.libs {
		if (p == l || strings.HasPrefix(p, l+"/")) && len(l) > len(best) {
			best = l
		}
	}
	return best
}

func (v *linkView) rootPlaced(p string) error {
	dir := path.Dir(p)
	err, ok := v.placed[dir]
	if !ok {
		err = ebpf.RootPlacedName(p)
		v.placed[dir] = err
	}
	if err != nil {
		return fmt.Errorf("it leads outside the lib_dirs to %s, a name root does not decide: %w", p, err)
	}
	return nil
}

// kind is what p, below a lib_dir, is once the plan is applied. resolve walks top-down, so p's
// parent is a directory in that view.
func (v *linkView) kind(p string) (entryKind, string, error) {
	if o := v.ops[p]; o != nil {
		switch o.kind {
		case opDelete:
			return entryMissing, "", nil
		case opDir:
			return entryDir, "", nil
		case opLink:
			return entryLink, o.target, nil
		default:
			return entryFile, "", nil
		}
	}
	// Every entry of a replaced directory has an op: one without is gone.
	if o := v.ops[path.Dir(p)]; o != nil && o.opaque {
		return entryMissing, "", nil
	}
	fi, err := os.Lstat(p)
	switch {
	case errors.Is(err, fs.ErrNotExist) || errors.Is(err, syscall.ENOTDIR):
		return entryMissing, "", nil
	case err != nil:
		return 0, "", err
	case fi.Mode()&os.ModeSymlink != 0:
		target, err := os.Readlink(p)
		return entryLink, target, err
	case fi.IsDir():
		return entryDir, "", nil
	}
	return entryFile, "", nil
}
