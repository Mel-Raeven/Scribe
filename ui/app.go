package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Mode represents the current application mode.
type Mode int

const (
	ModeNormal  Mode = iota // sidebar focused
	ModeEdit                // editor active
	ModePreview             // preview active
	ModePrompt              // new file / delete prompt
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeEdit:
		return "EDIT"
	case ModePreview:
		return "PREVIEW"
	case ModePrompt:
		return "PROMPT"
	}
	return ""
}

const sidebarWidth = 28

// clearMsgCmd clears the status bar message after a short delay.
func clearMsgCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return clearMsgMsg{}
	})
}

type clearMsgMsg struct{}

// promptMode distinguishes what the prompt is being used for.
type promptMode int

const (
	promptNew    promptMode = iota
	promptDelete
	promptRename
)

// App is the root Bubble Tea model.
type App struct {
	root      string
	startFile string // file to open on first render, cleared after use
	mode      Mode
	sidebar   *Sidebar
	editor    *Editor
	preview   *Preview
	statusbar *StatusBar

	width  int
	height int

	// Prompt state
	prompt     textinput.Model
	promptMode promptMode
	promptMsg  string

	openFile      string // currently open file path
	sidebarHidden bool   // true while a file is open and sidebar is dismissed

	// Nav-key rate limiter for edit mode (prevents ghost scrolling after
	// releasing a held key). pendingNavKey holds the most-recent nav key
	// received while the tick loop is running; only one move is applied
	// per tick, so queued events overwrite each other instead of stacking.
	pendingNavKey string
	navTickActive bool
}

// New creates a new App model. No I/O is performed here — everything
// happens asynchronously in Init so the TUI opens immediately.
func New(root, startFile, glamourStyle string) *App {
	InitTheme(glamourStyle == "dark")
	sb := newSidebar(root)
	sb.setFocused(true)

	ti := textinput.New()
	ti.CharLimit = 256

	return &App{
		root:      root,
		startFile: startFile,
		mode:      ModeNormal,
		sidebar:   sb,
		editor:    newEditor(),
		preview:   newPreview(glamourStyle, root),
		statusbar: newStatusBar(),
		prompt:    ti,
	}
}

// Init implements tea.Model. Kicks off async I/O so the first frame renders
// immediately with no blocking.
func (a *App) Init() tea.Cmd {
	cmds := []tea.Cmd{
		a.sidebar.loadRootCmd(), // expand root dir in a goroutine
	}
	if a.startFile != "" {
		cmds = append(cmds, a.loadFileCmd(a.startFile))
	}
	return tea.Batch(cmds...)
}

// loadFileCmd reads a file in a goroutine and sends FileLoadedMsg.
type FileLoadedMsg struct {
	Path    string
	Content string
}

func (a *App) loadFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		data, _ := os.ReadFile(path)
		return FileLoadedMsg{Path: path, Content: string(data)}
	}
}

// Update implements tea.Model.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		cmds = append(cmds, a.recalcSizes()...)

	case RendererReadyMsg:
		a.preview.applyRenderer(msg)

	case clearMsgMsg:
		a.statusbar.ClearMessage()

	// Async sidebar load completed
	case SidebarLoadedMsg:
		a.sidebar.applyLoad(msg)

	// Async file load completed (from Init startFile)
	case FileLoadedMsg:
		cmds = append(cmds, a.applyFileLoad(msg.Path, msg.Content)...)

	case SavedMsg:
		a.editor.OnSaved()
		a.statusbar.SetMessage(StyleStatusSaved.Render(fmt.Sprintf("  saved: %s", filepath.Base(msg.Path))))
		a.preview.SetContent(msg.Path, a.editor.Content())
		cmds = append(cmds, clearMsgCmd())

	case SaveErrMsg:
		a.statusbar.SetMessage(StyleStatusError.Render(fmt.Sprintf("  error saving: %v", msg.Err)))
		cmds = append(cmds, clearMsgCmd())

	case navTickMsg:
		if a.mode == ModeEdit && a.pendingNavKey != "" {
			k := a.pendingNavKey
			a.pendingNavKey = ""
			if km, ok := navKeyMsg(k); ok {
				cmds = append(cmds, a.editor.Update(km), navTickCmd())
			}
		} else {
			a.navTickActive = false
		}

	case tea.KeyMsg:
		cmd := a.handleKey(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case tea.MouseMsg:
		if cmd := a.handleMouse(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleKey(msg tea.KeyMsg) tea.Cmd {
	if a.mode == ModePrompt {
		return a.handlePromptKey(msg)
	}
	if msg.String() == "ctrl+c" {
		return tea.Quit
	}
	// '?' toggles help in all non-prompt, non-edit modes; also in edit mode.
	if msg.String() == "?" && a.mode != ModeEdit {
		a.statusbar.ToggleHelp()
		return nil
	}
	switch a.mode {
	case ModeNormal:
		return a.handleNormalKey(msg)
	case ModeEdit:
		return a.handleEditKey(msg)
	case ModePreview:
		return a.handlePreviewKey(msg)
	}
	return nil
}

func (a *App) handleMouse(msg tea.MouseMsg) tea.Cmd {
	// Only handle wheel events; ignore clicks/motion.
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if a.openFile != "" {
			a.preview.ScrollUp()
		}
	case tea.MouseButtonWheelDown:
		if a.openFile != "" {
			a.preview.ScrollDown()
		}
	}
	return nil
}

func (a *App) handleNormalKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q", "Q":
		return tea.Quit
	case "up", "k":
		a.sidebar.MoveUp()
	case "down", "j":
		a.sidebar.MoveDown()
	case "enter", " ":
		path, isFile := a.sidebar.Toggle()
		if isFile {
			return a.openFileCmd(path)
		}
	case "tab":
		if a.openFile != "" {
			a.mode = ModePreview
			a.sidebar.setFocused(false)
		}
	case "e":
		if a.openFile != "" {
			a.enterEditMode()
			return a.editor.Focus()
		}
	case "p":
		if a.openFile != "" {
			a.mode = ModePreview
			a.sidebar.setFocused(false)
		}
	case "n":
		return a.startPrompt(promptNew, "New file name: ")
	case "d":
		if a.openFile != "" {
			return a.startPrompt(promptDelete, fmt.Sprintf("Delete %s? (y/N): ", filepath.Base(a.openFile)))
		}
	case "r":
		if a.openFile != "" {
			return a.startPrompt(promptRename, "Rename to: ")
		}
	}
	return nil
}

func (a *App) handleEditKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+s":
		return a.editor.Save()
	case "?":
		a.statusbar.ToggleHelp()
		return nil
	case "esc", "ctrl+p":
		a.editor.Blur()
		a.navTickActive = false
		a.pendingNavKey = ""
		if msg.String() == "ctrl+p" && a.openFile != "" {
			a.mode = ModePreview
		} else {
			return tea.Batch(a.showSidebar()...)
		}
		return nil
	default:
		if isNavKey(msg.String()) {
			if !a.navTickActive {
				// First key of a sequence: apply immediately and start tick loop.
				a.navTickActive = true
				a.pendingNavKey = ""
				return tea.Batch(a.editor.Update(msg), navTickCmd())
			}
			// Key is still held: record it; the tick will apply one move.
			a.pendingNavKey = msg.String()
			return nil
		}
		return a.editor.Update(msg)
	}
}

// navTickMsg is the rate-limiter tick for edit-mode navigation.
type navTickMsg struct{}

func navTickCmd() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(time.Time) tea.Msg {
		return navTickMsg{}
	})
}

// isNavKey reports whether a key string is a cursor-movement key.
func isNavKey(k string) bool {
	switch k {
	case "up", "down", "left", "right",
		"pgup", "pgdown", "home", "end",
		"ctrl+home", "ctrl+end":
		return true
	}
	return false
}

// navKeyMsg constructs the tea.KeyMsg for a nav key string.
func navKeyMsg(k string) (tea.KeyMsg, bool) {
	m := map[string]tea.KeyType{
		"up":       tea.KeyUp,
		"down":     tea.KeyDown,
		"left":     tea.KeyLeft,
		"right":    tea.KeyRight,
		"pgup":     tea.KeyPgUp,
		"pgdown":   tea.KeyPgDown,
		"home":     tea.KeyHome,
		"end":      tea.KeyEnd,
	}
	if t, ok := m[k]; ok {
		return tea.KeyMsg{Type: t}, true
	}
	return tea.KeyMsg{}, false
}

func (a *App) handlePreviewKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		return tea.Batch(a.showSidebar()...)
	case "e":
		a.enterEditMode()
		return a.editor.Focus()
	case "up", "k":
		a.preview.ScrollUp()
	case "down", "j":
		a.preview.ScrollDown()
	case "pgup", "ctrl+b":
		a.preview.PageUp()
	case "pgdown", "ctrl+f":
		a.preview.PageDown()
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handlePromptKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		a.mode = ModeNormal
		a.sidebar.setFocused(true)
		return nil
	case "enter":
		return a.submitPrompt()
	default:
		var cmd tea.Cmd
		a.prompt, cmd = a.prompt.Update(msg)
		return cmd
	}
}

func (a *App) startPrompt(pm promptMode, label string) tea.Cmd {
	a.promptMode = pm
	a.promptMsg = label
	a.prompt.Reset()
	a.prompt.Placeholder = ""
	a.mode = ModePrompt
	a.sidebar.setFocused(false)
	return a.prompt.Focus()
}

func (a *App) submitPrompt() tea.Cmd {
	val := strings.TrimSpace(a.prompt.Value())
	a.mode = ModeNormal
	a.sidebar.setFocused(true)
	a.prompt.Blur()

	switch a.promptMode {
	case promptNew:
		if val == "" {
			return nil
		}
		newPath := filepath.Join(a.root, val)
		if filepath.Ext(newPath) == "" {
			newPath += ".md"
		}
		_ = os.MkdirAll(filepath.Dir(newPath), 0o755)
		f, err := os.Create(newPath)
		if err != nil {
			a.statusbar.SetMessage(StyleStatusError.Render("  error: " + err.Error()))
			return clearMsgCmd()
		}
		f.Close()
		a.sidebar.Refresh()
		return a.openFileCmd(newPath)

	case promptDelete:
		if strings.ToLower(val) != "y" {
			return nil
		}
		if err := os.Remove(a.openFile); err != nil {
			a.statusbar.SetMessage(StyleStatusError.Render("  error: " + err.Error()))
			return clearMsgCmd()
		}
		a.openFile = ""
		a.sidebar.Refresh()
		a.statusbar.SetMessage(StyleStatusSaved.Render("  deleted"))
		return clearMsgCmd()

	case promptRename:
		if val == "" {
			return nil
		}
		newPath := filepath.Join(filepath.Dir(a.openFile), val)
		if filepath.Ext(newPath) == "" {
			newPath += filepath.Ext(a.openFile)
		}
		if err := os.Rename(a.openFile, newPath); err != nil {
			a.statusbar.SetMessage(StyleStatusError.Render("  error: " + err.Error()))
			return clearMsgCmd()
		}
		a.sidebar.Refresh()
		return a.openFileCmd(newPath)
	}
	return nil
}

// openFileCmd reads a file in a goroutine and applies it on arrival.
func (a *App) openFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		data, _ := os.ReadFile(path)
		return FileLoadedMsg{Path: path, Content: string(data)}
	}
}

// applyFileLoad wires the loaded content into the editor and preview.
// The file is read exactly once; both components receive the same string.
func (a *App) applyFileLoad(path, content string) []tea.Cmd {
	a.openFile = path
	// Feed content directly — no second ReadFile
	a.editor.LoadContent(path, content)
	a.preview.SetContent(path, content)
	a.mode = ModePreview
	a.sidebar.setFocused(false)
	a.sidebarHidden = true
	a.startFile = "" // clear one-time start file
	// Recalc sizes so preview gets full width immediately.
	// Return the cmds so the renderer build cmd is not dropped.
	return a.recalcSizes()
}

// showSidebar reveals the sidebar and returns to normal mode.
func (a *App) showSidebar() []tea.Cmd {
	a.sidebarHidden = false
	a.mode = ModeNormal
	a.sidebar.setFocused(true)
	return a.recalcSizes()
}

func (a *App) enterEditMode() {
	if a.openFile == "" {
		return
	}
	a.mode = ModeEdit
	a.sidebar.setFocused(false)
	a.editor.ScrollToLine(a.preview.ScrollY())
}

func (a *App) recalcSizes() []tea.Cmd {
	contentH := a.height - 2 // header + statusbar
	if contentH < 1 {
		contentH = 1
	}

	var mainW int
	if a.sidebarHidden {
		mainW = a.width
	} else {
		mainW = a.width - sidebarWidth - 1 // -1 for border
	}
	if mainW < 1 {
		mainW = 1
	}

	a.sidebar.setSize(sidebarWidth, contentH)
	a.editor.setSize(mainW, contentH)
	a.statusbar.setWidth(a.width)

	var cmds []tea.Cmd
	if cmd := a.preview.setSize(mainW, contentH); cmd != nil {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// View implements tea.Model.
func (a *App) View() string {
	if a.width == 0 {
		return ""
	}

	header := a.renderHeader()
	body := a.renderBody()
	status := a.renderStatus()

	return lipgloss.JoinVertical(lipgloss.Left, header, body, status)
}

func (a *App) renderHeader() string {
	path := a.root
	if home, err := os.UserHomeDir(); err == nil {
		if rel, err := filepath.Rel(home, path); err == nil && !strings.HasPrefix(rel, "..") {
			path = "~/" + rel
		}
	}

	modeStr := StyleHeaderMode.Render(a.mode.String())
	pathStr := StyleHeaderPath.Render("  scribe  " + path)

	pad := a.width - lipgloss.Width(pathStr) - lipgloss.Width(modeStr) - 2
	if pad < 0 {
		pad = 0
	}
	spacer := strings.Repeat(" ", pad)

	return StyleHeader.Width(a.width).Render(pathStr + spacer + modeStr)
}

func (a *App) renderBody() string {
	if a.sidebarHidden {
		return a.renderMain()
	}
	sidebarView := a.sidebar.View()
	mainView := a.renderMain()
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, mainView)
}

func (a *App) renderMain() string {
	var mainW int
	if a.sidebarHidden {
		mainW = a.width
	} else {
		mainW = a.width - sidebarWidth - 1
	}
	contentH := a.height - 2
	if mainW < 1 {
		mainW = 1
	}
	if contentH < 1 {
		contentH = 1
	}

	switch a.mode {
	case ModeEdit:
		return a.editor.View()
	case ModePreview:
		return StyleMainPane.Width(mainW).Height(contentH).Render(a.preview.View())
	case ModePrompt:
		return a.renderPromptInMain(mainW, contentH)
	default: // ModeNormal — show preview of last open file if any
		if a.openFile != "" {
			return StyleMainPane.Width(mainW).Height(contentH).Render(a.preview.View())
		}
		return a.renderEmpty(mainW, contentH)
	}
}

func (a *App) renderPromptInMain(w, h int) string {
	lines := make([]string, h)
	mid := h / 2

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(StyleHeaderMode.GetForeground()).
		Padding(0, 2).
		Width(w - 6)

	inner := a.promptMsg + "\n" + a.prompt.View()
	rendered := box.Render(inner)

	boxH := strings.Count(rendered, "\n") + 1
	startLine := mid - boxH/2
	if startLine < 0 {
		startLine = 0
	}

	for i := range lines {
		if i == startLine {
			lines[i] = rendered
		} else {
			lines[i] = ""
		}
	}
	return strings.Join(lines, "\n")
}

func (a *App) renderEmpty(w, h int) string {
	lines := make([]string, h)
	mid := h / 2
	msg := StyleStatusKey.Render("Select a file from the sidebar to get started")
	lines[mid] = msg
	return strings.Join(lines, "\n")
}

func (a *App) renderStatus() string {
	shortPath := ""
	if a.openFile != "" {
		shortPath = relPath(a.root, a.openFile)
	}
	return a.statusbar.View(a.mode, a.editor.Modified(), shortPath)
}

func relPath(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return filepath.Base(target)
	}
	return rel
}

// ---------------------------------------------------------------------------
// Test accessors — thin getters that expose internal state to external tests.
// ---------------------------------------------------------------------------

// GetMode returns the current application mode.
func (a *App) GetMode() Mode { return a.mode }

// SidebarCursor returns the sidebar's current cursor position.
func (a *App) SidebarCursor() int { return a.sidebar.cursor }

// SidebarNodeCount returns the number of visible nodes in the sidebar.
func (a *App) SidebarNodeCount() int { return len(a.sidebar.nodes) }

// SidebarFocused reports whether the sidebar currently has focus.
func (a *App) SidebarFocused() bool { return a.sidebar.focused }

// EditorModified reports whether the editor has unsaved changes.
func (a *App) EditorModified() bool { return a.editor.Modified() }

// Width returns the current terminal width known to the app.
func (a *App) Width() int { return a.width }
