package utils

import (
	"github.com/spf13/viper"
)

// Config - Application Properties
type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	DBURL      string `mapstructure:"DBURL"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
