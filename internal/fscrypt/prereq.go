package fscrypt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/fscrypt/filesystem"
)

// Prereq is a host requirement the installer can satisfy by running ONE
// privileged command as root. Every command the installer may run on the
// user's behalf is produced by FilesystemPrereqs below — this is the single
// place to audit the installer's auto-run command surface.
type Prereq struct {
	// Title is a short imperative summary ("Enable fscrypt on /dev/sda2").
	Title string
	// Reason explains what is currently wrong and why the command fixes it.
	Reason string
	// Argv is the exact command, executed as root, argument by argument.
	Argv []string
}

// Command renders Argv as a copy-pasteable shell string (display only).
func (p Prereq) Command() string { return strings.Join(p.Argv, " ") }

// FilesystemPrereqs inspects the filesystem backing path and returns the
// commands needed to make it usable for fscrypt encryption, in the order
// they must run. It returns nil when the filesystem is already ready.
//
// A non-nil error is TERMINAL: the filesystem type or the kernel cannot
// support fscrypt at all, so no command would help and the caller must
// abort. Fixable conditions never return an error — they come back as
// Prereq entries the caller can offer to run.
func (v *Vault) FilesystemPrereqs(path string) ([]Prereq, error) {
	mnt, err := filesystem.FindMount(path)
	if err != nil {
		return nil, fmt.Errorf("resolve filesystem of %s: %w", path, err)
	}

	var out []Prereq

	// 1. The on-disk "encrypt" feature flag (ext4) / equivalent (f2fs).
	if serr := mnt.CheckSupport(); serr != nil {
		var notEnabled *filesystem.ErrEncryptionNotEnabled
		var notSupported *filesystem.ErrEncryptionNotSupported
		switch {
		case errors.As(serr, &notEnabled):
			m := notEnabled.Mount
			argv := []string{"tune2fs", "-O", "encrypt", m.Device}
			if m.FilesystemType == "f2fs" {
				argv = []string{"fsck.f2fs", "-O", "encrypt", m.Device}
			}
			out = append(out, Prereq{
				Title: fmt.Sprintf("Enable the fscrypt 'encrypt' feature on %s (%s)", m.Device, m.FilesystemType),
				Reason: fmt.Sprintf(
					"%s was created without the encryption feature flag, so no directory on it can be encrypted.",
					m.Device),
				Argv: argv,
			})
		case errors.As(serr, &notSupported):
			return nil, classifySupportError(path, serr)
		default:
			return nil, fmt.Errorf("fscrypt support check for %s: %w", path, serr)
		}
	}

	// 2. The per-filesystem fscrypt metadata directory (`fscrypt setup`).
	if cerr := mnt.CheckSetup(nil); cerr != nil {
		var notSetup *filesystem.ErrNotSetup
		var notSupported *filesystem.ErrSetupNotSupported
		switch {
		case errors.As(cerr, &notSetup):
			out = append(out, Prereq{
				Title: fmt.Sprintf("Initialize %s for fscrypt", mnt.Path),
				Reason: fmt.Sprintf(
					"the filesystem mounted at %s has no fscrypt metadata directory — 'fscrypt setup' has never run on it.",
					mnt.Path),
				Argv: []string{"fscrypt", "setup", mnt.Path, "--all-users"},
			})
		case errors.As(cerr, &notSupported):
			return nil, classifySetupError(path, cerr)
		default:
			return nil, fmt.Errorf("fscrypt setup check for %s: %w", path, cerr)
		}
	}

	return out, nil
}
