# Daemon configuration

The daemon reads `/etc/app-listener/daemon.conf`. The install wizard writes it for you from the
[catalog](catalog.md). This page explains what each kind of section does at run time; the syntax
of every section and directive, with examples, is in the [daemon.conf reference](config-reference.md).

To change the config of a running daemon, prefer `sudo app-listener edit-protected config`
([edit-protected](edit-protected.md#editing-the-daemon-config)): the file is self-guarded, and that
path validates the new config by reloading it and puts the old one back if the reload fails.

## Watch sections

A [`[watch <path>]`](config-reference.md#watch-path) section protects one directory or file.

- **Whitelist only, default deny.** Only the listed binaries may access the tree; everything else
  gets `EPERM`. Identity is by inode, so a renamed or copied program is not the listed one.
- **Encryption.** A resource with [`need_encryption: true`](config-reference.md#need_encryption)
  (the default) must already carry an fscrypt policy. The daemon unlocks it only after its guard
  attached, and locks it at shutdown while the guard still denies
  ([details](daemon.md#shutdown-and-crashes)).
- **Event masks** ([events](config-reference.md#binary-lines-and-events)) are a least-privilege
  hint, **not a confinement boundary**: a masked binary that execs another whitelisted binary
  escapes its mask.

### Several guarded trees, one vault

[`watch:`](config-reference.md#watch) lines inside a section turn the section path into the
**encryption root** and each `watch:` path into a separately guarded tree. The trees share the
section's whitelist and its lock/unlock lifecycle: one key unlocks them all. Without `watch:` lines
the section path itself is the guarded tree. Use it to guard the token files of a directory an app
also keeps other state in, without guarding the whole directory.

### Runtime sections

Agents listen on sockets in `/run/user/<uid>`, which exists only while that user is logged in. A
[`runtime: true`](config-reference.md#runtime) section names one directory directly in it. The
daemon makes the directory at the user's login and guards it before any of the user's processes
start ([details](daemon.md#runtime-sections-at-login)); while the user is logged out the section is
listed but guards nothing.

## Port sections

A service on a loopback TCP port (`bitcoind` RPC, the `adb` server, `ollama`) answers any local
program that can reach it. A [`[port tcp/<n>]`](config-reference.md#port-tcpn) section guards one
such port: only its `CONNECT` lines may connect to it, and only its `MKNOD` lines may listen on it,
so another program can't take the port while the service is stopped and collect the credentials
its clients send. A listener made while the daemon was stopped is judged at its first `accept`
instead, and every accepted connection is judged at each message the program sends or receives on
it, so a program that took the port then (even one already waiting in `accept`) can't serve a
client either.

- **Loopback only.** Connections to `127.0.0.0/8`, `::1`, `0.0.0.0` and `::` (which reach this
  host's loopback) are judged; connections to other hosts are not. Even a `MKNOD` program is refused
  when it listens on any other address (`0.0.0.0`, `::`, a LAN address), since clients could then
  reach it through an address nobody judges. Configure the service to listen on `127.0.0.1` or `::1`.
- **TCP only, this host's network namespace only.** A container's own loopback is another host's.
  An SMC socket (`AF_SMC`, or `IPPROTO_SMC`) is refused the port to everyone, its `MKNOD` lines
  included: it listens through a TCP socket of the kernel's own that no hook sees. The journal
  names it `dest=smc-socket`.
- **Not in the ephemeral range**, where the kernel could hand the port to any program.
- **The connection is guarded, not just the connect.** A connection handed to another program
  (inherited, passed over a unix socket) is judged again for every message that program sends or
  receives, as for unix sockets. A `sendto` on it that names another guarded port is refused: the
  data would still go to the port it is connected to.

The daemon keeps one root-owned anchor file per port in `/run/app-listener/ports` (`tcp-<n>`); its
events name that path. Removing a port section and reloading stops guarding the port.

## Secret Service sections

Apps such as `gh` keep their token in the Secret Service (the keyring), which any of your programs
can ask for over the session bus. With a [`[secret_service]`](config-reference.md#secret_service)
block the keyring reaches the bus only through the daemon: every call that reads, writes or deletes
an item is judged before the keyring sees it ([how](daemon.md#the-secret-service-gate)).

A [`[secret <name>]`](config-reference.md#secret-name) section is a scope: the items its `match`
lines cover and the binaries that may use them.

- An item that several scopes cover needs every scope's admission.
- Changing an item's attributes or its type (the keyring reads the type back as `xdg:schema`) is
  judged under its old and new ones, so an item can't be moved into or out of a scope.
- A call that names the same key twice in a dict (two `Attributes` entries) is refused: the keyring
  acts on the first, so the gate couldn't judge what it runs.
- Creating an item that replaces existing ones (the keyring picks them by a subset of their
  attributes) is judged under the scopes of every item it could replace.
- An item no scope covers stays as open as the keyring keeps it.
- The daemon keeps one root-owned anchor file per scope in `/run/app-listener/secret-scopes`; its
  events name that path and the call.

## Library trust

A whitelisted binary may only map libraries the daemon trusts, which is what defeats
`LD_PRELOAD` and library injection. Three things are trusted without any configuration:

1. **Root-owned system libraries**: any `.so` owned by root, not world-writable, in a root-owned
   directory (covers `/usr/lib`, `/opt`, ... and survives package upgrades).
2. **Anything inside a guarded tree**: only whitelisted binaries can write there. A tree its own
   writer moved aside stops counting as soon as its guard lets go of it.
3. **Code a whitelisted process generated itself**, such as a GPU driver's JIT output.

Anything else goes in a [`[libraries "<name>"]`](config-reference.md#libraries-name) block, one per
application:

- [`allow_lib`](config-reference.md#allow_lib) trusts one file. The trust is **daemon-wide**: every
  whitelisted binary may load it, which is why it lives in its own block rather than under one
  watch section.
- [`lib_dir`](config-reference.md#lib_dir) is for trees whose files can't be listed (Steam's
  runtimes are rebuilt under random names at each launch). Only the block's
  [`lib_binary`](config-reference.md#lib_binary) programs may change it, which is what makes loading
  from it safe; they get no access to any watch section.

Moving a file from one guarded tree into another is refused to every process but the daemon,
with one exception: two `lib_dir` trees with exactly the same writers, where those writers may
rename a file (not a directory) from one into the other. That is how Steam installs an update from
its staging dir `package/` into `linux32/`, `ubuntu12_32/` and the other runtime dirs. The moved
file stays protected by the tree it lands in. A refused move is logged with
`resource=cross-resource`.

## Generic Electron apps

An Electron with `resources/default_app.asar` next to it (Arch's `/usr/lib/electronNN/electron`,
which wrapper scripts for `signal-desktop`, `obsidian`, ... start) runs whatever app its first
non-switch argument names. Whitelisted, it may be started only with an app listed in
[`[electron_apps]`](config-reference.md#electron_apps).

- Any other app, a relative path, default_app's `-r`/`-i`/`-repl`, or `--app=<file>` counts as a
  risky launch (`op=LAUNCH`, see [Binary trust](binary-trust.md#launch-scan)). No app at all
  (Chromium's child processes) is fine.
- The list is daemon-wide: an Electron whitelisted in two sections may run either section's app.

A bundled Electron (Discord, VS Code: `resources/app.asar` or `resources/app/`, no
`default_app.asar`) ignores its arguments when choosing the app and is not judged this way.

## Launch rules

A whitelisted CLI that prints or spends its secret on request (`gh auth token`) hands it to whoever
runs it. Launch rules mark such a launch: the process is refused every guarded secret, as a risky
runtime launch is ([Binary trust](binary-trust.md#launch-rules), which lists the built-in rules).
Add your own in [`[launch_rules]`](config-reference.md#launch_rules). A rule can only refuse more:
it is a tripwire, never a boundary (see
[Limitations](limitations.md#whitelisted-programs-that-disclose-on-request)).

## Process inspectors

A process that read a protected file is *tainted*: no program outside its resource's whitelist
may read its memory or even its `/proc/<pid>` metadata. Programs listed under
[`[inspectors]`](config-reference.md#inspectors) may read that metadata (`environ`, `fd`, `maps`,
`root`, `exe`) of **every** protected process, but never its memory (`ptrace`, `process_vm_readv`,
`/proc/<pid>/mem`) and never a protected file.

The default inspector is `xdg-desktop-portal`: it opens `/proc/<pid>/root` of every
screen-sharing caller and otherwise refuses the request, so Meet or Discord sharing fails with no
picker. Inspectors are protected like whitelisted binaries (no `LD_PRELOAD`), and attaching to or
tracing one needs a program **every** section whitelists.

**What an inspector can see** is everything the kernel gates at `PTRACE_MODE_READ`: each
protected process's environment variables (tokens passed via env included), open-file paths and
memory map, the contents of its *unguarded* open files reachable through `/proc/<pid>/fd/N`
(pipes, memfds, deleted files), and `perf_event_open` samples of its user stack. Protected files
stay denied through `/proc/<pid>/fd`.

## Daemon settings

[`max_inodes`](config-reference.md#max_inodes) in `[settings]` sizes `guard_inodes`, the kernel
map holding one row per inode of every guarded tree.

- Only worth lowering on a small-RAM host with few, small resources: 37 resources (browser
  profiles, Steam, Discord, AI agents) use about 100,000 rows. `sudo scripts/kernel-prof.sh` prints
  the running daemon's `guard_inodes_used`.
- When the map passes 75% the daemon logs `guard_inodes holds N of M rows`, after startup, after
  each reload and every 10 minutes. Once full, inodes it can't hold are guarded by the ancestor
  walk alone (`degraded coverage: guard_inodes map full`): raise the value.
- It is never lowered automatically.

## Reloading

`sudo systemctl reload app-listener-daemon` re-reads the file. New guards attach before old ones
detach, and a malformed config keeps the previous one running. A reload never drops a guarded
resource (stop, edit and start for that), except a `lib_dir` whose directory its writer deleted
(an app's old version dir after an update), which it drops as a restart would. Binaries added by
hand are vetted like any other: one that isn't root-placed needs `sudo app-listener trust`
([Binary trust](binary-trust.md#trust)).
