package datasource

import (
	"encoding/csv"
	"fmt"
	"os"
)

// TSVSource loads elements from a TSV (Tab-Separated Values) file
type TSVSource struct {
	filePath string
	column   int
	file     *os.File
}

// NewTSVSource creates a new TSV data source
func NewTSVSource(filePath string, column int) (*TSVSource, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	if column < 0 {
		return nil, fmt.Errorf("column index must be non-negative")
	}

	return &TSVSource{
		filePath: filePath,
		column:   column,
	}, nil
}

// Load reads elements from the TSV file
func (s *TSVSource) Load() ([]string, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open TSV file: %w", err)
	}
	s.file = file

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Use tab as separator
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read TSV file: %w", err)
	}

	if len(records) == 0 {
		return []string{}, nil
	}

	var elements []string
	for i, record := range records {
		if len(record) <= s.column {
			return nil, fmt.Errorf("row %d does not have column %d", i, s.column)
		}
		elements = append(elements, record[s.column])
	}

	return elements, nil
}

// Close closes the file handle
func (s *TSVSource) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
