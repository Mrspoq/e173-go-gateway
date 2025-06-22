package repository

import (
	"context"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)


// postgresModemRepository implements ModemRepository using a PostgreSQL database.
type postgresModemRepository struct {
	db *pgxpool.Pool
}

// NewPostgresModemRepository creates a new instance of postgresModemRepository.
func NewPostgresModemRepository(db *pgxpool.Pool) ModemRepository {
	return &postgresModemRepository{db: db}
}

// CreateModem inserts a new modem record into the database and returns the new modem's ID.
func (r *postgresModemRepository) CreateModem(ctx context.Context, modem *models.Modem) (int, error) {
	query := `
		INSERT INTO modems (
			device_path, imei, imsi, model, manufacturer, firmware_version, 
			signal_strength_dbm, network_operator_name, network_registration_status, 
			status, last_seen_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, created_at, updated_at` // Also return created_at and updated_at

	err := r.db.QueryRow(ctx, query,
		modem.DevicePath, modem.IMEI, modem.IMSI, modem.Model, modem.Manufacturer, modem.FirmwareVersion,
		modem.SignalStrengthDBM, modem.NetworkOperatorName, modem.NetworkRegistrationStatus,
		modem.Status, modem.LastSeenAt,
	).Scan(&modem.ID, &modem.CreatedAt, &modem.UpdatedAt) // Scan the returned id, created_at, updated_at

	if err != nil {
		return 0, err
	}
	return modem.ID, nil
}

// GetAllModems retrieves all modem records from the database.
func (r *postgresModemRepository) GetAllModems(ctx context.Context) ([]models.Modem, error) {
	query := `
		SELECT id, device_path, imei, imsi, model, manufacturer, firmware_version, 
		       signal_strength_dbm, network_operator_name, network_registration_status, 
		       status, last_seen_at, created_at, updated_at 
		FROM modems ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modems []models.Modem
	for rows.Next() {
		var m models.Modem
		err := rows.Scan(
			&m.ID, &m.DevicePath, &m.IMEI, &m.IMSI, &m.Model, &m.Manufacturer, &m.FirmwareVersion,
			&m.SignalStrengthDBM, &m.NetworkOperatorName, &m.NetworkRegistrationStatus,
			&m.Status, &m.LastSeenAt, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		modems = append(modems, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return modems, nil
}

// GetModemByID retrieves a single modem record from the database by its ID.
func (r *postgresModemRepository) GetModemByID(ctx context.Context, id int) (*models.Modem, error) {
	query := `
		SELECT id, device_path, imei, imsi, model, manufacturer, firmware_version, 
		       signal_strength_dbm, network_operator_name, network_registration_status, 
		       status, last_seen_at, created_at, updated_at 
		FROM modems 
		WHERE id = $1`

	var m models.Modem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.DevicePath, &m.IMEI, &m.IMSI, &m.Model, &m.Manufacturer, &m.FirmwareVersion,
		&m.SignalStrengthDBM, &m.NetworkOperatorName, &m.NetworkRegistrationStatus,
		&m.Status, &m.LastSeenAt, &m.CreatedAt, &m.UpdatedAt,
	)

	if err != nil {
		// Consider using a more specific error check, e.g., for pgx.ErrNoRows
		// and returning ErrNotFound from the repository package.
		return nil, err 
	}
	return &m, nil
}
