// Package tests contains integration tests for the scribe application.
// These tests exercise the App's public API end-to-end without access to
// internal package state.
package tests

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/melm/scribe/filetree"
	"github.com/melm/scribe/ui"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func makeTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range []string{"a.md", "b.md", "c.md", "d.txt"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("# "+f), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// newTestApp creates a fully-initialised App without a real terminal.
func newTestApp(t *testing.T, dir string) *ui.App {
	t.Helper()
	app := ui.New(dir, "", "dark")

	m, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app = m.(*ui.App)

	// Expand the root directly and deliver the sidebar-loaded message.
	root := &filetree.Node{Name: dir, Path: dir, IsDir: true}
	if err := filetree.ExpandNode(root); err != nil {
		t.Fatal(err)
	}
	m, _ = app.Update(ui.SidebarLoadedMsg{Nodes: filetree.Flatten(root)})
	app = m.(*ui.App)

	return app
}

func sendKey(t *testing.T, app *ui.App, key string) *ui.App {
	t.Helper()
	m, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	if updated, ok := m.(*ui.App); ok {
		return updated
	}
	return app
}

func sendSpecialKey(t *testing.T, app *ui.App, kt tea.KeyType) *ui.App {
	t.Helper()
	m, _ := app.Update(tea.KeyMsg{Type: kt})
	if updated, ok := m.(*ui.App); ok {
		return updated
	}
	return app
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestAppInitialisesInNormalMode(t *testing.T) {
	app := newTestApp(t, makeTestDir(t))
	if got := app.GetMode(); got != ui.ModeNormal {
		t.Errorf("want ModeNormal, got %v", got)
	}
}

func TestAppSidebarHasNodesAfterLoad(t *testing.T) {
	app := newTestApp(t, makeTestDir(t))
	if app.SidebarNodeCount() == 0 {
		t.Error("sidebar has no nodes after initialisation")
	}
}

func TestConsecutiveDownKeys(t *testing.T) {
	app := newTestApp(t, makeTestDir(t))
	total := app.SidebarNodeCount()
	if total < 3 {
		t.Fatalf("need ≥3 nodes, got %d", total)
	}
	for i := 1; i <= total-1; i++ {
		app = sendKey(t, app, "j")
		if got := app.SidebarCursor(); got != i {
			t.Errorf("press %d: want cursor=%d, got=%d", i, i, got)
		}
	}
}

func TestConsecutiveUpKeys(t *testing.T) {
	app := newTestApp(t, makeTestDir(t))
	total := app.SidebarNodeCount()
	// Manually position cursor at the last node via key presses.
	for i := 0; i < total-1; i++ {
		app = sendKey(t, app, "j")
	}
	for i := total - 2; i >= 0; i-- {
		app = sendKey(t, app, "k")
		if got := app.SidebarCursor(); got != i {
			t.Errorf("press: want cursor=%d, got=%d", i, got)
		}
	}
}

func TestInterleavedAsyncMessages(t *testing.T) {
	app := newTestApp(t, makeTestDir(t))

	app = sendKey(t, app, "j")
	if got := app.SidebarCursor(); got != 1 {
		t.Fatalf("after first j: want cursor=1, got=%d", got)
	}

	// An async renderer message must not disturb the sidebar cursor.
	m, _ := app.Update(ui.RendererReadyMsg{Renderer: nil, Width: app.Width()})
	app = m.(*ui.App)
	if got := app.SidebarCursor(); got != 1 {
		t.Errorf("after async msg: cursor changed unexpectedly to %d", got)
	}

	app = sendKey(t, app, "j")
	if got := app.SidebarCursor(); got != 2 {
		t.Fatalf("after second j: want cursor=2, got=%d", got)
	}
}

func TestModeTransitions(t *testing.T) {
	dir := makeTestDir(t)
	file := filepath.Join(dir, "test.md")
	if err := os.WriteFile(file, []byte("# Test"), 0o644); err != nil {
		t.Fatal(err)
	}

	app := newTestApp(t, dir)

	m, _ := app.Update(ui.FileLoadedMsg{Path: file, Content: "# Test"})
	app = m.(*ui.App)
	if got := app.GetMode(); got != ui.ModePreview {
		t.Errorf("after FileLoadedMsg: want ModePreview, got %v", got)
	}

	app = sendKey(t, app, "e")
	if got := app.GetMode(); got != ui.ModeEdit {
		t.Errorf("after 'e': want ModeEdit, got %v", got)
	}

	app = sendSpecialKey(t, app, tea.KeyEsc)
	if got := app.GetMode(); got != ui.ModeNormal {
		t.Errorf("after Esc: want ModeNormal, got %v", got)
	}
	if !app.SidebarFocused() {
		t.Error("sidebar should be focused in ModeNormal")
	}
}

func TestSaveMarksFileUnmodified(t *testing.T) {
	dir := makeTestDir(t)
	file := filepath.Join(dir, "save_test.md")
	if err := os.WriteFile(file, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	app := newTestApp(t, dir)
	m, _ := app.Update(ui.FileLoadedMsg{Path: file, Content: "hello"})
	app = m.(*ui.App)

	app = sendKey(t, app, "e")
	m, _ = app.Update(ui.SavedMsg{Path: file})
	app = m.(*ui.App)

	if app.EditorModified() {
		t.Error("editor should not be modified after SavedMsg")
	}
}
