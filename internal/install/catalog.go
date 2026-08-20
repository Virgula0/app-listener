// Package install implements the building blocks of the `app-listener
// install` wizard: the catalog of critical directories, user enumeration,
// the recursive tree copy used during fscrypt migrations, and the daemon
// config generation.
package install

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CandidateDir describes one critical directory type the installer probes.
type CandidateDir struct {
	// Name is a short human-readable label shown in the TUI.
	Name string
	// RelPath is the path relative to the user's home directory, e.g.
	// ".ssh" or ".config/opencode". Exactly one of RelPath and AbsPath
	// must be set.
	RelPath string
	// AbsPath is a system-level path (e.g. "/etc/wireguard") that is
	// probed once, independent of the selected users.
	AbsPath string
	// Whitelist maps each whitelisted binary path to the events that
	// binary may trigger on this directory. An empty (or nil) event list
	// means every event is allowed and the config line is emitted as a
	// bare path; a non-empty list restricts the binary to exactly those
	// events and is emitted as "<path> EV1,EV2". Keep it minimal: every
	// binary added here is a potential privilege-escalation path if it can
	// be abused while the directory is unlocked. Identity is verified by
	// inode, not by name. Whitelist entries that do not exist on the
	// target system are dropped when the config is generated; entries
	// containing a glob (*, ?, [) are expanded to every existing match.
	Whitelist map[string][]string
}

// Catalog is the master list of critical directories to search for each
// selected user. Add, remove or tweak entries here — this is the single
// place that drives discovery. Paths that do not exist for a user are
// simply not proposed.
var Catalog = []CandidateDir{
	// --- SSH and remote access -------------------------------------------------
	{Name: "SSH client configuration and keys", RelPath: ".ssh",
		// ssh rewrites known_hosts (WRITE plus the DELETE/RENAME/HARDLINK
		// its rotation needs); sshd and its modular helpers (sshd-auth,
		// sshd-session, sftp-server) only READ authorized_keys — Debian
		// installs them under /usr/lib/openssh/, Arch under /usr/lib/ssh/.
		// ssh-keysign/ssh-pkcs11-helper/ssh-session-cleanup never touch
		// .ssh (inert). git is deliberately NOT whitelisted: it could read
		// the keys.
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
	{Name: "GNU Privacy Guard keyring", RelPath: ".gnupg",
		Whitelist: map[string][]string{
			"/usr/bin/gpg": nil, "/usr/bin/gpg-agent": nil,
			"/usr/bin/gpgconf": nil, "/usr/bin/gpg-connect-agent": nil,
		}},

	// --- AI coding agents and CLI tools ----------------------------------------
	{Name: "opencode", RelPath: ".config/opencode",
		Whitelist: map[string][]string{
			"/usr/local/bin/opencode": nil, "/usr/bin/opencode": nil,
			"%HOME%/.local/bin/opencode": nil,
		}},
	{Name: "code CLI (GitHub)", RelPath: ".config/code-cli",
		Whitelist: map[string][]string{
			"/usr/bin/code-cli": nil, "/usr/local/bin/code-cli": nil,
			"%HOME%/.local/bin/code-cli": nil,
		}},
	{Name: "Claude Code", RelPath: ".claude",
		Whitelist: map[string][]string{
			"/usr/local/bin/claude": nil, "/usr/bin/claude": nil,
			"%HOME%/.local/bin/claude": nil,
		}},
	{Name: "Claude Code config", RelPath: ".config/claude",
		Whitelist: map[string][]string{
			"/usr/local/bin/claude": nil, "/usr/bin/claude": nil,
			"%HOME%/.local/bin/claude": nil,
		}},
	{Name: "Codeium", RelPath: ".codeium",
		Whitelist: map[string][]string{
			"/usr/local/bin/codeium": nil, "/usr/bin/codeium": nil,
			"%HOME%/.local/bin/codeium": nil,
		}},
	{Name: "Gemini CLI", RelPath: ".gemini",
		Whitelist: map[string][]string{
			"/usr/local/bin/gemini": nil, "/usr/bin/gemini": nil,
			"%HOME%/.local/bin/gemini": nil,
		}},
	{Name: "Cursor", RelPath: ".cursor",
		Whitelist: map[string][]string{
			"/usr/bin/cursor": nil, "/usr/local/bin/cursor": nil,
			"%HOME%/.local/bin/cursor": nil,
		}},
	{Name: "GitHub Copilot", RelPath: ".config/github-copilot",
		Whitelist: map[string][]string{
			"/usr/bin/git": nil, "/usr/local/bin/copilot": nil,
			"%HOME%/.local/bin/copilot": nil,
		}},

	// --- IDEs and editors ---------------------------------------------------------
	// IDE configs hold auth tokens, API keys and credentials (VS Code
	// globalStorage, JetBrains certs, Zed auth.json). Wrapper binaries
	// (/usr/bin/code, firefox, chrome, ...) are shell scripts; the
	// whitelist matches the executed ELF under /opt and /usr/lib.
	{Name: "VS Code config", RelPath: ".config/Code",
		Whitelist: map[string][]string{
			"/usr/bin/code": nil, "/usr/local/bin/code": nil,
			"%HOME%/.local/bin/code":       nil,
			"/opt/visual-studio-code/code": nil, "/usr/share/code/code": nil, "/opt/visual-studio-code/bin/code": nil,
			// chrome_crashpad_handler is a separate process Code spawns to
			// write crash dumps into .config/Code/Crashpad.
			"/opt/visual-studio-code/chrome_crashpad_handler": nil, "/usr/share/code/chrome_crashpad_handler": nil,
		}},
	{Name: "VS Code Insiders config", RelPath: ".config/Code - Insiders",
		Whitelist: map[string][]string{
			"/usr/bin/code-insiders": nil, "/usr/local/bin/code-insiders": nil,
			"%HOME%/.local/bin/code-insiders": nil, "/opt/visual-studio-code/bin/code": nil,
		}},
	{Name: "VS Code server (remote development)", RelPath: ".vscode-server",
		Whitelist: map[string][]string{
			"/usr/bin/code": nil, "/usr/local/bin/code": nil,
			"%HOME%/.local/bin/code":       nil,
			"/opt/visual-studio-code/code": nil, "/usr/share/code/code": nil,
			// chrome_crashpad_handler is a separate process Code spawns to
			// write crash dumps into .config/Code/Crashpad.
			"/opt/visual-studio-code/chrome_crashpad_handler": nil, "/usr/share/code/chrome_crashpad_handler": nil,
		}},
	{Name: "VSCodium config", RelPath: ".config/VSCodium",
		Whitelist: map[string][]string{"/usr/bin/codium": nil, "/usr/local/bin/codium": nil,
			"/opt/vscodium/chrome_crashpad_handler": nil, "/usr/share/vscodium/chrome_crashpad_handler": nil}},
	{Name: "JetBrains IDEs", RelPath: ".config/JetBrains",
		// IDE processes run as the bundled JetBrains Runtime: the process
		// exe is the JBR java, so the JBR and its fsnotifier helper must be
		// whitelisted. The globs cover standalone installs under /opt. IDE
		// installs via Toolbox (~/.local/share/JetBrains/Toolbox/apps) are
		// deliberately not guarded: they hold no credentials — the IDE config
		// and auth state live in this .config/JetBrains directory.
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
	{Name: "Zed editor", RelPath: ".config/zed",
		Whitelist: map[string][]string{
			"/usr/bin/zed": nil, "/usr/local/bin/zed": nil,
			"%HOME%/.local/bin/zed": nil,
		}},
	{Name: "Zed editor data", RelPath: ".local/share/zed",
		Whitelist: map[string][]string{
			"/usr/bin/zed": nil, "/usr/local/bin/zed": nil,
			"%HOME%/.local/bin/zed": nil,
		}},
	{Name: "Sublime Text", RelPath: ".config/sublime-text",
		// /usr/bin/subl is a wrapper; the real binary is under /opt.
		Whitelist: map[string][]string{"/usr/bin/subl": nil, "/usr/bin/sublime-text": nil, "/usr/bin/sublime_text": nil,
			"/opt/sublime_text/sublime_text": nil}},
	{Name: "Insomnia API client", RelPath: ".config/Insomnia",
		Whitelist: map[string][]string{"/usr/bin/insomnia": nil, "/opt/insomnia/insomnia": nil}},
	{Name: "Postman API client", RelPath: ".config/Postman",
		Whitelist: map[string][]string{"/usr/bin/postman": nil, "/opt/postman/postman": nil}},
	{Name: "DBeaver database client", RelPath: ".local/share/DBeaverData",
		// /usr/bin/dbeaver is a wrapper; the real binary lives in /usr/lib.
		Whitelist: map[string][]string{"/usr/bin/dbeaver": nil, "/usr/lib/dbeaver/dbeaver": nil}},
	{Name: "Android adb keys", RelPath: ".android",
		Whitelist: map[string][]string{"/usr/bin/adb": nil}},

	// --- Cloud, containers and dev tooling --------------------------------------
	{Name: "AWS credentials", RelPath: ".aws",
		Whitelist: map[string][]string{"/usr/bin/aws": nil}},
	{Name: "Google Cloud SDK", RelPath: ".config/gcloud",
		// gcloud/gsutil are Python launchers: the process exe is the
		// interpreter (deliberately not whitelisted), so the entries are
		// inert — kept for documentation.
		Whitelist: map[string][]string{"/usr/bin/gcloud": nil, "/usr/bin/gsutil": nil}},
	{Name: "Kubernetes kubeconfig", RelPath: ".kube",
		Whitelist: map[string][]string{"/usr/bin/kubectl": nil, "/usr/bin/helm": nil, "/usr/bin/oc": nil,
			"/usr/lib/docker/cli-plugins/docker-buildx":     nil,
			"/usr/libexec/docker/cli-plugins/docker-buildx": nil,
			"%HOME%/.docker/cli-plugins/docker-buildx":      nil}},
	{Name: "Docker config and credentials", RelPath: ".docker",
		Whitelist: map[string][]string{"/usr/bin/docker": nil, "/usr/bin/docker-credential-desktop": nil, "/usr/bin/docker-credential-pass": nil,
			"/usr/lib/docker/cli-plugins/docker-buildx":     nil,
			"/usr/libexec/docker/cli-plugins/docker-buildx": nil,
			"%HOME%/.docker/cli-plugins/docker-buildx":      nil}},
	{Name: "GitHub CLI", RelPath: ".config/gh",
		Whitelist: map[string][]string{
			"/usr/bin/gh":                       nil,
			"/usr/local/bin/gh":                 nil,
			"%HOME%/.local/bin/gh":              nil,
			"/home/linuxbrew/.linuxbrew/bin/gh": nil,
			"/snap/bin/gh":                      nil,
		}},
	{Name: "Azure CLI", RelPath: ".azure",
		Whitelist: map[string][]string{"/usr/bin/az": nil}},
	{Name: "Azure CLI config", RelPath: ".config/azure",
		Whitelist: map[string][]string{"/usr/bin/az": nil}},
	{Name: "rclone config", RelPath: ".config/rclone",
		Whitelist: map[string][]string{"/usr/bin/rclone": nil}},
	{Name: "Ollama config", RelPath: ".ollama",
		Whitelist: map[string][]string{"/usr/bin/ollama": nil}},
	{Name: "npm data", RelPath: ".npm",
		Whitelist: map[string][]string{"/usr/bin/npm": nil, "/usr/bin/npx": nil, "/usr/bin/node": nil}},
	// pip is a #!/usr/bin/python script (inert, same reason as gcloud);
	// kept for documentation.
	{Name: "pip config", RelPath: ".config/pip",
		Whitelist: map[string][]string{"/usr/bin/pip": nil, "/usr/bin/pip3": nil}},
	{Name: "Git config directory", RelPath: ".config/git",
		Whitelist: map[string][]string{"/usr/bin/git": nil}},

	// --- Password managers and secrets -------------------------------------------
	{Name: "password-store (pass)", RelPath: ".password-store",
		Whitelist: map[string][]string{"/usr/bin/pass": nil, "/usr/bin/gpg": nil}},
	{Name: "gopass", RelPath: ".local/share/gopass",
		Whitelist: map[string][]string{"/usr/bin/gopass": nil}},
	{Name: "Bitwarden", RelPath: ".config/Bitwarden",
		Whitelist: map[string][]string{"/usr/bin/bitwarden": nil, "/usr/local/bin/bitwarden": nil}},
	{Name: "1Password", RelPath: ".1password",
		Whitelist: map[string][]string{"/usr/bin/1password": nil, "/usr/bin/op": nil, "/usr/local/bin/op": nil}},
	{Name: "age keys", RelPath: ".config/age",
		Whitelist: map[string][]string{"/usr/bin/age": nil, "/usr/bin/age-keygen": nil}},
	{Name: "KeePassXC", RelPath: ".local/share/keepassxc",
		Whitelist: map[string][]string{"/usr/bin/keepassxc": nil}},
	{Name: "GNOME keyring", RelPath: ".local/share/keyrings",
		Whitelist: map[string][]string{"/usr/bin/gnome-keyring-daemon": nil, "/usr/bin/gnome-keyring": nil}},
	{Name: "GNOME keyring (legacy)", RelPath: ".keyring",
		Whitelist: map[string][]string{"/usr/bin/gnome-keyring-daemon": nil, "/usr/bin/gnome-keyring": nil}},
	{Name: "Chezmoi state", RelPath: ".local/share/chezmoi",
		Whitelist: map[string][]string{"/usr/bin/chezmoi": nil}},

	// --- Legacy dot-config files --------------------------------------------------
	// No single-file entries here: the daemon guards directories only (the
	// config loader skips non-directory paths), so a plain file like
	// ~/.netrc or ~/.wgetrc would be written to the config but silently
	// dropped at load. The Steam entry below documents the same restriction
	// for the legacy ssfn sentry files and the required glob semantics.

	// --- Wallets and crypto -------------------------------------------------------
	// Probed like any other entry: they surface only once the wallet is
	// installed. Solana keeps its keypair as a plaintext id.json.
	{Name: "Bitcoin Core", RelPath: ".bitcoin/wallets",
		Whitelist: map[string][]string{"/usr/bin/bitcoind": nil, "/usr/bin/bitcoin-qt": nil, "/usr/bin/bitcoin-cli": nil, "/usr/bin/bitcoin-tx": nil}},
	{Name: "Litecoin Core", RelPath: ".litecoin/wallets",
		Whitelist: map[string][]string{"/usr/bin/litecoind": nil, "/usr/bin/litecoin-qt": nil}},
	{Name: "Monero", RelPath: ".bitmonero/wallets",
		Whitelist: map[string][]string{"/usr/bin/monero-wallet-gui": nil, "/usr/bin/monero-wallet-cli": nil, "/usr/bin/monerod": nil}},
	{Name: "Electrum", RelPath: ".electrum",
		Whitelist: map[string][]string{"/usr/bin/electrum": nil, "/usr/local/bin/electrum": nil}},
	{Name: "Sparrow", RelPath: ".sparrow",
		Whitelist: map[string][]string{"/usr/bin/sparrow": nil}},
	{Name: "Wasabi Wallet", RelPath: ".walletwasabi",
		Whitelist: map[string][]string{"/usr/bin/wassabee": nil}},
	{Name: "Ledger Live", RelPath: ".config/Ledger Live",
		Whitelist: map[string][]string{"/usr/bin/ledger-live": nil, "/opt/Ledger Live/ledger-live": nil}},
	{Name: "Trezor Suite", RelPath: ".config/@trezor",
		Whitelist: map[string][]string{"/usr/bin/trezor-suite": nil, "/opt/Trezor Suite/trezor-suite": nil}},
	{Name: "Exodus", RelPath: ".config/Exodus",
		Whitelist: map[string][]string{"/usr/bin/exodus": nil, "/opt/Exodus/exodus": nil}},
	// id.json is a plaintext keypair — the highest-value target of the set.
	{Name: "Solana CLI", RelPath: ".config/solana",
		Whitelist: map[string][]string{"/usr/bin/solana": nil, "/usr/bin/solana-keygen": nil}},
	{Name: "Ethereum (geth)", RelPath: ".ethereum/keystore",
		Whitelist: map[string][]string{"/usr/bin/geth": nil}},
	{Name: "LND (Lightning Network Daemon)", RelPath: ".lnd",
		Whitelist: map[string][]string{"/usr/bin/lnd": nil, "/usr/bin/lncli": nil}},
	{Name: "Core Lightning", RelPath: ".lightning",
		Whitelist: map[string][]string{"/usr/bin/lightningd": nil, "/usr/bin/lightning-cli": nil}},
	// Foundry's binaries live inside the guarded dir next to its keystores,
	// so they get the Discord treatment via a glob.
	{Name: "Foundry", RelPath: ".foundry",
		Whitelist: map[string][]string{"%HOME%/.foundry/bin/*": nil}},
	{Name: "KDE Wallet", RelPath: ".local/share/kwalletd",
		Whitelist: map[string][]string{"/usr/bin/kwalletd5": nil, "/usr/bin/kwalletd6": nil, "/usr/bin/kwalletmanager5": nil, "/usr/bin/kwalletmanager6": nil}},

	// --- VPNs ---------------------------------------------------------------------
	{Name: "NordVPN config", RelPath: ".config/nordvpn",
		Whitelist: map[string][]string{"/usr/bin/nordvpn": nil}},
	{Name: "Mullvad VPN config", RelPath: ".config/mullvad",
		Whitelist: map[string][]string{"/usr/bin/mullvad": nil}},
	{Name: "ProtonVPN config", RelPath: ".config/protonvpn",
		Whitelist: map[string][]string{"/usr/bin/protonvpn": nil, "/usr/bin/protonvpn-cli": nil}},
	{Name: "OpenVPN config", RelPath: ".config/openvpn",
		Whitelist: map[string][]string{"/usr/bin/openvpn": nil}},
	{Name: "WireGuard config (per-user)", RelPath: ".config/wireguard",
		Whitelist: map[string][]string{"/usr/bin/wg": nil, "/usr/bin/wg-quick": nil}},

	// --- Browsers and messaging -----------------------------------------------------
	// Browser binaries are shell wrappers; the whitelist matches the
	// executed ELF under /usr/lib and /opt. Firefox's crash helper and the
	// Chrome-family chrome_crashpad_handler are separate processes that
	// write the Crash Reports directory.
	{Name: "Firefox profile", RelPath: ".mozilla/firefox",
		Whitelist: map[string][]string{
			"/usr/bin/firefox":             nil,
			"/usr/lib/firefox/firefox":     nil,
			"/usr/lib/firefox/crashhelper": nil, "/usr/lib/firefox/crashreporter": nil,
		}},
	{Name: "Google Chrome profile", RelPath: ".config/google-chrome",
		Whitelist: map[string][]string{
			"/usr/bin/google-chrome": nil, "/usr/bin/google-chrome-stable": nil,
			"/opt/google/chrome/chrome": nil, "/usr/lib/chromium/chrome": nil,
			"/opt/google/chrome/chrome_crashpad_handler": nil,
		}},
	{Name: "Chromium profile", RelPath: ".config/chromium",
		Whitelist: map[string][]string{
			"/usr/bin/chromium": nil, "/usr/bin/chromium-browser": nil,
			"/usr/lib/chromium/chrome": nil, "/usr/lib/chromium/chromium": nil,
			"/usr/lib/chromium/chrome_crashpad_handler": nil,
		}},
	{Name: "Brave profile", RelPath: ".config/BraveSoftware",
		Whitelist: map[string][]string{
			"/usr/bin/brave-browser": nil, "/usr/bin/brave": nil,
			"/opt/brave/brave": nil, "/opt/brave-bin/brave": nil,
			"/opt/brave-bin/chrome_crashpad_handler":         nil,
			"/usr/lib/brave-browser/chrome_crashpad_handler": nil,
		}},
	{Name: "Discord", RelPath: ".config/discord",
		Whitelist: map[string][]string{
			"%HOME%/.config/discord/*/Discord":                 nil,
			"%HOME%/.config/discord/*/chrome-sandbox":          nil,
			"%HOME%/.config/discord/*/chrome_crashpad_handler": nil,
		}},
	{Name: "Discord Canary", RelPath: ".config/discord-canary",
		Whitelist: map[string][]string{"/usr/bin/discord-canary": nil}},
	{Name: "Telegram", RelPath: ".local/share/TelegramDesktop",
		Whitelist: map[string][]string{"/usr/bin/telegram-desktop": nil, "/usr/bin/Telegram": nil}},
	{Name: "Signal", RelPath: ".config/Signal",
		Whitelist: map[string][]string{"/usr/bin/signal-desktop": nil}},
	{Name: "Signal Beta", RelPath: ".config/Signal Beta",
		Whitelist: map[string][]string{
			"/usr/bin/signal-desktop-beta":             nil,
			"/opt/Signal Beta/signal-desktop-beta":     nil,
			"/opt/Signal Beta/chrome_crashpad_handler": nil,
		}},
	{Name: "Element", RelPath: ".config/Element",
		Whitelist: map[string][]string{"/usr/bin/element-desktop": nil}},

	// --- Gaming ----------------------------------------------------------------------
	// Steam stores its login credentials (accounts and auth tokens) in
	// config/loginusers.vdf and config/config.vdf, so only the small config
	// directory is guarded: the rest of ~/.local/share/Steam (steamapps/
	// games, userdata/ screenshots and cloud saves, caches) holds no
	// credentials and is intentionally left out.
	//
	// The legacy Steam Guard sentry files (~/.local/share/Steam/ssfn*) are
	// deliberately NOT protected: the daemon guards directories only (the
	// config loader skips non-directory paths), and a sentry alone cannot
	// compromise an account. Once the daemon supports single files, add a
	// catalog entry for them — and the glob (ssfn*) must then be resolved
	// dynamically by the daemon at runtime (on start/reload), never baked
	// into the static config generated at install time, because sentry
	// filenames can change after login.
	{Name: "Steam client config (accounts/credentials)", RelPath: ".local/share/Steam/config",
		Whitelist: map[string][]string{
			"/usr/bin/steam": nil, "/usr/bin/steamwebhelper": nil,
			"%HOME%/.local/share/Steam/ubuntu12_32/steam": nil,
			// steamwebhelper moved from ubuntu12_32/ to ubuntu12_64/ and
			// steamrt64/ across versions; the glob covers all layouts.
			"%HOME%/.local/share/Steam/*/steamwebhelper": nil,
		}},
	{Name: "Steam (legacy home)", RelPath: ".steam",
		Whitelist: map[string][]string{"/usr/bin/steam": nil, "/usr/bin/steamwebhelper": nil}},

	// --- System-level paths (probed once, not per user) --------------------------------
	// Same whitelist as the original ssh-guard template:
	{Name: "WireGuard system config", AbsPath: "/etc/wireguard",
		Whitelist: map[string][]string{"/usr/bin/nmcli": nil}},
}

// PathFor returns the absolute candidate path for user, expanding the
// %USER%/%HOME% placeholders used by whitelist entries. System-level
// entries (AbsPath) ignore the user; per-user entries are joined with the
// user's home directory.
func (c CandidateDir) PathFor(home, user string) string {
	if c.AbsPath != "" {
		return expandPlaceholders(c.AbsPath, user, home)
	}
	return filepath.Join(home, expandPlaceholders(c.RelPath, user, home))
}

// BinaryRule is one whitelisted binary and the events it may perform on
// the protected path. An empty Events list means every event is allowed;
// otherwise only those events are permitted.
type BinaryRule struct {
	Path   string
	Events []string
}

// ExpandWhitelist returns the whitelist with %USER%/%HOME% placeholders
// replaced. The result is sorted by path so config generation is
// deterministic regardless of map iteration order, which Go does not
// guarantee.
func (c CandidateDir) ExpandWhitelist(user, home string) []BinaryRule {
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

// Discover probes the catalog for a user and returns every per-user path
// that actually exists (files or directories). System-level entries
// (AbsPath) are NOT probed here — use DiscoverSystem/DiscoverForUsers.
func Discover(user User) []Candidate {
	home := user.Home
	var out []Candidate
	for _, entry := range Catalog {
		if entry.AbsPath != "" {
			continue
		}
		path := entry.PathFor(home, user.Name)
		if _, err := os.Lstat(path); err == nil {
			out = append(out, Candidate{User: user, Entry: entry, Path: path})
		}
	}
	return out
}

// DiscoverSystem probes the system-level (absolute path) catalog entries
// once; the result is independent of any user.
func DiscoverSystem() []Candidate {
	var out []Candidate
	for i := range Catalog {
		entry := &Catalog[i]
		if entry.AbsPath == "" {
			continue
		}
		path := expandPlaceholders(entry.AbsPath, "", "")
		if _, err := os.Lstat(path); err == nil {
			out = append(out, Candidate{User: User{}, Entry: *entry, Path: path})
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

// FilterExistingWhitelist drops entries that do not exist on this system
// and expands globs (*, ?, [) to every existing match (e.g.
// %HOME%/.config/discord/*/Discord), re-evaluated on every install so app
// updates that move to a new directory are picked up; every match inherits
// the events of the glob pattern that produced it.
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
