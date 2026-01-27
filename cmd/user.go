package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
	"github.com/pietern/frankie/internal/models"
	"github.com/pietern/frankie/internal/output"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Show user information",
	Long:  `Display information about the logged in user.`,
	RunE:  runUser,
}

func init() {
	rootCmd.AddCommand(userCmd)
}

func runUser(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	resp, err := client.Execute(api.MeQuery, "Me", nil)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}

	var result models.MeResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	user := result.Me

	if getOutputFormat() == "json" {
		return output.JSON(user)
	}

	// Build display data
	keys := []string{
		"Email",
		"Name",
		"Country",
		"Address",
		"Trees",
		"CO2 Compensation",
		"Smart Charging",
		"Smart Trading",
		"Advanced Payment",
		"Member Since",
	}

	pairs := map[string]string{
		"Email":   user.Email,
		"Country": user.CountryCode,
		"Trees":   fmt.Sprintf("%d", user.TreesCount),
	}

	if user.HasCO2Compensation {
		pairs["CO2 Compensation"] = "Yes"
	} else {
		pairs["CO2 Compensation"] = "No"
	}

	if user.ExternalDetails != nil {
		if user.ExternalDetails.Person != nil {
			name := strings.TrimSpace(user.ExternalDetails.Person.FirstName + " " + user.ExternalDetails.Person.LastName)
			if name != "" {
				pairs["Name"] = name
			}
		}
		if user.ExternalDetails.Address != nil {
			if addr := user.ExternalDetails.Address.FormattedAddress(); addr != "" {
				pairs["Address"] = addr
			}
		}
	}

	if user.SmartCharging != nil {
		if user.SmartCharging.IsActivated {
			pairs["Smart Charging"] = "Active"
		} else if user.SmartCharging.IsAvailableInCountry {
			pairs["Smart Charging"] = "Available"
		} else {
			pairs["Smart Charging"] = "Not available"
		}
	}

	if user.SmartTrading != nil {
		if user.SmartTrading.IsActivated {
			pairs["Smart Trading"] = "Active"
		} else if user.SmartTrading.IsAvailableInCountry {
			pairs["Smart Trading"] = "Available"
		} else {
			pairs["Smart Trading"] = "Not available"
		}
	}

	if user.AdvancedPaymentAmount > 0 {
		pairs["Advanced Payment"] = fmt.Sprintf("€%.2f", user.AdvancedPaymentAmount)
	}

	if !user.CreatedAt.IsZero() {
		pairs["Member Since"] = user.CreatedAt.Format("January 2006")
	}

	output.KeyValueOrdered(keys, pairs)
	return nil
}
