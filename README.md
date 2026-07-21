# app-listener

A CLI/TUI tool that monitors file system operations (open, read, write, delete, rename, symlink, hardlink, mkdir) on a specified directory using eBPF.

## Requirements

- Linux kernel with eBPF support (CONFIG_BPF, CONFIG_KPROBES, CONFIG_DEBUG_INFO_BTF)
- clang + llvm (for BPF compilation)
- bpftool (optional, for regenerating vmlinux.h)
- Go 1.24+
- Root privileges (for eBPF program loading)

## Quick Start

```bash
# Build (generates BPF bindings then compiles Go binary)
make build

# Monitor a directory
sudo ./build/linux/app-listener monitor /tmp

# Monitor recursively with depth limit
sudo ./build/linux/app-listener monitor /var/log --recursive --depth 3
```

The binary is built at `build/linux/app-listener`. Run with `sudo` — eBPF requires root.

## Commands

| Command | Description |
|---------|-------------|
| `monitor <dir>` | Monitor a directory for file system events |

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-r, --recursive` | false | Monitor subdirectories recursively |
| `-d, --depth <n>` | 0 | Max directory depth (requires `--recursive`) |
| `-l, --log-level <level>` | info | Log level: debug, info, warn, error |

Exit the TUI with `q` or `Ctrl+C`.

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Generate BPF bindings + build Go binary |
| `make build-linux` | Build Go binary only (Linux amd64) |
| `make generate` | Regenerate BPF C → Go bindings |
| `make bpftool-headers` | Regenerate vmlinux.h from running kernel BTF |
| `make lint` | Run golangci-lint |
| `make clean` | Remove build artifacts |
| `make install-deps` | Download Go dependencies |
| `make run` | Quick run with `go run` |

## Docker

```bash
docker compose build
docker compose run --rm app-listener monitor /tmp
```

The Dockerfile does a multi-stage build. The runner image is Debian slim; eBPF requires the host kernel, so the container needs `--privileged` or at least access to `/sys/kernel/btf/` and `CAP_BPF`.

## Architecture

```
cmd/
  root.go              — cobra root command, --log-level flag
  functions/monitor/   — "monitor" subcommand, wires eBPF + TUI
  entity/              — shared flag types

internal/
  infrastructure/ebpf/
    bpf/monitor.bpf.c  — BPF C source (kprobes on do_sys_open, ksys_read, ksys_write, …)
    bpf/vmlinux.h      — generated CO-RE header
    monitor_bpf.go     — bpf2go generated Go bindings
    embeds/monitor_bpf.o — compiled BPF ELF (embedded)
    monitor.go         — loads BPF, attaches kprobes, reads ringbuf, filters
    event.go           — event struct types
    checker.go         — eBPF runtime availability check
  tui/
    model.go           — bubbletea TUI with viewport, colored events
```

### Key design decisions

- **kprobes on internal functions**: On hardened kernels (≥6.x with PUSH_AND_CLEAR_REGS), `__x64_sys_*` entry points have all GP registers zeroed, making `PT_REGS_PARM*` return 0. We hook `do_sys_open` / `do_sys_openat2` / `ksys_read` / `ksys_write` instead — regular C functions that receive real arguments in registers.
- **CO-RE**: Compiled against vmlinux.h for portability across kernel versions.
- **Embedded BPF**: The compiled BPF `.o` is embedded in the Go binary via `//go:embed` — no runtime compilation.
- **PID filter**: Events from the monitor's own process are discarded to avoid feedback loops.
