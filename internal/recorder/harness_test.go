// harness_test.go is the record-parity harness (integrationtests/record_parity_test.go): a real
// recording needs root and the monitor's kprobes. It records APPLISTENER_HARNESS_PATH
// (colon-separated targets) while `sh -c "$APPLISTENER_HARNESS_CMD"` runs, then prints what
// `config record` would propose: one MASK line per binary and the daemon.conf Plan.Render writes with
// every admissible binary approved and Masks on. APPLISTENER_HARNESS_BASE, when set, is the
// daemon.conf the targets are validated, classified and rendered against (a guarded target is
// extended; a binary it already whitelists prints COVERED instead of MASK).
package recorder

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

type harnessSuite struct {
	suite.Suite
}

func TestRecorderHarness(t *testing.T) {
	suite.Run(t, new(harnessSuite))
}

func (s *harnessSuite) TestRecord() {
	paths, script := os.Getenv("APPLISTENER_HARNESS_PATH"), os.Getenv("APPLISTENER_HARNESS_CMD")
	if paths == "" || script == "" {
		s.T().Skip("APPLISTENER_HARNESS_PATH/CMD not set — this harness is driven by the integration suite")
	}
	baseText, cfg := []byte(nil), &daemonconfig.Config{}
	if p := os.Getenv("APPLISTENER_HARNESS_BASE"); p != "" {
		var err error
		baseText, err = os.ReadFile(p)
		s.Require().NoError(err)
		cfg, err = daemonconfig.Parse(baseText)
		s.Require().NoError(err)
	}
	targets, err := ValidateTargets(strings.Split(paths, ":"), cfg, nil)
	s.Require().NoError(err)
	for _, t := range targets {
		if t.Existing != nil {
			fmt.Printf("EXISTING %s %s\n", t.Path, t.Existing.Section)
		}
	}
	rec, err := Start(targets)
	s.Require().NoError(err)
	defer rec.Stop()

	out, err := exec.Command("sh", "-c", script).CombinedOutput()
	fmt.Printf("CMD-BEGIN\n%sCMD-END\n", out)
	s.Require().NoError(err, "the recorded command failed")
	drain(rec)

	recording := Classify(rec.Stop(), cfg)
	sel := &Selection{Masks: true}
	for _, b := range recording.Binaries {
		if !b.Admissible() {
			fmt.Printf("REFUSED %s %s\n", b.Exe, b.Problem)
			continue
		}
		if b.Covered() {
			fmt.Printf("COVERED %s\n", b.Line)
			continue
		}
		fmt.Printf("MASK %s %s\n", b.Line, strings.Join(b.Mask, ","))
		sel.Binaries = append(sel.Binaries, b.Line)
	}
	plan, err := recording.Plan(sel)
	s.Require().NoError(err)
	text, err := plan.Render(baseText)
	s.Require().NoError(err)
	fmt.Printf("CONFIG-BEGIN\n%s\nCONFIG-END\n", text)
}

// drain waits until the ring buffer has stopped delivering: the command's last events may still be
// in flight when it exits.
func drain(rec *Recorder) {
	last, still := rec.Progress().Events, 0
	for deadline := time.Now().Add(15 * time.Second); time.Now().Before(deadline) && still < 5; {
		time.Sleep(200 * time.Millisecond)
		if n := rec.Progress().Events; n != last {
			last, still = n, 0
			continue
		}
		still++
	}
}
