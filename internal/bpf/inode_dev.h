// inode_dev.h: the device half of every inode key, as stat(2) reports it. Shared by the guard,
// trust and network-guard objects; userspace keys files by st_dev (KernelDev).
#ifndef INODE_DEV_H
#define INODE_DEV_H

#define BTRFS_SUPER_MAGIC 0x9123683EUL

// INODE_DEV_UNKNOWN: a btrfs inode whose subvolume can't be read. No userspace key (a KernelDev,
// below 2^32) has it, so such an inode is never whitelisted, trusted or guarded.
#define INODE_DEV_UNKNOWN (1ULL << 63)

// btrfs_layout: byte offsets userspace reads from btrfs's BTF (struct btrfs_inode's root and
// vfs_inode, struct btrfs_root's anon_dev). valid stays 0 until written.
struct btrfs_layout {
	__u32 valid;
	__u32 inode_root;
	__u32 inode_vfs;
	__u32 root_anon_dev;
};

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, struct btrfs_layout);
} btrfs_layout SEC(".maps");

// inode_dev returns the device stat(2) reports for the inode at address inode: sb->s_dev, except on
// btrfs, whose subvolumes share one superblock and repeat inode numbers. There stat reports the
// subvolume's anonymous device, and so must every key, or a user's file in one subvolume would
// carry the identity of another subvolume's binary. Global: the verifier checks it once, not at
// every call site, and requires inode to be a scalar: pass an address read with
// bpf_probe_read_kernel or ctx_ptr, never a typed ctx pointer.
__attribute__((noinline)) __u64 inode_dev(__u64 inode)
{
	struct super_block *sb = NULL;
	if (!inode)
		return INODE_DEV_UNKNOWN;
	bpf_probe_read_kernel(&sb, sizeof(sb), &((struct inode *)inode)->i_sb);
	if (!sb)
		return INODE_DEV_UNKNOWN;
	unsigned long magic = 0;
	bpf_probe_read_kernel(&magic, sizeof(magic), &sb->s_magic);
	dev_t dev = 0;
	if (magic != BTRFS_SUPER_MAGIC) {
		bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
		return dev;
	}
	__u32 zero = 0;
	struct btrfs_layout *l = bpf_map_lookup_elem(&btrfs_layout, &zero);
	if (!l || !l->valid)
		return INODE_DEV_UNKNOWN;
	__u64 root = 0;
	bpf_probe_read_kernel(&root, sizeof(root), (void *)(inode - l->inode_vfs + l->inode_root));
	if (!root)
		return INODE_DEV_UNKNOWN;
	bpf_probe_read_kernel(&dev, sizeof(dev), (void *)(root + l->root_anon_dev));
	return dev ? dev : INODE_DEV_UNKNOWN;
}

// ctx_ptr returns LSM argument i as a plain address. Read directly, an LSM program's pointer
// arguments are typed kernel pointers (PTR_TO_BTF_ID), which inode_dev's scalar argument refuses.
static __always_inline __u64 ctx_ptr(unsigned long long *ctx, int i)
{
	__u64 p = 0;
	bpf_probe_read_kernel(&p, sizeof(p), &ctx[i]);
	return p;
}

// sb_dev is the superblock's device: what mountinfo reports, the unit of mounts (vouched devs, the
// guarded-filesystem gate). On btrfs it is no file's key.
static __always_inline __u64 sb_dev(struct inode *inode)
{
	struct super_block *sb = NULL;
	if (!inode)
		return 0;
	bpf_probe_read_kernel(&sb, sizeof(sb), &inode->i_sb);
	if (!sb)
		return 0;
	dev_t dev = 0;
	bpf_probe_read_kernel(&dev, sizeof(dev), &sb->s_dev);
	return dev;
}

#endif
