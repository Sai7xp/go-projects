/*
* Created on 06 May 2024
* @author Sai Sumanth
 */
package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config - stores all configurations of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig - reads configuration from file or env variable
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Println("fatal error config file: %w", err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
