package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
)

var (
	operationName string
	variablesJSON string
	debugRequest  bool
)

var apiCmd = &cobra.Command{
	Use:   "api <query>",
	Short: "Execute raw GraphQL query",
	Long: `Execute a raw GraphQL query against the Frank Energie API.

Handles authentication and required headers automatically.

Examples:
  frankie api 'query Me { me { email } }'
  frankie api 'query Version { version }'

  # For queries with variables, use heredoc to avoid shell escaping issues:
  cat << 'QUERY' | frankie api --var '{"date":"2025-01-28","resolution":"HOUR"}' -
  query MarketPrices($date: String!, $resolution: PriceResolution!) {
    marketPrices(date: $date, resolution: $resolution) {
      electricityPrices { from till allInPrice }
    }
  }
  QUERY`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAPI,
}

func init() {
	apiCmd.Flags().StringVar(&operationName, "op", "", "operation name (auto-detected if not specified)")
	apiCmd.Flags().StringVar(&variablesJSON, "var", "", "variables as JSON object")
	apiCmd.Flags().BoolVar(&debugRequest, "debug", false, "print request JSON before sending")
	rootCmd.AddCommand(apiCmd)
}

func runAPI(cmd *cobra.Command, args []string) error {
	var query string

	if len(args) == 0 || args[0] == "-" {
		// Read from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
		query = strings.TrimSpace(string(data))
	} else {
		query = args[0]
	}

	if query == "" {
		return fmt.Errorf("query is required")
	}

	// Auto-detect operation name if not specified
	opName := operationName
	if opName == "" {
		// Try to extract from query: "query Name" or "mutation Name"
		re := regexp.MustCompile(`(?:query|mutation)\s+(\w+)`)
		if matches := re.FindStringSubmatch(query); len(matches) > 1 {
			opName = matches[1]
		}
	}

	// Parse variables if provided
	var variables map[string]interface{}
	if variablesJSON != "" {
		if err := json.Unmarshal([]byte(variablesJSON), &variables); err != nil {
			return fmt.Errorf("invalid variables JSON: %w", err)
		}
	}

	if debugRequest {
		reqDebug := map[string]interface{}{
			"query":         query,
			"operationName": opName,
			"variables":     variables,
		}
		debugJSON, _ := json.MarshalIndent(reqDebug, "", "  ")
		fmt.Fprintln(os.Stderr, "Request:")
		fmt.Fprintln(os.Stderr, string(debugJSON))
		fmt.Fprintln(os.Stderr)
	}

	client := api.NewClient()

	// Try to authenticate if we have stored credentials
	manager := auth.NewManager(client)
	_ = manager.EnsureAuthenticated() // Ignore error, proceed without auth if not logged in

	body, err := client.ExecuteRaw(query, opName, variables)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Pretty print the JSON
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// Fall back to raw output
		fmt.Println(string(body))
		return nil
	}

	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(string(body))
		return nil
	}

	fmt.Println(string(output))
	return nil
}
