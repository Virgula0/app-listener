package guard

import (
	"debug/elf"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// loaderName: the dynamic loader's names (glibc ld-linux-*.so.N, ld-2.31.so, musl ld-musl-*.so.1,
// ppc64 ld64.so.N, ld.so) and musl's libc, which is its loader.
func loaderName(name string) bool {
	return (strings.HasPrefix(name, "ld-") && strings.Contains(name, ".so")) ||
		strings.HasPrefix(name, "ld.so") || strings.HasPrefix(name, "ld64.so") ||
		strings.HasPrefix(name, "libc.musl-")
}

// LoaderRefusal errors when path is the dynamic loader, by its name, the name it resolves to, or
// its ELF DT_SONAME (a renamed copy). `ld.so <program>` runs any program with the loader as its exe
// inode, so whitelisting it whitelists every program on the host. nil for anything else, including
// a path unreadable now (as MulticallRefusal).
func LoaderRefusal(path string) error {
	names := []string{filepath.Base(path)}
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		names = append(names, filepath.Base(resolved))
	}
	names = append(names, soNames(path)...)
	for _, n := range names {
		if loaderName(n) {
			return fmt.Errorf("%s is the dynamic loader (%s): `%s <program>` runs any program as this "+
				"executable, so whitelisting it admits every program; whitelist the program itself", path, n, path)
		}
	}
	return nil
}

// soNames reads path's DT_SONAME through the confined O_PATH handle.
func soNames(path string) []string {
	f, err := ebpf.OpenConfined(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	r, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return nil
	}
	defer r.Close()
	e, err := elf.NewFile(r)
	if err != nil {
		return nil
	}
	sonames, _ := e.DynString(elf.DT_SONAME)
	return sonames
}
