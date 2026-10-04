package daemon

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/install"
)

func mkdirs(t *testing.T, paths ...string) {
	t.Helper()
	for _, p := range paths {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func bitOf(t *testing.T, r guard.GlobReservations, name string) uint64 {
	t.Helper()
	i := slices.Index(r.Patterns, name)
	if i < 0 {
		t.Fatalf("pattern %q not reserved; have %v", name, r.Patterns)
	}
	return 1 << uint(i)
}

func TestBuildGlobReservations_Steam(t *testing.T) {
	home := t.TempDir()
	steam := filepath.Join(home, ".local/share/Steam")
	common := filepath.Join(steam, "steamapps/common")
	// compatibilitytools.d is missing: its globs must fall back to the deepest existing dir.
	mkdirs(t, filepath.Join(steam, "config"), common)
	client := filepath.Join(steam, "ubuntu12_32/steam")
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(steam, "config"), Binaries: []daemonconfig.BinaryRule{{Path: client}}},
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	ws := bitOf(t, r, "wineserver")
	if r.Roots[common]&ws == 0 {
		t.Errorf("wineserver not reserved below %s: %v", common, r.Roots)
	}
	if r.Roots[steam]&ws == 0 {
		t.Errorf("compatibilitytools.d/*/files/bin/wineserver must fall back to %s: %v", steam, r.Roots)
	}
	if r.Roots[steam]&bitOf(t, r, "steam") == 0 {
		t.Errorf("*/steam not reserved below %s", steam)
	}
	if r.Writers[client]&ws == 0 {
		t.Errorf("the Steam client must be a writer of wineserver: %v", r.Writers)
	}
	if _, ok := r.Writers["/usr/bin/ssh"]; ok {
		t.Errorf("another entry's binary must not write Steam's reserved names")
	}
	wantChild := guard.GlobChild{Parent: filepath.Join(steam, "steamapps"), Name: "common"}
	found := false
	for _, c := range r.Children {
		if c.Parent == wantChild.Parent && c.Name == wantChild.Name && c.Bits&ws != 0 {
			found = true
		}
	}
	if !found {
		t.Errorf("existing root %s must be reserved in its parent: %+v", common, r.Children)
	}
}

func TestBuildGlobReservations_UnconfiguredEntryReservesNothing(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, ".local/share/Steam/steamapps/common"))
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	if len(r.Patterns) != 0 || len(r.Roots) != 0 || len(r.Writers) != 0 {
		t.Fatalf("Steam is not configured, so nothing may be reserved: %+v", r)
	}
}

// A glob whose root lies inside a guarded resource is already write-gated by that resource.
func TestBuildGlobReservations_InTreeRootSkipped(t *testing.T) {
	home := t.TempDir()
	foundry := filepath.Join(home, ".foundry")
	mkdirs(t, filepath.Join(foundry, "bin"))
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: foundry, Binaries: []daemonconfig.BinaryRule{{Path: filepath.Join(foundry, "bin/forge")}}},
	}}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	if len(r.Patterns) != 0 {
		t.Fatalf("foundry's bin/* is inside its guarded tree, nothing to reserve: %v", r.Patterns)
	}
}

func discordLibConfig(home string) (*daemonconfig.Config, string) {
	bin := filepath.Join(home, ".config/discord/0.0.1/Discord")
	return &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".config/discord/sentry"), Binaries: []daemonconfig.BinaryRule{{Path: bin}}},
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}, bin
}

func TestBuildGlobReservations_DiscordLibs(t *testing.T) {
	home := t.TempDir()
	discord := filepath.Join(home, ".config/discord")
	mkdirs(t, discord)
	cfg, bin := discordLibConfig(home)
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	for _, name := range []string{"*.so", "lib*", "*.node"} {
		bit := bitOf(t, r, name)
		if r.Roots[discord]&bit == 0 {
			t.Errorf("%s not reserved below %s: %v", name, discord, r.Roots)
		}
		if r.Writers[bin]&bit == 0 {
			t.Errorf("Discord must be a writer of %s: %v", name, r.Writers)
		}
		if r.Writers["/usr/bin/ssh"]&bit != 0 {
			t.Errorf("another entry's binary must not write (or load) Discord's %s", name)
		}
	}
	lib := bitOf(t, r, "lib*")
	found := false
	for _, c := range r.Children {
		if c.Parent == filepath.Join(home, ".config") && c.Name == "discord" && c.Bits&lib != 0 {
			found = true
		}
	}
	if !found {
		t.Errorf("the library root must be reserved in its parent: %+v", r.Children)
	}
}

// Falling back to ~/.config would reserve lib* for every app there.
func TestBuildGlobReservations_MissingLibRootNotWidened(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, ".config"))
	cfg, _ := discordLibConfig(home)
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	lib := bitOf(t, r, "lib*")
	for root, bits := range r.Roots {
		if bits&lib != 0 {
			t.Errorf("lib* reserved below %s although %s/.config/discord is missing", root, home)
		}
	}
}

func TestGlobBuilder_LibBitsPerEntry(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, "a"), filepath.Join(home, "b"))
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".a"), Binaries: []daemonconfig.BinaryRule{{Path: "/bin/a"}}},
		{Path: filepath.Join(home, ".b"), Binaries: []daemonconfig.BinaryRule{{Path: "/bin/b"}}},
	}}
	b := globBuilder{
		cfg:      cfg,
		r:        guard.GlobReservations{Roots: map[string]uint64{}, Writers: map[string]uint64{}},
		bitOf:    map[string]int{},
		children: map[[2]string]uint64{},
	}
	u := install.User{Name: "u", Home: home}
	b.addEntry(&install.CandidateDir{Name: "A", RelPaths: []string{".a"}, ReservedLibs: []string{"%HOME%/a/*.so"}}, u)
	b.addEntry(&install.CandidateDir{Name: "B", RelPaths: []string{".b"}, ReservedLibs: []string{"%HOME%/b/*.so"}}, u)

	if len(b.r.Patterns) != 2 {
		t.Fatalf("each entry's *.so needs its own bit: %v", b.r.Patterns)
	}
	if b.r.Writers["/bin/a"]&b.r.Roots[filepath.Join(home, "b")] != 0 {
		t.Errorf("A's writer may plant or load *.so below B's root: writers %v roots %v", b.r.Writers, b.r.Roots)
	}
}

// A lib_dir glob match is guarded by the next catalog refresh with whatever it already holds, so
// its wildcard component must be reserved for the entry's writers like a whitelist glob.
func TestBuildGlobReservations_SteamLibDirGlobs(t *testing.T) {
	home := t.TempDir()
	steam := filepath.Join(home, ".local/share/Steam")
	common := filepath.Join(steam, "steamapps/common")
	mkdirs(t, filepath.Join(steam, "config"), common)
	client := filepath.Join(steam, "ubuntu12_32/steam")
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(steam, "config"), Binaries: []daemonconfig.BinaryRule{{Path: client}}},
	}}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	for _, name := range []string{"Proton*", "SteamLinuxRuntime_*"} {
		bit := bitOf(t, r, name)
		if r.Roots[common]&bit == 0 {
			t.Errorf("%s not reserved below %s: %v", name, common, r.Roots)
		}
		if r.Writers[client]&bit == 0 {
			t.Errorf("the Steam client must be a writer of %s: %v", name, r.Writers)
		}
	}
}

func hasChild(r guard.GlobReservations, parent, name string, bit uint64) bool {
	return slices.ContainsFunc(r.Children, func(c guard.GlobChild) bool {
		return c.Parent == parent && c.Name == name && c.Bits&bit != 0
	})
}

func TestBuildGlobReservations_FixedBinaryAncestorsAndSymlinkTarget(t *testing.T) {
	home := t.TempDir()
	versions := filepath.Join(home, ".local/share/claude/versions")
	bin := filepath.Join(home, ".local/bin")
	mkdirs(t, filepath.Join(home, ".claude"), versions, bin)
	target := filepath.Join(versions, "2.1.3")
	if err := os.WriteFile(target, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(bin, "claude")); err != nil {
		t.Fatal(err)
	}
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".claude"), Binaries: []daemonconfig.BinaryRule{{Path: target}}},
	}}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	claude := bitOf(t, r, "claude")
	if r.Roots[bin]&claude == 0 {
		t.Errorf("the binary's name must be reserved in %s: %v", bin, r.Roots)
	}
	for _, c := range [][2]string{{home, ".local"}, {filepath.Join(home, ".local"), "bin"}} {
		if !hasChild(r, c[0], c[1], claude) {
			t.Errorf("%s/%s must be reserved: %+v", c[0], c[1], r.Children)
		}
	}
	ver := bitOf(t, r, "2.1.3")
	for _, c := range [][2]string{{filepath.Join(home, ".local/share"), "claude"}, {filepath.Join(home, ".local/share/claude"), "versions"}} {
		if !hasChild(r, c[0], c[1], ver) {
			t.Errorf("symlink target directory %s/%s must be reserved: %+v", c[0], c[1], r.Children)
		}
	}
	if r.Writers[target]&(claude|ver) != claude|ver {
		t.Errorf("the whitelisted binary must write its own reservations: %v", r.Writers)
	}
}

func TestBuildGlobReservations_BunTmpdir(t *testing.T) {
	home := t.TempDir()
	bunDir := filepath.Join(home, install.BunTmpRelDir)
	mkdirs(t, filepath.Join(home, ".config/opencode"), bunDir)
	bin := filepath.Join(home, ".local/bin/opencode")
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".config/opencode"), Binaries: []daemonconfig.BinaryRule{{Path: bin}}},
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	bit := bitOf(t, r, install.BunReservedName)
	if r.Roots[bunDir]&bit == 0 {
		t.Errorf("%s not reserved below %s: %v", install.BunReservedName, bunDir, r.Roots)
	}
	if r.Writers[bin]&bit == 0 {
		t.Errorf("the configured Bun binary must be a writer of its extraction: %v", r.Writers)
	}
	if r.Writers["/usr/bin/ssh"]&bit != 0 {
		t.Errorf("a non-Bun binary must not write the Bun extraction: %v", r.Writers)
	}
}

func TestBuildGlobReservations_BunTmpdirMissingDirReservesNoRoot(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, ".config/opencode")) // the tmp dir is absent
	bin := filepath.Join(home, ".local/bin/opencode")
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".config/opencode"), Binaries: []daemonconfig.BinaryRule{{Path: bin}}},
	}}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	if _, ok := r.Roots[filepath.Join(home, install.BunTmpRelDir)]; ok {
		t.Errorf("a missing Bun tmp dir must reserve no root: %v", r.Roots)
	}
}

func TestBuildGlobReservations_BunTmpdirUnconfiguredReservesNothing(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, install.BunTmpRelDir)) // dir exists but opencode is not configured
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	if slices.Contains(r.Patterns, install.BunReservedName) {
		t.Errorf("no configured Bun app: .bun-* must not be reserved: %v", r.Patterns)
	}
}

func steamLibDirConfig(home string) (*daemonconfig.Config, string) {
	steam := filepath.Join(home, ".local/share/Steam")
	client := filepath.Join(steam, "ubuntu12_32/steam")
	return &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(steam, "config"), Binaries: []daemonconfig.BinaryRule{{Path: client}}},
		{Path: filepath.Join(home, ".ssh"), Binaries: []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}}},
	}}, client
}

// The catalog refresh adopts a wildcard-free lib dir once it exists, so a missing one must be
// reserved as well as a present one (a non-writer creating compatibilitytools.d with a library).
func TestBuildGlobReservations_SteamFixedLibDirs(t *testing.T) {
	home := t.TempDir()
	steam := filepath.Join(home, ".local/share/Steam")
	mkdirs(t, filepath.Join(steam, "config"), filepath.Join(steam, "ubuntu12_64"))
	cfg, client := steamLibDirConfig(home)
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	bit := bitOf(t, r, "ubuntu12_32")
	for _, c := range [][2]string{
		{home, ".local"}, {filepath.Join(home, ".local"), "share"}, {filepath.Join(home, ".local/share"), "Steam"},
		{steam, "ubuntu12_64"}, {steam, "compatibilitytools.d"}, {steam, "linux64"},
	} {
		if !hasChild(r, c[0], c[1], bit) {
			t.Errorf("%s/%s must be reserved for Steam's writers: %+v", c[0], c[1], r.Children)
		}
	}
	for root, bits := range r.Roots {
		if bits&bit != 0 {
			t.Errorf("lib dirs are exact children, but %s reserves their bit at any depth", root)
		}
	}
	if r.Writers[client]&bit == 0 {
		t.Errorf("the Steam client must be a writer of its lib dirs: %v", r.Writers)
	}
	if r.Writers["/usr/bin/ssh"]&bit != 0 {
		t.Errorf("another entry's binary must not create Steam's lib dirs: %v", r.Writers)
	}
}

func TestBuildGlobReservations_FixedLibDirStopsAtFirstMissing(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, ".local"))
	cfg, _ := steamLibDirConfig(home)
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	bit := bitOf(t, r, "ubuntu12_32")
	if !hasChild(r, filepath.Join(home, ".local"), "share", bit) {
		t.Errorf("the first missing component must be reserved in its parent: %+v", r.Children)
	}
	for _, c := range r.Children {
		if c.Bits&bit != 0 && c.Parent != home && c.Parent != filepath.Join(home, ".local") {
			t.Errorf("reserved %s/%s below a missing dir", c.Parent, c.Name)
		}
	}
}

// SetGlobReservations refuses more than 64 keys: every catalog entry configured at once must fit.
func TestBuildGlobReservations_WholeCatalogFits(t *testing.T) {
	home := t.TempDir()
	cfg := &daemonconfig.Config{}
	for i := range install.Catalog {
		e := &install.Catalog[i]
		if e.IsSystem() {
			continue
		}
		for j, p := range e.PathsFor(home, "u") {
			cfg.Resources = append(cfg.Resources, daemonconfig.Resource{
				Path:     p,
				Binaries: []daemonconfig.BinaryRule{{Path: filepath.Join("/opt", e.Name, string(rune('a'+j)))}},
			})
		}
	}
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	if len(r.Patterns) > 64 {
		t.Errorf("%d reserved keys, SetGlobReservations accepts 64: %q", len(r.Patterns), r.Patterns)
	}
	t.Logf("%d of 64 reserved keys used", len(r.Patterns))
}

func TestBuildGlobReservations_BunRootsListed(t *testing.T) {
	withDir, without := t.TempDir(), t.TempDir()
	mkdirs(t, filepath.Join(withDir, install.BunTmpRelDir))
	var cfg daemonconfig.Config
	for _, h := range []string{withDir, without} {
		cfg.Resources = append(cfg.Resources, daemonconfig.Resource{Path: filepath.Join(h, ".config/opencode"),
			Binaries: []daemonconfig.BinaryRule{{Path: filepath.Join(h, ".local/bin/opencode")}}})
	}
	_, bun := buildGlobReservations(&cfg, []install.User{{Name: "a", Home: withDir}, {Name: "b", Home: without}})
	if len(bun) != 1 || bun[0].Home != withDir {
		t.Errorf("only the user with a reserved Bun tmp dir must be vetted: %+v", bun)
	}
}

// Only Proton* is reserved at the root: inside a match lacking files/lib (the EasyAntiCheat
// runtime) the tail must be reserved too, or a non-writer creates it and the refresh adopts it.
func TestBuildGlobReservations_WildcardLibDirTail(t *testing.T) {
	home := t.TempDir()
	steam := filepath.Join(home, ".local/share/Steam")
	common := filepath.Join(steam, "steamapps/common")
	eac := filepath.Join(common, "Proton EasyAntiCheat Runtime")
	proton := filepath.Join(common, "Proton 9.0")
	game := filepath.Join(common, "SomeGame")
	mkdirs(t, filepath.Join(steam, "config"), eac, filepath.Join(proton, "files/lib"), game)
	cfg, client := steamLibDirConfig(home)
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	bit := bitOf(t, r, "ubuntu12_32")
	for _, c := range [][2]string{
		{eac, "files"}, {proton, "files"}, {filepath.Join(proton, "files"), "lib"},
	} {
		if !hasChild(r, c[0], c[1], bit) {
			t.Errorf("%s/%s must be reserved for Steam's writers: %+v", c[0], c[1], r.Children)
		}
	}
	for _, c := range r.Children {
		if c.Parent == filepath.Join(eac, "files") || c.Parent == game {
			t.Errorf("reserved %s/%s: below a missing dir or outside every Proton* match", c.Parent, c.Name)
		}
	}
	if r.Writers[client]&bit == 0 {
		t.Errorf("the Steam client must be a writer of the Proton tails: %v", r.Writers)
	}
}

// Discord's updater inputs and version dirs are reserved for its binaries (issue #79): no other
// process can create, replace or write settings.json/installer.db, or create an app-* dir.
func TestBuildGlobReservations_DiscordUpdaterInputs(t *testing.T) {
	home := t.TempDir()
	discord := filepath.Join(home, ".config/discord")
	mkdirs(t, discord)
	cfg, bin := discordLibConfig(home)
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})

	for _, name := range []string{"settings.json", "installer.db", "app-*"} {
		bit := bitOf(t, r, name)
		if r.Roots[discord]&bit == 0 {
			t.Errorf("%s not reserved below %s: %v", name, discord, r.Roots)
		}
		if r.Writers[bin]&bit == 0 {
			t.Errorf("Discord must write its %s", name)
		}
		if r.Writers["/usr/bin/ssh"]&bit != 0 {
			t.Errorf("another entry's binary must not write Discord's %s", name)
		}
	}
}

// Falling back to ~/.config would make every app's settings.json Discord's.
func TestBuildGlobReservations_MissingInputRootNotWidened(t *testing.T) {
	home := t.TempDir()
	mkdirs(t, filepath.Join(home, ".config"))
	cfg, _ := discordLibConfig(home)
	r, _ := buildGlobReservations(cfg, []install.User{{Name: "u", Home: home}})
	bit := bitOf(t, r, "settings.json")
	for root, bits := range r.Roots {
		if bits&bit != 0 {
			t.Errorf("settings.json reserved below %s although %s/.config/discord is missing", root, home)
		}
	}
}
