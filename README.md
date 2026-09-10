# app-listener

Monitor or guard file system and network operations with eBPF — the daemon protects critical directories (SSH keys, credentials, browser profiles, AI-agent tokens) against credential info-stealers and supply-chain attacks.

![demo.gif](./media/demo.gif)

## Contents

- [Compatibility](#compatibility)
- [Quick Start](#quick-start)
- [How it works](#how-it-works)
- [Modes](#modes)
  - [monitor](#monitor--observe)
  - [guard](#guard--block-file-access)
  - [network-monitor](#network-monitor--watch-network)
  - [network-guard](#network-guard--block-network)
  - [daemon](#daemon--fscrypt--whitelist-lifecycle)
  - [install.sh / install / uninstall / update / edit-protected](#installsh--install--uninstall--update--edit-protected)
- [Debug](#debug)
- [Makefile targets](#makefile-targets)
- [Docker](#docker)

> Vibe-coding experiment, coded mainly with the free DeepSeek V4 Flash and Claude Code. Not for production systems.

## Compatibility

Run `make check-compatibility` first — it performs every static check (kernel version, kernel `.config`, BTF, BPF-LSM activation, fscrypt prerequisites, BPF sysctls) and tells you plainly whether the host can run app-listener. The one-line installer below runs it for you and refuses to install when it fails. (`make build` runs the toolchain in Docker, so nothing has to be installed on the host to build it.)

| Distribution | Min. kernel | Support | Activation |
|---|---|---|---|
| **Arch Linux** | rolling (6.x) | ✅ **Fully supported** — all 23 LSM hooks attach | Add `bpf` to the `lsm=` list on your boot entry's kernel cmdline, reboot |
| **Ubuntu 24.04 LTS** | 6.8 | ✅ **Fully supported** | Append `lsm=landlock,lockdown,yama,integrity,apparmor,bpf` to `GRUB_CMDLINE_LINUX_DEFAULT`, `sudo update-grub`, reboot |
| **Ubuntu 22.04 LTS** | 5.15 (GA) | 🟡 **Partially supported** on the 5.15 GA kernel — a few later-kernel LSM hooks (notably `file_truncate`) are absent, so `ftruncate(2)` on a pre-opened fd is not denied (path `truncate(2)` still is); every other hook works. Install the HWE kernel (`linux-generic-hwe-22.04`, 6.x) for full support | same as 24.04 |

**Kernel floor:** `monitor` needs ≥ 5.8 (BPF ring buffer); `guard` / `network-guard` / `daemon` need ≥ 5.10 (BPF-LSM). Ubuntu 20.04 (kernel 5.4) is not supported. Only the two hooks `file_open` and `file_permission` are mandatory — every other LSM hook is best-effort: a kernel that lacks one logs a warning and keeps enforcing the rest.

**Stock Ubuntu and cloud images compile `CONFIG_BPF_LSM=y` but do not activate it** — without the cmdline change the LSM hooks attach but never deny. `linux/amd64` is the released target; `linux/arm64` cross-compiles but is untested.

## Quick Start

**1. Install the binary** (`check-compatibility` + signed release from GitHub → `/usr/local/sbin/app-listener`):

```bash
curl -fsSL https://raw.githubusercontent.com/Virgula0/app-listener/main/scripts/install.sh | sudo bash
```

Installs the latest **stable** release. For pre-release builds, append `-s -- --channel prerelease`.

<details><summary>Build from source instead</summary>

```bash
# Toolchain runs in a rootful Docker container (host BTF mounted); output at
# build/linux/app-listener. No Docker? clang/LLVM + bpftool + Go 1.26+ + GCC, then `make build-host`.
make build
sudo ./build/linux/app-listener install
```
</details>

**2. Protect directories with the daemon** — interactive, root, never automatic. Sets up fscrypt + whitelist + systemd units. Revert with `sudo app-listener uninstall`.

```bash
sudo app-listener install
```

One line per mode:

```bash
sudo app-listener monitor -w /tmp                                  # observe file ops
sudo app-listener guard /tmp -w /usr/bin/cat                       # block all but cat
sudo app-listener network-monitor /usr/bin/bash                    # watch bash network ops
sudo app-listener network-guard -w /usr/lib/firefox/firefox --auto-infra
sudo app-listener daemon --genkey                                  # fscrypt master key
sudo app-listener daemon --headless --blocked-only                 # protected daemon
sudo systemctl reload app-listener-daemon                          # re-resolve whitelist inodes
sudo app-listener update --yes                                     # self-update from GitHub
```

Exit any TUI with `q` or `Ctrl+C`.

## How it works

eBPF programs, compiled and embedded in the binary, attach at three levels:

- **kprobes** (`monitor`) — observe all I/O regardless of syscall path (io_uring, splice, sendfile, mmap) plus metadata ops (chmod/truncate/stat/access/readlink/mknod).
- **LSM hooks** (`guard`, `network-guard`, `daemon`) — the only kernel mechanism that can **deny**; 23 hooks covering open, read/write, mmap, unlink/rename/symlink/link/mkdir/rmdir/mknod, attributes, stat/access/readlink, mount, ptrace and exec.
- **tracepoints/kretprobes** (`network-monitor`) — TCP/UDP/DNS operations.

**Binary identity is by exe inode**, never by name: renaming a binary or comm-spoofing cannot bypass policy. Keep whitelisted binaries **root-owned** — an attacker who can modify a binary's contents owns its identity anyway.

Any mode can mirror its TUI into a browser with `--serve[=host:port]` (loopback by default): the local terminal TUI keeps running and the same event stream is shared read-only over WebSockets. `--user`/`--password` add HTTP Basic Auth (both required together). No TLS — put a reverse proxy in front when binding off-loopback. Mutually exclusive with `--headless` and `--gui`.

## Modes

### monitor — observe

Traces file operations under watched paths; nothing is blocked. Covers data I/O (`OPEN/READ/WRITE/MMAP`), tree changes, and metadata: `ATTR` (chmod/chown/utimes/truncate/setxattr), `STAT` (stat/access/readlink), `MKNOD`.

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
| `--headless` | `false` | No TUI; log to stderr |
| `--gui` | `false` | Desktop GUI (needs a `-tags gui` build — `make build-linux GUI=1`) |

### guard — block file access

Denies file operations on the guarded path by process identity, in **blacklist** or **whitelist** mode (default whitelist; omitted `-w` blocks everything).

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
| `-e, --events <list>` | all | Event type filter (same set as monitor) |
| `--headless` | `false` | No TUI; log `GUARD\|` events to stderr |

**Exec-open attribution (whitelist mode)**: executing a binary is an OPEN performed by the *launcher* — a shell wrapper runs through its interpreter, which the whitelist deliberately excludes. Opens are attributed to the **binary being executed** instead, so whitelisted binaries *inside* the guarded tree (e.g. Discord under `~/.config/discord`) work from any shell or wrapper. The exec fd is never exposed to the launcher; helper binaries an app spawns must be whitelisted explicitly. Blacklist mode always attributes to the launcher.

### network-monitor — watch network

Traces network operations (TCP, UDP, DNS) of the listed binaries only.

```bash
sudo ./build/linux/app-listener network-monitor /usr/bin/curl /usr/bin/wget -e CONNECT,ACCEPT,DNS
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

Denies socket operations by binary identity (blacklist or whitelist) via LSM `socket_connect/bind/listen/sendmsg/recvmsg` hooks.

```bash
sudo ./build/linux/app-listener network-guard -b /usr/bin/curl -e CONNECT,SEND
sudo ./build/linux/app-listener network-guard -w /usr/bin/vim --auto-infra   # only vim + system infra
```

| Flag | Default | Description |
|------|---------|-------------|
| `-b, --blacklist <binary>` | — | Block network ops for these binaries only (exclusive with `-w`) |
| `-w, --whitelist <binary>` | — | Block network ops for **all** binaries except these (default deny) |
| `--auto-infra` | `false` | Auto-allowlist running infra daemons (resolved, NetworkManager…) — otherwise DNS breaks for everyone |
| `--unsafe` | `false` | Also block AF_UNIX (X11, D-Bus, systemd) — may break the desktop |
| `--no-throttle` | `false` | Disable rate limiting (1 event/type/process per 250 ms default) |
| `-e, --events <list>` | all | `CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS,BIND,LISTEN` |
| `--headless` | `false` | No TUI; log `NETGUARD\|` events to stderr |

**Pick the real executable, not a wrapper** — identity is the exe inode, so whitelisting a `#!/bin/sh` wrapper matches nothing:

```bash
ls -l /proc/$(pgrep -n firefox)/exe    # running process → real binary
readlink -f /usr/bin/firefox           # follows symlinks, NOT shell wrappers
```

### daemon — fscrypt + whitelist lifecycle

Config-driven daemon protecting any number of directories with the guard's whitelist engine plus an fscrypt encryption lifecycle: resources are unlocked at startup and locked again on shutdown **while the guards remain attached** — never an unprotected window.

```bash
sudo ./build/linux/app-listener daemon --genkey   # create the fscrypt master key
sudo ./build/linux/app-listener daemon            # /etc/app-listener/daemon.conf → daemon-samples/daemon.conf
```

| Flag | Default | Description |
|------|---------|-------------|
| `--config <path>` | — | Config resolution: flag → `/etc/app-listener/daemon.conf` → `daemon-samples/daemon.conf` |
| `--headless` | `false` | Log `DAEMON …` to stderr (journald when a systemd service) |
| `--blocked-only` | `false` | Print only denied attempts (presentational) |
| `--genkey` | `false` | Generate `/etc/app-listener/fscrypt.key` and exit (regeneration asks for confirmation) |
| `--pprof <addr>` | — | Serve `net/http/pprof` on a loopback address for profiling |

Config grammar (full template in `daemon-samples/daemon.conf`):

```text
[watch /home/alice/.ssh]          # one section per protected directory
need_encryption: true             # default true; false skips the fscrypt lifecycle
/usr/bin/ssh READ,WRITE           # whitelisted binary, restricted to these events
/usr/bin/ssh-agent                # bare path = all events allowed
```

- **Whitelist only, default deny**; identity by inode.
- **Per-binary event masks**: unlisted events are denied; `READ`/`WRITE`/`MMAP` imply `OPEN`. Masks are a least-privilege hint, **not a confinement boundary** — a masked binary that execs another whitelisted binary escapes its mask.
- **SIGHUP reload** (`systemctl reload`): recomputes every binary's inode identity atomically — new guards attach before old ones detach; a malformed config keeps the previous one running.
- **Live edit control socket**: when `/etc/app-listener/edit-auth.hash` exists (an edit-protected password was set), the daemon opens `/run/app-listener-daemon.control` (`0600`, root + app-listener-binary peer only) so `edit-protected` can make an authenticated live change — see the [edit-protected](#installsh--install--uninstall--update--edit-protected) section.
- **fscrypt lifecycle**: `need_encryption: true` resources must already carry an fscrypt policy. Shutdown deprovisions keys in two passes while guards still deny access; hooks detach only after every vault is keyless. A hard `SIGKILL` cannot be caught, but the guard's LSM links are pinned to `/sys/fs/bpf` so the trees stay enforced until `ExecStopPost` locks the vaults.

### install.sh / install / uninstall / update / edit-protected

**`scripts/install.sh`** (the `curl … | sudo bash` one-liner) — runs `check-compatibility` and aborts if it fails; downloads the latest release of `--channel` (`release` [default] / `prerelease`) from GitHub; verifies the Ed25519 signature of the checksum against the embedded release key, the checksum against the binary, and the GitHub asset digest; then atomically installs `/usr/local/sbin/app-listener` + the PATH symlink. It does **not** install the daemon — it prints the reminder to run `sudo app-listener install` yourself.

**install** — TUI wizard, in safe order: stop a running daemon → build → generate the fscrypt key (existing kept) → pick users → probe a built-in catalog of critical directories (`internal/install/catalog.go`: SSH, GPG, AI agents, browsers, VPNs, password stores…) → encrypt selected directories (backup first, verified against the master key) → deploy systemd units, the package-manager catalog-refresh hook, per-user ssh-agent unit, binary and config.

**install --diff-catalog** — the *incremental* wizard: after a catalog update or a newly installed app, it lists the critical directories that now exist on the host but are **not** yet in `daemon.conf`, lets you pick which to add in the same picker as the full install, appends them to the config and encrypts them (backup first) — every existing section is left byte-for-byte intact. Stops the daemon for the cycle, restarts it on the merged config. `--update-catalog-only` refreshes existing sections' whitelists; `--diff-catalog` adds new sections — run both to fully re-sync. Requires a previous installation; interactive only.

**install --update-catalog-only** (the package hooks + the boot-time unit) — re-expands every catalog-matched whitelist and rewrites the config. Default: stops the daemon, unlocks each vault under an ephemeral self-only guard. `--live` (daemon running): no stop, no lock churn — applied via SIGHUP. It runs automatically from:
> - **pacman** — `/etc/pacman.d/hooks/50-app-listener-reload.hook` (`PostTransaction`);
> - **apt/dpkg** — `/etc/apt/apt.conf.d/95app-listener-reload` (`DPkg::Post-Invoke`);
> - **every boot** — `app-listener-catalog-refresh.service` runs a `--live` refresh once, after the daemon, catching package changes made while no hook fired (offline installs, direct `dpkg -i`/`pacman -U`, a failed hook). dnf/zypper have no bundled hook — the boot unit covers them, or run it by hand.

> **Self-updating apps (Discord, VS Code helpers…)** change their own binaries outside the package manager, so the whitelist (pinned to inodes) goes stale and the app breaks until refreshed. The package hooks do **not** fire for these. Fix, no downtime:
> ```bash
> sudo app-listener install --update-catalog-only --live --yes
> ```
> The catalog narrows watches for these apps: only the sensitive subtrees are guarded (Discord's `Local Storage/`, `Cookies`, …), the vault root the updater writes to stays unguarded.
>
> **fscrypt prerequisite**: each filesystem must be initialized (`fscrypt setup --all-users`) and support encryption (ext4: `tune2fs -O encrypt <dev>`). The installer checks this before migrating anything and, for a fixable gap, shows the exact command + reason and offers to run it for you (it is already root) — decline and it aborts, as before. Every command it may run is listed in `internal/fscrypt/prereq.go`.

**uninstall** — refuses while the daemon runs; re-scans the catalog; decrypts in place by default; deletes the master key only with `--delete-key`; at the end, lists any `.app_listener.backup` migration copies and offers to delete them (all preselected, one confirmation — plain unencrypted copies, harmless to keep).

**update** — self-updates from the latest signed `pre-YYYYMMDD-<sha>` GitHub pre-release (Ed25519 signature + checksum + asset digest all verified before anything is written).

**edit-protected** — edit a protected directory in a two-pane editor. `Ctrl+S` saves atomically; binaries, symlinks and files > 2 MiB refused.

Two modes, chosen automatically:

- **offline** (default) — refuses while the daemon runs; re-scans the catalog, unlocks one vault with the master key, edits, re-locks it on exit (a vault never stays open).
- **live** — when an *edit-protected password* was chosen during `sudo app-listener install` **and** the daemon is running: you enter the password **once, first**; only then does the daemon (over a local root-only control socket, `/run/app-listener-daemon.control`) return the list of guarded directories to pick from. It briefly grants write access to the one you pick, you edit, and the grant is dropped. The daemon keeps running and the fscrypt vaults are never touched. An unauthenticated caller learns **nothing** about which directories are protected. Grants are one-at-a-time, capped at 30 min, and revoked on disconnect or a SIGHUP reload; the socket locks out after 5 failed attempts.

The password is **separate from the fscrypt master key** — it only authenticates `edit-protected`. Its PBKDF2 hash lives at `/etc/app-listener/edit-auth.hash` (`0600` root), guarded by the running daemon the same way as `fscrypt.key` (readable/writable only by the app-listener binary).

```bash
sudo app-listener edit-protected                       # auto: live if a password is set + daemon up, else offline
sudo app-listener edit-protected --set-password         # set/rotate the password (refused if it was set at install — re-run install to rotate that)
sudo app-listener edit-protected --clear-password       # remove it (disables live mode)

# non-interactive live write (automation): password from $APP_LISTENER_EDIT_PASSWORD
echo "new contents" | sudo APP_LISTENER_EDIT_PASSWORD=… \
  app-listener edit-protected --resource /home/alice/.ssh --put config
```

Before exiting, both modes audit the edited tree against `daemon.conf` and warn (advisory, acknowledged once) about anything that would sit outside the daemon's protection: a file created outside every guarded watch path of a grouped vault, a new symlink, a world-readable new secret, a freshly dropped executable.

## Debug

```bash
systemctl is-active app-listener-daemon              # expect: active
journalctl -u app-listener-daemon -f                 # follow live
sudo journalctl -u app-listener-daemon -f | grep -i denied

sudo systemctl reload app-listener-daemon            # SIGHUP after a package update replaced a binary

# manual guard check
sudo systemctl stop app-listener-daemon
sudo /usr/local/sbin/app-listener daemon --headless --verbose 3
#   another terminal: ssh -T git@github.com → must WORK
#                     cat ~/.ssh/id_ed25519  → must be DENIED

sudo fscrypt status /home/alice/.ssh                 # Encrypted / Not encrypted
```

A wrong/old key fails immediately with "invalid wrapping key" — the daemon never silently generates a new one. Backups live at `<dir>.app_listener.backup`; `sudo app-listener install --restore-backups` restores them.

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Dockerized build (rootful): regenerate BPF bindings + build to `build/linux/app-listener` |
| `make build-host` | On-host build (needs clang/LLVM, bpftool, Go, GCC) |
| `make build-linux` | Build the Go binary only (`GUI=1` links the desktop GUI) |
| `make check-compatibility` | Static host check — can it run app-listener? |
| `make test` / `make lint` | Unit tests / golangci-lint |
| `make test-integration` | Docker integration + bypass suite (rootful Docker) |
| `make generate` | Regenerate BPF bindings |
| `make clean` | Remove build artifacts |

## Docker

```bash
docker compose build
docker compose run --rm app-listener monitor -w /tmp
```

Multi-stage build, Debian-slim runner. eBPF needs the host kernel — run privileged or with `/sys/kernel/btf/` + `CAP_BPF`. `make build` uses a separate toolchain image (`docker/builder.Dockerfile`), not this one.
