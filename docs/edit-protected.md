# edit-protected

Once a directory is protected, your shell and your editor can't open it either. `edit-protected`
is the supported way to change a protected file, to give another program temporary access, or to
change the daemon's own config.

```bash
sudo app-listener edit-protected                        # live if a password is set and the daemon runs, else offline
sudo app-listener edit-protected --set-password         # set or rotate the password
sudo app-listener edit-protected --clear-password       # remove it (disables live mode)
sudo app-listener edit-protected --forward -w <bin> [-e EVENTS]   # temporary access for another binary
sudo app-listener edit-protected --timeout-session 10m  # revoke a live session after 10 idle minutes
sudo app-listener edit-protected --edit-config          # edit the running daemon.conf; saving reloads it
```

It opens the directory in a two-pane editor. `Ctrl+S` saves atomically. Binaries, symlinks and
files over 2 MiB are refused.

## Offline and live mode

The mode is chosen automatically.

- **Offline** (default): refused while the daemon runs. Re-scans the catalog, unlocks one vault
  with the master key, lets you edit, and locks it again on exit. A vault never stays open.
- **Live**: when an edit-protected password was chosen during `sudo app-listener install` **and**
  the daemon is running. Nothing is stopped and the fscrypt vaults are never touched:
  1. you enter the password, **first**;
  2. only then does the daemon return the list of guarded directories to pick from;
  3. it briefly grants write access to the one you pick;
  4. you edit, and the grant is dropped.

An unauthenticated caller learns **nothing** about which directories are protected. Grants are
one at a time and are revoked on disconnect, on a SIGHUP reload, or after `--timeout-session`
(default `30m`, up to `24h`; e.g. `10m`, `45s`) without activity. Every editor keypress resets
that timer. The control socket locks out after 5 failed attempts.

### The password and the control socket

The password is **separate from the fscrypt master key**: it only authenticates `edit-protected`.
Its PBKDF2 hash is stored at `/etc/app-listener/edit-auth.hash` (`0600` root), guarded by the
running daemon like `fscrypt.key` (readable and writable only by the app-listener binary).

When that file exists, the daemon opens `/run/app-listener-daemon.control` (`0600`). It accepts
only root peers whose executable is the app-listener binary itself.

`--set-password` is refused if the password was set at install time; re-run `install` to rotate
that one.

### Non-interactive writes

```bash
echo "new contents" | sudo APP_LISTENER_EDIT_PASSWORD=... \
  app-listener edit-protected --resource /home/alice/.ssh --put config
```

`--put <file>` writes a file inside `--resource` from `--content-file` or stdin, authenticating
with `$APP_LISTENER_EDIT_PASSWORD`.

## Temporary access for other programs: `--forward`

Live mode only. Opens no editor: it gives *other* binaries temporary access to the guarded
directories you pick, with `guard`'s `-w`/`-b`/`-e` syntax.

```bash
# let graphify (a uv-managed python script) create its skill under a guarded ~/.claude
sudo app-listener edit-protected --forward \
  -w ~/.local/share/uv/python/cpython-3.12.14-linux-x86_64-gnu/bin/python3.12 \
  -e OPEN,READ,STAT,MKDIR,WRITE
```

- `-w` admits binaries the resource doesn't whitelist; `-b` denies whitelisted ones; `-e` limits
  the rule to those events (default: all). `-w` and `-b` are mutually exclusive.
- The password is asked first, then a multi-select picks the directories (or pass `--resource`
  once per directory).
- While the grant is active, the terminal shows the granted binaries' accesses to those
  directories, allowed and denied, in the same view as `guard`. Without a terminal they're printed
  as `GUARD|` lines.
- The access lasts until you quit (`q` or `Ctrl+C`). It is also revoked when the client dies, on
  a daemon reload, when a granted binary's content changes on disk, or after `--timeout-session`
  with neither a granted binary touching the directories nor a keypress.

Before admitting a binary, know that:

- A `-w` grant covers **every user** running that binary, like a `daemon.conf` whitelist line.
- Admitting an **interpreter** (python, node, bash) admits every script it runs. A script's
  process is its interpreter, so the script itself is not what the guard sees.
- A granted **Node/Electron, Chromium or JVM runtime** is judged at launch like a whitelisted one
  ([launch scan](binary-trust.md#launch-scan)): started with env or flags that load its caller's
  code (`ELECTRON_RUN_AS_NODE`, `NODE_OPTIONS=--require`, `--inspect`, `--remote-debugging-port`,
  `JAVA_TOOL_OPTIONS=-javaagent`, ...), it is refused the directories (`TRUST DENIED op=LAUNCH`). This covers an app that embeds a
  runtime (VS Code, Discord); it doesn't cover `node script.js` or a JVM given a class path, which
  run the script they're given.
- A binary that isn't root-placed (anything under a home directory) can be rewritten in place by
  its owner. You're asked to confirm it (`--yes` skips the question and is required without a
  terminal), and the daemon revokes the grant if its content changes.
- Revocation re-checks already-open files on their next read or write, but memory a process
  already mapped stays mapped.
- Identity is the binary's inode, resolved by the daemon. `-w` leaves a binary already
  whitelisted on a resource untouched, and `-b` leaves one that isn't untouched. The daemon's own
  binary, inspectors and tamper-demoted binaries are refused.
- If the daemon is killed mid-grant, its pinned guards keep enforcing *with* the temporary
  allows. `daemon --lockdown` (`ExecStopPost`) and the next start strip them using
  `/etc/app-listener/temp-grants.json`, which is written before any allow.

## Editing the daemon config

```bash
sudo app-listener edit-protected --edit-config
sudo APP_LISTENER_EDIT_PASSWORD=... app-listener edit-protected --edit-config --content-file new.conf
```

Live mode only, same password. Opens the running daemon's `daemon.conf` in the editor. `Ctrl+S`
shows the diff and asks for confirmation; then the daemon writes the file and reloads, exactly
like `systemctl reload`.

- If the reload fails (no `[watch]` section, more guard slots needed than are free, ...), the
  daemon keeps running the previous configuration, puts the previous `daemon.conf` back and tells
  you why. Edit again or discard.
- If the file changed while you were editing, nothing is overwritten.
- Refused while a grant session is active.
- Binaries added this way are vetted like any reload: one that isn't root-placed still needs
  [`trust-binaries`](binary-trust.md#trust-binaries).

## Editor keys

The editors (edit-protected, `--edit-config`, the installer's config step) and the
diff/changelog viewers open at the first line.

| Key | Action |
|---|---|
| `Ctrl+S` | Save |
| `Tab` / `Shift+Tab` | Indent / dedent by four spaces |
| `PgUp` / `PgDn` | Page |
| `Ctrl+←` / `Ctrl+→` | Jump between words and punctuation, across lines |
| `Ctrl+F` | Find box: *Search* jumps to the next match, *Search all* marks every match |
| `Ctrl+R` | Find and replace (editors only): *Replace* / *Replace all* |
| `Tab`, `Enter`, `Esc` (in a box) | Move between fields, activate, close |

Matching is literal and smart-case: case-sensitive only when the query has an upper-case letter.
Text is syntax-highlighted by file name or `#!` line (Go, Python, C, shell, JSON, YAML, TOML, INI,
Markdown, diffs, ... and `daemon.conf` itself); buffers over 512 KiB, or files with no recognized
language, are not colored.

## Exit audit

Before exiting, both modes check the edited tree against `daemon.conf` and warn (advisory,
acknowledged once) about anything that would sit outside the daemon's protection: a file created
outside every guarded watch path of a grouped vault, a new symlink, a world-readable new secret, a
freshly dropped executable.
