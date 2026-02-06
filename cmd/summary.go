package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/models"
	"github.com/pietern/frankie/internal/output"
)

var summarySite string

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show month summary",
	Long:  `Display the current month's cost summary and meter reading status.`,
	RunE:  runSummary,
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().StringVarP(&summarySite, "site", "s", "", "site reference (optional if you have one site)")
}

func runSummary(cmd *cobra.Command, args []string) error {
	client, err := newAuthenticatedClient()
	if err != nil {
		return err
	}

	// Resolve site reference
	siteRef, err := resolveSiteReference(client, summarySite)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"siteReference": siteRef,
	}

	resp, err := client.Execute(api.MonthSummaryQuery, "MonthSummary", variables)
	if err != nil {
		return fmt.Errorf("failed to fetch summary: %w", err)
	}

	var result models.MonthSummaryResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	summary := result.MonthSummary
	if summary == nil {
		fmt.Println("No summary data available")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(summary)
	}

	// Display summary
	keys := []string{
		"Actual Costs",
		"Expected Costs (to date)",
		"Expected Costs (month)",
		"Last Meter Reading",
		"Completeness",
		"Gas Included",
	}

	completeness := fmt.Sprintf("%.0f%%", summary.MeterReadingDayCompleteness*100)
	gasIncluded := "Yes"
	if summary.GasExcluded {
		gasIncluded = "No"
	}

	pairs := map[string]string{
		"Actual Costs":             fmt.Sprintf("€%.2f", summary.ActualCostsUntilLastMeterReadingDate),
		"Expected Costs (to date)": fmt.Sprintf("€%.2f", summary.ExpectedCostsUntilLastMeterReadingDate),
		"Expected Costs (month)":   fmt.Sprintf("€%.2f", summary.ExpectedCosts),
		"Last Meter Reading":       formatDate(summary.LastMeterReadingDate),
		"Completeness":             completeness,
		"Gas Included":             gasIncluded,
	}

	output.KeyValueOrdered(keys, pairs)
	return nil
}
