package guard_test

import (
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/install"
)

func encodeCatalog(t *testing.T, name string) (map[string]uint64, guard.GuardTrustArgvSet) {
	t.Helper()
	words, set, err := guard.EncodeLaunchRules(install.CatalogLaunchRules(name))
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return words, set
}

func TestCatalogLaunchRulesFit(t *testing.T) {
	var k guard.GuardTrustArgvKey
	for i := range install.Catalog {
		for name := range install.Catalog[i].LaunchRules {
			words, _ := encodeCatalog(t, name)
			for w := range words {
				if len(w) > daemonconfig.MaxLaunchWordLen || len(w) >= len(k.Name) {
					t.Errorf("%s: word %q too long", name, w)
				}
			}
		}
	}
}

func TestSSHConfigFileRefusedInEveryCluster(t *testing.T) {
	_, set := encodeCatalog(t, "ssh")
	for word, want := range map[string]bool{
		"-F": true, "-F/tmp/c": true, "-4F/tmp/c": true, "-vvqF": true, "-TF": true,
		"-oF": false, "-oForwardAgent=yes": false, "-L8080:F": false, "-iF": false, "--F": false,
		"-": false, "host": false, "-v": false,
	} {
		if got := guard.ShortRefused(set, word); got != want {
			t.Errorf("ssh %q: refused = %v, want %v", word, got, want)
		}
	}
}

func TestSSHServerBinariesAreRootOnly(t *testing.T) {
	for _, n := range []string{"sshd", "sshd-session", "sshd-auth"} {
		_, set := encodeCatalog(t, n)
		if set.Flags&guard.ArgvRootOnly == 0 || set.Lens != 0 {
			t.Errorf("%s: flags %#x lens %#x, want root-only and no argv rule", n, set.Flags, set.Lens)
		}
	}
}

// gpg runs any unambiguous abbreviation of a long option (GnuPG 2.4: `--export-secret-ke`, `--enarmo`).
func TestGPGSecretExportsRefusedInEverySpelling(t *testing.T) {
	words, _ := encodeCatalog(t, "gpg")
	for _, w := range []string{"--export-se", "--export-secret-", "--export-secret-k", "--export-secret-keys",
		"--export-secret-su", "--export-secret-subkeys", "--export-secret-ss", "--export-secret-ssh-key",
		"--enar", "--enarmo", "--enarmor", "--enarmou", "--enarmour", "--sto", "--store"} {
		if words[w] == 0 {
			t.Errorf("gpg %s is not refused", w)
		}
	}
	for _, w := range []string{"--export", "--export-s", "--export-ssh-key", "--encrypt", "--ena", "--enable-dsa2",
		"--st", "--status-fd", "--decrypt", "-e", "-d", "--sign"} {
		if words[w] != 0 {
			t.Errorf("gpg %s is refused", w)
		}
	}
}
