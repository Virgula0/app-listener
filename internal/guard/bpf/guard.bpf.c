#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <errno.h>

#define MAX_PATH 256
// VFS access-mask bits for bpf_lsm_inode_permission/file_permission. This kernel (7.0-hardened)
// differs from mainline: MAY_OPEN is 0x20 (not 0x40), MAY_ACCESS 0x10. Keep in sync with the target
// kernel's include/linux/fs.h.
#define MAY_EXEC   0x00000001
#define MAY_WRITE  0x00000002
#define MAY_READ   0x00000004
#define MAY_APPEND 0x00000008
#define MAY_ACCESS 0x00000010
#define MAY_OPEN   0x00000020
#define MAY_CHDIR  0x00000040
#define ATTR_SIZE 0x00000008
#define ATTR_FILE 0x00002000
#define GUARD_BLOCK 1
#define GUARD_ALLOW 2
// Root-gated allow: honored only for uid 0. For the owner's own binary (the daemon), so it isn't a
// universal key any local user could execute.
#define GUARD_ALLOW_ROOT 3
#define __FMODE_EXEC 0x20 // set in file->f_flags by kernel exec opens (do_open_execat)
#define FMODE_WRITE 0x2   // file->f_mode: file is open for write (stable across kernels)
#define PROT_WRITE 0x2    // mmap prot
#define MAP_SHARED 0x1    // mmap flags

// guard_config[0] modes: 0/1 = blacklist/whitelist. GUARD_MODE_READONLY (2): every process may
// READ, only whitelisted binaries may modify (write, truncate, rename, unlink, chmod, mkdir, mknod,
// xattr, writable mmap, mount-over). For /etc/app-listener (world-readable daemon.conf).
#define GUARD_MODE_BLACKLIST 0
#define GUARD_MODE_WHITELIST 1
#define GUARD_MODE_READONLY 2

// guard_event.reason: GUARD_REASON_NONE = ordinary path-keyed decision. GUARD_REASON_RAW_DEVICE =
// raw block-device gate denial (is_guarded_block_device): the target is the backing device, not a
// watched path; userspace must NOT attribute it to a resource.
#define GUARD_REASON_NONE 0
#define GUARD_REASON_RAW_DEVICE 1
// Process-gate denials: no file. `path` = the OTHER task's comm, `fd` = its tgid (tainted ptrace
// target, or unwhitelisted exec tracer).
#define GUARD_REASON_PTRACE 2
#define GUARD_REASON_TRACED_EXEC 3
#define GUARD_REASON_PROC_MEM 4

enum event_type {
	EVENT_OPEN,
	EVENT_READ,
	EVENT_WRITE,
	EVENT_DELETE,
	EVENT_RENAME,
	EVENT_SYMLINK,
	EVENT_HARDLINK,
	EVENT_MKDIR,
	EVENT_MMAP,
	EVENT_ATTR,
	EVENT_STAT,
	EVENT_MKNOD,
};

struct guard_event {
	__u32 pid;
	__u32 uid;
	__u32 gid;
	__u32 type;
	__u32 fd;
	__u32 blocked;
	__u32 reason; // GUARD_REASON_* — how userspace should attribute the event
	char comm[16];
	char path[MAX_PATH];
	char dest[MAX_PATH];
};

struct inode_key {
	__u64 dev;
	__u64 ino;
};

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 5);  // [0]=mode, [1]=recursive, [2]=depth,
	                         // [3]=watch root dev, [4]=watch root ino
	__type(key, __u32);
	__type(value, __u64);
} guard_config SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 131072);
	__type(key, struct inode_key);
	__type(value, __u8);
} guard_inodes SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 256);
	__type(key, struct inode_key);
	__type(value, __u8);
} guard_exe_actions SEC(".maps");

// Per-binary allowed-event bitmask (bit i = enum event_type i). Whitelist mode only; a missing
// entry = all events allowed (plain guard mode never populates this map).
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 64);
	__type(key, struct inode_key);
	__type(value, __u32);
} guard_exe_events SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, char[MAX_PATH]);
} guard_path SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 8);
	__type(key, __u64);  // filesystem device (sb->s_dev) of the guarded path
	__type(value, __u8);
} guard_fs_sbdevs SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 64);
	__type(key, __u32);
	__type(value, __u8);
} guard_fs_devices SEC(".maps"); // block devices hosting guarded filesystems

// Degradation counters: observability only, never consulted by a decision.
//   [0] guard_inodes full: runtime discovery dropped an add (falls back to the ancestor walk)
//   [1] guard_tainted_pids full: a taint stamp was lost (future ptrace/process_vm_readv on it isn't denied)
struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 2);
	__type(key, __u32);
	__type(value, __u64);
} guard_degrade SEC(".maps");

static __always_inline void count_degrade(__u32 slot)
{
	__u64 one = 1;
	__u64 zero = 0;
	__u64 *cell = bpf_map_lookup_elem(&guard_degrade, &slot);
	if (!cell) {
		bpf_map_update_elem(&guard_degrade, &slot, &zero, BPF_NOEXIST);
		return;
	}
	__sync_fetch_and_add(cell, one);
}

// Processes holding guarded file content in memory. Once tainted, process_vm_readv, ptrace and
// /proc/<pid>/mem against them are blocked unless the caller is whitelisted.
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 4096);
	__type(key, __u32);   // PID (TGID) of the tainted process
	__type(value, __u8);  // 1 = tainted
} guard_tainted_pids SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, char[MAX_PATH]);
} tmp_buf SEC(".maps"); // only used by guard_path_symlink for oldname_buf

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	// 2 MiB (~3.5k events). Enforcement is synchronous, so overflow drops only telemetry, never a
	// decision. Per guard, so size matters with dozens of guards.
	__uint(max_entries, 1 << 21);
} rb SEC(".maps");

char LICENSE[] SEC("license") = "GPL";

// Read the exe file inode of the current process for anti-spoof verification.
static __always_inline int get_current_exe_inode(struct inode_key *ik)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (!task)
		return 0;

	struct mm_struct *mm;
	bpf_probe_read_kernel(&mm, sizeof(mm), &task->mm);
	if (!mm)
		return 0;

	struct file *exe_file;
	bpf_probe_read_kernel(&exe_file, sizeof(exe_file), &mm->exe_file);
	if (!exe_file)
		return 0;
	struct inode *exe_inode;
	bpf_probe_read_kernel(&exe_inode, sizeof(exe_inode), &exe_file->f_inode);
	if (!exe_inode)
		return 0;

	bpf_probe_read_kernel(&ik->ino, sizeof(ik->ino), &exe_inode->i_ino);

	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &exe_inode->i_sb);
	if (!sb)
		return 0;

	dev_t dev;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	ik->dev = dev;

	return 1;
}

// get_task_exe_inode resolves the exe identity of an ARBITRARY task (protects whitelisted victims
// in ptrace_access_check). 0 for kernel threads (no mm, no exe identity).
static __always_inline int get_task_exe_inode(struct task_struct *task, struct inode_key *ik)
{
	if (!task)
		return 0;

	struct mm_struct *mm;
	bpf_probe_read_kernel(&mm, sizeof(mm), &task->mm);
	if (!mm)
		return 0;

	struct file *exe_file;
	bpf_probe_read_kernel(&exe_file, sizeof(exe_file), &mm->exe_file);
	if (!exe_file)
		return 0;

	struct inode *exe_inode;
	bpf_probe_read_kernel(&exe_inode, sizeof(exe_inode), &exe_file->f_inode);
	if (!exe_inode)
		return 0;

	bpf_probe_read_kernel(&ik->ino, sizeof(ik->ino), &exe_inode->i_ino);

	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &exe_inode->i_sb);
	if (!sb)
		return 0;

	dev_t dev;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	ik->dev = dev;

	return 1;
}

static __always_inline int read_inode_guard(struct inode *inode)
{
	if (!inode)
		return 0;

	struct inode_key ikey = {};
	bpf_probe_read_kernel(&ikey.ino, sizeof(ikey.ino), &inode->i_ino);

	dev_t dev;
	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	ikey.dev = dev;

	__u8 *found = bpf_map_lookup_elem(&guard_inodes, &ikey);
	return found != NULL ? 1 : 0;
}

static __always_inline struct inode *get_inode_from_path(void *ctx_path)
{
	struct path *p = (struct path *)ctx_path;
	if (!p) return NULL;
	struct dentry *d;
	bpf_probe_read_kernel(&d, sizeof(d), &p->dentry);
	if (!d) return NULL;
	struct inode *i;
	bpf_probe_read_kernel(&i, sizeof(i), &d->d_inode);
	return i;
}

static __always_inline struct dentry *get_dentry_from_path(void *ctx_path)
{
	struct path *p = (struct path *)ctx_path;
	if (!p)
		return NULL;
	struct dentry *d;
	bpf_probe_read_kernel(&d, sizeof(d), &p->dentry);
	return d;
}

// Auto-add a newly created directory to guard_inodes per the guard's recursive setting.
static __always_inline int should_add_new_dir(void)
{
	__u32 key = 1;
	__u64 *recursive = bpf_map_lookup_elem(&guard_config, &key);
	if (!recursive || *recursive == 0)
		return 0;
	return 1;
}

// is_readonly_mode: callers skip taint-tracking on allowed reads in GUARD_MODE_READONLY
// (world-readable content, not secrets; tainting readers would bloat guard_tainted_pids and
// over-restrict ptrace).
static __always_inline int is_readonly_mode(void)
{
	__u32 key = 0;
	__u64 *mode = bpf_map_lookup_elem(&guard_config, &key);
	return mode != NULL && *mode == GUARD_MODE_READONLY;
}

// event_is_read_class: in GUARD_MODE_READONLY, events that leak no modification are allowed for
// everyone. OPEN and MMAP count only without write intent (FMODE_WRITE / MAP_SHARED|PROT_WRITE,
// passed as write_intent to check_and_emit_ex).
static __always_inline int event_is_read_class(__u32 type, bool write_intent)
{
	switch (type) {
	case EVENT_READ:
	case EVENT_STAT:
		return 1;
	case EVENT_OPEN:
	case EVENT_MMAP:
		return !write_intent;
	default:
		return 0;
	}
}

// Walk up from `parent` and report whether its subtree lies in the guarded region: rooted at the
// farthest inode present in guard_inodes, bounded by the depth limit (a parent exactly at the limit
// is NOT guarded, mirroring the eager scan and discover_guarded_parent). Protects the whole subtree
// at any depth: covers the window between hook attach and the eager scan, deep files, and a full
// inode map.
//
// Each step (advance to parent, test guard membership and watch-root reachability) is written
// inline in the loop, not as a helper: a ref-parameter helper plus per-iteration early exits kept
// the loop head from converging and blew the 1M-insn verifier budget even at 16 steps.
static __always_inline int guarded_ancestor_within_limit(struct dentry *parent)
{
	if (!parent)
		return 0;

	// Non-recursive guards map only the root and its direct children; deeper levels are
	// deliberately unguarded.
	if (!should_add_new_dir())
		return 0;

	struct inode *parent_inode;
	bpf_probe_read_kernel(&parent_inode, sizeof(parent_inode), &parent->d_inode);
	if (!parent_inode)
		return 0;

	// Perf gate: only files on a filesystem hosting a guarded root can have a guarded ancestor.
	// guard_fs_sbdevs (populated unconditionally by Go, incl. anon tmpfs/overlayfs devs) lets other
	// filesystems skip the walk with one lookup.
	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &parent_inode->i_sb);
	if (sb) {
		__u64 dev = 0;
		bpf_probe_read_kernel(&dev, sizeof(dev_t), &sb->s_dev);
		if (!bpf_map_lookup_elem(&guard_fs_sbdevs, &dev))
			return 0;
	}

	__u32 depth_key = 2;
	__u64 *depth = bpf_map_lookup_elem(&guard_config, &depth_key);
	bool has_limit = depth != NULL && *depth > 0;

	// Root confinement: a guarded ancestor counts only if the chain also reaches the watch root
	// inode (guard_config[4]); an inode number reused outside the tree has no root in its chain.
	// The device is implied by the gate above (the root shares the superblock of every in-tree
	// parent). Fails closed if the root isn't configured (legacy pairing).
	bool root_configured = false;
	__u64 root_ino = 0;
	{
		__u32 rik = 4;
		__u64 *ri = bpf_map_lookup_elem(&guard_config, &rik);
		if (ri) {
			root_configured = true;
			root_ino = *ri;
		}
	}

	// Walk up to 16 levels: the largest bound the verifier accepts for this loop body (20+ steps
	// exceed the 1M-insn budget). Deeper coverage comes from lazy discovery at mkdir/open
	// (discover_guarded_parent maps every runtime-created chain a level at a time), so this loop
	// only bridges the last hops. With a depth limit: a guarded ancestor at/beyond the limit means
	// NOT guarded (a farther one is only more distant); finding none decides the same.
	struct dentry *ancestor = parent;
	__u64 farthest_steps = 0;
	bool found_guarded = false;
	bool rooted = !root_configured;

	for (int i = 0; i < 16; i++) {
		struct dentry *next;
		bpf_probe_read_kernel(&next, sizeof(next), &ancestor->d_parent);
		if (!next || next == ancestor)
			break;
		ancestor = next;

		struct inode *anc_inode;
		bpf_probe_read_kernel(&anc_inode, sizeof(anc_inode), &ancestor->d_inode);
		if (!anc_inode)
			break;
		if (anc_inode == parent_inode)
			break;

		__u64 anc_ino = 0;
		bpf_probe_read_kernel(&anc_ino, sizeof(anc_ino), &anc_inode->i_ino);
		if (root_configured && anc_ino == root_ino)
			rooted = true;

		// Depth decision once after the loop: an ancestor at/beyond the limit = NOT guarded (same
		// as lazy discovery).
		if (read_inode_guard(anc_inode)) {
			farthest_steps = i + 1;
			found_guarded = true;
		}
	}

	if (!found_guarded || !rooted)
		return 0;
	if (has_limit && farthest_steps >= *depth)
		return 0;
	return 1;
}

// Checks the file's own inode and its parent directory's.
#define S_IFMT  00170000
#define S_IFBLK 0060000
#define S_IFDIR 0040000

// Is the inode a block device whose dev_t matches a guarded filesystem (debugfs/dd raw reads bypass
// VFS)?
static __always_inline int is_guarded_block_device(struct inode *inode)
{
	if (!inode)
		return 0;

	umode_t mode;
	bpf_probe_read_kernel(&mode, sizeof(mode), &inode->i_mode);

	// Check if this is a block device (S_ISBLK)
	if ((mode & S_IFMT) != S_IFBLK)
		return 0;

	__u32 rdev;
	bpf_probe_read_kernel(&rdev, sizeof(rdev), &inode->i_rdev);

	__u8 *val = bpf_map_lookup_elem(&guard_fs_devices, &rdev);
	return val != NULL;
}

// inode_is_watch_root: exact (dev, ino) match with guard_config[3]/[4], no dentry walk (the same
// equality root_in_chain treats as decisive). Used by path_rename to deny clobbering the watch
// root, which the destination-parent check can't see for a single-file watch (its parent isn't in
// guard_inodes).
static __always_inline int inode_is_watch_root(struct inode *inode)
{
	if (!inode)
		return 0;

	__u32 dev_key = 3;
	__u32 ino_key = 4;
	__u64 *root_dev = bpf_map_lookup_elem(&guard_config, &dev_key);
	__u64 *root_ino = bpf_map_lookup_elem(&guard_config, &ino_key);
	if (!root_dev || !root_ino)
		return 0;

	__u64 ino = 0;
	bpf_probe_read_kernel(&ino, sizeof(ino), &inode->i_ino);
	if (ino != *root_ino)
		return 0;

	struct super_block *sb = NULL;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	dev_t dev = 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	return (__u64)dev == *root_dev;
}

// Mark the current process tainted (it holds guarded content in memory): non-whitelisted
// ptrace/process_vm_readv against it is blocked.
static __always_inline void mark_tainted(void)
{
	__u32 pid = bpf_get_current_pid_tgid() >> 32;
	__u8 val = 1;
	if (bpf_map_update_elem(&guard_tainted_pids, &pid, &val, BPF_ANY))
		count_degrade(1);
}

// Add an inode to guard_inodes (runtime discovery from inode_mkdir/path_rename).
static __always_inline void add_inode_to_guard(struct inode *inode)
{
	if (!inode)
		return;

	struct inode_key ikey = {};
	bpf_probe_read_kernel(&ikey.ino, sizeof(ikey.ino), &inode->i_ino);

	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return;

	dev_t dev;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	ikey.dev = dev;

	__u8 v = 1;
	if (bpf_map_update_elem(&guard_inodes, &ikey, &v, BPF_ANY))
		count_degrade(0);
}

// Evict an inode from guard_inodes when path_unlink/path_rmdir lets a delete through: the kernel
// frees the inode right after, and its number can go to an unrelated file. Evicting now (not
// waiting for ReconcileInodes, up to inodeGCEvery ~1h) closes that window. A miss (never a member)
// is silent.
static __always_inline void evict_inode_from_guard(struct inode *inode)
{
	if (!inode)
		return;

	struct inode_key ikey = {};
	bpf_probe_read_kernel(&ikey.ino, sizeof(ikey.ino), &inode->i_ino);

	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return;

	dev_t dev;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	ikey.dev = dev;

	bpf_map_delete_elem(&guard_inodes, &ikey);
}

// Returns the tainted target pid if dentry is /proc/<pid>/mem of one, else 0. Needed explicitly:
// /proc inodes aren't in guard_inodes.
static __always_inline __u32 is_proc_mem_of_tainted(struct dentry *dentry)
{
	if (!dentry)
		return 0;

	// Check if this dentry's name is "mem"
	const unsigned char *name_ptr;
	bpf_probe_read_kernel(&name_ptr, sizeof(name_ptr), &dentry->d_name.name);
	if (!name_ptr)
		return 0;

	char name[8];
	long ret = bpf_probe_read_kernel_str(name, sizeof(name), name_ptr);
	if (ret <= 0)
		return 0;

	// Must be exactly "mem"
	if (name[0] != 'm' || name[1] != 'e' || name[2] != 'm' || name[3] != '\0')
		return 0;

	// Check parent dentry (the PID directory)
	struct dentry *parent;
	bpf_probe_read_kernel(&parent, sizeof(parent), &dentry->d_parent);
	if (!parent || parent == dentry)
		return 0;

	const unsigned char *pname_ptr;
	bpf_probe_read_kernel(&pname_ptr, sizeof(pname_ptr), &parent->d_name.name);
	if (!pname_ptr)
		return 0;

	char pname[16];
	ret = bpf_probe_read_kernel_str(pname, sizeof(pname), pname_ptr);
	if (ret <= 0)
		return 0;

	// Parse parent name as a numeric PID
	__u32 pid = 0;
	for (int i = 0; i < 12; i++) {
		char c = pname[i];
		if (c >= '0' && c <= '9') {
			pid = pid * 10 + (c - '0');
		} else if (c == '\0') {
			break;
		} else {
			return 0;  // not a numeric directory
		}
	}
	if (pid == 0)
		return 0;

	__u8 *val = bpf_map_lookup_elem(&guard_tainted_pids, &pid);
	return val != NULL ? pid : 0;
}

// ptrace_may_access() mode bits (include/linux/ptrace.h).
#define PTRACE_MODE_READ 0x01
#define PTRACE_MODE_ATTACH 0x02

// emit_process_denial reports a process-gate denial (formerly a silent -EPERM, undiagnosable for
// e.g. browser engines or Wine). `other` is the other task if at hand, else other_tgid names it.
static __always_inline int emit_process_denial(__u32 reason, __u32 type, struct task_struct *other, __u32 other_tgid, __u32 mode)
{
	struct guard_event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (e) {
		e->pid = bpf_get_current_pid_tgid() >> 32;
		e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
		e->gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
		e->type = type;
		e->blocked = 1;
		e->reason = reason;
		bpf_get_current_comm(e->comm, sizeof(e->comm));
		e->path[0] = '\0';
		if (other) {
			bpf_probe_read_kernel(&other_tgid, sizeof(other_tgid), &other->tgid);
			bpf_probe_read_kernel_str(e->path, 16, &other->comm);
		}
		e->fd = other_tgid;
		// dest = access mode: ATTACH reaches memory (ptrace, process_vm_readv, /proc/<pid>/mem);
		// READ only metadata (environ, fd, maps).
		if (mode & PTRACE_MODE_ATTACH) {
			e->dest[0] = 'A'; e->dest[1] = 'T'; e->dest[2] = 'T';
			e->dest[3] = 'A'; e->dest[4] = 'C'; e->dest[5] = 'H'; e->dest[6] = '\0';
		} else if (mode & PTRACE_MODE_READ) {
			e->dest[0] = 'R'; e->dest[1] = 'E'; e->dest[2] = 'A';
			e->dest[3] = 'D'; e->dest[4] = '\0';
		} else {
			e->dest[0] = '\0';
		}
		bpf_ringbuf_submit(e, 0);
	}
	return -EPERM;
}


// /proc/<pid>/fd/<n>: opening it goes through file_open with the real file's dentry, so
// is_guarded_access already catches it; only /proc/<pid>/mem needs an explicit check.

// True when the accessed file's dentry chain (own inode included) reaches the watch root
// (guard_config[3..4]). A map entry only guards an access when the chain is rooted: a stale entry
// whose inode number was reused elsewhere (ext4) has no root in its chain and must not deny.
//
// Fails closed in both ambiguous cases (root not configured; bound exhausted before the fs root):
// the access is treated as rooted.
//
// The device is read once from the file's superblock; the watch root shares it, so a file on
// another device can't have the root in its chain and the check is skipped. Per-ancestor probes are
// inode-number only (probing every ancestor's superblock blew the verifier budget).
static __always_inline int root_in_chain(struct dentry *dentry, struct inode *file_inode, int bound)
{
	__u32 dev_key = 3;
	__u32 ino_key = 4;
	__u64 *root_dev = bpf_map_lookup_elem(&guard_config, &dev_key);
	__u64 *root_ino = bpf_map_lookup_elem(&guard_config, &ino_key);
	if (!root_dev || !root_ino)
		return 1;

	__u64 want_dev = *root_dev;
	__u64 want_ino = *root_ino;

	bool on_dev = false;
	__u64 ino = 0;
	if (file_inode) {
		bpf_probe_read_kernel(&ino, sizeof(ino), &file_inode->i_ino);
		struct super_block *sb = NULL;
		bpf_probe_read_kernel(&sb, sizeof(sb), &file_inode->i_sb);
		if (sb) {
			dev_t dev = 0;
			bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
			if ((__u64)dev == want_dev && ino == want_ino)
				return 1;
			on_dev = ((__u64)dev == want_dev);
		}
	}
	// A chain from a file on a different device can never contain the watch root.
	if (file_inode && !on_dev)
		return 0;
	if (!dentry)
		return 0;

	// Walk up the chain (the file's own inode was checked above). All ancestors share the file's
	// device, so only inodes are probed. Bound per hook: 32 for single-walk hooks (matching
	// file_open), 16 for the verifier-tight two-walk hooks (unlink, link, rename); both exceed real
	// path depth (deepest seen ~12). If exhausted, the final ancestor's parent is probed for the fs
	// root, so a chain exactly at the bound still resolves; only deeper chains fail closed.
	struct dentry *ancestor = dentry;
	for (int i = 0; i < bound; i++) {
		struct dentry *next;
		bpf_probe_read_kernel(&next, sizeof(next), &ancestor->d_parent);
		if (!next || next == ancestor)
			return 0;  // filesystem root reached, root inode not seen
		ancestor = next;

		struct inode *anc_inode;
		bpf_probe_read_kernel(&anc_inode, sizeof(anc_inode), &ancestor->d_inode);
		if (!anc_inode)
			return 0;
		if (anc_inode == file_inode)
			return 0;

		__u64 anc_ino = 0;
		bpf_probe_read_kernel(&anc_ino, sizeof(anc_ino), &anc_inode->i_ino);
		if (anc_ino == want_ino)
			return 1;
	}

	// A chain exactly `bound` deep lands on the fs root at the last iteration; recognizing it
	// (parent NULL/itself) needs one more hop, done here. Such a root ends the walk without the
	// watch root: not rooted.
	struct dentry *final_next;
	bpf_probe_read_kernel(&final_next, sizeof(final_next), &ancestor->d_parent);
	if (!final_next || final_next == ancestor)
		return 0;

	// Bound exhausted without reaching the filesystem root: fail closed.
	return 1;
}

// Shared gate for direct map-hit denies in path_*/inode_* hooks: an entry only denies when the
// chain reaches the watch root (same stale-entry containment as is_guarded_access, for inode
// numbers reused outside the tree).
static __always_inline int guarded_map_hit(struct dentry *dentry, struct inode *inode, int bound)
{
	if (!dentry || !read_inode_guard(inode))
		return 0;
	return root_in_chain(dentry, inode, bound);
}

static __always_inline int is_guarded_access(struct dentry *dentry, struct inode *file_inode)
{
	if (read_inode_guard(file_inode)) {
		// Own inode in the map. A stale entry (number reused outside the tree) must not deny, so
		// confine to the watch root. An unrooted chain is outside the tree and none of its
		// ancestors can be genuinely guarded either, so this is decisive.
		if (root_in_chain(dentry, file_inode, 32))
			return 1;
		return 0;
	}

	if (!dentry)
		return 0;

	struct dentry *parent;
	bpf_probe_read_kernel(&parent, sizeof(parent), &dentry->d_parent);
	if (!parent || parent == dentry)
		return 0;

	struct inode *parent_inode;
	bpf_probe_read_kernel(&parent_inode, sizeof(parent_inode), &parent->d_inode);
	if (parent_inode && parent_inode != file_inode && read_inode_guard(parent_inode)) {
		// Same stale-entry rule for the direct parent: only a rooted chain counts.
		if (root_in_chain(dentry, file_inode, 32))
			return 1;
		return 0;
	}

	// Neither the file nor its parent is in the map (runtime-created dir, file deeper than the
	// eager scan, or map filled mid-scan): walk the dentry chain; inside the guarded region =
	// guarded. guarded_ancestor_within_limit applies root confinement itself.
	return guarded_ancestor_within_limit(parent);
}

static __always_inline long read_path(struct dentry *dentry, char *buf, int buf_size)
{
    char tmp[64];
    int pos = MAX_PATH;

    pos--;
    buf[pos & (MAX_PATH - 1)] = '\0';

    struct dentry *d = dentry;

    for (int i = 0; i < 32; i++) {
        if (!d)
            break;

        struct dentry *parent;
        bpf_probe_read_kernel(&parent, sizeof(parent), &d->d_parent);
        if (!parent || parent == d)
            break;

        const unsigned char *name_ptr;
        bpf_probe_read_kernel(&name_ptr, sizeof(name_ptr), &d->d_name.name);
        if (!name_ptr)
            break;

        long ret = bpf_probe_read_kernel_str(tmp, sizeof(tmp), name_ptr);
        if (ret <= 1)
            break;

        int name_len = ret - 1;

        if (pos < name_len)
            break;

        pos -= name_len;

        // Cache the masked position to help the verifier understand the baseline
        int start_pos = pos & (MAX_PATH - 1);

        // Removed #pragma unroll here as well
        for (int j = 0; j < 64; j++) {
            if (j >= name_len) break;
            buf[(start_pos + j) & (MAX_PATH - 1)] = tmp[j & 63];
        }

        if (pos <= 0)
            break;

        pos--;
        buf[pos & (MAX_PATH - 1)] = '/';

        d = parent;
    }

    // Return the offset directly (no bounds check needed here, already bounded)
    return pos;
}

static void fill_path(struct dentry *dentry, char *out)
{
    if (!dentry) {
        out[0] = '\0';
        return;
    }
    __u32 key = 0;
    char *buf = bpf_map_lookup_elem(&tmp_buf, &key);
    if (!buf) {
        out[0] = '\0';
        return;
    }

    long off = read_path(dentry, buf, MAX_PATH);

    // 'off' is guaranteed to be >= 0 based on our new read_path implementation
    if (off < MAX_PATH) {
        off &= (MAX_PATH - 1);
        bpf_probe_read_kernel_str(out, MAX_PATH, buf + off);
    } else {
        out[0] = '\0';
    }
}

// is_allow_action: GUARD_ALLOW always allows; GUARD_ALLOW_ROOT only for uid 0 (the guard owner).
// NULL denies in whitelist mode.
static __always_inline int is_allow_action(const __u8 *action)
{
	if (!action)
		return 0;
	if (*action == GUARD_ALLOW)
		return 1;
	if (*action == GUARD_ALLOW_ROOT)
		return (bpf_get_current_uid_gid() & 0xffffffff) == 0;
	return 0;
}

static __always_inline int check_and_emit_ex(__u32 type, struct dentry *dentry, const char *dest_str, bool dest_is_user, struct dentry *dest_dentry, bool is_exec_open, bool quiet_allow, bool write_intent, __u32 reason)
{
	__u32 key = 0;
	__u64 *mode = bpf_map_lookup_elem(&guard_config, &key);
	if (!mode)
		return -EPERM;

	// Save mode value early to avoid verifier pointer-tracking concerns
	__u64 mode_val = *mode;

	// Kernel threads (no userspace mm) act for user processes (io_uring workers, SQPOLL) and have
	// no exe inode: block unconditionally.
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	struct mm_struct *current_mm = NULL;
	if (task)
		bpf_probe_read_kernel(&current_mm, sizeof(current_mm), &task->mm);

	int is_blocked;
	if (!current_mm) {
		is_blocked = 1;
	} else {
		struct inode_key exe_ik = {};
		__u8 *action = NULL;
		if (get_current_exe_inode(&exe_ik))
			action = bpf_map_lookup_elem(&guard_exe_actions, &exe_ik);

		if (mode_val == GUARD_MODE_BLACKLIST) {
			is_blocked = action != NULL && *action == GUARD_BLOCK;
		} else if (mode_val == GUARD_MODE_READONLY) {
			// Read-only guard: every process may read (allow not logged: a config file is read
			// constantly); modifying ops are gated on the whitelist as in whitelist mode.
			if (event_is_read_class(type, write_intent)) {
				is_blocked = 0;
				quiet_allow = true;
			} else {
				is_blocked = !is_allow_action(action);
			}
		} else {
			is_blocked = !is_allow_action(action);

			// Exec-open attribution (whitelist mode only): executing a binary is an OPEN performed
			// by the *launcher* (a shell wrapper runs via /bin/sh, whose exe inode isn't
			// whitelisted), so attribute exec-opens to the *target* and allowlist that. The exec fd
			// is never exposed to the caller (execve destroys its address space) and afterwards the
			// process IS the target, so later accesses are checked normally. Blacklist mode keeps
			// caller attribution.
			//
			// Two discriminators cover the exec chain, both attributed by the *accessed file's*
			// inode (only a whitelisted file passes):
			//   - the target's own open carries __FMODE_EXEC in f_flags (task->in_execve is set
			//     only later, in bprm_execve);
			//   - the kernel's load-time accesses (prepare_binprm kernel_read, binfmt mmap, in-tree
			//     interpreter opens) run with task->in_execve set.
			if (is_blocked &&
			    ((type == EVENT_OPEN && is_exec_open) ||
			     BPF_CORE_READ_BITFIELD_PROBED(task, in_execve))) {
				struct inode *target;
				bpf_probe_read_kernel(&target, sizeof(target), &dentry->d_inode);
				if (target) {
					struct inode_key target_ik = {};
					__u8 *target_action;
					bpf_probe_read_kernel(&target_ik.ino, sizeof(target_ik.ino),
							      &target->i_ino);
					struct super_block *tsb;
					bpf_probe_read_kernel(&tsb, sizeof(tsb), &target->i_sb);
					if (tsb) {
						dev_t tdev;
						bpf_probe_read_kernel(&tdev, sizeof(tdev), &tsb->s_dev);
						target_ik.dev = tdev;
					}
					target_action = bpf_map_lookup_elem(&guard_exe_actions, &target_ik);
					if (target_action) {
						if (is_allow_action(target_action)) {
							exe_ik = target_ik;
							is_blocked = 0;
						}
					}
				}
			}

			// Per-binary event mask: only listed event types are permitted; no mask entry = all
			// events.
			if (!is_blocked) {
				__u32 *mask = bpf_map_lookup_elem(&guard_exe_events, &exe_ik);
				if (mask && type < 32 && !(*mask & (1U << type)))
					is_blocked = 1;
			}
		}
	}

	// quiet_allow: allowed accesses aren't reported (inode_permission fires on every guarded access
	// and would flood the log with allowed STATs). Denials are always reported.
	if (!is_blocked && quiet_allow)
		return 0;

	struct guard_event *e;
	e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (e) {
		e->pid = bpf_get_current_pid_tgid() >> 32;
		e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
		e->gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
		e->type = type;
		e->fd = 0;
		e->blocked = is_blocked ? 1 : 0;
		e->reason = reason;
		bpf_get_current_comm(e->comm, sizeof(e->comm));

		fill_path(dentry, e->path);

		if (dest_str) {
			if (dest_is_user)
				bpf_probe_read_user_str(e->dest, sizeof(e->dest), dest_str);
			else
				bpf_probe_read_kernel_str(e->dest, sizeof(e->dest), dest_str);
		} else if (dest_dentry) {
			fill_path(dest_dentry, e->dest);
		} else {
			e->dest[0] = '\0';
		}

		bpf_ringbuf_submit(e, 0);
	}

	return is_blocked ? -EPERM : 0;
}

// Wrapper for ordinary call sites: GUARD_REASON_NONE. The raw block-device gate calls
// check_and_emit_ex with GUARD_REASON_RAW_DEVICE.
#define check_and_emit(type, dentry, dest_str, dest_is_user, dest_dentry, is_exec_open, quiet_allow) \
	check_and_emit_ex((type), (dentry), (dest_str), (dest_is_user), (dest_dentry), (is_exec_open), (quiet_allow), false, GUARD_REASON_NONE)

// Lazy directory discovery for file_open: for a file in a directory not yet in guard_inodes, find
// the farthest guarded ancestor and add the file's immediate parent, only if strictly within the
// depth limit (a boundary directory at exactly the limit is not added, so the guard never extends a
// level deeper than configured).
static __always_inline void discover_guarded_parent(struct dentry *dentry)
{
	if (!dentry)
		return;

	struct dentry *parent;
	bpf_probe_read_kernel(&parent, sizeof(parent), &dentry->d_parent);
	if (!parent || parent == dentry)
		return;

	struct inode *parent_inode;
	bpf_probe_read_kernel(&parent_inode, sizeof(parent_inode), &parent->d_inode);
	if (!parent_inode)
		return;

	if (read_inode_guard(parent_inode))
		return;

	// Respect non-recursive mode: only discover parents when recursive=1
	if (!should_add_new_dir())
		return;

	// Read the depth limit. 0 means unlimited.
	__u32 depth_key = 2;
	__u64 *depth = bpf_map_lookup_elem(&guard_config, &depth_key);
	bool has_limit = depth != NULL && *depth > 0;
	__u64 max_steps = has_limit ? *depth : 16;

	// Walk up tracking the farthest guarded ancestor; keep going after the first hit, since the
	// farthest one gives the parent's absolute depth.
	struct dentry *ancestor = parent;
	__u64 farthest_steps = 0;
	bool found_guarded = false;

	for (int i = 0; i < 16 && i < max_steps; i++) {
		struct dentry *next;
		bpf_probe_read_kernel(&next, sizeof(next), &ancestor->d_parent);
		if (!next || next == ancestor)
			break;
		ancestor = next;

		struct inode *anc_inode;
		bpf_probe_read_kernel(&anc_inode, sizeof(anc_inode), &ancestor->d_inode);
		if (!anc_inode)
			break;
		if (anc_inode == parent_inode)
			break;

		if (read_inode_guard(anc_inode)) {
			farthest_steps = i + 1;
			found_guarded = true;
		}
	}

	// Add the parent only if a guarded ancestor was found and the parent is strictly within the
	// depth limit (a parent at exactly `depth` is a boundary dir; adding it would guard one level
	// too deep). Learning is also confined to the watch root so a stale entry (inode reused outside
	// the tree) can't teach the guard another tree's inodes (root_in_chain fails closed on
	// truncation, which just preserves legacy learning; deny-side confinement still applies).
	if (found_guarded && (!has_limit || farthest_steps < *depth) &&
	    root_in_chain(parent, parent_inode, 32))
		add_inode_to_guard(parent_inode);
}

SEC("lsm/file_open")
int guard_file_open(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	if (!file)
		return 0;

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	// Block /proc/<pid>/mem access when the target PID is tainted
	__u32 mem_pid = is_proc_mem_of_tainted(dentry);
	if (mem_pid)
		return emit_process_denial(GUARD_REASON_PROC_MEM, EVENT_OPEN, NULL, mem_pid, PTRACE_MODE_ATTACH);

	// Lazy discovery as in path_mkdir: opening a file whose parent isn't in guard_inodes adds the
	// parent one level at a time. With mkdir-time discovery this maps any runtime-created chain
	// down to the file's parent, so a 16-step walk plus a map hit covers any real depth.
	discover_guarded_parent(dentry);

	if (is_guarded_access(dentry, inode)) {
		__u32 f_flags = 0;
		bpf_probe_read_kernel(&f_flags, sizeof(f_flags), &file->f_flags);
		bool is_exec_open = (f_flags & __FMODE_EXEC) != 0;
		__u32 f_mode = 0;
		bpf_probe_read_kernel(&f_mode, sizeof(f_mode), &file->f_mode);
		bool is_write_open = (f_mode & FMODE_WRITE) != 0;
		int ret = check_and_emit_ex(EVENT_OPEN, dentry, NULL, false, NULL, is_exec_open, false, is_write_open, GUARD_REASON_NONE);
		if (ret == 0 && !is_readonly_mode())
			mark_tainted();
		return ret;
	}

	// Block access to the block device hosting a guarded filesystem (debugfs, dd, fsck read files
	// via raw metadata, bypassing VFS). Coarse by nature: the target is the whole device, hence
	// GUARD_REASON_RAW_DEVICE so userspace doesn't blame a resource.
	if (is_guarded_block_device(inode))
		return check_and_emit_ex(EVENT_OPEN, dentry, NULL, false, NULL, false, false, false, GUARD_REASON_RAW_DEVICE);

	return 0;
}

SEC("lsm/mmap_file")
int guard_mmap_file(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	if (!file)
		return 0;

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (is_guarded_access(dentry, inode)) {
		unsigned long prot = (unsigned long)ctx[2];
		unsigned long flags = (unsigned long)ctx[3];
		bool is_write_map = (prot & PROT_WRITE) && (flags & MAP_SHARED);
		int ret = check_and_emit_ex(EVENT_MMAP, dentry, NULL, false, NULL, false, false, is_write_map, GUARD_REASON_NONE);
		if (ret == 0 && !is_readonly_mode())
			mark_tainted();
		return ret;
	}

	return 0;
}

SEC("lsm/file_permission")
int guard_file_permission(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	int mask = (int)ctx[1];
	if (!file)
		return 0;

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (!dentry)
		return 0;

	// Block reads of /proc/<pid>/mem of a tainted PID (catches fds opened before the taint).
	__u32 mem_pid = is_proc_mem_of_tainted(dentry);
	if (mem_pid)
		return emit_process_denial(GUARD_REASON_PROC_MEM, EVENT_READ, NULL, mem_pid, PTRACE_MODE_ATTACH);

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (!is_guarded_access(dentry, inode))
		return 0;

	__u32 event_type = (mask & MAY_WRITE) ? EVENT_WRITE : EVENT_READ;
	int ret = check_and_emit(event_type, dentry, NULL, false, NULL, false, false);
	if (ret == 0 && !is_readonly_mode())
		mark_tainted();
	return ret;
}

// ftruncate(2) on an fd that predates the guard (or was dup'ed in) passes neither file_open nor
// file_permission; without this hook a blacklisted process holding such an fd could truncate
// guarded files. Path-based truncate(2) is caught by file_open.
SEC("lsm/file_truncate")
int guard_file_truncate(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	if (!file)
		return 0;

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (is_guarded_access(dentry, inode)) {
		int ret = check_and_emit(EVENT_WRITE, dentry, NULL, false, NULL, false, false);
		if (ret == 0 && !is_readonly_mode())
			mark_tainted();
		return ret;
	}

	return 0;
}

SEC("lsm/path_unlink")
int guard_path_unlink(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);

	// Captured up front: guarded_map_hit already computes this, and reusing it lets the single
	// evict-on-allow at the bottom skip a second guard_inodes lookup.
	bool inode_was_guarded = read_inode_guard(inode);

	// i_nlink in this pre-hook is the count BEFORE the unlink: 1 = last name, inode about to be
	// freed; more = another hardlink survives at this exact (dev, ino). guard_inodes is keyed on
	// (dev, ino) alone, so the entry also protects surviving names; evicting on a non-final unlink
	// would strip them until the next re-scan.
	unsigned int nlink = 0;
	if (inode)
		bpf_probe_read_kernel(&nlink, sizeof(nlink), &inode->i_nlink);

	struct dentry *parent = get_dentry_from_path((void *)ctx[0]);
	struct inode *parent_inode = parent ? get_inode_from_path((void *)ctx[0]) : NULL;

	int ret;
	if (guarded_map_hit(dentry, inode, 16)) {
		ret = check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);
	} else if (parent_inode && parent_inode != inode && guarded_map_hit(parent, parent_inode, 12)) {
		ret = check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);
	} else if (parent && guarded_ancestor_within_limit(parent)) {
		// Deep-file coverage: parent not in the map but the file may be inside the guarded region
		// (ancestor walk).
		ret = check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);
	} else {
		ret = 0;
	}

	// Delete is going through (never on a denied one): evict the inode now rather than wait for the
	// periodic sweep (see evict_inode_from_guard). Final link only, per the nlink note.
	if (ret == 0 && inode_was_guarded && nlink <= 1)
		evict_inode_from_guard(inode);
	return ret;
}

SEC("lsm/path_rename")
int guard_path_rename(unsigned long long *ctx)
{
	struct dentry *old_dentry = (struct dentry *)ctx[1];
	struct dentry *new_dentry = (struct dentry *)ctx[3];
	if (!old_dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &old_dentry->d_inode);

	// Source side, one rooted walk: the file and its parent are both on old_dentry's chain, so a
	// single root check confines both map hits (two 32-step walks would exceed the verifier
	// budget).
	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	struct dentry *old_parent = get_dentry_from_path((void *)ctx[0]);
	if ((read_inode_guard(inode) ||
	     (parent_inode && parent_inode != inode && read_inode_guard(parent_inode))) &&
	    root_in_chain(old_dentry, inode, 16)) {
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry, false, false);
	}

	// Deep-file coverage, source side.
	if (old_parent && guarded_ancestor_within_limit(old_parent))
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry, false, false);

	// Destination-target side: a rename onto an existing file silently unlinks the victim
	// (security_path_unlink never fires for a rename target). If the victim IS the watch root, the
	// rename swaps guarded content for an unguarded inode, a bypass for a single-file watch whose
	// parent isn't in guard_inodes. Deny on identity with the watch root; RENAME_EXCHANGE with the
	// root as destination lands here too.
	if (new_dentry) {
		struct inode *victim;
		bpf_probe_read_kernel(&victim, sizeof(victim), &new_dentry->d_inode);
		if (victim && inode_is_watch_root(victim))
			return check_and_emit(EVENT_RENAME, new_dentry, NULL, false, NULL, false, false);
	}

	// Also check the destination parent: blocks renaming files from outside INTO a guarded dir, and
	// into depth-boundary directories added to guard_inodes.
	struct dentry *new_parent = get_dentry_from_path((void *)ctx[2]);
	if (new_parent) {
		struct inode *new_parent_inode = get_inode_from_path((void *)ctx[2]);
		if (new_parent_inode && new_parent_inode != inode && guarded_map_hit(new_parent, new_parent_inode, 12)) {
			int ret = check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry, false, false);
			if (ret != 0)
				return ret;  // blocked — reject the rename
			// Allowed: a renamed directory's inode is added to guard_inodes so its new location is
			// guarded.
			if (should_add_new_dir()) {
				umode_t old_mode;
				bpf_probe_read_kernel(&old_mode, sizeof(old_mode), &inode->i_mode);
				if ((old_mode & S_IFMT) == S_IFDIR)
					add_inode_to_guard(inode);
			}
			return 0;
		}
	}

	// Deep-file coverage, destination side: moving INTO a guarded region whose parent isn't in the
	// map.
	if (new_parent && guarded_ancestor_within_limit(new_parent)) {
		int ret = check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry, false, false);
		if (ret != 0)
			return ret;
		if (should_add_new_dir()) {
			umode_t old_mode;
			bpf_probe_read_kernel(&old_mode, sizeof(old_mode), &inode->i_mode);
			if ((old_mode & S_IFMT) == S_IFDIR)
				add_inode_to_guard(inode);
		}
		return 0;
	}

	return 0;
}

SEC("lsm/path_symlink")
int guard_path_symlink(unsigned long long *ctx)
{
    struct dentry *dentry = (struct dentry *)ctx[1];
    const char *old_name = (const char *)ctx[2];

    // 1. Check if the symlink is being created INSIDE a guarded directory
    struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
    struct dentry *parent = get_dentry_from_path((void *)ctx[0]);
    if (parent && parent_inode && guarded_map_hit(parent, parent_inode, 32)) {
        return check_and_emit(EVENT_SYMLINK, dentry, old_name, false, NULL, false, false);
    }

    // Deep-file coverage: symlink created inside a guarded region whose parent isn't in the map.
    if (parent && guarded_ancestor_within_limit(parent)) {
        return check_and_emit(EVENT_SYMLINK, dentry, old_name, false, NULL, false, false);
    }

    // 2. Symlink TARGET into the guarded path. Skipped in GUARD_MODE_READONLY: the check stops an
    // alias reaching content the creator couldn't read, but a read-only tree is world-readable, and
    // denying breaks tooling (Steam re-points ~/.steam/bin* with `ln -s`). Symlinks created INSIDE
    // the tree stay gated by (1). Secret files in whitelist mode (incl. fscrypt.key,
    // edit-auth.hash) keep this check.
    if (is_readonly_mode())
        return 0;

    __u32 key = 0;
    char *oldname_buf = bpf_map_lookup_elem(&tmp_buf, &key);
    if (!oldname_buf)
        return 0;

    long ret = bpf_probe_read_kernel_str(oldname_buf, MAX_PATH, old_name);
    if (ret <= 0)
        return 0;

    char *stored_path = bpf_map_lookup_elem(&guard_path, &key);
    if (!stored_path)
        return 0;

    // Prefix match logic to protect contents of a watched directory
    bool match = true;
    for (int i = 0; i < MAX_PATH; i++) {
        if (stored_path[i] == '\0') {
            // End of the guarded path: a match if the symlink target ends here, continues as a
            // sub-directory, or the stored path had a trailing slash.
            if (oldname_buf[i] == '\0' || oldname_buf[i] == '/' || (i > 0 && stored_path[i-1] == '/')) {
                break;
            }
            // Otherwise, it's just a similarly named folder (e.g., /bypasses_fake)
            match = false;
            break;
        }
        if (oldname_buf[i] != stored_path[i]) {
            match = false;
            break;
        }
    }

    if (match) {
        return check_and_emit(EVENT_SYMLINK, dentry, old_name, false, NULL, false, false);
    }

    return 0;
}

SEC("lsm/path_link")
int guard_path_link(unsigned long long *ctx)
{
	struct dentry *old_dentry = (struct dentry *)ctx[0];
	struct dentry *new_dentry = (struct dentry *)ctx[2];
	if (!old_dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &old_dentry->d_inode);

	if (guarded_map_hit(old_dentry, inode, 16)) {
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry, false, false);
	}

	// Source-side confinement: link(2) needs no read access to the source, so a file created after
	// guard start (only its parent in guard_inodes) could be hardlinked out and read via the escape
	// name. Hook args carry no source parent, so walk from old_dentry, as path_rename's source side
	// does.
	struct dentry *old_parent;
	bpf_probe_read_kernel(&old_parent, sizeof(old_parent), &old_dentry->d_parent);
	if (old_parent && old_parent != old_dentry) {
		struct inode *old_parent_inode;
		bpf_probe_read_kernel(&old_parent_inode, sizeof(old_parent_inode), &old_parent->d_inode);
		if (old_parent_inode && old_parent_inode != inode &&
		    guarded_map_hit(old_parent, old_parent_inode, 16))
			return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry, false, false);
	}
	// Deep-file coverage, source side.
	if (old_parent && guarded_ancestor_within_limit(old_parent))
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry, false, false);

	struct dentry *new_parent = get_dentry_from_path((void *)ctx[1]);
	if (new_parent) {
		struct inode *parent_inode = get_inode_from_path((void *)ctx[1]);
		if (parent_inode && parent_inode != inode && guarded_map_hit(new_parent, parent_inode, 12)) {
			return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry, false, false);
		}
	}

	// Deep-file coverage, destination side: linking INTO a guarded region whose parent isn't in the
	// map.
	if (new_parent && guarded_ancestor_within_limit(new_parent))
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry, false, false);

	return 0;
}

SEC("lsm/path_mkdir")
int guard_path_mkdir(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];

	// Lazy discovery as in file_open: each runtime mkdir adds the parent one level at a time, so
	// any runtime-created chain is fully mapped at any depth with a short ancestor walk.
	discover_guarded_parent(dentry);

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);

	struct dentry *parent = get_dentry_from_path((void *)ctx[0]);

	if (parent && guarded_map_hit(parent, parent_inode, 32)) {
		return check_and_emit(EVENT_MKDIR, dentry, NULL, false, NULL, false, false);
	}

	// Deep-file coverage: mkdir inside a guarded region whose parent isn't in the map.
	if (parent && guarded_ancestor_within_limit(parent))
		return check_and_emit(EVENT_MKDIR, dentry, NULL, false, NULL, false, false);

	return 0;
}

// Path/attribute bypass coverage: these ops create no struct file, so file_open/file_permission
// don't fire:
//   - truncate(2) -> path_truncate
//   - chmod/chown/utimes -> notify_change -> inode_setattr
//   - setxattr/removexattr -> inode_setxattr/inode_removexattr
//   - mknod(2) -> path_mknod
//   - rmdir(2) -> path_rmdir

SEC("lsm/path_truncate")
int guard_path_truncate(unsigned long long *ctx)
{
	struct dentry *dentry = get_dentry_from_path((void *)ctx[0]);
	if (!dentry)
		return 0;

	struct inode *inode = get_inode_from_path((void *)ctx[0]);
	if (is_guarded_access(dentry, inode))
		return check_and_emit(EVENT_ATTR, dentry, NULL, false, NULL, false, false);

	return 0;
}

SEC("lsm/inode_setattr")
int guard_inode_setattr(unsigned long long *ctx)
{
	// ctx[0]=mnt_idmap, [1]=dentry, [2]=iattr. Size changes are skipped:
	// path_truncate/file_truncate own them and already deny (ATTR_FILE follows ftruncate(2));
	// re-emitting would duplicate events.
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);

	if (!is_guarded_access(dentry, inode))
		return 0;

	struct iattr *attr = (struct iattr *)ctx[2];
	if (attr) {
		__u32 ia_valid;
		bpf_probe_read_kernel(&ia_valid, sizeof(ia_valid), &attr->ia_valid);
		if (ia_valid & (ATTR_SIZE | ATTR_FILE))
			return 0;
	}

	return check_and_emit(EVENT_ATTR, dentry, NULL, false, NULL, false, false);
}

SEC("lsm/inode_setxattr")
int guard_inode_setxattr(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);

	if (is_guarded_access(dentry, inode))
		return check_and_emit(EVENT_ATTR, dentry, NULL, false, NULL, false, false);

	return 0;
}

SEC("lsm/inode_removexattr")
int guard_inode_removexattr(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);

	if (is_guarded_access(dentry, inode))
		return check_and_emit(EVENT_ATTR, dentry, NULL, false, NULL, false, false);

	return 0;
}

SEC("lsm/path_mknod")
int guard_path_mknod(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	struct dentry *parent = get_dentry_from_path((void *)ctx[0]);
	if (parent && parent_inode && guarded_map_hit(parent, parent_inode, 32))
		return check_and_emit(EVENT_MKNOD, dentry, NULL, false, NULL, false, false);

	if (parent && guarded_ancestor_within_limit(parent))
		return check_and_emit(EVENT_MKNOD, dentry, NULL, false, NULL, false, false);

	return 0;
}

SEC("lsm/path_rmdir")
int guard_path_rmdir(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;

	// rmdir only succeeds on an EMPTY directory, so evicting this dentry's own entry on allow is
	// complete: descendants were each evicted by their own unlink/rmdir.
	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);
	bool inode_was_guarded = read_inode_guard(inode);

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	struct dentry *parent = get_dentry_from_path((void *)ctx[0]);

	int ret;
	if (parent && parent_inode && guarded_map_hit(parent, parent_inode, 32)) {
		ret = check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);
	} else if (parent && guarded_ancestor_within_limit(parent)) {
		ret = check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);
	} else {
		ret = 0;
	}

	if (ret == 0 && inode_was_guarded)
		evict_inode_from_guard(inode);
	return ret;
}

// stat/statx/lstat on a guarded path leaks size, mtime, mode, owner, nlink and existence: denied
// like any access. The hook gets a `struct path *` (ctx[0]), so is_guarded_access confines to the
// watch root like inode_readlink/path_truncate.
//
// Tools that stat before their real syscall (cp, rm, mv, ls) are now denied here first (event STAT,
// not OPEN/DELETE/RENAME). The daemon's root-gated self mask includes EVENT_STAT so its inode scan
// and fscrypt lifecycle keep working.
SEC("lsm/inode_getattr")
int guard_inode_getattr(unsigned long long *ctx)
{
	struct dentry *dentry = get_dentry_from_path((void *)ctx[0]);
	if (!dentry)
		return 0;

	struct inode *inode = get_inode_from_path((void *)ctx[0]);

	if (!is_guarded_access(dentry, inode))
		return 0;

	// Allow stat of the watch-root DIRECTORY itself: leaks almost nothing (the path is in the
	// config) and denying breaks `mkdir -p`, path probes and `ls` of the parent. A single-file
	// watch root IS the secret, so its stat stays denied.
	if (inode && inode_is_watch_root(inode)) {
		umode_t mode = 0;
		bpf_probe_read_kernel(&mode, sizeof(mode), &inode->i_mode);
		if ((mode & S_IFMT) == S_IFDIR)
			return 0;
	}

	return check_and_emit(EVENT_STAT, dentry, NULL, false, NULL, false, false);
}

SEC("lsm/inode_readlink")
int guard_inode_readlink(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[0];
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);

	if (is_guarded_access(dentry, inode))
		return check_and_emit(EVENT_STAT, dentry, NULL, false, NULL, false, false);

	return 0;
}

// getxattr/listxattr on a guarded file leak its xattrs; neither creates a struct file or calls
// inode_getattr, so file_open/file_permission/inode_getattr all miss them. Denied like STAT
// (EVENT_STAT) with is_guarded_access's watch-root and whitelist confinement; the daemon's self
// mask includes EVENT_STAT.
SEC("lsm/inode_getxattr")
int guard_inode_getxattr(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[0];
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);

	if (is_guarded_access(dentry, inode))
		return check_and_emit(EVENT_STAT, dentry, NULL, false, NULL, false, false);

	return 0;
}

SEC("lsm/inode_listxattr")
int guard_inode_listxattr(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[0];
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);

	if (is_guarded_access(dentry, inode))
		return check_and_emit(EVENT_STAT, dentry, NULL, false, NULL, false, false);

	return 0;
}

// inode_permission handles pure may-checks (access(2)/faccessat(2), MAY_ACCESS). Its context is
// only (inode, mask): no dentry, so the watch-root chain can't be walked and a raw map hit can't be
// confined (it false-denied inode numbers reused in other trees, and the audit denial was a log lie
// since the syscall was never blocked). access/faccessat leak no content, so it passes everything
// through; real reads/writes/execs/deletes are enforced in file_open and the path_*/inode_* hooks,
// which carry a dentry. An X_OK probe (MAY_EXEC without MAY_OPEN) is deliberately not denied. The
// MAY_ACCESS selection below keeps this hook from intercepting opens, execs, truncates or setxattr.
SEC("lsm/inode_permission")
int guard_inode_permission(unsigned long long *ctx)
{
	struct inode *inode = (struct inode *)ctx[0];
	if (!inode)
		return 0;

	int mask = (int)ctx[1];
	// Only genuine access(2)/faccessat probes carry MAY_ACCESS.
	if (!(mask & MAY_ACCESS))
		return 0;
	if (!(mask & (MAY_READ | MAY_WRITE)))
		return 0;
	if (mask & (MAY_EXEC | MAY_OPEN))
		return 0;

	umode_t mode;
	bpf_probe_read_kernel(&mode, sizeof(mode), &inode->i_mode);
	if ((mode & S_IFMT) == S_IFDIR)
		return 0;

	return 0;
}

SEC("lsm/sb_mount")
int guard_sb_mount(unsigned long long *ctx)
{
	// ctx[1] is the mount point path (struct path *)
	struct path *mount_path = (struct path *)ctx[1];
	if (!mount_path)
		return 0;

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &mount_path->dentry);
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);
	if (!inode)
		return 0;

	if (guarded_map_hit(dentry, inode, 32)) {
		// Not write_intent: in GUARD_MODE_READONLY a mount over the guarded dir passes. Blocking it
		// broke systemd's namespace setup for the daemon unit's helpers (ExecReload=/bin/kill with
		// ReadWritePaths=/etc/app-listener bind-mounts the dir, giving 226/NAMESPACE). mount(2)
		// needs CAP_SYS_ADMIN, outside what an LSM file guard can contain; whitelist/blacklist
		// still deny (mount-shadow regression test).
		return check_and_emit(EVENT_OPEN, dentry, NULL, false, NULL, false, false);
	}

	return 0;
}

// Block ptrace/process_vm_readv/writev against any tainted process (one holding guarded content in
// memory) unless the caller is whitelisted.
SEC("lsm/ptrace_access_check")
int guard_ptrace_access_check(unsigned long long *ctx)
{
	// LSM hooks pass a pointer to the arg struct in r1; the declared-args form maps child to r1 and
	// mode to r2, but the kernel only fills r1, so child->tgid would read past the struct. Read
	// args through ctx like the other hooks.
	struct task_struct *child = (struct task_struct *)ctx[0];
	__u32 mode = (__u32)ctx[1];

	if (!child)
		return 0;

	__u32 pid;
	bpf_probe_read_kernel(&pid, sizeof(pid), &child->tgid);

	__u8 *val = bpf_map_lookup_elem(&guard_tainted_pids, &pid);
	if (!val)
		return 0;  // child is not tainted — allow

	// The owner's process (the daemon: the guard's only GUARD_ALLOW_ROOT entry) is tainted because
	// it reads what it guards, but its secrets (fscrypt key, edit-auth hash) live in MEMORY,
	// reached only by PTRACE_MODE_ATTACH. Metadata-only access (READ: /proc/<pid>/environ, fd,
	// maps), which journald does to label every log line, is allowed. Identity is the target's exe
	// inode, never its pid.
	if (!(mode & PTRACE_MODE_ATTACH)) {
		struct inode_key child_ik = {};
		if (get_task_exe_inode(child, &child_ik)) {
			__u8 *child_action = bpf_map_lookup_elem(&guard_exe_actions, &child_ik);
			if (child_action && *child_action == GUARD_ALLOW_ROOT)
				return 0;
		}
	}

	// Child is tainted.  Check if the caller is whitelisted.
	struct inode_key exe_ik = {};
	if (!get_current_exe_inode(&exe_ik))
		return emit_process_denial(GUARD_REASON_PTRACE, EVENT_READ, child, 0, mode);

	__u8 *action = bpf_map_lookup_elem(&guard_exe_actions, &exe_ik);
	if (!is_allow_action(action))
		return emit_process_denial(GUARD_REASON_PTRACE, EVENT_READ, child, 0, mode);

	return 0;
}

SEC("lsm/bprm_check_security")
int guard_bprm_check_security(unsigned long long *ctx)
{
	// Executing a whitelisted binary while ptraced hands the tracer the victim's whole memory: the
	// tracer may have attached BEFORE the victim's exe identity was whitelisted (so the attach
	// check couldn't deny) and PEEKDATA fires no further hook. Allow the exec only if the tracer is
	// whitelisted.
	//
	// Read-only guards (lib_dir trees, /etc/app-listener) are skipped like every other taint site:
	// their contents are world-readable code, and tainting their writers would lock most of e.g.
	// Steam's process tree out of inspection.
	if (is_readonly_mode())
		return 0;

	struct linux_binprm *bprm = (struct linux_binprm *)ctx[0];
	if (!bprm)
		return 0;

	struct file *file = NULL;
	bpf_probe_read_kernel(&file, sizeof(file), &bprm->file);
	if (!file)
		return 0;

	struct inode *inode = NULL;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	if (!inode)
		return 0;

	struct inode_key target_ik = {};
	bpf_probe_read_kernel(&target_ik.ino, sizeof(target_ik.ino), &inode->i_ino);

	struct super_block *sb = NULL;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;

	dev_t dev = 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	target_ik.dev = dev;

	__u8 *taction = bpf_map_lookup_elem(&guard_exe_actions, &target_ik);
	if (!is_allow_action(taction))
		return 0;  // target is not whitelisted: not our concern

	// Executing a whitelisted image taints the process IMMEDIATELY, so a tracer attaching after
	// exec is denied even before the victim touches the tree. (Exec-open attribution in
	// check_and_emit handles the open itself.)
	mark_tainted();

	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (!task)
		return 0;

	// Only TRACED processes race: an untraced exec is the normal exec-open attribution path (a
	// non-whitelisted launcher may run a whitelisted in-tree binary; the exec fd is never exposed
	// to it).
	__u32 ptrace_flags = 0;
	bpf_probe_read_kernel(&ptrace_flags, sizeof(ptrace_flags), &task->ptrace);
	if (!ptrace_flags)
		return 0;

	// For a traced task ->parent is the tracer (real_parent is the biological parent). A
	// non-whitelisted tracer must not gain a whitelisted image's memory via exec (it attached
	// before the victim's identity existed and can PEEKDATA the content once read).
	struct task_struct *tracer = NULL;
	bpf_probe_read_kernel(&tracer, sizeof(tracer), &task->parent);
	if (!tracer)
		return emit_process_denial(GUARD_REASON_TRACED_EXEC, EVENT_OPEN, NULL, 0, 0);

	struct inode_key tracer_ik = {};
	if (!get_task_exe_inode(tracer, &tracer_ik))
		return emit_process_denial(GUARD_REASON_TRACED_EXEC, EVENT_OPEN, tracer, 0, 0);  // no resolvable identity: fail closed

	__u8 *tracer_action = bpf_map_lookup_elem(&guard_exe_actions, &tracer_ik);
	if (!is_allow_action(tracer_action))
		return emit_process_denial(GUARD_REASON_TRACED_EXEC, EVENT_OPEN, tracer, 0, 0);

	return 0;
}

// guard_inode_free forgets an inode the moment the kernel destroys a DELETED file. guard_inodes is
// keyed on (dev, ino) and filesystems recycle inode numbers, so a stale entry eventually matches an
// unrelated file. unlink/rmdir evict on the final link, but a rename OVER a tracked file frees the
// victim with no unlink (how apps save, e.g. Steam's registry.vdf). Waiting for the hourly
// ReconcileInodes caused false DENYs: past root_in_chain's bounded walk a collision fails closed,
// and a recycled watch-root number exact-matched the root.
//
// i_nlink == 0 is the load-bearing check. inode_free_security fires for EVERY inode leaving memory,
// incl. live files evicted from the cache (memory pressure, drop_caches, fscrypt lock). Forgetting
// a live file would unguard it, and a live watch root the entire tree. Only a nameless file is
// really gone.
//
// Separate program so it adds nothing to the verifier budget of rename/unlink/link
// (guard_path_rename sits at the 1M-insn limit, issue #45). Best-effort: without it stale entries
// wait for ReconcileInodes.
SEC("lsm/inode_free_security")
int guard_inode_free(unsigned long long *ctx)
{
	struct inode *inode = (struct inode *)ctx[0];
	if (!inode)
		return 0;

	unsigned int nlink = 1;
	bpf_probe_read_kernel(&nlink, sizeof(nlink), &inode->i_nlink);
	if (nlink != 0)
		return 0; // live file leaving the cache, not a deletion

	struct super_block *sb = NULL;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	__u64 dev = 0;
	bpf_probe_read_kernel(&dev, sizeof(dev_t), &sb->s_dev);
	if (!bpf_map_lookup_elem(&guard_fs_sbdevs, &dev))
		return 0;

	__u64 ino = 0;
	bpf_probe_read_kernel(&ino, sizeof(ino), &inode->i_ino);
	struct inode_key ikey = {};
	ikey.dev = dev;
	ikey.ino = ino;
	bpf_map_delete_elem(&guard_inodes, &ikey);

	// A deleted watch root protects nothing and its number is about to be reused: forget it (0
	// matches no inode) so another file can't exact-match the root. Userspace re-anchors on its
	// next sweep.
	__u32 dev_key = 3, root_key = 4;
	__u64 *root_dev = bpf_map_lookup_elem(&guard_config, &dev_key);
	__u64 *root_ino = bpf_map_lookup_elem(&guard_config, &root_key);
	if (root_dev && root_ino && *root_dev == dev && *root_ino == ino) {
		__u64 zero = 0;
		bpf_map_update_elem(&guard_config, &root_key, &zero, BPF_ANY);
	}
	return 0;
}

// guard_bprm_committed clears taint when the process exec'd an image this guard doesn't whitelist.
// Taint = holds this resource's secrets in memory; it's inherited across fork but must not survive
// exec (the address space is discarded), else every descendant of e.g. the Steam client (Proton,
// helpers, the game) stays "holding credentials" and desktop tools (GameMode, PipeWire, Wine's
// wineserver) are refused.
//
// Runs at bprm_committed_creds, after the point of no return: clearing earlier
// (bprm_check_security) would fail open, since a failed exec leaves the OLD secret-holding image
// running untainted.
//
// A whitelisted image stays tainted (marked in bprm_check_security; it may read the resource from
// its first instruction). What still crosses an exec: open fds (stay guarded; each read is
// whitelist-checked) and argv/envp (data the whitelisted parent chose to pass, outside taint's
// scope).
SEC("lsm/bprm_committed_creds")
int guard_bprm_committed(unsigned long long *ctx)
{
	if (is_readonly_mode())
		return 0; // read-only guards never taint

	__u32 tgid = bpf_get_current_pid_tgid() >> 32;
	if (!bpf_map_lookup_elem(&guard_tainted_pids, &tgid))
		return 0;

	struct linux_binprm *bprm = (struct linux_binprm *)ctx[0];
	if (!bprm)
		return 0;
	struct file *file = NULL;
	bpf_probe_read_kernel(&file, sizeof(file), &bprm->file);
	if (!file)
		return 0; // unresolvable image: keep the taint (fail closed)
	struct inode *inode = NULL;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	if (!inode)
		return 0;
	struct inode_key ik = {};
	bpf_probe_read_kernel(&ik.ino, sizeof(ik.ino), &inode->i_ino);
	struct super_block *sb = NULL;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	dev_t dev = 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	ik.dev = dev;

	__u8 *action = bpf_map_lookup_elem(&guard_exe_actions, &ik);
	if (is_allow_action(action))
		return 0; // a whitelisted image may read the resource: stays tainted

	bpf_map_delete_elem(&guard_tainted_pids, &tgid);
	return 0;
}

// guard_sched_process_fork taints the child of a tainted parent (it inherits the address space,
// incl. guarded content; else process_vm_readv on the child reads the secret). Runs at
// sched_process_fork, NOT security_task_alloc: at task_alloc the child's pid/tgid are still the
// PARENT's (assigned later in copy_process), so only the parent would be re-stamped. Also fires for
// thread creation (child->tgid == parent->tgid): a harmless no-op.
SEC("tp_btf/sched_process_fork")
int guard_sched_process_fork(unsigned long long *ctx)
{
	struct task_struct *parent = (struct task_struct *)ctx[0];
	struct task_struct *child = (struct task_struct *)ctx[1];
	if (!parent || !child)
		return 0;

	__u32 parent_tgid = 0;
	bpf_probe_read_kernel(&parent_tgid, sizeof(parent_tgid), &parent->tgid);
	if (!bpf_map_lookup_elem(&guard_tainted_pids, &parent_tgid))
		return 0; // parent not tainted

	__u32 child_tgid = 0;
	bpf_probe_read_kernel(&child_tgid, sizeof(child_tgid), &child->tgid);

	__u8 v = 1;
	if (bpf_map_update_elem(&guard_tainted_pids, &child_tgid, &v, BPF_ANY))
		count_degrade(1);
	return 0;
}

// Clear the tainted-PID entry on exit so a reused PID isn't stale-tainted.
SEC("lsm/task_free")
int guard_task_free(unsigned long long *ctx)
{
	struct task_struct *task = (struct task_struct *)ctx[0];

	if (!task)
		return 0;

	// guard_tainted_pids is keyed on tgid: clear only when the GROUP LEADER exits.
	// security_task_free fires per task incl. threads; deleting on any thread's exit (pid != tgid)
	// would strip a running multithreaded process's taint and reopen process_vm_readv/ptrace.
	__u32 pid, tgid;
	bpf_probe_read_kernel(&pid, sizeof(pid), &task->pid);
	bpf_probe_read_kernel(&tgid, sizeof(tgid), &task->tgid);
	if (pid != tgid)
		return 0;

	bpf_map_delete_elem(&guard_tainted_pids, &tgid);
	return 0;
}
