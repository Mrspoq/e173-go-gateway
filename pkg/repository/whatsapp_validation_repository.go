package repository

import (
    "context"
    "database/sql"
    "encoding/json"
    "time"
    
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/e173-gateway/e173_go_gateway/pkg/models"
)

// WhatsAppValidationRepository handles database operations for WhatsApp validation cache
type WhatsAppValidationRepository interface {
    GetValidation(ctx context.Context, phoneNumber string) (*models.ValidationResult, error)
    SaveValidation(ctx context.Context, result *models.ValidationResult) error
    CleanupExpired(ctx context.Context) error
    GetStats(ctx context.Context) (map[string]interface{}, error)
}

// PostgresWhatsAppValidationRepository implements WhatsAppValidationRepository using PostgreSQL
type PostgresWhatsAppValidationRepository struct {
    db *pgxpool.Pool
}

// NewPostgresWhatsAppValidationRepository creates a new PostgreSQL-based repository
func NewPostgresWhatsAppValidationRepository(db *pgxpool.Pool) WhatsAppValidationRepository {
    return &PostgresWhatsAppValidationRepository{db: db}
}

// GetValidation retrieves a cached validation result
func (r *PostgresWhatsAppValidationRepository) GetValidation(ctx context.Context, phoneNumber string) (*models.ValidationResult, error) {
    query := `
        SELECT has_whatsapp, profile_name, is_business, confidence, checked_at, expires_at, raw_response
        FROM whatsapp_validation_cache
        WHERE phone_number = $1 AND expires_at > NOW()
        LIMIT 1
    `
    
    var result models.ValidationResult
    var rawResponse sql.NullString
    var profileName sql.NullString
    var expiresAt time.Time
    
    err := r.db.QueryRow(ctx, query, phoneNumber).Scan(
        &result.HasWhatsApp,
        &profileName,
        &result.IsBusinessAccount,
        &result.Confidence,
        &result.LastUpdated,
        &expiresAt,
        &rawResponse,
    )
    
    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, nil // No cached result
        }
        return nil, err
    }
    
    result.PhoneNumber = phoneNumber
    result.Source = "database_cache"
    if profileName.Valid {
        result.ProfileName = profileName.String
    }
    
    return &result, nil
}

// SaveValidation stores a validation result in the cache
func (r *PostgresWhatsAppValidationRepository) SaveValidation(ctx context.Context, result *models.ValidationResult) error {
    // Convert result to JSON for raw_response field
    rawJSON, err := json.Marshal(result)
    if err != nil {
        return err
    }
    
    query := `
        INSERT INTO whatsapp_validation_cache (
            phone_number, has_whatsapp, profile_name, is_business, 
            confidence, checked_at, expires_at, raw_response
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (phone_number) DO UPDATE SET
            has_whatsapp = EXCLUDED.has_whatsapp,
            profile_name = EXCLUDED.profile_name,
            is_business = EXCLUDED.is_business,
            confidence = EXCLUDED.confidence,
            checked_at = EXCLUDED.checked_at,
            expires_at = EXCLUDED.expires_at,
            raw_response = EXCLUDED.raw_response
    `
    
    profileName := sql.NullString{String: result.ProfileName, Valid: result.ProfileName != ""}
    expiresAt := result.LastUpdated.Add(24 * time.Hour) // 24 hour cache
    
    _, err = r.db.Exec(ctx, query,
        result.PhoneNumber,
        result.HasWhatsApp,
        profileName,
        result.IsBusinessAccount,
        result.Confidence,
        result.LastUpdated,
        expiresAt,
        string(rawJSON),
    )
    
    return err
}

// CleanupExpired removes expired cache entries
func (r *PostgresWhatsAppValidationRepository) CleanupExpired(ctx context.Context) error {
    query := `DELETE FROM whatsapp_validation_cache WHERE expires_at < NOW()`
    _, err := r.db.Exec(ctx, query)
    return err
}

// GetStats returns cache statistics
func (r *PostgresWhatsAppValidationRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Total entries
    var totalEntries int64
    err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM whatsapp_validation_cache").Scan(&totalEntries)
    if err != nil {
        return nil, err
    }
    stats["total_entries"] = totalEntries
    
    // Valid entries (not expired)
    var validEntries int64
    err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM whatsapp_validation_cache WHERE expires_at > NOW()").Scan(&validEntries)
    if err != nil {
        return nil, err
    }
    stats["valid_entries"] = validEntries
    
    // WhatsApp users
    var whatsappUsers int64
    err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM whatsapp_validation_cache WHERE has_whatsapp = true AND expires_at > NOW()").Scan(&whatsappUsers)
    if err != nil {
        return nil, err
    }
    stats["whatsapp_users"] = whatsappUsers
    
    // Non-WhatsApp users  
    var nonWhatsappUsers int64
    err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM whatsapp_validation_cache WHERE has_whatsapp = false AND expires_at > NOW()").Scan(&nonWhatsappUsers)
    if err != nil {
        return nil, err
    }
    stats["non_whatsapp_users"] = nonWhatsappUsers
    
    // Business accounts
    var businessAccounts int64
    err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM whatsapp_validation_cache WHERE is_business = true AND expires_at > NOW()").Scan(&businessAccounts)
    if err != nil {
        return nil, err
    }
    stats["business_accounts"] = businessAccounts
    
    // Cache hit rate (if we tracked hits/misses)
    stats["cache_expiry"] = "24 hours"
    
    return stats, nil
}