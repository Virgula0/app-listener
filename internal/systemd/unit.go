package systemd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// DaemonUnitCommand is the hidden subcommand that prints the daemon unit a binary ships.
	DaemonUnitCommand = "daemon-unit"
	// noMetadataFlag is the ExecStart flag `install --allow-metadata-output` drops.
	noMetadataFlag = " --no-log-metadata-blocks"
	// shippedUnitTimeout bounds the DaemonUnitCommand run of a binary about to be deployed.
	shippedUnitTimeout = 30 * time.Second
	// maxUnitSize caps what DaemonUnitCommand may print.
	maxUnitSize = 64 << 10
)

// DaemonUnitWithMetadataOutput drops --no-log-metadata-blocks from the daemon unit's ExecStart.
func DaemonUnitWithMetadataOutput(unit []byte) []byte {
	return bytes.ReplaceAll(unit, []byte(noMetadataFlag), nil)
}

// InstalledDaemonUnit returns the on-disk daemon unit (present false when it is not installed).
func InstalledDaemonUnit() (unit []byte, present bool, err error) {
	unit, err = os.ReadFile(daemonUnitFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("reading %s: %w", daemonUnitFile, err)
	}
	return unit, true, nil
}

// ShippedDaemonUnit runs binPath's DaemonUnitCommand and returns the unit that build bundles. Only
// call it on a binary already verified for deployment: it is executed as root.
func ShippedDaemonUnit(binPath string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shippedUnitTimeout)
	defer cancel()
	var stdout limitedBuffer
	cmd := exec.CommandContext(ctx, binPath, DaemonUnitCommand)
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s: %w", binPath, DaemonUnitCommand, err)
	}
	if stdout.overflow {
		return nil, fmt.Errorf("%s %s printed more than %d bytes", binPath, DaemonUnitCommand, maxUnitSize)
	}
	unit := stdout.buf.Bytes()
	if !bytes.Contains(unit, []byte("[Service]")) || !bytes.Contains(unit, []byte("\nExecStart=")) {
		return nil, fmt.Errorf("%s %s did not print a service unit", binPath, DaemonUnitCommand)
	}
	return unit, nil
}

// DesiredDaemonUnit adapts the shipped unit to the installed one: an installed ExecStart without
// --no-log-metadata-blocks (install --allow-metadata-output) keeps that choice.
func DesiredDaemonUnit(installed, shipped []byte) []byte {
	if !bytes.Contains(installed, []byte(noMetadataFlag)) {
		return DaemonUnitWithMetadataOutput(shipped)
	}
	return shipped
}

// ConfirmUnitFunc decides whether the installed unit at path is replaced by desired.
type ConfirmUnitFunc func(path string, installed, desired []byte) (bool, error)

// PlanDaemonUnitUpdate diffs the daemon unit binPath ships against the installed one and returns
// the unit to write, or nil to keep the installed one (not installed, identical, unreadable from
// binPath, or declined by confirm).
func PlanDaemonUnitUpdate(binPath string, confirm ConfirmUnitFunc) ([]byte, error) {
	installed, present, err := InstalledDaemonUnit()
	if err != nil || !present {
		return nil, err
	}
	shipped, err := ShippedDaemonUnit(binPath)
	if err != nil {
		log.Warnf("cannot read the daemon unit the new binary ships (%v): keeping the installed %s", err, daemonUnitFile)
		return nil, nil
	}
	desired := DesiredDaemonUnit(installed, shipped)
	if bytes.Equal(installed, desired) {
		log.Infof("%s already matches the new binary's unit", daemonUnitFile)
		return nil, nil
	}
	overwrite, err := confirm(daemonUnitFile, installed, desired)
	if err != nil {
		return nil, err
	}
	if !overwrite {
		log.Warnf("keeping the installed %s", daemonUnitFile)
		return nil, nil
	}
	return desired, nil
}

// WriteDaemonUnit atomically replaces the daemon unit (mode 0644); the caller runs daemon-reload.
func WriteDaemonUnit(unit []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(daemonUnitFile), ".app-listener-unit-*")
	if err != nil {
		return fmt.Errorf("staging the daemon unit: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(unit); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(0o644); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), daemonUnitFile); err != nil {
		return fmt.Errorf("replacing %s: %w", daemonUnitFile, err)
	}
	log.Infof("updated %s", daemonUnitFile)
	return nil
}

// limitedBuffer keeps the first maxUnitSize bytes written and flags anything beyond. Not an
// embedded bytes.Buffer: its promoted ReadFrom would let io.Copy bypass the cap.
type limitedBuffer struct {
	buf      bytes.Buffer
	overflow bool
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if room := maxUnitSize - b.buf.Len(); n > room {
		b.overflow = true
		p = p[:max(room, 0)]
	}
	b.buf.Write(p)
	return n, nil
}
