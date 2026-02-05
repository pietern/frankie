package auth

import (
	"fmt"
	"time"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/models"
)

const (
	// TokenRefreshMargin is the time before expiration when we should refresh
	TokenRefreshMargin = 5 * time.Minute
)

// Manager handles authentication operations
type Manager struct {
	client *api.Client
}

// NewManager creates a new auth manager
func NewManager(client *api.Client) *Manager {
	return &Manager{client: client}
}

// Login authenticates with email and password
func (m *Manager) Login(email, password string) error {
	authToken, refreshToken, err := m.client.Login(email, password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	expiresAt, _ := ParseJWTExpiration(authToken)

	creds := &models.Credentials{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	if err := SaveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

// Logout clears stored credentials
func (m *Manager) Logout() error {
	return DeleteCredentials()
}

// GetValidToken returns a valid auth token, refreshing if necessary
func (m *Manager) GetValidToken() (string, error) {
	creds, err := LoadCredentials()
	if err != nil {
		return "", err
	}
	if creds == nil {
		return "", fmt.Errorf("not logged in")
	}

	// Check if token needs refresh
	if IsTokenExpired(creds.AuthToken, TokenRefreshMargin) {
		newCreds, err := m.RefreshToken(creds)
		if err != nil {
			return "", fmt.Errorf("token refresh failed: %w", err)
		}
		creds = newCreds
	}

	return creds.AuthToken, nil
}

// RefreshToken renews the auth token using the refresh token
func (m *Manager) RefreshToken(creds *models.Credentials) (*models.Credentials, error) {
	newAuthToken, newRefreshToken, err := m.client.RenewToken(creds.AuthToken, creds.RefreshToken)
	if err != nil {
		return nil, err
	}

	expiresAt, _ := ParseJWTExpiration(newAuthToken)

	newCreds := &models.Credentials{
		AuthToken:    newAuthToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}

	if err := SaveCredentials(newCreds); err != nil {
		return nil, fmt.Errorf("failed to save refreshed credentials: %w", err)
	}

	return newCreds, nil
}

// EnsureAuthenticated loads credentials and sets up the client
func (m *Manager) EnsureAuthenticated() error {
	token, err := m.GetValidToken()
	if err != nil {
		return err
	}
	m.client.SetAuthToken(token)
	return nil
}

// IsLoggedIn checks if user has valid credentials
func IsLoggedIn() bool {
	creds, err := LoadCredentials()
	if err != nil || creds == nil {
		return false
	}
	return !IsTokenExpired(creds.AuthToken, 0)
}
