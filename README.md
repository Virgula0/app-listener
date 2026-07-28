# app-listener

Monitor or guard file system operations (open, read, write, delete, rename, symlink, hardlink, mkdir, mmap) using eBPF.

## How it works

app-listener uses eBPF programs that attach to kernel hooks and emit events into a ring buffer. The userspace process reads the ring buffer, filters events against the watched paths, and either logs them (monitor mode) or denies the operation (guard mode).

**eBPF layer**: C programs compiled to BPF bytecode via clang/LLVM, embedded in the Go binary (`//go:embed`). They attach to:
- **kprobes** — `vfs_open`, `vfs_read`, `vfs_write`, `vfs_unlink`, `vfs_rename`, `vfs_symlink`, `vfs_link`, `vfs_mkdir`, `vfs_rmdir`, `security_mmap_file`, and splice/sendfile/copy_file_range variants (monitor mode)
- **LSM hooks** — `file_open`, `file_permission`, `mmap_file`, `path_unlink`, `path_rename`, `path_link`, `path_mkdir` (guard mode, blocks access)

**Userspace layer**: The Go binary opens the ring buffer, decodes raw events into typed `FileEvent` structs, applies path matching and recursion/depth filtering, and either prints to the TUI/log or enforces access policy.

## Requirements

- **Linux kernel** 5.x+ with `CONFIG_BPF`, `CONFIG_KPROBES`, `CONFIG_DEBUG_INFO_BTF` (monitor mode)
- **Guard mode** additionally requires `CONFIG_BPF_LSM` (available since 5.4+, LSM-enabled kernels)
- clang + LLVM (only needed to regenerate BPF bindings)
- bpftool (optional, for regenerating `vmlinux.h`)
- Go 1.24+
- **Root privileges** (required for eBPF program loading)

## Compatibility

| Aspect | Supported |
|--------|-----------|
| **Kernels** | 5.4+ (monitor), 5.10+ (guard with LSM). Pre-compiled BPF `.o` uses CO-RE — works across kernel versions without recompilation as long as BTF is available (`/sys/kernel/btf/vmlinux`). |
| **Architectures** | `linux/amd64` (primary, CI-tested). `linux/arm64` (cross-compiled, requires QEMU binfmt for testing). Other architectures need BPF regeneration with target-specific clang. |
| **Cgroups** | Works in privileged containers with `/sys/kernel/btf` bind-mounted + `CAP_BPF`, `CAP_SYS_ADMIN`. |

The embedded BPF `.o` is compiled for `x86_64`. For other architectures, run `make bpftool-headers` on the target kernel, then `make generate` to cross-compile the BPF programs.

## Modes

### monitor — observe events

Traces all file operations under watched paths and prints them to the TUI or headless log. No access is blocked.

```bash
sudo ./build/linux/app-listener monitor -w /tmp
sudo ./build/linux/app-listener monitor -w /var/log --recursive --depth 3
sudo ./build/linux/app-listener monitor -w /path/to/file.txt
```

| Flag | Default | Description |
|------|---------|-------------|
| `-w, --watch <path>` | required | Path to monitor (repeatable) |
| `-r, --recursive` | `false` | Monitor subdirectories recursively |
| `-d, --depth <n>` | `0` | Max directory depth (requires `--recursive`; `0` = unlimited) |
| `-e, --events <list>` | all | Event filter: comma-separated (`OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP`) |
| `--headless` | `false` | No TUI, log to stderr (for scripting/testing) |

### guard — block access

Intercepts file operations on guarded paths via LSM hooks and denies access based on process identity (binary hash). Supports blacklist and whitelist modes.

```bash
# Block all access to /secret (empty whitelist)
sudo ./build/linux/app-listener guard /secret

# Only allow cat
sudo ./build/linux/app-listener guard /secret -w /usr/bin/cat

# Block only rm
sudo ./build/linux/app-listener guard /secret -b /usr/bin/rm

# Guard with depth limit
sudo ./build/linux/app-listener guard /data --recursive --depth 2 -b /usr/bin/cat
```

| Flag | Default | Description |
|------|---------|-------------|
| `<path>` | required | Single file or directory to guard |
| `-b, --blacklist <binary>` | — | Binary paths to block (repeatable, mutually exclusive with `-w`) |
| `-w, --whitelist <binary>` | — | Binary paths to allow (repeatable, mutually exclusive with `-b`). When omitted, all binaries are blocked. |
| `-r, --recursive` | `true` | Guard subdirectories recursively |
| `-d, --depth <n>` | `0` | Max directory depth (`0` = unlimited) |
| `-e, --events <list>` | all | Event type filter (same as monitor) |
| `--headless` | `false` | No TUI, log GUARD\| events to stderr |

Guard mode uses binary **SHA256 hashes** to identify processes — not path names — so renaming a blacklisted binary does not bypass the policy.

## Quick Start

```bash
# Build (generates BPF bindings then compiles Go binary)
make build

# Monitor mode
sudo ./build/linux/app-listener monitor -w /tmp

# Guard mode — block everything
sudo ./build/linux/app-listener guard /tmp

# Guard mode — only cat is allowed
sudo ./build/linux/app-listener guard /tmp -w /usr/bin/cat
```

Exit the TUI with `q` or `Ctrl+C`.

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Generate BPF bindings + build Go binary |
| `make build-linux` | Build Go binary only (Linux amd64) |
| `make generate` | Regenerate all BPF C → Go bindings |
| `make generate-monitor` | Regenerate monitor BPF bindings |
| `make generate-guard` | Regenerate guard BPF bindings |
| `make bpftool-headers` | Regenerate `vmlinux.h` from running kernel BTF |
| `make lint` | Run golangci-lint |
| `make test` | Run unit tests (non-integration) |
| `make test-integration` | Build exploit binaries + run Docker integration tests |
| `make clean` | Remove build artifacts and generated BPF files |
| `make install-deps` | Download Go dependencies |
| `make run` | Quick run with `go run` |

## Docker

```bash
docker compose build
docker compose run --rm app-listener monitor -w /tmp
```

The Dockerfile does a multi-stage build. The runner image is Debian slim; eBPF requires the host kernel, so the container needs `--privileged` or access to `/sys/kernel/btf/` + `CAP_BPF`.

## Architecture

```
cmd/
  root.go              — cobra root command, --log-level flag
  functions/monitor/   — "monitor" subcommand, wires eBPF + TUI/headless
  functions/guard/     — "guard" subcommand, wires eBPF LSM + TUI/headless
  entity/              — shared flag types (Recursive, Depth, Headless, EventTypes)

internal/
  infrastructure/
    event.go           — FileEvent, BpfEvent types shared by both modes
    target.go          — Target type (path resolution)
    checker.go         — eBPF runtime availability check
    bpf/               — shared BPF helpers (if any)

  monitor/
    bpf/monitor.bpf.c  — BPF C source: kprobes on vfs_*, do_splice, security_mmap_file, …
    bpf/vmlinux.h      — generated CO-RE header (shared)
    embeds/monitor_bpf.o — compiled BPF ELF (embedded)
    monitor_bpf.go     — bpf2go generated Go bindings
    monitor.go         — loads BPF, attaches kprobes, reads ringbuf, filters path/depth/event types

  guard/
    bpf/guard.bpf.c    — BPF C source: LSM hooks + BPF maps for inode/policy tracking
    bpf/vmlinux.h      — generated CO-RE header
    embeds/guard_bpf.o — compiled BPF ELF (embedded)
    guard_bpf.go       — bpf2go generated Go bindings
    guard.go           — loads BPF, attaches LSM programs, manages inode/comm maps, enforces policy

  tui/
    model.go           — bubbletea TUI with viewport, colored events (monitor)
    guardmodel.go      — bubbletea TUI for guard mode

integrationtests/
  main_test.go         — suite setup, container helpers, log diff utilities
  exploit_test.go      — common exploit runner + individual TestExploit_* tests
  guard_test.go        — guard integration tests (blocking, blacklist, depth, runtime, bypasses)
  monitor_test.go      — monitor integration tests (event coverage, filters, depth, recursive)
  exploits/            — C exploit binaries (splice, sendfile, mmap, io_uring, …)
```

### Key design decisions

- **kprobes on VFS functions**: Hooking `vfs_open`/`vfs_read`/`vfs_write` catches all I/O regardless of syscall path (io_uring, splice, sendfile, mmap, …). On hardened kernels where register-zeroing makes `__x64_sys_*` probes unreliable, these VFS-level kprobes continue to work.
- **LSM hooks for guard**: `file_open`, `file_permission`, `mmap_file`, etc. are the only kernel mechanism that can **deny** an operation. Guard mode reverts to LSM blocking when the policy denies access.
- **CO-RE**: Compiled against `vmlinux.h` for portability across kernel versions without per-kernel recompilation.
- **Embedded BPF**: The compiled BPF `.o` is embedded in the Go binary — no runtime compilation or external dependencies.
- **PID filter**: Events from the monitor's own process are discarded to avoid feedback loops.
- **Binary identity via SHA256**: Guard identifies processes by cryptographic hash of the executable, not by path — renaming a binary does not bypass policy.
