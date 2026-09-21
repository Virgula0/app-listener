package daemon

import (
	"slices"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// A binary that is not in guard_trusted_files gets NO library allowlist at all: trust_mmap returns
// early unless the mapping process's exe carries TRUSTED_BINARY, so such a process can map any .so
// (LD_PRELOAD) while the per-resource guard still grants it full access to the secrets.
//
// PendingBinaries are whitelisted binaries that were unreadable at parse time. Config parsing
// always runs with the vaults still locked, so every whitelisted binary living INSIDE an encrypted
// watch tree lands there permanently (the Foundry catalog entry is exactly this shape: resource
// ~/.foundry, whitelist ~/.foundry/bin/*). The guard admits them post-unlock via
// ResolvePendingBinaries, so the trust set must cover them too — PendingLibs already is.
func TestBuildTrustedSetIncludesPendingBinaries(t *testing.T) {
	cfg := &daemonconfig.Config{
		Resources: []daemonconfig.Resource{{
			Path:            "/home/tester/.foundry",
			NeedEncryption:  true,
			Binaries:        []daemonconfig.BinaryRule{{Path: "/usr/bin/forge"}},
			PendingBinaries: []daemonconfig.BinaryRule{{Path: "/home/tester/.foundry/bin/cast"}},
		}},
	}

	binaries, _, _ := buildTrustedSet(cfg)

	if !slices.Contains(binaries, "/usr/bin/forge") {
		t.Fatalf("resolved binary missing from the trusted set: %v", binaries)
	}
	if !slices.Contains(binaries, "/home/tester/.foundry/bin/cast") {
		t.Fatalf("deferred binary missing from the trusted set, so it gets no library allowlist "+
			"and no write-protection: %v", binaries)
	}
}
