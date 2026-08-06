// Package install implements the building blocks of the `app-listener
// install` wizard: the catalog of critical directories, user enumeration,
// the recursive tree copy used during fscrypt migrations, and the daemon
// config generation.
package install

import (
	"os"
	"path/filepath"
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
	// Whitelist lists the absolute binary paths that are allowed to
	// access this directory by default. Keep it minimal: every binary
	// added here is a potential privilege-escalation path if it can be
	// abused while the directory is unlocked. Identity is verified by
	// inode, not by name. Whitelist entries that do not exist on the
	// target system are dropped when the config is generated; entries
	// containing a glob (*, ?, [) are expanded to every existing match.
	Whitelist []string
}

// Catalog is the master list of critical directories to search for each
// selected user. Add, remove or tweak entries here — this is the single
// place that drives discovery. Paths that do not exist for a user are
// simply not proposed.
var Catalog = []CandidateDir{
	// --- SSH and remote access -------------------------------------------------
	{Name: "SSH client configuration and keys", RelPath: ".ssh",
		Whitelist: []string{
			"/usr/bin/ssh", "/usr/bin/ssh-agent", "/usr/bin/ssh-keygen",
			"/usr/bin/ssh-add", "/usr/bin/scp", "/usr/bin/sftp",
			// git is deliberately NOT whitelisted (see the original
			// ssh-guard config): a compromised git can read the keys.
		}},
	{Name: "GNU Privacy Guard keyring", RelPath: ".gnupg",
		Whitelist: []string{
			"/usr/bin/gpg", "/usr/bin/gpg-agent", "/usr/bin/gpgconf",
			"/usr/bin/gpg-connect-agent",
		}},

	// --- AI coding agents and CLI tools ----------------------------------------
	{Name: "opencode", RelPath: ".config/opencode",
		Whitelist: []string{
			"/usr/local/bin/opencode", "/usr/bin/opencode",
			"%HOME%/.local/bin/opencode",
		}},
	{Name: "code CLI (GitHub)", RelPath: ".config/code-cli",
		Whitelist: []string{
			"/usr/bin/code-cli", "/usr/local/bin/code-cli",
			"%HOME%/.local/bin/code-cli",
		}},
	{Name: "Claude Code", RelPath: ".claude",
		Whitelist: []string{
			"/usr/local/bin/claude", "/usr/bin/claude",
			"%HOME%/.local/bin/claude",
		}},
	{Name: "Claude Code config", RelPath: ".config/claude",
		Whitelist: []string{
			"/usr/local/bin/claude", "/usr/bin/claude",
			"%HOME%/.local/bin/claude",
		}},
	{Name: "Codeium", RelPath: ".codeium",
		Whitelist: []string{
			"/usr/local/bin/codeium", "/usr/bin/codeium",
			"%HOME%/.local/bin/codeium",
		}},
	{Name: "Gemini CLI", RelPath: ".gemini",
		Whitelist: []string{
			"/usr/local/bin/gemini", "/usr/bin/gemini",
			"%HOME%/.local/bin/gemini",
		}},
	{Name: "Cursor", RelPath: ".cursor",
		Whitelist: []string{
			"/usr/bin/cursor", "/usr/local/bin/cursor",
			"%HOME%/.local/bin/cursor",
		}},
	{Name: "Cursor agent", RelPath: ".cursor-agent",
		Whitelist: []string{
			"/usr/bin/cursor", "/usr/local/bin/cursor",
			"%HOME%/.local/bin/cursor",
		}},
	{Name: "GitHub Copilot", RelPath: ".config/github-copilot",
		Whitelist: []string{
			"/usr/bin/git", "/usr/local/bin/copilot",
			"%HOME%/.local/bin/copilot",
		}},

	// --- IDEs and editors ---------------------------------------------------------
	// IDE configuration directories commonly hold auth tokens, API keys and
	// credentials (VS Code globalStorage, JetBrains options/certs, Zed's
	// auth.json, ...).
	//
	// NOTE: /usr/bin/code (Arch/Debian) is a shell wrapper, not the real
	// binary; the whitelist matches the *executed* ELF, so the actual
	// Electron binary must be listed. Same for Firefox, Chrome/Chromium
	// wrappers, etc.
	{Name: "VS Code config", RelPath: ".config/Code",
		Whitelist: []string{
			"/usr/bin/code", "/usr/local/bin/code",
			"%HOME%/.local/bin/code",
			"/opt/visual-studio-code/code", "/usr/share/code/code", "/opt/visual-studio-code/bin/code",
			// chrome_crashpad_handler is a separate process Code spawns to
			// write crash dumps into .config/Code/Crashpad.
			"/opt/visual-studio-code/chrome_crashpad_handler", "/usr/share/code/chrome_crashpad_handler",
		}},
	{Name: "VS Code Insiders config", RelPath: ".config/Code - Insiders",
		Whitelist: []string{
			"/usr/bin/code-insiders", "/usr/local/bin/code-insiders",
			"%HOME%/.local/bin/code-insiders", "/opt/visual-studio-code/bin/code",
		}},
	{Name: "VS Code server (remote development)", RelPath: ".vscode-server",
		Whitelist: []string{
			"/usr/bin/code", "/usr/local/bin/code",
			"%HOME%/.local/bin/code",
			"/opt/visual-studio-code/code", "/usr/share/code/code",
			// chrome_crashpad_handler is a separate process Code spawns to
			// write crash dumps into .config/Code/Crashpad.
			"/opt/visual-studio-code/chrome_crashpad_handler", "/usr/share/code/chrome_crashpad_handler",
		}},
	{Name: "VSCodium config", RelPath: ".config/VSCodium",
		Whitelist: []string{"/usr/bin/codium", "/usr/local/bin/codium",
			"/opt/vscodium/chrome_crashpad_handler", "/usr/share/vscodium/chrome_crashpad_handler"}},
	{Name: "JetBrains IDEs", RelPath: ".config/JetBrains",
		// Toolbox-installed IDE binaries live under
		// ~/.local/share/JetBrains/Toolbox/apps/<product>/bin/<product> with
		// product-specific names: add them manually in the editor step when
		// the IDE is launched from Toolbox.
		//
		// Standalone installs (e.g. ~/.goland for GoLand) run on the
		// bundled JetBrains Runtime: the process exe is the JBR java, not
		// the launcher ELF, so the JBR and its fsnotifier helper must be
		// whitelisted for the IDE to access its own config.
		Whitelist: []string{
			"/usr/bin/idea", "/usr/bin/pycharm", "/usr/bin/webstorm",
			"/usr/bin/clion", "/usr/bin/goland", "/usr/bin/phpstorm",
			"/usr/bin/rider", "/usr/bin/datagrip", "/usr/bin/rustrover",
			"/usr/bin/android-studio",
			"%HOME%/.goland/jbr/bin/java", "%HOME%/.goland/bin/fsnotifier",
		}},
	{Name: "JetBrains Toolbox", RelPath: ".local/share/JetBrains/Toolbox",
		Whitelist: []string{
			"%HOME%/.local/share/JetBrains/Toolbox/bin/jetbrains-toolbox",
		}},
	{Name: "Zed editor", RelPath: ".config/zed",
		Whitelist: []string{
			"/usr/bin/zed", "/usr/local/bin/zed",
			"%HOME%/.local/bin/zed",
		}},
	{Name: "Zed editor data", RelPath: ".local/share/zed",
		Whitelist: []string{
			"/usr/bin/zed", "/usr/local/bin/zed",
			"%HOME%/.local/bin/zed",
		}},
	{Name: "Sublime Text", RelPath: ".config/sublime-text",
		// /usr/bin/subl is a wrapper; the real binary is under /opt.
		Whitelist: []string{"/usr/bin/subl", "/usr/bin/sublime-text", "/usr/bin/sublime_text",
			"/opt/sublime_text/sublime_text"}},
	{Name: "Insomnia API client", RelPath: ".config/Insomnia",
		Whitelist: []string{"/usr/bin/insomnia", "/opt/insomnia/insomnia"}},
	{Name: "Postman API client", RelPath: ".config/Postman",
		Whitelist: []string{"/usr/bin/postman", "/opt/postman/postman"}},
	{Name: "DBeaver database client", RelPath: ".local/share/DBeaverData",
		// /usr/bin/dbeaver is a wrapper; the real binary lives in /usr/lib.
		Whitelist: []string{"/usr/bin/dbeaver", "/usr/lib/dbeaver/dbeaver"}},
	{Name: "Android adb keys", RelPath: ".android",
		Whitelist: []string{"/usr/bin/adb"}},

	// --- Cloud, containers and dev tooling --------------------------------------
	{Name: "AWS credentials", RelPath: ".aws",
		Whitelist: []string{"/usr/bin/aws"}},
	{Name: "Google Cloud SDK", RelPath: ".config/gcloud",
		// /usr/bin/gcloud (and gsutil) are Python launchers whose process
		// exe is the interpreter, which the whitelist deliberately does
		// not list. The entries are kept for documentation; they are
		// inert.
		Whitelist: []string{"/usr/bin/gcloud", "/usr/bin/gsutil"}},
	{Name: "Kubernetes kubeconfig", RelPath: ".kube",
		Whitelist: []string{"/usr/bin/kubectl", "/usr/bin/helm", "/usr/bin/oc"}},
	{Name: "Docker config and credentials", RelPath: ".docker",
		Whitelist: []string{"/usr/bin/docker", "/usr/bin/docker-credential-desktop", "/usr/bin/docker-credential-pass"}},
	{Name: "GitHub CLI", RelPath: ".config/gh",
		Whitelist: []string{"/usr/bin/gh"}},
	{Name: "Azure CLI", RelPath: ".azure",
		Whitelist: []string{"/usr/bin/az"}},
	{Name: "Azure CLI config", RelPath: ".config/azure",
		Whitelist: []string{"/usr/bin/az"}},
	{Name: "rclone config", RelPath: ".config/rclone",
		Whitelist: []string{"/usr/bin/rclone"}},
	{Name: "Ollama config", RelPath: ".ollama",
		Whitelist: []string{"/usr/bin/ollama"}},
	{Name: "npm config file", RelPath: ".npmrc",
		Whitelist: []string{"/usr/bin/npm", "/usr/bin/npx", "/usr/bin/node"}},
	{Name: "npm data", RelPath: ".npm",
		Whitelist: []string{"/usr/bin/npm", "/usr/bin/npx", "/usr/bin/node"}},
	// /usr/bin/pip is a #!/usr/bin/python script, so its inode never
	// matches a running process (the exe is the python interpreter, which
	// is deliberately NOT whitelisted — too broad). The entries are kept
	// for documentation; they are inert.
	{Name: "pip config", RelPath: ".config/pip",
		Whitelist: []string{"/usr/bin/pip", "/usr/bin/pip3"}},
	{Name: "Git config", RelPath: ".gitconfig",
		Whitelist: []string{"/usr/bin/git"}},
	{Name: "Git config directory", RelPath: ".config/git",
		Whitelist: []string{"/usr/bin/git"}},

	// --- Password managers and secrets -------------------------------------------
	{Name: "password-store (pass)", RelPath: ".password-store",
		Whitelist: []string{"/usr/bin/pass", "/usr/bin/gpg"}},
	{Name: "gopass", RelPath: ".local/share/gopass",
		Whitelist: []string{"/usr/bin/gopass"}},
	{Name: "Bitwarden", RelPath: ".config/Bitwarden",
		Whitelist: []string{"/usr/bin/bitwarden", "/usr/local/bin/bitwarden"}},
	{Name: "1Password", RelPath: ".1password",
		Whitelist: []string{"/usr/bin/1password", "/usr/bin/op", "/usr/local/bin/op"}},
	{Name: "age keys", RelPath: ".config/age",
		Whitelist: []string{"/usr/bin/age", "/usr/bin/age-keygen"}},
	{Name: "KeePassXC", RelPath: ".local/share/keepassxc",
		Whitelist: []string{"/usr/bin/keepassxc"}},
	{Name: "GNOME keyring", RelPath: ".local/share/keyrings",
		Whitelist: []string{"/usr/bin/gnome-keyring-daemon", "/usr/bin/gnome-keyring"}},
	{Name: "GNOME keyring (legacy)", RelPath: ".keyring",
		Whitelist: []string{"/usr/bin/gnome-keyring-daemon", "/usr/bin/gnome-keyring"}},
	{Name: "Chezmoi state", RelPath: ".local/share/chezmoi",
		Whitelist: []string{"/usr/bin/chezmoi"}},
	{Name: "Netrc credentials file", RelPath: ".netrc",
		Whitelist: []string{"/usr/bin/curl", "/usr/bin/wget", "/usr/bin/git"}},
	{Name: "Wget config", RelPath: ".wgetrc",
		Whitelist: []string{"/usr/bin/wget"}},

	// --- VPNs ---------------------------------------------------------------------
	{Name: "NordVPN config", RelPath: ".config/nordvpn",
		Whitelist: []string{"/usr/bin/nordvpn"}},
	{Name: "Mullvad VPN config", RelPath: ".config/mullvad",
		Whitelist: []string{"/usr/bin/mullvad"}},
	{Name: "ProtonVPN config", RelPath: ".config/protonvpn",
		Whitelist: []string{"/usr/bin/protonvpn", "/usr/bin/protonvpn-cli"}},
	{Name: "OpenVPN config", RelPath: ".config/openvpn",
		Whitelist: []string{"/usr/bin/openvpn"}},
	{Name: "WireGuard config (per-user)", RelPath: ".config/wireguard",
		Whitelist: []string{"/usr/bin/wg", "/usr/bin/wg-quick"}},

	// --- Browsers and messaging -----------------------------------------------------
	// /usr/bin/firefox (Arch), /usr/bin/google-chrome*, /usr/bin/brave,
	// /usr/bin/chromium are shell wrappers; the whitelist matches the
	// executed ELF, so the real binaries under /usr/lib/<browser> and
	// /opt are listed as well. Firefox's crash handler (crashhelper,
	// crashreporter) is a separate binary firefox spawns for the
	// Crash Reports directory: it must be whitelisted for crashes to be
	// recorded.
	{Name: "Firefox profile", RelPath: ".mozilla/firefox",
		Whitelist: []string{
			"/usr/bin/firefox",
			"/usr/lib/firefox/firefox",
			"/usr/lib/firefox/crashhelper", "/usr/lib/firefox/crashreporter",
		}},
	{Name: "Google Chrome profile", RelPath: ".config/google-chrome",
		Whitelist: []string{
			"/usr/bin/google-chrome", "/usr/bin/google-chrome-stable",
			"/opt/google/chrome/chrome", "/usr/lib/chromium/chrome",
		}},
	{Name: "Chromium profile", RelPath: ".config/chromium",
		Whitelist: []string{
			"/usr/bin/chromium", "/usr/bin/chromium-browser",
			"/usr/lib/chromium/chrome", "/usr/lib/chromium/chromium",
		}},
	{Name: "Brave profile", RelPath: ".config/BraveSoftware",
		Whitelist: []string{
			"/usr/bin/brave-browser", "/usr/bin/brave",
			"/opt/brave/brave", "/opt/brave-bin/brave",
		}},
	{Name: "Discord", RelPath: ".config/discord",
		// Discord installs into a versioned directory and re-links
		// .config/discord/Discord on every update. The glob resolves to
		// every installed version and is re-expanded on each install run,
		// so an app update does not leave the whitelist pointing at a
		// stale inode.
		Whitelist: []string{"%HOME%/.config/discord/*/Discord"}},
	{Name: "Discord Canary", RelPath: ".config/discord-canary",
		Whitelist: []string{"/usr/bin/discord-canary"}},
	{Name: "Telegram", RelPath: ".local/share/TelegramDesktop",
		Whitelist: []string{"/usr/bin/telegram-desktop", "/usr/bin/Telegram"}},
	{Name: "Signal", RelPath: ".config/Signal",
		Whitelist: []string{"/usr/bin/signal-desktop"}},
	{Name: "Element", RelPath: ".config/Element",
		Whitelist: []string{"/usr/bin/element-desktop"}},

	// --- Gaming ----------------------------------------------------------------------
	{Name: "Steam data", RelPath: ".local/share/Steam",
		Whitelist: []string{
			"/usr/bin/steam", "/usr/bin/steamwebhelper",
			"%HOME%/.local/share/Steam/ubuntu12_32/steam",
			"%HOME%/.local/share/Steam/ubuntu12_32/steamwebhelper",
		}},
	{Name: "Steam (legacy home)", RelPath: ".steam",
		Whitelist: []string{"/usr/bin/steam", "/usr/bin/steamwebhelper"}},

	// --- System-level paths (probed once, not per user) --------------------------------
	// Same whitelist as the original ssh-guard template:
	{Name: "WireGuard system config", AbsPath: "/etc/wireguard",
		Whitelist: []string{"/usr/bin/nmcli"}},
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

// ExpandWhitelist returns the whitelist with %USER% placeholders replaced.
func (c CandidateDir) ExpandWhitelist(user, home string) []string {
	out := make([]string, 0, len(c.Whitelist))
	for _, bin := range c.Whitelist {
		out = append(out, expandPlaceholders(bin, user, home))
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

// FilterExistingWhitelist drops whitelist entries that do not exist on
// this system, so the generated config only references real binaries.
// Entries containing a glob metacharacter (*, ? or [) are expanded with
// filepath.Glob and every matching path becomes its own whitelist entry;
// this covers versioned application directories such as
// %HOME%/.config/discord/*/Discord and is re-evaluated on every install,
// so app updates that move to a new directory are picked up.
func (c *Candidate) FilterExistingWhitelist() []string {
	var out []string
	for _, bin := range c.Entry.ExpandWhitelist(c.User.Name, c.User.Home) {
		if strings.ContainsAny(bin, "*?[") {
			matches, err := filepath.Glob(bin)
			if err != nil {
				continue // malformed pattern: skip the whole entry
			}
			for _, m := range matches {
				if _, err := os.Stat(m); err == nil {
					out = append(out, m)
				}
			}
			continue
		}
		if _, err := os.Stat(bin); err == nil {
			out = append(out, bin)
		}
	}
	return out
}
