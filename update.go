package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tasksLoadedMsg:
		m.tasks = msg.tasks
		m.rollover = m.pendingRollover()
		if len(m.rollover) > 0 {
			m.mode = modeRollover
		}
		return m, nil

	case errMsg:
		// TODO: show error in UI
		return m, nil
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if key.Matches(msg, keys.Quit) && m.mode == modeNormal {
		return m, tea.Quit
	}

	switch m.mode {
	case modeNormal:
		return m.handleNormalMode(msg)
	case modeAdding:
		return m.handleAddingMode(msg)
	case modeEditing:
		return m.handleEditingMode(msg)
	case modeRollover:
		return m.handleRolloverMode(msg)
	case modeDatePicker:
		return m.handleDatePickerMode(msg)
	case modeConfirmDelete:
		return m.handleConfirmDeleteMode(msg)
	case modeHelp:
		m.mode = modeNormal
		return m, nil
	}

	return m, nil
}

func (m model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	tasks := m.tasksForDate(m.currentDate)

	// Handle number keys 1-9 for quick jump
	if num, ok := isNumberKey(msg.String()); ok {
		if num <= len(tasks) {
			m.cursor = num - 1
		}
		return m, nil
	}

	switch {
	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(msg, keys.Down):
		if m.cursor < len(tasks)-1 {
			m.cursor++
		}

	case key.Matches(msg, keys.Left):
		m.currentDate = m.currentDate.AddDate(0, 0, -1)
		m.cursor = 0

	case key.Matches(msg, keys.Right):
		m.currentDate = m.currentDate.AddDate(0, 0, 1)
		m.cursor = 0

	case key.Matches(msg, keys.Today):
		m.currentDate = normalizeDate(timeNow())
		m.cursor = 0

	case key.Matches(msg, keys.GoTo):
		m.mode = modeDatePicker
		m.textInput.Reset()
		m.textInput.Placeholder = "2025-01-15 or +7"
		m.textInput.Focus()
		return m, nil

	case key.Matches(msg, keys.Add):
		m.mode = modeAdding
		m.textInput.Reset()
		m.textInput.Placeholder = "What needs to be done?"
		m.textInput.Focus()
		return m, nil

	case key.Matches(msg, keys.Edit):
		if len(tasks) > 0 && m.cursor < len(tasks) {
			m.mode = modeEditing
			m.editingID = tasks[m.cursor].ID
			m.textInput.SetValue(tasks[m.cursor].Title)
			m.textInput.Focus()
		}
		return m, nil

	case key.Matches(msg, keys.Toggle):
		if len(tasks) > 0 && m.cursor < len(tasks) {
			taskID := tasks[m.cursor].ID
			m.toggleTask(taskID)
			m.sortTasksForDate(m.currentDate)
			// Find where the task ended up after sorting
			newTasks := m.tasksForDate(m.currentDate)
			for i, t := range newTasks {
				if t.ID == taskID {
					m.cursor = i
					break
				}
			}
			return m, m.saveTasks
		}

	case key.Matches(msg, keys.Delete):
		if len(tasks) > 0 && m.cursor < len(tasks) {
			if tasks[m.cursor].Done {
				m.deleteTask(tasks[m.cursor].ID)
				if m.cursor >= len(m.tasksForDate(m.currentDate)) {
					m.cursor = max(0, m.cursor-1)
				}
				return m, m.saveTasks
			}
			m.mode = modeConfirmDelete
		}

	case key.Matches(msg, keys.MoveUp):
		if len(tasks) > 1 && m.cursor > 0 {
			// Only allow moving within incomplete tasks
			if !tasks[m.cursor].Done && !tasks[m.cursor-1].Done {
				m.swapTasks(tasks[m.cursor].ID, tasks[m.cursor-1].ID)
				m.cursor--
				return m, m.saveTasks
			}
		}

	case key.Matches(msg, keys.MoveDown):
		if len(tasks) > 1 && m.cursor < len(tasks)-1 {
			// Only allow moving within incomplete tasks
			if !tasks[m.cursor].Done && !tasks[m.cursor+1].Done {
				m.swapTasks(tasks[m.cursor].ID, tasks[m.cursor+1].ID)
				m.cursor++
				return m, m.saveTasks
			}
		}

	case key.Matches(msg, keys.MoveNext):
		if len(tasks) > 0 && m.cursor < len(tasks) {
			m.moveTaskToNextDay(tasks[m.cursor].ID)
			if m.cursor >= len(m.tasksForDate(m.currentDate)) {
				m.cursor = max(0, m.cursor-1)
			}
			return m, m.saveTasks
		}

	case key.Matches(msg, keys.Help):
		m.mode = modeHelp
	}

	return m, nil
}

func (m model) handleAddingMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Escape):
		m.mode = modeNormal
		m.textInput.Reset()
		return m, nil

	case key.Matches(msg, keys.Enter):
		title := strings.TrimSpace(m.textInput.Value())
		if title != "" {
			task := NewTask(title, m.currentDate)
			m.tasks = append(m.tasks, task)
			// Position cursor at the new task (last incomplete task)
			newTasks := m.tasksForDate(m.currentDate)
			for i, t := range newTasks {
				if t.ID == task.ID {
					m.cursor = i
					break
				}
			}
		}
		m.mode = modeNormal
		m.textInput.Reset()
		return m, m.saveTasks

	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m model) handleEditingMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Escape):
		m.mode = modeNormal
		m.textInput.Reset()
		m.editingID = ""
		return m, nil

	case key.Matches(msg, keys.Enter):
		title := strings.TrimSpace(m.textInput.Value())
		if title != "" {
			m.updateTaskTitle(m.editingID, title)
		}
		m.mode = modeNormal
		m.textInput.Reset()
		m.editingID = ""
		return m, m.saveTasks

	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m model) handleRolloverMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		today := normalizeDate(timeNow())
		for i := range m.tasks {
			for _, r := range m.rollover {
				if m.tasks[i].ID == r.ID {
					m.tasks[i].DueDate = today
				}
			}
		}
		m.rollover = nil
		m.mode = modeNormal
		return m, m.saveTasks

	case "n", "N":
		m.rollover = nil
		m.mode = modeNormal

	case "v", "V":
		m.rollover = nil
		m.mode = modeNormal
	}

	return m, nil
}

func (m model) handleDatePickerMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Escape):
		m.mode = modeNormal
		m.textInput.Reset()
		return m, nil

	case key.Matches(msg, keys.Enter):
		input := strings.TrimSpace(m.textInput.Value())
		if date, ok := parseDate(input, m.currentDate); ok {
			m.currentDate = date
			m.cursor = 0
		}
		m.mode = modeNormal
		m.textInput.Reset()
		return m, nil

	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m model) handleConfirmDeleteMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		tasks := m.tasksForDate(m.currentDate)
		if m.cursor < len(tasks) {
			m.deleteTask(tasks[m.cursor].ID)
			if m.cursor >= len(m.tasksForDate(m.currentDate)) {
				m.cursor = max(0, m.cursor-1)
			}
		}
		m.mode = modeNormal
		return m, m.saveTasks

	case "n", "N", "esc":
		m.mode = modeNormal
	}

	return m, nil
}

func (m *model) toggleTask(id string) {
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			m.tasks[i].Done = !m.tasks[i].Done
			return
		}
	}
}

func (m *model) deleteTask(id string) {
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return
		}
	}
}

func (m *model) updateTaskTitle(id, title string) {
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			m.tasks[i].Title = title
			return
		}
	}
}

func (m *model) swapTasks(id1, id2 string) {
	var i1, i2 int = -1, -1
	for i := range m.tasks {
		if m.tasks[i].ID == id1 {
			i1 = i
		}
		if m.tasks[i].ID == id2 {
			i2 = i
		}
	}
	if i1 >= 0 && i2 >= 0 {
		m.tasks[i1], m.tasks[i2] = m.tasks[i2], m.tasks[i1]
	}
}

func (m *model) moveTaskToNextDay(id string) {
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			m.tasks[i].DueDate = m.currentDate.AddDate(0, 0, 1)
			return
		}
	}
}

// sortTasksForDate sorts tasks for a specific date: incomplete first, then completed
func (m *model) sortTasksForDate(date time.Time) {
	// Get indices of tasks for this date
	type taskIndex struct {
		index int
		done  bool
	}
	var indices []taskIndex
	for i, t := range m.tasks {
		if sameDay(t.DueDate, date) {
			indices = append(indices, taskIndex{i, t.Done})
		}
	}

	// Sort: incomplete first, then completed
	sort.SliceStable(indices, func(i, j int) bool {
		if indices[i].done != indices[j].done {
			return !indices[i].done // incomplete (false) comes before completed (true)
		}
		return false // maintain relative order within groups
	})

	// Rebuild the task slice with sorted order for this date
	var newTasks []Task
	dateTaskIndices := make(map[int]bool)
	for _, idx := range indices {
		dateTaskIndices[idx.index] = true
	}

	// First, add non-date tasks in original order, inserting date tasks at first date task position
	inserted := false
	for i, t := range m.tasks {
		if dateTaskIndices[i] {
			if !inserted {
				// Insert all date tasks here
				for _, idx := range indices {
					newTasks = append(newTasks, m.tasks[idx.index])
				}
				inserted = true
			}
		} else {
			newTasks = append(newTasks, t)
		}
	}
	// If date tasks were at the end
	if !inserted {
		for _, idx := range indices {
			newTasks = append(newTasks, m.tasks[idx.index])
		}
	}

	m.tasks = newTasks
}

func (m model) saveTasks() tea.Msg {
	if err := m.store.Save(m.tasks); err != nil {
		return errMsg{err}
	}
	return nil
}

func parseDate(input string, current time.Time) (time.Time, bool) {
	// Relative: +7, -3
	if strings.HasPrefix(input, "+") || strings.HasPrefix(input, "-") {
		days, err := strconv.Atoi(input)
		if err == nil {
			return current.AddDate(0, 0, days), true
		}
	}

	// Absolute: 2025-01-15
	t, err := time.Parse("2006-01-02", input)
	if err == nil {
		return normalizeDate(t), true
	}

	return time.Time{}, false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
