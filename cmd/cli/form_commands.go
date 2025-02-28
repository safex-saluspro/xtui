package cli

import (
	"github.com/spf13/cobra"
	"github.com/faelmori/xtui/components"
	"github.com/faelmori/xtui/types"
	"github.com/faelmori/xtui/wrappers"
	tea "github.com/charmbracelet/bubbletea"
	"time"
	"fmt"
)

func FormsCmdsList() []*cobra.Command {
	inputCmd := InputFormCommand()
	loaderCmd := LoaderFormCommand()
	splitCmd := SplitFormCommand()

	return []*cobra.Command{
		inputCmd,
		loaderCmd,
		splitCmd,
	}
}

func InputFormCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "input-form",
		Aliases: []string{"input", "formInput", "inputForm", "formInput", "form-input"},
		Short:   "Form inputs for any command",
		Long:    "Form inputs screen, interactive mode, for any command with flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := types.TuizConfigz{
				Tt: "Dynamic Form",
				Fds: types.TuizFields{
					Tt: "Dynamic Fields",
					Fds: []types.TuizInputz{
						types.TuizInput{
							Ph:  "Username",
							Tp:  "text",
							Val: "",
							Req: true,
							Min: 3,
							Max: 20,
							Err: "Username is required and must be between 3 and 20 characters.",
							Vld: func(value string) error {
								if len(value) < 3 || len(value) > 20 {
									return fmt.Errorf("Username must be between 3 and 20 characters.")
								}
								return nil
							},
						},
						types.TuizInput{
							Ph:  "Password",
							Tp:  "password",
							Val: "",
							Req: true,
							Min: 6,
							Max: 20,
							Err: "Password is required and must be between 6 and 20 characters.",
							Vld: func(value string) error {
								if len(value) < 6 || len(value) > 20 {
									return fmt.Errorf("Password must be between 6 and 20 characters.")
								}
								return nil
							},
						},
					},
				},
			}
			_, err := components.KbdzInputs(config)
			return err
		},
	}

	return cmd
}

func LoaderFormCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "loader-form",
		Aliases: []string{"loader", "formLoader", "loaderForm", "formLoader", "form-loader"},
		Short:   "Form loader for any command",
		Long:    "Form loader screen, interactive mode, for any command with flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			messages := make(chan tea.Msg)
			go func() {
				messages <- wrappers.KbdzLoaderMsg{Message: "Loading dynamic properties..."}
				time.Sleep(2 * time.Second)
				messages <- wrappers.KbdzLoaderMsg{Message: "Dynamic properties loaded successfully."}
				time.Sleep(1 * time.Second)
				messages <- wrappers.KbdzLoaderCloseMsg{}
			}()
			return wrappers.StartLoader(messages)
		},
	}

	return cmd
}

func SplitFormCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "split-form",
		Aliases: []string{"splitInput", "splitInputForm", "inputSplit", "inputSplitForm", "split-input-form"},
		Short:   "Split form inputs for any command",
		Long:    "Split form inputs screen, interactive mode, for any command with flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return components.SplitScreenNew(args...)
		},
	}

	return cmd
}
