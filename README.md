# app-listener

Monitor or guard file system operations (open, read, write, delete, rename, symlink, hardlink, mkdir, mmap) and network operations (TCP connect/accept/close, UDP send/recv, DNS) using eBPF.

## Why

This can be seen as a more serious rewrite version of [arch-supply-chain-hardening](https://github.com/Virgula0/arch-app-armor-hardening), switching to `eBPF` and introducing a lot more features.

Tested on :

- `Arch Linux`
- Filesystem format: `ext4`
- For `amd64`
- On the kernel version `7.0.11-hardened2-1-hardened`

Different kernel versions, patches and file systems may produce undesired results, security bypasses or general bugs.
Before proceeding, it is important to know that this is a vibe-coding-like experiment and should not be used in any way to protect production-ready systems. It was mainly coded using the free `DeepSeek V4 Flash`.

## How it works

app-listener uses eBPF programs that attach to kernel hooks and emit events into a ring buffer. The userspace process reads the ring buffer, filters events against the watched paths, and either logs them (monitor mode) or denies the operation (guard mode).

**eBPF layer**: C programs compiled to BPF bytecode via clang/LLVM, embedded in the Go binary (`//go:embed`). They attach to:
- **kprobes** — `vfs_open`, `vfs_read`, `vfs_write`, `vfs_unlink`, `vfs_rename`, `vfs_symlink`, `vfs_link`, `vfs_mkdir`, `vfs_rmdir`, `security_mmap_file`, and splice/sendfile/copy_file_range variants (monitor mode)
- **LSM hooks** — `file_open`, `file_permission`, `mmap_file`, `path_unlink`, `path_rename`, `path_link`, `path_mkdir` (guard mode, blocks access)
- **Tracepoints** — `syscalls/sys_enter_connect`, `sys_enter_accept4`, `sys_enter_sendto`, `sys_enter_recvfrom`, `sys_enter_sendmsg`, `sys_enter_recvmsg`, `sys_enter_close` (network-monitor mode)
- **kretprobe** — `inet_csk_accept` for TCP accept client address capture (network-monitor mode)

**Userspace layer**: The Go binary opens the ring buffer, decodes raw events into typed `FileEvent` structs, applies path matching and recursion/depth filtering, and either prints to the TUI/log or enforces access policy.

## Requirements

- **Linux kernel** 5.x+ with `CONFIG_BPF`, `CONFIG_KPROBES`, `CONFIG_DEBUG_INFO_BTF` (monitor/network-monitor mode)
- **Guard and network-guard modes** additionally require `CONFIG_BPF_LSM` (available since 5.4+, LSM-enabled kernels)
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

### network-monitor — watch network operations

Traces network operations (TCP, UDP, ICMP, DNS, etc.) of specific binaries using eBPF tracepoints. Only events from watched binaries are shown.

```bash
# Watch all network ops from bash
sudo ./build/linux/app-listener network-monitor /usr/bin/bash

# Watch specific binaries
sudo ./build/linux/app-listener network-monitor /usr/bin/curl /usr/bin/wget

# Filter by event type
sudo ./build/linux/app-listener network-monitor /usr/sbin/nginx -e CONNECT,ACCEPT,DNS

# Headless logging
sudo ./build/linux/app-listener network-monitor /usr/bin/tcpdump --headless
```

| Flag | Default | Description |
|------|---------|-------------|
| `<binary>` | required | Binary paths to watch (one or more positional args) |
| `-e, --events <list>` | all | Event filter: comma-separated (`CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS`) |
| `--headless` | `false` | No TUI, log NETEVENT\| events to stderr |

Binary identity is verified by **exe inode** (the BPF program reads `current->mm->exe_file->f_inode`), preventing comm-spoofing via `prctl(PR_SET_NAME)`.

#### Network event types

| Event | Description | BPF Hook |
|-------|-------------|----------|
| `CONNECT` | Outbound connect (TCP, UDP connect, raw) | `tracepoint/syscalls/sys_enter_connect` |
| `ACCEPT` | Inbound connection accepted (TCP) | `kretprobe/inet_csk_accept`, `tracepoint/syscalls/sys_enter_accept[4]` |
| `SEND` | Data sent (any protocol) | `tracepoint/syscalls/sys_enter_sendto`, `sys_enter_sendmsg` |
| `RECV` | Data received (any protocol) | `tracepoint/syscalls/sys_enter_recvfrom`, `sys_enter_recvmsg` |
| `CLOSE` | Socket close | `tracepoint/syscalls/sys_enter_close` (with socket fd detection) |
| `DNS` | DNS query (port 53/853) | Detected via `sys_enter_connect`/`sys_enter_sendto` to port 53 |

### network-guard — block or allow network operations

Blocks or allows network operations (TCP, UDP, DNS, etc.) of specific binaries by attaching eBPF LSM programs to the connect, bind, listen, sendmsg and recvmsg hooks. Supports blacklist and whitelist modes.

```bash
# Block all network ops from curl
sudo ./build/linux/app-listener network-guard -b /usr/bin/curl

# Block curl and wget, only CONNECT + SEND
sudo ./build/linux/app-listener network-guard -b /usr/bin/curl /usr/bin/wget -e CONNECT,SEND

# Whitelist: only vim may use the network, everything else is blocked
sudo ./build/linux/app-listener network-guard -w /usr/bin/vim

# Whitelist + keep system DNS/network daemons working
sudo ./build/linux/app-listener network-guard -w /usr/bin/firefox --auto-infra

# Headless
sudo ./build/linux/app-listener network-guard -b /usr/bin/wget --headless
```

| Flag | Default | Description |
|------|---------|-------------|
| `-b, --blacklist <binary>` | — | Block AF_INET/AF_INET6 operations for the listed binaries only; everything else is allowed. Repeatable, mutually exclusive with `-w`. |
| `-w, --whitelist <binary>` | — | Whitelist mode: block AF_INET/AF_INET6 for **all** binaries except the listed ones (default deny). Repeatable, mutually exclusive with `-b`. |
| `--unsafe` | `false` | Also block AF_UNIX sockets (used by X11, D-Bus, systemd activation). Only valid with `-w`, and may break desktop applications. Requires confirmation. |
| `--auto-infra` | `false` | Automatically allowlist running essential system network daemons (DNS resolver, network manager, etc.). Only valid with `-w`. |
| `-e, --events <list>` | all | Event filter: comma-separated (`CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS,BIND,LISTEN`) |
| `--headless` | `false` | No TUI, log NETGUARD\| events to stderr |

Binary identity is verified by the **exe inode** (`current->mm->exe_file->f_inode`), the same mechanism used by network-monitor — comm-spoofing via `prctl(PR_SET_NAME)` cannot bypass the policy.

#### Picking binaries: use the real executable, not a wrapper

Because identity is matched by exe inode, the path you pass must be the **actual binary** the process runs — not a shell-script wrapper. On many distributions `/usr/bin/firefox` is a `#!/bin/sh` script that just does `exec /usr/lib/firefox/firefox "$@"`; the running processes' exe is `/usr/lib/firefox/firefox`, so whitelisting (or blacklisting) `/usr/bin/firefox` matches nothing and the default action applies to all its traffic.

Before starting the guard, resolve the real binary first:

```bash
# While the application is running:
ls -l /proc/$(pgrep -n firefox)/exe          # → /usr/lib/firefox/firefox
readlink -f /usr/bin/firefox                 # follows symlink chains

# Then use the real path:
sudo ./build/linux/app-listener network-guard -w /usr/lib/firefox/firefox --auto-infra
```

Note that `readlink -f` resolves symlinks but **not** shell wrappers (a `#!/bin/sh` script is still a script). For wrapper scripts, read the `exec` target inside the file and use that path.

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

# Network monitor mode — watch bash network ops
sudo ./build/linux/app-listener network-monitor /usr/bin/bash

# Network-guard mode — block only curl
sudo ./build/linux/app-listener network-guard -b /usr/bin/curl

# Network-guard mode — whitelist firefox, keep DNS working
sudo ./build/linux/app-listener network-guard -w /usr/lib/firefox/firefox --auto-infra
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
| `make generate-networkmonitor` | Regenerate network-monitor BPF bindings |
| `make generate-networkguard` | Regenerate network-guard BPF bindings |
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
  functions/networkmonitor/ — "network-monitor" subcommand, wires eBPF tracepoints + TUI/headless
  functions/networkguard/  — "network-guard" subcommand, wires eBPF LSM hooks + TUI/headless
  entity/              — shared flag types (Recursive, Depth, Headless, EventTypes)

internal/
  infrastructure/
    event.go           — FileEvent, BpfEvent types shared by both modes
    networkevent.go    — NetEvent, NetBpfEvent types for network-monitor mode
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
    bpf/guard.bpf.c    — BPF C source: LSM hooks for file_open, file_permission, mmap, …
    bpf/vmlinux.h      — shared CO-RE header
    embeds/guard_bpf.o — compiled BPF ELF (embedded)
    guard_bpf.go       — bpf2go generated Go bindings
    guard.go           — loads BPF, attaches LSM programs, manages inode/comm maps, enforces policy

  networkmonitor/
    bpf/networkmonitor.bpf.c — BPF C source: tracepoints for connect, accept, sendto, recvfrom, …
    bpf/vmlinux.h            — shared CO-RE header
    embeds/networkmonitor_bpf.o — compiled BPF ELF (embedded)
    networkmonitor_bpf.go    — bpf2go generated Go bindings
    networkmonitor.go        — loads BPF, attaches tracepoints, manages exe inode map, binary filtering

  networkguard/
    bpf/networkguard.bpf.c   — BPF C source: LSM hooks for socket_connect, socket_bind, socket_listen, sendmsg, recvmsg
    bpf/vmlinux.h            — shared CO-RE header
    embeds/guardnet_bpf.o    — compiled BPF ELF (embedded)
    guardnet_bpf.go          — bpf2go generated Go bindings
    networkguard.go          — loads BPF, attaches LSM programs, enforces blacklist/whitelist policy, --auto-infra discovery

  tui/
    model.go           — bubbletea TUI with viewport, colored events (monitor)
    guardmodel.go      — bubbletea TUI for guard mode
    network_model.go   — bubbletea TUI for network-monitor mode

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
- **Binary identity via exe inode**: Network-monitor identifies processes by the exe file inode read directly from `current->mm->exe_file->f_inode` in BPF, preventing comm-spoofing.
- **LSM socket hooks for network-guard**: blocking is enforced by returning `-EPERM` from `socket_connect`, `socket_bind`, `socket_listen`, `socket_sendmsg` and `socket_recvmsg` — the only kernel mechanism that can deny a socket operation.
- **Safety boundary in whitelist mode**: by default whitelist mode only guards AF_INET/AF_INET6, leaving AF_UNIX (D-Bus, X11, systemd activation) untouched so desktop environments keep working. `--unsafe` extends guarding to all address families, at the risk of breaking the desktop.
- **`--auto-infra` keeps the system resolvable**: whitelist mode denies AF_INET/AF_INET6 for every process not explicitly allowed, including essential system daemons. Without an exception, `systemd-resolved` — which performs upstream DNS lookups on behalf of all processes (via the `resolve` NSS module) — would be blocked, breaking name resolution even for whitelisted apps. `--auto-infra` discovers running infra daemons (`systemd-resolved`, `NetworkManager`, `systemd-networkd`) via `/proc/*/exe` and allowlists them automatically.
