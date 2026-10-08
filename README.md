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
| **Arch Linux** | rolling (6.x) | ✅ **Fully supported** — every LSM hook attaches | Add `bpf` to the `lsm=` list on your boot entry's kernel cmdline, reboot |
| **Ubuntu 24.04 LTS** | 6.8 | ✅ **Fully supported** | Append `lsm=landlock,lockdown,yama,integrity,apparmor,bpf` to `GRUB_CMDLINE_LINUX_DEFAULT`, `sudo update-grub`, reboot |
| **Ubuntu 22.04 LTS** | 6.x (HWE) | 🟡 **HWE kernel only** — on the 5.15 GA kernel `guard` / `daemon` do not load (they need `bpf_loop`, kernel ≥ 5.17); `monitor` and `network-guard` still work there. Install the HWE kernel (`linux-generic-hwe-22.04`, 6.x) for full support | same as 24.04 |

**Kernel floor:** `monitor` needs ≥ 5.8 (BPF ring buffer); `network-guard` needs ≥ 5.10 (BPF-LSM); `guard` / `daemon` need ≥ 5.17 (BPF-LSM plus the `bpf_loop` helper). Ubuntu 20.04 (kernel 5.4) is not supported. Mandatory hooks: `file_open`, `file_permission`, and `inode_unlink`, `inode_rename`, `bprm_committed_creds` plus the `sched_process_fork` tracepoint (replaced-binary tracking), and the `sched_process_exec` tracepoint (launch scan of whitelisted runtimes) — every other LSM hook is best-effort: a kernel that lacks one logs a warning and keeps enforcing the rest.

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
sudo app-listener daemon --headless --blocked-only --no-log-metadata-blocks   # protected daemon
sudo systemctl reload app-listener-daemon                          # re-resolve whitelist inodes
sudo app-listener update --yes                                     # self-update from GitHub
```

Exit any TUI with `q` or `Ctrl+C`.

## How it works

eBPF programs, compiled and embedded in the binary, attach at three levels:

- **kprobes** (`monitor`) — observe all I/O regardless of syscall path (io_uring, splice, sendfile, mmap) plus metadata ops (chmod/truncate/stat/access/readlink/mknod).
- **LSM hooks** (`guard`, `network-guard`, `daemon`) — the only kernel mechanism that can **deny**; ~30 hooks covering open, read/write, mmap, unlink/rename/symlink/link/mkdir/rmdir/mknod, attributes, stat/access/readlink, mount, btrfs copy ioctls, ptrace and exec.
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
| `--no-log-metadata-blocks` | `false` | Don't log denied metadata-only process inspections (`op=PTRACE mode=READ`, e.g. the compositor or audio server reading `/proc/<pid>` of Steam); still denied. Without it they are logged once per caller/target and summarized per minute. Memory and file denials are always logged. Used by the systemd unit. Probes the kernel itself doesn't audit (`ps`/`pgrep` reading `/proc/<pid>/stat`) are never logged. A refusal for a process holding several resources' content reads `resource=multiple`, or the comma-separated paths of its taint set |
| `--genkey` | `false` | Generate `/etc/app-listener/fscrypt.key` and exit (regeneration asks for confirmation) |
| `--check` | `false` | Preflight: check BPF-LSM is active and load every guard eBPF program into this kernel's verifier, attaching nothing; exit non-zero on any rejection |
| `--verifier-only` | `false` | With `--check`: skip the BPF-LSM activation check (CI hosts that build BPF-LSM in without enabling it). A pass then says the programs verify, not that the host enforces |
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
- **Binary replacement**: a whitelisted binary at a user-writable path may be modified only by a binary of a resource that whitelists it (never a shell, interpreter, `git`, `cp`, …). A new inode at a whitelisted path is re-admitted, with no reload, if one of those updaters created it, or if it is a root-owned file in a root-owned directory on a filesystem root mounted without `nosuid` (a package upgrade); anything else stays denied (`not re-admitted`), at the next reload or start too, until confirmed with `trust-binaries` (binary ledger below). A new root-owned match of the catalog's absolute patterns (`/opt/*/jbr/bin/java`) is admitted the same way.
- **Binary ledger** (`/etc/app-listener/binaries.db`): nothing guards a whitelisted path while the daemon is stopped, booting or being reinstalled, so every start and reload checks each whitelisted binary against the content hash recorded for its line (the link path for a symlink line). A binary is admitted when only root could have placed it (root-owned file and directories, outside every user home, on a mount without `nosuid`), when its hash matches, or when its resource's updater created it while the daemon ran (recorded automatically, as are live replacements). Anything else — a binary changed while the daemon was not running in a home or a user-owned directory, or a new catalog match planted meanwhile — is refused: `DAEMON binary-unconfirmed resource=… path=… reason=new|changed sha256=…` in the journal, no access, and no updater or reserved-name writer rights either. Confirm it with `sudo app-listener trust-binaries` (below). The first start after upgrading from a version without the ledger records what it finds once (with a warning); lines the install wizard just added are recorded at the next start. The ledger is self-guarded with the rest of `/etc/app-listener` and removed by `uninstall`.
- **Boot ordering**: the unit is `Type=notify` and ordered `Before=systemd-user-sessions.service cron.service crond.service cronie.service atd.service`, so logins, lingering user managers, cron and `at` jobs wait until every guard is attached (READY=1) — no user code runs while the whitelisted paths are unguarded. Ordering only: a failed daemon start does not lock users out.
- **Launch scan** (whitelisted Node/Electron, Chromium and JVM runtimes, detected from the binary's own strings or the name `java`): identity is the executable, so `ELECTRON_RUN_AS_NODE=1 code -e …` or `NODE_OPTIONS=--require=/tmp/x.js discord` is the whitelisted binary running its caller's code. At exec, before any of its code runs, the daemon judges the new process's environment and arguments: `ELECTRON_RUN_AS_NODE`, `NODE_PATH`, `NODE_OPTIONS` / `JAVA_TOOL_OPTIONS` / `_JAVA_OPTIONS` / `JDK_JAVA_OPTIONS` holding anything but a short allowlist of harmless options (heap size, warnings, TLS/CA settings, font and UI scaling), and `NODE_REPL_EXTERNAL_MODULE`, `--inspect*`, `--eval`, `--print`, `--require`, `--import`, `--loader`, `--env-file*`, `--remote-debugging-*`, `--load-extension`, the child-process launchers `--renderer-cmd-prefix` / `--utility-cmd-prefix` / `--zygote-cmd-prefix` / `--gpu-launcher` / `--browser-subprocess-path`, and VS Code's `--extensions-dir`, `--extensionDevelopmentPath`, `--extensionTestsPath`, `--install-extension` (also spelled with Chromium's single dash, or with Node's `_` for `-`). `--user-data-dir` is not judged: Electron apps pass it, with their own profile, to the processes they start. Such a launch runs, but the process (and its forks; an exec is judged afresh) never reads a protected file nor changes an app's read-only code dir: `TRUST DENIED  op=LAUNCH comm=… path=…`. It also loses the binary's updater and writer rights: it can't replace the app's binaries, create a reserved name (a new version's binary, `settings.json`), and nothing it creates is admitted as updater-made (`op=PLANT` / `op=WRITEBLOCK`). Start the app without that variable or flag. A process that mapped untrusted code before its binary was whitelisted is refused the same way (`op=PRELOADED`). Such a process also gets nothing else below a whitelist-mode tree: no `stat`, directory listing or xattr.
- **Generic Electron** (`[electron_apps]` block): an Electron with `resources/default_app.asar` beside it (Arch's `/usr/lib/electronNN/electron`, which `signal-desktop`, `obsidian`, … wrapper scripts start) runs whatever app its first non-switch argument names. Whitelisted, it may be started only with an app listed under `[electron_apps]` — an absolute path written exactly as the wrapper passes it (≤ 43 bytes), naming a root-owned file (an `.asar` or a JS entry point; a directory is refused) in root-owned directories on a filesystem mounted without `nosuid` (else `electron app … refused`). Any other app, a relative path, or default_app's `-r`/`-i`/`-repl` or `--app=<file>` (which picks the app itself) is a risky launch (`op=LAUNCH`); no app at all (Chromium's child processes) is not. Apps share the daemon-wide list: an Electron whitelisted for two sections may run either section's app. A bundled Electron (Discord, VS Code: `resources/app.asar` or `resources/app/`, no `default_app.asar`) ignores its arguments when choosing the app and is not judged this way.

  ```ini
  [electron_apps]
  /usr/lib/signal-desktop/app.asar
  ```
- **What identity cannot stop**: a whitelisted program that discloses its secret on request still does — `gh auth token`, `git credential fill`, `ssh-add -L`, an app's own export or debug command. The guard authorizes *which executable* opens the file, not *why*: any process that can run the whitelisted binary with arguments of its choosing gets whatever that binary prints. Whitelist only what needs the resource, prefer apps that keep secrets in a running agent, and never whitelist a general tool — a shell, interpreter (`node`, `python`, `java`), `git`, `cp`, `curl` … reads the resource on behalf of any script. The installer and the daemon warn about every resource that whitelists one.
- **Multicall binaries**: uutils coreutils (Ubuntu 25.10+/26.04) is one file behind every applet name. Its applets are told apart: the kernel attests at exec which applet runs (the exec'd basename and `argv[0]` must both name it), so whitelisting `/usr/bin/head` admits `head` only, never `cat` or `dd`. Whitelist an applet by its own link (`/usr/bin/<applet>` or `/usr/lib/cargo/bin/coreutils/<applet>`), not a symlink of another name. BusyBox and toybox run applets without an exec, so they can't be told apart, and any other binary that calls itself a multi-call binary (a future uutils diffutils/findutils) is treated the same until it is vetted: `guard -w/-b`, `network-guard -w/-b` and `edit-protected --forward` refuse them (as they refuse a uutils path that isn't an applet's own name), the installer and `edit-protected --edit-config` refuse a new such line, the catalog never writes one, and the daemon drops one already in daemon.conf (CRITICAL in the log) while keeping everything else enforced. Use a dedicated binary (on Ubuntu the GNU `gnu*` tools, or `coreutils-from-gnu`). `network-monitor` watches an applet the same way; a multicall it can't tell apart is watched by its file, with a warning. Whitelisting a single binary with several hard links logs a warning: all its names share one identity, which is fine for one program under several names (`perl`, `e2fsck`) but not for a multicall the classifier doesn't recognize.
  - **Limitation — processes started before the guard:** the applet is attested at exec, so a multicall process already running when a guard starts (or the daemon restarts) has no attested applet (a `systemctl reload` keeps them; a restart doesn't). A whitelist denies it and a blacklist blocks it until it execs again; `network-monitor` doesn't report it. Restart such long-running applet processes (a `tail -f`, a `sleep` loop) after starting the guard.
- **Replaced binaries**: a process started before its binary was replaced keeps its access, through a reload too. Re-running the replaced image afterwards (an fd held across the update, `/proc/<pid>/exe`) does not, and neither does a file that later reuses the old inode number (`exe_supersede.h`).
- **Catalog refresh**: the daemon watches the catalog's directories (inotify, a hint only) and re-expands the catalog whitelists at startup and whenever a binary is written below a catalog pattern — the same rules as `install --update-catalog-only`, keeping lines whose binary still exists. `daemon.conf` is rewritten only while it still holds the configuration the daemon runs, then reloaded in-process (`DAEMON catalog-refresh resource=… admitted=… dropped=…` in the journal). No package-manager hook or boot unit is involved.
- **SIGHUP reload** (`systemctl reload`): recomputes every binary's inode identity atomically — new guards attach before old ones detach; a malformed config keeps the previous one running.
- **Process inspectors** (`[inspectors]` block): a process that read a protected file is *tainted* — no program outside its resource's whitelist may read its memory or even its `/proc/<pid>` metadata. Programs listed under `[inspectors]` may read that metadata (`environ`, `fd`, `maps`, `root`, `exe`) of **every** protected process, never its memory (`ptrace`, `process_vm_readv`, `/proc/<pid>/mem`) nor a protected file. Default without the block: `xdg-desktop-portal`, which opens `/proc/<pid>/root` of each screen-sharing caller and otherwise refuses the request (Meet/Discord sharing fails with no picker); an empty `[inspectors]` block disables it. Only a root-owned binary in root-owned directories on a filesystem mounted without `nosuid` is accepted (else `inspector … refused`); it is trusted-binary protected (no `LD_PRELOAD`), and attaching to or tracing an inspector needs a program **every** section whitelists. Exposure: whatever the kernel gates at `PTRACE_MODE_READ` — each protected process's environment variables (tokens passed via env included), open-file paths and memory map, the contents of its *unguarded* open files reachable by reopening `/proc/<pid>/fd/N` (pipes, memfds, deleted files), and `perf_event_open` samples of its user stack. Protected files stay denied through `/proc/<pid>/fd`. Never list a shell, interpreter, `cat`, `ps`, `perf` or anything that does what its caller asks.
- **Filesystem-wide gates** (coarse by nature, judged globally so only the daemon's own binaries pass): the raw block device backing a guarded filesystem is blocked (`debugfs`/`dd`/`fsck` can't reparse the on-disk image around the VFS), and on **btrfs** the copy ioctls — snapshot (`SNAP_CREATE`/`_V2`), `SEND` and `TREE_SEARCH` — are denied on the whole superblock hosting a guarded root (a snapshot is a new subvolume with the source's inode numbers, so no per-inode key would match the copy). Both block *every* path on that device/filesystem, guarded or not: while the daemon runs, raw-imaging backups and btrfs snapshot tools (`snapper`, `timeshift`, `btrbk`, `snap-pac`, `btrfs subvolume snapshot/send/list`) stop working on any disk or btrfs filesystem that hosts a guarded resource — intended. Logged as `resource=raw-block-device` / `resource=btrfs-ioctl`, never a watched path.
- **`trust-binaries`**: lists the binaries the daemon refused (line, resolved path, resource, sha256, owner, mtime) and asks to confirm each; `--yes` confirms all, `--list` only lists, explicit paths confirm those lines as they are now (a line added to `daemon.conf` by hand). Only the exact content that was refused is confirmed — a file that changed again is skipped and judged afresh. A running daemon is reloaded afterwards (`--no-reload` to skip).

  ```bash
  sudo app-listener trust-binaries --list
  sudo app-listener trust-binaries                       # review and confirm one by one
  sudo app-listener trust-binaries /opt/myapp/bin/myapp  # confirm a hand-added line
  ```
- **Live edit control socket**: when `/etc/app-listener/edit-auth.hash` exists (an edit-protected password was set), the daemon opens `/run/app-listener-daemon.control` (`0600`, root + app-listener-binary peer only) so `edit-protected` can make an authenticated live change — see the [edit-protected](#installsh--install--uninstall--update--edit-protected) section.
- **fscrypt lifecycle**: `need_encryption: true` resources must already carry an fscrypt policy. Shutdown deprovisions keys in two passes while guards still deny access; hooks detach only after every vault is keyless. A hard `SIGKILL` cannot be caught, but the guard's LSM links are pinned to `/sys/fs/bpf` so the trees stay enforced until `ExecStopPost` locks the vaults.

### install.sh / install / uninstall / update / edit-protected

**`scripts/install.sh`** (the `curl … | sudo bash` one-liner) — runs `check-compatibility` and aborts if it fails; downloads the latest release of `--channel` (`release` [default] / `prerelease`) from GitHub; verifies the Ed25519 signature of the checksum against the embedded release key, the checksum against the binary, and the GitHub asset digest; then atomically installs `/usr/local/sbin/app-listener` + the PATH symlink. It does **not** install the daemon — it prints the reminder to run `sudo app-listener install` yourself.

**install** — TUI wizard, in safe order: stop a running daemon → build → generate the fscrypt key (existing kept) → pick users → probe a built-in catalog of critical directories (`internal/install/catalog.go`: SSH, GPG, AI agents, browsers, VPNs, password stores…) → encrypt selected directories (backup first, verified against the master key) → deploy systemd units, binary and config (removing the pacman/apt catalog-refresh hooks and the boot-time refresh unit an earlier version installed: the daemon refreshes the catalog itself). The per-user ssh-agent unit is offered only when that user's `~/.ssh` ends up guarded (one question per user, naming the user and the unit path, asked before encryption), together with `AddKeysToAgent yes` at the top of the user's `~/.ssh/config` (created if missing, an existing `AddKeysToAgent` kept; written before `~/.ssh` is encrypted, so `ssh` loads a key into the agent on first use), and with a marked `SSH_AUTH_SOCK` block appended to that user's shell startup file (`~/.zshrc` / `~/.bashrc` / fish `conf.d`; login shell's is created if missing, an existing own `SSH_AUTH_SOCK` line or an already-set variable is never overridden, and `uninstall` removes the block); skipped entirely otherwise. The daemon never loads keys itself. When a selected app is Bun-based (e.g. **opencode**), Bun extracts a bundled native library to `$TMPDIR/.bun-<uid>-<hash>.so` and `dlopen`s it — on world-writable `/tmp` the trust guard denies that load (it can't tell the app's own extraction from a planted `.so`). The wizard offers, per user, to redirect it: it creates a private `~/.cache/app-listener/bun` (`0700`) and adds a marked shell-function wrapper for each launcher (`~/.zshrc` / `~/.bashrc` / fish `conf.d`) that runs the command with `$TMPDIR` pointed there. The daemon then reserves the `.bun-*` name below that dir for the app's binaries (`guard_trust.bpf.c` #3), so the app loads its own extraction while no other process can plant or load one. A dir the daemon has not reserved yet (at every start, and when it first appears) may already hold a plant, so the daemon first replaces it with an empty one — the app re-extracts on its next launch; ordinary temp files the app's subprocesses write to that dir are unaffected (only the reserved name is gated). Interactive shells only — a Bun app launched from a desktop menu still uses `/tmp`; `uninstall` removes the wrappers.

**install --diff-catalog** — the *incremental* wizard: after a catalog update or a newly installed app, it lists the critical directories that now exist on the host but are **not** yet in `daemon.conf`, lets you pick which to add in the same picker as the full install, appends them to the config and encrypts them (backup first) — every existing section is left byte-for-byte intact. Stops the daemon for the cycle, restarts it on the merged config. `--update-catalog-only` refreshes existing sections' whitelists; `--diff-catalog` adds new sections — run both to fully re-sync. Requires a previous installation; interactive only.

**install --allow-metadata-output** — installs the daemon unit without `--no-log-metadata-blocks`, so denied metadata-only process inspections (`op=PTRACE mode=READ`) are logged. Use it to diagnose an app that breaks with nothing logged; re-run `install` without it to quiet them again.

**install --update-catalog-only** (manual / debug) — re-expands every catalog-matched whitelist and rewrites the config, dropping every line the catalog does not produce. Default: stops the daemon, unlocks each vault under an ephemeral self-only guard. `--live` (daemon running): no stop, no lock churn — applied via SIGHUP. A running daemon does this on its own (see the daemon's catalog refresh above), so you normally never run it.
>
> The refresh also adopts an app's fixed library dirs as read-only `lib_dir`s once they exist (Steam's `ubuntu12_*`, `linux*`, `steamrt64`, `compatibilitytools.d`) and trusts their libraries for that app. So a non-app process can't plant one ahead of it, the daemon reserves those dir names for the app's own binaries, present or not: e.g. only Steam can create `~/.local/share/Steam/compatibilitytools.d`, not your shell or ProtonUp-Qt — the same limit an adopted lib dir already has, since it is read-only to everything but the app.

> **Self-updating apps (Discord, VS Code helpers…)** change their own binaries outside the package manager. The daemon follows them: an update written in place of the binary by the app's own updater is re-admitted directly, and a new version directory (Discord's `~/.config/discord/0.0.N/`) is admitted by the catalog refresh as soon as the updater finishes writing the binary — no restart, no manual step. A process that creates such a directory without being one of the app's binaries cannot plant the binary name there (reserved glob names), so nothing it does is admitted. The catalog narrows watches for these apps: only the sensitive subtrees are guarded (Discord's `Local Storage/`, `Cookies`, …), the vault root the updater writes to stays unguarded. The code they run is not left writable either: Discord's `app-*` version dirs (inside its vault, guarded once it is unlocked), `~/.vscode`, `~/.vscode-insiders`, `~/.vscode-oss` (extensions), `~/.local/share/JetBrains` (plugins) and Steam's update staging `~/.local/share/Steam/package` are read-only `lib_dir`s only the app's own binaries may change, and Discord's updater inputs (`settings.json`, which names the update endpoint, and `installer.db`) are reserved names only Discord's binaries may write — so no other process can make the genuine updater install chosen code. There is no npm entry: npm runs as `node`, and whitelisting `node` would let `node -e` read the tree.
>
> **fscrypt prerequisite**: each filesystem must be initialized (`fscrypt setup --all-users`) and support encryption (ext4: `tune2fs -O encrypt <dev>`). The installer checks this before migrating anything and, for a fixable gap, shows the exact command + reason and offers to run it for you (it is already root) — decline and it aborts, as before. Every command it may run is listed in `internal/fscrypt/prereq.go`.

**uninstall** — refuses while the daemon runs; re-scans the catalog; decrypts in place by default; deletes the master key only with `--delete-key`; at the end, lists any `.app_listener.backup` migration copies and offers to delete them (all preselected, one confirmation — plain unencrypted copies, harmless to keep).

**update** — self-updates from the latest signed `pre-YYYYMMDD-<sha>` GitHub pre-release (Ed25519 signature + checksum + asset digest all verified before anything is written).

**edit-protected** — edit a protected directory in a two-pane editor. `Ctrl+S` saves atomically; binaries, symlinks and files > 2 MiB refused.

Two modes, chosen automatically:

- **offline** (default) — refuses while the daemon runs; re-scans the catalog, unlocks one vault with the master key, edits, re-locks it on exit (a vault never stays open).
- **live** — when an *edit-protected password* was chosen during `sudo app-listener install` **and** the daemon is running: you enter the password **once, first**; only then does the daemon (over a local root-only control socket, `/run/app-listener-daemon.control`) return the list of guarded directories to pick from. It briefly grants write access to the one you pick, you edit, and the grant is dropped. The daemon keeps running and the fscrypt vaults are never touched. An unauthenticated caller learns **nothing** about which directories are protected. Grants are one-at-a-time and revoked on disconnect, on a SIGHUP reload, or after `--timeout-session` (default `30m`, up to `24h`; e.g. `10m`, `45s`) without activity. Every editor keypress resets that timer. The socket locks out after 5 failed attempts.

The password is **separate from the fscrypt master key** — it only authenticates `edit-protected`. Its PBKDF2 hash lives at `/etc/app-listener/edit-auth.hash` (`0600` root), guarded by the running daemon the same way as `fscrypt.key` (readable/writable only by the app-listener binary).

```bash
sudo app-listener edit-protected                       # auto: live if a password is set + daemon up, else offline
sudo app-listener edit-protected --set-password         # set/rotate the password (refused if it was set at install — re-run install to rotate that)
sudo app-listener edit-protected --clear-password       # remove it (disables live mode)
sudo app-listener edit-protected --forward -w <bin> [-e EVENTS]  # temporary access for other binaries (see below)
sudo app-listener edit-protected --timeout-session 10m  # revoke a live session after 10 idle minutes (default 30m)
sudo app-listener edit-protected --edit-config         # edit the running daemon.conf; saving reloads it (live mode)

# non-interactive live write (automation): password from $APP_LISTENER_EDIT_PASSWORD
echo "new contents" | sudo APP_LISTENER_EDIT_PASSWORD=… \
  app-listener edit-protected --resource /home/alice/.ssh --put config
```

**`--forward`** (live mode only) opens no editor. It gives *other* binaries temporary access to the guarded directories you pick, using `guard`'s `-w`/`-b`/`-e` syntax: `-w` admits binaries the resource doesn't whitelist, `-b` denies whitelisted ones, and `-e` limits the rule to those events (default: all). `-w` and `-b` are mutually exclusive. The password is asked first, then a multi-select picks the directories (or pass `--resource` once per directory). While the grant is active, the terminal shows the granted binaries' accesses to those directories, allowed and denied, in the same view as `guard`; without a terminal they're printed as `GUARD|` lines. The access lasts until you quit (`q` or Ctrl+C). It is also revoked when the client dies, on a daemon reload, when a granted binary's content changes on disk, or after `--timeout-session` passes with neither a granted binary touching the directories nor a keypress.

```bash
# let graphify (a uv-managed python script) create its skill under a guarded ~/.claude
sudo app-listener edit-protected --forward \
  -w ~/.local/share/uv/python/cpython-3.12.14-linux-x86_64-gnu/bin/python3.12 \
  -e OPEN,READ,STAT,MKDIR,WRITE
```

Things to know before admitting a binary:
- A `-w` grant covers **every user** running that binary, like a `daemon.conf` whitelist line.
- Admitting an **interpreter** (python, node, bash) admits every script it runs. A script's process is its interpreter, so the script itself is not what the guard sees.
- A granted **Node/Electron, Chromium or JVM runtime** is judged at launch like a whitelisted one: started with env or flags that load its caller's code (`ELECTRON_RUN_AS_NODE`, `NODE_OPTIONS=--require`, `--inspect`, `--remote-debugging-port`, `JAVA_TOOL_OPTIONS=-javaagent`, …), it is refused the directories (`TRUST DENIED op=LAUNCH`). This covers an app that embeds a runtime (VS Code, Discord); it doesn't cover `node script.js` or a JVM given a class path, which run the script they're given.
- A binary that isn't root-placed (anything under a home directory) can be rewritten in place by its owner. You're asked to confirm it (`--yes` skips the question and is required without a terminal), and the daemon revokes the grant if its content changes.
- Revocation re-checks already-open files on their next read or write, but memory a process already mapped stays mapped.
- Identity is the binary's inode, resolved by the daemon. A binary already whitelisted on a resource is left untouched by `-w`, and one that isn't is untouched by `-b`. The daemon's own binary, inspectors and tamper-demoted binaries are refused.
- If the daemon is killed mid-grant, its pinned guards keep enforcing *with* the temporary allows. `daemon --lockdown` (ExecStopPost) and the next start strip them using `/etc/app-listener/temp-grants.json`, written before any allow.

**`--edit-config`** (live mode only, same password) opens the running daemon's `daemon.conf` in the editor. Ctrl+S shows the diff and asks for confirmation; then the daemon writes the file and reloads, exactly like `systemctl reload`. If the reload fails (for example the new configuration has no `[watch]` section, or needs more guard slots than are free), the daemon keeps running the previous configuration, puts the previous `daemon.conf` back, and tells you why; you can edit again or discard. If the file changed while you were editing, nothing is overwritten. `--content-file <file>` replaces the configuration without the editor (password from `$APP_LISTENER_EDIT_PASSWORD`). Binaries added this way are vetted like any reload: a binary that isn't root-placed still needs `app-listener trust-binaries`.

The editors (edit-protected, `--edit-config`, the installer's config step) and the diff/changelog viewers open at the first line. Tab / Shift+Tab indent and dedent by four spaces, PgUp/PgDn page, and Ctrl+←/→ jump between words and punctuation, across lines. Ctrl+F opens a find box (Search jumps to the next match, Search all marks every match) and, in the editors, Ctrl+R a find-and-replace box (Replace / Replace all); Tab moves between fields and buttons, Enter activates, Esc closes. Matching is literal and smart-case: case-sensitive only when the query has an upper-case letter. Matches stay in their syntax color, underlined; the current one is reversed. Line numbers sit right-aligned in a fixed-width gutter, so text starts on the same column on every line. Text is syntax-highlighted by file name or `#!` line (Go, Python, C, shell, JSON, YAML, TOML, INI, Markdown, diffs, … and `daemon.conf` itself); buffers over 512 KiB, or files with no recognized language, are not colored.

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
