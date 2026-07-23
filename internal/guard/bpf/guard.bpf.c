#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <errno.h>

#define MAX_PATH 256
#define MAX_COMM_LEN 16
#define MAY_READ  0x00000001
#define MAY_WRITE 0x00000002

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
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, char[MAX_PATH]);
} guard_path SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, char[MAX_PATH]);
} tmp_buf SEC(".maps"); // only used by guard_path_symlink for oldname_buf

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

// Check whether the operation targets a guarded path by examining
// both the file's own inode and its parent directory's inode.
static __always_inline int is_guarded_access(struct dentry *dentry, struct inode *file_inode)
{
	if (read_inode_guard(file_inode))
		return 1;

	if (!dentry)
		return 0;

	struct dentry *parent;
	bpf_probe_read_kernel(&parent, sizeof(parent), &dentry->d_parent);
	if (!parent || parent == dentry)
		return 0;

	struct inode *parent_inode;
	bpf_probe_read_kernel(&parent_inode, sizeof(parent_inode), &parent->d_inode);
	if (parent_inode && parent_inode != file_inode && read_inode_guard(parent_inode))
		return 1;

	return 0;
}

static __always_inline long read_path(struct dentry *dentry, char *buf, int buf_size)
{
    char tmp[64];
    int pos = MAX_PATH;

    pos--;
    buf[pos & (MAX_PATH - 1)] = '\0';

    struct dentry *d = dentry;

    // Notice we removed #pragma unroll. Modern verifiers handle bounded loops
    // beautifully without exploding the instruction limit.
    for (int i = 0; i < 32; i++) {
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

        if (pos < name_len)
            break;

        pos -= name_len;

        // Cache the masked position to help the verifier understand the baseline
        int start_pos = pos & (MAX_PATH - 1);

        // Removed #pragma unroll here as well
        for (int j = 0; j < 64; j++) {
            if (j >= name_len) break;
            buf[(start_pos + j) & (MAX_PATH - 1)] = tmp[j & 63];
        }

        if (pos <= 0)
            break;

        pos--;
        buf[pos & (MAX_PATH - 1)] = '/';

        d = parent;
    }

    // Return the offset directly (no bounds check needed here, already bounded)
    return pos;
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

    // 'off' is guaranteed to be >= 0 based on our new read_path implementation
    if (off < MAX_PATH) {
        off &= (MAX_PATH - 1);
        bpf_probe_read_kernel_str(out, MAX_PATH, buf + off);
    } else {
        out[0] = '\0';
    }
}

static __always_inline int check_and_emit(__u32 type, struct dentry *dentry, const char *dest_str, bool dest_is_user, struct dentry *dest_dentry)
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
			if (dest_is_user)
				bpf_probe_read_user_str(e->dest, sizeof(e->dest), dest_str);
			else
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

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (is_guarded_access(dentry, inode))
		return check_and_emit(EVENT_OPEN, dentry, NULL, false, NULL);

	return 0;
}

SEC("lsm/mmap_file")
int guard_mmap_file(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	if (!file)
		return 0;

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (is_guarded_access(dentry, inode))
		return check_and_emit(EVENT_MMAP, dentry, NULL, false, NULL);

	return 0;
}

SEC("lsm/file_permission")
int guard_file_permission(unsigned long long *ctx)
{
	struct file *file = (struct file *)ctx[0];
	int mask = (int)ctx[1];
	if (!file)
		return 0;

	// Only intercept read and write permission checks
	if (!(mask & (MAY_READ | MAY_WRITE)))
		return 0;

	struct dentry *dentry;
	bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
	if (!dentry)
		return 0;

	struct inode *inode;
	bpf_probe_read_kernel(&inode, sizeof(inode), &file->f_inode);

	if (!is_guarded_access(dentry, inode))
		return 0;

	__u32 event_type = (mask & MAY_WRITE) ? EVENT_WRITE : EVENT_READ;
	return check_and_emit(event_type, dentry, NULL, false, NULL);
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
		return check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_DELETE, dentry, NULL, false, NULL);
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
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_RENAME, old_dentry, NULL, false, new_dentry);
	}

	return 0;
}

SEC("lsm/path_symlink")
int guard_path_symlink(unsigned long long *ctx)
{
    struct dentry *dentry = (struct dentry *)ctx[1];
    const char *old_name = (const char *)ctx[2];

    // 1. Check if the symlink is being created INSIDE a guarded directory
    struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
    if (read_inode_guard(parent_inode)) {
        return check_and_emit(EVENT_SYMLINK, dentry, old_name, false, NULL);
    }

    // 2. Check if the symlink TARGET points INTO the guarded path
    __u32 key = 0;
    char *oldname_buf = bpf_map_lookup_elem(&tmp_buf, &key);
    if (!oldname_buf)
        return 0;

    long ret = bpf_probe_read_kernel_str(oldname_buf, MAX_PATH, old_name);
    if (ret <= 0)
        return 0;

    char *stored_path = bpf_map_lookup_elem(&guard_path, &key);
    if (!stored_path)
        return 0;

    // Prefix match logic to protect contents of a watched directory
    bool match = true;
    for (int i = 0; i < MAX_PATH; i++) {
        if (stored_path[i] == '\0') {
            // We reached the end of the guarded path.
            // It is a match if the symlink target ends here exactly,
            // continues as a sub-directory, or if the stored path had a trailing slash.
            if (oldname_buf[i] == '\0' || oldname_buf[i] == '/' || (i > 0 && stored_path[i-1] == '/')) {
                break;
            }
            // Otherwise, it's just a similarly named folder (e.g., /bypasses_fake)
            match = false;
            break;
        }
        if (oldname_buf[i] != stored_path[i]) {
            match = false;
            break;
        }
    }

    if (match) {
        return check_and_emit(EVENT_SYMLINK, dentry, old_name, false, NULL);
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
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry);
	}

	struct inode *parent_inode = get_inode_from_path((void *)ctx[1]);
	if (parent_inode && parent_inode != inode && read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_HARDLINK, old_dentry, NULL, false, new_dentry);
	}

	return 0;
}

SEC("lsm/path_mkdir")
int guard_path_mkdir(unsigned long long *ctx)
{
	struct inode *parent_inode = get_inode_from_path((void *)ctx[0]);
	struct dentry *dentry = (struct dentry *)ctx[1];

	if (read_inode_guard(parent_inode)) {
		return check_and_emit(EVENT_MKDIR, dentry, NULL, false, NULL);
	}

	return 0;
}
