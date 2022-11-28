package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version bool
	rootCmd = &cobra.Command{
		Use:   "task",
		Short: "CLI Task Manager",
		Long:  "CLI Task Manager",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				versionCmd.Run(cmd, nil)
			}
		}}
)

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Print the version of Task.")
}

// Execute executes the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("Execute() Failed: %w", err)
	}

	return nil
}
