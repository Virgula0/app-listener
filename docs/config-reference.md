# daemon.conf reference

Every section and directive `/etc/app-listener/daemon.conf` accepts: what it does at run time,
its syntax, the values it takes, what the daemon refuses, and an example.

To change the config of a running daemon, use `sudo app-listener edit-protected config`
([edit-protected](edit-protected.md#editing-the-daemon-config)): the file is self-guarded, and that
command applies the new config by reloading it and puts the old one back if the reload fails. The
install wizard writes the file for you from the [catalog](catalog.md).

- [File format](#file-format)
- [Sections at a glance](#sections-at-a-glance)
- [Binary lines and events](#binary-lines-and-events)
- [`[watch <path>]`](#watch-path): a protected directory or file
  ([`need_encryption`](#need_encryption), [`watch:`](#watch), [`runtime:`](#runtime))
- [`[port tcp/<n>]`](#port-tcpn): a loopback TCP service
- [`[secret_service]`](#secret_service) and [`[secret <name>]`](#secret-name): keyring items
- [`[libraries "<name>"]`](#libraries-name): library trust
  ([`allow_lib`](#allow_lib), [`lib_dir`](#lib_dir), [`lib_binary`](#lib_binary))
- [`[electron_apps]`](#electron_apps): apps a generic Electron may run
- [`[launch_rules]`](#launch_rules): risky arguments and environment
  ([`deny_args`, `warn_args`](#deny_args-warn_args), [`deny_env`, `warn_env`](#deny_env-warn_env),
  [`allow_env`](#allow_env))
- [`[inspectors]`](#inspectors): programs that may inspect protected processes
- [`[settings]`](#settings): daemon-wide settings ([`max_inodes`](#max_inodes))
- [Reloading](#reloading)
- [A complete example](#a-complete-example)

## File format

```ini
# A comment takes a whole line.
[watch /home/alice/.ssh]
need_encryption: true
/usr/bin/ssh READ,WRITE,STAT
"/opt/My Tools/bin/agent" READ
```

- **Lines.** One directive per line. Leading and trailing whitespace is ignored, and so are blank
  lines.
- **Comments.** A line whose first non-blank character is `#` is a comment. There are no comments
  at the end of a line: `/usr/bin/ssh READ # client` reads `#` as an event and is refused.
- **Sections.** A line starting with `[` is a section header, and every line after it belongs to
  that section until the next header. A bracketed line that is not one of the headers below is an
  error, so a typo can't merge the binaries after it into the previous section. A line before the
  first header is ignored, with a warning.
- **Paths** are absolute: `~` and environment variables are not expanded. Wrap a path in double
  quotes when it contains a space, so the rest of the line isn't read as its event list. Quotes
  are allowed around any path (headers included); the installer always writes them.
- **Order.** Sections may come in any order. The installer writes the `[watch]` sections first and
  the blocks after them.
- **Errors.** Any refused line stops the daemon from starting, naming the line number. A reload
  with a refused line keeps the running config.

## Sections at a glance

| Header | Guards | Takes |
|---|---|---|
| [`[watch <path>]`](#watch-path) | a directory (recursively) or one file | binary lines, `need_encryption`, `watch:`, `runtime:` |
| [`[narrow <path>]`](#narrow-path) | a subtree of a `[watch]` tree, by a narrower whitelist | binary lines |
| [`[port tcp/<n>]`](#port-tcpn) | a TCP port on loopback | binary lines (`CONNECT`, `MKNOD`) |
| [`[secret_service]`](#secret_service) | (names the keyring daemon and its helpers) | `provider`, `prompter`, `manager` |
| [`[secret <name>]`](#secret-name) | keyring items | `match`, binary lines (`READ`, `WRITE`, `DELETE`) |
| [`[libraries "<name>"]`](#libraries-name) | libraries whitelisted binaries may load | `allow_lib`, `lib_dir`, `lib_binary` |
| [`[electron_apps]`](#electron_apps) | apps a generic Electron may run | app paths |
| [`[launch_rules]`](#launch_rules) | launches of whitelisted binaries | `deny_args`, `warn_args`, `deny_env`, `warn_env`, `allow_env` |
| [`[inspectors]`](#inspectors) | `/proc` access to protected processes | binary paths |
| [`[settings]`](#settings) | (daemon-wide values) | `max_inodes` |

## Binary lines and events

```ini
/usr/bin/ssh-agent
/usr/bin/ssh READ,WRITE,STAT
/usr/bin/git read, stat
"/home/alice/.local/share/Steam/steamapps/common/Proton - Experimental/files/bin/wineserver" READ
```

A binary line whitelists one program: `<binary> [EVENT[,EVENT...]]`.

- **No event list** allows every event the section takes.
- **An event list** allows only those events: every other one is denied to that binary. Events are
  separated by commas, case-insensitive, and may have spaces after the commas.
- `READ`, `WRITE` and `MMAP` imply `OPEN`.
- **Symlinks** are resolved when the config loads, so you may list either a link or its target;
  the log shows `binary symlink resolved: <link> -> <target>`.
- **A binary that can't be read yet** (it sits inside a locked vault, or isn't installed) is
  deferred: it stays denied until the daemon can resolve it, never dropped.
- The identity checked is the file's inode, so renaming another program to a whitelisted name
  grants nothing, and a renamed or copied program is not the listed one.
- **An event list is a least-privilege hint, not a confinement boundary**: a masked binary that
  execs another whitelisted binary escapes its mask.

| Event | Allows |
|---|---|
| `OPEN` | opening a file or directory; executing a program stored in the tree |
| `READ` | reading, listing a directory, `splice`/`sendfile`/`copy_file_range` out of a file |
| `WRITE` | writing, truncating an open file, copying into a file |
| `MMAP` | mapping a file into memory |
| `DELETE` | removing a file or directory, including the guarded directory itself once it is empty |
| `RENAME` | renaming or moving into, inside or out of the tree |
| `SYMLINK` | creating a symlink |
| `HARDLINK` | creating a hard link into, inside or out of the tree |
| `MKDIR` | creating a directory |
| `MKNOD` | creating a file, FIFO or socket; on a port, listening on it |
| `ATTR` | `chmod`, `chown`, timestamps, truncating by path, setting or removing extended attributes |
| `STAT` | `stat`, `readlink`, reading or listing extended attributes |
| `CONNECT` | connecting, sending or receiving on a unix socket bound inside the tree; on a port, connecting to it |

Which events a section takes: `[watch]` all of them, `[port]` `CONNECT` and `MKNOD`, `[secret]`
`READ`, `WRITE` and `DELETE` (with their own meaning there). Any other event is refused.

## `[watch <path>]`

```ini
[watch /home/alice/.ssh]
need_encryption: true
/usr/bin/ssh READ,WRITE,STAT
/usr/bin/ssh-agent
/usr/bin/ssh-add READ,STAT
```

Protects `<path>`: only the binary lines of the section may touch it, everything else gets
`EPERM`. Whitelist only, default deny.

- **Targets.** A directory, guarded recursively, or a single regular file with its own fscrypt
  policy. A symlink, a file with more than one hard link and special files (FIFOs, sockets,
  devices) are refused, since another path to the same inode would share the rule. A refused or
  missing target skips the whole section, with a warning, and the daemon starts without it.
- **No nesting.** No guarded path (a section path, a `watch:` path, a `lib_dir`) may sit inside
  another one, and none may be declared twice: a file follows exactly one whitelist. Guard sibling
  directories instead.
- **Takes** binary lines, [`need_encryption`](#need_encryption), [`watch:`](#watch) and
  [`runtime:`](#runtime). The library directives are accepted here too, in older configs, but
  belong in a [`[libraries]`](#libraries-name) block.

### `need_encryption`

```ini
[watch /home/alice/.gnupg]
need_encryption: false
/usr/bin/gpg
```

`need_encryption: true` (the default) or `need_encryption: false`; the colon is required.

- `true`: the target must already carry an fscrypt policy, or the daemon refuses to start. The
  daemon unlocks it only after its guard attached, and locks it at shutdown while the guard still
  denies ([details](daemon.md#shutdown-and-crashes)).
- `false`: guarded only, never unlocked or locked.

### `watch:`

```ini
[watch /home/alice/.codex]
need_encryption: true
watch: /home/alice/.codex/auth.json
watch: /home/alice/.codex/config.toml
/usr/bin/codex
```

`watch: <path>` (or `watch <path>`) turns the section path into the encryption root and guards
`<path>` as its own tree, sharing the section's binary lines and its lock/unlock lifecycle: one
key unlocks them all. With `watch:` lines, only those paths are guarded, not the rest of the
section path. Use it to guard the token files of a directory an app also keeps other state in,
without guarding the whole directory.

- `<path>` must be strictly inside the section path, after `..` is resolved. A symlink between the
  section path and `<path>` is refused.
- A `watch:` path hidden by a locked vault is checked again after the daemon unlocks it.
- Each path once per section.

### `runtime:`

```ini
[watch /run/user/1000/gnupg]
runtime: true
need_encryption: false
/usr/bin/gpg-agent MKNOD,DELETE,MKDIR,CONNECT,STAT
/usr/bin/gpg MKDIR,CONNECT,STAT
```

Agents listen on sockets in `/run/user/<uid>`, which exists only while that user is logged in.
`runtime: true` (or `runtime true`; no other value) marks a section for one such directory. The
daemon makes it at that user's login and guards it before any of the user's processes start
([details](daemon.md#runtime-sections-at-login)); while the user is logged out the section is
listed but guards nothing.

- The path is one directory directly in `/run/user/<uid>`.
- Never encrypted: `need_encryption: true` is refused (leave it out or write `false`).
- Takes only binary lines: no `watch:` or library directives.
- Give every line an explicit event list, and `DELETE` only to a program that removes its own
  sockets: a program a caller can steer into `rmdir` would let the directory be recreated
  unguarded.

### `[narrow <path>]`

```ini
[watch /home/alice/.gnupg]
need_encryption: true
/usr/bin/gpg
/usr/bin/gpg-agent
/usr/bin/gpgconf READ,STAT,CONNECT

[narrow /home/alice/.gnupg/private-keys-v1.d]
/usr/bin/gpg-agent
```

A `[narrow <path>]` section guards a directory or file inside a tree an earlier `[watch]` section
guards, by its own binary lines instead of the tree's. The files below `<path>` follow only the
narrow's lines: above, `gpg` may read the whole keyring but not the secret keys, which only
`gpg-agent` handles. A binary that can be made to print what it reads (`gpg --enarmor`, `-e -r`)
can't then print what it never reads.

- Comes after the section of its tree, and `<path>` is strictly inside one of that section's
  guarded paths (the section path, or a `watch:` path). A lib_dir, runtime, port or secret section
  can't be narrowed, nor a narrow (one level).
- Only takes access away: each line must be a line of the parent section, spelled the same, with
  events the parent's line also grants. When the parent's line has an event list, the narrow's
  line needs one too.
- Takes only binary lines. Its encryption and lifecycle are its parent's.
- `<path>` must exist (inside a locked vault, it is checked again after the unlock). A missing
  one is an error, never a fallback to the parent's whitelist: create it or remove the section.
- The narrow's root can't be removed or renamed by a binary of the parent only, and moving a file
  between the narrow and the rest of the tree is the daemon's alone (as between two resources).
- `edit-protected` of the parent widens its narrows too. `edit-protected forward` grants only the
  resource named, and `edit-protected config record` refuses a tree with a narrow.
- `install catalog refresh` regenerates the catalog's narrows and drops, from a hand-written one,
  the lines its parent no longer has.

### Installer placeholder

```ini
[watch]
path = /home/alice/.ssh
/usr/bin/ssh
```

A bare `[watch]` header followed by `path = <dir>` (or `path=<dir>`) is the form the install
wizard's editor writes before you pick a directory. It means the same as `[watch <dir>]`; a
`[watch]` with no `path` guards nothing and its lines are ignored, with a warning.

## `[port tcp/<n>]`

```ini
[port tcp/8332]
/usr/bin/bitcoind MKNOD
/usr/bin/bitcoin-cli CONNECT
```

A service on a loopback TCP port (`bitcoind` RPC, the `adb` server, `ollama`) answers any local
program that can reach it. This section guards one such port: only its `CONNECT` lines may connect
to it, and only its `MKNOD` lines may listen on it, so another program can't take the port while
the service is stopped and collect the credentials its clients send. A listener made while the
daemon was stopped is judged at its first `accept` instead, and every accepted connection is
judged at each message the program sends or receives on it, so a program that took the port then
(even one already waiting in `accept`) can't serve a client either.

- **Loopback only.** Connections to `127.0.0.0/8`, `::1`, `0.0.0.0` and `::` (which reach this
  host's loopback) are judged; connections to other hosts are not. Even a `MKNOD` program is refused
  when it listens on any other address (`0.0.0.0`, `::`, a LAN address), since clients could then
  reach it through an address nobody judges. Configure the service to listen on `127.0.0.1` or `::1`.
- **TCP only, this host's network namespace only.** A container's own loopback is another host's.
  An SMC socket (`AF_SMC`, or `IPPROTO_SMC`) is refused the port to everyone, its `MKNOD` lines
  included: it listens through a TCP socket of the kernel's own that no hook sees. The journal
  names it `dest=smc-socket`.
- **The connection is guarded, not just the connect.** A connection handed to another program
  (inherited, passed over a unix socket) is judged again for every message that program sends or
  receives, as for unix sockets. A `sendto` on it that names another guarded port is refused: the
  data would still go to the port it is connected to.
- The daemon keeps one root-owned anchor file per port in `/run/app-listener/ports` (`tcp-<n>`); its
  events name that path. Removing a port section and reloading stops guarding the port.

Syntax:

- `<n>` is 1 to 65535, after `tcp/`: TCP only.
- Takes only binary lines, with `CONNECT`, `MKNOD` or no event list (both). `need_encryption:
  true` and every other directive are refused.
- One section per port.
- The daemon refuses to start with a port in the ephemeral range (`net.ipv4.ip_local_port_range`,
  by default 32768-60999), where the kernel could hand the port to any program, unless it is listed
  in `net.ipv4.ip_local_reserved_ports`.

## `[secret_service]`

```ini
[secret_service]
provider /usr/bin/gnome-keyring-daemon
prompter /usr/libexec/gcr-prompter
manager /usr/bin/seahorse
```

Apps such as `gh` keep their token in the Secret Service (the keyring), which any of your programs
can ask for over the session bus. With a `[secret_service]` block the keyring reaches the bus only
through the daemon: every call that reads, writes or deletes an item is judged before the keyring
sees it ([how](daemon.md#the-secret-service-gate)), against the [`[secret]`](#secret-name) scopes.

The block names the keyring daemon and the programs around it. Each line is
`<keyword> <absolute path>`, once per path.

- `provider`: the keyring daemon. Required as soon as the config has a `[secret]` section, a
  `prompter` or a `manager`. gnome-keyring is the one supported.
- `prompter`: the program that shows the keyring's password dialogs.
- `manager`: a program that may manage the keyring as a whole (delete a collection, set an alias,
  change the master password).
- Prompter and manager lines must name root-owned binaries.

## `[secret <name>]`

```ini
[secret gh]
match service=gh:*
/usr/bin/gh

[secret mytool]
match application=mytool
match service=mytool-* user=alice
/usr/local/bin/mytool READ
/usr/local/bin/mytool-login
```

A scope: the keyring items it covers and the binaries that may use them, judged through the
[`[secret_service]`](#secret_service) gate.

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

Syntax:

- **`<name>`**: lowercase letters, digits, `.`, `_` and `-`, starting with a letter or digit, at
  most 64 characters. One section per name.
- **`match key=value ...`**: covers an item whose attributes have every listed value. A value
  ending in `*` matches any value starting with the rest; `*` anywhere else, a bare `*`, an empty
  key or value and a repeated key are refused. A section needs at least one `match` line and covers
  an item any of them matches. `secret-tool search --all <key> <value>` shows an item's attributes.
- **Binary lines** take `READ` (get the secret), `WRITE` (create or change an item, its secret,
  label or attributes) and `DELETE`, or no event list for all three.
- Takes nothing else: `need_encryption: true`, `watch:`, `runtime:` and library directives are
  refused.

## `[libraries "<name>"]`

```ini
[libraries "Steam (alice)"]
allow_lib /home/alice/.local/share/Steam/plugin.so
lib_dir /home/alice/.local/share/Steam/steamrt64
lib_binary /home/alice/.local/share/Steam/ubuntu12_32/steam
```

A whitelisted binary may only map libraries the daemon trusts, which is what defeats
`LD_PRELOAD` and library injection. Three things are trusted without any configuration:

1. **Root-owned system libraries**: any `.so` owned by root, not world-writable, in a root-owned
   directory (covers `/usr/lib`, `/opt`, ... and survives package upgrades).
2. **Anything inside a guarded tree**: only whitelisted binaries can write there. A tree its own
   writer moved aside stops counting as soon as its guard lets go of it.
3. **Code a whitelisted process generated itself**, such as a GPU driver's JIT output.

Anything else goes in a `[libraries]` block, one per application. The header is
`[libraries]`, `[libraries <name>]` or `[libraries "<name>"]`; the name is a label the installer
uses to find the block again. A block takes only the three directives below (each also written
`<directive>: <path>`); a binary line or anything else is refused.

### `allow_lib`

```ini
allow_lib /home/alice/.local/share/Steam/plugin.so
```

Trusts one library file at a fixed path that an app `dlopen`s, for every whitelisted binary of
every section, and write-protects it. The trust is daemon-wide, which is why it lives in its own
block rather than under one watch section. A library a binary links against (its `DT_NEEDED`
dependencies) is found by the daemon and needs no line. A file that can't be read yet is deferred
like a binary line.

### `lib_dir`

```ini
lib_dir /home/alice/.local/share/Steam/steamrt64
```

For trees whose files can't be listed (Steam's runtimes are rebuilt under random names at each
launch). Guards a library directory, recursively and read-only: everyone may still read it, but
only the block's `lib_binary` programs may write, rename, delete or truncate inside it, which is
what makes loading from it safe. Never encrypted. Each directory once per block. A directory that
doesn't exist is dropped with a warning.

Moving a file from one guarded tree into another is refused to every process but the daemon,
with one exception: two `lib_dir` trees with exactly the same writers, where those writers may
rename a file (not a directory) from one into the other. That is how Steam installs an update from
its staging dir `package/` into `linux32/`, `ubuntu12_32/` and the other runtime dirs. The moved
file stays protected by the tree it lands in. A refused move is logged with
`resource=cross-resource`.

### `lib_binary`

```ini
lib_binary /home/alice/.local/share/Steam/ubuntu12_32/steam
```

A program that may write this block's `lib_dir` trees, and nothing else (no other block, no watch
section). An event list is refused.

## `[electron_apps]`

```ini
[electron_apps]
/usr/lib/signal-desktop/app.asar
/usr/lib/obsidian/app.asar
```

An Electron with `resources/default_app.asar` next to it (Arch's `/usr/lib/electronNN/electron`,
which wrapper scripts for `signal-desktop`, `obsidian`, ... start) runs whatever app its first
non-switch argument names. Whitelisted, it may be started only with an app listed here.

- Any other app, a relative path, default_app's `-r`/`-i`/`-repl`, or `--app=<file>` counts as a
  risky launch (`op=LAUNCH`, see [Binary trust](binary-trust.md#launch-scan)). No app at all
  (Chromium's child processes) is fine.
- A bundled Electron (Discord, VS Code: `resources/app.asar` or `resources/app/`, no
  `default_app.asar`) ignores its arguments when choosing the app and is not judged this way.

One path per line:

- An absolute path, written exactly as the wrapper script passes it, at most 43 bytes; at most 32
  entries, each once.
- It must name a root-owned file (an `.asar` or a JS entry point; a directory is refused) in
  root-owned directories on a filesystem mounted without `nosuid`; otherwise the daemon logs
  `electron app ... refused` and the entry grants nothing.
- The list is daemon-wide: an Electron whitelisted in two sections may run either section's app.

## `[launch_rules]`

```ini
[launch_rules]
deny_args /usr/bin/gh repo delete
warn_args /usr/bin/kubectl config view --raw
deny_env /usr/bin/kubectl HTTPS_PROXY
warn_env /usr/bin/terraform TF_LOG
allow_env SSL_CERT_FILE=/etc/corp/ca-bundle.pem
```

A whitelisted CLI that prints or spends its secret on request (`gh auth token`) hands it to whoever
runs it. A launch rule marks such a launch: the process is refused every guarded secret, as a risky
runtime launch is ([Binary trust](binary-trust.md#launch-rules), which lists the catalog's rules).
These rules come on top of the catalog's, and can only refuse more: a rule is a tripwire, never a
boundary (see [Limitations](limitations.md#whitelisted-programs-that-disclose-on-request)).

The `<binary>` of a rule is an absolute path
of a whitelisted line, as written or its symlink's target; a rule naming any other path does
nothing and is logged. Each rule once.

### `deny_args`, `warn_args`

```ini
deny_args /usr/bin/gh auth token
warn_args /usr/bin/kubectl config view --raw
```

`deny_args <binary> <word>...` marks a launch whose arguments hold the 1 to 4 words in this order,
anything between them. `warn_args` only logs the match (`ARGV comm=... word=...`).

- A word is at most 43 bytes (an abbreviation spelled out in full).
- A word starting with `-` is a flag, compared up to `=`: write it without a value (`--raw`, not
  `--raw=true`, which is refused).
- Words are compared byte for byte: `--Raw` is another word, so a rule fits only a program whose
  parser is case-sensitive too. Some catalog rules ignore case or take either dash count where
  their program's parser does (Firefox's flags), each rule on its own.
- `--name[rest]` matches every abbreviation from `--name` to `--namerest`, for a program that runs
  any unambiguous prefix of a long option (gpg, getopt_long): `deny_args /usr/bin/gpg
  --export-se[cret-keys]` refuses `--export-se`, `--export-sec`, ..., `--export-secret-keys`. Only a
  flag abbreviates, and `--name` must hold more than the dashes.
- At most 16 argument rules per binary, counting the catalog's for its name and every line naming
  any of its spellings (a symlink and its target). A config over that is refused: the daemon won't
  start on it, and a reload or `edit-protected config` keeps the running one.

### `deny_env`, `warn_env`

```ini
deny_env /usr/bin/kubectl HTTPS_PROXY
```

`deny_env <binary> <NAME>` marks a launch with the variable `NAME` set, whatever its value;
`warn_env` only logs it. `NAME` is one variable: letters, digits and `_`, not starting with a digit,
at most 43 bytes.

### `allow_env`

```ini
allow_env SSL_CERT_FILE=/etc/corp/ca-bundle.pem
```

`allow_env <NAME>=<path>`: no env rule on `NAME` (built-in or yours) matches when the value is
exactly `<path>`, for every binary. Use it for a corporate CA bundle passed in `SSL_CERT_FILE`.

- `<path>` is one clean absolute path, without `:` or whitespace; the whole `NAME=<path>` is at
  most 63 bytes.
- The daemon accepts it only when root decides what the path holds: every directory on the way and
  the file (or, for a directory, everything in it) root-owned and writable by no one else, outside
  user homes, on a mount without `nosuid`. Otherwise it logs `allow_env ... refused` and launches
  setting it stay marked.

## `[inspectors]`

```ini
[inspectors]
/usr/lib/xdg-desktop-portal
```

A process that read a protected file is *tainted*: no program outside its resource's whitelist
may read its memory or even its `/proc/<pid>` metadata. The programs listed here may read that
metadata (`environ`, `fd`, `maps`, `root`, `exe`) of **every** protected process, but never its
memory (`ptrace`, `process_vm_readv`, `/proc/<pid>/mem`) and never a protected file. Inspectors are
protected like whitelisted binaries (no `LD_PRELOAD`), and attaching to or tracing one needs a
program **every** section whitelists.

**What an inspector can see** is everything the kernel gates at `PTRACE_MODE_READ`: each
protected process's environment variables (tokens passed via env included), open-file paths and
memory map, the contents of its *unguarded* open files reachable through `/proc/<pid>/fd/N`
(pipes, memfds, deleted files), and `perf_event_open` samples of its user stack. Protected files
stay denied through `/proc/<pid>/fd`.

One absolute path per line, nothing after it.

- **Without the block**, the default is `xdg-desktop-portal` (`/usr/lib/...` or
  `/usr/libexec/...`, whichever is installed): it opens `/proc/<pid>/root` of every
  screen-sharing caller and otherwise refuses the request, so Meet or Discord sharing fails with no
  picker.
- **An empty block** removes the default.
- At most 16 entries, each once. Only a root-owned binary in root-owned directories on a
  filesystem mounted without `nosuid` is accepted; otherwise the daemon logs `inspector ...
  refused`.
- Never list a shell, an interpreter, `cat`, `ps`, `perf`, `gdb` or anything that does what its
  caller asks.

## `[settings]`

```ini
[settings]
max_inodes: 65536
```

Daemon-wide `key: value` lines (the colon is required). An unknown or repeated key, or a value out
of range, is refused.

### `max_inodes`

Sizes `guard_inodes`, the kernel map holding one row per inode of every guarded tree: 32768 to
4194304, default 262144 (about 22 MiB of preallocated kernel memory). A reload can't resize the
live map: a changed value takes effect at the next restart, and the reload logs so.

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

## A complete example

```ini
# SSH keys: encrypted at rest, only the OpenSSH tools.
[watch /home/alice/.ssh]
need_encryption: true
/usr/bin/ssh READ,WRITE,STAT
/usr/bin/ssh-add READ,STAT
/usr/bin/ssh-agent

# Codex: two files of the directory, one vault.
[watch /home/alice/.codex]
need_encryption: true
watch: /home/alice/.codex/auth.json
watch: /home/alice/.codex/config.toml
/usr/bin/codex

# The ssh-agent socket directory, made at alice's login.
[watch /run/user/1000/ssh-agent]
runtime: true
need_encryption: false
/usr/bin/ssh-agent MKNOD,DELETE,CONNECT,STAT
/usr/bin/ssh CONNECT,STAT
/usr/bin/ssh-add CONNECT,STAT

# bitcoind's RPC port.
[port tcp/8332]
/usr/bin/bitcoind MKNOD
/usr/bin/bitcoin-cli CONNECT

# The gh token in the keyring.
[secret_service]
provider /usr/bin/gnome-keyring-daemon
prompter /usr/libexec/gcr-prompter

[secret gh]
match service=gh:*
/usr/bin/gh

# Steam's runtimes, writable only by Steam.
[libraries "Steam (alice)"]
lib_dir /home/alice/.local/share/Steam/steamrt64
lib_binary /home/alice/.local/share/Steam/ubuntu12_32/steam

[electron_apps]
/usr/lib/signal-desktop/app.asar

[launch_rules]
deny_args /usr/bin/gh repo delete
allow_env SSL_CERT_FILE=/etc/corp/ca-bundle.pem

[inspectors]
/usr/lib/xdg-desktop-portal

[settings]
max_inodes: 131072
```

The fully commented template the daemon falls back to is
[`daemon-samples/daemon.conf`](../daemon-samples/daemon.conf).
