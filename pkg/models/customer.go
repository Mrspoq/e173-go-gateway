package models

import (
	"time"
)

// Customer represents a customer in the CRM system
type Customer struct {
	ID                      int64      `json:"id" db:"id"`
	CustomerCode            string     `json:"customer_code" db:"customer_code"`
	CompanyName             *string    `json:"company_name" db:"company_name"`
	ContactPerson           *string    `json:"contact_person" db:"contact_person"`
	Email                   *string    `json:"email" db:"email"`
	Phone                   *string    `json:"phone" db:"phone"`
	Address                 *string    `json:"address" db:"address"`
	City                    *string    `json:"city" db:"city"`
	State                   *string    `json:"state" db:"state"`
	Country                 *string    `json:"country" db:"country"`
	PostalCode              *string    `json:"postal_code" db:"postal_code"`
	BillingAddress          *string    `json:"billing_address" db:"billing_address"`
	BillingCity             *string    `json:"billing_city" db:"billing_city"`
	BillingState            *string    `json:"billing_state" db:"billing_state"`
	BillingCountry          *string    `json:"billing_country" db:"billing_country"`
	BillingPostalCode       *string    `json:"billing_postal_code" db:"billing_postal_code"`
	AccountStatus           string     `json:"account_status" db:"account_status"`
	CreditLimit             float64    `json:"credit_limit" db:"credit_limit"`
	CurrentBalance          float64    `json:"current_balance" db:"current_balance"`
	MonthlyLimit            *float64   `json:"monthly_limit" db:"monthly_limit"`
	Timezone                string     `json:"timezone" db:"timezone"`
	PreferredCurrency       string     `json:"preferred_currency" db:"preferred_currency"`
	AutoRechargeEnabled     bool       `json:"auto_recharge_enabled" db:"auto_recharge_enabled"`
	AutoRechargeThreshold   *float64   `json:"auto_recharge_threshold" db:"auto_recharge_threshold"`
	AutoRechargeAmount      *float64   `json:"auto_recharge_amount" db:"auto_recharge_amount"`
	Notes                   *string    `json:"notes" db:"notes"`
	CreatedBy               *int64     `json:"created_by" db:"created_by"`
	AssignedTo              *int64     `json:"assigned_to" db:"assigned_to"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
}

// Customer status constants
const (
	CustomerStatusActive     = "active"
	CustomerStatusSuspended  = "suspended"
	CustomerStatusTerminated = "terminated"
	CustomerStatusPending    = "pending"
)

// DisplayName returns the best display name for the customer
func (c *Customer) DisplayName() string {
	if c.CompanyName != nil && *c.CompanyName != "" {
		return *c.CompanyName
	}
	if c.ContactPerson != nil && *c.ContactPerson != "" {
		return *c.ContactPerson
	}
	return c.CustomerCode
}

// IsActive returns true if the customer account is active
func (c *Customer) IsActive() bool {
	return c.AccountStatus == CustomerStatusActive
}

// NeedsAutoRecharge returns true if auto-recharge should be triggered
func (c *Customer) NeedsAutoRecharge() bool {
	return c.AutoRechargeEnabled && 
		   c.AutoRechargeThreshold != nil && 
		   c.CurrentBalance <= *c.AutoRechargeThreshold
}

// Payment represents a payment transaction
type Payment struct {
	ID                int64      `json:"id" db:"id"`
	CustomerID        int64      `json:"customer_id" db:"customer_id"`
	PaymentReference  string     `json:"payment_reference" db:"payment_reference"`
	PaymentType       string     `json:"payment_type" db:"payment_type"`
	Amount            float64    `json:"amount" db:"amount"`
	Currency          string     `json:"currency" db:"currency"`
	Description       *string    `json:"description" db:"description"`
	PaymentMethod     *string    `json:"payment_method" db:"payment_method"`
	TransactionID     *string    `json:"transaction_id" db:"transaction_id"`
	GatewayResponse   *string    `json:"gateway_response" db:"gateway_response"` // JSON string
	Status            string     `json:"status" db:"status"`
	ProcessedAt       *time.Time `json:"processed_at" db:"processed_at"`
	ProcessedBy       *int64     `json:"processed_by" db:"processed_by"`
	Notes             *string    `json:"notes" db:"notes"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// Payment type constants
const (
	PaymentTypeCredit       = "credit"
	PaymentTypeDebit        = "debit"
	PaymentTypeRefund       = "refund"
	PaymentTypeAdjustment   = "adjustment"
	PaymentTypeAutoRecharge = "auto_recharge"
)

// Payment status constants
const (
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"
	PaymentStatusCancelled = "cancelled"
)

// IsCompleted returns true if the payment is successfully completed
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

// IsCredit returns true if this is a credit transaction (adds to balance)
func (p *Payment) IsCredit() bool {
	return p.PaymentType == PaymentTypeCredit || p.PaymentType == PaymentTypeAutoRecharge
}

// RatePlan represents a billing rate plan
type RatePlan struct {
	ID                     int64      `json:"id" db:"id"`
	PlanName               string     `json:"plan_name" db:"plan_name"`
	PlanCode               string     `json:"plan_code" db:"plan_code"`
	Description            *string    `json:"description" db:"description"`
	Currency               string     `json:"currency" db:"currency"`
	RatePerMinute          float64    `json:"rate_per_minute" db:"rate_per_minute"`
	RatePerSecond          float64    `json:"rate_per_second" db:"rate_per_second"`
	MinimumBillingSeconds  int        `json:"minimum_billing_seconds" db:"minimum_billing_seconds"`
	ConnectionFee          float64    `json:"connection_fee" db:"connection_fee"`
	DailyCap               *float64   `json:"daily_cap" db:"daily_cap"`
	MonthlyCap             *float64   `json:"monthly_cap" db:"monthly_cap"`
	EffectiveFrom          time.Time  `json:"effective_from" db:"effective_from"`
	EffectiveUntil         *time.Time `json:"effective_until" db:"effective_until"`
	IsActive               bool       `json:"is_active" db:"is_active"`
	CreatedBy              *int64     `json:"created_by" db:"created_by"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
}

// CalculateCallCost calculates the cost for a call duration in seconds
func (rp *RatePlan) CalculateCallCost(durationSeconds int) float64 {
	if durationSeconds <= 0 {
		return rp.ConnectionFee
	}
	
	billableSeconds := durationSeconds
	if billableSeconds < rp.MinimumBillingSeconds {
		billableSeconds = rp.MinimumBillingSeconds
	}
	
	return rp.ConnectionFee + (float64(billableSeconds) * rp.RatePerSecond)
}

// IsValidForDate returns true if the rate plan is valid for the given date
func (rp *RatePlan) IsValidForDate(date time.Time) bool {
	if !rp.IsActive {
		return false
	}
	
	if date.Before(rp.EffectiveFrom) {
		return false
	}
	
	if rp.EffectiveUntil != nil && date.After(*rp.EffectiveUntil) {
		return false
	}
	
	return true
}

// CustomerRatePlan represents the assignment of a rate plan to a customer
type CustomerRatePlan struct {
	ID             int64      `json:"id" db:"id"`
	CustomerID     int64      `json:"customer_id" db:"customer_id"`
	RatePlanID     int64      `json:"rate_plan_id" db:"rate_plan_id"`
	EffectiveFrom  time.Time  `json:"effective_from" db:"effective_from"`
	EffectiveUntil *time.Time `json:"effective_until" db:"effective_until"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedBy      *int64     `json:"created_by" db:"created_by"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}
