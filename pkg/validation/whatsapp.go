package validation

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    
    "github.com/e173-gateway/e173_go_gateway/pkg/models"
)

// WhatsAppValidator handles WhatsApp Business API validation
type WhatsAppValidator struct {
    apiKey    string
    baseURL   string
    client    *http.Client
    cache     *ValidationCache
}

// WhatsAppAPIResponse represents the API response structure
type WhatsAppAPIResponse struct {
    Valid       bool   `json:"valid"`
    WhatsApp    bool   `json:"whatsapp"`
    Business    bool   `json:"business"`
    Name        string `json:"name"`
    CountryCode string `json:"country_code"`
    Carrier     string `json:"carrier"`
}

// ValidationCache handles caching of validation results
type ValidationCache struct {
    results map[string]*models.ValidationResult
    expiry  time.Duration
}

// NewWhatsAppValidator creates a new WhatsApp validator instance
func NewWhatsAppValidator(apiKey string) *WhatsAppValidator {
    return &WhatsAppValidator{
        apiKey:  apiKey,
        baseURL: "https://api.whatsapp.com/v1/contacts", // Official WhatsApp Business API
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        cache: &ValidationCache{
            results: make(map[string]*models.ValidationResult),
            expiry:  24 * time.Hour, // Cache for 24 hours
        },
    }
}

// ValidateNumber checks if a phone number has WhatsApp
func (w *WhatsAppValidator) ValidateNumber(phoneNumber string) (*models.ValidationResult, error) {
    // Check cache first
    if cached := w.cache.Get(phoneNumber); cached != nil {
        return cached, nil
    }

    // Make API request
    result, err := w.makeAPIRequest(phoneNumber)
    if err != nil {
        return nil, fmt.Errorf("WhatsApp API request failed: %w", err)
    }

    // Cache the result
    w.cache.Set(phoneNumber, result)

    return result, nil
}

// IsLikelyRealPerson determines if a number belongs to a real person
func (w *WhatsAppValidator) IsLikelyRealPerson(phoneNumber string) (bool, float64, error) {
    result, err := w.ValidateNumber(phoneNumber)
    if err != nil {
        return false, 0, err
    }

    // Scoring logic for real person detection
    confidence := 0.0

    if result.HasWhatsApp {
        confidence += 0.7 // Having WhatsApp is strong indicator
    }

    if result.ProfileName != "" {
        confidence += 0.2 // Having a profile name helps
    }

    if result.IsBusinessAccount {
        confidence -= 0.1 // Business accounts might be automated
    }

    // In Africa, 90%+ of smartphone users have WhatsApp
    // So NOT having WhatsApp is suspicious for a real person
    if !result.HasWhatsApp {
        confidence = 0.1 // Very low confidence for real person
    }

    return confidence >= 0.6, confidence, nil
}

// makeAPIRequest performs the actual API call to WhatsApp
func (w *WhatsAppValidator) makeAPIRequest(phoneNumber string) (*models.ValidationResult, error) {
    // Prepare request payload
    payload := map[string]interface{}{
        "contacts": []map[string]string{
            {"input": phoneNumber},
        },
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", w.baseURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Authorization", "Bearer "+w.apiKey)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "E173-Gateway/1.0")

    // Make request
    resp, err := w.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Handle different response codes
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
    }

    // Parse response
    var apiResp WhatsAppAPIResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    // Convert to our format
    result := &models.ValidationResult{
        PhoneNumber:       phoneNumber,
        HasWhatsApp:       apiResp.WhatsApp,
        IsBusinessAccount: apiResp.Business,
        ProfileName:       apiResp.Name,
        LastUpdated:       time.Now(),
        Source:            "whatsapp_api",
    }

    // Calculate confidence based on response
    result.Confidence = w.calculateConfidence(apiResp)

    return result, nil
}

// calculateConfidence determines confidence score based on API response
func (w *WhatsAppValidator) calculateConfidence(resp WhatsAppAPIResponse) float64 {
    confidence := 0.5 // Base confidence

    if resp.Valid {
        confidence += 0.3
    }

    if resp.WhatsApp {
        confidence += 0.4
    }

    if resp.Name != "" {
        confidence += 0.1
    }

    if confidence > 1.0 {
        confidence = 1.0
    }

    return confidence
}

// Cache methods

// Get retrieves a cached validation result
func (c *ValidationCache) Get(phoneNumber string) *models.ValidationResult {
    result, exists := c.results[phoneNumber]
    if !exists {
        return nil
    }

    // Check if expired
    if time.Since(result.LastUpdated) > c.expiry {
        delete(c.results, phoneNumber)
        return nil
    }

    return result
}

// Set stores a validation result in cache
func (c *ValidationCache) Set(phoneNumber string, result *models.ValidationResult) {
    c.results[phoneNumber] = result
}

// Cleanup removes expired entries from cache
func (c *ValidationCache) Cleanup() {
    now := time.Now()
    for phone, result := range c.results {
        if now.Sub(result.LastUpdated) > c.expiry {
            delete(c.results, phone)
        }
    }
}

// GetStats returns cache statistics
func (c *ValidationCache) GetStats() map[string]interface{} {
    return map[string]interface{}{
        "total_entries": len(c.results),
        "cache_expiry":  c.expiry.String(),
    }
}

// Alternative API implementations for different providers

// NumVerifyValidator uses NumVerify API as backup
type NumVerifyValidator struct {
    apiKey string
    client *http.Client
}

// NewNumVerifyValidator creates NumVerify validator (backup option)
func NewNumVerifyValidator(apiKey string) *NumVerifyValidator {
    return &NumVerifyValidator{
        apiKey: apiKey,
        client: &http.Client{Timeout: 5 * time.Second},
    }
}

// ValidateNumber validates using NumVerify API
func (n *NumVerifyValidator) ValidateNumber(phoneNumber string) (*models.ValidationResult, error) {
    url := fmt.Sprintf("http://apilayer.net/api/validate?access_key=%s&number=%s", 
        n.apiKey, phoneNumber)

    resp, err := n.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result struct {
        Valid       bool   `json:"valid"`
        CountryCode string `json:"country_code"`
        Carrier     string `json:"carrier"`
        LineType    string `json:"line_type"`
    }

    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    confidence := 0.3
    if result.Valid {
        confidence = 0.7
    }

    return &models.ValidationResult{
        PhoneNumber: phoneNumber,
        HasWhatsApp: result.Valid && result.LineType == "mobile",
        Confidence:  confidence,
        LastUpdated: time.Now(),
        Source:      "numverify",
    }, nil
}
