package models

import "time"

// ValidationResult contains WhatsApp validation results
type ValidationResult struct {
    PhoneNumber       string    `json:"phone_number"`
    HasWhatsApp       bool      `json:"has_whatsapp"`
    IsBusinessAccount bool      `json:"is_business_account"`
    ProfileName       string    `json:"profile_name"`
    Confidence        float64   `json:"confidence"`
    LastUpdated       time.Time `json:"last_updated"`
    Source            string    `json:"source"`
}

// PrivateWhatsAppResponse represents the private API response structure
type PrivateWhatsAppResponse struct {
    Status   bool   `json:"status"`    // true if API responded
    Valid    bool   `json:"valid"`     // true if number has WhatsApp
    WaID     string `json:"wa_id"`     // WhatsApp ID (empty if invalid)
    ChatLink string `json:"chat_link"` // WhatsApp chat link (empty if invalid)
}