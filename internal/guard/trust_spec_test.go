package guard

import (
	"os"
	"testing"
)

func TestKernelHasFunc(t *testing.T) {
	if _, err := os.Stat("/sys/kernel/btf/vmlinux"); err != nil {
		t.Skip("no kernel BTF")
	}
	if !KernelHasFunc("vfs_read") {
		t.Fatal("vfs_read not found in vmlinux BTF")
	}
	if KernelHasFunc("app_listener_no_such_function") {
		t.Fatal("a missing function was reported present: trust_memfd_alloc would fail the trust object")
	}
}
