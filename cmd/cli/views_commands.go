package cli

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/faelmori/xtui/components"
	t "github.com/faelmori/xtui/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"testing"
)

func ViewsCmdsList() []*cobra.Command {
	tableCmd := tableViewCmd()

	return []*cobra.Command{
		tableCmd,
	}
}

func tableViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "table",
		Aliases: []string{"tbl-view", "view-table"},
		Annotations: GetDescriptions(
			[]string{
				"Table view for any command",
				"Table view screen, interactive mode, for any command with flags",
			},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return NavigateAndExecuteViewCommand(cmd, args)
		},
	}

	return cmd
}

func TestTableViewCmd(t *testing.T) {
	cmd := tableViewCmd()
	if cmd.Use != "table" {
		t.Errorf("expected 'table', got '%s'", cmd.Use)
	}
	if cmd.Short != "Table view for any command" {
		t.Errorf("expected 'Table view for any command', got '%s'", cmd.Short)
	}
	if cmd.Long != "Table view screen, interactive mode, for any command with flags" {
		t.Errorf("expected 'Table view screen, interactive mode, for any command with flags', got '%s'", cmd.Long)
	}
}

func NavigateAndExecuteViewCommand(cmd *cobra.Command, args []string) error {
	// Detect command and its flags
	commandName := cmd.Name()
	flags := cmd.Flags()

	// Display command selection and flag definition in a table view
	tableConfig := createTableConfig(commandName, flags)
	customStyles := map[string]lipgloss.Color{
		"Info":    lipgloss.Color("#75FBAB"),
		"Warning": lipgloss.Color("#FDFF90"),
		"Error":   lipgloss.Color("#FF7698"),
		"Debug":   lipgloss.Color("#929292"),
	}
	if err := components.StartTableScreen(tableConfig, customStyles); err != nil {
		return err
	}

	// Set flag values based on table input
	for key, value := range tableConfig.Fields {
		if err := cmd.Flags().Set(key, value.Value()); err != nil {
			return err
		}
	}

	// Execute the command
	return cmd.Execute()
}

func createTableConfig(commandName string, flags *pflag.FlagSet) t.FormConfig {
	var tableFields []t.Field

	flags.VisitAll(func(flag *pflag.Flag) {
		tableFields = append(tableFields, t.InputField{
			Ph:  flag.Name,
			Tp:  "text",
			Val: flag.Value.String(),
		})
	})

	return t.FormConfig{
		Title:  fmt.Sprintf("Configure %s Command", commandName),
		Fields: tableFields,
	}
}
