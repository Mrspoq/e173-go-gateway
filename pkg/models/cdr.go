package models

import (
	"time"

	"github.com/google/uuid" // For the ID field
)

// Cdr represents a Call Detail Record in the system.
// This struct aligns with the schema in 000003_create_call_detail_records.up.sql
type Cdr struct {
	ID                   uuid.UUID  `json:"id"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Channel              *string    `json:"channel,omitempty"`
	UniqueID             string     `json:"unique_id"` // Asterisk's UniqueID
	CallerIDNum          *string    `json:"caller_id_num,omitempty"`
	CallerIDName         *string    `json:"caller_id_name,omitempty"`
	ConnectedLineNum     *string    `json:"connected_line_num,omitempty"`
	ConnectedLineName    *string    `json:"connected_line_name,omitempty"`
	AccountCode          *string    `json:"account_code,omitempty"`
	Cause                *string    `json:"cause,omitempty"`      // Hangup cause code
	CauseTxt             *string    `json:"cause_txt,omitempty"`  // Hangup cause text
	StartTime            *time.Time `json:"start_time,omitempty"` // Dial start time
	AnswerTime           *time.Time `json:"answer_time,omitempty"`
	EndTime              *time.Time `json:"end_time,omitempty"`
	Duration             *int       `json:"duration,omitempty"`          // Total call duration in seconds
	BillableSeconds      *int       `json:"billable_seconds,omitempty"` // Billable duration in seconds
	ModemID              *int       `json:"modem_id,omitempty"`
	SimCardID            *int       `json:"sim_card_id,omitempty"`
	CallDirection        *string    `json:"call_direction,omitempty"` // "inbound", "outbound"
	SipCustomerID        *int       `json:"sip_customer_id,omitempty"`
	Cost                 *float64   `json:"cost,omitempty"`
	CustomerPrice        *float64   `json:"customer_price,omitempty"`
	IsSpam               *bool      `json:"is_spam,omitempty"`
	SpamDetectionMethod  *string    `json:"spam_detection_method,omitempty"`
	Context              *string    `json:"context,omitempty"`
	Extension            *string    `json:"extension,omitempty"`
	Priority             *int       `json:"priority,omitempty"`
	RawEventData         []byte     `json:"raw_event_data,omitempty"` // Storing as JSONB, so []byte for raw JSON
	Disposition          *string    `json:"disposition,omitempty"`    // Added based on AMIService logic
}

// Constants for CallDirection (can be moved or kept here)
const (
	CallDirectionInbound  = "inbound"
	CallDirectionOutbound = "outbound"
	CallDirectionUnknown  = "unknown"
)

// Constants for CallDisposition
const (
	CallDispositionAnswered = "ANSWERED"
	CallDispositionNoAnswer = "NO ANSWER"
	CallDispositionBusy     = "BUSY"
	CallDispositionFailed   = "FAILED"
)
