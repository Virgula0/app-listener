#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>
#include <errno.h>
#include "inode_dev.h"

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
	// Process gates: always reported, never in the -e filter (guard_net_events).
	NET_PTRACE,
	NET_TRACED_EXEC,
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
	__u32 reason; // NETG_REASON_*: why an ALLOW row did not admit the process
};

struct inode_key {
	__u64 dev;
	__u64 ino;
};

// Per-applet identity for uutils multicall binaries: a multicall exe's key carries its attested
// applet in dev bits 32-47, so whitelisting one applet does not admit the other ~108. See mc_tag.h.
#include "mc_tag.h"
// Superseded keys and the code-suspect mark (trust_code_suspect), with this object's own maps.
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wunused-function" // helpers only the guard objects call
#include "exe_supersede.h"
#pragma clang diagnostic pop

// net_guard_event.reason. A suspect process's reason is its trust_code_suspect value.
#define NETG_REASON_CODE SUSPECT_PRELOAD // mapped code its exe doesn't vouch for, or traced at start
#define NETG_REASON_ENV SUSPECT_LAUNCH   // started with LD_PRELOAD / LD_AUDIT
#define NETG_REASON_SUPERSEDED 3         // a replaced binary run after its replacement

// mc_multicall value bits (presence is all mc_tag.h reads). MC_HAS_ALLOW: one of its applets is
// whitelisted, so its processes are protected and judged before their applet is attested.
#define MC_PRESENT 1
#define MC_HAS_ALLOW 2

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 8);
	__type(key, __u32);
	__type(value, __u64);
} guard_net_events SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 4);
	__type(key, __u32);
	__type(value, __u64);
} guard_net_config SEC(".maps");
// config[0] = default_action (0=allow, 1=block)
// config[1] = blocking_enabled (0=events only, 1=real blocking)
// config[2] = unsafe_families (0=AF_INET/AF_INET6 only, 1=all families)
// config[3] = throttle_enabled (0=emit every event, 1=rate-limit per (type, verdict, comm))

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 64);
	__type(key, struct inode_key);
	__type(value, __u8);
} guard_net_exe_actions SEC(".maps");

struct throttle_key {
	__u32 type;
	__u32 verdict; // blocked | reason << 1: an allowed event never hides a refusal of the same comm
	char comm[16];
};

// Per-(type, verdict, comm) throttle: the LSM hooks are global, so in whitelist mode every blocked
// socket op of noisy host processes would flood the ring buffer (dropping important events on
// overflow). One event per key per interval.
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 512);
	__type(key, struct throttle_key);
	__type(value, __u64);
} guard_net_throttle SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} guard_net_rb SEC(".maps");

char LICENSE[] SEC("license") = "GPL";

// exe_inode_of returns task's main exe inode (mm->exe_file); NULL for kernel threads.
static __always_inline struct inode *exe_inode_of(struct task_struct *task)
{
	if (!task)
		return NULL;
	struct mm_struct *mm;
	bpf_probe_read_kernel(&mm, sizeof(mm), &task->mm);
	if (!mm)
		return NULL;
	struct file *exe_file;
	bpf_probe_read_kernel(&exe_file, sizeof(exe_file), &mm->exe_file);
	if (!exe_file)
		return NULL;
	struct inode *exe_inode = NULL;
	bpf_probe_read_kernel(&exe_inode, sizeof(exe_inode), &exe_file->f_inode);
	return exe_inode;
}

static __always_inline int fill_inode_key(struct inode *inode, struct inode_key *k)
{
	if (!inode)
		return 0;
	bpf_probe_read_kernel(&k->ino, sizeof(k->ino), &inode->i_ino);
	k->dev = inode_dev((__u64)inode);
	return 1;
}

static __always_inline int get_task_exe_key(struct task_struct *task, struct inode_key *ik)
{
	if (!fill_inode_key(exe_inode_of(task), ik))
		return 0;
	// For a multicall exe, fold in the applet this process attested at exec: a whitelist row for one
	// applet must not admit the others. Non-multicall exes keep their plain key (mc_tag_bits == 0).
	ik->dev |= mc_tag_bits(task, ik);
	return 1;
}

enum watch_action {
	WATCH_NONE,
	WATCH_ALLOW,
	WATCH_BLOCK,
};

// netg_suspect_lost[0]: code-suspect marks a fork could not carry to the child (trust_code_suspect
// full). The unmarked child would run its parent's code admitted, and this hook can't refuse the
// fork nor signal the child: from then on no ALLOW row admits anyone (userspace reports it).
struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, __u64);
} netg_suspect_lost SEC(".maps");

static __always_inline void suspect_lost(void)
{
	__u32 z = 0;
	__u64 *lost = bpf_map_lookup_elem(&netg_suspect_lost, &z);
	if (lost)
		__sync_fetch_and_add(lost, 1);
}

// admit_refusal: why task, whose exe has an ALLOW row, is not admitted: it runs code its exe doesn't
// vouch for (trust_code_suspect), or it exec'd a replaced binary after the replacement
// (exe_supersede.h). 0 = admitted.
static __noinline __u32 admit_refusal(struct task_struct *task, struct inode_key *exe)
{
	__u32 tgid = 0;
	if (!task)
		return NETG_REASON_CODE;
	bpf_probe_read_kernel(&tgid, sizeof(tgid), &task->tgid);
	__u8 *why = bpf_map_lookup_elem(&trust_code_suspect, &tgid);
	if (why)
		return *why ? *why : NETG_REASON_CODE;
	__u32 z = 0;
	__u64 *lost = bpf_map_lookup_elem(&netg_suspect_lost, &z);
	if (!lost || *lost)
		return NETG_REASON_CODE;
	return exe_refused(task, exe) ? NETG_REASON_SUPERSEDED : 0;
}

// task_admitted: task's exe has an ALLOW row that admits it.
static __always_inline int task_admitted(struct task_struct *task)
{
	struct inode_key ik = {};
	if (!get_task_exe_key(task, &ik))
		return 0;
	__u8 *action = bpf_map_lookup_elem(&guard_net_exe_actions, &ik);
	return action && *action == GUARD_ALLOW && !admit_refusal(task, &ik);
}

// task_protected: task's exe has an ALLOW row, admitting it or not, or is a multicall one of whose
// applets has (its own applet is attested only at sched_process_exec, after the image is mapped).
// Its memory and mapped code are what admits it, so they are guarded.
static __always_inline int task_protected(struct task_struct *task)
{
	struct inode_key ik = {};
	if (!fill_inode_key(exe_inode_of(task), &ik))
		return 0;
	__u32 *mc = bpf_map_lookup_elem(&mc_multicall, &ik);
	if (mc)
		return (*mc & MC_HAS_ALLOW) != 0;
	__u8 *action = bpf_map_lookup_elem(&guard_net_exe_actions, &ik);
	return action && *action == GUARD_ALLOW;
}

static __always_inline int check_watched(__u32 *reason)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	struct inode_key exe_ik = {};
	if (!get_task_exe_key(task, &exe_ik))
		return WATCH_NONE;

	__u8 *action = bpf_map_lookup_elem(&guard_net_exe_actions, &exe_ik);
	if (action && *action != GUARD_ALLOW)
		return WATCH_BLOCK;
	if (action) {
		*reason = admit_refusal(task, &exe_ik);
		if (!*reason)
			return WATCH_ALLOW;
	}

	// Fall back to default action (a refused ALLOW row too: never admitted, in either mode).
	// config[0]: 0 = allow (blacklist mode), 1 = block (whitelist mode).
	__u32 key = 0;
	__u64 *default_action = bpf_map_lookup_elem(&guard_net_config, &key);
	if (default_action && *default_action == 1)
		return WATCH_BLOCK;
	return WATCH_NONE;
}

static __always_inline int is_socket_guarded(struct socket *sock)
{
	if (!sock)
		return 0;

	// Check unsafe families flag (config[2]).
	__u32 config_key = 2;
	__u64 *unsafe = bpf_map_lookup_elem(&guard_net_config, &config_key);
	if (unsafe && *unsafe == 1)
		return 1; // unsafe: guard all socket families

	// Safe mode: only guard AF_INET and AF_INET6.
	struct sock *sk;
	bpf_probe_read_kernel(&sk, sizeof(sk), &sock->sk);
	if (!sk)
		return 0;

	short family;
	bpf_probe_read_kernel(&family, sizeof(family),
			      &sk->__sk_common.skc_family);
	return (family == AF_INET || family == AF_INET6) ? 1 : 0;
}

static __always_inline int should_block(void)
{
	__u32 key = 1;
	__u64 *enabled = bpf_map_lookup_elem(&guard_net_config, &key);
	return (enabled && *enabled == 1) ? 1 : 0;
}

static __always_inline int is_event_type_allowed(__u32 type)
{
	__u32 key = type;
	__u64 *val = bpf_map_lookup_elem(&guard_net_events, &key);
	if (!val)
		return 0;
	return *val != 0;
}

#define EVENT_THROTTLE_NS 250000000ULL /* 250ms */

static __always_inline long emit_event(struct net_guard_event *e)
{
	// Rate-limit per (type, verdict, comm) so a busy host can't overflow the ring buffer and drop events of
	// interest. Disabled via config[3] for full fidelity (--no-throttle).
	__u32 ckey = 3;
	__u64 *throttle_on = bpf_map_lookup_elem(&guard_net_config, &ckey);
	if (!throttle_on || *throttle_on != 0) {
		struct throttle_key tk = {
			.type = e->type,
			.verdict = e->blocked | e->reason << 1,
		};
		__builtin_memcpy(tk.comm, e->comm, sizeof(e->comm));

		__u64 now = bpf_ktime_get_ns();
		__u64 *last = bpf_map_lookup_elem(&guard_net_throttle, &tk);
		if (last) {
			if (now - *last < EVENT_THROTTLE_NS)
				return 0;
		}
		bpf_map_update_elem(&guard_net_throttle, &tk, &now, BPF_ANY);
	}

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

	if (!is_socket_guarded(sock))
		return 0;

	__u32 type = NET_CONNECT;
	// The -e event filter is presentation only: it must gate what is REPORTED, never whether the
	// operation is enforced. Run the watch/allow decision unconditionally; filter emit_event below.
	__u32 reason = 0;
	int action = check_watched(&reason);
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
	e.reason = reason;
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

	if (is_event_type_allowed(type))
		emit_event(&e);
	if (action == WATCH_BLOCK && should_block())
		return -EPERM;
	return 0;
}

SEC("lsm/socket_bind")
int guard_net_socket_bind(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	if (!is_socket_guarded(sock))
		return 0;

	__u32 type = NET_BIND;
	// The -e event filter is presentation only: it must gate what is REPORTED, never whether the
	// operation is enforced. Run the watch/allow decision unconditionally; filter emit_event below.
	__u32 reason = 0;
	int action = check_watched(&reason);
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
	e.reason = reason;
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

	if (is_event_type_allowed(type))
		emit_event(&e);
	if (action == WATCH_BLOCK && should_block())
		return -EPERM;
	return 0;
}

SEC("lsm/socket_listen")
int guard_net_socket_listen(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	if (!is_socket_guarded(sock))
		return 0;

	__u32 type = NET_LISTEN;
	// The -e event filter is presentation only: it must gate what is REPORTED, never whether the
	// operation is enforced. Run the watch/allow decision unconditionally; filter emit_event below.
	__u32 reason = 0;
	int action = check_watched(&reason);
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
	e.reason = reason;
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

	if (is_event_type_allowed(type))
		emit_event(&e);
	if (action == WATCH_BLOCK && should_block())
		return -EPERM;
	return 0;
}

SEC("lsm/socket_sendmsg")
int guard_net_socket_sendmsg(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	if (!is_socket_guarded(sock))
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

	// The -e event filter is presentation only: it must gate what is REPORTED, never whether the
	// operation is enforced. Run the watch/allow decision unconditionally; filter emit_event below.
	__u32 reason = 0;
	int action = check_watched(&reason);
	if (action == WATCH_NONE)
		return 0;

	struct net_guard_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = type;
	e.size = size;
	e.blocked = (action == WATCH_BLOCK) ? 1 : 0;
	e.reason = reason;
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

	if (is_event_type_allowed(type))
		emit_event(&e);
	if (action == WATCH_BLOCK && should_block())
		return -EPERM;
	return 0;
}

SEC("lsm/socket_recvmsg")
int guard_net_socket_recvmsg(unsigned long long *ctx)
{
	struct socket *sock = (struct socket *)ctx[0];
	if (!sock)
		return 0;

	if (!is_socket_guarded(sock))
		return 0;

	__u32 type = NET_RECV;
	// The -e event filter is presentation only: it must gate what is REPORTED, never whether the
	// operation is enforced. Run the watch/allow decision unconditionally; filter emit_event below.
	__u32 reason = 0;
	int action = check_watched(&reason);
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
	e.reason = reason;
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

	if (is_event_type_allowed(type))
		emit_event(&e);
	if (action == WATCH_BLOCK && should_block())
		return -EPERM;
	return 0;
}

// Multicall applet attestation (mc_tag.h): the exec tracepoint stamps a multicall process with the
// applet it was invoked as (or MC_TAG_NONE), fork inherits it, task_free drops it. get_task_exe_key
// reads the stamp back into the exe key. Mandatory for whitelist mode to tell uutils applets apart;
// a multicall process that exec'd before this object attached keys as NONE (denied in whitelist mode).
SEC("tp_btf/sched_process_exec")
int netg_exec_applet(unsigned long long *ctx)
{
	mc_attest(ctx);
	return 0;
}

// Process lifecycle of the exec stamps (exe_supersede.h, mc_tag.h) and the code-suspect mark. Fork:
// the child runs the parent's image and shares its mapped code. At sched_process_fork, not
// task_alloc: the child's tgid is assigned after security_task_alloc.
SEC("tp_btf/sched_process_fork")
int netg_sched_fork(unsigned long long *ctx)
{
	struct task_struct *parent = (struct task_struct *)ctx[0];
	struct task_struct *child = (struct task_struct *)ctx[1];
	if (!parent || !child)
		return 0;
	mc_stamp_fork(parent, child);
	exe_stamp_fork(parent, child);
	__u32 ptgid = 0, ctgid = 0;
	bpf_probe_read_kernel(&ptgid, sizeof(ptgid), &parent->tgid);
	bpf_probe_read_kernel(&ctgid, sizeof(ctgid), &child->tgid);
	if (ptgid == ctgid)
		return 0;
	__u8 *why = bpf_map_lookup_elem(&trust_code_suspect, &ptgid);
	if (!why) {
		bpf_map_delete_elem(&trust_code_suspect, &ctgid); // a reused pid's leftover
		return 0;
	}
	__u8 v = *why;
	if (bpf_map_update_elem(&trust_code_suspect, &ctgid, &v, BPF_ANY))
		suspect_lost();
	return 0;
}

SEC("lsm/task_free")
int netg_task_free(unsigned long long *ctx)
{
	struct task_struct *task = (struct task_struct *)ctx[0];
	if (!task)
		return 0;
	// Keyed on tgid: drop only when the GROUP LEADER exits (task_free fires per thread).
	__u32 pid = 0, tgid = 0;
	bpf_probe_read_kernel(&pid, sizeof(pid), &task->pid);
	bpf_probe_read_kernel(&tgid, sizeof(tgid), &task->tgid);
	if (pid != tgid)
		return 0;
	mc_stamp_free(task, tgid);
	exe_stamp_free(task, tgid);
	bpf_map_delete_elem(&trust_code_suspect, &tgid);
	return 0;
}

// Exec replaces every mapping. Runs before the new image's segments and libraries are mapped, so
// those are judged afresh (netg_mmap_file), and stamps the exec for exe_refused.
SEC("lsm/bprm_committed_creds")
int netg_bprm_committed(unsigned long long *ctx)
{
	exe_stamp_exec();
	__u32 tgid = bpf_get_current_pid_tgid() >> 32;
	bpf_map_delete_elem(&trust_code_suspect, &tgid);
	return 0;
}

// Runtime-generated code (GPU driver JIT into a memfd/O_TMPFILE/mkstemp file it maps executable):
// trusted only for the very image that created the inode while no other thread group has opened it
// for writing, as the trust guard's guard_jit_origin. Recorded only for a protected creator. A
// foreign write-open also marks the creator code-suspect: its existing mapping would run the new
// bytes. LRU: an evicted entry is no provenance, which judges the mapping untrusted.
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
} netg_jit_origin SEC(".maps");

// Reused inode numbers. A whitelisted key that loses its last link is superseded (exe_supersede.h):
// only processes exec'd before keep its row. Once the inode is freed its rows go, so the number's
// next file inherits nothing; a multicall's applet rows become unreachable with its presence row.
static __always_inline int net_member(struct inode_key *k)
{
	return bpf_map_lookup_elem(&guard_net_exe_actions, k) || bpf_map_lookup_elem(&mc_multicall, k);
}

static __always_inline int net_supersede(struct inode *inode)
{
	struct inode_key k = {};
	if (!exe_last_link_key(inode, &k) || !net_member(&k))
		return 0;
	return exe_note_supersede(&k, 0) ? -EPERM : 0;
}

SEC("lsm/inode_unlink")
int netg_inode_unlink(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;
	struct inode *inode = NULL;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);
	return net_supersede(inode);
}

// A rename over a file unlinks the victim with no inode_unlink. No flags: for RENAME_EXCHANGE it
// runs once per side, marking a whitelisted side that keeps its link (stricter).
SEC("lsm/inode_rename")
int netg_inode_rename(unsigned long long *ctx)
{
	struct dentry *new_dentry = (struct dentry *)ctx[3];
	if (!new_dentry)
		return 0;
	struct inode *victim = NULL;
	bpf_probe_read_kernel(&victim, sizeof(victim), &new_dentry->d_inode);
	return net_supersede(victim);
}

// i_nlink == 0 is the load-bearing check: inode_free_security also fires for live inodes leaving the
// cache, whose rows must stay. A freed inode has no process left running it, so its rows go now,
// rows before the mark: a number reused before this hook runs is refused by the mark meanwhile.
SEC("lsm/inode_free_security")
int netg_inode_free(unsigned long long *ctx)
{
	struct inode *inode = (struct inode *)ctx_ptr(ctx, 0);
	if (!inode)
		return 0;
	unsigned int nlink = 1;
	bpf_probe_read_kernel(&nlink, sizeof(nlink), &inode->i_nlink);
	struct inode_key k = {};
	if (nlink != 0 || !fill_inode_key(inode, &k))
		return 0;
	if (net_member(&k)) {
		bpf_map_delete_elem(&guard_net_exe_actions, &k);
		bpf_map_delete_elem(&mc_multicall, &k);
	}
	bpf_map_delete_elem(&exe_superseded, &k);
	bpf_map_delete_elem(&netg_jit_origin, &k);
	return 0;
}

// ---- Code integrity (whitelist mode) ----
// An ALLOW row admits an exe's code, not whatever runs inside its process. A protected process
// (task_protected) that maps executable code its exe doesn't vouch for (LD_PRELOAD/LD_AUDIT,
// dlopen, LD_LIBRARY_PATH) or is started with LD_PRELOAD/LD_AUDIT set is marked code-suspect
// (trust_code_suspect) and loses its row until it execs. Its memory is refused to every process the
// whitelist doesn't admit (ptrace, process_vm_writev, /proc/<pid>/mem, pidfd_getfd), and so is a
// traced exec of it. A mark that can't be recorded refuses the mapping, or kills the process where
// the hook can't refuse.

// netg_vouched_devs / netg_dead_devs / netg_dev_seq: as the trust guard's guard_vouched_devs (see
// guard_trust.bpf.c): sb_dev of every non-FUSE superblock with a mount without nosuid in pid 1's
// namespace, synced by userspace (guard.MountVouch); only there do root ownership bits mean root.
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 4096);
	__type(key, __u64);
	__type(value, __u64);
} netg_vouched_devs SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 16384);
	__type(key, __u64);
	__type(value, __u64);
} netg_dead_devs SEC(".maps");

#define DEV_SEQ_NOW 0
#define DEV_SEQ_LOST 1

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 2);
	__type(key, __u32);
	__type(value, __u64);
} netg_dev_seq SEC(".maps");

#define PROT_EXEC 0x4
#define FMODE_WRITE 0x2
#define O_CREAT 00000100
#define O_EXCL 00000200
#define __O_TMPFILE 020000000
#define S_IWGRP 00020
#define S_IWOTH 00002
#define FUSE_SUPER_MAGIC 0x65735546
#define PTRACE_MODE_ATTACH 0x02
#define NO_GID 0xffffffff
#define BPF_IS_ERR_PTR(p) ((unsigned long)(p) > (unsigned long)-4096)

static __always_inline int mark_suspect(__u32 tgid, __u8 why)
{
	if (bpf_map_lookup_elem(&trust_code_suspect, &tgid))
		return 0; // already marked; the first reason stays
	return bpf_map_update_elem(&trust_code_suspect, &tgid, &why, BPF_ANY) ? -1 : 0;
}

static __always_inline __u32 inode_uid(struct inode *inode)
{
	__u32 uid = 0xffffffff;
	bpf_probe_read_kernel(&uid, sizeof(uid), &inode->i_uid); // kuid_t { uid_t val; }
	return uid;
}

// inode_ro_for: owned by uid and modifiable in place by no one else: not other-writable, and
// group-writable only for group gw_gid (NO_GID: none).
static __always_inline int inode_ro_for(struct inode *inode, __u32 uid, __u32 gw_gid)
{
	if (!inode || inode_uid(inode) != uid)
		return 0;
	umode_t mode = 0;
	bpf_probe_read_kernel(&mode, sizeof(mode), &inode->i_mode);
	if (mode & S_IWOTH)
		return 0;
	if (mode & S_IWGRP) {
		__u32 gid = NO_GID;
		bpf_probe_read_kernel(&gid, sizeof(gid), &inode->i_gid); // kgid_t { gid_t val; }
		return gid == gw_gid;
	}
	return 1;
}

static __always_inline struct super_block *inode_sb(struct inode *inode)
{
	struct super_block *sb = NULL;
	if (inode)
		bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	return sb;
}

// owner_vouched: root vouches for the ownership bits on inode's superblock: not FUSE (its server
// reports any owner), vouched (netg_vouched_devs) and not shut down since that was synced.
static __always_inline int owner_vouched(struct inode *inode)
{
	struct super_block *sb = inode_sb(inode);
	if (!sb)
		return 0;
	unsigned long magic = 0;
	bpf_probe_read_kernel(&magic, sizeof(magic), &sb->s_magic);
	if (magic == FUSE_SUPER_MAGIC)
		return 0;
	__u64 dev = sb_dev(inode);
	__u64 *synced = bpf_map_lookup_elem(&netg_vouched_devs, &dev);
	if (!dev || !synced)
		return 0;
	__u64 *died = bpf_map_lookup_elem(&netg_dead_devs, &dev);
	return !died || *died <= *synced;
}

static __always_inline struct inode *parent_inode(struct dentry *dentry)
{
	struct dentry *parent = NULL;
	struct inode *pinode = NULL;
	if (!dentry)
		return NULL;
	bpf_probe_read_kernel(&parent, sizeof(parent), &dentry->d_parent);
	if (!parent || parent == dentry)
		return NULL;
	bpf_probe_read_kernel(&pinode, sizeof(pinode), &parent->d_inode);
	return pinode;
}

// code_vouched: exe vouches for the code in inode (its dentry: the directory it sits in):
//   - a system library: root-owned, in a root-owned directory, neither modifiable by a non-root
//     user, on a superblock root vouches for (the trust guard's is_system_trusted rule);
//   - for an exe a regular user owns, a library of that user modifiable by nobody the exe isn't (the
//     same owner, or the exe's group when the exe is group-writable too), in such a directory or a
//     root one, on the exe's superblock or a vouched one. Whoever can change it can already rewrite
//     the exe in place, which keeps its identity: no new admission.
static __noinline int code_vouched(struct inode *inode, struct dentry *dentry, struct inode *exe)
{
	struct inode *dir = parent_inode(dentry);
	if (!inode || !dir || !exe)
		return 0;
	int vouched = owner_vouched(inode);
	if (vouched && inode_ro_for(inode, 0, 0) && inode_ro_for(dir, 0, 0))
		return 1;
	__u32 uid = inode_uid(exe);
	if (uid == 0 || (!vouched && inode_sb(inode) != inode_sb(exe)))
		return 0;
	umode_t mode = 0;
	__u32 gid = NO_GID;
	bpf_probe_read_kernel(&mode, sizeof(mode), &exe->i_mode);
	bpf_probe_read_kernel(&gid, sizeof(gid), &exe->i_gid);
	__u32 gw = (mode & S_IWGRP) ? gid : NO_GID;
	return inode_ro_for(inode, uid, gw) && (inode_ro_for(dir, uid, gw) || inode_ro_for(dir, 0, 0));
}

struct proc_image {
	__u64 mm;
	struct inode_key exe;
	__u32 tgid;
};

static __always_inline int fill_current_image(struct proc_image *img)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	struct mm_struct *mm = NULL;
	if (!task)
		return 0;
	bpf_probe_read_kernel(&mm, sizeof(mm), &task->mm);
	if (!mm || !fill_inode_key(exe_inode_of(task), &img->exe))
		return 0;
	img->mm = (__u64)mm;
	img->tgid = bpf_get_current_pid_tgid() >> 32;
	return 1;
}

static __always_inline void jit_record_creation(struct inode *inode)
{
	struct inode_key k = {};
	struct proc_image img = {};
	if (!fill_inode_key(inode, &k) || !fill_current_image(&img))
		return;
	if (!task_protected((struct task_struct *)bpf_get_current_task()))
		return;
	struct jit_origin o = {};
	o.mm = img.mm;
	o.exe = img.exe;
	o.tgid = img.tgid;
	bpf_map_update_elem(&netg_jit_origin, &k, &o, BPF_ANY);
}

// jit_taint_foreign_write: an open for writing by anything but the creating image taints the inode
// for good and marks the creator (see netg_jit_origin). Nonzero when that mark could not be recorded.
static __always_inline int jit_taint_foreign_write(struct inode *inode)
{
	struct inode_key k = {};
	if (!fill_inode_key(inode, &k))
		return 0;
	struct jit_origin *o = bpf_map_lookup_elem(&netg_jit_origin, &k);
	if (!o)
		return 0;
	struct proc_image img = {};
	if (fill_current_image(&img) && o->tgid == img.tgid && o->mm == img.mm &&
	    o->exe.dev == img.exe.dev && o->exe.ino == img.exe.ino)
		return 0; // the creator writing its own JIT output
	o->tainted = 1;
	return mark_suspect(o->tgid, NETG_REASON_CODE);
}

static __always_inline int jit_self_created(struct inode *inode)
{
	struct inode_key k = {};
	struct proc_image img = {};
	if (!fill_inode_key(inode, &k))
		return 0;
	struct jit_origin *o = bpf_map_lookup_elem(&netg_jit_origin, &k);
	if (!o || o->tainted || !fill_current_image(&img))
		return 0;
	return o->tgid == img.tgid && o->mm == img.mm && o->exe.dev == img.exe.dev &&
	       o->exe.ino == img.exe.ino;
}

// judge_code: file is about to be mapped executable in the current process. A protected process
// mapping code its exe doesn't vouch for is marked code-suspect; -EPERM if the mark can't be
// recorded (fail closed: the mapping is refused).
static __noinline int judge_code(struct file *file)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	struct inode *inode = NULL;
	struct dentry *dentry = NULL;
	if (!file || !task_protected(task))
		return 0;
	__u32 tgid = bpf_get_current_pid_tgid() >> 32;
	if (bpf_map_lookup_elem(&trust_code_suspect, &tgid))
		return 0;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	struct inode *exe = exe_inode_of(task);
	if (!inode || inode == exe)
		return 0; // a process's own image is its identity
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (code_vouched(inode, dentry, exe) || jit_self_created(inode))
		return 0;
	return mark_suspect(tgid, NETG_REASON_CODE) ? -EPERM : 0;
}

SEC("lsm/mmap_file")
int netg_mmap_file(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	unsigned long prot = (unsigned long)ctx[2];
	if (!file || !(prot & PROT_EXEC))
		return 0;
	return judge_code(file);
}

// mmap without PROT_EXEC, then mprotect(PROT_EXEC): the same code, judged here.
SEC("lsm/file_mprotect")
int netg_file_mprotect(unsigned long long *ctx)
{
	struct vm_area_struct *vma = (struct vm_area_struct *)ctx[0];
	unsigned long prot = (unsigned long)ctx[2];
	if (!vma || !(prot & PROT_EXEC))
		return 0;
	struct file *file = NULL;
	bpf_probe_read_kernel(&file, sizeof(file), &vma->vm_file);
	if (!file)
		return 0; // anonymous memory: no file to inject through
	return judge_code(file);
}

// JIT provenance bookkeeping. Only flags that guarantee creation count: O_CREAT|O_EXCL fails on an
// existing name, __O_TMPFILE always makes a new inode (bare O_CREAT opens a planted file).
SEC("lsm/file_open")
int netg_file_open(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	if (!file)
		return 0;
	__u32 f_mode = 0, f_flags = 0;
	struct inode *inode = NULL;
	bpf_probe_read_kernel(&f_mode, sizeof(f_mode), &file->f_mode);
	bpf_probe_read_kernel(&f_flags, sizeof(f_flags), &file->f_flags);
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	if ((f_flags & __O_TMPFILE) == __O_TMPFILE || ((f_flags & O_CREAT) && (f_flags & O_EXCL)))
		jit_record_creation(inode);
	if ((f_mode & FMODE_WRITE) && jit_taint_foreign_write(inode))
		return -EPERM;
	return 0;
}

// memfd_create() allocates its file without security_file_open(): provenance is recorded where the
// kernel makes it. Best-effort: without it a memfd is never vouched for (fail closed).
SEC("fexit/memfd_alloc_file")
int netg_memfd_alloc(unsigned long long *ctx)
{
	// fexit context: the traced function's arguments, then its return value.
	struct file *file = (struct file *)ctx[2];
	if (!file || BPF_IS_ERR_PTR(file))
		return 0;
	struct inode *inode = NULL;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	jit_record_creation(inode);
	return 0;
}

// netg_sb_delete runs before a dying superblock's dev can be reused: its vouch is tombstoned with
// the shutdown sequence (see the trust guard's trust_sb_delete).
SEC("lsm/sb_delete")
int netg_sb_delete(unsigned long long *ctx)
{
	struct super_block *sb = (struct super_block *)ctx[0];
	if (!sb)
		return 0;
	dev_t d = 0;
	bpf_probe_read_kernel(&d, sizeof(d), &sb->s_dev);
	__u64 dev = d;
	// Tombstone before the delete: a sync that read mountinfo while this sb lived can only put the
	// dev back once the tombstone already marks it stale.
	__u32 k = DEV_SEQ_NOW;
	__u64 *seq = bpf_map_lookup_elem(&netg_dev_seq, &k);
	__u64 now = 0;
	if (seq) {
		__sync_fetch_and_add(seq, 1);
		now = *seq; // >= this death's number: a racing death only makes the tombstone stricter
	}
	if (!seq || bpf_map_update_elem(&netg_dead_devs, &dev, &now, BPF_ANY)) {
		k = DEV_SEQ_LOST;
		__u64 *lost = bpf_map_lookup_elem(&netg_dev_seq, &k);
		if (lost)
			__sync_fetch_and_add(lost, 1);
	}
	bpf_map_delete_elem(&netg_vouched_devs, &dev);
	return 0;
}

// LD_PRELOAD / LD_AUDIT make ld.so run code its caller chose before the image's own: set non-empty
// at the exec of a protected image, they mark it code-suspect whatever library they name. The env
// strings are read one bpf_loop step at a time (state in a per-CPU map: see the trust guard's
// launch_step); an unreadable or unfinished scan marks too.
#define ENV_BUF 256 // power of two
#define ENV_STEPS 16384

struct env_scan {
	__u64 pos;
	__u64 end;
	__u32 cont; // the next read continues an over-long string
	__u32 risky;
	__u32 done;
	char buf[ENV_BUF];
};

struct {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, struct env_scan);
} netg_env_scan SEC(".maps");

// env_loader: buf (n bytes read, NUL included) is LD_PRELOAD= or LD_AUDIT= with a value.
static __always_inline int env_loader(const char *b, long n)
{
	if (b[0] != 'L' || b[1] != 'D' || b[2] != '_')
		return 0;
	if (n > 12 && b[3] == 'P' && b[4] == 'R' && b[5] == 'E' && b[6] == 'L' && b[7] == 'O' &&
	    b[8] == 'A' && b[9] == 'D' && b[10] == '=')
		return 1;
	return n > 10 && b[3] == 'A' && b[4] == 'U' && b[5] == 'D' && b[6] == 'I' && b[7] == 'T' &&
	       b[8] == '=';
}

static long env_step(__u64 i, void *ctx)
{
	__u32 z = 0;
	struct env_scan *s = bpf_map_lookup_elem(&netg_env_scan, &z);
	if (!s)
		return 1;
	if (s->pos >= s->end) {
		s->done = 1;
		return 1;
	}
	long n = bpf_probe_read_user_str(s->buf, ENV_BUF, (void *)s->pos);
	if (n <= 0) {
		s->risky = 1;
		return 1;
	}
	if (!s->cont && env_loader(s->buf, n)) {
		s->risky = 1;
		return 1;
	}
	s->cont = n == ENV_BUF;
	s->pos += s->cont ? ENV_BUF - 1 : n;
	return 0;
}

SEC("tp_btf/sched_process_exec")
int netg_exec_env(unsigned long long *ctx)
{
	struct task_struct *p = (struct task_struct *)ctx[0];
	if (!p || !task_protected((struct task_struct *)bpf_get_current_task()))
		return 0;
	__u32 z = 0;
	struct env_scan *s = bpf_map_lookup_elem(&netg_env_scan, &z);
	struct mm_struct *mm = BPF_CORE_READ(p, mm);
	if (!s || !mm)
		return 0;
	s->pos = BPF_CORE_READ(mm, env_start);
	s->end = BPF_CORE_READ(mm, env_end);
	s->cont = 0;
	s->risky = 0;
	s->done = 0;
	long ret = bpf_loop(ENV_STEPS, env_step, NULL, 0);
	s = bpf_map_lookup_elem(&netg_env_scan, &z);
	if (s && ret >= 0 && s->done && !s->risky)
		return 0;
	if (mark_suspect(bpf_get_current_pid_tgid() >> 32, NETG_REASON_ENV))
		bpf_send_signal(9); // SIGKILL: unmarkable, and this hook can't refuse the exec
	return 0;
}

// emit_gate reports a refused access to a protected process: comm/pid are the caller's, fd the
// target's tgid (0: none).
static __always_inline void emit_gate(__u32 type, __u32 target)
{
	struct net_guard_event e = {};
	e.pid = bpf_get_current_pid_tgid() >> 32;
	e.tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
	e.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
	e.gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
	e.type = type;
	e.fd = target;
	e.blocked = 1;
	bpf_get_current_comm(e.comm, sizeof(e.comm));
	e.cgroup_id = bpf_get_current_cgroup_id();
	emit_event(&e);
}

// Memory access to a protected process (ptrace attach, process_vm_readv/writev, /proc/<pid>/mem,
// pidfd_getfd: PTRACE_MODE_ATTACH) would let the caller run code as it: only a process the
// whitelist admits may. Metadata-only access (PTRACE_MODE_READ) stays open. A thread of the same
// group never reaches this hook.
SEC("lsm/ptrace_access_check")
int netg_ptrace_access_check(unsigned long long *ctx)
{
	struct task_struct *child = (struct task_struct *)ctx[0];
	__u32 mode = (__u32)ctx[1];
	if (!child || !(mode & PTRACE_MODE_ATTACH) || !task_protected(child))
		return 0;
	if (task_admitted((struct task_struct *)bpf_get_current_task()))
		return 0;
	__u32 tgid = 0;
	bpf_probe_read_kernel(&tgid, sizeof(tgid), &child->tgid);
	emit_gate(NET_PTRACE, tgid);
	return -EPERM;
}

// A tracer attached before the exec (PTRACE_TRACEME, or to a process that then execs a protected
// image) gets the new image's memory with no further check: the exec is refused unless the tracer
// is admitted. ->parent is the tracer of a traced task.
SEC("lsm/bprm_check_security")
int netg_bprm_check(unsigned long long *ctx)
{
	struct linux_binprm *bprm = (struct linux_binprm *)ctx[0];
	struct file *file = NULL;
	struct inode *inode = NULL;
	if (!bprm)
		return 0;
	bpf_probe_read_kernel(&file, sizeof(file), &bprm->file);
	if (file)
		bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	struct inode_key k = {};
	if (!fill_inode_key(inode, &k))
		return 0;
	__u32 *mc = bpf_map_lookup_elem(&mc_multicall, &k);
	__u8 *action = bpf_map_lookup_elem(&guard_net_exe_actions, &k);
	if (!(mc && (*mc & MC_HAS_ALLOW)) && !(action && *action == GUARD_ALLOW))
		return 0;
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	__u32 ptrace_flags = 0;
	if (!task)
		return 0;
	bpf_probe_read_kernel(&ptrace_flags, sizeof(ptrace_flags), &task->ptrace);
	if (!ptrace_flags)
		return 0;
	struct task_struct *tracer = NULL;
	bpf_probe_read_kernel(&tracer, sizeof(tracer), &task->parent);
	if (tracer && task_admitted(tracer))
		return 0;
	__u32 ttgid = 0;
	if (tracer)
		bpf_probe_read_kernel(&ttgid, sizeof(ttgid), &tracer->tgid);
	emit_gate(NET_TRACED_EXEC, ttgid);
	return -EPERM;
}

// Processes older than the hooks are judged from /proc (preexisting.go), which names them by their
// pid in the reader's pid namespace; the marks key on the kernel's tgid. Inside a pid namespace the
// two differ and userspace can't see the global one: netg_procs reports both for every thread group
// the reader's namespace sees.
struct netg_proc {
	__u32 nspid;
	__u32 tgid;
	__u64 start; // leader start_time: tells a reused tgid from the process a pass already judged
};

// numbers[] is a flexible array: index it by hand.
static __always_inline struct upid upid_at(struct pid *p, __u32 level)
{
	struct upid up = {};
	bpf_probe_read_kernel(&up, sizeof(up), (void *)&p->numbers[0] + (__u64)level * sizeof(up));
	return up;
}

SEC("iter/task")
int netg_procs(struct bpf_iter__task *ctx)
{
	struct task_struct *task = ctx->task;
	if (!task)
		return 0;
	struct netg_proc r = {};
	__u32 tid = 0;
	bpf_probe_read_kernel(&tid, sizeof(tid), &task->pid);
	bpf_probe_read_kernel(&r.tgid, sizeof(r.tgid), &task->tgid);
	if (tid != r.tgid)
		return 0;
	// The reader's namespace is the one its own pid was allocated in: numbers[level].
	struct task_struct *me = (struct task_struct *)bpf_get_current_task();
	struct pid *mine = NULL, *theirs = NULL;
	bpf_probe_read_kernel(&mine, sizeof(mine), &me->thread_pid);
	bpf_probe_read_kernel(&theirs, sizeof(theirs), &task->thread_pid);
	if (!mine || !theirs)
		return 0;
	unsigned int lvl = 0, tlvl = 0;
	bpf_probe_read_kernel(&lvl, sizeof(lvl), &mine->level);
	bpf_probe_read_kernel(&tlvl, sizeof(tlvl), &theirs->level);
	if (lvl > 32 || tlvl < lvl) // MAX_PID_NS_LEVEL; a shallower process is outside the namespace
		return 0;
	struct upid ns = upid_at(mine, lvl), up = upid_at(theirs, lvl);
	if (!ns.ns || up.ns != ns.ns)
		return 0;
	r.nspid = up.nr;
	bpf_probe_read_kernel(&r.start, sizeof(r.start), &task->start_time);
	bpf_seq_write(ctx->meta->seq, &r, sizeof(r));
	return 0;
}
