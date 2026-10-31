# daemon: permanent protection with encryption at rest

The daemon is what `sudo app-listener install` sets up. It runs the [guard](guard.md) whitelist
engine over any number of directories, from one config file, and adds an fscrypt encryption
lifecycle:

- while the daemon runs, each protected directory is unlocked but readable only by the programs
  its section whitelists;
- while it doesn't (shut down, booting, uninstalled without decrypting), the directory is
  encrypted.

The rule it never breaks: **a protected directory is never readable without a live guard
attached.** Directories are unlocked only after their guard is in place, and locked again on
shutdown *while the guards still deny*.

```bash
sudo app-listener daemon genkey   # create the fscrypt master key
sudo app-listener daemon            # run with /etc/app-listener/daemon.conf
```

You rarely run it by hand: `install` deploys it as the `app-listener-daemon` systemd service.

## Commands and flags

| Command | Description |
|---|---|
| `daemon` | Run the daemon (the systemd service's `ExecStart`) |
| `daemon check [--verifier-only]` | Check BPF-LSM is active and load every guard eBPF program into this kernel's verifier, attaching nothing; exit non-zero on any rejection. `--verifier-only` skips the BPF-LSM activation check (CI hosts that build BPF-LSM in without enabling it): a pass then says the programs verify, not that the host enforces |
| `daemon lockdown` | Force-lock every encryption root in the config and exit. The systemd unit runs it as `ExecStopPost`, after every exit (clean stop, crash, SIGKILL, startup timeout) |
| `daemon genkey` | Generate `/etc/app-listener/fscrypt.key` and exit (regenerating asks for confirmation) |
| `daemon login <uid>` | Ask the running daemon to make and guard that user's [runtime sections](#runtime-sections-at-login); exit 0 with no daemon running, non-zero when it refused or timed out (75 s) |
| `daemon logout <uid>` | Ask the running daemon to drop that user's runtime sections once their directory is gone; never fails |

`--config <path>` applies to `daemon`, `check`, `lockdown` and `genkey`. Config resolution order: this flag, then
`/etc/app-listener/daemon.conf`, then `daemon-samples/daemon.conf`.

Flags of `daemon` itself:

| Flag | Default | Description |
|------|---------|-------------|
| `--headless` | `false` | Log `DAEMON ...` lines to stderr (journald when run as a systemd service) |
| `--blocked-only` | `false` | Print only denied attempts (presentation only) |
| `--no-log-metadata-blocks` | `false` | Don't log denied metadata-only process inspections; see [below](#quiet-metadata-denials) |
| `--pprof <addr>` | none | Serve `net/http/pprof` on a loopback address for profiling |
| `--serve[=host:port]` | off | Mirror the TUI to a browser |

The flags `daemon --check`, `--verifier-only`, `--lockdown` and `--genkey` are aliases of the
subcommands and stay listed in `--help`: installed systemd units, `scripts/install.sh` and
`check-compatibility.sh` run binaries of other versions with them.

The config format is covered in [Daemon configuration](daemon-config.md).

## Lifecycle

### Startup order

Attach the guards, then unlock the vaults, then populate the protected inodes, then resolve
binaries that didn't exist yet, then re-sync. Nothing is ever unlocked before its guard is
attached.

A vault holding a file with a hard link outside it is never unlocked: the guard covers only names
inside the tree, and a link can be made while the daemon is stopped (fscrypt allows linking a
locked file into an unencrypted directory). The start (or reload) fails and names the files; remove
their outside names (`find / -xdev -samefile <file>`) and start again.

### Boot ordering

The unit is `Type=notify` and ordered
`Before=systemd-user-sessions.service cron.service crond.service cronie.service atd.service`. Logins,
lingering user managers, cron and `at` jobs wait until every guard is attached (`READY=1`), so no
user code runs while the whitelisted paths are unguarded. This is ordering only: a failed daemon
start does not lock users out.

### Runtime sections at login

A [runtime section](daemon-config.md#runtime-sections) guards a directory in `/run/user/<uid>`,
which systemd mounts at the user's first login (`user-runtime-dir@<uid>.service`) and removes at
the last logout. The installer adds a drop-in,
`/etc/systemd/system/user-runtime-dir@.service.d/app-listener.conf`, that runs
`app-listener daemon login <uid>` right after the mount and `daemon logout <uid>` after the
unmount. Both talk to the daemon over `/run/app-listener-daemon.login` (root only, and only this
binary).

At login the daemon:

1. creates each missing runtime directory as root, mode 0700, so the user can't enter it yet, with
   a root-owned `.app-listener/` inside that keeps it non-empty: no whitelisted program with
   `DELETE` can then remove the guarded directory and let it be recreated unguarded;
2. reloads, so a guard anchors on that directory;
3. checks that the guard holds the very directory it made, then hands it to the user (`chown`).

Only then does `user@<uid>.service` start the user's processes. If any step fails, the login
fails. If no daemon is running, the login goes ahead with nothing guarded there. The drop-in is
ordered after the daemon unit, so a login during the daemon's start waits for it.

The mount is fresh at login, so a runtime directory already in it was made by one of the user's
processes before the daemon, and may hold its listeners: the daemon refuses it and the login fails
(remove the directory and log in again). Only a directory a running guard already holds is
accepted. A directory present when the daemon itself starts (the user logged in while no daemon
ran) is guarded as found, and may hold sockets bound before the guard. The daemon plants its
sentinel in it then (and at every reload), so an agent removing its own socket can't leave it
empty.

At logout the directory is gone with its mount, and the daemon reloads to drop the resource. It
drops one only once its mount is gone: a filesystem the user mounts over `/run/user/<uid>` (FUSE,
through the setuid `fusermount3`) hides the directory, and every reload is refused until it is
unmounted, as is a daemon start that can't look inside. Like
any reload, a login or logout that adds or drops a runtime section ends a live `edit-protected`
session.

### The Secret Service gate

With a [`[secret_service]` block](daemon-config.md#secret-service-sections) the daemon:

- installs a session bus policy (`/etc/dbus-1/session.d/app-listener-secrets.conf`) under which
  only root may own `org.freedesktop.secrets` or monitor the bus, and the Secret Service interfaces
  are reachable only through that name;
- installs a drop-in (`/etc/systemd/user/gnome-keyring-daemon.service.d/app-listener.conf`) that
  points gnome-keyring at a relay socket, `/run/app-listener/secret-relay/<uid>.bus`, guarded so
  only the provider may connect;
- at each login, before the user's processes start, registers `/run/user/<uid>` with the guard and
  opens that user's relay. It reaches the user's bus as root for the keyring, through the socket
  named `/run/user/<uid>/bus` (never a symlink), and the guard lets only a root-placed
  `dbus-broker` or `dbus-daemon` that loaded no untrusted code write on those connections: a
  listener the user puts at that name can't answer the relay, so it can't forge a caller.

Every method call to the keyring passes the gate first. Calls it doesn't know are refused; reading,
writing or deleting an item a [secret scope](daemon-config.md#secret-service-sections) covers needs
a caller that scope's whitelist admits. The caller is the process the bus reports for the
connection, checked by inode like any guard decision, and refused if it ran another program since it
connected to the bus, if another process ever sent on one of its bus connections (an inherited or
passed socket: the bus would credit that process's calls to it), or if it started before the daemon
opened that user's relay (its connections may predate the guard, which then never saw them). Each decision is logged like a guard event, under the scope's anchor.

Both files are written to match the configuration at every start and reload, kept when the daemon
stops (the keyring then can't reach the bus, so the Secret Service is off rather than unguarded) and
removed by `uninstall`. Users logged in when they are first installed must log in again.

### Reload (`systemctl reload`)

```bash
sudo systemctl reload app-listener-daemon
```

SIGHUP recomputes every binary's inode identity atomically: new guards attach before old ones
detach. A malformed config keeps the previous one running.

### Shutdown and crashes

`need_encryption: true` resources must already carry an fscrypt policy (the install wizard sets
it up). On shutdown, keys are removed in two passes while the guards still deny access, and the
hooks detach only after every vault is keyless.

A hard `SIGKILL` can't be caught, but the guards' LSM links are pinned under `/sys/fs/bpf`, so the
trees stay enforced until `ExecStopPost` (`daemon lockdown`) locks the vaults.

A wrong or old master key fails immediately with "invalid wrapping key". The daemon never silently
generates a new one.

## Self-protection

The daemon guards its own files in `/etc/app-listener` (config, fscrypt key, edit-protected
password hash, binary ledger): only the app-listener binary can read or write them while it runs.

The password hash lives in the binary ledger (`binaries.db`), so the ledger and its journal are
sealed like the key: no other process, root included, can read them. If the fscrypt key or the
ledger can't be sealed, the daemon refuses to start (exit status 78) before unlocking anything. The `/etc/app-listener` directory guard alone is best effort: if it
fails to attach, the daemon logs a `CRITICAL` line and keeps running.

## Quiet metadata denials

Some desktop components routinely look at other processes: a compositor or audio server reading
`/proc/<pid>` of Steam, for example. When the target holds protected content, those reads are
denied and logged as `op=PTRACE mode=READ`.

- Without `--no-log-metadata-blocks`, they are logged once per caller/target pair and summarized
  per minute.
- With it (the systemd unit's default) they are still denied, just not logged.
- Memory and file denials are always logged.
- Probes the kernel itself doesn't audit (`ps`/`pgrep` reading `/proc/<pid>/stat`) are never
  logged.

A refusal for a process holding several resources' content reads `resource=multiple`, or the
comma-separated paths of its taint set. To diagnose an app that breaks with nothing in the log,
reinstall with `install --allow-metadata-output` (see [Troubleshooting](troubleshooting.md)).

## Filesystem-wide gates

Two protections are coarse by nature and apply to the whole device or filesystem, so only the
daemon's own binaries pass:

- **Raw block device.** The block device behind a guarded filesystem is blocked, so `debugfs`,
  `dd` or `fsck` can't read the on-disk image around the VFS.
- **btrfs copy ioctls.** On btrfs, snapshot (`SNAP_CREATE`/`_V2`), `SEND` and `TREE_SEARCH` are
  denied on the whole filesystem that hosts a guarded root. A snapshot is a new subvolume with the
  source's inode numbers, so no per-inode rule would match the copy.

Side effect, intended: while the daemon runs, raw-imaging backups and btrfs snapshot tools
(`snapper`, `timeshift`, `btrbk`, `snap-pac`, `btrfs subvolume snapshot/send/list`) stop working on
any disk or btrfs filesystem that hosts a guarded resource. These denials are logged as
`resource=raw-block-device` or `resource=btrfs-ioctl`.

## Process inspectors and taint

A process that read a protected file is *tainted*: no program outside its resource's whitelist may
read its memory or its `/proc/<pid>` metadata. Otherwise any process could read the secret back
out of the app that loaded it. Programs that must inspect every process (screen-sharing portals)
are listed under `[inspectors]`, see [Daemon configuration](daemon-config.md#process-inspectors).

## App maintainers (Steam runtime updates)

Steam's launcher, `~/.local/share/Steam/steam.sh`, unpacks and pins its runtimes inside
`ubuntu12_32` and `steamrt64` with `tar`, `mkdir` and `ln`. Those are read-only `lib_dir`s, and
general tools are not their writers (whitelisting them would let any process plant a library
there), so the launch after a runtime update fails. The daemon runs Valve's own `steam.sh` for it,
unmodified, and keeps only what it can validate. Nothing about what the script does is built into
app-listener: a new Steam version needs no new app-listener.

**When.** At start and reload, within about 15 seconds of an entry being created, removed or
renamed directly in one of Steam's lib_dirs (Steam wrote an update), and every 30 minutes (the pins
follow host library upgrades). Never while a Steam binary runs. A failed run is retried after 30
minutes, or sooner once the lib_dirs change again.

**Which script.** Only the `steam.sh` the binary ledger confirms, like a whitelisted binary:
recorded when Steam's updater writes it while the daemon runs, or by `install`, or confirmed with
[`trust`](binary-trust.md#trust). It is a reserved name, so only Steam's binaries can write it
while the daemon runs. A `steam.sh` changed while the daemon was stopped is not run:

```
daemon: .../steam.sh was not run: it is not confirmed in the binary ledger; confirm it with `sudo app-listener trust` if Steam updated it
```

**How it runs.** The script runs as the Steam user in private mount, pid, IPC and network
namespaces, with a clean environment and no new privileges. It sees an empty private home holding
only the Steam directory, plus private `/tmp` and `/run/user/<uid>`: it reaches none of the user's
other files or sockets, and whatever it writes outside the Steam directory is discarded. The
`steam.sh` it runs is a root-owned copy of the file the ledger checked, refused unless its sha256
still matches, bound over its path whatever the path names by then.
Steam's own binaries (`steam`, `steamwebhelper`, ...) are replaced with `true` inside the sandbox,
so the script does its maintenance and then "starts" nothing. Everything it writes lands in an
overlay whose writable layer is root-only, under `/var/lib/app-listener/maintain`. That directory
must not be on overlayfs itself (a container's root, a live system): the kernel refuses an overlay's
writable layer there, and every run fails until it is moved to a real filesystem or a tmpfs.

**What is kept.** Only changes below Steam's lib_dirs that are on the same filesystem as
`steam.sh`, and only when every one passes:

- directories and deletions are kept;
- a symlink is kept only if it leads, through every hop as the tree will be after the commit, into
  one of Steam's lib_dirs, or to a name only root decides, outside every user home and not under
  `/proc`: a root-owned file no one else may write, a directory only root can fill all the way
  down (at most 4096 entries; `/etc/ssl/certs` passes, `/tmp` or a home does not), on a
  filesystem mounted at least once without `nosuid` and not FUSE (a user's udisks2 image presents
  any owner its author wrote), or a missing name in a directory only root can write that is not
  one root later creates for a user (a new home beside `/home/*`, `/run/user/<uid>`). Nothing at
  or below `/media`, `/run/media` or `/mnt` passes, existing or not: a user's own filesystem may
  be mounted there. A `..` may only lead a target
  (`../lib/libfoo.so`). Same-uid malware can write the overlay while `steam.sh` runs (through
  `/proc/<pid>/cwd`), so a symlink out to the user's files would hand Steam's whitelisted
  binaries content no lib_dir vouches for;
- a regular file is kept only if its exact content is already in that lib_dir, as a file or as a
  member of a tar archive there (plain, gzip, bzip2, xz or zstd; `xz` and `zstd` run as `nobody`),
  or it is a path record of at most 4 KiB: lines that are empty or an absolute path of letters,
  digits and `/._+@~:,-` only (what `setup.sh` records about the pins), written without exec bits;
- setuid, setgid and sticky bits are dropped; devices, FIFOs and sockets refuse the run;
- nothing may replace a whitelisted binary of any resource, or a directory above one.

Otherwise nothing is committed. Files are written from what was validated, never re-read from
the overlay afterwards, and owned by the Steam user. Writes happen only while the daemon holds a
write grant on that lib_dir and a reserved-name writer row for its own binary as root (the runtime
holds files named like Steam's binaries, such as `game-bin/steam`). The grant is refused while an
`edit-protected` session is active. The log then reads:

```
daemon: committed N change(s) .../steam.sh made to its lib_dirs; start Steam again if it failed to start
```

The first launch after a runtime update still fails, because `steam.sh` runs again right away and
its own writes are refused. Start Steam again once that line appears.

## Related

- [Daemon configuration](daemon-config.md): the config file.
- [Binary trust](binary-trust.md): binary updates, the ledger, `trust`, launch scan.
- [Self-updating apps and the catalog](catalog.md): catalog refresh.
- [edit-protected](edit-protected.md): editing protected files while the daemon runs, and the
  live control socket.
- [Troubleshooting](troubleshooting.md).
