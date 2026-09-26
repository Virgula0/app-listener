package daemon

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/guard"
)

func TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.StandardLogger().Out)

	const discord = "/h/.config/discord/app/Discord"
	uid := errors.New("not root-owned (uid 1000)")
	rejected := map[string]*libRejection{
		"/h/steam/steamrt64/libcef.so":         {why: uid, bins: []string{"/h/steam/steamrt64/steamwebhelper"}},
		"/h/.config/discord/app/libffmpeg.so":  {why: uid, bins: []string{discord}},
		"/h/.local/share/evil/libplanted.so.1": {why: uid, bins: []string{"/usr/bin/a", "/usr/bin/b"}},
		"/h/steam/steamrt64/libother.so":       {why: uid, bins: []string{"/usr/bin/ssh"}},
	}
	r := guard.GlobReservations{
		Patterns: []string{"*.so"},
		Roots:    map[string]uint64{"/h/.config/discord": 1},
		Writers:  map[string]uint64{discord: 1},
	}
	dirs := []guard.TrustedDir{{Path: "/h/steam/steamrt64", Loaders: []string{"/h/steam/steamrt64/steamwebhelper"}}}
	warnUntrustedLibs(rejected, dirs, r)

	out := buf.String()
	if strings.Contains(out, "libcef.so") || strings.Contains(out, "libffmpeg.so") {
		t.Errorf("guarded-tree and reserved libraries are trusted by the kernel, no warning: %s", out)
	}
	if strings.Count(out, "libplanted.so.1") != 1 || !strings.Contains(out, "2 whitelisted binary(ies)") {
		t.Errorf("an uncovered library must be reported once, with its loaders counted: %s", out)
	}
	if !strings.Contains(out, "libother.so") {
		t.Errorf("a lib_dir library is trusted only for the dir's writers, ssh must be warned about: %s", out)
	}
}
