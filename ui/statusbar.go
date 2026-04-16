package ui

import (
	"path/filepath"
	"strings"
)

// StatusBar renders the bottom status line.
type StatusBar struct {
	width    int
	message  string // transient message (saved, error, etc.)
	showHelp bool   // toggled by '?'
}

func newStatusBar() *StatusBar {
	return &StatusBar{}
}

func (s *StatusBar) setWidth(w int) {
	s.width = w
}

func (s *StatusBar) SetMessage(msg string) {
	s.message = msg
}

func (s *StatusBar) ClearMessage() {
	s.message = ""
}

func (s *StatusBar) ToggleHelp() {
	s.showHelp = !s.showHelp
}

// View renders the status bar.
func (s *StatusBar) View(mode Mode, modified bool, filePath string) string {
	// Transient message (save confirmation, error, etc.) takes priority.
	if s.message != "" {
		return StyleStatusBar.Width(s.width).Render(s.message)
	}

	if s.showHelp {
		return s.helpView(mode, modified)
	}
	return s.defaultView(mode, modified, filePath)
}

// defaultView: minimal one-liner — mode badge · file · [?] help
func (s *StatusBar) defaultView(mode Mode, modified bool, filePath string) string {
	// Right side: file name (+ modified marker)
	right := ""
	if filePath != "" {
		name := filepath.Base(filePath)
		if modified {
			name += " " + StyleStatusModified.Render("●")
		}
		right = StyleStatusKey.Render(" " + name + " ")
	}

	helpHint := hint("?", "help")

	leftWidth := s.width - visibleWidth(right)
	if leftWidth < 0 {
		leftWidth = 0
	}
	left := StyleStatusBar.Width(leftWidth).Render(helpHint)
	return strings.TrimSuffix(left+right, "\n")
}

// helpView: full keybind reference for the current mode.
func (s *StatusBar) helpView(mode Mode, modified bool) string {
	var hints string
	switch mode {
	case ModeNormal:
		hints = hint("↑↓", "navigate") + "  " +
			hint("Enter", "open/expand") + "  " +
			hint("e", "edit") + "  " +
			hint("p", "preview") + "  " +
			hint("n", "new file") + "  " +
			hint("d", "delete") + "  " +
			hint("r", "rename") + "  " +
			hint("q", "quit") + "  " +
			hint("?", "close help")
	case ModeEdit:
		modTag := ""
		if modified {
			modTag = "  " + StyleStatusModified.Render("[modified]")
		}
		hints = hint("Ctrl+S", "save") + "  " +
			hint("Ctrl+P", "preview") + "  " +
			hint("Esc", "sidebar") + "  " +
			hint("?", "close help") + modTag
	case ModePreview:
		hints = hint("↑↓/PgUp/PgDn", "scroll") + "  " +
			hint("e", "edit") + "  " +
			hint("Esc", "sidebar") + "  " +
			hint("?", "close help")
	}

	return StyleStatusBar.Width(s.width).Render(hints)
}

func hint(key, label string) string {
	return StyleStatusKey.Render("["+key+"]") + " " +
		StyleStatusBar.Render(label)
}

// visibleWidth returns the printable character count of a string,
// stripping ANSI escape sequences.
func visibleWidth(s string) int {
	inEsc := false
	w := 0
	for _, r := range s {
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		if r == '\x1b' {
			inEsc = true
			continue
		}
		w++
	}
	return w
}
