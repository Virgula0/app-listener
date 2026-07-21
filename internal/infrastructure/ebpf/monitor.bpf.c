#include <linux/types.h>
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

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

SEC("kprobe/__x64_sys_openat")
int trace_openat(struct pt_regs *ctx)
{
	const char *filename = (const char *)PT_REGS_PARM2(ctx);
	return emit_event(EVENT_OPEN, 0, filename, NULL);
}

SEC("kprobe/__x64_sys_open")
int trace_open(struct pt_regs *ctx)
{
	const char *filename = (const char *)PT_REGS_PARM1(ctx);
	return emit_event(EVENT_OPEN, 0, filename, NULL);
}

SEC("kprobe/__x64_sys_read")
int trace_read(struct pt_regs *ctx)
{
	__u32 fd = (__u32)PT_REGS_PARM1(ctx);
	return emit_event(EVENT_READ, fd, NULL, NULL);
}

SEC("kprobe/__x64_sys_write")
int trace_write(struct pt_regs *ctx)
{
	__u32 fd = (__u32)PT_REGS_PARM1(ctx);
	return emit_event(EVENT_WRITE, fd, NULL, NULL);
}

SEC("kprobe/__x64_sys_unlinkat")
int trace_unlinkat(struct pt_regs *ctx)
{
	const char *pathname = (const char *)PT_REGS_PARM2(ctx);
	return emit_event(EVENT_DELETE, 0, pathname, NULL);
}

SEC("kprobe/__x64_sys_renameat2")
int trace_renameat2(struct pt_regs *ctx)
{
	const char *oldname = (const char *)PT_REGS_PARM2(ctx);
	const char *newname = (const char *)PT_REGS_PARM4(ctx);
	return emit_event(EVENT_RENAME, 0, oldname, newname);
}

SEC("kprobe/__x64_sys_symlinkat")
int trace_symlinkat(struct pt_regs *ctx)
{
	const char *target = (const char *)PT_REGS_PARM1(ctx);
	const char *linkpath = (const char *)PT_REGS_PARM3(ctx);
	return emit_event(EVENT_SYMLINK, 0, target, linkpath);
}

SEC("kprobe/__x64_sys_linkat")
int trace_linkat(struct pt_regs *ctx)
{
	const char *oldpath = (const char *)PT_REGS_PARM2(ctx);
	const char *newpath = (const char *)PT_REGS_PARM4(ctx);
	return emit_event(EVENT_HARDLINK, 0, oldpath, newpath);
}

SEC("kprobe/__x64_sys_mkdirat")
int trace_mkdirat(struct pt_regs *ctx)
{
	const char *pathname = (const char *)PT_REGS_PARM2(ctx);
	return emit_event(EVENT_MKDIR, 0, pathname, NULL);
}

char LICENSE[] SEC("license") = "GPL";
