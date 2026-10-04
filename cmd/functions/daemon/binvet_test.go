package daemon

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/binledger"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
)

func newTestVetter(t *testing.T) (*binaryVetter, *binledger.Ledger, *bytes.Buffer) {
	t.Helper()
	l, err := binledger.Open(filepath.Join(t.TempDir(), "binaries.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	var out bytes.Buffer
	v, err := newBinaryVetter(l, &out)
	if err != nil {
		t.Fatal(err)
	}
	return v, l, &out
}

func testBinary(t *testing.T, name, content string) (string, *os.File) {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	return p, f
}

func vetResource(rules ...daemonconfig.BinaryRule) (*daemonconfig.Config, *daemonconfig.Resource) {
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{{Path: "/home/u/.config/gh", Binaries: rules}}}
	return cfg, &cfg.Resources[0]
}

func TestVetterBootstrapRecordsThenChecks(t *testing.T) {
	v, l, out := newTestVetter(t)
	bin, f := testBinary(t, "gh", "v1")
	rule := daemonconfig.BinaryRule{Path: bin}
	cfg, res := vetResource(rule)
	key := guard.GuardInodeKey{Dev: 1, Ino: 10}

	v.begin(cfg)
	if !v.judge(res, rule, bin, f, key, [32]byte{1}) {
		t.Fatal("bootstrap refused a binary found at first start")
	}
	v.endBootstrap()
	if sum, ok, _ := l.Lookup(bin); !ok || sum != [32]byte{1} {
		t.Fatalf("bootstrap did not record the binary: %x %v", sum, ok)
	}

	// Next start: same content admitted, changed content refused and listed.
	v.begin(cfg)
	if !v.judge(res, rule, bin, f, key, [32]byte{1}) {
		t.Fatal("recorded content refused")
	}
	v.begin(cfg)
	if v.judge(res, rule, bin, f, key, [32]byte{2}) {
		t.Fatal("content changed while the daemon was down was admitted")
	}
	if !strings.Contains(out.String(), "DAEMON binary-unconfirmed") || !strings.Contains(out.String(), "reason=changed") {
		t.Fatalf("refusal not reported: %q", out.String())
	}
	if p, _ := l.ListPending(); len(p) != 1 || p[0].SHA256 != [32]byte{2} || p[0].Resources != res.Path {
		t.Fatalf("pending = %+v", p)
	}
}

func TestVetterRefusesNewLineAfterBootstrap(t *testing.T) {
	v, _, out := newTestVetter(t)
	v.endBootstrap()
	bin, f := testBinary(t, "planted", "evil")
	rule := daemonconfig.BinaryRule{Path: bin}
	cfg, res := vetResource(rule)
	v.begin(cfg)
	if v.judge(res, rule, bin, f, guard.GuardInodeKey{Dev: 1, Ino: 11}, [32]byte{3}) {
		t.Fatal("a binary never recorded and with no provenance was admitted")
	}
	if !strings.Contains(out.String(), "reason=new") {
		t.Fatalf("refusal not reported as new: %q", out.String())
	}
}

// A refused binary must not resolve for the trust guard: with updater bits it could write a new
// inode at a confirmed path and have it admitted by provenance.
func TestVetterResolverOnlyApproved(t *testing.T) {
	v, l, _ := newTestVetter(t)
	v.endBootstrap()
	good, gf := testBinary(t, "good", "g")
	bad, bf := testBinary(t, "bad", "b")
	if err := l.Record(good, [32]byte{7}, binledger.SourceInstall); err != nil {
		t.Fatal(err)
	}
	rg, rb := daemonconfig.BinaryRule{Path: good}, daemonconfig.BinaryRule{Path: bad}
	cfg, res := vetResource(rg, rb)
	v.begin(cfg)
	gk := guard.GuardInodeKey{Dev: 1, Ino: 20}
	if !v.judge(res, rg, good, gf, gk, [32]byte{7}) || v.judge(res, rb, bad, bf, guard.GuardInodeKey{Dev: 1, Ino: 21}, [32]byte{8}) {
		t.Fatal("unexpected verdicts")
	}
	if dev, ino, err := v.resolve(good); err != nil || dev != gk.Dev || ino != gk.Ino {
		t.Fatalf("approved binary resolves to %d:%d (%v)", dev, ino, err)
	}
	if _, _, err := v.resolve(bad); !errors.Is(err, guard.ErrUnconfirmed) {
		t.Fatalf("refused binary resolves: %v", err)
	}
}

// One path judged on two different inodes within a generation was swapped between the opens:
// the second is refused and neither keeps write rights.
func TestVetterSwapWithinGenerationLosesRights(t *testing.T) {
	v, l, _ := newTestVetter(t)
	v.endBootstrap()
	bin, f := testBinary(t, "app", "a")
	if err := l.Record(bin, [32]byte{5}, binledger.SourceInstall); err != nil {
		t.Fatal(err)
	}
	rule := daemonconfig.BinaryRule{Path: bin}
	cfg, res := vetResource(rule)
	v.begin(cfg)
	if !v.judge(res, rule, bin, f, guard.GuardInodeKey{Dev: 1, Ino: 30}, [32]byte{5}) {
		t.Fatal("first judgement refused")
	}
	if v.judge(res, rule, bin, f, guard.GuardInodeKey{Dev: 1, Ino: 31}, [32]byte{5}) {
		t.Fatal("a second inode at the same path was admitted")
	}
	if _, _, err := v.resolve(bin); err == nil {
		t.Fatal("swapped path kept updater rights")
	}
}

// A link line is keyed on the link: re-pointing it at another confirmed binary is refused unless
// that content is the one recorded for the link.
func TestVetterLinkKeyedOnLink(t *testing.T) {
	v, l, _ := newTestVetter(t)
	v.endBootstrap()
	target, f := testBinary(t, "gh", "gh-content")
	link := filepath.Join(t.TempDir(), "claude")
	if err := l.Record(target, [32]byte{0xa}, binledger.SourceInstall); err != nil {
		t.Fatal(err)
	}
	if err := l.Record(link, [32]byte{0xb}, binledger.SourceInstall); err != nil {
		t.Fatal(err)
	}
	rule := daemonconfig.BinaryRule{Path: target, Link: link}
	cfg, res := vetResource(rule)
	v.begin(cfg)
	if v.judge(res, rule, target, f, guard.GuardInodeKey{Dev: 1, Ino: 40}, [32]byte{0xa}) {
		t.Fatal("link re-pointed at another confirmed binary was admitted")
	}
}

func TestVetterRecordLiveUnderLinkLines(t *testing.T) {
	v, l, _ := newTestVetter(t)
	v.endBootstrap()
	target, f := testBinary(t, "2.0.1", "claude-new")
	link := filepath.Join(t.TempDir(), "claude")
	cfg, _ := vetResource(daemonconfig.BinaryRule{Path: target, Link: link})
	v.begin(cfg)
	v.recordLive(target, f, guard.GuardInodeKey{Dev: 1, Ino: 50})
	for _, line := range []string{target, link} {
		if _, ok, _ := l.Lookup(line); !ok {
			t.Fatalf("live re-admission not recorded under %s", line)
		}
	}
}
