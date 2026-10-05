package codeedit

import (
	"github.com/alecthomas/chroma/v2"
)

// daemonConfLexer tokenizes daemon.conf (internal/daemonconfig): `# comments`, `[watch …]` /
// `[libraries "…"]` / `[inspectors]` / `[electron_apps]` headers, `key:` / `allow_lib`-style
// directives, binary paths and their event lists (READ,WRITE,…). Paths are consumed whole first so
// an upper-case path component is not colored as an event.
var daemonConfLexer = chroma.MustNewLexer(
	&chroma.Config{Name: "daemon.conf", Filenames: []string{"daemon.conf"}},
	func() chroma.Rules {
		return chroma.Rules{
			"root": {
				{Pattern: `^[^\S\n]*#.*`, Type: chroma.Comment},
				{
					Pattern: `^([^\S\n]*)(\[)(watch|libraries|inspectors|electron_apps)\b([^\]\n]*)(\])`,
					Type: chroma.ByGroups(chroma.Text, chroma.Punctuation, chroma.Keyword,
						chroma.LiteralString, chroma.Punctuation),
				},
				{Pattern: `^[^\S\n]*\[[^\]\n]*\]`, Type: chroma.Keyword},
				{
					Pattern: `^([^\S\n]*)([a-z_]+)(:?)(?=[^\S\n]|$)`,
					Type:    chroma.ByGroups(chroma.Text, chroma.NameAttribute, chroma.Punctuation),
				},
				{Pattern: `"[^"\n]*"?`, Type: chroma.LiteralString},
				{Pattern: `[/~][^\s,"]*`, Type: chroma.Name},
				{Pattern: `\b(true|false)\b`, Type: chroma.KeywordConstant},
				{Pattern: `\b[A-Z][A-Z_]*\b`, Type: chroma.NameConstant},
				{Pattern: `,`, Type: chroma.Punctuation},
				{Pattern: `\s+`, Type: chroma.Text},
				{Pattern: `[^\s,"/~]+`, Type: chroma.Text},
			},
		}
	},
)
