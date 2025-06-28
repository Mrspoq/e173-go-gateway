package spam

import (
    "regexp"
    "strconv"
    "strings"
    "time"
)

// SpamPatternDetector analyzes call patterns for spam detection
type SpamPatternDetector struct {
    database *CallPatternDB
}

// CallPattern represents calling behavior analysis
type CallPattern struct {
    PhoneNumber       string    `json:"phone_number"`
    TotalCalls        int       `json:"total_calls"`
    CallsLast24H      int       `json:"calls_last_24h"`
    CallsLastHour     int       `json:"calls_last_hour"`
    AverageCallLength float64   `json:"avg_call_length"`
    UniqueDestinations int      `json:"unique_destinations"`
    SequentialPattern bool      `json:"sequential_pattern"`
    SpamScore         float64   `json:"spam_score"`
    LastCallTime      time.Time `json:"last_call_time"`
}

// SpamDetectionResult contains the analysis result
type SpamDetectionResult struct {
    IsSpam       bool    `json:"is_spam"`
    Confidence   float64 `json:"confidence"`
    Reasons      []string `json:"reasons"`
    SpamScore    float64 `json:"spam_score"`
    Action       string  `json:"action"` // "block", "route_to_ai", "allow"
}

// NewSpamPatternDetector creates a new spam detector
func NewSpamPatternDetector(db *CallPatternDB) *SpamPatternDetector {
    return &SpamPatternDetector{
        database: db,
    }
}

// AnalyzeNumber performs comprehensive spam analysis on a phone number
func (s *SpamPatternDetector) AnalyzeNumber(phoneNumber string) (*SpamDetectionResult, error) {
    pattern, err := s.database.GetCallPattern(phoneNumber)
    if err != nil {
        return nil, err
    }

    result := &SpamDetectionResult{
        Reasons: make([]string, 0),
    }

    // Rule 1: Sequential number pattern detection
    if s.isSequentialNumber(phoneNumber) {
        result.SpamScore += 0.4
        result.Reasons = append(result.Reasons, "Sequential number pattern detected")
    }

    // Rule 2: High frequency calling (multiple calls in short time)
    if pattern.CallsLastHour > 10 {
        result.SpamScore += 0.5
        result.Reasons = append(result.Reasons, "High frequency calling pattern")
    }

    // Rule 3: Short call duration pattern
    if pattern.AverageCallLength < 10 && pattern.TotalCalls > 5 {
        result.SpamScore += 0.3
        result.Reasons = append(result.Reasons, "Consistently short call durations")
    }

    // Rule 4: Multiple destinations from same source
    if pattern.UniqueDestinations > 20 && pattern.CallsLast24H > 50 {
        result.SpamScore += 0.6
        result.Reasons = append(result.Reasons, "Calling multiple destinations rapidly")
    }

    // Rule 5: Time pattern analysis (calls outside normal hours)
    if s.isOutsideBusinessHours(pattern.LastCallTime) && pattern.CallsLast24H > 10 {
        result.SpamScore += 0.2
        result.Reasons = append(result.Reasons, "Calling outside business hours")
    }

    // Rule 6: Repetitive calling to same destinations
    if s.hasRepetitivePattern(phoneNumber) {
        result.SpamScore += 0.3
        result.Reasons = append(result.Reasons, "Repetitive calling pattern")
    }

    // Calculate final confidence and decision
    result.Confidence = result.SpamScore
    if result.SpamScore > 0.8 {
        result.IsSpam = true
        result.Action = "block"
    } else if result.SpamScore > 0.5 {
        result.IsSpam = true
        result.Action = "route_to_ai" // Monetize suspicious calls
    } else {
        result.IsSpam = false
        result.Action = "allow"
    }

    return result, nil
}

// isSequentialNumber detects if phone number follows sequential patterns
func (s *SpamPatternDetector) isSequentialNumber(phoneNumber string) bool {
    // Remove country code and formatting
    digits := s.extractDigits(phoneNumber)
    if len(digits) < 8 {
        return false
    }

    // Check for sequential patterns in last 4-6 digits
    lastDigits := digits[len(digits)-6:]
    
    // Pattern 1: Consecutive numbers (1234, 5678)
    consecutive := 0
    for i := 1; i < len(lastDigits); i++ {
        if lastDigits[i] == lastDigits[i-1]+1 {
            consecutive++
        }
    }
    if consecutive >= 3 {
        return true
    }

    // Pattern 2: Repeating digits (1111, 2222)
    repeating := 0
    for i := 1; i < len(lastDigits); i++ {
        if lastDigits[i] == lastDigits[i-1] {
            repeating++
        }
    }
    if repeating >= 4 {
        return true
    }

    // Pattern 3: Simple patterns (1212, 3434)
    if s.hasSimplePattern(lastDigits) {
        return true
    }

    return false
}

// hasSimplePattern detects simple repeating patterns
func (s *SpamPatternDetector) hasSimplePattern(digits []int) bool {
    if len(digits) < 4 {
        return false
    }

    // Check for ABAB pattern
    if digits[0] == digits[2] && digits[1] == digits[3] {
        return true
    }

    // Check for AABB pattern
    if digits[0] == digits[1] && digits[2] == digits[3] {
        return true
    }

    return false
}

// extractDigits converts phone number to digit array
func (s *SpamPatternDetector) extractDigits(phoneNumber string) []int {
    re := regexp.MustCompile(`\d`)
    matches := re.FindAllString(phoneNumber, -1)
    
    digits := make([]int, len(matches))
    for i, match := range matches {
        digits[i], _ = strconv.Atoi(match)
    }
    
    return digits
}

// isOutsideBusinessHours checks if call time is suspicious
func (s *SpamPatternDetector) isOutsideBusinessHours(callTime time.Time) bool {
    hour := callTime.Hour()
    // Business hours: 8 AM to 8 PM
    return hour < 8 || hour > 20
}

// hasRepetitivePattern checks for repetitive calling behavior
func (s *SpamPatternDetector) hasRepetitivePattern(phoneNumber string) bool {
    // Get call history for this number
    history, err := s.database.GetCallHistory(phoneNumber, 24) // Last 24 hours
    if err != nil {
        return false
    }

    if len(history) < 5 {
        return false
    }

    // Check for calls to same destination multiple times
    destinationCounts := make(map[string]int)
    for _, call := range history {
        destinationCounts[call.Destination]++
    }

    // If calling same number more than 3 times in 24h, suspicious
    for _, count := range destinationCounts {
        if count > 3 {
            return true
        }
    }

    return false
}

// UpdatePatternFromCall updates the pattern database with new call data
func (s *SpamPatternDetector) UpdatePatternFromCall(phoneNumber, destination string, duration int) error {
    return s.database.UpdateCallPattern(phoneNumber, destination, duration)
}

// GetFilterLevel returns recommended filter level based on patterns
func (s *SpamPatternDetector) GetFilterLevel(phoneNumber string) string {
    result, err := s.AnalyzeNumber(phoneNumber)
    if err != nil {
        return "medium" // Default
    }

    if result.SpamScore > 0.8 {
        return "maximum"
    } else if result.SpamScore > 0.4 {
        return "medium"
    }
    
    return "basic"
}

// Placeholder for database interface
type CallPatternDB struct {
    // TODO: Implement with actual database connection
}

type CallRecord struct {
    Destination string
    Duration    int
    Timestamp   time.Time
}

func (db *CallPatternDB) GetCallPattern(phoneNumber string) (*CallPattern, error) {
    // TODO: Implement database query
    return &CallPattern{
        PhoneNumber: phoneNumber,
        TotalCalls:  0,
        CallsLast24H: 0,
        CallsLastHour: 0,
    }, nil
}

func (db *CallPatternDB) GetCallHistory(phoneNumber string, hours int) ([]*CallRecord, error) {
    // TODO: Implement database query
    return []*CallRecord{}, nil
}

func (db *CallPatternDB) UpdateCallPattern(phoneNumber, destination string, duration int) error {
    // TODO: Implement database update
    return nil
}
