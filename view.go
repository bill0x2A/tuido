package main

import (
	"fmt"
	"strings"
	"time"
)

func (m model) View() string {
	switch m.mode {
	case modeRollover:
		return m.viewRollover()
	case modeHelp:
		return m.viewHelp()
	case modeConfirmDelete:
		return m.viewConfirmDelete()
	case modeDatePicker:
		return m.viewDatePicker()
	default:
		return m.viewNormal()
	}
}

func (m model) viewNormal() string {
	var b strings.Builder

	// Use default width if not set yet
	width := m.width
	if width < 20 {
		width = 80
	}

	// Header with date
	dateStr := m.currentDate.Format("Monday, January 2, 2006")
	header := dateStyle.Render(dateStr)
	if sameDay(m.currentDate, normalizeDate(timeNow())) {
		header += " (today)"
	}
	b.WriteString(titleStyle.Render("tuido") + "\n\n")
	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("─", min(50, width-4)) + "\n\n")

	// Tasks
	tasks := m.tasksForDate(m.currentDate)
	if len(tasks) == 0 {
		empty := helpStyle.Render("No tasks for this day. Press 'a' to add one.")
		b.WriteString(empty + "\n")
	} else {
		for i, task := range tasks {
			b.WriteString(m.renderTask(task, i == m.cursor) + "\n")
		}
	}

	// Input field if adding/editing
	if m.mode == modeAdding || m.mode == modeEditing {
		b.WriteString("\n" + m.textInput.View() + "\n")
	}

	// Footer
	b.WriteString("\n" + m.viewFooter())

	return borderStyle.Width(width - 4).Render(b.String())
}

func (m model) renderTask(t Task, selected bool) string {
	var checkbox string
	if t.Done {
		checkbox = checkboxDone
	} else {
		checkbox = checkboxEmpty
	}

	title := t.Title
	if t.Done {
		title = doneStyle.Render(title)
	}

	var cursor string
	if selected {
		cursor = cursorStyle
	} else {
		cursor = "  "
	}

	line := fmt.Sprintf("%s%s  %s", cursor, checkbox, title)

	if t.IsRolledOver() {
		line += "  " + rolloverIcon
	}

	return line
}

func (m model) viewFooter() string {
	hints := []string{
		"←→ days",
		"a add",
		"⏎ toggle",
		"? help",
	}
	return helpStyle.Render(strings.Join(hints, "  "))
}

func (m model) viewRollover() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d incomplete tasks from previous days.\n", len(m.rollover)))
	b.WriteString("Roll them over to today?\n\n")
	b.WriteString("[Y] Yes     [N] No     [V] View\n")
	return dialogStyle.Render(b.String())
}

func (m model) viewHelp() string {
	help := `Navigation:
  ←/→      Previous/next day
  t        Jump to today
  g        Go to date
  j/k ↑/↓  Move selection

Tasks:
  a        Add task
  e        Edit task
  ⏎/space  Toggle done
  d        Delete task
  J/K      Reorder task

General:
  ?        Toggle help
  q        Quit

Press any key to close help`
	return dialogStyle.Render(help)
}

func (m model) viewConfirmDelete() string {
	tasks := m.tasksForDate(m.currentDate)
	if m.cursor >= len(tasks) {
		return ""
	}
	task := tasks[m.cursor]
	msg := fmt.Sprintf("Delete task?\n\n\"%s\"\n\n[Y] Yes  [N] No", task.Title)
	return dialogStyle.Render(msg)
}

func (m model) viewDatePicker() string {
	var b strings.Builder
	b.WriteString("Go to date:\n\n")
	b.WriteString(m.textInput.View() + "\n\n")
	b.WriteString(helpStyle.Render("Format: 2025-01-15 or +7 / -3"))
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
