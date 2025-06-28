package repository

import (
    "context"
    "time"
    
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/e173-gateway/e173_go_gateway/pkg/models"
)

// SimpleWhatsAppValidationRepository implements caching using the existing schema
type SimpleWhatsAppValidationRepository struct {
    db *pgxpool.Pool
}

// NewSimpleWhatsAppValidationRepository creates a new repository for the existing schema
func NewSimpleWhatsAppValidationRepository(db *pgxpool.Pool) WhatsAppValidationRepository {
    return &SimpleWhatsAppValidationRepository{db: db}
}

// GetValidation retrieves a cached validation result from the existing table
func (r *SimpleWhatsAppValidationRepository) GetValidation(ctx context.Context, phoneNumber string) (*models.ValidationResult, error) {
    query := `
        SELECT has_whatsapp, confidence_score, checked_at
        FROM whatsapp_validation_cache
        WHERE phone_number = $1 AND expires_at > NOW()
        LIMIT 1
    `
    
    var result models.ValidationResult
    var confidence *float64
    
    err := r.db.QueryRow(ctx, query, phoneNumber).Scan(
        &result.HasWhatsApp,
        &confidence,
        &result.LastUpdated,
    )
    
    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, nil // No cached result
        }
        return nil, err
    }
    
    result.PhoneNumber = phoneNumber
    result.Source = "database_cache"
    
    // Handle nullable confidence score
    if confidence != nil {
        result.Confidence = *confidence
    } else {
        result.Confidence = 0.5 // Default confidence
    }
    
    // These fields don't exist in the current schema
    result.IsBusinessAccount = false
    result.ProfileName = ""
    
    return &result, nil
}

// SaveValidation stores a validation result in the existing cache table
func (r *SimpleWhatsAppValidationRepository) SaveValidation(ctx context.Context, result *models.ValidationResult) error {
    query := `
        INSERT INTO whatsapp_validation_cache (
            phone_number, has_whatsapp, confidence_score, checked_at, expires_at
        ) VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (phone_number) DO UPDATE SET
            has_whatsapp = EXCLUDED.has_whatsapp,
            confidence_score = EXCLUDED.confidence_score,
            checked_at = EXCLUDED.checked_at,
            expires_at = EXCLUDED.expires_at
    `
    
    expiresAt := result.LastUpdated.Add(24 * time.Hour) // 24 hour cache
    
    _, err := r.db.Exec(ctx, query,
        result.PhoneNumber,
        result.HasWhatsApp,
        result.Confidence,
        result.LastUpdated,
        expiresAt,
    )
    
    return err
}

// CleanupExpired removes expired cache entries
func (r *SimpleWhatsAppValidationRepository) CleanupExpired(ctx context.Context) error {
    query := `DELETE FROM whatsapp_validation_cache WHERE expires_at < NOW()`
    _, err := r.db.Exec(ctx, query)
    return err
}

// GetStats returns cache statistics for the existing schema
func (r *SimpleWhatsAppValidationRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
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
    
    // Average confidence score
    var avgConfidence *float64
    err = r.db.QueryRow(ctx, "SELECT AVG(confidence_score) FROM whatsapp_validation_cache WHERE expires_at > NOW()").Scan(&avgConfidence)
    if err == nil && avgConfidence != nil {
        stats["average_confidence"] = *avgConfidence
    }
    
    stats["cache_expiry"] = "24 hours"
    
    return stats, nil
}