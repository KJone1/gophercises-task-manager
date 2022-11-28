package main

import (
	"os"

	"github.com/KJone1/gophercises-task-manager/src/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := cmd.Execute(); err != nil {
		log.Fatal().Msg("Failed to execute")
	}
}
