//go:build tools

package tools

import (
	_ "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/lipgloss"
	_ "github.com/cilium/ebpf"
	_ "github.com/cilium/ebpf/asm"
	_ "github.com/sirupsen/logrus"
)
