# Desktop GUI

The `app-listener-gui` release asset (or a `make build GUI=1` build) adds a desktop window to the
same binary. It covers the interactive modes in one place; the daemon has no GUI and stays a
systemd service, and `uninstall` and `edit-protected` stay terminal-only.

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

## Workspaces

| Workspace | Covers |
|---|---|
| Files | `monitor` and `guard`, on several folders or files side by side. |
| Network | `network-monitor` and `network-guard`. |
| Install | `install`: the wizard and its maintenance flags. |

## Files

The Files workspace runs `monitor` and `guard` together on each folder or file you add. Targets
are listed on the left; each has its own monitor, its own guard and its own policy.

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

- **Activity**: every monitor event and guard verdict on the target, newest first, with paths
  shown relative to it. Filter by text or verdict, pause the stream, or clear it. Selecting an
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

## Install

The Install workspace runs `install` with every prompt as a dialog. The page shows the installed
binary (and whether it is identical to the build staged below), whether the daemon runs, and this
window's version.

### The daemon binary

The daemon never runs the GUI build. The installer deploys a headless build instead, staged first in
a root-only (`0700`) temporary folder that is deleted when the window closes. **Source** picks where
it comes from:

| Source | Headless build |
|---|---|
| This version's release (release windows only, the default there) | The headless asset of this window's own tag. |
| Latest stable release | The newest `vX.Y.Z` release, as `update` picks it. |
| Latest pre-release | The newest `pre-YYYYMMDD-<sha>` build of main, as `update --channel pre-release` picks it. |
| Local file | A headless binary you built (`make build`), copied in. It is not signature-verified, and the page says so. |

Every download runs the checks `update` runs: Ed25519 signature, checksum, GitHub digest, ELF
architecture and `--version`. Whatever the source, the build must be an app-listener binary without
the `gui` tag or a fyne dependency (read from its Go build info); the installer checks that again
itself right before deploying, and there is no fallback to copying the window's own executable.
The page then shows the staged build's tag (or "local build"), SHA-256 and path.

When the staged build is another version than the window (always the case for a development
window), the page says so: the window's installer writes `daemon.conf`, and before the service is
enabled the deployed build checks that config with `daemon --check`; the install stops if it can't
load it.

### Actions

| Action | CLI equivalent | Needs the staged build |
|---|---|---|
| Install wizard | `install` | yes |
| Refresh catalog whitelists (Live: keep the daemon running) | `install --update-catalog-only` (`--live`) | no |
| Protect new catalog folders | `install --diff-catalog` | no |
| Restore backups | `install --restore-backups` | no |
| Delete backups | `install --delete-post-backups` | no |

**Log denied metadata inspections** is `--allow-metadata-output`. There is no `--yes`: every
confirmation is asked.

The wizard needs the installed binary to be identical to the staged build, as the terminal
wizard needs it to be the one running: if they differ, update it from a terminal first
(`app-listener update`, or `install --binary-only`). The GUI has no binary-only action.

Starting an action switches to its run view: the steps reached so far, a progress bar for long
copies (migration, backups, the download), and the log. A canceled prompt stops the run where it
is; steps already completed stay applied, as in the terminal. Closing the window during a run asks
first.

No action starts while a guard runs in this window (a Files target or the network guard): the
installer moves and encrypts protected folders, restarts the daemon and downloads from GitHub,
and none of those guards admit it. Turn them off first.
