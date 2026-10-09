package guard

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	logtest "github.com/sirupsen/logrus/hooks/test"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestAppletKeyTag(t *testing.T) {
	real := GuardInodeKey{Dev: 0x800003, Ino: 42}
	k := real.withTag(7)
	if !k.tagged() || k.Real() != real || k.Dev>>mcTagShift != 7 {
		t.Fatalf("withTag(7) = %+v", k)
	}
	if real.tagged() || real.Real() != real {
		t.Fatal("a real key must be untouched by Real")
	}
	// INODE_DEV_UNKNOWN (bit 63) is no tag and survives Real.
	unknown := GuardInodeKey{Dev: 1 << 63, Ino: 1}
	if unknown.tagged() || unknown.Real() != unknown {
		t.Fatal("bit 63 must not read as a tag")
	}
}

func TestAppletTagStablePerName(t *testing.T) {
	a, err := appletTag("multicall-test-a")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := appletTag("multicall-test-b")
	again, _ := appletTag("multicall-test-a")
	if a == b || a != again || a == 0 || a >= mcTagNone {
		t.Fatalf("tags a=%d b=%d again=%d", a, b, again)
	}
}

func TestWithRealKeysAddsFileKeyOfApplets(t *testing.T) {
	real := GuardInodeKey{Dev: 3, Ino: 9}
	single := GuardInodeKey{Dev: 3, Ino: 10}
	allows := map[uint32]map[GuardInodeKey]uint8{
		1: {real.withTag(4): GUARD_ALLOW, single: GUARD_ALLOW},
		2: {single: GUARD_ALLOW},
	}
	got := withRealKeys(allows)
	if _, ok := got[1][real]; !ok {
		t.Fatal("resource 1 allows an applet of real: its real key must join the union")
	}
	if _, ok := got[2][real]; ok {
		t.Fatal("resource 2 allows no applet of real")
	}
	if _, ok := allows[1][real]; ok {
		t.Fatal("withRealKeys must not modify its input")
	}
}

func TestMcNameKey(t *testing.T) {
	k := mcNameKey(GuardInodeKey{Dev: 1, Ino: 2}, "sha256sum")
	var b strings.Builder
	for _, c := range k.Name {
		if c == 0 {
			break
		}
		b.WriteByte(byte(c))
	}
	if b.String() != "sha256sum" || k.Name[len(k.Name)-1] != 0 {
		t.Fatalf("name key %q", b.String())
	}
}

func TestExeKeyRefusesOpaqueMulticall(t *testing.T) {
	dir := t.TempDir()
	bb := filepath.Join(dir, "busybox")
	marker := []byte("v xoBysuB") // reversed: spelled forward it would be in this test binary
	slices.Reverse(marker)
	if err := os.WriteFile(bb, append([]byte("\x7fELF...."), marker...), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(bb, filepath.Join(dir, "cat")); err != nil {
		t.Fatal(err)
	}
	plain := filepath.Join(dir, "tool")
	if err := os.WriteFile(plain, []byte("\x7fELF plain"), 0o755); err != nil {
		t.Fatal(err)
	}
	for path, refused := range map[string]bool{filepath.Join(dir, "cat"): true, plain: false} {
		dev, ino, err := ebpf.StatInode(path)
		if err != nil {
			t.Fatal(err)
		}
		real := GuardInodeKey{Dev: dev, Ino: ino}
		got, err := ExeKey(path, nil, real)
		if IsMulticallRefusal(err) != refused {
			t.Fatalf("%s: err %v, want refused=%v", path, err, refused)
		}
		if !refused && got != real {
			t.Fatalf("%s: a single binary's key must be its inode, got %+v", path, got)
		}
		if refused && !strings.Contains(err.Error(), "busybox") {
			t.Fatalf("refusal must name the siblings: %v", err)
		}
	}
}

func TestRefuseMulticallLines(t *testing.T) {
	dir := t.TempDir()
	bb := filepath.Join(dir, "wget")
	marker := []byte("v xoBysuB") // reversed: see TestExeKeyRefusesOpaqueMulticall
	slices.Reverse(marker)
	if err := os.WriteFile(bb, append([]byte("\x7fELF...."), marker...), 0o755); err != nil {
		t.Fatal(err)
	}
	plain := filepath.Join(dir, "tool")
	if err := os.WriteFile(plain, []byte("\x7fELF plain"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := MulticallRefusal(bb); !IsMulticallRefusal(err) {
		t.Fatalf("busybox: %v", err)
	}
	for _, p := range []string{plain, filepath.Join(dir, "missing")} {
		if err := MulticallRefusal(p); err != nil {
			t.Fatalf("%s: %v", p, err)
		}
	}
	cfgOf := func(paths ...string) *daemonconfig.Config {
		r := daemonconfig.Resource{Path: dir}
		for _, p := range paths {
			r.Binaries = append(r.Binaries, daemonconfig.BinaryRule{Path: p})
		}
		return &daemonconfig.Config{Resources: []daemonconfig.Resource{r}}
	}
	if err := RefuseMulticallLines(cfgOf(plain, bb), nil); !IsMulticallRefusal(err) {
		t.Fatalf("a new busybox line must be refused: %v", err)
	}
	if err := RefuseMulticallLines(cfgOf(plain, bb), cfgOf(bb)); err != nil {
		t.Fatalf("a line prev already has must stay: %v", err)
	}
	libOnly := cfgOf()
	libOnly.Resources[0].Binaries = []daemonconfig.BinaryRule{{Path: bb, LibBinary: true}}
	if err := RefuseMulticallLines(libOnly, nil); err != nil {
		t.Fatalf("a lib_binary writer is no whitelist identity: %v", err)
	}
}

func TestExeKeyRefusesUnrecognizedMulticall(t *testing.T) {
	p := filepath.Join(t.TempDir(), "cmp")
	marker := []byte("yranib llac-itlum") // reversed: see TestExeKeyRefusesOpaqueMulticall
	slices.Reverse(marker)
	if err := os.WriteFile(p, append([]byte("\x7fELF...."), marker...), 0o755); err != nil {
		t.Fatal(err)
	}
	dev, ino, err := ebpf.StatInode(p)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ExeKey(p, nil, GuardInodeKey{Dev: dev, Ino: ino})
	if !IsMulticallRefusal(err) || !strings.Contains(err.Error(), "is an unrecognized multicall binary") {
		t.Fatalf("an unvetted multicall must be refused: %v", err)
	}
}

func TestExeKeyRefusesGNUCoreutilsSingle(t *testing.T) {
	dir := t.TempDir()
	marker := []byte("gorp-slitueroc--") // reversed: see TestExeKeyRefusesOpaqueMulticall
	slices.Reverse(marker)
	bin := filepath.Join(dir, "coreutils")
	if err := os.WriteFile(bin, append([]byte("\x7fELF...."), marker...), 0o755); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "cat")
	content := "#!" + bin + " " + string(marker) + "-shebang=cat\n"
	if err := os.WriteFile(script, []byte(content), 0o555); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{bin, script} {
		dev, ino, err := ebpf.StatInode(p)
		if err != nil {
			t.Fatal(err)
		}
		_, err = ExeKey(p, nil, GuardInodeKey{Dev: dev, Ino: ino})
		if !IsMulticallRefusal(err) || !strings.Contains(err.Error(), "is a gnu-coreutils multicall binary") ||
			!strings.Contains(err.Error(), "dnf swap coreutils-single coreutils") {
			t.Fatalf("%s: GNU coreutils-single must be refused: %v", p, err)
		}
	}
}

func TestWarnMultiLinkedOncePerInode(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "perl")
	if err := os.WriteFile(p, []byte("\x7fELF plain"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(p, filepath.Join(dir, "perl5.40")); err != nil {
		t.Fatal(err)
	}
	single := filepath.Join(dir, "tool")
	if err := os.WriteFile(single, []byte("\x7fELF plain"), 0o755); err != nil {
		t.Fatal(err)
	}
	hook := logtest.NewGlobal()
	defer hook.Reset()
	for _, path := range []string{p, p, single} {
		dev, ino, err := ebpf.StatInode(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := ExeKey(path, nil, GuardInodeKey{Dev: dev, Ino: ino}); err != nil {
			t.Fatalf("%s: a single binary is never refused: %v", path, err)
		}
	}
	var warned []string
	for _, e := range hook.AllEntries() {
		if strings.Contains(e.Message, "hard links") {
			warned = append(warned, e.Message)
		}
	}
	if len(warned) != 1 || !strings.Contains(warned[0], "also: perl5.40") {
		t.Fatalf("want one warning naming perl5.40, got %q", warned)
	}
}
