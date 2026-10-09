# How it works

## Where it sits: DAC, MAC and app-listener

Standard Linux permissions are *discretionary* access control (DAC): owner, group, mode bits.
They decide by **user**, so every program you run as yourself can read every file you own. Your
editor, a package's install script and an AI agent's shell commands are all equal.

*Mandatory* access control (MAC) adds rules the file's owner can't change. AppArmor and SELinux
do this with per-program profiles that describe everything a program may touch, which is
powerful but takes real effort to write and maintain.

app-listener takes the MAC enforcement point (an LSM hook, which the kernel consults after DAC
has said yes) and keeps the policy as simple as DAC's: a protected directory plus the programs
allowed in it. It only restricts what you list, so everything else on the system keeps working
unchanged. It is a BPF LSM of its own: it does not use or require AppArmor or SELinux, and runs
alongside them when they are active.

Applied per directory, this segregates applications from each other: each app's config and
secrets are readable by that app only, not by the other programs running as the same user.

## eBPF at three levels

The eBPF programs are compiled and embedded in the single `app-listener` binary. Each mode
attaches at the level that fits its job:

| Level | Modes | What it gives |
|---|---|---|
| **kprobes** | `monitor` | Observes all I/O regardless of the syscall path (io_uring, splice, sendfile, mmap) plus metadata ops (chmod, truncate, stat, access, readlink, mknod). |
| **LSM hooks** | `guard`, `network-guard`, `daemon` | The only kernel mechanism that can **deny**. About 30 hooks cover open, read/write, mmap, unlink/rename/symlink/link/mkdir/rmdir/mknod, attributes, stat/access/readlink, mount, btrfs copy ioctls, ptrace and exec. |
| **tracepoints / kretprobes** | `network-monitor` | TCP, UDP and DNS operations. |

Because enforcement happens in the kernel, it applies to every way of reaching a file: a plain
`open`, io_uring, `open_by_handle_at`, `copy_file_range`, raw block devices, hard links, renames
and so on. The repository ships a corpus of bypass attempts
([`integrationtests/exploits`](../integrationtests/exploits/README.md)) that the integration tests
run against every enforcing mode.

## Fail closed

An enforcing mode that can't enforce refuses to start. Before attaching, each one checks that
BPF-LSM is in the active LSM list (otherwise hooks attach but deny nothing), and a mandatory hook
that won't attach is a fatal error. See [Compatibility](compatibility.md#mandatory-and-best-effort-hooks).

## Identity is the executable's inode

Policies name programs by path, but app-listener resolves each path to the `dev:ino` of the
executable file and decides on that, never on a name, `comm`, `argv` or anything else the
process controls. Renaming a binary, copying it elsewhere or spoofing its process name does not
move it past policy.

Consequences worth knowing:

- **Keep whitelisted binaries root-owned.** Whoever can modify a binary's contents owns its
  identity anyway. [Binary trust](binary-trust.md) covers how the daemon handles binaries that
  live in user-writable places and binaries that get updated.
- **Whitelist the real executable, not a wrapper script.** A `#!/bin/sh` launcher runs as its
  interpreter. Find the real binary with:

  ```bash
  ls -l /proc/$(pgrep -n firefox)/exe    # running process: the real binary
  readlink -f /usr/bin/firefox           # follows symlinks, NOT shell wrappers
  ```

- **Identity says which program, not why.** A whitelisted program that prints its secret on
  request still does. See [Limitations](limitations.md).

## Output

Every mode has three presentations:

- **TUI** (default). Quit with `q` or `Ctrl+C`.
- **`--headless`**: no TUI, one line per event on stderr, with a stable prefix per mode that
  scripts and the integration tests parse: `EVENT|` (monitor), `GUARD|` (guard), `NETEVENT|`
  (network-monitor), `NETGUARD|` (network-guard), `DAEMON ...` (daemon; journald under systemd).
- **`--gui`** (monitor only): a desktop window, available in builds made with `make build-linux GUI=1`.

### Global flags

| Flag | Description |
|---|---|
| `-v, --verbose <0..4>` | Verbosity of the headless stream and `--dump-log`: `0` errors only, `1` essential, `2` current display, `3` plus internal details, `4` everything. Requires `--headless`. |
| `--dump-log <file>` | Mirror every log line into this file as plain text (asks before overwriting). |
| `--gui` | Desktop GUI instead of the TUI (monitor only, `-tags gui` build). |

## Watching from a browser: `--serve`

Any mode can mirror its TUI into a browser:

```bash
sudo app-listener guard /secret -w /usr/bin/cat --serve                 # 127.0.0.1:9999
sudo app-listener monitor -w /tmp --serve=0.0.0.0:8080 --user me --password s3cret
```

- The local TUI keeps running; the same event stream is shared **read-only** over WebSockets.
- Binds to loopback (`127.0.0.1:9999`) unless you give an address.
- `--user` and `--password` add HTTP Basic Auth and must be given together.
- No TLS. Put a reverse proxy in front when binding off-loopback.
- Mutually exclusive with `--headless` and `--gui`.
