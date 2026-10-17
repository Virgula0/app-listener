package install

import (
	"errors"
	"testing"

	"github.com/Virgula0/app-listener/internal/wizard"
)

type cancelUI struct{ wizard.UI }

func (cancelUI) EditText(string, string, string) (string, error) { return "", wizard.ErrCanceled }

func TestEditTextMapsUICancel(t *testing.T) {
	defer wizard.Use(cancelUI{})()
	if _, err := EditText("Review", "daemon.conf", "x"); !errors.Is(err, ErrEditCanceled) {
		t.Fatalf("got %v, want ErrEditCanceled", err)
	}
}
