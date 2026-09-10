#!/usr/bin/env bash
#
# check-compatibility — static host compatibility check for app-listener.
#
# It performs every check that can be done WITHOUT loading eBPF or running
# the integration suite: kernel version, kernel .config, BTF, the BPF LSM
# activation, fscrypt prerequisites and BPF sysctls. Build tooling is NOT
# checked — `make build` runs the whole toolchain in a Docker container.
# The final verdict is binary: the host can run app-listener, or it cannot
# (with the exact blockers listed). Warnings never change the verdict.
#
# Usage:
#   make check-compatibility        (or: bash scripts/check-compatibility.sh)
#   bash scripts/check-compatibility.sh --binary /usr/local/sbin/app-listener
#
#   --binary <path>   additionally load every guard eBPF program in <path> into
#                     THIS kernel's verifier (runs "<path> daemon --check").
#                     Catches a prebuilt-release / kernel mismatch (issue #45)
#                     at install time instead of at "systemctl start". Needs
#                     root and an active BPF LSM; skipped with a note otherwise.
#   --binary-only     run only the checks the eBPF probe depends on (kernel
#                     version, BTF, BPF-LSM activation) plus the probe itself.
#                     Used by scripts/install.sh to re-check just the download.
#
# Exit code: 0 = installable, 1 = a hard requirement is missing (or the
# --binary program's guard eBPF is rejected by this kernel).
#
# Test hooks (not for end users):
#   CHECK_KERNEL / CHECK_BTF_PATH / CHECK_LSM_PATH / CHECK_CONFIG_PATH /
#   CHECK_CMDLINE_PATH / CHECK_OS_RELEASE — override the probed sources.

set -uo pipefail

CHECK_BINARY=""
BINARY_ONLY=0
while [ $# -gt 0 ]; do
	case "$1" in
	--binary)      CHECK_BINARY="${2:-}"; shift 2 || { echo "--binary needs a path" >&2; exit 2; } ;;
	--binary=*)    CHECK_BINARY="${1#*=}"; shift ;;
	--binary-only) BINARY_ONLY=1; shift ;;
	-h | --help)   sed -n '2,30p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
	*)             echo "unknown argument: $1 (try --help)" >&2; exit 2 ;;
	esac
done

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[0;33m'; BOLD='\033[1m'; NC='\033[0m'

PASS=0; WARN=0; FAIL=0
FAILED_ITEMS=()

pass() { PASS=$((PASS + 1)); printf "  ${GREEN}[OK]${NC}    %s\n" "$*"; }
warn() { WARN=$((WARN + 1)); printf "  ${YELLOW}[WARN]${NC}  %s\n" "$*"; }
fail() { FAIL=$((FAIL + 1)); FAILED_ITEMS+=("$*"); printf "  ${RED}[FAIL]${NC}  %s\n" "$*"; }
note() { printf "          %s\n" "$*"; }
section() { printf "\n${BOLD}%s${NC}\n" "$*"; }

version_at_least() { [ "$(printf '%s\n' "$1" "$2" | sort -V | head -1)" = "$1" ]; }

##############################################################################
# Sources
##############################################################################

KERNEL_RELEASE="${CHECK_KERNEL:-$(uname -r)}"
KVER="$(printf '%s' "$KERNEL_RELEASE" | grep -oE '^[0-9]+\.[0-9]+' || true)"
BTF_PATH="${CHECK_BTF_PATH:-/sys/kernel/btf/vmlinux}"
LSM_PATH="${CHECK_LSM_PATH:-/sys/kernel/security/lsm}"
CMDLINE_PATH="${CHECK_CMDLINE_PATH:-/proc/cmdline}"
OS_RELEASE="${CHECK_OS_RELEASE:-/etc/os-release}"

DISTRO_ID=""; DISTRO_PRETTY=""
if [ -r "$OS_RELEASE" ]; then
	DISTRO_ID="$(sed -n 's/^ID=//p' "$OS_RELEASE" | tr -d '"')"
	DISTRO_PRETTY="$(sed -n 's/^PRETTY_NAME=//p' "$OS_RELEASE" | tr -d '"')"
fi

# read_config <CONFIG_NAME> — echoes the value (y/m/…) or empty; sets
# CONFIG_SOURCE the first time a config file is found.
CONFIG_SOURCE=""
CONFIG_CACHE=""
load_config() {
	if [ -n "${CHECK_CONFIG_PATH:-}" ] && [ -r "$CHECK_CONFIG_PATH" ]; then
		CONFIG_SOURCE="$CHECK_CONFIG_PATH"
		CONFIG_CACHE="$(cat "$CHECK_CONFIG_PATH")"
	elif [ -r /proc/config.gz ]; then
		CONFIG_SOURCE="/proc/config.gz"
		CONFIG_CACHE="$(zcat /proc/config.gz 2>/dev/null)"
	elif [ -r "/boot/config-$KERNEL_RELEASE" ]; then
		CONFIG_SOURCE="/boot/config-$KERNEL_RELEASE"
		CONFIG_CACHE="$(cat "/boot/config-$KERNEL_RELEASE")"
	fi
}
read_config() { printf '%s\n' "$CONFIG_CACHE" | sed -n "s/^$1=//p" | head -1; }

# require_config <NAME> <purpose> <hardness: fail|warn>
require_config() {
	local name="$1" purpose="$2" hard="$3" val
	[ -z "$CONFIG_SOURCE" ] && return 0
	val="$(read_config "$name")"
	if [ "$val" = "y" ] || [ "$val" = "m" ]; then
		pass "$name=$val ($purpose)"
	elif [ "$hard" = "fail" ]; then
		fail "$name is not set — $purpose (rebuild the kernel with $name=y)"
	else
		warn "$name is not set — $purpose"
	fi
}

lsm_bpf_instructions() {
	case "$DISTRO_ID" in
	ubuntu | debian)
		note "Ubuntu/Debian: append 'bpf' to the LSM list on the kernel cmdline:"
		note "  echo 'GRUB_CMDLINE_LINUX_DEFAULT=\"\$GRUB_CMDLINE_LINUX_DEFAULT lsm=landlock,lockdown,yama,integrity,apparmor,bpf\"' | sudo tee /etc/default/grub.d/bpf-lsm.cfg"
		note "  sudo update-grub && sudo reboot"
		note "  cloud images also override the cmdline in /etc/default/grub.d/50-cloudimg-settings.cfg — edit there too" ;;
	arch)
		note "Arch: add 'lsm=...,bpf' to the kernel cmdline of your boot entry:"
		note "  systemd-boot: append to the 'options' line in /boot/loader/entries/*.conf"
		note "  GRUB: add to GRUB_CMDLINE_LINUX_DEFAULT in /etc/default/grub, then grub-mkconfig -o /boot/grub/grub.cfg"
		note "  then reboot" ;;
	*)
		note "Add 'bpf' to the kernel's 'lsm=' cmdline parameter (keep the existing entries) and reboot." ;;
	esac
	note "verify after boot:  grep -o bpf $LSM_PATH"
}

##############################################################################

printf "${BOLD}app-listener — host compatibility${NC}\n"
[ -n "$DISTRO_PRETTY" ] && printf "  detected: %s, kernel %s\n" "$DISTRO_PRETTY" "$KERNEL_RELEASE"
load_config
if [ -n "$CONFIG_SOURCE" ]; then
	printf "  kernel config: %s\n" "$CONFIG_SOURCE"
else
	printf "  kernel config: ${YELLOW}not found${NC} (no /proc/config.gz, no /boot/config-%s) — config checks skipped\n" "$KERNEL_RELEASE"
fi

##############################################################################
section "Kernel version"
##############################################################################

if [ -z "$KVER" ]; then
	fail "cannot parse kernel version from '$KERNEL_RELEASE'"
else
	if version_at_least "5.8" "$KVER"; then
		pass "kernel $KERNEL_RELEASE (>= 5.8, monitor mode)"
	else
		fail "kernel $KERNEL_RELEASE is too old — monitor needs >= 5.8 (BPF ring buffer)"
	fi
	if version_at_least "5.10" "$KVER"; then
		pass "kernel $KERNEL_RELEASE (>= 5.10, guard / network-guard / daemon)"
	else
		fail "kernel $KERNEL_RELEASE < 5.10 — guard / network-guard / daemon need BPF-LSM bpf_link"
	fi
	version_at_least "6.2" "$KVER" || warn "kernel < 6.2 — the file_truncate LSM hook is unavailable (ftruncate on a pre-opened fd is not denied; path truncate still is)"
fi

##############################################################################
section "BTF (required to load the pre-compiled CO-RE eBPF programs)"
##############################################################################

if [ -r "$BTF_PATH" ]; then
	pass "BTF present ($BTF_PATH)"
else
	fail "BTF not found at $BTF_PATH — kernel needs CONFIG_DEBUG_INFO_BTF=y"
fi

if [ "$BINARY_ONLY" -eq 0 ]; then
##############################################################################
section "Kernel configuration"
##############################################################################

if [ -z "$CONFIG_SOURCE" ]; then
	warn "kernel .config not readable — cannot verify CONFIG_* options statically"
else
	require_config CONFIG_BPF_SYSCALL     "eBPF core"                              fail
	require_config CONFIG_DEBUG_INFO_BTF  "CO-RE BTF"                              fail
	require_config CONFIG_BPF_LSM         "guard / network-guard / daemon"        fail
	require_config CONFIG_KPROBES         "monitor (VFS kprobes)"                 fail
	require_config CONFIG_BPF_EVENTS      "attaching BPF to kprobes/tracepoints"  fail
	require_config CONFIG_SECURITY_PATH   "path-based LSM hooks (unlink/rename/mkdir/symlink/…)"  warn
	require_config CONFIG_FS_ENCRYPTION   "fscrypt (daemon + installer)"          warn
fi
fi  # BINARY_ONLY

##############################################################################
section "BPF LSM activation (guard / network-guard / daemon)"
##############################################################################

lsm_list=""
lsm_readable=0
if [ -r "$LSM_PATH" ]; then
	lsm_readable=1
	lsm_list="$(tr -d '[:space:]' < "$LSM_PATH")"
	pass "securityfs readable ($LSM_PATH: $lsm_list)"
else
	fail "cannot read $LSM_PATH — securityfs not mounted (mount -t securityfs securityfs /sys/kernel/security)"
fi

BPF_LSM_ACTIVE=0
if [ "$lsm_readable" -eq 1 ] && printf '%s' "$lsm_list" | tr ',' '\n' | grep -qx 'bpf'; then
	BPF_LSM_ACTIVE=1
	pass "the BPF LSM is ACTIVE — guard modes can deny"
elif [ "$lsm_readable" -eq 1 ]; then
	fail "the BPF LSM is NOT active ($LSM_PATH has no 'bpf') — guard modes would attach but never deny"
	if [ -n "$CONFIG_SOURCE" ] && [ "$(read_config CONFIG_BPF_LSM)" != "y" ]; then
		note "this kernel also lacks CONFIG_BPF_LSM=y — a kernel rebuild or a different kernel is required."
	else
		lsm_bpf_instructions
	fi
	if [ -r "$CMDLINE_PATH" ] && ! grep -qE 'lsm=' "$CMDLINE_PATH"; then
		note "(current cmdline has no 'lsm=' parameter, so the kernel is using its compiled-in list)"
	fi
fi

##############################################################################
section "Guard eBPF — does THIS build load on THIS kernel? (issue #45)"
##############################################################################

if [ -z "$CHECK_BINARY" ]; then
	note "skipped — pass '--binary <path-to-app-listener>' to load every guard"
	note "eBPF program into this kernel's verifier. scripts/install.sh and"
	note "'sudo app-listener install' run this automatically against the binary"
	note "they are about to install."
elif [ ! -f "$CHECK_BINARY" ]; then
	warn "guard eBPF probe skipped — '$CHECK_BINARY' is not a file"
elif [ "$BPF_LSM_ACTIVE" -ne 1 ]; then
	warn "guard eBPF probe skipped — the BPF LSM is not active on this kernel (see above)"
elif [ "$(id -u)" -ne 0 ]; then
	warn "guard eBPF probe skipped — re-run as root to load the programs into the verifier"
else
	if probe_out="$("$CHECK_BINARY" daemon --check 2>&1)"; then
		pass "every guard eBPF program is accepted by this kernel's verifier"
	else
		fail "this build's guard eBPF is REJECTED by this kernel's verifier — a prebuilt-release / kernel mismatch (issue #45), not a fault in your host. Rebuild from source so the eBPF is compiled against this kernel: git clone https://github.com/Virgula0/app-listener && cd app-listener && make build && sudo ./build/linux/app-listener install"
		printf '%s\n' "$probe_out" | sed 's/^/          | /'
	fi
fi

if [ "$BINARY_ONLY" -eq 0 ]; then
##############################################################################
section "fscrypt (daemon + installer)"
##############################################################################

if command -v fscrypt >/dev/null 2>&1; then
	pass "fscrypt CLI found ($(command -v fscrypt))"
else
	warn "fscrypt not found — the installer's encryption lifecycle needs it (Ubuntu: apt install fscrypt; Arch: pacman -S fscrypt)"
fi

home_fs="$(findmnt -no FSTYPE -T "${HOME:-/root}" 2>/dev/null || true)"
case "$home_fs" in
ext4 | f2fs) pass "home filesystem is $home_fs (fscrypt-capable)" ;;
"")          warn "could not determine the filesystem backing \$HOME — fscrypt needs ext4 or f2fs" ;;
*)           warn "home filesystem is '$home_fs' — fscrypt needs ext4 or f2fs; run the daemon with need_encryption: false, or move protected dirs onto ext4/f2fs" ;;
esac

bpffs="$(findmnt -no FSTYPE /sys/fs/bpf 2>/dev/null || true)"
[ -z "$bpffs" ] && bpffs="$(stat -f -c %T /sys/fs/bpf 2>/dev/null || true)"
case "$bpffs" in
bpf | bpffs | bpf_fs)
	pass "/sys/fs/bpf is a bpffs mount (daemon pins LSM links there so they survive a SIGKILL)" ;;
*)
	warn "/sys/fs/bpf is not a bpffs mount — the daemon will try to mount one; if it cannot, guards still enforce but do not survive a SIGKILL" ;;
esac

##############################################################################
section "BPF runtime"
##############################################################################

ubpf="$(cat /proc/sys/kernel/unprivileged_bpf_disabled 2>/dev/null || true)"
case "$ubpf" in
1 | 2) pass "unprivileged BPF disabled (=$ubpf) — run app-listener as root (it is meant to)" ;;
*)     pass "unprivileged_bpf_disabled=${ubpf:-n/a}" ;;
esac

jit="$(cat /proc/sys/net/core/bpf_jit_enable 2>/dev/null || true)"
if [ "$jit" = "0" ]; then
	warn "net.core.bpf_jit_enable=0 — some programs may fail to load; set it to 1"
else
	pass "BPF JIT enabled (net.core.bpf_jit_enable=${jit:-n/a})"
fi

[ "$(id -u)" -ne 0 ] && warn "not running this check as root — 'app-listener install' and every mode need root (sudo)"

# No build-toolchain checks: `make build` runs the whole toolchain inside a
# Docker container, so nothing (Go, clang, bpftool, GCC) has to be installed
# on the host to build or install app-listener.
fi  # BINARY_ONLY

##############################################################################
section "Verdict"
##############################################################################

if [ "$FAIL" -eq 0 ]; then
	printf "\n  ${GREEN}${BOLD}app-listener CAN be installed on this host.${NC}"
	[ "$WARN" -gt 0 ] && printf "  ${YELLOW}(%d warning(s) above — read them; none block installation.)${NC}" "$WARN"
	printf "\n"
	exit 0
fi

printf "\n  ${RED}${BOLD}app-listener CANNOT be installed on this host.${NC}  %d blocker(s):\n" "$FAIL"
for item in "${FAILED_ITEMS[@]}"; do
	printf "    ${RED}•${NC} %s\n" "$item"
done
printf "  Fix the blockers (reboot where asked) and re-run this check.\n"
exit 1
