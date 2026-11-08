package guard

import "testing"

func TestUnixRoleByName(t *testing.T) {
	for path, want := range map[string]uint8{
		"/usr/bin/ssh-add": unixConfined, "/usr/bin/ssh-agent": unixAgent, "/usr/bin/gpg-agent": unixAgent,
		"/usr/lib/gcr-ssh-agent": unixAgent, "/opt/1Password/1password": unixAgent,
		"/usr/bin/ssh": 0, "/usr/lib/systemd/systemd": 0, "/usr/bin/gpg-connect-agent": 0,
	} {
		if got := unixRole(path); got != want {
			t.Errorf("%s: role %d, want %d", path, got, want)
		}
	}
}
