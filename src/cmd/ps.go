package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	showAll bool
	psCmd   = &cobra.Command{
		Use:   "ps",
		Short: "List task.",
		Long:  "List tasks added to task manager.",
		Run: func(cmd *cobra.Command, args []string) {
			// Configure table params
			t := table.NewWriter()
			t.SetStyle(table.StyleDouble)
			t.SetOutputMirror(os.Stdout)

			if showAll {
				t.AppendHeader(table.Row{"ID", "Task", "Sub Tasks", "Tag", "Date Added", "Status"})
				t.SetColumnConfigs([]table.ColumnConfig{
					{
						Name:  "Sub Tasks",
						Align: text.AlignLeft,
					},
				})
			} else {
				t.AppendHeader(table.Row{"Task", "Sub Tasks", "Tag", "Date Added"})
				t.SetColumnConfigs([]table.ColumnConfig{
					{
						Name:  "Sub Tasks",
						Align: text.AlignCenter,
					},
				})
			}

			// It will be created if it doesn't exist.
			var filePermission fs.FileMode = 0600
			timeout := 10 * time.Second
			db, err := bolt.Open("tasks.db", filePermission, &bolt.Options{Timeout: timeout})
			if err != nil {
				log.Error().Msgf("Failed to open DB: %w", err)
				os.Exit(1)
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
						log.Error().Msgf("Failed to unmarshal: %w", err)
						os.Exit(1)
					}

					if showAll {
						date := fmt.Sprintf("%d/%d/%d", task.Date.Day, task.Date.Month, task.Date.Year)
						subtasks := strings.Join(task.Subtasks, ", ")
						var status string
						if task.IsDone {
							status = "Done"
						} else {
							status = "In Progress"
						}
						row := (table.Row{task.ID, task.Task, subtasks, task.Tag, date, status})

						t.AppendRow(row)
						t.AppendSeparator()

						return err
					}

					if task.IsDone {
						return err
					}

					date := fmt.Sprintf("%d/%d/%d", task.Date.Day, task.Date.Month, task.Date.Year)
					var subtasksCount int
					if len(task.Subtasks) == 0 {
						subtasksCount = 0
					} else {
						subtasksCount = len(task.Subtasks)
					}
					row := (table.Row{task.Task, subtasksCount, task.Tag, date})

					t.AppendRow(row)
					t.AppendSeparator()

					return err
				})

				return err
			})
			if err != nil {
				log.Error().Msgf("Failed to list tasks: %w", err)
			}

			t.Render()
		},
	}
)

func init() {
	psCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Print all tasks")
	rootCmd.AddCommand(psCmd)
}
