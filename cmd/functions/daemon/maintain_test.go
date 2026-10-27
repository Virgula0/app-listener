package daemon

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/install"
)

func TestMaintainedAppsNeedTheUsersReadOnlyLibDir(t *testing.T) {
	home := t.TempDir()
	steamRoot := filepath.Join(home, ".local/share/Steam")
	libDir := filepath.Join(steamRoot, "ubuntu12_32")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatal(err)
	}
	steam := filepath.Join(libDir, "steam")
	if err := os.WriteFile(steam, nil, 0o755); err != nil {
		t.Fatal(err)
	}
	user := install.User{Name: "u", UID: uint32(os.Getuid()), GID: uint32(os.Getgid()), Home: home} //nolint:gosec // test uid
	cfg := func(readOnly bool) *daemonconfig.Config {
		return &daemonconfig.Config{Resources: []daemonconfig.Resource{
			{Path: libDir, ReadOnly: readOnly, Binaries: []daemonconfig.BinaryRule{{Path: steam}}},
			{Path: filepath.Join(home, "elsewhere"), ReadOnly: readOnly, Binaries: []daemonconfig.BinaryRule{{Path: steam}}},
		}}
	}

	apps := maintainedApps(cfg(true), []install.User{user})
	if len(apps) != 1 {
		t.Fatalf("apps %+v", apps)
	}
	a := apps[0]
	if a.Script != filepath.Join(steamRoot, "steam.sh") || len(a.LibDirs) != 1 || a.LibDirs[0].Path != libDir ||
		a.owner != a.LibDirs[0] {
		t.Fatalf("app %+v", a)
	}
	if len(a.masks) != 1 || a.masks[0] != steam {
		t.Fatalf("the app's binaries below the script's directory are masked: %v", a.masks)
	}
	if apps := maintainedApps(cfg(false), []install.User{user}); len(apps) != 0 {
		t.Fatalf("a whitelist resource is no lib_dir: %+v", apps)
	}
	other := user
	other.UID++
	if apps := maintainedApps(cfg(true), []install.User{other}); len(apps) != 0 {
		t.Fatalf("another user's lib_dir: %+v", apps)
	}
}
