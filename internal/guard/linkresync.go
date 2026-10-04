package guard

import (
	"errors"
	"fmt"
	"maps"
	"os"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
)

// A whitelist line that is a symlink (an updater's ~/.local/bin/claude -> versions/<v>) is admitted
// by the file it resolves to. An updater re-pointing it changes no config line and replaces no
// inode at a whitelisted path, so the re-sync follows the link itself: the new target passes the
// same replacement check as an in-place update, judged by the link's name and owners, and the old
// target is superseded so only processes exec'd before keep it.

// WithBinaryLinks records, per whitelist line that is a symlink, the target it was admitted by.
func WithBinaryLinks(links map[string]string) GuardOption {
	return func(g *Guard) {
		if g.links == nil {
			g.links = make(map[string]string, len(links))
		}
		maps.Copy(g.links, links)
	}
}

// resyncLinks follows every whitelisted link to its current target.
func (g *Guard) resyncLinks(exeEvents map[string][]ebpf.EventType) (int, error) {
	g.mu.Lock()
	links := maps.Clone(g.links)
	g.mu.Unlock()
	changed := 0
	for link, target := range links {
		ok, err := g.resyncLink(link, target, exeEvents)
		if err != nil {
			return changed, err
		}
		if ok {
			changed++
		}
	}
	return changed, nil
}

// resyncLink admits the file link now resolves to in place of target. One confined open: the inode
// judged, hashed and admitted is the one the link reached, whatever it points at afterwards.
func (g *Guard) resyncLink(link, target string, exeEvents map[string][]ebpf.EventType) (bool, error) {
	f, err := ebpf.OpenConfined(link)
	if err != nil {
		return false, nil // gone mid-update, or re-pointed out of its tree: the deployed key stays
	}
	defer f.Close()
	resolved, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil || resolved == target {
		return false, nil // same name: resyncOne judges a replacement there
	}
	dev, ino, err := ebpf.StatFile(f)
	if err != nil {
		return false, nil
	}
	key := GuardInodeKey{Dev: dev, Ino: ino}

	g.mu.Lock()
	old, deployed := g.deployed[target]
	if !deployed {
		old = g.retired[target]
	}
	g.mu.Unlock()
	if key == old {
		return false, nil
	}
	if !replacementAllowed(link, f, old, key) {
		if g.refused.firstTime(link, key) {
			log.Warnf("guard %s: %s now points at %s (inode %d), which was neither created by its updater nor a "+
				"system file at a root-placed name — not admitted; reload the daemon after verifying it", g.path,
				logging.SanitizeText(link), logging.SanitizeText(resolved), key.Ino)
		}
		return false, nil
	}
	entry, err := ebpf.ComputeBinaryEntryFile(f, resolved)
	if err != nil {
		return false, nil
	}
	liftSuperseded(key)
	if err := g.putBinaryKey(key, target, exeEvents); err != nil {
		return false, err
	}
	g.retargetLink(link, target, entry, key, old, deployed)
	supersedeUnnamed(old)
	log.Infof("guard %s: %s re-pointed to %s (inode %d): admitted, the previous target kept only for processes "+
		"started before", g.path, logging.SanitizeText(link), logging.SanitizeText(resolved), key.Ino)
	return true, nil
}

// retargetLink moves link's bookkeeping from target to entry (admitted as key). target's own
// bookkeeping goes unless another line still names it.
func (g *Guard) retargetLink(link, target string, entry BinaryEntry, key, old GuardInodeKey, deployed bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	resolved := entry.Path
	g.links[link] = resolved
	replaced, named := false, false
	for i := range g.binaries {
		if g.binaries[i].Path != target {
			continue
		}
		if !replaced {
			g.binaries[i], replaced = entry, true
			continue
		}
		named = true
	}
	if !replaced {
		g.binaries = append(g.binaries, entry)
	}
	for _, t := range g.links {
		named = named || t == target
	}
	if g.exeEvents == nil {
		g.exeEvents = make(map[string][]ebpf.EventType)
	}
	if types, ok := g.exeEvents[target]; ok {
		g.exeEvents[resolved] = types
	}
	g.canonicalPaths[resolved] = resolved
	g.deployed[resolved] = key
	delete(g.retired, resolved)
	g.keyPaths[key] = resolved
	if g.binaryVerifyStates != nil {
		g.binaryVerifyStates[resolved] = &binaryVerifyState{key: key, hash: entry.Hash}
	}
	if named {
		return
	}
	delete(g.deployed, target)
	delete(g.retired, target)
	if g.binaryVerifyStates != nil {
		delete(g.binaryVerifyStates, target)
	}
	if deployed {
		g.keyPaths[old] = link // carried across a reload while the link is whitelisted (carriesLine)
	}
}

// supersedeUnnamed marks k superseded once no guard names it by a whitelist line: as if its last
// link had gone, it keeps admitting processes exec'd before and refuses every later exec
// (exe_supersede.h). Its rows go when its inode does (PruneSuperseded).
func supersedeUnnamed(k GuardInodeKey) {
	if k == (GuardInodeKey{}) || sharedEngine.namesExe(k) {
		return
	}
	m, err := supersedeMaps()
	if err != nil {
		log.Warnf("guard: superseding inode %d: %v", k.Ino, err)
		return
	}
	supersede.mu.Lock()
	defer supersede.mu.Unlock()
	var now uint64
	if err = m[GuardMapExeSeq].Lookup(uint32(0), &now); err != nil {
		log.Warnf("guard: superseding inode %d: reading the exec counter: %v", k.Ino, err)
		return
	}
	// now+1: a process whose exec stamped now ran before the mark. An existing mark keeps its seq.
	err = m[GuardMapExeSuperseded].Update(k, GuardExeSupersede{Seq: now + 1}, cilium.UpdateNoExist)
	if err != nil && !errors.Is(err, cilium.ErrKeyExist) {
		log.Warnf("guard: superseding inode %d: %v — new processes still run the previous version", k.Ino, err)
	}
}

// namesExe reports whether any guard's whitelist line still resolves to exe.
func (e *engine) namesExe(exe GuardInodeKey) bool {
	e.mu.Lock()
	var guards []*Guard
	for _, g := range &e.slots {
		if g != nil {
			guards = append(guards, g)
		}
	}
	e.mu.Unlock()
	for _, g := range guards {
		g.mu.Lock()
		named := false
		for _, k := range g.deployed {
			named = named || k == exe
		}
		g.mu.Unlock()
		if named {
			return true
		}
	}
	return false
}
