#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>
#include "inode_dev.h"

#define MAX_PATH 256
#define MON_PROT_EXEC 0x4

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
	// EVENT_EXEC names a new executable (path = its file, dest = bprm->filename) so userspace can
	// map exe_dev/exe_ino to a path after a short-lived process exited. Never a file access.
	EVENT_EXEC,
};

struct event {
	__u32 pid;
	__u32 uid;
	__u32 gid;
	__u32 type;
	__u32 fd;
	__u32 pad;
	char comm[16];
	char path[MAX_PATH];
	char dest[MAX_PATH];
	// Identity of the dentry the kernel paths were walked from; dev == 0 means path/dest are
	// user strings. read_path stops at the superblock root, so userspace maps (dev, path) to
	// the mount through /proc/self/mountinfo, whose major:minor is this same kernel s_dev.
	__u32 dev;
	__u32 mode;
	__u64 ino;
	// The calling process's executable (mm->exe_file), keyed as the guard keys it (inode_dev);
	// 0/0 for kernel threads.
	__u64 exe_dev;
	__u64 exe_ino;
	// prot of an mmap/mprotect event (MON_PROT_EXEC = a code mapping); 0 otherwise.
	__u32 prot;
	__u32 pad2;
};

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} rb SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, char[MAX_PATH]);
} tmp_buf SEC(".maps");

static __always_inline void fill_identity(struct event *e, struct dentry *dentry)
{
	e->dev = 0;
	e->mode = 0;
	e->ino = 0;
	if (!dentry)
		return;
	e->dev = BPF_CORE_READ(dentry, d_sb, s_dev);
	struct inode *inode = BPF_CORE_READ(dentry, d_inode);
	if (!inode)
		return;
	e->ino = BPF_CORE_READ(inode, i_ino);
	e->mode = BPF_CORE_READ(inode, i_mode);
}

struct exe_id {
	__u64 dev;
	__u64 ino;
};

// current_exe reads the caller's exe identity. Called before the ringbuf reservation: inode_dev is
// a global function.
static __always_inline void current_exe(struct exe_id *id)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	struct mm_struct *mm = NULL;
	struct file *exe = NULL;
	struct inode *inode = NULL;

	if (!task || bpf_probe_read_kernel(&mm, sizeof(mm), &task->mm) || !mm)
		return;
	if (bpf_probe_read_kernel(&exe, sizeof(exe), &mm->exe_file) || !exe)
		return;
	if (bpf_probe_read_kernel(&inode, sizeof(inode), &exe->f_inode) || !inode)
		return;
	bpf_probe_read_kernel(&id->ino, sizeof(id->ino), &inode->i_ino);
	id->dev = inode_dev((__u64)inode);
}

static __always_inline struct event *reserve_event(__u32 type, __u32 fd, __u32 prot)
{
	struct exe_id id = {};
	current_exe(&id);

	struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e)
		return NULL;

	e->pid = bpf_get_current_pid_tgid() >> 32;
	e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e->gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e->type = type;
	e->fd = fd;
	e->pad = 0;
	e->exe_dev = id.dev;
	e->exe_ino = id.ino;
	e->prot = prot;
	e->pad2 = 0;
	bpf_get_current_comm(e->comm, sizeof(e->comm));
	return e;
}

static __always_inline int emit_user(__u32 type, __u32 fd, const char *path, const char *dest, __u32 prot)
{
	struct event *e = reserve_event(type, fd, prot);
	if (!e)
		return 0;

	fill_identity(e, NULL);

	if (path)
		bpf_probe_read_user_str(e->path, sizeof(e->path), path);
	else
		e->path[0] = '\0';

	if (dest)
		bpf_probe_read_user_str(e->dest, sizeof(e->dest), dest);
	else
		e->dest[0] = '\0';

	bpf_ringbuf_submit(e, 0);
	return 0;
}

static __always_inline int emit_event(__u32 type, __u32 fd, const char *path, const char *dest)
{
	return emit_user(type, fd, path, dest, 0);
}

static __always_inline int emit_kern(__u32 type, __u32 fd, const char *path, const char *dest,
				     struct dentry *dentry, __u32 prot)
{
	struct event *e = reserve_event(type, fd, prot);
	if (!e)
		return 0;

	fill_identity(e, dentry);

	if (path)
		bpf_probe_read_kernel_str(e->path, sizeof(e->path), path);
	else
		e->path[0] = '\0';

	if (dest)
		bpf_probe_read_kernel_str(e->dest, sizeof(e->dest), dest);
	else
		e->dest[0] = '\0';

	bpf_ringbuf_submit(e, 0);
	return 0;
}

static __always_inline int emit_event_kern(__u32 type, __u32 fd, const char *path, const char *dest,
					   struct dentry *dentry)
{
	return emit_kern(type, fd, path, dest, dentry, 0);
}

// VFS-level kprobes: catch open/read/write from the common syscall families (normal syscalls,
// execve, sendfile, splice, copy_file_range, open_by_handle_at). io_uring bypasses the vfs_*
// helpers and is handled by the io_read/io_write kprobes below.
//
// The dentry tree is walked manually because bpf_d_path isn't available in non-sleepable
// kprobe/fentry programs (vfs_open/read/write aren't sleepable).
//
// vfs_open: ALL opens. vfs_read: ALL reads. vfs_write: ALL writes.

static __always_inline long read_path(struct dentry *dentry, char *buf, int buf_size)
{
	char tmp[64];
	int pos = buf_size - 1;
	buf[pos & (MAX_PATH - 1)] = '\0';

	// A superblock root (e.g. a tmpfs mountpoint) is named "/": emit it. Pipe, socket and
	// anon-inode dentries are parentless too but named otherwise; they stay dropped.
	struct dentry *root_parent;
	struct qstr root_q;
	if (dentry && !bpf_probe_read_kernel(&root_parent, sizeof(root_parent), &dentry->d_parent) &&
	    root_parent == dentry &&
	    !bpf_probe_read_kernel(&root_q, sizeof(root_q), &dentry->d_name) && root_q.len == 1) {
		char root_name = 0;
		bpf_probe_read_kernel(&root_name, 1, root_q.name);
		if (root_name == '/') {
			buf[(pos - 1) & (MAX_PATH - 1)] = '/';
			return pos - 1;
		}
	}

	#pragma unroll
	for (int i = 0; i < 48; i++) {
		if (!dentry)
			break;

		struct dentry *d_parent;
		if (bpf_probe_read_kernel(&d_parent, sizeof(d_parent), &dentry->d_parent))
			break;
		if (d_parent == dentry || !d_parent)
			break;

		struct qstr d_name;
		if (bpf_probe_read_kernel(&d_name, sizeof(d_name), &dentry->d_name))
			break;

		const unsigned char *name_ptr;
		if (bpf_probe_read_kernel(&name_ptr, sizeof(name_ptr), &d_name.name))
			break;

		long ret = bpf_probe_read_kernel_str(tmp, sizeof(tmp), name_ptr);
		if (ret <= 1)
			break;

		int name_len = ret - 1;

		if (pos < name_len)
			break;

		pos -= name_len;

		#pragma unroll
		for (int j = 0; j < 64; j++) {
			if (j >= name_len) break;
			buf[(pos + j) & (MAX_PATH - 1)] = tmp[j];
		}

		if (pos <= 0)
			break;

		pos--;
		buf[pos & (MAX_PATH - 1)] = '/';

		dentry = d_parent;
	}

	return pos < buf_size - 1 ? pos : -1;
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

	if (off >= 0 && off < MAX_PATH) {
		bpf_probe_read_kernel_str(out, MAX_PATH, buf + off);
	} else {
		out[0] = '\0';
	}
}

SEC("kprobe/vfs_open")
int trace_vfs_open(struct pt_regs *ctx)
{
	struct path *path = (struct path *)PT_REGS_PARM1(ctx);
	if (!path)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &path->dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_OPEN, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_read")
int trace_vfs_read(struct pt_regs *ctx)
{
	struct file *file = (struct file *)PT_REGS_PARM1(ctx);
	if (!file)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_write")
int trace_vfs_write(struct pt_regs *ctx)
{
	struct file *file = (struct file *)PT_REGS_PARM1(ctx);
	if (!file)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_WRITE, 0, buf + off, NULL, dentry);
}

// Additional VFS kprobes for syscalls that bypass vfs_read/vfs_write:
//   vfs_readv             - readv, preadv (modern kernels)
//   vfs_copy_file_range   - copy_file_range
//   do_splice             - splice (when in is a file)
//   do_splice_direct      - sendfile (via do_sendfile -> do_splice_direct)
//   security_mmap_file    - mmap (kprobe alternative to tracepoint)

SEC("kprobe/vfs_readv")
int trace_vfs_readv(struct pt_regs *ctx)
{
	struct file *file = (struct file *)PT_REGS_PARM1(ctx);
	if (!file)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_copy_file_range")
int trace_vfs_copy_file_range(struct pt_regs *ctx)
{
	struct file *file_in = (struct file *)PT_REGS_PARM1(ctx);
	if (!file_in)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file_in->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

SEC("kprobe/do_splice")
int trace_do_splice(struct pt_regs *ctx)
{
	struct file *in = (struct file *)PT_REGS_PARM1(ctx);
	if (!in)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &in->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

SEC("kprobe/do_splice_direct")
int trace_do_splice_direct(struct pt_regs *ctx)
{
	struct file *in = (struct file *)PT_REGS_PARM1(ctx);
	if (!in)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &in->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

SEC("kprobe/splice_file_range")
int trace_splice_file_range(struct pt_regs *ctx)
{
	struct file *in = (struct file *)PT_REGS_PARM1(ctx);
	if (!in)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &in->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

// do_sendfile catches sendfile on kernels where it no longer calls
// do_splice_direct/splice_file_range. The file is resolved from the task's fd table via CO-RE.
SEC("kprobe/do_sendfile")
int trace_do_sendfile(struct pt_regs *ctx)
{
	int in_fd = (int)PT_REGS_PARM2(ctx);

	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (!task)
		return 0;

	struct file **fd_array;
	struct files_struct *files;
	struct fdtable *fdt;
	int max_fds;

	if (bpf_probe_read_kernel(&files, sizeof(files), &task->files))
		return 0;
	if (!files)
		return 0;

	if (bpf_probe_read_kernel(&fdt, sizeof(fdt), &files->fdt))
		return 0;
	if (!fdt)
		return 0;

	if (bpf_probe_read_kernel(&max_fds, sizeof(max_fds), &fdt->max_fds))
		return 0;
	if (in_fd < 0 || in_fd >= max_fds)
		return 0;

	if (bpf_probe_read_kernel(&fd_array, sizeof(fd_array), &fdt->fd))
		return 0;
	if (!fd_array)
		return 0;

	struct file *file;
	if (bpf_probe_read_kernel(&file, sizeof(file), &fd_array[in_fd]))
		return 0;
	if (!file)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

// vfs_iter_read: catch-all for read-like operations on newer kernels.
SEC("kprobe/vfs_iter_read")
int trace_vfs_iter_read(struct pt_regs *ctx)
{
	struct file *file = (struct file *)PT_REGS_PARM1(ctx);
	if (!file)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

// io_uring read/write. io_uring does NOT go through vfs_read/vfs_readv/vfs_iter_read (its path is
// io_read() -> call_read_iter() -> f_op->read_iter(), mirrored for writes), so every vfs_* kprobe
// above misses a READV/WRITEV on a watched file. io_read()/io_write() take the request as arg1;
// io_kiocb.file (first member of an anonymous union) is the already-resolved target. Runs inline in
// the submitter, or in an io-wq worker sharing the owner's tgid, so the pid is the real process
// either way.
SEC("kprobe/io_read")
int trace_io_read(struct pt_regs *ctx)
{
	struct io_kiocb *req = (struct io_kiocb *)PT_REGS_PARM1(ctx);
	if (!req)
		return 0;

	struct file *file;
	if (bpf_probe_read_kernel(&file, sizeof(file), &req->file))
		return 0;
	if (!file)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL, dentry);
}

SEC("kprobe/io_write")
int trace_io_write(struct pt_regs *ctx)
{
	struct io_kiocb *req = (struct io_kiocb *)PT_REGS_PARM1(ctx);
	if (!req)
		return 0;

	struct file *file;
	if (bpf_probe_read_kernel(&file, sizeof(file), &req->file))
		return 0;
	if (!file)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_WRITE, 0, buf + off, NULL, dentry);
}

SEC("kprobe/security_mmap_file")
int trace_security_mmap_file(struct pt_regs *ctx)
{
	struct file *file = (struct file *)PT_REGS_PARM1(ctx);
	__u32 prot = (__u32)PT_REGS_PARM2(ctx);
	if (!file)
		return 0;

	struct dentry *dentry;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_kern(EVENT_MMAP, 0, buf + off, NULL, dentry, prot);
}

// security_file_mprotect: a file mapping made executable after mmap (some loaders map, then
// mprotect PROT_EXEC). Only that transition is emitted.
SEC("kprobe/security_file_mprotect")
int trace_security_file_mprotect(struct pt_regs *ctx)
{
	struct vm_area_struct *vma = (struct vm_area_struct *)PT_REGS_PARM1(ctx);
	__u32 prot = (__u32)PT_REGS_PARM3(ctx);
	if (!vma || !(prot & MON_PROT_EXEC))
		return 0;

	struct file *file = NULL;
	if (bpf_probe_read_kernel(&file, sizeof(file), &vma->vm_file) || !file)
		return 0;
	struct dentry *dentry = NULL;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry) || !dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_kern(EVENT_MMAP, 0, buf + off, NULL, dentry, prot);
}

// Record probes (Monitor.AttachRecordProbes): the guard's own LSM chokepoints for what the vfs_*
// probes above miss, so a recording learns every op the guard will judge. path_mknod: O_CREAT
// opens and regular-file mknod (vfs_mknod sees only devices/FIFOs/sockets). file_permission:
// writev, pwritev, splice/sendfile/copy_file_range destinations. file_truncate: ftruncate and
// O_TRUNC opens, which the guard judges as WRITE. inode_getxattr/listxattr: judged as STAT.
#define MON_MAY_WRITE 0x2

static __always_inline int emit_dentry(__u32 type, struct dentry *dentry)
{
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(type, 0, buf + off, NULL, dentry);
}

SEC("kprobe/security_path_mknod")
int trace_security_path_mknod(struct pt_regs *ctx)
{
	return emit_dentry(EVENT_MKNOD, (struct dentry *)PT_REGS_PARM2(ctx));
}

// stat(2) and fstat: NewMonitor attaches this, not vfs_getattr, which vfs_statx_path inlines on
// some kernels (7.0); vfs_getattr is security_inode_getattr's only caller.
SEC("kprobe/security_inode_getattr")
int trace_security_inode_getattr(struct pt_regs *ctx)
{
	struct path *path = (struct path *)PT_REGS_PARM1(ctx);
	struct dentry *dentry = NULL;
	if (!path || bpf_probe_read_kernel(&dentry, sizeof(dentry), &path->dentry))
		return 0;
	return emit_dentry(EVENT_STAT, dentry);
}

SEC("kprobe/security_inode_getxattr")
int trace_security_inode_getxattr(struct pt_regs *ctx)
{
	return emit_dentry(EVENT_STAT, (struct dentry *)PT_REGS_PARM1(ctx));
}

SEC("kprobe/security_inode_listxattr")
int trace_security_inode_listxattr(struct pt_regs *ctx)
{
	return emit_dentry(EVENT_STAT, (struct dentry *)PT_REGS_PARM1(ctx));
}

static __always_inline int emit_file_write(struct file *file)
{
	struct dentry *dentry = NULL;
	if (!file || bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry))
		return 0;
	return emit_dentry(EVENT_WRITE, dentry);
}

SEC("kprobe/security_file_permission")
int trace_security_file_permission(struct pt_regs *ctx)
{
	if (!((int)PT_REGS_PARM2(ctx) & MON_MAY_WRITE))
		return 0;
	return emit_file_write((struct file *)PT_REGS_PARM1(ctx));
}

SEC("kprobe/security_file_truncate")
int trace_security_file_truncate(struct pt_regs *ctx)
{
	return emit_file_write((struct file *)PT_REGS_PARM1(ctx));
}

// sched_process_exec runs after the exec succeeded: current's exe identity is the new image.
// tp_btf needs no tracefs (often absent in containers).
SEC("tp_btf/sched_process_exec")
int trace_sched_process_exec(unsigned long long *ctx)
{
	struct linux_binprm *bprm = (struct linux_binprm *)ctx[2];
	if (!bprm)
		return 0;
	struct file *file = NULL;
	if (bpf_probe_read_kernel(&file, sizeof(file), &bprm->file) || !file)
		return 0;
	struct dentry *dentry = NULL;
	if (bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry) || !dentry)
		return 0;
	const char *filename = NULL;
	bpf_probe_read_kernel(&filename, sizeof(filename), &bprm->filename);

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_kern(EVENT_EXEC, 0, buf + off, filename, dentry, 0);
}

// VFS kprobes for metadata operations. They replace tracepoint/syscalls/sys_enter_* because
// tracepoints need debugfs/tracefs, often absent in containers.
//   vfs_mkdir   - mkdir, mkdirat
//   vfs_rmdir   - rmdir, unlinkat(AT_REMOVEDIR)
//   vfs_unlink  - unlink, unlinkat
//   vfs_rename  - rename, renameat2
//   vfs_symlink - symlink, symlinkat
//   vfs_link    - link, linkat
//
// Handlers needing two path buffers (rename, link) use fill_path() with the tmp_buf percpu array to
// stay under the 512-byte BPF stack limit.

SEC("kprobe/vfs_mkdir")
int trace_vfs_mkdir(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM3(ctx);
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_MKDIR, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_rmdir")
int trace_vfs_rmdir(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM3(ctx);
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_DELETE, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_unlink")
int trace_vfs_unlink(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM3(ctx);
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_DELETE, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_rename")
int trace_vfs_rename(struct pt_regs *ctx)
{
	struct renamedata *rd = (struct renamedata *)PT_REGS_PARM1(ctx);
	if (!rd)
		return 0;

	struct dentry *old_dentry, *new_dentry;
	if (bpf_probe_read_kernel(&old_dentry, sizeof(old_dentry), &rd->old_dentry))
		return 0;
	if (bpf_probe_read_kernel(&new_dentry, sizeof(new_dentry), &rd->new_dentry))
		return 0;
	if (!old_dentry || !new_dentry)
		return 0;

	struct event *e = reserve_event(EVENT_RENAME, 0, 0);
	if (!e)
		return 0;

	fill_path(old_dentry, e->path);
	fill_path(new_dentry, e->dest);
	fill_identity(e, old_dentry);

	bpf_ringbuf_submit(e, 0);
	return 0;
}

SEC("kprobe/vfs_symlink")
int trace_vfs_symlink(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM3(ctx);
	if (!dentry)
		return 0;

	const char *symname = (const char *)PT_REGS_PARM4(ctx);

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_SYMLINK, 0, buf + off, symname, dentry);
}

SEC("kprobe/vfs_link")
int trace_vfs_link(struct pt_regs *ctx)
{
	struct dentry *old_dentry = (struct dentry *)PT_REGS_PARM1(ctx);
	struct dentry *new_dentry = (struct dentry *)PT_REGS_PARM4(ctx);
	if (!old_dentry || !new_dentry)
		return 0;

	struct event *e = reserve_event(EVENT_HARDLINK, 0, 0);
	if (!e)
		return 0;

	fill_path(old_dentry, e->path);
	fill_path(new_dentry, e->dest);
	fill_identity(e, old_dentry);

	bpf_ringbuf_submit(e, 0);
	return 0;
}

// Metadata-operation coverage: path-based ops that never create a struct file, so never hit
// vfs_open. Enum values and order MUST match the Go EventType enum
// (internal/infrastructure/event.go).
//   notify_change            - chmod/chown/utimes/truncate size (vfs_truncate -> do_truncate -> notify_change); also O_TRUNC opens (the size does change)
//   vfs_setxattr/removexattr - extended attributes
//   vfs_mknod                - device nodes / FIFOs / sockets
//   vfs_getattr              - stat/lstat/statx/fstat (same chokepoint as the guard's inode_getattr hook)
//   vfs_readlink             - symlink body reads
//   do_faccessat             - access(2)/faccessat(2); the user-supplied pathname is emitted directly (arg 2)

SEC("kprobe/notify_change")
int trace_notify_change(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM2(ctx);
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_ATTR, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_setxattr")
int trace_vfs_setxattr(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM2(ctx);
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_ATTR, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_removexattr")
int trace_vfs_removexattr(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM2(ctx);
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_ATTR, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_mknod")
int trace_vfs_mknod(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM3(ctx);
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_MKNOD, 0, buf + off, NULL, dentry);
}

SEC("kprobe/vfs_readlink")
int trace_vfs_readlink(struct pt_regs *ctx)
{
	struct dentry *dentry = (struct dentry *)PT_REGS_PARM1(ctx);
	if (!dentry)
		return 0;

	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off < 0)
		return 0;

	return emit_event_kern(EVENT_STAT, 0, buf + off, NULL, dentry);
}

SEC("kprobe/do_faccessat")
int trace_do_faccessat(struct pt_regs *ctx)
{
	const char *pathname = (const char *)PT_REGS_PARM2(ctx);
	if (!pathname)
		return 0;

	return emit_event(EVENT_STAT, 0, pathname, NULL);
}

// Memory-mapped I/O: reads happen with no read-family syscall. Keeps the mmap tracepoint because it
// also captures the fd and works alongside security_mmap_file.

SEC("tracepoint/syscalls/sys_enter_mmap")
int trace_mmap(struct trace_event_raw_sys_enter *ctx)
{
	__s32 fd = (__s32)ctx->args[4];
	if (fd >= 0)
		return emit_user(EVENT_MMAP, fd, NULL, NULL, (__u32)ctx->args[2]);
	return 0;
}

char LICENSE[] SEC("license") = "GPL";
