package util

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	DBDriver    string `mapstructure:"DB_DRIVER"`
	DBSource    string `mapstructure:"DB_SOURCE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(fileName, filePath string) (config Config, err error) {
	viper.AddConfigPath(filePath)
	viper.SetConfigName(fileName)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return config, fmt.Errorf("config file not found in path %s", filePath)
		}
	}

	viper.AutomaticEnv()
	return config, viper.Unmarshal(&config)
}
