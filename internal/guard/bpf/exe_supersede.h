// exe_supersede.h: a whitelisted binary's identity over time, shared by guard.bpf.c and
// guard_trust.bpf.c. Userspace hands both objects the same three maps.
//
// An exe key is (dev, ino), and both outlive the file. Once the last link goes, a process can still
// exec the image through a held fd or /proc/<pid>/exe; once the inode is freed, the filesystem hands
// its number to the next file created, by any user. A key that lost its last link is superseded: it
// keeps admitting processes that exec'd it before, and refuses every later exec, the reused number's
// included. One counter (exe_seq) orders execs against supersedes.
//
// Includer defines struct inode_key first.
#ifndef EXE_SUPERSEDE_H
#define EXE_SUPERSEDE_H

struct exe_supersede {
	__u64 seq;   // exe_seq when the last link went
	__u64 freed; // the inode is gone: userspace may drop the key's rows
};

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 4096);
	__type(key, struct inode_key);
	__type(value, struct exe_supersede);
} exe_superseded SEC(".maps");

// start is the thread-group leader's start_time: a pid can be reused before task_free runs for its
// previous owner, whose row must not pass for the new process's.
struct exec_stamp {
	__u64 seq;
	__u64 start;
};

// exe_stamps: tgid -> exe_seq at its last exec. Inherited on fork, rewritten on exec, cleared on
// leader exit. A tgid with no row exec'd before the guard attached.
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 65536);
	__type(key, __u32);
	__type(value, struct exec_stamp);
} exe_stamps SEC(".maps");

#define EXE_SEQ_NOW 0
#define EXE_SEQ_LOST 1 // stamps that could not be recorded: a missing row proves nothing any more

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 2);
	__type(key, __u32);
	__type(value, __u64);
} exe_seq SEC(".maps");

static __always_inline __u64 leader_start(struct task_struct *task)
{
	struct task_struct *leader = NULL;
	__u64 start = 0;
	bpf_probe_read_kernel(&leader, sizeof(leader), &task->group_leader);
	if (leader)
		bpf_probe_read_kernel(&start, sizeof(start), &leader->start_time);
	return start;
}

// exe_refused: exe is superseded and task's image was exec'd at or after that, or its stamp can't
// be trusted. Callers fold the result into a lookup key instead of branching on it.
static __noinline int exe_refused(struct task_struct *task, struct inode_key *exe)
{
	struct exe_supersede *s = bpf_map_lookup_elem(&exe_superseded, exe);
	if (!s)
		return 0;
	if (!task)
		return 1;
	__u32 tgid = 0;
	bpf_probe_read_kernel(&tgid, sizeof(tgid), &task->tgid);
	struct exec_stamp *st = bpf_map_lookup_elem(&exe_stamps, &tgid);
	if (st)
		return st->start != leader_start(task) || st->seq >= s->seq;
	__u32 k = EXE_SEQ_LOST;
	__u64 *lost = bpf_map_lookup_elem(&exe_seq, &k);
	return !lost || *lost != 0;
}

static __noinline int exe_is_superseded(struct inode_key *exe)
{
	return bpf_map_lookup_elem(&exe_superseded, exe) != NULL;
}

static __always_inline void exe_stamp_lost(__u32 tgid)
{
	bpf_map_delete_elem(&exe_stamps, &tgid);
	__u32 k = EXE_SEQ_LOST;
	__u64 *lost = bpf_map_lookup_elem(&exe_seq, &k);
	if (lost)
		__sync_fetch_and_add(lost, 1);
}

// exe_stamp_exec records the current task's exec. The increment comes first, so an exec after a
// supersede always stamps later than it; a racing increment only makes the stamp later (stricter).
static __always_inline void exe_stamp_exec(void)
{
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	__u32 tgid = bpf_get_current_pid_tgid() >> 32;
	__u32 k = EXE_SEQ_NOW;
	__u64 *seq = bpf_map_lookup_elem(&exe_seq, &k);
	if (seq && task) {
		__sync_fetch_and_add(seq, 1);
		struct exec_stamp st = {};
		st.seq = *seq;
		st.start = leader_start(task);
		if (!bpf_map_update_elem(&exe_stamps, &tgid, &st, BPF_ANY))
			return;
	}
	exe_stamp_lost(tgid);
}

// exe_stamp_fork gives a new process its parent's stamp: it runs the same image. A parent row left
// by another task stamps the child as late as possible.
static __always_inline void exe_stamp_fork(struct task_struct *parent, struct task_struct *child)
{
	__u32 ptgid = 0, ctgid = 0;
	bpf_probe_read_kernel(&ptgid, sizeof(ptgid), &parent->tgid);
	bpf_probe_read_kernel(&ctgid, sizeof(ctgid), &child->tgid);
	if (ptgid == ctgid)
		return; // a thread shares its group's row
	struct exec_stamp *pst = bpf_map_lookup_elem(&exe_stamps, &ptgid);
	if (!pst) {
		bpf_map_delete_elem(&exe_stamps, &ctgid); // a reused pid's leftover
		return;
	}
	struct exec_stamp st = {};
	st.seq = pst->start == leader_start(parent) ? pst->seq : ~0ULL;
	st.start = leader_start(child);
	if (bpf_map_update_elem(&exe_stamps, &ctgid, &st, BPF_ANY))
		exe_stamp_lost(ctgid);
}

// exe_stamp_free drops a leader's row, unless a new process already reuses its pid.
static __always_inline void exe_stamp_free(struct task_struct *task, __u32 tgid)
{
	struct exec_stamp *st = bpf_map_lookup_elem(&exe_stamps, &tgid);
	if (st && st->start == leader_start(task))
		bpf_map_delete_elem(&exe_stamps, &tgid);
}

static __always_inline __u64 exe_seq_now(void)
{
	__u32 k = EXE_SEQ_NOW;
	__u64 *seq = bpf_map_lookup_elem(&exe_seq, &k);
	return seq ? *seq : 0;
}

// exe_note_supersede marks k superseded. An existing row keeps its seq: a reused number stays
// refused from the first supersede on. Nonzero when it could not be recorded.
static __always_inline int exe_note_supersede(struct inode_key *k, __u64 freed)
{
	struct exe_supersede *cur = bpf_map_lookup_elem(&exe_superseded, k);
	if (cur) {
		cur->freed = freed;
		return 0;
	}
	struct exe_supersede v = {};
	v.seq = exe_seq_now();
	v.freed = freed;
	long err = bpf_map_update_elem(&exe_superseded, k, &v, BPF_NOEXIST);
	return err && err != -EEXIST;
}

// exe_last_link_key fills the key of a regular file about to lose its last link. Called from
// inode_unlink/inode_rename, which run after the VFS permission checks: a caller that may not
// delete the file never gets here. i_nlink is the count before the drop.
static __always_inline int exe_last_link_key(struct inode *inode, struct inode_key *k)
{
	if (!inode)
		return 0;
	umode_t mode = 0;
	unsigned int nlink = 0;
	bpf_probe_read_kernel(&mode, sizeof(mode), &inode->i_mode);
	bpf_probe_read_kernel(&nlink, sizeof(nlink), &inode->i_nlink);
	if ((mode & 00170000) != 0100000 || nlink != 1)
		return 0;
	bpf_probe_read_kernel(&k->ino, sizeof(k->ino), &inode->i_ino);
	struct super_block *sb = NULL;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	dev_t dev = 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	k->dev = dev;
	return 1;
}

// exe_note_freed marks k's inode gone (inode_free_security, i_nlink 0). A member key that lost its
// last link unseen (a path no inode hook covers) is superseded here, before its number is reused.
static __always_inline void exe_note_freed(struct inode_key *k, int member)
{
	struct exe_supersede *cur = bpf_map_lookup_elem(&exe_superseded, k);
	if (cur) {
		cur->freed = 1;
		return;
	}
	if (member)
		exe_note_supersede(k, 1);
}

#endif
