package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/auth"
	"github.com/pietern/frankie/internal/config"
	"github.com/pietern/frankie/internal/output"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display the current authentication status, including login state and token expiration.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

// StatusInfo holds the status information for JSON output
type StatusInfo struct {
	LoggedIn     bool       `json:"logged_in"`
	Email        string     `json:"email,omitempty"`
	TokenExpiry  *time.Time `json:"token_expiry,omitempty"`
	TokenExpired bool       `json:"token_expired,omitempty"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	creds, err := auth.LoadCredentials()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	if creds == nil {
		return showNotLoggedIn()
	}

	// Parse token info
	expiry, _ := auth.ParseJWTExpiration(creds.AuthToken)
	email := extractEmailFromToken(creds.AuthToken)
	expired := auth.IsTokenExpired(creds.AuthToken, 0)

	if getOutputFormat() == "json" {
		info := StatusInfo{
			LoggedIn:     !expired,
			Email:        email,
			TokenExpired: expired,
		}
		if !expiry.IsZero() {
			info.TokenExpiry = &expiry
		}
		return output.JSON(info)
	}

	// Table output
	if expired {
		fmt.Println("Session expired")
		fmt.Println()
		fmt.Println("Run 'frankie login' to authenticate.")
		return nil
	}

	fmt.Println("Logged in")
	fmt.Println()

	keys := []string{"Email", "Token", "Expiry"}
	pairs := map[string]string{}

	if email != "" {
		pairs["Email"] = email
	}

	if !expiry.IsZero() {
		pairs["Token"] = formatTimeRemaining(expiry)
		pairs["Expiry"] = expiry.Local().Format("2006-01-02 15:04:05")
	}

	output.KeyValueOrdered(keys, pairs)
	return nil
}

func showNotLoggedIn() error {
	if getOutputFormat() == "json" {
		return output.JSON(StatusInfo{LoggedIn: false})
	}

	fmt.Println("Not logged in")
	fmt.Println()
	fmt.Printf("Credentials file: %s\n", config.GetCredentialsPath())
	fmt.Println()
	fmt.Println("Run 'frankie login' to authenticate.")
	return nil
}

// extractEmailFromToken attempts to extract email from JWT claims
func extractEmailFromToken(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ""
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}

	// Try common email claim fields
	for _, field := range []string{"email", "sub"} {
		if val, ok := claims[field]; ok {
			if str, ok := val.(string); ok && strings.Contains(str, "@") {
				return str
			}
		}
	}

	return ""
}

// formatTimeRemaining formats the time until expiry
func formatTimeRemaining(expiry time.Time) string {
	remaining := time.Until(expiry)

	if remaining <= 0 {
		return "expired"
	}

	hours := remaining.Hours()
	if hours >= 24 {
		days := int(hours / 24)
		if days == 1 {
			return "expires in 1 day"
		}
		return fmt.Sprintf("expires in %d days", days)
	}

	if hours >= 1 {
		return fmt.Sprintf("expires in %.1f hours", hours)
	}

	minutes := int(remaining.Minutes())
	if minutes == 1 {
		return "expires in 1 minute"
	}
	return fmt.Sprintf("expires in %d minutes", minutes)
}
