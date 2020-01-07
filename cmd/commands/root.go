package commands

import (
	"fmt"
	"github.com/normegil/dionysos/internal/configuration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Root() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "dionysos",
		Short: "Dionysos is a software designed to manage household stock of foods and other supply.",
		Long:  `Dionysos is a software designed to manage household stock of foods and other supply.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); nil != err {
				log.Fatal().Err(err).Msg("could not print help message")
			}
		},
	}

	loggingColorKey := configuration.KeyColorizedLogging
	rootCmd.PersistentFlags().Bool(loggingColorKey.CommandLine.Name, false, loggingColorKey.Description)
	if err := viper.BindPFlag(loggingColorKey.Name, rootCmd.PersistentFlags().Lookup(loggingColorKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", loggingColorKey.Name, err)
	}

	listenCmd, err := listen()
	if err != nil {
		return nil, fmt.Errorf("creating 'listen' command: %w", err)
	}
	rootCmd.AddCommand(listenCmd)

	return rootCmd, nil
}
