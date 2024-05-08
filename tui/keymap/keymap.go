package keymap

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Reload key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		Reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload"),
		),
	}
}
