# Self-updating apps and the catalog

## The catalog

The install wizard doesn't ask you to know where every app keeps its secrets. It ships a built-in
catalog (`internal/install/catalog.go`) of sensitive directories with a ready-made whitelist for
each, and offers the ones that exist on your host. It covers, among others:

- **Keys and credentials:** SSH, GnuPG, age, AWS, Google Cloud, Azure, Kubernetes, Docker,
  GitHub CLI, adb keys, rclone, pip, Git
- **AI agents and editors:** Claude Code, Gemini CLI, opencode, Codeium, Cursor, GitHub Copilot,
  VS Code (and Insiders, server, VSCodium), JetBrains IDEs, Zed, Sublime Text
- **API and database clients:** Insomnia, Postman, DBeaver
- **Password managers:** pass, gopass, Bitwarden, 1Password, KeePassXC, GNOME keyring, KDE Wallet
- **Browsers:** Firefox, Chrome, Chromium, Brave
- **Messaging and apps:** Discord, Telegram, Signal, Element, Steam
- **VPNs:** NordVPN, Mullvad, ProtonVPN, OpenVPN, WireGuard
- **Crypto wallets:** Bitcoin Core, Litecoin Core, Monero, Electrum, Sparrow, Wasabi, Ledger Live,
  Trezor Suite, Exodus, Solana CLI, geth, LND, Core Lightning, Foundry

The catalog narrows each watch to the sensitive subtrees, so apps that update themselves inside
their own config directory keep working (see below).

To add catalog directories that appeared after you installed (a newly installed app), run
`sudo app-listener install --diff-catalog`.

## Catalog refresh

The daemon keeps the catalog whitelists current by itself. It watches the catalog's directories
(inotify, as a hint only) and re-expands the catalog whitelists at startup and whenever a binary
is written below a catalog pattern, with the same rules as `install --update-catalog-only`,
keeping lines whose binary still exists.

`daemon.conf` is rewritten only while it still holds the configuration the daemon runs, then
reloaded in-process. The journal shows:

```
DAEMON catalog-refresh resource=... admitted=... dropped=...
```

No package-manager hook or boot unit is involved.

## Self-updating apps

Discord, VS Code helpers and similar apps change their own binaries outside the package manager.
The daemon follows them:

- An update written in place of the binary by the app's own updater is re-admitted directly.
- A new version directory (Discord's `~/.config/discord/0.0.N/`) is admitted by the catalog
  refresh as soon as the updater finishes writing the binary. No restart, no manual step.
- A process that creates such a directory without being one of the app's binaries can't plant the
  binary name there (the names are reserved), so nothing it does is admitted.

Only the sensitive subtrees are guarded (Discord's `Local Storage/`, `Cookies`, ...); the root the
updater writes to stays unguarded. The code these apps run is not left writable either. These are
read-only `lib_dir`s that only the app's own binaries may change:

- Discord's `app-*` version dirs (inside its vault, guarded once unlocked)
- `~/.vscode`, `~/.vscode-insiders`, `~/.vscode-oss` (extensions)
- `~/.local/share/JetBrains` (plugins)
- Steam's update staging `~/.local/share/Steam/package`

Discord's updater inputs (`settings.json`, which names the update endpoint, and `installer.db`)
are reserved names only Discord's binaries may write, so no other process can make the genuine
updater install chosen code.

The refresh also adopts an app's fixed library dirs as read-only `lib_dir`s once they exist
(Steam's `ubuntu12_*`, `linux*`, `steamrt64`, `compatibilitytools.d`) and trusts their libraries
for that app. So that a non-app process can't plant one ahead of it, the daemon reserves those dir
names for the app's own binaries, present or not: only Steam can create
`~/.local/share/Steam/compatibilitytools.d`, not your shell or ProtonUp-Qt.

There is no npm entry: npm runs as `node`, and whitelisting `node` would let `node -e` read the
tree.
