# app-listener

Monitor or guard file system operations (open, read, write, delete, rename, symlink, hardlink, mkdir, mmap) and network operations (TCP/UDP/DNS) using eBPF. The daemon protects critical directories (SSH keys, credentials, browser profiles, AI-agent tokens) against credential info-stealers and supply-chain attacks.

![demo.gif](./media/demo.gif)

## Contents

- [Why](#why)
- [Quick Start](#quick-start)
- [Compatibility](#compatibility)
- [How it works](#how-it-works)
- [Modes](#modes)
  - [monitor — observe](#monitor--observe)
  - [guard — block file access](#guard--block-file-access)
  - [network-monitor — watch network](#network-monitor--watch-network)
  - [network-guard — block network](#network-guard--block-network)
  - [daemon — fscrypt + whitelist lifecycle](#daemon--fscrypt--whitelist-lifecycle)
  - [install / uninstall / update / edit-protected](#install--uninstall--update--edit-protected)
- [Debug](#debug)
- [Makefile targets](#makefile-targets)
- [Docker](#docker)
- [Architecture](#architecture)
- [Key design decisions](#key-design-decisions)

## Why

A serious eBPF rewrite of [arch-supply-chain-hardening](https://github.com/Virgula0/arch-app-armor-hardening). Tested on **Arch Linux, ext4, amd64**, kernel `7.0.11-hardened2-1-hardened`. Other kernels, patches and filesystems may produce bugs or bypasses.

> This is a vibe-coding experiment (coded mainly with the free DeepSeek V4 Flash). Do not use it to protect production systems.

## Quick Start

```bash
# 1. Build (regenerates BPF bindings from the running kernel, then compiles)
make bpftool-headers
make build

# 2. Interactive installer (root): builds the binary, generates the fscrypt
#    master key, discovers critical directories, encrypts the selected ones
#    with backups, installs the systemd unit + pacman reload hook, enables
#    the daemon. Revert with `sudo ./build/linux/app-listener uninstall`.
sudo ./build/linux/app-listener install
```

Other modes in one line each:

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
| Kernel | 5.4+ (monitor), 5.10+ (guard/daemon, needs `CONFIG_BPF_LSM`). Pre-compiled BPF uses CO-RE — works across kernels as long as BTF is present (`/sys/kernel/btf/vmlinux`). |
| Architecture | `linux/amd64` (CI-tested); `linux/arm64` cross-compiled. Others: regenerate BPF with target clang. |
| Containers | Privileged, with `/sys/kernel/btf` + `CAP_BPF`/`CAP_SYS_ADMIN`. |
| Host tools | clang/LLVM + bpftool only to regenerate bindings; Go 1.24+; root for eBPF loading. |

Run `make check-compatibility` **before** building — it verifies kernel, BTF, the **BPF LSM** (`bpf` in `/sys/kernel/security/lsm`, mandatory for guard/daemon, absent by default on stock Ubuntu and cloud kernels), sysctls, root, and build tools.

> **Ubuntu / cloud VMs**: stock Ubuntu kernels build `CONFIG_BPF_LSM=y` but do not activate it. Add `lsm=landlock,lockdown,yama,integrity,apparmor,bpf` to the kernel cmdline (`/etc/default/grub.d/50-cloudimg-settings.cfg` on cloud images) and reboot, or guard modes attach but never deny.

## How it works

eBPF C programs, compiled to BPF bytecode and embedded in the Go binary (`//go:embed`), attach to:

- **kprobes** (`monitor`): `vfs_open/read/write`, `vfs_unlink/rename/symlink/link/mkdir/rmdir`, `security_mmap_file`, splice/sendfile/copy_file_range variants
- **LSM hooks** (`guard`, `daemon`): `file_open`, `file_permission`, `mmap_file`, `path_unlink/rename/link/mkdir` — the only kernel mechanism that can **deny** an operation
- **tracepoints/kretprobes** (`network-monitor`): connect, accept4, sendto/msg, recvfrom/msg, close, `inet_csk_accept`

Userspace decodes ring-buffer events, applies path/recursion/depth/type filters, and prints (monitor) or enforces (guard) policy.

**Binary identity is by exe inode** (`current->mm->exe_file->f_inode`), never by name: renaming a binary does not bypass policy, and `prctl(PR_SET_NAME)` comm-spoofing cannot fool it.

## Modes

### monitor — observe

Traces file operations under watched paths; nothing is blocked.

```bash
sudo ./build/linux/app-listener monitor -w /var/log --recursive --depth 3
sudo ./build/linux/app-listener monitor -w /path/to/file.txt
```

| Flag | Default | Description |
|------|---------|-------------|
| `-w, --watch <path>` | required | Path to monitor (repeatable) |
| `-r, --recursive` | `false` | Recurse into subdirectories |
| `-d, --depth <n>` | `0` | Max depth (needs `--recursive`; `0` = unlimited) |
| `-e, --events <list>` | all | `OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP` |
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
| `-e, --events <list>` | all | Event type filter |
| `--headless` | `false` | No TUI; log `GUARD\|` events to stderr |

**Exec-open attribution (whitelist mode)**: executing a binary is an OPEN performed by the *launcher*, not the binary itself — a shell wrapper runs through its interpreter (`/bin/sh`), which the whitelist deliberately excludes. Two discriminator sources cover the full exec chain, both keyed on the **binary being executed**: the target's own open carries `__FMODE_EXEC` in `file->f_flags` (set by `do_open_execat`; `in_execve` is set *after* the exec file is opened, so it cannot mark the open), and the kernel's load-time accesses of the exec file (`prepare_binprm`'s read, binfmt mappings, in-tree interpreter opens) happen while `in_execve` is set. Either way the accessed file itself must be whitelisted, so a whitelisted binary *inside* the guarded tree (e.g. Discord under `~/.config/discord`) can be launched from any shell or wrapper script, and its in-tree helpers (libraries, `.pak`, the exec'd `chrome-sandbox`) resolve under the same attribution. Once running, the process *is* the whitelisted binary, so all its reads/writes/renames are checked against the whitelist normally. The exec fd is never exposed to the launcher, so this grants no way to read guarded content. Separate helper binaries an app spawns must be whitelisted explicitly; blacklist mode always attributes to the launcher.

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
- **Per-binary event masks**: unlisted events are denied (EPERM); `READ`/`WRITE`/`MMAP` implicitly allow `OPEN`. Valid: `OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP`.
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

**install** (`sudo ./build/linux/app-listener install`) — TUI wizard, in safe order:

0. stops a running daemon (a daemon outside systemd is a fatal error) → 1. builds the binary if missing → 2. generates `/etc/app-listener/fscrypt.key` if missing (existing keys are kept) → 3. picks users to protect (root + real users, all preselected; no ssh-agent unit for root) → 4. probes the built-in catalog (`internal/install/catalog.go`: SSH, GPG, AI agents, browsers, VPNs, password stores…) per user plus system entries → 5. manual directories (non-existent path = fatal) → 6. embedded editor for the generated `daemon.conf` (`Ctrl+S` save, `Esc` abort) → 7. verifies the encryption state of **every** configured directory and tests each against the master key → 8. asks per not-yet-encrypted directory (already-encrypted ones are never asked): migrate to `<dir>.app_listener.backup`, recreate, encrypt in place, copy contents back with modes/owners/timestamps → 9. deploys systemd unit, pacman reload hook, per-user ssh-agent unit, binary and config; differing existing files show a unified diff; daemon ends enabled+active (changed config → `systemctl reload`, not restart) → 10. orphan cleanup of fscrypt metadata (only `app-listener-key-*` protector signatures) → 11. per-backup deletion prompt (kept by default).

Safety: an existing backup **aborts** the migration; an encrypted directory declared `need_encryption: false` is fatal; a failed policy application rolls back (backup removed, original restored); empty whitelists still deny everything (with a warning); the unit runs with `ProtectSystem=yes`, `PrivateTmp=yes`, `NoNewPrivileges=yes`.

Maintenance: `sudo app-listener install --restore-backups` restores backups (aborts while the daemon runs); `--delete-post-backups` deletes them; both are TUI with progress bars.

> **Prerequisite**: each filesystem must be fscrypt-initialized (`sudo fscrypt setup --all-users`) and support encryption (ext4: `sudo tune2fs -O encrypt <dev>`; f2fs: `sudo fsck.f2fs -O encrypt <dev>`). The installer verifies both before asking anything.

**uninstall** (`sudo ./build/linux/app-listener uninstall`) — refuses while the daemon runs; re-scans the catalog (never trusts the config); tests every encrypted directory against the key (mismatch = fatal); asks per directory whether to decrypt in place (default **no**; decrypted via a temp sibling, never half-decrypted); removes orphan fscrypt metadata; reverts systemd unit/hook/binary/config; deletes the master key **only with `--delete-key`** (otherwise every still-encrypted directory becomes permanently unlockable-by-nobody).

**update** (`sudo app-listener update [--yes]`) — self-update from the latest signed `pre-YYYYMMDD-<sha>` GitHub pre-release. Flow: version check (up-to-date exits) → download binary + sha256 + Ed25519 signature → verify signature against the **embedded** public key, checksum, and GitHub's asset digest (any failure aborts before writing) → changelog viewer + confirmation (`--yes` skips both) → atomic replace of the running binary, symlink ensured, daemon restarted.

**edit-protected** (`sudo ./build/linux/app-listener edit-protected`) — edit one fscrypt-encrypted catalog directory in a two-pane vim-style editor without touching the install. Refuses while the daemon runs (the vault is what it guards); re-scans the catalog; the chosen vault must unlock with the master key; only **one** directory is unlocked at a time; keybindings include `a`/`A` new file/dir, `d` delete (confirmed), `m` chmod, `c` chown, `Ctrl+S` save (atomic write, preserves mode/owner); refuses binaries (NUL scan), files > 2 MiB and symlinks; on exit the vault is re-locked with the daemon's own two-pass teardown.

## Debug

### Daemon status and logs

```bash
systemctl is-enabled app-listener-daemon    # expect: enabled
systemctl is-active  app-listener-daemon    # expect: active
journalctl -u app-listener-daemon -f        # follow live (errors are here)
sudo journalctl -u app-listener-daemon -f | grep -i denied   # only guard decisions
```

### Reload after a package upgrade

The whitelist is matched by inode, so `pacman -Syu` replacing a whitelisted binary locks it out until reload:

```bash
sudo systemctl reload app-listener-daemon   # SIGHUP: recomputes identities atomically
sudo systemctl kill -s HUP app-listener-daemon   # same (done automatically by the pacman hook)
```

### Test the guard manually

```bash
sudo systemctl stop app-listener-daemon
sudo /usr/local/sbin/app-listener daemon --headless -l debug
# in another terminal: ssh-add -l / ssh -T git@github.com → must WORK
#                      cat /home/alice/.ssh/id_ed25519 → must be DENIED (DAEMON [DENIED])
```

### fscrypt: check, unlock, lock

```bash
sudo fscrypt status /home/alice/.ssh     # "Encrypted" / "Not encrypted"
sudo fscrypt unlock /home/alice/.ssh --key=/etc/app-listener/fscrypt.key
sudo fscrypt lock /home/alice/.ssh       # daemon does this automatically at shutdown
```

A wrong/old key fails immediately with "invalid wrapping key" — the daemon never silently generates a new one. A `need_encryption: true` resource without a policy makes the daemon refuse to start.

### Migration backups and master key

```bash
ls -ld /home/alice/.ssh.app_listener.backup                     # backups are kept by default
sudo rm -rf /home/alice/.ssh.app_listener.backup                # safe once the daemon is verified active
sudo mv /home/alice/.ssh.app_listener.backup /home/alice/.ssh   # disaster restore (stop daemon first)
sudo ls -l /etc/app-listener/fscrypt.key                        # must be 0600
sudo /usr/local/sbin/app-listener daemon --genkey               # regenerate? only with confirmation
```

### ssh-agent (per user)

```bash
systemctl --user status ssh-agent    # run as the user, not root
systemctl --user start ssh-agent     # or just relogin
```

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Regenerate BPF bindings (+ `vmlinux.h` from running kernel BTF) + build |
| `make build-linux` | Build the Go binary only |
| `make generate` / `make generate-{monitor,guard,networkmonitor,networkguard}` | Regenerate BPF bindings |
| `make bpftool-headers` | Regenerate shared `internal/bpf/vmlinux.h` |
| `make lint` / `make test` | golangci-lint / unit tests |
| `make test-integration` | Build exploit binaries + Docker integration suite |
| `make check-compatibility` | Verify the host before building/installing |
| `make clean` / `make install-deps` / `make run` | Artifacts / deps / `go run` |

## Docker

```bash
docker compose build
docker compose run --rm app-listener monitor -w /tmp
```

Multi-stage build; runner is Debian slim. eBPF needs the host kernel — run privileged or with `/sys/kernel/btf/` + `CAP_BPF`.

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

- **VFS kprobes for monitoring**: `vfs_open/read/write` catch all I/O regardless of syscall path (io_uring, splice, sendfile, mmap); they survive hardened kernels where `__x64_sys_*` probes are unreliable.
- **LSM hooks for guard**: the only kernel mechanism that can *deny*; guard reverts to blocking when the policy denies.
- **Binary identity via exe inode**: read in BPF from `current->mm->exe_file->f_inode` — renaming/symlinking a binary and comm-spoofing cannot bypass policy. Exec-opens attribute to the binary being executed (whitelist mode) so in-tree whitelisted binaries work behind shell wrappers.
- **CO-RE + embedded BPF**: one `vmlinux.h` regenerated from the running kernel; `.o` files embedded — no runtime compilation, portable across kernels with BTF.
- **Per-binary event masks in the guard engine** (`guard_exe_events`): `READ,WRITE`-style daemon restrictions enforced in BPF; missing entry = all events, so plain guard mode is unaffected; `OPEN` is implied by `READ`/`WRITE`/`MMAP`.
- **Daemon reload without a protection gap**: new LSM programs attach before old ones detach; the kernel denies when *any* denies — protection is the intersection.
- **TOC-TOU-safe fscrypt teardown**: two-pass deprovision (plain + force-flush EBUSY retry) while guards still deny, hooks detached only after every vault is keyless — stronger than ssh-guard.
- **No `chattr`**: dropped for the LSM engine (finer granularity, no immutable/append race).
- **Taint tracking**: processes that touched guarded content are tracked; `process_vm_readv`, ptrace and `/proc/<pid>/mem` against them are blocked unless the caller is whitelisted; taint is inherited across fork/exec.
- **Master key is non-destructive**: `O_EXCL` creation, regeneration needs explicit confirmation (it invalidates every provisioned directory).
