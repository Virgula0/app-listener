#!/usr/bin/env bash
#
# kernel-prof.sh: measure the in-kernel cost of a running app-listener daemon,
# the part `make pprof` cannot see.
#
#   1. per-program BPF run time (kernel.bpf_stats_enabled), which is charged to
#      EVERY process that hits a guarded LSM hook, not to the daemon;
#   2. kernel memory held by the BPF maps (bytes_memlock, not in the daemon's RSS);
#   3. a `perf record` of the daemon with user + kernel stacks (its sys time);
#   4. the daemon's utime/stime over the window, RSS, threads, open fds.
#
# Usage (keep using the machine normally while it samples):
#   sudo bash scripts/kernel-prof.sh [seconds]      default 120
#
# Output: build/pprof/kernel-<timestamp>/ (bpf-hooks.txt, bpf-maps.txt,
# perf.flat.txt, perf.callers.txt, proc.txt, raw progs.json + perf.data).
# Pair it with `make pprof` (Go-side CPU/heap) for the same window.
#
set -euo pipefail

SECS=${1:-120}
[[ $SECS =~ ^[0-9]+$ ]] || { echo "usage: $0 [seconds]" >&2; exit 2; }
[[ $EUID -eq 0 ]] || { echo "ERROR: needs root (bpftool, perf, bpf_stats sysctl)" >&2; exit 1; }
for t in bpftool perf jq column; do
	command -v "$t" >/dev/null || { echo "ERROR: $t not found" >&2; exit 1; }
done

REPO=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
OUT=$REPO/build/pprof/kernel-$(date +%Y%m%d-%H%M%S)

PID=$(systemctl show -p MainPID --value app-listener-daemon.service 2>/dev/null || true)
if [[ -z $PID || $PID == 0 ]]; then
	PID=$(pgrep -f 'app-listener daemon' | head -1 || true)
fi
[[ -n $PID && -d /proc/$PID ]] || { echo "ERROR: no running app-listener daemon" >&2; exit 1; }

mkdir -p "$OUT"
echo ">> daemon pid $PID, sampling ${SECS}s into $OUT"

# bpf_stats adds a little per-call overhead; always restore the previous value.
prev=$(cat /proc/sys/kernel/bpf_stats_enabled)
trap 'sysctl -qw kernel.bpf_stats_enabled="$prev"' EXIT
sysctl -qw kernel.bpf_stats_enabled=1

read -r u0 s0 < <(awk '{print $14, $15}' /proc/"$PID"/stat)
perf record -q -F 499 -g -p "$PID" -o "$OUT/perf.data" -- sleep "$SECS" 2>/dev/null
read -r u1 s1 < <(awk '{print $14, $15}' /proc/"$PID"/stat)
bpftool -j prog show > "$OUT/progs.json"
sysctl -qw kernel.bpf_stats_enabled="$prev"

# run_time_ns counts from enable time, i.e. this window only.
jq -r --argjson secs "$SECS" '
  [.[] | select((.run_cnt // 0) > 0)] | sort_by(-.run_time_ns) | .[]
  | [.id, .type, .name, .run_cnt, .run_time_ns, (.run_time_ns / .run_cnt | floor),
     (.run_time_ns / ($secs * 1e9) * 100 * 1000 | floor / 1000)] | @tsv' "$OUT/progs.json" \
	| (echo -e "id\ttype\tname\tcalls\ttotal_ns\tns/call\t%one_cpu"; cat) | column -t > "$OUT/bpf-hooks.txt"

bpftool -j map show | jq -r '
  sort_by(-(.bytes_memlock // 0)) | .[]
  | [.id, .type, .name, .max_entries, .bytes_memlock] | @tsv' \
	| (echo -e "id\ttype\tname\tmax_entries\tbytes_memlock"; cat) | column -t > "$OUT/bpf-maps.txt"

# Reports are rendered here, as root: the installed binary is root-only, so a later
# `perf report` as a user cannot resolve the Go symbols.
perf report -i "$OUT/perf.data" --stdio --no-children --sort dso,sym -g none --percent-limit 0.5 \
	2>/dev/null > "$OUT/perf.flat.txt" || true
perf report -i "$OUT/perf.data" --stdio --children --sort sym -g none --percent-limit 1 \
	2>/dev/null > "$OUT/perf.children.txt" || true
perf report -i "$OUT/perf.data" --stdio --no-children --sort sym -g caller,0.5,callee --percent-limit 2 \
	2>/dev/null | head -600 > "$OUT/perf.callers.txt" || true

{
	echo "window_s=$SECS utime_ticks=$((u1 - u0)) stime_ticks=$((s1 - s0)) clk_tck=$(getconf CLK_TCK)"
	awk -v u=$((u1 - u0)) -v s=$((s1 - s0)) -v hz="$(getconf CLK_TCK)" -v w="$SECS" \
		'BEGIN { printf "daemon_cpu_pct_one_core=%.2f\n", (u + s) / hz / w * 100 }'
	grep -E 'VmRSS|VmHWM|RssAnon|RssFile|Threads' /proc/"$PID"/status
	echo "open_fds=$(find /proc/"$PID"/fd -mindepth 1 -maxdepth 1 | wc -l)"
	awk 'NR > 1 { s += $5 } END { printf "bpf_maps_memlock_total_bytes=%d\n", s }' "$OUT/bpf-maps.txt"
} > "$OUT/proc.txt"

[[ -n ${SUDO_USER:-} ]] && chown -R "$SUDO_USER": "$OUT"
chmod -R a+rX "$OUT"
echo ">> done: $OUT"
cat "$OUT/proc.txt"
head -20 "$OUT/bpf-hooks.txt"
