# tuido

A minimal TUI todo app with day-based task management.

![tuido demo](https://img.shields.io/badge/go-%3E%3D1.21-blue)

## Features

- **Day-based tasks** - Navigate between days with arrow keys
- **Task rollover** - Incomplete tasks from past days prompt for rollover
- **Vim-style navigation** - Use `j`/`k` to move, `J`/`K` to reorder
- **Persistent storage** - Tasks saved to `~/.local/share/tuido/tasks.json`
- **Clean TUI** - Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss)

## Installation

### From source (requires Go 1.21+)

```bash
go install github.com/yourusername/tuido@latest
```

Or clone and build:

```bash
git clone https://github.com/yourusername/tuido
cd tuido
make install
```

### Manual install

```bash
git clone https://github.com/yourusername/tuido
cd tuido
go build -o tuido
sudo mv tuido /usr/local/bin/
```

## Usage

```bash
tuido
```

## Keybindings

### Navigation
| Key | Action |
|-----|--------|
| `j` / `k` or `↑` / `↓` | Move selection up/down |
| `←` / `→` | Previous/next day |
| `t` | Jump to today |
| `g` | Go to specific date |

### Tasks
| Key | Action |
|-----|--------|
| `a` | Add new task |
| `e` | Edit selected task |
| `Enter` / `Space` | Toggle done |
| `d` | Delete task |
| `J` / `K` | Move task up/down |

### General
| Key | Action |
|-----|--------|
| `?` | Show help |
| `q` | Quit |

## Data Storage

Tasks are stored in JSON format at:
- Linux/macOS: `~/.local/share/tuido/tasks.json`
- With `XDG_DATA_HOME` set: `$XDG_DATA_HOME/tuido/tasks.json`

## License

MIT
