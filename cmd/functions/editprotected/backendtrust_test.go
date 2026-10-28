package editprotected

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/editbackend"
)

func TestBackendTrustAddNeedsAnInspection(t *testing.T) {
	b, _ := newTestBackend(t)
	item := []editbackend.TrustItem{{Line: "/usr/bin/true", SHA256: strings.Repeat("0", 64)}}
	if m := call(b, editbackend.Request{Op: editbackend.OpTrustAdd, Trust: item}); m.Err == "" {
		t.Fatal("trust_add before trust_inspect: accepted")
	}
	b.state, b.password, b.inspected = stTrusting, "pw", []string{"/usr/bin/false"}
	if m := call(b, editbackend.Request{Op: editbackend.OpTrustAdd, Trust: item}); !strings.Contains(m.Err, "not inspected") {
		t.Fatalf("trust_add of a line never inspected: err=%q", m.Err)
	}
	b.inspected = []string{"/usr/bin/true"}
	item[0].SHA256 = "zz"
	if m := call(b, editbackend.Request{Op: editbackend.OpTrustAdd, Trust: item}); !strings.Contains(m.Err, "malformed sha256") {
		t.Fatalf("malformed sha256: err=%q", m.Err)
	}
}

// A caller-named path is hashed by this binary, which reads guarded trees: never before AUTH while
// the daemon runs.
func TestBackendTrustInspectNeedsAuthWhileLive(t *testing.T) {
	b, _ := newTestBackend(t)
	b.live = true
	m := call(b, editbackend.Request{Op: editbackend.OpTrustInspect, Binaries: []string{"/usr/bin/true"}})
	if !strings.Contains(m.Err, "authenticate first") || b.state != stNew {
		t.Fatalf("inspect before auth: err=%q state=%v", m.Err, b.state)
	}
}

func TestBackendTrustInspectAfterAuth(t *testing.T) {
	b, _ := newTestBackend(t)
	b.live, b.password, b.state = true, "pw", stAuthed
	bin := filepath.Join(t.TempDir(), "tool")
	mustWrite(t, bin, "#!/bin/sh\n")
	if m := call(b, editbackend.Request{Op: editbackend.OpTrustInspect, Binaries: []string{"rel/tool"}}); m.Err == "" {
		t.Fatal("relative path accepted")
	}
	m := call(b, editbackend.Request{Op: editbackend.OpTrustInspect, Binaries: []string{bin}})
	if m.Err != "" || len(m.Trust) != 1 || b.state != stTrusting {
		t.Fatalf("inspect: %+v state=%v", m, b.state)
	}
	sum := sha256.Sum256([]byte("#!/bin/sh\n"))
	if r := m.Trust[0]; r.Line != bin || r.SHA256 != hex.EncodeToString(sum[:]) || r.Skip != "" || !r.HasStat {
		t.Fatalf("row = %+v", r)
	}
	for _, op := range []string{editbackend.OpTrustList, editbackend.OpTrustConfirm, editbackend.OpSelect} {
		if m := call(b, editbackend.Request{Op: op}); m.Err == "" {
			t.Errorf("%s after trust_inspect: accepted", op)
		}
	}
}
