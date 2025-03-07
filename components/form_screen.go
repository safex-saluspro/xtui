package components

import (
	"fmt"
	"github.com/faelmori/logz"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	. "github.com/faelmori/xtui/types"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))

	focusedButton = focusedStyle.Render("[ Proceed ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Proceed"))
)

var inputResult map[string]string

type FormModel struct {
	Title        string
	FocusIndex   int
	Inputs       []textinput.Model
	CursorMode   cursor.Mode
	Fields       []FormField
	ErrorMessage string
}

func initialFormModel(config Config) FormModel {
	cfg := &config
	var inputs []FormField

	for _, field := range cfg.Fields.Inputs() {
		inputs = append(inputs, field.(FormField))
	}

	availableProperties := getAvailableProperties()
	if len(availableProperties) > 0 {
		inputs = adaptInputsToProperties(inputs, availableProperties)
	}

	m := FormModel{
		Title:        cfg.Title,
		FocusIndex:   0,
		CursorMode:   cursor.CursorBlink,
		Fields:       config.Fields.Inputs(),
		Inputs:       make([]textinput.Model, len(inputs)),
		ErrorMessage: "",
	}

	var t textinput.Model
	for i, field := range inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Placeholder = field.Placeholder()
		t.SetValue(field.Value())

		if field.Type() == PASSWORD {
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		m.Inputs[i] = t
	}

	return m
}

func (m *FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+r":
			m.CursorMode++
			if m.CursorMode > cursor.CursorHide {
				m.CursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := range m.Inputs {
				cmds[i] = m.Inputs[i].Cursor.SetMode(m.CursorMode)
			}
			return m, tea.Batch(cmds...)

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.FocusIndex == len(m.Inputs) {
				return m, m.submit()
			}

			if s == "up" || s == "shift+tab" {
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
					m.Inputs[i].PromptStyle = focusedStyle
					m.Inputs[i].TextStyle = focusedStyle
					continue
				}
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = noStyle
				m.Inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *FormModel) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("\n%s\n\n", m.Title))

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.FocusIndex == len(m.Inputs) {
		button = &focusedButton
	}
	_, _ = fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if m.ErrorMessage != "" {
		b.WriteString(errorStyle.Render(m.ErrorMessage))
		b.WriteString("\n\n")
	}

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.CursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (m *FormModel) submit() tea.Cmd {
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

func ShowForm(config Config) (map[string]string, error) {
	inputResult = make(map[string]string)
	var newConfig Config
	var newFields = config.Fields.Inputs()
	if newFields == nil {
		iNewConfig := FormConfig{
			Title:  config.Title,
			Fields: nil,
		}
		newConfig = Config{
			Title: iNewConfig.Title,
			Fields: FormFields{
				Title:  iNewConfig.Title,
				Fields: config.GetFields().Inputs(),
			},
		}
	}
	initialModel := initialFormModel(newConfig)
	_, resultModelErr := tea.NewProgram(&initialModel).Run()
	if resultModelErr != nil {
		logz.Error("Error running form model.", map[string]interface{}{
			"context": "ShowForm",
			"error":   resultModelErr,
		})
		return nil, resultModelErr
	}
	return inputResult, nil
}

func (m *FormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func getAvailableProperties() map[string]string {
	return map[string]string{
		"property1": "value1",
		"property2": "value2",
	}
}

func adaptInputsToProperties(inputs []FormField, properties map[string]string) []FormField {
	adaptedInputs := inputs
	for key, value := range properties {
		adaptedInputs = append(adaptedInputs, InputField{
			Ph:  key,
			Tp:  "text",
			Val: value,
			Req: false,
			Min: 0,
			Max: 100,
			Err: "",
			Vld: func(value string) error { return nil },
		})
	}
	return adaptedInputs
}
