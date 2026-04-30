package datasource

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCSVSource(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		column   int
		wantErr  bool
	}{
		{
			name:     "valid config",
			filePath: "test.csv",
			column:   0,
			wantErr:  false,
		},
		{
			name:     "empty file path",
			filePath: "",
			column:   0,
			wantErr:  true,
		},
		{
			name:     "negative column",
			filePath: "test.csv",
			column:   -1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCSVSource(tt.filePath, tt.column)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCSVSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCSVSourceLoad(t *testing.T) {
	// Create a temporary CSV file
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	content := "name,age,city\nAlice,30,NYC\nBob,25,LA\nCharlie,35,SF"
	if err := os.WriteFile(csvFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test CSV file: %v", err)
	}

	tests := []struct {
		name     string
		column   int
		expected []string
		wantErr  bool
	}{
		{
			name:     "load column 0",
			column:   0,
			expected: []string{"name", "Alice", "Bob", "Charlie"},
			wantErr:  false,
		},
		{
			name:     "load column 2",
			column:   2,
			expected: []string{"city", "NYC", "LA", "SF"},
			wantErr:  false,
		},
		{
			name:     "column out of range",
			column:   5,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, err := NewCSVSource(csvFile, tt.column)
			if err != nil {
				t.Fatalf("failed to create CSV source: %v", err)
			}
			defer source.Close()

			elements, err := source.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(elements) != len(tt.expected) {
					t.Errorf("expected %d elements, got %d", len(tt.expected), len(elements))
					return
				}
				for i, elem := range elements {
					if elem != tt.expected[i] {
						t.Errorf("at index %d: expected %q, got %q", i, tt.expected[i], elem)
					}
				}
			}
		})
	}
}

func TestTSVSourceLoad(t *testing.T) {
	// Create a temporary TSV file
	tmpDir := t.TempDir()
	tsvFile := filepath.Join(tmpDir, "test.tsv")

	content := "name\tage\tcity\nAlice\t30\tNYC\nBob\t25\tLA"
	if err := os.WriteFile(tsvFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test TSV file: %v", err)
	}

	source, err := NewTSVSource(tsvFile, 0)
	if err != nil {
		t.Fatalf("failed to create TSV source: %v", err)
	}
	defer source.Close()

	elements, err := source.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expected := []string{"name", "Alice", "Bob"}
	if len(elements) != len(expected) {
		t.Errorf("expected %d elements, got %d", len(expected), len(elements))
	}
}

func TestNewDatabaseSource(t *testing.T) {
	tests := []struct {
		name    string
		config  DatabaseConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: DatabaseConfig{
				Driver: "postgres",
				Table:  "users",
				Column: "name",
			},
			wantErr: false,
		},
		{
			name: "empty driver",
			config: DatabaseConfig{
				Table:  "users",
				Column: "name",
			},
			wantErr: true,
		},
		{
			name: "empty table",
			config: DatabaseConfig{
				Driver: "postgres",
				Column: "name",
			},
			wantErr: true,
		},
		{
			name: "empty column",
			config: DatabaseConfig{
				Driver: "postgres",
				Table:  "users",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDatabaseSource(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatabaseSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
