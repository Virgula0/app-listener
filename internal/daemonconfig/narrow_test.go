package daemonconfig

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// narrowTree makes <tmp>/gnupg with private-keys-v1.d and gpg.conf.
func narrowTree(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "gnupg")
	if err := os.MkdirAll(filepath.Join(root, "private-keys-v1.d"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "gpg.conf"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

func narrowParentSection(root string, encrypted bool) string {
	return "[watch " + root + "]\nneed_encryption: " + map[bool]string{true: "true", false: "false"}[encrypted] +
		"\n/usr/bin/gpg\n/usr/bin/gpg-agent\n/usr/bin/gpgconf READ,STAT,CONNECT\n"
}

func TestNarrowSectionParsed(t *testing.T) {
	root := narrowTree(t)
	keys := filepath.Join(root, "private-keys-v1.d")
	conf := filepath.Join(root, "gpg.conf")
	cfg, err := Load(writeConfig(t, narrowParentSection(root, false)+
		"\n[narrow "+keys+"]\n/usr/bin/gpg-agent\n"+
		"\n[narrow \""+conf+"\"]\n/usr/bin/gpg READ,STAT\n/usr/bin/gpgconf READ\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 3 {
		t.Fatalf("resources = %+v", cfg.Resources)
	}
	k, c := cfg.Resources[1], cfg.Resources[2]
	if k.Path != keys || k.NarrowOf != root || k.EncryptionRoot != root || k.NeedEncryption || len(k.Binaries)+len(k.PendingBinaries) != 1 {
		t.Errorf("keys narrow = %+v", k)
	}
	if c.Path != conf || c.NarrowOf != root || len(c.Binaries)+len(c.PendingBinaries) != 2 {
		t.Errorf("conf narrow = %+v", c)
	}
}

func TestNarrowSectionRefused(t *testing.T) {
	root := narrowTree(t)
	keys := filepath.Join(root, "private-keys-v1.d")
	parent := narrowParentSection(root, false)
	other := t.TempDir()
	for name, text := range map[string]string{
		"not in the parent":     parent + "[narrow " + keys + "]\n/usr/bin/ssh\n",
		"wider than the parent": parent + "[narrow " + keys + "]\n/usr/bin/gpgconf READ,WRITE\n",
		"all events on a list":  parent + "[narrow " + keys + "]\n/usr/bin/gpgconf\n",
		"before its parent":     "[narrow " + keys + "]\n/usr/bin/gpg-agent\n" + parent,
		"outside any tree":      parent + "[narrow " + other + "]\n/usr/bin/gpg-agent\n",
		"a directive":           parent + "[narrow " + keys + "]\nneed_encryption: false\n/usr/bin/gpg-agent\n",
		"a watch path":          parent + "[narrow " + keys + "]\nwatch: " + keys + "/x\n",
		"a lib_dir":             parent + "[narrow " + keys + "]\nlib_dir " + keys + "\n",
		"missing":               parent + "[narrow " + filepath.Join(root, "gpg-agent.conf") + "]\n/usr/bin/gpg-agent\n",
		"relative":              parent + "[narrow gnupg/x]\n/usr/bin/gpg-agent\n",
		"narrow in narrow": parent + "[narrow " + keys + "]\n/usr/bin/gpg-agent\n" +
			"[narrow " + keys + "/sub]\n/usr/bin/gpg-agent\n",
		"twice":            parent + "[narrow " + keys + "]\n/usr/bin/gpg-agent\n[narrow " + keys + "]\n/usr/bin/gpg-agent\n",
		"a nested [watch]": parent + "[watch " + keys + "]\nneed_encryption: false\n/usr/bin/gpg-agent\n",
	} {
		t.Run(name, func(t *testing.T) {
			if name == "narrow in narrow" {
				if err := os.Mkdir(filepath.Join(keys, "sub"), 0o700); err != nil && !os.IsExist(err) {
					t.Fatal(err)
				}
			}
			if _, err := Load(writeConfig(t, text)); err == nil {
				t.Fatalf("must be refused:\n%s", text)
			}
		})
	}
}

// Inside a locked vault the target's name is unreadable: deferred, re-validated after unlock.
func TestNarrowSectionPendingInLockedVault(t *testing.T) {
	root := narrowTree(t)
	hidden := filepath.Join(root, "not-yet-visible")
	cfg, err := Load(writeConfig(t, narrowParentSection(root, true)+"[narrow "+hidden+"]\n/usr/bin/gpg-agent\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	n := cfg.Resources[1]
	if !n.PathPending || !n.NeedEncryption || n.EncryptionRoot != root {
		t.Fatalf("narrow = %+v", n)
	}
	if err := ResolvePendingPaths(cfg); err == nil || !strings.Contains(err.Error(), hidden) {
		t.Fatalf("a narrow still missing after unlock must be a hard error, got %v", err)
	}
	if err := os.Mkdir(hidden, 0o700); err != nil {
		t.Fatal(err)
	}
	cfg, _ = Load(writeConfig(t, narrowParentSection(root, true)+"[narrow "+hidden+"]\n/usr/bin/gpg-agent READ\n"))
	if err := ResolvePendingPaths(cfg); err != nil || cfg.Resources[1].PathPending {
		t.Fatalf("visible after unlock: %v %+v", err, cfg.Resources[1])
	}
	if ev := slices.Concat(cfg.Resources[1].Binaries, cfg.Resources[1].PendingBinaries); len(ev) != 1 ||
		!slices.Equal(ev[0].Events, []ebpf.EventType{ebpf.EventRead}) {
		t.Fatalf("lines = %+v", ev)
	}
}
