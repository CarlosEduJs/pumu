package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Item represents a selectable item in the multi-select list.
type Item struct {
	Label    string
	Detail   string // e.g. formatted size
	Selected bool
}

// Result holds the outcome of the multi-select interaction.
type Result struct {
	Items    []Item
	Canceled bool
}

type model struct {
	title    string
	items    []Item
	cursor   int
	showHelp bool
	done     bool
	canceled bool
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			PaddingBottom(1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E")).
			Bold(true)

	unselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	itemLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB"))

	itemLabelDimStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9CA3AF"))

	detailStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA"))

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			PaddingTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			PaddingTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	helpTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true).
			Underline(true).
			PaddingBottom(1)
)

func initialModel(title string, items []Item) model {
	return model{
		title:    title,
		items:    items,
		cursor:   0,
		showHelp: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}
	case "home", "g":
		m.cursor = 0
	case "end", "G":
		m.cursor = len(m.items) - 1
	case " ":
		m.items[m.cursor].Selected = !m.items[m.cursor].Selected
	case "a":
		for i := range m.items {
			m.items[i].Selected = true
		}
	case "n":
		for i := range m.items {
			m.items[i].Selected = false
		}
	case "i":
		for i := range m.items {
			m.items[i].Selected = !m.items[i].Selected
		}
	case "?":
		m.showHelp = !m.showHelp
	case "enter":
		m.done = true
		return m, tea.Quit
	case "q", "esc", "ctrl+c":
		m.canceled = true
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n")

	// Items
	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▸ ")
		}

		checkbox := unselectedStyle.Render("[ ]")
		if item.Selected {
			checkbox = selectedStyle.Render("[✓]")
		}

		label := itemLabelDimStyle.Render(item.Label)
		if i == m.cursor {
			label = itemLabelStyle.Render(item.Label)
		}

		detail := ""
		if item.Detail != "" {
			detail = " " + detailStyle.Render(item.Detail)
		}

		b.WriteString(fmt.Sprintf("%s%s %s%s\n", cursor, checkbox, label, detail))
	}

	// Status bar
	selected := 0
	for _, item := range m.items {
		if item.Selected {
			selected++
		}
	}
	b.WriteString(statusBarStyle.Render(
		fmt.Sprintf("  %d/%d selected", selected, len(m.items)),
	))
	b.WriteString("\n")

	// Help
	if m.showHelp {
		b.WriteString("\n")
		b.WriteString(helpTitleStyle.Render("Keyboard Shortcuts"))
		b.WriteString("\n")
		helpItems := []struct{ key, desc string }{
			{"↑/k", "move up"},
			{"↓/j", "move down"},
			{"g/G", "go to first/last"},
			{"space", "toggle item"},
			{"a", "select all"},
			{"n", "deselect all"},
			{"i", "invert selection"},
			{"enter", "confirm"},
			{"q/esc", "cancel"},
		}
		for _, h := range helpItems {
			b.WriteString(fmt.Sprintf("  %s %s\n",
				helpKeyStyle.Render(fmt.Sprintf("%-8s", h.key)),
				helpDescStyle.Render(h.desc),
			))
		}
	} else {
		b.WriteString(helpStyle.Render("  press ? for help"))
		b.WriteString("\n")
	}

	return b.String()
}

// RunMultiSelect launches an interactive multi-select prompt and returns the result.
// All items are pre-selected by default.
func RunMultiSelect(title string, items []Item) (Result, error) {
	m := initialModel(title, items)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return Result{}, fmt.Errorf("failed to run multi-select: %w", err)
	}

	fm, ok := finalModel.(model)
	if !ok {
		return Result{}, fmt.Errorf("unexpected model type")
	}

	return Result{
		Items:    fm.items,
		Canceled: fm.canceled,
	}, nil
}
