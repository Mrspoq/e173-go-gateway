package models

import (
	"time"
)

// User represents a system user with role-based access control
type User struct {
	ID                   int64      `json:"id" db:"id"`
	Username             string     `json:"username" db:"username"`
	Email                string     `json:"email" db:"email"`
	PasswordHash         string     `json:"-" db:"password_hash"` // Never expose password hash in JSON
	FirstName            *string    `json:"first_name" db:"first_name"`
	LastName             *string    `json:"last_name" db:"last_name"`
	Role                 string     `json:"role" db:"role"`
	IsActive             bool       `json:"is_active" db:"is_active"`
	Is2FAEnabled         bool       `json:"is_2fa_enabled" db:"is_2fa_enabled"`
	TwoFASecret          *string    `json:"-" db:"two_fa_secret"` // Never expose 2FA secret
	LastLoginAt          *time.Time `json:"last_login_at" db:"last_login_at"`
	FailedLoginAttempts  int        `json:"failed_login_attempts" db:"failed_login_attempts"`
	LockedUntil          *time.Time `json:"locked_until" db:"locked_until"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// UserRole constants
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleManager    = "manager"
	RoleEmployee   = "employee"
	RoleViewer     = "viewer"
)

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName != nil && u.LastName != nil {
		return *u.FirstName + " " + *u.LastName
	}
	if u.FirstName != nil {
		return *u.FirstName
	}
	if u.LastName != nil {
		return *u.LastName
	}
	return u.Username
}

// IsLocked returns true if the user account is currently locked
func (u *User) IsLocked() bool {
	return u.LockedUntil != nil && u.LockedUntil.After(time.Now())
}

// HasRole checks if the user has the specified role or higher
func (u *User) HasRole(role string) bool {
	roleHierarchy := map[string]int{
		RoleViewer:     1,
		RoleEmployee:   2,
		RoleManager:    3,
		RoleAdmin:      4,
		RoleSuperAdmin: 5,
	}
	
	userLevel, userExists := roleHierarchy[u.Role]
	requiredLevel, requiredExists := roleHierarchy[role]
	
	if !userExists || !requiredExists {
		return false
	}
	
	return userLevel >= requiredLevel
}

// UserSession represents an active user session
type UserSession struct {
	ID               int64      `json:"id" db:"id"`
	UserID           int64      `json:"user_id" db:"user_id"`
	SessionToken     string     `json:"-" db:"session_token"` // Never expose session token
	IPAddress        *string    `json:"ip_address" db:"ip_address"`
	UserAgent        *string    `json:"user_agent" db:"user_agent"`
	ExpiresAt        time.Time  `json:"expires_at" db:"expires_at"`
	IsActive         bool       `json:"is_active" db:"is_active"`
	LastActivityAt   time.Time  `json:"last_activity_at" db:"last_activity_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// IsExpired returns true if the session has expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// AuditLog represents a system audit log entry
type AuditLog struct {
	ID           int64       `json:"id" db:"id"`
	UserID       *int64      `json:"user_id" db:"user_id"`
	Action       string      `json:"action" db:"action"`
	EntityType   *string     `json:"entity_type" db:"entity_type"`
	EntityID     *int64      `json:"entity_id" db:"entity_id"`
	OldValues    *string     `json:"old_values" db:"old_values"` // JSON string
	NewValues    *string     `json:"new_values" db:"new_values"` // JSON string
	IPAddress    *string     `json:"ip_address" db:"ip_address"`
	UserAgent    *string     `json:"user_agent" db:"user_agent"`
	Success      bool        `json:"success" db:"success"`
	ErrorMessage *string     `json:"error_message" db:"error_message"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
}

// UserNotification represents a notification for a specific user
type UserNotification struct {
	ID               int64      `json:"id" db:"id"`
	UserID           int64      `json:"user_id" db:"user_id"`
	NotificationType string     `json:"notification_type" db:"notification_type"`
	Title            string     `json:"title" db:"title"`
	Message          string     `json:"message" db:"message"`
	Priority         string     `json:"priority" db:"priority"`
	IsRead           bool       `json:"is_read" db:"is_read"`
	ReadAt           *time.Time `json:"read_at" db:"read_at"`
	ActionURL        *string    `json:"action_url" db:"action_url"`
	Metadata         *string    `json:"metadata" db:"metadata"` // JSON string
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// NotificationPriority constants
const (
	PriorityLow      = "low"
	PriorityNormal   = "normal"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)
