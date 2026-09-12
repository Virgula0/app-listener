package install

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// TestCatalogSanity verifies the master list is well-formed: every entry
// has a name, exactly one non-empty of RelPaths/AbsPaths, absolute
// whitelist binaries, WatchRelPaths only with a single RelPaths entry, and
// no two entries probe the same path.
func TestCatalogSanity(t *testing.T) {
	seen := make(map[string]string)
	for _, entry := range Catalog {
		if entry.Name == "" {
			t.Errorf("entry with empty name: %+v", entry)
		}
		rel, abs := len(entry.RelPaths), len(entry.AbsPaths)
		if (rel == 0) == (abs == 0) {
			t.Errorf("entry %q must set exactly one of RelPaths or AbsPaths (non-empty)", entry.Name)
		}
		for _, p := range entry.RelPaths {
			if p == "" || filepath.IsAbs(p) {
				t.Errorf("entry %q: RelPaths entry %q must be a non-empty home-relative path", entry.Name, p)
			}
		}
		for _, p := range entry.AbsPaths {
			if !filepath.IsAbs(p) {
				t.Errorf("entry %q: AbsPaths entry %q must be absolute", entry.Name, p)
			}
		}
		if len(entry.WatchRelPaths) > 0 && len(entry.RelPaths) != 1 {
			t.Errorf("entry %q: WatchRelPaths is only valid with exactly one RelPaths entry (has %d)", entry.Name, len(entry.RelPaths))
		}
		for _, key := range slices.Concat(entry.RelPaths, entry.AbsPaths) {
			if prev, ok := seen[key]; ok {
				t.Errorf("duplicate catalog path %q (%s and %s)", key, prev, entry.Name)
			}
			seen[key] = entry.Name
		}
		for bin := range entry.Whitelist {
			if !isAbsoluteCandidate(bin) {
				t.Errorf("entry %q: whitelist path %q is not absolute (use %%HOME%%/ or %%USER%%/ for home-relative paths)", entry.Name, bin)
			}
		}
	}
	if len(Catalog) < 30 {
		t.Errorf("catalog unexpectedly small: %d entries", len(Catalog))
	}
}

// TestCatalogMergedResources pins the multi-location entries introduced for
// issue #51: one Name, one shared whitelist, several watch roots. The old
// per-location names must be gone.
func TestCatalogMergedResources(t *testing.T) {
	want := map[string][]string{
		"Claude Code":   {".claude", ".config/claude"},
		"Azure CLI":     {".azure", ".config/azure"},
		"GNOME keyring": {".local/share/keyrings", ".keyring"},
		"Zed editor":    {".config/zed", ".local/share/zed"},
		"Steam":         {".local/share/Steam/config", ".local/share/Steam/userdata/*/config/localconfig.vdf", ".steam/registry.vdf"},
	}
	for name, paths := range want {
		var e *CandidateDir
		for i := range Catalog {
			if Catalog[i].Name == name {
				e = &Catalog[i]
				break
			}
		}
		if e == nil {
			t.Errorf("merged entry %q not found in catalog", name)
			continue
		}
		if !reflect.DeepEqual(e.RelPaths, paths) {
			t.Errorf("%s: RelPaths = %v, want %v", name, e.RelPaths, paths)
		}
	}
	gone := []string{
		"Claude Code config", "Azure CLI config", "GNOME keyring (legacy)",
		"Zed editor data", "Steam (legacy home)",
		"Steam client config (accounts/credentials)",
	}
	for i := range Catalog {
		if slices.Contains(gone, Catalog[i].Name) {
			t.Errorf("entry %q should have been merged away", Catalog[i].Name)
		}
	}
}

// isAbsoluteCandidate reports whether a whitelist path is system-absolute
// or expands to an absolute path via the %HOME%/%USER% placeholders.
func isAbsoluteCandidate(bin string) bool {
	return strings.HasPrefix(bin, "/") ||
		strings.HasPrefix(bin, "%HOME%/") ||
		strings.HasPrefix(bin, "%USER%/")
}

// TestCatalogHasWireGuardSystemEntry verifies the /etc/wireguard entry
// from the original ssh-guard template config is present at system level.
func TestCatalogHasWireGuardSystemEntry(t *testing.T) {
	found := false
	for _, entry := range Catalog {
		if slices.Contains(entry.AbsPaths, "/etc/wireguard") {
			found = true
			for bin := range entry.Whitelist {
				if bin == "/usr/bin/nmcli" {
					return
				}
			}
			t.Errorf("/etc/wireguard entry missing /usr/bin/nmcli whitelist: %v", entry.Whitelist)
		}
	}
	if !found {
		t.Error("/etc/wireguard is not in the catalog")
	}
}

// TestPathsFor verifies placeholder expansion in candidate paths, for both
// a single-location and a multi-location entry.
func TestPathsFor(t *testing.T) {
	single := CandidateDir{RelPaths: []string{".config/opencode"}}
	if got := single.PathsFor("/home/alice", "alice"); len(got) != 1 || got[0] != "/home/alice/.config/opencode" {
		t.Errorf("PathsFor = %q", got)
	}
	multi := CandidateDir{RelPaths: []string{".claude", ".config/claude"}}
	want := []string{"/home/alice/.claude", "/home/alice/.config/claude"}
	if got := multi.PathsFor("/home/alice", "alice"); !reflect.DeepEqual(got, want) {
		t.Errorf("PathsFor = %q, want %q", got, want)
	}
}

// TestExpandWhitelist verifies %USER% placeholders are replaced and the
// input whitelist is not mutated, and that events are carried through.
func TestExpandWhitelist(t *testing.T) {
	entry := CandidateDir{Whitelist: map[string][]string{
		"/home/%USER%/.local/bin/opencode": {"READ", "WRITE"},
		"/usr/bin/ssh":                     nil,
	}}
	got := entry.ExpandWhitelist("bob", "/home/bob")
	if len(got) != 2 {
		t.Fatalf("len = %d", len(got))
	}
	if got[0].Path != "/home/bob/.local/bin/opencode" {
		t.Errorf("got[0].Path = %q", got[0].Path)
	}
	if len(got[0].Events) != 2 {
		t.Errorf("events not preserved: %v", got[0].Events)
	}
	if entry.Whitelist["/home/%USER%/.local/bin/opencode"][0] != "READ" {
		t.Errorf("source whitelist mutated: %v", entry.Whitelist)
	}
}

// TestSSHWhitelistIncludesDaemon verifies the .ssh entry whitelists sshd
// (the OpenSSH server binary), which must be allowed to read a user's
// authorized_keys to accept public-key logins while the directory is
// guarded.
func TestSSHWhitelistIncludesDaemon(t *testing.T) {
	var entry *CandidateDir
	for i := range Catalog {
		if slices.Contains(Catalog[i].RelPaths, ".ssh") {
			entry = &Catalog[i]
			break
		}
	}
	if entry == nil {
		t.Fatal(".ssh entry not found in catalog")
	}
	found := false
	for bin := range entry.Whitelist {
		if bin == "/usr/bin/sshd" || bin == "/usr/sbin/sshd" {
			found = true
			events := entry.Whitelist[bin]
			if len(events) != 1 || events[0] != "READ" {
				t.Errorf("sshd must be READ-only, got %v", events)
			}
		}
	}
	if !found {
		t.Fatal(".ssh whitelist does not include sshd: public-key logins would be denied while guarded")
	}
	if events := entry.Whitelist["/usr/bin/ssh"]; len(events) != 5 || events[0] != "READ" || events[1] != "WRITE" || events[2] != "DELETE" || events[3] != "RENAME" || events[4] != "HARDLINK" {
		t.Errorf("/usr/bin/ssh must be READ,WRITE,DELETE,RENAME,HARDLINK, got %v", events)
	}
}

// TestBuildxWhitelistedInKubeAndDocker verifies that docker-buildx (the
// CLI plugin at /usr/lib/docker/cli-plugins) is whitelisted in both the
// .kube and .docker entries: buildx reads ~/.kube/config for its
// kubernetes driver and ~/.docker/config.json plus the buildx store
// (buildx/.lock, instance state) for its own bookkeeping. Without these
// entries, every docker buildx invocation is denied while the dirs are
// guarded.
func TestBuildxWhitelistedInKubeAndDocker(t *testing.T) {
	const buildx = "/usr/lib/docker/cli-plugins/docker-buildx"
	seen := map[string]bool{}
	for i := range Catalog {
		for _, rel := range Catalog[i].RelPaths {
			if rel != ".kube" && rel != ".docker" {
				continue
			}
			if _, ok := Catalog[i].Whitelist[buildx]; !ok {
				t.Errorf("%s entry does not whitelist docker-buildx", rel)
			}
			seen[rel] = true
		}
	}
	for _, rel := range []string{".kube", ".docker"} {
		if !seen[rel] {
			t.Errorf("%s entry not found in catalog", rel)
		}
	}
}

// TestSSHWhitelistIncludesModularDaemon verifies the .ssh entry whitelists
// the modular OpenSSH server helpers (OpenSSH >= 9.8 splits sshd into a
// dispatcher plus sshd-auth/sshd-session, the processes that actually open
// authorized_keys). Both the Debian (/usr/lib/openssh) and Arch
// (/usr/lib/ssh) layouts must be present and READ-only, like sshd itself.
func TestSSHWhitelistIncludesModularDaemon(t *testing.T) {
	var entry *CandidateDir
	for i := range Catalog {
		if slices.Contains(Catalog[i].RelPaths, ".ssh") {
			entry = &Catalog[i]
			break
		}
	}
	if entry == nil {
		t.Fatal(".ssh entry not found in catalog")
	}
	required := []string{
		"/usr/lib/openssh/sshd-auth", "/usr/lib/ssh/sshd-auth",
		"/usr/lib/openssh/sshd-session", "/usr/lib/ssh/sshd-session",
		"/usr/lib/openssh/sftp-server", "/usr/lib/ssh/sftp-server",
	}
	for _, bin := range required {
		events, ok := entry.Whitelist[bin]
		if !ok {
			t.Errorf(".ssh whitelist is missing %s: public-key logins and scp reads would be denied while guarded", bin)
			continue
		}
		if len(events) != 1 || events[0] != "READ" {
			t.Errorf("%s must be READ-only, got %v", bin, events)
		}
	}
}

// TestGHWhitelistIncludesCommonPaths verifies the .config/gh entry covers
// the common gh install layouts, in particular the user-local
// ~/.local/bin path used by Webi and manual installs (gh would otherwise
// be denied access to its own config while the directory is guarded).
func TestGHWhitelistIncludesCommonPaths(t *testing.T) {
	var entry *CandidateDir
	for i := range Catalog {
		if slices.Contains(Catalog[i].RelPaths, ".config/gh") {
			entry = &Catalog[i]
			break
		}
	}
	if entry == nil {
		t.Fatal(".config/gh entry not found in catalog")
	}
	required := []string{
		"/usr/bin/gh",
		"/usr/local/bin/gh",
		"%HOME%/.local/bin/gh",
		"/home/linuxbrew/.linuxbrew/bin/gh",
		"/snap/bin/gh",
	}
	for _, bin := range required {
		if _, ok := entry.Whitelist[bin]; !ok {
			t.Errorf(".config/gh whitelist is missing %s: gh at that path would be denied access to its config while guarded", bin)
		}
	}
}

// TestDiscoverOnlyExisting verifies that only paths that really exist in
// the user's home are proposed.
func TestDiscoverOnlyExisting(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, ".config", "opencode"), 0o700); err != nil {
		t.Fatal(err)
	}

	user := User{Name: "tester", Home: home}
	got := Discover(user)

	paths := make(map[string]bool)
	for _, c := range got {
		paths[c.Path] = true
	}
	if !paths[filepath.Join(home, ".ssh")] {
		t.Error("expected .ssh to be discovered")
	}
	if !paths[filepath.Join(home, ".config", "opencode")] {
		t.Error("expected .config/opencode to be discovered")
	}
	if paths[filepath.Join(home, ".aws")] {
		t.Error("unexpected .aws: directory does not exist")
	}
}

// TestDiscoverPerUserHomes verifies the catalog is probed for EVERY
// selected user: two users with the same critical directory (e.g. .ssh)
// each get their own per-user candidate path.
func TestDiscoverPerUserHomes(t *testing.T) {
	homeA := t.TempDir()
	homeB := t.TempDir()
	for _, home := range []string{homeA, homeB} {
		if err := os.MkdirAll(filepath.Join(home, ".ssh"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	users := []User{{Name: "alice", Home: homeA}, {Name: "bob", Home: homeB}}

	got := map[string]User{}
	for _, c := range DiscoverForUsers(users) {
		if slices.Contains(c.Entry.RelPaths, ".ssh") {
			got[c.Path] = c.User
		}
	}
	if u, ok := got[filepath.Join(homeA, ".ssh")]; !ok || u.Name != "alice" {
		t.Errorf("missing alice's .ssh (got %v)", got)
	}
	if u, ok := got[filepath.Join(homeB, ".ssh")]; !ok || u.Name != "bob" {
		t.Errorf("missing bob's .ssh (got %v)", got)
	}
}

// TestDiscoverForUsersDeduplicates verifies the same path is proposed only
// once, even when the same user is passed twice.
func TestDiscoverForUsersDeduplicates(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".gnupg"), 0o700); err != nil {
		t.Fatal(err)
	}
	user := User{Name: "tester", Home: home}
	got := DiscoverForUsers([]User{user, user})
	seen := map[string]bool{}
	for _, c := range got {
		if seen[c.Path] {
			t.Errorf("duplicate candidate %s", c.Path)
		}
		seen[c.Path] = true
	}
	if !seen[filepath.Join(home, ".gnupg")] {
		t.Error("expected .gnupg to be discovered")
	}
}

// TestDiscoverSystemOnlyAbsolute verifies system-level entries are probed
// outside the per-user loop and never leak into per-user discovery.
func TestDiscoverSystemOnlyAbsolute(t *testing.T) {
	home := t.TempDir()
	user := User{Name: "tester", Home: home}

	for _, c := range Discover(user) {
		if c.Entry.IsSystem() {
			t.Errorf("per-user Discover returned system entry %s", c.Path)
		}
	}
	for _, c := range DiscoverSystem() {
		if !c.Entry.IsSystem() {
			t.Errorf("DiscoverSystem returned per-user entry %s", c.Path)
		}
		if !slices.Contains(c.Entry.AbsPaths, c.Path) {
			t.Errorf("DiscoverSystem path %q not among AbsPaths %v", c.Path, c.Entry.AbsPaths)
		}
	}
}

// TestFilterExistingWhitelist verifies that non-existent binaries are
// dropped from the generated whitelist.
func TestFilterExistingWhitelist(t *testing.T) {
	existing := filepath.Join(t.TempDir(), "existing-bin")
	if err := os.WriteFile(existing, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	c := Candidate{
		User: User{Name: "tester", Home: "/home/tester"},
		Entry: CandidateDir{Whitelist: map[string][]string{
			existing:                      nil,
			"/definitely/not/here/binary": nil,
		}},
	}
	got := c.FilterExistingWhitelist()
	if len(got) != 1 || got[0].Path != existing {
		t.Errorf("FilterExistingWhitelist = %v, want [%s]", got, existing)
	}
	if len(got[0].Events) != 0 {
		t.Errorf("existing binary should have no event restriction: %v", got[0].Events)
	}
}

// TestExpandWhitelistUsesPasswdHome verifies the %HOME% fix: whitelist
// paths must expand to the real /etc/passwd home directory, never to a
// path derived from the username (a login name whose home differs, e.g.
// user "pwn3r" living in /home/angelo, must not produce /home/pwn3r).
func TestExpandWhitelistUsesPasswdHome(t *testing.T) {
	home := t.TempDir()
	c := Candidate{
		User: User{Name: "pwn3r", Home: home},
		Entry: CandidateDir{Whitelist: map[string][]string{
			"%HOME%/.local/bin/opencode":       nil,
			"%HOME%/.config/discord/*/Discord": {"READ", "WRITE"},
		}},
	}
	got := c.Entry.ExpandWhitelist(c.User.Name, c.User.Home)
	if len(got) != 2 {
		t.Fatalf("ExpandWhitelist returned %d entries, want 2: %v", len(got), got)
	}
	for _, r := range got {
		if strings.HasPrefix(r.Path, "/home/pwn3r") {
			t.Errorf("whitelist path %q derived from the username instead of the passwd home", r.Path)
		}
		if !strings.HasPrefix(r.Path, home) {
			t.Errorf("whitelist path %q does not live under the passwd home %s", r.Path, home)
		}
	}
	// Sorted output: "%HOME%/.config/discord/*/Discord" sorts before
	// "%HOME%/.local/bin/opencode", and it is the entry with events.
	if len(got[0].Events) != 2 || got[0].Events[0] != "READ" {
		t.Errorf("glob pattern events not preserved: %v", got[0].Events)
	}
}

// TestCatalogWhitelistNoUsernameHomeGuard prevents a regression to the
// /home/%USER% pattern, which breaks on systems where the login name does
// not match the home directory basename.
func TestCatalogWhitelistNoUsernameHomePrefix(t *testing.T) {
	for _, entry := range Catalog {
		for bin := range entry.Whitelist {
			if strings.HasPrefix(bin, "/home/%USER%") {
				t.Errorf("entry %q: whitelist %q must use %%HOME%%/, never /home/%%USER%%/", entry.Name, bin)
			}
		}
	}
}

// TestFilterExistingWhitelistGlob verifies that glob patterns in the
// whitelist are expanded to every existing match (versioned app dirs like
// ~/.config/discord/app-*/Discord) and that non-matching patterns produce
// no entries.
func TestFilterExistingWhitelistGlob(t *testing.T) {
	home := t.TempDir()
	mk := func(rel ...string) {
		p := filepath.Join(append([]string{home}, rel...)...)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("ELF"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	mk(".config/discord/app-1.0.150/Discord")
	mk(".config/discord/app-2.0.0/Discord")

	c := Candidate{
		User: User{Name: "tester", Home: home},
		Entry: CandidateDir{Whitelist: map[string][]string{
			"%HOME%/.config/discord/*/Discord":                 {"READ", "WRITE"},
			"%HOME%/.config/discord/no-such-version-*/Discord": nil,
		}},
	}
	got := c.FilterExistingWhitelist()
	want := []string{
		filepath.Join(home, ".config/discord/app-1.0.150/Discord"),
		filepath.Join(home, ".config/discord/app-2.0.0/Discord"),
	}
	if len(got) != len(want) {
		t.Fatalf("FilterExistingWhitelist = %v, want %v", got, want)
	}
	for i := range want {
		if got[i].Path != want[i] {
			t.Errorf("entry[%d] = %q, want %q", i, got[i].Path, want[i])
		}
		if len(got[i].Events) != 2 {
			t.Errorf("entry[%d] events = %v, want the glob pattern's events inherited", i, got[i].Events)
		}
	}
}

// TestFilterExistingWhitelistMixed verifies a whitelist mixing a concrete
// path and a glob keeps the existing concrete path and every glob match.
func TestFilterExistingWhitelistMixed(t *testing.T) {
	home := t.TempDir()
	bin := filepath.Join(home, "real-bin")
	if err := os.WriteFile(bin, []byte("ELF"), 0o755); err != nil {
		t.Fatal(err)
	}
	app := filepath.Join(home, "apps", "v1", "tool")
	if err := os.MkdirAll(filepath.Dir(app), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(app, []byte("ELF"), 0o755); err != nil {
		t.Fatal(err)
	}

	c := Candidate{
		User: User{Name: "tester", Home: home},
		Entry: CandidateDir{Whitelist: map[string][]string{
			bin:                   nil,
			"%HOME%/apps/*/tool":  {"READ"},
			"/definitely/missing": nil,
		}},
	}
	got := c.FilterExistingWhitelist()
	// Sorted output: the glob match apps/v1/tool sorts before real-bin.
	if len(got) != 2 {
		t.Fatalf("FilterExistingWhitelist = %v, want 2 entries", got)
	}
	if got[0].Path != app || len(got[0].Events) != 1 || got[0].Events[0] != "READ" {
		t.Errorf("glob match = %+v, want %s with [READ]", got[0], app)
	}
	if got[1].Path != bin || len(got[1].Events) != 0 {
		t.Errorf("concrete binary = %+v, want %s with no events", got[1], bin)
	}
}

// TestDiscordNarrowedWatches verifies the Discord catalog entry generates
// the narrowed watch set: the sensitive subtrees (token stores, crash
// dumps) are guarded while the update workspace (app dirs, installer.db)
// stays out of the guarded set for the self-updating wrapper.
func TestDiscordNarrowedWatches(t *testing.T) {
	var discord *CandidateDir
	for i := range Catalog {
		if Catalog[i].Name == "Discord" {
			discord = &Catalog[i]
			break
		}
	}
	if discord == nil {
		t.Fatal("Discord catalog entry not found")
	}
	if len(discord.RelPaths) != 1 || discord.RelPaths[0] != ".config/discord" {
		t.Errorf("RelPaths = %q, want [.config/discord] (the encryption root)", discord.RelPaths)
	}
	want := []string{
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
	}
	if !reflect.DeepEqual(discord.WatchRelPaths, want) {
		t.Errorf("WatchRelPaths = %v, want %v", discord.WatchRelPaths, want)
	}
	// The update workspace must NOT be watched: planting binaries there is
	// the updater's job (the wrapper writes the vault root freely).
	for _, watched := range discord.WatchRelPaths {
		if strings.Contains(watched, "app-") || strings.Contains(watched, "installer") {
			t.Errorf("update workspace %q must not be watched", watched)
		}
	}
}

// TestExtraWatchPathsForSkipsMissing is the regression test for a daemon
// startup crash: a fresh Discord install only creates some of the catalog's
// WatchRelPaths sub-directories (e.g. "IndexedDB" may not exist yet), and
// ExtraWatchPathsFor used to return every one of them unconditionally. The
// generated config then declared a `watch:` path that never resolves, and
// the daemon's ResolvePendingPaths pass (daemonconfig.go) treats a grouped
// watch path still missing after its encryption root is available as a
// fatal error — by design, per ResolvePendingPaths' doc comment, "silently
// dropping it would leave a declared-protected directory unguarded". A
// non-existent sub-path must therefore never reach the generated config in
// the first place, exactly like a plain (non-grouped) RelPaths candidate
// that Discover simply does not propose when missing.
func TestExtraWatchPathsForSkipsMissing(t *testing.T) {
	home := t.TempDir()
	root := filepath.Join(home, ".config", "discord")
	existing := filepath.Join(root, "Local Storage")
	if err := os.MkdirAll(existing, 0o755); err != nil {
		t.Fatal(err)
	}
	// "Cookies" and every other WatchRelPaths entry deliberately left absent.

	entry := CandidateDir{
		RelPaths: []string{".config/discord"},
		WatchRelPaths: []string{
			".config/discord/Local Storage",
			".config/discord/Cookies",
		},
	}
	got := entry.ExtraWatchPathsFor(home, "someuser")
	want := []string{existing}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtraWatchPathsFor = %v, want %v (missing sub-paths must be dropped)", got, want)
	}
}
