package cli

import (
	"github.com/spf13/cobra"
)

func ServicesCmds() []*cobra.Command {
	runAsDaemon := runAsDaemonCmd()

	return []*cobra.Command{runAsDaemon}
}

func runAsDaemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "daemonize",
		Aliases: []string{"daemon", "runDaemon", "run-daemon", "runAsDaemon"},
		Short:   "Run command as daemon",
		//RunE: func(cmd *cobra.Command, args []string) error {
		//	execCmdAsDaemon, _ := cmd.Flags().GetString("exec")
		//	if execCmdAsDaemon == "" && len(args) == 0 {
		//		return fmt.Errorf("no command to execute as daemon")
		//	}
		//	argsList := append(args, execCmdAsDaemon)
		//	return Daemonize(argsList...)
		//},
	}

	cmd.Flags().StringP("exec", "e", "", "Command to execute as daemon")
	cmd.Flags().StringP("name", "n", "", "Daemon name")
	cmd.Flags().StringP("pidfile", "p", "", "PID file path")
	cmd.Flags().BoolP("quiet", "q", false, "Quiet mode")

	return cmd
}
