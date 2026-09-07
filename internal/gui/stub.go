//go:build !gui

package gui

import ebpf "github.com/Virgula0/app-listener/internal/infrastructure"

// Available reports whether this binary was built with the fyne desktop GUI
// (the `gui` build tag). It is false in the default build: the daemon and
// the other subcommands share the single app-listener binary and must not
// link the X11/OpenGL/image toolkit. Callers of Run must check this first.
const Available = false

// Run is unreachable in the default build (monitor checks Available and
// errors out); it exists only so the package compiles without the `gui` tag.
func Run(events <-chan ebpf.FileEvent, paths []string, recursive bool, depth int) {}
