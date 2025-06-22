package models

import (
	"time"
	"strconv"
	"encoding/json"
)

// SystemConfig represents a system configuration setting
type SystemConfig struct {
	ID          int64     `json:"id" db:"id"`
	ConfigKey   string    `json:"config_key" db:"config_key"`
	ConfigValue *string   `json:"config_value" db:"config_value"`
	ConfigType  string    `json:"config_type" db:"config_type"`
	Description *string   `json:"description" db:"description"`
	IsEncrypted bool      `json:"is_encrypted" db:"is_encrypted"`
	IsSystem    bool      `json:"is_system" db:"is_system"`
	Category    *string   `json:"category" db:"category"`
	UpdatedBy   *int64    `json:"updated_by" db:"updated_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Config type constants
const (
	ConfigTypeString  = "string"
	ConfigTypeInteger = "integer"
	ConfigTypeFloat   = "float"
	ConfigTypeBoolean = "boolean"
	ConfigTypeJSON    = "json"
)

// Config category constants
const (
	ConfigCategoryGeneral       = "general"
	ConfigCategorySecurity      = "security"
	ConfigCategoryBilling       = "billing"
	ConfigCategoryRouting       = "routing"
	ConfigCategoryNotifications = "notifications"
)

// GetStringValue returns the config value as a string
func (sc *SystemConfig) GetStringValue() string {
	if sc.ConfigValue == nil {
		return ""
	}
	return *sc.ConfigValue
}

// GetIntValue returns the config value as an integer
func (sc *SystemConfig) GetIntValue() int {
	if sc.ConfigValue == nil {
		return 0
	}
	val, err := strconv.Atoi(*sc.ConfigValue)
	if err != nil {
		return 0
	}
	return val
}

// GetFloatValue returns the config value as a float64
func (sc *SystemConfig) GetFloatValue() float64 {
	if sc.ConfigValue == nil {
		return 0.0
	}
	val, err := strconv.ParseFloat(*sc.ConfigValue, 64)
	if err != nil {
		return 0.0
	}
	return val
}

// GetBoolValue returns the config value as a boolean
func (sc *SystemConfig) GetBoolValue() bool {
	if sc.ConfigValue == nil {
		return false
	}
	val, err := strconv.ParseBool(*sc.ConfigValue)
	if err != nil {
		return false
	}
	return val
}

// GetJSONValue unmarshals the config value into the provided interface
func (sc *SystemConfig) GetJSONValue(v interface{}) error {
	if sc.ConfigValue == nil {
		return nil
	}
	return json.Unmarshal([]byte(*sc.ConfigValue), v)
}

// NotificationTemplate represents a template for system notifications
type NotificationTemplate struct {
	ID              int64     `json:"id" db:"id"`
	TemplateName    string    `json:"template_name" db:"template_name"`
	TemplateType    string    `json:"template_type" db:"template_type"`
	SubjectTemplate *string   `json:"subject_template" db:"subject_template"`
	BodyTemplate    string    `json:"body_template" db:"body_template"`
	Variables       *string   `json:"variables" db:"variables"` // JSON string
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedBy       *int64    `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Template type constants
const (
	TemplateTypeEmail    = "email"
	TemplateTypeSMS      = "sms"
	TemplateTypeWebhook  = "webhook"
	TemplateTypeInternal = "internal"
)

// Render renders the template with the provided variables
func (nt *NotificationTemplate) Render(variables map[string]interface{}) (subject, body string) {
	// Simple variable replacement - can be enhanced with a proper template engine
	body = nt.BodyTemplate
	if nt.SubjectTemplate != nil {
		subject = *nt.SubjectTemplate
	}
	
	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		valueStr := ""
		if value != nil {
			valueStr = value.(string)
		}
		if nt.SubjectTemplate != nil {
			subject = replaceAll(subject, placeholder, valueStr)
		}
		body = replaceAll(body, placeholder, valueStr)
	}
	
	return subject, body
}

// Simple string replacement function
func replaceAll(text, old, new string) string {
	// This is a simplified version - in production you'd use strings.ReplaceAll or a template engine
	result := ""
	for i := 0; i < len(text); {
		if i+len(old) <= len(text) && text[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(text[i])
			i++
		}
	}
	return result
}

// CDRBilling represents billing information calculated for a CDR
type CDRBilling struct {
	ID                int64      `json:"id" db:"id"`
	CDRID             int64      `json:"cdr_id" db:"cdr_id"`
	CustomerID        *int64     `json:"customer_id" db:"customer_id"`
	RatePlanID        *int64     `json:"rate_plan_id" db:"rate_plan_id"`
	RatePerSecond     float64    `json:"rate_per_second" db:"rate_per_second"`
	ConnectionFee     float64    `json:"connection_fee" db:"connection_fee"`
	BillableSeconds   int        `json:"billable_seconds" db:"billable_seconds"`
	CallCost          float64    `json:"call_cost" db:"call_cost"`
	MarkupPercent     float64    `json:"markup_percent" db:"markup_percent"`
	TotalCost         float64    `json:"total_cost" db:"total_cost"`
	Currency          string     `json:"currency" db:"currency"`
	BillingStatus     string     `json:"billing_status" db:"billing_status"`
	ProcessedAt       *time.Time `json:"processed_at" db:"processed_at"`
	ProcessedBy       *int64     `json:"processed_by" db:"processed_by"`
	ErrorMessage      *string    `json:"error_message" db:"error_message"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// Billing status constants
const (
	BillingStatusPending   = "pending"
	BillingStatusProcessed = "processed"
	BillingStatusFailed    = "failed"
	BillingStatusSkipped   = "skipped"
)

// Stats represents real-time system statistics
type Stats struct {
	TotalModems        int     `json:"total_modems"`
	OnlineModems       int     `json:"online_modems"`
	OfflineModems      int     `json:"offline_modems"`
	TotalSIMCards      int     `json:"total_sim_cards"`
	ActiveSIMCards     int     `json:"active_sim_cards"`
	TotalCustomers     int     `json:"total_customers"`
	ActiveCustomers    int     `json:"active_customers"`
	TotalCallsToday    int     `json:"total_calls_today"`
	TotalCallsThisWeek int     `json:"total_calls_this_week"`
	SpamCallsToday     int     `json:"spam_calls_today"`
	TotalRevenue       float64 `json:"total_revenue"`
	RevenueToday       float64 `json:"revenue_today"`
	ActiveChannels     int     `json:"active_channels"`
	QueuedCalls        int     `json:"queued_calls"`
}

// SystemHealth represents overall system health metrics
type SystemHealth struct {
	Status               string    `json:"status"` // "healthy", "warning", "critical"
	DatabaseConnected    bool      `json:"database_connected"`
	AsteriskConnected    bool      `json:"asterisk_connected"`
	ModemsOnline         int       `json:"modems_online"`
	ModemsOffline        int       `json:"modems_offline"`
	LastHealthCheck      time.Time `json:"last_health_check"`
	ErrorCount24h        int       `json:"error_count_24h"`
	AverageResponseTime  float64   `json:"average_response_time_ms"`
	DiskUsagePercent     float64   `json:"disk_usage_percent"`
	MemoryUsagePercent   float64   `json:"memory_usage_percent"`
	CPUUsagePercent      float64   `json:"cpu_usage_percent"`
}
