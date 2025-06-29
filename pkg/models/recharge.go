package models

import (
	"time"
	"database/sql"
)

// RechargeCode represents a recharge voucher code
type RechargeCode struct {
	ID              int64          `json:"id" db:"id"`
	SimCardID       int64          `json:"sim_card_id" db:"sim_card_id"`
	Code            string         `json:"code" db:"code"`
	Amount          float64        `json:"amount" db:"amount"`
	Operator        string         `json:"operator" db:"operator"`
	Status          string         `json:"status" db:"status"` // pending, used, failed, expired
	UsedAt          *time.Time     `json:"used_at" db:"used_at"`
	ExpiryDate      *time.Time     `json:"expiry_date" db:"expiry_date"`
	ResponseMessage sql.NullString `json:"response_message" db:"response_message"`
	CreatedBy       int64          `json:"created_by" db:"created_by"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// RechargeBatch represents a batch of recharge operations
type RechargeBatch struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	TotalCodes  int       `json:"total_codes" db:"total_codes"`
	UsedCodes   int       `json:"used_codes" db:"used_codes"`
	TotalAmount float64   `json:"total_amount" db:"total_amount"`
	Status      string    `json:"status" db:"status"` // draft, processing, completed, failed
	CreatedBy   int64     `json:"created_by" db:"created_by"`
	StartedAt   *time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// RechargeHistory tracks all recharge attempts
type RechargeHistory struct {
	ID             int64          `json:"id" db:"id"`
	SimCardID      int64          `json:"sim_card_id" db:"sim_card_id"`
	RechargeCodeID *int64         `json:"recharge_code_id" db:"recharge_code_id"`
	BatchID        *int64         `json:"batch_id" db:"batch_id"`
	PhoneNumber    string         `json:"phone_number" db:"phone_number"`
	Amount         float64        `json:"amount" db:"amount"`
	BalanceBefore  sql.NullFloat64 `json:"balance_before" db:"balance_before"`
	BalanceAfter   sql.NullFloat64 `json:"balance_after" db:"balance_after"`
	Method         string         `json:"method" db:"method"` // ussd, sms, api
	Status         string         `json:"status" db:"status"` // success, failed, pending
	ErrorMessage   sql.NullString `json:"error_message" db:"error_message"`
	Attempts       int            `json:"attempts" db:"attempts"`
	ProcessedBy    int64          `json:"processed_by" db:"processed_by"`
	ProcessedAt    time.Time      `json:"processed_at" db:"processed_at"`
}

// Recharge status constants
const (
	RechargeStatusPending = "pending"
	RechargeStatusUsed    = "used"
	RechargeStatusFailed  = "failed"
	RechargeStatusExpired = "expired"
	
	BatchStatusDraft      = "draft"
	BatchStatusProcessing = "processing"
	BatchStatusCompleted  = "completed"
	BatchStatusFailed     = "failed"
	
	RechargeMethodUSSD = "ussd"
	RechargeMethodSMS  = "sms"
	RechargeMethodAPI  = "api"
)