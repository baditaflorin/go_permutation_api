package datasource

import (
	"encoding/csv"
	"fmt"
	"os"
)

// CSVSource loads elements from a CSV file
type CSVSource struct {
	filePath string
	column   int
	file     *os.File
}

// NewCSVSource creates a new CSV data source
func NewCSVSource(filePath string, column int) (*CSVSource, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	if column < 0 {
		return nil, fmt.Errorf("column index must be non-negative")
	}

	return &CSVSource{
		filePath: filePath,
		column:   column,
	}, nil
}

// Load reads elements from the CSV file
func (s *CSVSource) Load() ([]string, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	s.file = file

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
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
func (s *CSVSource) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
