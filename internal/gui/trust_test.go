//go:build gui

package gui

import (
	"slices"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/Virgula0/app-listener/internal/editbackend"
)

func trustShell(t *testing.T, link *fakeLink) (*shell, *trustPage) {
	t.Helper()
	a := test.NewApp()
	t.Cleanup(a.Quit)
	a.Settings().SetTheme(&appTheme{variant: theme.VariantDark})
	w := a.NewWindow("app-listener")
	s := newShell(a, w, &Options{Version: "dev"})
	s.syncBackground = true
	s.startBackend = func() (backendLink, error) { return link, nil }
	w.SetContent(s.view())
	w.Resize(fyne.NewSize(1280, 820))
	s.show(pageTrust)
	return s, s.trust
}

// confirmOpen presses the confirm button of the open confirmation dialog.
func confirmOpen(t *testing.T, s *shell) {
	t.Helper()
	for d := range s.dialogs {
		if c, ok := d.(interface{ Confirm() }); ok {
			c.Confirm()
			return
		}
	}
	t.Fatal("no confirmation dialog open")
}

func lastRequest(t *testing.T, link *fakeLink, op string) editbackend.Request {
	t.Helper()
	link.mu.Lock()
	defer link.mu.Unlock()
	for i := len(link.reqs) - 1; i >= 0; i-- {
		if link.reqs[i].Op == op {
			return link.reqs[i]
		}
	}
	t.Fatalf("no %s request in %d", op, len(link.reqs))
	return editbackend.Request{}
}

func TestTrustPageConfirmsTheCheckedRowsWithTheirContent(t *testing.T) {
	link := newFakeLink(true)
	refused := strings.Repeat("a", 64)
	link.trust = []editbackend.TrustRow{
		{Line: "/home/alice/.local/bin/gh", Resolved: "/home/alice/.local/bin/gh", SHA256: refused,
			Resources: "/home/alice/.config/gh"},
		{Line: "/opt/tool/bin/tool", Resolved: "/opt/tool/bin/tool", SHA256: strings.Repeat("b", 64),
			Resources: "/home/alice/.ssh", Skip: "changed again since the daemon refused it"},
	}
	s, p := trustShell(t, link)
	if !s.nav[pageTrust].badge {
		t.Fatal("no sidebar badge with refused binaries")
	}
	if p.checks[0] == nil || p.checks[1] != nil || p.trustBtn.Disabled() {
		t.Fatalf("checks = %v, trust disabled=%v", p.checks, p.trustBtn.Disabled())
	}
	p.checks[0].SetChecked(true)
	p.reload.SetChecked(false)
	p.confirmSelected()
	confirmOpen(t, s)
	req := lastRequest(t, link, editbackend.OpTrustConfirm)
	want := []editbackend.TrustItem{{Line: "/home/alice/.local/bin/gh", SHA256: refused}}
	if !slices.Equal(req.Trust, want) || !req.NoReload {
		t.Fatalf("trust_confirm = %+v no_reload=%v", req.Trust, req.NoReload)
	}
}

func TestTrustPageNoBadgeWhenNothingRefused(t *testing.T) {
	s, p := trustShell(t, newFakeLink(true))
	if s.nav[pageTrust].badge || !p.trustBtn.Disabled() {
		t.Fatalf("badge=%v trust enabled=%v", s.nav[pageTrust].badge, !p.trustBtn.Disabled())
	}
}

func TestTrustPageAddAuthenticatesBeforeInspecting(t *testing.T) {
	link := newFakeLink(true)
	s, p := trustShell(t, link)
	p.addLine("relative/tool")
	p.addLine("/home/alice/bin/tool")
	if !slices.Equal(p.added, []string{"/home/alice/bin/tool"}) {
		t.Fatalf("added = %v", p.added)
	}
	link.reqs = nil
	p.reviewAdded()
	if p.link == nil {
		t.Fatal("no backend held for the review")
	}
	if ops := link.ops(); len(ops) != 0 {
		t.Fatalf("sent before the password: %v", ops)
	}
	confirmOpen(t, s) // password
	confirmOpen(t, s) // review
	if ops := link.ops(); !slices.Equal(ops[:3], []string{editbackend.OpAuth, editbackend.OpTrustInspect,
		editbackend.OpTrustAdd}) {
		t.Fatalf("ops = %v", ops)
	}
	if p.link != nil || len(p.added) != 0 {
		t.Fatalf("after add: link=%v added=%v", p.link, p.added)
	}
}

func TestTrustPageAddRefusedWhileTheDaemonRunsWithoutItsSocket(t *testing.T) {
	_, p := trustShell(t, newFakeLink(false)) // the fake's daemon runs
	if !p.add.btn.Disabled() || p.addNote() == "" {
		t.Fatal("add by hand offered while only the password-less path applies to a running daemon")
	}
}
