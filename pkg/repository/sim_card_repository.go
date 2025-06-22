package repository

import (
	"context"
	// "errors" // ErrNotFound is now in repository.go
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/e173-gateway/e173_go_gateway/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


// postgresSIMCardRepository is the PostgreSQL implementation of SIMCardRepository.
type postgresSIMCardRepository struct {
	db *pgxpool.Pool
}

// NewPostgresSIMCardRepository creates a new instance of postgresSIMCardRepository.
func NewPostgresSIMCardRepository(db *pgxpool.Pool) SIMCardRepository {
	return &postgresSIMCardRepository{db: db}
}

func (r *postgresSIMCardRepository) CreateSIMCard(ctx context.Context, sim *models.SIMCard) (int64, error) {
	logger := logging.Logger.WithContext(ctx)
	query := `
		INSERT INTO sim_cards (
			modem_id, iccid, imsi, msisdn, operator_name, network_country_code,
			balance, balance_currency, balance_last_checked_at,
			data_allowance_mb, data_used_mb, status,
			pin1, puk1, pin2, puk2,
			activation_date, expiry_date, recharge_history, notes,
			cell_id, lac, psc, rscp, ecio, bts_info_history
			-- created_at and updated_at have defaults
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		) RETURNING id`

	var id int64
	err := r.db.QueryRow(ctx, query,
		sim.ModemID, sim.ICCID, sim.IMSI, sim.MSISDN, sim.OperatorName, sim.NetworkCountryCode,
		sim.Balance, sim.BalanceCurrency, sim.BalanceLastCheckedAt,
		sim.DataAllowanceMB, sim.DataUsedMB, sim.Status,
		sim.PIN1, sim.PUK1, sim.PIN2, sim.PUK2,
		sim.ActivationDate, sim.ExpiryDate, sim.RechargeHistory, sim.Notes,
		sim.CellID, sim.LAC, sim.PSC, sim.RSCP, sim.ECIO, sim.BTSInfoHistory,
	).Scan(&id)

	if err != nil {
		logger.WithError(err).Error("Error creating SIM card in database")
		return 0, err
	}
	logger.WithField("sim_id", id).Info("Successfully created SIM card")
	return id, nil
}

func (r *postgresSIMCardRepository) GetSIMCardByID(ctx context.Context, id int64) (*models.SIMCard, error) {
	logger := logging.Logger.WithContext(ctx)
	query := `
		SELECT
			id, modem_id, iccid, imsi, msisdn, operator_name, network_country_code,
			balance, balance_currency, balance_last_checked_at,
			data_allowance_mb, data_used_mb, status,
			pin1, puk1, pin2, puk2,
			activation_date, expiry_date, recharge_history, notes,
			cell_id, lac, psc, rscp, ecio, bts_info_history,
			created_at, updated_at
		FROM sim_cards
		WHERE id = $1`

	sim := &models.SIMCard{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sim.ID, &sim.ModemID, &sim.ICCID, &sim.IMSI, &sim.MSISDN, &sim.OperatorName, &sim.NetworkCountryCode,
		&sim.Balance, &sim.BalanceCurrency, &sim.BalanceLastCheckedAt,
		&sim.DataAllowanceMB, &sim.DataUsedMB, &sim.Status,
		&sim.PIN1, &sim.PUK1, &sim.PIN2, &sim.PUK2,
		&sim.ActivationDate, &sim.ExpiryDate, &sim.RechargeHistory, &sim.Notes,
		&sim.CellID, &sim.LAC, &sim.PSC, &sim.RSCP, &sim.ECIO, &sim.BTSInfoHistory,
		&sim.CreatedAt, &sim.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithField("sim_id", id).Warn("SIM card not found by ID")
			return nil, ErrNotFound // Or return nil, nil or a specific error type
		}
		logger.WithError(err).WithField("sim_id", id).Error("Error getting SIM card by ID from database")
		return nil, err
	}
	logger.WithField("sim_id", id).Info("Successfully retrieved SIM card by ID")
	return sim, nil
}

func (r *postgresSIMCardRepository) GetSIMCardByICCID(ctx context.Context, iccid string) (*models.SIMCard, error) {
	logger := logging.Logger.WithContext(ctx)
	query := `
		SELECT
			id, modem_id, iccid, imsi, msisdn, operator_name, network_country_code,
			balance, balance_currency, balance_last_checked_at,
			data_allowance_mb, data_used_mb, status,
			pin1, puk1, pin2, puk2,
			activation_date, expiry_date, recharge_history, notes,
			cell_id, lac, psc, rscp, ecio, bts_info_history,
			created_at, updated_at
		FROM sim_cards
		WHERE iccid = $1`

	sim := &models.SIMCard{}
	err := r.db.QueryRow(ctx, query, iccid).Scan(
		&sim.ID, &sim.ModemID, &sim.ICCID, &sim.IMSI, &sim.MSISDN, &sim.OperatorName, &sim.NetworkCountryCode,
		&sim.Balance, &sim.BalanceCurrency, &sim.BalanceLastCheckedAt,
		&sim.DataAllowanceMB, &sim.DataUsedMB, &sim.Status,
		&sim.PIN1, &sim.PUK1, &sim.PIN2, &sim.PUK2,
		&sim.ActivationDate, &sim.ExpiryDate, &sim.RechargeHistory, &sim.Notes,
		&sim.CellID, &sim.LAC, &sim.PSC, &sim.RSCP, &sim.ECIO, &sim.BTSInfoHistory,
		&sim.CreatedAt, &sim.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithField("iccid", iccid).Warn("SIM card not found by ICCID")
			return nil, ErrNotFound
		}
		logger.WithError(err).WithField("iccid", iccid).Error("Error getting SIM card by ICCID from database")
		return nil, err
	}
	logger.WithField("iccid", iccid).Info("Successfully retrieved SIM card by ICCID")
	return sim, nil
}

func (r *postgresSIMCardRepository) GetAllSIMCards(ctx context.Context) ([]models.SIMCard, error) {
	logger := logging.Logger.WithContext(ctx)
	query := `
		SELECT
			id, modem_id, iccid, imsi, msisdn, operator_name, network_country_code,
			balance, balance_currency, balance_last_checked_at,
			data_allowance_mb, data_used_mb, status,
			pin1, puk1, pin2, puk2,
			activation_date, expiry_date, recharge_history, notes,
			cell_id, lac, psc, rscp, ecio, bts_info_history,
			created_at, updated_at
		FROM sim_cards
		ORDER BY id ASC` // Or any other default ordering you prefer

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.WithError(err).Error("Error getting all SIM cards from database")
		return nil, err
	}
	defer rows.Close()

	var sims []models.SIMCard
	for rows.Next() {
		var sim models.SIMCard
		err := rows.Scan(
			&sim.ID, &sim.ModemID, &sim.ICCID, &sim.IMSI, &sim.MSISDN, &sim.OperatorName, &sim.NetworkCountryCode,
			&sim.Balance, &sim.BalanceCurrency, &sim.BalanceLastCheckedAt,
			&sim.DataAllowanceMB, &sim.DataUsedMB, &sim.Status,
			&sim.PIN1, &sim.PUK1, &sim.PIN2, &sim.PUK2,
			&sim.ActivationDate, &sim.ExpiryDate, &sim.RechargeHistory, &sim.Notes,
			&sim.CellID, &sim.LAC, &sim.PSC, &sim.RSCP, &sim.ECIO, &sim.BTSInfoHistory,
			&sim.CreatedAt, &sim.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("Error scanning SIM card row")
			// Decide if you want to return partial results or an error
			return nil, err 
		}
		sims = append(sims, sim)
	}

	if rows.Err() != nil {
		logger.WithError(rows.Err()).Error("Error iterating over SIM card rows")
		return nil, rows.Err()
	}

	logger.WithField("count", len(sims)).Info("Successfully retrieved all SIM cards")
	return sims, nil
}

func (r *postgresSIMCardRepository) UpdateSIMCard(ctx context.Context, sim *models.SIMCard) error {
	logger := logging.Logger.WithContext(ctx).WithField("sim_id", sim.ID)
	query := `
		UPDATE sim_cards SET
			modem_id = $1,
			iccid = $2,
			imsi = $3,
			msisdn = $4,
			operator_name = $5,
			network_country_code = $6,
			balance = $7,
			balance_currency = $8,
			balance_last_checked_at = $9,
			data_allowance_mb = $10,
			data_used_mb = $11,
			status = $12,
			pin1 = $13,
			puk1 = $14,
			pin2 = $15,
			puk2 = $16,
			activation_date = $17,
			expiry_date = $18,
			recharge_history = $19,
			notes = $20,
			cell_id = $21,
			lac = $22,
			psc = $23,
			rscp = $24,
			ecio = $25,
			bts_info_history = $26
			-- updated_at is handled by a trigger
		WHERE id = $27`

	commandTag, err := r.db.Exec(ctx, query,
		sim.ModemID, sim.ICCID, sim.IMSI, sim.MSISDN, sim.OperatorName, sim.NetworkCountryCode,
		sim.Balance, sim.BalanceCurrency, sim.BalanceLastCheckedAt,
		sim.DataAllowanceMB, sim.DataUsedMB, sim.Status,
		sim.PIN1, sim.PUK1, sim.PIN2, sim.PUK2,
		sim.ActivationDate, sim.ExpiryDate, sim.RechargeHistory, sim.Notes,
		sim.CellID, sim.LAC, sim.PSC, sim.RSCP, sim.ECIO, sim.BTSInfoHistory,
		sim.ID, // For the WHERE clause
	)

	if err != nil {
		logger.WithError(err).Error("Error updating SIM card in database")
		return err
	}

	if commandTag.RowsAffected() == 0 {
		logger.Warn("No SIM card found with given ID to update, or data was the same")
		return ErrNotFound // Or a more specific error if needed
	}

	logger.Info("Successfully updated SIM card")
	return nil
}

func (r *postgresSIMCardRepository) DeleteSIMCard(ctx context.Context, id int64) error {
	logger := logging.Logger.WithContext(ctx).WithField("sim_id", id)
	query := `DELETE FROM sim_cards WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.WithError(err).Error("Error deleting SIM card from database")
		return err
	}

	if commandTag.RowsAffected() == 0 {
		logger.Warn("No SIM card found with given ID to delete")
		return ErrNotFound
	}

	logger.Info("Successfully deleted SIM card")
	return nil
}
