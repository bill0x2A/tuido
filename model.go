package main

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	modeNormal mode = iota
	modeAdding
	modeEditing
	modeRollover
	modeDatePicker
	modeConfirmDelete
	modeHelp
)

type model struct {
	tasks       []Task
	store       *Store
	currentDate time.Time
	cursor      int
	mode        mode
	textInput   textinput.Model
	editingID   string
	rollover    []Task // tasks pending rollover
	width       int
	height      int
}

func newModel(store *Store) model {
	ti := textinput.New()
	ti.Placeholder = "Task description..."
	ti.CharLimit = 256
	ti.Width = 40

	return model{
		tasks:       []Task{},
		store:       store,
		currentDate: normalizeDate(time.Now()),
		cursor:      0,
		mode:        modeNormal,
		textInput:   ti,
	}
}

func (m model) Init() tea.Cmd {
	return m.loadTasks
}

func (m model) loadTasks() tea.Msg {
	tasks, err := m.store.Load()
	if err != nil {
		return errMsg{err}
	}
	return tasksLoadedMsg{tasks}
}

type tasksLoadedMsg struct {
	tasks []Task
}

type errMsg struct {
	err error
}

func (m model) tasksForDate(date time.Time) []Task {
	var result []Task
	for _, t := range m.tasks {
		if sameDay(t.DueDate, date) {
			result = append(result, t)
		}
	}
	return result
}

func (m model) pendingRollover() []Task {
	var result []Task
	today := normalizeDate(time.Now())
	for _, t := range m.tasks {
		if !t.Done && t.DueDate.Before(today) {
			result = append(result, t)
		}
	}
	return result
}

// Update handles all messages - stub for now, will be replaced in update.go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}
