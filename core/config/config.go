package config

import (
	"github.com/spf13/viper"
	"strings"
)

func LoadConfig(filename string, configPath []string) error {

	viper.SetEnvPrefix(filename)
	viper.AutomaticEnv()
	viper.SetConfigName(filename)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AddConfigPath(".")
	for _, c := range configPath {
		viper.AddConfigPath(c)
	}

	return viper.ReadInConfig() // Find and read the config file
}
