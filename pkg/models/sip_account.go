package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// SIPAccount represents a SIP account for a customer
type SIPAccount struct {
	ID                    int64      `json:"id" db:"id"`
	CustomerID            int64      `json:"customer_id" db:"customer_id"`
	AccountName           string     `json:"account_name" db:"account_name"`
	Username              string     `json:"username" db:"username"`
	Password              string     `json:"-" db:"password"` // Never expose in JSON
	Domain                string     `json:"domain" db:"domain"`
	Extension             string     `json:"extension" db:"extension"`
	CallerID              string     `json:"caller_id" db:"caller_id"`
	CallerIDName          string     `json:"caller_id_name" db:"caller_id_name"`
	Context               string     `json:"context" db:"context"`
	Transport             string     `json:"transport" db:"transport"` // UDP, TCP, TLS
	NATSupport            bool       `json:"nat_support" db:"nat_support"`
	DirectMediaSupport    bool       `json:"direct_media_support" db:"direct_media_support"`
	EncryptionEnabled     bool       `json:"encryption_enabled" db:"encryption_enabled"`
	CodecsAllowed         string     `json:"codecs_allowed" db:"codecs_allowed"` // Comma-separated list
	MaxConcurrentCalls    int        `json:"max_concurrent_calls" db:"max_concurrent_calls"`
	CurrentActiveCalls    int        `json:"current_active_calls" db:"current_active_calls"`
	Status                string     `json:"status" db:"status"`
	LastRegisteredIP      *string    `json:"last_registered_ip" db:"last_registered_ip"`
	LastRegisteredAt      *time.Time `json:"last_registered_at" db:"last_registered_at"`
	LastCallAt            *time.Time `json:"last_call_at" db:"last_call_at"`
	TotalCalls            int64      `json:"total_calls" db:"total_calls"`
	TotalMinutes          int64      `json:"total_minutes" db:"total_minutes"`
	Notes                 *string    `json:"notes" db:"notes"`
	CreatedBy             *int64     `json:"created_by" db:"created_by"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`

	// Relations
	Customer *Customer `json:"customer,omitempty" db:"-"`
}

// SIP account status constants
const (
	SIPAccountStatusActive    = "active"
	SIPAccountStatusSuspended = "suspended"
	SIPAccountStatusDisabled  = "disabled"
	SIPAccountStatusPending   = "pending"
)

// Transport constants
const (
	TransportUDP = "UDP"
	TransportTCP = "TCP"
	TransportTLS = "TLS"
)

// Default codecs
const DefaultCodecs = "g711u,g711a,g729,g722"

// GenerateSecurePassword generates a cryptographically secure password
func GenerateSecurePassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// GetSIPURI returns the full SIP URI for this account
func (s *SIPAccount) GetSIPURI() string {
	return fmt.Sprintf("sip:%s@%s", s.Username, s.Domain)
}

// IsActive returns true if the account is active
func (s *SIPAccount) IsActive() bool {
	return s.Status == SIPAccountStatusActive
}

// CanMakeCall returns true if the account can make a new call
func (s *SIPAccount) CanMakeCall() bool {
	return s.IsActive() && s.CurrentActiveCalls < s.MaxConcurrentCalls
}

// IsRegistered returns true if the account is currently registered
func (s *SIPAccount) IsRegistered() bool {
	if s.LastRegisteredAt == nil {
		return false
	}
	// Consider registered if last registration was within 5 minutes
	return time.Since(*s.LastRegisteredAt) < 5*time.Minute
}

// SIPAccountPermission represents permissions for a SIP account
type SIPAccountPermission struct {
	ID                     int64     `json:"id" db:"id"`
	SIPAccountID           int64     `json:"sip_account_id" db:"sip_account_id"`
	AllowInternational     bool      `json:"allow_international" db:"allow_international"`
	AllowPremiumNumbers    bool      `json:"allow_premium_numbers" db:"allow_premium_numbers"`
	AllowEmergencyCalls    bool      `json:"allow_emergency_calls" db:"allow_emergency_calls"`
	AllowedCountries       *string   `json:"allowed_countries" db:"allowed_countries"` // Comma-separated country codes
	BlockedCountries       *string   `json:"blocked_countries" db:"blocked_countries"` // Comma-separated country codes
	AllowedPrefixes        *string   `json:"allowed_prefixes" db:"allowed_prefixes"`   // Comma-separated prefixes
	BlockedPrefixes        *string   `json:"blocked_prefixes" db:"blocked_prefixes"`   // Comma-separated prefixes
	TimeRestrictions       *string   `json:"time_restrictions" db:"time_restrictions"` // JSON object with time rules
	DailyCallLimit         *int      `json:"daily_call_limit" db:"daily_call_limit"`
	DailyMinuteLimit       *int      `json:"daily_minute_limit" db:"daily_minute_limit"`
	MonthlyCallLimit       *int      `json:"monthly_call_limit" db:"monthly_call_limit"`
	MonthlyMinuteLimit     *int      `json:"monthly_minute_limit" db:"monthly_minute_limit"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// SIPRegistration represents a SIP registration event
type SIPRegistration struct {
	ID               int64     `json:"id" db:"id"`
	SIPAccountID     int64     `json:"sip_account_id" db:"sip_account_id"`
	ContactURI       string    `json:"contact_uri" db:"contact_uri"`
	SourceIP         string    `json:"source_ip" db:"source_ip"`
	SourcePort       int       `json:"source_port" db:"source_port"`
	UserAgent        *string   `json:"user_agent" db:"user_agent"`
	ExpiresSeconds   int       `json:"expires_seconds" db:"expires_seconds"`
	RegisteredAt     time.Time `json:"registered_at" db:"registered_at"`
	ExpiredAt        time.Time `json:"expired_at" db:"expired_at"`
	UnregisteredAt   *time.Time `json:"unregistered_at" db:"unregistered_at"`
	IsActive         bool      `json:"is_active" db:"is_active"`
}

// SIPAccountUsage represents usage statistics for a SIP account
type SIPAccountUsage struct {
	ID                    int64     `json:"id" db:"id"`
	SIPAccountID          int64     `json:"sip_account_id" db:"sip_account_id"`
	Date                  time.Time `json:"date" db:"date"`
	TotalCalls            int       `json:"total_calls" db:"total_calls"`
	SuccessfulCalls       int       `json:"successful_calls" db:"successful_calls"`
	FailedCalls           int       `json:"failed_calls" db:"failed_calls"`
	TotalMinutes          int       `json:"total_minutes" db:"total_minutes"`
	IncomingCalls         int       `json:"incoming_calls" db:"incoming_calls"`
	OutgoingCalls         int       `json:"outgoing_calls" db:"outgoing_calls"`
	InternationalCalls    int       `json:"international_calls" db:"international_calls"`
	AverageCallDuration   int       `json:"average_call_duration" db:"average_call_duration"`
	PeakConcurrentCalls   int       `json:"peak_concurrent_calls" db:"peak_concurrent_calls"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}