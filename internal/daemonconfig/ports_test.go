package daemonconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestParsePortSection(t *testing.T) {
	sh := filepath.Join(t.TempDir(), "srv")
	if err := os.WriteFile(sh, []byte("#!/bin/sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Parse([]byte("[port tcp/8332]\n" + sh + " MKNOD\n/usr/bin/true CONNECT\n/usr/bin/false\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Resources) != 1 {
		t.Fatalf("want 1 resource, got %d", len(cfg.Resources))
	}
	r := cfg.Resources[0]
	if r.Port != 8332 || r.Path != PortAnchor(8332) || r.NeedEncryption {
		t.Fatalf("got port %d path %q encrypted %v", r.Port, r.Path, r.NeedEncryption)
	}
	if len(r.Binaries) != 3 || r.Binaries[0].Events[0] != ebpf.EventMknod || r.Binaries[1].Events[0] != ebpf.EventConnect ||
		len(r.Binaries[2].Events) != 0 {
		t.Fatalf("binaries: %+v", r.Binaries)
	}
}

func TestParsePortSectionRefusals(t *testing.T) {
	for name, text := range map[string]string{
		"udp":        "[port udp/53]\n",
		"zero":       "[port tcp/0]\n",
		"too big":    "[port tcp/70000]\n",
		"not a port": "[port tcp/x]\n",
		"file event": "[port tcp/8332]\n/usr/bin/true READ\n",
		"encrypted":  "[port tcp/8332]\nneed_encryption: true\n",
		"watch path": "[port tcp/8332]\nwatch: /tmp\n",
		"allow_lib":  "[port tcp/8332]\nallow_lib /usr/lib/libc.so.6\n",
		"runtime":    "[port tcp/8332]\nruntime: true\n",
		"duplicate":  "[port tcp/8332]\n/usr/bin/true\n[port tcp/8332]\n/usr/bin/true\n",
	} {
		if _, err := Parse([]byte(text)); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestPortSectionTakesNoFilesystemPath(t *testing.T) {
	cfg, err := Parse([]byte("[port tcp/5037]\n/usr/bin/true CONNECT\n\n[watch " + t.TempDir() + "]\n/usr/bin/true\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Resources) != 2 || !strings.HasPrefix(cfg.Resources[0].Path, PortAnchorDir+"/") {
		t.Fatalf("resources: %+v", cfg.Resources)
	}
	if cfg.Resources[1].Port != 0 {
		t.Fatalf("a watch section got a port: %+v", cfg.Resources[1])
	}
}
