# Desktop GUI

The `app-listener-gui` release asset (or a `make build GUI=1` build) adds a desktop window to the
same binary. It covers the interactive modes in one place; the daemon has no GUI and stays a
systemd service. Installing from the window also installs this desktop build as
`/usr/local/sbin/app-listener-gui` (see [Install](#install)).

```bash
sudo app-listener-gui --gui
```

- `--gui` opens the window only without a subcommand. Any other flag is then ignored, and a
  warning names each one: everything a flag would set is a control in the window. With a
  subcommand (`sudo app-listener-gui monitor --gui`), the subcommand runs in the terminal and
  `--gui` is ignored with a warning.
- It needs root, like every mode it drives (each one loads eBPF programs).
- X11 or Wayland. Under `sudo` or `pkexec` it reuses the invoking user's display session
  (`DISPLAY`, `WAYLAND_DISPLAY`, `XAUTHORITY`), taken only from that user's own processes and
  files.
- The default build answers `--gui` with an error naming the GUI build.
- On a tiling window manager (Hyprland, sway, i3) the window opens floating, sized to fit its
  content and the output (1280x820 when it fits, at most 90% of the output) and centered. The
  WINDOW card on Home turns this off; unchecking it tiles the window again. The WM is reached
  over its own IPC socket (`HYPRLAND_INSTANCE_SIGNATURE`, `SWAYSOCK`, `I3SOCK`, or i3's socket
  under `XDG_RUNTIME_DIR`), taken like the display variables and only if the invoking user owns it.
- Pages and dialogs shrink to the window: a page smaller than its content scrolls, and a dialog
  never grows past the window.

## Home

The landing page checks the host and links the workspaces:

| Check | Meaning |
|---|---|
| Privileges | Running as root. |
| eBPF | The kernel loads BPF programs: monitors can run. |
| BPF LSM | `bpf` is in the active LSM list: guards can deny. Without it, guards are disabled and monitors still work. |
| Daemon | Whether the protection daemon is running. |

The sun/moon button at the bottom of the sidebar switches between the dark and light themes, and
the button next to the logo collapses the sidebar to its icons. Both choices are remembered. The
USAGE lines at the bottom of the sidebar show the whole process's CPU (share of one core),
resident memory and events per second across workspaces, refreshed every second.

### Tables

Every event and process table has a search field (Ctrl+F focuses the one on screen) and an
**Export** button that writes the rows as currently filtered to a CSV or JSON file. The file is
written without following symlinks in its path and is owned by the user who ran `sudo`. Clicking
a cell selects its whole row.

### File and folder pickers

Every Browse, File... and Folder... button opens the same picker, starting in the home of the user
who ran `sudo`. Its search field has the focus when it opens (Ctrl+F returns to it) and filters the
folder shown as you type; several words must all appear in the name. **Subfolders** searches the
whole tree below the folder instead (symlinks are not followed, `/proc`, `/sys` and `/dev` are
skipped, and the search stops at 1000 matches). Up/Down move through the results, Enter opens the
highlighted folder or picks the highlighted file, Escape cancels. Clicking a folder opens it; a
folder picker picks the folder shown. The path field above takes a typed absolute path.

## Workspaces

| Workspace | Covers |
|---|---|
| Files | `monitor` and `guard`, on several folders or files side by side. |
| Network | `network-monitor` and `network-guard`. |
| Protected | `edit-protected`: edit, temporary access, daemon config, password. |
| Update | `update`, from the stable or the pre-release channel. |
| Install | `install`: the wizard and its maintenance flags. |

## Files

The Files workspace runs `monitor` and `guard` together on each folder or file you add. Targets
are listed on the left (the button beside Add collapses the list to a narrow rail); each has its
own monitor, its own guard and its own policy.

**Add** asks for what `monitor` takes on the command line:

| Control | CLI equivalent |
|---|---|
| Path (typed, or picked with Folder... / File...) | `-w` |
| Include subfolders | `-r` |
| Max depth (0 = every level) | `-d` |
| Show events | `-e` |

Targets may not overlap: the same path, or one inside another, is refused. The monitor starts
right away and the guard starts off.

A target has three tabs:

- **Activity**: every monitor event and guard verdict on the target, newest first, with absolute
  paths. Filter by text or verdict, pause the stream, or clear it. Selecting an
  event offers to whitelist or blacklist the executable that produced it.
- **Processes**: one row per executable that touched the target, keyed by its inode, with event
  and denial counts and its place in the policy.
- **Guard**: the policy. Pick Off, Whitelist (`-w`) or Blacklist (`-b`), add executables by
  full path (typed, or picked with Browse...), then Apply. Nothing is enforced until Apply.

Whitelisting or blacklisting a process resolves `/proc/<pid>/exe` and admits that file only if
it is still the same inode the process ran; if the path now names another file, the action is
refused. The dialog shows the path, inode and SHA-256 before anything changes.

How Apply changes a running guard:

| From | To | Effect |
|---|---|---|
| Off | Whitelist or Blacklist | A new guard starts. |
| Whitelist | Whitelist | Updated in place: removed executables lose access before new ones gain it. |
| Blacklist | Blacklist | Updated in place: new executables are denied before removed ones are released. |
| Blacklist | Whitelist or Off | The new guard (if any) starts before the old one stops. |
| Whitelist | Blacklist or Off | The whitelist stops first, so the path is briefly unprotected. The dialog warns before applying. |

If a new guard fails to start, the target is left unguarded and the error is shown under its
header. Events caused by the GUI itself are not listed. Removing a target, or closing the window,
stops its guard and then its monitor.

## Network

The Network workspace runs `network-monitor` and `network-guard` together. Unlike Files, it has
one monitor and one guard: the network guard applies to the whole host.

| Tab | Covers |
|---|---|
| Watched | The programs `network-monitor` traces (its positional binaries, added by typed full path or Browse...) and the events shown (`-e`). Changes apply at once. |
| Activity | Monitor events and guard verdicts in one table (event, process, protocol, local and remote address, size). Selecting an event offers to watch, whitelist or blacklist its executable. |
| Processes | One row per executable seen on the network, keyed by its inode, with whether it is watched and its place in the guard policy. |
| Guard | The `network-guard` policy, applied with Apply. |

Guard controls and their CLI equivalents:

| Control | CLI equivalent |
|---|---|
| Mode Whitelist / Blacklist (Off: no guard) | `-w` / `-b` |
| Keep network daemons online (whitelist only) | `--auto-infra` |
| Block local sockets too (whitelist only) | `--unsafe` |
| Rate-limit events (on by default) | `--no-throttle` when off |
| Report events | `-e` |

Before Apply, every listed executable is hashed again; one whose content changed since it was
added is refused until it is removed and added again. The binaries `--auto-infra` admits are shown
under the list, marked as added automatically.

Apply always starts the new guard before stopping the old one. While both run, an operation needs
both to allow it, so a policy change never lets through what either policy blocks. The new guard
keeps the old one's code-suspect marks (see [network-guard](network-guard.md#code-integrity--w)),
and judges already running whitelisted programs from `/proc` as a fresh `network-guard` start does:
restart them for the full guarantee. If the new guard fails to start, the old one keeps running and
the error is shown above the tabs.

With **Block local sockets too**, the GUI's own executable is added to the whitelist (shown as
"this window"), since blocking AF_UNIX would otherwise cut the window off its display server. If
the window still stops responding, Ctrl+C in the terminal that started it stops the guard with it.

Closing the window stops the guard, then the monitor.

## Protected

The Protected workspace is `edit-protected`. The window never touches a protected file or the
daemon's control socket itself: each action runs the installed headless binary
(`/usr/local/sbin/app-listener edit-protected --backend`) and drives it over a private pipe. That
process is the one the daemon admits: it authenticates, holds the grant, and reads and writes the
files; the window only shows them. Before the password is sent, the window checks that the binary
is root-placed, that it is the binary the running daemon executes (after an `update` without a
daemon restart they differ, and the page says to restart it), and that the started process runs
that same file. The backend gets a fixed environment, not the window's. It never runs the GUI build,
which is never given access to a guarded directory either (**Temporary access** refuses to admit it).

The page shows the edit mode (live or offline, chosen like the terminal), whether the daemon runs,
and the password's state.

| Action | CLI equivalent |
|---|---|
| Edit protected files | `edit-protected` |
| Temporary access | `edit-protected --forward -w/-b [-e]` |
| Daemon configuration | `edit-protected --edit-config` |
| Change the password (asks the current one; it can't be removed) | `--set-password` |
| Live session idle timeout | `--timeout-session` (5 minutes to 4 hours) |

Every live action asks for the password first; only then does the daemon list the guarded
directories. A wrong password can be retried (the daemon's lockout applies).

**Edit protected files** opens one directory in a tree and an editor, live (the daemon grants write
access to it for the session) or offline (the daemon must be stopped; the directory is unlocked
with the master key and locked again when the session ends). When only the offline flow applies
and a daemon runs (no password, or its control socket is missing), the card says to stop it and
**Start** stays disabled; so does **Set password** for an install with no password yet. Ctrl+S saves. Binaries, symlinks and
files over 2 MiB are refused; **New file**, **New folder**, **Permissions**, **Owner** (a login
user) and **Delete** work as in the terminal editor. Typing and using the session's controls is
activity for the idle timeout. If the daemon ends the session (idle timeout, reload), the page
says so and nothing more is written. **End session** narrows the grant (or locks the directory) and
shows the post-edit audit.

**Temporary access** takes the binaries (absolute paths), admit or deny, and the event types; then
the password, then the directories. Admitting a binary that isn't root-placed asks first, as the
terminal does. The active grant shows the granted binaries' accesses until **Stop and revoke**, the
idle timeout, or a reload.

**Daemon configuration** fetches `daemon.conf` into an editor; **Review and apply** (or Ctrl+S)
shows the diff and, once confirmed, the daemon writes and reloads it. If the file changed
meanwhile, nothing is written and the next review diffs against the current file. A configuration
that fails to reload is not kept.

One session runs at a time, and none while an installer or updater run is in progress or a guard
runs in this window (it would deny the backend's reads and writes). Closing the window during a
session asks first, then ends it: the grant is revoked, an offline directory is locked again. If
the window dies, the backend sees its pipe close and does the same.

The window shows the files it edits: on X11 any program in your session can read the window or
type into it. Prefer a Wayland session when editing secrets.

## Update

The Update workspace runs `update` in the window with every prompt as a dialog. The page shows the
installed version (and the desktop build's), whether the daemon runs, and the newest release of the
selected **Channel** (looked up on GitHub when the page opens, on a channel change and on **Check**),
with whether `update` would install it. A newer release also shows an **A new version is available**
banner at the top of the page, with an **Update now** button.

The window checks for a newer release by itself too: when it opens, every 6 hours, and after each
run. The installed version's own channel is checked (`pre-*` builds track pre-releases, anything
else stable). A newer release puts a dot on the **Update** entry of the sidebar. A failed lookup
(offline, or a network whitelist running in this window) leaves the dot as it was.

| Channel | CLI equivalent |
|---|---|
| Latest stable release (default for an installed `vX.Y.Z`) | `update` |
| Latest pre-release (default for an installed `pre-*` build) | `update --channel pre-release` |

There is no local-file source and no `--yes`. **Update** runs the same steps as the terminal: the
release's headless asset is downloaded and verified (Ed25519 signature, checksum, GitHub digest,
ELF architecture, `--version`, no `gui` build tag); the changelog since the installed version opens
in a dialog and you confirm; a changed daemon unit shows its diff and asks; then the daemon is
stopped, its folders verified locked, the binary replaced atomically and the daemon started again.
Nothing changes before that last step.

The release's desktop build is downloaded and verified with the headless one (same checks, but it
must carry the `gui` tag) and then replaces `/usr/local/sbin/app-listener-gui`, or is installed
there when missing. **Also install or update the desktop build** (on by default) controls this;
unchecked, only the daemon is updated and that file is left as is. An open window keeps running
its old version until you restart it; the window's own executable is never copied. A release
without a desktop build for this architecture updates the daemon only and logs a warning. When the
daemon is already the newest, **Update** only brings the desktop build to the daemon's version
(with the option checked; nothing to do when it already is).

The run view, the close confirmation and the guard rule are the Install workspace's (below): no
update starts while a guard runs in this window, nor while an installer run is in progress.

## Install

The Install workspace runs `install` with every prompt as a dialog. The page shows the installed
binary (and whether it is identical to the build staged below), whether the daemon runs, and the
installed desktop build (and whether it is on the daemon's version).

### The release

The daemon never runs the GUI build. The installer deploys a release's headless build instead, and
installs the same release's desktop build next to it, so the host gets both:

| Path | Build | Run it |
|---|---|---|
| `/usr/local/sbin/app-listener` (`/usr/local/bin` symlink) | headless, the daemon | `sudo app-listener <mode>` |
| `/usr/local/sbin/app-listener-gui` (`/usr/local/bin` symlink) | desktop | `sudo app-listener-gui --gui` |

Both are staged first in a root-only (`0700`) temporary folder that is deleted when the window
closes. **Source** picks the release they are downloaded from; only signed releases are offered,
never a local file:

| Source | Release |
|---|---|
| This version's release (release windows only, the default there) | The headless asset of this window's own tag. |
| Latest stable release | The newest `vX.Y.Z` release, as `update` picks it. |
| Latest pre-release | The newest `pre-YYYYMMDD-<sha>` build of main, as `update --channel pre-release` picks it. |

Every download runs the checks `update` runs on each build: Ed25519 signature, signed asset name,
checksum, GitHub digest, ELF architecture and `--version`. The headless build must be an
app-listener binary without the `gui` tag or a fyne dependency, the desktop build one with them
(read from the Go build info); the installer checks both again itself right before deploying, and
there is no fallback to copying the window's own executable. The page then shows the staged tag
and each build's SHA-256 and path. The release builds its desktop asset best-effort: a release
without one stages the headless build alone, the page says so, and the wizard installs only the
daemon. To deploy a build of your own, use `install --binary-only` from a terminal.

When the staged build is another version than the window (always the case for a development
window), the page says so: the window's installer writes `daemon.conf`, and before the service is
enabled the deployed build checks that config with `daemon --check`; the install stops if it can't
load it.

### Actions

| Action | CLI equivalent | Needs the staged build |
|---|---|---|
| Install wizard | `install` | yes |
| Refresh catalog whitelists (Live: keep the daemon running, on by default) | `install --refresh-catalog-only` (`--live`) | no |
| Protect new catalog folders | `install --diff-catalog` | no |
| Restore backups | `install --restore-backups` | no |
| Delete backups | `install --delete-post-backups` | no |

**Also install the desktop build (app-listener-gui)** (on by default) installs the staged
release's desktop build once the wizard completes. Unchecked, the wizard installs only the daemon.
It is disabled when the staged release has no desktop build.

**Log denied metadata inspections** is `--allow-metadata-output`. There is no `--yes`: every
confirmation is asked. Restore and Delete backups list the backups found, then ask once more naming
each one; Restore is refused while the daemon runs, Delete is not.

A card whose action the installer would refuse says why and keeps **Start** disabled: Restore
while any daemon runs (`sudo systemctl stop app-listener-daemon` first, which locks its folders
again); the wizard, Protect new catalog folders and a non-live refresh stop a systemd daemon
themselves, so only a daemon running outside systemd disables them; a live refresh needs the
daemon running under systemd. The install wizard is also disabled once an installation has
completed (`/etc/app-listener/install-complete`, see [installation](installation.md#the-install-wizard)):
uninstall first to run it again. The window never stops the daemon itself.

The wizard needs the installed binary to be identical to the staged build, as the terminal
wizard needs it to be the one running: if they differ, update it first (the Update workspace, or
`install --binary-only` from a terminal). The GUI has no binary-only action. Once the wizard
completes, it installs the staged desktop build at `/usr/local/sbin/app-listener-gui`, replacing
any one there (when the option above is checked); a canceled or failed wizard installs none.

Starting an action switches to its run view: the steps reached so far, a progress bar for long
copies (migration, backups, the download), and the log. A canceled prompt stops the run where it
is; steps already completed stay applied, as in the terminal. Closing the window during a run asks
first.

No action starts while a guard runs in this window (a Files target or the network guard): the
installer moves and encrypts protected folders, restarts the daemon and downloads from GitHub,
and none of those guards admit it. Turn them off first.

### Uninstall

The **Uninstall** card at the bottom of the page runs `uninstall` with every prompt as a dialog, on
the same run view (and under the same rules) as the actions above:

| Option | CLI equivalent |
|---|---|
| Keep the desktop build (app-listener-gui), off by default | `--keep-gui` |
| Also delete the fscrypt master key, off by default | `--delete-key` |

Like the terminal, it is refused while the daemon runs (under systemd or not): the button is
disabled until you stop it (`sudo systemctl stop app-listener-daemon`, which locks its folders
again). The window doesn't stop the daemon itself. **Uninstall** first shows a summary to confirm,
then each encrypted folder is asked about (decrypt permanently, or keep encrypted), and the
per-user ssh-agent setup and the backup cleanup ask as in the terminal. Removing the desktop build
deletes `/usr/local/sbin/app-listener-gui`; the open window keeps running until you close it.
