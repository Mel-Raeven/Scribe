package ui

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

// RendererReadyMsg is sent when glamour finishes building a renderer.
type RendererReadyMsg struct {
	Renderer *glamour.TermRenderer
	Width    int
}

// glamourStyleJSON returns a custom glamour style JSON for the given theme.
// Both Catppuccin Mocha (dark) and Latte (light) — headings rendered without
// the raw ## prefix that the built-in glamour styles show.
func glamourStyleJSON(dark bool) []byte {
	if dark {
		return []byte(`{
  "document": { "block_prefix": "\n", "block_suffix": "\n", "color": "#cdd6f4", "margin": 2 },
  "block_quote": { "indent": 1, "indent_token": "│ ", "color": "#a6adc8" },
  "paragraph": {},
  "list": { "level_indent": 2 },
  "heading": { "block_suffix": "\n", "bold": true },
  "h1": { "prefix": " ", "suffix": " ", "color": "#1e1e2e", "background_color": "#89b4fa", "bold": true },
  "h2": { "color": "#89b4fa", "bold": true, "underline": true },
  "h3": { "color": "#74c7ec", "bold": true },
  "h4": { "color": "#89dceb", "bold": true },
  "h5": { "color": "#a6e3a1", "bold": true },
  "h6": { "color": "#a6adc8", "bold": false },
  "text": {},
  "strikethrough": { "crossed_out": true },
  "emph": { "italic": true },
  "strong": { "bold": true },
  "hr": { "color": "#45475a", "format": "\n--------\n" },
  "item": { "block_prefix": "• " },
  "enumeration": { "block_prefix": ". " },
  "task": { "ticked": "[✓] ", "unticked": "[ ] " },
  "link": { "color": "#89b4fa", "underline": true },
  "link_text": { "color": "#cba6f7", "bold": true },
  "image": { "color": "#f38ba8", "underline": true },
  "image_text": { "color": "#6c7086", "format": "Image: {{.text}} →" },
  "code": { "prefix": " ", "suffix": " ", "color": "#f38ba8", "background_color": "#313244" },
  "code_block": { "color": "#cdd6f4", "margin": 2, "chroma": {
    "text":                  { "color": "#cdd6f4" },
    "error":                 { "color": "#f38ba8" },
    "comment":               { "color": "#6c7086", "italic": true },
    "comment_preproc":       { "color": "#89b4fa" },
    "keyword":               { "color": "#cba6f7", "bold": true },
    "keyword_reserved":      { "color": "#cba6f7" },
    "keyword_namespace":     { "color": "#f38ba8" },
    "keyword_type":          { "color": "#f9e2af" },
    "operator":              { "color": "#89dceb" },
    "punctuation":           { "color": "#cdd6f4" },
    "name":                  { "color": "#cdd6f4" },
    "name_builtin":          { "color": "#fab387" },
    "name_tag":              { "color": "#f38ba8" },
    "name_attribute":        { "color": "#89b4fa" },
    "name_class":            { "color": "#f9e2af", "bold": true },
    "name_constant":         { "color": "#fab387" },
    "name_decorator":        { "color": "#89b4fa" },
    "name_function":         { "color": "#89b4fa" },
    "literal_number":        { "color": "#fab387" },
    "literal_string":        { "color": "#a6e3a1" },
    "literal_string_escape": { "color": "#f2cdcd" },
    "generic_deleted":       { "color": "#f38ba8" },
    "generic_emph":          { "italic": true },
    "generic_inserted":      { "color": "#a6e3a1" },
    "generic_strong":        { "bold": true },
    "generic_subheading":    { "color": "#a6adc8" },
    "background":            { "background_color": "#1e1e2e" }
  } },
  "table": {},
  "definition_list": {},
  "definition_term": {},
  "definition_description": { "block_prefix": "\n🠶 " },
  "html_block": {},
  "html_span": {}
}`)
	}
	return []byte(`{
  "document": { "block_prefix": "\n", "block_suffix": "\n", "color": "#4c4f69", "margin": 2 },
  "block_quote": { "indent": 1, "indent_token": "│ ", "color": "#5c5f77" },
  "paragraph": {},
  "list": { "level_indent": 2 },
  "heading": { "block_suffix": "\n", "bold": true },
  "h1": { "prefix": " ", "suffix": " ", "color": "#eff1f5", "background_color": "#1e66f5", "bold": true },
  "h2": { "color": "#1e66f5", "bold": true, "underline": true },
  "h3": { "color": "#209fb5", "bold": true },
  "h4": { "color": "#04a5e5", "bold": true },
  "h5": { "color": "#40a02b", "bold": true },
  "h6": { "color": "#8c8fa1", "bold": false },
  "text": {},
  "strikethrough": { "crossed_out": true },
  "emph": { "italic": true },
  "strong": { "bold": true },
  "hr": { "color": "#9ca0b0", "format": "\n--------\n" },
  "item": { "block_prefix": "• " },
  "enumeration": { "block_prefix": ". " },
  "task": { "ticked": "[✓] ", "unticked": "[ ] " },
  "link": { "color": "#1e66f5", "underline": true },
  "link_text": { "color": "#8839ef", "bold": true },
  "image": { "color": "#d20f39", "underline": true },
  "image_text": { "color": "#8c8fa1", "format": "Image: {{.text}} →" },
  "code": { "prefix": " ", "suffix": " ", "color": "#d20f39", "background_color": "#e6e9ef" },
  "code_block": { "color": "#4c4f69", "margin": 2, "chroma": {
    "text":                  { "color": "#4c4f69" },
    "error":                 { "color": "#d20f39" },
    "comment":               { "color": "#9ca0b0", "italic": true },
    "comment_preproc":       { "color": "#1e66f5" },
    "keyword":               { "color": "#8839ef", "bold": true },
    "keyword_reserved":      { "color": "#8839ef" },
    "keyword_namespace":     { "color": "#d20f39" },
    "keyword_type":          { "color": "#df8e1d" },
    "operator":              { "color": "#04a5e5" },
    "punctuation":           { "color": "#4c4f69" },
    "name":                  { "color": "#4c4f69" },
    "name_builtin":          { "color": "#fe640b" },
    "name_tag":              { "color": "#d20f39" },
    "name_attribute":        { "color": "#1e66f5" },
    "name_class":            { "color": "#df8e1d", "bold": true },
    "name_constant":         { "color": "#fe640b" },
    "name_decorator":        { "color": "#1e66f5" },
    "name_function":         { "color": "#1e66f5" },
    "literal_number":        { "color": "#fe640b" },
    "literal_string":        { "color": "#40a02b" },
    "literal_string_escape": { "color": "#179299" },
    "generic_deleted":       { "color": "#d20f39" },
    "generic_emph":          { "italic": true },
    "generic_inserted":      { "color": "#40a02b" },
    "generic_strong":        { "bold": true },
    "generic_subheading":    { "color": "#6c6f85" },
    "background":            { "background_color": "#e6e9ef" }
  } },
  "table": {},
  "definition_list": {},
  "definition_term": {},
  "definition_description": { "block_prefix": "\n🠶 " },
  "html_block": {},
  "html_span": {}
}`)
}

// buildRendererCmd constructs a glamour renderer in a goroutine.
// style should be "dark" or "light" — detected before Bubble Tea starts.
func buildRendererCmd(width int, style string) tea.Cmd {
	return func() tea.Msg {
		r, err := glamour.NewTermRenderer(
			glamour.WithStylesFromJSONBytes(glamourStyleJSON(style == "dark")),
			glamour.WithWordWrap(width),
		)
		if err != nil {
			return nil
		}
		return RendererReadyMsg{Renderer: r, Width: width}
	}
}

// Preview renders content using Glamour.
type Preview struct {
	content       string
	rendered      string
	filePath      string
	vaultRoot     string
	width         int
	height        int
	scrollY       int
	renderer      *glamour.TermRenderer
	rendererWidth int    // width the current renderer was built for
	style         string // "dark" or "light"
}

func newPreview(style, vaultRoot string) *Preview {
	return &Preview{style: style, vaultRoot: vaultRoot}
}

// setSize records the new dimensions. Returns a Cmd if the renderer needs
// to be rebuilt (width changed); nil otherwise.
// One column is reserved for the scrollbar, so the glamour renderer is
// built for (w-1) and content is wrapped at that width.
func (p *Preview) setSize(w, h int) tea.Cmd {
	p.height = h
	if w == p.width {
		return nil // height-only change — renderer still valid
	}
	p.width = w
	contentW := p.contentWidth()
	if p.rendererWidth == contentW && p.renderer != nil {
		return nil // already have a renderer for this width
	}
	return buildRendererCmd(contentW, p.style)
}

// contentWidth returns the width available for rendered text (full width
// minus the 1-column scrollbar).
func (p *Preview) contentWidth() int {
	w := p.width - 1
	if w < 1 {
		w = 1
	}
	return w
}

// applyRenderer stores a freshly-built renderer and re-renders any buffered content.
func (p *Preview) applyRenderer(msg RendererReadyMsg) {
	// Discard if a resize happened after this cmd was dispatched
	if msg.Width != p.contentWidth() {
		return
	}
	p.renderer = msg.Renderer
	p.rendererWidth = msg.Width
	if p.content != "" {
		p.render()
	}
}

// SetContent updates the raw content and re-renders if renderer is ready.
// If the renderer isn't ready yet, the raw text is shown until it arrives.
func (p *Preview) SetContent(path, content string) {
	p.filePath = path
	p.content = content
	p.scrollY = 0
	if p.renderer != nil {
		p.render()
	} else {
		p.rendered = content
	}
}

// wrapLines hard-wraps each line of content to at most width runes, splitting
// at character boundaries. This prevents long lines in plain-text files from
// overflowing the terminal layout when rendered inside a lipgloss box.
func wrapLines(content string, width int) string {
	if width <= 0 {
		return content
	}
	lines := strings.Split(content, "\n")
	var out []string
	for _, line := range lines {
		runes := []rune(line)
		for len(runes) > width {
			out = append(out, string(runes[:width]))
			runes = runes[width:]
		}
		out = append(out, string(runes))
	}
	return strings.Join(out, "\n")
}

func (p *Preview) render() {
	if p.renderer == nil || p.content == "" {
		p.rendered = p.content
		return
	}

	ext := strings.ToLower(filepath.Ext(p.filePath))
	if ext != ".md" {
		// For plain-text files, hard-wrap long lines so they don't overflow
		// the lipgloss layout container and push other panes off-screen.
		p.rendered = wrapLines(p.content, p.contentWidth())
		return
	}

	rendered, err := p.renderer.Render(p.content)
	if err != nil {
		p.rendered = p.content
		return
	}
	p.rendered = rendered
}

// ScrollY returns the current scroll offset (line index of the top visible line).
func (p *Preview) ScrollY() int {
	return p.scrollY
}

func (p *Preview) ScrollUp() {
	if p.scrollY > 0 {
		p.scrollY--
	}
}

func (p *Preview) ScrollDown() {
	lines := strings.Split(p.rendered, "\n")
	if p.scrollY < len(lines)-p.height {
		p.scrollY++
	}
}

func (p *Preview) PageUp() {
	p.scrollY -= p.height
	if p.scrollY < 0 {
		p.scrollY = 0
	}
}

func (p *Preview) PageDown() {
	lines := strings.Split(p.rendered, "\n")
	p.scrollY += p.height
	maxScroll := len(lines) - p.height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if p.scrollY > maxScroll {
		p.scrollY = maxScroll
	}
}

func (p *Preview) View() string {
	if p.filePath == "" {
		return ""
	}
	lines := strings.Split(p.rendered, "\n")
	totalLines := len(lines)
	start := p.scrollY
	end := start + p.height
	if end > totalLines {
		end = totalLines
	}
	if start > totalLines {
		start = totalLines
	}
	visible := lines[start:end]

	contentW := p.contentWidth()
	barChars := scrollbarChars(p.scrollY, totalLines, p.height)

	result := make([]string, len(visible))
	for i, line := range visible {
		sc := " "
		if i < len(barChars) {
			sc = barChars[i]
		}
		// Pad each line to contentW so the scrollbar column is always flush right.
		padded := lipgloss.NewStyle().Width(contentW).Render(line)
		result[i] = padded + sc
	}
	return strings.Join(result, "\n")
}

// scrollbarChars builds a slice of single-character strings (one per visible
// row) representing the scrollbar column. When content fits without scrolling,
// a subtle track character is shown in a muted colour.
func scrollbarChars(scrollY, totalLines, height int) []string {
	bar := make([]string, height)

	if totalLines <= height {
		// Content fits — show a subtle track to fill the reserved column.
		for i := range bar {
			bar[i] = StyleScrollbarTrack.Render("▏")
		}
		return bar
	}

	// Proportional thumb: size scales with the visible fraction.
	thumbSize := height * height / totalLines
	if thumbSize < 1 {
		thumbSize = 1
	}

	maxScroll := totalLines - height
	thumbPos := 0
	if maxScroll > 0 {
		thumbPos = scrollY * (height - thumbSize) / maxScroll
	}

	for i := range bar {
		if i >= thumbPos && i < thumbPos+thumbSize {
			bar[i] = StyleScrollbarThumb.Render("█")
		} else {
			bar[i] = StyleScrollbarTrack.Render("▏")
		}
	}
	return bar
}
