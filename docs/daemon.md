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

`--config <path>` applies to all four. Config resolution order: this flag, then
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

## Steam runtime updates

Steam's launcher unpacks its runtime (`ubuntu12_32/steam-runtime.tar.xz`) and runs the runtime's
`setup.sh` with `mkdir`, `tar` and `ln`. `ubuntu12_32` is a read-only `lib_dir`, and those general
tools are not its writers (whitelisting them would let any process plant a library there), so the
first launch after a runtime update fails. The daemon does that job instead, within about 15
seconds:

- **Unpack.** Only when the trust guard proves Steam's own updater created both
  `steam-runtime.tar.xz` and its `.checksum` while the daemon ran, and nothing else has written them
  since. The archive is checked against the checksum, decompressed by `xz` running as `nobody`, and
  extracted as the Steam user. Paths outside `steam-runtime/`, devices, FIFOs and setuid bits are
  refused or stripped. The previous tree is kept as `steam-runtime.old`.
- **Pin.** `setup.sh --force` runs as the Steam user in a private mount namespace, over an overlay
  whose writable layer is a root-only tmpfs under `/run/app-listener/steamrt`. Its output is
  committed only if it touches nothing but `pinned_libs_{32,64}` and `libcurl_compat_{32,64}`, and
  holds only symlinks plus the `done`, `has_pins` and `system_*` records. A copied library or
  anything else is refused (`refusing setup.sh's output`). The pins are redone when a host library
  they point at changes; a failed run is retried after 30 minutes.

Writes happen only while the daemon holds a write grant on `ubuntu12_32` for itself, which is
refused while an `edit-protected` session is active. Start Steam again once the log reads
`unpacked the Steam runtime update`.

An archive written before the daemon started (Steam updated while the daemon was stopped) has no
provenance and stays packed: `not proven written by Steam's updater since this daemon started`.
Unpack it by hand with the daemon stopped:

```bash
sudo systemctl stop app-listener-daemon
cd ~/.local/share/Steam/ubuntu12_32
md5sum -c steam-runtime.tar.xz.checksum
mv steam-runtime steam-runtime.old && tar -xJf steam-runtime.tar.xz
cp steam-runtime.tar.xz.checksum steam-runtime/checksum
./steam-runtime/setup.sh --force
sudo systemctl start app-listener-daemon
```

Run the commands between stop and start as the Steam user, not root.

## Related

- [Daemon configuration](daemon-config.md): the config file.
- [Binary trust](binary-trust.md): binary updates, the ledger, `trust`, launch scan.
- [Self-updating apps and the catalog](catalog.md): catalog refresh.
- [edit-protected](edit-protected.md): editing protected files while the daemon runs, and the
  live control socket.
- [Troubleshooting](troubleshooting.md).
