// Package install implements the building blocks of the `app-listener
// install` wizard: the critical-directory catalog, user enumeration,
// fscrypt tree migration and daemon config generation.
package install

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CandidateDir describes one critical directory type the installer probes.
type CandidateDir struct {
	// Name is a short human-readable label shown in the TUI. It identifies
	// the resource: two locations of the same application (Steam's
	// ".local/share/Steam/config" and ".steam") share one entry, one Name.
	Name string
	// RelPaths are the home-relative locations of this resource (e.g.
	// ".ssh"). On a given host more than one may exist; each existing one is
	// protected as its own watch + fscrypt root, all sharing this entry's
	// Whitelist. Set exactly one of RelPaths and AbsPaths, non-empty.
	RelPaths []string
	// AbsPaths are system-level locations (e.g. "/etc/wireguard") probed
	// once, independent of the selected users. Same multi-location
	// semantics as RelPaths.
	AbsPaths []string
	// Whitelist maps binary paths to allowed events (nil = all events,
	// emitted bare; otherwise "<path> EV1,EV2"); matching is by inode.
	// Keep it minimal: each entry is a potential privilege-escalation path
	// if abused while the directory is unlocked. When an entry has several
	// RelPaths the whitelist is shared by all of them (the union).
	Whitelist map[string][]string
	// WatchRelPaths optionally replaces the single [watch <RelPath>] section
	// with MULTIPLE guarded subtrees inside that RelPath, all sharing this
	// entry's whitelist and fscrypt lifecycle (the section path stays the
	// encryption root). Used for self-updating applications whose update
	// workspace must stay outside the guarded set while the user data
	// remains protected. Empty means the plain per-RelPath watch (historical
	// behavior). Only valid together with exactly one RelPaths entry. The
	// paths are relative to the user's home, supporting the same %HOME%
	// expansion as whitelist entries.
	WatchRelPaths []string
}

// IsSystem reports whether this is a system-level (AbsPaths) entry, probed
// once regardless of the selected users.
func (c *CandidateDir) IsSystem() bool { return len(c.AbsPaths) > 0 }

// Catalog is the master list of critical directories probed for each selected
// user — tweak entries here to drive discovery. Non-existent paths are
// simply not proposed.
var Catalog = []CandidateDir{
	// --- SSH and remote access -------------------------------------------------
	{Name: "SSH client configuration and keys", RelPaths: []string{".ssh"},
		// ssh needs WRITE (known_hosts rotation); sshd and its helpers
		// (sshd-auth/-session, sftp-server) only READ authorized_keys; the
		// keysign/pkcs11/cleanup helpers are inert. git excluded: reads keys.
		Whitelist: map[string][]string{
			"/usr/bin/ssh":                         {"READ", "WRITE", "DELETE", "RENAME", "HARDLINK"},
			"/usr/bin/ssh-add":                     nil,
			"/usr/bin/ssh-agent":                   nil,
			"/usr/bin/ssh-keygen":                  nil,
			"/usr/bin/scp":                         nil,
			"/usr/bin/sftp":                        nil,
			"/usr/bin/sshd":                        {"READ"},
			"/usr/sbin/sshd":                       {"READ"},
			"/usr/lib/openssh/sshd-auth":           {"READ"},
			"/usr/lib/ssh/sshd-auth":               {"READ"},
			"/usr/lib/openssh/sshd-session":        {"READ"},
			"/usr/lib/ssh/sshd-session":            {"READ"},
			"/usr/lib/openssh/ssh-sk-helper":       {"READ"},
			"/usr/lib/ssh/ssh-sk-helper":           {"READ"},
			"/usr/lib/openssh/sftp-server":         {"READ"},
			"/usr/lib/ssh/sftp-server":             {"READ"},
			"/usr/lib/openssh/ssh-keysign":         nil,
			"/usr/lib/ssh/ssh-keysign":             nil,
			"/usr/lib/openssh/ssh-pkcs11-helper":   nil,
			"/usr/lib/ssh/ssh-pkcs11-helper":       nil,
			"/usr/lib/openssh/ssh-session-cleanup": nil,
		}},
	{Name: "GNU Privacy Guard keyring", RelPaths: []string{".gnupg"},
		Whitelist: map[string][]string{
			"/usr/bin/gpg": nil, "/usr/bin/gpg-agent": nil,
			"/usr/bin/gpgconf": nil, "/usr/bin/gpg-connect-agent": nil,
		}},

	// --- AI coding agents and CLI tools ----------------------------------------
	{Name: "opencode", RelPaths: []string{".config/opencode"},
		Whitelist: map[string][]string{
			"/usr/local/bin/opencode": nil, "/usr/bin/opencode": nil,
			"%HOME%/.local/bin/opencode": nil,
		}},
	{Name: "code CLI (GitHub)", RelPaths: []string{".config/code-cli"},
		Whitelist: map[string][]string{
			"/usr/bin/code-cli": nil, "/usr/local/bin/code-cli": nil,
			"%HOME%/.local/bin/code-cli": nil,
		}},
	// One resource, two locations: the CLI state (~/.claude) and its XDG
	// config (~/.config/claude). Same whitelist, each its own watch root.
	{Name: "Claude Code", RelPaths: []string{".claude", ".config/claude"},
		Whitelist: map[string][]string{
			"/usr/local/bin/claude":    nil,
			"/usr/bin/claude":          nil,
			"%HOME%/.local/bin/claude": nil,
		}},
	{Name: "Codeium", RelPaths: []string{".codeium"},
		Whitelist: map[string][]string{
			"/usr/local/bin/codeium": nil, "/usr/bin/codeium": nil,
			"%HOME%/.local/bin/codeium": nil,
		}},
	{Name: "Gemini CLI", RelPaths: []string{".gemini"},
		Whitelist: map[string][]string{
			"/usr/local/bin/gemini": nil, "/usr/bin/gemini": nil,
			"%HOME%/.local/bin/gemini": nil,
		}},
	{Name: "Cursor", RelPaths: []string{".cursor"},
		Whitelist: map[string][]string{
			"/usr/bin/cursor": nil, "/usr/local/bin/cursor": nil,
			"%HOME%/.local/bin/cursor": nil,
		}},
	{Name: "GitHub Copilot", RelPaths: []string{".config/github-copilot"},
		Whitelist: map[string][]string{
			"/usr/bin/git": nil, "/usr/local/bin/copilot": nil,
			"%HOME%/.local/bin/copilot": nil,
		}},

	// --- IDEs and editors ---------------------------------------------------------
	// IDE configs hold auth tokens and credentials; wrapper binaries
	// (/usr/bin/code, firefox, ...) are shell scripts, so the whitelist
	// matches the executed ELF under /opt and /usr/lib.
	{Name: "VS Code config", RelPaths: []string{".config/Code"},
		Whitelist: map[string][]string{
			"/usr/bin/code": nil, "/usr/local/bin/code": nil,
			"%HOME%/.local/bin/code":       nil,
			"/opt/visual-studio-code/code": nil, "/usr/share/code/code": nil, "/opt/visual-studio-code/bin/code": nil,
			// Code's separate crash-dump writer process.
			"/opt/visual-studio-code/chrome_crashpad_handler": nil, "/usr/share/code/chrome_crashpad_handler": nil,
			// Code's extension-signature verifier, spawned per extension
			// install/update to check CachedExtensionVSIXs/*.sigzip files.
			"/opt/visual-studio-code/resources/app/node_modules/@vscode/vsce-sign/bin/vsce-sign": nil,
			"/usr/share/code/resources/app/node_modules/@vscode/vsce-sign/bin/vsce-sign":         nil,
		}},
	{Name: "VS Code Insiders config", RelPaths: []string{".config/Code - Insiders"},
		Whitelist: map[string][]string{
			"/usr/bin/code-insiders": nil, "/usr/local/bin/code-insiders": nil,
			"%HOME%/.local/bin/code-insiders": nil, "/opt/visual-studio-code/bin/code": nil,
		}},
	{Name: "VS Code server (remote development)", RelPaths: []string{".vscode-server"},
		Whitelist: map[string][]string{
			"/usr/bin/code": nil, "/usr/local/bin/code": nil,
			"%HOME%/.local/bin/code":       nil,
			"/opt/visual-studio-code/code": nil, "/usr/share/code/code": nil,
			// Code's separate crash-dump writer process.
			"/opt/visual-studio-code/chrome_crashpad_handler": nil, "/usr/share/code/chrome_crashpad_handler": nil,
		}},
	{Name: "VSCodium config", RelPaths: []string{".config/VSCodium"},
		Whitelist: map[string][]string{"/usr/bin/codium": nil, "/usr/local/bin/codium": nil,
			"/opt/vscodium/chrome_crashpad_handler": nil, "/usr/share/vscodium/chrome_crashpad_handler": nil}},
	{Name: "JetBrains IDEs", RelPaths: []string{".config/JetBrains"},
		// IDEs run as the bundled JBR java, so the JBR and its fsnotifier
		// helper are whitelisted (globs cover /opt installs). Toolbox apps
		// stay unguarded: no credentials — those live in .config/JetBrains.
		Whitelist: map[string][]string{
			"/usr/bin/idea": nil, "/usr/bin/pycharm": nil, "/usr/bin/webstorm": nil,
			"/usr/bin/clion": nil, "/usr/bin/goland": nil, "/usr/bin/phpstorm": nil,
			"/usr/bin/rider": nil, "/usr/bin/datagrip": nil, "/usr/bin/rustrover": nil,
			"/usr/bin/android-studio":       nil,
			"/opt/*/jbr/bin/java":           nil,
			"/opt/*/bin/fsnotifier":         nil,
			"%HOME%/.goland/jbr/bin/java":   nil,
			"%HOME%/.goland/bin/fsnotifier": nil,
		}},
	// Config (~/.config/zed) and data (~/.local/share/zed, holds the auth
	// token) — one resource, one whitelist, two watch roots.
	{Name: "Zed editor", RelPaths: []string{".config/zed", ".local/share/zed"},
		Whitelist: map[string][]string{
			"/usr/bin/zed": nil, "/usr/local/bin/zed": nil,
			"%HOME%/.local/bin/zed": nil,
		}},
	{Name: "Sublime Text", RelPaths: []string{".config/sublime-text"},
		// /usr/bin/subl is a wrapper; the real binary is under /opt.
		Whitelist: map[string][]string{"/usr/bin/subl": nil, "/usr/bin/sublime-text": nil, "/usr/bin/sublime_text": nil,
			"/opt/sublime_text/sublime_text": nil}},
	{Name: "Insomnia API client", RelPaths: []string{".config/Insomnia"},
		Whitelist: map[string][]string{"/usr/bin/insomnia": nil, "/opt/insomnia/insomnia": nil}},
	{Name: "Postman API client", RelPaths: []string{".config/Postman"},
		Whitelist: map[string][]string{"/usr/bin/postman": nil, "/opt/postman/postman": nil}},
	{Name: "DBeaver database client", RelPaths: []string{".local/share/DBeaverData"},
		// /usr/bin/dbeaver is a wrapper; the real binary lives in /usr/lib.
		Whitelist: map[string][]string{"/usr/bin/dbeaver": nil, "/usr/lib/dbeaver/dbeaver": nil}},
	{Name: "Android adb keys", RelPaths: []string{".android"},
		Whitelist: map[string][]string{"/usr/bin/adb": nil}},

	// --- Cloud, containers and dev tooling --------------------------------------
	{Name: "AWS credentials", RelPaths: []string{".aws"},
		Whitelist: map[string][]string{"/usr/bin/aws": nil}},
	{Name: "Google Cloud SDK", RelPaths: []string{".config/gcloud"},
		// Python launchers: the exe is the interpreter (not whitelisted),
		// so both entries are inert — kept for documentation.
		Whitelist: map[string][]string{"/usr/bin/gcloud": nil, "/usr/bin/gsutil": nil}},
	{Name: "Kubernetes kubeconfig", RelPaths: []string{".kube"},
		Whitelist: map[string][]string{"/usr/bin/kubectl": nil, "/usr/bin/helm": nil, "/usr/bin/oc": nil,
			"/usr/lib/docker/cli-plugins/docker-buildx":     nil,
			"/usr/libexec/docker/cli-plugins/docker-buildx": nil,
			"%HOME%/.docker/cli-plugins/docker-buildx":      nil}},
	{Name: "Docker config and credentials", RelPaths: []string{".docker"},
		Whitelist: map[string][]string{"/usr/bin/docker": nil, "/usr/bin/docker-credential-desktop": nil, "/usr/bin/docker-credential-pass": nil,
			"/usr/lib/docker/cli-plugins/docker-buildx":     nil,
			"/usr/libexec/docker/cli-plugins/docker-buildx": nil,
			"%HOME%/.docker/cli-plugins/docker-buildx":      nil}},
	{Name: "GitHub CLI", RelPaths: []string{".config/gh"},
		Whitelist: map[string][]string{
			"/usr/bin/gh":                       nil,
			"/usr/local/bin/gh":                 nil,
			"%HOME%/.local/bin/gh":              nil,
			"/home/linuxbrew/.linuxbrew/bin/gh": nil,
			"/snap/bin/gh":                      nil,
		}},
	{Name: "Azure CLI", RelPaths: []string{".azure", ".config/azure"},
		Whitelist: map[string][]string{"/usr/bin/az": nil}},
	{Name: "rclone config", RelPaths: []string{".config/rclone"},
		Whitelist: map[string][]string{"/usr/bin/rclone": nil}},
	{Name: "Ollama config", RelPaths: []string{".ollama"},
		Whitelist: map[string][]string{"/usr/bin/ollama": nil}},
	{Name: "npm data", RelPaths: []string{".npm"},
		Whitelist: map[string][]string{"/usr/bin/npm": nil, "/usr/bin/npx": nil, "/usr/bin/node": nil}},
	// pip is a #!/usr/bin/python script (inert like gcloud); documentation only.
	{Name: "pip config", RelPaths: []string{".config/pip"},
		Whitelist: map[string][]string{"/usr/bin/pip": nil, "/usr/bin/pip3": nil}},
	{Name: "Git config directory", RelPaths: []string{".config/git"},
		Whitelist: map[string][]string{"/usr/bin/git": nil}},

	// --- Password managers and secrets -------------------------------------------
	{Name: "password-store (pass)", RelPaths: []string{".password-store"},
		Whitelist: map[string][]string{"/usr/bin/pass": nil, "/usr/bin/gpg": nil}},
	{Name: "gopass", RelPaths: []string{".local/share/gopass"},
		Whitelist: map[string][]string{"/usr/bin/gopass": nil}},
	{Name: "Bitwarden", RelPaths: []string{".config/Bitwarden"},
		Whitelist: map[string][]string{"/usr/bin/bitwarden": nil, "/usr/local/bin/bitwarden": nil}},
	{Name: "1Password", RelPaths: []string{".1password"},
		Whitelist: map[string][]string{"/usr/bin/1password": nil, "/usr/bin/op": nil, "/usr/local/bin/op": nil}},
	{Name: "age keys", RelPaths: []string{".config/age"},
		Whitelist: map[string][]string{"/usr/bin/age": nil, "/usr/bin/age-keygen": nil}},
	{Name: "KeePassXC", RelPaths: []string{".local/share/keepassxc"},
		Whitelist: map[string][]string{"/usr/bin/keepassxc": nil}},
	{Name: "GNOME keyring", RelPaths: []string{".local/share/keyrings", ".keyring"},
		Whitelist: map[string][]string{"/usr/bin/gnome-keyring-daemon": nil, "/usr/bin/gnome-keyring": nil}},
	{Name: "Chezmoi state", RelPaths: []string{".local/share/chezmoi"},
		Whitelist: map[string][]string{"/usr/bin/chezmoi": nil}},

	// --- Legacy dot-config files --------------------------------------------------
	// No single-file entries: the daemon guards directories only (non-dir
	// paths are skipped at load), so files like ~/.netrc would be silently
	// dropped. Steam below documents the same restriction for its ssfn files.

	// --- Wallets and crypto -------------------------------------------------------
	// Probed like any other entry: they surface only once the wallet is
	// installed. Solana keeps its keypair as a plaintext id.json.
	{Name: "Bitcoin Core", RelPaths: []string{".bitcoin/wallets"},
		Whitelist: map[string][]string{"/usr/bin/bitcoind": nil, "/usr/bin/bitcoin-qt": nil, "/usr/bin/bitcoin-cli": nil, "/usr/bin/bitcoin-tx": nil}},
	{Name: "Litecoin Core", RelPaths: []string{".litecoin/wallets"},
		Whitelist: map[string][]string{"/usr/bin/litecoind": nil, "/usr/bin/litecoin-qt": nil}},
	{Name: "Monero", RelPaths: []string{".bitmonero/wallets"},
		Whitelist: map[string][]string{"/usr/bin/monero-wallet-gui": nil, "/usr/bin/monero-wallet-cli": nil, "/usr/bin/monerod": nil}},
	{Name: "Electrum", RelPaths: []string{".electrum"},
		Whitelist: map[string][]string{"/usr/bin/electrum": nil, "/usr/local/bin/electrum": nil}},
	{Name: "Sparrow", RelPaths: []string{".sparrow"},
		Whitelist: map[string][]string{"/usr/bin/sparrow": nil}},
	{Name: "Wasabi Wallet", RelPaths: []string{".walletwasabi"},
		Whitelist: map[string][]string{"/usr/bin/wassabee": nil}},
	{Name: "Ledger Live", RelPaths: []string{".config/Ledger Live"},
		Whitelist: map[string][]string{"/usr/bin/ledger-live": nil, "/opt/Ledger Live/ledger-live": nil}},
	{Name: "Trezor Suite", RelPaths: []string{".config/@trezor"},
		Whitelist: map[string][]string{"/usr/bin/trezor-suite": nil, "/opt/Trezor Suite/trezor-suite": nil}},
	{Name: "Exodus", RelPaths: []string{".config/Exodus"},
		Whitelist: map[string][]string{"/usr/bin/exodus": nil, "/opt/Exodus/exodus": nil}},
	// id.json is a plaintext keypair — the highest-value target of the set.
	{Name: "Solana CLI", RelPaths: []string{".config/solana"},
		Whitelist: map[string][]string{"/usr/bin/solana": nil, "/usr/bin/solana-keygen": nil}},
	{Name: "Ethereum (geth)", RelPaths: []string{".ethereum/keystore"},
		Whitelist: map[string][]string{"/usr/bin/geth": nil}},
	{Name: "LND (Lightning Network Daemon)", RelPaths: []string{".lnd"},
		Whitelist: map[string][]string{"/usr/bin/lnd": nil, "/usr/bin/lncli": nil}},
	{Name: "Core Lightning", RelPaths: []string{".lightning"},
		Whitelist: map[string][]string{"/usr/bin/lightningd": nil, "/usr/bin/lightning-cli": nil}},
	// Binaries live inside the guarded dir next to the keystores: glob needed.
	{Name: "Foundry", RelPaths: []string{".foundry"},
		Whitelist: map[string][]string{"%HOME%/.foundry/bin/*": nil}},
	{Name: "KDE Wallet", RelPaths: []string{".local/share/kwalletd"},
		Whitelist: map[string][]string{"/usr/bin/kwalletd5": nil, "/usr/bin/kwalletd6": nil, "/usr/bin/kwalletmanager5": nil, "/usr/bin/kwalletmanager6": nil}},

	// --- VPNs ---------------------------------------------------------------------
	{Name: "NordVPN config", RelPaths: []string{".config/nordvpn"},
		Whitelist: map[string][]string{"/usr/bin/nordvpn": nil}},
	{Name: "Mullvad VPN config", RelPaths: []string{".config/mullvad"},
		Whitelist: map[string][]string{"/usr/bin/mullvad": nil}},
	{Name: "ProtonVPN config", RelPaths: []string{".config/protonvpn"},
		Whitelist: map[string][]string{"/usr/bin/protonvpn": nil, "/usr/bin/protonvpn-cli": nil}},
	{Name: "OpenVPN config", RelPaths: []string{".config/openvpn"},
		Whitelist: map[string][]string{"/usr/bin/openvpn": nil}},
	{Name: "WireGuard config (per-user)", RelPaths: []string{".config/wireguard"},
		Whitelist: map[string][]string{"/usr/bin/wg": nil, "/usr/bin/wg-quick": nil}},

	// --- Browsers and messaging -----------------------------------------------------
	// Browser binaries are shell wrappers; the whitelist matches the
	// executed ELF under /usr/lib and /opt. Firefox's crashhelper and the
	// chrome_crashpad_handler family write the Crash Reports dirs separately.
	{Name: "Firefox profile", RelPaths: []string{".mozilla/firefox", ".config/mozilla/firefox"},
		Whitelist: map[string][]string{
			"/usr/bin/firefox":             nil,
			"/usr/lib/firefox/firefox":     nil,
			"/usr/lib/firefox/crashhelper": nil, "/usr/lib/firefox/crashreporter": nil,
		}},
	{Name: "Google Chrome profile", RelPaths: []string{".config/google-chrome"},
		Whitelist: map[string][]string{
			"/usr/bin/google-chrome": nil, "/usr/bin/google-chrome-stable": nil,
			"/opt/google/chrome/chrome": nil, "/usr/lib/chromium/chrome": nil,
			"/opt/google/chrome/chrome_crashpad_handler": nil,
		}},
	{Name: "Chromium profile", RelPaths: []string{".config/chromium"},
		Whitelist: map[string][]string{
			"/usr/bin/chromium": nil, "/usr/bin/chromium-browser": nil,
			"/usr/lib/chromium/chrome": nil, "/usr/lib/chromium/chromium": nil,
			"/usr/lib/chromium/chrome_crashpad_handler": nil,
		}},
	{Name: "Brave profile", RelPaths: []string{".config/BraveSoftware"},
		Whitelist: map[string][]string{
			"/usr/bin/brave-browser": nil, "/usr/bin/brave": nil,
			"/opt/brave/brave": nil, "/opt/brave-bin/brave": nil,
			"/opt/brave-bin/chrome_crashpad_handler":         nil,
			"/usr/lib/brave-browser/chrome_crashpad_handler": nil,
		}},
	{Name: "Discord", RelPaths: []string{".config/discord"},
		// The Arch package is a launcher stub: the app self-updates into
		// versioned app-<ver> dirs and its wrapper/updater writes the vault
		// root (installer.db, downloads) — so only the SENSIBLE subtrees are
		// guarded and the update workspace stays out of the guarded set.
		WatchRelPaths: []string{
			".config/discord/Local Storage",
			".config/discord/Session Storage",
			".config/discord/IndexedDB",
			".config/discord/WebStorage",
			".config/discord/Service Worker",
			".config/discord/Cookies",
			".config/discord/Local State",
			".config/discord/Crashpad",
			".config/discord/blob_storage",
			".config/discord/sentry",
		},
		Whitelist: map[string][]string{
			"%HOME%/.config/discord/*/Discord": nil,
			// Versioned app dir; the wildcard covers every release.
			"%HOME%/.config/discord/*/chrome-sandbox":          nil,
			"%HOME%/.config/discord/*/chrome_crashpad_handler": nil,
		}},
	{Name: "Discord Canary", RelPaths: []string{".config/discord-canary"},
		Whitelist: map[string][]string{"/usr/bin/discord-canary": nil}},
	{Name: "Telegram", RelPaths: []string{".local/share/TelegramDesktop"},
		Whitelist: map[string][]string{"/usr/bin/telegram-desktop": nil, "/usr/bin/Telegram": nil}},
	{Name: "Signal", RelPaths: []string{".config/Signal"},
		Whitelist: map[string][]string{"/usr/bin/signal-desktop": nil}},
	{Name: "Signal Beta", RelPaths: []string{".config/Signal Beta"},
		Whitelist: map[string][]string{
			"/usr/bin/signal-desktop-beta":             nil,
			"/opt/Signal Beta/signal-desktop-beta":     nil,
			"/opt/Signal Beta/chrome_crashpad_handler": nil,
		}},
	{Name: "Element", RelPaths: []string{".config/Element"},
		Whitelist: map[string][]string{"/usr/bin/element-desktop": nil}},

	// --- Gaming ----------------------------------------------------------------------
	// Steam keeps login credentials (accounts, auth tokens) in
	// config/*.vdf under the data dir and in registry.vdf under the legacy
	// ~/.steam home — one resource, both locations guarded, one whitelist.
	// The rest of the Steam tree holds no secrets; legacy ssfn* sentries
	// stay unprotected (dirs-only guarding, low value alone).
	{Name: "Steam", RelPaths: []string{".local/share/Steam/config", ".local/share/Steam/userdata/*/config/localconfig.vdf", ".steam/registry.vdf"},
		Whitelist: map[string][]string{
			"/usr/bin/steam": nil, "/usr/bin/steamwebhelper": nil,
			"/usr/lib/steam/steam":                        nil,
			"%HOME%/.local/share/Steam/ubuntu12_32/steam": nil,
			// Webhelper moved across versioned dirs; the glob covers all layouts.
			"%HOME%/.local/share/Steam/*/steamwebhelper": nil,
			"/usr/bin/lsof":                       nil,
			"%HOME%/.steam/root/steamapps/common": nil, // protons installs allowed
		}},

	// --- System-level paths (probed once, not per user; ssh-guard template) ---
	{Name: "WireGuard system config", AbsPaths: []string{"/etc/wireguard"},
		Whitelist: map[string][]string{"/usr/bin/nmcli": nil}},
}

// PathsFor returns every absolute candidate path for user, expanding
// %USER%/%HOME%; AbsPaths entries ignore the user entirely. Order follows
// RelPaths / AbsPaths.
func (c *CandidateDir) PathsFor(home, user string) []string {
	if c.IsSystem() {
		out := make([]string, len(c.AbsPaths))
		for i, p := range c.AbsPaths {
			out[i] = expandPlaceholders(p, user, home)
		}
		return out
	}
	out := make([]string, len(c.RelPaths))
	for i, p := range c.RelPaths {
		out[i] = filepath.Join(home, expandPlaceholders(p, user, home))
	}
	return out
}

// ExtraWatchPathsFor expands the entry's WatchRelPaths into absolute paths
// for the given user (placeholders resolved, same as PathsFor).
func (c *CandidateDir) ExtraWatchPathsFor(home, user string) []string {
	if len(c.WatchRelPaths) == 0 {
		return nil
	}
	out := make([]string, 0, len(c.WatchRelPaths))
	for _, rel := range c.WatchRelPaths {
		out = append(out, filepath.Join(home, expandPlaceholders(rel, user, home)))
	}
	return out
}

// BinaryRule is one whitelisted binary and its allowed events (empty list
// means every event).
type BinaryRule struct {
	Path   string
	Events []string
}

// ExpandWhitelist expands %USER%/%HOME% and returns rules sorted by path so
// config generation is deterministic despite map ordering.
func (c *CandidateDir) ExpandWhitelist(user, home string) []BinaryRule {
	paths := make([]string, 0, len(c.Whitelist))
	for bin := range c.Whitelist {
		paths = append(paths, bin)
	}
	sort.Strings(paths)
	out := make([]BinaryRule, 0, len(paths))
	for _, bin := range paths {
		out = append(out, BinaryRule{
			Path:   expandPlaceholders(bin, user, home),
			Events: c.Whitelist[bin],
		})
	}
	return out
}

func expandPlaceholders(s, user, home string) string {
	s = strings.ReplaceAll(s, "%USER%", user)
	return strings.ReplaceAll(s, "%HOME%", home)
}

// Candidate is a discovered, existing path belonging to a specific user.
type Candidate struct {
	User  User
	Entry CandidateDir
	Path  string
}

// Discover probes the catalog for a user, returning existing per-user paths.
// AbsPaths entries are skipped (use DiscoverSystem/DiscoverForUsers). An
// entry with several RelPaths yields one Candidate per existing location.
func Discover(user User) []Candidate {
	home := user.Home
	var out []Candidate
	for i := range Catalog {
		entry := &Catalog[i]
		if entry.IsSystem() {
			continue
		}
		for _, path := range entry.PathsFor(home, user.Name) {
			if _, err := os.Lstat(path); err == nil {
				out = append(out, Candidate{User: user, Entry: *entry, Path: path})
			}
		}
	}
	return out
}

// DiscoverSystem probes the system-level (AbsPaths) entries, user-independent.
func DiscoverSystem() []Candidate {
	var out []Candidate
	for i := range Catalog {
		entry := &Catalog[i]
		if !entry.IsSystem() {
			continue
		}
		for _, path := range entry.PathsFor("", "") {
			if _, err := os.Lstat(path); err == nil {
				out = append(out, Candidate{User: User{}, Entry: *entry, Path: path})
			}
		}
	}
	return out
}

// DiscoverForUsers probes every selected user's catalog paths plus the
// system-level entries, deduplicated by path.
func DiscoverForUsers(users []User) []Candidate {
	var out []Candidate
	seen := make(map[string]bool)
	system := DiscoverSystem()
	for i := range system {
		c := &system[i]
		seen[c.Path] = true
		out = append(out, *c)
	}
	for _, u := range users {
		found := Discover(u)
		for i := range found {
			c := &found[i]
			if seen[c.Path] {
				continue
			}
			seen[c.Path] = true
			out = append(out, *c)
		}
	}
	return out
}

// FilterExistingWhitelist drops missing entries and expands globs (*, ?, [)
// to existing matches, re-evaluated every install so relocated apps are
// picked up; matches inherit the pattern's events.
func (c *Candidate) FilterExistingWhitelist() []BinaryRule {
	var out []BinaryRule
	for _, rule := range c.Entry.ExpandWhitelist(c.User.Name, c.User.Home) {
		if strings.ContainsAny(rule.Path, "*?[") {
			matches, err := filepath.Glob(rule.Path)
			if err != nil {
				continue // malformed pattern: skip the whole entry
			}
			for _, m := range matches {
				if _, err := os.Stat(m); err == nil {
					out = append(out, BinaryRule{Path: m, Events: rule.Events})
				}
			}
			continue
		}
		if _, err := os.Stat(rule.Path); err == nil {
			out = append(out, rule)
		}
	}
	return out
}
