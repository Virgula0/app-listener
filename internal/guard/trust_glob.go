package guard

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Name kinds and limits — must match guard_trust.bpf.c.
const (
	globExact    uint8 = 1
	globPrefix   uint8 = 2
	globSuffix   uint8 = 3
	globNameMax        = 32
	globAffixMax       = 8
	globMaxBits        = 64
)

// GlobReservations is protection #3 (guard_trust.bpf.c): names a catalog whitelist glob fixes,
// reserved below the glob's fixed root so that only the owning entry's binaries may bind them.
// Bit i of every mask stands for Patterns[i].
type GlobReservations struct {
	// Patterns are last glob components: "wineserver", "pv-*" or "*-capsule-capture-libs".
	Patterns []string
	// Roots maps a directory to the bits reserved at any depth below it.
	Roots map[string]uint64
	// Children reserves a name directly inside an existing directory.
	Children []GlobChild
	// Writers maps an executable to the bits it may bind.
	Writers map[string]uint64
}

// GlobChild reserves Name directly inside Parent for Bits.
type GlobChild struct {
	Parent string
	Name   string
	Bits   uint64
}

// ParseGlobName splits a last glob component into its BPF kind and fixed text. Only an exact name,
// "prefix*" and "*suffix" are expressible; ok is false for anything else.
func ParseGlobName(p string) (kind uint8, text string, ok bool) {
	switch {
	case !strings.ContainsAny(p, "*?[\\"):
		kind, text = globExact, p
	case strings.HasSuffix(p, "*") && !strings.ContainsAny(p[:len(p)-1], "*?[\\"):
		kind, text = globPrefix, p[:len(p)-1]
	case strings.HasPrefix(p, "*") && !strings.ContainsAny(p[1:], "*?[\\"):
		kind, text = globSuffix, p[1:]
	default:
		return 0, "", false
	}
	if text == "" || len(text) >= globNameMax || strings.Contains(text, "/") {
		return 0, "", false
	}
	return kind, text, true
}

func globText(s string) [globNameMax]int8 {
	var out [globNameMax]int8
	for i := 0; i < len(s) && i < globNameMax; i++ {
		out[i] = int8(s[i]) //nolint:gosec // byte reinterpreted as a C char
	}
	return out
}

// compileGlobNames builds the guard_glob_names entries (bit i = patterns[i]) and the distinct
// prefix/suffix shapes for guard_glob_affixes.
func compileGlobNames(patterns []string) (map[GuardTrustGlobNameKey]uint64, []GuardTrustGlobAffix, error) {
	if len(patterns) > globMaxBits {
		return nil, nil, fmt.Errorf("%d reserved glob names, at most %d fit", len(patterns), globMaxBits)
	}
	names := make(map[GuardTrustGlobNameKey]uint64)
	var affixes []GuardTrustGlobAffix
	for i, p := range patterns {
		kind, text, ok := ParseGlobName(p)
		if !ok {
			return nil, nil, fmt.Errorf("reserved glob name %q is not expressible (exact, prefix* or *suffix, < %d bytes)",
				p, globNameMax)
		}
		n := uint8(len(text)) //nolint:gosec // ParseGlobName bounds it below globNameMax
		names[GuardTrustGlobNameKey{Kind: kind, Len: n, S: globText(text)}] |= uint64(1) << i
		if a := (GuardTrustGlobAffix{Kind: kind, Len: n}); kind != globExact && !slices.Contains(affixes, a) {
			affixes = append(affixes, a)
		}
	}
	if len(affixes) > globAffixMax {
		return nil, nil, fmt.Errorf("%d prefix/suffix shapes, at most %d fit", len(affixes), globAffixMax)
	}
	return names, affixes, nil
}

// resolveBits stats every path and ORs its bits under the inode key; unresolvable paths are skipped.
func resolveBits(paths map[string]uint64) map[GuardTrustInodeKey]uint64 {
	out := make(map[GuardTrustInodeKey]uint64)
	for p, bits := range paths {
		if k, ok := statKey(p); ok {
			out[k] |= bits
		}
	}
	return out
}

func resolveChildren(children []GlobChild) (map[GuardTrustGlobChildKey]uint64, error) {
	out := make(map[GuardTrustGlobChildKey]uint64)
	for _, c := range children {
		if kind, _, ok := ParseGlobName(c.Name); !ok || kind != globExact {
			return nil, fmt.Errorf("reserved component %q is not a plain name", c.Name)
		}
		if k, ok := statKey(c.Parent); ok {
			out[GuardTrustGlobChildKey{Parent: k, S: globText(c.Name)}] |= c.Bits
		}
	}
	return out, nil
}

// SetGlobReservations (re)applies protection #3. Entries are written before stale ones are
// deleted, so a reload never passes through an empty (unenforced) policy. Unresolvable paths are
// skipped: a missing root or writer reserves nothing / grants nothing, and the next reload retries.
func (t *TrustGuard) SetGlobReservations(r GlobReservations) error {
	names, affixes, err := compileGlobNames(r.Patterns)
	if err != nil {
		return err
	}
	children, err := resolveChildren(r.Children)
	if err != nil {
		return err
	}
	roots, writers := resolveBits(r.Roots), resolveBits(r.Writers)

	// Writers and names first: a root that lands before its writers would deny the app itself.
	if err := syncMap(t.objs.GuardGlobWriters, writers); err != nil {
		return fmt.Errorf("glob writers: %w", err)
	}
	if err := syncMap(t.objs.GuardGlobNames, names); err != nil {
		return fmt.Errorf("glob names: %w", err)
	}
	for i := range uint32(globAffixMax) {
		var a GuardTrustGlobAffix
		if int(i) < len(affixes) {
			a = affixes[i]
		}
		if err := t.objs.GuardGlobAffixes.Put(i, a); err != nil {
			return fmt.Errorf("glob affix %d: %w", i, err)
		}
	}
	if err := syncMap(t.objs.GuardGlobRoots, roots); err != nil {
		return fmt.Errorf("glob roots: %w", err)
	}
	if err := syncMap(t.objs.GuardGlobChildren, children); err != nil {
		return fmt.Errorf("glob children: %w", err)
	}
	log.Infof("trust guard: %d reserved glob name(s) under %d root(s), %d writer(s)",
		len(r.Patterns), len(roots), len(writers))
	return nil
}

func statKey(path string) (GuardTrustInodeKey, bool) {
	dev, ino, err := ebpf.StatInode(path)
	if err != nil {
		log.Warnf("trust guard: skipping unresolvable %s: %v", path, err)
		return GuardTrustInodeKey{}, false
	}
	return GuardTrustInodeKey{Dev: dev, Ino: ino}, true
}

// syncMap makes m hold exactly want: puts first, then deletes keys want lacks.
func syncMap[K comparable](m *cilium.Map, want map[K]uint64) error {
	for k, v := range want {
		if err := m.Put(k, v); err != nil {
			return err
		}
	}
	var (
		k     K
		v     uint64
		stale []K
	)
	it := m.Iterate()
	for it.Next(&k, &v) {
		if _, ok := want[k]; !ok {
			stale = append(stale, k)
		}
	}
	if err := it.Err(); err != nil {
		return err
	}
	for i := range stale {
		if err := m.Delete(stale[i]); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	return nil
}
