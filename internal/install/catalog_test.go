package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCatalogSanity verifies the master list is well-formed: every entry
// has a name, exactly one of RelPath/AbsPath, absolute whitelist binaries,
// and no two entries probe the same path.
func TestCatalogSanity(t *testing.T) {
	seen := make(map[string]string)
	for _, entry := range Catalog {
		if entry.Name == "" {
			t.Errorf("entry %q/%q has an empty name", entry.RelPath, entry.AbsPath)
		}
		if (entry.RelPath == "") == (entry.AbsPath == "") {
			t.Errorf("entry %q must set exactly one of RelPath or AbsPath", entry.Name)
		}
		if entry.RelPath != "" && filepath.IsAbs(entry.RelPath) {
			t.Errorf("entry %q: RelPath must be relative to the home directory, got %q", entry.Name, entry.RelPath)
		}
		if entry.AbsPath != "" && !filepath.IsAbs(entry.AbsPath) {
			t.Errorf("entry %q: AbsPath must be absolute, got %q", entry.Name, entry.AbsPath)
		}
		key := entry.RelPath
		if key == "" {
			key = entry.AbsPath
		}
		if prev, ok := seen[key]; ok {
			t.Errorf("duplicate catalog path %q (%s and %s)", key, prev, entry.Name)
		}
		seen[key] = entry.Name
		for _, bin := range entry.Whitelist {
			if !isAbsoluteCandidate(bin) {
				t.Errorf("entry %q: whitelist path %q is not absolute (use %%HOME%%/ or %%USER%%/ for home-relative paths)", entry.Name, bin)
			}
		}
	}
	if len(Catalog) < 30 {
		t.Errorf("catalog unexpectedly small: %d entries", len(Catalog))
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
		if entry.AbsPath == "/etc/wireguard" {
			found = true
			for _, bin := range entry.Whitelist {
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

// TestPathFor verifies placeholder expansion in candidate paths.
func TestPathFor(t *testing.T) {
	entry := CandidateDir{RelPath: ".config/opencode"}
	if got := entry.PathFor("/home/alice", "alice"); got != "/home/alice/.config/opencode" {
		t.Errorf("PathFor = %q", got)
	}
}

// TestExpandWhitelist verifies %USER% placeholders are replaced and the
// input slice is not mutated.
func TestExpandWhitelist(t *testing.T) {
	entry := CandidateDir{Whitelist: []string{"/home/%USER%/.local/bin/opencode", "/usr/bin/ssh"}}
	got := entry.ExpandWhitelist("bob", "/home/bob")
	if len(got) != 2 {
		t.Fatalf("len = %d", len(got))
	}
	if got[0] != "/home/bob/.local/bin/opencode" {
		t.Errorf("got[0] = %q", got[0])
	}
	if entry.Whitelist[0] != "/home/%USER%/.local/bin/opencode" {
		t.Errorf("source whitelist mutated: %q", entry.Whitelist[0])
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
		if c.Entry.RelPath == ".ssh" {
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
		if c.Entry.AbsPath != "" {
			t.Errorf("per-user Discover returned system entry %s", c.Path)
		}
	}
	for _, c := range DiscoverSystem() {
		if c.Entry.AbsPath == "" {
			t.Errorf("DiscoverSystem returned per-user entry %s", c.Path)
		}
		if c.Path != c.Entry.AbsPath {
			t.Errorf("DiscoverSystem path %q != AbsPath %q", c.Path, c.Entry.AbsPath)
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
		Entry: CandidateDir{Whitelist: []string{
			existing,
			"/definitely/not/here/binary",
		}},
	}
	got := c.FilterExistingWhitelist()
	if len(got) != 1 || got[0] != existing {
		t.Errorf("FilterExistingWhitelist = %v, want [%s]", got, existing)
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
		Entry: CandidateDir{Whitelist: []string{
			"%HOME%/.local/bin/opencode",
			"%HOME%/.config/discord/*/Discord",
		}},
	}
	got := c.Entry.ExpandWhitelist(c.User.Name, c.User.Home)
	if len(got) != 2 {
		t.Fatalf("ExpandWhitelist returned %d entries, want 2: %v", len(got), got)
	}
	for _, p := range got {
		if strings.HasPrefix(p, "/home/pwn3r") {
			t.Errorf("whitelist path %q derived from the username instead of the passwd home", p)
		}
		if !strings.HasPrefix(p, home) {
			t.Errorf("whitelist path %q does not live under the passwd home %s", p, home)
		}
	}
}

// TestCatalogWhitelistNoUsernameHomeGuard prevents a regression to the
// /home/%USER% pattern, which breaks on systems where the login name does
// not match the home directory basename.
func TestCatalogWhitelistNoUsernameHomePrefix(t *testing.T) {
	for _, entry := range Catalog {
		for _, bin := range entry.Whitelist {
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
		Entry: CandidateDir{Whitelist: []string{
			"%HOME%/.config/discord/*/Discord",
			"%HOME%/.config/discord/no-such-version-*/Discord",
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
		if got[i] != want[i] {
			t.Errorf("entry[%d] = %q, want %q", i, got[i], want[i])
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
		Entry: CandidateDir{Whitelist: []string{
			bin,
			"%HOME%/apps/*/tool",
			"/definitely/missing",
		}},
	}
	got := c.FilterExistingWhitelist()
	want := []string{bin, app}
	if len(got) != len(want) {
		t.Fatalf("FilterExistingWhitelist = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
