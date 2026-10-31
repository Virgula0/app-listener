# Self-updating apps and the catalog

## The catalog

The install wizard doesn't ask you to know where every app keeps its secrets. It ships a built-in
catalog (`internal/install/catalog.go`) of sensitive directories with a ready-made whitelist for
each, and offers the ones that exist on your host. It covers, among others:

- **Keys and credentials:** SSH, GnuPG (and the `ssh-agent` and `gpg-agent` sockets in
  `/run/user/<uid>`, guarded as [runtime sections](daemon-config.md#runtime-sections)), age, AWS, Google Cloud, Azure, Kubernetes, Docker,
  GitHub CLI, adb keys, rclone, pip, Git
- **AI agents and editors:** Claude Code, Gemini CLI and Antigravity (CLI, 2.0 and IDE), OpenAI
  Codex CLI, pi, opencode, Codeium, Cursor, GitHub Copilot (and Copilot CLI), VS Code (and Insiders, server, VSCodium), JetBrains IDEs, Zed, Sublime Text
- **API and database clients:** Insomnia, Postman, DBeaver
- **Password managers:** pass, gopass, Bitwarden, 1Password, KeePassXC, GNOME keyring, KDE Wallet
- **Browsers:** Firefox, Chrome, Chromium, Brave
- **Messaging and apps:** Discord, Telegram, Signal, Element, Steam
- **VPNs:** NordVPN, Mullvad, ProtonVPN, OpenVPN, WireGuard
- **Crypto wallets:** Bitcoin Core and Litecoin Core (wallets and RPC cookie), Monero, Electrum, Sparrow, Wasabi, Ledger Live,
  Trezor Suite, Exodus, Solana CLI, geth, LND, Core Lightning, Foundry
- **Keyring items** ([secret scopes](daemon-config.md#secret-service-sections)): the `gh`, Codex
  and Antigravity CLI tokens, `git`'s libsecret credentials and the Chrome and Chromium Safe Storage keys, plus gnome-keyring's
  sockets as a runtime section. Offered when gnome-keyring and one of the binaries are installed.
- **Loopback services** ([port sections](daemon-config.md#port-sections)): Bitcoin Core and
  Litecoin Core RPC (mainnet, testnet3 and regtest ports), the `adb` server, the Ollama API,
  `opencode serve`. Offered when one of their binaries is installed.

The catalog narrows each watch to the sensitive subtrees, so apps that update themselves inside
their own config directory keep working (see below).

### What an entry can't cover

- **Tokens in the keyring.** These apps keep their main token in the Secret Service keyring by
  default, so their entry protects only the files: `gh`, Copilot CLI, Gemini CLI, JetBrains IDEs, Zed,
  ProtonVPN, and anything stored in the GNOME keyring or KDE Wallet. Any of your programs can ask
  the keyring ([Secrets served over IPC](limitations.md#secrets-served-over-ipc)) unless a secret
  scope covers the item: the catalog has scopes for `gh`, Codex, the Antigravity CLI, git and
  Chromium browsers only, on gnome-keyring. Codex's keyring mode also caches its token in the
  session's kernel keyring, which no scope covers: keep `cli_auth_credentials_store = "file"` (the
  default) and the guarded `auth.json`. `install catalog refresh` leaves secret sections as written. The installer
  says when `gh` or `docker` could keep the token in their guarded file instead
  (`gh auth login --insecure-storage`; no `credsStore` in `~/.docker/config.json`), and names your
  git credential helper.
- **Script launchers.** A `#!` script runs as its interpreter, so a whitelist line naming one never
  matches, and the app is denied like any other program. Google Cloud SDK, Azure CLI, pip, Electrum
  and ProtonVPN are Python apps, and Gemini CLI and pi (from npm or pi's installer) are Node apps:
  their entries lock the guarded files against everyone, the app included, unless you whitelist an
  ELF they end up running (pi's release tarball ships one). Never whitelist the interpreter. The installer and
  the daemon warn about every script line.
- **opencode snapshots.** opencode 2 keeps provider credentials in
  `~/.local/share/opencode/opencode.db`, which the entry guards. opencode runs `git` on its snapshot
  repositories in the same directory, and `git` is never whitelisted (it would print the database), so
  snapshots fail while the directory is guarded.
- **Cursor logins.** The Cursor entry guards `~/.cursor` (MCP and extension config). The IDE keeps
  its login in `~/.config/Cursor/User/globalStorage/state.vscdb` and `cursor-agent` in
  `~/.config/cursor/auth.json`; neither is in the catalog. Add a section for them yourself,
  whitelisting the executable that actually runs, not a launcher script. `cursor-agent` runs on a
  bundled `node`, which is an interpreter: whitelisting it admits any script.
- **`~/.gemini` is shared.** Gemini CLI and every Antigravity product (`agy`, the 2.0 hub's
  `language_server`, the IDE) keep state there, so one entry covers it and each whitelisted binary
  can read the others' tokens. It guards the token files and the config the agents run (MCP servers,
  hooks, plugins), not the whole directory. The binaries are matched at the paths the Arch packages
  and `agy`'s installer use (`/opt/Antigravity`, `/opt/antigravity-ide`, `~/.local/bin/agy`); add
  yours if you unpacked the tarball elsewhere.
- **Files created after install.** The Codex, pi and `~/.gemini` entries guard named files
  (`auth.json`, `oauth_creds.json`, ...) inside their directory, and only those that exist when the
  section is written; `install catalog refresh` updates whitelists, never `watch:` lines. Log in to
  the app before you install, or add the missing `watch: <path>` line later with
  `sudo app-listener edit-protected config` ([several guarded trees](daemon-config.md#several-guarded-trees-one-vault)).
  If none of them exists yet, the section guards the whole directory, which breaks Codex's
  standalone updates in `~/.codex/packages`.
- **Copilot CLI `/ide`.** The VS Code Copilot extension writes lock files to `~/.copilot/ide`.
  Whitelist the VS Code binary on `~/.copilot` for that integration.
- **Bitcoin and Litecoin RPC cookie.** The cookie grants every RPC call, including spending from an
  unlocked wallet. By default the node writes it to the datadir root, outside any guarded tree, as
  `.cookie.tmp`, then renames it. To guard it, create `~/.bitcoin/rpc` (`~/.litecoin/rpc`), add
  `rpccookiefile=rpc/.cookie` to the node's config, and run `sudo app-listener install catalog diff`:
  the "RPC cookie" entry whitelists the node, its GUI and its CLI. Whitelist there any other RPC client
  that reads the cookie (electrs, LND, Core Lightning's `bcli` runs `bitcoin-cli`), not on
  `wallets`. The node's config file stays writable, so malware can add an `rpcauth` login of its own
  and wait for a restart.

- **Other clients of a guarded port.** A port entry whitelists the service and its own client: a
  web UI for Ollama, electrs or LND talking to `bitcoind`, an IDE plugin talking to `opencode serve`
  needs its own `CONNECT` line. Signet and testnet4 RPC ports sit in the default ephemeral range,
  which the daemon refuses, and `install catalog refresh` leaves port sections as written.

To add catalog directories that appeared after you installed (a newly installed app), run
`sudo app-listener install catalog diff`.

## Catalog refresh

The daemon keeps the catalog whitelists current by itself. It watches the catalog's directories
(inotify, as a hint only) and re-expands the catalog whitelists at startup and whenever a binary
is written below a catalog pattern, with the same rules as `install catalog refresh`,
keeping lines whose binary still exists. In the app's `[libraries]` block it also drops the
`lib_dir` and `lib_binary` lines a catalog pattern produced whose path is gone (an old version
directory the updater removed). Hand-written lines stay, and nothing is dropped from an encrypted
section that shows no library directory at all (its vault looks locked).

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
  refresh as soon as the updater finishes writing the binary, also when it deletes the old version
  first. No restart, no manual step.
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
(Steam's `ubuntu12_*`, `linux*`, `steamrt64`, `compatibilitytools.d`, and `steamui`, the
JavaScript `steamwebhelper` runs) and trusts their libraries for that app. So that a non-app process can't plant one ahead of it, the daemon reserves those dir
names for the app's own binaries, present or not: only Steam can create
`~/.local/share/Steam/compatibilitytools.d`, not your shell or ProtonUp-Qt.

`steam.sh` unpacks and pins Steam's runtimes inside those dirs with general tools the read-only
dirs refuse, so the daemon runs it sandboxed and commits only what it validates (see
[App maintainers](daemon.md#app-maintainers-steam-runtime-updates)). It is a reserved name too:
only Steam's binaries may write it.

There is no npm entry: npm runs as `node`, and whitelisting `node` would let `node -e` read the
tree.
