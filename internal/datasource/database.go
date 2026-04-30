package datasource

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"    // SQLite driver
)

var safeIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func secureIdentifier(name string) (string, error) {
	if !safeIdentifier.MatchString(name) {
		return "", fmt.Errorf("identifier %q contains unsafe characters", name)
	}
	return fmt.Sprintf(`"%s"`, name), nil
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Driver   string // postgres, mysql, sqlite3
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Table    string
	Column   string
	SSLMode  string
}

// DatabaseSource loads elements from a database
type DatabaseSource struct {
	config DatabaseConfig
	db     *sql.DB
}

// NewDatabaseSource creates a new database data source
func NewDatabaseSource(config DatabaseConfig) (*DatabaseSource, error) {
	if config.Driver == "" {
		return nil, fmt.Errorf("database driver cannot be empty")
	}
	if config.Table == "" {
		return nil, fmt.Errorf("table name cannot be empty")
	}
	if config.Column == "" {
		return nil, fmt.Errorf("column name cannot be empty")
	}

	return &DatabaseSource{
		config: config,
	}, nil
}

// Load reads elements from the database
func (s *DatabaseSource) Load() ([]string, error) {
	connStr, err := s.buildConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(s.config.Driver, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	s.db = db

	// Configure connection pooling
	ConfigureConnectionPool(db)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Use security package to safely quote identifiers (ADR-0003: SQL injection prevention)
	col, err := secureIdentifier(s.config.Column)
	if err != nil {
		return nil, fmt.Errorf("unsafe column name: %w", err)
	}
	tbl, err := secureIdentifier(s.config.Table)
	if err != nil {
		return nil, fmt.Errorf("unsafe table name: %w", err)
	}

	query := fmt.Sprintf("SELECT %s FROM %s", col, tbl)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var elements []string
	for rows.Next() {
		var element string
		if err := rows.Scan(&element); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		elements = append(elements, element)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return elements, nil
}

// buildConnectionString builds a connection string based on the driver
func (s *DatabaseSource) buildConnectionString() (string, error) {
	switch s.config.Driver {
	case "postgres":
		sslMode := s.config.SSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			s.config.Host, s.config.Port, s.config.User, s.config.Password, s.config.Database, sslMode), nil

	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			s.config.User, s.config.Password, s.config.Host, s.config.Port, s.config.Database), nil

	case "sqlite3":
		return s.config.Database, nil

	default:
		return "", fmt.Errorf("unsupported database driver: %s", s.config.Driver)
	}
}

// Close closes the database connection
func (s *DatabaseSource) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
