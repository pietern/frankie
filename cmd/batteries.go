package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
	"github.com/pietern/frankie/internal/models"
	"github.com/pietern/frankie/internal/output"
)

var (
	batteriesDeviceID  string
	batteriesStartDate string
	batteriesEndDate   string
)

var batteriesCmd = &cobra.Command{
	Use:   "batteries",
	Short: "Show smart batteries",
	Long:  `Display smart battery information and trading results.`,
	RunE:  runBatteries,
}

var batteriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all batteries",
	RunE:  runBatteriesList,
}

var batteriesDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Show battery details",
	RunE:  runBatteriesDetails,
}

var batteriesSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Show battery trading sessions",
	RunE:  runBatteriesSessions,
}

func init() {
	rootCmd.AddCommand(batteriesCmd)
	batteriesCmd.AddCommand(batteriesListCmd)
	batteriesCmd.AddCommand(batteriesDetailsCmd)
	batteriesCmd.AddCommand(batteriesSessionsCmd)

	batteriesDetailsCmd.Flags().StringVarP(&batteriesDeviceID, "device", "d", "", "device ID")
	batteriesSessionsCmd.Flags().StringVarP(&batteriesDeviceID, "device", "d", "", "device ID")
	batteriesSessionsCmd.Flags().StringVar(&batteriesStartDate, "start", "", "start date (YYYY-MM-DD)")
	batteriesSessionsCmd.Flags().StringVar(&batteriesEndDate, "end", "", "end date (YYYY-MM-DD)")
}

func runBatteries(cmd *cobra.Command, args []string) error {
	// Default to list
	return runBatteriesList(cmd, args)
}

func runBatteriesList(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	resp, err := client.Execute(api.SmartBatteriesQuery, "SmartBatteries", nil)
	if err != nil {
		if errors.Is(err, api.ErrSmartTradingNotEnabled) {
			fmt.Println("Smart trading is not enabled for your account")
			return nil
		}
		return fmt.Errorf("failed to fetch batteries: %w", err)
	}

	var result models.SmartBatteriesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	batteries := result.SmartBatteries
	if len(batteries) == 0 {
		fmt.Println("No smart batteries found")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(batteries)
	}

	headers := []string{"ID", "Brand", "Capacity", "Max Charge", "Max Discharge", "Provider"}
	var rows [][]string

	for _, b := range batteries {
		rows = append(rows, []string{
			b.ID,
			b.Brand,
			fmt.Sprintf("%.1f kWh", b.Capacity),
			fmt.Sprintf("%.1f kW", b.MaxChargePower),
			fmt.Sprintf("%.1f kW", b.MaxDischargePower),
			b.Provider,
		})
	}

	output.Table(headers, rows)

	return nil
}

func runBatteriesDetails(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	deviceID := batteriesDeviceID
	if deviceID == "" {
		// Try to get the first battery
		var err error
		deviceID, err = getFirstBatteryID(client)
		if err != nil {
			if errors.Is(err, api.ErrSmartTradingNotEnabled) {
				fmt.Println("Smart trading is not enabled for your account")
				return nil
			}
			return err
		}
		if deviceID == "" {
			fmt.Println("No batteries found")
			return nil
		}
	}

	variables := map[string]interface{}{
		"deviceId": deviceID,
	}

	resp, err := client.Execute(api.SmartBatteryDetailsQuery, "SmartBattery", variables)
	if err != nil {
		return fmt.Errorf("failed to fetch battery details: %w", err)
	}

	var result models.SmartBatteryDetailsResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if getOutputFormat() == "json" {
		return output.JSON(result)
	}

	if result.SmartBattery != nil {
		fmt.Printf("Battery: %s (%s)\n", result.SmartBattery.Brand, result.SmartBattery.ID)
		fmt.Printf("Capacity: %.1f kWh\n", result.SmartBattery.Capacity)
		if result.SmartBattery.Settings != nil {
			fmt.Printf("Mode: %s\n", result.SmartBattery.Settings.BatteryMode)
		}
	}

	if result.SmartBatterySummary != nil {
		fmt.Println()
		fmt.Printf("State of Charge: %.0f%%\n", result.SmartBatterySummary.LastKnownStateOfCharge)
		fmt.Printf("Status: %s\n", result.SmartBatterySummary.LastKnownStatus)
		fmt.Printf("Total Result: €%.2f\n", result.SmartBatterySummary.TotalResult)
		fmt.Printf("Last Update: %s\n", formatDate(result.SmartBatterySummary.LastUpdate))
	}

	return nil
}

func runBatteriesSessions(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	deviceID := batteriesDeviceID
	if deviceID == "" {
		var err error
		deviceID, err = getFirstBatteryID(client)
		if err != nil {
			if errors.Is(err, api.ErrSmartTradingNotEnabled) {
				fmt.Println("Smart trading is not enabled for your account")
				return nil
			}
			return err
		}
		if deviceID == "" {
			fmt.Println("No batteries found")
			return nil
		}
	}

	// Default date range: last 30 days
	endDate := batteriesEndDate
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	startDate := batteriesStartDate
	if startDate == "" {
		t, _ := time.Parse("2006-01-02", endDate)
		startDate = t.AddDate(0, 0, -30).Format("2006-01-02")
	}

	variables := map[string]interface{}{
		"deviceId":  deviceID,
		"startDate": startDate,
		"endDate":   endDate,
	}

	resp, err := client.Execute(api.SmartBatterySessionsQuery, "SmartBatterySessions", variables)
	if err != nil {
		return fmt.Errorf("failed to fetch battery sessions: %w", err)
	}

	var result models.SmartBatterySessionsResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	sessions := result.SmartBatterySessions
	if sessions == nil {
		fmt.Println("No session data available")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(sessions)
	}

	fmt.Printf("Battery Sessions (%s to %s)\n", startDate, endDate)
	fmt.Printf("Total Result: €%.2f\n\n", sessions.PeriodTotalResult)

	if len(sessions.Sessions) == 0 {
		fmt.Println("No sessions in this period")
		return nil
	}

	headers := []string{"Date", "Result", "Cumulative", "Status"}
	var rows [][]string

	for _, s := range sessions.Sessions {
		rows = append(rows, []string{
			formatDate(s.Date),
			fmt.Sprintf("€%.2f", s.Result),
			fmt.Sprintf("€%.2f", s.CumulativeResult),
			s.Status,
		})
	}

	output.Table(headers, rows)

	return nil
}

func getFirstBatteryID(client *api.Client) (string, error) {
	resp, err := client.Execute(api.SmartBatteriesQuery, "SmartBatteries", nil)
	if err != nil {
		return "", err
	}

	var result models.SmartBatteriesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", err
	}

	if len(result.SmartBatteries) > 0 {
		return result.SmartBatteries[0].ID, nil
	}

	return "", nil
}
