package common

import (
	"fmt"

	"github.com/spf13/viper"
)

type WebServiceConfig struct {
	Port int `mapstructure:"port"`
}

type LoggerConfig struct {
	AuditLogFile string `mapstructure:"audit-log-file"`
	LogFile      string `mapstructure:"log-file"`
	LogLevel     string `mapstructure:"log-level"`
}

type AgentConfig struct {
	WindowsRootFolder string `mapstructure:"windows-root-folder"`
}

type Configuration struct {
	WebService WebServiceConfig `mapstructure:"webservice"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Agent      AgentConfig      `mapstructure:"agent"`
}

var Config Configuration

// InitializeConfig loads the configuration from the specified file and unmarshals it.
func InitializeConfig(filePath string) error {
	// Set the configuration file name.
	viper.SetConfigFile(filePath)
	// Set the configuration file type.
	viper.SetConfigType("ini")
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
		return err
	}

	// Unmarshal the configuration values into a Config struct.
	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
		return err
	}

	return nil
}

// GetConfig returns the configuration.
func GetConfig() Configuration {
	return Config
}
