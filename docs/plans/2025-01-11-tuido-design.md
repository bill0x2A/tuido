# tuido - Personal TUI Todo App

A minimal, clean TUI todo application built with Go and Bubble Tea.

## Core Features

- Add tasks with a title
- Toggle tasks done/undone
- Reorder tasks with hotkeys
- Browse previous and future days
- Incomplete tasks roll over with a prompt on startup
- Visual distinction between new and rolled-over tasks

## Technology

- **Language:** Go
- **TUI Framework:** Bubble Tea
- **Storage:** Single JSON file at `~/.local/share/tuido/tasks.json`

## Data Model

```go
type Task struct {
    ID        string    `json:"id"`        // UUID
    Title     string    `json:"title"`
    Done      bool      `json:"done"`
    CreatedAt time.Time `json:"created_at"` // Original creation date
    DueDate   time.Time `json:"due_date"`   // The date this task belongs to
}

type Store struct {
    Tasks []Task `json:"tasks"`
}
```

A task's `DueDate` is the day it appears on. When a task rolls over, its `DueDate` updates to today while `CreatedAt` stays the same. Rollover indicator: `DueDate != CreatedAt.Date()`.

## Project Structure

```
tuido/
├── main.go           # Entry point, Bubble Tea setup
├── model.go          # App state, Update logic
├── view.go           # Rendering
├── task.go           # Task type and methods
├── store.go          # JSON file read/write
├── keys.go           # Key bindings
└── styles.go         # Lip Gloss styles
```

## UI Layout

### Main View (Today's Tasks)

```
┌─ tuido ──────────────────────────────────────┐
│                                              │
│  Saturday, January 11, 2025          [t][g]  │
│  ─────────────────────────────────────────   │
│                                              │
│    ○  Buy groceries                     ↻    │
│    ○  Review pull request                    │
│  ▸ ○  Write design doc                  ↻    │
│    ●̶  ̶F̶i̶x̶ ̶l̶o̶g̶i̶n̶ ̶b̶u̶g̶                          │
│                                              │
│                                              │
│                                              │
│  ←→ navigate days  a add  ⏎ toggle  ? help   │
└──────────────────────────────────────────────┘
```

### Visual Elements

- `▸` - Cursor/selection indicator
- `○` - Incomplete task
- `●` with strikethrough text - Completed task
- `↻` - Rolled over from a previous day
- Muted footer showing key hints

### Date Picker (triggered by `g`)

Simple popup: type a date like `2025-01-15` or use relative shortcuts (`-1` for yesterday, `+7` for next week). Press Enter to jump, Escape to cancel.

### Rollover Prompt (on startup)

```
┌─ Rollover ───────────────────────────────────┐
│                                              │
│  3 incomplete tasks from previous days.      │
│  Roll them over to today?                    │
│                                              │
│         [Y] Yes     [N] No     [V] View      │
│                                              │
└──────────────────────────────────────────────┘
```

`V` shows the list so you can decide.

## Key Bindings

### Navigation

| Key | Action |
|-----|--------|
| `←` / `→` | Previous / next day |
| `t` | Jump to today |
| `g` | Open date picker |
| `j` / `↓` | Move selection down |
| `k` / `↑` | Move selection up |

### Task Management

| Key | Action |
|-----|--------|
| `a` | Add new task (inline text input appears) |
| `Enter` / `Space` | Toggle task done/undone |
| `d` | Delete task (with confirmation if not done) |
| `e` | Edit task title |
| `J` (shift) | Move task down in list |
| `K` (shift) | Move task up in list |

### General

| Key | Action |
|-----|--------|
| `?` | Toggle help overlay |
| `q` / `Ctrl+C` | Quit |

## Interaction Flows

### Adding a Task

1. Press `a`
2. Text input appears at bottom of task list
3. Type task title
4. Press `Enter` to save, `Escape` to cancel
5. New task appears, cursor moves to it

### Editing a Task

1. Select task with `j`/`k`
2. Press `e`
3. Inline editor appears with existing text
4. Edit and press `Enter` to save, `Escape` to cancel

## Behavior

### Rollover Logic

- On startup, find all tasks where `DueDate < today` and `Done == false`
- If any exist, show the rollover prompt
- "Yes" updates their `DueDate` to today
- "No" leaves them on their original days (visible when browsing back)
- "View" shows the list, then asks again

### Empty States

- No tasks today: "No tasks for today. Press 'a' to add one."
- Browsing an empty past day: "No tasks on this day."

### Persistence

- Save to disk after every mutation (add, edit, delete, toggle, reorder, rollover)
- Atomic write: write to temp file, then rename (prevents corruption)

### Task Ordering

- Tasks maintain explicit order (stored as array order in JSON)
- New tasks append to bottom of incomplete tasks (above completed)
- Reordering with `J`/`K` only moves within the incomplete section

### Browsing History

- No limit on how far back you can browse
- Days with no tasks show empty state
- Future dates accessible for planning ahead

### Startup Sequence

1. Load JSON file (create if missing)
2. Check for rollover candidates
3. Show rollover prompt if any, otherwise go straight to today's view

## Visual Style

- Minimal and clean
- Subtle borders, muted colors
- Single accent color for highlights
- Focus on content over decoration
