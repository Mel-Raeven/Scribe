package ui

import "github.com/charmbracelet/lipgloss"

// Styles holds all lipgloss styles for the UI. Initialised by InitTheme.
var (
	StyleHeader        lipgloss.Style
	StyleHeaderPath    lipgloss.Style
	StyleHeaderMode    lipgloss.Style
	StyleSidebar       lipgloss.Style
	StyleSidebarTitle  lipgloss.Style
	StyleTreeItem      lipgloss.Style
	StyleTreeItemSelected lipgloss.Style
	StyleTreeDir       lipgloss.Style
	StyleTreeDirSelected lipgloss.Style
	StyleStatusBar     lipgloss.Style
	StyleStatusKey     lipgloss.Style
	StyleStatusSaved   lipgloss.Style
	StyleStatusModified lipgloss.Style
	StyleStatusError   lipgloss.Style
	StyleMainPane        lipgloss.Style
	StyleScrollbarTrack  lipgloss.Style
	StyleScrollbarThumb  lipgloss.Style
)

// InitTheme sets all package-level styles to either the dark (Catppuccin Mocha)
// or light (Catppuccin Latte) palette. Must be called before the TUI renders.
func InitTheme(dark bool) {
	var (
		colorBase      lipgloss.Color
		colorSurface   lipgloss.Color
		colorOverlay   lipgloss.Color
		colorMuted     lipgloss.Color
		colorSubtle    lipgloss.Color
		colorText      lipgloss.Color
		colorAccent    lipgloss.Color
		colorGreen     lipgloss.Color
		colorYellow    lipgloss.Color
		colorRed       lipgloss.Color
		colorHighlight lipgloss.Color
	)

	if dark {
		// Catppuccin Mocha
		colorBase      = lipgloss.Color("#1e1e2e")
		colorSurface   = lipgloss.Color("#313244")
		colorOverlay   = lipgloss.Color("#45475a")
		colorMuted     = lipgloss.Color("#6c7086")
		colorSubtle    = lipgloss.Color("#a6adc8")
		colorText      = lipgloss.Color("#cdd6f4")
		colorAccent    = lipgloss.Color("#89b4fa")
		colorGreen     = lipgloss.Color("#a6e3a1")
		colorYellow    = lipgloss.Color("#f9e2af")
		colorRed       = lipgloss.Color("#f38ba8")
		colorHighlight = lipgloss.Color("#585b70")
	} else {
		// Catppuccin Latte
		colorBase      = lipgloss.Color("#eff1f5")
		colorSurface   = lipgloss.Color("#e6e9ef")
		colorOverlay   = lipgloss.Color("#9ca0b0")
		colorMuted     = lipgloss.Color("#8c8fa1")
		colorSubtle    = lipgloss.Color("#6c6f85")
		colorText      = lipgloss.Color("#4c4f69")
		colorAccent    = lipgloss.Color("#1e66f5")
		colorGreen     = lipgloss.Color("#40a02b")
		colorYellow    = lipgloss.Color("#df8e1d")
		colorRed       = lipgloss.Color("#d20f39")
		colorHighlight = lipgloss.Color("#ccd0da")
	}

	_ = colorBase // used via terminal bg implicitly

	StyleHeader = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorText).
		Padding(0, 1)

	StyleHeaderPath = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorSubtle)

	StyleHeaderMode = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorAccent)

	StyleSidebar = lipgloss.NewStyle().
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorOverlay)

	StyleSidebarTitle = lipgloss.NewStyle().
		Foreground(colorMuted).
		Padding(0, 1).
		Bold(true)

	StyleTreeItem = lipgloss.NewStyle().
		Foreground(colorText).
		Padding(0, 0)

	StyleTreeItemSelected = lipgloss.NewStyle().
		Background(colorHighlight).
		Foreground(colorAccent).
		Bold(true)

	StyleTreeDir = lipgloss.NewStyle().
		Foreground(colorYellow)

	StyleTreeDirSelected = lipgloss.NewStyle().
		Background(colorHighlight).
		Foreground(colorYellow).
		Bold(true)

	StyleStatusBar = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorMuted).
		Padding(0, 1)

	StyleStatusKey = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorAccent).
		Bold(true)

	StyleStatusSaved = lipgloss.NewStyle().
		Foreground(colorGreen)

	StyleStatusModified = lipgloss.NewStyle().
		Foreground(colorYellow)

	StyleStatusError = lipgloss.NewStyle().
		Foreground(colorRed)

	StyleMainPane = lipgloss.NewStyle()

	StyleScrollbarTrack = lipgloss.NewStyle().
		Foreground(colorOverlay)

	StyleScrollbarThumb = lipgloss.NewStyle().
		Foreground(colorAccent)
}
