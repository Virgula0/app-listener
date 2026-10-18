package editprotected

import (
	"errors"
	"fmt"

	"github.com/Virgula0/app-listener/internal/editbackend"
)

// recordStart is runRecord after the password: the paths arrive only now, after AUTH, and are
// validated against the daemon's configuration before anything watches them.
func (b *backend) recordStart(req *editbackend.Request) (*editbackend.Message, error) {
	if b.session != nil {
		b.session.End()
		b.session = nil
	}
	flow, err := newRecordFlow(b.password, req.Resources)
	if err != nil {
		return nil, err
	}
	if err := flow.start(); err != nil {
		return nil, err
	}
	b.record, b.state = flow, stRecording
	paths := make([]string, len(flow.targets))
	for i, t := range flow.targets {
		paths[i] = t.Path
	}
	return &editbackend.Message{Resources: paths}, nil
}

func (b *backend) recordStatus(*editbackend.Request) (*editbackend.Message, error) {
	p := b.record.rec.Progress()
	return &editbackend.Message{RecordProgress: &p}, nil
}

func (b *backend) recordStop(*editbackend.Request) (*editbackend.Message, error) {
	rec := b.record.stop()
	b.state = stRecorded
	return &editbackend.Message{Recording: rec, EncryptBlockers: b.record.encryptBlockers()}, nil
}

func (b *backend) recordPlan(req *editbackend.Request) (*editbackend.Message, error) {
	if req.Selection == nil {
		return nil, errors.New("no selection")
	}
	base, text, err := b.record.preview(req.Selection)
	if err != nil {
		return nil, err
	}
	return &editbackend.Message{Base: base, Data: []byte(text)}, nil
}

// recordApply applies the last record_plan, reporting the encryption copy as progress.
func (b *backend) recordApply(*editbackend.Request) (*editbackend.Message, error) {
	out, err := b.record.apply(func(step string, fraction float64) {
		b.send(&editbackend.Message{Kind: editbackend.KindProgress,
			Progress: &editbackend.Progress{Step: step, Fraction: fraction}})
	})
	if err != nil {
		return nil, err
	}
	b.state = stApplied
	return &editbackend.Message{Notes: out.Notes, Backups: out.Backups}, nil
}

func (b *backend) recordDeleteBackup(req *editbackend.Request) (*editbackend.Message, error) {
	if err := b.record.deleteBackup(req.Path); err != nil {
		return nil, fmt.Errorf("not removed: %w", err)
	}
	return nil, nil
}
