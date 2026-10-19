package install

import "testing"

func TestNamedGeneralTool(t *testing.T) {
	for _, p := range []string{
		"/usr/bin/python3.12", "/usr/bin/gawk-5.4.0", "/bin/nc.openbsd", "/usr/bin/vim.basic", "/usr/bin/vim.gtk3",
		"/usr/bin/tclsh8.6", "/usr/bin/sort", "/usr/bin/sqlite3", "/usr/bin/batcat", "/usr/bin/zcat",
	} {
		if !namedGeneralTool(p) {
			t.Errorf("%s must be a general tool", p)
		}
	}
	for _, p := range []string{
		"/usr/bin/ssh", "/usr/lib/ssh/ssh-keysign", "/usr/bin/gpg-agent", "/usr/bin/git-lfs-helper-x",
		"/usr/bin/lua-language-server", "/usr/bin/code-cli", "/opt/discord/Discord", "/usr/bin/claude",
	} {
		if namedGeneralTool(p) {
			t.Errorf("%s is not a general tool", p)
		}
	}
}
