package models

import "time"

// Credentials holds the authentication tokens
type Credentials struct {
	AuthToken    string    `json:"auth_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
