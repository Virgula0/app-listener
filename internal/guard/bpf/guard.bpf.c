#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <errno.h>

#define MAX_PATH 256
// VFS access-mask bits as passed to bpf_lsm_inode_permission /
// bpf_lsm_file_permission.  NOTE: this kernel (7.0-hardened) reshuffled
// the mainline layout — MAY_OPEN is 0x20 here (not 0x40) and MAY_ACCESS
// is 0x10.  Keep in sync with include/linux/fs.h of the target kernel.
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
#define __FMODE_EXEC 0x20 // set in file->f_flags by kernel exec opens (do_open_execat)

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
	__uint(max_entries, 3);  // [0]=mode, [1]=recursive, [2]=depth
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

// Per-binary allowed-event bitmask (bit i = enum event_type i). Only
// consulted in whitelist mode; a missing entry means all events are
// allowed, so the plain guard mode (which never populates this map)
// behaves exactly as before.
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

// Tracks processes that have guarded file content in their address space.
// Once tainted, process_vm_readv, ptrace, and /proc/<pid>/mem access
// against this PID are blocked unless the caller is whitelisted.
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
	__uint(max_entries, 1 << 24);
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

// Check if a newly created directory should be auto-added to guard_inodes
// based on the guard's recursive setting.
static __always_inline int should_add_new_dir(void)
{
	__u32 key = 1;
	__u64 *recursive = bpf_map_lookup_elem(&guard_config, &key);
	if (!recursive || *recursive == 0)
		return 0;
	return 1;
}

// Walk up from `parent` toward the root and report whether the subtree
// under it lies inside the guarded region.  The guarded region is the
// subtree rooted at the farthest (closest-to-root) inode present in
// guard_inodes, bounded by the configured depth limit; a parent that sits
// exactly at the limit boundary is not guarded, which mirrors the eager
// inode scan and discover_guarded_parent exactly.  This makes the guarded
// root in the map protect the entire subtree at any depth: it closes the
// window between LSM hook attach and the eager inode scan, keeps deep
// files (deeper than the discovery bound) protected, and preserves open/
// read/write enforcement even when the inode map fills up mid-scan.
static __always_inline int guarded_ancestor_within_limit(struct dentry *parent)
{
	if (!parent)
		return 0;

	// Non-recursive guards only ever add the guarded root and its direct
	// children to the map; deeper levels are deliberately unprotected.
	if (!should_add_new_dir())
		return 0;

	struct inode *parent_inode;
	bpf_probe_read_kernel(&parent_inode, sizeof(parent_inode), &parent->d_inode);
	if (!parent_inode)
		return 0;

	// Performance gate: only files on a filesystem that hosts a guarded
	// root can have a guarded ancestor.  guard_fs_sbdevs holds the
	// filesystem device of the guarded path (populated unconditionally by
	// the Go side, including anonymous tmpfs/overlayfs devices); every
	// other filesystem skips the walk entirely, keeping the per-access
	// cost at a single hash lookup.
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

	// Track the farthest guarded ancestor (closest to root) and its
	// distance from `parent`, like discover_guarded_parent.  The 128-step
	// bound is far beyond any realistic Linux path depth and is the
	// largest the verifier accepts for this loop body (a 160-step bound
	// already exceeds the 1M-instruction complexity budget).  With a
	// depth limit the result is decided early: a guarded ancestor at or
	// beyond the limit means the subtree is NOT guarded (a farther one
	// would only be more distant), and finding none within the limit
	// decides the same.
	struct dentry *ancestor = parent;
	__u64 farthest_steps = 0;
	bool found_guarded = false;

	for (int i = 0; i < 128; i++) {
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
			if (has_limit && farthest_steps >= *depth)
				return 0;  // guarded ancestor at/beyond the limit boundary
		}

		if (has_limit && !found_guarded && i + 1 >= *depth)
			return 0;  // nothing guarded within the limit
	}

	if (!found_guarded)
		return 0;
	// With a limit, no early exit fired, so every guarded ancestor is
	// strictly inside the limit.
	return 1;
}

// Check whether the operation targets a guarded path by examining
// both the file's own inode and its parent directory's inode.
#define S_IFMT  00170000
#define S_IFBLK 0060000
#define S_IFDIR 0040000

// Check if the inode is a block device whose dev_t matches a
// guarded filesystem.  This catches debugfs, dd, and any raw
// block device reader that bypasses VFS access control.
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

// Mark the current process as tainted — it now has guarded file content
// in its address space.  Subsequent ptrace/process_vm_readv attempts
// from non-whitelisted processes will be blocked.
static __always_inline void mark_tainted(void)
{
	__u32 pid = bpf_get_current_pid_tgid() >> 32;
	__u8 val = 1;
	bpf_map_update_elem(&guard_tainted_pids, &pid, &val, BPF_ANY);
}

// Add an inode to guard_inodes.  Used by inode_mkdir and path_rename
// to auto-discover new directories at runtime.
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
	bpf_map_update_elem(&guard_inodes, &ikey, &v, BPF_ANY);
}

// Check if the accessed file is /proc/<pid>/mem for a tainted PID.
// The guard can't protect /proc filesystem inodes (they're not in the
// guarded inode map), so we must explicitly detect this vector.
static __always_inline int is_proc_mem_of_tainted(struct dentry *dentry)
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
	return val != NULL;
}

// Check /proc/<pid>/fd/<n> — opening it creates a file struct pointing
// to the actual file, which goes through security_file_open with the
// real file's dentry/inode, so is_guarded_access catches it.
// We only need to explicitly guard /proc/<pid>/mem here.

static __always_inline int is_guarded_access(struct dentry *dentry, struct inode *file_inode)
{
	if (read_inode_guard(file_inode))
		return 1;

	if (!dentry)
		return 0;

	struct dentry *parent;
	bpf_probe_read_kernel(&parent, sizeof(parent), &dentry->d_parent);
	if (!parent || parent == dentry)
		return 0;

	struct inode *parent_inode;
	bpf_probe_read_kernel(&parent_inode, sizeof(parent_inode), &parent->d_inode);
	if (parent_inode && parent_inode != file_inode && read_inode_guard(parent_inode))
		return 1;

	// The file's own inode and its direct parent are not in the map (a
	// runtime-created directory, a file deeper than the eager-scan depth,
	// or the inode map filled up mid-scan).  Walk up the dentry chain: if
	// the subtree lies inside the guarded region, the access is guarded.
	return guarded_ancestor_within_limit(parent);
}

static __always_inline long read_path(struct dentry *dentry, char *buf, int buf_size)
{
    char tmp[64];
    int pos = MAX_PATH;

    pos--;
    buf[pos & (MAX_PATH - 1)] = '\0';

    struct dentry *d = dentry;

    // Notice we removed #pragma unroll. Modern verifiers handle bounded loops
    // beautifully without exploding the instruction limit.
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

static __always_inline int check_and_emit(__u32 type, struct dentry *dentry, const char *dest_str, bool dest_is_user, struct dentry *dest_dentry, bool is_exec_open, bool quiet_allow)
{
	__u32 key = 0;
	__u64 *mode = bpf_map_lookup_elem(&guard_config, &key);
	if (!mode)
		return -EPERM;

	// Save mode value early to avoid verifier pointer-tracking concerns
	__u64 mode_val = *mode;

	// Detect kernel threads (no userspace mm).  These execute on behalf of
	// user processes (e.g., io_uring workers, SQPOLL threads) and cannot be
	// identified via exe inode.  Block them unconditionally.
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

		if (mode_val == 0) {
			is_blocked = action != NULL && *action == GUARD_BLOCK;
		} else {
			is_blocked = action == NULL || *action != GUARD_ALLOW;

			// Exec-open attribution (whitelist mode only): executing a
			// binary is an OPEN of that binary performed by the
			// *launcher*, not by the binary itself — a shell wrapper
			// runs through its interpreter (/bin/sh), whose exe inode
			// the whitelist deliberately does not list. Attribute
			// exec-opens to the binary being executed: allowlist the
			// *target*. The exec fd is never exposed to the caller
			// (execve destroys the caller's address space), so this
			// grants the launcher no way to read guarded content;
			// afterwards the process identity IS the target, so every
			// later access is checked against the whitelist as usual.
			// Blacklist mode keeps caller attribution (blocking listed
			// binaries from exec'ing guarded files).
			//
			// Two discriminator sources cover the full exec chain:
			//  - the exec target's own open carries __FMODE_EXEC in
			//    file->f_flags (task->in_execve is set only afterwards,
			//    in bprm_execve);
			//  - the kernel's load-time accesses of the exec file
			//    (prepare_binprm's kernel_read, binfmt mmap, in-tree
			//    interpreter opens) happen with task->in_execve set.
			// Both are attributed by the *accessed file's* inode: only
			// a file that is itself whitelisted passes.
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
						if (*target_action == GUARD_ALLOW) {
							exe_ik = target_ik;
							is_blocked = 0;
						}
					}
				}
			}

			// Per-binary event restriction: if the binary has an allowed
			// event mask, only the listed event types are permitted. A
			// missing mask entry means all events are allowed.
			if (!is_blocked) {
				__u32 *mask = bpf_map_lookup_elem(&guard_exe_events, &exe_ik);
				if (mask && type < 32 && !(*mask & (1U << type)))
					is_blocked = 1;
			}
		}
	}

	// Allowed accesses with quiet_allow set are not reported: this is
	// used by inode_permission, which fires for every access to a
	// guarded file and would otherwise flood the log with allowed-STAT
	// events. Denied accesses are always reported.
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

// Lazy directory discovery for file_open.
// When a file is accessed inside a newly-created directory whose inode is
// not yet in guard_inodes, walk up the dentry tree to find the farthest
// guarded ancestor (closest to root) and add the immediate parent of the
// file to guard_inodes — but only if the parent is strictly within the
// configured depth limit.  Boundary directories (those at the exact depth
// limit) are not added, which prevents the guard from extending a level
// deeper than the user specified.
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

	// Walk up from parent, tracking the farthest (closest to root) guarded
	// ancestor found.  We continue walking even after finding a guarded
	// ancestor because we need the farthest one to correctly compute the
	// absolute depth of the parent directory.
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

	// Only add the parent if we found a guarded ancestor and the parent is
	// strictly within the depth limit.  A parent at exactly `depth` levels
	// from the farthest guarded ancestor is a boundary directory that should
	// NOT be added — otherwise files one level deeper become guarded.
	if (found_guarded && (!has_limit || farthest_steps < *depth))
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
	if (is_proc_mem_of_tainted(dentry))
		return -EPERM;

	// Lazy directory discovery: if a file's parent is newly created and not
	// yet in guard_inodes but a higher ancestor is guarded, add the parent.
	// This prevents bypasses where a whitelisted binary creates a directory
	// at runtime and an attacker reads files inside.
	discover_guarded_parent(dentry);

	if (is_guarded_access(dentry, inode)) {
		__u32 f_flags = 0;
		bpf_probe_read_kernel(&f_flags, sizeof(f_flags), &file->f_flags);
		bool is_exec_open = (f_flags & __FMODE_EXEC) != 0;
		int ret = check_and_emit(EVENT_OPEN, dentry, NULL, false, NULL, is_exec_open, false);
		if (ret == 0)
			mark_tainted();
		return ret;
	}

	// Block access to the block device that hosts a guarded filesystem.
	// Without this, tools like debugfs, dd, and fsck can read guarded
	// files by opening the raw block device and interpreting filesystem
	// metadata directly (bypassing VFS entirely).
	if (is_guarded_block_device(inode))
		return check_and_emit(EVENT_OPEN, dentry, NULL, false, NULL, false, false);

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
		int ret = check_and_emit(EVENT_MMAP, dentry, NULL, false, NULL, false, false);
		if (ret == 0)
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

	// Block read access to /proc/<pid>/mem when the target PID is tainted.
	// This catches reads on an already-open fd (opened before taint).
	if (is_proc_mem_of_tainted(dentry))
		return -EPERM;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (!is_guarded_access(dentry, inode))
		return 0;

	__u32 event_type = (mask & MAY_WRITE) ? EVENT_WRITE : EVENT_READ;
	int ret = check_and_emit(event_type, dentry, NULL, false, NULL, false, false);
	if (ret == 0)
		mark_tainted();
	return ret;
}

// ftruncate(2) on an fd that predates the guard (opened before the hooks
// attached, or dup'ed into the process) never passes through file_open,
// and the kernel does not invoke file_permission for it either — without
// this hook a blacklisted process holding such an fd could truncate
// guarded files.  path-based truncate(2) is already caught by file_open.
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
		if (ret == 0)
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

	if (read_inode_guard(inode)) {
		return check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);
	}

	// Deep-file coverage: the parent is not in the map but the file may
	// still be inside the guarded region (ancestor walk).
	struct dentry *parent = get_dentry_from_path((void *)ctx[0]);
	if (parent && guarded_ancestor_within_limit(parent))
		return check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);

	return 0;
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

	if (read_inode_guard(inode)) {
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry, false, false);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry, false, false);
	}

	// Deep-file coverage on the source side: the file may be inside the
	// guarded region even though its parent is not in the map.
	struct dentry *old_parent = get_dentry_from_path((void *)ctx[0]);
	if (old_parent && guarded_ancestor_within_limit(old_parent))
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry, false, false);

	// Also check the destination parent directory.  This prevents renaming
	// files from outside the guarded area INTO a guarded directory, and
	// blocks renames into depth-boundary directories that were added to
	// the guard_inodes map.
	struct inode *new_parent_inode = get_inode_from_path((void *)ctx[2]);
	if (new_parent_inode && new_parent_inode != inode && read_inode_guard(new_parent_inode)) {
		int ret = check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry, false, false);
		if (ret != 0)
			return ret;  // blocked — reject the rename
		// Allowed — if the source is a directory, add its inode
		// to guard_inodes so its new location is guarded.
		if (should_add_new_dir()) {
			umode_t old_mode;
			bpf_probe_read_kernel(&old_mode, sizeof(old_mode), &inode->i_mode);
			if ((old_mode & S_IFMT) == S_IFDIR)
				add_inode_to_guard(inode);
		}
		return 0;
	}

	// Deep-file coverage on the destination side: moving INTO a guarded
	// region whose parent is not in the map (runtime-created deep dir).
	struct dentry *new_parent = get_dentry_from_path((void *)ctx[2]);
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
    if (read_inode_guard(parent_inode)) {
        return check_and_emit(EVENT_SYMLINK, dentry, old_name, false, NULL, false, false);
    }

    // Deep-file coverage: the symlink may be created inside a guarded
    // region whose parent is not in the map.
    struct dentry *parent = get_dentry_from_path((void *)ctx[0]);
    if (parent && guarded_ancestor_within_limit(parent)) {
        return check_and_emit(EVENT_SYMLINK, dentry, old_name, false, NULL, false, false);
    }

    // 2. Check if the symlink TARGET points INTO the guarded path
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
            // We reached the end of the guarded path.
            // It is a match if the symlink target ends here exactly,
            // continues as a sub-directory, or if the stored path had a trailing slash.
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

	if (read_inode_guard(inode)) {
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry, false, false);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[1]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry, false, false);
	}

	// Deep-file coverage on the destination side: linking INTO a guarded
	// region whose parent is not in the map.
	struct dentry *new_parent = get_dentry_from_path((void *)ctx[1]);
	if (new_parent && guarded_ancestor_within_limit(new_parent))
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry, false, false);

	return 0;
}

SEC("lsm/path_mkdir")
int guard_path_mkdir(unsigned long long *ctx)
{
	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	struct dentry *dentry = (struct dentry *)ctx[1];

	if (read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_MKDIR, dentry, NULL, false, NULL, false, false);
	}

	// Deep-file coverage: mkdir inside a guarded region whose parent is
	// not in the map.
	struct dentry *parent = get_dentry_from_path((void *)ctx[0]);
	if (parent && guarded_ancestor_within_limit(parent))
		return check_and_emit(EVENT_MKDIR, dentry, NULL, false, NULL, false, false);

	return 0;
}

// ------------------------------------------------------------------
// Path/attribute-based bypass coverage.  These operations never
// create a struct file, so file_open / file_permission never fire:
//   - truncate(2)            -> security_path_truncate (path_truncate)
//   - chmod/chown/utimes     -> notify_change -> inode_setattr
//   - setxattr/removexattr   -> inode_setxattr / inode_removexattr
//   - mknod(2)               -> security_path_mknod (path_mknod)
//   - rmdir(2)               -> security_path_rmdir (path_rmdir)
// ------------------------------------------------------------------

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
	// ctx[0] = mnt_idmap, ctx[1] = dentry, ctx[2] = iattr.
	// Skip size changes: truncate events are owned by path_truncate /
	// file_truncate, which run before notify_change and already deny;
	// emitting again would duplicate events (ATTR_FILE follows
	// ftruncate(2), also covered by file_truncate).
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
	if (parent_inode && read_inode_guard(parent_inode))
		return check_and_emit(EVENT_MKNOD, dentry, NULL, false, NULL, false, false);

	struct dentry *parent = get_dentry_from_path((void *)ctx[0]);
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

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	if (parent_inode && read_inode_guard(parent_inode))
		return check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);

	struct dentry *parent = get_dentry_from_path((void *)ctx[0]);
	if (parent && guarded_ancestor_within_limit(parent))
		return check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL, false, false);

	return 0;
}

// Observe-only: never blocks.  stat(2)/statx(2) are probe-only and
// ubiquitous — most tools (e.g. coreutils cat/mkdir) stat a guarded file
// before their real syscall, so denying here would pre-empt file_open and
// turn every blocked read into a STAT event instead of the OPEN event the
// enforcement is defined by.  Real enforcement stays in file_open and the
// path_*/inode_* hooks; metadata probing is observable via monitor mode.
SEC("lsm/inode_getattr")
int guard_inode_getattr(unsigned long long *ctx)
{
	return 0;
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

// Block pure inode-level may-checks — access(2), faccessat(2) — against
// guarded files.  On this kernel every syscall path already carries its
// MAY_* bits: access(2) passes mode|MAY_ACCESS, open/exec passes
// MAY_OPEN|acc_mode, truncate passes plain MAY_WRITE (inside
// vfs_truncate).  Requiring MAY_ACCESS therefore selects exactly the
// access(2)-style probes: opens, execs, truncate and setxattr fall
// through to their own dedicated hooks (file_open, path_truncate,
// inode_setxattr), which keep their event identity.  Directory
// traversal is never blocked (a directory's MAY_EXEC must stay
// available for every path walk into the guarded tree).  Allowed
// accesses emit nothing (quiet_allow): this hook fires for every
// guarded-file may-check and would otherwise flood the log.  An X_OK
// probe (MAY_EXEC without MAY_OPEN) leaks no content and is
// deliberately not denied — it must not shadow open/exec handling.
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

	if (!read_inode_guard(inode))
		return 0;

	return check_and_emit(EVENT_STAT, NULL, NULL, false, NULL, false, true);
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

	if (read_inode_guard(inode)) {
		return check_and_emit(EVENT_OPEN, dentry, NULL, false, NULL, false, false);
	}

	return 0;
}

// Block ptrace / process_vm_readv / process_vm_writev against any process
// that has guarded file content in its address space (tainted PID).
// Once a process opens/reads a guarded file, no non-whitelisted process
// can ptrace it or read its memory via process_vm_readv.
SEC("lsm/ptrace_access_check")
int guard_ptrace_access_check(struct task_struct *child, unsigned int mode)
{
	if (!child)
		return 0;

	__u32 pid;
	bpf_probe_read_kernel(&pid, sizeof(pid), &child->tgid);

	__u8 *val = bpf_map_lookup_elem(&guard_tainted_pids, &pid);
	if (!val)
		return 0;  // child is not tainted — allow

	// Child is tainted.  Check if the caller is whitelisted.
	struct inode_key exe_ik = {};
	if (!get_current_exe_inode(&exe_ik))
		return -EPERM;

	__u8 *action = bpf_map_lookup_elem(&guard_exe_actions, &exe_ik);
	if (!action || *action != GUARD_ALLOW)
		return -EPERM;

	return 0;
}

// When a tainted process forks, the child inherits the parent's address
// space (including file mappings).  Mark the child as tainted too so that
// process_vm_readv on the child is also blocked.
SEC("lsm/task_alloc")
int guard_task_alloc(struct task_struct *task, unsigned long clone_flags)
{
	if (!task)
		return 0;

	__u32 parent_pid = bpf_get_current_pid_tgid() >> 32;

	__u8 *val = bpf_map_lookup_elem(&guard_tainted_pids, &parent_pid);
	if (!val)
		return 0;  // parent not tainted

	__u32 child_pid;
	bpf_probe_read_kernel(&child_pid, sizeof(child_pid), &task->tgid);

	__u8 v = 1;
	bpf_map_update_elem(&guard_tainted_pids, &child_pid, &v, BPF_ANY);
	return 0;
}

// Clean up the tainted PID entry when the process exits, preventing
// stale entries when the PID is reused by a different process.
SEC("lsm/task_free")
int guard_task_free(struct task_struct *task)
{
	if (!task)
		return 0;

	__u32 pid;
	bpf_probe_read_kernel(&pid, sizeof(pid), &task->tgid);

	bpf_map_delete_elem(&guard_tainted_pids, &pid);
	return 0;
}
