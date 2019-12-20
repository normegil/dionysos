package main

import (
	"github.com/normegil/dionysos/cmd/commands"
	zLog "github.com/rs/zerolog/log"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		zLog.Fatal().Err(err).Msg("Could not execute command")
	}
}
