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

/*
 * Internal kernel function kprobes – args are passed in registers correctly
 * on all kernel versions.
 */

SEC("kprobe/do_sys_open")
int trace_do_sys_open(struct pt_regs *ctx)
{
	const char *filename = (const char *)PT_REGS_PARM2(ctx);
	return emit_event(EVENT_OPEN, 0, filename, NULL);
}

SEC("kprobe/do_sys_openat2")
int trace_do_sys_openat2(struct pt_regs *ctx)
{
	const char *filename = (const char *)PT_REGS_PARM2(ctx);
	return emit_event(EVENT_OPEN, 0, filename, NULL);
}

SEC("kprobe/ksys_read")
int trace_ksys_read(struct pt_regs *ctx)
{
	__u32 fd = (__u32)PT_REGS_PARM1(ctx);
	return emit_event(EVENT_READ, fd, NULL, NULL);
}

SEC("kprobe/ksys_write")
int trace_ksys_write(struct pt_regs *ctx)
{
	__u32 fd = (__u32)PT_REGS_PARM1(ctx);
	return emit_event(EVENT_WRITE, fd, NULL, NULL);
}

/*
 * Tracepoints for syscall entry – work on all kernels including hardened
 * 6.x where __x64_sys_* kprobes have zeroed registers.
 * The trace_event_raw_sys_enter backing struct is provided by vmlinux.h (BTF).
 */

SEC("tracepoint/syscalls/sys_enter_open")
int trace_open(struct trace_event_raw_sys_enter *ctx)
{
	const char *filename = (const char *)ctx->args[0];
	return emit_event(EVENT_OPEN, 0, filename, NULL);
}

SEC("tracepoint/syscalls/sys_enter_openat")
int trace_openat(struct trace_event_raw_sys_enter *ctx)
{
	const char *filename = (const char *)ctx->args[1];
	return emit_event(EVENT_OPEN, 0, filename, NULL);
}

SEC("tracepoint/syscalls/sys_enter_unlinkat")
int trace_unlinkat(struct trace_event_raw_sys_enter *ctx)
{
	const char *pathname = (const char *)ctx->args[1];
	return emit_event(EVENT_DELETE, 0, pathname, NULL);
}

SEC("tracepoint/syscalls/sys_enter_renameat2")
int trace_renameat2(struct trace_event_raw_sys_enter *ctx)
{
	const char *oldname = (const char *)ctx->args[1];
	const char *newname = (const char *)ctx->args[3];
	return emit_event(EVENT_RENAME, 0, oldname, newname);
}

SEC("tracepoint/syscalls/sys_enter_symlinkat")
int trace_symlinkat(struct trace_event_raw_sys_enter *ctx)
{
	const char *target = (const char *)ctx->args[0];
	const char *linkpath = (const char *)ctx->args[2];
	return emit_event(EVENT_SYMLINK, 0, target, linkpath);
}

SEC("tracepoint/syscalls/sys_enter_linkat")
int trace_linkat(struct trace_event_raw_sys_enter *ctx)
{
	const char *oldpath = (const char *)ctx->args[1];
	const char *newpath = (const char *)ctx->args[3];
	return emit_event(EVENT_HARDLINK, 0, oldpath, newpath);
}

SEC("tracepoint/syscalls/sys_enter_mkdirat")
int trace_mkdirat(struct trace_event_raw_sys_enter *ctx)
{
	const char *pathname = (const char *)ctx->args[1];
	return emit_event(EVENT_MKDIR, 0, pathname, NULL);
}

char LICENSE[] SEC("license") = "GPL";
