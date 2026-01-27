package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	AppName = "frankie"
)

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".config/frankie"
	}
	return filepath.Join(home, ".config", AppName)
}

// GetCredentialsPath returns the path to the credentials file
func GetCredentialsPath() string {
	return filepath.Join(GetConfigDir(), "credentials.json")
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.yaml")
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	configDir := GetConfigDir()
	return os.MkdirAll(configDir, 0700)
}

// Config holds application configuration
type Config struct {
	Output      string `mapstructure:"output"`
	DefaultSite string `mapstructure:"default_site"`
	Country     string `mapstructure:"country"`
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	return &Config{
		Output:      viper.GetString("output"),
		DefaultSite: viper.GetString("default_site"),
		Country:     viper.GetString("country"),
	}
}

// SetDefaults sets default configuration values
func SetDefaults() {
	viper.SetDefault("output", "table")
	viper.SetDefault("country", "NL")
}

// SaveConfig saves the current configuration to file
func SaveConfig() error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}
	return viper.WriteConfigAs(GetConfigPath())
}
