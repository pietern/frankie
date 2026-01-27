package config

import (
	"os"
	"path/filepath"
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

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	return os.MkdirAll(GetConfigDir(), 0700)
}
