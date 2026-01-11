package main

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight)

	dateStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	taskStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			PaddingLeft(0).
			Foreground(highlight)

	doneStyle = lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(subtle)

	rolloverIcon = lipgloss.NewStyle().
			Foreground(subtle).
			Render("↻")

	cursorStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Render("▸ ")

	checkboxEmpty = lipgloss.NewStyle().
			Foreground(subtle).
			Render("○")

	checkboxDone = lipgloss.NewStyle().
			Foreground(special).
			Render("●")

	helpStyle = lipgloss.NewStyle().
			Foreground(subtle)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(1, 2)

	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 2).
			Width(50)
)
