package install

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// TestUncoveredCandidates: a catalog directory already in the config is not
// proposed; one that exists on disk but is absent from the config is.
func TestUncoveredCandidates(t *testing.T) {
	home := t.TempDir()
	for _, rel := range []string{".ssh", ".gnupg"} {
		if err := os.MkdirAll(filepath.Join(home, rel), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	users := []inst.User{{Name: "tester", Home: home}}
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".ssh"), NeedEncryption: true},
	}}

	fresh := uncoveredCandidates(cfg, users)

	got := map[string]bool{}
	for i := range fresh {
		got[fresh[i].Path] = true
	}
	if got[filepath.Join(home, ".ssh")] {
		t.Error(".ssh is already in the config — must not be proposed")
	}
	if !got[filepath.Join(home, ".gnupg")] {
		t.Errorf(".gnupg exists and is not configured — must be proposed (got %v)", got)
	}
}

// TestDiffMergePreservesExistingAndParses: appending the generated sections
// to the installed text keeps the original as a verbatim prefix and the
// result parses with both resources.
func TestDiffMergePreservesExistingAndParses(t *testing.T) {
	dirA := t.TempDir()
	dirB := t.TempDir()
	old := "[watch " + dirA + "]\n\n/usr/bin/x\n\nneed_encryption: false\n"

	picked := []inst.Candidate{{Entry: inst.CandidateDir{Name: "n"}, Path: dirB}}
	added := inst.GenerateSections(sectionsFromCandidates(picked))
	full := old + added + "\n"

	if full[:len(old)] != old {
		t.Fatalf("existing config not preserved verbatim as prefix:\n%q", full)
	}
	cfg, err := validateConfigText(full)
	if err != nil {
		t.Fatalf("merged config does not parse: %v\n%s", err, full)
	}
	if len(cfg.Resources) != 2 {
		t.Errorf("merged config has %d resources, want 2", len(cfg.Resources))
	}
}

func TestAddedResourcesSkipsAlreadyGuardedSSH(t *testing.T) {
	prev := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: "/home/u/.ssh", NeedEncryption: true},
	}}
	merged := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: "/home/u/.ssh", NeedEncryption: true},
		{Path: "/home/u/.config/gh", NeedEncryption: true},
	}}

	got := addedResources(prev, merged)

	if len(got.Resources) != 1 || got.Resources[0].Path != "/home/u/.config/gh" {
		t.Fatalf("addedResources = %+v, want only /home/u/.config/gh", got.Resources)
	}
}
