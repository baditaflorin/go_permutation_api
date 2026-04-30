package datasource

import (
	"fmt"
)

// DataSource is an interface for loading elements from various sources
type DataSource interface {
	Load() ([]string, error)
	Close() error
}

// SourceType represents the type of data source
type SourceType string

const (
	SourceTypeCSV      SourceType = "csv"
	SourceTypeTSV      SourceType = "tsv"
	SourceTypeDatabase SourceType = "database"
)

// Config holds configuration for creating a data source
type Config struct {
	Type     SourceType
	FilePath string // for CSV/TSV
	Column   int    // column index for CSV/TSV (0-based)

	// Database configuration
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBTable    string
	DBColumn   string
	DBSSLMode  string
}

// New creates a new data source based on the configuration
func New(cfg Config) (DataSource, error) {
	switch cfg.Type {
	case SourceTypeCSV:
		return NewCSVSource(cfg.FilePath, cfg.Column)
	case SourceTypeTSV:
		return NewTSVSource(cfg.FilePath, cfg.Column)
	case SourceTypeDatabase:
		return NewDatabaseSource(DatabaseConfig{
			Driver:   cfg.DBDriver,
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			User:     cfg.DBUser,
			Password: cfg.DBPassword,
			Database: cfg.DBName,
			Table:    cfg.DBTable,
			Column:   cfg.DBColumn,
			SSLMode:  cfg.DBSSLMode,
		})
	default:
		return nil, fmt.Errorf("unsupported source type: %s", cfg.Type)
	}
}
