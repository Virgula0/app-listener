#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <errno.h>

#define MAX_PATH 256
#define MAX_COMM_LEN 16

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

struct guard_event {
	__u32 pid;
	__u32 uid;
	__u32 gid;
	__u32 type;
	__u32 fd;
	__u32 blocked;
	char comm[MAX_COMM_LEN];
	char path[MAX_PATH];
	char dest[MAX_PATH];
};

struct comm_key {
	char comm[MAX_COMM_LEN];
};

struct inode_key {
	__u64 dev;
	__u64 ino;
};

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, __u64);
} guard_config SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 8192);
	__type(key, struct inode_key);
	__type(value, __u8);
} guard_inodes SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 64);
	__type(key, struct comm_key);
	__type(value, __u8);
} guard_comms SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} rb SEC(".maps");

char LICENSE[] SEC("license") = "GPL";

static __always_inline int read_inode_guard(struct inode *inode)
{
	if (!inode)
		return 0;

	struct inode_key ikey = {};
	bpf_probe_read_kernel(&ikey.ino, sizeof(ikey.ino), &inode->i_ino);

	dev_t dev;
	struct super_block *sb;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	ikey.dev = dev;

	__u8 *found = bpf_map_lookup_elem(&guard_inodes, &ikey);
	return found != NULL ? 1 : 0;
}

static __always_inline struct inode *get_inode_from_path(void *ctx_path)
{
	struct path *p = (struct path *)ctx_path;
	if (!p) return NULL;
	struct dentry *d;
	bpf_probe_read_kernel(&d, sizeof(d), &p->dentry);
	if (!d) return NULL;
	struct inode *i;
	bpf_probe_read_kernel(&i, sizeof(i), &d->d_inode);
	return i;
}

static __always_inline long read_path(struct dentry *dentry, char *buf, int buf_size)
{
	char tmp[64];
	int pos = buf_size;
	buf[--pos] = '\0';

	struct dentry *d = dentry;

	#pragma unroll
	for (int i = 0; i < 48; i++) {
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

		pos -= name_len;
		if (pos < 0)
			break;

		#pragma unroll
		for (int j = 0; j < 64; j++) {
			if (j >= name_len) break;
			buf[pos + j] = tmp[j];
		}

		if (pos <= 0)
			break;
		buf[--pos] = '/';

		d = parent;
	}

	return pos < buf_size - 1 ? pos : -1;
}

static __always_inline void fill_path(struct dentry *dentry, char *out)
{
	if (!dentry) {
		out[0] = '\0';
		return;
	}
	char buf[MAX_PATH] = {};
	long off = read_path(dentry, buf, sizeof(buf));
	if (off >= 0 && off < sizeof(buf)) {
		// Use bpf_probe_read_kernel_str to safely read from our own stack
		// This avoids the verifier's variable-offset stack read restrictions
		bpf_probe_read_kernel_str(out, MAX_PATH, buf + off);
	} else {
		out[0] = '\0';
	}
}

static __always_inline int check_and_emit(__u32 type, struct dentry *dentry, const char *dest_str, struct dentry *dest_dentry)
{
	__u32 key = 0;
	__u64 *mode = bpf_map_lookup_elem(&guard_config, &key);
	if (!mode)
		return 0;

	struct comm_key ck = {};
	bpf_get_current_comm(&ck, sizeof(ck));

	__u8 *found = bpf_map_lookup_elem(&guard_comms, &ck);

	int is_blocked;
	if (*mode == 0)
		is_blocked = found != NULL;
	else
		is_blocked = found == NULL;

	struct guard_event *e;
	e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (e) {
		e->pid = bpf_get_current_pid_tgid() >> 32;
		e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
		e->gid = (bpf_get_current_uid_gid() >> 32) & 0xFFFFFFFF;
		e->type = type;
		e->fd = 0;
		e->blocked = is_blocked ? 1 : 0;
		bpf_get_current_comm(e->comm, sizeof(e->comm));
		
		fill_path(dentry, e->path);

		if (dest_str) {
			bpf_probe_read_kernel_str(e->dest, sizeof(e->dest), dest_str);
		} else if (dest_dentry) {
			fill_path(dest_dentry, e->dest);
		} else {
			e->dest[0] = '\0';
		}

		bpf_ringbuf_submit(e, 0);
	}

	return is_blocked ? -EPERM : 0;
}

SEC("lsm/file_open")
int guard_file_open(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	if (!file)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (read_inode_guard(inode)) {
		struct dentry *dentry;
		bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
		return check_and_emit(EVENT_OPEN, dentry, NULL, NULL);
	}
	return 0;
}

SEC("lsm/mmap_file")
int guard_mmap_file(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	if (!file)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (read_inode_guard(inode)) {
		struct dentry *dentry;
		bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
		return check_and_emit(EVENT_MMAP, dentry, NULL, NULL);
	}
	return 0;
}

SEC("lsm/path_unlink")
int guard_path_unlink(unsigned long long *ctx)
{
	struct dentry *dentry = (struct dentry *)ctx[1];
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &dentry->d_inode);

	if (read_inode_guard(inode)) {
		return check_and_emit(EVENT_DELETE, dentry, NULL, NULL);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_DELETE, dentry, NULL, NULL);
	}

	return 0;
}

SEC("lsm/path_rename")
int guard_path_rename(unsigned long long *ctx)
{
	struct dentry *old_dentry = (struct dentry *)ctx[1];
	struct dentry *new_dentry = (struct dentry *)ctx[3];
	if (!old_dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &old_dentry->d_inode);

	if (read_inode_guard(inode)) {
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, new_dentry);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, new_dentry);
	}

	return 0;
}

SEC("lsm/path_symlink")
int guard_path_symlink(unsigned long long *ctx)
{
	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	struct dentry *dentry = (struct dentry *)ctx[1];
	const char *old_name = (const char *)ctx[2];

	if (read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_SYMLINK, dentry, old_name, NULL);
	}
	
	return 0;
}

SEC("lsm/path_link")
int guard_path_link(unsigned long long *ctx)
{
	struct dentry *old_dentry = (struct dentry *)ctx[0];
	struct dentry *new_dentry = (struct dentry *)ctx[2];
	if (!old_dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &old_dentry->d_inode);

	if (read_inode_guard(inode)) {
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, new_dentry);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[1]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, new_dentry);
	}
	
	return 0;
}

SEC("lsm/path_mkdir")
int guard_path_mkdir(unsigned long long *ctx)
{
	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	struct dentry *dentry = (struct dentry *)ctx[1];

	if (read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_MKDIR, dentry, NULL, NULL);
	}

	return 0;
}
