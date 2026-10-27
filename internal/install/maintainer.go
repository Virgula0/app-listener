package install

import (
	"strings"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// Maintainer is one user's copy of a catalog entry's Maintainer script, with the configured
// read-only lib_dirs below the script's directory that its runs may change.
type Maintainer struct {
	Entry   *CandidateDir
	User    User
	Script  string
	LibDirs []*daemonconfig.Resource
}

// Maintainers lists, per user, every catalog Maintainer with at least one lib_dir of its entry
// configured in cfg below the script's directory.
func Maintainers(cfg *daemonconfig.Config, users []User) []Maintainer {
	var out []Maintainer
	for _, u := range users {
		for i := range Catalog {
			e := &Catalog[i]
			script := e.ExpandMaintainer(u.Name, u.Home)
			if script == "" {
				continue
			}
			root := script[:strings.LastIndex(script, "/")+1]
			m := Maintainer{Entry: e, User: u, Script: script}
			for j := range cfg.Resources {
				r := &cfg.Resources[j]
				if r.ReadOnly && !r.PathPending && strings.HasPrefix(r.Path, root) && e.OwnsPath(u.Name, u.Home, r.Path) {
					m.LibDirs = append(m.LibDirs, r)
				}
			}
			if len(m.LibDirs) > 0 {
				out = append(out, m)
			}
		}
	}
	return out
}
