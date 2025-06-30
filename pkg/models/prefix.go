package models

import (
	"time"
)

// Prefix represents a phone number prefix with routing information
type Prefix struct {
	ID            string    `json:"id" db:"id"`
	Prefix        string    `json:"prefix" db:"prefix"`
	Country       string    `json:"country" db:"country"`
	Operator      string    `json:"operator" db:"operator"`
	GatewayID     string    `json:"gateway_id" db:"gateway_id"`
	RatePerMinute float64   `json:"rate_per_minute" db:"rate_per_minute"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// PrefixRoute represents a routing rule for a prefix
type PrefixRoute struct {
	ID              string    `json:"id" db:"id"`
	PrefixID        string    `json:"prefix_id" db:"prefix_id"`
	GatewayID       string    `json:"gateway_id" db:"gateway_id"`
	Priority        int       `json:"priority" db:"priority"`
	Weight          int       `json:"weight" db:"weight"`
	MaxConcurrent   int       `json:"max_concurrent" db:"max_concurrent"`
	CurrentActive   int       `json:"current_active" db:"current_active"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}