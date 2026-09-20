//go:build !gui

package gui

import ebpf "github.com/Virgula0/app-listener/internal/infrastructure"

// Available reports whether this binary has the fyne GUI (`gui` build tag); false in the default
// build, since the daemon and other subcommands share the binary and must not link the X11/OpenGL
// toolkit. Callers of Run must check it first.
const Available = false

// Run is unreachable in the default build (monitor checks Available and
// errors out); it exists only so the package compiles without the `gui` tag.
func Run(events <-chan ebpf.FileEvent, paths []string, recursive bool, depth int) {}
