package editprotected

import (
	"errors"
	"fmt"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/recorder"
)

// socketStore is recorder.Store over the control socket: one authenticated connection per call,
// as --edit-config makes them.
type socketStore struct{ password string }

func (s socketStore) Fetch() ([]byte, error) { return fetchConfig(s.password) }

func (s socketStore) Put(base, next []byte) error {
	err := putConfig(s.password, base, next)
	if errors.Is(err, errConfigChanged) {
		return fmt.Errorf("%w: %w", recorder.ErrStale, err)
	}
	return err
}

func (s socketStore) Guarded() ([]string, error) {
	session, err := dialLiveSession()
	if err != nil {
		return nil, err
	}
	defer session.End()
	return session.Authenticate(s.password)
}

// recordFlow is `--edit-config --record` shared by the TUI and the GUI backend. Nothing touches a
// caller-named path before newRecordFlow has authenticated: this binary keeps a root READ baseline
// on every guarded tree.
type recordFlow struct {
	store     socketStore
	baseCfg   *daemonconfig.Config
	targets   []recorder.Target
	rec       *recorder.Recorder
	recording *recorder.Recording
	plan      *recorder.Plan
	outcome   *recorder.Outcome
}

// newRecordFlow authenticates (fetching daemon.conf), then validates paths against it.
func newRecordFlow(password string, paths []string) (*recordFlow, error) {
	store := socketStore{password: password}
	base, err := store.Fetch()
	if err != nil {
		return nil, err
	}
	cfg, err := daemonconfig.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("the running daemon.conf does not parse: %w", err)
	}
	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	homes := make([]string, len(users))
	for i, u := range users {
		homes[i] = u.Home
	}
	targets, err := recorder.ValidateTargets(paths, cfg, homes)
	if err != nil {
		return nil, err
	}
	return &recordFlow{store: store, baseCfg: cfg, targets: targets}, nil
}

func (f *recordFlow) start() error {
	rec, err := recorder.Start(f.targets)
	if err != nil {
		return err
	}
	f.rec = rec
	return nil
}

// stop ends the recording and classifies it.
func (f *recordFlow) stop() *recorder.Recording {
	if f.recording == nil && f.rec != nil {
		f.recording = recorder.Classify(f.rec.Stop(), f.baseCfg)
	}
	return f.recording
}

// encryptBlockers says, per target, why it can't be encrypted live ("" = it can).
func (f *recordFlow) encryptBlockers() map[string]string {
	vault := fscrypt.New()
	out := map[string]string{}
	for _, t := range f.targets {
		if t.Existing != nil {
			out[t.Path] = "it is already in daemon.conf: change its need_encryption with --edit-config"
			continue
		}
		out[t.Path] = encryptBlocker(vault, t.Path)
	}
	return out
}

func encryptBlocker(vault *fscrypt.Vault, path string) string {
	enc, err := vault.IsEncrypted(path)
	switch {
	case err != nil:
		return err.Error()
	case enc:
		return "it already carries an encryption policy"
	}
	prereqs, err := vault.FilesystemPrereqs(path)
	if err != nil {
		return err.Error()
	}
	if len(prereqs) > 0 {
		return "its filesystem isn't prepared for fscrypt (" + prereqs[0].Reason + "): run `" +
			prereqs[0].Command() + "` first"
	}
	return ""
}

// preview validates sel and renders it over the daemon's current daemon.conf.
func (f *recordFlow) preview(sel *recorder.Selection) (base []byte, text string, err error) {
	if f.recording == nil {
		return nil, "", errors.New("no recording to review")
	}
	plan, err := f.recording.Plan(sel) // first: only the recording's own targets reach encryptBlocker
	if err != nil {
		return nil, "", err
	}
	vault := fscrypt.New()
	for _, p := range plan.Encrypt {
		if why := encryptBlocker(vault, p); why != "" {
			return nil, "", fmt.Errorf("%s can't be encrypted: %s", p, why)
		}
	}
	if base, err = f.store.Fetch(); err != nil {
		return nil, "", err
	}
	if text, err = plan.Render(base); err != nil {
		return nil, "", err
	}
	f.plan = plan
	return base, text, nil
}

// apply performs the last previewed plan.
func (f *recordFlow) apply(progress recorder.ProgressFunc) (*recorder.Outcome, error) {
	if f.plan == nil {
		return nil, errors.New("review the changes first")
	}
	if f.outcome != nil {
		return nil, errors.New("this recording was already applied")
	}
	out, err := recorder.Apply(f.store, f.plan, progress)
	if err != nil {
		return nil, err
	}
	f.outcome = out
	return out, nil
}

func (f *recordFlow) deleteBackup(backup string) error {
	if f.outcome == nil {
		return errors.New("nothing was applied")
	}
	return recorder.DeleteBackup(f.outcome, backup)
}

// close ends a recording still running (its probes detach with it).
func (f *recordFlow) close() {
	if f.rec != nil {
		f.rec.Stop()
	}
}
