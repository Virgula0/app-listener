# Troubleshooting

## Is the daemon running and enforcing?

```bash
systemctl is-active app-listener-daemon              # expect: active
journalctl -u app-listener-daemon -f                 # follow live
sudo journalctl -u app-listener-daemon -f | grep -i denied
sudo fscrypt status /home/alice/.ssh                 # Encrypted / Not encrypted
```

## Check the guard by hand

```bash
sudo systemctl stop app-listener-daemon
sudo /usr/local/sbin/app-listener daemon --headless --verbose 3
# in another terminal:
#   ssh -T git@github.com    must WORK
#   cat ~/.ssh/id_ed25519    must be DENIED
```

`--headless --blocked-only` prints only denials.

## A program I use is denied

1. Find the denial in the journal; it names the program and the resource.
2. If the program should have access, add its **real executable** (not a wrapper script) to the
   resource's section with `sudo app-listener edit-protected --edit-config`
   ([how to find it](how-it-works.md#identity-is-the-executables-inode)). For a one-off, use
   `edit-protected --forward -w <binary>` instead.
3. If the journal says `binary-unconfirmed` or `not re-admitted`, the binary changed while it
   wasn't watched, or was replaced by something other than its updater. Check it, then
   `sudo app-listener trust-binaries` ([Binary trust](binary-trust.md#trust-binaries)).
4. `TRUST DENIED op=LAUNCH` means the app was started with an environment variable or flag that
   loads foreign code ([launch scan](binary-trust.md#launch-scan)). Start it without.

## After a package update replaced a binary

Root-placed binaries replaced by a package upgrade are re-admitted automatically. If something
still looks stale:

```bash
sudo systemctl reload app-listener-daemon
```

## An app breaks and nothing is logged

Usually a denied metadata-only inspection of another process (`op=PTRACE mode=READ`): a component
reading `/proc/<pid>` of a process that holds protected content. The systemd unit doesn't log
these by default ([details](daemon.md#quiet-metadata-denials)). Reinstall with them visible:

```bash
sudo app-listener install --allow-metadata-output
```

Re-run `install` without the flag to quiet them again. If the reader is a desktop service that
legitimately needs every process's metadata (a screen-sharing portal), list it under
[`[inspectors]`](daemon-config.md#process-inspectors).

## Snapshots or disk backups stopped working

Expected on filesystems that host a guarded resource, see
[Filesystem-wide gates](daemon.md#filesystem-wide-gates).

## Encryption key errors

A wrong or old master key fails immediately with "invalid wrapping key". The daemon never silently
generates a new one.

## Backups

The install wizard backs up every directory before encrypting it, at `<dir>.app_listener.backup`.

```bash
sudo app-listener install --restore-backups        # undo the migration from the backups
sudo app-listener install --delete-post-backups    # delete the backups
```

`uninstall` also offers to delete them at the end.
