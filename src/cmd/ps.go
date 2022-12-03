package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	psCmd = &cobra.Command{
		Use:   "ps",
		Short: "List task.",
		Long:  "List tasks added to task manager.",
		Run: func(cmd *cobra.Command, args []string) {
			t := table.NewWriter()
			t.SetStyle(table.StyleDouble)
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Task", "Sub Tasks", "Tag", "Date Added"})
			// It will be created if it doesn't exist.
			db, err := bolt.Open("tasks.db", 0600, &bolt.Options{Timeout: 10 * time.Second})
			if err != nil {
				log.Fatal().Msgf("Failed to open DB: %w", err)
			}
			defer db.Close()

			err = db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte("Tasks"))
				if bucket == nil {
					log.Info().Msg("Nothing to show.")
					os.Exit(0)
				}

				err = bucket.ForEach(func(k, v []byte) error {

					task := Task{}
					if err = json.Unmarshal(v, &task); err != nil {
						log.Error().Msgf("failed to unmarshal: %w", err)
					}

					date := fmt.Sprintf("%d/%d/%d", task.Date.Day, task.Date.Month, task.Date.Year)
					t.AppendRow(table.Row{task.Task, task.Subtasks, task.Tag, date})
					t.AppendSeparator()
					return err
				})
				t.Render()
				return err
			})
			if err != nil {
				log.Error().Msgf("failed to list tasks: %w", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(psCmd)
}
