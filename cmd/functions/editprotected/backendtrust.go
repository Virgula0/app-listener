package editprotected

import (
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Virgula0/app-listener/internal/binledger"
	"github.com/Virgula0/app-listener/internal/bintrust"
	"github.com/Virgula0/app-listener/internal/editbackend"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// maxTrustItems bounds one trust request (the ledger lists one row per refused whitelist line).
const maxTrustItems = 512

// trustList is `trust list`. It opens only paths the ledger names (whitelist lines the daemon
// refused), never the GUI's, so it runs before AUTH like the CLI, which needs only root.
func (b *backend) trustList(*editbackend.Request) (*editbackend.Message, error) {
	l, err := bintrust.OpenLedger()
	if err != nil {
		return nil, err
	}
	defer l.Close()
	rows, err := bintrust.Pending(l)
	if err != nil {
		return nil, err
	}
	return &editbackend.Message{Trust: trustRows(rows)}, nil
}

// trustConfirm is `trust` after the operator's choice: each line is recorded only while the
// daemon's refusal still names the content the GUI showed and the file still holds it.
func (b *backend) trustConfirm(req *editbackend.Request) (*editbackend.Message, error) {
	return b.trustRecord(req, bintrust.ConfirmPending)
}

// trustInspect describes lines typed by hand (`trust add`). Hashing a caller-named path with this
// binary's root READ baseline on guarded trees needs AUTH, or the daemon stopped (as offline_scan).
func (b *backend) trustInspect(req *editbackend.Request) (*editbackend.Message, error) {
	if b.password == "" {
		if err := b.trustOffline(); err != nil {
			return nil, err
		}
	}
	if len(req.Binaries) == 0 || len(req.Binaries) > maxTrustItems {
		return nil, fmt.Errorf("give between 1 and %d paths", maxTrustItems)
	}
	for _, p := range req.Binaries {
		if err := checkTrustLine(p); err != nil {
			return nil, err
		}
	}
	if b.session != nil {
		b.session.End() // AUTH was the proof; no grant follows
		b.session = nil
	}
	b.inspected, b.state = slices.Clone(req.Binaries), stTrusting
	return &editbackend.Message{Trust: trustRows(bintrust.Given(req.Binaries, "(added by hand)"))}, nil
}

// trustAdd records lines of the last trust_inspect while they still hold the content shown.
func (b *backend) trustAdd(req *editbackend.Request) (*editbackend.Message, error) {
	if b.password == "" {
		if err := b.trustOffline(); err != nil {
			return nil, err
		}
	}
	for _, it := range req.Trust {
		if !slices.Contains(b.inspected, it.Line) {
			return nil, fmt.Errorf("%s was not inspected in this session", it.Line)
		}
	}
	return b.trustRecord(req, bintrust.ConfirmGiven)
}

// trustOffline: without AUTH, only while no daemon guards anything.
func (b *backend) trustOffline() error {
	if b.live {
		return errors.New("the daemon is running with a password: authenticate first")
	}
	return protected.RequireDaemonStopped()
}

func (b *backend) trustRecord(req *editbackend.Request,
	confirm func(*binledger.Ledger, string, [32]byte) error) (*editbackend.Message, error) {
	if len(req.Trust) == 0 || len(req.Trust) > maxTrustItems {
		return nil, fmt.Errorf("give between 1 and %d binaries", maxTrustItems)
	}
	sums := make([][32]byte, len(req.Trust))
	for i, it := range req.Trust {
		if err := checkTrustLine(it.Line); err != nil {
			return nil, err
		}
		raw, err := hex.DecodeString(it.SHA256)
		if err != nil || len(raw) != len(sums[i]) {
			return nil, fmt.Errorf("%s: malformed sha256", it.Line)
		}
		copy(sums[i][:], raw)
	}
	l, err := bintrust.OpenLedger()
	if err != nil {
		return nil, err
	}
	defer l.Close()
	done := &editbackend.TrustDone{}
	for i, it := range req.Trust {
		if err := confirm(l, it.Line, sums[i]); err != nil {
			if done.Skipped == nil {
				done.Skipped = map[string]string{}
			}
			done.Skipped[it.Line] = err.Error()
			continue
		}
		done.Confirmed = append(done.Confirmed, it.Line)
	}
	if len(done.Confirmed) > 0 && !req.NoReload {
		if err := systemd.ReloadDaemonIfActive(); err != nil {
			done.ReloadErr = fmt.Sprintf("%v: run `systemctl reload %s`", err, systemd.DaemonServiceName)
		} else {
			done.Reloaded = true
		}
	}
	return &editbackend.Message{TrustDone: done}, nil
}

func checkTrustLine(p string) error {
	switch {
	case !utf8.ValidString(p), strings.ContainsRune(p, 0), len(p) > 4096:
		return errors.New("malformed path")
	case !filepath.IsAbs(p):
		return fmt.Errorf("%s is not an absolute path", p)
	}
	return nil
}

func trustRows(rows []bintrust.Row) []editbackend.TrustRow {
	out := make([]editbackend.TrustRow, 0, len(rows))
	for i := range rows {
		r := &rows[i]
		w := editbackend.TrustRow{Line: r.Line, Resolved: r.Resolved, Resources: r.Resources, SeenAt: r.SeenAt,
			SHA256: hex.EncodeToString(r.SHA256[:]), Skip: r.Skip()}
		if st := r.Stat; st != nil {
			w.HasStat, w.UID, w.Mode = true, st.Uid, st.Mode&0o7777
			w.ModTime = time.Unix(st.Mtim.Sec, 0)
		}
		out = append(out, w)
	}
	return out
}
