#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>

#define MAX_COMM_LEN 16
#define AF_INET      2
#define AF_INET6     10
#define IPPROTO_TCP  6
#define IPPROTO_UDP  17
#define IPPROTO_ICMP 1
#define S_IFMT       0170000
#define S_IFSOCK     0140000
#define DNS_PORT     53
#define DNS_PORT_TLS 853

enum event_type {
	NET_CONNECT,
	NET_ACCEPT,
	NET_SEND,
	NET_RECV,
	NET_CLOSE,
	NET_DNS,
	NET_BIND,
	NET_LISTEN,
};

struct net_event {
	__u32 pid;
	__u32 uid;
	__u32 gid;
	__u32 type;
	__u32 proto;
	__u32 size;
	__u32 fd;
	__u32 af;
	__u32 saddr[4];
	__u32 daddr[4];
	__u16 sport;
	__u16 dport;
	char comm[MAX_COMM_LEN];
	__u32 tid;
	__u64 netns;
	__u64 cgroup_id;
};

struct inode_key {
	__u64 dev;
	__u64 ino;
};

struct comm_key {
	char comm[MAX_COMM_LEN];
};

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 64);
	__type(key, struct inode_key);
	__type(value, __u8);
} watch_exe_inodes SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 64);
	__type(key, struct comm_key);
	__type(value, __u8);
} watch_comms SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} rb SEC(".maps");

char LICENSE[] SEC("license") = "GPL";

static __always_inline __u64 get_netns(void)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (!task) return 0;

	struct nsproxy *np;
	if (bpf_probe_read_kernel(&np, sizeof(np), &task->nsproxy) || !np)
		return 0;

	struct net *net_ns;
	if (bpf_probe_read_kernel(&net_ns, sizeof(net_ns), &np->net_ns) || !net_ns)
		return 0;

	unsigned int inum;
	if (bpf_probe_read_kernel(&inum, sizeof(inum), &net_ns->ns.inum))
		return 0;

	return inum;
}

static __always_inline int is_watched_binary(void)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (!task)
		return 0;

	char comm[16];
	bpf_get_current_comm(comm, sizeof(comm));

	struct comm_key ck = {};
	__builtin_memcpy(ck.comm, comm, sizeof(ck.comm));
	__u8 *found = bpf_map_lookup_elem(&watch_comms, &ck);
	if (!found)
		return 0;

	struct mm_struct *mm;
	bpf_probe_read_kernel(&mm, sizeof(mm), &task->mm);
	if (!mm)
		return 1;

	struct file *exe_file;
	bpf_probe_read_kernel(&exe_file, sizeof(exe_file), &mm->exe_file);
	if (!exe_file)
		return 1;

	struct inode *exe_inode;
	bpf_probe_read_kernel(&exe_inode, sizeof(exe_inode), &exe_file->f_inode);
	if (!exe_inode)
		return 1;

	struct inode_key ik = {};
	bpf_probe_read_kernel(&ik.ino, sizeof(ik.ino), &exe_inode->i_ino);

	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &exe_inode->i_sb);
	if (!sb)
		return 1;

	dev_t dev;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	ik.dev = dev;

	found = bpf_map_lookup_elem(&watch_exe_inodes, &ik);
	if (!found)
		return 0;

	return 1;
}

static __always_inline int get_socket_proto(int fd)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (!task)
		return 0;

	struct files_struct *files;
	bpf_probe_read_kernel(&files, sizeof(files), &task->files);
	if (!files)
		return 0;

	struct fdtable *fdt;
	bpf_probe_read_kernel(&fdt, sizeof(fdt), &files->fdt);
	if (!fdt)
		return 0;

	int max_fds;
	bpf_probe_read_kernel(&max_fds, sizeof(max_fds), &fdt->max_fds);
	if (fd < 0 || fd >= max_fds)
		return 0;

	struct file **fd_array;
	bpf_probe_read_kernel(&fd_array, sizeof(fd_array), &fdt->fd);
	if (!fd_array)
		return 0;

	struct file *file;
	bpf_probe_read_kernel(&file, sizeof(file), &fd_array[fd]);
	if (!file)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	if (!inode)
		return 0;

	umode_t mode;
	bpf_probe_read_kernel(&mode, sizeof(mode), &inode->i_mode);
	if ((mode & S_IFMT) != S_IFSOCK)
		return 0;

	struct socket *sock;
	bpf_probe_read_kernel(&sock, sizeof(sock), &file->private_data);
	if (!sock)
		return 0;

	struct sock *sk;
	bpf_probe_read_kernel(&sk, sizeof(sk), &sock->sk);
	if (!sk)
		return 0;

	unsigned short sk_protocol;
	bpf_probe_read_kernel(&sk_protocol, sizeof(sk_protocol), &sk->sk_protocol);
	return sk_protocol;
}

static __always_inline void read_inet_addr(const struct sockaddr *addr, int addrlen,
					   __u32 *af, __u32 *ip, __u16 *port)
{
	*af = 0;
	*port = 0;
	ip[0] = ip[1] = ip[2] = ip[3] = 0;

	if (!addr || addrlen < 2)
		return;

	__u16 family;
	bpf_probe_read_user(&family, sizeof(family), &addr->sa_family);
	*af = family;

	if (family == AF_INET && addrlen >= sizeof(struct sockaddr_in)) {
		struct sockaddr_in sin;
		bpf_probe_read_user(&sin, sizeof(sin), addr);
		*port = sin.sin_port;
		ip[0] = sin.sin_addr.s_addr;
	} else if (family == AF_INET6 && addrlen >= sizeof(struct sockaddr_in6)) {
		struct sockaddr_in6 sin6;
		bpf_probe_read_user(&sin6, sizeof(sin6), addr);
		*port = sin6.sin6_port;
		ip[0] = sin6.sin6_addr.in6_u.u6_addr32[0];
		ip[1] = sin6.sin6_addr.in6_u.u6_addr32[1];
		ip[2] = sin6.sin6_addr.in6_u.u6_addr32[2];
		ip[3] = sin6.sin6_addr.in6_u.u6_addr32[3];
	}
}

static __always_inline long emit_event(struct net_event *e)
{
	struct net_event *out;
	out = bpf_ringbuf_reserve(&rb, sizeof(*out), 0);
	if (!out)
		return 0;

	__builtin_memcpy(out, e, sizeof(*out));
	bpf_ringbuf_submit(out, 0);
	return 1;
}

SEC("tracepoint/syscalls/sys_enter_connect")
int trace_connect(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];
	struct sockaddr *addr = (struct sockaddr *)ctx->args[1];
	int addrlen = (int)ctx->args[2];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_CONNECT;
	e.fd = fd;
	e.size = 0;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	read_inet_addr(addr, addrlen, &e.af, e.daddr, &e.dport);

	if (e.dport == __builtin_bswap16(DNS_PORT) ||
	    e.dport == __builtin_bswap16(DNS_PORT_TLS))
		e.type = NET_DNS;

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_bind")
int trace_bind(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];
	struct sockaddr *addr = (struct sockaddr *)ctx->args[1];
	int addrlen = (int)ctx->args[2];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_BIND;
	e.fd = fd;
	e.size = 0;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	read_inet_addr(addr, addrlen, &e.af, e.saddr, &e.sport);

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_listen")
int trace_listen(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];
	int backlog = (int)ctx->args[1];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_LISTEN;
	e.fd = fd;
	e.size = backlog;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_accept")
int trace_accept(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_ACCEPT;
	e.fd = fd;
	e.size = 0;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_accept4")
int trace_accept4(struct trace_event_raw_sys_enter *ctx)
{

	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_ACCEPT;
	e.fd = fd;
	e.size = 0;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_sendto")
int trace_sendto(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];
	unsigned long len = ctx->args[2];
	struct sockaddr *addr = (struct sockaddr *)ctx->args[4];
	int addrlen = (int)ctx->args[5];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_SEND;
	e.fd = fd;
	e.size = (__u32)len;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	if (addr && addrlen > 0)
		read_inet_addr(addr, addrlen, &e.af, e.daddr, &e.dport);

	if (e.dport == __builtin_bswap16(DNS_PORT) ||
	    e.dport == __builtin_bswap16(DNS_PORT_TLS))
		e.type = NET_DNS;

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_recvfrom")
int trace_recvfrom(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];
	unsigned long len = ctx->args[2];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_RECV;
	e.fd = fd;
	e.size = (__u32)len;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_sendmsg")
int trace_sendmsg(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];
	struct msghdr *msg = (struct msghdr *)ctx->args[1];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_SEND;
	e.fd = fd;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	if (msg) {
		void *msg_name;
		int msg_namelen;
		bpf_probe_read_user(&msg_name, sizeof(msg_name), &msg->msg_name);
		bpf_probe_read_user(&msg_namelen, sizeof(msg_namelen), &msg->msg_namelen);
		if (msg_name && msg_namelen > 0)
			read_inet_addr((struct sockaddr *)msg_name, msg_namelen,
				       &e.af, e.daddr, &e.dport);
	}

	if (e.dport == __builtin_bswap16(DNS_PORT) ||
	    e.dport == __builtin_bswap16(DNS_PORT_TLS))
		e.type = NET_DNS;

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_recvmsg")
int trace_recvmsg(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_RECV;
	e.fd = fd;
	e.proto = get_socket_proto(fd);
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	emit_event(&e);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_close")
int trace_close(struct trace_event_raw_sys_enter *ctx)
{
	if (!is_watched_binary())
		return 0;

	int fd = ctx->args[0];
	int proto = get_socket_proto(fd);
	if (proto == 0)
		return 0;

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_CLOSE;
	e.fd = fd;
	e.proto = proto;
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	emit_event(&e);
	return 0;
}

SEC("kretprobe/inet_csk_accept")
int BPF_KRETPROBE(trace_inet_csk_accept)
{
	if (!is_watched_binary())
		return 0;

	struct sock *newsk = (struct sock *)PT_REGS_RC(ctx);
	if (!newsk)
		return 0;

	struct net_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = NET_ACCEPT;
	e.proto = IPPROTO_TCP;
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.netns = get_netns();
	e.cgroup_id = bpf_get_current_cgroup_id();

	unsigned short dport;
	__u32 saddr, daddr;

	bpf_probe_read_kernel(&dport, sizeof(dport), &newsk->__sk_common.skc_dport);
	bpf_probe_read_kernel(&saddr, sizeof(saddr), &newsk->__sk_common.skc_rcv_saddr);
	bpf_probe_read_kernel(&daddr, sizeof(daddr), &newsk->__sk_common.skc_daddr);

	struct inet_sock *inet = (struct inet_sock *)newsk;
	unsigned short sport;
	bpf_probe_read_kernel(&sport, sizeof(sport), &inet->inet_sport);

	e.af = AF_INET;
	e.sport = sport;
	e.dport = dport;
	e.saddr[0] = saddr;
	e.daddr[0] = daddr;

	emit_event(&e);
	return 0;
}
