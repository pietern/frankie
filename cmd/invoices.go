package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
	"github.com/pietern/frankie/internal/models"
	"github.com/pietern/frankie/internal/output"
)

var (
	invoicesSite string
	invoicesAll  bool
)

var invoicesCmd = &cobra.Command{
	Use:   "invoices",
	Short: "Show invoices",
	Long:  `Display invoice information including current, previous, and upcoming periods.`,
	RunE:  runInvoices,
}

func init() {
	rootCmd.AddCommand(invoicesCmd)
	invoicesCmd.Flags().StringVarP(&invoicesSite, "site", "s", "", "site reference (optional if you have one site)")
	invoicesCmd.Flags().BoolVarP(&invoicesAll, "all", "a", false, "show all invoices")
}

func runInvoices(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	// Resolve site reference
	siteRef, err := resolveSiteReference(client, invoicesSite)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"siteReference": siteRef,
	}

	resp, err := client.Execute(api.InvoicesQuery, "Invoices", variables)
	if err != nil {
		return fmt.Errorf("failed to fetch invoices: %w", err)
	}

	var result models.InvoicesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	invoices := result.Invoices
	if invoices == nil {
		fmt.Println("No invoice data available")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(invoices)
	}

	if invoicesAll && len(invoices.AllInvoices) > 0 {
		displayAllInvoices(invoices.AllInvoices)
	} else {
		displayInvoiceSummary(invoices)
	}

	return nil
}

func displayInvoiceSummary(invoices *models.Invoices) {
	fmt.Println("Invoice Summary")
	fmt.Println()

	if invoices.PreviousPeriodInvoice != nil {
		inv := invoices.PreviousPeriodInvoice
		fmt.Printf("Previous: %s - €%.2f\n", inv.PeriodDescription, inv.TotalAmount)
	}

	if invoices.CurrentPeriodInvoice != nil {
		inv := invoices.CurrentPeriodInvoice
		fmt.Printf("Current:  %s - €%.2f\n", inv.PeriodDescription, inv.TotalAmount)
	}

	if invoices.UpcomingPeriodInvoice != nil {
		inv := invoices.UpcomingPeriodInvoice
		fmt.Printf("Upcoming: %s - €%.2f\n", inv.PeriodDescription, inv.TotalAmount)
	}

	if len(invoices.AllInvoices) > 0 {
		fmt.Printf("\nTotal invoices: %d (use --all to see all)\n", len(invoices.AllInvoices))
	}
}

func displayAllInvoices(invoices []models.Invoice) {
	headers := []string{"Period", "Date", "Amount"}
	var rows [][]string

	for _, inv := range invoices {
		rows = append(rows, []string{
			inv.PeriodDescription,
			formatDate(inv.StartDate),
			fmt.Sprintf("€%.2f", inv.TotalAmount),
		})
	}

	table := output.Table(headers)
	table.AppendBulk(rows)
	table.Render()
}
