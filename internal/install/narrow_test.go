package install

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var gnupgEntry = CandidateDir{Name: "gnupg", Narrowed: map[string]map[string][]string{
	"private-keys-v1.d": {"gpg-agent": nil},
	"openpgp-revocs.d":  {"gpg": {"MKNOD", "WRITE"}},
}}

// gnupgHome makes <tmp>/.gnupg with only private-keys-v1.d (openpgp-revocs.d missing).
func gnupgHome(t *testing.T) (User, string) {
	t.Helper()
	home := t.TempDir()
	dir := filepath.Join(home, ".gnupg")
	if err := os.MkdirAll(filepath.Join(dir, "private-keys-v1.d"), 0o700); err != nil {
		t.Fatal(err)
	}
	return User{Name: "u", Home: home, UID: uint32(os.Getuid()), GID: uint32(os.Getgid())}, dir //nolint:gosec // test ids
}

var gnupgAllow = []BinaryRule{{Path: "/usr/bin/gpg"}, {Path: "/usr/bin/gpg-agent"},
	{Path: "/usr/bin/gpgconf", Events: []string{"READ", "STAT"}}}

func TestGeneratedNarrowsLoad(t *testing.T) {
	_, dir := gnupgHome(t)
	narrows := gnupgEntry.NarrowSectionsFor(dir, gnupgAllow)
	if len(narrows) != 1 || narrows[0].Path != filepath.Join(dir, "private-keys-v1.d") {
		t.Fatalf("a missing narrow must be skipped: %+v", narrows)
	}
	cfg := mustParse(t, GenerateConf([]Section{{Path: dir, Allow: gnupgAllow, Narrows: narrows}}, nil))
	if len(cfg.Resources) != 2 || cfg.Resources[1].NarrowOf != dir {
		t.Fatalf("resources = %+v", cfg.Resources)
	}
	lines := slices.Concat(cfg.Resources[1].Binaries, cfg.Resources[1].PendingBinaries)
	if len(lines) != 1 || filepath.Base(lines[0].Path) != "gpg-agent" {
		t.Fatalf("narrow lines = %+v", lines)
	}
}

func TestSeedNarrowedMakesMissingDirs(t *testing.T) {
	u, dir := gnupgHome(t)
	if err := gnupgEntry.SeedNarrowed(u, dir); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Lstat(filepath.Join(dir, "openpgp-revocs.d"))
	if err != nil || !fi.IsDir() || fi.Mode().Perm() != 0o700 {
		t.Fatalf("seeded: %v %v", fi, err)
	}
	if got := gnupgEntry.NarrowSectionsFor(dir, gnupgAllow); len(got) != 2 {
		t.Fatalf("narrows = %+v", got)
	}
	if err := gnupgEntry.SeedNarrowed(u, "/etc/x"); err != nil {
		t.Fatalf("a system location seeds nothing: %v", err)
	}
}

func TestNarrowAllowIntersectsEvents(t *testing.T) {
	got := narrowAllow(gnupgAllow, map[string][]string{"gpgconf": {"READ", "WRITE"}, "gpg": {"MKNOD"}, "ssh": nil})
	want := []BinaryRule{{Path: "/usr/bin/gpg", Events: []string{"MKNOD"}},
		{Path: "/usr/bin/gpgconf", Events: []string{"READ"}}}
	if !slices.EqualFunc(got, want, func(a, b BinaryRule) bool { return a.Path == b.Path && slices.Equal(a.Events, b.Events) }) {
		t.Fatalf("got %+v", got)
	}
	if got := narrowAllow(gnupgAllow, map[string][]string{"gpgconf": {"WRITE"}}); len(got) != 0 {
		t.Fatalf("no shared event = no line: %+v", got)
	}
}

// A refresh that drops a parent line must not leave a [narrow] the parser refuses.
func TestSyncNarrowsKeepsConfigLoadable(t *testing.T) {
	u, dir := gnupgHome(t)
	hand := filepath.Join(dir, "sshcontrol.d")
	if err := os.Mkdir(hand, 0o700); err != nil {
		t.Fatal(err)
	}
	old := GenerateConf([]Section{{Path: dir, Allow: gnupgAllow,
		Narrows: gnupgEntry.NarrowSectionsFor(dir, gnupgAllow)}}, nil) +
		"\n[narrow \"" + hand + "\"]\n# hand-written\n\"/usr/bin/gpg-agent\"\n\"/usr/bin/gpgconf\" READ,WRITE\n"
	mustParse(t, strings.Replace(old, "READ,WRITE", "READ", 1))
	if err := gnupgEntry.SeedNarrowed(u, dir); err != nil {
		t.Fatal(err)
	}

	fresh := []BinaryRule{{Path: "/usr/bin/gpg", Events: []string{"MKNOD", "STAT"}}, {Path: "/usr/bin/gpgconf",
		Events: []string{"READ"}}}
	text, err := SetSectionWhitelist(old, dir, fresh)
	if err != nil {
		t.Fatal(err)
	}
	text = SyncNarrows(text, dir, &gnupgEntry, fresh)
	cfg := mustParse(t, text)
	byPath := map[string][]string{}
	for _, r := range cfg.Resources {
		for _, b := range slices.Concat(r.Binaries, r.PendingBinaries) {
			ev := make([]string, len(b.Events))
			for i, e := range b.Events {
				ev[i] = e.String()
			}
			byPath[r.Path] = append(byPath[r.Path], filepath.Base(b.Path)+":"+strings.Join(ev, ","))
		}
	}
	for p, want := range map[string][]string{
		filepath.Join(dir, "private-keys-v1.d"): nil,
		filepath.Join(dir, "openpgp-revocs.d"):  {"gpg:MKNOD"},
		hand:                                    {"gpgconf:READ"},
	} {
		if !slices.Equal(byPath[p], want) {
			t.Errorf("%s: %v, want %v\n%s", p, byPath[p], want, text)
		}
	}
	if !strings.Contains(text, "# hand-written") {
		t.Errorf("a hand-written narrow's comments are kept:\n%s", text)
	}
	if again := SyncNarrows(text, dir, &gnupgEntry, fresh); again != text {
		t.Errorf("not idempotent:\n%s\n---\n%s", text, again)
	}
}
