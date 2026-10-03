package daemon

import (
	"errors"
	"io/fs"
	"slices"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
)

func TestResolveInspectorsKeepsOnlyAdmitted(t *testing.T) {
	cfg := &daemonconfig.Config{Inspectors: []string{"/usr/lib/portal", "/home/u/portal", "/usr/libexec/gone"}}
	admit := func(p string) (guard.GuardInodeKey, error) {
		switch p {
		case "/usr/lib/portal":
			return guard.GuardInodeKey{Dev: 1, Ino: 7}, nil
		case "/usr/libexec/gone":
			return guard.GuardInodeKey{}, fs.ErrNotExist
		}
		return guard.GuardInodeKey{}, errors.New("is not root-owned")
	}
	got := resolveInspectors(cfg, admit)
	if len(got) != 1 || got[0].Path != "/usr/lib/portal" || got[0].Key.Ino != 7 {
		t.Fatalf("only the admitted inspector may be granted: %+v", got)
	}
}

func TestInspectorPathsOnlyRootPlaced(t *testing.T) {
	cfg := &daemonconfig.Config{Inspectors: []string{"/usr/lib/portal", "/home/u/portal"}}
	got := inspectorPaths(cfg, func(p string) bool { return p == "/usr/lib/portal" })
	if !slices.Equal(got, []string{"/usr/lib/portal"}) {
		t.Fatalf("a user's file must not become TRUSTED_BINARY (write-protected) as an inspector: %v", got)
	}
}

func TestCatalogWatchPlanWatchesInspectors(t *testing.T) {
	cfg := &daemonconfig.Config{Inspectors: []string{"/usr/lib/xdg-desktop-portal"}}
	for _, p := range catalogWatchPlan(cfg, nil) {
		if p.key() == newWatchPattern("/usr/lib/xdg-desktop-portal", false, false).key() {
			return
		}
	}
	t.Fatal("an upgraded inspector must trigger a re-sync, or the new version is refused until a reload")
}
