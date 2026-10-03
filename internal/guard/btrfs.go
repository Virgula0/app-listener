package guard

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

var btrfsLayout struct {
	mu sync.Mutex
	ok bool
}

// ensureBtrfsLayout fills the btrfs_layout map the guard and trust objects share, once the
// kernel's BTF describes btrfs (the module's appears only once it is loaded). Until then inode_dev
// keys every btrfs inode INODE_DEV_UNKNOWN, which no row matches.
func ensureBtrfsLayout() bool {
	btrfsLayout.mu.Lock()
	defer btrfsLayout.mu.Unlock()
	if btrfsLayout.ok {
		return true
	}
	maps, err := supersedeMaps()
	if err != nil {
		return false
	}
	l, err := ebpf.ResolveBtrfsLayout()
	if err != nil {
		log.Debugf("guard: btrfs layout unresolved: %v", err)
		return false
	}
	v := GuardBtrfsLayout{Valid: l.Valid, InodeRoot: l.InodeRoot, InodeVfs: l.InodeVfs, RootAnonDev: l.RootAnonDev}
	if err := maps[GuardMapBtrfsLayout].Put(uint32(0), v); err != nil {
		log.Warnf("guard: storing the btrfs layout: %v", err)
		return false
	}
	btrfsLayout.ok = true
	return true
}

// checkBtrfsKeys refuses a guard with its root or a binary on btrfs while the layout is unknown:
// no key there would match, leaving the tree unguarded and a blacklisted binary unblocked.
func checkBtrfsKeys(root string, bins []BinaryEntry) error {
	if ensureBtrfsLayout() {
		return nil
	}
	paths := make([]string, 0, len(bins)+1)
	paths = append(paths, root)
	for _, b := range bins {
		paths = append(paths, b.Path)
	}
	if p := ebpf.FirstOnBtrfs(paths); p != "" {
		return fmt.Errorf("%s is on btrfs, but the kernel's BTF does not describe btrfs (struct btrfs_inode, "+
			"btrfs_root): its inodes cannot be identified, refusing to start", p)
	}
	return nil
}
