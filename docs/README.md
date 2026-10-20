# app-listener documentation

The [project README](../README.md) is the short tour. These pages are the full reference: every
subcommand, flag, config option and known limitation.

## Getting started

- [Compatibility](compatibility.md): supported distributions, kernel versions, enabling BPF-LSM.
- [Installation](installation.md): the one-line installer, building from source, the install
  wizard, `update`, `uninstall`, Docker.
- [How it works](how-it-works.md): where it sits next to AppArmor/SELinux, the eBPF attach model, binary identity, output formats, `--serve`.
- [Desktop GUI](gui.md): `app-listener-gui`, the window covering the file and network modes, the updater and the installer.

## Modes

| Mode | What it does | Page |
|---|---|---|
| `monitor` | Shows every file operation under a path. Blocks nothing. | [monitor](monitor.md) |
| `guard` | Blocks file access to a path, by program. | [guard](guard.md) |
| `network monitor` | Shows the network activity of chosen programs. | [network monitor](network-monitor.md) |
| `network guard` | Blocks network access, by program. | [network guard](network-guard.md) |
| `daemon` | Protects many directories permanently, with encryption at rest. | [daemon](daemon.md) |

## Command tree

Every command takes `--help`. Options that only make sense together live on their own
subcommand, so `--help` at each level lists only what applies there.

```text
app-listener
├── monitor PATH...                     observe file events          monitor.md
├── guard PATH [-w BIN... | -b BIN...]  deny file access             guard.md
├── network (alias net)
│   ├── monitor BIN...                  observe network events       network-monitor.md
│   └── guard [-w BIN... | -b BIN...]   deny network access          network-guard.md
├── daemon                              run the daemon               daemon.md
│   ├── check [--verifier-only]         load the eBPF into the verifier
│   ├── lockdown                        lock every vault and exit
│   └── genkey                          create the fscrypt master key
├── edit-protected                      edit a protected directory   edit-protected.md
│   ├── forward (-w BIN | -b BIN)...    temporary access for other binaries
│   ├── config [--from FILE]            edit the running daemon.conf
│   │   └── record PATH...              learn who uses PATH, then guard it
│   ├── put FILE --resource DIR         write one file, no editor
│   └── password                        change the edit-protected password
├── trust                               confirm refused binaries     binary-trust.md
│   ├── list
│   └── add PATH...
├── install                             the install wizard           installation.md
│   ├── binary                          deploy only the binary
│   ├── catalog refresh [--live]        re-expand the catalog whitelists
│   ├── catalog diff                    guard new catalog directories
│   ├── backups restore                 undo the encryption
│   ├── backups delete
│   └── facultatives                    install optional per-user setup
├── uninstall
│   └── facultatives                    remove optional per-user setup
└── update [--channel C]
```

The desktop window is the separate `app-listener-gui` binary: `sudo app-listener-gui`.

## The daemon in depth

- [Daemon](daemon.md): flags, lifecycle, boot ordering, reloads, encryption, filesystem-wide gates.
- [Daemon configuration](daemon-config.md): `daemon.conf` grammar, event masks, `[electron_apps]`,
  `[inspectors]`.
- [Binary trust](binary-trust.md): how whitelisted binaries stay trustworthy across updates, the
  binary ledger, `trust`, launch scan, multicall binaries.
- [Self-updating apps and the catalog](catalog.md): the built-in list of sensitive directories and
  how apps like Discord or VS Code keep updating while guarded.
- [edit-protected](edit-protected.md): changing protected files and the daemon config safely.
- [Limitations](limitations.md): what the daemon cannot stop, and side effects to expect.

## Operating

- [Troubleshooting](troubleshooting.md): checking the daemon, reading denials, apps that break
  silently, backups.
- [Development](development.md): Makefile targets, tests, Docker images.
