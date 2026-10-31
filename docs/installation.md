# Installation

Installing app-listener is two separate steps on purpose:

1. **Put the binary on the system** (`install.sh` or a source build). This changes nothing about
   how your files are protected.
2. **Set up the daemon** (`sudo app-listener install`). This is interactive and is never done
   automatically, because it encrypts directories and starts enforcing.

The standalone modes (`monitor`, `guard`, `network monitor`, `network guard`) need only step 1.

## One-line installer

```bash
curl -fsSL https://raw.githubusercontent.com/Virgula0/app-listener/main/scripts/install.sh | sudo bash
```

Installs the latest **stable** release. For pre-release builds:

```bash
curl -fsSL https://raw.githubusercontent.com/Virgula0/app-listener/main/scripts/install.sh | sudo bash -s -- --channel prerelease
```

What `scripts/install.sh` does, in order:

1. Runs `check-compatibility` and aborts, touching nothing, if the host can't run the guard or
   daemon modes.
2. Downloads the latest release of `--channel` (`release`, the default, or `prerelease`) for the
   host's architecture (`app-listener` on x86_64, `app-listener-arm64` on aarch64).
3. Verifies it: the Ed25519 signature of the checksum against the embedded release key, that the
   signature was made for that asset name, the checksum against the binary, the GitHub asset
   digest, and the ELF machine type.
4. Atomically installs `/usr/local/sbin/app-listener` plus the `/usr/local/bin/app-listener` PATH
   symlink.

It does **not** install the daemon. It prints a reminder to run `sudo app-listener install`.

## Build from source

```bash
make build                                  # toolchain in a rootful Docker container
sudo ./build/linux/app-listener install
```

`make build` dumps `vmlinux.h` from the host BTF, regenerates the BPF bindings and builds
`build/linux/app-listener`. Nothing but Docker is needed on the host.

Without Docker, install a recent clang/LLVM, `bpftool`, Go 1.26+ and GCC, then run
`make build-host`. Builds are native only: build arm64 on an arm64 host. See
[Development](development.md) for every target.

## The install wizard

```bash
sudo app-listener install
```

A TUI wizard that runs in a safe order:

1. Stops a running daemon (a daemon running outside systemd is a fatal error: stop it first).
2. Deploys the binary to `/usr/local/sbin/app-listener` when it isn't there yet, by copying the
   running executable (an existing install is left untouched: upgrade it with `install binary` or
   `update`), and generates the fscrypt master key (an existing key is kept).
3. Lets you pick the users to protect.
4. Probes the [built-in catalog](catalog.md) of sensitive directories (SSH, GPG, AI agents,
   browsers, VPNs, password stores, crypto wallets, ...) and lets you pick which to guard.
   Further directories can be added by hand, and the generated `daemon.conf` is shown in an editor
   for review.
5. Encrypts the selected directories. Each is backed up first (an existing
   `.app_listener.backup` aborts: an older backup is never overwritten) and every already
   encrypted directory must unlock with the master key.
6. Asks for the edit-protected password (mandatory; see [edit-protected](edit-protected.md)). When
   resuming an interrupted install, an existing password is kept unless you choose to replace it.
7. Deploys the systemd units, the binary and the config.
8. Runs `daemon check` on the deployed binary before enabling the service: if this kernel's
   verifier rejects the guard eBPF, the service is not enabled. Files already installed are
   compared with the bundled ones: a differing one shows its diff and asks before overwriting.
9. Cleans the orphaned fscrypt metadata of directories that no longer exist (only the
   installer's own `app-listener-key-*` protectors), then offers to delete the migration backups.

Aborting (`Esc`) any step stops the installer; steps that already completed stay done.

An interrupted wizard (a canceled prompt, a failed step) is resumed by running it again: the
existing master key is kept, already encrypted directories are verified and identical files are
skipped. Once a wizard completes (daemon running and enabled, orphaned fscrypt metadata cleaned), it
records the installation in `/etc/app-listener/install-complete` (version and time), and from then
on the wizard refuses to run. Change a completed installation with `install catalog diff`,
`install catalog refresh`, `edit-protected config` or [`update`](#update), or run
[`uninstall`](#uninstall) first to start over. `install binary` and the other subcommands below
are not affected. An installation completed before this record existed has none, so its wizard still runs
until one completes.

Earlier versions installed pacman/apt hooks and a boot unit to refresh the catalog; the wizard,
`install binary` and `update` remove them, since the daemon now refreshes the catalog itself.

### SSH agent

When a user's `~/.ssh` ends up guarded, the wizard offers (one question per user, naming the user
and the unit path, asked before encryption):

- a per-user `ssh-agent` systemd unit. It binds `$XDG_RUNTIME_DIR/ssh-agent/socket`, which the
  catalog's "SSH agent socket" entry guards (a [runtime section](daemon-config.md#runtime-sections)),
  so only the programs whitelisted there (`ssh`, `ssh-add`, `ssh-keygen`) may use a loaded key.
  The wizard says so when that entry is not selected. A unit from an earlier release binds
  `$XDG_RUNTIME_DIR/ssh-agent.socket`, unguarded: `install facultatives` offers the new one and
  rewrites the shell block, and `install catalog diff` adds the runtime section,
- `AddKeysToAgent yes` at the top of `~/.ssh/config` (created if missing; an existing
  `AddKeysToAgent` is kept), written before `~/.ssh` is encrypted, so `ssh` loads a key into the
  agent on first use,
- a marked `SSH_AUTH_SOCK` block in the user's shell startup file (`~/.zshrc`, `~/.bashrc` or fish
  `conf.d`). The login shell's file is created if missing. An existing `SSH_AUTH_SOCK` line or an
  already-set variable is never overridden. `uninstall` removes the block.

None of this is offered when `~/.ssh` is not guarded. The daemon never loads keys itself.

### Bun-based apps (opencode)

Bun extracts a bundled native library to `$TMPDIR/.bun-<uid>-<hash>.so` and `dlopen`s it. On a
world-writable `/tmp` the trust guard denies that load, since it can't tell the app's own
extraction from a planted `.so`.

The wizard offers, per user, to redirect it: it creates a private `~/.cache/app-listener/bun`
(`0700`) and adds a marked shell-function wrapper for each launcher (`~/.zshrc`, `~/.bashrc` or
fish `conf.d`) that runs the command with `$TMPDIR` pointed there. The daemon then reserves the
`.bun-*` name below that dir for the app's binaries, so the app loads its own extraction while no
other process can plant or load one.

A dir the daemon has not reserved yet (at every start, and when it first appears) may already hold
a plant, so the daemon first replaces it with an empty one; the app re-extracts on its next launch.
Ordinary temp files that the app's subprocesses write there are unaffected (only the reserved name
is gated). This only works for interactive shells: a Bun app launched from a desktop menu still
uses `/tmp`. `uninstall` removes the wrappers.

### Facultatives

The SSH agent and the Bun redirect above are the wizard's facultatives: optional per-user setup,
not needed for protection. `install facultatives` installs them without the wizard:

```bash
sudo app-listener install facultatives
sudo app-listener uninstall facultatives
```

It lists, in one multi-select, the SSH agent setup for every non-root user and the Bun redirect for
every user with a Bun app configured or present in the home, whether or not `daemon.conf` guards
`~/.ssh` or the app. Already installed ones are marked. Only the picked ones are installed; nothing
else changes, and it works on a completed installation, with the daemon running or not:

- **SSH agent**: the per-user unit and the `SSH_AUTH_SOCK` block, as in the wizard. When
  `daemon.conf` guards `~/.ssh`, its `config` is left alone (the daemon denies the edit, and a
  stopped daemon leaves the vault locked): add `AddKeysToAgent yes` with
  [edit-protected](edit-protected.md). Otherwise `AddKeysToAgent yes` is added too.
- **Bun redirect**: the private dir and the launcher wrappers. The daemon reserves the dir only when
  it loads its trust set, so for a Bun app it guards, a running daemon is reloaded (`systemctl
  reload`, not a restart).

`uninstall facultatives` is the inverse: it lists the installed facultatives (an ssh-agent unit
matching the bundled one or an `SSH_AUTH_SOCK` block, the Bun wrappers) and removes the picked ones.
A guarded `~/.ssh/config` keeps its `AddKeysToAgent` block, the Bun dir stays (as with `uninstall`),
and an ssh-agent already running keeps running until logout. It also works with the daemon running.

### fscrypt prerequisites

Each filesystem must be initialized for fscrypt (`fscrypt setup --all-users`) and support
encryption (ext4: `tune2fs -O encrypt <dev>`). The wizard checks this before migrating anything.
For a fixable gap it shows the exact command and why, and offers to run it for you; decline and
it aborts. Every command it may run is listed in `internal/fscrypt/prereq.go`.

### Install subcommands and flags

| Command or flag | Description |
|---|---|
| `install --allow-metadata-output` | Installs the daemon unit **without** `--no-log-metadata-blocks`, so denied metadata-only process inspections (`op=PTRACE mode=READ`) are logged. Use it to diagnose an app that breaks with nothing logged; re-run `install` without it to quiet them again. |
| `install catalog diff` | Incremental wizard: lists the catalog directories that exist on the host but are **not** in `daemon.conf` yet, lets you pick which to add, appends them and encrypts them (backup first). Existing sections stay byte-for-byte intact. Refused while the daemon runs: stop it first (`systemctl stop app-listener-daemon`, or the window's Daemon page). At the end it asks whether to start the daemon on the merged config. Requires a previous installation; interactive only. |
| `install catalog refresh` | Re-expands every catalog-matched whitelist and rewrites the config, dropping every line the catalog does not produce. Stops the daemon and unlocks each vault under an ephemeral self-only guard. The running daemon already does this on its own (see [catalog refresh](catalog.md#catalog-refresh)), so you normally never need it. |
| `install catalog refresh --live` | Refresh without stopping the daemon (no lock churn), applied via SIGHUP. Requires the daemon to be running. |
| `install catalog refresh -y, --yes` | Writes without the confirmation prompt (scripts). |
| `install binary` | Moves the freshly built binary to the install path and recreates the PATH symlink. No wizard, config or fscrypt. If the daemon service is installed, diffs its unit against the one the new build ships (replaced only if you accept, as with [`update`](#update)) and asks whether to restart it now; a skipped restart applies the accepted unit on the next one. |
| `install backups restore` | Undoes the fscrypt migration: lists the `.app_listener.backup` directories, and after a confirmation naming each one deletes the encrypted copies and moves the backups back. Refused while the daemon runs. |
| `install facultatives` | Installs the optional per-user setup you pick (SSH agent, Bun redirect), for every user and whether or not `daemon.conf` guards it. No wizard; the daemon keeps running. See [Facultatives](#facultatives). |
| `install backups delete` | Deletes the `.app_listener.backup` directories (listed, then a confirmation naming each one). Works with the daemon running: the backups sit outside every guarded tree. |

`install catalog refresh` refreshes the whitelists of existing sections; `install catalog diff` adds new
sections. Run both to fully re-sync with the catalog.

## Update

```bash
sudo app-listener update            # shows the changelog and asks first
sudo app-listener update --yes      # no changelog viewer, no prompt
sudo app-listener update --channel pre-release
```

Self-updates from the latest signed release of `--channel` (`stable`, the default, or
`pre-release`), fetching the asset for the running binary's architecture. The Ed25519 signature,
signed asset name, checksum, GitHub asset digest and ELF machine are all verified before anything
is written.

The desktop window (`app-listener-gui`) runs the same update from its Update workspace, with the changelog and
the confirmations as dialogs (see [Desktop GUI](gui.md#update)).

When the desktop build is installed (`/usr/local/sbin/app-listener-gui`, deployed by the window's
installer), `update` replaces it with the new release's desktop asset, downloaded and verified
like the binary before anything is written, after the daemon is back up. A release without one
leaves the old desktop build in place with a warning. With the daemon already up to date, an
installed desktop build of another version is brought to the daemon's release. Updating from the
window also installs the desktop build when it is missing.

The changelog covers every release of the channel between the installed version and the new one,
newest first, with a link to the full compare. If the installed version is not among the fetched
releases (a local `dev` build, or older than the 100 newest), it shows every fetched release up to
the new one and says so.

When the daemon service is installed, `update` also compares the systemd unit the new release ships
(`/etc/systemd/system/app-listener-daemon.service`) with the installed one. If they differ, it shows
the diff and asks whether to replace it; declining keeps the installed unit and still updates the
binary. The new unit is written while the daemon is stopped, so the restart already runs under it.
With `--yes` the new unit is accepted and the diff is logged. The shipped unit is written as is:
local edits to the installed one (an `install --allow-metadata-output` unit, an added `--pprof`)
are not carried over, so re-run `install --allow-metadata-output` after updating to keep metadata
denials logged.

## Uninstall

```bash
sudo app-listener uninstall
sudo app-listener uninstall --delete-key
sudo app-listener uninstall --keep-gui
```

- Refuses while the daemon runs.
- Removes the completed-install record first, so an interrupted uninstall can be followed by a
  fresh `install`.
- Re-scans the catalog and asks, per encrypted directory, whether to decrypt it in place for good
  (default: keep it encrypted).
- Removes the binary, the desktop build (`app-listener-gui`, kept with `--keep-gui`), their PATH
  symlinks, the units and the config.
- Keeps the master key unless `--delete-key` is given. Without the key, any directory still
  encrypted can never be unlocked again, so with directories left encrypted `--delete-key` asks
  once more (default: keep the key).

The desktop window runs the same uninstall from its Install workspace (see
[Desktop GUI](gui.md#uninstall)).
- `uninstall facultatives` removes only the optional per-user setup you pick, and works with the
  daemon running (see [Facultatives](#facultatives)).
- At the end, lists any `.app_listener.backup` migration copies and offers to delete them (all
  preselected, one confirmation). They are plain unencrypted copies and harmless to keep.

## Docker

```bash
sudo make deploy
docker exec -it app-listener app-listener monitor /tmp
```

`make deploy` builds the binary inside Docker, on the toolchain image (`docker/builder.Dockerfile`),
so nothing is compiled or written on the host, and starts an idle container to `docker exec` into
(`/app` is the working directory and on `PATH`). It needs rootful Docker: a rootless container can't
load eBPF, so `deploy` refuses it (`sudo` drops a rootless `DOCKER_HOST`). The container is
privileged and gets the host's `/sys/kernel/security` read-only, so the guard modes still refuse to
start on a kernel without BPF-LSM. `$MONITOR_DIR` (default `./data`) is mounted read-only at
`/data`. `make deploy-down` stops it.
