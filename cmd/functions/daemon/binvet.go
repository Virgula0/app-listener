package daemon

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/binledger"
	"github.com/Virgula0/app-listener/internal/constants"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
)

// binaryVetter decides which whitelisted binaries a start or reload may admit (issue #80). Nothing
// guards a whitelisted path while the daemon is not enforcing (stopped, booting, install), so what
// sits there now is admitted only when:
//   - only root could have placed it (ebpf.SystemPlacedInode);
//   - its content hash is the one the ledger records for its whitelist line;
//   - an updater of its resource created it while this daemon ran (guard_bin_origin);
//   - or the ledger was never initialized (first start after an upgrade), or the install wizard
//     just added the line (binledger.Grant): recorded on that start only.
//
// Anything else is refused and listed for `app-listener trust-binaries`. A refused binary gains no
// updater or reserved-name writer rights either (SetBinaryResolver): with them it could write a new
// inode at a confirmed path and have it admitted by provenance.
type binaryVetter struct {
	ledger *binledger.Ledger
	out    io.Writer

	mu        sync.Mutex
	bootstrap bool
	trust     *guard.TrustGuard
	// decided/approved: this generation's verdicts, by resolved path and by every path naming the
	// approved inode (the config's Path and the resolved one).
	decided  map[string]vetDecision
	approved map[string]guard.GuardInodeKey
	// linksByTarget: whitelist symlink lines per target, so an in-place update of a target is
	// recorded under the line the ledger keys on.
	linksByTarget map[string][]string
}

type vetDecision struct {
	line string
	key  guard.GuardInodeKey
	hash [32]byte
	ok   bool
}

func newBinaryVetter(l *binledger.Ledger, out io.Writer) (*binaryVetter, error) {
	initialized, err := l.Initialized()
	if err != nil {
		return nil, fmt.Errorf("reading the binary ledger: %w", err)
	}
	if !initialized {
		log.Warnf("daemon: the binary ledger %s is new: recording every whitelisted binary found now as "+
			"confirmed (one time). From the next start, a binary that changed while the daemon was not "+
			"running is refused until `app-listener trust-binaries` confirms it", binledger.DefaultPath)
	}
	return &binaryVetter{ledger: l, out: out, bootstrap: !initialized}, nil
}

// begin starts a generation for cfg: verdicts of the previous one no longer apply.
func (v *binaryVetter) begin(cfg *daemonconfig.Config) {
	links := map[string][]string{}
	for i := range cfg.Resources {
		for _, list := range [][]daemonconfig.BinaryRule{cfg.Resources[i].Binaries, cfg.Resources[i].PendingBinaries} {
			for _, b := range list {
				if b.Link != "" {
					links[b.Path] = append(links[b.Path], b.Link)
				}
			}
		}
	}
	v.mu.Lock()
	v.decided = map[string]vetDecision{}
	v.approved = map[string]guard.GuardInodeKey{}
	v.linksByTarget = links
	v.mu.Unlock()
}

func (v *binaryVetter) setTrust(t *guard.TrustGuard) {
	v.mu.Lock()
	v.trust = t
	v.mu.Unlock()
}

// endBootstrap closes the one-time recording windows once startup admitted the whitelist: the
// first-start bootstrap and the install's grants.
func (v *binaryVetter) endBootstrap() {
	v.mu.Lock()
	was := v.bootstrap
	v.bootstrap = false
	v.mu.Unlock()
	if err := v.ledger.DropGrants(); err != nil {
		log.Errorf("daemon: dropping unused install grants from the binary ledger: %v", err)
	}
	if !was {
		return
	}
	if err := v.ledger.MarkInitialized(); err != nil {
		log.Errorf("daemon: could not mark the binary ledger initialized (%v) — the next start records "+
			"the whitelist again instead of checking it", err)
	}
}

// resolve is the trust guard's binary resolver: only an approved path resolves, to the inode
// approved.
func (v *binaryVetter) resolve(path string) (dev, ino uint64, err error) {
	v.mu.Lock()
	k, ok := v.approved[path]
	v.mu.Unlock()
	if !ok {
		return 0, 0, guard.ErrUnconfirmed
	}
	return k.Dev, k.Ino, nil
}

// judge decides rule, resolved to the inode f holds (key, content hash), for res.
func (v *binaryVetter) judge(res *daemonconfig.Resource, rule daemonconfig.BinaryRule, resolved string,
	f *os.File, key guard.GuardInodeKey, hash [32]byte) bool {
	line := binledger.Line(rule)
	if err := multicallAdmissible(rule, f, key); err != nil {
		// Dropped, never fatal: a line whitelisted before a distro switched to a multicall must not
		// keep the daemon from protecting everything else.
		log.Errorf("daemon: CRITICAL: %s whitelist line %s dropped — it stays denied: %v",
			logging.SanitizeText(res.Path), logging.SanitizeText(line), err)
		return false
	}
	v.mu.Lock()
	prev, seen := v.decided[resolved]
	v.mu.Unlock()
	if seen && prev.line == line && prev.key == key && prev.hash == hash {
		return prev.ok
	}

	ok := v.decide(res, line, resolved, f, key, hash)

	v.mu.Lock()
	defer v.mu.Unlock()
	if seen && prev.ok && (prev.key != key || prev.hash != hash) {
		// One path, two inodes or contents within a generation: swapped between two opens. Each
		// guard keeps the inode it judged; neither earns write rights.
		ok = false
		delete(v.approved, rule.Path)
		delete(v.approved, resolved)
	}
	v.decided[resolved] = vetDecision{line: line, key: key, hash: hash, ok: ok}
	if ok {
		v.approved[rule.Path] = key
		v.approved[resolved] = key
	}
	return ok
}

func (v *binaryVetter) decide(res *daemonconfig.Resource, line, resolved string, f *os.File,
	key guard.GuardInodeKey, hash [32]byte) bool {
	if ebpf.SystemPlacedInode(line, key.Dev, key.Ino) {
		return true
	}
	recorded, found, err := v.ledger.Lookup(line)
	if err != nil {
		log.Errorf("daemon: binary ledger lookup for %s failed (%v) — refusing it", logging.SanitizeText(line), err)
		return false
	}
	if found && recorded == hash {
		return true
	}
	v.mu.Lock()
	bootstrap, trust := v.bootstrap, v.trust
	v.mu.Unlock()
	switch {
	case bootstrap && !found:
		v.record(line, hash, binledger.SourceBootstrap)
		return true
	case !found && v.takeGrant(line):
		v.record(line, hash, binledger.SourceInstall)
		return true
	case trust != nil && trust.CreatedByUpdater(f, key, ownerPaths(res)):
		v.record(line, hash, binledger.SourceUpdater)
		log.Infof("daemon: admitted %s: created by its resource's updater while the daemon ran", logging.SanitizeText(line))
		return true
	}
	v.refuse(res, line, resolved, hash, found)
	return false
}

// multicallAdmissible refuses a line reaching a multicall whose applets can't be told apart, or a
// uutils applet named through a symlink of another name: its process is keyed by the name exec'd.
func multicallAdmissible(rule daemonconfig.BinaryRule, f *os.File, key guard.GuardInodeKey) error {
	exe, err := guard.ExeKey(rule.Path, f, key)
	if err != nil {
		return err
	}
	if exe != key && rule.Link != "" && filepath.Base(rule.Link) != filepath.Base(rule.Path) {
		return fmt.Errorf("%s names multicall applet %s: whitelist the applet's own link",
			rule.Link, filepath.Base(rule.Path))
	}
	return nil
}

// takeGrant: the install wizard added line since the last start (binledger.Grant). Only during
// startup; endBootstrap drops what is left.
func (v *binaryVetter) takeGrant(line string) bool {
	ok, err := v.ledger.TakeGrant(line)
	if err != nil {
		log.Errorf("daemon: reading install grant for %s: %v", logging.SanitizeText(line), err)
	}
	return ok
}

func (v *binaryVetter) record(line string, hash [32]byte, src binledger.Source) {
	if err := v.ledger.Record(line, hash, src); err != nil {
		log.Errorf("daemon: recording %s in the binary ledger: %v — it will be refused at the next start "+
			"until confirmed", logging.SanitizeText(line), err)
	}
}

func (v *binaryVetter) refuse(res *daemonconfig.Resource, line, resolved string, hash [32]byte, changed bool) {
	reason := "new"
	if changed {
		reason = "changed"
	}
	if err := v.ledger.AddPending(&binledger.Pending{Line: line, Resolved: resolved, SHA256: hash,
		Resources: res.Path, SeenAt: time.Now()}); err != nil {
		log.Errorf("daemon: recording unconfirmed binary %s: %v", logging.SanitizeText(line), err)
	}
	fmt.Fprintf(v.out, "%sDAEMON binary-unconfirmed  resource=%s  path=%s  reason=%s  sha256=%x  — refused; "+
		"confirm it with `sudo app-listener trust-binaries` if you installed or updated it\n", syslogWarning,
		logging.SanitizeText(res.Path), logging.SanitizeText(line), reason, hash)
}

// recordLive records a replacement the trust guard approved live (updater-created or a system file)
// under every line naming path.
func (v *binaryVetter) recordLive(path string, f *os.File, key guard.GuardInodeKey) {
	entry, err := ebpf.ComputeBinaryEntryFile(f, path)
	if err != nil {
		log.Warnf("daemon: hashing re-admitted %s for the binary ledger: %v — it will be refused at the next "+
			"start until confirmed", logging.SanitizeText(path), err)
		return
	}
	v.mu.Lock()
	lines := append([]string{path}, v.linksByTarget[path]...)
	v.approved[path] = key
	v.mu.Unlock()
	for _, l := range lines {
		v.record(l, entry.Hash, binledger.SourceUpdater)
	}
}

// ownerPaths are the lines res whitelists: the trust guard's owner bits of the previous plan name
// the resources a new line's creator must update.
func ownerPaths(res *daemonconfig.Resource) []string {
	var out []string
	for _, list := range [][]daemonconfig.BinaryRule{res.Binaries, res.PendingBinaries} {
		for _, b := range list {
			out = append(out, b.Path)
			if b.Link != "" {
				out = append(out, b.Link)
			}
		}
	}
	return out
}

// openBinaryVetter opens the binary ledger and starts cfg's generation. The ledger is always
// binledger.DefaultPath, creating /etc/app-listener when a config elsewhere is used: it never lives
// beside a config whose directory a user may write.
func openBinaryVetter(cfg *daemonconfig.Config) (*binaryVetter, error) {
	if err := os.MkdirAll(filepath.Dir(binledger.DefaultPath), 0o755); err != nil {
		return nil, fmt.Errorf("%w: creating the binary ledger's directory: %w", constants.ErrCriticalStartup, err)
	}
	l, err := binledger.Open(binledger.DefaultPath)
	if err != nil {
		return nil, fmt.Errorf("%w: opening the binary ledger: %w", constants.ErrCriticalStartup, err)
	}
	v, err := newBinaryVetter(l, os.Stderr)
	if err != nil {
		l.Close()
		return nil, fmt.Errorf("%w: %w", constants.ErrCriticalStartup, err)
	}
	v.begin(cfg)
	return v, nil
}

func (v *binaryVetter) close() {
	if err := v.ledger.Close(); err != nil {
		log.Warnf("daemon: closing the binary ledger: %v", err)
	}
}
