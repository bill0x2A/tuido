package main

import (
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAppFlow(t *testing.T) {
	// Setup
	dir := t.TempDir()
	store := NewStore(filepath.Join(dir, "tasks.json"))

	// Mock time
	fixedTime := time.Date(2025, 1, 11, 12, 0, 0, 0, time.Local)
	oldTimeNow := timeNow
	timeNow = func() time.Time { return fixedTime }
	defer func() { timeNow = oldTimeNow }()

	m := newModel(store)
	m.width = 80
	m.height = 24

	// Simulate tasks loaded
	updated, _ := m.Update(tasksLoadedMsg{tasks: []Task{}})
	m = updated.(model)

	// Verify initial state
	if m.mode != modeNormal {
		t.Errorf("expected normal mode, got %v", m.mode)
	}

	// Add a task
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(model)
	if m.mode != modeAdding {
		t.Error("expected adding mode after pressing 'a'")
	}

	// Type task title
	for _, r := range "Test task" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = updated.(model)
	}

	// Press enter to save
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	if m.mode != modeNormal {
		t.Error("expected normal mode after saving task")
	}

	tasks := m.tasksForDate(m.currentDate)
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Title != "Test task" {
		t.Errorf("expected title 'Test task', got %q", tasks[0].Title)
	}

	// Toggle task done
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	tasks = m.tasksForDate(m.currentDate)
	if !tasks[0].Done {
		t.Error("task should be done after toggle")
	}
}
