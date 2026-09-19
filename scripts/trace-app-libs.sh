#!/usr/bin/env bash
#
# trace-app-libs.sh — capture the dynamic libraries the daemon's whitelisted
# binaries load, to populate `allow_lib` for the trusted-library allowlist.
#
# It reads /etc/app-listener/daemon.conf to show which binaries are in scope,
# then follows the running daemon's Phase-1 observe output (the `LIBLOAD`
# lines the trust guard emits for every executable library a whitelisted
# binary maps that is NOT already in its static closure). Those events are
# already filtered to whitelisted binaries by the kernel (exe-inode), so this
# script simply aggregates them.
#
# Run it, then exercise every app you want to make strict (open Discord, VS
# Code, Claude, the browsers, …). Press Ctrl-C to stop: it writes a
# de-duplicated log next to this script that you can hand back for populating
# the catalog.
#
# Requirements: the app-listener daemon must be RUNNING and built with the
# trust observer (it logs `trust guard: observing library loads …` at start).
# No extra tooling — just journalctl.
set -euo pipefail

UNIT="app-listener-daemon"
CONF="${1:-/etc/app-listener/daemon.conf}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="${SCRIPT_DIR}/lib-trace-${STAMP}.log"
RAW="$(mktemp)"
FOLLOW_PID=""

cleanup() {
	[ -n "${FOLLOW_PID}" ] && kill "${FOLLOW_PID}" 2>/dev/null || true
	rm -f "${RAW}" 2>/dev/null || true
}
trap cleanup EXIT

# --- scope: the whitelisted binaries in daemon.conf ------------------------
list_conf_binaries() {
	[ -r "${CONF}" ] || { echo "(cannot read ${CONF})"; return; }
	# Binary lines are the section-body lines that are not comments, blank,
	# section headers, or the need_encryption/watch/allow_lib/path directives.
	awk '
		/^[[:space:]]*#/       { next }
		/^[[:space:]]*$/        { next }
		/^[[:space:]]*\[watch/  { next }
		/^[[:space:]]*need_encryption:/ { next }
		/^[[:space:]]*watch:/   { next }
		/^[[:space:]]*watch[[:space:]]/ { next }
		/^[[:space:]]*allow_lib/{ next }
		/^[[:space:]]*lib_dir/  { next }
		/^[[:space:]]*path[[:space:]]*=/ { next }
		{
			line=$0
			sub(/^[[:space:]]+/, "", line)
			# strip a trailing " EV1,EV2" event list and surrounding quotes
			sub(/[[:space:]].*$/, "", line)
			gsub(/"/, "", line)
			if (line != "") print line
		}
	' "${CONF}" | sort -u
}

if ! systemctl is-active --quiet "${UNIT}" 2>/dev/null; then
	echo "WARNING: ${UNIT} is not active — start it (with the observe build) before tracing." >&2
fi

echo "app-listener library tracer"
echo "  config : ${CONF}"
echo "  output : ${OUT}"
echo
echo "Whitelisted binaries in scope:"
list_conf_binaries | sed 's/^/  - /'
echo
echo "Now exercise the apps you want to make strict. Press Ctrl-C when done."
echo

# --- follow the daemon's LIBLOAD stream from now onward ---------------------
# -o cat: MESSAGE only (no journal prefix). Start at 'now' so only this
# session's loads are captured.
( journalctl -u "${UNIT}" -f -o cat --since "$(date '+%Y-%m-%d %H:%M:%S')" 2>/dev/null \
	| grep --line-buffered 'LIBLOAD' >> "${RAW}" ) &
FOLLOW_PID=$!

# Live counter until interrupted. Ctrl-C stops the follower so the loop below
# exits and the report is written.
count=0
trap 'echo; echo "stopping…"; kill "${FOLLOW_PID}" 2>/dev/null || true' INT
while kill -0 "${FOLLOW_PID}" 2>/dev/null; do
	sleep 2
	if [ -s "${RAW}" ]; then
		count=$(wc -l < "${RAW}")
		printf '\r  captured %s LIBLOAD events…' "${count}"
	fi
done 2>/dev/null || true
# The INT trap above interrupts the sleep loop; fall through to writing output.
printf '\r'

# --- de-duplicate and write the report -------------------------------------
# comm  : text between "comm=" and "  pid=" (comm may contain spaces).
# path  : text between "path=" and " — a whitelisted binary tried" (the em
#         dash is one UTF-8 char, matched by "." under a UTF-8 locale).
comm_of() { sed -n 's/.* comm=\(.*\)  pid=[0-9].*/\1/p' "${RAW}"; }
path_of() { LC_ALL=C.UTF-8 sed -n 's/.*path=\(.*\) . a whitelisted binary tried.*/\1/p' "${RAW}"; }

PAIRS="$(mktemp)"
paste -d'\t' <(comm_of) <(path_of) | awk -F'\t' 'NF==2 && $2!=""' | sort -u > "${PAIRS}"

{
	echo "# app-listener library trace — ${STAMP}"
	echo "# Config: ${CONF}"
	echo "# De-duplicated dynamic libraries loaded by whitelisted binaries."
	echo "#"
	echo "# Section 1: every unique library path (drop these into allow_lib)."
	echo "# Section 2: the same, grouped by the loading process (comm) so you"
	echo "#            can attribute each to the right catalog entry."
	echo
	echo "## unique library paths"
	cut -f2 "${PAIRS}" | sort -u
	echo
	echo "## grouped by loader (comm)"
	awk -F'\t' '
		{ libs[$1] = libs[$1] "\n    " $2 }
		END { for (c in libs) print "[" c "]" libs[c] }
	' "${PAIRS}" | sort
} > "${OUT}"

rm -f "${PAIRS}"

echo "done — $(cut -f2 "${OUT}" 2>/dev/null | grep -c '^' || true) lines"
echo "wrote ${OUT}"
echo "Pass that file back to populate allow_lib."
