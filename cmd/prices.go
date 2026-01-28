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

const (
	resolution15Min = 15
	resolution60Min = 60
	resolutionPT15M = "PT15M"
	resolutionPT60M = "PT60M"

	// Day-ahead prices are published around 12:55 CET
	tomorrowPricesAvailableHour = 13
)

var (
	pricesDate       string
	pricesBelgium    bool
	pricesSite       string
	pricesShowGas    bool
	pricesResolution int
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
	pricesCmd.Flags().IntVarP(&pricesResolution, "resolution", "r", resolution60Min, "price resolution in minutes (15 or 60, 15 requires login)")
}

func runPrices(cmd *cobra.Command, args []string) error {
	client := api.NewClient()

	// Determine dates to fetch
	dates := getPriceDates()

	// Collect all prices
	type priceResult struct {
		date   string
		prices *models.MarketPrices
	}
	var results []priceResult

	for _, date := range dates {
		var prices *models.MarketPrices
		var err error

		if pricesSite != "" {
			// Customer-specific prices (requires auth)
			// Resolve partial site reference (only once)
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
			var resolution string
			switch pricesResolution {
			case resolution15Min:
				resolution = resolutionPT15M
				// 15-minute resolution requires authentication
				manager := auth.NewManager(client)
				if err := manager.EnsureAuthenticated(); err != nil {
					return fmt.Errorf("15-minute resolution requires login: %w", err)
				}
			case resolution60Min:
				resolution = resolutionPT60M
			default:
				return fmt.Errorf("invalid resolution: %d (must be %d or %d)", pricesResolution, resolution15Min, resolution60Min)
			}
			prices, err = fetchPublicPrices(client, date, resolution)
		}

		if err != nil {
			return err
		}

		if prices != nil {
			results = append(results, priceResult{date: date, prices: prices})
		}
	}

	if len(results) == 0 {
		return fmt.Errorf("no prices available")
	}

	if getOutputFormat() == "json" {
		if len(results) == 1 {
			return output.JSON(results[0].prices)
		}
		// Multiple dates: merge prices into single response
		merged := &models.MarketPrices{}
		for _, r := range results {
			merged.ElectricityPrices = append(merged.ElectricityPrices, r.prices.ElectricityPrices...)
			merged.GasPrices = append(merged.GasPrices, r.prices.GasPrices...)
		}
		return output.JSON(merged)
	}

	// Merge all prices
	var allPrices models.MarketPrices
	for _, r := range results {
		allPrices.ElectricityPrices = append(allPrices.ElectricityPrices, r.prices.ElectricityPrices...)
		allPrices.GasPrices = append(allPrices.GasPrices, r.prices.GasPrices...)
	}

	// Display prices
	if pricesShowGas {
		return displayGasPrices(&allPrices)
	}
	return displayElectricityPrices(&allPrices)
}

// getPriceDates returns the dates to fetch prices for.
// If a specific date was requested, returns only that date.
// Otherwise returns today, and tomorrow if after 13:00 CET.
func getPriceDates() []string {
	if pricesDate != "" {
		return []string{pricesDate}
	}

	// Load CET timezone
	cet, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		// Fallback to local time if timezone fails
		cet = time.Local
	}

	now := time.Now().In(cet)
	today := now.Format("2006-01-02")
	tomorrow := now.AddDate(0, 0, 1).Format("2006-01-02")

	// Day-ahead prices are published around 13:00 CET
	if now.Hour() >= tomorrowPricesAvailableHour {
		return []string{today, tomorrow}
	}
	return []string{today}
}

func fetchPublicPrices(client *api.Client, date, resolution string) (*models.MarketPrices, error) {
	variables := map[string]interface{}{
		"date":       date,
		"resolution": resolution,
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

func displayElectricityPrices(prices *models.MarketPrices) error {
	if len(prices.ElectricityPrices) == 0 {
		fmt.Println("No electricity prices available")
		return nil
	}

	fmt.Println("Electricity prices\n")

	headers := []string{"Date", "Time", "Market", "Total", "All-In"}
	var rows [][]string

	for _, p := range prices.ElectricityPrices {
		dateStr := p.From.Local().Format("2006-01-02")
		timeStr := p.From.Local().Format("15:04")
		allIn := p.AllInPrice
		if allIn == 0 {
			allIn = p.TotalPrice() // Fallback for Belgium prices
		}
		rows = append(rows, []string{
			dateStr,
			timeStr,
			fmt.Sprintf("€%.4f", p.MarketPrice),
			fmt.Sprintf("€%.4f", p.TotalPrice()),
			fmt.Sprintf("€%.4f", allIn),
		})
	}

	output.Table(headers, rows)

	return nil
}

func displayGasPrices(prices *models.MarketPrices) error {
	if len(prices.GasPrices) == 0 {
		fmt.Println("No gas prices available")
		return nil
	}

	fmt.Println("Gas prices\n")

	headers := []string{"Date", "Time", "Market", "Total", "All-In"}
	var rows [][]string

	for _, p := range prices.GasPrices {
		dateStr := p.From.Local().Format("2006-01-02")
		timeStr := p.From.Local().Format("15:04")
		allIn := p.AllInPrice
		if allIn == 0 {
			allIn = p.TotalPrice() // Fallback for Belgium prices
		}
		rows = append(rows, []string{
			dateStr,
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
