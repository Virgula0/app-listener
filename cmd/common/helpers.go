package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

type RawTarget struct {
	AbsPath string
	IsDir   bool
}

func ResolveTargets(paths []string) ([]RawTarget, error) {
	seen := make(map[string]bool)
	var targets []RawTarget

	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			log.Errorf("Invalid path %q: %v", p, err)
			return nil, err
		}

		if seen[abs] {
			return nil, fmt.Errorf("duplicate watch path: %s", abs)
		}
		seen[abs] = true

		info, err := os.Stat(abs)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Errorf("Path does not exist: %s", abs)
			} else {
				log.Errorf("Cannot access path %s: %v", abs, err)
			}
			return nil, err
		}

		targets = append(targets, RawTarget{AbsPath: abs, IsDir: info.IsDir()})
	}

	return targets, nil
}

func ValidateTargets(targets []RawTarget) error {
	for i, a := range targets {
		for j, b := range targets {
			if i == j {
				continue
			}
			if err := validatePair(a, b); err != nil {
				return err
			}
		}
	}
	return nil
}

func validatePair(a, b RawTarget) error {
	if a.IsDir && b.IsDir {
		if isSubDir(a.AbsPath, b.AbsPath) {
			return fmt.Errorf(
				"%q is a subdirectory of %q \u2014 remove the more specific path",
				b.AbsPath, a.AbsPath)
		}
		if isSubDir(b.AbsPath, a.AbsPath) {
			return fmt.Errorf(
				"%q is a subdirectory of %q \u2014 remove the more specific path",
				a.AbsPath, b.AbsPath)
		}
	}

	if !a.IsDir && b.IsDir && pathWithinDir(a.AbsPath, b.AbsPath) {
		return fmt.Errorf(
			"%q is already covered by directory %q \u2014 remove the redundant file path",
			a.AbsPath, b.AbsPath)
	}
	if a.IsDir && !b.IsDir && pathWithinDir(b.AbsPath, a.AbsPath) {
		return fmt.Errorf(
			"%q is already covered by directory %q \u2014 remove the redundant file path",
			b.AbsPath, a.AbsPath)
	}
	return nil
}

func isSubDir(child, parent string) bool {
	return strings.HasPrefix(child+"/", parent+"/")
}

func pathWithinDir(path, dir string) bool {
	return strings.HasPrefix(path+"/", dir+"/")
}

func BuildEBPFTargets(rawTargets []RawTarget) []ebpf.Target {
	targets := make([]ebpf.Target, len(rawTargets))
	for i, rt := range rawTargets {
		t := ebpf.Target{Path: rt.AbsPath, IsDir: rt.IsDir}
		if rt.IsDir {
			t.Dir = rt.AbsPath
		} else {
			t.Dir = filepath.Dir(rt.AbsPath)
			t.File = filepath.Base(rt.AbsPath)
		}
		targets[i] = t
	}
	return targets
}

func MakeDisplayPaths(rawTargets []RawTarget) []string {
	paths := make([]string, len(rawTargets))
	for i, rt := range rawTargets {
		paths[i] = rt.AbsPath
	}
	return paths
}

func CheckEBPF() error {
	log.Info("Checking eBPF availability...")
	if err := ebpf.Check(); err != nil {
		log.Errorf("eBPF check failed: %v", err)
		return err
	}
	log.Info("eBPF available")
	return nil
}

func ParseEventsFlag(eventsFlag []string) ([]ebpf.EventType, error) {
	if len(eventsFlag) == 0 {
		return ebpf.EventTypes(), nil
	}

	var parsed []ebpf.EventType
	for _, s := range eventsFlag {
		et, ok := ebpf.ParseEventType(strings.TrimSpace(s))
		if !ok {
			return nil, fmt.Errorf("unknown event type %q (valid: OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP)", s)
		}
		parsed = append(parsed, et)
	}
	return parsed, nil
}

func WarnIgnoredFlags(rawTargets []RawTarget, recursive bool, depth int) {
	for _, rt := range rawTargets {
		if !rt.IsDir && recursive {
			log.Warnf("--recursive is ignored when monitoring a single file (%s)", rt.AbsPath)
		}
		if !rt.IsDir && depth > 0 {
			log.Warnf("--depth is ignored when monitoring a single file (%s)", rt.AbsPath)
		}
	}
}
