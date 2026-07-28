#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>
#include <errno.h>

#define AF_INET      2
#define AF_INET6     10
#define IPPROTO_TCP  6
#define IPPROTO_UDP  17
#define GUARD_BLOCK  1
#define GUARD_ALLOW  2

enum net_event_type {
	NET_CONNECT,
	NET_ACCEPT,
	NET_SEND,
	NET_RECV,
	NET_CLOSE,
	NET_DNS,
	NET_BIND,
	NET_LISTEN,
};

struct net_guard_event {
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
	char comm[16];
	__u32 tid;
	__u64 netns;
	__u64 cgroup_id;
	__u32 blocked;
};

struct inode_key {
	__u64 dev;
	__u64 ino;
};

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 8);
	__type(key, __u32);
	__type(value, __u64);
} guard_net_events SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 64);
	__type(key, struct inode_key);
	__type(value, __u8);
} guard_net_exe_actions SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} guard_net_rb SEC(".maps");

char LICENSE[] SEC("license") = "GPL";

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

enum watch_action {
	WATCH_NONE,
	WATCH_ALLOW,
	WATCH_BLOCK,
};

static __always_inline int check_watched(void)
{
	struct inode_key exe_ik = {};
	if (!get_current_exe_inode(&exe_ik))
		return WATCH_NONE;

	__u8 *action = bpf_map_lookup_elem(&guard_net_exe_actions, &exe_ik);
	if (!action)
		return WATCH_NONE;

	return (*action == GUARD_ALLOW) ? WATCH_ALLOW : WATCH_BLOCK;
}

static __always_inline int is_event_type_allowed(__u32 type)
{
	__u32 key = type;
	__u64 *val = bpf_map_lookup_elem(&guard_net_events, &key);
	if (!val)
		return 0;
	return *val != 0;
}

static __always_inline long emit_event(struct net_guard_event *e)
{
	struct net_guard_event *out;
	out = bpf_ringbuf_reserve(&guard_net_rb, sizeof(*out), 0);
	if (!out)
		return 0;

	__builtin_memcpy(out, e, sizeof(*out));
	bpf_ringbuf_submit(out, 0);
	return 1;
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

SEC("lsm/socket_connect")
int guard_net_socket_connect(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	__u32 type = NET_CONNECT;
	if (!is_event_type_allowed(type))
		return 0;

	int action = check_watched();
	if (action == WATCH_NONE)
		return 0;

	struct sock *sk;
	bpf_probe_read_kernel(&sk, sizeof(sk), &sock->sk);

	struct net_guard_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = type;
	e.blocked = (action == WATCH_BLOCK) ? 1 : 0;
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.cgroup_id = bpf_get_current_cgroup_id();

	if (sk) {
		unsigned short protocol;
		bpf_probe_read_kernel(&protocol, sizeof(protocol), &sk->sk_protocol);
		e.proto = protocol;
	}

	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (task) {
		struct nsproxy *np;
		bpf_probe_read_kernel(&np, sizeof(np), &task->nsproxy);
		if (np) {
			struct net *net_ns;
			bpf_probe_read_kernel(&net_ns, sizeof(net_ns), &np->net_ns);
			if (net_ns) {
				unsigned int inum;
				bpf_probe_read_kernel(&inum, sizeof(inum), &net_ns->ns.inum);
				e.netns = inum;
			}
		}
	}

	struct sockaddr *addr = (struct sockaddr *)ctx[1];
	int addrlen = (int)(long)ctx[2];
	read_inet_addr(addr, addrlen, &e.af, e.daddr, &e.dport);

	emit_event(&e);
	return (action == WATCH_BLOCK) ? -EPERM : 0;
}

SEC("lsm/socket_bind")
int guard_net_socket_bind(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	__u32 type = NET_BIND;
	if (!is_event_type_allowed(type))
		return 0;

	int action = check_watched();
	if (action == WATCH_NONE)
		return 0;

	struct sock *sk;
	bpf_probe_read_kernel(&sk, sizeof(sk), &sock->sk);

	struct net_guard_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = type;
	e.blocked = (action == WATCH_BLOCK) ? 1 : 0;
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.cgroup_id = bpf_get_current_cgroup_id();

	if (sk) {
		unsigned short protocol;
		bpf_probe_read_kernel(&protocol, sizeof(protocol), &sk->sk_protocol);
		e.proto = protocol;
	}

	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (task) {
		struct nsproxy *np;
		bpf_probe_read_kernel(&np, sizeof(np), &task->nsproxy);
		if (np) {
			struct net *net_ns;
			bpf_probe_read_kernel(&net_ns, sizeof(net_ns), &np->net_ns);
			if (net_ns) {
				unsigned int inum;
				bpf_probe_read_kernel(&inum, sizeof(inum), &net_ns->ns.inum);
				e.netns = inum;
			}
		}
	}

	struct sockaddr *addr = (struct sockaddr *)ctx[1];
	int addrlen = (int)(long)ctx[2];
	read_inet_addr(addr, addrlen, &e.af, e.saddr, &e.sport);

	emit_event(&e);
	return (action == WATCH_BLOCK) ? -EPERM : 0;
}

SEC("lsm/socket_listen")
int guard_net_socket_listen(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	__u32 type = NET_LISTEN;
	if (!is_event_type_allowed(type))
		return 0;

	int action = check_watched();
	if (action == WATCH_NONE)
		return 0;

	struct sock *sk;
	bpf_probe_read_kernel(&sk, sizeof(sk), &sock->sk);

	struct net_guard_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = type;
	e.size = (__u32)(long)ctx[1];
	e.blocked = (action == WATCH_BLOCK) ? 1 : 0;
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.cgroup_id = bpf_get_current_cgroup_id();

	if (sk) {
		unsigned short protocol;
		bpf_probe_read_kernel(&protocol, sizeof(protocol), &sk->sk_protocol);
		e.proto = protocol;
	}

	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (task) {
		struct nsproxy *np;
		bpf_probe_read_kernel(&np, sizeof(np), &task->nsproxy);
		if (np) {
			struct net *net_ns;
			bpf_probe_read_kernel(&net_ns, sizeof(net_ns), &np->net_ns);
			if (net_ns) {
				unsigned int inum;
				bpf_probe_read_kernel(&inum, sizeof(inum), &net_ns->ns.inum);
				e.netns = inum;
			}
		}
	}

	emit_event(&e);
	return (action == WATCH_BLOCK) ? -EPERM : 0;
}

SEC("lsm/socket_sendmsg")
int guard_net_socket_sendmsg(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	__u32 type;
	int size = (int)(long)ctx[2];

	struct sock *sk;
	bpf_probe_read_kernel(&sk, sizeof(sk), &sock->sk);

	if (sk) {
		unsigned short dport;
		bpf_probe_read_kernel(&dport, sizeof(dport), &sk->__sk_common.skc_dport);
		if (dport == __builtin_bswap16(53) || dport == __builtin_bswap16(853))
			type = NET_DNS;
		else
			type = NET_SEND;
	} else {
		type = NET_SEND;
	}

	if (!is_event_type_allowed(type))
		return 0;

	int action = check_watched();
	if (action == WATCH_NONE)
		return 0;

	struct net_guard_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = type;
	e.size = size;
	e.blocked = (action == WATCH_BLOCK) ? 1 : 0;
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.cgroup_id = bpf_get_current_cgroup_id();

	if (sk) {
		unsigned short protocol;
		bpf_probe_read_kernel(&protocol, sizeof(protocol), &sk->sk_protocol);
		e.proto = protocol;
	}

	struct msghdr *msg = (struct msghdr *)ctx[1];
	if (msg) {
		void *msg_name;
		int msg_namelen;
		bpf_probe_read_user(&msg_name, sizeof(msg_name), &msg->msg_name);
		bpf_probe_read_user(&msg_namelen, sizeof(msg_namelen), &msg->msg_namelen);
		if (msg_name && msg_namelen > 0)
			read_inet_addr((struct sockaddr *)msg_name, msg_namelen,
				       &e.af, e.daddr, &e.dport);
	}

	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (task) {
		struct nsproxy *np;
		bpf_probe_read_kernel(&np, sizeof(np), &task->nsproxy);
		if (np) {
			struct net *net_ns;
			bpf_probe_read_kernel(&net_ns, sizeof(net_ns), &np->net_ns);
			if (net_ns) {
				unsigned int inum;
				bpf_probe_read_kernel(&inum, sizeof(inum), &net_ns->ns.inum);
				e.netns = inum;
			}
		}
	}

	emit_event(&e);
	return (action == WATCH_BLOCK) ? -EPERM : 0;
}

SEC("lsm/socket_recvmsg")
int guard_net_socket_recvmsg(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	__u32 type = NET_RECV;
	if (!is_event_type_allowed(type))
		return 0;

	int action = check_watched();
	if (action == WATCH_NONE)
		return 0;

	struct sock *sk;
	bpf_probe_read_kernel(&sk, sizeof(sk), &sock->sk);

	struct net_guard_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = type;
	e.blocked = (action == WATCH_BLOCK) ? 1 : 0;
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.cgroup_id = bpf_get_current_cgroup_id();

	if (sk) {
		unsigned short protocol;
		bpf_probe_read_kernel(&protocol, sizeof(protocol), &sk->sk_protocol);
		e.proto = protocol;
	}

	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (task) {
		struct nsproxy *np;
		bpf_probe_read_kernel(&np, sizeof(np), &task->nsproxy);
		if (np) {
			struct net *net_ns;
			bpf_probe_read_kernel(&net_ns, sizeof(net_ns), &np->net_ns);
			if (net_ns) {
				unsigned int inum;
				bpf_probe_read_kernel(&inum, sizeof(inum), &net_ns->ns.inum);
				e.netns = inum;
			}
		}
	}

	emit_event(&e);
	return (action == WATCH_BLOCK) ? -EPERM : 0;
}
