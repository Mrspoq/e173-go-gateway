package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "time"

	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/jmoiron/sqlx"
)

// RechargeRepository defines the interface for recharge operations
type RechargeRepository interface {
	// Recharge code operations
	CreateRechargeCode(ctx context.Context, code *models.RechargeCode) error
	GetRechargeCodeByID(ctx context.Context, id int64) (*models.RechargeCode, error)
	GetRechargeCodeByCode(ctx context.Context, code, operator string) (*models.RechargeCode, error)
	UpdateRechargeCode(ctx context.Context, code *models.RechargeCode) error
	ListRechargeCodes(ctx context.Context, simCardID int64, status string) ([]*models.RechargeCode, error)
	
	// Batch operations
	CreateBatch(ctx context.Context, batch *models.RechargeBatch) error
	GetBatchByID(ctx context.Context, id int64) (*models.RechargeBatch, error)
	UpdateBatch(ctx context.Context, batch *models.RechargeBatch) error
	ListBatches(ctx context.Context, status string) ([]*models.RechargeBatch, error)
	
	// History operations
	CreateRechargeHistory(ctx context.Context, history *models.RechargeHistory) error
	GetRechargeHistory(ctx context.Context, simCardID int64, limit int) ([]*models.RechargeHistory, error)
	
	// Bulk operations
	CreateBulkRechargeCodes(ctx context.Context, codes []*models.RechargeCode) error
	GetSimCardsForAutoRecharge(ctx context.Context, threshold float64) ([]*models.SIMCard, error)
}

type rechargeRepository struct {
	db *sqlx.DB
}

// NewRechargeRepository creates a new recharge repository
func NewRechargeRepository(db *sqlx.DB) RechargeRepository {
	return &rechargeRepository{db: db}
}

// CreateRechargeCode creates a new recharge code
func (r *rechargeRepository) CreateRechargeCode(ctx context.Context, code *models.RechargeCode) error {
	query := `
		INSERT INTO recharge_codes (sim_card_id, code, amount, operator, status, expiry_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
		
	err := r.db.QueryRowContext(ctx, query,
		code.SimCardID,
		code.Code,
		code.Amount,
		code.Operator,
		code.Status,
		code.ExpiryDate,
		code.CreatedBy,
	).Scan(&code.ID, &code.CreatedAt, &code.UpdatedAt)
	
	return err
}

// GetRechargeCodeByID retrieves a recharge code by ID
func (r *rechargeRepository) GetRechargeCodeByID(ctx context.Context, id int64) (*models.RechargeCode, error) {
	var code models.RechargeCode
	query := `SELECT * FROM recharge_codes WHERE id = $1`
	
	err := r.db.GetContext(ctx, &code, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	
	return &code, err
}

// GetRechargeCodeByCode retrieves a recharge code by code and operator
func (r *rechargeRepository) GetRechargeCodeByCode(ctx context.Context, code, operator string) (*models.RechargeCode, error) {
	var rechargeCode models.RechargeCode
	query := `SELECT * FROM recharge_codes WHERE code = $1 AND operator = $2`
	
	err := r.db.GetContext(ctx, &rechargeCode, query, code, operator)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	
	return &rechargeCode, err
}

// UpdateRechargeCode updates a recharge code
func (r *rechargeRepository) UpdateRechargeCode(ctx context.Context, code *models.RechargeCode) error {
	query := `
		UPDATE recharge_codes 
		SET status = $2, used_at = $3, response_message = $4
		WHERE id = $1`
		
	_, err := r.db.ExecContext(ctx, query,
		code.ID,
		code.Status,
		code.UsedAt,
		code.ResponseMessage,
	)
	
	return err
}

// ListRechargeCodes lists recharge codes with optional filters
func (r *rechargeRepository) ListRechargeCodes(ctx context.Context, simCardID int64, status string) ([]*models.RechargeCode, error) {
	var codes []*models.RechargeCode
	query := `SELECT * FROM recharge_codes WHERE 1=1`
	args := []interface{}{}
	argCount := 0
	
	if simCardID > 0 {
		argCount++
		query += fmt.Sprintf(" AND sim_card_id = $%d", argCount)
		args = append(args, simCardID)
	}
	
	if status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}
	
	query += " ORDER BY created_at DESC"
	
	err := r.db.SelectContext(ctx, &codes, query, args...)
	return codes, err
}

// CreateBatch creates a new recharge batch
func (r *rechargeRepository) CreateBatch(ctx context.Context, batch *models.RechargeBatch) error {
	query := `
		INSERT INTO recharge_batches (name, description, total_codes, total_amount, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
		
	err := r.db.QueryRowContext(ctx, query,
		batch.Name,
		batch.Description,
		batch.TotalCodes,
		batch.TotalAmount,
		batch.Status,
		batch.CreatedBy,
	).Scan(&batch.ID, &batch.CreatedAt, &batch.UpdatedAt)
	
	return err
}

// GetBatchByID retrieves a batch by ID
func (r *rechargeRepository) GetBatchByID(ctx context.Context, id int64) (*models.RechargeBatch, error) {
	var batch models.RechargeBatch
	query := `SELECT * FROM recharge_batches WHERE id = $1`
	
	err := r.db.GetContext(ctx, &batch, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	
	return &batch, err
}

// UpdateBatch updates a recharge batch
func (r *rechargeRepository) UpdateBatch(ctx context.Context, batch *models.RechargeBatch) error {
	query := `
		UPDATE recharge_batches 
		SET status = $2, used_codes = $3, started_at = $4, completed_at = $5
		WHERE id = $1`
		
	_, err := r.db.ExecContext(ctx, query,
		batch.ID,
		batch.Status,
		batch.UsedCodes,
		batch.StartedAt,
		batch.CompletedAt,
	)
	
	return err
}

// ListBatches lists recharge batches with optional status filter
func (r *rechargeRepository) ListBatches(ctx context.Context, status string) ([]*models.RechargeBatch, error) {
	var batches []*models.RechargeBatch
	query := `SELECT * FROM recharge_batches`
	args := []interface{}{}
	
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}
	
	query += " ORDER BY created_at DESC"
	
	err := r.db.SelectContext(ctx, &batches, query, args...)
	return batches, err
}

// CreateRechargeHistory creates a new recharge history record
func (r *rechargeRepository) CreateRechargeHistory(ctx context.Context, history *models.RechargeHistory) error {
	query := `
		INSERT INTO recharge_history (
			sim_card_id, recharge_code_id, batch_id, phone_number, 
			amount, balance_before, balance_after, method, status, 
			error_message, attempts, processed_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, processed_at`
		
	err := r.db.QueryRowContext(ctx, query,
		history.SimCardID,
		history.RechargeCodeID,
		history.BatchID,
		history.PhoneNumber,
		history.Amount,
		history.BalanceBefore,
		history.BalanceAfter,
		history.Method,
		history.Status,
		history.ErrorMessage,
		history.Attempts,
		history.ProcessedBy,
	).Scan(&history.ID, &history.ProcessedAt)
	
	return err
}

// GetRechargeHistory retrieves recharge history for a SIM card
func (r *rechargeRepository) GetRechargeHistory(ctx context.Context, simCardID int64, limit int) ([]*models.RechargeHistory, error) {
	var history []*models.RechargeHistory
	query := `
		SELECT * FROM recharge_history 
		WHERE sim_card_id = $1 
		ORDER BY processed_at DESC 
		LIMIT $2`
	
	err := r.db.SelectContext(ctx, &history, query, simCardID, limit)
	return history, err
}

// CreateBulkRechargeCodes creates multiple recharge codes in a transaction
func (r *rechargeRepository) CreateBulkRechargeCodes(ctx context.Context, codes []*models.RechargeCode) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO recharge_codes (sim_card_id, code, amount, operator, status, expiry_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	for _, code := range codes {
		err = stmt.QueryRowContext(ctx,
			code.SimCardID,
			code.Code,
			code.Amount,
			code.Operator,
			code.Status,
			code.ExpiryDate,
			code.CreatedBy,
		).Scan(&code.ID, &code.CreatedAt, &code.UpdatedAt)
		
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

// GetSimCardsForAutoRecharge retrieves SIM cards that need auto-recharge
func (r *rechargeRepository) GetSimCardsForAutoRecharge(ctx context.Context, threshold float64) ([]*models.SIMCard, error) {
	var sims []*models.SIMCard
	query := `
		SELECT * FROM sim_cards 
		WHERE auto_recharge_enabled = true 
		AND status = $1 
		AND balance < auto_recharge_threshold
		AND (last_recharge_at IS NULL OR last_recharge_at < NOW() - INTERVAL '1 hour')
		ORDER BY balance ASC`
	
	err := r.db.SelectContext(ctx, &sims, query, "active")
	return sims, err
}