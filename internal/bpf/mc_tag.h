// mc_tag.h: per-applet identity for uutils multicall binaries, for the network objects
// (networkguard.bpf.c, networkmonitor.bpf.c). The guard object has its own copy in guard.bpf.c +
// exe_supersede.h and is intentionally NOT refactored onto this header (its verifier cost is pinned).
//
// One inode runs ~109 programs picked by the name it was exec'd as, so a plain dev:ino collapses
// them. The kernel attests the applet at exec (argv[0] and AT_EXECFN must agree and name an applet of
// that build's table) and stamps the process; a multicall exe's key then carries the applet in dev
// bits 32-47 (dev_t is 32 bits, so a tagged key never names a real file). Anything the kernel cannot
// attest (exec -a, fexecve, a script, a suffix name) stays MC_TAG_NONE, which no row admits.
//
// Includer must define `struct inode_key { __u64 dev; __u64 ino; }` and include inode_dev.h first.
#ifndef MC_TAG_H
#define MC_TAG_H

#define MC_TAG_SHIFT 32
#define MC_TAG_MASK  (0xffffULL << MC_TAG_SHIFT)
#define MC_TAG_NONE  0xffff
#define MC_NAME      16   // applet name, NUL-padded; MulticallNameMax in userspace
#define MC_PATH_MAX  4096 // bprm->filename / argv[0] read bound (power of two for the verifier mask)

// start is the thread-group leader's start_time: a reused tgid must not inherit the previous owner's
// tag. tag is the applet attested at the last exec, MC_TAG_NONE if none was.
struct mc_exec_stamp {
	__u64 start;
	__u32 tag;
	__u32 pad;
};

// mc_exe_stamps: tgid -> its multicall exec stamp. Written only for a multicall exec, inherited on
// fork, dropped when the process execs a non-multicall image or its leader exits. A multicall process
// with no row (exec'd before this object attached, or the map was full) keys as MC_TAG_NONE.
// MC_STAMPS_MAP_TYPE: an object that can't attach task_free (the monitor: no BPF-LSM requirement)
// uses an LRU map so leaked rows are evicted; an evicted live row degrades to NONE, never a tag.
#ifndef MC_STAMPS_MAP_TYPE
#define MC_STAMPS_MAP_TYPE BPF_MAP_TYPE_HASH
#endif
struct {
	__uint(type, MC_STAMPS_MAP_TYPE);
	__uint(max_entries, 65536);
	__type(key, __u32);
	__type(value, struct mc_exec_stamp);
} mc_exe_stamps SEC(".maps");

// mc_multicall: inodes whose applets are told apart (uutils). Presence only; userspace writes an
// inode's names (mc_multicall_names) before this row and drops neither while running.
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 64);
	__type(key, struct inode_key);
	__type(value, __u32);
} mc_multicall SEC(".maps");

struct mc_name_key {
	struct inode_key exe; // the multicall's real key (no tag)
	char name[MC_NAME];
};

// mc_multicall_names: that build's applets, by the basename that runs each, to their tag (>=1).
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 2048);
	__type(key, struct mc_name_key);
	__type(value, __u32);
} mc_multicall_names SEC(".maps");

struct mc_scratch_val {
	char path[MC_PATH_MAX];
	struct mc_name_key argv0;
	struct mc_name_key execfn;
};

struct {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, struct mc_scratch_val);
} mc_scratch SEC(".maps");

static __always_inline __u64 mc_leader_start(struct task_struct *task)
{
	struct task_struct *leader = NULL;
	__u64 start = 0;
	bpf_probe_read_kernel(&leader, sizeof(leader), &task->group_leader);
	if (leader)
		bpf_probe_read_kernel(&start, sizeof(start), &leader->start_time);
	return start;
}

// mc_tag_bits: what a multicall exe's key ORs into dev — the applet its process attested at exec,
// MC_TAG_NONE if none was. 0 for every other exe, whose key stays its plain inode.
static __noinline __u64 mc_tag_bits(struct task_struct *task, struct inode_key *exe)
{
	if (!bpf_map_lookup_elem(&mc_multicall, exe))
		return 0;
	__u64 tag = MC_TAG_NONE;
	__u32 tgid = 0;
	if (task) {
		bpf_probe_read_kernel(&tgid, sizeof(tgid), &task->tgid);
		struct mc_exec_stamp *st = bpf_map_lookup_elem(&mc_exe_stamps, &tgid);
		if (st && st->tag && st->start == mc_leader_start(task))
			tag = st->tag;
	}
	return tag << MC_TAG_SHIFT;
}

// mc_basename copies the last '/'-separated component of the string read into name (n: the read's
// return, NUL included), NUL-padded. 0 when empty, too long, or the read may have been truncated.
static __always_inline int mc_basename(const char *s, long n, char *name)
{
	if (n <= 1 || n >= MC_PATH_MAX)
		return 0;
	long len = n - 1;
	long start = -1;
	for (int k = 1; k <= MC_NAME; k++) {
		long idx = len - k;
		if (idx < 0) {
			start = 0;
			break;
		}
		if (s[idx & (MC_PATH_MAX - 1)] == '/') {
			start = idx + 1;
			break;
		}
	}
	long blen = len - start;
	if (start < 0 || blen <= 0 || blen >= MC_NAME)
		return 0;
	for (int k = 0; k < MC_NAME; k++)
		name[k] = k < blen ? s[(start + k) & (MC_PATH_MAX - 1)] : 0;
	return 1;
}

// mc_stamp_fork gives a new process its parent's stamp: it runs the same image. Only for a parent
// that is itself a live multicall process (has a row); a thread shares its group's row.
static __always_inline void mc_stamp_fork(struct task_struct *parent, struct task_struct *child)
{
	__u32 ptgid = 0, ctgid = 0;
	bpf_probe_read_kernel(&ptgid, sizeof(ptgid), &parent->tgid);
	bpf_probe_read_kernel(&ctgid, sizeof(ctgid), &child->tgid);
	if (ptgid == ctgid)
		return;
	struct mc_exec_stamp *pst = bpf_map_lookup_elem(&mc_exe_stamps, &ptgid);
	if (!pst) {
		bpf_map_delete_elem(&mc_exe_stamps, &ctgid); // a reused pid's leftover
		return;
	}
	struct mc_exec_stamp st = {};
	int same = pst->start == mc_leader_start(parent);
	st.tag = same ? pst->tag : MC_TAG_NONE;
	st.start = mc_leader_start(child);
	bpf_map_update_elem(&mc_exe_stamps, &ctgid, &st, BPF_ANY);
}

// mc_stamp_free drops a leader's row, unless a new process already reuses its tgid.
static __always_inline void mc_stamp_free(struct task_struct *task, __u32 tgid)
{
	struct mc_exec_stamp *st = bpf_map_lookup_elem(&mc_exe_stamps, &tgid);
	if (st && st->start == mc_leader_start(task))
		bpf_map_delete_elem(&mc_exe_stamps, &tgid);
}

// mc_attest is the sched_process_exec body: for a multicall image it stamps the process with the
// attested applet (or NONE); for any other image it drops a stale stamp so the tag never outlives
// its exec. ctx is the raw tracepoint args: ctx[0]=task, ctx[2]=linux_binprm.
static __always_inline void mc_attest(unsigned long long *ctx)
{
	struct task_struct *p = (struct task_struct *)ctx[0];
	struct linux_binprm *bprm = (struct linux_binprm *)ctx[2];
	if (!p || !bprm)
		return;
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	__u32 tgid = bpf_get_current_pid_tgid() >> 32;

	struct file *file = NULL;
	bpf_probe_read_kernel(&file, sizeof(file), &bprm->file);
	struct inode *inode = NULL;
	if (file)
		bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);
	struct inode_key real = {};
	if (inode) {
		bpf_probe_read_kernel(&real.ino, sizeof(real.ino), &inode->i_ino);
		real.dev = inode_dev((__u64)inode);
	}
	if (!inode || !bpf_map_lookup_elem(&mc_multicall, &real)) {
		bpf_map_delete_elem(&mc_exe_stamps, &tgid); // exec'd a non-multicall image: no tag
		return;
	}

	struct mc_exec_stamp st = {};
	st.start = mc_leader_start(task);
	st.tag = MC_TAG_NONE;

	// A script or binfmt_misc image runs with an argv the kernel rewrote: filename != interp.
	const char *filename = NULL, *interp = NULL;
	bpf_probe_read_kernel(&filename, sizeof(filename), &bprm->filename);
	bpf_probe_read_kernel(&interp, sizeof(interp), &bprm->interp);
	if (filename && filename == interp) {
		__u32 z = 0;
		struct mc_scratch_val *s = bpf_map_lookup_elem(&mc_scratch, &z);
		struct mm_struct *mm = BPF_CORE_READ(p, mm);
		if (s && mm) {
			s->execfn.exe = real;
			s->argv0.exe = real;
			long n = bpf_probe_read_kernel_str(s->path, MC_PATH_MAX, filename);
			if (mc_basename(s->path, n, s->execfn.name)) {
				unsigned long arg_start = BPF_CORE_READ(mm, arg_start);
				n = bpf_probe_read_user_str(s->path, MC_PATH_MAX, (void *)arg_start);
				if (mc_basename(s->path, n, s->argv0.name)) {
					__u64 *a = (__u64 *)s->argv0.name, *b = (__u64 *)s->execfn.name;
					if (a[0] == b[0] && a[1] == b[1]) {
						__u32 *tag = bpf_map_lookup_elem(&mc_multicall_names, &s->execfn);
						if (tag)
							st.tag = *tag;
					}
				}
			}
		}
	}
	bpf_map_update_elem(&mc_exe_stamps, &tgid, &st, BPF_ANY);
}

#endif
