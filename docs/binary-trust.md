# Binary trust

A whitelist entry is only as good as the file it names. If anyone can rewrite a whitelisted
binary, or feed it code at launch, the whitelist protects nothing. This page covers how the daemon
keeps whitelisted binaries trustworthy: across updates, across restarts, and at exec time.

## Binary replacement

Binaries get replaced all the time (package upgrades, self-updating apps). The rules:

- A whitelisted binary at a **user-writable path** may be modified only by a binary of a resource
  that whitelists it, never by a shell, interpreter, `git`, `cp`, ...
- A **new inode at a whitelisted path** is re-admitted with no reload when either:
  - one of those updaters created it, or
  - it is a root-owned file in a root-owned directory, on a filesystem mounted without `nosuid`
    (a package upgrade).
- Anything else stays denied (`not re-admitted`), at the next reload or start too, until you
  confirm it with [`trust`](#trust).
- A new root-owned match of the catalog's absolute patterns (`/opt/*/jbr/bin/java`) is admitted
  the same way.

A process started before its binary was replaced keeps its access, through a reload too, without
a gap: the reload's new guard holds the old binary's entry before it takes over the directory.
Re-running the replaced image afterwards (an fd held across the update, `/proc/<pid>/exe`) does
not, and neither does a file that later reuses the old inode number.

Every resource whitelisting the path judges the new inode, even one whose app hasn't run since. A
self-updating app such as Steam restarts as soon as it has written its update, which frees the old
binary the update was recorded against; before forgetting a freed binary, the daemon re-syncs
every guard, so none of them is left judging the update after its updater is gone. A freed binary
is also kept for at least one 30 second prune cycle, and while a catalog refresh is pending or
running (up to 5 minutes), so an app that deletes its old version right after writing the new one
(Discord) still has its update judged against the updater that wrote it.

## Binary ledger

Nothing guards a whitelisted path while the daemon is stopped, booting or being reinstalled. So
every start and reload checks each whitelisted binary against the content hash recorded for its
line in `/etc/app-listener/binaries.db` (for a symlink line, the link path).

A binary is admitted when:

- only root could have placed it (root-owned file and directories, outside every user home, on a
  mount without `nosuid`), or
- its hash matches the ledger, or
- its resource's updater created it while the daemon ran (recorded automatically, as are live
  replacements). A new version directory counts too when the old version is already gone by the
  time the refresh adds the new one: the updater is judged by the resource, not by a path.

Anything else, for example a binary changed while the daemon was not running in a home or a
user-owned directory, or a new catalog match planted meanwhile, is refused:

```
DAEMON binary-unconfirmed resource=... path=... reason=new|changed sha256=...
```

A refused binary gets no access, and no updater or reserved-name writer rights either.

The first start after upgrading from a version without the ledger records what it finds once
(with a warning). Lines the install wizard just added are recorded at the next start. The ledger
is self-guarded with the rest of `/etc/app-listener` and removed by `uninstall`.

### Hash cache

To spare reading every binary at each start (gigabytes for Steam or Electron apps), the ledger
also remembers each binary's hash under its file identity: device, inode, size, modification,
change and birth times, and filesystem id. A start whose binary still shows exactly that identity
reuses the hash instead of reading the file; the ledger check above still decides admission. Any
change to the content moves the change time, which no user can set back, so a changed binary is
always read again.

Only files on ext4, xfs, btrfs or f2fs, on a mount without `nosuid`, are cached: FUSE, network
and overlay filesystems, and images a user mounted, report whatever timestamps their author
chose. Before a hash is cached the file is flushed to disk, so a process writing it through a
shared memory mapping can't change it later without moving its change time. A file changed less
than a second before it is read is never cached. Rows unused by a start are dropped. The ready
line in the journal says how many hashes came from the cache.

## trust

```bash
sudo app-listener trust list                      # only list the refused binaries
sudo app-listener trust                           # review and confirm one by one
sudo app-listener trust add /opt/myapp/bin/myapp  # confirm a hand-added line
```

Lists the binaries the daemon refused (line, resolved path, resource, sha256, owner, mtime) and
asks you to confirm each. An app maintainer script the daemon refused to run (Steam's `steam.sh`,
see [App maintainers](daemon.md#app-maintainers-steam-runtime-updates)) is listed and confirmed the
same way.

| Command or flag | Description |
|---|---|
| `trust list` | Only list the refused binaries |
| `trust add PATH...` | Confirm these lines as they are now (a line added to `daemon.conf` by hand) |
| `-y, --yes` | `trust` and `trust add`: confirm every binary without asking |
| `--no-reload` | `trust` and `trust add`: don't reload a running daemon afterwards |

Only the exact content that was refused is confirmed: a file that changed again is skipped and
judged afresh.

The desktop window's Trust page does the same (see [gui.md](gui.md#trust)).

## Launch scan

Identity is the executable, so `ELECTRON_RUN_AS_NODE=1 code -e ...` or
`NODE_OPTIONS=--require=/tmp/x.js discord` is the whitelisted binary running **its caller's** code.
For whitelisted Node/Electron, Chromium, JVM and Bun runtimes (detected from the binary's own
strings, or the name `java`), the daemon judges the new process's environment and arguments at
exec, before any of its code runs. Bun covers `bun build --compile` apps such as opencode and the
native `claude`.

A launch is **risky** when it has any of:

- environment: `ELECTRON_RUN_AS_NODE`, `NODE_PATH`, `NODE_REPL_EXTERNAL_MODULE`, or
  `NODE_OPTIONS` / `JAVA_TOOL_OPTIONS` / `_JAVA_OPTIONS` / `JDK_JAVA_OPTIONS` holding anything but
  a short allowlist of harmless options (heap size, warnings, TLS/CA settings, font and UI
  scaling);
- Bun environment: `BUN_BE_BUN` (a compiled app becomes the `bun` CLI), `BUN_OPTIONS`
  (`--preload`), `BUN_INSPECT*`, `BUN_RUNTIME_TRANSPILER_CACHE_PATH`;
- code-loading flags: `--inspect*`, `--eval`, `--print`, `--require`, `--import`, `--loader`,
  `--env-file*`, `--remote-debugging-*`, `--load-extension`;
- child-process launchers: `--renderer-cmd-prefix`, `--utility-cmd-prefix`,
  `--zygote-cmd-prefix`, `--gpu-launcher`, `--browser-subprocess-path`;
- VS Code extension flags: `--extensions-dir`, `--extensionDevelopmentPath`,
  `--extensionTestsPath`, `--install-extension`.

Flags are matched also with Chromium's single dash and with Node's `_` for `-`. `--user-data-dir`
is not judged: Electron apps pass it, with their own profile, to the processes they start.

A risky launch still runs, but the process (and its forks; an exec is judged afresh):

- never reads a protected file nor changes an app's read-only code dir
  (`TRUST DENIED  op=LAUNCH comm=... path=...`);
- loses the binary's updater and writer rights: it can't replace the app's binaries or create a
  reserved name (a new version's binary, `settings.json`), and nothing it creates is admitted as
  updater-made (`op=PLANT` / `op=WRITEBLOCK`);
- gets nothing else below a whitelist-mode tree either: no `stat`, directory listing or xattr.

Fix: start the app without that variable or flag. A process that mapped untrusted code before its
binary was whitelisted is refused the same way (`op=PRELOADED`). Libraries from a `lib_dir` of the
updater that created the binary don't count: a self-updating app such as Steam starts its new binary
before the daemon admits it, and everything that process runs is code its updater chose. That new
binary's launch is judged at exec too, before the daemon knows what runtime it is: by every
runtime's rules above (so a flag only another runtime treats as risky marks it too), except the
generic Electron app rule below.

Generic Electron binaries that run an app named on the command line are covered by
[`[electron_apps]`](daemon-config.md#generic-electron-apps-electron_apps).

## Launch rules

Some whitelisted CLIs print or spend their secret when asked. The launch scan also judges any
whitelisted binary with launch rules, and a match marks the process the same way
(`TRUST DENIED  op=ARGV`, and `op=ARGV ... word=<the matching word>` at exec). Built in, by the
whitelisted file's name:

- `gh`: `auth token`, `auth status --show-token`, `auth status -t`, `alias set`, `alias import`;
  `auth git-credential` is only logged, since git runs it on every push.
- `aws`: `configure export-credentials`, `configure get`.
- `bitcoin-cli`, `litecoin-cli`: `dumpprivkey`, `dumpwallet`, `walletpassphrase`, `send`,
  `sendall`, `sendmany`, `sendtoaddress`, and the `-stdin` flag (it reads the method from stdin).
- `lightning-cli`: `pay`, `xpay`, `renepay`, `keysend`, `withdraw`, `sendpay`, `sendonion`,
  `fundchannel`, `multifundchannel`, `close`, `exposesecret`.
- Every whitelisted Go binary: `SSL_CERT_FILE` or `SSL_CERT_DIR` set, which would make it trust a
  certificate its caller chose and send its token to that server.

A rule's words match in order anywhere in the arguments, with anything between them, so a flag's
value can't split them (`gh auth --hostname x token`). Add your own with
[`[launch_rules]`](daemon-config.md#launch-rules-launch_rules).

The rules are a tripwire. They don't see combined short flags (`-th`), aliases, a script that
feeds the CLI its arguments another way, or a token the CLI reads from the keyring rather than a
guarded file: with keyring storage, `gh auth token` still prints it.

## The dynamic loader

`/lib64/ld-linux-x86-64.so.2 <program>` (any `ld-*.so*`, or a copy with its `DT_SONAME`) runs
`<program>` with the loader as its executable, so a whitelisted loader admits every program run
through it. `install`, `edit-protected config` and `config record` refuse a new loader line; `config record` lists a
program run that way as not whitelistable. Whitelist the program itself.

## Multicall binaries

uutils coreutils (Ubuntu 25.10+ and 26.04) is one file behind every applet name. app-listener
tells its applets apart: the kernel attests at exec which applet runs (the exec'd basename and
`argv[0]` must both name it), so whitelisting `/usr/bin/head` admits `head` only, never `cat` or
`dd`. Whitelist an applet by its own link (`/usr/bin/<applet>` or
`/usr/lib/cargo/bin/coreutils/<applet>`), not a symlink of another name.

Refused, because their applets can't be told apart:

- **BusyBox and toybox**, which run applets without an exec;
- any other binary that calls itself multi-call (a future uutils diffutils or findutils) until it
  is vetted;
- **GNU coreutils built as one binary** (Fedora's `coreutils-single`, whose applets are
  `#!/usr/bin/coreutils --coreutils-prog-shebang=...` scripts);
- a uutils path that isn't an applet's own name;
- a uutils build with no root-placed `coreutils` link beside it, in `/usr/bin` or in `/bin` (a
  copy, say): Ubuntu's uutils only lists its applets when run under that name.

Where they are refused: `guard -w/-b`, `network guard -w/-b` and `edit-protected forward` reject
them; the installer and `edit-protected config` refuse a new such line; the catalog never
writes one; and the daemon drops one already in `daemon.conf` (CRITICAL in the log) while keeping
everything else enforced.

Use a dedicated binary instead: on Ubuntu the GNU `gnu*` tools or `coreutils-from-gnu`; on Fedora
`dnf swap coreutils-single coreutils`.

`network monitor` watches an applet the same way; a multicall it can't tell apart is watched by
its file, with a warning.

Whitelisting a single binary with several hard links logs a warning: all its names share one
identity. That is fine for one program under several names (`perl`, `e2fsck`), but not for a
multicall the classifier doesn't recognize. The warning needs another name for it in the same
directory, where a multicall's applet links sit; links only in other directories (Steam's
pressure-vessel hard-links its tools into each container it builds) are logged at debug level.

**Processes started before the guard.** The applet is attested at exec, so a multicall process
already running when a guard starts (or the daemon restarts; a `systemctl reload` keeps them) has
no attested applet. A whitelist denies it and a blacklist blocks it until it execs again, and
`network monitor` doesn't report it. Restart long-running applet processes (a `tail -f`, a `sleep`
loop) after starting the guard.
