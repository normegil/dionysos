package commands

import (
	"github.com/normegil/dionysos/internal/configuration"
	logCfg "github.com/normegil/dionysos/internal/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string

func init() {
	logCfg.Init()

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file")

	loggingColorKey := configuration.KeyColorizedLogging
	RootCmd.PersistentFlags().Bool(loggingColorKey.CommandLine.Name, false, loggingColorKey.Description)
	if err := viper.BindPFlag(loggingColorKey.Name, RootCmd.PersistentFlags().Lookup(loggingColorKey.CommandLine.Name)); err != nil {
		log.Fatal().Err(err).Str("paramName", loggingColorKey.Name).Msg("Binding parameter")
	}
}

var RootCmd = &cobra.Command{
	Use:   "dionysos",
	Short: "Dionysos is a software designed to manage household stock of foods and other supply.",
	Long:  `Dionysos is a software designed to manage household stock of foods and other supply.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Could not execute root command")
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/dionysos")
		viper.AddConfigPath("$XDG_CONFIG_HOME/dionysos")
		viper.AddConfigPath("$HOME/.dionysos")
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
