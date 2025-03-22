package cli

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/faelmori/xtui/components"
	t "github.com/faelmori/xtui/types"
	"github.com/spf13/cobra"
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
			config := t.FormConfig{
				Title: "Sample Table",
				Fields: []t.Field{
					t.InputField{
						Ph:  "Column1",
						Tp:  "text",
						Val: "Value1",
					},
					t.InputField{
						Ph:  "Column2",
						Tp:  "text",
						Val: "Value2",
					},
				},
			}
			customStyles := map[string]lipgloss.Color{
				"Info":    lipgloss.Color("#75FBAB"),
				"Warning": lipgloss.Color("#FDFF90"),
				"Error":   lipgloss.Color("#FF7698"),
				"Debug":   lipgloss.Color("#929292"),
			}
			return components.StartTableScreen(config, customStyles)
		},
	}

	return cmd
}

func TestTableViewCmd(t *testing.T) {
	cmd := tableViewCmd()
	if cmd.Use != "table-view" {
		t.Errorf("expected 'table-view', got '%s'", cmd.Use)
	}
	if cmd.Short != "Table view for any command" {
		t.Errorf("expected 'Table view for any command', got '%s'", cmd.Short)
	}
	if cmd.Long != "Table view screen, interactive mode, for any command with flags" {
		t.Errorf("expected 'Table view screen, interactive mode, for any command with flags', got '%s'", cmd.Long)
	}
}
