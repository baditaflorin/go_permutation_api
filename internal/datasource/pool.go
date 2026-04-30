package datasource

import (
	"database/sql"
	"time"
)

// ConfigureConnectionPool sets up database connection pooling
func ConfigureConnectionPool(db *sql.DB) {
	// Maximum number of open connections
	db.SetMaxOpenConns(25)

	// Maximum number of idle connections
	db.SetMaxIdleConns(5)

	// Maximum lifetime of a connection
	db.SetConnMaxLifetime(5 * time.Minute)

	// Maximum idle time for a connection
	db.SetConnMaxIdleTime(1 * time.Minute)
}
