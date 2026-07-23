#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

#define MAX_PATH 256

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

static __always_inline int emit_event(__u32 type, __u32 fd, const char *path, const char *dest)
{
	struct event *e;

	e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e)
		return 0;

	e->pid = bpf_get_current_pid_tgid() >> 32;
	e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e->gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e->type = type;
	e->fd = fd;

	bpf_get_current_comm(e->comm, sizeof(e->comm));

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

static __always_inline int emit_event_kern(__u32 type, __u32 fd, const char *path, const char *dest)
{
	struct event *e;

	e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e)
		return 0;

	e->pid = bpf_get_current_pid_tgid() >> 32;
	e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e->gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e->type = type;
	e->fd = fd;

	bpf_get_current_comm(e->comm, sizeof(e->comm));

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

/*
 * VFS-level kprobes – catch ALL open/read/write regardless of which
 * syscall or subsystem triggered them (normal syscalls, io_uring,
 * execve, sendfile, splice, copy_file_range, open_by_handle_at, etc.).
 *
 * We walk the dentry tree manually instead of using bpf_d_path
 * because bpf_d_path is not available in non-sleepable kprobe/fentry
 * programs and vfs_open/vfs_read/vfs_write are not sleepable.
 *
 * vfs_open       – ALL file opens.
 * vfs_read       – ALL file reads.
 * vfs_write      – ALL file writes.
 */

static __always_inline long read_path(struct dentry *dentry, char *buf, int buf_size)
{
	char tmp[64];
	int pos = buf_size - 1;
	buf[pos & (MAX_PATH - 1)] = '\0';

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

	return emit_event_kern(EVENT_OPEN, 0, buf + off, NULL);
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

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL);
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

	return emit_event_kern(EVENT_WRITE, 0, buf + off, NULL);
}

/*
 * Additional VFS kprobes for syscalls that bypass vfs_read/vfs_write.
 *
 *   vfs_readv             – readv, preadv (modern kernels)
 *   vfs_copy_file_range   – copy_file_range
 *   do_splice             – splice (when in is a file)
 *   do_splice_direct      – sendfile (via do_sendfile → do_splice_direct)
 *   security_mmap_file    – mmap (kprobe alternative to tracepoint)
 */

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

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL);
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

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL);
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

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL);
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

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL);
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

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL);
}

/*
 * do_sendfile catches the sendfile syscall on kernels where
 * sendfile no longer calls do_splice_direct/splice_file_range.
 * We resolve the file from the task's fd table using CO-RE.
 */
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

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL);
}

/*
 * vfs_iter_read is a catch-all for read-like operations on newer kernels.
 */
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

	return emit_event_kern(EVENT_READ, 0, buf + off, NULL);
}

SEC("kprobe/security_mmap_file")
int trace_security_mmap_file(struct pt_regs *ctx)
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

	return emit_event_kern(EVENT_MMAP, 0, buf + off, NULL);
}

/*
 * VFS-level kprobes for filesystem metadata operations.
 *
 * These replace tracepoint/syscalls/sys_enter_* because tracepoints
 * require debugfs/tracefs to be mounted, which is not always available
 * in containers. VFS kprobes work wherever kprobes are supported.
 *
 * vfs_mkdir     – catches both mkdir and mkdirat syscalls
 * vfs_rmdir     – catches rmdir and unlinkat(AT_REMOVEDIR)
 * vfs_unlink    – catches both unlink and unlinkat syscalls
 * vfs_rename    – catches both rename and renameat2 syscalls
 * vfs_symlink   – catches both symlink and symlinkat syscalls
 * vfs_link      – catches both link and linkat syscalls
 *
 * For handlers that need two path buffers (rename, link), we use
 * fill_path() with the tmp_buf percpu array to avoid exceeding the
 * BPF stack limit (512 bytes).
 */

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

	return emit_event_kern(EVENT_MKDIR, 0, buf + off, NULL);
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

	return emit_event_kern(EVENT_DELETE, 0, buf + off, NULL);
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

	return emit_event_kern(EVENT_DELETE, 0, buf + off, NULL);
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

	struct event *e;
	e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e)
		return 0;

	e->pid = bpf_get_current_pid_tgid() >> 32;
	e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e->gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e->type = EVENT_RENAME;
	e->fd = 0;
	bpf_get_current_comm(e->comm, sizeof(e->comm));

	fill_path(old_dentry, e->path);
	fill_path(new_dentry, e->dest);

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

	return emit_event_kern(EVENT_SYMLINK, 0, buf + off, symname);
}

SEC("kprobe/vfs_link")
int trace_vfs_link(struct pt_regs *ctx)
{
	struct dentry *old_dentry = (struct dentry *)PT_REGS_PARM1(ctx);
	struct dentry *new_dentry = (struct dentry *)PT_REGS_PARM4(ctx);
	if (!old_dentry || !new_dentry)
		return 0;

	struct event *e;
	e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e)
		return 0;

	e->pid = bpf_get_current_pid_tgid() >> 32;
	e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e->gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e->type = EVENT_HARDLINK;
	e->fd = 0;
	bpf_get_current_comm(e->comm, sizeof(e->comm));

	fill_path(old_dentry, e->path);
	fill_path(new_dentry, e->dest);

	bpf_ringbuf_submit(e, 0);
	return 0;
}

/*
 * Memory-mapped I/O – maps a file into memory so reads happen without
 * any read-family syscall. Keeps the mmap tracepoint because it also
 * captures the file descriptor and works alongside security_mmap_file.
 */

SEC("tracepoint/syscalls/sys_enter_mmap")
int trace_mmap(struct trace_event_raw_sys_enter *ctx)
{
	__s32 fd = (__s32)ctx->args[4];
	if (fd >= 0)
		return emit_event(EVENT_MMAP, fd, NULL, NULL);
	return 0;
}

char LICENSE[] SEC("license") = "GPL";
