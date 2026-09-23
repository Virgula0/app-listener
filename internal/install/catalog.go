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
	// Name is the TUI label and identifies the resource: two locations of one app (Steam's
	// ".local/share/Steam/config" and ".steam") share one entry and Name.
	Name string
	// RelPaths are home-relative locations. More than one may exist on a host; each existing one is
	// its own watch + fscrypt root, all sharing this entry's Whitelist. Set exactly one of
	// RelPaths/AbsPaths.
	RelPaths []string
	// AbsPaths are system-level locations (e.g. "/etc/wireguard"), probed once regardless of users.
	// Same multi-location semantics as RelPaths.
	AbsPaths []string
	// Whitelist maps binary paths to allowed events (nil = all, emitted bare; else "<path>
	// EV1,EV2"); matched by inode. Keep it minimal: each entry is a potential privilege-escalation
	// path while the directory is unlocked. Shared (union) across an entry's RelPaths.
	Whitelist map[string][]string
	// WatchRelPaths replaces the single [watch <RelPath>] with MULTIPLE guarded subtrees inside
	// that RelPath, sharing the whitelist and fscrypt lifecycle (the section path stays the
	// encryption root). For self-updating apps whose update workspace must stay outside the guarded
	// set. Empty = plain per-RelPath watch. Only valid with exactly one RelPaths entry;
	// home-relative, same %HOME% expansion as whitelist entries.
	WatchRelPaths []string
	// Libs are extra libraries the entry's whitelisted binaries may load (`allow_lib`). Only needed
	// for user-writable libraries at FIXED paths an app dlopen()s: root-owned system libs are
	// auto-trusted and in-tree ones trusted by location. %HOME%/%USER% expand. Per-launch ephemeral
	// paths (e.g. Steam's runtime) can't be listed; guard their containing directory instead.
	Libs []string
	// LibDirRelPaths are home-relative library DIRECTORIES guarded read-only (`lib_dir`): everyone
	// may read, only this entry's whitelisted binaries may write. A directory, not a Libs file
	// list, is the only workable rule for trees assembled per launch under changing names (Steam's
	// pressure-vessel `var/tmp-XXXXXX`): the trust object walks a loaded library's ancestors, so
	// one entry covers the subtree forever. Globs (*, ?, [) expand against the filesystem; missing
	// paths are dropped. Use the SHALLOWEST directory holding only libraries and app data the
	// binary owns.
	LibDirRelPaths []string
	// LibDirWriters are binaries allowed to WRITE this entry's LibDirRelPaths and nothing else
	// (`lib_binary`; unlike Whitelist they never reach the protected secret paths). For an app's
	// runtime-maintenance tools: Steam's pressure-vessel helpers rebuild a merged /usr under
	// `var/tmp-XXXXXX` each launch and symlink host GPU drivers into the runtime, but must not read
	// Steam login credentials in the same entry. Same %HOME%/%USER%/glob expansion as Whitelist;
	// only regular files survive.
	LibDirWriters []string
	// ReservedLibs are %HOME% library patterns "<dir>/<name>" (name exact, prefix* or *suffix) for
	// apps that install their own code outside every guarded tree. The name is reserved at any depth
	// below dir for the entry's whitelisted binaries (guard_trust.bpf.c #3), which may then load such
	// files with no allow_lib. dir must exist when the daemon (re)loads, or nothing is reserved.
	ReservedLibs []string
}

// IsSystem reports a system-level (AbsPaths) entry, probed once regardless of users.
func (c *CandidateDir) IsSystem() bool { return len(c.AbsPaths) > 0 }

// Catalog is the master list of critical directories probed per selected user; missing paths are
// simply not proposed.
var Catalog = []CandidateDir{
	// --- SSH and remote access -------------------------------------------------
	{Name: "SSH client configuration and keys", RelPaths: []string{".ssh"},
		// ssh needs WRITE (known_hosts rotation); sshd and helpers (sshd-auth/-session,
		// sftp-server) only READ authorized_keys; keysign/pkcs11/cleanup helpers are inert. git
		// excluded: it reads keys.
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
	// One resource, two locations (~/.claude and ~/.config/claude): same whitelist, each its own
	// watch root.
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
	// IDE configs hold auth tokens; wrapper binaries (/usr/bin/code, firefox, ...) are shell
	// scripts, so the whitelist matches the executed ELF under /opt and /usr/lib.
	{Name: "VS Code config", RelPaths: []string{".config/Code"},
		Whitelist: map[string][]string{
			"/usr/bin/code": nil, "/usr/local/bin/code": nil,
			"%HOME%/.local/bin/code":       nil,
			"/opt/visual-studio-code/code": nil, "/usr/share/code/code": nil, "/opt/visual-studio-code/bin/code": nil,
			// Code's separate crash-dump writer process.
			"/opt/visual-studio-code/chrome_crashpad_handler": nil, "/usr/share/code/chrome_crashpad_handler": nil,
			// Code's extension-signature verifier (checks CachedExtensionVSIXs/*.sigzip on
			// install/update).
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
		// IDEs run as the bundled JBR java, so the JBR and its fsnotifier are whitelisted (globs
		// cover /opt installs). Toolbox apps stay unguarded (no credentials; those live in
		// .config/JetBrains).
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
	// Config (~/.config/zed) and data (~/.local/share/zed, holds the auth token): one resource, two
	// watch roots.
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
		// Python launchers: the exe is the interpreter (not whitelisted), so both entries are
		// inert; documentation only.
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
	// No single-file entries: the daemon guards directories only (non-dir paths are skipped at
	// load), so e.g. ~/.netrc would be silently dropped (same for Steam's ssfn files).

	// --- Wallets and crypto -------------------------------------------------------
	// Surface only once the wallet is installed. Solana keeps its keypair as a plaintext id.json.
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
	// Browser binaries are shell wrappers, so the whitelist matches the executed ELF under /usr/lib
	// and /opt. Firefox's crashhelper and the chrome_crashpad_handler family write Crash Reports
	// separately.
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
		// WatchRelPaths encrypts/decrypts the whole ~/.config/discord with fscrypt but only watches
		// and guards the listed subpaths (other accesses are ignored, and the guarded set is
		// faster). Bonus: on shutdown the entire directory is re-encrypted, not just the watched
		// groups. Not feasible for Steam (too big, may contain games).
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
		},
		// Self-updated native modules (discord_voice.node, ...) and bundled libs (libffmpeg.so).
		ReservedLibs: []string{
			"%HOME%/.config/discord/*.so",
			"%HOME%/.config/discord/lib*",
			"%HOME%/.config/discord/*.node",
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
	// Steam keeps login credentials in config/*.vdf under the data dir and registry.vdf under the
	// legacy ~/.steam: one resource, both locations guarded, one whitelist. The rest of the tree
	// holds no secrets; legacy ssfn* sentries stay unprotected (dirs-only guarding, low value
	// alone).
	{Name: "Steam", RelPaths: []string{".local/share/Steam/config", ".local/share/Steam/userdata/*/config/localconfig.vdf", ".steam/registry.vdf"},
		Whitelist: map[string][]string{
			"/usr/bin/steam":                                                                     nil,
			"/usr/bin/steamwebhelper":                                                            nil,
			"/usr/lib/steam/steam":                                                               nil,
			"%HOME%/.local/share/Steam/*/steam":                                                  nil,
			"%HOME%/.local/share/Steam/*/steamwebhelper":                                         nil,
			"%HOME%/.local/share/Steam/*/gameoverlayui":                                          nil,
			"%HOME%/.local/share/Steam/steamapps/common/*/files/bin/wineserver":                  nil,
			"%HOME%/.local/share/Steam/steamapps/common/*/files/lib/wine/*/wine64-preloader":     nil,
			"%HOME%/.local/share/Steam/compatibilitytools.d/*/files/lib/wine/*/wine64-preloader": nil,
			// Wine's WoW64 mode (PROTON_USE_WOW64=1, default in newer Proton) runs every process
			// under wine-preloader instead of wine64-preloader: same role, same access.
			"%HOME%/.local/share/Steam/steamapps/common/*/files/lib/wine/*/wine-preloader":     nil,
			"%HOME%/.local/share/Steam/compatibilitytools.d/*/files/lib/wine/*/wine-preloader": nil,
			// Custom Proton builds (GE-Proton) ship their own wineserver, counterpart of the
			// steamapps/common one above. Without it a GE game (tainted because its
			// wine64-preloader is whitelisted) can't be reached by its own server (ptrace ATTACH
			// denied).
			"%HOME%/.local/share/Steam/compatibilitytools.d/*/files/bin/wineserver": nil,
			"/usr/bin/lsof": nil,
			"/usr/bin/ps":   nil,
		},
		// Steam loads hundreds of libraries that aren't root-owned (so none is auto-trusted): its
		// shipped runtime (ubuntu12_*, linux64), Proton/Wine builds under compatibilitytools.d, and
		// the pressure-vessel container runtimes. Guarded read-only (all may read, only Steam may
		// write), which makes them safe to load and blocks planting an LD_PRELOAD payload.
		//
		// The runtime dirs matter most: pressure-vessel assembles a merged /usr under a random
		// `var/tmp-XXXXXX` each launch, so no file list could cover it. Guarding the STABLE parent
		// does (the trust object walks a loaded library's ancestors). steamapps/common/* at large
		// is deliberately not listed: it's the game library (hundreds of GiB of non-library data).
		LibDirRelPaths: []string{
			".local/share/Steam/ubuntu12_32",
			".local/share/Steam/ubuntu12_64",
			".local/share/Steam/linux32",
			".local/share/Steam/linux64",
			".local/share/Steam/steamrt64",
			".local/share/Steam/compatibilitytools.d",
			// Container runtimes only (soldier, sniper, 4; hence the underscore): their whitelisted
			// pressure-vessel binaries load libraries from inside them. NOT the legacy scout
			// runtime ("SteamLinuxRuntime", no suffix): it holds only libraries for non-whitelisted
			// native games, so guarding it protected nothing, while its entry point must `ln -fns
			// var/steam-runtime/amd64` each launch and the refusal made native scout games exit
			// instantly.
			".local/share/Steam/steamapps/common/SteamLinuxRuntime_*",
			// Valve's Proton builds (steamapps/common/Proton 11.0, Proton - Experimental, ...):
			// their whitelisted wine64-preloader loads Wine from files/lib. ONLY files/lib: the
			// `proton` script (python3, which can't be whitelisted) rewrites dist.lock and builds
			// files/share/default_pfx each launch, so guarding the whole folder would refuse every
			// Valve-Proton game. GE-Proton needs no entry (compatibilitytools.d is guarded above).
			".local/share/Steam/steamapps/common/Proton*/files/lib",
		},
		// pressure-vessel OWNS the runtime trees above and rewrites them every launch: it hardlinks
		// the runtime's ~6.6k files into a fresh `var/tmp-XXXXXX`, capsule-capture-libs symlinks
		// host GPU drivers (libvdpau_nvidia.so, Vulkan layers, libc) into `overrides/lib/*`, and
		// pv-locale-gen builds a locale archive, all INSIDE the read-only tree. Without these
		// writers the container has no GL/Vulkan and games don't start.
		//
		// Whole helper directories are listed, not individual tools (the set grows with every
		// runtime update); they're the tree's own vendor binaries, write-protected by this guard.
		// They are lib_binary (not Whitelist) entries so they can't reach the Steam credentials in
		// this entry. Only WRITING tools are listed: pressure-vessel-wrap/-unruntime (build the tmp
		// tree, take `.ref` locks), capsule-capture-libs, and the pv-*/srt-* helpers
		// (pv-locale-gen, pv-adverb, srt-logger, srt-bwrap). Read-only probes (check-gl,
		// inspect-library, detect-platform, wflinfo, true) need nothing.
		//
		// Steam's own binaries write these trees too: the client self-updates ubuntu12_32/64,
		// linux64, steamrt64; steamwebhelper keeps CEF state beside itself (ubuntu12_64); the
		// overlay and Proton's wine live and write inside them. A [libraries] lib_dir inherits no
		// whitelist, so these are the unrestricted Whitelist entries the old nested form let write,
		// minus lsof (never writes) and shell-script launchers (they run as their interpreter's
		// inode).
		LibDirWriters: []string{
			"%HOME%/.local/share/Steam/*/steam",
			"%HOME%/.local/share/Steam/*/steamwebhelper",
			"%HOME%/.local/share/Steam/*/gameoverlayui",
			"%HOME%/.local/share/Steam/steamapps/common/*/files/bin/wineserver",
			"%HOME%/.local/share/Steam/steamapps/common/*/files/lib/wine/*/wine64-preloader",
			"%HOME%/.local/share/Steam/compatibilitytools.d/*/files/lib/wine/*/wine64-preloader",
			"%HOME%/.local/share/Steam/steamapps/common/*/files/lib/wine/*/wine-preloader",
			"%HOME%/.local/share/Steam/compatibilitytools.d/*/files/lib/wine/*/wine-preloader",
			"%HOME%/.local/share/Steam/compatibilitytools.d/*/files/bin/wineserver",
			"%HOME%/.local/share/Steam/steamrt64/pv-runtime/*/pressure-vessel/bin/pressure-vessel-*",
			"%HOME%/.local/share/Steam/steamrt64/pv-runtime/*/pressure-vessel/libexec/steam-runtime-tools-0/*-capsule-capture-libs",
			"%HOME%/.local/share/Steam/steamrt64/pv-runtime/*/pressure-vessel/libexec/steam-runtime-tools-0/pv-*",
			"%HOME%/.local/share/Steam/steamrt64/pv-runtime/*/pressure-vessel/libexec/steam-runtime-tools-0/srt-*",
			"%HOME%/.local/share/Steam/steamapps/common/SteamLinuxRuntime_*/pressure-vessel*/bin/pressure-vessel-*",
			"%HOME%/.local/share/Steam/steamapps/common/SteamLinuxRuntime_*/pressure-vessel*/libexec/steam-runtime-tools-0/*-capsule-capture-libs",
			"%HOME%/.local/share/Steam/steamapps/common/SteamLinuxRuntime_*/pressure-vessel*/libexec/steam-runtime-tools-0/pv-*",
			"%HOME%/.local/share/Steam/steamapps/common/SteamLinuxRuntime_*/pressure-vessel*/libexec/steam-runtime-tools-0/srt-*",
		}},

	// --- System-level paths (probed once, not per user; ssh-guard template) ---
	{Name: "WireGuard system config", AbsPaths: []string{"/etc/wireguard"},
		Whitelist: map[string][]string{"/usr/bin/nmcli": nil}},
}

// PathsFor returns every absolute candidate path for user (expanding %USER%/%HOME%; AbsPaths ignore
// user), in RelPaths/AbsPaths order.
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

// ExtraWatchPathsFor expands WatchRelPaths to absolute paths for user, keeping only existing ones
// (same gate as Discover for RelPaths). A missing sub-path must never reach the generated config:
// the daemon treats every grouped `watch:` as a tree that must resolve after unlock and fails
// closed (fatal) otherwise, so it would crash the daemon rather than under-protect.
func (c *CandidateDir) ExtraWatchPathsFor(home, user string) []string {
	if len(c.WatchRelPaths) == 0 {
		return nil
	}
	out := make([]string, 0, len(c.WatchRelPaths))
	for _, rel := range c.WatchRelPaths {
		p := filepath.Join(home, expandPlaceholders(rel, user, home))
		if _, err := os.Lstat(p); err != nil {
			continue
		}
		out = append(out, p)
	}
	return out
}

// BinaryRule is one whitelisted binary and its allowed events (empty = every event).
type BinaryRule struct {
	Path   string
	Events []string
}

// ExpandWhitelist expands %USER%/%HOME% and returns rules sorted by path (deterministic config
// generation).
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

// ExpandLibs expands %USER%/%HOME% in allow_lib paths, sorted.
func (c *CandidateDir) ExpandLibs(user, home string) []string {
	out := make([]string, 0, len(c.Libs))
	for _, lib := range c.Libs {
		out = append(out, expandPlaceholders(lib, user, home))
	}
	sort.Strings(out)
	return out
}

// ExpandLibDirs expands LibDirRelPaths for user into absolute directories: placeholders resolved,
// globs matched, non-directories and missing paths dropped (a nonexistent tree protects nothing,
// and versioned runtime dirs come and go), sorted and de-duplicated.
func (c *CandidateDir) ExpandLibDirs(user, home string) []string {
	if len(c.LibDirRelPaths) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	out := make([]string, 0, len(c.LibDirRelPaths))
	for _, rel := range c.LibDirRelPaths {
		pattern := filepath.Join(home, expandPlaceholders(rel, user, home))
		matches := []string{pattern}
		if strings.ContainsAny(pattern, "*?[") {
			m, err := filepath.Glob(pattern)
			if err != nil {
				continue
			}
			matches = m
		}
		for _, p := range matches {
			if seen[p] {
				continue
			}
			info, err := os.Lstat(p)
			if err != nil || !info.IsDir() {
				continue
			}
			seen[p] = true
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out
}

// ExpandLibDirWriters expands %USER%/%HOME% in LibDirWriters, matches globs, and keeps only
// existing REGULAR files, sorted and de-duplicated. Globbing a whole helper directory is the
// intended use (tool dirs gain binaries each update), so directories and dangling symlinks are
// dropped rather than written as writers that can never match.
func (c *CandidateDir) ExpandLibDirWriters(user, home string) []string {
	if len(c.LibDirWriters) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	out := make([]string, 0, len(c.LibDirWriters))
	for _, bin := range c.LibDirWriters {
		pattern := expandPlaceholders(bin, user, home)
		matches := []string{pattern}
		if strings.ContainsAny(pattern, "*?[") {
			m, err := filepath.Glob(pattern)
			if err != nil {
				continue // malformed pattern: skip the whole entry
			}
			matches = m
		}
		for _, p := range matches {
			if seen[p] {
				continue
			}
			// Stat, not Lstat: a symlinked helper resolves like any whitelisted binary (the daemon
			// records the target's inode).
			info, err := os.Stat(p)
			if err != nil || !info.Mode().IsRegular() {
				continue
			}
			seen[p] = true
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out
}

// LibraryBlockFor builds this entry's `[libraries "<name> (<user>)"]` block for one user. The user
// is in the name because library trees are per-user (own runtime, own writers) and the installer
// finds the block by name on refresh.
func (c *CandidateDir) LibraryBlockFor(user, home string) LibraryBlock {
	name := c.Name
	if user != "" {
		name += " (" + user + ")"
	}
	return LibraryBlock{
		Name:       name,
		Libs:       c.ExpandLibs(user, home),
		LibDirs:    c.ExpandLibDirs(user, home),
		LibWriters: c.ExpandLibDirWriters(user, home),
	}
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

// Discover probes the catalog for a user, returning existing per-user paths (AbsPaths skipped: see
// DiscoverSystem/DiscoverForUsers). An entry with several RelPaths yields one Candidate per
// existing location.
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

// DiscoverForUsers probes every selected user's paths plus system-level entries, deduplicated by
// path.
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

// FilterExistingWhitelist drops missing entries and expands globs (*, ?, [) to existing matches,
// re-evaluated every install so relocated apps are found; matches inherit the pattern's events.
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
