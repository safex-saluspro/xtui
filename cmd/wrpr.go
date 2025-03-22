package main

import (
	"fmt"
	"github.com/faelmori/xtui/cmd/cli"
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
	return ""
}

// ShortDescription retorna uma descrição curta do módulo ui.
func (m *XTui) ShortDescription() string {
	return "Terminal UI"
}

// LongDescription retorna uma descrição longa do módulo ui.
func (m *XTui) LongDescription() string {
	return "Terminal XTUI module. It allows you to interact with the terminal using a graphical interface."
}

// Usage retorna a forma de uso do módulo ui.
func (m *XTui) Usage() string {
	return "xui [command] [args]"
}

// Examples retorna exemplos de uso do módulo ui.
func (m *XTui) Examples() []string {
	return []string{"xtui [command] [args]", "xtui logz -o 'file.log'", "xtui deps -o 'install'", "xtui tcp-status '127.0.0.1:8080'"}
}

// Active verifica se o módulo ui está ativo.
func (m *XTui) Active() bool {
	return true
}

// Module retorna o nome do módulo ui.
func (m *XTui) Module() string {
	return "xtui"
}

// Execute executa o comando especificado para o módulo ui.
func (m *XTui) Execute() error {
	return m.Command().Execute()
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
	var cd string
	var opts []string

	c := &cobra.Command{
		Use:         m.Module(),
		Aliases:     []string{m.Alias()},
		Example:     m.concatenateExamples(),
		Annotations: cli.GetDescriptions([]string{m.ShortDescription(), m.LongDescription()}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch cd {
			case "logz":
				return LogViewer(opts...)
			case "deps":
				return InstallDependenciesWithUI(opts...)
			case "tcp-status":
				return TcpStatus(args...)
			}

			return fmt.Errorf("error: %s", opts[0])
		},
	}

	c.Flags().StringArrayVarP(&opts, "opts", "o", []string{}, "Options")
	c.Flags().StringVarP(&cd, "cmd", "c", "logz", "Log file viewer")

	// Adiciona os comandos relacionados ao módulo

	pkgCmdRoot := &cobra.Command{
		Use:     "pkg",
		Aliases: []string{"package", "packages"},
		Annotations: cli.GetDescriptions(
			[]string{
				"Package management",
				"Package installation, removal, and management with friendly UI and much more",
			}, false,
		),
		RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	pkgCmdRoot.AddCommand(cli.PkgCmdsList()...)
	c.AddCommand(pkgCmdRoot)

	appCmdRoot := &cobra.Command{
		Use:     "deps",
		Aliases: []string{"dep", "dependencies"},
		Annotations: cli.GetDescriptions(
			[]string{
				"Dependencies management",
				"Install, remove, and manage dependencies with friendly UI and much more",
			}, false,
		),
		RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	appCmdRoot.AddCommand(cli.AppsCmdsList()...)
	c.AddCommand(appCmdRoot)

	formCmdRoot := &cobra.Command{
		Use:     "forms",
		Aliases: []string{"frm", "form"},
		Annotations: cli.GetDescriptions(
			[]string{
				"Terminal forms builder",
				"Build terminal forms with validation, input types, and much more",
			}, false,
		),
		RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	formCmdRoot.AddCommand(cli.FormsCmdsList()...)
	c.AddCommand(formCmdRoot)

	dataCmdRoot := &cobra.Command{
		Use:     "forms",
		Aliases: []string{"frm", "form"},
		Annotations: cli.GetDescriptions(
			[]string{
				"Terminal forms builder",
				"Build terminal forms with validation, input types, and much more",
			}, false,
		),
		RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	dataCmdRoot.AddCommand(cli.ViewsCmdsList()...)
	c.AddCommand(dataCmdRoot)

	setUsageDefinition(c)
	for _, subCmd := range c.Commands() {
		setUsageDefinition(subCmd)
	}

	return c
}
