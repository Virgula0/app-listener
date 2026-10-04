package guard

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestGuardLSMHooksIoctlCompatOnlyWhenPresent(t *testing.T) {
	var o GuardObjects
	for _, compat := range []bool{false, true} {
		got := map[string]bool{}
		for _, h := range guardLSMHooks(&o, compat) {
			got[h.hook] = true
		}
		if !got["file_ioctl"] {
			t.Fatalf("compat=%v: file_ioctl not attached", compat)
		}
		if got["file_ioctl_compat"] != compat {
			t.Fatalf("compat=%v: file_ioctl_compat attached=%v", compat, got["file_ioctl_compat"])
		}
	}
}

func TestIoctlHooksAreRequired(t *testing.T) {
	for _, h := range []string{"file_ioctl", "file_ioctl_compat"} {
		if !requiredHooks[h] {
			t.Fatalf("%s is optional: a failed attach would leave btrfs snapshots ungated", h)
		}
	}
}

func TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(t *testing.T) {
	spec, compat, err := guardSpec()
	if err != nil {
		t.Fatal(err)
	}
	if compat != kernelHasFunc(ioctlCompatTarget) {
		t.Fatalf("compat=%v disagrees with the kernel's BTF", compat)
	}
	p := spec.Programs[GuardProgGuardFileIoctlCompat]
	if p == nil {
		t.Fatal("guard_file_ioctl_compat missing from the spec")
	}
	if inert := p.AttachTo == ""; inert == compat {
		t.Fatalf("compat=%v but program attach target %q", compat, p.AttachTo)
	}
}

// A suspect mark rides in its own byte: a gate's label beside it survives.
func TestParseGuardEventSuspectLabels(t *testing.T) {
	for _, tc := range []struct {
		reason     uint32
		op, fsGate string
	}{
		{0, "OPEN", ""},
		{suspectPreload << guardReasonSuspectShift, "PRELOADED", ""},
		{suspectLaunch << guardReasonSuspectShift, "LAUNCH", ""},
		{suspectPreload<<guardReasonSuspectShift | guardReasonRawDevice, "PRELOADED", RawDeviceResourceLabel},
	} {
		var buf bytes.Buffer
		if err := binary.Write(&buf, binary.LittleEndian, bpfGuardEvent{Reason: tc.reason, Blocked: 1}); err != nil {
			t.Fatal(err)
		}
		ev, _, ok := parseGuardEvent(buf.Bytes())
		if !ok {
			t.Fatalf("reason %#x: decode failed", tc.reason)
		}
		if ev.Op() != tc.op || ev.FsGate != tc.fsGate || ev.Process != "" {
			t.Fatalf("reason %#x: Op=%q FsGate=%q Process=%q, want %q %q \"\"", tc.reason, ev.Op(), ev.FsGate,
				ev.Process, tc.op, tc.fsGate)
		}
	}
}

func TestParseGuardEventFsGateLabels(t *testing.T) {
	for reason, want := range map[uint32]string{
		0:                     "",
		guardReasonRawDevice:  RawDeviceResourceLabel,
		guardReasonBtrfsIoctl: BtrfsIoctlResourceLabel,
		guardReasonPtrace:     "",
	} {
		var buf bytes.Buffer
		if err := binary.Write(&buf, binary.LittleEndian, bpfGuardEvent{Reason: reason, Blocked: 1}); err != nil {
			t.Fatal(err)
		}
		ev, _, ok := parseGuardEvent(buf.Bytes())
		if !ok {
			t.Fatalf("reason %d: decode failed", reason)
		}
		if ev.FsGate != want {
			t.Fatalf("reason %d: FsGate=%q, want %q", reason, ev.FsGate, want)
		}
	}
}
