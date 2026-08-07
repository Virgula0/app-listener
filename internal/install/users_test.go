package install

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParsePasswd keeps root (UID 0) and the real users, filters other
// system users, users without a home directory and users whose home does
// not exist, and fails on malformed lines.
func TestParsePasswd(t *testing.T) {
	home := t.TempDir()
	realHome := filepath.Join(home, "alice")
	if err := os.MkdirAll(realHome, 0o700); err != nil {
		t.Fatal(err)
	}
	rootHome := filepath.Join(home, "root")
	if err := os.MkdirAll(rootHome, 0o700); err != nil {
		t.Fatal(err)
	}

	data := []byte(
		"root:x:0:0:root:" + rootHome + ":/bin/bash\n" +
			"daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin\n" +
			"alice:x:1000:1000:Alice:" + realHome + ":/bin/bash\n" +
			"bob:x:1001:1001:Bob:/home/bob:/bin/bash\n" + // home does not exist
			"nohome:x:1002:1002:No Home::/bin/bash\n",
	)
	users, err := parsePasswd(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("got %d users, want 2 (root and alice): %+v", len(users), users)
	}
	if users[0].Name != "root" || users[0].UID != 0 || users[0].GID != 0 || users[0].Home != rootHome {
		t.Errorf("unexpected root entry: %+v", users[0])
	}
	if users[1].Name != "alice" || users[1].UID != 1000 || users[1].GID != 1000 || users[1].Home != realHome {
		t.Errorf("unexpected alice entry: %+v", users[1])
	}
}

// TestParsePasswdMalformed rejects garbage lines.
func TestParsePasswdMalformed(t *testing.T) {
	if _, err := parsePasswd([]byte("not-a-passwd-line\n")); err == nil {
		t.Fatal("expected error for malformed line")
	}
	if _, err := parsePasswd([]byte("alice:x:notanumber:1000:Alice:/home/alice:/bin/bash\n")); err == nil {
		t.Fatal("expected error for malformed UID")
	}
	if _, err := parsePasswd([]byte("alice:x:1000:notanumber:Alice:/home/alice:/bin/bash\n")); err == nil {
		t.Fatal("expected error for malformed GID")
	}
}
