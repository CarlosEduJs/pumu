// Package ui provides interactive TUI components for pumu using Bubble Tea.
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
	filter   string
	filtered []int
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

	filterLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F59E0B")).
				Bold(true)

	filterValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E5E7EB"))

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
		filter:   "",
		filtered: nil,
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
	case "up", "k", "down", "j", "home", "g", "end", "G":
		m.handleNavigation(msg.String())
	case " ", "a", "n", "i":
		m.handleSelection(msg.String())
	case "backspace":
		m.handleFilterBackspace()
	case "ctrl+u":
		m.handleFilterClear()
	case "esc":
		if m.filter != "" {
			m.handleFilterClear()
			return m, nil
		}
		m.canceled = true
		m.done = true
		return m, tea.Quit
	case "?":
		m.showHelp = !m.showHelp
	case "enter":
		m.done = true
		return m, tea.Quit
	case "q", "ctrl+c":
		m.canceled = true
		m.done = true
		return m, tea.Quit
	default:
		if msg.Type == tea.KeyRunes {
			m.handleFilterAppend(msg.String())
		}
	}

	return m, nil
}

func (m *model) handleNavigation(key string) {
	visibleCount := m.visibleCount()
	if visibleCount == 0 {
		return
	}

	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < visibleCount-1 {
			m.cursor++
		}
	case "home", "g":
		m.cursor = 0
	case "end", "G":
		m.cursor = visibleCount - 1
	}
}

func (m *model) handleSelection(key string) {
	switch key {
	case " ":
		if index, ok := m.currentIndex(); ok {
			m.items[index].Selected = !m.items[index].Selected
		}
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
	}
}

func (m *model) handleFilterAppend(value string) {
	if value == "" || value == "?" {
		return
	}
	m.filter += value
	m.updateFilter()
}

func (m *model) handleFilterBackspace() {
	if len(m.filter) == 0 {
		return
	}
	m.filter = m.filter[:len(m.filter)-1]
	m.updateFilter()
}

func (m *model) handleFilterClear() {
	m.filter = ""
	m.updateFilter()
}

func (m *model) updateFilter() {
	if strings.TrimSpace(m.filter) == "" {
		m.filtered = nil
		m.cursor = 0
		return
	}
	query := strings.ToLower(m.filter)
	filtered := make([]int, 0, len(m.items))
	for i, item := range m.items {
		label := strings.ToLower(item.Label)
		detail := strings.ToLower(item.Detail)
		if strings.Contains(label, query) || (detail != "" && strings.Contains(detail, query)) {
			filtered = append(filtered, i)
		}
	}
	m.filtered = filtered
	if m.cursor >= len(m.filtered) {
		m.cursor = 0
	}
}

func (m *model) currentIndex() (int, bool) {
	if len(m.filtered) == 0 {
		if len(m.items) == 0 {
			return 0, false
		}
		return m.cursor, true
	}
	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return 0, false
	}
	return m.filtered[m.cursor], true
}

func (m model) visibleCount() int {
	if len(m.filtered) > 0 {
		return len(m.filtered)
	}
	return len(m.items)
}

func (m model) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n")

	// Filter
	filterValue := m.filter
	if filterValue == "" {
		filterValue = "(type to filter)"
	}
	b.WriteString(filterLabelStyle.Render("Filter:"))
	b.WriteString(" ")
	b.WriteString(filterValueStyle.Render(filterValue))
	b.WriteString("\n\n")

	// Items
	indices := m.filtered
	if len(indices) == 0 {
		indices = make([]int, len(m.items))
		for i := range m.items {
			indices[i] = i
		}
	}

	for i, index := range indices {
		item := m.items[index]
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

		fmt.Fprintf(&b, "%s%s %s%s\n", cursor, checkbox, label, detail)
	}

	// Status bar
	selected := 0
	for _, item := range m.items {
		if item.Selected {
			selected++
		}
	}
	visible := m.visibleCount()
	b.WriteString(statusBarStyle.Render(
		fmt.Sprintf("  %d/%d selected · %d shown", selected, len(m.items), visible),
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
			{"type", "filter list"},
			{"backspace", "delete filter"},
			{"ctrl+u", "clear filter"},
			{"enter", "confirm"},
			{"q/esc", "cancel"},
		}
		for _, h := range helpItems {
			fmt.Fprintf(&b, "  %s %s\n",
				helpKeyStyle.Render(fmt.Sprintf("%-8s", h.key)),
				helpDescStyle.Render(h.desc),
			)
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
