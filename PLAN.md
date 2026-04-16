# note — Terminal Markdown Note-Taking App

## Overview

A fast, terminal-based note-taking app written in Go. Open it from any directory
like `note` or `note .` and immediately start editing markdown files with a live
rendered preview and a file-tree sidebar.

---

## Tech Stack

| Concern              | Library                                                              |
| -------------------- | -------------------------------------------------------------------- |
| TUI framework        | [Bubble Tea](https://github.com/charmbracelet/bubbletea)             |
| Styling / layout     | [Lipgloss](https://github.com/charmbracelet/lipgloss)                |
| Markdown rendering   | [Glamour](https://github.com/charmbracelet/glamour)                  |
| Markdown linting     | [markdownlint](https://github.com/DavidAnson/markdownlint) via a Go wrapper, or a lightweight custom rule set |
| File watching        | [fsnotify](https://github.com/fsnotify/fsnotify)                     |
| Config               | Plain TOML file at `~/.config/note/config.toml`                      |

> All dependencies compile into a single static binary — startup will be near-instant.

---

## CLI Interface

```
note              # open in current directory
note .            # same
note ~/notes      # open a specific directory
note README.md    # open a specific file directly
```

Install the binary to `$PATH` (e.g. `/usr/local/bin/note`) so it works from anywhere.

---

## Layout

```
┌─────────────────────────────────────────────────────┐
│  note  ~/projects/my-notes                [NORMAL]  │
├───────────────┬─────────────────────────────────────┤
│ Sidebar       │ Editor / Preview                    │
│               │                                     │
│ > notes/      │  # My Heading                       │
│   ├ ideas.md  │                                     │
│   ├ todo.md   │  Some rendered **markdown** here    │
│   └ drafts/   │  with syntax highlighting.          │
│     └ wip.md  │                                     │
│               │  - list item                        │
│               │  - another item                     │
│               │                                     │
├───────────────┴─────────────────────────────────────┤
│ [e] edit  [p] preview  [n] new  [d] delete  [?] help│
└─────────────────────────────────────────────────────┘
```

The layout has three regions:
1. **Header bar** — current path and mode indicator
2. **Sidebar** — file tree, navigable with arrow keys
3. **Main pane** — toggles between editor mode and rendered preview mode
4. **Status bar** — keybinding hints and lint warnings

---

## Modes

| Mode      | Description                                        |
| --------- | -------------------------------------------------- |
| `NORMAL`  | Navigate sidebar, browse files                     |
| `PREVIEW` | Read-only glamour-rendered markdown view           |
| `EDIT`    | Raw markdown editor (uses a Bubble Tea textarea)   |

Switching modes:
- `e` → EDIT mode
- `p` → PREVIEW mode
- `Esc` → back to NORMAL / sidebar focus

---

## Features

### Must-Have (v1)
- [ ] File tree sidebar with keyboard navigation (arrows, `Enter` to open, `Tab` to toggle focus)
- [ ] Markdown rendering via Glamour (PREVIEW mode)
- [ ] Basic text editor pane (EDIT mode) — wraps Bubble Tea's textarea component
- [ ] New file (`n`), delete file (`d`), rename file (`r`)
- [ ] Save on write (`Ctrl+S` or auto-save on mode switch)
- [ ] Single binary installable to `$PATH`

### Should-Have (v2)
- [ ] Markdown linting — highlight issues inline in EDIT mode (e.g. missing blank lines, broken links)
- [ ] Fuzzy file search (`/` or `Ctrl+P` to open a picker)
- [ ] Directory creation
- [ ] Configurable theme (dark/light) via `~/.config/note/config.toml`

### Nice-to-Have (v3)
- [ ] Split preview: edit on left, live rendered preview on right (updates as you type)
- [ ] Tag/frontmatter support
- [ ] Backlinks panel (which other notes link to this one)
- [ ] Export to HTML or PDF via pandoc (if installed)

---

## File Structure

```
note/
├── main.go
├── go.mod
├── go.sum
├── cmd/
│   └── root.go          # CLI entrypoint (cobra or manual flag parsing)
├── ui/
│   ├── app.go           # top-level Bubble Tea model
│   ├── sidebar.go       # file tree component
│   ├── editor.go        # editor pane component
│   ├── preview.go       # glamour preview pane component
│   ├── statusbar.go     # bottom status bar
│   └── styles.go        # lipgloss styles / theme
├── fs/
│   ├── tree.go          # directory reading and file tree model
│   └── watcher.go       # fsnotify watcher for live reload
├── lint/
│   └── lint.go          # markdown lint rules
└── config/
    └── config.go        # config loading from ~/.config/note/config.toml
```

---

## Startup Performance Goals

- Cold start to interactive UI: **< 100ms** (realistic for a compiled Go binary)
- File tree for a directory with 1000 files: rendered **< 50ms**
- Markdown render for a 500-line file: **< 30ms**

Achieved by:
- No JIT / interpreter overhead — single compiled binary
- Lazy loading of file contents (only read a file when it is opened)
- Glamour renders once on open; re-renders only on save or mode switch

---

## Installation

```sh
# Build and install
go build -o note .
sudo mv note /usr/local/bin/

# Or via go install (once published)
go install github.com/yourname/note@latest
```

---

## Open Questions (tweak these)

1. **Editor depth** — use a simple Bubble Tea textarea, or embed a more capable editor (e.g. a minimal vim-mode)?
  
Right now normal text is fina, but ideally this could be added in the future

2. **Linting strictness** — which markdown rules matter most to you (heading levels, line length, blank lines around headings, etc.)?

Skip linting for now

3. **Auto-save** — save continuously while typing, or only on `Ctrl+S`?

Only ctrl+s

4. **File types** — only `.md`, or also `.txt` and others?

I want it to be able to support .md, .txt, and other basic text stuff

5. **App name / binary name** — `note` conflicts with the POSIX `note` command on some systems; alternatives: `nts`, `mdn`, `jot`
app name could be scribe? maybe ill change this
