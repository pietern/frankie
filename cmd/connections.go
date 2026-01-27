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

var connectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "List energy connections",
	Long:  `Display all energy connections (electricity and gas meters) with EAN codes and contract details.`,
	RunE:  runConnections,
}

func init() {
	rootCmd.AddCommand(connectionsCmd)
}

func runConnections(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	manager := auth.NewManager(client)

	if err := manager.EnsureAuthenticated(); err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	resp, err := client.Execute(api.MeQuery, "Me", nil)
	if err != nil {
		return fmt.Errorf("failed to fetch connections: %w", err)
	}

	var result models.MeResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	connections := result.Me.Connections

	if len(connections) == 0 {
		fmt.Println("No connections found")
		return nil
	}

	if getOutputFormat() == "json" {
		return output.JSON(connections)
	}

	headers := []string{"Segment", "EAN", "Grid Operator", "Product", "Status", "Meter Type"}
	var rows [][]string

	for _, conn := range connections {
		gridOperator := ""
		productName := ""

		if conn.ExternalDetails != nil {
			gridOperator = conn.ExternalDetails.GridOperator
			if conn.ExternalDetails.Contract != nil {
				productName = conn.ExternalDetails.Contract.ProductName
			}
		}

		rows = append(rows, []string{
			conn.Segment,
			conn.EAN,
			gridOperator,
			productName,
			conn.Status,
			conn.MeterType,
		})
	}

	table := output.Table(headers)
	table.AppendBulk(rows)
	table.Render()

	return nil
}
