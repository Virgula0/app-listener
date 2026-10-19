package recorder

import (
	"os"
	"regexp"
	"slices"
	"testing"
)

// hookCoverage maps every LSM hook of guard.bpf.c to the monitor probes that report what it judges
// to a recording, and the recordops ops (integrationtests/exploits/recordops.c) that prove it in
// TestRecord_GuardParity. A hook with no probes says why nothing can be learned for it.
var hookCoverage = map[string]struct {
	probes []string
	ops    []string
	none   string
}{
	"file_open":         {probes: []string{"vfs_open"}, ops: []string{"open"}},
	"file_permission":   {probes: []string{"vfs_read", "vfs_readv", "security_file_permission"}, ops: []string{"read", "readv", "pread", "readdir", "writev", "pwritev", "splice_in", "splice_out", "sendfile_in", "sendfile_out", "cfr_in", "cfr_out"}},
	"mmap_file":         {probes: []string{"security_mmap_file"}, ops: []string{"mmap_read", "mmap_write"}},
	"file_truncate":     {probes: []string{"security_file_truncate"}, ops: []string{"ftruncate", "otrunc"}},
	"path_truncate":     {probes: []string{"notify_change"}, ops: []string{"truncate"}},
	"inode_setattr":     {probes: []string{"notify_change"}, ops: []string{"chmod", "chown", "utimes"}},
	"inode_setxattr":    {probes: []string{"vfs_setxattr"}, ops: []string{"setxattr"}},
	"inode_removexattr": {probes: []string{"vfs_removexattr"}, ops: []string{"removexattr"}},
	"inode_getxattr":    {probes: []string{"security_inode_getxattr"}, ops: []string{"getxattr"}},
	"inode_listxattr":   {probes: []string{"security_inode_listxattr"}, ops: []string{"listxattr"}},
	"inode_getattr":     {probes: []string{"security_inode_getattr"}, ops: []string{"stat"}},
	"inode_readlink":    {probes: []string{"vfs_readlink"}, ops: []string{"readlink"}},
	"path_mknod":        {probes: []string{"security_path_mknod"}, ops: []string{"create", "mknod_reg", "mkfifo"}},
	"path_mkdir":        {probes: []string{"vfs_mkdir"}, ops: []string{"mkdir"}},
	"path_rmdir":        {probes: []string{"vfs_rmdir"}, ops: []string{"rmdir"}},
	"path_unlink":       {probes: []string{"vfs_unlink"}, ops: []string{"unlink"}},
	"path_rename":       {probes: []string{"vfs_rename"}, ops: []string{"rename", "rename_out", "rename_in"}},
	"path_link":         {probes: []string{"vfs_link"}, ops: []string{"link", "link_out", "link_in"}},
	"path_symlink":      {probes: []string{"vfs_symlink"}, ops: []string{"symlink"}},

	"bprm_check_security":  {none: "exec of an image in the tree: judged as the image's OPEN, which its own opens and mappings record"},
	"file_ioctl":           {none: "btrfs copy ioctls, refused superblock-wide: never proposed"},
	"file_ioctl_compat":    {none: "btrfs copy ioctls (compat), refused superblock-wide: never proposed"},
	"sb_mount":             {none: "mounting over a tree needs CAP_SYS_ADMIN: never proposed"},
	"ptrace_access_check":  {none: "process gate on tainted processes, not a file op"},
	"task_prctl":           {none: "process gate: denies prctl(PR_SET_MM) exe_file rewrites, not a file op"},
	"inode_free_security":  {none: "bookkeeping (evicts freed inodes)"},
	"bprm_committed_creds": {none: "bookkeeping (supersede and taint on exec)"},
	"task_free":            {none: "bookkeeping (exited processes)"},
	"inode_unlink":         {none: "bookkeeping (superseded binary keys)"},
	"inode_rename":         {none: "bookkeeping (superseded binary keys)"},
}

func readSource(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(b)
}

// TestGuardHookCoverage: a guard hook the recorder can't see makes --record propose a mask the guard
// then denies the program under. Adding a hook means adding its monitor probe (attached by
// AttachRecordProbes or NewMonitor) and a recordops op, or saying here why it can't be learned.
func TestGuardHookCoverage(t *testing.T) {
	guardSrc := readSource(t, "../guard/bpf/guard.bpf.c")
	monitorSrc := readSource(t, "../monitor/bpf/monitor.bpf.c")
	monitorGo := readSource(t, "../monitor/monitor.go")
	opsSrc := readSource(t, "../../integrationtests/exploits/recordops.c")
	parityTest := readSource(t, "../../integrationtests/record_parity_test.go")

	var hooks []string
	for _, m := range regexp.MustCompile(`SEC\("lsm/([a-z_0-9]+)"\)`).FindAllStringSubmatch(guardSrc, -1) {
		hooks = append(hooks, m[1])
	}
	if len(hooks) == 0 {
		t.Fatal("no lsm/ hooks found in guard.bpf.c")
	}
	for _, h := range hooks {
		c, ok := hookCoverage[h]
		switch {
		case !ok:
			t.Errorf("guard hook %s has no recorder coverage entry", h)
		case len(c.probes) == 0 && c.none == "":
			t.Errorf("guard hook %s: no probes and no reason", h)
		case len(c.probes) > 0 && len(c.ops) == 0:
			t.Errorf("guard hook %s: no recordops op proves its probes", h)
		}
		for _, p := range c.probes {
			if !regexp.MustCompile(`SEC\("kprobe/` + p + `"\)`).MatchString(monitorSrc) {
				t.Errorf("%s: monitor.bpf.c has no kprobe/%s", h, p)
			}
			if !regexp.MustCompile(`"` + p + `"\}`).MatchString(monitorGo) {
				t.Errorf("%s: monitor.go never attaches %s", h, p)
			}
		}
		for _, op := range c.ops {
			if !regexp.MustCompile(`strcmp\(op, "` + op + `"\)`).MatchString(opsSrc) {
				t.Errorf("%s: recordops.c has no op %q", h, op)
			}
			if !regexp.MustCompile(`\{"` + op + `", "[A-Z]+"\}`).MatchString(parityTest) {
				t.Errorf("%s: TestRecord_GuardParity does not run op %q", h, op)
			}
		}
	}
	for h := range hookCoverage {
		if !slices.Contains(hooks, h) {
			t.Errorf("coverage entry %s names no guard hook", h)
		}
	}
}
