package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// OutputFormat represents the output format for CLI
type OutputFormat string

const (
	FormatPlain OutputFormat = "plain"
	FormatJSON  OutputFormat = "json"
	FormatCSV   OutputFormat = "csv"
)

// OutputWriter handles writing permutations in different formats
type OutputWriter struct {
	format    OutputFormat
	csvWriter *csv.Writer
	perms     [][]string
}

// NewOutputWriter creates a new output writer
func NewOutputWriter(format OutputFormat) *OutputWriter {
	w := &OutputWriter{
		format: format,
		perms:  [][]string{},
	}

	if format == FormatCSV {
		w.csvWriter = csv.NewWriter(os.Stdout)
	}

	return w
}

// Write writes a permutation
func (w *OutputWriter) Write(perm []string) error {
	switch w.format {
	case FormatPlain:
		fmt.Println(strings.Join(perm, ","))
	case FormatJSON:
		// Collect for batch output
		w.perms = append(w.perms, perm)
	case FormatCSV:
		if err := w.csvWriter.Write(perm); err != nil {
			return err
		}
	}
	return nil
}

// Flush flushes any buffered output
func (w *OutputWriter) Flush() error {
	switch w.format {
	case FormatJSON:
		data, err := json.MarshalIndent(w.perms, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case FormatCSV:
		w.csvWriter.Flush()
		return w.csvWriter.Error()
	}
	return nil
}
