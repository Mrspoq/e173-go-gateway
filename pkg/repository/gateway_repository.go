package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresGatewayRepository implements the GatewayRepository interface for PostgreSQL.
type PostgresGatewayRepository struct {
	db *pgxpool.Pool
}

// NewPostgresGatewayRepository creates a new instance of PostgresGatewayRepository.
func NewPostgresGatewayRepository(db *pgxpool.Pool) *PostgresGatewayRepository {
	return &PostgresGatewayRepository{db: db}
}

// CreateGateway creates a new gateway record in the database.
func (r *PostgresGatewayRepository) CreateGateway(ctx context.Context, gateway *models.Gateway) error {
	query := `
		INSERT INTO gateways (
			id, name, description, location, ami_host, ami_port,
			ami_user, ami_pass, status, enabled, last_seen, 
			last_error, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)`

	if gateway.ID == "" {
		gateway.ID = uuid.New().String()
	}
	if gateway.CreatedAt.IsZero() {
		gateway.CreatedAt = time.Now()
	}
	if gateway.UpdatedAt.IsZero() {
		gateway.UpdatedAt = time.Now()
	}
	if gateway.Status == "" {
		gateway.Status = models.GatewayStatusOffline
	}
	if gateway.AMIPort == "" {
		gateway.AMIPort = "5038"
	}

	_, err := r.db.Exec(ctx, query,
		gateway.ID, gateway.Name, gateway.Description, gateway.Location, 
		gateway.AMIHost, gateway.AMIPort, gateway.AMIUser, gateway.AMIPass,
		gateway.Status, gateway.Enabled, gateway.LastSeen, gateway.LastError,
		gateway.CreatedAt, gateway.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("PostgresGatewayRepository.CreateGateway: failed to create gateway: %w", err)
	}
	return nil
}

// GetGatewayByID retrieves a gateway by its ID.
func (r *PostgresGatewayRepository) GetGatewayByID(ctx context.Context, id string) (*models.Gateway, error) {
	query := `
		SELECT 
			id, name, description, location, ami_host, ami_port,
			ami_user, ami_pass, status, enabled, last_seen, 
			last_error, created_at, updated_at
		FROM gateways
		WHERE id = $1`

	gateway := &models.Gateway{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&gateway.ID, &gateway.Name, &gateway.Description, &gateway.Location, &gateway.AMIHost, &gateway.AMIPort,
		&gateway.AMIUser, &gateway.AMIPass, &gateway.Status, &gateway.Enabled, &gateway.LastSeen, &gateway.LastError,
		&gateway.CreatedAt, &gateway.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("PostgresGatewayRepository.GetGatewayByID: failed to scan gateway: %w", err)
	}
	return gateway, nil
}

// ListGateways retrieves all gateways from the database.
func (r *PostgresGatewayRepository) ListGateways(ctx context.Context) ([]*models.Gateway, error) {
	query := `
		SELECT 
			id, name, description, location, ami_host, ami_port,
			ami_user, ami_pass, status, enabled, last_seen, 
			last_error, created_at, updated_at
		FROM gateways
		ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PostgresGatewayRepository.ListGateways: failed to query gateways: %w", err)
	}
	defer rows.Close()

	var gateways []*models.Gateway
	for rows.Next() {
		gateway := &models.Gateway{}
		err := rows.Scan(
			&gateway.ID, &gateway.Name, &gateway.Description, &gateway.Location, &gateway.AMIHost, &gateway.AMIPort,
			&gateway.AMIUser, &gateway.AMIPass, &gateway.Status, &gateway.Enabled, &gateway.LastSeen, &gateway.LastError,
			&gateway.CreatedAt, &gateway.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("PostgresGatewayRepository.ListGateways: failed to scan gateway: %w", err)
		}
		gateways = append(gateways, gateway)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PostgresGatewayRepository.ListGateways: row iteration error: %w", err)
	}

	return gateways, nil
}

// UpdateGateway updates an existing gateway in the database.
func (r *PostgresGatewayRepository) UpdateGateway(ctx context.Context, gateway *models.Gateway) error {
	query := `
		UPDATE gateways
		SET name = $2, description = $3, location = $4, ami_host = $5, ami_port = $6,
			ami_user = $7, ami_pass = $8, status = $9, enabled = $10, last_seen = $11, 
			last_error = $12, updated_at = $13
		WHERE id = $1`

	gateway.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query,
		gateway.ID, gateway.Name, gateway.Description, gateway.Location, gateway.AMIHost, gateway.AMIPort,
		gateway.AMIUser, gateway.AMIPass, gateway.Status, gateway.Enabled, gateway.LastSeen, gateway.LastError,
		gateway.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("PostgresGatewayRepository.UpdateGateway: failed to update gateway: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateGatewayHeartbeat updates the last seen timestamp and sets status to online.
func (r *PostgresGatewayRepository) UpdateGatewayHeartbeat(ctx context.Context, id string) error {
	query := `
		UPDATE gateways
		SET last_seen = $2, status = 'online', updated_at = $3
		WHERE id = $1`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, id, now, now)
	if err != nil {
		return fmt.Errorf("PostgresGatewayRepository.UpdateGatewayHeartbeat: failed to update heartbeat: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteGateway deletes a gateway by its ID.
func (r *PostgresGatewayRepository) DeleteGateway(ctx context.Context, id string) error {
	query := `DELETE FROM gateways WHERE id = $1`
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("PostgresGatewayRepository.DeleteGateway: failed to delete gateway: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetGatewayStats retrieves statistics about gateways.
func (r *PostgresGatewayRepository) GetGatewayStats(ctx context.Context) (total, online, offline int, err error) {
	// Total gateways
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM gateways").Scan(&total)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("PostgresGatewayRepository.GetGatewayStats: failed to get total gateways: %w", err)
	}

	// Online gateways (seen within last 5 minutes)
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM gateways WHERE status = 'online' AND last_seen > NOW() - INTERVAL '5 minutes'").Scan(&online)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("PostgresGatewayRepository.GetGatewayStats: failed to get online gateways: %w", err)
	}

	// Offline gateways
	offline = total - online

	return total, online, offline, nil
}
