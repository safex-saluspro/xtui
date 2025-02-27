package wrappers

import (
	"bufio"
	"fmt"
	"github.com/faelmori/logz"

	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	treeHeight int
	logHeight  int
)

type LogzViewerModel struct {
	logs         []string
	moduleColors map[string]string
	mu           sync.Mutex
	treeView     string
	scrollOffset int
	autoScroll   bool
}

func (m *LogzViewerModel) Init() tea.Cmd {
	return tea.Batch(streamLogs(m.moduleColors, &m.mu), updateTreeView(&m.mu), logzViewerTickCmd())
}

func logzViewerTickCmd() tea.Cmd {
	return tea.Tick(time.Second/3, func(t time.Time) tea.Msg {
		return t
	})
}

func (m *LogzViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.scrollOffset > 0 {
				m.scrollOffset--
				m.autoScroll = false
			}
		case "down":
			if m.scrollOffset < len(m.logs)-1 {
				m.scrollOffset++
				if m.scrollOffset == len(m.logs)-1 {
					m.autoScroll = true
				}
			}
		}
	case time.Time:
		return m, tea.Batch(updateTreeView(&m.mu), logzViewerTickCmd())
	case string:
		m.mu.Lock()
		m.logs = append(m.logs, msg)
		if len(m.logs) > 100 {
			m.logs = m.logs[len(m.logs)-100:]
		}
		if m.autoScroll {
			lines, _ := strconv.Atoi(os.Getenv("LINES"))
			logHeight = lines - treeHeight - 2
			if len(m.logs) > logHeight {
				m.scrollOffset = len(m.logs) - logHeight
			}
		}
		m.mu.Unlock()
		return m, streamLogs(m.moduleColors, &m.mu)
	case logzViewerTreeViewMsg:
		m.mu.Lock()
		m.treeView = string(msg)
		m.mu.Unlock()
		return m, updateTreeView(&m.mu)
	}
	return m, nil
}

func (m *LogzViewerModel) View() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	logView := strings.Join(m.logs[m.scrollOffset:], "\n")
	treeView := lipgloss.NewStyle().Height(treeHeight).Render(m.treeView)
	return lipgloss.JoinVertical(lipgloss.Top, treeView, logView)
}

func parseAnsiColors(text string, moduleColors map[string]string) string {
	for module, color := range moduleColors {
		text = strings.ReplaceAll(text, module, lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(module))
	}
	return text
}

func streamLogs(moduleColors map[string]string, mu *sync.Mutex) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("kbx", "log", "--show=all", "-f")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logz.Panic("failed to get stdout pipe", "streamLogs")
		}

		if err := cmd.Start(); err != nil {
			logz.Panic("failed to start command: %v", err)
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			logLine := scanner.Text()
			mu.Lock()
			coloredLogLine := parseAnsiColors(logLine, moduleColors)
			mu.Unlock()
			return coloredLogLine
		}
		if err := scanner.Err(); err != nil {
			logz.Panic("error reading stdout: %v", err)
		}
		return ""
	}
}

type logzViewerTreeViewMsg string

func updateTreeView(mu *sync.Mutex) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("tree", os.Getenv("HOME")+"/.cache/kubex", "-s", "--du", "-C", "-h", "-P", "*.log")
		stdout, err := cmd.Output()
		if err != nil {
			logz.Panic("failed to execute command: %v", err)
		}
		treeView := string(stdout)
		treeHeight = len(strings.Split(treeView, "\n")) + 3
		mu.Lock()
		defer mu.Unlock()
		return logzViewerTreeViewMsg(treeView)
	}
}

func LogzViewer(args ...string) error {
	moduleColors := map[string]string{
		"module1": "1", // Defina as cores para cada módulo
		"module2": "2",
		// Adicione mais módulos conforme necessário
	}
	p := tea.NewProgram(&LogzViewerModel{moduleColors: moduleColors}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %v", err)
	}
	return nil
}
