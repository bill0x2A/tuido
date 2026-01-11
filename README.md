# tuido

A minimal TUI todo app with day-based task management.

![tuido demo](https://img.shields.io/badge/go-%3E%3D1.21-blue)

## Features

- **Day-based tasks** - Navigate between days with arrow keys or `h`/`l`
- **Quick jump** - Press `1-9` to jump directly to a task
- **Task rollover** - Incomplete tasks from past days prompt for rollover
- **Move to tomorrow** - Press `>` or `n` to bump a task to the next day
- **Auto-sort** - Completed tasks automatically move to the bottom
- **Vim-style navigation** - Use `j`/`k` to move, `J`/`K` or `Shift+↑/↓` to reorder
- **Persistent storage** - Tasks saved to `~/.local/share/tuido/tasks.json`
- **Refined TUI** - Clean, intentional design with a cohesive color palette

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
| `1-9` | Jump to task # |
| `←` / `→` or `h` / `l` | Previous/next day |
| `t` | Jump to today |
| `g` | Go to specific date |

### Tasks
| Key | Action |
|-----|--------|
| `a` | Add new task |
| `e` | Edit selected task |
| `x` / `Enter` / `Space` | Toggle done |
| `d` | Delete task |
| `>` / `n` | Move task to tomorrow |
| `J` / `K` or `Shift+↑/↓` | Reorder task |

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
