package fscrypt

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
)

// maxReportedOutsideLinks caps how many offending files the refusal names.
const maxReportedOutsideLinks = 5

// refuseOutsideHardlinks errors when a regular file under root has more links than names inside
// root: another name lives outside the tree. The guard confines a hit to names whose dentry chain
// reaches the watch root, so that outside name reads the plaintext once the key is provisioned. A
// guard denies link() of a guarded file, but not while none is attached (daemon stopped, update),
// and fscrypt only checks an encrypted link DESTINATION, so a locked vault's file links out fine.
// Called before provisioning, with the guard already attached: no new name can appear meanwhile.
func refuseOutsideHardlinks(root string) error {
	bad, err := outsideHardlinks(root)
	if err != nil {
		return fmt.Errorf("auditing hard links under %s: %w", root, err)
	}
	if len(bad) == 0 {
		return nil
	}
	total := len(bad)
	if total > maxReportedOutsideLinks {
		bad = bad[:maxReportedOutsideLinks]
	}
	return fmt.Errorf("refusing to unlock: %d file(s) have a hard link outside %s, which would expose "+
		"their plaintext unguarded: %s. Find each one's other names (find / -xdev -samefile <file>), "+
		"remove those outside the vault, then start the daemon again",
		total, root, strings.Join(bad, ", "))
}

type linkedFile struct {
	path         string
	names, nlink uint64
}

// outsideHardlinks lists root's regular files with fewer names inside root than links. A file that
// vanishes mid-walk is skipped; any other error fails the audit.
func outsideHardlinks(root string) ([]string, error) {
	type key struct{ dev, ino uint64 }
	files := map[key]*linkedFile{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) && p != root {
				return nil
			}
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		st, err := lstatOf(d)
		if err != nil || st == nil || st.Nlink <= 1 {
			return err
		}
		k := key{uint64(st.Dev), st.Ino} //nolint:unconvert // Dev's width differs per arch
		if f := files[k]; f != nil {
			f.names++
		} else {
			files[k] = &linkedFile{path: p, names: 1, nlink: uint64(st.Nlink)} //nolint:unconvert // per arch
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var bad []string
	for _, f := range files {
		if f.names < f.nlink {
			bad = append(bad, fmt.Sprintf("%s (%d of %d links inside)", f.path, f.names, f.nlink))
		}
	}
	slices.Sort(bad)
	return bad, nil
}

// lstatOf is d's lstat; nil, nil when the entry vanished.
func lstatOf(d fs.DirEntry) (*syscall.Stat_t, error) {
	info, err := d.Info()
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	st, _ := info.Sys().(*syscall.Stat_t)
	return st, nil
}
