package main

import "github.com/charmbracelet/lipgloss"

// ASCII art logo
var logo = `┌┬┐┬ ┬┬┌┬┐┌─┐
 │ │ ││ │││ │
 ┴ └─┘┴─┴┘└─┘`

// Color palette - warm, muted tones with purposeful accents
var (
	// Base colors
	base      = lipgloss.AdaptiveColor{Light: "#1a1a2e", Dark: "#eaeaea"}
	muted     = lipgloss.AdaptiveColor{Light: "#6b7280", Dark: "#6b7280"}
	subtle    = lipgloss.AdaptiveColor{Light: "#9ca3af", Dark: "#4b5563"}
	surface   = lipgloss.AdaptiveColor{Light: "#f8fafc", Dark: "#1e1e2e"}
	overlay   = lipgloss.AdaptiveColor{Light: "#e2e8f0", Dark: "#313244"}

	// Accent colors
	accent    = lipgloss.AdaptiveColor{Light: "#6366f1", Dark: "#a5b4fc"}  // Indigo
	success   = lipgloss.AdaptiveColor{Light: "#10b981", Dark: "#6ee7b7"}  // Emerald
	warning   = lipgloss.AdaptiveColor{Light: "#f59e0b", Dark: "#fcd34d"}  // Amber
	rose      = lipgloss.AdaptiveColor{Light: "#f43f5e", Dark: "#fb7185"}  // Rose

	// App title - bold accent
	titleStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	// Logo style
	logoStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			MarginBottom(1)

	// Date header - prominent but not overwhelming
	dateStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(base)

	// "today" badge
	todayBadge = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#1e1e2e"}).
			Background(accent).
			Padding(0, 1).
			MarginLeft(1).
			Bold(true)

	// Separator line
	separatorStyle = lipgloss.NewStyle().
			Foreground(subtle)

	// Task styles
	taskStyle = lipgloss.NewStyle()

	taskTitleStyle = lipgloss.NewStyle().
			Foreground(base)

	taskTitleDoneStyle = lipgloss.NewStyle().
				Foreground(muted).
				Strikethrough(true)

	taskTitleSelectedStyle = lipgloss.NewStyle().
				Foreground(accent).
				Bold(true)

	// Checkboxes - using nice Unicode
	checkboxEmpty    = lipgloss.NewStyle().Foreground(subtle).Render("○")
	checkboxDone     = lipgloss.NewStyle().Foreground(success).Render("●")
	checkboxSelected = lipgloss.NewStyle().Foreground(accent).Render("◉")

	// Cursor
	cursorStyle    = lipgloss.NewStyle().Foreground(accent).Render("▸ ")
	cursorEmpty    = "  "

	// Icons
	rolloverIcon = lipgloss.NewStyle().Foreground(warning).Render(" ↻")
	doneIcon     = lipgloss.NewStyle().Foreground(success).Render(" ✓")

	// Empty state
	emptyStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true).
			MarginTop(1).
			MarginBottom(1)

	// Footer / help hints
	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(muted).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(subtle)

	// Main container
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(1, 2)

	// Dialog boxes
	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(1, 2).
			Width(54)

	dialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accent).
				MarginBottom(1)

	dialogKeyStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	// Section headers in help
	sectionStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			MarginTop(1)

	// Task count badge
	countStyle = lipgloss.NewStyle().
			Foreground(muted).
			MarginLeft(1)
)
