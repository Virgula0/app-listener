# Compatibility

## Check your host first

```bash
make check-compatibility
```

It runs every static check (kernel version, kernel `.config`, BTF, BPF-LSM activation, fscrypt
prerequisites, BPF sysctls) and says plainly whether the host can run app-listener. The
[one-line installer](installation.md#one-line-installer) runs it for you and refuses to install when
it fails. `make build` runs the toolchain in Docker, so nothing has to be installed on the host to
build it.

To also test a specific binary, run the script directly as root:

```bash
sudo bash scripts/check-compatibility.sh --binary /usr/local/sbin/app-listener
```

It loads that binary's guard eBPF into the running kernel's verifier (`daemon --check`), which
catches a prebuilt binary that this kernel rejects.

## Distributions

| Distribution | Min. kernel | Support | Activation |
|---|---|---|---|
| **Arch Linux** | rolling (6.x) | ✅ **Fully supported**: every LSM hook attaches | Add `bpf` to the `lsm=` list on your boot entry's kernel cmdline, reboot |
| **Ubuntu 24.04 LTS** | 6.8 | ✅ **Fully supported** | Append `lsm=landlock,lockdown,yama,integrity,apparmor,bpf` to `GRUB_CMDLINE_LINUX_DEFAULT`, `sudo update-grub`, reboot |
| **Ubuntu 22.04 LTS** | 6.x (HWE) | 🟡 **HWE kernel only**. On the 5.15 GA kernel `guard` and `daemon` do not load (they need `bpf_loop`, kernel ≥ 5.17); `monitor` and `network-guard -b` still work there (`network-guard -w` needs 5.17 too). Install the HWE kernel (`linux-generic-hwe-22.04`, 6.x) for full support | same as 24.04 |

Ubuntu 20.04 (kernel 5.4) is not supported.

## BPF-LSM must be active, not just compiled in

**Stock Ubuntu and cloud images compile `CONFIG_BPF_LSM=y` but do not activate it.** Without the
cmdline change above, LSM hooks attach but never deny anything. Every enforcing mode checks this at
startup and refuses to run instead of pretending to protect you.

Check what is active:

```bash
cat /sys/kernel/security/lsm    # must contain "bpf"
```

## Kernel floors per mode

| Mode | Minimum kernel | Why |
|---|---|---|
| `monitor` | 5.8 | BPF ring buffer |
| `network-monitor` | 5.8 | BPF ring buffer |
| `network-guard -b` | 5.10 | BPF-LSM |
| `network-guard -w` | 5.17 | its code-integrity hooks use `bpf_loop` |
| `guard`, `daemon` | 5.17 | BPF-LSM plus the `bpf_loop` helper |

### Mandatory and best-effort hooks

These must attach or the guard refuses to start:

- `file_open`, `file_permission`
- `inode_unlink`, `inode_rename`, `bprm_committed_creds` and the `sched_process_fork` tracepoint
  (replaced-binary tracking)
- the `sched_process_exec` tracepoint (launch scan of whitelisted runtimes)

Every other LSM hook is best-effort: a kernel that lacks one logs a warning and keeps enforcing the
rest.

## Architectures

Releases ship for `linux/amd64` (`x86_64`) and `linux/arm64` (`aarch64`); `install.sh` and `update`
pick the host's. On arm64 the LSM modes (`guard`, `daemon`, `network-guard`) need kernel ≥ 6.0 (the
arm64 BPF trampoline).

Source builds are native only: the eBPF is compiled against the build host's BTF, so build arm64 on
an arm64 host.
