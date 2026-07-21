//go:build tools

package tools

import (
	_ "github.com/charmbracelet/bubbles"
	_ "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/lipgloss"
	_ "github.com/cilium/ebpf"
	_ "github.com/cilium/ebpf/asm"
	_ "github.com/cilium/ebpf/cmd/bpf2go"
	_ "github.com/cilium/ebpf/link"
	_ "github.com/cilium/ebpf/ringbuf"
	_ "github.com/cilium/ebpf/rlimit"
	_ "github.com/sirupsen/logrus"
	_ "github.com/spf13/cobra"
)
