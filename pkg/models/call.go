package models

import (
	"time"
)

// Call represents a SIP call for filtering
type Call struct {
	ID           string    `json:"id"`
	SourceNumber string    `json:"source_number"`
	DestNumber   string    `json:"dest_number"`
	GatewayID    string    `json:"gateway_id"`
	CallTime     time.Time `json:"call_time"`
}