package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
	"github.com/pietern/frankie/internal/models"
	"github.com/pietern/frankie/internal/output"
)

var (
	pricesDate     string
	pricesBelgium  bool
	pricesSite     string
	pricesShowGas  bool
)

var pricesCmd = &cobra.Command{
	Use:   "prices",
	Short: "Show energy prices",
	Long:  `Display current electricity and gas market prices.`,
	RunE:  runPrices,
}

func init() {
	rootCmd.AddCommand(pricesCmd)
	pricesCmd.Flags().StringVarP(&pricesDate, "date", "d", "", "date to show prices for (YYYY-MM-DD, default: today)")
	pricesCmd.Flags().BoolVar(&pricesBelgium, "be", false, "show Belgium prices instead of Netherlands")
	pricesCmd.Flags().StringVarP(&pricesSite, "site", "s", "", "site reference for customer-specific prices")
	pricesCmd.Flags().BoolVar(&pricesShowGas, "gas", false, "show gas prices instead of electricity")
}

func runPrices(cmd *cobra.Command, args []string) error {
	client := api.NewClient()

	// Set date
	date := pricesDate
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var prices *models.MarketPrices
	var err error

	if pricesSite != "" {
		// Customer-specific prices (requires auth)
		// Resolve partial site reference
		var siteRef string
		siteRef, err = resolveSiteReference(client, pricesSite)
		if err != nil {
			return err
		}
		prices, err = fetchCustomerPrices(client, date, siteRef)
	} else if pricesBelgium {
		// Belgium prices
		client.SetCountry("BE")
		prices, err = fetchBelgiumPrices(client, date)
	} else {
		// Netherlands public prices
		prices, err = fetchPublicPrices(client, date)
	}

	if err != nil {
		return err
	}

	if prices == nil {
		return fmt.Errorf("no prices available for %s", date)
	}

	if getOutputFormat() == "json" {
		return output.JSON(prices)
	}

	// Display prices
	if pricesShowGas {
		return displayGasPrices(prices, date)
	}
	return displayElectricityPrices(prices, date)
}

func fetchPublicPrices(client *api.Client, date string) (*models.MarketPrices, error) {
	variables := map[string]interface{}{
		"date":       date,
		"resolution": "PT60M",
	}

	resp, err := client.Execute(api.MarketPricesQuery, "MarketPrices", variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}

	var result models.MarketPricesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.MarketPrices, nil
}

func fetchBelgiumPrices(client *api.Client, date string) (*models.MarketPrices, error) {
	variables := map[string]interface{}{
		"date": date,
	}

	resp, err := client.Execute(api.BelgiumMarketPricesQuery, "MarketPrices", variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Belgium prices: %w", err)
	}

	var result models.MarketPricesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.MarketPrices, nil
}

func fetchCustomerPrices(client *api.Client, date, siteRef string) (*models.MarketPrices, error) {
	manager := auth.NewManager(client)
	if err := manager.EnsureAuthenticated(); err != nil {
		return nil, fmt.Errorf("not logged in: %w", err)
	}

	variables := map[string]interface{}{
		"date":          date,
		"siteReference": siteRef,
	}

	resp, err := client.Execute(api.CustomerMarketPricesQuery, "MarketPrices", variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customer prices: %w", err)
	}

	var result models.CustomerMarketPricesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.CustomerMarketPrices, nil
}

func displayElectricityPrices(prices *models.MarketPrices, date string) error {
	if len(prices.ElectricityPrices) == 0 {
		fmt.Printf("No electricity prices available for %s\n", date)
		return nil
	}

	fmt.Printf("Electricity prices for %s\n\n", date)

	// Show average if available
	if prices.AverageElectricityPrices != nil {
		avg := prices.AverageElectricityPrices
		fmt.Printf("Average: €%.4f/%s (all-in: €%.4f)\n\n",
			avg.AverageMarketPrice, avg.PerUnit, avg.AverageAllInPrice)
	}

	headers := []string{"Time", "Market", "Total", "All-In"}
	var rows [][]string

	for _, p := range prices.ElectricityPrices {
		timeStr := p.From.Local().Format("15:04")
		allIn := p.AllInPrice
		if allIn == 0 {
			allIn = p.TotalPrice() // Fallback for Belgium prices
		}
		rows = append(rows, []string{
			timeStr,
			fmt.Sprintf("€%.4f", p.MarketPrice),
			fmt.Sprintf("€%.4f", p.TotalPrice()),
			fmt.Sprintf("€%.4f", allIn),
		})
	}

	output.Table(headers, rows)

	return nil
}

func displayGasPrices(prices *models.MarketPrices, date string) error {
	if len(prices.GasPrices) == 0 {
		fmt.Printf("No gas prices available for %s\n", date)
		return nil
	}

	fmt.Printf("Gas prices for %s\n\n", date)

	headers := []string{"Time", "Market", "Total", "All-In"}
	var rows [][]string

	for _, p := range prices.GasPrices {
		timeStr := p.From.Local().Format("15:04")
		allIn := p.AllInPrice
		if allIn == 0 {
			allIn = p.TotalPrice() // Fallback for Belgium prices
		}
		rows = append(rows, []string{
			timeStr,
			fmt.Sprintf("€%.4f", p.MarketPrice),
			fmt.Sprintf("€%.4f", p.TotalPrice()),
			fmt.Sprintf("€%.4f", allIn),
		})
	}

	output.Table(headers, rows)

	return nil
}

// resolveSiteReference resolves a partial site reference to the full reference
// It supports matching by:
// - Full reference (exact match)
// - Postal code (e.g., "8147RJ")
// - Postal code with house number (e.g., "8147RJ 26")
func resolveSiteReference(client *api.Client, partial string) (string, error) {
	manager := auth.NewManager(client)
	if err := manager.EnsureAuthenticated(); err != nil {
		return "", fmt.Errorf("not logged in: %w", err)
	}

	// Fetch sites
	resp, err := client.Execute(api.UserSitesQuery, "UserSites", nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch sites: %w", err)
	}

	var result models.UserSitesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", fmt.Errorf("failed to parse sites: %w", err)
	}

	sites := result.UserSites
	if len(sites) == 0 {
		return "", fmt.Errorf("no sites found")
	}

	// Normalize the partial reference for matching
	partial = strings.ToUpper(strings.TrimSpace(partial))

	// Try to find a matching site
	var matches []models.Site
	for _, site := range sites {
		ref := strings.ToUpper(site.Reference)
		if ref == partial {
			// Exact match
			return site.Reference, nil
		}
		if strings.HasPrefix(ref, partial) {
			matches = append(matches, site)
		}
	}

	if len(matches) == 1 {
		return matches[0].Reference, nil
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous site reference '%s' matches %d sites", partial, len(matches))
	}

	// If only one site, use it regardless of the input
	if len(sites) == 1 {
		return sites[0].Reference, nil
	}

	return "", fmt.Errorf("no site found matching '%s'", partial)
}
