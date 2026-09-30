// harness_test.go holds the editor checks that need root (chown to another uid). Like
// internal/fscrypt's TestFscryptHarness, the integration suite compiles this package's tests into
// a static binary (tuiTestAmd64Bin, integrationtests/main_test.go) and runs one subtest at a time
// in a container with APPLISTENER_TUI_HARNESS=1; a plain `go test` skips them.
package tui

import (
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"
)

type editorHarness struct {
	suite.Suite
}

func TestEditorHarness(t *testing.T) {
	suite.Run(t, new(editorHarness))
}

func (s *editorHarness) SetupTest() {
	if os.Getenv("APPLISTENER_TUI_HARNESS") == "" {
		s.T().Skip("APPLISTENER_TUI_HARNESS not set — driven by the integration suite, not a plain `go test`")
	}
	s.Require().Zero(os.Geteuid(), "the editor harness must run as root")
}

// chown of another user's file is the escalation: root hands a system file to the vault's user.
func (s *editorHarness) TestChownRefusesSymlinkedParent() {
	m, stored, victim := swappedParent(s.T(), "config")
	m.chownPath = stored
	m.chownUsers = []chownUser{{label: "nobody", uid: 65534, gid: 65534}}
	m.applyChown()
	info, err := os.Lstat(victim)
	s.Require().NoError(err)
	st, ok := info.Sys().(*syscall.Stat_t)
	s.Require().True(ok)
	s.Require().NotEqualf(uint32(65534), st.Uid,
		"chown reached outside the vault through a swapped parent dir (status %q)", m.status)
}

// Control: a chown inside the vault still applies.
func (s *editorHarness) TestChownInsideVault() {
	root := s.T().TempDir()
	target := root + "/config"
	s.Require().NoError(os.WriteFile(target, []byte("x"), 0o600))
	m := newFileEditModel(root)
	defer m.vault.close()
	m.chownPath = target
	m.chownUsers = []chownUser{{label: "nobody", uid: 65534, gid: 65534}}
	m.applyChown()
	info, err := os.Lstat(target)
	s.Require().NoError(err)
	st, ok := info.Sys().(*syscall.Stat_t)
	s.Require().True(ok)
	s.Require().Equalf(uint32(65534), st.Uid, "chown inside the vault must apply (status %q)", m.status)
}
