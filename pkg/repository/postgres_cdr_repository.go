package repository

import (
	"context"
	"fmt"
	"database/sql"

	"github.com/e173-gateway/e173_go_gateway/pkg/models" // Corrected model import path
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresCdrRepository implements CdrRepository for PostgreSQL.
// Note: Renamed from PostgresCDRRepository to PostgresCdrRepository for consistency with interface name.
type PostgresCdrRepository struct {
	db *pgxpool.Pool
}

// NewPostgresCdrRepository creates a new PostgresCdrRepository.
func NewPostgresCdrRepository(db *pgxpool.Pool) *PostgresCdrRepository {
	return &PostgresCdrRepository{db: db}
}

// CreateCdr inserts a new Cdr into the database.
func (r *PostgresCdrRepository) CreateCdr(ctx context.Context, cdr *models.Cdr) error {
	query := `
		INSERT INTO call_detail_records (
			asterisk_unique_id, source_number, destination_number,
			call_start_time, call_answer_time, call_end_time,
			duration_seconds, billable_duration_seconds, modem_id, sim_card_id,
			call_direction, customer_id, total_cost, cost_per_minute,
			is_spam, spam_reason, disposition, hangup_cause, recorded_audio_path
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
			$11, $12, $13, $14, $15, $16, $17, $18, $19
		) RETURNING id`

	var returnedID int64
	err := r.db.QueryRow(ctx, query,
		cdr.UniqueID, cdr.CallerIDNum, cdr.ConnectedLineNum,
		cdr.StartTime, cdr.AnswerTime, cdr.EndTime,
		cdr.Duration, cdr.BillableSeconds, cdr.ModemID, cdr.SimCardID,
		cdr.CallDirection, cdr.SipCustomerID, cdr.Cost, cdr.CustomerPrice,
		cdr.IsSpam, cdr.SpamDetectionMethod, cdr.Disposition, cdr.Cause, nil,
	).Scan(&returnedID)

	if err != nil {
		return fmt.Errorf("PostgresCdrRepository.CreateCdr: failed to insert CDR: %w", err)
	}
	// Optional: Log the returnedID if needed, or assign to cdr.ID if the model field should be updated post-insert.
	// cdr.ID = returnedID 
	return nil
}

// UpdateCdr updates an existing Cdr in the database.
func (r *PostgresCdrRepository) UpdateCdr(ctx context.Context, cdr *models.Cdr) error {
	query := `
		UPDATE call_detail_records
		SET updated_at = CURRENT_TIMESTAMP,
			asterisk_unique_id = $2, source_number = $3, destination_number = $4,
			call_start_time = $5, call_answer_time = $6, call_end_time = $7,
			duration_seconds = $8, billable_duration_seconds = $9, modem_id = $10, sim_card_id = $11,
			call_direction = $12, customer_id = $13, total_cost = $14, cost_per_minute = $15,
			is_spam = $16, spam_reason = $17, disposition = $18, hangup_cause = $19
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query,
		cdr.ID,
		cdr.UniqueID, cdr.CallerIDNum, cdr.ConnectedLineNum,
		cdr.StartTime, cdr.AnswerTime, cdr.EndTime,
		cdr.Duration, cdr.BillableSeconds, cdr.ModemID, cdr.SimCardID,
		cdr.CallDirection, cdr.SipCustomerID, cdr.Cost, cdr.CustomerPrice,
		cdr.IsSpam, cdr.SpamDetectionMethod, cdr.Disposition, cdr.Cause,
	)

	if err != nil {
		return fmt.Errorf("PostgresCdrRepository.UpdateCdr: failed to update CDR: %w", err)
	}
	return nil
}

// GetCdrByID retrieves a Cdr by its ID.
func (r *PostgresCdrRepository) GetCdrByID(ctx context.Context, id uuid.UUID) (*models.Cdr, error) {
	// Since the database uses bigint for ID, we need to convert or handle this differently
	// For now, we'll return an error as the schema mismatch needs to be resolved
	return nil, fmt.Errorf("PostgresCdrRepository.GetCdrByID: ID type mismatch - database uses bigint, model uses UUID")
}

// GetCdrByDatabaseID retrieves a Cdr by its database ID (bigint).
func (r *PostgresCdrRepository) GetCdrByDatabaseID(ctx context.Context, id int64) (*models.Cdr, error) {
	query := `
		SELECT
			id, created_at, updated_at, asterisk_unique_id, source_number, destination_number,
			call_start_time, call_answer_time, call_end_time,
			duration_seconds, billable_duration_seconds, modem_id, sim_card_id,
			call_direction, customer_id, total_cost, cost_per_minute,
			is_spam, spam_reason, disposition, hangup_cause
		FROM call_detail_records
		WHERE id = $1`

	cdr := &models.Cdr{}
	var dbID int64
	var modemID, simCardID, customerID sql.NullInt64
	var sourceNum, destNum, callDir, spamReason, disposition, hangupCause sql.NullString
	var callAnswerTime sql.NullTime
	var duration, billable sql.NullInt32
	var totalCost, costPerMin sql.NullFloat64
	var isSpam sql.NullBool

	err := r.db.QueryRow(ctx, query, id).Scan(
		&dbID, &cdr.CreatedAt, &cdr.UpdatedAt, &cdr.UniqueID, &sourceNum, &destNum,
		&cdr.StartTime, &callAnswerTime, &cdr.EndTime,
		&duration, &billable, &modemID, &simCardID,
		&callDir, &customerID, &totalCost, &costPerMin,
		&isSpam, &spamReason, &disposition, &hangupCause,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("PostgresCdrRepository.GetCdrByDatabaseID: failed to scan CDR: %w", err)
	}

	// Map nullable fields
	if sourceNum.Valid {
		cdr.CallerIDNum = &sourceNum.String
	}
	if destNum.Valid {
		cdr.ConnectedLineNum = &destNum.String
	}
	if callDir.Valid {
		cdr.CallDirection = &callDir.String
	}
	if callAnswerTime.Valid {
		cdr.AnswerTime = &callAnswerTime.Time
	}
	if duration.Valid {
		d := int(duration.Int32)
		cdr.Duration = &d
	}
	if billable.Valid {
		b := int(billable.Int32)
		cdr.BillableSeconds = &b
	}
	if modemID.Valid {
		m := int(modemID.Int64)
		cdr.ModemID = &m
	}
	if simCardID.Valid {
		s := int(simCardID.Int64)
		cdr.SimCardID = &s
	}
	if customerID.Valid {
		c := int(customerID.Int64)
		cdr.SipCustomerID = &c
	}
	if totalCost.Valid {
		cdr.Cost = &totalCost.Float64
	}
	if costPerMin.Valid {
		cdr.CustomerPrice = &costPerMin.Float64
	}
	if isSpam.Valid {
		cdr.IsSpam = &isSpam.Bool
	}
	if spamReason.Valid {
		cdr.SpamDetectionMethod = &spamReason.String
	}
	if disposition.Valid {
		cdr.Disposition = &disposition.String
	}
	if hangupCause.Valid {
		cdr.Cause = &hangupCause.String
	}

	return cdr, nil
}

// GetCdrByAsteriskUniqueID retrieves a Cdr by its Asterisk Unique ID.
// This method can be added to the CdrRepository interface if needed.
func (r *PostgresCdrRepository) GetCdrByAsteriskUniqueID(ctx context.Context, asteriskUniqueID string) (*models.Cdr, error) {
	query := `
		SELECT
			id, created_at, updated_at, asterisk_unique_id, source_number, destination_number,
			call_start_time, call_answer_time, call_end_time,
			duration_seconds, billable_duration_seconds, modem_id, sim_card_id,
			call_direction, customer_id, total_cost, cost_per_minute,
			is_spam, spam_reason, disposition, hangup_cause
		FROM call_detail_records
		WHERE asterisk_unique_id = $1`

	cdr := &models.Cdr{}
	var dbID int64
	var modemID, simCardID, customerID sql.NullInt64
	var sourceNum, destNum, callDir, spamReason, disposition, hangupCause sql.NullString
	var callAnswerTime sql.NullTime
	var duration, billable sql.NullInt32
	var totalCost, costPerMin sql.NullFloat64
	var isSpam sql.NullBool

	err := r.db.QueryRow(ctx, query, asteriskUniqueID).Scan(
		&dbID, &cdr.CreatedAt, &cdr.UpdatedAt, &cdr.UniqueID, &sourceNum, &destNum,
		&cdr.StartTime, &callAnswerTime, &cdr.EndTime,
		&duration, &billable, &modemID, &simCardID,
		&callDir, &customerID, &totalCost, &costPerMin,
		&isSpam, &spamReason, &disposition, &hangupCause,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("PostgresCdrRepository.GetCdrByAsteriskUniqueID: failed to scan CDR: %w", err)
	}

	// Map nullable fields
	if sourceNum.Valid {
		cdr.CallerIDNum = &sourceNum.String
	}
	if destNum.Valid {
		cdr.ConnectedLineNum = &destNum.String
	}
	if callDir.Valid {
		cdr.CallDirection = &callDir.String
	}
	if callAnswerTime.Valid {
		cdr.AnswerTime = &callAnswerTime.Time
	}
	if duration.Valid {
		d := int(duration.Int32)
		cdr.Duration = &d
	}
	if billable.Valid {
		b := int(billable.Int32)
		cdr.BillableSeconds = &b
	}
	if modemID.Valid {
		m := int(modemID.Int64)
		cdr.ModemID = &m
	}
	if simCardID.Valid {
		s := int(simCardID.Int64)
		cdr.SimCardID = &s
	}
	if customerID.Valid {
		c := int(customerID.Int64)
		cdr.SipCustomerID = &c
	}
	if totalCost.Valid {
		cdr.Cost = &totalCost.Float64
	}
	if costPerMin.Valid {
		cdr.CustomerPrice = &costPerMin.Float64
	}
	if isSpam.Valid {
		cdr.IsSpam = &isSpam.Bool
	}
	if spamReason.Valid {
		cdr.SpamDetectionMethod = &spamReason.String
	}
	if disposition.Valid {
		cdr.Disposition = &disposition.String
	}
	if hangupCause.Valid {
		cdr.Cause = &hangupCause.String
	}

	return cdr, nil
}

// GetRecentCDRs retrieves the most recent CDRs from the database, limited by the specified count.
func (r *PostgresCdrRepository) GetRecentCDRs(ctx context.Context, limit int) ([]*models.Cdr, error) {
	query := `
		SELECT 
			id, created_at, updated_at, asterisk_unique_id, source_number, destination_number,
			call_start_time, call_answer_time, call_end_time,
			duration_seconds, billable_duration_seconds, modem_id, sim_card_id,
			call_direction, customer_id, total_cost, cost_per_minute,
			is_spam, spam_reason, disposition, hangup_cause
		FROM call_detail_records
		ORDER BY call_start_time DESC
		LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("PostgresCdrRepository.GetRecentCDRs: failed to query recent CDRs: %w", err)
	}
	defer rows.Close()

	var cdrs []*models.Cdr
	for rows.Next() {
		cdr := &models.Cdr{}
		var dbID int64
		var modemID, simCardID, customerID sql.NullInt64
		var sourceNum, destNum, callDir, spamReason, disposition, hangupCause sql.NullString
		var callAnswerTime sql.NullTime
		var duration, billable sql.NullInt32
		var totalCost, costPerMin sql.NullFloat64
		var isSpam sql.NullBool

		err := rows.Scan(
			&dbID, &cdr.CreatedAt, &cdr.UpdatedAt, &cdr.UniqueID, &sourceNum, &destNum,
			&cdr.StartTime, &callAnswerTime, &cdr.EndTime,
			&duration, &billable, &modemID, &simCardID,
			&callDir, &customerID, &totalCost, &costPerMin,
			&isSpam, &spamReason, &disposition, &hangupCause,
		)
		if err != nil {
			return nil, fmt.Errorf("PostgresCdrRepository.GetRecentCDRs: failed to scan CDR: %w", err)
		}

		// Map nullable fields
		if sourceNum.Valid {
			cdr.CallerIDNum = &sourceNum.String
		}
		if destNum.Valid {
			cdr.ConnectedLineNum = &destNum.String
		}
		if callDir.Valid {
			cdr.CallDirection = &callDir.String
		}
		if callAnswerTime.Valid {
			cdr.AnswerTime = &callAnswerTime.Time
		}
		if duration.Valid {
			d := int(duration.Int32)
			cdr.Duration = &d
		}
		if billable.Valid {
			b := int(billable.Int32)
			cdr.BillableSeconds = &b
		}
		if modemID.Valid {
			m := int(modemID.Int64)
			cdr.ModemID = &m
		}
		if simCardID.Valid {
			s := int(simCardID.Int64)
			cdr.SimCardID = &s
		}
		if customerID.Valid {
			c := int(customerID.Int64)
			cdr.SipCustomerID = &c
		}
		if totalCost.Valid {
			cdr.Cost = &totalCost.Float64
		}
		if costPerMin.Valid {
			cdr.CustomerPrice = &costPerMin.Float64
		}
		if isSpam.Valid {
			cdr.IsSpam = &isSpam.Bool
		}
		if spamReason.Valid {
			cdr.SpamDetectionMethod = &spamReason.String
		}
		if disposition.Valid {
			cdr.Disposition = &disposition.String
		}
		if hangupCause.Valid {
			cdr.Cause = &hangupCause.String
		}

		cdrs = append(cdrs, cdr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PostgresCdrRepository.GetRecentCDRs: row iteration error: %w", err)
	}

	return cdrs, nil
}

// GetCallStats retrieves call statistics from the database.
func (r *PostgresCdrRepository) GetCallStats(ctx context.Context) (int, int, int, error) {
	var totalCalls, answeredCalls, missedCalls int

	// Total calls
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM call_detail_records").Scan(&totalCalls)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("PostgresCdrRepository.GetCallStats: failed to get total calls: %w", err)
	}

	// Answered calls (calls with answer time)
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM call_detail_records WHERE call_answer_time IS NOT NULL").Scan(&answeredCalls)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("PostgresCdrRepository.GetCallStats: failed to get answered calls: %w", err)
	}

	// Missed calls (calls without answer time)
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM call_detail_records WHERE call_answer_time IS NULL").Scan(&missedCalls)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("PostgresCdrRepository.GetCallStats: failed to get missed calls: %w", err)
	}

	return totalCalls, answeredCalls, missedCalls, nil
}

// GetSpamStats retrieves spam statistics from the database.
func (r *PostgresCdrRepository) GetSpamStats(ctx context.Context) (int, int, error) {
	var totalSpamCalls, blockedNumbers int

	// Total spam calls
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM call_detail_records WHERE is_spam = true").Scan(&totalSpamCalls)
	if err != nil {
		return 0, 0, fmt.Errorf("PostgresCdrRepository.GetSpamStats: failed to get total spam calls: %w", err)
	}

	// For blocked numbers, we would need a separate table or check spam_reason
	// For now, return distinct spam source numbers as a proxy
	err = r.db.QueryRow(ctx, "SELECT COUNT(DISTINCT source_number) FROM call_detail_records WHERE is_spam = true AND source_number IS NOT NULL").Scan(&blockedNumbers)
	if err != nil {
		return 0, 0, fmt.Errorf("PostgresCdrRepository.GetSpamStats: failed to get blocked numbers: %w", err)
	}

	return totalSpamCalls, blockedNumbers, nil
}
