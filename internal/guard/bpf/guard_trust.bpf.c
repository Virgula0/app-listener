// guard_trust.bpf.c — the trusted-binary / trusted-library object.
//
// A single, daemon-wide object (attached ONCE, not per watched resource) that
// owns guard_trusted_files: the inodes of every whitelisted application binary
// (TRUSTED_BINARY) and every user-writable library the operator explicitly
// trusts for them (TRUSTED_LIB, via allow_lib). Keeping it separate from the
// per-resource guards means it never adds to their verifier-tight programs
// (issue #45) and its handful of programs attach just once, well under the
// per-function trampoline limit.
//
// It closes two gaps that inode-keyed access control cannot, because both act
// on files that live OUTSIDE every guarded tree:
//
//   #1 writer attribution — a whitelisted binary that lives at a USER-WRITABLE
//      path (a home-directory app: Claude, Discord) may only be replaced or
//      modified by another whitelisted binary (the app's own updater), never
//      by unrelated same-user malware. Root-owned system binaries (/usr/bin/…)
//      are deliberately NOT protected: a normal user cannot modify them anyway,
//      and protecting them would break the package manager's own upgrades.
//
//   #2 library load allowlist — a whitelisted process may only map executable
//      code from a library that a non-whitelisted process could not have
//      written, because that is exactly what makes the library trustworthy. A
//      library qualifies when it is:
//        - an explicit TRUSTED_LIB (allow_lib), or
//        - a root-owned file in a root-only-writable directory — an ordinary
//          system library under /usr/lib, /opt, … (auto-trusted; no operator
//          config, and it survives package updates), or
//        - inside a GUARDED resource tree — a directory the daemon watches in
//          whitelist mode, where only whitelisted binaries may create or modify
//          files. A library there cannot have been planted by an attacker, so
//          it is safe to load even though the path is user-writable in the
//          filesystem sense. This is how self-contained apps with bundled,
//          per-launch libraries (Steam/Proton/Wine) are handled: guard their
//          library directories like any other resource and their libraries
//          become loadable with no per-file allow_lib and no exceptions.
//      LD_PRELOAD of an attacker .so under /tmp or an unguarded $HOME path is
//      refused: it is neither system-owned nor inside a guarded tree.
//
// Both #1 and #2 are ALWAYS enforced once the programs attach — there is no
// observe/enforce toggle. They are best-effort to attach; a rejected program
// never affects the per-resource enforcement guards.
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <errno.h>

#define MAX_PATH 256
#define PROT_EXEC 0x4
#define FMODE_WRITE 0x2
#define S_IWGRP 00020
#define S_IWOTH 00002

// Open flags used for provenance. Only flags that GUARANTEE the open created
// the file are accepted: O_CREAT|O_EXCL fails if the name already exists, and
// __O_TMPFILE always makes a brand-new anonymous inode. A bare O_CREAT is
// deliberately NOT enough — it succeeds on a pre-existing file, so honouring
// it would let an attacker's planted .so acquire this process's provenance.
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

// guard_trusted_dirs holds the (dev, ino) of every guarded resource ROOT (the
// directories the daemon watches in whitelist mode). A library whose path
// passes through one of these roots lives in a write-protected tree — only
// whitelisted binaries may create or modify files there — so it is trustworthy
// to load. Populated by userspace from the daemon's resource list.
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1024);
	__type(key, struct inode_key);
	__type(value, __u8);
} guard_trusted_dirs SEC(".maps");

// guard_jit_origin is the provenance ledger for runtime-generated code. GPU
// drivers (NVIDIA above all) JIT shaders and kernels into a file they create
// themselves and then map executable: a memfd, an O_TMPFILE, or an
// mkstemp()'d file. Such an inode is anonymous, user-owned and named at
// random, so none of the three ordinary trust sources can ever cover it —
// but "this very process made it, and nothing else has written it" is a
// property the kernel can attest, and it is strictly stronger than a path
// rule: an attacker who hands a whitelisted process a memfd of their own
// (the LD_PRELOAD=/proc/self/fd/N trick) is a DIFFERENT thread group, so the
// mapping is refused.
//
// An entry is created only by a whitelisted process performing an open that
// is guaranteed to be a creation (see the flag macros above) and is tainted
// for good as soon as any other thread group opens that inode for writing.
//
// Residual risk, stated plainly: the taint is driven by file_open, so foreign
// bytes can only slip in through a writable fd the creator itself handed out
// (SCM_RIGHTS or a fork that outlives the record) — never through an attacker
// opening the inode, which is the reachable shape of the attack. A whitelisted
// binary that deliberately passes a writable fd to its own JIT arena is a
// confused deputy in the same class as any other whitelist entry.
//
// The map is an LRU on purpose. Eviction loses provenance, and no provenance
// means "deny" in every consumer, so pressure degrades towards refusing JIT
// mappings — never towards trusting a stale inode number that the kernel has
// since recycled for somebody else's file.

// jit_origin identifies the creator by its whole process IMAGE, not just its pid.
// execve() keeps the thread group id while replacing the program, so a thread
// group alone would let a whitelisted process create an inode and then exec a
// different whitelisted binary that inherits the fd — provenance would follow
// the pid across an image it never produced. Requiring the mm (allocated
// afresh by execve, distinct per fork) and the executable's inode as well
// makes the check "the same program, in the same thread group, running the
// same image that created this inode".
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

// fill_inode_key resolves inode's (dev, ino) with the SAME encoding the rest of
// the guard uses (dev from sb->s_dev). Returns 0 if unreadable.
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

// current_exe_flags returns the trusted-file flags of the current process's
// main executable (mm->exe_file), 0 for kernel threads or an unresolvable exe.
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

// caller_is_app reports whether the current process's executable is one of the
// whitelisted application binaries (a valid modifier/updater for #1).
static __always_inline int caller_is_app(void)
{
	return (current_exe_flags() & TRUSTED_BINARY) != 0;
}

// inode_is_root_ro reports whether inode is owned by root and cannot be
// modified in place by a normal user: not other-writable, and group-writable
// only when the group is root itself (0775 root:root is common in shipped
// packages — e.g. VS Code's bundled .node modules — and a member of group
// root is already root-equivalent, so it grants no non-root user anything).
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

// is_system_trusted reports whether the file at dentry is a trusted system file:
// root-owned and not group/other-writable, AND sitting directly in a directory
// that is likewise root-owned and not group/other-writable. That two-level
// check is enough to distinguish /usr/lib, /opt/… (auto-trusted) from anything
// a non-root user could replace — a user-writable directory lets its owner swap
// even a root-owned file by unlink+create, so the parent must be root-owned too.
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

// under_guarded_tree reports whether the file at dentry sits inside a guarded
// resource tree: an ancestor (or the file itself) is a guarded root in
// guard_trusted_dirs. The whole tree shares the file's device (a guarded root
// is an ancestor on the same filesystem), so the device is read once and only
// the inode number is compared per ancestor — the same shape as the
// per-resource guard's root_in_chain walk.
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

// fill_current_image resolves the running process image: thread group, mm, and
// the inode of the executable together with its trusted flags. Returns 0 when
// any of it is unreadable (a kernel thread, an exiting task), which every
// caller treats as "no provenance".
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

// jit_record_creation claims inode for the CURRENT process image. Called only
// for an open/allocation that certainly created the inode, and only when the
// creator is a whitelisted binary — an attacker's own files therefore have no
// entry at all, which is what makes a missing entry safe to treat as "deny".
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

// jit_taint_foreign_write marks an inode as no longer self-produced because
// something other than the exact image that created it opened it for writing.
// This is what closes the /proc/<pid>/fd/N route: a same-uid attacker can
// reach a whitelisted process's memfd through procfs, but only by opening it,
// and the open is what poisons the mapping. The taint is permanent — an inode
// that has held foreign bytes never becomes self-produced again.
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

// jit_self_created reports whether inode was created by the process image
// asking to map it and has never been opened for writing by another thread
// group. Absent provenance answers no, so an evicted entry, a file created
// before the daemon started and anything the daemon never saw created all stay
// denied.
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

// #2 — library load allowlist. A whitelisted binary may map executable code
// only from a trusted library (explicit allow_lib) or an auto-trusted system
// file (root-owned, root-only-writable dir). An untrusted exec-map is the exact
// shape of an LD_PRELOAD / dlopen of attacker code; it is logged always and
// denied when enforce_libs is set.
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

// --- #1: writer attribution for user-writable whitelisted binaries ---------
// A protected target is a TRUSTED_BINARY/TRUSTED_LIB inode that is NOT an
// auto-trusted system file — i.e. it lives at a user-writable path, where
// same-user malware could otherwise replace it and have the daemon re-admit
// the replacement. Only another whitelisted binary (the app's own updater) may
// modify it; every other caller is denied. Root-owned system binaries are not
// protected here (the package manager must be able to upgrade them, and a
// non-root attacker cannot touch them anyway).
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
	// Deny renaming a protected binary AWAY (old_dentry) and renaming anything
	// OVER a protected binary (new_dentry victim) by a non-app caller.
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

	// Provenance bookkeeping for runtime-generated code. Recording comes
	// first: a creating open is a write-open too, and the creator must not
	// taint its own brand-new inode.
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

// memfd_create() allocates its file without ever calling security_file_open(),
// so a memfd's provenance has to be recorded where the kernel makes it. This
// is the first link in NVIDIA's JIT fallback chain (memfd → O_TMPFILE →
// mkstemp), and the only one an LSM hook cannot see.
//
// Best-effort like every other trust program: on a kernel where the symbol is
// inlined or renamed the attach fails, memfd exec-maps keep being denied, and
// the driver falls through to its file-backed fallbacks, which file_open does
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
