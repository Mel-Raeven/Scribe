package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/melm/scribe/filetree"
)

// Sidebar manages the file tree panel.
type Sidebar struct {
	root    *filetree.Node
	nodes   []*filetree.Node
	cursor  int
	width   int
	height  int
	focused bool
}

// newSidebar returns an empty sidebar immediately (no I/O).
// Call loadRoot() from a goroutine / tea.Cmd to populate it.
func newSidebar(rootPath string) *Sidebar {
	return &Sidebar{
		root: &filetree.Node{
			Name:  rootPath,
			Path:  rootPath,
			IsDir: true,
			// Expanded starts false so ExpandNode actually reads the directory.
			// It sets Expanded = true after the ReadDir completes.
		},
	}
}

// SidebarLoadedMsg is sent when the async root expand finishes.
type SidebarLoadedMsg struct {
	Nodes []*filetree.Node
	Err   error
}

// loadRootCmd expands the root directory in a goroutine and sends SidebarLoadedMsg.
func (s *Sidebar) loadRootCmd() func() tea.Msg {
	root := s.root
	return func() tea.Msg {
		if err := filetree.ExpandNode(root); err != nil {
			return SidebarLoadedMsg{Err: err}
		}
		return SidebarLoadedMsg{Nodes: filetree.Flatten(root)}
	}
}

// applyLoad applies the result of an async load.
func (s *Sidebar) applyLoad(msg SidebarLoadedMsg) {
	if msg.Err != nil || msg.Nodes == nil {
		return
	}
	s.nodes = msg.Nodes
}

func (s *Sidebar) setSize(w, h int) {
	s.width = w
	s.height = h
}

func (s *Sidebar) setFocused(f bool) {
	s.focused = f
}

// SelectedNode returns the currently highlighted node, or nil.
func (s *Sidebar) SelectedNode() *filetree.Node {
	if len(s.nodes) == 0 {
		return nil
	}
	return s.nodes[s.cursor]
}

// MoveUp moves the cursor up.
func (s *Sidebar) MoveUp() {
	if s.cursor > 0 {
		s.cursor--
	}
}

// MoveDown moves the cursor down.
func (s *Sidebar) MoveDown() {
	if s.cursor < len(s.nodes)-1 {
		s.cursor++
	}
}

// Toggle expands or collapses a directory, or returns the path of a file to open.
// Returns (filePath, wasFile).
func (s *Sidebar) Toggle() (string, bool) {
	node := s.SelectedNode()
	if node == nil {
		return "", false
	}
	if !node.IsDir {
		return node.Path, true
	}
	// Toggle expand/collapse
	if node.Expanded {
		node.Expanded = false
	} else {
		_ = filetree.ExpandNode(node)
	}
	s.refresh()
	return "", false
}

// refresh re-flattens the tree after expand/collapse changes.
func (s *Sidebar) refresh() {
	currentPath := ""
	if node := s.SelectedNode(); node != nil {
		currentPath = node.Path
	}
	s.nodes = filetree.Flatten(s.root)
	// Restore cursor to same path if possible
	for i, n := range s.nodes {
		if n.Path == currentPath {
			s.cursor = i
			return
		}
	}
	if s.cursor >= len(s.nodes) {
		s.cursor = max(0, len(s.nodes)-1)
	}
}

// Refresh re-reads the filesystem (e.g. after creating/deleting a file).
func (s *Sidebar) Refresh() {
	s.root.Expanded = false
	_ = filetree.ExpandNode(s.root)
	// Re-expand previously expanded dirs - for simplicity just collapse all
	s.refresh()
}

// View renders the sidebar.
func (s *Sidebar) View() string {
	if s.width == 0 || s.height == 0 {
		return ""
	}

	var sb strings.Builder

	// Scroll window so cursor is always visible
	visibleHeight := s.height
	scrollOffset := 0
	if s.cursor >= visibleHeight {
		scrollOffset = s.cursor - visibleHeight + 1
	}

	rendered := 0
	for i, node := range s.nodes {
		if i < scrollOffset {
			continue
		}
		if rendered >= visibleHeight {
			break
		}

		line := s.renderNode(node, i == s.cursor)
		sb.WriteString(line)
		sb.WriteRune('\n')
		rendered++
	}

	// Pad remaining lines
	for rendered < visibleHeight {
		sb.WriteString(strings.Repeat(" ", s.width))
		sb.WriteRune('\n')
		rendered++
	}

	content := sb.String()
	// Trim trailing newline — lipgloss Join handles spacing
	content = strings.TrimSuffix(content, "\n")

	return StyleSidebar.Width(s.width).Height(s.height).Render(content)
}

func (s *Sidebar) renderNode(node *filetree.Node, selected bool) string {
	// Indentation: 2 spaces per depth level (depth 1 = top level children)
	indent := strings.Repeat("  ", node.Depth-1)

	var prefix string
	var nameStyle lipgloss.Style

	if node.IsDir {
		if node.Expanded {
			prefix = "▾ "
		} else {
			prefix = "▸ "
		}
		if selected && s.focused {
			nameStyle = StyleTreeDirSelected
		} else {
			nameStyle = StyleTreeDir
		}
	} else {
		prefix = "  "
		if selected && s.focused {
			nameStyle = StyleTreeItemSelected
		} else {
			nameStyle = StyleTreeItem
		}
	}

	name := node.Name
	// Truncate long names to fit sidebar width
	maxNameLen := s.width - len(indent) - len(prefix) - 2
	if maxNameLen < 1 {
		maxNameLen = 1
	}
	if len(name) > maxNameLen {
		name = name[:maxNameLen-1] + "…"
	}

	raw := fmt.Sprintf("%s%s%s", indent, prefix, name)

	if selected && s.focused {
		// Pad to full width for highlight bar
		padded := raw + strings.Repeat(" ", max(0, s.width-len(raw)))
		return nameStyle.Render(padded)
	}
	return nameStyle.Render(raw)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
