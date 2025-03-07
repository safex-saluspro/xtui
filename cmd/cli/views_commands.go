package cli

import "github.com/spf13/cobra"

func ViewsCmdsList() []*cobra.Command {
	tableCmd := tableViewCmd()

	return []*cobra.Command{
		tableCmd,
	}
}

func tableViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "table-view",
		Aliases: []string{"table", "tableView", "table-view", "viewTable", "view-table"},
		Short:   "Table view for any command",
		Long:    "Table view screen, interactive mode, for any command with flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
