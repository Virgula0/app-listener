#!/usr/bin/env bash
#
# check-compatibility — exhaustive host compatibility check for
# app-listener (build, install, and Docker integration tests).
#
# Usage:
#   make check-compatibility       (or directly: bash scripts/check-compatibility.sh)
#
# Exit code: 0 when everything needed is present, 1 when a hard
# requirement is missing. Warnings never change the exit code.
#
# Internal test hooks (not meant for end users):
#   CHECK_LSM_PATH   alternate path for the LSM list file
#   CHECK_BTF_PATH   alternate path for the BTF vmlinux file
#   CHECK_KERNEL     alternate uname -r output

set -uo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BOLD='\033[1m'
NC='\033[0m'

PASS=0
WARN=0
FAIL=0

pass() { PASS=$((PASS + 1)); printf "  ${GREEN}[OK]${NC}   %s\n" "$*"; }
warn() { WARN=$((WARN + 1)); printf "  ${YELLOW}[WARN]${NC} %s\n" "$*"; }
fail() { FAIL=$((FAIL + 1)); printf "  ${RED}[FAIL]${NC}  %s\n" "$*"; }

section() { printf "\n${BOLD}%s${NC}\n" "$*"; }

version_at_least() {
	# version_at_least <required> <got> — true when got >= required (sort -V)
	[ "$(printf '%s\n' "$1" "$2" | sort -V | head -1)" = "$1" ]
}

##############################################################################
# Kernel: version, BTF, securityfs, BPF LSM, BPF sysctls
##############################################################################

KERNEL_RELEASE="${CHECK_KERNEL:-$(uname -r)}"
KERNEL_MAJOR="$(printf '%s' "$KERNEL_RELEASE" | cut -d. -f1)"
KERNEL_MINOR="$(printf '%s' "$KERNEL_RELEASE" | cut -d. -f2)"
BTF_PATH="${CHECK_BTF_PATH:-/sys/kernel/btf/vmlinux}"
LSM_PATH="${CHECK_LSM_PATH:-/sys/kernel/security/lsm}"

section "Kernel"
if [ -z "$KERNEL_MAJOR" ]; then
	fail "cannot determine kernel version (uname -r returned empty)"
elif [ "$KERNEL_MAJOR" -ge 5 ]; then
	pass "kernel $KERNEL_RELEASE (>= 5.x required)"
	if [ "$KERNEL_MAJOR" -eq 5 ] && [ "${KERNEL_MINOR:-0}" -lt 10 ]; then
		warn "kernel < 5.10: guard/network-guard/daemon LSM hooks may be unreliable (monitor still works)"
	fi
else
	fail "kernel $KERNEL_RELEASE is too old: 5.x or newer is required"
fi

section "BTF (needed to load the pre-compiled CO-RE eBPF programs)"
if [ -r "$BTF_PATH" ]; then
	pass "BTF present ($BTF_PATH)"
else
	fail "BTF not found at $BTF_PATH — kernel needs CONFIG_DEBUG_INFO_BTF=y (stock Ubuntu/Arch kernels ship it since 5.4)"
fi

lsm_exists=0
if [ -r "$LSM_PATH" ]; then
	lsm_exists=1
	pass "securityfs readable ($LSM_PATH)"
else
	fail "cannot read $LSM_PATH (is securityfs mounted? run: mount -t securityfs securityfs /sys/kernel/security)"
fi

section "BPF LSM (required by guard, network-guard and daemon modes)"
if [ "$lsm_exists" -eq 1 ] && tr ',' '\n' < "$LSM_PATH" | grep -qw 'bpf'; then
	pass "bpf is in the active LSM list ($(tr '\n' ',' < "$LSM_PATH" | sed 's/,$//' | tr -d '\n'))"
else
	fail "the BPF LSM is not active: $LSM_PATH does not list 'bpf'"
	printf "        ${RED}Without it the LSM hooks attach but NEVER fire: guard modes\n"
	printf "        ${RED}would seem to work while denying nothing. Enable it and reboot:\n"
	printf "        ${RED}  1) echo 'GRUB_CMDLINE_LINUX_DEFAULT=\"quiet splash lsm=landlock,lockdown,yama,integrity,apparmor,bpf\"' | sudo tee /etc/default/grub.d/bpf-lsm.cfg\n"
	printf "        ${RED}  2) sudo update-grub && sudo reboot\n"
	printf "        ${RED}  3) verify after boot: grep -o bpf $LSM_PATH\n"
	printf "        ${RED}Note: Ubuntu cloud images (EC2/Azure/GCE) override the cmdline in\n"
	printf "        ${RED}/etc/default/grub.d/50-cloudimg-settings.cfg — add 'lsm=...bpf' there too.\n"
fi

section "BPF runtime"
unprivileged_bpf="$(cat /proc/sys/kernel/unprivileged_bpf_disabled 2>/dev/null || true)"
case "$unprivileged_bpf" in
1) warn "kernel.unprivileged_bpf_disabled=1: eBPF loads only as privileged user (root / CAP_BPF)" ;;
2) warn "kernel.unprivileged_bpf_disabled=2: unprivileged BPF fully disabled — run app-listener as root" ;;
*) pass "unprivileged BPF is not disabled (value ${unprivileged_bpf:-n/a})" ;;
esac

bpf_jit="$(cat /proc/sys/net/core/bpf_jit_enable 2>/dev/null || true)"
if [ "$bpf_jit" = "0" ]; then
	warn "net.core.bpf_jit_enable=0: JIT disabled — eBPF programs may fail to load; set it to 1"
else
	pass "BPF JIT enabled (net.core.bpf_jit_enable=${bpf_jit:-n/a})"
fi

if [ "$(id -u)" -ne 0 ]; then
	warn "not running as root: loading eBPF programs needs CAP_SYS_ADMIN+CAP_BPF (use sudo)"
fi

section "Toolchain (needed to build)"
if command -v go >/dev/null 2>&1; then
	go_ver="$(go version 2>/dev/null | sed -E 's/^.*go([0-9]+\.[0-9]+\.[0-9]+).*$/\1/')"
	go_require="$(sed -n 's/^go[[:space:]]*//p' go.mod 2>/dev/null | head -1)"
	if [ -z "$go_require" ] || version_at_least "$go_require" "$go_ver"; then
		pass "Go $go_ver (>= $go_require from go.mod)"
	else
		fail "Go $go_ver installed but go.mod requires $go_require — install Go $go_require from go.dev"
	fi
else
	fail "Go not found — needed for any build (install Go 1.26+ from go.dev)"
fi

if command -v gcc >/dev/null 2>&1; then
	pass "gcc found (CGO static build + exploit tests)"
else
	fail "gcc not found — 'make build' and 'make test-integration' will fail (install gcc + libc6-dev)"
fi

if command -v make >/dev/null 2>&1; then
	pass "make found"
else
	fail "make not found — the Makefile targets will not run (install make)"
fi

section "Docker (needed for 'make test-integration')"
if command -v docker >/dev/null 2>&1; then
	if docker info >/dev/null 2>&1; then
		pass "docker CLI + daemon reachable"
	else
		warn "docker found but daemon not reachable (is it running? is your user in the docker group?)"
	fi
else
	warn "docker not found — 'make test-integration' requires it (apt install docker.io)"
fi

section "Optional tools"
for tool in clang bpftool git; do
	if command -v "$tool" >/dev/null 2>&1; then
		pass "$tool found"
	else
		case "$tool" in
		clang) warn "clang not found — only needed to regenerate BPF bindings (make generate)" ;;
		bpftool) warn "bpftool not found — only needed to regenerate vmlinux.h (make bpftool-headers)" ;;
		git) warn "git not found — only needed for development workflows" ;;
		esac
	fi
done

section "LLVM suite (needed with clang to regenerate BPF bindings)"
llvm_tools="llvm-objdump llvm-readelf llvm-strip llvm-link llvm-nm llvm-ar llvm-size ld.lld"
for tool in $llvm_tools; do
	if command -v "$tool" >/dev/null 2>&1; then
		pass "$tool found"
	else
		warn "$tool not found — only needed for 'make generate' (Ubuntu: apt install llvm lld; Arch: pacman -S llvm lld)"
	fi
done

section "Summary"
if [ "$FAIL" -gt 0 ]; then
	printf "\n  ${RED}%d problem(s) found, %d warning(s)${NC}.\n" "$FAIL" "$WARN"
	printf "  Fix the FAIL items above, reboot when asked, then re-run this check.\n"
	exit 1
fi
if [ "$WARN" -gt 0 ]; then
	printf "\n  ${GREEN}All hard checks passed, %d warning(s)${NC} — read them above (none block, but they matter).\n" "$WARN"
else
	printf "\n  ${GREEN}All checks passed — your host is compatible.${NC}\n"
fi
exit 0