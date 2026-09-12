// harness_test.go is a root-capable exec harness for the integration
// suite, mirroring internal/guard/guard_test.go's TestGuardUnitTest /
// integrationtests' guardTestAmd64Bin pattern: this package's own real
// Vault.Encrypt/IsEncrypted (kernel fscrypt policies for directories, the
// userspace file-vault AEAD for regular files — never a reimplementation)
// needs root plus, for a directory target, an actual fscrypt-capable
// filesystem, neither of which a normal `go test` sandbox provides. The
// integration suite compiles this file into a static test binary (see
// integrationtests/main_test.go's fscryptTestAmd64Bin) and execs one
// subtest at a time inside a privileged container, exactly like
// guardTestAmd64Bin.
//
// Every subtest is parameterized entirely through the
// APPLISTENER_HARNESS_PATH env var so the integration suite can drive it
// with `-test.run TestFscryptHarness/<name>`.
package fscrypt

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/fscrypt/filesystem"
	"github.com/stretchr/testify/suite"
)

type harnessSuite struct {
	suite.Suite
}

func TestFscryptHarness(t *testing.T) {
	suite.Run(t, new(harnessSuite))
}

// harnessPath returns APPLISTENER_HARNESS_PATH, or skips the subtest when
// unset — mirroring internal/guard/guard_test.go's `os.Getuid() != 0` skip
// for its own root-only BPF subtests. This package's harness suite is
// discovered and run by plain `make test` / `go test ./internal/fscrypt/...`
// like any other test in the package; without this it would fail there
// instead of only running (with the var set) inside the integration suite's
// compiled fscrypt-amd64 binary.
func (s *harnessSuite) harnessPath() string {
	p := os.Getenv("APPLISTENER_HARNESS_PATH")
	if p == "" {
		s.T().Skip("APPLISTENER_HARNESS_PATH not set — this harness is driven by the integration suite, not a plain `go test`")
	}
	return p
}

// TestMigrate performs the one-time "install"-time migration of
// APPLISTENER_HARNESS_PATH from plaintext to encrypted-at-rest, exactly
// what a real `install` run does before the daemon ever attaches a guard:
// a directory gets a real kernel fscrypt policy (the filesystem
// prerequisites — the `encrypt` feature flag and the per-filesystem
// fscrypt metadata directory — are applied programmatically here via the
// google/fscrypt library directly, the same steps the `fscrypt setup` CLI
// performs, so the test container needs no extra package installed); a
// regular file gets sealed into the userspace file-vault format
// (filevault.go). Requires /etc/app-listener/fscrypt.key to already exist
// (MasterKeyFile, 32 bytes).
func (s *harnessSuite) TestMigrate() {
	path := s.harnessPath()

	if !isRegularFileTarget(path) {
		mnt, err := filesystem.FindMount(path)
		s.Require().NoError(err, "resolving the filesystem backing %s", path)
		s.Require().NoError(mnt.CheckSupport(),
			"filesystem backing %s is not marked for fscrypt encryption (mkfs -O encrypt missing?)", path)
		if setupErr := mnt.CheckSetup(nil); setupErr != nil {
			s.Require().NoError(mnt.Setup(filesystem.WorldWritable),
				"initializing fscrypt metadata on %s", mnt.Path)
		}
	}
	s.Require().NoError(EnsureSystemSetup())

	v := New()
	s.Require().NoError(v.Encrypt(path), "encrypting %s", path)
}

// TestIsEncrypted reports the current on-disk state of
// APPLISTENER_HARNESS_PATH on stdout ("ENCRYPTED" or "PLAINTEXT") without
// mutating anything — used by the integration tests to confirm a resource
// never spends any observable instant readable-in-plaintext outside the
// ordering the daemon's attach -> unlock -> populate contract promises.
func (s *harnessSuite) TestIsEncrypted() {
	path := s.harnessPath()
	v := New()
	encrypted, err := v.IsEncrypted(path)
	s.Require().NoError(err, "checking encryption state of %s", path)
	if encrypted {
		fmt.Println("ENCRYPTED")
	} else {
		fmt.Println("PLAINTEXT")
	}
}
