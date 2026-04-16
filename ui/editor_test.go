package ui

import (
	"testing"
)

func TestEditorStartsAtTopOnLoad(t *testing.T) {
	e := newEditor()
	e.setSize(80, 24)
	e.LoadContent("/tmp/test.md", "line0\nline1\nline2\nline3\nline4")
	if got := e.CursorLine(); got != 0 {
		t.Fatalf("expected cursor at line 0 after LoadContent, got line %d", got)
	}
}

func TestScrollToLine(t *testing.T) {
	e := newEditor()
	e.setSize(80, 24)
	e.LoadContent("/tmp/test.md", "line0\nline1\nline2\nline3\nline4")
	e.ScrollToLine(3)
	if got := e.CursorLine(); got != 3 {
		t.Fatalf("expected cursor at line 3 after ScrollToLine(3), got line %d", got)
	}
}
