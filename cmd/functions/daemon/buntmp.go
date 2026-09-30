package daemon

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/install"
)

type inodeID [2]uint64

// reserveGlobs applies r with every Bun tmp dir root pinned to a vetted inode. glob_lib_trusted
// trusts any .bun-* below the root for the Bun app, whoever wrote it, so a dir this daemon has not
// reserved before may hold a plant that predates its reservation: openBunTmp replaces it first. The
// root goes to SetGlobReservations as /proc/self/fd/N, so the inode reserved is the one vetted, not
// whatever a rename put at the path meanwhile. A dir that can't be vetted stays unreserved (its
// .bun-* libraries untrusted, the Bun app denied).
func (m *trustManager) reserveGlobs(r guard.GlobReservations, bun []install.User) error {
	pinned := r
	pinned.Roots = maps.Clone(r.Roots)
	vetted := make(map[inodeID]bool)
	var open []*bunTmp
	defer func() {
		for _, t := range open {
			t.f.Close()
		}
	}()
	for _, u := range bun {
		path := install.BunTmpDir(u.Home)
		bits := pinned.Roots[path]
		delete(pinned.Roots, path)
		t, err := openBunTmp(u, m.bunRoots)
		if err != nil {
			log.Errorf("trust guard: CRITICAL: cannot vet %s (%v) — its .bun-* libraries stay untrusted, so the "+
				"Bun app's native library is denied until a reload succeeds", path, err)
			continue
		}
		open = append(open, t)
		pinned.Roots[fmt.Sprintf("/proc/self/fd/%d", t.f.Fd())] |= bits
		vetted[t.id] = true
	}
	if err := m.tg.SetGlobReservations(pinned); err != nil {
		m.bunRoots = nil // the kernel may hold a partial set: vet every dir again next time
		return err
	}
	m.bunRoots = vetted
	for _, t := range open {
		if !t.renewed {
			continue
		}
		// Only now, reserved, may the user write it; until then it is root's and 0700.
		if err := t.f.Chown(int(t.user.UID), int(t.user.GID)); err != nil {
			log.Errorf("trust guard: handing %s to %s: %v — the Bun app cannot extract there", t.f.Name(),
				t.user.Name, err)
		}
	}
	return nil
}

type bunTmp struct {
	f       *os.File
	id      inodeID
	user    install.User
	renewed bool
}

// openBunTmp opens u's Bun tmp dir. One not in known (not reserved by the last apply) is renamed away
// and deleted, and an empty dir made in its place, owned by the daemon and 0700 so nothing can be
// planted before it is reserved. Replacing, not sweeping: renames inside a reserved root are
// allowed, so a same-user process could shuffle a plant past a recursive sweep. A lost extraction
// is re-extracted on the Bun app's next launch; a running one keeps its mapping.
func openBunTmp(u install.User, known map[inodeID]bool) (*bunTmp, error) {
	home, err := os.OpenRoot(u.Home)
	if err != nil {
		return nil, err
	}
	defer home.Close()
	parent, err := home.OpenRoot(filepath.Dir(install.BunTmpRelDir))
	if err != nil {
		return nil, err
	}
	defer parent.Close()
	name := filepath.Base(install.BunTmpRelDir)

	if f, oerr := openDirNoFollow(parent, name); oerr == nil {
		if id, ierr := fileID(f); ierr == nil && known[id] {
			return &bunTmp{f: f, id: id, user: u}, nil
		}
		f.Close()
	}
	f, err := renewDir(parent, name)
	if err != nil {
		return nil, err
	}
	id, err := fileID(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return &bunTmp{f: f, id: id, user: u, renewed: true}, nil
}

// renewDir replaces name in parent with an empty dir owned by the daemon, and opens it.
func renewDir(parent *os.Root, name string) (*os.File, error) {
	stale := fmt.Sprintf("%s.stale.%d", name, time.Now().UnixNano())
	if err := parent.Rename(name, stale); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("moving the unvetted dir away: %w", err)
	}
	if err := parent.RemoveAll(stale); err != nil {
		log.Warnf("trust guard: removing %s/%s: %v (untrusted: no longer a reserved root)", parent.Name(), stale, err)
	}
	if err := parent.Mkdir(name, 0o700); err != nil {
		return nil, err
	}
	f, err := openDirNoFollow(parent, name)
	if err != nil {
		return nil, err
	}
	if err := checkFreshDir(f); err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func openDirNoFollow(r *os.Root, name string) (*os.File, error) {
	return r.OpenFile(name, os.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW, 0)
}

// checkFreshDir: f is the empty dir just made, not one a same-user process swapped in (it can't
// create one owned by the daemon's euid with no group/other access).
func checkFreshDir(f *os.File) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || !fi.IsDir() || int(st.Uid) != os.Geteuid() || fi.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("%s is not the dir the daemon just created", f.Name())
	}
	if _, err := f.Readdirnames(1); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%s is not empty", f.Name())
	}
	return nil
}

func fileID(f *os.File) (inodeID, error) {
	fi, err := f.Stat()
	if err != nil {
		return inodeID{}, err
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return inodeID{}, fmt.Errorf("%s: no inode", f.Name())
	}
	return inodeID{st.Dev, st.Ino}, nil
}
