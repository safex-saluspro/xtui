package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	. "github.com/faelmori/xtui/services"
)

func ServicesCmds() []*cobra.Command {
	runAsDaemon := RunAsDaemonCommand()

	return []*cobra.Command{runAsDaemon}
}

func RunAsDaemonCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "daemonize",
		Aliases: []string{"daemon", "runDaemon", "run-daemon", "runAsDaemon"},
		Short:   "Run command as daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			execCmdAsDaemon, _ := cmd.Flags().GetString("exec")
			if execCmdAsDaemon == "" && len(args) == 0 {
				return fmt.Errorf("no command to execute as daemon")
			}
			argsList := append(args, execCmdAsDaemon)

			// Dynamic adaptation logic
			availableProperties := getAvailableProperties()
			if len(availableProperties) > 0 {
				adaptedArgs := adaptArgsToProperties(argsList, availableProperties)
				return Daemonize(nil, adaptedArgs...)
			}

			return Daemonize(nil, argsList...)
		},
	}

	cmd.Flags().StringP("exec", "e", "", "Command to execute as daemon")
	cmd.Flags().StringP("name", "n", "", "Daemon name")
	cmd.Flags().StringP("pidfile", "p", "", "PID file path")
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
