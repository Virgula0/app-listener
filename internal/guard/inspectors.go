package guard

import (
	"errors"
	"fmt"
	"slices"

	cilium "github.com/cilium/ebpf"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// maxInspectors mirrors guard_inspectors' max_entries in guard.bpf.c.
const maxInspectors = 64

// Inspector is one admitted [inspectors] line: its path and the inode it resolved to.
type Inspector struct {
	Path string
	Key  GuardInodeKey
}

// SetInspectors replaces the daemon's process inspectors (guard_inspectors): exe inodes that may
// read any tainted process's /proc metadata, never its memory. Keys must come from
// TrustGuard.AdmitInspector; listed is every configured path. A previous key whose inode is
// superseded but alive is kept while its path is still listed, so an inspector exec'd before a
// package upgrade works until it exits (exe_refused refuses every later exec of it).
func SetInspectors(admitted []Inspector, listed []string) error {
	return sharedEngine.setInspectors(admitted, listed)
}

func (e *engine) setInspectors(admitted []Inspector, listed []string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	want := make(map[GuardInodeKey]string, len(admitted))
	for _, in := range admitted {
		want[in.Key] = in.Path
	}
	for k, path := range e.inspectors {
		if _, ok := want[k]; ok || !slices.Contains(listed, path) {
			continue
		}
		if mark, ok := supersededMark(k); ok && mark.Freed == 0 {
			want[k] = path
		}
	}
	if len(want) > maxInspectors {
		return fmt.Errorf("%d inspector inodes exceed the limit of %d", len(want), maxInspectors)
	}
	if e.started {
		if err := e.syncInspectorsLocked(want); err != nil {
			return err
		}
	}
	e.inspectors = want
	return nil
}

// syncInspectorsLocked writes want over guard_inspectors: puts first, then deletes.
func (e *engine) syncInspectorsLocked(want map[GuardInodeKey]string) error {
	m := e.objs.GuardInspectors
	for k := range want {
		if err := m.Put(k, uint8(1)); err != nil {
			return fmt.Errorf("granting inspector inode %d: %w", k.Ino, err)
		}
	}
	var stale []GuardInodeKey
	var k GuardInodeKey
	var v uint8
	it := m.Iterate()
	for it.Next(&k, &v) {
		if _, ok := want[k]; !ok {
			stale = append(stale, k)
		}
	}
	if err := it.Err(); err != nil {
		return fmt.Errorf("reading inspectors: %w", err)
	}
	for _, k := range stale {
		if err := m.Delete(k); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return fmt.Errorf("revoking inspector inode %d: %w", k.Ino, err)
		}
	}
	return nil
}

// isInspectorLocked reports whether exe is a current inspector key.
func (e *engine) isInspectorLocked(exe GuardInodeKey) bool {
	_, ok := e.inspectors[exe]
	return ok
}

// AdmitInspector resolves path to an inspector key: a system file at a name only root could have
// placed, on a superblock root vouches for (systemFile), and TRUSTED_BINARY before it is returned, so
// trust_mmap refuses it untrusted code (LD_PRELOAD) from its first exec on.
func (t *TrustGuard) AdmitInspector(path string) (GuardInodeKey, error) {
	f, err := ebpf.OpenSystemPlaced(path)
	if err != nil {
		return GuardInodeKey{}, err
	}
	defer f.Close()
	dev, ino, err := ebpf.StatFile(f)
	if err != nil {
		return GuardInodeKey{}, err
	}
	k := GuardInodeKey{Dev: dev, Ino: ino}
	if !t.systemFile(path, f, k) {
		return GuardInodeKey{}, errors.New("not on a filesystem root vouches for (nosuid or user mount)")
	}
	liftSuperseded(k)
	tk := GuardTrustInodeKey{Dev: dev, Ino: ino}
	var flags uint8
	if t.objs.GuardTrustedFiles.Lookup(tk, &flags) == nil && flags&trustedBinary != 0 {
		return k, nil
	}
	if err := t.objs.GuardTrustedFiles.Put(tk, flags|trustedBinary); err != nil {
		return GuardInodeKey{}, fmt.Errorf("trusting the binary: %w", err)
	}
	return k, nil
}
