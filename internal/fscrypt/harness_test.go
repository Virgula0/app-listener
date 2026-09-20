// harness_test.go is a root-capable exec harness for the integration suite (like
// internal/guard/guard_test.go's TestGuardUnitTest / guardTestAmd64Bin): this package's real
// Vault.Encrypt/IsEncrypted (kernel fscrypt for directories, file-vault AEAD for regular files;
// never a reimplementation) needs root and, for directories, a fscrypt-capable filesystem, which a
// normal `go test` sandbox lacks. The integration suite compiles it into a static binary
// (fscryptTestAmd64Bin, integrationtests/main_test.go) and execs one subtest at a time in a
// privileged container.
//
// Every subtest is parameterized via APPLISTENER_HARNESS_PATH, driven with `-test.run
// TestFscryptHarness/<name>`.
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

// harnessPath returns APPLISTENER_HARNESS_PATH, or skips the subtest when unset (like
// guard_test.go's `os.Getuid() != 0` skip): plain `make test` / `go test ./internal/fscrypt/...`
// discovers this suite and would fail without it, instead of running only inside the integration
// suite's fscrypt-amd64 binary.
func (s *harnessSuite) harnessPath() string {
	p := os.Getenv("APPLISTENER_HARNESS_PATH")
	if p == "" {
		s.T().Skip("APPLISTENER_HARNESS_PATH not set — this harness is driven by the integration suite, not a plain `go test`")
	}
	return p
}

// The one-time install-time migration of APPLISTENER_HARNESS_PATH to encrypted-at-rest, as a real
// `install` does before the daemon attaches a guard: a directory gets a kernel fscrypt policy
// (prerequisites, the `encrypt` feature flag and per-filesystem metadata dir, applied via the
// google/fscrypt library as `fscrypt setup` would, so the container needs no extra package); a
// regular file is sealed into the file-vault format (filevault.go). Requires
// /etc/app-listener/fscrypt.key (MasterKeyFile, 32 bytes).
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
