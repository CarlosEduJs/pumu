package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple
	ColorSuccess   = lipgloss.Color("#22C55E") // Green
	ColorWarning   = lipgloss.Color("#FACC15") // Yellow
	ColorDanger    = lipgloss.Color("#EF4444") // Red
	ColorInfo      = lipgloss.Color("#3B82F6") // Blue
	ColorText      = lipgloss.Color("#E5E7EB") // Dim White
	ColorSubtext   = lipgloss.Color("#9CA3AF") // Gray
	ColorHighlight = lipgloss.Color("#A78BFA") // Light Purple

	// Text Styles
	BoldStyle   = lipgloss.NewStyle().Bold(true)
	ItalicStyle = lipgloss.NewStyle().Italic(true)

	// Context Styles
	AlertStyle   = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	InfoStyle    = lipgloss.NewStyle().Foreground(ColorInfo)
	PrimaryStyle = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	SubtextStyle = lipgloss.NewStyle().Foreground(ColorSubtext)
	DetailStyle  = lipgloss.NewStyle().Foreground(ColorInfo)

	// Layout Styles
	TableColumnStyle = lipgloss.NewStyle().PaddingRight(2)
)

// RenderRow aligns a set of strings into columns with fixed widths.
func RenderRow(items []string, widths []int) string {
	var row []string
	for i, item := range items {
		if i < len(widths) {
			row = append(row, lipgloss.NewStyle().Width(widths[i]).MaxWidth(widths[i]).Render(item))
		} else {
			row = append(row, item)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, row...)
}
