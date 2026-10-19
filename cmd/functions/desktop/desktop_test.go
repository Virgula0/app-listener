package desktop

import (
	"testing"

	"github.com/Virgula0/app-listener/internal/gui"
)

func TestEveryInstallModeMapped(t *testing.T) {
	for m := gui.InstallWizard; m <= gui.InstallDeleteBackups; m++ {
		if _, ok := installModes[m]; !ok {
			t.Errorf("gui.InstallMode %d has no installer mode", m)
		}
	}
}
