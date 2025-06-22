package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// CreateSQLXAdapter creates a *sqlx.DB from a *pgxpool.Pool for compatibility
// with enterprise repositories that expect sqlx interface
func CreateSQLXAdapter(pool *pgxpool.Pool) (*sqlx.DB, error) {
	// Create a sql.DB using pgx/stdlib driver from the pool
	stdDB := stdlib.OpenDBFromPool(pool)
	
	// Wrap with sqlx
	sqlxDB := sqlx.NewDb(stdDB, "pgx")
	
	// Test the connection
	if err := sqlxDB.Ping(); err != nil {
		return nil, err
	}
	
	return sqlxDB, nil
}

// CloseAdapter properly closes the sqlx adapter
func CloseAdapter(sqlxDB *sqlx.DB) error {
	return sqlxDB.Close()
}
