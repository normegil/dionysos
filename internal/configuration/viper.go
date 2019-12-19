package configuration

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func NewConfiguration() (*viper.Viper, error) {
	v := viper.New()

	def := defaults()
	for key, val := range def {
		v.SetDefault(key.String(), val)
	}

	v.SetConfigType("yaml")
	v.SetConfigName("dionysos")
	pathToCfg := paths()
	for _, path := range pathToCfg {
		v.AddConfigPath(path)
	}

	v.SetEnvPrefix("DIONYSOS_")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}
	return v, nil
}
