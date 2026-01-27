package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
	"github.com/pietern/frankie/internal/models"
	"github.com/pietern/frankie/internal/output"
)

var (
	usageSite  string
	usageDate  string
	usageType  string
)

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show energy usage and costs",
	Long:  `Display energy usage and costs for a specific period.`,
	RunE:  runUsage,
}

func init() {
	rootCmd.AddCommand(usageCmd)
	usageCmd.Flags().StringVarP(&usageSite, "site", "s", "", "site reference (optional if you have one site)")
	usageCmd.Flags().StringVarP(&usageDate, "date", "d", "", "date (YYYY-MM-DD, default: today)")
	usageCmd.Flags().StringVarP(&usageType, "type", "t", "", "type: electricity, gas, or feedin (default: all)")
}

func runUsage(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	// Resolve site reference
	siteRef, err := resolveSiteReference(client, usageSite)
	if err != nil {
		return err
	}

	// Set date
	date := usageDate
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	variables := map[string]interface{}{
		"siteReference": siteRef,
		"date":          date,
	}

	resp, err := client.Execute(api.PeriodUsageAndCostsQuery, "PeriodUsageAndCosts", variables)
	if err != nil {
		return fmt.Errorf("failed to fetch usage: %w", err)
	}

	var result models.PeriodUsageAndCostsResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	usage := result.PeriodUsageAndCosts
	if usage == nil {
		fmt.Println("No usage data available")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(usage)
	}

	// Display based on type filter
	switch usageType {
	case "electricity", "elec", "e":
		displayEnergyUsage("Electricity", usage.Electricity, date)
	case "gas", "g":
		displayEnergyUsage("Gas", usage.Gas, date)
	case "feedin", "feed", "f":
		displayEnergyUsage("Feed-in", usage.FeedIn, date)
	default:
		// Show summary of all
		displayUsageSummary(usage, date)
	}

	return nil
}

func displayUsageSummary(usage *models.PeriodUsageAndCosts, date string) {
	fmt.Printf("Usage summary for %s\n\n", date)

	headers := []string{"Type", "Usage", "Costs"}
	var rows [][]string

	if usage.Electricity != nil {
		rows = append(rows, []string{
			"Electricity",
			fmt.Sprintf("%.2f %s", usage.Electricity.UsageTotal, usage.Electricity.Unit),
			fmt.Sprintf("€%.2f", usage.Electricity.CostsTotal),
		})
	}

	if usage.Gas != nil {
		rows = append(rows, []string{
			"Gas",
			fmt.Sprintf("%.2f %s", usage.Gas.UsageTotal, usage.Gas.Unit),
			fmt.Sprintf("€%.2f", usage.Gas.CostsTotal),
		})
	}

	if usage.FeedIn != nil && usage.FeedIn.UsageTotal != 0 {
		rows = append(rows, []string{
			"Feed-in",
			fmt.Sprintf("%.2f %s", usage.FeedIn.UsageTotal, usage.FeedIn.Unit),
			fmt.Sprintf("€%.2f", usage.FeedIn.CostsTotal),
		})
	}

	if len(rows) == 0 {
		fmt.Println("No usage data available")
		return
	}

	table := output.Table(headers)
	table.AppendBulk(rows)
	table.Render()
}

func displayEnergyUsage(name string, category *models.EnergyCategory, date string) {
	if category == nil || len(category.Items) == 0 {
		fmt.Printf("No %s usage data available for %s\n", name, date)
		return
	}

	fmt.Printf("%s usage for %s\n", name, date)
	fmt.Printf("Total: %.2f %s (€%.2f)\n\n", category.UsageTotal, category.Unit, category.CostsTotal)

	headers := []string{"Time", "Usage", "Costs"}
	var rows [][]string

	for _, item := range category.Items {
		rows = append(rows, []string{
			formatTime(item.From),
			fmt.Sprintf("%.3f %s", item.Usage, item.Unit),
			fmt.Sprintf("€%.4f", item.Costs),
		})
	}

	table := output.Table(headers)
	table.AppendBulk(rows)
	table.Render()
}

func formatTime(isoTime string) string {
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		return isoTime
	}
	return t.Local().Format("15:04")
}
