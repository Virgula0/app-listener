# daemon.conf reference

Every section and directive `/etc/app-listener/daemon.conf` accepts, with its syntax, the values
it takes, what the daemon refuses, and an example. What each section does at run time (how a
whitelist is enforced, what library trust covers, how the Secret Service gate works) is in
[Daemon configuration](daemon-config.md).

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
  grants nothing.

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
`EPERM`.

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
  daemon unlocks it after the guard attaches and locks it again at shutdown.
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
`<path>` as its own tree, sharing the section's binary lines and its lock/unlock lifecycle. With
`watch:` lines, only those paths are guarded, not the rest of the section path.

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

`runtime: true` (or `runtime true`; no other value) marks a section for a directory that exists
only while its user is logged in. The daemon makes it at that user's login and guards it before
any of the user's processes start ([details](daemon.md#runtime-sections-at-login)).

- The path is one directory directly in `/run/user/<uid>`.
- Never encrypted: `need_encryption: true` is refused (leave it out or write `false`).
- Takes only binary lines: no `watch:` or library directives.
- Give every line an explicit event list, and `DELETE` only to a program that removes its own
  sockets: a program a caller can steer into `rmdir` would let the directory be recreated
  unguarded.

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

Guards a service on a loopback TCP port: `CONNECT` lines may connect to it, `MKNOD` lines may
listen on it ([what is judged](daemon-config.md#port-sections)).

- `<n>` is 1 to 65535, after `tcp/`: TCP only.
- Takes only binary lines, with `CONNECT`, `MKNOD` or no event list (both). `need_encryption:
  true` and every other directive are refused.
- One section per port.
- The daemon refuses to start with a port in the ephemeral range (`net.ipv4.ip_local_port_range`,
  by default 32768-60999) unless it is listed in `net.ipv4.ip_local_reserved_ports`.

## `[secret_service]`

```ini
[secret_service]
provider /usr/bin/gnome-keyring-daemon
prompter /usr/libexec/gcr-prompter
manager /usr/bin/seahorse
```

Names the keyring daemon and the programs around it. Each line is `<keyword> <absolute path>`,
once per path.

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

A scope: the keyring items it covers and the binaries that may use them
([how calls are judged](daemon-config.md#secret-service-sections)).

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

Libraries the daemon trusts beyond the ones it trusts by itself
([library trust](daemon-config.md#library-trust)). One block per application. The header is
`[libraries]`, `[libraries <name>]` or `[libraries "<name>"]`; the name is a label the installer
uses to find the block again. A block takes only the three directives below (each also written
`<directive>: <path>`); a binary line or anything else is refused.

### `allow_lib`

```ini
allow_lib /home/alice/.local/share/Steam/plugin.so
```

Trusts one library file at a fixed path that an app `dlopen`s, for every whitelisted binary of
every section, and write-protects it. A library a binary links against (its `DT_NEEDED`
dependencies) is found by the daemon and needs no line. A file that can't be read yet is deferred
like a binary line.

### `lib_dir`

```ini
lib_dir /home/alice/.local/share/Steam/steamrt64
```

Guards a library directory, recursively and read-only: everyone may still read it, but only the
block's `lib_binary` programs may write, rename, delete or truncate inside it, and only they may
load from it. Never encrypted. Each directory once per block. A directory that doesn't exist is
dropped with a warning.

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

The apps a whitelisted generic Electron (one with `resources/default_app.asar`) may be started
with ([why](daemon-config.md#generic-electron-apps)). One path per line.

- An absolute path, written exactly as the wrapper script passes it, at most 43 bytes; at most 32
  entries, each once.
- It must name a root-owned file (an `.asar` or a JS entry point; a directory is refused) in
  root-owned directories on a filesystem mounted without `nosuid`; otherwise the daemon logs
  `electron app ... refused` and the entry grants nothing.
- The list is daemon-wide.

## `[launch_rules]`

```ini
[launch_rules]
deny_args /usr/bin/gh repo delete
warn_args /usr/bin/kubectl config view --raw
deny_env /usr/bin/kubectl HTTPS_PROXY
warn_env /usr/bin/terraform TF_LOG
allow_env SSL_CERT_FILE=/etc/corp/ca-bundle.pem
```

Marks risky launches of whitelisted binaries, on top of the built-in rules
([what a mark does](binary-trust.md#launch-rules)). The `<binary>` of a rule is an absolute path
of a whitelisted line, as written or its symlink's target; a rule naming any other path does
nothing and is logged. Each rule once.

### `deny_args`, `warn_args`

```ini
deny_args /usr/bin/gh auth token
warn_args /usr/bin/kubectl config view --raw
```

`deny_args <binary> <word>...` marks a launch whose arguments hold the 1 to 4 words in this order,
anything between them. `warn_args` only logs the match (`ARGV comm=... word=...`).

- A word is at most 43 bytes.
- A word starting with `-` is a flag, compared up to `=`: write it without a value (`--raw`, not
  `--raw=true`, which is refused).
- At most 16 argument rules per binary, built-ins included.

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

Programs that may read the `/proc/<pid>` metadata (`environ`, `fd`, `maps`, `root`, `exe`) of
every protected process, never its memory or a protected file
([what they can see](daemon-config.md#process-inspectors)). One absolute path per line, nothing
after it.

- **Without the block**, the default is `xdg-desktop-portal` (`/usr/lib/...` or
  `/usr/libexec/...`, whichever is installed).
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
live map: a changed value takes effect at the next restart, and the reload logs so. When to change
it: [Daemon settings](daemon-config.md#daemon-settings).

## A complete example

```ini
# SSH keys: encrypted at rest, only the OpenSSH tools.
[watch /home/alice/.ssh]
need_encryption: true
/usr/bin/ssh READ,WRITE,STAT
/usr/bin/ssh-add READ,STAT
/usr/bin/ssh-agent
/usr/bin/scp READ,STAT

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
