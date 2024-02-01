package cmd

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gondos/util"
)

func init() {
	godotenv.Load()
	// configure global logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339

	if util.IsDevel() {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
		log.Debug().Msg("Running in Development")
	}

	// this is important to know what timezone is currently used
	log.Info().Msgf("Used time: %s", time.Now())
}
