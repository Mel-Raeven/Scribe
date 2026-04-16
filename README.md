# scribe

A fast, keyboard-driven terminal note-taking app with a markdown preview, file tree sidebar, and built-in editor.

## Features

- Markdown preview with syntax highlighting (Catppuccin Mocha / Latte themes)
- Side-by-side file tree and editor/preview pane
- Scrollable preview with scrollbar indicator and mouse wheel support
- Full in-app editor — create, rename, and delete files without leaving the TUI
- Supports `.md`, `.txt`, `.text`, `.rst`, `.log` files
- Instant startup — all I/O is asynchronous

## Install

Build from source:

```sh
git clone https://github.com/Mel-Raeven/Scribe
cd scribe
make install
```

## Usage

```sh
scribe              # open current directory
scribe ~/notes      # open a specific directory
scribe README.md    # open a specific file directly
```

## Keybindings

### Normal mode (sidebar focused)

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` / `Space` | Open file / expand directory |
| `e` | Edit open file |
| `p` | Preview open file |
| `Tab` | Switch to preview |
| `n` | New file |
| `r` | Rename open file |
| `d` | Delete open file |
| `q` | Quit |
| `?` | Toggle help |

### Preview mode

| Key | Action |
|-----|--------|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `PgUp` / `Ctrl+B` | Page up |
| `PgDn` / `Ctrl+F` | Page down |
| Mouse wheel | Scroll up / down |
| `e` | Switch to editor |
| `Esc` | Back to sidebar |
| `?` | Toggle help |

### Edit mode

| Key | Action |
|-----|--------|
| `Ctrl+S` | Save |
| `Ctrl+P` | Switch to preview |
| `Esc` | Back to sidebar |
| `?` | Toggle help |

## License

MIT
