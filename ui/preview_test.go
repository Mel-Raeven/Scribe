package ui

import (
	"os"
	"testing"
)

// TestRendererDoesNotReadStdin verifies that buildRendererCmd can run without
// access to stdin — it uses a JSON style instead of WithAutoStyle() to avoid
// terminal detection that races with Bubble Tea's keyboard reader.
func TestRendererDoesNotReadStdin(t *testing.T) {
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()

	origStdin := os.Stdin
	os.Stdin = devNull
	defer func() { os.Stdin = origStdin }()

	cmd := buildRendererCmd(80, "dark")
	msg := cmd()

	if msg == nil {
		t.Fatal("buildRendererCmd returned nil — renderer construction failed")
	}
	if _, ok := msg.(RendererReadyMsg); !ok {
		t.Fatalf("expected RendererReadyMsg, got %T", msg)
	}
}
