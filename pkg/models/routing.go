package models

import (
	"time"
	"github.com/lib/pq"
)

// RoutingRule represents a call routing rule
type RoutingRule struct {
	ID                    int64           `json:"id" db:"id"`
	RuleName              string          `json:"rule_name" db:"rule_name"`
	RuleOrder             int             `json:"rule_order" db:"rule_order"`
	PrefixPattern         string          `json:"prefix_pattern" db:"prefix_pattern"`
	DestinationPattern    *string         `json:"destination_pattern" db:"destination_pattern"`
	CallerIDPattern       *string         `json:"caller_id_pattern" db:"caller_id_pattern"`
	RouteToModemID        *int64          `json:"route_to_modem_id" db:"route_to_modem_id"`
	RouteToPool           *string         `json:"route_to_pool" db:"route_to_pool"`
	MaxChannels           int             `json:"max_channels" db:"max_channels"`
	TimeRestrictions      *string         `json:"time_restrictions" db:"time_restrictions"` // JSON string
	CustomerRestrictions  pq.Int64Array   `json:"customer_restrictions" db:"customer_restrictions"`
	CostMarkupPercent     float64         `json:"cost_markup_percent" db:"cost_markup_percent"`
	IsActive              bool            `json:"is_active" db:"is_active"`
	Notes                 *string         `json:"notes" db:"notes"`
	CreatedBy             *int64          `json:"created_by" db:"created_by"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`
}

// MatchesNumber checks if this routing rule matches the given number
func (rr *RoutingRule) MatchesNumber(number string) bool {
	if !rr.IsActive {
		return false
	}
	
	// Simple prefix matching for now - can be enhanced with regex later
	if len(number) >= len(rr.PrefixPattern) {
		return number[:len(rr.PrefixPattern)] == rr.PrefixPattern
	}
	
	return false
}

// Blacklist represents a blacklisted number or pattern
type Blacklist struct {
	ID                 int64      `json:"id" db:"id"`
	NumberPattern      string     `json:"number_pattern" db:"number_pattern"`
	BlacklistType      string     `json:"blacklist_type" db:"blacklist_type"`
	Reason             *string    `json:"reason" db:"reason"`
	AutoAdded          bool       `json:"auto_added" db:"auto_added"`
	DetectionMethod    *string    `json:"detection_method" db:"detection_method"`
	BlockInbound       bool       `json:"block_inbound" db:"block_inbound"`
	BlockOutbound      bool       `json:"block_outbound" db:"block_outbound"`
	TemporaryUntil     *time.Time `json:"temporary_until" db:"temporary_until"`
	ViolationCount     int        `json:"violation_count" db:"violation_count"`
	LastViolationAt    time.Time  `json:"last_violation_at" db:"last_violation_at"`
	CreatedBy          *int64     `json:"created_by" db:"created_by"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// Blacklist type constants
const (
	BlacklistTypeNumber  = "number"
	BlacklistTypePattern = "pattern"
	BlacklistTypePrefix  = "prefix"
)

// Detection method constants
const (
	DetectionShortCall     = "short_call"
	DetectionHighFrequency = "high_frequency"
	DetectionManual        = "manual"
	DetectionSpamReport    = "spam_report"
)

// IsActive returns true if the blacklist entry is currently active
func (bl *Blacklist) IsActive() bool {
	if bl.TemporaryUntil == nil {
		return true // Permanent blacklist
	}
	return time.Now().Before(*bl.TemporaryUntil)
}

// MatchesNumber checks if this blacklist entry matches the given number
func (bl *Blacklist) MatchesNumber(number string) bool {
	if !bl.IsActive() {
		return false
	}
	
	switch bl.BlacklistType {
	case BlacklistTypeNumber:
		return number == bl.NumberPattern
	case BlacklistTypePrefix:
		if len(number) >= len(bl.NumberPattern) {
			return number[:len(bl.NumberPattern)] == bl.NumberPattern
		}
		return false
	case BlacklistTypePattern:
		// Simple wildcard matching - can be enhanced with regex
		if bl.NumberPattern == "*" {
			return true
		}
		// More sophisticated pattern matching can be added here
		return number == bl.NumberPattern
	}
	
	return false
}

// ShouldBlock returns true if this entry should block the given direction
func (bl *Blacklist) ShouldBlock(direction string) bool {
	if !bl.IsActive() {
		return false
	}
	
	switch direction {
	case "inbound":
		return bl.BlockInbound
	case "outbound":
		return bl.BlockOutbound
	default:
		return bl.BlockInbound || bl.BlockOutbound
	}
}

// SIMPool represents a pool of SIM cards for load balancing
type SIMPool struct {
	ID                   int64     `json:"id" db:"id"`
	PoolName             string    `json:"pool_name" db:"pool_name"`
	Description          *string   `json:"description" db:"description"`
	LoadBalanceMethod    string    `json:"load_balance_method" db:"load_balance_method"`
	MaxChannelsPerSIM    int       `json:"max_channels_per_sim" db:"max_channels_per_sim"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	CreatedBy            *int64    `json:"created_by" db:"created_by"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// Load balance method constants
const (
	LoadBalanceRoundRobin = "round_robin"
	LoadBalanceLeastUsed  = "least_used"
	LoadBalanceRandom     = "random"
	LoadBalanceFailover   = "failover"
)

// SIMPoolAssignment represents the assignment of a SIM card to a pool
type SIMPoolAssignment struct {
	ID         int64     `json:"id" db:"id"`
	SIMPoolID  int64     `json:"sim_pool_id" db:"sim_pool_id"`
	SIMCardID  int64     `json:"sim_card_id" db:"sim_card_id"`
	Priority   int       `json:"priority" db:"priority"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
	AssignedBy *int64    `json:"assigned_by" db:"assigned_by"`
}

// CallRoutingResult represents the result of call routing logic
type CallRoutingResult struct {
	Success          bool     `json:"success"`
	RouteToModemID   *int64   `json:"route_to_modem_id,omitempty"`
	RouteToPool      *string  `json:"route_to_pool,omitempty"`
	SelectedSIMID    *int64   `json:"selected_sim_id,omitempty"`
	RoutingRuleID    *int64   `json:"routing_rule_id,omitempty"`
	IsBlocked        bool     `json:"is_blocked"`
	BlockReason      *string  `json:"block_reason,omitempty"`
	CostMarkup       float64  `json:"cost_markup"`
	ErrorMessage     *string  `json:"error_message,omitempty"`
}
