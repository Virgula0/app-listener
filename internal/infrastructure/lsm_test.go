package ebpf

import (
	"os"
	"path/filepath"
	"testing"
)

func writeLSM(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "lsm")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write lsm fixture: %v", err)
	}
	return p
}

func TestBpfLsmListed(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"bpf present", "lockdown,capability,landlock,yama,apparmor,ima,evm,bpf", true},
		{"bpf leading", "bpf,landlock,yama", true},
		{"bpf only", "bpf", true},
		{"no bpf", "lockdown,capability,landlock,yama,apparmor,ima,evm", false},
		{"empty", "", false},
		{"whitespace padding", "  bpf , landlock ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bpfLsmListed(writeLSM(t, tt.content))
			if err != nil {
				t.Fatalf("bpfLsmListed: %v", err)
			}
			if got != tt.want {
				t.Errorf("bpfLsmListed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBpfLsmListedMissingFile(t *testing.T) {
	if _, err := bpfLsmListed(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("expected error for missing lsm file, got nil")
	}
}

func TestCheckBPFLSM(t *testing.T) {
	t.Run("bpf present passes", func(t *testing.T) {
		t.Setenv(assumeBpfLsmEnv, "")
		lsmFile := writeLSM(t, "bpf,landlock,yama")
		if err := CheckBPFLSMAt(lsmFile); err != nil {
			t.Errorf("CheckBPFLSMAt: %v", err)
		}
	})

	t.Run("bpf missing fails", func(t *testing.T) {
		t.Setenv(assumeBpfLsmEnv, "")
		lsmFile := writeLSM(t, "lockdown,capability,landlock,yama,apparmor,ima,evm")
		if err := CheckBPFLSMAt(lsmFile); err == nil {
			t.Error("expected error when bpf is missing, got nil")
		}
	})

	t.Run("unreadable file fails", func(t *testing.T) {
		t.Setenv(assumeBpfLsmEnv, "")
		if err := CheckBPFLSMAt(filepath.Join(t.TempDir(), "missing")); err == nil {
			t.Error("expected error for unreadable lsm file, got nil")
		}
	})

	t.Run("env override bypasses", func(t *testing.T) {
		t.Setenv(assumeBpfLsmEnv, "1")
		if err := CheckBPFLSMAt(filepath.Join(t.TempDir(), "missing")); err != nil {
			t.Errorf("env override should bypass the check, got: %v", err)
		}
	})
}
