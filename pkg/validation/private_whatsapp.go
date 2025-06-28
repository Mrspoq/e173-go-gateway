package validation

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "regexp"
    "time"
    
    "github.com/e173-gateway/e173_go_gateway/pkg/models"
)

// PrivateWhatsAppValidator uses your private WhatsApp validation API
type PrivateWhatsAppValidator struct {
    apiKey     string
    baseURL    string
    client     *http.Client
    cache      *ValidationCache
}

// NewPrivateWhatsAppValidator creates validator for your private API
func NewPrivateWhatsAppValidator(apiKey string) *PrivateWhatsAppValidator {
    return &PrivateWhatsAppValidator{
        apiKey:  apiKey,
        baseURL: "https://bulkvalidation.wa-validator.xyz/v2/validate/wa_id",
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        cache: &ValidationCache{
            results: make(map[string]*models.ValidationResult),
            expiry:  24 * time.Hour, // Cache for 24 hours
        },
    }
}

// ValidateNumber checks if a phone number has WhatsApp using your private API
func (w *PrivateWhatsAppValidator) ValidateNumber(phoneNumber string) (*models.ValidationResult, error) {
    // Check cache first
    if cached := w.cache.Get(phoneNumber); cached != nil {
        return cached, nil
    }

    // Make API request to your private API
    result, err := w.makePrivateAPIRequest(phoneNumber)
    if err != nil {
        return nil, fmt.Errorf("Private WhatsApp API request failed: %w", err)
    }

    // Cache the result
    w.cache.Set(phoneNumber, result)

    return result, nil
}

// makePrivateAPIRequest calls your specific API endpoint
func (w *PrivateWhatsAppValidator) makePrivateAPIRequest(phoneNumber string) (*models.ValidationResult, error) {
    // Clean phone number (remove + and spaces)
    cleanNumber := w.cleanPhoneNumber(phoneNumber)
    
    // Build URL with your API format
    url := fmt.Sprintf("%s?number=%s", w.baseURL, cleanNumber)
    
    // Create HTTP request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set your authorization header
    req.Header.Set("Authorization", "Bearer "+w.apiKey)
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

    // Parse your API response format
    var apiResp models.PrivateWhatsAppResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    // Convert to our standard format
    result := &models.ValidationResult{
        PhoneNumber:       phoneNumber,
        HasWhatsApp:       apiResp.Valid && apiResp.Status,
        IsBusinessAccount: false, // Your API doesn't provide this
        ProfileName:       "",    // Your API doesn't provide this
        LastUpdated:       time.Now(),
        Source:            "private_whatsapp_api",
    }

    // Calculate confidence based on your API response
    if apiResp.Status && apiResp.Valid && apiResp.WaID != "" {
        result.Confidence = 0.95 // High confidence when API says valid
    } else if apiResp.Status && !apiResp.Valid {
        result.Confidence = 0.90 // High confidence when API says invalid
    } else {
        result.Confidence = 0.50 // Lower confidence if API didn't respond properly
    }

    return result, nil
}

// cleanPhoneNumber removes formatting for your API
func (w *PrivateWhatsAppValidator) cleanPhoneNumber(phoneNumber string) string {
    // Your API expects: 34674944456 (no + sign)
    cleaned := phoneNumber
    
    // Remove + sign
    if cleaned[0] == '+' {
        cleaned = cleaned[1:]
    }
    
    // Remove any spaces, dashes, parentheses
    cleaned = regexp.MustCompile(`[^\d]`).ReplaceAllString(cleaned, "")
    
    return cleaned
}

// IsLikelyRealPerson uses your API to determine if number belongs to real person
func (w *PrivateWhatsAppValidator) IsLikelyRealPerson(phoneNumber string) (bool, float64, error) {
    result, err := w.ValidateNumber(phoneNumber)
    if err != nil {
        return false, 0, err
    }

    // In Africa/Middle East, 90%+ of real people have WhatsApp
    // So having WhatsApp = likely real person
    // Not having WhatsApp = suspicious (could be automated system)
    
    confidence := result.Confidence
    isReal := result.HasWhatsApp
    
    return isReal, confidence, nil
}

// BatchValidate validates multiple numbers efficiently
func (w *PrivateWhatsAppValidator) BatchValidate(phoneNumbers []string) (map[string]*models.ValidationResult, error) {
    results := make(map[string]*models.ValidationResult)
    
    // Process in parallel with goroutines (rate limited)
    semaphore := make(chan struct{}, 5) // Max 5 concurrent requests
    resultChan := make(chan struct {
        number string
        result *models.ValidationResult
        err    error
    }, len(phoneNumbers))
    
    // Start validation for each number
    for _, number := range phoneNumbers {
        go func(num string) {
            semaphore <- struct{}{} // Acquire semaphore
            defer func() { <-semaphore }() // Release semaphore
            
            result, err := w.ValidateNumber(num)
            resultChan <- struct {
                number string
                result *models.ValidationResult
                err    error
            }{num, result, err}
        }(number)
    }
    
    // Collect results
    for i := 0; i < len(phoneNumbers); i++ {
        result := <-resultChan
        if result.err == nil {
            results[result.number] = result.result
        }
    }
    
    return results, nil
}

// GetCacheStats returns statistics about the validation cache
func (w *PrivateWhatsAppValidator) GetCacheStats() map[string]interface{} {
    return w.cache.GetStats()
}

// ClearCache removes all cached validation results
func (w *PrivateWhatsAppValidator) ClearCache() {
    w.cache.results = make(map[string]*models.ValidationResult)
}
