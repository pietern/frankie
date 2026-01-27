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

var vehiclesCmd = &cobra.Command{
	Use:   "vehicles",
	Short: "Show smart vehicles",
	Long:  `Display smart vehicle information and charging status.`,
	RunE:  runVehicles,
}

func init() {
	rootCmd.AddCommand(vehiclesCmd)
}

func runVehicles(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	resp, err := client.Execute(api.EnodeVehiclesQuery, "EnodeVehicles", nil)
	if err != nil {
		if errors.Is(err, api.ErrSmartChargingNotEnabled) {
			fmt.Println("Smart charging is not enabled for your account")
			return nil
		}
		return fmt.Errorf("failed to fetch vehicles: %w", err)
	}

	var result models.EnodeVehiclesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	vehicles := result.EnodeVehicles
	if len(vehicles) == 0 {
		fmt.Println("No smart vehicles found")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(vehicles)
	}

	headers := []string{"Brand", "Model", "Status", "Battery", "Range", "Charging", "Rate"}
	var rows [][]string

	for _, v := range vehicles {
		brand := ""
		model := ""
		if v.Information != nil {
			brand = v.Information.Brand
			model = v.Information.Model
		}

		status := "Offline"
		if v.IsReachable {
			status = "Online"
		}

		battery := "-"
		rangeKm := "-"
		charging := "-"
		rate := "-"
		if v.ChargeState != nil {
			if v.ChargeState.BatteryLevel > 0 {
				battery = fmt.Sprintf("%.0f%%", v.ChargeState.BatteryLevel)
			}
			if v.ChargeState.Range > 0 {
				rangeKm = fmt.Sprintf("%.0f km", v.ChargeState.Range)
			}
			if v.ChargeState.IsCharging {
				charging = "Yes"
				if v.ChargeState.ChargeRate > 0 {
					rate = fmt.Sprintf("%.1f kW", v.ChargeState.ChargeRate)
				}
			} else {
				charging = "No"
			}
		}

		rows = append(rows, []string{brand, model, status, battery, rangeKm, charging, rate})
	}

	table := output.Table(headers)
	table.AppendBulk(rows)
	table.Render()

	return nil
}
