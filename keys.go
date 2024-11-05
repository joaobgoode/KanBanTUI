package main

import "github.com/charmbracelet/bubbles/key"

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

type keyMap struct {
	New        key.Binding
	Edit       key.Binding
	Delete     key.Binding
	Up         key.Binding
	Down       key.Binding
	Right      key.Binding
	Left       key.Binding
	Enter      key.Binding
	Prev       key.Binding
	Help       key.Binding
	Quit       key.Binding
	Back       key.Binding
	Todo       key.Binding
	InProgress key.Binding
	Done       key.Binding
	Cancel     key.Binding
	Projects   key.Binding
	Filtering  key.Binding
}

var keys = keyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("delete", "delete"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/l", "move left"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Prev: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "previous"),
	),
	Todo: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "ToDo"),
	),
	InProgress: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "InProgress"),
	),
	Done: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "Done"),
	),
	Projects: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Projects"),
	),
	Filtering: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
}
