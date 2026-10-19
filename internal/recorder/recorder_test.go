package recorder

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/binledger"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/install"
)

type fakeSource struct{ ch chan ebpf.FileEvent }

func (f *fakeSource) Events() <-chan ebpf.FileEvent { return f.ch }
func (f *fakeSource) Stop()                         {}

func record(t *testing.T, targets []Target, self exeKey, evs ...ebpf.FileEvent) *Raw {
	t.Helper()
	src := &fakeSource{ch: make(chan ebpf.FileEvent, len(evs))}
	for _, ev := range evs {
		src.ch <- ev
	}
	close(src.ch)
	r := newRecorder(src, targets, self)
	r.loop()
	return r.Stop()
}

func TestObserveKeysByKernelExeAndSkipsSelf(t *testing.T) {
	targets := []Target{{Path: "/home/u/test", IsDir: true}, {Path: "/home/u/token", IsDir: false}}
	self := exeKey{1, 99}
	raw := record(t, targets, self,
		ebpf.FileEvent{Type: ebpf.EventMmap, Prot: unix.PROT_READ | unix.PROT_EXEC, Path: "/usr/lib/libc.so.6",
			ExeDev: 1, ExeIno: 10, Exe: "/usr/bin/cat"},
		ebpf.FileEvent{Type: ebpf.EventMmap, Prot: unix.PROT_READ, Path: "/usr/share/data", ExeDev: 1, ExeIno: 10},
		ebpf.FileEvent{Type: ebpf.EventRead, Path: "/home/u/test/a", Comm: "cat", ExeDev: 1, ExeIno: 10,
			Exe: "/usr/bin/cat", UID: 1000},
		// Same inode, spoofed comm: still one accessor.
		ebpf.FileEvent{Type: ebpf.EventOpen, Path: "/home/u/test", Comm: "sshd", ExeDev: 1, ExeIno: 10, UID: 1000},
		ebpf.FileEvent{Type: ebpf.EventRead, Path: "/home/u/token", ExeDev: 1, ExeIno: 99, Exe: "/usr/local/bin/app-listener"},
		ebpf.FileEvent{Type: ebpf.EventRead, Path: "/home/u/token"}, // kernel thread
		ebpf.FileEvent{Type: ebpf.EventRead, Path: "/home/u/tokenx", ExeDev: 1, ExeIno: 11},
		ebpf.FileEvent{Type: ebpf.EventRename, Path: "/home/u/x", Dest: "/home/u/token", ExeDev: 1, ExeIno: 12,
			Exe: "/usr/bin/mv"},
	)
	if len(raw.Accessors) != 2 {
		t.Fatalf("accessors = %+v", raw.Accessors)
	}
	cat := raw.Accessors[0]
	if cat.Exe != "/usr/bin/cat" || cat.Count != 2 || !slices.Equal(cat.Comms, []string{"cat", "sshd"}) {
		t.Fatalf("cat = %+v", cat)
	}
	if !slices.Equal(cat.Libs, []string{"/usr/lib/libc.so.6"}) {
		t.Fatalf("libs = %v: only code mappings count", cat.Libs)
	}
	if mv := raw.Accessors[1]; mv.Exe != "/usr/bin/mv" || !slices.Equal(mv.Targets, []string{"/home/u/token"}) {
		t.Fatalf("a rename onto a target is an access to it: %+v", mv)
	}
	if raw.Events != 3 {
		t.Fatalf("events = %d", raw.Events)
	}
}

func TestMonitorTargetsNameFilesByAbsolutePath(t *testing.T) {
	got := monitorTargets([]Target{{Path: "/home/u/.netrc"}})
	if got[0].File != "/home/u/.netrc" {
		t.Fatalf("a basename would match every .netrc on the system: %+v", got[0])
	}
}

// privateTempRoots lets fixtures live in t.TempDir (under /tmp, a refused root).
func privateTempRoots(t *testing.T) {
	saved := worldWritableRoots
	worldWritableRoots = nil
	t.Cleanup(func() { worldWritableRoots = saved })
}

// fixtureExe writes an executable in a private dir and returns it with its kernel key.
func fixtureExe(t *testing.T) (string, exeKey) {
	t.Helper()
	privateTempRoots(t)
	dir := t.TempDir()
	p := filepath.Join(dir, "tool")
	if err := os.WriteFile(p, []byte("#!/bin/true\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	dev, ino, err := ebpf.StatInode(p)
	if err != nil {
		t.Fatal(err)
	}
	return p, exeKey{dev, ino}
}

func TestClassifyVerifiesTheInodeThatRan(t *testing.T) {
	exe, k := fixtureExe(t)
	targets := []Target{{Path: "/home/u/test", IsDir: true}}
	raw := &Raw{Targets: targets, Accessors: []RawAccessor{
		{Dev: k.dev, Ino: k.ino, Exe: exe, Ops: []ebpf.EventType{ebpf.EventOpen}, Targets: []string{"/home/u/test"}},
		{Dev: k.dev, Ino: k.ino + 1, Exe: exe, Targets: []string{"/home/u/test"}},
		{Dev: 1, Ino: 2, Targets: []string{"/home/u/test"}},
		{Dev: 1, Ino: 3, Exe: "/usr/bin/gone (deleted)", Targets: []string{"/home/u/test"}},
	}}
	rec := Classify(raw, nil)
	if len(rec.Binaries) != 4 {
		t.Fatalf("binaries = %+v", rec.Binaries)
	}
	ok := rec.Binaries[0]
	if !ok.Admissible() || ok.Line != exe || ok.SystemPlaced || ok.Suggested {
		t.Fatalf("a user-owned executable is whitelistable but never preselected: %+v", ok)
	}
	content, _ := os.ReadFile(exe)
	if sum := sha256.Sum256(content); ok.SHA256 != hex.EncodeToString(sum[:]) {
		t.Fatalf("hash %s is not the file's", ok.SHA256)
	}
	if !slices.Equal(ok.Mask, []string{"OPEN", "READ", "STAT"}) {
		t.Fatalf("mask = %v", ok.Mask)
	}
	for _, b := range rec.Binaries[1:] {
		if b.Admissible() || b.Problem == "" {
			t.Fatalf("not the inode that ran, unnamed or deleted must not be whitelistable: %+v", b)
		}
	}
}

func TestClassifyRefusesWorldWritableLocations(t *testing.T) {
	privateTempRoots(t)
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o777); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "dropper")
	if err := os.WriteFile(p, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	dev, ino, _ := ebpf.StatInode(p)
	rec := Classify(&Raw{Accessors: []RawAccessor{{Dev: dev, Ino: ino, Exe: p}}}, nil)
	if rec.Binaries[0].Admissible() || !strings.Contains(rec.Binaries[0].Problem, "anyone can write") {
		t.Fatalf("%+v", rec.Binaries[0])
	}
}

func TestClassifyLibs(t *testing.T) {
	privateTempRoots(t)
	dir := t.TempDir()
	plugin := filepath.Join(dir, "plugin.so")
	if err := os.WriteFile(plugin, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	inTarget := "/home/u/test/lib.so"
	base := mustParse(t, "")
	base.SharedAllowLibs = []string{"/home/u/already.so"}
	libs, trusted := classifyLibs([]string{plugin, inTarget, "/home/u/already.so", dir + "/gone.so", "/exe"},
		"/exe", []Target{{Path: "/home/u/test", IsDir: true}}, base)
	if trusted != 2 {
		t.Fatalf("trusted = %d: the target and an existing allow_lib cover their libraries", trusted)
	}
	byPath := map[string]Lib{}
	for _, l := range libs {
		byPath[l.Path] = l
	}
	if l := byPath[plugin]; l.Problem != "" || l.Why == "" {
		t.Fatalf("a user-owned plugin needs allow_lib: %+v", l)
	}
	if byPath[dir+"/gone.so"].Problem == "" {
		t.Fatalf("a vanished library is never proposed: %+v", libs)
	}
	if err := os.Chmod(dir, 0o777); err != nil {
		t.Fatal(err)
	}
	if libs, _ := classifyLibs([]string{plugin}, "/exe", nil, base); libs[0].Problem == "" {
		t.Fatalf("a library in a world-writable directory is never proposed: %+v", libs)
	}
	worldWritableRoots = []string{"/tmp"}
	if libs, _ := classifyLibs([]string{plugin}, "/exe", nil, base); libs[0].Problem == "" {
		t.Fatalf("a library under /tmp is never proposed: %+v", libs)
	}
	if _, ok := byPath["/exe"]; ok {
		t.Fatal("the executable's own mapping is not a library")
	}
}

func mustParse(t *testing.T, text string) *daemonconfig.Config {
	t.Helper()
	cfg, err := daemonconfig.Parse([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestValidateTargets(t *testing.T) {
	home := t.TempDir()
	guarded := filepath.Join(home, ".ssh")
	secret := filepath.Join(home, "secret")
	other := filepath.Join(home, "other")
	for _, d := range []string{guarded, secret, other} {
		if err := os.Mkdir(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	link := filepath.Join(home, "link")
	if err := os.Symlink(secret, link); err != nil {
		t.Fatal(err)
	}
	hard := filepath.Join(home, "hard")
	if err := os.WriteFile(filepath.Join(secret, "f"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(filepath.Join(secret, "f"), hard); err != nil {
		t.Fatal(err)
	}
	base := mustParse(t, "[watch "+guarded+"]\n/usr/bin/ssh\nneed_encryption: false\n")
	homes := []string{home}

	got, err := ValidateTargets([]string{secret + "/", other}, base, homes)
	if err != nil || len(got) != 2 || got[0].Path != secret || !got[0].IsDir {
		t.Fatalf("got %+v, %v", got, err)
	}
	for name, paths := range map[string][]string{
		"relative":       {"home/u"},
		"system root":    {"/"},
		"system tree":    {"/usr/lib/x"},
		"daemon config":  {"/etc/app-listener"},
		"whole home":     {home},
		"inside guarded": {filepath.Join(guarded, "k")},
		"around guarded": {filepath.Dir(guarded)},
		"symlink":        {link},
		"hard link":      {hard},
		"overlap":        {secret, filepath.Join(secret, "f")},
		"missing":        {filepath.Join(home, "nope")},
		"backup":         {secret + fscrypt.BackupSuffix},
	} {
		if _, err := ValidateTargets(paths, base, homes); err == nil {
			t.Errorf("%s: %v accepted", name, paths)
		}
	}
}

func TestValidateTargetsExtendsAGuardedPath(t *testing.T) {
	home := t.TempDir()
	plain, root, libs := filepath.Join(home, ".ssh"), filepath.Join(home, "app"), filepath.Join(home, "libs")
	m1, m2 := filepath.Join(root, "one"), filepath.Join(root, "two")
	for _, d := range []string{plain, root, m1, m2, libs} {
		if err := os.Mkdir(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	base := mustParse(t, "[watch "+plain+"]\n/usr/bin/ls OPEN,READ\nneed_encryption: false\n\n"+
		"[watch "+root+"]\nwatch: "+m1+"\nwatch: "+m2+"\n/usr/bin/cat\nlib_dir "+libs+"\nneed_encryption: false\n")
	got, err := ValidateTargets([]string{plain, m1}, base, []string{home})
	if err != nil {
		t.Fatal(err)
	}
	ex := got[0].Existing
	if ex == nil || ex.Section != plain || !slices.Equal(ex.Lines(), []string{"/usr/bin/ls OPEN,READ"}) ||
		len(ex.SharedWith) != 0 {
		t.Fatalf("plain section: %+v", ex)
	}
	ex = got[1].Existing
	if ex == nil || ex.Section != root || !ex.Has("/usr/bin/cat") || !slices.Equal(ex.SharedWith, []string{m2}) {
		t.Fatalf("group member: %+v", ex)
	}
	for name, p := range map[string]string{"encryption root": root, "lib_dir": libs,
		"inside a member": filepath.Join(m1, "x")} {
		if _, err := ValidateTargets([]string{p}, base, nil); err == nil {
			t.Errorf("%s: %s accepted", name, p)
		}
	}
}

func TestPlanExtendsAnExistingSection(t *testing.T) {
	rec, a, b := testRecording(t)
	rec.Targets[0].Existing = &Existing{Section: a, Allow: []ExistingRule{{Line: "/usr/bin/ssh"}}}
	rec.Binaries[0].Whitelisted = []string{a}
	for name, sel := range map[string]Selection{
		"covered binary":   {Binaries: []string{"/usr/bin/ssh"}},
		"encrypt existing": {Binaries: []string{"/usr/bin/git"}, Encrypt: []string{a}},
	} {
		if _, err := rec.Plan(&sel); err == nil {
			t.Errorf("%s accepted", name)
		}
	}
	rec.Targets[1].Existing = &Existing{Section: b}
	if _, err := rec.Plan(&Selection{}); err == nil {
		t.Fatal("a plan adding nothing must not apply")
	}
	rec.Targets[1].Existing = nil
	p, err := rec.Plan(&Selection{Binaries: []string{"/usr/bin/git"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Sections) != 1 || p.Sections[0].Path != b || len(p.Extensions) != 1 || p.Extensions[0].Section != a ||
		len(p.Extensions[0].Allow) != 1 || p.Extensions[0].Allow[0].Path != "/usr/bin/git" {
		t.Fatalf("plan: %+v", p)
	}
	base := "[watch " + a + "]\n/usr/bin/ssh\n\nneed_encryption: false\n"
	text, err := p.Render([]byte(base))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(text, "[watch "+a+"]\n/usr/bin/ssh\n\"/usr/bin/git\"\n\nneed_encryption: false\n") {
		t.Fatalf("the existing section keeps its lines and gains the new one:\n%s", text)
	}
	p.Sections = nil
	again, err := p.Render([]byte(text))
	if err != nil || !slices.Equal(install.ParseSectionWhitelist(again, a), []string{"/usr/bin/ssh", "/usr/bin/git"}) {
		t.Fatalf("re-rendering over a base that has the line never duplicates it: %v\n%s", err, again)
	}
	if _, err := p.Render([]byte("")); err == nil {
		t.Fatal("an extension whose section vanished must not apply")
	}
}

func TestNoteWhitelisted(t *testing.T) {
	targets := []Target{{Path: "/a", Existing: &Existing{Allow: []ExistingRule{{Line: "/usr/bin/x",
		Events: []string{"READ"}}}}}, {Path: "/b"}}
	b := &Binary{Line: "/usr/bin/x", Targets: []string{"/a", "/b"}, Ops: []string{"OPEN", "READ", "STAT", "WRITE"}}
	noteWhitelisted(b, targets)
	if !slices.Equal(b.Whitelisted, []string{"/a"}) || b.Covered() {
		t.Fatalf("whitelisted for /a only: %+v", b)
	}
	if !slices.ContainsFunc(b.Notes, func(n string) bool { return strings.Contains(n, "but seen doing WRITE too") }) {
		t.Fatalf("only WRITE is beyond READ (it implies OPEN and STAT): %q", b.Notes)
	}
	b = &Binary{Line: "/usr/bin/x", Targets: []string{"/a"}}
	finishBinary(b, exeVerdict{placed: true}, targets)
	if !b.Covered() || b.Suggested {
		t.Fatalf("a binary whitelisted everywhere it went starts deselected: %+v", b)
	}
}

func testRecording(t *testing.T) (*Recording, string, string) {
	t.Helper()
	home := t.TempDir()
	a, b := filepath.Join(home, "a"), filepath.Join(home, "b")
	for _, d := range []string{a, b} {
		if err := os.Mkdir(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	sum := sha256.Sum256(nil)
	rec := &Recording{Targets: []Target{{Path: a, IsDir: true}, {Path: b, IsDir: true}}, Binaries: []Binary{
		{Line: "/usr/bin/ssh", Targets: []string{a}, Mask: []string{"OPEN", "READ", "STAT"},
			SHA256: hex.EncodeToString(sum[:]), SystemPlaced: true},
		{Line: "/usr/bin/git", Targets: []string{a, b}, SHA256: hex.EncodeToString(sum[:]),
			Libs: []Lib{{Path: "/home/u/plugin.so", Why: "x"}, {Path: "/tmp/evil.so", Problem: "y"}}},
		{Exe: "/usr/bin/gone", Problem: "deleted", Targets: []string{a}},
	}}
	return rec, a, b
}

func TestPlanAdmitsOnlyWhatWasRecorded(t *testing.T) {
	rec, a, _ := testRecording(t)
	for name, sel := range map[string]Selection{
		"unrecorded binary":   {Binaries: []string{"/usr/bin/curl"}},
		"inadmissible binary": {Binaries: []string{"/usr/bin/gone"}},
		"lib of unselected":   {Binaries: []string{"/usr/bin/ssh"}, Libs: []string{"/home/u/plugin.so"}},
		"refused lib":         {Binaries: []string{"/usr/bin/git"}, Libs: []string{"/tmp/evil.so"}},
		"foreign encrypt":     {Encrypt: []string{"/home/u/.ssh"}},
	} {
		if _, err := rec.Plan(&sel); err == nil {
			t.Errorf("%s accepted", name)
		}
	}
	p, err := rec.Plan(&Selection{Binaries: []string{"/usr/bin/ssh", "/usr/bin/git"}, Libs: []string{"/home/u/plugin.so"},
		Encrypt: []string{a}, Masks: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Sections) != 2 || len(p.Sections[0].Allow) != 2 || len(p.Sections[1].Allow) != 1 {
		t.Fatalf("each target lists only the binaries that touched it: %+v", p.Sections)
	}
	if !p.Sections[0].Encrypt || p.Sections[1].Encrypt {
		t.Fatalf("encrypt flags: %+v", p.Sections)
	}
	if ssh := p.Sections[0].Allow[1]; ssh.Path != "/usr/bin/ssh" || !slices.Equal(ssh.Events, []string{"OPEN", "READ", "STAT"}) {
		t.Fatalf("masked rule = %+v", ssh)
	}
}

func TestRenderChecksTheSectionsMaterialize(t *testing.T) {
	rec, a, b := testRecording(t)
	p, err := rec.Plan(&Selection{Binaries: []string{"/usr/bin/ssh"}, Encrypt: []string{a}})
	if err != nil {
		t.Fatal(err)
	}
	base := []byte("[libraries \"x\"]\nallow_lib /usr/lib/x.so\n")
	text, err := p.Render(base)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(text, "[watch "+`"`+a+`"`) > strings.Index(text, "[libraries") {
		t.Fatalf("watch sections go before the library blocks:\n%s", text)
	}
	cfg := mustParse(t, text)
	if len(cfg.Resources) != 2 {
		t.Fatalf("resources = %+v", cfg.Resources)
	}
	if err := os.Remove(b); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Render(base); err == nil {
		t.Fatal("a section the parser skips must not apply as if it guarded something")
	}
}

type fakeStore struct {
	conf    []byte
	stale   int
	putErr  error
	guarded []string
	puts    int
}

func (s *fakeStore) Fetch() ([]byte, error) { return s.conf, nil }
func (s *fakeStore) Put(base, next []byte) error {
	s.puts++
	if s.stale > 0 {
		s.stale--
		s.conf = append(append([]byte(nil), s.conf...), "# concurrent edit\n"...)
		return ErrStale
	}
	if string(base) != string(s.conf) {
		return ErrStale
	}
	if s.putErr != nil {
		return s.putErr
	}
	s.conf = next
	return nil
}
func (s *fakeStore) Guarded() ([]string, error) { return s.guarded, nil }

type fakeVault struct {
	prereqs   []fscrypt.Prereq
	encrypted []string
	restored  []string
}

func (v *fakeVault) FilesystemPrereqs(string) ([]fscrypt.Prereq, error) { return v.prereqs, nil }
func (v *fakeVault) EncryptWithProgress(p string, _ func(int64, int64)) error {
	v.encrypted = append(v.encrypted, p)
	return nil
}
func (v *fakeVault) RestoreBackup(p string) error {
	v.restored = append(v.restored, p)
	return nil
}

type fakeLedger struct{ rows map[string][32]byte }

func (l *fakeLedger) RecordIfAbsent(line string, sum [32]byte, _ binledger.Source) (bool, error) {
	if _, ok := l.rows[line]; ok {
		return false, nil
	}
	l.rows[line] = sum
	return true, nil
}
func (l *fakeLedger) Forget(line string, sum [32]byte, _ binledger.Source) error {
	if l.rows[line] == sum {
		delete(l.rows, line)
	}
	return nil
}
func (l *fakeLedger) Lookup(line string) ([32]byte, bool, error) {
	s, ok := l.rows[line]
	return s, ok, nil
}

func testPlan(t *testing.T) (*Plan, string) {
	rec, a, _ := testRecording(t)
	p, err := rec.Plan(&Selection{Binaries: []string{"/usr/bin/ssh", "/usr/bin/git"}, Encrypt: []string{a}})
	if err != nil {
		t.Fatal(err)
	}
	return p, a
}

func TestApplyPinsUserBinariesAndRetriesStale(t *testing.T) {
	p, a := testPlan(t)
	store := &fakeStore{stale: 1, guarded: []string{p.Sections[0].Path}}
	vault, ledger := &fakeVault{}, &fakeLedger{rows: map[string][32]byte{}}
	ap := &applier{store: store, vault: vault, ledger: ledger, systemSetup: func() error { return nil }}
	out, err := ap.apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if store.puts != 2 || !strings.Contains(string(store.conf), "# concurrent edit") {
		t.Fatalf("a concurrent edit is kept and the plan re-rendered over it (puts=%d):\n%s", store.puts, store.conf)
	}
	if _, ok := ledger.rows["/usr/bin/git"]; !ok {
		t.Fatal("a binary that isn't root-placed is pinned to the recorded content")
	}
	if _, ok := ledger.rows["/usr/bin/ssh"]; ok {
		t.Fatal("a root-placed binary needs no ledger row")
	}
	if !slices.Equal(out.Backups, []string{a + fscrypt.BackupSuffix}) {
		t.Fatalf("backups = %v", out.Backups)
	}
	if len(out.Notes) != 1 || !strings.Contains(out.Notes[0], p.Sections[1].Path) {
		t.Fatalf("a section the daemon doesn't report is flagged: %v", out.Notes)
	}
}

func TestApplyRollsBackOnRefusal(t *testing.T) {
	p, a := testPlan(t)
	store := &fakeStore{putErr: errors.New("reload failed")}
	vault := &fakeVault{}
	other := [32]byte{9}
	ledger := &fakeLedger{rows: map[string][32]byte{"/usr/bin/other": other}}
	ap := &applier{store: store, vault: vault, ledger: ledger, systemSetup: func() error { return nil }}
	if _, err := ap.apply(p); err == nil {
		t.Fatal("refused put succeeded")
	}
	if !slices.Equal(vault.restored, []string{a}) {
		t.Fatalf("encrypted targets are restored from their backups: %v", vault.restored)
	}
	if _, ok := ledger.rows["/usr/bin/git"]; ok || ledger.rows["/usr/bin/other"] != other {
		t.Fatalf("only the pins this apply wrote are forgotten: %v", ledger.rows)
	}
}

func TestApplyKeepsConfirmedLedgerVerdict(t *testing.T) {
	p, _ := testPlan(t)
	confirmed := [32]byte{7}
	ledger := &fakeLedger{rows: map[string][32]byte{"/usr/bin/git": confirmed}}
	ap := &applier{store: &fakeStore{}, vault: &fakeVault{}, ledger: ledger, systemSetup: func() error { return nil }}
	out, err := ap.apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if ledger.rows["/usr/bin/git"] != confirmed {
		t.Fatal("a recording must never re-confirm a line the ledger holds for other content")
	}
	if !slices.ContainsFunc(out.Notes, func(n string) bool { return strings.Contains(n, "trust") }) {
		t.Fatalf("notes = %v", out.Notes)
	}
}

func TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(t *testing.T) {
	p, _ := testPlan(t)
	vault := &fakeVault{prereqs: []fscrypt.Prereq{{Reason: "no metadata", Argv: []string{"fscrypt", "setup"}}}}
	ledger := &fakeLedger{rows: map[string][32]byte{}}
	store := &fakeStore{}
	ap := &applier{store: store, vault: vault, ledger: ledger, systemSetup: func() error { return nil }}
	if _, err := ap.apply(p); err == nil || !strings.Contains(err.Error(), "fscrypt setup") {
		t.Fatalf("err = %v", err)
	}
	if len(vault.encrypted)+len(ledger.rows)+store.puts != 0 {
		t.Fatal("nothing may change before the filesystem check")
	}
}

func TestDeleteBackupOnlyItsOwn(t *testing.T) {
	dir := t.TempDir()
	b := filepath.Join(dir, "a"+fscrypt.BackupSuffix)
	if err := os.Mkdir(b, 0o700); err != nil {
		t.Fatal(err)
	}
	o := &Outcome{Backups: []string{b}}
	if err := DeleteBackup(o, filepath.Join(dir, "a")); err == nil {
		t.Fatal("deleted a path that is not one of its backups")
	}
	if err := DeleteBackup(o, b); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(b); !os.IsNotExist(err) {
		t.Fatal("backup kept")
	}
}

func TestProgressNamesMulticallApplets(t *testing.T) {
	r := newRecorder(&fakeSource{}, []Target{{Path: "/t", IsDir: true}}, exeKey{})
	r.observe(&ebpf.FileEvent{Type: ebpf.EventRead, Path: "/t/x", ExeDev: 1, ExeIno: 2, Exe: "/usr/bin/coreutils",
		ExecName: "/usr/bin/cat"})
	p := r.Progress()
	if p.Binaries != 1 || p.Recent[0].Exe != "/usr/bin/coreutils as cat" || time.Since(p.Recent[0].Last) > time.Minute {
		t.Fatalf("%+v", p)
	}
}

func TestFileReadersStartDeselected(t *testing.T) {
	for _, line := range []string{"/usr/bin/cat", "/usr/bin/dd", "/usr/bin/head", "/usr/bin/base64",
		"/usr/bin/xxd", "/usr/bin/sed", "/usr/bin/awk", "/usr/bin/grep", "/usr/bin/openssl", "/usr/bin/nc",
		"/usr/bin/bzip2", "/usr/bin/base32"} {
		b := &Binary{Line: line, Exe: line}
		finishBinary(b, exeVerdict{placed: true}, nil)
		if !b.GeneralTool || b.Suggested {
			t.Errorf("%s: GeneralTool=%v Suggested=%v", line, b.GeneralTool, b.Suggested)
		}
	}
}

func TestShellHelpersStartDeselected(t *testing.T) {
	for line, want := range map[string]bool{
		"/usr/bin/starship": true, "/home/u/.cache/gitstatus/gitstatusd-linux-x86_64": true,
		"/usr/bin/zoxide": true, "/usr/bin/ls": false, "/usr/bin/starshipd": false,
	} {
		b := &Binary{Line: line, Exe: line}
		finishBinary(b, exeVerdict{placed: true}, nil)
		if b.ShellHelper != want || b.Suggested == want {
			t.Errorf("%s: ShellHelper=%v Suggested=%v", line, b.ShellHelper, b.Suggested)
		}
	}
}
