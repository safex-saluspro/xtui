package wrappers

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	kbdzLoaderSpinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("215"))
	kbdzLoaderHelpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	kbdzLoaderDotStyle      = kbdzLoaderHelpStyle.UnsetMargins()
	kbdzLoaderDurationStyle = kbdzLoaderDotStyle
	kbdzLoaderAppStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
	kbdzLoaderErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	kbdzLoaderTitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	kbdzLoaderMessageStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))

	// Novos estilos
	kbdzLoaderSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	kbdzLoaderWarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Italic(true)
	kbdzLoaderInfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Underline(true)
)

type KbdzLoaderMsg struct {
	Message string
}

type KbdzLoaderCloseMsg struct{}

type kbdzLoaderModel struct {
	spinner  spinner.Model
	messages []KbdzLoaderMsg
	quitting bool
	err      error
}

func kbdzLoaderNewModel() kbdzLoaderModel {
	s := spinner.New()
	s.Style = kbdzLoaderSpinnerStyle
	return kbdzLoaderModel{
		spinner:  s,
		messages: []KbdzLoaderMsg{},
	}
}

func (m kbdzLoaderModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m kbdzLoaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case KbdzLoaderMsg:
		m.messages = append(m.messages, msg)
		return m, nil
	case KbdzLoaderCloseMsg:
		m.quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m kbdzLoaderModel) View() string {
	var s string

	if m.quitting {
		s += kbdzLoaderTitleStyle.Render("Done!")
	} else {
		s += kbdzLoaderTitleStyle.Render(m.spinner.View() + " Working...")
	}

	s += "\n\n"

	for _, msg := range m.messages {
		// Aplicar diferentes estilos às mensagens
		if strings.Contains(msg.Message, "Error: ") {
			s += kbdzLoaderErrorStyle.Render(msg.Message) + "\n"
		} else if strings.Contains(msg.Message, "Success: ") {
			s += kbdzLoaderSuccessStyle.Render(msg.Message) + "\n"
		} else if strings.Contains(msg.Message, "Warning: ") {
			s += kbdzLoaderWarningStyle.Render(msg.Message) + "\n"
		} else {
			s += kbdzLoaderMessageStyle.Render(msg.Message) + "\n"
		}
	}

	if !m.quitting {
		s += kbdzLoaderHelpStyle.Render("Press any key to quit")
	}

	if m.quitting {
		s += "\n"
	}

	return kbdzLoaderAppStyle.Render(s)
}

func StartLoader(messages chan tea.Msg) error {
	p := tea.NewProgram(kbdzLoaderNewModel())

	go func() {
		for msg := range messages {
			p.Send(msg)
		}
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
