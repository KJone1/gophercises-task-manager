package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var (
	subtask = []string{}
	tag     string
	addCmd  = &cobra.Command{
		Use:   "add",
		Short: "Add new task.",
		Long:  "Add new task to task manager.",
		Run: func(cmd *cobra.Command, args []string) {
			task := strings.Join(args, " ")
			// subtask := strings.Join(args, " ")
			if len(args) == 0 {
				log.Error().Msg("No arguments provided for \"Add\" command")
				os.Exit(1)
			}
			if len(subtask) == 0 {
				fmt.Printf("Added task \"%s\" to [%s] list.\n", task, tag)

				return
			}
			fmt.Printf("Added task \"%s\" to [%s] list with %d subtasks.\n", task, tag, len(subtask))
			for i, k := range subtask {
				fmt.Println(i, k)
			}
		},
	}
)

func init() {
	addCmd.Flags().StringSliceVarP(&subtask, "subtask", "s", nil, "Add subtask")
	addCmd.Flags().StringVarP(&tag, "tag", "t", "Inbox", "Tag your task")
	rootCmd.AddCommand(addCmd)
}
