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

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "List user sites",
	Long:  `Display all sites (delivery addresses) linked to your account.`,
	RunE:  runSites,
}

func init() {
	rootCmd.AddCommand(sitesCmd)
}

func runSites(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	resp, err := client.Execute(api.UserSitesQuery, "UserSites", nil)
	if err != nil {
		return fmt.Errorf("failed to fetch sites: %w", err)
	}

	var result models.UserSitesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	sites := result.UserSites

	if len(sites) == 0 {
		fmt.Println("No sites found")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(sites)
	}

	headers := []string{"Reference", "Address", "Status", "Segments", "Start", "Last Reading"}
	var rows [][]string

	for _, site := range sites {
		address := ""
		if site.Address != nil {
			address = site.Address.FormattedAddress()
		}

		segments := strings.Join(site.Segments, ", ")

		startDate := ""
		if site.DeliveryStartDate != "" {
			startDate = formatDate(site.DeliveryStartDate)
		}

		lastReading := ""
		if site.LastMeterReadingDate != "" {
			lastReading = formatDate(site.LastMeterReadingDate)
		}

		rows = append(rows, []string{
			site.Reference,
			address,
			site.Status,
			segments,
			startDate,
			lastReading,
		})
	}

	table := output.Table(headers)
	table.AppendBulk(rows)
	table.Render()

	return nil
}

// formatDate formats an ISO date string to a shorter format
func formatDate(isoDate string) string {
	if len(isoDate) >= 10 {
		return isoDate[:10]
	}
	return isoDate
}
