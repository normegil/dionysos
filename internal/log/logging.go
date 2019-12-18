package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func Configure() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}
