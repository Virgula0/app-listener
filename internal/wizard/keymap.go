package wizard

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// MultiSelectKeymap is the shared key map for the wizards' preselected
// multi-selects: Ctrl+K toggles every entry, space/x toggles one, and the
// field legend shows both. Apply it to the Form, not the field —
// huh.NewForm overwrites each field's keymap with the form default.
func MultiSelectKeymap() *huh.KeyMap {
	keys := huh.NewDefaultKeyMap()
	keys.MultiSelect.SelectAll = key.NewBinding(
		key.WithKeys("ctrl+k"), key.WithHelp("ctrl+k", "select/deselect all"))
	keys.MultiSelect.SelectNone = key.NewBinding(
		key.WithKeys("ctrl+k"), key.WithHelp("ctrl+k", "select/deselect all"))
	keys.MultiSelect.Toggle = key.NewBinding(
		key.WithKeys(" ", "x"), key.WithHelp("space/x", "select one"))
	return keys
}
