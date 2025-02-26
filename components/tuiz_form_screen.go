package components

import (
	"fmt"
	"github.com/faelmori/kbx/mods/logz"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	. "github.com/faelmori/xtui/types"
)

var (
	kbdzInputsFocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	kbdzInputsBlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	kbdzInputsCursorStyle         = kbdzInputsFocusedStyle
	kbdzInputsNoStyle             = lipgloss.NewStyle()
	kbdzInputsHelpStyle           = kbdzInputsBlurredStyle
	kbdzInputsCursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	kbdzInputsErrorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))

	kbdzInputsFocusedButton = kbdzInputsFocusedStyle.Render("[ Proceed ]")
	kbdzInputsBlurredButton = fmt.Sprintf("[ %s ]", kbdzInputsBlurredStyle.Render("Proceed"))
)

var inputResult map[string]string

type FormScreenModel struct {
	Title        string
	FocusIndex   int
	Inputs       []textinput.Model
	CursorMode   cursor.Mode
	Fields       []TuizInputz
	ErrorMessage string
}

func kbdzInputsInitialModel(config TuizConfigz) FormScreenModel {
	cfg := config
	var inputs []TuizInput
	for _, field := range cfg.Fds.Inputs() {
		inputs = append(inputs, field.(TuizInput))
	}

	m := FormScreenModel{
		Title:        cfg.Title(),
		FocusIndex:   0,
		CursorMode:   cursor.CursorBlink,
		Fields:       config.Fds.Inputs(),
		Inputs:       make([]textinput.Model, len(inputs)),
		ErrorMessage: "",
	}

	var t textinput.Model
	for i, field := range inputs {
		t = textinput.New()
		t.Cursor.Style = kbdzInputsCursorStyle
		t.CharLimit = 32
		t.Placeholder = field.Placeholder()
		t.SetValue(field.Value())

		if field.Tp == PASSWORD {
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		if i == 0 {
			t.Focus()
			t.PromptStyle = kbdzInputsFocusedStyle
			t.TextStyle = kbdzInputsFocusedStyle
		}

		m.Inputs[i] = t
	}

	return m
}

func (m *FormScreenModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *FormScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case CTRLC, ESC:
			return m, tea.Quit
		case CTRLR:
			m.CursorMode++
			if m.CursorMode > cursor.CursorHide {
				m.CursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := range m.Inputs {
				cmds[i] = m.Inputs[i].Cursor.SetMode(m.CursorMode)
			}
			return m, tea.Batch(cmds...)

		case TAB, SHIFTTAB, ENTER, UP, DOWN:
			s := msg.String()

			if s == ENTER && m.FocusIndex == len(m.Inputs) {
				return m, m.submit()
			}

			if s == UP || s == SHIFTTAB {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex > len(m.Inputs) {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = len(m.Inputs)
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					cmds[i] = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = kbdzInputsFocusedStyle
					m.Inputs[i].TextStyle = kbdzInputsFocusedStyle
					continue
				}
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = kbdzInputsNoStyle
				m.Inputs[i].TextStyle = kbdzInputsNoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *FormScreenModel) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("\n%s\n\n", m.Title))

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &kbdzInputsBlurredButton
	if m.FocusIndex == len(m.Inputs) {
		button = &kbdzInputsFocusedButton
	}
	_, _ = fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if m.ErrorMessage != "" {
		b.WriteString(kbdzInputsErrorStyle.Render(m.ErrorMessage))
		b.WriteString("\n\n")
	}

	b.WriteString(kbdzInputsHelpStyle.Render("cursor mode is "))
	b.WriteString(kbdzInputsCursorModeHelpStyle.Render(m.CursorMode.String()))
	b.WriteString(kbdzInputsHelpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (m *FormScreenModel) submit() tea.Cmd {
	for i, input := range m.Inputs {
		value := input.Value()
		field := m.Fields[i]

		if field.Required() && value == "" {
			m.ErrorMessage = field.ErrorMessage()
			return nil
		}
		if field.MinLength() > 0 && len(value) < field.MinLength() {
			m.ErrorMessage = field.ErrorMessage()
			return nil
		}
		if field.MaxLength() > 0 && len(value) > field.MaxLength() {
			m.ErrorMessage = field.ErrorMessage()
			return nil
		}
		if field.Validator()(value) != nil {
			if err := field.Validator()(value); err != nil {
				m.ErrorMessage = err.Error()
				return nil
			}
		}

		inputResult[fmt.Sprintf("field%d", i)] = value
	}

	m.ErrorMessage = ""
	return tea.Quit
}

func KbdzInputs(config TuizConfig) (map[string]string, error) {
	inputResult = make(map[string]string)
	var newConfig TuizConfigz
	var newFields = config.Fields()
	if newFields == nil {
		newConfig = TuizConfigz{
			Tt:  config.Title(),
			Fds: nil,
		}
	} else {
		newConfig = TuizConfigz{
			Tt:  config.Title(),
			Fds: newFields.(TuizFieldz),
		}
	}
	initialModel := kbdzInputsInitialModel(newConfig)
	_, resultModelErr := tea.NewProgram(&initialModel).Run()
	if resultModelErr != nil {
		return nil, logz.ErrorLog("Error running inputs model: "+resultModelErr.Error(), "ui")
	}
	return inputResult, nil
}

func (m *FormScreenModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
