package models

import (
	"time"
)

// BlacklistEntry represents a blacklisted phone number
type BlacklistEntry struct {
	ID          int       `json:"id" db:"id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Reason      string    `json:"reason" db:"reason"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}