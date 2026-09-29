package daemon

import (
	"errors"
	"path/filepath"
	"slices"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// The closure vets a lib's symlink-resolved target but hands back the unresolved path, which
// SetTrusted re-stats later: re-pointing a symlink in a user-writable RPATH dir in between makes a
// user-owned inode TRUSTED_LIB. Anything the closure would vouch for the kernel already trusts.
func TestBuildTrustedSetExcludesLibraryClosure(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "app")
	lib := filepath.Join(dir, "libpoc.so")
	refusedLib := filepath.Join(dir, "libother.so")
	orig := libraryClosure
	defer func() { libraryClosure = orig }()
	libraryClosure = func(string) ([]string, map[string]error, error) {
		return []string{lib}, map[string]error{refusedLib: errors.New("not root-owned (uid 1000)")}, nil
	}
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{{
		Path:      "/home/tester/.ssh",
		Binaries:  []daemonconfig.BinaryRule{{Path: bin}},
		AllowLibs: []string{"/opt/reviewed/libok.so"},
	}}}

	_, libs, _, rejected := buildTrustedSet(cfg)

	if slices.Contains(libs, lib) {
		t.Fatalf("closure member %s reached the trusted set; SetTrusted re-stats it after the "+
			"ownership check: %v", lib, libs)
	}
	if !slices.Contains(libs, "/opt/reviewed/libok.so") {
		t.Errorf("allow_lib entry missing from the trusted set: %v", libs)
	}
	if rejected[refusedLib] == nil || !slices.Contains(rejected[refusedLib].bins, bin) {
		t.Errorf("closure refusals must still reach warnUntrustedLibs: %v", rejected)
	}
}
