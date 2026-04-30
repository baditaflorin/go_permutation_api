package datasource

import (
	"context"
	"database/sql"
	"time"
)

// HealthCheck checks if a database connection is healthy
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	// Try a simple query
	var result int
	err := db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return err
	}

	return nil
}

// GetDatabaseStats returns database connection pool stats
func GetDatabaseStats(db *sql.DB) map[string]interface{} {
	stats := db.Stats()
	return map[string]interface{}{
		"max_open_connections":   stats.MaxOpenConnections,
		"open_connections":       stats.OpenConnections,
		"in_use":                 stats.InUse,
		"idle":                   stats.Idle,
		"wait_count":             stats.WaitCount,
		"wait_duration":          stats.WaitDuration.String(),
		"max_idle_closed":        stats.MaxIdleClosed,
		"max_idle_time_closed":   stats.MaxIdleTimeClosed,
		"max_lifetime_closed":    stats.MaxLifetimeClosed,
	}
}
