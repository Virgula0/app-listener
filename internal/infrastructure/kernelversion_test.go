package ebpf

import (
	"reflect"
	"testing"

	"github.com/cilium/ebpf/features"
)

func TestParseKernelRelease(t *testing.T) {
	for release, want := range map[string]ciliumVersion{
		"7.0.11-arch1-1":     {7, 0, 11},
		"6.17.0-5-generic":   {6, 17, 0},
		"5.15.300":           {5, 15, 255},
		"6.8":                {6, 8, 0},
		"6.8-rc1":            {6, 8, 0},
		"6.12.rc3-something": {6, 12, 0},
	} {
		got, err := parseKernelRelease(release)
		if err != nil || got != want {
			t.Errorf("parseKernelRelease(%q) = %v, %v; want %v", release, got, err, want)
		}
	}
	for _, bad := range []string{"", "linux", "6", "x.1.2", "6.y"} {
		if _, err := parseKernelRelease(bad); err == nil {
			t.Errorf("parseKernelRelease(%q) accepted", bad)
		}
	}
}

// A sentinel proves the linkname binds cilium's own variable with a matching layout: an unbound
// copy would leave cilium's vDSO lookup in place, which returns the same real version.
func TestCiliumKernelVersionOverride(t *testing.T) {
	orig := ciliumKernelVersion
	if reflect.ValueOf(orig).Pointer() != reflect.ValueOf(unameKernelVersion).Pointer() {
		t.Fatal("init did not install the uname override")
	}
	t.Cleanup(func() { ciliumKernelVersion = orig })
	ciliumKernelVersion = func() (ciliumVersion, error) { return ciliumVersion{1, 2, 3}, nil }
	code, err := features.LinuxVersionCode()
	if err != nil || code != 0x010203 {
		t.Fatalf("features.LinuxVersionCode() = %#x, %v; want 0x10203 from the override", code, err)
	}
	ciliumKernelVersion = orig
	if _, err := features.LinuxVersionCode(); err != nil {
		t.Fatalf("uname override: %v", err)
	}
}
