# app-listener

Monitor or guard file system and network operations with eBPF — the daemon protects critical directories (SSH keys, credentials, browser profiles, AI-agent tokens) against credential info-stealers and supply-chain attacks.

![demo.gif](./media/demo.gif)

> [!IMPORTANT]
> **Verify before you install.** The only way to be sure that **100% of the operations will work on your current system** is to run the integration suite:
>
> ```bash
> make test-integration
> ```
>
> It exercises every mode and bypass vector inside Docker containers. **Rootful Docker is required** — the suite loads real eBPF programs, which rootless or remote Docker daemons cannot run. All tests must pass.

## Contents

- [Why](#why)
- [Quick Start](#quick-start)
- [Compatibility](#compatibility)
- [How it works](#how-it-works)
- [Modes](#modes)
- [Debug](#debug)
- [Makefile targets](#makefile-targets)
- [Docker](#docker)
- [Architecture](#architecture)
- [Key design decisions](#key-design-decisions)

## Why

A serious eBPF rewrite of [arch-supply-chain-hardening](https://github.com/Virgula0/arch-app-armor-hardening). Tested on **Arch Linux, ext4, amd64**, kernel `7.0.11-hardened2-1-hardened` — other kernels, patches and filesystems may produce bugs or bypasses (exactly what `make test-integration` checks on your machine).

> This is a vibe-coding experiment (coded mainly with the free DeepSeek V4 Flash). Do not use it to protect production systems.

## Quick Start

```bash
# 1. Build (regenerates BPF bindings from the running kernel, then compiles)
#    Requires a ROOTFUL Docker daemon: the build runs in an isolated container
#    with the host's BTF vmlinux mounted, and outputs build/linux/app-listener
#    owned by your user. No Docker? Install clang/LLVM, bpftool, Go 1.26+ and
#    GCC, then run `make build-host` instead.
make build

# 2. Interactive installer (root): builds the binary, generates the fscrypt
#    master key, discovers critical directories, encrypts the selected ones
#    with backups, installs the systemd unit + pacman reload hook, enables
#    the daemon. Revert with `sudo ./build/linux/app-listener uninstall`.
sudo ./build/linux/app-listener install
```

One line per mode:

```bash
sudo ./build/linux/app-listener monitor -w /tmp                                  # observe file ops
sudo ./build/linux/app-listener guard /tmp -w /usr/bin/cat                       # block all but cat
sudo ./build/linux/app-listener network-monitor /usr/bin/bash                    # watch bash network ops
sudo ./build/linux/app-listener network-guard -w /usr/lib/firefox/firefox --auto-infra
sudo ./build/linux/app-listener daemon --genkey                                  # fscrypt master key
sudo ./build/linux/app-listener daemon --headless --blocked-only                 # protected daemon
sudo systemctl reload app-listener-daemon                                       # re-resolve whitelist inodes
sudo app-listener update --yes                                                   # self-update from GitHub
```

Exit any TUI with `q` or `Ctrl+C`.

## Compatibility

| Aspect | Requirement |
|--------|-------------|
| Kernel | 5.4+ (monitor), 5.10+ (guard/daemon, needs `CONFIG_BPF_LSM`). CO-RE BPF — needs BTF (`/sys/kernel/btf/vmlinux`). |
| Architecture | `linux/amd64` (CI-tested); `linux/arm64` cross-compiled. |
| Tools | `make build` needs a **rootful Docker** daemon. `make build-host` instead needs clang/LLVM + bpftool + Go 1.26+ + GCC installed locally. |

Run `make check-compatibility` **before** building — it verifies kernel, BTF, the **BPF LSM** (`bpf` in `/sys/kernel/security/lsm`, mandatory for guard/daemon, absent by default on stock Ubuntu/cloud kernels), sysctls, root and build tools.

> **Ubuntu / cloud VMs**: stock kernels build `CONFIG_BPF_LSM=y` but don't activate it. Add `lsm=landlock,lockdown,yama,integrity,apparmor,bpf` to the kernel cmdline and reboot, or guard modes attach but never deny.

## How it works

eBPF programs, compiled and embedded in the binary, attach at three levels:

- **kprobes** (`monitor`) — observe all I/O regardless of syscall path (io_uring, splice, sendfile, mmap) plus metadata operations (chmod/truncate/stat/access/readlink/mknod).
- **LSM hooks** (`guard`, `daemon`) — the only kernel mechanism that can **deny**; 22 hooks covering open, read/write permissions, mmap, unlink/rename/symlink/link/mkdir/rmdir/mknod, attributes (chmod/chown/utimes/truncate/setxattr), stat/access/readlink probes, mount and ptrace.
- **tracepoints/kretprobes** (`network-monitor`) — TCP/UDP/DNS operations.

**Binary identity is by exe inode**, never by name: renaming a binary or comm-spoofing cannot bypass policy. Keep whitelisted binaries **root-owned** — an attacker who can modify a binary's contents owns its identity anyway.

## Modes

### Browser TUI

`monitor`, `guard`, `network-monitor`, `network-guard`, and `daemon` can mirror their TUI into a browser. With `--serve` the normal local TUI keeps running on the terminal **and** the same event stream is shared, read-only, over WebSockets; quitting the local TUI (`q`, `ctrl+c` or SIGINT/SIGTERM) also stops the browser endpoint. A served TUI accepts one browser viewer at a time so browser dimensions map deterministically to its Bubble Tea viewport — the terminal and the browser size independently.

```bash
sudo ./build/linux/app-listener monitor -w /tmp --serve
sudo ./build/linux/app-listener guard /secret --serve=192.168.1.10:8080
sudo ./build/linux/app-listener daemon --serve --user admin --password secret
```

`--serve` binds to `127.0.0.1:9999`; use `--serve=host:port` to choose another address. `--user` and `--password` must be supplied together and protect both the page and WebSocket with HTTP Basic Auth. They cannot be used without `--serve`. Serving is mutually exclusive with `--headless` and `--gui`, and requires an interactive terminal (systemd-style services should stay on `--headless`).

The built-in server does not provide TLS. When binding outside loopback, place it behind a TLS reverse proxy because TUI data and Basic Auth credentials otherwise cross the network unencrypted.

### monitor — observe

Traces file operations under watched paths; nothing is blocked. In addition to data I/O (`OPEN/READ/WRITE/MMAP`) and tree changes it observes **metadata operations**: `ATTR` (chmod/chown/utimes/truncate/setxattr), `STAT` (stat/access/readlink) and `MKNOD` — the same path-based operations the guard denies.

```bash
sudo ./build/linux/app-listener monitor -w /var/log --recursive --depth 3
sudo ./build/linux/app-listener monitor -w /path/to/file.txt
```

| Flag | Default | Description |
|------|---------|-------------|
| `-w, --watch <path>` | required | Path to monitor (repeatable) |
| `-r, --recursive` | `false` | Recurse into subdirectories |
| `-d, --depth <n>` | `0` | Max depth (needs `--recursive`; `0` = unlimited) |
| `-e, --events <list>` | all | `OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP,ATTR,STAT,MKNOD` |
| `--serve[=<host:port>]` | disabled | Mirror the TUI into a browser and keep the local TUI (`127.0.0.1:9999` when no address is given) |
| `--user <name>` | — | HTTP Basic Auth username (requires non-empty `--password` and `--serve`) |
| `--password <password>` | — | HTTP Basic Auth password (requires non-empty `--user` and `--serve`) |
| `--headless` | `false` | No TUI; log to stderr |

### guard — block file access

Denies file operations on the guarded path by process identity, in **blacklist** or **whitelist** mode (default whitelist: omitted `-w` blocks everything).

```bash
sudo ./build/linux/app-listener guard /secret                 # block everything
sudo ./build/linux/app-listener guard /secret -w /usr/bin/cat # allow only cat
sudo ./build/linux/app-listener guard /secret -b /usr/bin/rm  # block only rm
```

| Flag | Default | Description |
|------|---------|-------------|
| `<path>` | required | File or directory to guard |
| `-w, --whitelist <binary>` | — | Binaries allowed (repeatable; mutually exclusive with `-b`) |
| `-b, --blacklist <binary>` | — | Binaries blocked (repeatable; mutually exclusive with `-w`) |
| `-r, --recursive` | `true` | Recurse into subdirectories |
| `-d, --depth <n>` | `0` | Max depth (`0` = unlimited) |
| `-e, --events <list>` | all | Event type filter (`OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP,ATTR,STAT,MKNOD`) |
| `--serve[=<host:port>]` | disabled | Mirror the TUI into a browser and keep the local TUI (`127.0.0.1:9999` when no address is given) |
| `--user <name>` | — | HTTP Basic Auth username (requires non-empty `--password` and `--serve`) |
| `--password <password>` | — | HTTP Basic Auth password (requires non-empty `--user` and `--serve`) |
| `--headless` | `false` | No TUI; log `GUARD\|` events to stderr |

**Exec-open attribution (whitelist mode)**: executing a binary is an OPEN performed by the *launcher* — a shell wrapper runs through its interpreter, which the whitelist deliberately excludes. Opens are attributed to the **binary being executed** instead, so whitelisted binaries *inside* the guarded tree (e.g. Discord under `~/.config/discord`) work from any shell or wrapper, and their in-tree helpers resolve under the same attribution. The exec fd is never exposed to the launcher, so this grants no way to read guarded content; separate helper binaries an app spawns must be whitelisted explicitly. Blacklist mode always attributes to the launcher.

### network-monitor — watch network

Traces network operations (TCP, UDP, DNS, …) of the listed binaries only.

```bash
sudo ./build/linux/app-listener network-monitor /usr/bin/curl /usr/bin/wget -e CONNECT,ACCEPT,DNS
sudo ./build/linux/app-listener network-monitor /usr/bin/tcpdump --headless
```

| Flag | Default | Description |
|------|---------|-------------|
| `<binary>` | required | Binaries to watch (positional, repeatable) |
| `-e, --events <list>` | all | `CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS` |
| `--serve[=<host:port>]` | disabled | Mirror the TUI into a browser and keep the local TUI (`127.0.0.1:9999` when no address is given) |
| `--user <name>` | — | HTTP Basic Auth username (requires non-empty `--password` and `--serve`) |
| `--password <password>` | — | HTTP Basic Auth password (requires non-empty `--user` and `--serve`) |
| `--headless` | `false` | No TUI; log `NETEVENT\|` events to stderr |

| Event | Meaning | Hook |
|-------|---------|------|
| `CONNECT` | Outbound connect | `sys_enter_connect` |
| `ACCEPT` | Inbound accepted (TCP) | `kretprobe/inet_csk_accept`, `sys_enter_accept[4]` |
| `SEND` / `RECV` | Data sent / received | `sys_enter_sendto`,`sendmsg` / `recvfrom`,`recvmsg` |
| `CLOSE` | Socket close | `sys_enter_close` |
| `DNS` | Query to port 53/853 | `sys_enter_connect`/`sendto` |

### network-guard — block network

Denies socket operations by binary identity, blacklist or whitelist mode, via LSM `socket_connect/bind/listen/sendmsg/recvmsg` hooks.

```bash
sudo ./build/linux/app-listener network-guard -b /usr/bin/curl /usr/bin/wget -e CONNECT,SEND
sudo ./build/linux/app-listener network-guard -w /usr/bin/vim --auto-infra   # only vim + system infra
```

| Flag | Default | Description |
|------|---------|-------------|
| `-b, --blacklist <binary>` | — | Block network ops for these binaries only (repeatable; exclusive with `-w`) |
| `-w, --whitelist <binary>` | — | Block network ops for **all** binaries except these (default deny) |
| `--auto-infra` | `false` | Auto-allowlist running infra daemons (resolved, NetworkManager, …) — otherwise DNS breaks for everyone |
| `--unsafe` | `false` | Also block AF_UNIX (X11, D-Bus, systemd) — may break the desktop |
| `--no-throttle` | `false` | Disable rate limiting (1 event/type/process per 250 ms default) |
| `-e, --events <list>` | all | `CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS,BIND,LISTEN` |
| `--serve[=<host:port>]` | disabled | Mirror the TUI into a browser and keep the local TUI (`127.0.0.1:9999` when no address is given) |
| `--user <name>` | — | HTTP Basic Auth username (requires non-empty `--password` and `--serve`) |
| `--password <password>` | — | HTTP Basic Auth password (requires non-empty `--user` and `--serve`) |
| `--headless` | `false` | No TUI; log `NETGUARD\|` events to stderr |

**Pick the real executable, not a wrapper**: identity is the exe inode, so whitelisting a `#!/bin/sh` wrapper (many distros ship `/usr/bin/firefox` as one) matches nothing. Resolve first:

```bash
ls -l /proc/$(pgrep -n firefox)/exe    # running process → real binary
readlink -f /usr/bin/firefox           # follows symlinks, NOT shell wrappers
# then whitelist the real path, e.g. /usr/lib/firefox/firefox
```

### daemon — fscrypt + whitelist lifecycle

Config-driven daemon protecting any number of directories with the guard's whitelist engine, plus fscrypt encryption lifecycle: resources are unlocked at startup and locked again on shutdown **while the guards remain attached** — never an unprotected window. Successor of [ssh-guard](https://github.com/Virgula0/arch-app-armor-hardening) (same philosophy, LSM instead of fanotify, no `chattr`).

```bash
sudo ./build/linux/app-listener daemon --genkey   # create the fscrypt master key
sudo ./build/linux/app-listener daemon            # default config: /etc/app-listener/daemon.conf → daemon-samples/daemon.conf
sudo ./build/linux/app-listener daemon --config /etc/ssh-guard/config --headless --blocked-only
```

| Flag | Default | Description |
|------|---------|-------------|
| `--config <path>` | — | Config resolution: flag → `/etc/app-listener/daemon.conf` → `daemon-samples/daemon.conf` |
| `--serve[=<host:port>]` | disabled | Mirror the TUI into a browser and keep the local TUI (`127.0.0.1:9999` when no address is given) |
| `--user <name>` | — | HTTP Basic Auth username (requires non-empty `--password` and `--serve`) |
| `--password <password>` | — | HTTP Basic Auth password (requires non-empty `--user` and `--serve`) |
| `--headless` | `false` | Log `DAEMON [DENIED]\|` to stderr (journald when a systemd service) |
| `--blocked-only` | `false` | Print only denied attempts (presentational) |
| `--genkey` | `false` | Generate `/etc/app-listener/fscrypt.key` and exit; regeneration asks for confirmation (it invalidates every encrypted directory) |

Config grammar (full template in `daemon-samples/daemon.conf`):

```text
[watch /home/alice/.ssh]          # one section per protected directory
need_encryption: true             # default true; false skips the fscrypt lifecycle
/usr/bin/ssh READ,WRITE           # whitelisted binary, restricted to these events
/usr/bin/ssh-agent                # bare path = all events allowed
```

- **Whitelist only, default deny**; identity by inode — renaming a binary does not grant access.
- **Per-binary event masks**: unlisted events are denied (EPERM); `READ`/`WRITE`/`MMAP` implicitly allow `OPEN`. Valid: `OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP, ATTR, STAT, MKNOD`.
- **Tolerance**: missing paths/binaries are skipped with a warning; malformed directives in a valid section fail fast.
- **SIGHUP reload** (`systemctl reload`, or the pacman `PostTransaction` hook): recomputes every binary's inode identity atomically — new guards attach before old ones detach, protection is never weaker; a malformed config keeps the previous one running.
- **fscrypt lifecycle**: `need_encryption: true` resources must already carry an fscrypt policy or the daemon refuses to start. Shutdown deprovisions keys in two passes (plain, then force-flush with an EBUSY retry loop) while guards still deny access; hooks detach only after every vault is keyless.

Minimal systemd unit:

```ini
[Unit]
Description=app-listener daemon (fscrypt + eBPF LSM whitelist)

[Service]
ExecStart=/usr/local/bin/app-listener daemon --headless
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

### install / uninstall / update / edit-protected

**install** — TUI wizard, in safe order: stops a running daemon → builds the binary → generates the fscrypt key (existing keys kept) → picks users to protect → probes a built-in catalog of critical directories (`internal/install/catalog.go`: SSH, GPG, AI agents, browsers, VPNs, password stores…) → encrypts selected directories (backup first, then verified against the master key) → deploys systemd unit, pacman reload hook, per-user ssh-agent unit, binary and config. An existing backup **aborts** the migration; a failed policy application rolls back; empty whitelists still deny everything; the unit runs with `ProtectSystem=yes`, `PrivateTmp=yes`, `NoNewPrivileges=yes`.

> **Prerequisite**: each filesystem must be fscrypt-initialized (`sudo fscrypt setup --all-users`) and support encryption (ext4: `sudo tune2fs -O encrypt <dev>`; f2fs: `sudo fsck.f2fs -O encrypt <dev>`). The installer verifies both before asking anything.

**uninstall** — refuses while the daemon runs; re-scans the catalog (never trusts the config); tests every encrypted directory against the key; decrypts in place by default (never half-decrypted); reverts systemd unit/hook/binary/config; deletes the master key **only with `--delete-key`**.

**update** — self-updates from the latest signed `pre-YYYYMMDD-<sha>` GitHub pre-release. The Ed25519 signature, checksum and GitHub's asset digest are all verified before anything is written; the running binary is replaced atomically and the daemon restarted.

**edit-protected** — edit one fscrypt-encrypted catalog directory in a two-pane vim-style editor without touching the install. Refuses while the daemon runs; one vault is unlocked at a time and re-locked with the daemon's own two-pass teardown on exit. `Ctrl+S` saves atomically (preserving mode/owner); binaries, symlinks and files > 2 MiB are refused.

## Debug

```bash
systemctl is-enabled app-listener-daemon    # expect: enabled
systemctl is-active  app-listener-daemon    # expect: active
journalctl -u app-listener-daemon -f        # follow live (errors are here)
sudo journalctl -u app-listener-daemon -f | grep -i denied   # only guard decisions
```

The whitelist is matched by inode, so `pacman -Syu` replacing a whitelisted binary locks it out until reload:

```bash
sudo systemctl reload app-listener-daemon   # SIGHUP: recomputes identities atomically
sudo systemctl kill -s HUP app-listener-daemon   # same (done automatically by the pacman hook)
```

Manual guard check:

```bash
sudo systemctl stop app-listener-daemon
sudo /usr/local/sbin/app-listener daemon --headless --verbose 3
# another terminal: ssh-add -l / ssh -T git@github.com → must WORK
#                   cat /home/alice/.ssh/id_ed25519 → must be DENIED
```

fscrypt lifecycle:

```bash
sudo fscrypt status /home/alice/.ssh     # "Encrypted" / "Not encrypted"
sudo fscrypt unlock /home/alice/.ssh --key=/etc/app-listener/fscrypt.key
sudo fscrypt lock /home/alice/.ssh       # daemon does this automatically at shutdown
```

A wrong/old key fails immediately with "invalid wrapping key" — the daemon never silently generates a new one. Backups live at `<dir>.app_listener.backup` (kept by default); `sudo app-listener install --restore-backups` restores them (aborts while the daemon runs). The ssh-agent service is per user: `systemctl --user start ssh-agent`.

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Dockerized build: isolated rootful container (host `/sys/kernel/btf/vmlinux` mounted read-only) regenerates BPF bindings + builds to `build/linux/app-listener`, owned by your user |
| `make build-host` | On-host build (needs clang/LLVM, bpftool, Go, GCC installed): regenerates BPF bindings + builds |
| `make build-image` | Build the `app-listener-builder` toolchain image (auto-run by `make build`) |
| `make build-linux` | Build the Go binary only |
| `make generate` / `make generate-{monitor,guard,networkmonitor,networkguard}` | Regenerate BPF bindings |
| `make bpftool-headers` | Regenerate shared `internal/bpf/vmlinux.h` |
| `make test-integration` | Builds exploit binaries + Docker integration suite (rootful Docker required) |
| `make check-compatibility` | Verify the host before building/installing |
| `make lint` / `make test` | golangci-lint / unit tests |
| `make clean` / `make install-deps` / `make run` | Artifacts / deps / `go run` |

## Docker

```bash
docker compose build
docker compose run --rm app-listener monitor -w /tmp
```

Multi-stage build; runner is Debian slim. eBPF needs the host kernel — run privileged or with `/sys/kernel/btf/` + `CAP_BPF`.

`make build` uses a dedicated, separately-built toolchain image (`docker/builder.Dockerfile`, built by `make build-image`); it does **not** touch this runtime multi-stage image.

## Architecture

```
cmd/functions/            monitor, guard, networkmonitor, networkguard, daemon, install, uninstall, update, edit-protected
cmd/common, cmd/printers  shared eBPF availability check, logo
internal/infrastructure/  FileEvent/NetEvent types, path resolution, eBPF availability check
internal/monitor/         bpf/monitor.bpf.c + kprobe loading, ringbuf reading, path/depth/type filters
internal/guard/           bpf/guard.bpf.c (LSM hooks) + policy maps (guard_inodes, guard_exe_actions, guard_exe_events, taint)
internal/networkmonitor/  tracepoint programs + exe-inode binary filter
internal/networkguard/    LSM socket hooks, blacklist/whitelist, --auto-infra discovery
internal/daemonconfig/    [watch <dir>] config parser
internal/fscrypt/         vault: master key, policy detection, unlock/lock lifecycle
internal/repository/      ports (GuardRepository, Vault…) + ErrKeyBusy/ErrKeyMissing sentinels
internal/usecase/         daemon orchestration (TOC-TOU-safe shutdown, atomic SIGHUP reload) + mode usecases
internal/tui/             bubbletea TUIs per mode, shared event-line formatting
integrationtests/         Docker suite: monitor/guard/network tests, C exploit binaries in exploits/
```

## Key design decisions

- **VFS kprobes for monitoring** — catch all I/O regardless of syscall path (io_uring, splice, sendfile, mmap); they survive hardened kernels where `__x64_sys_*` probes are unreliable. Metadata ops (chmod/truncate/stat/access/readlink/mknod) via `notify_change`, `vfs_setxattr/removexattr`, `vfs_getattr`, `vfs_readlink`, `do_faccessat`, `vfs_mknod`.
- **LSM hooks for guard** — the only kernel mechanism that can *deny*; guard reverts to blocking when the policy denies.
- **Metadata bypass coverage** — 22 hooks deny path-based ops that never create a struct file: truncate/chmod/chown/utimes/setxattr → `ATTR`, mknod → `MKNOD`, stat/access/readlink → `STAT`. `inode_permission` fires only on `access(2)`/`faccessat` probes (`MAY_ACCESS`), because plain read/write opens pass through it too and blocking them would steal identity from `file_open`/`path_truncate`/`inode_setxattr`.
- **Identity by exe inode** — renaming/symlinking a binary and comm-spoofing cannot bypass policy; exec-opens attribute to the binary being executed so whitelisted binaries work behind shell wrappers.
- **CO-RE + embedded BPF** — one `vmlinux.h` regenerated from the running kernel; `.o` files embedded — no runtime compilation, portable across kernels with BTF.
- **Per-binary event masks** — daemon.conf `READ,WRITE`-style restrictions enforced in BPF; a missing entry = all events (plain guard mode unaffected); `OPEN` implied by `READ`/`WRITE`/`MMAP`.
- **Reload without a protection gap** — new LSM programs attach before old ones detach; the kernel denies when *any* denies — protection is the intersection.
- **TOC-TOU-safe fscrypt teardown** — two-pass deprovision (plain + force-flush EBUSY retry) while guards still deny, hooks detached only after every vault is keyless — stronger than ssh-guard.
- **Taint tracking** — processes that touched guarded content are tracked; `process_vm_readv`, ptrace and `/proc/<pid>/mem` against them are blocked unless the caller is whitelisted; taint is inherited across fork/exec.
- **No `chattr`** — dropped for the LSM engine (finer granularity, no immutable/append race).
- **Non-destructive master key** — `O_EXCL` creation; regeneration needs explicit confirmation (it invalidates every provisioned directory).
