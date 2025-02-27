package cli

import (
	"fmt"
	"github.com/faelmori/kbx/mods/logz"
	. "github.com/faelmori/xtui/wrappers"
	"github.com/spf13/cobra"
	"strings"
)

func AppsCmdsList() []*cobra.Command {
	instAppsCmd := InstallApplicationsCommand()

	return []*cobra.Command{instAppsCmd}
}

func InstallApplicationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app-install",
		Aliases: []string{"install", "appInstall", "installApp", "aptInstall", "apt-install", "depInstall", "dep-install"},
		Short:   "Install applications and dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			depList, _ := cmd.Flags().GetStringArray("application")
			path, _ := cmd.Flags().GetString("path")
			yes, _ := cmd.Flags().GetBool("yes")
			quiet, _ := cmd.Flags().GetBool("quiet")
			if len(depList) == 0 && len(args) == 0 {
				return logz.ErrorLog("Empty applications list", "ui")
			}
			newArgs := []string{strings.Join(depList, " "), path, fmt.Sprintf("%t", yes), fmt.Sprintf("%t", quiet)}
			args = append(args, newArgs...)

			// Dynamic adaptation logic
			availableProperties := getAvailableProperties()
			if len(availableProperties) > 0 {
				adaptedArgs := adaptArgsToProperties(args, availableProperties)
				return InstallDependenciesWithUI(adaptedArgs...)
			}

			return InstallDependenciesWithUI(args...)
		},
	}

	cmd.Flags().StringArrayP("application", "a", []string{}, "Applications list to install")
	cmd.Flags().StringP("path", "p", "", "Apps installation path")
	cmd.Flags().BoolP("yes", "y", false, "Automatic yes to prompts")
	cmd.Flags().BoolP("quiet", "q", false, "Quiet mode")

	return cmd
}

// Helper function to get available properties
func getAvailableProperties() map[string]string {
	// Implement logic to fetch available properties
	return map[string]string{
		"property1": "value1",
		"property2": "value2",
	}
}

// Helper function to adapt arguments based on available properties
func adaptArgsToProperties(args []string, properties map[string]string) []string {
	// Implement logic to adapt arguments based on properties
	adaptedArgs := args
	for key, value := range properties {
		adaptedArgs = append(adaptedArgs, fmt.Sprintf("--%s=%s", key, value))
	}
	return adaptedArgs
}
