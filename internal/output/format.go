package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Format represents the output format type
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Table creates a new table writer with default styling
func Table(headers []string) *tablewriter.Table {
	return TableTo(os.Stdout, headers)
}

// TableTo creates a new table writer to a specific writer
func TableTo(w io.Writer, headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetHeader(headers)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)
	return table
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
		table := Table(headers)
		table.AppendBulk(rows)
		table.Render()
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
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
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
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	for _, k := range keys {
		if v, ok := pairs[k]; ok {
			fmt.Fprintf(w, "%-*s  %s\n", maxKeyLen, k+":", v)
		}
	}
}
