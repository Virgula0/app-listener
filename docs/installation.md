# Installation

Installing app-listener is two separate steps on purpose:

1. **Put the binary on the system** (`install.sh` or a source build). This changes nothing about
   how your files are protected.
2. **Set up the daemon** (`sudo app-listener install`). This is interactive and is never done
   automatically, because it encrypts directories and starts enforcing.

The standalone modes (`monitor`, `guard`, `network-monitor`, `network-guard`) need only step 1.

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

1. Stops a running daemon.
2. Generates the fscrypt master key (an existing key is kept).
3. Lets you pick the users to protect.
4. Probes the [built-in catalog](catalog.md) of sensitive directories (SSH, GPG, AI agents,
   browsers, VPNs, password stores, crypto wallets, ...) and lets you pick which to guard.
5. Encrypts the selected directories. Each is backed up first and the backup is verified against
   the master key.
6. Deploys the systemd units, the binary and the config.
7. Runs `daemon --check` on the deployed binary before enabling the service: if this kernel's
   verifier rejects the guard eBPF, the service is not enabled.

Earlier versions installed pacman/apt hooks and a boot unit to refresh the catalog; the wizard
removes them, since the daemon now refreshes the catalog itself.

### SSH agent

When a user's `~/.ssh` ends up guarded, the wizard offers (one question per user, naming the user
and the unit path, asked before encryption):

- a per-user `ssh-agent` systemd unit,
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

### fscrypt prerequisites

Each filesystem must be initialized for fscrypt (`fscrypt setup --all-users`) and support
encryption (ext4: `tune2fs -O encrypt <dev>`). The wizard checks this before migrating anything.
For a fixable gap it shows the exact command and why, and offers to run it for you; decline and
it aborts. Every command it may run is listed in `internal/fscrypt/prereq.go`.

### Install flags

| Flag | Description |
|---|---|
| `--diff-catalog` | Incremental wizard: lists the catalog directories that exist on the host but are **not** in `daemon.conf` yet, lets you pick which to add, appends them and encrypts them (backup first). Existing sections stay byte-for-byte intact. Stops the daemon for the cycle and restarts it on the merged config. Requires a previous installation; interactive only. |
| `--update-catalog-only` | Re-expands every catalog-matched whitelist and rewrites the config, dropping every line the catalog does not produce. Stops the daemon and unlocks each vault under an ephemeral self-only guard. The running daemon already does this on its own (see [catalog refresh](catalog.md#catalog-refresh)), so you normally never need it. |
| `--live` | With `--update-catalog-only`: refresh without stopping the daemon (no lock churn), applied via SIGHUP. Requires the daemon to be running. |
| `--allow-metadata-output` | Installs the daemon unit **without** `--no-log-metadata-blocks`, so denied metadata-only process inspections (`op=PTRACE mode=READ`) are logged. Use it to diagnose an app that breaks with nothing logged; re-run `install` without it to quiet them again. |
| `--binary-only` | Moves the freshly built binary to the install path and recreates the PATH symlink. No wizard, config, fscrypt or units. If the daemon service is installed, asks whether to restart it now. |
| `--restore-backups` | Undoes the fscrypt migration: lists the `.app_listener.backup` directories, and after one confirmation deletes the encrypted copies and moves the backups back. Refused while the daemon runs. |
| `--delete-post-backups` | Deletes the `.app_listener.backup` directories (listed, confirmed once). |
| `-y, --yes` | Skips confirmation prompts (for `--update-catalog-only` in scripts). |

`--update-catalog-only` refreshes the whitelists of existing sections; `--diff-catalog` adds new
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

## Uninstall

```bash
sudo app-listener uninstall
sudo app-listener uninstall --delete-key
```

- Refuses while the daemon runs.
- Re-scans the catalog and decrypts every protected directory in place.
- Keeps the master key unless `--delete-key` is given. Without the key, any directory still
  encrypted can never be unlocked again.
- At the end, lists any `.app_listener.backup` migration copies and offers to delete them (all
  preselected, one confirmation). They are plain unencrypted copies and harmless to keep.

## Docker

```bash
docker compose build
docker compose run --rm app-listener monitor -w /tmp
```

Multi-stage build on a Debian-slim runner. eBPF needs the host kernel: run it privileged, or with
`/sys/kernel/btf/` mounted and `CAP_BPF`. `make build` uses a separate toolchain image
(`docker/builder.Dockerfile`), not this one.
