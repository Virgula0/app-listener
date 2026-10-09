# Daemon configuration

The daemon reads `/etc/app-listener/daemon.conf`. The install wizard writes it for you from the
[catalog](catalog.md); this page is for reading it or editing it by hand. The fully commented
template is [`daemon-samples/daemon.conf`](../daemon-samples/daemon.conf).

To change the config of a running daemon, prefer `sudo app-listener edit-protected --edit-config`
([edit-protected](edit-protected.md#editing-the-daemon-config)): the file is self-guarded, and that
path validates the new config by reloading it and puts the old one back if the reload fails.

## Watch sections

```ini
[watch /home/alice/.ssh]          # one section per protected directory
need_encryption: true             # default true; false skips the fscrypt lifecycle
/usr/bin/ssh READ,WRITE           # whitelisted binary, restricted to these events
/usr/bin/ssh-agent                # bare path = all events allowed
```

- **Whitelist only, default deny.** Only the listed binaries may access the tree; everything else
  gets `EPERM`. Identity is by inode.
- **Absolute paths only.** `~` is not expanded.
- **Symlinked binaries** are resolved to their real path at load time, so either the link or its
  target may be listed.
- **Quoting.** Wrap a path in double quotes when it contains spaces, so it isn't read as the event
  list: `"/home/alice/Steam/.../Proton - Experimental/files/bin/wineserver"`. The installer always
  quotes the paths it writes.
- **Targets.** A section watches a directory (recursively) or a single regular file that carries
  its own fscrypt policy. Symlinks, files with more than one hard link and special files are
  refused, because an alias would extend the rule to another path to the same inode.
- **No nesting.** Guarded paths may not sit inside one another: a file follows exactly one
  whitelist. Guard sibling directories instead.
- **`need_encryption: true`** (the default) requires the directory to already carry an fscrypt
  policy; otherwise the daemon refuses to start. Use `false` for a resource you only want guarded.

### Several guarded trees, one vault

`watch: <path>` lines inside a section turn the section path into the **encryption root** and each
`watch:` path into a separately guarded tree that shares the section's whitelist and lock/unlock
lifecycle. Without `watch:` lines the section path itself is the guarded tree.

### Per-binary event masks

Valid events: `OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP, ATTR, STAT, MKNOD`.

Unlisted events are denied for that binary. `READ`, `WRITE` and `MMAP` imply `OPEN`.

Masks are a least-privilege hint, **not a confinement boundary**: a masked binary that execs
another whitelisted binary escapes its mask.

## Library trust: `[libraries "<name>"]`

A whitelisted binary may only map libraries the daemon trusts, which is what defeats
`LD_PRELOAD` and library injection. Three things are trusted without any configuration:

1. **Root-owned system libraries**: any `.so` owned by root, not world-writable, in a root-owned
   directory (covers `/usr/lib`, `/opt`, ... and survives package upgrades).
2. **Anything inside a guarded tree**: only whitelisted binaries can write there.
3. **Code a whitelisted process generated itself**, such as a GPU driver's JIT output.

Anything else goes in a `[libraries]` block, one per application, after the watch sections:

```ini
[libraries "Steam (alice)"]
allow_lib /home/alice/.local/share/Steam/plugin.so
lib_dir /home/alice/.local/share/Steam/steamrt64
lib_binary /home/alice/.local/share/Steam/ubuntu12_32/steam
```

| Directive | Meaning |
|---|---|
| `allow_lib <path>` | Trusts one library file at a fixed, user-writable path that the app `dlopen`s. The file also becomes write-protected. Trust is **daemon-wide**: every whitelisted binary may load it. Static ELF dependencies are resolved automatically and need no listing. |
| `lib_dir <path>` | Guards a whole library directory read-only: everyone may still read it, but only the block's `lib_binary` entries may write, rename, delete or truncate inside it, and only they may load from it. Never encrypted. For trees whose files can't be listed (Steam's runtimes are rebuilt under random names at each launch). |
| `lib_binary <path>` | Allows one binary to write this block's `lib_dir` trees, and nothing else (not other blocks, not any watch section). Event masks are rejected here. |

The block name is a label the installer uses to find the block again. Older configs that nest
these directives inside a `[watch]` section are still accepted.

## Generic Electron apps: `[electron_apps]`

An Electron with `resources/default_app.asar` next to it (Arch's `/usr/lib/electronNN/electron`,
which wrapper scripts for `signal-desktop`, `obsidian`, ... start) runs whatever app its first
non-switch argument names. Whitelisted, it may be started only with an app listed here:

```ini
[electron_apps]
/usr/lib/signal-desktop/app.asar
```

- Each entry is an absolute path written exactly as the wrapper passes it, at most 43 bytes, at
  most 32 entries.
- It must name a root-owned file (an `.asar` or a JS entry point; a directory is refused) in
  root-owned directories on a filesystem mounted without `nosuid`; otherwise the daemon logs
  `electron app ... refused`.
- Any other app, a relative path, default_app's `-r`/`-i`/`-repl`, or `--app=<file>` counts as a
  risky launch (`op=LAUNCH`, see [Binary trust](binary-trust.md#launch-scan)). No app at all
  (Chromium's child processes) is fine.
- The list is daemon-wide: an Electron whitelisted in two sections may run either section's app.

A bundled Electron (Discord, VS Code: `resources/app.asar` or `resources/app/`, no
`default_app.asar`) ignores its arguments when choosing the app and is not judged this way.

## Process inspectors

A process that read a protected file is *tainted*: no program outside its resource's whitelist
may read its memory or even its `/proc/<pid>` metadata. Programs listed under `[inspectors]` may
read that metadata (`environ`, `fd`, `maps`, `root`, `exe`) of **every** protected process, but
never its memory (`ptrace`, `process_vm_readv`, `/proc/<pid>/mem`) and never a protected file.

```ini
[inspectors]
/usr/lib/xdg-desktop-portal
```

- **Default without the block:** `xdg-desktop-portal` (`/usr/lib/...` or `/usr/libexec/...`;
  a path not installed is skipped). It opens `/proc/<pid>/root` of every screen-sharing caller
  and otherwise refuses the request, so Meet or Discord sharing fails with no picker.
- **An empty `[inspectors]` block** disables the default.
- At most 16 entries. Only a root-owned binary in root-owned directories on a filesystem mounted
  without `nosuid` is accepted (else `inspector ... refused`). Inspectors are protected like
  whitelisted binaries (no `LD_PRELOAD`), and attaching to or tracing one needs a program
  **every** section whitelists.

**What an inspector can see** is everything the kernel gates at `PTRACE_MODE_READ`: each
protected process's environment variables (tokens passed via env included), open-file paths and
memory map, the contents of its *unguarded* open files reachable through `/proc/<pid>/fd/N`
(pipes, memfds, deleted files), and `perf_event_open` samples of its user stack. Protected files
stay denied through `/proc/<pid>/fd`. Never list a shell, an interpreter, `cat`, `ps`, `perf`,
`gdb` or anything that does what its caller asks.

## Reloading

`sudo systemctl reload app-listener-daemon` re-reads the file. New guards attach before old ones
detach, and a malformed config keeps the previous one running. Binaries added by hand are vetted
like any other: one that isn't root-placed needs `sudo app-listener trust-binaries`
([Binary trust](binary-trust.md#trust-binaries)).
