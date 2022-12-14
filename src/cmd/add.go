package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

type Task struct {
	ID       string
	IsDone   bool
	Task     string
	Subtasks []string
	Tag      string
	Date     struct {
		Year  int
		Month time.Month
		Day   int
	}
}

func Save(db *bolt.DB, key string, value Task) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		encoded, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to parse task: %w", err)
		}

		err = bucket.Put([]byte(key), (encoded))
		if err != nil {
			return fmt.Errorf("failed to Write task to DB: %w", err)
		}

		return nil
	})
}

var (
	newtask = Task{}
	addCmd  = &cobra.Command{
		Use:   "add",
		Short: "Add new task.",
		Long:  "Add new task to task manager.",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				log.Error().Msg("No arguments provided for \"add\" command")
				os.Exit(1)
			}

			newtask.Task = strings.Join(args, " ")
			newtask.Date.Year, newtask.Date.Month, newtask.Date.Day = time.Now().Date()
			newtask.IsDone = false

			useAlphabet := "12345678"
			idLength := 4
			var err error
			newtask.ID, err = nanoid.Generate(useAlphabet, idLength)
			if err != nil {
				log.Error().Msgf("Failed to add task: %w", err)
				os.Exit(1)
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

			err = Save(db, newtask.ID, newtask)
			if err != nil {
				log.Error().Msgf("Failed to save task: %w", err)
				os.Exit(1)
			}

			switch len(newtask.Subtasks) {
			case 0:
				log.Info().Msgf("Added task \"%s\" to [%s] list.\n", newtask.Task, newtask.Tag)
			case 1:
				log.Info().Msgf("Added task \"%s\" and 1 subtask to [%s] list.  \n", newtask.Task, newtask.Tag)
			default:
				log.Info().Msgf("Added task \"%s\" and %d subtasks to [%s] list.  \n", newtask.Task, len(newtask.Subtasks), newtask.Tag)
			}

		}}
)

func init() {
	addCmd.Flags().StringSliceVarP(&newtask.Subtasks, "subtask", "s", nil, "Add subtask")
	addCmd.Flags().StringVarP(&newtask.Tag, "tag", "t", "Inbox", "Tag your task")
	rootCmd.AddCommand(addCmd)
}
