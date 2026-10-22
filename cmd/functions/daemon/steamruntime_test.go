package daemon

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/install"
)

func TestKeptRuntimesNeedsTheUsersReadOnlyLibDir(t *testing.T) {
	home := t.TempDir()
	libDir := filepath.Join(home, ".local/share/Steam/ubuntu12_32")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatal(err)
	}
	steam := filepath.Join(libDir, "steam")
	user := install.User{Name: "u", UID: uint32(os.Getuid()), GID: uint32(os.Getgid()), Home: home} //nolint:gosec // test uid
	cfg := func(readOnly bool) *daemonconfig.Config {
		return &daemonconfig.Config{Resources: []daemonconfig.Resource{{Path: libDir, ReadOnly: readOnly,
			Binaries: []daemonconfig.BinaryRule{{Path: steam}}}}}
	}

	kept := keptRuntimes(cfg(true), []install.User{user})
	if len(kept) != 1 {
		t.Fatalf("kept %+v", kept)
	}
	if k := kept[0]; k.rt.Dir != libDir || k.rt.Name != "steam-runtime" || len(k.owners) != 1 || k.owners[0] != steam {
		t.Fatalf("kept %+v", k)
	}
	if kept := keptRuntimes(cfg(false), []install.User{user}); len(kept) != 0 {
		t.Fatalf("a whitelist resource is no lib_dir: %+v", kept)
	}
	other := user
	other.UID++
	if kept := keptRuntimes(cfg(true), []install.User{other}); len(kept) != 0 {
		t.Fatalf("another user's lib_dir: %+v", kept)
	}
}
