package database

import (
	"context"

	// "github.com/e173-gateway/e173_go_gateway/pkg/config" // No longer needed directly if URL passed
	"github.com/e173-gateway/e173_go_gateway/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDBPool initializes and returns a new database connection pool.
func NewDBPool(databaseURL string) (*pgxpool.Pool, error) { // Renamed and changed signature
	logging.Logger.Info("Attempting to connect to database...")

	parsedConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		logging.Logger.Errorf("Unable to parse database URL: %v", err)
		return nil, err
	}

	// You can configure pool settings here if needed, e.g.:
	// parsedConfig.MaxConns = 10
	// parsedConfig.MinConns = 2
	// parsedConfig.MaxConnLifetime = time.Hour
	// parsedConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), parsedConfig)
	if err != nil {
		logging.Logger.Errorf("Unable to create connection pool: %v", err)
		return nil, err
	}

	// Ping the database to verify connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close() // Close the pool if ping fails
		logging.Logger.Errorf("Unable to connect to database (ping failed): %v", err)
		return nil, err
	}
	// logging.Logger.Info("Successfully connected to the database.") // This log can be in main
	return pool, nil
}

// CloseDB is no longer strictly necessary here if main handles its pool closure.
// func CloseDB(pool *pgxpool.Pool) {
// 	if pool != nil {
// 		pool.Close()
// 		logging.Logger.Info("Database connection pool closed.")
// 	}
// }
