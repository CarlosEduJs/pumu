package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgressMsg struct {
	Increment float64
}

type ProgressDoneMsg struct{}

type progressModel struct {
	progress progress.Model
	message  string
	total    int
	current  int
	percent  float64
	done     bool
}

func initialProgressModel(message string, total int) progressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
	)
	return progressModel{
		progress: p,
		message:  message,
		total:    total,
	}
}

func (m progressModel) Init() tea.Cmd {
	return nil
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 10
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		return m, nil
	case ProgressMsg:
		m.current++
		m.percent = float64(m.current) / float64(m.total)
		cmd := m.progress.SetPercent(m.percent)
		if m.current >= m.total {
			m.done = true
			return m, tea.Sequence(cmd, tea.Quit)
		}
		return m, cmd
	case ProgressDoneMsg:
		m.done = true
		return m, tea.Quit
	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		if p, ok := progressModel.(progress.Model); ok {
			m.progress = p
		}
		return m, cmd
	}
	return m, nil
}

func (m progressModel) View() string {
	if m.done {
		return ""
	}
	pad := strings.Repeat(" ", 2)
	text := PrimaryStyle.Render(fmt.Sprintf("%s [%d/%d]", m.message, m.current, m.total))

	// percentage
	percentStr := lipgloss.NewStyle().Foreground(ColorSubtext).Render(fmt.Sprintf("%3.0f%%", m.percent*100))

	return "\n" + pad + text + "\n" + pad + m.progress.View() + " " + percentStr + "\n"
}

// TrackProgress starts a progress bar program. It accepts a channel to receive iteration ticks.
func TrackProgress(title string, totalItems int, tickChan <-chan struct{}) error {
	if totalItems <= 0 {
		return nil
	}

	p := tea.NewProgram(initialProgressModel(title, totalItems))

	// Goroutine to translate channel tickets to Bubble Tea messages
	go func() {
		for range tickChan {
			p.Send(ProgressMsg{Increment: 1.0})
		}
		p.Send(ProgressDoneMsg{})
	}()

	_, err := p.Run()
	return err
}
