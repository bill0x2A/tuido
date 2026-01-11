package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Today     key.Binding
	GoTo      key.Binding
	Add       key.Binding
	Edit      key.Binding
	Toggle    key.Binding
	Delete    key.Binding
	MoveUp    key.Binding
	MoveDown  key.Binding
	MoveNext  key.Binding // Move task to next day
	Help      key.Binding
	Quit      key.Binding
	Enter     key.Binding
	Escape    key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev day"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next day"),
	),
	Today: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "today"),
	),
	GoTo: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "go to date"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("enter", " ", "x"),
		key.WithHelp("⏎/x", "toggle"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	MoveUp: key.NewBinding(
		key.WithKeys("K", "shift+up"),
		key.WithHelp("K/⇧↑", "move up"),
	),
	MoveDown: key.NewBinding(
		key.WithKeys("J", "shift+down"),
		key.WithHelp("J/⇧↓", "move down"),
	),
	MoveNext: key.NewBinding(
		key.WithKeys(">", "n"),
		key.WithHelp(">/n", "→ tomorrow"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
	),
}

// Number keys for quick jump
// 1-9 = tasks 1-9
// 0 = task 10
// !@#$%^&*( (shift+1-9-0) = tasks 11-20
func isNumberKey(s string) (int, bool) {
	if len(s) != 1 {
		return 0, false
	}

	// 1-9 for tasks 1-9
	if s[0] >= '1' && s[0] <= '9' {
		return int(s[0] - '0'), true
	}

	// 0 for task 10
	if s[0] == '0' {
		return 10, true
	}

	// Shift+number for tasks 11-20
	shiftMap := map[byte]int{
		'!': 11, // shift+1
		'@': 12, // shift+2
		'#': 13, // shift+3
		'$': 14, // shift+4
		'%': 15, // shift+5
		'^': 16, // shift+6
		'&': 17, // shift+7
		'*': 18, // shift+8
		'(': 19, // shift+9
		')': 20, // shift+0
	}
	if num, ok := shiftMap[s[0]]; ok {
		return num, true
	}

	return 0, false
}
