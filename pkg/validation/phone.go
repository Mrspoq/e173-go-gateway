package validation

import (
    "fmt"
    "regexp"
    "strings"
)

// GooglePhoneValidator validates phone numbers using Google's libphonenumber logic
type GooglePhoneValidator struct {
    regionCode string
}

// PhoneNumberInfo contains parsed phone number information
type PhoneNumberInfo struct {
    Original        string `json:"original"`
    Formatted       string `json:"formatted"`
    CountryCode     string `json:"country_code"`
    NationalNumber  string `json:"national_number"`
    IsValid         bool   `json:"is_valid"`
    IsMobile        bool   `json:"is_mobile"`
    Carrier         string `json:"carrier"`
    Region          string `json:"region"`
    NumberType      string `json:"number_type"`
}

// NewGooglePhoneValidator creates a new phone validator
func NewGooglePhoneValidator(defaultRegion string) *GooglePhoneValidator {
    return &GooglePhoneValidator{
        regionCode: defaultRegion,
    }
}

// ValidatePhoneNumber validates and parses a phone number
func (g *GooglePhoneValidator) ValidatePhoneNumber(phoneNumber string) (*PhoneNumberInfo, error) {
    // Clean the phone number
    cleaned := g.cleanPhoneNumber(phoneNumber)
    
    // Parse the number
    info := &PhoneNumberInfo{
        Original: phoneNumber,
    }
    
    // Basic validation
    if !g.isValidFormat(cleaned) {
        info.IsValid = false
        return info, nil
    }
    
    // Extract country code and national number
    countryCode, nationalNumber := g.parseNumber(cleaned)
    info.CountryCode = countryCode
    info.NationalNumber = nationalNumber
    info.Formatted = g.formatNumber(countryCode, nationalNumber)
    
    // Determine region and carrier
    info.Region = g.getRegionFromCountryCode(countryCode)
    info.Carrier = g.getCarrierFromNumber(countryCode, nationalNumber)
    
    // Determine if it's mobile
    info.IsMobile = g.isMobileNumber(countryCode, nationalNumber)
    info.NumberType = g.getNumberType(countryCode, nationalNumber)
    
    // Final validation
    info.IsValid = g.validateParsedNumber(info)
    
    return info, nil
}

// cleanPhoneNumber removes formatting and normalizes the number
func (g *GooglePhoneValidator) cleanPhoneNumber(phoneNumber string) string {
    // Remove common formatting characters
    cleaned := strings.ReplaceAll(phoneNumber, " ", "")
    cleaned = strings.ReplaceAll(cleaned, "-", "")
    cleaned = strings.ReplaceAll(cleaned, "(", "")
    cleaned = strings.ReplaceAll(cleaned, ")", "")
    cleaned = strings.ReplaceAll(cleaned, ".", "")
    
    // Handle different international prefixes
    if strings.HasPrefix(cleaned, "00") {
        cleaned = "+" + cleaned[2:]
    }
    
    return cleaned
}

// isValidFormat checks basic format requirements
func (g *GooglePhoneValidator) isValidFormat(phoneNumber string) bool {
    // Must start with + for international numbers
    if !strings.HasPrefix(phoneNumber, "+") {
        return false
    }
    
    // Must contain only digits after +
    digits := phoneNumber[1:]
    matched, _ := regexp.MatchString(`^\d+$`, digits)
    if !matched {
        return false
    }
    
    // Length validation (7-15 digits for international numbers)
    if len(digits) < 7 || len(digits) > 15 {
        return false
    }
    
    return true
}

// parseNumber extracts country code and national number
func (g *GooglePhoneValidator) parseNumber(phoneNumber string) (string, string) {
    if !strings.HasPrefix(phoneNumber, "+") {
        return "", phoneNumber
    }
    
    digits := phoneNumber[1:]
    
    // Known country codes (simplified for major regions)
    countryCodes := map[string]int{
        "1":   1,  // US/Canada
        "44":  2,  // UK
        "49":  2,  // Germany
        "33":  2,  // France
        "234": 3,  // Nigeria
        "27":  2,  // South Africa
        "254": 3,  // Kenya
        "256": 3,  // Uganda
        "255": 3,  // Tanzania
        "91":  2,  // India
        "86":  2,  // China
    }
    
    // Try to match country codes
    for code, length := range countryCodes {
        if strings.HasPrefix(digits, code) && len(digits) >= len(code)+7 {
            return code, digits[length:]
        }
    }
    
    // Default: assume first 1-3 digits are country code
    if len(digits) >= 10 {
        return digits[:3], digits[3:]
    } else if len(digits) >= 8 {
        return digits[:2], digits[2:]
    }
    
    return digits[:1], digits[1:]
}

// formatNumber formats the number in international format
func (g *GooglePhoneValidator) formatNumber(countryCode, nationalNumber string) string {
    return fmt.Sprintf("+%s %s", countryCode, g.formatNationalNumber(countryCode, nationalNumber))
}

// formatNationalNumber formats the national part based on country
func (g *GooglePhoneValidator) formatNationalNumber(countryCode, nationalNumber string) string {
    switch countryCode {
    case "234": // Nigeria
        if len(nationalNumber) == 10 {
            return fmt.Sprintf("%s %s %s", nationalNumber[:3], nationalNumber[3:6], nationalNumber[6:])
        }
    case "1": // US/Canada
        if len(nationalNumber) == 10 {
            return fmt.Sprintf("(%s) %s-%s", nationalNumber[:3], nationalNumber[3:6], nationalNumber[6:])
        }
    case "44": // UK
        if len(nationalNumber) >= 10 {
            return fmt.Sprintf("%s %s %s", nationalNumber[:4], nationalNumber[4:7], nationalNumber[7:])
        }
    }
    
    // Default formatting
    if len(nationalNumber) >= 6 {
        mid := len(nationalNumber) / 2
        return fmt.Sprintf("%s %s", nationalNumber[:mid], nationalNumber[mid:])
    }
    
    return nationalNumber
}

// getRegionFromCountryCode maps country codes to regions
func (g *GooglePhoneValidator) getRegionFromCountryCode(countryCode string) string {
    regions := map[string]string{
        "1":   "US",
        "44":  "GB",
        "49":  "DE",
        "33":  "FR",
        "234": "NG",
        "27":  "ZA",
        "254": "KE",
        "256": "UG",
        "255": "TZ",
        "91":  "IN",
        "86":  "CN",
    }
    
    if region, exists := regions[countryCode]; exists {
        return region
    }
    
    return "UNKNOWN"
}

// isMobileNumber determines if the number is mobile based on patterns
func (g *GooglePhoneValidator) isMobileNumber(countryCode, nationalNumber string) bool {
    switch countryCode {
    case "234": // Nigeria mobile patterns
        if len(nationalNumber) == 10 {
            mobilePrefix := nationalNumber[:3]
            mobilePrefixes := []string{"803", "806", "813", "814", "816", "903", "906", "915", "907", "908", "809", "818", "817", "909", "908", "802", "808", "812", "701", "805", "815", "705", "707", "708", "802", "901", "904", "905", "702", "703", "704", "706"}
            for _, prefix := range mobilePrefixes {
                if mobilePrefix == prefix {
                    return true
                }
            }
        }
    case "1": // US/Canada - more complex rules needed
        return len(nationalNumber) == 10
    case "44": // UK mobile usually start with 7
        return len(nationalNumber) >= 10 && nationalNumber[0] == '7'
    }
    
    // Default assumption for unknown patterns
    return true
}

// getNumberType determines the type of phone number
func (g *GooglePhoneValidator) getNumberType(countryCode, nationalNumber string) string {
    if g.isMobileNumber(countryCode, nationalNumber) {
        return "MOBILE"
    }
    
    // Additional type detection logic
    switch countryCode {
    case "234":
        if len(nationalNumber) == 8 || len(nationalNumber) == 9 {
            return "FIXED_LINE"
        }
    case "1":
        if len(nationalNumber) == 10 {
            // Could be either mobile or fixed line - need more analysis
            return "FIXED_LINE_OR_MOBILE"
        }
    }
    
    return "UNKNOWN"
}

// getCarrierFromNumber determines carrier based on number patterns
func (g *GooglePhoneValidator) getCarrierFromNumber(countryCode, nationalNumber string) string {
    switch countryCode {
    case "234": // Nigeria
        if len(nationalNumber) >= 3 {
            prefix := nationalNumber[:3]
            carriers := map[string]string{
                "803": "MTN",
                "806": "MTN", 
                "813": "MTN",
                "814": "MTN",
                "816": "MTN",
                "903": "MTN",
                "906": "MTN",
                "802": "Airtel",
                "808": "Airtel",
                "812": "Airtel", 
                "701": "Airtel",
                "805": "Glo",
                "815": "Glo",
                "705": "Glo",
                "807": "Glo",
                "809": "9mobile",
                "818": "9mobile",
                "817": "9mobile",
                "909": "9mobile",
            }
            if carrier, exists := carriers[prefix]; exists {
                return carrier
            }
        }
    }
    
    return "UNKNOWN"
}

// validateParsedNumber performs final validation on parsed number
func (g *GooglePhoneValidator) validateParsedNumber(info *PhoneNumberInfo) bool {
    // Must have valid country code
    if info.CountryCode == "" {
        return false
    }
    
    // Must have national number
    if info.NationalNumber == "" {
        return false
    }
    
    // Region must be known for validation
    if info.Region == "UNKNOWN" {
        return false
    }
    
    // Length validation based on region
    switch info.Region {
    case "NG": // Nigeria
        return len(info.NationalNumber) >= 8 && len(info.NationalNumber) <= 10
    case "US": // United States
        return len(info.NationalNumber) == 10
    case "GB": // United Kingdom
        return len(info.NationalNumber) >= 10 && len(info.NationalNumber) <= 11
    }
    
    // Default validation
    return len(info.NationalNumber) >= 7 && len(info.NationalNumber) <= 12
}

// IsValidMobile checks if a number is a valid mobile number
func (g *GooglePhoneValidator) IsValidMobile(phoneNumber string) bool {
    info, err := g.ValidatePhoneNumber(phoneNumber)
    if err != nil {
        return false
    }
    
    return info.IsValid && info.IsMobile
}

// GetCarrier returns the carrier for a phone number
func (g *GooglePhoneValidator) GetCarrier(phoneNumber string) string {
    info, err := g.ValidatePhoneNumber(phoneNumber)
    if err != nil {
        return "UNKNOWN"
    }
    
    return info.Carrier
}
