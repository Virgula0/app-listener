#!/bin/sh
# UserPromptSubmit hook: rebuild graphify-out/graph.json when uncommitted code differs from
# what it was built from, so this prompt's graphify queries see the working tree. The git
# hooks only cover commit/checkout. AST-only (no LLM); ~3s when stale, ~50ms otherwise.
# stdout is injected as prompt context; always exits 0 so a prompt is never blocked.
cd "${CLAUDE_PROJECT_DIR:-.}" 2>/dev/null || exit 0
command -v graphify >/dev/null 2>&1 || exit 0
[ -f graphify-out/graph.json ] || exit 0

stamp=graphify-out/.prompt_stamp
key=$(
	{
		git rev-parse HEAD
		git diff HEAD -- '*.go' '*.c' '*.h' '*.sh'
		git ls-files -z --others --exclude-standard -- '*.go' '*.c' '*.h' '*.sh' |
			xargs -0 -r stat -c '%n %Y'
	} 2>/dev/null | sha256sum | cut -d' ' -f1
)
[ "$(cat "$stamp" 2>/dev/null)" = "$key" ] && exit 0

if out=$(graphify update . 2>&1); then
	printf '%s\n' "$key" >"$stamp"
	echo "graphify: graph.json rebuilt for the current working tree."
else
	echo "graphify: rebuild FAILED, graph may be stale for edited files: $(printf '%s' "$out" | tail -n 1)"
fi
exit 0
