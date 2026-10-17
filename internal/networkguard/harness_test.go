// harness_test.go holds network guard checks that load BPF (root). Like internal/tui's
// TestEditorHarness, the integration suite compiles this package's tests into a static binary
// (netguardTestAmd64Bin) and runs one subtest at a time in a container with
// APPLISTENER_NETGUARD_HARNESS=1; a plain `go test` skips them.
package networkguard

import (
	"errors"
	"os"
	"testing"

	cilium "github.com/cilium/ebpf"
	"github.com/stretchr/testify/suite"
)

type netguardHarness struct {
	suite.Suite
}

func TestNetGuardHarness(t *testing.T) {
	suite.Run(t, new(netguardHarness))
}

func (s *netguardHarness) SetupTest() {
	if os.Getenv("APPLISTENER_NETGUARD_HARNESS") == "" {
		s.T().Skip("APPLISTENER_NETGUARD_HARNESS not set — driven by the integration suite, not a plain `go test`")
	}
	s.Require().Zero(os.Geteuid(), "the network guard harness must run as root")
}

func (s *netguardHarness) whitelist(opts ...Option) (*NetGuard, error) {
	self, err := os.Executable()
	s.Require().NoError(err)
	entry, err := ComputeBinaryEntry(self)
	s.Require().NoError(err)
	return NewNetGuard(ModeWhitelist, []BinaryEntry{entry}, nil, false, true, opts...)
}

// No live process: pid_max is at most 2^22. A dead tgid's mark stays until a fork reuses it.
const deadTgid = uint32(0x7ffffff0)

// A guard replacing a running one keeps its code-suspect marks: a swap must not clear what /proc
// can no longer show.
func (s *netguardHarness) TestSwapInheritsSuspectMarks() {
	g1, err := s.whitelist()
	s.Require().NoError(err)
	defer g1.Stop()
	s.Require().NoError(g1.objs.TrustCodeSuspect.Put(deadTgid, uint8(reasonEnv)))

	g2, err := s.whitelist(WithSuspectsFrom(g1))
	s.Require().NoError(err)
	defer g2.Stop()
	var why uint8
	s.Require().NoError(g2.objs.TrustCodeSuspect.Lookup(deadTgid, &why), "the mark was not carried over")
	s.Require().Equal(uint8(reasonEnv), why)

	g3, err := s.whitelist()
	s.Require().NoError(err)
	defer g3.Stop()
	s.Require().True(errors.Is(g3.objs.TrustCodeSuspect.Lookup(deadTgid, &why), cilium.ErrKeyNotExist),
		"control: a guard started without WithSuspectsFrom has no mark")
}

func (s *netguardHarness) TestInheritFromStoppedGuardRefused() {
	g1, err := s.whitelist()
	s.Require().NoError(err)
	g1.Stop()
	g2, err := s.whitelist(WithSuspectsFrom(g1))
	if err == nil {
		g2.Stop()
	}
	s.Require().Error(err, "marks that can't be read must refuse the start, not drop them")
}
