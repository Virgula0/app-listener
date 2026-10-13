package guard

import (
	"testing"
	"time"
)

func TestUnchangedPerShared(t *testing.T) {
	now := time.Now()
	deployed := GuardInodeKey{Dev: 1, Ino: 10}
	g := &Guard{deployed: map[string]GuardInodeKey{"/opt/app/bin": deployed}}
	reset := func() { sweepShared = sweepMemo{ids: map[string]sweepID{}} }
	defer reset()

	for _, tc := range []struct {
		name      string
		canonical string
		key       GuardInodeKey
		at        time.Time
		want      bool
	}{
		{"fresh resolution of the deployed inode", "/opt/app/bin", deployed, now, true},
		{"replaced inode: must re-check itself", "/opt/app/bin", GuardInodeKey{Dev: 1, Ino: 11}, now, false},
		{"expired resolution", "/opt/app/bin", deployed, now.Add(-sweepShareTTL - time.Millisecond), false},
		{"resolution from the future (clock step)", "/opt/app/bin", deployed, now.Add(time.Second), false},
		{"canonical path not deployed here", "/opt/other/bin", deployed, now, false},
	} {
		reset()
		sweepShared.record("/usr/bin/app", tc.canonical, tc.key, tc.at)
		if got := g.unchangedPerShared("/usr/bin/app", now); got != tc.want {
			t.Fatalf("%s: unchangedPerShared = %v, want %v", tc.name, got, tc.want)
		}
	}

	reset()
	if g.unchangedPerShared("/usr/bin/app", now) {
		t.Fatal("no shared resolution must never skip")
	}
}

// A multicall applet's deployed key carries its applet tag; the shared resolution is the real
// inode, compared like resyncOne does (Real()).
func TestUnchangedPerSharedComparesRealKey(t *testing.T) {
	defer func() { sweepShared = sweepMemo{ids: map[string]sweepID{}} }()
	real := GuardInodeKey{Dev: 1, Ino: 10}
	tagged := GuardInodeKey{Dev: real.Dev | 5<<32, Ino: real.Ino}
	if tagged.Real() != real {
		t.Skip("applet tag encoding changed; update this fixture")
	}
	g := &Guard{deployed: map[string]GuardInodeKey{"/usr/bin/coreutils": tagged}}
	now := time.Now()
	sweepShared.record("/usr/bin/ls", "/usr/bin/coreutils", real, now)
	if !g.unchangedPerShared("/usr/bin/ls", now) {
		t.Fatal("tagged deployed key must match its real inode")
	}
}
