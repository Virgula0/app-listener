//go:build !gui

package gui

// Available reports whether this binary has the fyne GUI (`gui` build tag); false in the default
// build, since the daemon and other subcommands share the binary and must not link the X11/OpenGL
// toolkit. Callers of Run must check it first.
const Available = false

// Run is unreachable in the default build (the --gui entry checks Available and errors out); it
// exists only so the package compiles without the `gui` tag.
func Run(Options) error { return nil }
