# guard: block file access by program

Denies file operations on a path, deciding by the identity of the program that asks. It runs in
the foreground until you quit: useful for a quick lockdown, for testing a whitelist before putting
it in the [daemon](daemon.md), or for guarding a path for the length of a session.

![guard](../media/guard.gif)

```bash
sudo app-listener guard /secret                    # block everything
sudo app-listener guard /secret -w /usr/bin/cat    # allow only cat
sudo app-listener guard /secret -b /usr/bin/rm     # block only rm
```

## Whitelist and blacklist

- **Whitelist** (`-w`, the default): only the listed binaries get through. Without any `-w`,
  everything is blocked.
- **Blacklist** (`-b`): only the listed binaries are blocked.

`-w` and `-b` are mutually exclusive. Binaries are matched by inode, see
[How it works](how-it-works.md#identity-is-the-executables-inode).

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `<path>` | required | File or directory to guard |
| `-w, --whitelist <binary>` | none | Binaries allowed (repeatable) |
| `-b, --blacklist <binary>` | none | Binaries blocked (repeatable) |
| `-r, --recursive` | `true` | Recurse into subdirectories |
| `-d, --depth <n>` | `0` | Max depth (`0` = unlimited) |
| `-e, --events <list>` | all | Event types to enforce (same set as [monitor](monitor.md#events)) |
| `--headless` | `false` | No TUI; log `GUARD\|` lines to stderr |
| `--serve[=host:port]` | off | Mirror the TUI to a browser |

## Executables inside the guarded tree

Executing a binary is an `OPEN` performed by the *launcher*, and a shell wrapper runs through its
interpreter, which a whitelist deliberately excludes. In whitelist mode, opens for exec are
therefore attributed to the **binary being executed**. That lets whitelisted binaries that live
*inside* the guarded tree (for example Discord under `~/.config/discord`) start from any shell or
wrapper.

The exec fd is never exposed to the launcher. Helper binaries that an app spawns must be
whitelisted explicitly. Blacklist mode always attributes to the launcher.

## Multicall binaries

uutils coreutils is one file behind every applet name; app-listener tells its applets apart, so
`-w /usr/bin/head` admits `head` only. BusyBox, toybox and other unvetted multicall binaries are
refused. Details in [Binary trust](binary-trust.md#multicall-binaries).

## Replaced binaries

A process started before its binary was replaced keeps its access. Re-running the replaced image
afterwards, or a file that later reuses the old inode number, does not.
