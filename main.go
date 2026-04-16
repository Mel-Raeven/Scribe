package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/melm/scribe/ui"
	"github.com/muesli/termenv"
)

func main() {
	root := "."

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "--help" || arg == "-h" {
			fmt.Println("scribe [path]")
			fmt.Println("  Open scribe in the given directory or file. Defaults to current directory.")
			os.Exit(0)
		}
		root = arg
	}

	abs, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving path: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "path does not exist: %s\n", abs)
		os.Exit(1)
	}

	// If a file is passed directly, use its parent dir as root and pre-open it.
	startFile := ""
	if !info.IsDir() {
		startFile = abs
		abs = filepath.Dir(abs)
	}

	// Detect dark/light background before Bubble Tea takes over stdin.
	glamourStyle := "dark"
	if !termenv.HasDarkBackground() {
		glamourStyle = "light"
	}

	p := tea.NewProgram(
		ui.New(abs, startFile, glamourStyle),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "scribe: %v\n", err)
		os.Exit(1)
	}
}
