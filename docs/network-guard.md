# network-guard: block network access by program

Denies socket operations by binary identity through the LSM `socket_connect`, `socket_bind`,
`socket_listen`, `socket_sendmsg` and `socket_recvmsg` hooks. Two typical uses: cut a single
tool off the network (`-b`), or let only a browser and the system's own network daemons out
(`-w`).

```bash
sudo app-listener network-guard -b /usr/bin/curl -e CONNECT,SEND
sudo app-listener network-guard -w /usr/lib/firefox/firefox --auto-infra
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-b, --blacklist <binary>` | none | Block network ops for these binaries only (exclusive with `-w`) |
| `-w, --whitelist <binary>` | none | Block network ops for **all** binaries except these (default deny) |
| `--auto-infra` | `false` | Auto-allow running infrastructure daemons (systemd-resolved, NetworkManager, ...). Without it DNS breaks for everyone in whitelist mode |
| `--unsafe` | `false` | Also block AF_UNIX (X11, D-Bus, systemd); may break the desktop |
| `--no-throttle` | `false` | Disable rate limiting (default: 1 event per type, verdict and process name per 250 ms) |
| `-e, --events <list>` | all | `CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS,BIND,LISTEN` |
| `--headless` | `false` | No TUI; log `NETGUARD\|` lines to stderr |
| `--serve[=host:port]` | off | Mirror the TUI to a browser |

**Pick the real executable, not a wrapper.** Identity is the exe inode, so whitelisting a
`#!/bin/sh` wrapper matches nothing:

```bash
ls -l /proc/$(pgrep -n firefox)/exe    # running process: the real binary
readlink -f /usr/bin/firefox           # follows symlinks, NOT shell wrappers
```

## Code integrity (`-w`)

A whitelist row admits the binary's **code**, not whatever ends up running inside its process. A
whitelisted process loses network access until it execs again when:

- it was started with `LD_PRELOAD` or `LD_AUDIT` (`NETGUARD|...|true|ld-preload`), or
- it maps executable code its binary doesn't vouch for, such as a `dlopen` or `LD_LIBRARY_PATH`
  library (`|code-suspect`).

System libraries (root-owned, in a root-owned directory, on a mount without `nosuid`) are fine,
and so is a library owned by the binary's own non-root owner that nobody else can modify. Forked
children inherit the mark.

Other processes can't inject into a whitelisted one either. Attaching to it, writing into its
memory, or opening its `/proc/<pid>/mem` (`ptrace`, `process_vm_writev`, `pidfd_getfd`) is refused
to every process the whitelist doesn't admit (`NETGUARD|PTRACE|...`), and so is exec'ing a
whitelisted binary under such a tracer (`NETGUARD|TRACED_EXEC|...`).

Processes already running when the guard starts are judged once from `/proc` (environment,
tracer, mapped libraries), but code written into their memory earlier can't be seen. **Restart
whitelisted apps after starting the guard.**

### What code integrity cannot stop

A whitelisted interpreter (`python3`, `node`, `bash`) runs whatever script its caller passes, and
a JIT or interpreter inside a whitelisted app runs code from anonymous memory, which is never
judged. Whitelist the application, not a general-purpose runtime.

## Replaced binaries

As in [guard](guard.md#replaced-binaries): a process started before its binary was replaced keeps
its access; re-running the replaced image, or a file that later reuses its inode number, does not.

## Multicall binaries

`-w` and `-b` refuse multicall binaries that can't be told apart. See
[Binary trust](binary-trust.md#multicall-binaries).
