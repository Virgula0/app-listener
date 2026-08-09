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
		// ssh writes into known_hosts, so it gets READ,WRITE plus the
		// DELETE/RENAME/HARDLINK its rotation needs: the old known_hosts
		// is linked (or renamed) to .old, a temp file is renamed into
		// place, and stale .old/temp files are unlinked. sshd only reads
		// authorized_keys (READ-only).
		Whitelist: map[string][]string{
			"/usr/bin/ssh":        {"READ", "WRITE", "DELETE", "RENAME", "HARDLINK"},
			"/usr/bin/ssh-add":    nil,
			"/usr/bin/ssh-agent":  nil,
			"/usr/bin/ssh-keygen": nil,
			"/usr/bin/scp":        nil,
			"/usr/bin/sftp":       nil,
			// sshd reads the user's authorized_keys to accept public-key
			// logins; /usr/sbin/sshd covers Debian-style layouts (the
			// whitelist keeps whichever exists per machine). READ-only:
			// the daemon must never write into a user's .ssh.
			"/usr/bin/sshd":  {"READ"},
			"/usr/sbin/sshd": {"READ"},
			// OpenSSH >= 9.8 splits sshd into a dispatcher plus modular
			// helpers: the process that actually opens authorized_keys is
			// sshd-auth (pubkey authentication) and sshd-session (session
			// setup), not the dispatcher. Ubuntu/Debian install them under
			// /usr/lib/openssh/, Arch under /usr/lib/ssh/; the whitelist
			// keeps whichever exists per machine. sftp-server may fetch
			// keys via scp/sftp but is READ-only like sshd: pushing a
			// written authorized_keys into a guarded .ssh over the network
			// stays denied. ssh-keysign/ssh-pkcs11-helper never touch .ssh
			// (entries inert); ssh-session-cleanup is a shell script, so
			// its entry is inert too (the executed ELF is the
			// interpreter).
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
			// git is deliberately NOT whitelisted (see the original
			// ssh-guard config): a compromised git can read the keys.
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
	{Name: "Cursor agent", RelPath: ".cursor-agent",
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
	// IDE configuration directories commonly hold auth tokens, API keys and
	// credentials (VS Code globalStorage, JetBrains options/certs, Zed's
	// auth.json, ...).
	//
	// NOTE: /usr/bin/code (Arch/Debian) is a shell wrapper, not the real
	// binary; the whitelist matches the *executed* ELF, so the actual
	// Electron binary must be listed. Same for Firefox, Chrome/Chromium
	// wrappers, etc.
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
		// Toolbox-installed IDE binaries live under
		// ~/.local/share/JetBrains/Toolbox/apps/<product>/bin/<product> with
		// product-specific names: add them manually in the editor step when
		// the IDE is launched from Toolbox.
		//
		// Standalone installs (e.g. ~/.goland for GoLand) run on the
		// bundled JetBrains Runtime: the process exe is the JBR java, not
		// the launcher ELF, so the JBR and its fsnotifier helper must be
		// whitelisted for the IDE to access its own config.
		Whitelist: map[string][]string{
			"/usr/bin/idea": nil, "/usr/bin/pycharm": nil, "/usr/bin/webstorm": nil,
			"/usr/bin/clion": nil, "/usr/bin/goland": nil, "/usr/bin/phpstorm": nil,
			"/usr/bin/rider": nil, "/usr/bin/datagrip": nil, "/usr/bin/rustrover": nil,
			"/usr/bin/android-studio":     nil,
			"%HOME%/.goland/jbr/bin/java": nil, "%HOME%/.goland/bin/fsnotifier": nil,
		}},
	{Name: "JetBrains Toolbox", RelPath: ".local/share/JetBrains/Toolbox",
		Whitelist: map[string][]string{
			"%HOME%/.local/share/JetBrains/Toolbox/bin/jetbrains-toolbox": nil,
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
		// /usr/bin/gcloud (and gsutil) are Python launchers whose process
		// exe is the interpreter, which the whitelist deliberately does
		// not list. The entries are kept for documentation; they are
		// inert.
		Whitelist: map[string][]string{"/usr/bin/gcloud": nil, "/usr/bin/gsutil": nil}},
	{Name: "Kubernetes kubeconfig", RelPath: ".kube",
		Whitelist: map[string][]string{"/usr/bin/kubectl": nil, "/usr/bin/helm": nil, "/usr/bin/oc": nil}},
	{Name: "Docker config and credentials", RelPath: ".docker",
		Whitelist: map[string][]string{"/usr/bin/docker": nil, "/usr/bin/docker-credential-desktop": nil, "/usr/bin/docker-credential-pass": nil}},
	{Name: "GitHub CLI", RelPath: ".config/gh",
		// gh reads and writes config.yml and hosts.yml (the OAuth tokens)
		// under ~/.config/gh. The whitelist covers the common install
		// layouts and keeps whichever exists per machine: official
		// .deb/.rpm/zypper repos and Arch/Alpine/... package gh to
		// /usr/bin; the tarball and install scripts to /usr/local/bin;
		// Webi and manual user installs to ~/.local/bin; Homebrew on
		// Linux to /home/linuxbrew/.linuxbrew/bin; Snap to /snap/bin.
		// asdf/mise shims are shell scripts (the executed ELF is the
		// interpreter), so they are inert as whitelist entries.
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
	{Name: "npm config file", RelPath: ".npmrc",
		Whitelist: map[string][]string{"/usr/bin/npm": nil, "/usr/bin/npx": nil, "/usr/bin/node": nil}},
	{Name: "npm data", RelPath: ".npm",
		Whitelist: map[string][]string{"/usr/bin/npm": nil, "/usr/bin/npx": nil, "/usr/bin/node": nil}},
	// /usr/bin/pip is a #!/usr/bin/python script, so its inode never
	// matches a running process (the exe is the python interpreter, which
	// is deliberately NOT whitelisted — too broad). The entries are kept
	// for documentation; they are inert.
	{Name: "pip config", RelPath: ".config/pip",
		Whitelist: map[string][]string{"/usr/bin/pip": nil, "/usr/bin/pip3": nil}},
	{Name: "Git config", RelPath: ".gitconfig",
		Whitelist: map[string][]string{"/usr/bin/git": nil}},
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
	{Name: "Netrc credentials file", RelPath: ".netrc",
		Whitelist: map[string][]string{"/usr/bin/curl": nil, "/usr/bin/wget": nil, "/usr/bin/git": nil}},
	{Name: "Wget config", RelPath: ".wgetrc",
		Whitelist: map[string][]string{"/usr/bin/wget": nil}},

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
	// /usr/bin/firefox (Arch), /usr/bin/google-chrome*, /usr/bin/brave,
	// /usr/bin/chromium are shell wrappers; the whitelist matches the
	// executed ELF, so the real binaries under /usr/lib/<browser> and
	// /opt are listed as well. Firefox's crash handler (crashhelper,
	// crashreporter) is a separate binary firefox spawns for the
	// Crash Reports directory: it must be whitelisted for crashes to be
	// recorded.
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
		}},
	{Name: "Chromium profile", RelPath: ".config/chromium",
		Whitelist: map[string][]string{
			"/usr/bin/chromium": nil, "/usr/bin/chromium-browser": nil,
			"/usr/lib/chromium/chrome": nil, "/usr/lib/chromium/chromium": nil,
		}},
	{Name: "Brave profile", RelPath: ".config/BraveSoftware",
		Whitelist: map[string][]string{
			"/usr/bin/brave-browser": nil, "/usr/bin/brave": nil,
			"/opt/brave/brave": nil, "/opt/brave-bin/brave": nil,
		}},
	{Name: "Discord", RelPath: ".config/discord",
		// Discord installs into a versioned directory and re-links
		// .config/discord/Discord on every update. The glob resolves to
		// every installed version and is re-expanded on each install run,
		// so an app update does not leave the whitelist pointing at a
		// stale inode.
		Whitelist: map[string][]string{"%HOME%/.config/discord/*/Discord": nil}},
	{Name: "Discord Canary", RelPath: ".config/discord-canary",
		Whitelist: map[string][]string{"/usr/bin/discord-canary": nil}},
	{Name: "Telegram", RelPath: ".local/share/TelegramDesktop",
		Whitelist: map[string][]string{"/usr/bin/telegram-desktop": nil, "/usr/bin/Telegram": nil}},
	{Name: "Signal", RelPath: ".config/Signal",
		Whitelist: map[string][]string{"/usr/bin/signal-desktop": nil}},
	{Name: "Element", RelPath: ".config/Element",
		Whitelist: map[string][]string{"/usr/bin/element-desktop": nil}},

	// --- Gaming ----------------------------------------------------------------------
	{Name: "Steam data", RelPath: ".local/share/Steam",
		Whitelist: map[string][]string{
			"/usr/bin/steam": nil, "/usr/bin/steamwebhelper": nil,
			"%HOME%/.local/share/Steam/ubuntu12_32/steam":          nil,
			"%HOME%/.local/share/Steam/ubuntu12_32/steamwebhelper": nil,
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

// FilterExistingWhitelist drops whitelist entries that do not exist on
// this system, so the generated config only references real binaries.
// Entries containing a glob metacharacter (*, ? or [) are expanded with
// filepath.Glob and every matching path becomes its own whitelist entry;
// this covers versioned application directories such as
// %HOME%/.config/discord/*/Discord and is re-evaluated on every install,
// so app updates that move to a new directory are picked up. Every match
// inherits the events of the glob pattern that produced it.
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
