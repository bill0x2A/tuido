package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var content string

	switch m.mode {
	case modeRollover:
		content = m.viewRollover()
	case modeHelp:
		content = m.viewHelp()
	case modeConfirmDelete:
		content = m.viewConfirmDelete()
	case modeDatePicker:
		content = m.viewDatePicker()
	default:
		content = m.viewNormal()
	}

	// Center content in terminal
	return m.centerView(content)
}

func (m model) centerView(content string) string {
	width := m.width
	height := m.height
	if width < 20 {
		width = 80
	}
	if height < 10 {
		height = 24
	}

	// Get content dimensions
	contentWidth := lipgloss.Width(content)
	contentHeight := lipgloss.Height(content)

	// Calculate padding
	padX := (width - contentWidth) / 2
	padY := (height - contentHeight) / 2

	if padX < 0 {
		padX = 0
	}
	if padY < 0 {
		padY = 0
	}

	// Apply horizontal centering
	style := lipgloss.NewStyle().
		PaddingLeft(padX).
		PaddingTop(padY)

	return style.Render(content)
}

func (m model) viewNormal() string {
	var b strings.Builder

	// Use default width if not set yet
	width := m.width
	if width < 20 {
		width = 80
	}

	// Constrain content width for readability
	contentWidth := min(60, width-8)
	if contentWidth < 30 {
		contentWidth = 30
	}

	// Logo
	b.WriteString(logoStyle.Render(logo))
	b.WriteString("\n\n")

	// Date header with optional "today" badge
	dateStr := m.currentDate.Format("Monday, January 2, 2006")
	header := dateStyle.Render(dateStr)
	if sameDay(m.currentDate, normalizeDate(timeNow())) {
		header += todayBadge.Render("TODAY")
	}
	b.WriteString(header)
	b.WriteString("\n")

	// Separator
	sepWidth := min(contentWidth-4, 50)
	if sepWidth < 10 {
		sepWidth = 10
	}
	b.WriteString(separatorStyle.Render(strings.Repeat("─", sepWidth)))
	b.WriteString("\n\n")

	// Tasks
	tasks := m.tasksForDate(m.currentDate)
	if len(tasks) == 0 {
		b.WriteString(emptyStyle.Render("No tasks for this day. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		// Count incomplete vs complete
		incomplete := 0
		for _, t := range tasks {
			if !t.Done {
				incomplete++
			}
		}

		for i, task := range tasks {
			b.WriteString(m.renderTask(task, i, i == m.cursor))
			b.WriteString("\n")
		}

		// Task count summary
		if len(tasks) > 0 {
			complete := len(tasks) - incomplete
			summary := fmt.Sprintf("%d of %d complete", complete, len(tasks))
			b.WriteString("\n")
			b.WriteString(countStyle.Render(summary))
			b.WriteString("\n")
		}
	}

	// Input field if adding/editing
	if m.mode == modeAdding || m.mode == modeEditing {
		b.WriteString("\n")
		b.WriteString(m.textInput.View())
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(m.viewFooter())

	return containerStyle.Width(contentWidth).Render(b.String())
}

func (m model) renderTask(t Task, index int, selected bool) string {
	var parts []string

	// Number key hint
	// 1-9 for tasks 1-9, 0 for task 10, then blank
	num := ""
	if index < 9 {
		if selected {
			num = helpKeyStyle.Render(fmt.Sprintf("%d", index+1))
		} else {
			num = helpDescStyle.Render(fmt.Sprintf("%d", index+1))
		}
	} else if index == 9 {
		if selected {
			num = helpKeyStyle.Render("0")
		} else {
			num = helpDescStyle.Render("0")
		}
	} else {
		num = " "
	}
	parts = append(parts, num)

	// Cursor
	if selected {
		parts = append(parts, cursorStyle)
	} else {
		parts = append(parts, cursorEmpty)
	}

	// Checkbox
	if t.Done {
		parts = append(parts, checkboxDone)
	} else if selected {
		parts = append(parts, checkboxSelected)
	} else {
		parts = append(parts, checkboxEmpty)
	}

	// Title
	title := t.Title
	if t.Done {
		title = taskTitleDoneStyle.Render(title)
	} else if selected {
		title = taskTitleSelectedStyle.Render(title)
	} else {
		title = taskTitleStyle.Render(title)
	}
	parts = append(parts, " "+title)

	// Rollover indicator
	if t.IsRolledOver() {
		parts = append(parts, rolloverIcon)
	}

	return strings.Join(parts, "")
}

func (m model) viewFooter() string {
	hints := []string{
		helpKeyStyle.Render("←→") + helpDescStyle.Render(" days"),
		helpKeyStyle.Render("a") + helpDescStyle.Render(" add"),
		helpKeyStyle.Render("x") + helpDescStyle.Render(" done"),
		helpKeyStyle.Render(">") + helpDescStyle.Render(" tomorrow"),
		helpKeyStyle.Render("?") + helpDescStyle.Render(" help"),
	}
	return strings.Join(hints, "  ")
}

func (m model) viewRollover() string {
	var b strings.Builder

	b.WriteString(dialogTitleStyle.Render("Rollover Tasks"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("%d incomplete tasks from previous days.\n", len(m.rollover)))
	b.WriteString("Roll them over to today?\n\n")

	options := []string{
		dialogKeyStyle.Render("[Y]") + " Yes, move to today",
		dialogKeyStyle.Render("[N]") + " No, leave them",
		dialogKeyStyle.Render("[V]") + " View tasks",
	}
	b.WriteString(strings.Join(options, "    "))

	return dialogStyle.Render(b.String())
}

func (m model) viewHelp() string {
	var b strings.Builder

	b.WriteString(dialogTitleStyle.Render("Keyboard Shortcuts"))
	b.WriteString("\n")

	// Navigation section
	b.WriteString(sectionStyle.Render("Navigation"))
	b.WriteString("\n")
	navItems := []string{
		helpKeyStyle.Render("  ←/→ h/l  ") + "Previous/next day",
		helpKeyStyle.Render("  t        ") + "Jump to today",
		helpKeyStyle.Render("  g        ") + "Go to specific date",
		helpKeyStyle.Render("  j/k ↑/↓  ") + "Move selection",
		helpKeyStyle.Render("  1-0      ") + "Jump to task 1-10",
		helpKeyStyle.Render("  !-⇧0     ") + "Jump to task 11-20",
	}
	b.WriteString(strings.Join(navItems, "\n"))
	b.WriteString("\n")

	// Tasks section
	b.WriteString(sectionStyle.Render("Tasks"))
	b.WriteString("\n")
	taskItems := []string{
		helpKeyStyle.Render("  a        ") + "Add new task",
		helpKeyStyle.Render("  e        ") + "Edit task",
		helpKeyStyle.Render("  x/⏎/space") + "Toggle done",
		helpKeyStyle.Render("  d        ") + "Delete task",
		helpKeyStyle.Render("  >/n      ") + "Move to tomorrow",
		helpKeyStyle.Render("  J/K ⇧↑/↓ ") + "Reorder task",
	}
	b.WriteString(strings.Join(taskItems, "\n"))
	b.WriteString("\n")

	// General section
	b.WriteString(sectionStyle.Render("General"))
	b.WriteString("\n")
	generalItems := []string{
		helpKeyStyle.Render("  ?        ") + "Toggle help",
		helpKeyStyle.Render("  q/esc    ") + "Quit",
	}
	b.WriteString(strings.Join(generalItems, "\n"))
	b.WriteString("\n\n")

	b.WriteString(helpDescStyle.Render("Press any key to close"))

	return dialogStyle.Render(b.String())
}

func (m model) viewConfirmDelete() string {
	tasks := m.tasksForDate(m.currentDate)
	if m.cursor >= len(tasks) {
		return ""
	}
	task := tasks[m.cursor]

	var b strings.Builder
	b.WriteString(dialogTitleStyle.Render("Delete Task?"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("\"%s\"\n\n", task.Title))
	b.WriteString(dialogKeyStyle.Render("[Y]") + " Yes  ")
	b.WriteString(dialogKeyStyle.Render("[N]") + " No")

	return dialogStyle.Render(b.String())
}

func (m model) viewDatePicker() string {
	var b strings.Builder

	b.WriteString(dialogTitleStyle.Render("Go to Date"))
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpDescStyle.Render("Format: 2025-01-15 or +7 / -3"))

	return dialogStyle.Render(b.String())
}

// timeNow is a variable for testing
var timeNow = func() time.Time {
	return time.Now()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
