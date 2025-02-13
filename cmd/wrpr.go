package cmd

import (
	"fmt"
	. "github.com/faelmori/xtui/cmd/cli"
	. "github.com/faelmori/xtui/components"
	. "github.com/faelmori/xtui/services"
	. "github.com/faelmori/xtui/wrappers"
	"github.com/spf13/cobra"
)

// XTui representa a estrutura do módulo ui.
type XTui struct{}

// RegX registra e retorna uma nova instância de XTui.
func RegX() *XTui {
	return &XTui{}
}

// Alias retorna o alias do módulo ui.
func (m *XTui) Alias() string {
	return "tuiz"
}

// ShortDescription retorna uma descrição curta do módulo ui.
func (m *XTui) ShortDescription() string {
	return "Terminal UI"
}

// LongDescription retorna uma descrição longa do módulo ui.
func (m *XTui) LongDescription() string {
	return "Terminal UI module. It allows you to interact with the terminal using a graphical interface."
}

// Usage retorna a forma de uso do módulo ui.
func (m *XTui) Usage() string {
	return "kbx ui [command] [args]"
}

// Examples retorna exemplos de uso do módulo ui.
func (m *XTui) Examples() []string {
	return []string{"kbx ui -c logz", "kbx ui -c kbx-deps"}
}

// Active verifica se o módulo ui está ativo.
func (m *XTui) Active() bool {
	return true
}

// Module retorna o nome do módulo ui.
func (m *XTui) Module() string {
	return "ui"
}

// Execute executa o comando especificado para o módulo ui.
func (m *XTui) Execute(commandArgs []string) error {

	return fmt.Errorf("comando ou flag inválido")
}

// concatenateExamples concatena os exemplos de uso do módulo.
func (m *XTui) concatenateExamples() string {
	examples := ""
	for _, example := range m.Examples() {
		examples += string(example) + "\n  "
	}
	return examples
}

// Command retorna o comando cobra para o módulo.
func (m *XTui) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     m.Module(),
		Aliases: []string{m.Alias()},
		Example: m.concatenateExamples(),
		Short:   m.ShortDescription(),
		Long:    m.LongDescription(),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdFlagValue, _ := cmd.Flags().GetString("cmd")
			optsFlagValue, _ := cmd.Flags().GetStringArray("opts")

			optsFlagValue = append(optsFlagValue, args...)

			switch cmdFlagValue {
			case "logz":
				return LogzViewer(optsFlagValue...)
			case "split-screen":
				return SplitScreenNew(optsFlagValue...)
			case "deps":
				return InstallDepsWithUI(optsFlagValue...)
			case "tcp-status":
				return TcpStatus(args...)
				//case "daemon":
				//	return Daemonize(args...)
			}

			return fmt.Errorf("error: %s", optsFlagValue[0])
		},
	}

	cmd.Flags().StringArrayP("opts", "o", []string{}, "Options")
	cmd.Flags().StringP("cmd", "c", "logz", "Log file viewer")

	appsListCmds := AppsCmdsList()
	cmd.AddCommand(appsListCmds...)

	srvCmds := ServicesCmds()
	cmd.AddCommand(srvCmds...)

	formsListCmds := FormsCmdsList()
	cmd.AddCommand(formsListCmds...)

	viewsListCmds := ViewsCmdsList()
	cmd.AddCommand(viewsListCmds...)

	return cmd
}
