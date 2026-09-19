#!/usr/bin/env bash
#
# trace-app-libs.sh — record what the daemon refuses while you use an app, so
# the missing permissions can be added to the catalog / daemon.conf.
#
# It follows the running daemon's journal and captures BOTH denial streams:
#
#   TRUST DENIED  op=LIBLOAD   a whitelisted binary tried to exec-map a
#                              library the daemon does not trust
#                              → fixed by allow_lib / lib_dir
#   DAEMON DENIED op=<OP>      a process was refused an operation on a
#                              guarded resource
#                              → fixed by a whitelist line in the resource's
#                                [watch] section, or by a lib_binary line in
#                                the application's [libraries] block when the
#                                resource is a lib_dir (a read-only runtime
#                                tree)
#
# Each denial is classified against /etc/app-listener/daemon.conf. A denial
# whose path lies OUTSIDE the resource it is reported under is flagged
# separately: that is not a missing permission but a guard-engine false
# positive (a recycled inode number), and must be reported, never
# whitelisted around.
#
# Run it, exercise the app (launch Steam, start a game, …), then press Ctrl-C.
# It writes a de-duplicated report next to this script.
#
# Replay instead of following live — analyse a session that already happened:
#   journalctl -u app-listener-daemon -o cat --since "12:45" > session.txt
#   TRACE_REPLAY=session.txt scripts/trace-app-libs.sh
#
# Requirements: the app-listener daemon must be RUNNING, and you need read
# access to its journal (root, or membership of the systemd-journal group).
set -euo pipefail

UNIT="app-listener-daemon"
CONF="${1:-/etc/app-listener/daemon.conf}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="${SCRIPT_DIR}/app-trace-${STAMP}.log"
RAW="$(mktemp)"
FOLLOW_PID=""

cleanup() {
	[ -n "${FOLLOW_PID}" ] && kill "${FOLLOW_PID}" 2>/dev/null || true
	rm -f "${RAW}" "${RAW}".* 2>/dev/null || true
}
trap cleanup EXIT

# --- scope, read from daemon.conf ------------------------------------------
# unquote_path: strip the directive keyword and surrounding double quotes.
conf_directive_values() { # $1 = directive keyword (lib_dir, lib_binary)
	[ -r "${CONF}" ] || return 0
	sed -n "s/^[[:space:]]*$1[[:space:]:]*\"\\{0,1\\}\\([^\"]*\\)\"\\{0,1\\}[[:space:]]*\$/\\1/p" "${CONF}" | sort -u
}

list_conf_binaries() {
	[ -r "${CONF}" ] || { echo "(cannot read ${CONF})"; return; }
	# Whitelist lines are watch-section body lines that are not comments,
	# blanks, headers or directives. [libraries] blocks hold only
	# directives, so everything inside them is skipped.
	awk '
		/^[[:space:]]*\[libraries/ { inlib = 1; next }
		/^[[:space:]]*\[/          { inlib = 0; next }
		inlib                      { next }
		/^[[:space:]]*#/           { next }
		/^[[:space:]]*$/           { next }
		/^[[:space:]]*need_encryption:/ { next }
		/^[[:space:]]*watch[[:space:]:]/ { next }
		/^[[:space:]]*(allow_lib|lib_dir|lib_binary)[[:space:]:]/ { next }
		/^[[:space:]]*path[[:space:]]*=/ { next }
		{
			line = $0
			sub(/^[[:space:]]+/, "", line)
			if (line ~ /^"/) { sub(/^"/, "", line); sub(/".*$/, "", line) }
			else             { sub(/[[:space:]].*$/, "", line) }
			if (line != "") print line
		}
	' "${CONF}" | sort -u
}

if [ -z "${TRACE_REPLAY:-}" ] && ! systemctl is-active --quiet "${UNIT}" 2>/dev/null; then
	echo "WARNING: ${UNIT} is not active — start it before tracing." >&2
fi
if [ -z "${TRACE_REPLAY:-}" ] && [ -z "$(journalctl -u "${UNIT}" -n 1 -o cat 2>/dev/null)" ]; then
	echo "WARNING: cannot read the ${UNIT} journal — run with sudo or join the systemd-journal group." >&2
fi

LIBDIRS="${RAW}.libdirs"
conf_directive_values lib_dir > "${LIBDIRS}"

echo "app-listener denial tracer"
echo "  config : ${CONF}"
echo "  output : ${OUT}"
echo
echo "Whitelisted binaries (watch sections):"
list_conf_binaries | sed 's/^/  - /'
echo "Library-tree writers (lib_binary):"
conf_directive_values lib_binary | sed 's/^/  - /'
echo "Library directories (lib_dir):"
sed 's/^/  - /' "${LIBDIRS}"
echo
[ -n "${TRACE_REPLAY:-}" ] || { echo "Now exercise the apps (launch Steam, start a game, …). Press Ctrl-C when done."; echo; }

# --- follow both denial streams from now onward ----------------------------
# -o cat: MESSAGE only. Start at 'now' so only this session is captured.
if [ -n "${TRACE_REPLAY:-}" ]; then
	grep -E 'TRUST DENIED|DAEMON DENIED|level=error' "${TRACE_REPLAY}" > "${RAW}" || true
else
( journalctl -u "${UNIT}" -f -o cat --since "$(date '+%Y-%m-%d %H:%M:%S')" 2>/dev/null \
	| grep --line-buffered -E 'TRUST DENIED|DAEMON DENIED|level=error' >> "${RAW}" ) &
FOLLOW_PID=$!

trap 'echo; echo "stopping…"; kill "${FOLLOW_PID}" 2>/dev/null || true' INT
while kill -0 "${FOLLOW_PID}" 2>/dev/null; do
	sleep 2
	if [ -s "${RAW}" ]; then
		libs=$(grep -c 'op=LIBLOAD' "${RAW}" || true)
		denies=$(grep -c 'DAEMON DENIED' "${RAW}" || true)
		printf '\r  captured %s LIBLOAD, %s DAEMON DENIED events…' "${libs}" "${denies}"
	fi
done 2>/dev/null || true
printf '\r'
fi

# --- LIBLOAD: library paths by loader ---------------------------------------
# comm : between "comm=" and "  pid=" (comm may contain spaces).
# path : between "path=" and " — a whitelisted binary tried" (the em dash is
#        one UTF-8 char, matched by "." under a UTF-8 locale).
LIBPAIRS="${RAW}.libs"
paste -d'\t' \
	<({ grep 'op=LIBLOAD' "${RAW}" || true; } | sed -n 's/.* comm=\(.*\)  pid=[0-9].*/\1/p') \
	<({ grep 'op=LIBLOAD' "${RAW}" || true; } | LC_ALL=C.UTF-8 sed -n 's/.*path=\(.*\) . a whitelisted binary tried.*/\1/p') |
	awk -F'\t' 'NF==2' | sort -u > "${LIBPAIRS}"

# --- DAEMON DENIED: op, comm, exe, uid, resource, path as TSV --------------
DENY="${RAW}.deny"
{ grep 'DAEMON DENIED' "${RAW}" || true; } |
	sed -n 's/.*DAEMON DENIED  op=\([A-Z_]*\)  comm=\(.*\)  commFullPath=\(.*\)  pid=[0-9]*  uid=\([^ ]*\)  resource=\(.*\)  path=\(.*\)$/\1\t\2\t\3\t\4\t\5\t\6/p' \
	> "${DENY}"

{
	echo "# app-listener denial trace — ${STAMP}"
	echo "# Config: ${CONF}"
	echo "#"
	echo "# Section 1: libraries a whitelisted binary was refused (allow_lib / lib_dir)."
	echo "# Section 2: operations refused on guarded resources, per resource and"
	echo "#            binary, with the config line that WOULD grant them. Review"
	echo "#            each one: a suggestion is not a recommendation."
	echo "# Section 3: denials reported under a resource their path is not in —"
	echo "#            engine false positives. Report them; never whitelist them."
	echo
	if grep -q 'level=error' "${RAW}"; then
		echo "## 0. DAEMON ERRORS — fix these first (a failed reload keeps the OLD config running)"
		{ grep 'level=error' "${RAW}" || true; } | sed -n 's/.*msg="\(.*\)"$/  \1/p' | sort | uniq -c
		echo
	fi
	echo "## 1. LIBLOAD — unique library paths"
	cut -f2 "${LIBPAIRS}" | sed 's/^$/(anonymous: memfd \/ O_TMPFILE — runtime-generated code)/' | sort -u
	echo
	echo "## 1b. LIBLOAD — grouped by loader (comm)"
	awk -F'\t' '{ libs[$1] = libs[$1] "\n    " ($2 == "" ? "(anonymous)" : $2) } END { n = asorti(libs, keys); for (i = 1; i <= n; i++) print "[" keys[i] "]" libs[keys[i]] }' "${LIBPAIRS}"
	echo
	echo "## 2. DAEMON DENIED — grouped by resource"
	awk -F'\t' -v libdirs="${LIBDIRS}" '
		BEGIN { while ((getline d < libdirs) > 0) if (d != "") isLib[d] = 1 }
		function inside(p, r) { return p == r || index(p, r "/") == 1 }
		{
			op = $1; comm = $2; exe = $3; uid = $4; res = $5; path = $6
			who = (exe == "~" || exe == "") ? "comm:" comm : exe
			if (op == "PTRACE" || op == "TRACED_EXEC" || op == "PROC_MEM") {
				# path: "pid=N comm=NAME[ mode=READ|ATTACH]" — group by target
				# NAME and mode, not pid (short-lived helpers repeat a lot).
				tcomm = path; sub(/^pid=[0-9]+ comm=/, "", tcomm)
				tmode = "-"
				if (match(tcomm, / mode=[A-Z]+$/)) { tmode = substr(tcomm, RSTART + 6); tcomm = substr(tcomm, 1, RSTART - 1) }
				gate[res SUBSEP who SUBSEP op SUBSEP tmode SUBSEP tcomm]++
				next
			}
			if (!inside(path, res)) {
				outside[res SUBSEP who SUBSEP op] = path
				next
			}
			k = res SUBSEP who
			if (!(k in seen)) { seen[k] = 1; order[++n] = k }
			ops[k, op]++
			if (!(k SUBSEP op in sample)) sample[k SUBSEP op] = path
			oplist[k] = (index(" " oplist[k] " ", " " op " ") ? oplist[k] : oplist[k] " " op)
			users[k] = uid
		}
		END {
			for (i = 1; i <= n; i++) {
				split(order[i], kv, SUBSEP); res = kv[1]; who = kv[2]
				kind = (res in isLib) ? "lib_dir" : "watch"
				if (res != last) { printf "\n[%s] %s\n", kind, res; last = res }
				printf "  %s  (uid %s)\n", who, users[order[i]]
				m = split(substr(oplist[order[i]], 2), o, " ")
				evs = ""
				for (j = 1; j <= m; j++) {
					printf "      %-9s x%-4d e.g. %s\n", o[j], ops[order[i], o[j]], sample[order[i] SUBSEP o[j]]
					evs = evs (evs == "" ? "" : ",") o[j]
				}
				if (who ~ /^comm:/)
					printf "    → exe unresolved (the process exited first): identify \"%s\" before granting anything\n", substr(who, 6)
				else if (kind == "lib_dir")
					printf "    → would be granted by, in the application [libraries] block:  lib_binary \"%s\"\n", who
				else
					printf "    → would be granted by, in this [watch] section:  \"%s\" %s\n", who, evs
			}
			ng = 0
			for (k in gate) { gk[++ng] = k }
			# sort by resource, then by count descending
			for (i = 2; i <= ng; i++) {
				v = gk[i]; split(v, vk, SUBSEP); j = i - 1
				while (j > 0) {
					split(gk[j], jk, SUBSEP)
					if (jk[1] < vk[1] || (jk[1] == vk[1] && gate[gk[j]] >= gate[v])) break
					gk[j + 1] = gk[j]; j--
				}
				gk[j + 1] = v
			}
			if (ng) {
				print "\n## 2b. PROCESS GATES — a process was refused access to a tainted one (no file involved)"
				print "  mode=ATTACH: memory (ptrace, process_vm_readv, /proc/<pid>/mem)"
				print "  mode=READ:   metadata only (/proc/<pid>/environ, fd, maps, ...)"
				print "  TRACED_EXEC: a whitelisted binary exec'"'"'d while traced by a non-whitelisted tracer"
				lastres = ""
				for (i = 1; i <= ng; i++) {
					split(gk[i], kv, SUBSEP)
					if (kv[1] != lastres) { printf "\n  [%s]\n", kv[1]; lastres = kv[1] }
					printf "    x%-5d %-11s %-6s caller=%s  target=%s\n", gate[gk[i]], kv[3], kv[4], kv[2], kv[5]
				}
				print "  → each target read this resource'"'"'s secrets (or inherited them); the caller is not whitelisted for it."
			}
			any = 0
			for (k in outside) {
				if (!any) { print "\n## 3. ENGINE FALSE POSITIVES — path is outside the reported resource"; any = 1 }
				split(k, kv, SUBSEP)
				printf "  [%s] %s  %s  path=%s\n", kv[3], kv[1], kv[2], outside[k]
			}
			if (any) print "  → a recycled inode number matched this resource. Not a missing permission: report it."
		}
	' "${DENY}"
} > "${OUT}"

echo "done — $(wc -l < "${LIBPAIRS}") library load(s), $(wc -l < "${DENY}") access denial(s) captured"
echo "wrote ${OUT}"
