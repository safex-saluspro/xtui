package services

import (
	"flag"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"io"
	"log"
	"os"
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	mainStyle = lipgloss.NewStyle().MarginLeft(1)
)

type DaemonizeModel struct {
	spinner  spinner.Model
	quitting bool
	initFunc func() tea.Msg
}

func DaemonizeNewModel(initFunc func() tea.Msg) DaemonizeModel {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))

	return DaemonizeModel{
		spinner:  sp,
		initFunc: initFunc,
	}
}

func (m DaemonizeModel) Init() tea.Cmd {
	log.Println("Starting work...")
	return tea.Batch(
		m.spinner.Tick,
		m.initFunc,
	)
}

func (m DaemonizeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case error:
		log.Printf("Error: %v", msg)
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m DaemonizeModel) View() string {
	s := "\n" +
		m.spinner.View() + " Running daemon...\n\n"

	s += helpStyle("\nPress any key to exit\n")

	if m.quitting {
		s += "\n"
	}

	return mainStyle.Render(s)
}

func Daemonize(initFunc func() tea.Msg, args ...string) error {
	var (
		daemonMode bool
		showHelp   bool
		opts       []tea.ProgramOption
	)

	flag.BoolVar(&daemonMode, "d", false, "run as a daemon")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.Parse()

	if showHelp {
		flag.Usage()
		return nil
	}

	if daemonMode || !isatty.IsTerminal(os.Stdout.Fd()) {
		opts = []tea.ProgramOption{tea.WithoutRenderer()}
	} else {
		log.SetOutput(io.Discard)
	}

	p := tea.NewProgram(DaemonizeNewModel(initFunc), opts...)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
