package common

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Webservice struct {
		Port string
	}
}

func GetConfig() (config *Config, err error) {
	// Set the configuration file name.
	viper.SetConfigName("config")

	// Set the configuration file type.
	viper.SetConfigType("ini") // or "json", "toml", etc.

	// Set the configuration file search paths.
	viper.AddConfigPath(".")

	// Enable automatic configuration file searching and reading.
	if err := viper.ReadInConfig(); err != nil {
		// Handle the error if the configuration file is not found.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found.")
		} else {
			fmt.Printf("Error reading config file: %s\n", err)
		}
		return nil, err
	}

	// Unmarshal the configuration values into a Config struct.
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
		return nil, err
	}

	return config, nil
}
