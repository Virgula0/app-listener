package install

import (
	"strings"
	"testing"
)

// TestUnifiedDiff verifies the diff between two configs contains both
// sides with unified markers and context.
func TestUnifiedDiff(t *testing.T) {
	existing := "line one\nold line\nshared\n"
	desired := "line one\nnew line\nshared\n"

	diff := UnifiedDiff(existing, desired)
	for _, want := range []string{"--- existing", "+++ new", "@@", "-old line", "+new line"} {
		if !strings.Contains(diff, want) {
			t.Errorf("diff missing %q:\n%s", want, diff)
		}
	}
	if strings.Contains(diff, "+shared") || strings.Contains(diff, "-shared") {
		t.Errorf("unchanged lines must not appear as changes:\n%s", diff)
	}
}

// TestUnifiedDiffIdentical verifies identical inputs produce no hunks.
func TestUnifiedDiffIdentical(t *testing.T) {
	content := "a\nb\nc\n"
	diff := UnifiedDiff(content, content)
	if strings.Contains(diff, "@@") {
		t.Errorf("identical content produced hunks:\n%s", diff)
	}
}
