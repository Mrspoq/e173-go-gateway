package models

import (
	"time"
)

// Modem represents the structure of our 'modems' table.
type Modem struct {
	ID                          int        `json:"id" db:"id"`
	DevicePath                  string     `json:"device_path" db:"device_path"`
	IMEI                        *string    `json:"imei,omitempty" db:"imei"` // Use pointer for nullable fields
	IMSI                        *string    `json:"imsi,omitempty" db:"imsi"`
	Model                       *string    `json:"model,omitempty" db:"model"`
	Manufacturer                *string    `json:"manufacturer,omitempty" db:"manufacturer"`
	FirmwareVersion             *string    `json:"firmware_version,omitempty" db:"firmware_version"`
	SignalStrengthDBM           *int       `json:"signal_strength_dbm,omitempty" db:"signal_strength_dbm"`
	NetworkOperatorName         *string    `json:"network_operator_name,omitempty" db:"network_operator_name"`
	NetworkRegistrationStatus   *string    `json:"network_registration_status,omitempty" db:"network_registration_status"`
	Status                      string     `json:"status" db:"status"`
	LastSeenAt                  *time.Time `json:"last_seen_at,omitempty" db:"last_seen_at"`
	CreatedAt                   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at" db:"updated_at"`
}
