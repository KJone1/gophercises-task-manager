package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Task Manager",
		Short: "CLI Task Manager",
		Long:  "CLI Task Manager",
	}
)

// Execute executes the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("Execute() Failed: %w", err)
	}

	return nil
}
