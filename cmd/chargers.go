package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
	"github.com/pietern/frankie/internal/models"
	"github.com/pietern/frankie/internal/output"
)

var chargersCmd = &cobra.Command{
	Use:   "chargers",
	Short: "Show smart chargers",
	Long:  `Display smart charger information and status.`,
	RunE:  runChargers,
}

func init() {
	rootCmd.AddCommand(chargersCmd)
}

func runChargers(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	resp, err := client.Execute(api.EnodeChargersQuery, "EnodeChargers", nil)
	if err != nil {
		if errors.Is(err, api.ErrSmartChargingNotEnabled) {
			fmt.Println("Smart charging is not enabled for your account")
			return nil
		}
		return fmt.Errorf("failed to fetch chargers: %w", err)
	}

	var result models.EnodeChargersResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	chargers := result.EnodeChargers
	if len(chargers) == 0 {
		fmt.Println("No smart chargers found")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(chargers)
	}

	headers := []string{"Brand", "Model", "Status", "Smart", "Plugged In", "Charging", "Rate"}
	var rows [][]string

	for _, c := range chargers {
		brand := ""
		model := ""
		if c.Information != nil {
			brand = c.Information.Brand
			model = c.Information.Model
		}

		status := "Offline"
		if c.IsReachable {
			status = "Online"
		}

		smart := "No"
		if c.CanSmartCharge {
			smart = "Yes"
		}

		pluggedIn := "-"
		charging := "-"
		rate := "-"
		if c.ChargeState != nil {
			if c.ChargeState.IsPluggedIn {
				pluggedIn = "Yes"
			} else {
				pluggedIn = "No"
			}
			if c.ChargeState.IsCharging {
				charging = "Yes"
				rate = fmt.Sprintf("%.1f kW", c.ChargeState.ChargeRate)
			} else {
				charging = "No"
			}
		}

		rows = append(rows, []string{brand, model, status, smart, pluggedIn, charging, rate})
	}

	table := output.Table(headers)
	table.AppendBulk(rows)
	table.Render()

	return nil
}
