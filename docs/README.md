# app-listener documentation

The [project README](../README.md) is the short tour. These pages are the full reference: every
subcommand, flag, config option and known limitation.

## Getting started

- [Compatibility](compatibility.md): supported distributions, kernel versions, enabling BPF-LSM.
- [Installation](installation.md): the one-line installer, building from source, the install
  wizard, `update`, `uninstall`, Docker.
- [How it works](how-it-works.md): the eBPF attach model, binary identity, output formats, `--serve`.

## Modes

| Mode | What it does | Page |
|---|---|---|
| `monitor` | Shows every file operation under a path. Blocks nothing. | [monitor](monitor.md) |
| `guard` | Blocks file access to a path, by program. | [guard](guard.md) |
| `network-monitor` | Shows the network activity of chosen programs. | [network-monitor](network-monitor.md) |
| `network-guard` | Blocks network access, by program. | [network-guard](network-guard.md) |
| `daemon` | Protects many directories permanently, with encryption at rest. | [daemon](daemon.md) |

## The daemon in depth

- [Daemon](daemon.md): flags, lifecycle, boot ordering, reloads, encryption, filesystem-wide gates.
- [Daemon configuration](daemon-config.md): `daemon.conf` grammar, event masks, `[electron_apps]`,
  `[inspectors]`.
- [Binary trust](binary-trust.md): how whitelisted binaries stay trustworthy across updates, the
  binary ledger, `trust-binaries`, launch scan, multicall binaries.
- [Self-updating apps and the catalog](catalog.md): the built-in list of sensitive directories and
  how apps like Discord or VS Code keep updating while guarded.
- [edit-protected](edit-protected.md): changing protected files and the daemon config safely.
- [Limitations](limitations.md): what the daemon cannot stop, and side effects to expect.

## Operating

- [Troubleshooting](troubleshooting.md): checking the daemon, reading denials, apps that break
  silently, backups.
- [Development](development.md): Makefile targets, tests, Docker images.
