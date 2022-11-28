package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	verbose    bool
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of Task.",
		Long:  `All software has versions. This is Task's.`,
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				fmt.Println("CLI Task Manager Version 0.0.1")

				return
			}
			fmt.Println("Task - v0.0.1")
		},
	}
)

func init() {
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose version")
	rootCmd.AddCommand(versionCmd)
}
