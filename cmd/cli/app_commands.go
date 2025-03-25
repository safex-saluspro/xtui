package cli

import (
	"fmt"
	"github.com/faelmori/logz"
	"github.com/faelmori/xtui/types"
	. "github.com/faelmori/xtui/wrappers"
	"github.com/spf13/cobra"
	"strings"
)

func AppsCmdsList() []*cobra.Command {
	return []*cobra.Command{
		InstallApplicationsCommand(),
	}
}

func InstallApplicationsCommand() *cobra.Command {
	var depList []string
	var path string
	var yes, quiet bool

	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i", "ins", "add"},
		Annotations: GetDescriptions(
			[]string{
				"Install applications and dependencies",
				"Install applications from a file or a repository and add, them to the system"},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(depList) == 0 && len(args) == 0 {
				logz.Error("Empty applications list", map[string]interface{}{
					"context": "InstallApplicationsCommand",
					"error":   "no applications to install",
				})
				return fmt.Errorf("no applications to install")

			}
			newArgs := []string{strings.Join(depList, " "), path, fmt.Sprintf("%t", yes), fmt.Sprintf("%t", quiet)}
			args = append(args, newArgs...)

			availableProperties := getAvailableProperties()
			if len(availableProperties) > 0 {
				adaptedArgs := adaptArgsToProperties(args, availableProperties)
				return InstallDependenciesWithUI(adaptedArgs...)
			}

			// Notification: Starting installation
			DisplayNotification("Starting installation of applications", "info")

			err := InstallDependenciesWithUI(args...)

			if err != nil {
				// Notification: Error during installation
				DisplayNotification(fmt.Sprintf("Error during installation: %s", err.Error()), "error")
				return err
			}

			// Notification: Successful installation
			DisplayNotification("Applications installed successfully", "info")

			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&depList, "application", "a", []string{}, "Applications list to install")
	cmd.Flags().StringVarP(&path, "path", "p", "", "Apps installation path")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic yes to prompts")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")

	return cmd
}

func getAvailableProperties() map[string]string {
	return map[string]string{
		"property1": "value1",
		"property2": "value2",
	}
}

func adaptArgsToProperties(args []string, properties map[string]string) []string {
	adaptedArgs := args
	for key, value := range properties {
		adaptedArgs = append(adaptedArgs, fmt.Sprintf("--%s=%s", key, value))
	}
	return adaptedArgs
}

func NavigateAndExecuteCommand(cmd *cobra.Command, args []string) error {
	// Detect command and its flags
	commandName := cmd.Name()
	flags := cmd.Flags()

	// Display command selection and flag definition in a form
	formConfig := createFormConfig(commandName, flags)
	formResult, err := ShowFormWithNotification(formConfig)
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

func createFormConfig(commandName string, flags *pflag.FlagSet) Config {
	var formFields []FormField

	flags.VisitAll(func(flag *pflag.Flag) {
		formFields = append(formFields, InputField{
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

	return Config{
		Title:  fmt.Sprintf("Configure %s Command", commandName),
		Fields: FormFields{Fields: formFields},
	}
}
