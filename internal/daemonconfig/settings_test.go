package daemonconfig

import (
	"strconv"
	"testing"
)

func TestSettingsMaxInodesParsed(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, "[settings]\nmax_inodes: 65536\n\n[watch "+dir+"]\nneed_encryption: false\n/usr/bin/ssh\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.MaxInodes != 65536 {
		t.Fatalf("MaxInodes = %d", cfg.MaxInodes)
	}
	if len(cfg.Resources) != 1 || len(cfg.Resources[0].Binaries) != 1 {
		t.Fatalf("a [settings] block must not touch the watch sections: %+v", cfg.Resources)
	}
}

func TestSettingsDefaultMaxInodesIsZero(t *testing.T) {
	cfg, err := Load(writeConfig(t, "[watch "+t.TempDir()+"]\nneed_encryption: false\n/usr/bin/ssh\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.MaxInodes != 0 {
		t.Fatalf("MaxInodes = %d without a [settings] block", cfg.MaxInodes)
	}
}

func TestSettingsRejectsMalformedLines(t *testing.T) {
	for name, line := range map[string]string{
		"below floor":  "max_inodes: " + strconv.Itoa(MaxInodesFloor-1),
		"above cap":    "max_inodes: " + strconv.Itoa(MaxInodesCeiling+1),
		"zero":         "max_inodes: 0",
		"negative":     "max_inodes: -1",
		"not a number": "max_inodes: lots",
		"suffix":       "max_inodes: 64k",
		"no value":     "max_inodes:",
		"no colon":     "max_inodes 65536",
		"unknown key":  "max_inode: 65536",
		"duplicate":    "max_inodes: 65536\nmax_inodes: 65536",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Load(writeConfig(t, "[settings]\n"+line+"\n")); err == nil {
				t.Fatalf("%q must be refused", line)
			}
		})
	}
}

func TestSettingsBoundsAccepted(t *testing.T) {
	for _, n := range []int{MaxInodesFloor, MaxInodesCeiling} {
		cfg, err := Load(writeConfig(t, "[settings]\nmax_inodes: "+strconv.Itoa(n)+"\n"))
		if err != nil || cfg.MaxInodes != uint32(n) {
			t.Fatalf("max_inodes %d: %v (%d)", n, err, cfg.MaxInodes)
		}
	}
}
