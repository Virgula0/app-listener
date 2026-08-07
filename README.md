# app-listener

Monitor or guard file system operations (open, read, write, delete, rename, symlink, hardlink, mkdir, mmap) and network operations (TCP connect/accept/close, UDP send/recv, DNS) using eBPF.
The Daemon is particularly useful to protect most important directories on the filesystem against credential info-stealers largely used by supply-chain attacks.

![demo.gif](./media/demo.gif)

## Contents

- [Why](#why)
- [Quick Start](#quick-start)
- [How it works](#how-it-works)
- [Requirements](#requirements)
- [Compatibility](#compatibility)
- [Modes](#modes)
  - [monitor — observe events](#monitor--observe-events)
  - [guard — block access](#guard--block-access)
  - [network-monitor — watch network operations](#network-monitor--watch-network-operations)
  - [network-guard — block or allow network operations](#network-guard--block-or-allow-network-operations)
  - [daemon — ssh-guard style multi-directory whitelist with fscrypt lifecycle](#daemon--ssh-guard-style-multi-directory-whitelist-with-fscrypt-lifecycle)
  - [install — interactive daemon installer (root only)](#install--interactive-daemon-installer-root-only)
  - [uninstall — interactive daemon uninstaller (root only)](#uninstall--interactive-daemon-uninstaller-root-only)
- [Debug](#debug)
- [Makefile targets](#makefile-targets)
- [Docker](#docker)
- [Architecture](#architecture)
  - [Key design decisions](#key-design-decisions)

## Why

This can be seen as a more serious rewrite version of [arch-supply-chain-hardening](https://github.com/Virgula0/arch-app-armor-hardening), switching to `eBPF` and introducing a lot more features.

Tested on :

- `Arch Linux`
- Filesystem format: `ext4`
- For `amd64`
- On the kernel version `7.0.11-hardened2-1-hardened`

Different kernel versions, patches and file systems may produce undesired results, security bypasses or general bugs.
Before proceeding, it is important to know that this is a vibe-coding-like experiment and should not be used in any way to protect production-ready systems. It was mainly coded using the free `DeepSeek V4 Flash`.

## Quick Start

Want to protect your SSH keys, AI agent credentials, browser profiles and VPN
configs with the daemon? Use the interactive installer:

```bash
# 1. Build (generates BPF bindings then compiles the Go binary)
make build

# 2. Run the installer wizard (root only, TUI): it builds the binary if
#    missing, generates the fscrypt master key, discovers critical
#    directories, encrypts the selected ones with backups, installs the
#    systemd unit + pacman reload hook and enables the daemon
sudo ./build/linux/app-listener install

# When done, revert it with the uninstaller (root only, TUI): it refuses
# while the daemon is running, re-scans the catalog, asks which encrypted
# directories to permanently decrypt (default: no, with progress), removes
# the installed services/binary/config, and deletes the fscrypt master key
# ONLY with --delete-key (otherwise every still-encrypted directory can
# never be unlocked again)
sudo ./build/linux/app-listener uninstall
```

All other modes:

```bash
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

# Generate the daemon's fscrypt master key
sudo ./build/linux/app-listener daemon --genkey

# Daemon mode — protect encrypted directories, reload after package updates
sudo ./build/linux/app-listener daemon --headless --blocked-only
sudo systemctl reload app-listener-daemon   # or: kill -HUP $(cat /run/app-listener-daemon.pid)
```

Exit the TUI with `q` or `Ctrl+C`.

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

Guard mode identifies processes by the **exe inode** (read from `current->mm->exe_file->f_inode` in BPF) — not path names — so renaming a blacklisted binary does not bypass the policy.

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
| `--no-throttle` | `false` | Emit every network event without rate limiting (default: 1 event per type+process per 250ms, which protects the ring buffer from being flooded by noisy host processes in whitelist mode). |
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

### daemon — ssh-guard style multi-directory whitelist with fscrypt lifecycle

A config-file driven daemon that protects any number of directories with the guard's whitelist engine, and manages their fscrypt encryption lifecycle: encrypted resources are unlocked at startup, locked again on shutdown **while the guards remain attached**, so there is never an unprotected window. A successor of the [ssh-guard](https://github.com/Virgula0/arch-app-armor-hardening) daemon: same philosophy, LSM engine instead of fanotify, and no `chattr`/`exclude_chattr` (the LSM engine provides the same granularity).

```bash
# Generate the fscrypt master key (asks for confirmation if one already exists)
sudo ./build/linux/app-listener daemon --genkey

# Run with the default config (/etc/app-listener/daemon.conf, else daemon-samples/daemon.conf)
sudo ./build/linux/app-listener daemon

# Run with an explicit config, printing events to journald (systemd style)
sudo ./build/linux/app-listener daemon --config /etc/ssh-guard/config --headless --blocked-only
```

| Flag | Default | Description |
|------|---------|-------------|
| `--config <path>` | — | Config file. Resolved as: `--config` flag → `/etc/app-listener/daemon.conf` → `daemon-samples/daemon.conf` in the working directory |
| `--headless` | `false` | No TUI, print `DAEMON [DENIED]\|` lines to stderr (captured by journald when run as a systemd service). Lines carry a syslog priority marker (`<4>` warning / `<6>` info), so journalctl colors denied attempts yellow like ssh-guard's syslog output |
| `--blocked-only` | `false` | Print only blocked (denied) attempts; allowed events are suppressed. Presentational only — the guard behavior never changes |
| `--genkey` | `false` | Generate the fscrypt master key file and exit. If the key already exists, asks `Regenerate? [y/N]` on the terminal first — regenerating invalidates every fscrypt directory provisioned with the old key |

Config file grammar (see `daemon-samples/daemon.conf` for the full template):

```text
[watch /home/alice/.ssh]          # one section per protected directory
need_encryption: true             # default true; false skips the fscrypt lifecycle
/usr/bin/ssh READ,WRITE           # whitelisted binary, restricted to these events
/usr/bin/ssh-agent                # bare path = all events allowed
```

- **Whitelist only** (default deny): only the listed binaries may access the watched directory; identity is by binary inode, so renaming a binary to `ssh` does not grant access.
- **Per-binary event restrictions with blocking semantics**: unlisted events are denied with EPERM. Listing `READ`, `WRITE` or `MMAP` implicitly allows `OPEN` (a binary must open a file before reading it). Valid events: `OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP`.
- **Missing paths and binaries are skipped with a warning** (like ssh-guard); the directives of a skipped section are ignored with it. Malformed directives inside a valid section fail fast — a security configuration must not be silently misread.
- **SIGHUP reloads the configuration** (`systemctl reload`, `kill -HUP <pid>`): every binary's identity is recomputed by inode, so after a package upgrade replaced a whitelisted binary (`pacman -Syu`) a reload restores access for the new binary immediately. The new guards attach before the old ones detach — during the transition the kernel stacks both LSM programs and denies when *any* of them denies — so protection is never weaker than either configuration. A malformed config (or a failed reload) keeps the previous configuration running.
- **fscrypt lifecycle**: resources with `need_encryption: true` must already carry an fscrypt policy (run the fscrypt migration first) or the daemon refuses to start. On shutdown the keys are deprovisioned in two passes (first pass, then a force-flush pass with an EBUSY retry loop) while the guards still deny access, and only then are the LSM hooks detached — strictly stronger than ssh-guard's teardown, which removed its marks before the final lock pass.

Example systemd unit:

```ini
[Unit]
Description=app-listener daemon (fscrypt + eBPF LSM whitelist)

[Service]
ExecStart=/usr/local/bin/app-listener daemon --headless
ExecReload=/bin/kill -HUP $MAINPID
NotifyAccess=all
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

### install — interactive daemon installer (root only)

A TUI wizard that installs and configures the whole daemon stack above, in
the exact order the setup must happen to be safe:

```bash
sudo ./build/linux/app-listener install
```

The wizard (abort any step with `Esc`; completed steps stay completed):

0. **Daemon guard** — a running daemon is stopped before anything else: an
   active systemd unit is stopped (`systemctl stop app-listener-daemon`;
   the end of the install re-enables and re-starts it), a daemon process
   running outside systemd (e.g. a manual `app-listener daemon --headless`)
   is a **fatal error** — stop it manually first — and a daemon that is not
   running needs no action.
1. **Build** — if `build/linux/app-listener` is missing it runs `go build`
   in the repository (the repo must be the working directory); fails fast
   otherwise.
2. **Master key** — generates `/etc/app-listener/fscrypt.key` if missing;
   an existing key is kept (regenerating it would invalidate every
   already-encrypted directory).
3. **Users** — asks which users to protect: root plus the real users from
   `/etc/passwd` (UID ≥ 1000, with a home directory), all preselected
   (space/x toggles one, Ctrl+K toggles all/none, Enter continues). Root
   is protected like any other user (`/root/.ssh` is the most valuable
   target for a credential stealer), but no per-user ssh-agent unit is
   installed for root.
4. **Directory catalog** — probes the built-in catalog of critical
   directories (SSH, GPG, AI agents, browsers, VPNs, password stores, ...)
   **for every selected user** (each user's `.ssh`, `.gnupg`, ... is
   resolved against that user's home) plus the system-level entries
   (e.g. `/etc/wireguard`), and asks which discovered ones to protect. The
   catalog lives in `internal/install/catalog.go` — edit it to add or
   remove directories without touching the wizard.
5. **Manual directories** — additional directories can be added by hand;
   a path that does not exist is a fatal error (typos must not silently
   protect nothing).
6. **Editor** — shows the generated `daemon.conf` (whitelists are filtered
   to the binaries actually present, `need_encryption` mirrors the
   previous steps) in an embedded editor. Manual edits are preserved in
   every later step. `Ctrl+S` saves, `Esc` aborts the install.
7. **Key verification** — checks the encryption state of **every**
   configured directory, always: an already-encrypted directory whose
   section declares `need_encryption: false` is a fatal error (an
   encrypted directory must never be left unmanaged). Every
   already-encrypted directory is then tested against the master key; a
   directory that does not unlock is a fatal error (fixing it with the
   old key is impossible anyway, since the key is regenerated only with
   confirmation).
8. **Encryption** — asks per directory whether fscrypt encryption is
   required, but **only for directories that are not yet encrypted**:
   already-encrypted ones are never asked, never migrated and never get a
   backup. Migrated directories are moved to `<dir>.app_listener.backup`,
   recreated, encrypted in place with the master key, and the contents
   are copied back with permissions, owners and timestamps preserved.
9. **Deploy** — installs the systemd unit, the pacman reload hook
   (`PostTransaction` → `systemctl kill -s HUP`) and a per-user
   `ssh-agent` user unit (for the selected users) from the embedded
   `daemon-samples/`; copies the binary to `/usr/local/sbin/app-listener`
   and the config to `/etc/app-listener/daemon.conf` (0700/0600). Every
   file that already exists is compared with the bundled version:
   identical files are skipped, differing ones show the unified diff in
   the TUI and ask whether to overwrite (same for an existing config).
   Then the daemon is brought to the enabled-and-running state: a changed
   config is delivered to an already-running daemon with `systemctl
   reload` (SIGHUP) instead of a restart, a stopped daemon is started,
    and `is-enabled`/`is-active` are verified as `enabled`/`active`.
10. **Orphan cleanup** — removes fscrypt metadata (`/.fscrypt` policy and
    raw-key protector files) left behind by encrypted directories that no
    longer exist. The live set is scoped exactly like the catalog: every
    discoverable watch directory for **every** user, the system entries
    and the final config's resources — a directory that still exists
    (including one encrypted outside the config, like a manually protected
    `~/.ssh`) keeps its metadata; only policy/protector pairs carrying the
    installer's `app-listener-key-*` signature are ever deleted, so login
    protectors and other fscrypt setups are untouched.
11. **Backup cleanup** — asks per backup whether to delete the
    `.app_listener.backup` directories left by the migration (kept by
    default so nothing is lost after an interrupted migration).

The whole migration can be undone with `sudo app-listener install
--restore-backups`: the found `.app_listener.backup` directories are shown
in a TUI list (all preselected) and, after a single confirmation, each
encrypted directory is deleted and its backup moved back, restoring the
original unencrypted content. It **aborts** while the daemon is running —
stop it first with `systemctl stop app-listener-daemon`. The backups can
also be deleted standalone with `sudo app-listener install
--delete-post-backups` (same list + confirmation flow). Both options show
a TUI progress bar while running, and during the encryption step of a
regular install a TUI progress bar shows the copy progress.

Re-runs are safe: the installer detects a completed or interrupted
installation and resumes it (existing master key kept, already encrypted
directories verified, identical files skipped, backups never overwritten,
config reloaded instead of blindly rewritten).

Safety properties:

- A directory whose `.app_listener.backup` already exists **aborts** the
  migration: an older backup is never overwritten, and the migration can
  be resumed manually from the backup.
- The encryption state of every configured directory is always checked:
  a directory that is already encrypted but declared
  `need_encryption: false` is a **fatal error** — the installer exits
  instead of leaving an encrypted directory unmanaged.
- Already-encrypted directories are never asked about encryption and are
  never migrated, so no backup is created for them.
- If the fscrypt policy cannot be applied (unsupported filesystem, wrong
  kernel, ...), the migration rolls back: the backup is removed and the
  original directory restored untouched.
- Whitelists are filtered against the binaries actually present; a
  protected directory whose whitelist would end up empty still gets an
  explicit denial for everything, and the wizard logs a warning telling
  you to add binaries in the editor.
- The daemon unit runs with `ProtectSystem=yes`, `PrivateTmp=yes` and
  `NoNewPrivileges=yes`; the binary itself runs as root like its
  predecessor (eBPF loading requires it).

### uninstall — interactive daemon uninstaller (root only)

A TUI wizard that reverts everything `install` deployed, in the safe order:

```bash
sudo ./build/linux/app-listener uninstall
```

The wizard (abort any step with `Esc`; completed steps stay completed):

0. **Daemon guard** — fatally refuses while the daemon is running: an
   active systemd unit (`systemctl stop app-listener-daemon` first) or a
   daemon process running outside systemd (`kill <pid>`) both abort the
   uninstall. The installer stops the daemon because it re-enables it at
   the end; the uninstaller has nothing to re-enable, so it refuses
   instead.
1. **Catalog re-scan** — probes the built-in catalog for every local user
   (plus the system-level entries), exactly like the installer. The current
   `/etc/app-listener/daemon.conf` is deliberately **not** consulted: the
   uninstaller protects whatever is actually encrypted on the system now.
2. **Key verification** — every encrypted catalog directory is tested
   against the master key; a directory that does not unlock is a **fatal
   error**. Removing the daemon (and possibly the key with `--delete-key`)
   while a directory cannot be unlocked would lock it forever.
3. **Decryption** — asks per encrypted directory (default: **no**) whether
   the fscrypt protection must be permanently removed, then decrypts the
   confirmed ones in place with a TUI progress bar showing the copy
   progress. The contents are copied into a temporary plaintext sibling
   (`.app_listener.decrypt`), the encrypted directory is only removed after
   the copy completed, and the plaintext copy is renamed into place — a
   failure (or `Esc`) leaves the directory encrypted and untouched, never
   half-decrypted.
4. **Orphan cleanup** — removes the fscrypt policy/protector metadata left
   behind by the decrypted directories (same `app-listener-key-*`
   signature rule as the installer); still-encrypted directories keep their
   metadata.
5. **System revert** — after a confirmation, removes the systemd daemon
   unit, the pacman reload hook, the binary at `/usr/local/sbin` and the
   config at `/etc/app-listener/daemon.conf` (`systemctl disable` +
   `daemon-reload`). The per-user `ssh-agent` units installed by the
   installer are only reverted after a separate confirmation (default:
   **no**) and only when their content matches the bundled sample — a
   user's own modified unit is never touched.
6. **Master key** — the fscrypt key at `/etc/app-listener/fscrypt.key` is
   deleted **only** when `--delete-key` is passed:

   ```bash
   sudo ./build/linux/app-listener uninstall --delete-key
   ```

   The default keeps it: without the key, every directory that was left
   encrypted can never be unlocked again.

Migration backups (`.app_listener.backup`) are **not** touched by the
uninstaller — use `install --restore-backups` / `--delete-post-backups` for
those. The shared `/etc/fscrypt.conf` created by `fscrypt setup` is also
left in place.

## Debug

Everything lives in three places: the binary (`/usr/local/sbin/app-listener`),
the config (`/etc/app-listener/daemon.conf`) and the master key
(`/etc/app-listener/fscrypt.key`).

### Daemon status and logs

```bash
# Is it enabled across reboots and running right now?
systemctl is-enabled app-listener-daemon     # expect: enabled
systemctl is-active  app-listener-daemon     # expect: active

# Full status, recent errors and the journal tail
systemctl status app-listener-daemon
journalctl -u app-listener-daemon -e         # jump to the end (errors are here)
journalctl -u app-listener-daemon -b         # current boot only
journalctl -u app-listener-daemon -f         # follow live

# Only the guard decisions (whitelist hits and denials):
sudo journalctl -u app-listener-daemon -f | grep -i denied
```

### Reload after a package upgrade

The whitelist is matched by binary inode, so `pacman -Syu` replacing a
whitelisted binary locks it out until the config is reloaded:

```bash
sudo systemctl reload app-listener-daemon   # SIGHUP: recomputes identities atomically
sudo systemctl kill -s HUP app-listener-daemon   # same, done automatically by the pacman hook
```

A malformed config keeps the previous one running — the daemon only ever
rejects a reload, never applies a half-broken one.

### Test the guard manually

Stop the service and run the daemon in the foreground for instant feedback
(no systemd, events on stderr):

```bash
sudo systemctl stop app-listener-daemon
sudo /usr/local/sbin/app-listener daemon --headless -l debug

# In another terminal, this must WORK (whitelisted):
ssh-add -l
ssh -T git@github.com

# This must be DENIED (not whitelisted) and logged as DAEMON [DENIED]:
cat /home/alice/.ssh/id_ed25519
```

### fscrypt: check, unlock, lock

All watched directories must carry an fscrypt policy unlocked with the
master key in `/etc/app-listener/fscrypt.key`:

```bash
# Encryption state of a directory:
sudo fscrypt status /home/alice/.ssh      # "Encrypted" / "Not encrypted"

# Manually unlock with the daemon's master key:
sudo fscrypt unlock /home/alice/.ssh --key=/etc/app-listener/fscrypt.key

# Lock again (the daemon does this automatically at shutdown while the
# guards are still attached, so there is never an unprotected window):
sudo fscrypt lock /home/alice/.ssh

# Wrong/old key? The unlock fails immediately with an "invalid wrapping
# key" error — the daemon never silently generates a new key.
```

The daemon unlocks every `need_encryption: true` resource at startup and
locks it again on shutdown. A resource without a policy makes the daemon
refuse to start — check it with `fscrypt status` before complaining.

### Migration backups and master key

```bash
# Backups left by the installer (safe to delete once the daemon is
# verified active; kept by default on purpose):
ls -ld /home/alice/.ssh.app_listener.backup
sudo rm -rf /home/alice/.ssh.app_listener.backup

# Restore from a backup after a disaster (stop the daemon first):
sudo systemctl stop app-listener-daemon
sudo mv /home/alice/.ssh.app_listener.backup /home/alice/.ssh

# Master key location and permissions (must be 0600):
sudo ls -l /etc/app-listener/fscrypt.key

# Regenerating the key invalidates every provisioned policy — only do it
# with confirmation, then re-run the installer so directories are
# re-encrypted:
sudo /usr/local/sbin/app-listener daemon --genkey
```

### ssh-agent (per user)

```bash
systemctl --user status ssh-agent           # run as the user, not root
journalctl --user -u ssh-agent -b -f
systemctl --user start ssh-agent            # or just relogin
```

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
  functions/daemon/    — "daemon" subcommand: config file, --genkey, SIGHUP reload, headless/TUI
  common/              — shared eBPF availability check
  printers/            — logo printing

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

  daemonconfig/
    daemonconfig.go      — [watch <dir>] config parser (event lists, need_encryption, skip tolerance)

  fscrypt/
    fscrypt.go           — fscrypt vault: master key, policy detection, unlock/lock lifecycle

  repository/
    repository.go        — ports (GuardRepository, MonitorRepository, …)
    vault.go             — Vault port + ErrKeyBusy/ErrKeyMissing sentinels

  usecase/
    daemon.go            — daemon orchestration: start/stop, TOC-TOU-safe shutdown, atomic SIGHUP reload
    monitor.go, guard.go, networkmonitor.go, networkguard.go — other modes' orchestration

  tui/
    model.go           — bubbletea TUI with viewport, colored events (monitor)
    guardmodel.go      — bubbletea TUI for guard mode
    network_model.go   — bubbletea TUI for network-monitor mode
    daemon_model.go    — bubbletea TUI for daemon mode (per-resource status + merged events)
    eventline.go       — shared guard/daemon event line formatting

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
- **Binary identity via exe inode**: Guard and network-monitor identify processes by the exe file inode read directly from `current->mm->exe_file->f_inode` in BPF, not by path — renaming a binary does not bypass policy and comm-spoofing via `prctl(PR_SET_NAME)` cannot fool it.
- **LSM socket hooks for network-guard**: blocking is enforced by returning `-EPERM` from `socket_connect`, `socket_bind`, `socket_listen`, `socket_sendmsg` and `socket_recvmsg` — the only kernel mechanism that can deny a socket operation.
- **Safety boundary in whitelist mode**: by default whitelist mode only guards AF_INET/AF_INET6, leaving AF_UNIX (D-Bus, X11, systemd activation) untouched so desktop environments keep working. `--unsafe` extends guarding to all address families, at the risk of breaking the desktop.
- **`--auto-infra` keeps the system resolvable**: whitelist mode denies AF_INET/AF_INET6 for every process not explicitly allowed, including essential system daemons. Without an exception, `systemd-resolved` — which performs upstream DNS lookups on behalf of all processes (via the `resolve` NSS module) — would be blocked, breaking name resolution even for whitelisted apps. `--auto-infra` discovers running infra daemons (`systemd-resolved`, `NetworkManager`, `systemd-networkd`) via `/proc/*/exe` and allowlists them automatically.
- **Daemon reload without a protection gap (SIGHUP)**: reloading never detaches a guard whose replacement is not attached first. The new LSM programs are attached (and their ringbuf readers started) before the old ones are stopped; during the transition the kernel runs both programs and denies an access when any of them denies, so protection is strictly the intersection — never weaker than either configuration. A failed reload detaches the new guards and keeps the previous configuration running.
- **TOC-TOU-safe fscrypt teardown in the daemon**: shutdown locks every encrypted resource through a two-pass deprovision (first pass, then a force-flush pass with an EBUSY retry loop) **while the guards still deny access**, and detaches the LSM hooks only after every vault is keyless — no new open can pin an inode while the key is being revoked. Strictly stronger than ssh-guard, which removed its fanotify marks before the final lock pass.
- **Per-binary event masks in the guard engine**: the daemon's `READ,WRITE`-style restrictions are enforced in BPF via a per-exe-inode bitmask map (`guard_exe_events`); a missing entry means all events allowed, so plain guard mode is unaffected. `OPEN` is implicitly allowed whenever `READ`, `WRITE` or `MMAP` is listed — a binary cannot read without opening first.
- **No `chattr` in the daemon**: immutable/append flags and `exclude_chattr` were dropped entirely; the LSM whitelist engine replaces that mechanism with finer granularity (per-binary, per-event), avoiding the classic race between marking a file immutable and the attacker opening it.
- **Master key generation is non-destructive by default**: `daemon --genkey` creates the key with `O_EXCL` and never overwrites; when the key already exists the operator must explicitly confirm a regeneration (which atomically replaces the file) because it invalidates every fscrypt directory provisioned with the old key.
