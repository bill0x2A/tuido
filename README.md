# tuido

A minimal TUI todo app with day-based task management.

![tuido demo](https://img.shields.io/badge/go-%3E%3D1.21-blue)

## Features

- **Day-based tasks** - Navigate between days with arrow keys or `h`/`l`
- **Quick jump** - Press `1-0` to jump to tasks 1-10, `!@#$%^&*()` for 11-20
- **Task rollover** - Incomplete tasks from past days prompt for rollover
- **Move between days** - Press `>` to move task to tomorrow, `<` for yesterday
- **Auto-sort** - Completed tasks automatically move to the bottom
- **Vim-style navigation** - Use `j`/`k` to move, `J`/`K` or `Shift+↑/↓` to reorder
- **Persistent storage** - Tasks saved to `~/.local/share/tuido/tasks.json`

## Installation

### From source (requires Go 1.21+)

```bash
go install github.com/bill0x2a/tuido@latest
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
| `1-9` | Jump to task 1-9 |
| `0` | Jump to task 10 |
| `!@#$%^&*()` | Jump to task 11-20 (shift+number) |
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
| `>` | Move task to tomorrow |
| `<` | Move task to yesterday |
| `J` / `K` or `Shift+↑/↓` | Reorder task up/down |

### General
| Key | Action |
|-----|--------|
| `?` | Show help |
| `q` / `Esc` | Quit |

## Data Storage

Tasks are stored in JSON format at:
- Linux/macOS: `~/.local/share/tuido/tasks.json`
- With `XDG_DATA_HOME` set: `$XDG_DATA_HOME/tuido/tasks.json`

## License

MIT
