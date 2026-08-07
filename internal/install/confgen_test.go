package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func mustParse(t *testing.T, text string) *daemonconfig.Config {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "daemon.conf")
	if err := os.WriteFile(tmp, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := daemonconfig.Load(tmp)
	if err != nil {
		t.Fatalf("daemonconfig.Load: %v", err)
	}
	return cfg
}

// TestGenerateConfParses verifies the generated config is accepted by the
// real daemon parser and carries the expected sections and directives.
func TestGenerateConfParses(t *testing.T) {
	sshDir := filepath.Join(t.TempDir(), "ssh")
	openDir := filepath.Join(t.TempDir(), "opencode")
	if err := os.MkdirAll(sshDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(openDir, 0o700); err != nil {
		t.Fatal(err)
	}
	conf := GenerateConf([]Section{
		{Path: sshDir, Allow: []BinaryRule{
			{Path: "/usr/bin/ssh", Events: []string{"READ", "WRITE"}},
			{Path: "/usr/bin/ssh-agent"},
		}, Encrypt: true},
		{Path: openDir, Encrypt: false},
	})
	cfg := mustParse(t, conf)
	if len(cfg.Resources) != 2 {
		t.Fatalf("resources = %d, want 2", len(cfg.Resources))
	}
	if cfg.Resources[0].Path != sshDir || len(cfg.Resources[0].Binaries) != 2 {
		t.Errorf("first resource = %+v", cfg.Resources[0])
	}
	if !cfg.Resources[0].NeedEncryption {
		t.Error("first resource should need encryption")
	}
	if cfg.Resources[1].NeedEncryption {
		t.Error("second resource should not need encryption")
	}
	if cfg.Resources[1].Path != openDir {
		t.Errorf("second resource = %+v", cfg.Resources[1])
	}
}

// TestGenerateConfEmitsEvents verifies whitelisted binaries with event
// restrictions are emitted as "<path> EV1,EV2" and that the events survive
// the real daemon parser round-trip. Binaries without restrictions stay a
// bare path and are parsed back with an empty event list.
func TestGenerateConfEmitsEvents(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ssh")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	conf := GenerateConf([]Section{
		{Path: dir, Allow: []BinaryRule{
			{Path: "/usr/bin/ssh", Events: []string{"READ", "WRITE"}},
			{Path: "/usr/bin/ssh-agent"},
		}, Encrypt: true},
	})
	if !strings.Contains(conf, "/usr/bin/ssh READ,WRITE\n") {
		t.Fatalf("config missing READ,WRITE restriction:\n%s", conf)
	}
	if !strings.Contains(conf, "/usr/bin/ssh-agent\n") {
		t.Fatalf("config missing bare ssh-agent line:\n%s", conf)
	}
	cfg := mustParse(t, conf)
	res := cfg.Resources[0]
	if len(res.Binaries) != 2 {
		t.Fatalf("binaries = %d, want 2: %+v", len(res.Binaries), res.Binaries)
	}
	if res.Binaries[0].Path != "/usr/bin/ssh" || len(res.Binaries[0].Events) != 2 {
		t.Errorf("ssh rule = %+v, want READ,WRITE", res.Binaries[0])
	}
	if len(res.Binaries[1].Events) != 0 {
		t.Errorf("ssh-agent rule = %+v, want unrestricted (empty events)", res.Binaries[1])
	}
}

// TestSetNeedEncryptionReplace flips an existing directive, preserving the
// whitelist around it.
func TestSetNeedEncryptionReplace(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ssh")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	conf := GenerateConf([]Section{
		{Path: dir, Allow: []BinaryRule{{Path: "/usr/bin/ssh"}}, Encrypt: true},
	})
	updated, err := SetNeedEncryption(conf, dir, false)
	if err != nil {
		t.Fatal(err)
	}
	cfg := mustParse(t, updated)
	if cfg.Resources[0].NeedEncryption {
		t.Error("directive was not flipped to false")
	}
	if len(cfg.Resources[0].Binaries) != 1 {
		t.Error("whitelist was not preserved")
	}
}

// TestSetNeedEncryptionInsert adds the directive to a section that lacks
// one (default true in the parser, so flipping to false must insert it).
func TestSetNeedEncryptionInsert(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "secret")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	conf := GenerateConf([]Section{{Path: dir, Allow: nil, Encrypt: true}})
	noDirective := strings.Replace(conf, "need_encryption: true", "", 1)
	updated, err := SetNeedEncryption(noDirective, dir, false)
	if err != nil {
		t.Fatal(err)
	}
	cfg := mustParse(t, updated)
	if cfg.Resources[0].NeedEncryption {
		t.Error("directive was not inserted as false")
	}
}

// TestSetNeedEncryptionMissingSection reports an error for an unknown path.
func TestSetNeedEncryptionMissingSection(t *testing.T) {
	conf := GenerateConf([]Section{{Path: "/srv/secret"}})
	if _, err := SetNeedEncryption(conf, "/does/not/exist", true); err == nil {
		t.Fatal("expected error for missing section")
	}
}

// TestSetNeedEncryptionPreservesManualEdits verifies a comment added by the
// user in the embedded editor survives the directive update.
func TestSetNeedEncryptionPreservesManualEdits(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ssh")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	conf := GenerateConf([]Section{{Path: dir, Allow: []BinaryRule{{Path: "/usr/bin/ssh"}}}})
	withComment := strings.Replace(conf,
		"[watch "+dir+"]\n",
		"[watch "+dir+"]\n# keep this comment\n",
		1)
	updated, err := SetNeedEncryption(withComment, dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(updated, "# keep this comment") {
		t.Error("manual comment was lost")
	}
}
