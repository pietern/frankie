package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pietern/frankie/internal/config"
	"github.com/pietern/frankie/internal/models"
)

// LoadCredentials reads stored credentials from disk
func LoadCredentials() (*models.Credentials, error) {
	path := config.GetCredentialsPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds models.Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

// SaveCredentials writes credentials to disk with secure permissions
func SaveCredentials(creds *models.Credentials) error {
	if err := config.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	path := config.GetCredentialsPath()

	// Write with secure permissions (0600 = owner read/write only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

// DeleteCredentials removes stored credentials
func DeleteCredentials() error {
	path := config.GetCredentialsPath()

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete credentials: %w", err)
	}

	return nil
}

// CredentialsExist checks if credentials file exists
func CredentialsExist() bool {
	path := config.GetCredentialsPath()
	_, err := os.Stat(path)
	return err == nil
}

// GetConfigDir returns the config directory, creating it if necessary
func GetConfigDir() string {
	return filepath.Dir(config.GetCredentialsPath())
}
