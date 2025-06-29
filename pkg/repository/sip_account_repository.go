package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/jmoiron/sqlx"
)

// SIPAccountRepository handles database operations for SIP accounts
type SIPAccountRepository interface {
	CreateSIPAccount(ctx context.Context, account *models.SIPAccount) error
	GetSIPAccountByID(ctx context.Context, id int64) (*models.SIPAccount, error)
	GetSIPAccountByUsername(ctx context.Context, username string) (*models.SIPAccount, error)
	GetSIPAccountsByCustomerID(ctx context.Context, customerID int64) ([]*models.SIPAccount, error)
	UpdateSIPAccount(ctx context.Context, account *models.SIPAccount) error
	DeleteSIPAccount(ctx context.Context, id int64) error
	ListSIPAccounts(ctx context.Context, limit, offset int) ([]*models.SIPAccount, error)
	SearchSIPAccounts(ctx context.Context, query string, limit, offset int) ([]*models.SIPAccount, error)
	
	// Permissions
	GetSIPAccountPermissions(ctx context.Context, accountID int64) (*models.SIPAccountPermission, error)
	UpdateSIPAccountPermissions(ctx context.Context, permissions *models.SIPAccountPermission) error
	
	// Registration
	CreateRegistration(ctx context.Context, registration *models.SIPRegistration) error
	UpdateRegistrationStatus(ctx context.Context, accountID int64, ip string, registered bool) error
	GetActiveRegistrations(ctx context.Context, accountID int64) ([]*models.SIPRegistration, error)
	
	// Usage
	RecordUsage(ctx context.Context, usage *models.SIPAccountUsage) error
	GetUsageByDate(ctx context.Context, accountID int64, date time.Time) (*models.SIPAccountUsage, error)
	GetUsageStats(ctx context.Context, accountID int64, startDate, endDate time.Time) ([]*models.SIPAccountUsage, error)
}

type sipAccountRepository struct {
	db *sqlx.DB
}

// NewSIPAccountRepository creates a new SIP account repository
func NewSIPAccountRepository(db *sqlx.DB) SIPAccountRepository {
	return &sipAccountRepository{db: db}
}

func (r *sipAccountRepository) CreateSIPAccount(ctx context.Context, account *models.SIPAccount) error {
	query := `
		INSERT INTO sip_accounts (
			customer_id, account_name, username, password, domain, extension,
			caller_id, caller_id_name, context, transport, nat_support,
			direct_media_support, encryption_enabled, codecs_allowed,
			max_concurrent_calls, status, notes, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		account.CustomerID,
		account.AccountName,
		account.Username,
		account.Password,
		account.Domain,
		account.Extension,
		account.CallerID,
		account.CallerIDName,
		account.Context,
		account.Transport,
		account.NATSupport,
		account.DirectMediaSupport,
		account.EncryptionEnabled,
		account.CodecsAllowed,
		account.MaxConcurrentCalls,
		account.Status,
		account.Notes,
		account.CreatedBy,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("username already exists")
		}
		return fmt.Errorf("failed to create SIP account: %w", err)
	}

	// Create default permissions
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO sip_account_permissions (sip_account_id)
		VALUES ($1)
	`, account.ID)

	return err
}

func (r *sipAccountRepository) GetSIPAccountByID(ctx context.Context, id int64) (*models.SIPAccount, error) {
	var account models.SIPAccount
	query := `
		SELECT 
			id, customer_id, account_name, username, password, domain, extension,
			caller_id, caller_id_name, context, transport, nat_support,
			direct_media_support, encryption_enabled, codecs_allowed,
			max_concurrent_calls, current_active_calls, status,
			last_registered_ip, last_registered_at, last_call_at,
			total_calls, total_minutes, notes, created_by, created_at, updated_at
		FROM sip_accounts
		WHERE id = $1`

	err := r.db.GetContext(ctx, &account, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get SIP account: %w", err)
	}

	return &account, nil
}

func (r *sipAccountRepository) GetSIPAccountByUsername(ctx context.Context, username string) (*models.SIPAccount, error) {
	var account models.SIPAccount
	query := `
		SELECT 
			id, customer_id, account_name, username, password, domain, extension,
			caller_id, caller_id_name, context, transport, nat_support,
			direct_media_support, encryption_enabled, codecs_allowed,
			max_concurrent_calls, current_active_calls, status,
			last_registered_ip, last_registered_at, last_call_at,
			total_calls, total_minutes, notes, created_by, created_at, updated_at
		FROM sip_accounts
		WHERE username = $1`

	err := r.db.GetContext(ctx, &account, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get SIP account by username: %w", err)
	}

	return &account, nil
}

func (r *sipAccountRepository) GetSIPAccountsByCustomerID(ctx context.Context, customerID int64) ([]*models.SIPAccount, error) {
	var accounts []*models.SIPAccount
	query := `
		SELECT 
			id, customer_id, account_name, username, password, domain, extension,
			caller_id, caller_id_name, context, transport, nat_support,
			direct_media_support, encryption_enabled, codecs_allowed,
			max_concurrent_calls, current_active_calls, status,
			last_registered_ip, last_registered_at, last_call_at,
			total_calls, total_minutes, notes, created_by, created_at, updated_at
		FROM sip_accounts
		WHERE customer_id = $1
		ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &accounts, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get SIP accounts for customer: %w", err)
	}

	return accounts, nil
}

func (r *sipAccountRepository) UpdateSIPAccount(ctx context.Context, account *models.SIPAccount) error {
	query := `
		UPDATE sip_accounts SET
			account_name = $2,
			extension = $3,
			caller_id = $4,
			caller_id_name = $5,
			context = $6,
			transport = $7,
			nat_support = $8,
			direct_media_support = $9,
			encryption_enabled = $10,
			codecs_allowed = $11,
			max_concurrent_calls = $12,
			status = $13,
			notes = $14
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.AccountName,
		account.Extension,
		account.CallerID,
		account.CallerIDName,
		account.Context,
		account.Transport,
		account.NATSupport,
		account.DirectMediaSupport,
		account.EncryptionEnabled,
		account.CodecsAllowed,
		account.MaxConcurrentCalls,
		account.Status,
		account.Notes,
	)

	if err != nil {
		return fmt.Errorf("failed to update SIP account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *sipAccountRepository) DeleteSIPAccount(ctx context.Context, id int64) error {
	query := `DELETE FROM sip_accounts WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete SIP account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *sipAccountRepository) ListSIPAccounts(ctx context.Context, limit, offset int) ([]*models.SIPAccount, error) {
	var accounts []*models.SIPAccount
	query := `
		SELECT 
			s.id, s.customer_id, s.account_name, s.username, s.password, s.domain, s.extension,
			s.caller_id, s.caller_id_name, s.context, s.transport, s.nat_support,
			s.direct_media_support, s.encryption_enabled, s.codecs_allowed,
			s.max_concurrent_calls, s.current_active_calls, s.status,
			s.last_registered_ip, s.last_registered_at, s.last_call_at,
			s.total_calls, s.total_minutes, s.notes, s.created_by, s.created_at, s.updated_at
		FROM sip_accounts s
		ORDER BY s.created_at DESC
		LIMIT $1 OFFSET $2`

	err := r.db.SelectContext(ctx, &accounts, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list SIP accounts: %w", err)
	}

	return accounts, nil
}

func (r *sipAccountRepository) SearchSIPAccounts(ctx context.Context, searchQuery string, limit, offset int) ([]*models.SIPAccount, error) {
	var accounts []*models.SIPAccount
	query := `
		SELECT 
			s.id, s.customer_id, s.account_name, s.username, s.password, s.domain, s.extension,
			s.caller_id, s.caller_id_name, s.context, s.transport, s.nat_support,
			s.direct_media_support, s.encryption_enabled, s.codecs_allowed,
			s.max_concurrent_calls, s.current_active_calls, s.status,
			s.last_registered_ip, s.last_registered_at, s.last_call_at,
			s.total_calls, s.total_minutes, s.notes, s.created_by, s.created_at, s.updated_at
		FROM sip_accounts s
		WHERE 
			s.account_name ILIKE $1 OR
			s.username ILIKE $1 OR
			s.extension ILIKE $1 OR
			s.caller_id ILIKE $1 OR
			s.caller_id_name ILIKE $1
		ORDER BY s.created_at DESC
		LIMIT $2 OFFSET $3`

	searchPattern := "%" + searchQuery + "%"
	err := r.db.SelectContext(ctx, &accounts, query, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search SIP accounts: %w", err)
	}

	return accounts, nil
}

func (r *sipAccountRepository) GetSIPAccountPermissions(ctx context.Context, accountID int64) (*models.SIPAccountPermission, error) {
	var permissions models.SIPAccountPermission
	query := `
		SELECT 
			id, sip_account_id, allow_international, allow_premium_numbers,
			allow_emergency_calls, allowed_countries, blocked_countries,
			allowed_prefixes, blocked_prefixes, time_restrictions,
			daily_call_limit, daily_minute_limit, monthly_call_limit,
			monthly_minute_limit, created_at, updated_at
		FROM sip_account_permissions
		WHERE sip_account_id = $1`

	err := r.db.GetContext(ctx, &permissions, query, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create default permissions if not found
			permissions = models.SIPAccountPermission{
				SIPAccountID:        accountID,
				AllowEmergencyCalls: true,
			}
			_, err = r.db.ExecContext(ctx, `
				INSERT INTO sip_account_permissions (sip_account_id, allow_emergency_calls)
				VALUES ($1, $2)
			`, accountID, true)
			if err != nil {
				return nil, fmt.Errorf("failed to create default permissions: %w", err)
			}
			return &permissions, nil
		}
		return nil, fmt.Errorf("failed to get SIP account permissions: %w", err)
	}

	return &permissions, nil
}

func (r *sipAccountRepository) UpdateSIPAccountPermissions(ctx context.Context, permissions *models.SIPAccountPermission) error {
	query := `
		UPDATE sip_account_permissions SET
			allow_international = $2,
			allow_premium_numbers = $3,
			allow_emergency_calls = $4,
			allowed_countries = $5,
			blocked_countries = $6,
			allowed_prefixes = $7,
			blocked_prefixes = $8,
			time_restrictions = $9,
			daily_call_limit = $10,
			daily_minute_limit = $11,
			monthly_call_limit = $12,
			monthly_minute_limit = $13
		WHERE sip_account_id = $1`

	result, err := r.db.ExecContext(ctx, query,
		permissions.SIPAccountID,
		permissions.AllowInternational,
		permissions.AllowPremiumNumbers,
		permissions.AllowEmergencyCalls,
		permissions.AllowedCountries,
		permissions.BlockedCountries,
		permissions.AllowedPrefixes,
		permissions.BlockedPrefixes,
		permissions.TimeRestrictions,
		permissions.DailyCallLimit,
		permissions.DailyMinuteLimit,
		permissions.MonthlyCallLimit,
		permissions.MonthlyMinuteLimit,
	)

	if err != nil {
		return fmt.Errorf("failed to update SIP account permissions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Insert if not exists
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO sip_account_permissions (
				sip_account_id, allow_international, allow_premium_numbers,
				allow_emergency_calls, allowed_countries, blocked_countries,
				allowed_prefixes, blocked_prefixes, time_restrictions,
				daily_call_limit, daily_minute_limit, monthly_call_limit,
				monthly_minute_limit
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			permissions.SIPAccountID,
			permissions.AllowInternational,
			permissions.AllowPremiumNumbers,
			permissions.AllowEmergencyCalls,
			permissions.AllowedCountries,
			permissions.BlockedCountries,
			permissions.AllowedPrefixes,
			permissions.BlockedPrefixes,
			permissions.TimeRestrictions,
			permissions.DailyCallLimit,
			permissions.DailyMinuteLimit,
			permissions.MonthlyCallLimit,
			permissions.MonthlyMinuteLimit,
		)
		if err != nil {
			return fmt.Errorf("failed to insert SIP account permissions: %w", err)
		}
	}

	return nil
}

func (r *sipAccountRepository) CreateRegistration(ctx context.Context, registration *models.SIPRegistration) error {
	query := `
		INSERT INTO sip_registrations (
			sip_account_id, contact_uri, source_ip, source_port,
			user_agent, expires_seconds, registered_at, expired_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		registration.SIPAccountID,
		registration.ContactURI,
		registration.SourceIP,
		registration.SourcePort,
		registration.UserAgent,
		registration.ExpiresSeconds,
		registration.RegisteredAt,
		registration.ExpiredAt,
	).Scan(&registration.ID)

	if err != nil {
		return fmt.Errorf("failed to create registration: %w", err)
	}

	// Update account's last registration info
	_, err = r.db.ExecContext(ctx, `
		UPDATE sip_accounts
		SET last_registered_ip = $2, last_registered_at = $3
		WHERE id = $1`,
		registration.SIPAccountID,
		registration.SourceIP,
		registration.RegisteredAt,
	)

	return err
}

func (r *sipAccountRepository) UpdateRegistrationStatus(ctx context.Context, accountID int64, ip string, registered bool) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if !registered {
		// Mark all active registrations as unregistered
		now := time.Now()
		_, err = tx.ExecContext(ctx, `
			UPDATE sip_registrations
			SET is_active = false, unregistered_at = $2
			WHERE sip_account_id = $1 AND is_active = true`,
			accountID, now,
		)
		if err != nil {
			return fmt.Errorf("failed to update registration status: %w", err)
		}
	}

	// Update account registration info
	if registered {
		_, err = tx.ExecContext(ctx, `
			UPDATE sip_accounts
			SET last_registered_ip = $2, last_registered_at = $3
			WHERE id = $1`,
			accountID, ip, time.Now(),
		)
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE sip_accounts
			SET last_registered_ip = NULL, last_registered_at = NULL
			WHERE id = $1`,
			accountID,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to update account registration: %w", err)
	}

	return tx.Commit()
}

func (r *sipAccountRepository) GetActiveRegistrations(ctx context.Context, accountID int64) ([]*models.SIPRegistration, error) {
	var registrations []*models.SIPRegistration
	query := `
		SELECT 
			id, sip_account_id, contact_uri, source_ip, source_port,
			user_agent, expires_seconds, registered_at, expired_at,
			unregistered_at, is_active
		FROM sip_registrations
		WHERE sip_account_id = $1 AND is_active = true
		ORDER BY registered_at DESC`

	err := r.db.SelectContext(ctx, &registrations, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active registrations: %w", err)
	}

	return registrations, nil
}

func (r *sipAccountRepository) RecordUsage(ctx context.Context, usage *models.SIPAccountUsage) error {
	query := `
		INSERT INTO sip_account_usage (
			sip_account_id, date, total_calls, successful_calls,
			failed_calls, total_minutes, incoming_calls, outgoing_calls,
			international_calls, average_call_duration, peak_concurrent_calls
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (sip_account_id, date) DO UPDATE SET
			total_calls = sip_account_usage.total_calls + EXCLUDED.total_calls,
			successful_calls = sip_account_usage.successful_calls + EXCLUDED.successful_calls,
			failed_calls = sip_account_usage.failed_calls + EXCLUDED.failed_calls,
			total_minutes = sip_account_usage.total_minutes + EXCLUDED.total_minutes,
			incoming_calls = sip_account_usage.incoming_calls + EXCLUDED.incoming_calls,
			outgoing_calls = sip_account_usage.outgoing_calls + EXCLUDED.outgoing_calls,
			international_calls = sip_account_usage.international_calls + EXCLUDED.international_calls,
			average_call_duration = CASE 
				WHEN (sip_account_usage.total_calls + EXCLUDED.total_calls) > 0 
				THEN ((sip_account_usage.total_minutes * 60) + (EXCLUDED.total_minutes * 60)) / (sip_account_usage.total_calls + EXCLUDED.total_calls)
				ELSE 0
			END,
			peak_concurrent_calls = GREATEST(sip_account_usage.peak_concurrent_calls, EXCLUDED.peak_concurrent_calls)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		usage.SIPAccountID,
		usage.Date,
		usage.TotalCalls,
		usage.SuccessfulCalls,
		usage.FailedCalls,
		usage.TotalMinutes,
		usage.IncomingCalls,
		usage.OutgoingCalls,
		usage.InternationalCalls,
		usage.AverageCallDuration,
		usage.PeakConcurrentCalls,
	).Scan(&usage.ID)

	if err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}

	return nil
}

func (r *sipAccountRepository) GetUsageByDate(ctx context.Context, accountID int64, date time.Time) (*models.SIPAccountUsage, error) {
	var usage models.SIPAccountUsage
	query := `
		SELECT 
			id, sip_account_id, date, total_calls, successful_calls,
			failed_calls, total_minutes, incoming_calls, outgoing_calls,
			international_calls, average_call_duration, peak_concurrent_calls,
			created_at
		FROM sip_account_usage
		WHERE sip_account_id = $1 AND date = $2`

	err := r.db.GetContext(ctx, &usage, query, accountID, date.Format("2006-01-02"))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get usage by date: %w", err)
	}

	return &usage, nil
}

func (r *sipAccountRepository) GetUsageStats(ctx context.Context, accountID int64, startDate, endDate time.Time) ([]*models.SIPAccountUsage, error) {
	var usage []*models.SIPAccountUsage
	query := `
		SELECT 
			id, sip_account_id, date, total_calls, successful_calls,
			failed_calls, total_minutes, incoming_calls, outgoing_calls,
			international_calls, average_call_duration, peak_concurrent_calls,
			created_at
		FROM sip_account_usage
		WHERE sip_account_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC`

	err := r.db.SelectContext(ctx, &usage, query, 
		accountID, 
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage stats: %w", err)
	}

	return usage, nil
}