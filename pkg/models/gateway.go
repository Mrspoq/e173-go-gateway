package models

import (
	"time"
)

// Gateway represents a remote E173 gateway with Asterisk + dongles
type Gateway struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Location    string     `json:"location" db:"location"`
	AMIHost     string     `json:"ami_host" db:"ami_host"`
	AMIPort     string     `json:"ami_port" db:"ami_port"`
	AMIUser     string     `json:"ami_user" db:"ami_user"`
	AMIPass     string     `json:"-" db:"ami_pass"` // Don't expose password in JSON
	Status      string     `json:"status" db:"status"`   // online, offline, error
	Enabled     bool       `json:"enabled" db:"enabled"` // Can be disabled by admin
	LastSeen    *time.Time `json:"last_seen" db:"last_seen"`
	LastError   *string    `json:"last_error" db:"last_error"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	
	// Runtime stats (not stored in DB)
	ActiveCalls    int     `json:"active_calls" db:"-"`
	OnlineModems   int     `json:"online_modems" db:"-"`
	TotalModems    int     `json:"total_modems" db:"-"`
	UptimePercent  float64 `json:"uptime_percent" db:"-"`
}

// Gateway statuses
const (
	GatewayStatusOnline  = "online"
	GatewayStatusOffline = "offline"
	GatewayStatusError   = "error"
)
