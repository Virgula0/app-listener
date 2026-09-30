package wizard

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// MultiSelectKeymap is the shared keymap for preselected multi-selects: Ctrl+K toggles all, space/x
// one, and the legend shows both. Apply it to the Form, not the field (huh.NewForm overwrites each
// field's keymap).
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
