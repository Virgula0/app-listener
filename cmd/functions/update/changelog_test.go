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

func rel(tag string, pre bool, at, body string) githubRelease {
	return githubRelease{TagName: tag, Prerelease: pre, PublishedAt: at, Body: body}
}

// Server order is not publish order (the releases API sorts by creation).
var cumulativeFixture = []githubRelease{
	rel("pre-20261013-ddddddd", true, "2026-10-13T12:00:00Z", "notes d"),
	rel("pre-20261009-bbbbbbb", true, "2026-10-09T09:00:00Z", "notes b"),
	rel("pre-20261009-ccccccc", true, "2026-10-09T16:00:00Z", "notes c"),
	rel("v0.2.0", false, "2026-10-10T00:00:00Z", "stable 0.2.0"),
	rel("v0.1.1", false, "2026-09-20T00:00:00Z", "stable 0.1.1"),
	rel("pre-20261001-aaaaaaa", true, "2026-10-01T00:00:00Z", "notes a"),
	rel("v0.1.0", false, "2026-09-08T00:00:00Z", "stable 0.1.0"),
}

func tags(rs []githubRelease) []string {
	var out []string
	for i := range rs {
		out = append(out, rs[i].TagName)
	}
	return out
}

func TestReleasesSince(t *testing.T) {
	latestPre := cumulativeFixture[0]
	latestStable := cumulativeFixture[3]
	tests := []struct {
		name      string
		installed string
		latest    githubRelease
		channel   string
		want      []string
		wantFound bool
	}{
		{"pre-release span, newest first", "pre-20261001-aaaaaaa", latestPre, channelPreRelease,
			[]string{"pre-20261013-ddddddd", "pre-20261009-ccccccc", "pre-20261009-bbbbbbb"}, true},
		{"same-day ordering by time", "pre-20261009-bbbbbbb", latestPre, channelPreRelease,
			[]string{"pre-20261013-ddddddd", "pre-20261009-ccccccc"}, true},
		{"stable by semver", "v0.1.0", latestStable, channelStable,
			[]string{"v0.2.0", "v0.1.1"}, true},
		{"stable installed off the fetched page", "v0.0.9", latestStable, channelStable,
			[]string{"v0.2.0", "v0.1.1", "v0.1.0"}, true},
		{"stable to pre-release channel switch", "v0.1.1", latestPre, channelPreRelease,
			[]string{"pre-20261013-ddddddd", "pre-20261009-ccccccc", "pre-20261009-bbbbbbb", "pre-20261001-aaaaaaa"}, true},
		{"unknown installed: everything up to latest", "dev", latestStable, channelStable,
			[]string{"v0.2.0", "v0.1.1", "v0.1.0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			span, found := releasesSince(tt.installed, &tt.latest, cumulativeFixture, tt.channel)
			if got := tags(span); strings.Join(got, ",") != strings.Join(tt.want, ",") || found != tt.wantFound {
				t.Fatalf("releasesSince = %v (found %v), want %v (found %v)", got, found, tt.want, tt.wantFound)
			}
		})
	}
}

func TestCumulativeChangelog(t *testing.T) {
	latest := cumulativeFixture[0]
	got := cumulativeChangelog("pre-20261009-bbbbbbb", &latest, cumulativeFixture, channelPreRelease)
	for _, want := range []string{
		"Changes from pre-20261009-bbbbbbb to pre-20261013-ddddddd (2 releases)",
		"compare/pre-20261009-bbbbbbb...pre-20261013-ddddddd",
		"# pre-20261013-ddddddd", "notes d", "# pre-20261009-ccccccc", "notes c",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "notes b") || strings.Index(got, "notes d") > strings.Index(got, "notes c") {
		t.Fatalf("wrong span or order:\n%s", got)
	}

	one := cumulativeChangelog("pre-20261009-ccccccc", &latest, cumulativeFixture, channelPreRelease)
	if one != "notes d" {
		t.Fatalf("single-step update should show the release notes as-is, got:\n%s", one)
	}

	unknown := cumulativeChangelog("dev", &latest, cumulativeFixture, channelPreRelease)
	if !strings.Contains(unknown, "not among the fetched releases") || strings.Contains(unknown, "compare/") {
		t.Fatalf("unknown installed version not flagged:\n%s", unknown)
	}
}
