package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	frankieErrors "github.com/pietern/frankie/internal/errors"
)

var (
	outputFormat string
	verbose      bool
)

var rootCmd = &cobra.Command{
	Use:   "frankie",
	Short: "CLI tool for Frank Energie",
	Long:  `Frankie is a command-line interface for interacting with the Frank Energie API.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, frankieErrors.Format(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format: table or json")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func getOutputFormat() string {
	return outputFormat
}

func isVerbose() bool {
	return verbose
}
