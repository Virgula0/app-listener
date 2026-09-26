package install

import (
	"strings"
	"testing"
)

var steamBlock = LibraryBlock{
	Name:       "Steam (alice)",
	LibDirs:    []string{"/home/alice/.local/share/Steam/steamrt64"},
	LibWriters: []string{"/home/alice/.local/share/Steam/steamrt64/pv/bin/pressure-vessel-wrap"},
}

// TestGenerateConfPutsLibrariesAfterSections: the library directives are
// declared once, in their own block, after every watch section — never
// nested under one of the application's watch paths.
func TestGenerateConfPutsLibrariesAfterSections(t *testing.T) {
	conf := GenerateConf([]Section{
		{Path: "/home/alice/.local/share/Steam/config", Allow: []BinaryRule{{Path: "/usr/bin/steam"}}, Encrypt: true},
		{Path: "/home/alice/.steam/registry.vdf", Allow: []BinaryRule{{Path: "/usr/bin/steam"}}, Encrypt: true},
	}, []LibraryBlock{steamBlock, {Name: "Empty (alice)"}})

	header := `[libraries "Steam (alice)"]`
	hi := strings.Index(conf, header)
	if hi < 0 {
		t.Fatalf("missing %s:\n%s", header, conf)
	}
	if last := strings.LastIndex(conf, "[watch "); last > hi {
		t.Errorf("a watch section follows the library block:\n%s", conf)
	}
	for _, want := range steamBlock.DirectiveLines() {
		if n := strings.Count(conf, want); n != 1 {
			t.Errorf("%q appears %d times, want once:\n%s", want, n, conf)
		}
		if strings.Index(conf, want) < hi {
			t.Errorf("%q is declared outside its block:\n%s", want, conf)
		}
	}
	if strings.Contains(conf, "Empty (alice)") {
		t.Errorf("an empty block must not be rendered:\n%s", conf)
	}
}

func TestEnsureLibraryBlockAppendsThenMerges(t *testing.T) {
	base := "[watch \"/home/alice/.steam\"]\n\n\"/usr/bin/steam\"\n\nneed_encryption: true\n"

	out := EnsureLibraryBlock(base, &steamBlock)
	if !strings.Contains(out, `[libraries "Steam (alice)"]`) {
		t.Fatalf("block not appended:\n%s", out)
	}
	if again := EnsureLibraryBlock(out, &steamBlock); again != out {
		t.Errorf("merging the same block twice changed the config:\n%s\n---\n%s", out, again)
	}

	// A refresh brings a new writer; a hand-added line in the block survives.
	withHand := strings.Replace(out, "lib_binary", "lib_binary \"/home/alice/hand-added\"\nlib_binary", 1)
	grown := steamBlock
	grown.LibWriters = append(append([]string(nil), steamBlock.LibWriters...), "/home/alice/new-writer")
	merged := EnsureLibraryBlock(withHand, &grown)
	for _, want := range []string{`lib_binary "/home/alice/hand-added"`, `lib_binary "/home/alice/new-writer"`} {
		if !strings.Contains(merged, want) {
			t.Errorf("missing %q after merge:\n%s", want, merged)
		}
	}
}

// TestEnsureLibraryBlockChecksPresenceWithinTheBlock: a lib_dir already nested
// under a watch section by an older installer must STILL be added to the
// block — a block's lib_binary writers reach only the lib_dirs of that block.
func TestEnsureLibraryBlockChecksPresenceWithinTheBlock(t *testing.T) {
	legacy := "[watch \"/home/alice/.steam\"]\n\"/usr/bin/steam\"\n" +
		"lib_dir \"/home/alice/.local/share/Steam/steamrt64\"\nneed_encryption: true\n"
	out := EnsureLibraryBlock(legacy, &steamBlock)
	if n := strings.Count(out, `lib_dir "/home/alice/.local/share/Steam/steamrt64"`); n != 2 {
		t.Errorf("lib_dir must also be declared in the block (want 2 occurrences, got %d):\n%s", n, out)
	}
}

// TestMigrateLegacySectionLibDirectives is the refresh path on a config from
// an older installer: the catalog-generated lines leave the watch section and
// land in the block; a hand-written one stays where the operator put it.
func TestMigrateLegacySectionLibDirectives(t *testing.T) {
	legacy := "[watch \"/home/alice/.steam\"]\n\n\"/usr/bin/steam\"\n" +
		"lib_dir \"/home/alice/.local/share/Steam/steamrt64\"\n" +
		"allow_lib \"/home/alice/hand.so\"\n\nneed_encryption: true\n"
	out, err := RemoveSectionLibDirectives(legacy, "/home/alice/.steam", steamBlock.DirectiveLines())
	if err != nil {
		t.Fatalf("RemoveSectionLibDirectives: %v", err)
	}
	out = EnsureLibraryBlock(out, &steamBlock)

	section := out[:strings.Index(out, "[libraries")]
	if strings.Contains(section, "lib_dir") {
		t.Errorf("the catalog lib_dir must leave the watch section:\n%s", out)
	}
	if !strings.Contains(section, `allow_lib "/home/alice/hand.so"`) {
		t.Errorf("a hand-written directive must stay in its section:\n%s", out)
	}
	if n := strings.Count(out, `lib_dir "/home/alice/.local/share/Steam/steamrt64"`); n != 1 {
		t.Errorf("the lib_dir must end up declared exactly once, got %d:\n%s", n, out)
	}
}

// TestSectionEndStopsAtLibraries: a whitelist refresh of the LAST watch section
// must not treat a following [libraries] block as part of that section.
func TestSectionEndStopsAtLibraries(t *testing.T) {
	conf := GenerateConf([]Section{
		{Path: "/home/alice/.steam", Allow: []BinaryRule{{Path: "/usr/bin/steam"}}, Encrypt: true},
	}, []LibraryBlock{steamBlock})
	out, err := SetSectionWhitelist(conf, "/home/alice/.steam", []BinaryRule{{Path: "/usr/bin/steam"}, {Path: "/usr/bin/lsof"}})
	if err != nil {
		t.Fatalf("SetSectionWhitelist: %v", err)
	}
	hi := strings.Index(out, "[libraries")
	if hi < 0 {
		t.Fatalf("the library block was lost:\n%s", out)
	}
	if strings.Contains(out[hi:], "lsof") {
		t.Errorf("the refreshed whitelist spilled into the library block:\n%s", out)
	}
	for _, want := range steamBlock.DirectiveLines() {
		if !strings.Contains(out[hi:], want) {
			t.Errorf("the refresh dropped %q from the library block:\n%s", want, out)
		}
	}
}

func TestInsertSectionsBeforeLibraries(t *testing.T) {
	conf := GenerateConf([]Section{{Path: "/a", Encrypt: true}}, []LibraryBlock{steamBlock})
	out := InsertSectionsBeforeLibraries(conf, GenerateSections([]Section{{Path: "/b", Encrypt: true}}))
	if strings.Index(out, `[watch "/b"]`) > strings.Index(out, "[libraries") {
		t.Errorf("an appended section must go before the library blocks:\n%s", out)
	}
}
