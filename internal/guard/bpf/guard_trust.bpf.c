// guard_trust.bpf.c: the trusted-binary / trusted-library object.
//
// One daemon-wide object (attached ONCE, not per resource) owning guard_trusted_files: inodes of
// every whitelisted application binary (TRUSTED_BINARY) and every user-writable library the
// operator trusts for them (TRUSTED_LIB, via allow_lib). Kept separate so it adds nothing to the
// per-resource guards' verifier-tight programs (issue #45) and attaches once, under the
// per-function trampoline limit.
//
// It closes two gaps inode-keyed access control can't, both on files OUTSIDE every guarded tree:
//   #1 writer attribution: a whitelisted binary at a USER-WRITABLE path (home-directory apps: Claude, Discord) may only be replaced/modified by another whitelisted binary (the app's updater), never by same-user malware. Root-owned system binaries (/usr/bin/...) are deliberately NOT protected: users can't modify them, and protecting them would break package upgrades.
//   #2 library load allowlist: a whitelisted process may map executable code only from a library a non-whitelisted process couldn't have written. A library qualifies if it is:
//     - an explicit TRUSTED_LIB (allow_lib), or
//     - a root-owned file in a root-only-writable directory (system libs under /usr/lib, /opt; auto-trusted, survives package updates), or
//     - inside a GUARDED resource tree (whitelist mode: only whitelisted binaries create/modify files there, so nothing can be planted). This is how bundled per-launch libraries (Steam/Proton/Wine) work: guard their library dirs and they become loadable with no per-file allow_lib.
//   LD_PRELOAD of an attacker .so under /tmp or an unguarded $HOME path is refused.
//
// Both are ALWAYS enforced once attached (no observe/enforce toggle). Attach is best-effort; a
// rejected program never affects the per-resource guards.
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <errno.h>

#define MAX_PATH 256
#define PROT_EXEC 0x4
#define FMODE_WRITE 0x2
#define S_IWGRP 00020
#define S_IWOTH 00002

// Open flags for provenance. Only flags that GUARANTEE creation are accepted: O_CREAT|O_EXCL fails
// if the name exists, __O_TMPFILE always makes a new anonymous inode. Bare O_CREAT is NOT enough:
// it succeeds on an existing file, so an attacker's planted .so could acquire this process's
// provenance.
#define O_CREAT 00000100
#define O_EXCL 00000200
#define __O_TMPFILE 020000000

// IS_ERR_VALUE: the top 4 KiB of the address space is the kernel's error range.
#define BPF_IS_ERR_PTR(p) ((unsigned long)(p) > (unsigned long)-4096)

// guard_trusted_files value flags.
#define TRUSTED_BINARY 1 // a whitelisted application binary
#define TRUSTED_LIB 2    // a user-writable library explicitly trusted (allow_lib)

// trust_event.kind
#define TRUST_LIBLOAD 0    // a whitelisted binary mapped an untrusted library
#define TRUST_WRITEBLOCK 1 // a non-app tried to modify a protected binary

struct inode_key {
	__u64 dev;
	__u64 ino;
};

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 8192);
	__type(key, struct inode_key);
	__type(value, __u8);
} guard_trusted_files SEC(".maps");

// guard_trusted_dirs holds (dev, ino) of every guarded resource ROOT (dirs watched in whitelist
// mode). A library whose path passes through one lives in a write-protected tree (only whitelisted
// binaries may create/modify there), so it is trustworthy. Populated by userspace from the daemon's
// resource list.
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1024);
	__type(key, struct inode_key);
	__type(value, __u8);
} guard_trusted_dirs SEC(".maps");

// guard_jit_origin is the provenance ledger for runtime-generated code. GPU drivers (NVIDIA above
// all) JIT into a file they create themselves (memfd, O_TMPFILE, mkstemp) and map it executable.
// Such an inode is anonymous, user-owned and randomly named, so no ordinary trust source covers it,
// but "this very process made it and nothing else has written it" is kernel-attestable and stronger
// than a path rule: an attacker handing a whitelisted process their own memfd
// (LD_PRELOAD=/proc/self/fd/N) is a DIFFERENT thread group, so the mapping is refused.
//
// An entry is created only by a whitelisted process doing a guaranteed-creation open (flag macros
// above) and is tainted for good once any other thread group opens the inode for writing.
//
// Residual risk: taint is driven by file_open, so foreign bytes can only arrive through a writable
// fd the creator itself handed out (SCM_RIGHTS, or a fork outliving the record), never by an
// attacker opening the inode. A whitelisted binary that passes a writable fd to its own JIT arena
// is a confused deputy, like any other whitelist entry.
//
// LRU on purpose: eviction loses provenance, and no provenance means "deny" everywhere, so pressure
// degrades toward refusing JIT mappings, never toward trusting a recycled inode number.

// jit_origin identifies the creator by its whole process IMAGE, not just pid: execve() keeps the
// thread group id while replacing the program, so a tgid alone would let provenance follow the pid
// into a different whitelisted binary that inherits the fd. Requiring the mm (fresh per execve,
// distinct per fork) and the exe inode too means "the same program, thread group and image that
// created this inode".
struct jit_origin {
	__u64 mm;
	struct inode_key exe;
	__u32 tgid;
	__u8 tainted;
};

struct {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, 16384);
	__type(key, struct inode_key);
	__type(value, struct jit_origin);
} guard_jit_origin SEC(".maps");

struct trust_event {
	__u32 pid;
	__u32 uid;
	__u32 kind;
	char comm[16];
	char path[MAX_PATH];
};

struct {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, char[MAX_PATH]);
} trust_tmp SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 18); // 256 KiB — observe/deny events are low-rate
} trust_rb SEC(".maps");

char LICENSE[] SEC("license") = "GPL";

// fill_inode_key resolves inode's (dev, ino) with the same encoding as the rest of the guard (dev
// from sb->s_dev). 0 if unreadable.
static __always_inline int fill_inode_key(struct inode *inode, struct inode_key *k)
{
	if (!inode)
		return 0;
	bpf_probe_read_kernel(&k->ino, sizeof(k->ino), &inode->i_ino);
	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	dev_t dev = 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	k->dev = dev;
	return 1;
}

static __always_inline __u8 trusted_flags_of(struct inode *inode)
{
	struct inode_key k = {};
	if (!fill_inode_key(inode, &k))
		return 0;
	__u8 *f = bpf_map_lookup_elem(&guard_trusted_files, &k);
	return f ? *f : 0;
}

// current_exe_flags returns the trusted-file flags of the current process's main exe
// (mm->exe_file); 0 for kernel threads or an unresolvable exe.
static __always_inline __u8 current_exe_flags(void)
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
	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &exe_file->f_inode);
	return trusted_flags_of(inode);
}

// caller_is_app: the current exe is a whitelisted application binary (a valid modifier/updater for
// #1).
static __always_inline int caller_is_app(void)
{
	return (current_exe_flags() & TRUSTED_BINARY) != 0;
}

// inode_is_root_ro: owned by root and not modifiable in place by a normal user: not other-writable,
// and group-writable only if the group is root (0775 root:root is common in packages, e.g. VS
// Code's .node modules, and group root is already root-equivalent).
static __always_inline int inode_is_root_ro(struct inode *inode)
{
	if (!inode)
		return 0;
	__u32 uid = 0xffffffff;
	bpf_probe_read_kernel(&uid, sizeof(uid), &inode->i_uid); // kuid_t { uid_t val; }
	if (uid != 0)
		return 0;
	umode_t mode = 0;
	bpf_probe_read_kernel(&mode, sizeof(mode), &inode->i_mode);
	if (mode & S_IWOTH)
		return 0;
	if (mode & S_IWGRP) {
		__u32 gid = 0xffffffff;
		bpf_probe_read_kernel(&gid, sizeof(gid), &inode->i_gid); // kgid_t { gid_t val; }
		if (gid != 0)
			return 0;
	}
	return 1;
}

// is_system_trusted: the file is root-owned and not group/other-writable AND sits directly in a
// directory that is likewise. Two levels are enough to tell /usr/lib, /opt/... (auto-trusted) from
// anything a non-root user could replace: a user-writable directory lets its owner swap even a
// root-owned file by unlink+create, so the parent must be root-owned too.
static __always_inline int is_system_trusted(struct inode *inode, struct dentry *dentry)
{
	if (!inode_is_root_ro(inode))
		return 0;
	if (!dentry)
		return 0;
	struct dentry *parent;
	bpf_probe_read_kernel(&parent, sizeof(parent), &dentry->d_parent);
	if (!parent || parent == dentry)
		return 0;
	struct inode *pinode;
	bpf_probe_read_kernel(&pinode, sizeof(pinode), &parent->d_inode);
	return inode_is_root_ro(pinode);
}

// under_guarded_tree: an ancestor (or the file itself) is a guarded root in guard_trusted_dirs. The
// tree shares the file's device, so the device is read once and only inode numbers are compared per
// ancestor (same shape as root_in_chain).
static __always_inline int under_guarded_tree(struct inode *inode, struct dentry *dentry)
{
	if (!inode || !dentry)
		return 0;
	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	dev_t dev = 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);

	struct inode_key k = {};
	k.dev = dev;

	struct dentry *d = dentry;
	for (int i = 0; i < 32; i++) {
		if (!d)
			break;
		struct inode *di;
		bpf_probe_read_kernel(&di, sizeof(di), &d->d_inode);
		if (di) {
			bpf_probe_read_kernel(&k.ino, sizeof(k.ino), &di->i_ino);
			if (bpf_map_lookup_elem(&guard_trusted_dirs, &k))
				return 1;
		}
		struct dentry *parent;
		bpf_probe_read_kernel(&parent, sizeof(parent), &d->d_parent);
		if (!parent || parent == d)
			break;
		d = parent;
	}
	return 0;
}

// --- provenance ledger for runtime-generated code (guard_jit_origin) -------

// fill_current_image resolves the running image: tgid, mm, and exe inode with its trusted flags. 0
// if any is unreadable (kernel thread, exiting task); every caller treats that as "no provenance".
struct proc_image {
	__u64 mm;
	struct inode_key exe;
	__u32 tgid;
	__u8 flags;
};

static __always_inline int fill_current_image(struct proc_image *img)
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
	if (!fill_inode_key(exe_inode, &img->exe))
		return 0;
	img->mm = (__u64)mm;
	img->tgid = bpf_get_current_pid_tgid() >> 32;
	__u8 *f = bpf_map_lookup_elem(&guard_trusted_files, &img->exe);
	img->flags = f ? *f : 0;
	return 1;
}

// jit_record_creation claims inode for the CURRENT image. Called only for a certain creation and
// only by a whitelisted creator, so an attacker's files have no entry, which is why a missing entry
// safely means "deny".
static __always_inline void jit_record_creation(struct inode *inode)
{
	struct inode_key k = {};
	if (!fill_inode_key(inode, &k))
		return;
	struct proc_image img = {};
	if (!fill_current_image(&img))
		return;
	if (!(img.flags & TRUSTED_BINARY))
		return;
	struct jit_origin o = {};
	o.mm = img.mm;
	o.exe = img.exe;
	o.tgid = img.tgid;
	bpf_map_update_elem(&guard_jit_origin, &k, &o, BPF_ANY);
}

// jit_taint_foreign_write marks an inode no longer self-produced because something other than the
// exact creating image opened it for writing. This closes the /proc/<pid>/fd/N route: a same-uid
// attacker can reach a whitelisted process's memfd via procfs only by opening it, and the open
// poisons the mapping. Permanent.
static __always_inline void jit_taint_foreign_write(struct inode *inode)
{
	struct inode_key k = {};
	if (!fill_inode_key(inode, &k))
		return;
	struct jit_origin *o = bpf_map_lookup_elem(&guard_jit_origin, &k);
	if (!o)
		return;
	struct proc_image img = {};
	if (!fill_current_image(&img)) {
		o->tainted = 1; // unattributable writer: assume the worst
		return;
	}
	if (o->tgid == img.tgid && o->mm == img.mm &&
	    o->exe.dev == img.exe.dev && o->exe.ino == img.exe.ino)
		return; // the creator writing its own JIT output
	o->tainted = 1;
}

// jit_self_created: inode was created by the image asking to map it and never opened for writing by
// another thread group. Absent provenance answers no (evicted entries, files predating the daemon,
// anything never seen created stay denied).
static __always_inline int jit_self_created(struct inode *inode)
{
	struct inode_key k = {};
	if (!fill_inode_key(inode, &k))
		return 0;
	struct jit_origin *o = bpf_map_lookup_elem(&guard_jit_origin, &k);
	if (!o || o->tainted)
		return 0;
	struct proc_image img = {};
	if (!fill_current_image(&img))
		return 0;
	if (o->tgid != img.tgid || o->mm != img.mm)
		return 0;
	return o->exe.dev == img.exe.dev && o->exe.ino == img.exe.ino;
}

static __always_inline struct dentry *path_dentry(void *ctx_path)
{
	struct path *p = (struct path *)ctx_path;
	if (!p)
		return NULL;
	struct dentry *d = NULL;
	bpf_probe_read_kernel(&d, sizeof(d), &p->dentry);
	return d;
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
		int start_pos = pos & (MAX_PATH - 1);
		for (int j = 0; j < 64; j++) {
			if (j >= name_len)
				break;
			buf[(start_pos + j) & (MAX_PATH - 1)] = tmp[j & 63];
		}
		if (pos <= 0)
			break;
		pos--;
		buf[pos & (MAX_PATH - 1)] = '/';
		d = parent;
	}
	return pos;
}

static __always_inline void fill_path(struct dentry *dentry, char *out)
{
	if (!dentry) {
		out[0] = '\0';
		return;
	}
	__u32 key = 0;
	char *buf = bpf_map_lookup_elem(&trust_tmp, &key);
	if (!buf) {
		out[0] = '\0';
		return;
	}
	long off = read_path(dentry, buf, MAX_PATH);
	if (off < MAX_PATH) {
		off &= (MAX_PATH - 1);
		bpf_probe_read_kernel_str(out, MAX_PATH, buf + off);
	} else {
		out[0] = '\0';
	}
}

static __always_inline void emit(struct dentry *dentry, __u32 kind)
{
	struct trust_event *e = bpf_ringbuf_reserve(&trust_rb, sizeof(*e), 0);
	if (!e)
		return;
	e->pid = bpf_get_current_pid_tgid() >> 32;
	e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e->kind = kind;
	bpf_get_current_comm(e->comm, sizeof(e->comm));
	fill_path(dentry, e->path);
	bpf_ringbuf_submit(e, 0);
}

// #2 library load allowlist: a whitelisted binary may map executable code only from a trusted
// library (allow_lib) or an auto-trusted system file (root-owned, root-only-writable dir). An
// untrusted exec-map is the shape of an LD_PRELOAD/dlopen of attacker code: always logged, denied
// when enforce_libs is set.
SEC("lsm/mmap_file")
int trust_mmap(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	unsigned long prot = (unsigned long)ctx[2];
	if (!file)
		return 0;
	if (!(prot & PROT_EXEC))
		return 0;
	if (!(current_exe_flags() & TRUSTED_BINARY))
		return 0; // only whitelisted-application processes are subject to it

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	__u8 flags = trusted_flags_of(inode);
	if (flags & (TRUSTED_LIB | TRUSTED_BINARY))
		return 0; // explicitly trusted (or the binary re-mapping itself)

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (is_system_trusted(inode, dentry))
		return 0; // auto-trusted root-owned system library
	if (under_guarded_tree(inode, dentry))
		return 0; // inside a write-protected guarded tree — attacker cannot plant here
	if (jit_self_created(inode))
		return 0; // runtime-generated code: this process made it, nothing else wrote it

	emit(dentry, TRUST_LIBLOAD);
	return -EPERM;
}

// --- #1: writer attribution for user-writable whitelisted binaries ---
// A protected target is a TRUSTED_BINARY/TRUSTED_LIB inode that is NOT an auto-trusted system file,
// i.e. at a user-writable path where same-user malware could replace it and have the daemon
// re-admit the replacement. Only another whitelisted binary (the app's updater) may modify it.
// Root-owned system binaries are exempt (package manager upgrades; non-root can't touch them).
static __always_inline int protected_writable_target(struct inode *inode, struct dentry *dentry)
{
	if (!(trusted_flags_of(inode) & (TRUSTED_BINARY | TRUSTED_LIB)))
		return 0;
	if (is_system_trusted(inode, dentry))
		return 0; // root-owned system binary: not our concern
	return 1;
}

static __always_inline int deny_if_protected(struct inode *inode, struct dentry *dentry)
{
	if (!protected_writable_target(inode, dentry))
		return 0;
	if (caller_is_app())
		return 0; // the app's own updater may replace it
	emit(dentry, TRUST_WRITEBLOCK);
	return -EPERM;
}

SEC("lsm/path_unlink")
int trust_path_unlink(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;
	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);
	return deny_if_protected(inode, dentry);
}

SEC("lsm/path_rename")
int trust_path_rename(unsigned long long *ctx)
{
	// Deny a non-app caller renaming a protected binary AWAY (old_dentry) or anything OVER one
	// (new_dentry victim).
	struct dentry *old_dentry = (struct dentry *)ctx[1];
	struct dentry *new_dentry = (struct dentry *)ctx[3];

	if (old_dentry) {
		struct inode *oi;
		bpf_probe_read_kernel(&oi, sizeof(oi), &old_dentry->d_inode);
		int r = deny_if_protected(oi, old_dentry);
		if (r)
			return r;
	}
	if (new_dentry) {
		struct inode *vi;
		bpf_probe_read_kernel(&vi, sizeof(vi), &new_dentry->d_inode);
		int r = deny_if_protected(vi, new_dentry);
		if (r)
			return r;
	}
	return 0;
}

SEC("lsm/file_open")
int trust_file_open(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	if (!file)
		return 0;
	__u32 f_mode = 0;
	bpf_probe_read_kernel(&f_mode, sizeof(f_mode), &file->f_mode);
	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	// Provenance bookkeeping first: a creating open is also a write-open, and the creator must not
	// taint its own new inode.
	__u32 f_flags = 0;
	bpf_probe_read_kernel(&f_flags, sizeof(f_flags), &file->f_flags);
	if ((f_flags & __O_TMPFILE) || ((f_flags & O_CREAT) && (f_flags & O_EXCL)))
		jit_record_creation(inode); // no-op unless the creator is whitelisted
	if (f_mode & FMODE_WRITE)
		jit_taint_foreign_write(inode);

	if (!(f_mode & FMODE_WRITE))
		return 0; // only a write-open can modify the binary in place
	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (!dentry)
		return 0;
	return deny_if_protected(inode, dentry);
}

// memfd_create() allocates its file without security_file_open(), so a memfd's provenance is
// recorded where the kernel makes it: the first link in NVIDIA's JIT fallback chain (memfd,
// O_TMPFILE, mkstemp) and the only one an LSM hook can't see.
//
// Best-effort like every trust program: if the symbol is inlined/renamed the attach fails, memfd
// exec-maps stay denied, and the driver falls through to file-backed fallbacks that file_open does
// see.
SEC("fexit/memfd_alloc_file")
int trust_memfd_alloc(unsigned long long *ctx)
{
	// fexit context: the traced function's arguments, then its return value.
	struct file *file = (struct file *)ctx[2];
	if (!file || BPF_IS_ERR_PTR(file))
		return 0;
	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	jit_record_creation(inode); // no-op unless the creator is whitelisted
	return 0;
}

SEC("lsm/path_truncate")
int trust_path_truncate(unsigned long long *ctx)
{
	struct dentry *dentry = path_dentry((void *)ctx[0]);
	if (!dentry)
		return 0;
	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);
	return deny_if_protected(inode, dentry);
}
