package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/melm/scribe/filetree"
)

func makeSidebarDir(t *testing.T) *Sidebar {
	t.Helper()
	dir := t.TempDir()
	for _, f := range []string{"a.md", "b.md", "c.md", "d.txt"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("# "+f), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	sb := newSidebar(dir)
	sb.setSize(28, 20)
	if err := filetree.ExpandNode(sb.root); err != nil {
		t.Fatal(err)
	}
	sb.applyLoad(SidebarLoadedMsg{Nodes: filetree.Flatten(sb.root)})
	return sb
}

func TestSidebarLoadPopulatesNodes(t *testing.T) {
	sb := makeSidebarDir(t)
	if len(sb.nodes) == 0 {
		t.Fatal("expected sidebar nodes after load, got none")
	}
}

func TestSidebarCursorMoveDown(t *testing.T) {
	sb := makeSidebarDir(t)
	total := len(sb.nodes)
	if total < 2 {
		t.Fatalf("need at least 2 nodes, got %d", total)
	}
	for i := 1; i < total; i++ {
		sb.MoveDown()
		if sb.cursor != i {
			t.Errorf("after MoveDown #%d: want cursor=%d, got %d", i, i, sb.cursor)
		}
	}
}

func TestSidebarCursorMoveUp(t *testing.T) {
	sb := makeSidebarDir(t)
	total := len(sb.nodes)
	sb.cursor = total - 1
	for i := total - 2; i >= 0; i-- {
		sb.MoveUp()
		if sb.cursor != i {
			t.Errorf("after MoveUp: want cursor=%d, got %d", i, sb.cursor)
		}
	}
}

func TestSidebarCursorDoesNotGoOutOfBounds(t *testing.T) {
	sb := makeSidebarDir(t)

	for i := 0; i < 10; i++ {
		sb.MoveUp()
	}
	if sb.cursor != 0 {
		t.Errorf("cursor below 0: got %d", sb.cursor)
	}

	last := len(sb.nodes) - 1
	for i := 0; i < len(sb.nodes)+5; i++ {
		sb.MoveDown()
	}
	if sb.cursor != last {
		t.Errorf("cursor beyond last node: want %d, got %d", last, sb.cursor)
	}
}
