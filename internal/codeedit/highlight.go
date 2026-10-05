package codeedit

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"
)

// maxHighlightBytes bounds what is re-tokenized on every edit; larger buffers render uncolored.
const maxHighlightBytes = 512 << 10

// Highlighter renders a codeedit textarea: line numbers in a fixed-width gutter, text colored by
// the language of the file it holds (none: uncolored). A nil *Highlighter still draws the gutter.
type Highlighter struct {
	lexer chroma.Lexer // nil: no language

	// Token cache: text is the buffer classes was computed from; classes[line][rune] indexes
	// tokenStyles (0 = unstyled).
	text    string
	classes [][]uint8

	// Wrap cache, keyed on the same text plus the textarea's size; labels[displayRow] is the
	// gutter's line number ("" on a soft-wrapped continuation).
	wrapWidth, vpHeight int
	wraps               [][][]rune
	labels              []string

	// Gutter: the textarea's Prompt, then a number column digits wide.
	bar    string
	digits int

	// Find marks (Finder.View): underlined spans per line; current is also reversed.
	marks      map[int][]span
	current    span
	hasCurrent bool
}

// NewHighlighter picks the language from name (base name or extension; daemon.conf has its own
// lexer), then from a `#!` interpreter line in content.
func NewHighlighter(name, content string) *Highlighter {
	if l := lexerFor(name, content); l != nil {
		return &Highlighter{lexer: chroma.Coalesce(l)}
	}
	return &Highlighter{}
}

// Language is the detected language's name, "" when the text renders plain.
func (h *Highlighter) Language() string {
	if h == nil || h.lexer == nil {
		return ""
	}
	return h.lexer.Config().Name
}

func lexerFor(name, content string) chroma.Lexer {
	if filepath.Base(name) == "daemon.conf" {
		return daemonConfLexer
	}
	if l := lexers.Match(filepath.Base(name)); l != nil && l != lexers.Fallback {
		if l.Config().Name != "plaintext" {
			return l
		}
		return nil
	}
	if l := shebangLexer(content); l != nil {
		return l
	}
	return nil
}

// shebangLexer resolves `#!/usr/bin/env python3` / `#!/bin/bash` to the interpreter's lexer.
func shebangLexer(content string) chroma.Lexer {
	first, _, _ := strings.Cut(content, "\n")
	if !strings.HasPrefix(first, "#!") {
		return nil
	}
	fields := strings.Fields(first[2:])
	if len(fields) == 0 {
		return nil
	}
	interp := filepath.Base(fields[0])
	if interp == "env" {
		interp = ""
		for _, f := range fields[1:] {
			if !strings.HasPrefix(f, "-") && !strings.Contains(f, "=") {
				interp = filepath.Base(f)
				break
			}
		}
	}
	if interp == "" {
		return nil
	}
	if l := lexers.Get(interp); l != nil {
		return l
	}
	return lexers.Get(strings.TrimRight(interp, "0123456789."))
}

// tokenize refreshes the per-rune classes for text; no lexer, an oversized buffer or a lexer
// error leaves every rune unstyled.
func (h *Highlighter) tokenize(text string) {
	if text == h.text && h.classes != nil {
		return
	}
	h.text, h.wraps = text, nil
	nLines := strings.Count(text, "\n") + 1
	if h.lexer == nil || len(text) > maxHighlightBytes {
		h.classes = make([][]uint8, nLines)
		return
	}
	h.classes = make([][]uint8, 1, nLines)
	it, err := h.lexer.Tokenise(nil, text)
	if err != nil {
		h.classes = make([][]uint8, nLines)
		return
	}
	for tok := it(); tok != chroma.EOF; tok = it() {
		c := tokenClassOf(tok.Type)
		for _, r := range tok.Value {
			if r == '\n' {
				h.classes = append(h.classes, nil)
				continue
			}
			last := len(h.classes) - 1
			h.classes[last] = append(h.classes[last], c)
		}
	}
	// Lexers may append a trailing newline (EnsureNL): keep the buffer's own line count.
	for len(h.classes) < nLines {
		h.classes = append(h.classes, nil)
	}
	h.classes = h.classes[:nLines]
}

// tokenStyles is the palette, tuned for both light and dark terminals; index 0 is "no style".
var tokenStyles = []lipgloss.Style{
	{},
	fg("244", "245").Italic(true), // comment
	fg("162", "204"),              // keyword
	fg("31", "81"),                // type / builtin / namespace
	fg("97", "141"),               // constant / number / bool
	fg("28", "114"),               // function / class / inserted
	fg("130", "186"),              // string
	fg("25", "117"),               // tag / attribute / key
	fg("160", "203"),              // operator / deleted
	fg("133", "176"),              // preprocessor / decorator
	fg("162", "212").Bold(true),   // heading
	lipgloss.NewStyle().Bold(true),
	lipgloss.NewStyle().Italic(true),
	fg("166", "216"), // variable
}

const (
	clsNone uint8 = iota
	clsComment
	clsKeyword
	clsType
	clsConstant
	clsFunction
	clsString
	clsTag
	clsOperator
	clsPreproc
	clsHeading
	clsStrong
	clsEmph
	clsVariable
)

func fg(light, dark string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: light, Dark: dark})
}

// tokenClasses maps token types to palette entries; tokenClassOf falls back from the exact type to its
// sub-category, then its category.
var tokenClasses = map[chroma.TokenType]uint8{
	chroma.Comment:                clsComment,
	chroma.CommentPreproc:         clsPreproc,
	chroma.CommentPreprocFile:     clsString,
	chroma.Keyword:                clsKeyword,
	chroma.KeywordType:            clsType,
	chroma.KeywordConstant:        clsConstant,
	chroma.Name:                   clsNone,
	chroma.NameBuiltin:            clsType,
	chroma.NameBuiltinPseudo:      clsConstant,
	chroma.NameClass:              clsFunction,
	chroma.NameFunction:           clsFunction,
	chroma.NameFunctionMagic:      clsFunction,
	chroma.NameConstant:           clsConstant,
	chroma.NameNamespace:          clsType,
	chroma.NameDecorator:          clsPreproc,
	chroma.NameTag:                clsTag,
	chroma.NameAttribute:          clsTag,
	chroma.NameLabel:              clsTag,
	chroma.NameProperty:           clsTag,
	chroma.NameVariable:           clsVariable,
	chroma.NameEntity:             clsConstant,
	chroma.LiteralString:          clsString,
	chroma.LiteralStringEscape:    clsConstant,
	chroma.LiteralStringInterpol:  clsVariable,
	chroma.LiteralStringSymbol:    clsConstant,
	chroma.LiteralNumber:          clsConstant,
	chroma.LiteralDate:            clsConstant,
	chroma.Literal:                clsString,
	chroma.Operator:               clsOperator,
	chroma.OperatorWord:           clsKeyword,
	chroma.GenericInserted:        clsFunction,
	chroma.GenericDeleted:         clsOperator,
	chroma.GenericHeading:         clsHeading,
	chroma.GenericSubheading:      clsType,
	chroma.GenericStrong:          clsStrong,
	chroma.GenericEmph:            clsEmph,
	chroma.GenericPrompt:          clsComment,
	chroma.GenericUnderline:       clsEmph,
	chroma.LiteralStringHeredoc:   clsString,
	chroma.LiteralStringDoc:       clsComment,
	chroma.LiteralStringBacktick:  clsString,
	chroma.LiteralStringAffix:     clsKeyword,
	chroma.LiteralStringDelimiter: clsString,
}

func tokenClassOf(t chroma.TokenType) uint8 {
	for _, k := range []chroma.TokenType{t, t.SubCategory(), t.Category()} {
		if c, ok := tokenClasses[k]; ok {
			return c
		}
	}
	return clsNone
}
