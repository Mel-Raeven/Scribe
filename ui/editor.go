package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// SavedMsg is sent when a file is successfully saved.
type SavedMsg struct{ Path string }

// SaveErrMsg is sent when saving fails.
type SaveErrMsg struct{ Err error }

// Editor wraps a bubbles/textarea for editing file contents.
type Editor struct {
	textarea textarea.Model
	filePath string
	modified bool
	width    int
	height   int
}

func newEditor() *Editor {
	ta := textarea.New()
	ta.ShowLineNumbers = true
	ta.CharLimit = 0 // no limit
	ta.SetWidth(80)
	ta.SetHeight(24)

	// Styling
	ta.FocusedStyle.Base = StyleMainPane
	ta.BlurredStyle.Base = StyleMainPane
	ta.Prompt = ""

	return &Editor{textarea: ta}
}

func (e *Editor) setSize(w, h int) {
	e.width = w
	e.height = h
	e.textarea.SetWidth(w)
	e.textarea.SetHeight(h)
}

// LoadFile reads a file from disk into the textarea and positions the cursor
// at the top of the document.
func (e *Editor) LoadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	content := string(data)
	e.filePath = path
	e.modified = false
	e.textarea.SetValue(content)
	e.gotoTop()
	return content, nil
}

// LoadContent loads already-read content into the editor without hitting disk.
func (e *Editor) LoadContent(path, content string) {
	e.filePath = path
	e.modified = false
	e.textarea.SetValue(content)
	e.gotoTop()
}

// ScrollToLine moves the editor cursor to the given line (0-indexed).
func (e *Editor) ScrollToLine(line int) {
	e.gotoTop()
	wasFocused := e.textarea.Focused()
	if !wasFocused {
		e.textarea.Focus() //nolint — cmd is a cursor blink, safe to discard
	}
	for i := 0; i < line; i++ {
		e.textarea, _ = e.textarea.Update(tea.KeyMsg{Type: tea.KeyDown})
	}
	if !wasFocused {
		e.textarea.Blur()
	}
}

// gotoTop moves the textarea cursor to row 0, col 0.
// SetValue leaves the cursor at the end; this corrects it.
func (e *Editor) gotoTop() {
	wasFocused := e.textarea.Focused()
	if !wasFocused {
		e.textarea.Focus() //nolint — cmd is a cursor blink, safe to discard
	}
	e.textarea, _ = e.textarea.Update(tea.KeyMsg{Type: tea.KeyCtrlHome})
	if !wasFocused {
		e.textarea.Blur()
	}
}

// Save writes the current content to disk.
func (e *Editor) Save() tea.Cmd {
	if e.filePath == "" {
		return nil
	}
	content := e.textarea.Value()
	return func() tea.Msg {
		if err := os.WriteFile(e.filePath, []byte(content), 0o644); err != nil {
			return SaveErrMsg{Err: err}
		}
		return SavedMsg{Path: e.filePath}
	}
}

// Content returns the current editor text.
func (e *Editor) Content() string {
	return e.textarea.Value()
}

// Modified returns whether the buffer has unsaved changes.
func (e *Editor) Modified() bool {
	return e.modified
}

// FilePath returns the currently open file path.
func (e *Editor) FilePath() string {
	return e.filePath
}

// Focus gives the textarea keyboard focus.
func (e *Editor) Focus() tea.Cmd {
	return e.textarea.Focus()
}

// Blur removes keyboard focus.
func (e *Editor) Blur() {
	e.textarea.Blur()
}

// OnSaved resets the modified flag.
func (e *Editor) OnSaved() {
	e.modified = false
}

// CursorLine returns the current cursor line (0-indexed).
func (e *Editor) CursorLine() int {
	return e.textarea.Line()
}
func (e *Editor) Update(msg tea.Msg) tea.Cmd {
	prev := e.textarea.Value()
	var cmd tea.Cmd
	e.textarea, cmd = e.textarea.Update(msg)
	if e.textarea.Value() != prev {
		e.modified = true
	}
	return cmd
}

// View renders the editor.
func (e *Editor) View() string {
	if e.filePath == "" {
		return e.emptyView()
	}
	return e.textarea.View()
}

func (e *Editor) emptyView() string {
	lines := []string{}
	msg := "No file open — select a file from the sidebar"
	padTop := e.height / 2
	for i := 0; i < padTop; i++ {
		lines = append(lines, "")
	}
	lines = append(lines, StyleStatusKey.Render(msg))
	return strings.Join(lines, "\n")
}
