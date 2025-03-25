package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/faelmori/xtui/components"
	"github.com/faelmori/xtui/types"
	"github.com/faelmori/xtui/wrappers"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"testing"
	"time"
)

func FormsCmdsList() []*cobra.Command {
	inputCmd := InputFormCommand()
	loaderCmd := LoaderFormCommand()

	return []*cobra.Command{
		inputCmd,
		loaderCmd,
	}
}

func InputFormCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "input-form",
		Aliases: []string{"input", "formInput", "inputForm", "formInput", "form-input"},
		Short:   "Form inputs for any command",
		Long:    "Form inputs screen, interactive mode, for any command with flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return NavigateAndExecuteFormCommand(cmd, args)
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
				messages <- wrappers.LoaderMsg{Message: "Loading dynamic properties..."}
				time.Sleep(2 * time.Second)
				messages <- wrappers.LoaderMsg{Message: "Dynamic properties loaded successfully."}
				time.Sleep(1 * time.Second)
				messages <- wrappers.LoaderCloseMsg{}
			}()
			return wrappers.StartLoader(messages)
		},
	}

	return cmd
}

// Unit tests for InputFormCommand
func TestInputFormCommand(t *testing.T) {
	cmd := InputFormCommand()
	if cmd.Use != "input-form" {
		t.Errorf("expected 'input-form', got '%s'", cmd.Use)
	}
	if cmd.Short != "Form inputs for any command" {
		t.Errorf("expected 'Form inputs for any command', got '%s'", cmd.Short)
	}
	if cmd.Long != "Form inputs screen, interactive mode, for any command with flags" {
		t.Errorf("expected 'Form inputs screen, interactive mode, for any command with flags', got '%s'", cmd.Long)
	}
}

// Unit tests for LoaderFormCommand
func TestLoaderFormCommand(t *testing.T) {
	cmd := LoaderFormCommand()
	if cmd.Use != "loader-form" {
		t.Errorf("expected 'loader-form', got '%s'", cmd.Use)
	}
	if cmd.Short != "Form loader for any command" {
		t.Errorf("expected 'Form loader for any command', got '%s'", cmd.Short)
	}
	if cmd.Long != "Form loader screen, interactive mode, for any command with flags" {
		t.Errorf("expected 'Form loader screen, interactive mode, for any command with flags', got '%s'", cmd.Long)
	}
}

func NavigateAndExecuteFormCommand(cmd *cobra.Command, args []string) error {
	// Detect command and its flags
	commandName := cmd.Name()
	flags := cmd.Flags()

	// Display command selection and flag definition in a form
	formConfig := createFormConfig(commandName, flags)
	formResult, err := components.ShowFormWithNotification(formConfig)
	if err != nil {
		return err
	}

	// Set flag values based on form input
	for key, value := range formResult {
		if err := cmd.Flags().Set(key, value); err != nil {
			return err
		}
	}

	// Execute the command
	return cmd.Execute()
}

func createFormConfig(commandName string, flags *pflag.FlagSet) types.Config {
	var formFields []types.FormField

	flags.VisitAll(func(flag *pflag.Flag) {
		formFields = append(formFields, &types.InputField{
			Ph:  flag.Name,
			Tp:  "text",
			Val: flag.Value.String(),
			Req: false,
			Min: 0,
			Max: 100,
			Err: "",
			Vld: func(value string) error { return nil },
		})
	})

	return types.Config{
		Title:  fmt.Sprintf("Configure %s Command", commandName),
		Fields: types.FormFields{Fields: formFields},
	}
}
