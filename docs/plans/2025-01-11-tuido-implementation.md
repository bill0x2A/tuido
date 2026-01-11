# tuido Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a minimal TUI todo app with day-based task management and rollover functionality.

**Architecture:** Bubble Tea MVC pattern with a single JSON file store. Model holds app state (current date, tasks, input mode), View renders based on state, Update handles key events and state transitions.

**Tech Stack:** Go 1.21+, Bubble Tea, Lip Gloss, Bubbles (for text input)

---

### Task 1: Project Setup

**Files:**
- Create: `go.mod`
- Create: `main.go`

**Step 1: Initialize Go module**

Run: `go mod init tuido`
Expected: Creates go.mod file

**Step 2: Add dependencies**

Run: `go get github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss github.com/charmbracelet/bubbles github.com/google/uuid`
Expected: Downloads dependencies, updates go.mod and go.sum

**Step 3: Create minimal main.go**

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "tuido - press q to quit\n"
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

**Step 4: Verify it runs**

Run: `go run main.go`
Expected: Shows "tuido - press q to quit", pressing q exits cleanly

**Step 5: Commit**

```bash
git add go.mod go.sum main.go
git commit -m "feat: initial project setup with Bubble Tea"
```

---

### Task 2: Task Data Model

**Files:**
- Create: `task.go`
- Create: `task_test.go`

**Step 1: Write tests for Task type**

```go
package main

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	title := "Buy groceries"
	dueDate := time.Date(2025, 1, 11, 0, 0, 0, 0, time.Local)

	task := NewTask(title, dueDate)

	if task.Title != title {
		t.Errorf("expected title %q, got %q", title, task.Title)
	}
	if task.Done {
		t.Error("new task should not be done")
	}
	if task.ID == "" {
		t.Error("task should have an ID")
	}
	if !task.DueDate.Equal(dueDate) {
		t.Errorf("expected due date %v, got %v", dueDate, task.DueDate)
	}
}

func TestTaskIsRolledOver(t *testing.T) {
	yesterday := time.Date(2025, 1, 10, 0, 0, 0, 0, time.Local)
	today := time.Date(2025, 1, 11, 0, 0, 0, 0, time.Local)

	task := NewTask("Test task", yesterday)

	if task.IsRolledOver() {
		t.Error("task created today should not be rolled over")
	}

	task.DueDate = today
	if !task.IsRolledOver() {
		t.Error("task with different due date than created date should be rolled over")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test -v -run TestNewTask`
Expected: FAIL - NewTask not defined

**Step 3: Implement Task type**

```go
package main

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	DueDate   time.Time `json:"due_date"`
}

func NewTask(title string, dueDate time.Time) Task {
	now := time.Now()
	return Task{
		ID:        uuid.New().String(),
		Title:     title,
		Done:      false,
		CreatedAt: now,
		DueDate:   normalizeDate(dueDate),
	}
}

func (t Task) IsRolledOver() bool {
	return !sameDay(t.CreatedAt, t.DueDate)
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -v`
Expected: PASS

**Step 5: Commit**

```bash
git add task.go task_test.go
git commit -m "feat: add Task data model with rollover detection"
```

---

### Task 3: JSON Store

**Files:**
- Create: `store.go`
- Create: `store_test.go`

**Step 1: Write tests for Store**

```go
package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	store := NewStore(path)
	tasks, err := store.Load()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty tasks, got %d", len(tasks))
	}
}

func TestStoreSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	store := NewStore(path)
	today := time.Date(2025, 1, 11, 0, 0, 0, 0, time.Local)
	tasks := []Task{
		NewTask("Task 1", today),
		NewTask("Task 2", today),
	}

	err := store.Save(tasks)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(loaded))
	}
	if loaded[0].Title != "Task 1" {
		t.Errorf("expected title 'Task 1', got %q", loaded[0].Title)
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test -v -run TestStore`
Expected: FAIL - NewStore not defined

**Step 3: Implement Store**

```go
package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Store struct {
	path string
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func DefaultStorePath() string {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataDir, "tuido", "tasks.json")
}

func (s *Store) Load() ([]Task, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Store) Save(tasks []Task) error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	// Atomic write: temp file then rename
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -v -run TestStore`
Expected: PASS

**Step 5: Commit**

```bash
git add store.go store_test.go
git commit -m "feat: add JSON store with atomic writes"
```

---

### Task 4: Core Model Structure

**Files:**
- Modify: `main.go`
- Create: `model.go`
- Create: `keys.go`

**Step 1: Create key bindings**

```go
package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Today    key.Binding
	GoTo     key.Binding
	Add      key.Binding
	Edit     key.Binding
	Toggle   key.Binding
	Delete   key.Binding
	MoveUp   key.Binding
	MoveDown key.Binding
	Help     key.Binding
	Quit     key.Binding
	Enter    key.Binding
	Escape   key.Binding
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
		key.WithKeys("left"),
		key.WithHelp("←", "prev day"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "next day"),
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
		key.WithKeys("enter", " "),
		key.WithHelp("⏎/space", "toggle"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	MoveUp: key.NewBinding(
		key.WithKeys("K"),
		key.WithHelp("K", "move up"),
	),
	MoveDown: key.NewBinding(
		key.WithKeys("J"),
		key.WithHelp("J", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
	),
}
```

**Step 2: Create model.go with app state**

```go
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
```

**Step 3: Update main.go**

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	store := NewStore(DefaultStorePath())
	m := newModel(store)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

**Step 4: Verify it compiles**

Run: `go build`
Expected: Compiles without errors (will panic on run since Update/View not complete yet)

**Step 5: Commit**

```bash
git add main.go model.go keys.go
git commit -m "feat: add core model structure and key bindings"
```

---

### Task 5: Styles

**Files:**
- Create: `styles.go`

**Step 1: Create styles**

```go
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
```

**Step 2: Verify it compiles**

Run: `go build`
Expected: Compiles without errors

**Step 3: Commit**

```bash
git add styles.go
git commit -m "feat: add Lip Gloss styles"
```

---

### Task 6: View Rendering

**Files:**
- Create: `view.go`

**Step 1: Create view.go**

```go
package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

	// Header with date
	dateStr := m.currentDate.Format("Monday, January 2, 2006")
	header := dateStyle.Render(dateStr)
	if sameDay(m.currentDate, normalizeDate(timeNow())) {
		header += " (today)"
	}
	b.WriteString(titleStyle.Render("tuido") + "\n\n")
	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("─", min(50, m.width-4)) + "\n\n")

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

	return borderStyle.Width(m.width - 4).Render(b.String())
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
```

**Step 2: Add missing import**

Add `"time"` to imports in view.go.

**Step 3: Verify it compiles**

Run: `go build`
Expected: Compiles without errors

**Step 4: Commit**

```bash
git add view.go
git commit -m "feat: add view rendering for all modes"
```

---

### Task 7: Update Logic

**Files:**
- Create: `update.go`

**Step 1: Create update.go with core logic**

```go
package main

import (
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
		m.textInput.Placeholder = "Task description..."
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
			m.toggleTask(tasks[m.cursor].ID)
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
			m.swapTasks(tasks[m.cursor].ID, tasks[m.cursor-1].ID)
			m.cursor--
			return m, m.saveTasks
		}

	case key.Matches(msg, keys.MoveDown):
		if len(tasks) > 1 && m.cursor < len(tasks)-1 {
			m.swapTasks(tasks[m.cursor].ID, tasks[m.cursor+1].ID)
			m.cursor++
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
			m.cursor = len(m.tasksForDate(m.currentDate)) - 1
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
		// TODO: show rollover list view
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
```

**Step 2: Verify it compiles and runs**

Run: `go build && ./tuido`
Expected: App runs, shows today's date, can quit with 'q'

**Step 3: Commit**

```bash
git add update.go
git commit -m "feat: add update logic for all modes and interactions"
```

---

### Task 8: Integration Test

**Files:**
- Create: `main_test.go`

**Step 1: Write integration test**

```go
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
	m, _ = m.Update(tasksLoadedMsg{tasks: []Task{}}).(model)

	// Verify initial state
	if m.mode != modeNormal {
		t.Errorf("expected normal mode, got %v", m.mode)
	}

	// Add a task
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}).(model)
	if m.mode != modeAdding {
		t.Error("expected adding mode after pressing 'a'")
	}

	// Type task title
	for _, r := range "Test task" {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}).(model)
	}

	// Press enter to save
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}).(model)

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
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}).(model)
	tasks = m.tasksForDate(m.currentDate)
	if !tasks[0].Done {
		t.Error("task should be done after toggle")
	}
}
```

**Step 2: Run test**

Run: `go test -v -run TestAppFlow`
Expected: PASS

**Step 3: Commit**

```bash
git add main_test.go
git commit -m "test: add integration test for basic app flow"
```

---

### Task 9: Final Polish

**Files:**
- Modify: `view.go` (minor refinements if needed)

**Step 1: Run the app and verify all features**

Run: `go run .`
Expected: All features work as designed

**Step 2: Build release binary**

Run: `go build -o tuido`
Expected: Creates `tuido` binary

**Step 3: Final commit**

```bash
git add -A
git commit -m "chore: final polish and build"
```

---

## Summary

The implementation covers:
1. Project setup with Bubble Tea
2. Task data model with rollover detection
3. JSON persistence with atomic writes
4. Key bindings for all actions
5. Lip Gloss styling
6. View rendering for all modes
7. Update logic for all interactions
8. Integration tests

Total: 9 tasks, each with bite-sized TDD steps.
