package update

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestChangelogText(t *testing.T) {
	tests := []struct {
		name    string
		release githubRelease
		wantSub string
	}{
		{
			name: "published body is used as-is",
			release: githubRelease{
				TagName: "pre-20260808-aaaaaaa",
				Body:    "## What's changed\n* fix: something",
			},
			wantSub: "## What's changed",
		},
		{
			name: "empty body falls back to the tag",
			release: githubRelease{
				TagName: "pre-20260808-bbbbbbb",
			},
			wantSub: "pre-20260808-bbbbbbb",
		},
		{
			name: "whitespace-only body falls back too",
			release: githubRelease{
				TagName:     "pre-20260808-ccccccc",
				Body:        "  \n\t ",
				PublishedAt: "2026-08-08T10:00:00Z",
			},
			wantSub: "No changelog available",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changelogText(&tt.release)
			if !strings.Contains(got, tt.wantSub) {
				t.Fatalf("changelogText() = %q, want it to contain %q", got, tt.wantSub)
			}
		})
	}
}

func TestNewChangelogModelViewContainsTitleAndNotes(t *testing.T) {
	m := newChangelogModel("Release notes: pre-20260808-aaaaaaa", "fix: something\nfeat: another")
	view := m.View()
	if !strings.Contains(view, "Release notes: pre-20260808-aaaaaaa") {
		t.Fatalf("view does not contain the title: %q", view)
	}
	if !strings.Contains(view, "fix: something") {
		t.Fatalf("view does not contain the release notes: %q", view)
	}
}

func TestChangelogModelQuitKeys(t *testing.T) {
	cases := []struct {
		name string
		msg  tea.KeyMsg
	}{
		{"q", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}},
		{"esc", tea.KeyMsg{Type: tea.KeyEsc}},
		{"enter", tea.KeyMsg{Type: tea.KeyEnter}},
		{"ctrl+q", tea.KeyMsg{Type: tea.KeyCtrlQ}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := newChangelogModel("title", "notes")
			got, cmd := m.Update(tc.msg)
			if _, ok := got.(*changelogModel); !ok {
				t.Fatalf("Update(%s) returned %T, want *changelogModel", tc.name, got)
			}
			if cmd == nil {
				t.Fatalf("Update(%s) returned no command, want tea.Quit", tc.name)
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Fatalf("Update(%s) returned %T, want tea.QuitMsg", tc.name, cmd())
			}
		})
	}
}

func TestChangelogModelOtherKeysDoNotQuit(t *testing.T) {
	m := newChangelogModel("title", "notes")
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if _, ok := got.(*changelogModel); !ok {
		t.Fatalf("Update returned %T, want *changelogModel", got)
	}
}

func TestChangelogModelResize(t *testing.T) {
	m := newChangelogModel("title", "notes")
	got, cmd := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if cmd != nil {
		t.Fatalf("resize returned an unexpected command")
	}
	resized, ok := got.(*changelogModel)
	if !ok {
		t.Fatalf("resize returned %T, want *changelogModel", got)
	}
	if width := resized.textarea.Width(); width <= changelogMinWidth {
		t.Fatalf("textarea width = %d after resize, want wider than %d", width, changelogMinWidth)
	}
}
