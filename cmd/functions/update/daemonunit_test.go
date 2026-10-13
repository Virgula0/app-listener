package update

import (
	"bytes"
	"testing"

	inst "github.com/Virgula0/app-listener/internal/install"
)

func TestDaemonUnitCmdPrintsBundledUnit(t *testing.T) {
	want, err := inst.SampleContent("app-listener-daemon.service")
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	DaemonUnitCmd.SetOut(&out)
	defer DaemonUnitCmd.SetOut(nil)
	if err := DaemonUnitCmd.RunE(DaemonUnitCmd, nil); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out.Bytes(), want) {
		t.Fatalf("printed unit differs from the bundled one:\n%s", out.Bytes())
	}
	if !DaemonUnitCmd.Hidden {
		t.Fatal("daemon-unit must stay hidden")
	}
}
