package editprotected

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/x/term"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/internal/constants"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/editbackend"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/guard"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/vaultfs"
)

// backendState is where one backend process is in its single workflow. Each op is accepted only in
// the states listed in backendOps: a request can't skip AUTH, reach files before SELECT, or act
// after the daemon ended the session. This binary keeps a root-only OPEN/READ/STAT baseline on
// every guarded tree (daemonSelfBaselineEvents), so nothing that opens a caller-named path may run
// before AUTH (or, offline, with the daemon running).
type backendState int

const (
	stNew        backendState = iota // after hello: auth (live), offline scan/open, password ops
	stAuthed                         // AUTH accepted: select, forward or config
	stEditing                        // a vault is open (live SELECT or offline unlock)
	stForwarding                     // a FORWARD grant is held
	stConfig                         // daemon.conf fetched
	stRecording                      // a recording runs (record_start)
	stRecorded                       // stopped: review, plan, apply
	stApplied                        // applied: only its backups may be deleted
	stDone                           // only end
)

type backendOp struct {
	states []backendState
	run    func(*backend, *editbackend.Request) (*editbackend.Message, error)
}

var backendOps = map[string]backendOp{
	editbackend.OpAuth:          {[]backendState{stNew}, (*backend).auth},
	editbackend.OpOfflineScan:   {[]backendState{stNew}, (*backend).offlineScan},
	editbackend.OpOfflineOpen:   {[]backendState{stNew}, (*backend).offlineOpen},
	editbackend.OpPasswordState: {[]backendState{stNew}, (*backend).passwordState},
	editbackend.OpSetPassword:   {[]backendState{stNew}, (*backend).setPassword},
	editbackend.OpForwardCheck:  {[]backendState{stAuthed}, (*backend).forwardCheck},
	editbackend.OpSelect:        {[]backendState{stAuthed}, (*backend).selectResource},
	editbackend.OpForward:       {[]backendState{stAuthed}, (*backend).forward},
	editbackend.OpConfigGet:     {[]backendState{stAuthed, stConfig}, (*backend).configGet},
	editbackend.OpConfigPut:     {[]backendState{stConfig}, (*backend).configPut},
	editbackend.OpList:          {[]backendState{stEditing}, (*backend).list},
	editbackend.OpRead:          {[]backendState{stEditing}, (*backend).read},
	editbackend.OpWrite:         {[]backendState{stEditing}, (*backend).write},
	editbackend.OpCreate:        {[]backendState{stEditing}, (*backend).create},
	editbackend.OpChmod:         {[]backendState{stEditing}, (*backend).chmod},
	editbackend.OpChown:         {[]backendState{stEditing}, (*backend).chown},
	editbackend.OpRemove:        {[]backendState{stEditing}, (*backend).remove},
	editbackend.OpUsers:         {[]backendState{stEditing}, (*backend).users},
	editbackend.OpPing:          {[]backendState{stEditing, stForwarding}, (*backend).ping},

	editbackend.OpRecordStart:        {[]backendState{stAuthed}, (*backend).recordStart},
	editbackend.OpRecordStatus:       {[]backendState{stRecording}, (*backend).recordStatus},
	editbackend.OpRecordStop:         {[]backendState{stRecording}, (*backend).recordStop},
	editbackend.OpRecordPlan:         {[]backendState{stRecorded}, (*backend).recordPlan},
	editbackend.OpRecordApply:        {[]backendState{stRecorded}, (*backend).recordApply},
	editbackend.OpRecordDeleteBackup: {[]backendState{stApplied}, (*backend).recordDeleteBackup},
}

// errSessionEnded: the daemon revoked the live grant; nothing more is written.
var errSessionEnded = errors.New("the daemon ended the live session (idle timeout or reload) — saves made " +
	"before it are kept, later ones were not written")

// backend serves one GUI workflow over stdin/stdout (editbackend). It holds what the CLI holds in
// one run: the password, the control session, the open vault.
type backend struct {
	out   *editbackend.Writer
	outMu sync.Mutex
	// mu serializes requests with the shutdown a signal triggers.
	mu    sync.Mutex
	state backendState

	live      bool
	password  string
	timeout   time.Duration
	session   *liveSession
	pinger    func()
	resources []string
	resource  string
	vault     *vaultfs.Vault
	fsVault   *fscrypt.Vault // offline: the unlocked vault to re-lock
	record    *recordFlow
	ended     atomic.Bool
	// closing: our own cleanup closes the connection, which is not the daemon ending it.
	closing   atomic.Bool
	cleanedUp bool
}

// runBackend is `edit-protected --backend`: the GUI's channel, never a terminal's.
func runBackend() error {
	if term.IsTerminal(os.Stdin.Fd()) {
		return errors.New("--backend is the GUI's channel to this binary, not an interactive mode")
	}
	out, err := reserveStdout()
	if err != nil {
		return err
	}
	b := &backend{out: editbackend.NewWriter(out), live: LiveModeAvailable(), timeout: SessionTimeoutDefault}
	go b.onSignal()
	defer b.cleanup()
	return b.serve(editbackend.NewReader(os.Stdin))
}

// reserveStdout keeps fd 1 for the protocol only: a library's stray print would corrupt the stream,
// so fd 1 becomes stderr and the protocol writes to a private duplicate.
func reserveStdout() (*os.File, error) {
	fd, err := unix.FcntlInt(1, unix.F_DUPFD_CLOEXEC, 3)
	if err != nil {
		return nil, fmt.Errorf("reserving stdout: %w", err)
	}
	if err := unix.Dup2(2, 1); err != nil {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("reserving stdout: %w", err)
	}
	return os.NewFile(uintptr(fd), "backend-out"), nil
}

// onSignal ends the session (a live grant is revoked, an offline vault re-locked) and exits.
func (b *backend) onSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	s := <-sig
	log.Warnf("edit-protected backend: %v — ending the session", s)
	b.cleanup()
	os.Exit(1)
}

func (b *backend) send(m *editbackend.Message) {
	b.outMu.Lock()
	defer b.outMu.Unlock()
	if err := b.out.Write(m); err != nil {
		log.Warnf("edit-protected backend: writing to the GUI: %v", err)
	}
}

// serve answers requests until end, EOF (the GUI went away) or a malformed stream.
func (b *backend) serve(in *editbackend.Reader) error {
	var hello editbackend.Request
	if err := in.Read(&hello); err != nil {
		return fmt.Errorf("reading hello: %w", err)
	}
	if hello.Op != editbackend.OpHello {
		return fmt.Errorf("expected %q first", editbackend.OpHello)
	}
	b.send(b.hello(hello.ID))
	for {
		var req editbackend.Request
		err := in.Read(&req)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading a request: %w", err)
		}
		if req.Op == editbackend.OpEnd {
			b.send(b.end(req.ID))
			return nil
		}
		b.send(b.handle(&req))
	}
}

func (b *backend) hello(id uint64) *editbackend.Message {
	hasPw, err := HashConfigured()
	m := &editbackend.Message{Kind: editbackend.KindReply, ID: id, Hello: &editbackend.Hello{
		Proto: editbackend.Proto, Version: constants.Version, Live: b.live, HasPassword: hasPw,
	}}
	if err != nil {
		m.Err = err.Error()
	}
	return m
}

func (b *backend) handle(req *editbackend.Request) *editbackend.Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	reply := func(err error) *editbackend.Message {
		return &editbackend.Message{Kind: editbackend.KindReply, ID: req.ID, Err: err.Error()}
	}
	op, ok := backendOps[req.Op]
	switch {
	case !ok:
		return reply(fmt.Errorf("unknown op %q", req.Op))
	case !slices.Contains(op.states, b.state):
		return reply(fmt.Errorf("%s is not allowed at this point of the session", req.Op))
	case b.ended.Load() && b.state != stNew:
		return reply(errSessionEnded)
	}
	m, err := op.run(b, req)
	if err != nil {
		return reply(err)
	}
	if m == nil {
		m = &editbackend.Message{}
	}
	m.Kind, m.ID = editbackend.KindReply, req.ID
	return m
}

// end finishes the workflow like the CLI does on exit: narrow the grant (or re-lock), then audit.
func (b *backend) end(id uint64) *editbackend.Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	m := &editbackend.Message{Kind: editbackend.KindReply, ID: id}
	resource, editing, offline := b.resource, b.state == stEditing, b.fsVault != nil
	if editing && offline {
		// The offline CLI audits before its deferred re-lock: names are only readable unlocked.
		m.Findings = auditFindings(resource)
	}
	if err := b.cleanupErr(); err != nil {
		m.Err = err.Error()
	}
	if editing && !offline {
		m.Findings = auditFindings(resource) // after the grant is narrowed, as runLiveEdit does
	}
	return m
}

func (b *backend) cleanup() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.cleanupErr(); err != nil {
		log.Error(err)
	}
}

// cleanupErr releases everything: the vault, the live grant, an offline unlock. Idempotent.
func (b *backend) cleanupErr() error {
	if b.cleanedUp {
		return nil
	}
	b.cleanedUp = true
	b.closing.Store(true)
	b.state = stDone
	if b.vault != nil {
		b.vault.Close()
		b.vault = nil
	}
	if b.record != nil {
		b.record.close()
		b.record = nil
	}
	if b.session != nil {
		b.session.End()
		b.session = nil
	}
	if b.fsVault != nil {
		if err := protected.RelockResources(b.fsVault, b.resource); err != nil {
			return fmt.Errorf("could not fully re-lock %s: %w — run `fscrypt lock %s` manually as soon as possible",
				b.resource, err, b.resource)
		}
		b.fsVault = nil
	}
	return nil
}

// checkSecret refuses a password that would break the control protocol's line framing.
func checkSecret(pw string) error {
	if pw == "" {
		return errors.New("no password entered")
	}
	if strings.ContainsAny(pw, "\r\n\x00") || !utf8.ValidString(pw) {
		return errors.New("the password contains a line break or an invalid character")
	}
	return nil
}

func (b *backend) auth(req *editbackend.Request) (*editbackend.Message, error) {
	if !b.live {
		return nil, errors.New("live mode needs the daemon running with an edit-protected password configured " +
			"(the control socket " + ControlSocket + " is not present)")
	}
	if err := checkSecret(req.Password); err != nil {
		return nil, err
	}
	if req.TimeoutSec != 0 {
		// Bounded in seconds first: a huge value would overflow the Duration.
		if req.TimeoutSec < 1 || req.TimeoutSec > int64(SessionTimeoutMax/time.Second) {
			return nil, fmt.Errorf("the session timeout must be between 1s and %s", SessionTimeoutMax)
		}
		b.timeout = time.Duration(req.TimeoutSec) * time.Second
	}
	session, err := dialLiveSession()
	if err != nil {
		return nil, err
	}
	session.timeout = b.timeout
	resources, err := session.Authenticate(req.Password)
	if err != nil {
		session.End()
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	b.session, b.password, b.resources, b.state = session, req.Password, resources, stAuthed
	return &editbackend.Message{Resources: resources}, nil
}

// selectResource is runLiveEdit after the picker: widen the grant, pin the vault, watch for the
// daemon ending the session.
func (b *backend) selectResource(req *editbackend.Request) (*editbackend.Message, error) {
	if !slices.Contains(b.resources, req.Resource) {
		return nil, fmt.Errorf("%s is not a guarded directory in the running daemon", req.Resource)
	}
	if err := b.session.Select(req.Resource); err != nil {
		return nil, fmt.Errorf("activating live edit access to %s: %w", req.Resource, err)
	}
	b.resource, b.state = req.Resource, stEditing
	log.Infof("editing %s LIVE — the daemon keeps guarding it; write access is granted only for this session "+
		"and revoked after %s without input", req.Resource, b.timeout)
	_, ended := b.session.Watch()
	b.pinger = b.session.Pinger(pingInterval(b.timeout))
	go b.watchEnded(ended)
	return b.openVault()
}

func (b *backend) watchEnded(ended <-chan struct{}) {
	<-ended
	b.markEnded()
}

// markEnded tells the GUI once that the daemon revoked the session; later ops are refused.
func (b *backend) markEnded() {
	if b.closing.Load() || b.ended.Swap(true) {
		return
	}
	b.send(&editbackend.Message{Kind: editbackend.KindEnded, Err: errSessionEnded.Error()})
}

// openVault pins the selected resource; on failure only end is left (it narrows or re-locks).
func (b *backend) openVault() (*editbackend.Message, error) {
	v, err := vaultfs.Open(b.resource)
	if err != nil {
		b.state = stDone
		return nil, fmt.Errorf("opening %s: %w", b.resource, err)
	}
	b.vault = v
	return &editbackend.Message{SingleFile: v.SingleFile()}, nil
}

// offlineScan is runOfflineEdit up to the picker: only with the daemon stopped.
func (b *backend) offlineScan(*editbackend.Request) (*editbackend.Message, error) {
	if b.live {
		return nil, errors.New("the daemon is running with a password: edit live (auth) instead")
	}
	if err := protected.RequireDaemonStopped(); err != nil {
		return nil, err
	}
	encrypted, err := protected.ScanEncryptedCatalogDirs(fscrypt.New())
	if err != nil {
		return nil, err
	}
	if len(encrypted) == 0 {
		return nil, errors.New("no catalog directory is currently encrypted: nothing to edit (run the installer first)")
	}
	b.resources = encrypted
	return &editbackend.Message{Resources: encrypted}, nil
}

// offlineOpen unlocks ONE scanned directory with the master key; end (or EOF, or a signal)
// re-locks it.
func (b *backend) offlineOpen(req *editbackend.Request) (*editbackend.Message, error) {
	if b.live || !slices.Contains(b.resources, req.Resource) {
		return nil, fmt.Errorf("%s is not an encrypted catalog directory of the last scan", req.Resource)
	}
	if err := protected.RequireDaemonStopped(); err != nil {
		return nil, err
	}
	vault := fscrypt.New()
	// Never edit a directory whose policy doesn't unlock with the master key: writes would land in
	// an unreadable policy and the re-lock couldn't complete.
	if err := protected.VerifyEncryptedKeys(vault, []string{req.Resource}); err != nil {
		return nil, err
	}
	log.Infof("unlocking %s ...", req.Resource)
	if err := vault.Unlock(req.Resource); err != nil {
		return nil, fmt.Errorf("unlock %s: %w", req.Resource, err)
	}
	b.fsVault, b.resource, b.state = vault, req.Resource, stEditing
	log.Infof("editing %s (it is unlocked and accessible while the editor is open)", req.Resource)
	return b.openVault()
}

func (b *backend) ping(*editbackend.Request) (*editbackend.Message, error) {
	if b.pinger != nil {
		b.pinger()
	}
	return nil, nil
}

func (b *backend) passwordState(*editbackend.Request) (*editbackend.Message, error) {
	st := &editbackend.PwState{DaemonRunning: daemonRunning()}
	encoded, err := LoadHash()
	switch {
	case errors.Is(err, ErrNoPassword):
	case err != nil:
		return nil, err
	default:
		origin, oErr := OriginOf(encoded)
		if oErr != nil {
			return nil, fmt.Errorf("the stored password record is unreadable (%w)", oErr)
		}
		st.Configured, st.Origin = true, string(origin)
	}
	return &editbackend.Message{Password: st}, nil
}

func (b *backend) setPassword(req *editbackend.Request) (*editbackend.Message, error) {
	if err := checkSecret(req.NewPassword); err != nil {
		return nil, err
	}
	err := setPassword(func() (string, error) { return req.Password, checkSecret(req.Password) },
		func() (string, error) { return req.NewPassword, ValidatePassword(req.NewPassword) })
	return nil, err
}

func (b *backend) forwardCheck(req *editbackend.Request) (*editbackend.Message, error) {
	out := make([]editbackend.BinaryCheck, 0, len(req.Binaries))
	for _, p := range req.Binaries {
		c := editbackend.BinaryCheck{Path: p}
		if !filepath.IsAbs(p) || strings.ContainsRune(p, 0) {
			c.Err = "not an absolute path"
		} else if placed, runtimes, err := forwardCheck(p); err != nil {
			c.Err = err.Error()
		} else {
			c.SystemPlaced, c.Runtimes = placed, runtimes
		}
		out = append(out, c)
	}
	return &editbackend.Message{Binaries: out}, nil
}

// forward is runForward after its prompts: the daemon resolves and grants; the granted binaries'
// events stream to the GUI until end or the daemon ends the grant.
func (b *backend) forward(req *editbackend.Request) (*editbackend.Message, error) {
	if len(req.Binaries) == 0 || len(req.Resources) == 0 {
		return nil, errors.New("pick at least one guarded directory and one binary")
	}
	for _, r := range req.Resources {
		if !slices.Contains(b.resources, r) {
			return nil, fmt.Errorf("%s is not a guarded directory in the running daemon", r)
		}
	}
	for _, p := range req.Binaries {
		if !filepath.IsAbs(p) || strings.ContainsRune(p, 0) {
			return nil, fmt.Errorf("%q: binaries must be absolute paths", p)
		}
	}
	events := make([]string, 0, len(req.Events))
	for _, e := range req.Events {
		events = append(events, strings.ToUpper(strings.TrimSpace(e)))
	}
	if _, err := common.ParseEventsFlag(events); err != nil {
		return nil, err
	}
	notes, err := b.session.Forward(req.Resources, req.Binaries, req.Allow, events)
	if err != nil {
		return nil, fmt.Errorf("activating temporary access: %w", err)
	}
	b.resources, b.state = req.Resources, stForwarding
	evs, ended := b.session.Watch()
	b.pinger = b.session.Pinger(pingInterval(b.timeout))
	go b.streamEvents(evs, ended)
	return &editbackend.Message{Notes: notes}, nil
}

func (b *backend) streamEvents(evs <-chan guard.GuardEvent, ended <-chan struct{}) {
	for {
		select {
		case ev := <-evs:
			b.send(&editbackend.Message{Kind: editbackend.KindEvent, Event: &editbackend.Event{
				Time: time.Now(), PID: ev.PID, UID: ev.UID, Comm: logging.SanitizeText(ev.Comm), Op: ev.Op(),
				Path: logging.SanitizeText(ev.Path), Dest: logging.SanitizeText(ev.Dest), Allowed: !ev.Blocked,
			}})
		case <-ended:
			b.markEnded()
			return
		}
	}
}

// configGet returns daemon.conf: on the authenticated connection the first time, then on a fresh
// one (the CLI's fetchConfig).
func (b *backend) configGet(*editbackend.Request) (*editbackend.Message, error) {
	var data []byte
	var err error
	if b.session != nil {
		data, err = b.session.Config()
		b.session.End()
		b.session = nil
	} else {
		data, err = fetchConfig(b.password)
	}
	if err != nil {
		return nil, err
	}
	b.state = stConfig
	return &editbackend.Message{Data: data}, nil
}

// configPut is checkAndApply after the GUI's diff and confirmation: written only while daemon.conf
// still holds Base (putConfig); a stale Base comes back with the current file.
func (b *backend) configPut(req *editbackend.Request) (*editbackend.Message, error) {
	if len(req.Data) > editbackend.MaxWriteBytes {
		return nil, errors.New("the configuration is too large")
	}
	if _, err := daemonconfig.Parse(req.Data); err != nil {
		return nil, fmt.Errorf("the configuration does not parse: %w", err)
	}
	err := putConfig(b.password, req.Base, req.Data)
	if errors.Is(err, errConfigChanged) {
		fresh, ferr := fetchConfig(b.password)
		if ferr != nil {
			return nil, ferr
		}
		return &editbackend.Message{Err: err.Error(), Stale: true, Data: fresh}, nil
	}
	return nil, err
}

// abs maps a wire path onto the open vault; the vault confines it (os.Root).
func (b *backend) abs(rel string) (string, error) {
	c, err := editbackend.CleanRel(rel)
	if err != nil {
		return "", err
	}
	return filepath.Join(b.resource, c), nil
}

func (b *backend) list(req *editbackend.Request) (*editbackend.Message, error) {
	path, err := b.abs(req.Path)
	if err != nil {
		return nil, err
	}
	if b.vault.SingleFile() {
		info, serr := b.vault.Lstat(path)
		if serr != nil {
			return nil, serr
		}
		return &editbackend.Message{Entries: []editbackend.Entry{entryOf(info.Name(), info)}}, nil
	}
	dirents, err := b.vault.ReadDir(path)
	if err != nil {
		return nil, err
	}
	out := make([]editbackend.Entry, 0, len(dirents))
	for _, d := range dirents {
		info, err := d.Info()
		if err != nil {
			continue // vanished since the listing
		}
		out = append(out, entryOf(d.Name(), info))
	}
	return &editbackend.Message{Entries: out}, nil
}

func entryOf(name string, info os.FileInfo) editbackend.Entry {
	e := editbackend.Entry{
		Name: name, Dir: info.IsDir(), Symlink: info.Mode()&os.ModeSymlink != 0, Size: info.Size(),
		Mode: vaultfs.FileModeToUnixMode(info.Mode()), ModeStr: info.Mode().String(), ModTime: info.ModTime(),
	}
	if !utf8.ValidString(name) {
		e.Name, e.Invalid = strings.ToValidUTF8(name, "�"), true
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		e.UID, e.GID = st.Uid, st.Gid
	}
	return e
}

// read mirrors the TUI's preview: symlinks refused, binary and oversized files not sent.
func (b *backend) read(req *editbackend.Request) (*editbackend.Message, error) {
	path, err := b.abs(req.Path)
	if err != nil {
		return nil, err
	}
	info, err := b.vault.Lstat(path)
	if err != nil {
		return nil, err
	}
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return nil, fmt.Errorf("%s is a symlink — editing it is refused", path)
	case !info.Mode().IsRegular():
		return nil, fmt.Errorf("%s is not a regular file", path)
	case info.Size() > vaultfs.MaxEditableBytes:
		return &editbackend.Message{File: &editbackend.File{Size: info.Size(), TooLarge: true}}, nil
	}
	data, err := b.vault.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := &editbackend.File{Size: int64(len(data))}
	switch {
	case len(data) > vaultfs.MaxEditableBytes:
		f.TooLarge = true
	case vaultfs.IsBinaryContent(data):
		f.Binary = true
	default:
		f.Data = data
	}
	return &editbackend.Message{File: f}, nil
}

func (b *backend) write(req *editbackend.Request) (*editbackend.Message, error) {
	if len(req.Data) > editbackend.MaxWriteBytes {
		return nil, fmt.Errorf("the content is over %d MiB", editbackend.MaxWriteBytes>>20)
	}
	path, err := b.abs(req.Path)
	if err != nil {
		return nil, err
	}
	return nil, b.vault.Write(path, req.Data)
}

func (b *backend) create(req *editbackend.Request) (*editbackend.Message, error) {
	dir, err := b.abs(req.Path)
	if err != nil {
		return nil, err
	}
	if !utf8.ValidString(req.Name) || strings.ContainsRune(req.Name, 0) {
		return nil, errors.New("malformed name")
	}
	target, err := b.vault.Create(dir, req.Name, req.Dir)
	if err != nil {
		return nil, err
	}
	if err := b.vault.ApplyNewMeta(target, req.Dir); err != nil {
		return nil, fmt.Errorf("created %s, but: %w", target, err)
	}
	return nil, nil
}

func (b *backend) chmod(req *editbackend.Request) (*editbackend.Message, error) {
	mode, err := vaultfs.ParseMode(req.Mode)
	if err != nil {
		return nil, err
	}
	path, err := b.abs(req.Path)
	if err != nil {
		return nil, err
	}
	return nil, b.vault.Chmod(path, mode)
}

// chown accepts only a login user's own uid:gid, the TUI picker's choices.
func (b *backend) chown(req *editbackend.Request) (*editbackend.Message, error) {
	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	if !slices.ContainsFunc(users, func(u inst.User) bool {
		return int(u.UID) == req.UID && int(u.GID) == req.GID
	}) {
		return nil, fmt.Errorf("uid %d gid %d is not a login user", req.UID, req.GID)
	}
	path, err := b.abs(req.Path)
	if err != nil {
		return nil, err
	}
	return nil, b.vault.Chown(path, req.UID, req.GID)
}

func (b *backend) remove(req *editbackend.Request) (*editbackend.Message, error) {
	path, err := b.abs(req.Path)
	if err != nil {
		return nil, err
	}
	return nil, b.vault.Remove(path)
}

func (b *backend) users(*editbackend.Request) (*editbackend.Message, error) {
	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	out := make([]editbackend.User, 0, len(users))
	for _, u := range users {
		out = append(out, editbackend.User{Name: u.Name, UID: int(u.UID), GID: int(u.GID)})
	}
	return &editbackend.Message{Users: out}, nil
}
