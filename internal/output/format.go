package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// Format represents the output format type
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Table prints a table with default styling to stdout
func Table(headers []string, rows [][]string) {
	TableTo(os.Stdout, headers, rows)
}

// TableTo prints a table with default styling to a specific writer
func TableTo(w io.Writer, headers []string, rows [][]string) {
	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().PaddingRight(2)
			if row == table.HeaderRow {
				return style.Bold(true)
			}
			return style
		})
	fmt.Fprintln(w, t)
}

// JSON outputs data as formatted JSON
func JSON(data interface{}) error {
	return JSONTo(os.Stdout, data)
}

// JSONTo outputs data as formatted JSON to a specific writer
func JSONTo(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// Print outputs data in the specified format
func Print(format Format, headers []string, rows [][]string, jsonData interface{}) error {
	switch format {
	case FormatJSON:
		return JSON(jsonData)
	default:
		Table(headers, rows)
		return nil
	}
}

// KeyValue prints key-value pairs in a simple format
func KeyValue(pairs map[string]string) {
	KeyValueTo(os.Stdout, pairs)
}

// KeyValueTo prints key-value pairs to a specific writer
func KeyValueTo(w io.Writer, pairs map[string]string) {
	maxKeyLen := 0
	for k := range pairs {
		if len(k)+1 > maxKeyLen { // +1 for colon
			maxKeyLen = len(k) + 1
		}
	}

	for k, v := range pairs {
		fmt.Fprintf(w, "%-*s  %s\n", maxKeyLen, k+":", v)
	}
}

// KeyValueOrdered prints key-value pairs in order
func KeyValueOrdered(keys []string, pairs map[string]string) {
	KeyValueOrderedTo(os.Stdout, keys, pairs)
}

// KeyValueOrderedTo prints key-value pairs in order to a specific writer
func KeyValueOrderedTo(w io.Writer, keys []string, pairs map[string]string) {
	maxKeyLen := 0
	for _, k := range keys {
		if len(k)+1 > maxKeyLen { // +1 for colon
			maxKeyLen = len(k) + 1
		}
	}

	for _, k := range keys {
		if v, ok := pairs[k]; ok {
			fmt.Fprintf(w, "%-*s  %s\n", maxKeyLen, k+":", v)
		}
	}
}
