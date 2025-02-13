package services

// A simple program that makes a GET request and prints the response status.

import (
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var tcpStatusUrl string

type TcpStatusModel struct {
	status int
	err    error
}

type tcpStatusStatusMsg int

type tcpStatusErrMsg struct{ error }

func (e tcpStatusErrMsg) Error() string { return e.error.Error() }

func TcpStatus(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: tcp-status <url>")
	}
	tcpStatusUrl = args[0]
	p := tea.NewProgram(TcpStatusModel{})
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func (m TcpStatusModel) Init() tea.Cmd {
	return tcpStatusCheckServer
}

func (m TcpStatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}

	case tcpStatusStatusMsg:
		m.status = int(msg)
		return m, tea.Quit

	case tcpStatusErrMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m TcpStatusModel) View() string {
	s := fmt.Sprintf("Checking %s...", tcpStatusUrl)
	if m.err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.err)
	} else if m.status != 0 {
		s += fmt.Sprintf("%d %s", m.status, http.StatusText(m.status))
	}
	return s + "\n"
}

func tcpStatusCheckServer() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(tcpStatusUrl)
	if err != nil {
		return tcpStatusErrMsg{err}
	}
	defer res.Body.Close() // nolint:errcheck

	return tcpStatusStatusMsg(res.StatusCode)
}
