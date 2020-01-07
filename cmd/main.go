package main

import (
	"github.com/normegil/dionysos/cmd/commands"
	"github.com/normegil/dionysos/internal/configuration"
	logCfg "github.com/normegil/dionysos/internal/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var cfgFile string //nolint:gochecknoglobals // Satisfying cobra interface 'OnInitialize' require this global variable

func main() {
	logCfg.Init()
	cobra.OnInitialize(initConfig)

	root, err := commands.Root()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not execute command")
	}
	root.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file")
	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Could not execute command")
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/dionysos")
		viper.AddConfigPath("$XDG_CONFIG_HOME" + string(os.PathSeparator) + "dionysos")
		viper.AddConfigPath("$HOME" + string(os.PathSeparator) + ".dionysos")
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		viper.SetConfigName("dionysos")
	}

	viper.SetEnvPrefix("DIONYSOS_")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, isNotFound := err.(viper.ConfigFileNotFoundError); !isNotFound {
			log.Fatal().Err(err).Msg("could not read configuration")
		}
	}

	logCfg.Configure(viper.GetBool(configuration.KeyColorizedLogging.Name))
}
