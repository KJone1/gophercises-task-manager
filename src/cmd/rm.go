package cmd

import (
	"encoding/json"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove task.",
		Long:  "Remove task from task manager.",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				log.Error().Msg("No arguments provided for \"rm\" command")
				os.Exit(1)
			}

			taskID := strings.Join(args, " ")

			// It will be created if it doesn't exist.
			var filePermission fs.FileMode = 0600
			timeout := 10 * time.Second
			db, err := bolt.Open("tasks.db", filePermission, &bolt.Options{Timeout: timeout})
			if err != nil {
				log.Error().Msgf("Failed to open DB: %w", err)
				os.Exit(1)
			}
			defer db.Close()

			err = db.Update(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte("Tasks"))
				if bucket == nil {
					log.Info().Msg("Nothing to delete.")
					os.Exit(0)
				}

				readValue := bucket.Get([]byte(taskID))
				task := Task{}

				if err = json.Unmarshal(readValue, &task); err != nil {
					log.Error().Msgf("Failed to unmarshal: %w", err)
					os.Exit(1)
				}

				if err = bucket.Delete([]byte(taskID)); err != nil {
					log.Error().Msgf("Failed to delete task: %w", err)
					os.Exit(1)
				}

				log.Info().Msgf("Removed %s", task.Task)

				return err
			})
			if err != nil {
				log.Error().Msgf("Failed to update database: %w", err)
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(rmCmd)
}
