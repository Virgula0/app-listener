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
  confirm it with [`trust-binaries`](#trust-binaries).
- A new root-owned match of the catalog's absolute patterns (`/opt/*/jbr/bin/java`) is admitted
  the same way.

A process started before its binary was replaced keeps its access, through a reload too.
Re-running the replaced image afterwards (an fd held across the update, `/proc/<pid>/exe`) does
not, and neither does a file that later reuses the old inode number.

## Binary ledger

Nothing guards a whitelisted path while the daemon is stopped, booting or being reinstalled. So
every start and reload checks each whitelisted binary against the content hash recorded for its
line in `/etc/app-listener/binaries.db` (for a symlink line, the link path).

A binary is admitted when:

- only root could have placed it (root-owned file and directories, outside every user home, on a
  mount without `nosuid`), or
- its hash matches the ledger, or
- its resource's updater created it while the daemon ran (recorded automatically, as are live
  replacements).

Anything else, for example a binary changed while the daemon was not running in a home or a
user-owned directory, or a new catalog match planted meanwhile, is refused:

```
DAEMON binary-unconfirmed resource=... path=... reason=new|changed sha256=...
```

A refused binary gets no access, and no updater or reserved-name writer rights either.

The first start after upgrading from a version without the ledger records what it finds once
(with a warning). Lines the install wizard just added are recorded at the next start. The ledger
is self-guarded with the rest of `/etc/app-listener` and removed by `uninstall`.

## trust-binaries

```bash
sudo app-listener trust-binaries --list
sudo app-listener trust-binaries                       # review and confirm one by one
sudo app-listener trust-binaries /opt/myapp/bin/myapp  # confirm a hand-added line
```

Lists the binaries the daemon refused (line, resolved path, resource, sha256, owner, mtime) and
asks you to confirm each.

| Flag | Description |
|---|---|
| `--list` | Only list the refused binaries |
| `--yes` | Confirm every listed binary without asking |
| `--no-reload` | Don't reload a running daemon afterwards |
| `[path...]` | Confirm these lines as they are now (a line added to `daemon.conf` by hand) |

Only the exact content that was refused is confirmed: a file that changed again is skipped and
judged afresh.

## Launch scan

Identity is the executable, so `ELECTRON_RUN_AS_NODE=1 code -e ...` or
`NODE_OPTIONS=--require=/tmp/x.js discord` is the whitelisted binary running **its caller's** code.
For whitelisted Node/Electron, Chromium and JVM runtimes (detected from the binary's own strings,
or the name `java`), the daemon judges the new process's environment and arguments at exec,
before any of its code runs.

A launch is **risky** when it has any of:

- environment: `ELECTRON_RUN_AS_NODE`, `NODE_PATH`, `NODE_REPL_EXTERNAL_MODULE`, or
  `NODE_OPTIONS` / `JAVA_TOOL_OPTIONS` / `_JAVA_OPTIONS` / `JDK_JAVA_OPTIONS` holding anything but
  a short allowlist of harmless options (heap size, warnings, TLS/CA settings, font and UI
  scaling);
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
binary was whitelisted is refused the same way (`op=PRELOADED`).

Generic Electron binaries that run an app named on the command line are covered by
[`[electron_apps]`](daemon-config.md#generic-electron-apps-electron_apps).

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
- a uutils path that isn't an applet's own name.

Where they are refused: `guard -w/-b`, `network-guard -w/-b` and `edit-protected --forward` reject
them; the installer and `edit-protected --edit-config` refuse a new such line; the catalog never
writes one; and the daemon drops one already in `daemon.conf` (CRITICAL in the log) while keeping
everything else enforced.

Use a dedicated binary instead: on Ubuntu the GNU `gnu*` tools or `coreutils-from-gnu`; on Fedora
`dnf swap coreutils-single coreutils`.

`network-monitor` watches an applet the same way; a multicall it can't tell apart is watched by
its file, with a warning.

Whitelisting a single binary with several hard links logs a warning: all its names share one
identity. That is fine for one program under several names (`perl`, `e2fsck`), but not for a
multicall the classifier doesn't recognize.

**Processes started before the guard.** The applet is attested at exec, so a multicall process
already running when a guard starts (or the daemon restarts; a `systemctl reload` keeps them) has
no attested applet. A whitelist denies it and a blacklist blocks it until it execs again, and
`network-monitor` doesn't report it. Restart long-running applet processes (a `tail -f`, a `sleep`
loop) after starting the guard.
