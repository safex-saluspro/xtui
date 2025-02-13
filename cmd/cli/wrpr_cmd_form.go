package cli

import "github.com/spf13/cobra"

func FormsCmdsList() []*cobra.Command {
	inputCmd := inputsFormCmd()
	loaderCmd := loaderFormCmd()
	splitCmd := splitFormCmd()

	return []*cobra.Command{
		inputCmd,
		loaderCmd,
		splitCmd,
	}
}

func inputsFormCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "input-form",
		Aliases: []string{"input", "formInput", "inputForm", "formInput", "form-input"},
		Short:   "Form inputs for any command",
		Long:    "Form inputs screen, interactive mode, for any command with flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func loaderFormCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "loader-form",
		Aliases: []string{"loader", "formLoader", "loaderForm", "formLoader", "form-loader"},
		Short:   "Form loader for any command",
		Long:    "Form loader screen, interactive mode, for any command with flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func splitFormCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "split-form",
		Aliases: []string{"splitInput", "splitInputForm", "inputSplit", "inputSplitForm", "split-input-form"},
		Short:   "Split form inputs for any command",
		Long:    "Split form inputs screen, interactive mode, for any command with flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
