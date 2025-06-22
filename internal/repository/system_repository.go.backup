package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type SystemRepository interface {
	// System Config
	GetConfigByKey(key string) (*models.SystemConfig, error)
	SetConfig(key, value, configType string, userID int64) error
	GetConfigsByCategory(category string) ([]*models.SystemConfig, error)
	GetAllConfigs() ([]*models.SystemConfig, error)
	DeleteConfig(key string) error
	
	// User Sessions
	CreateSession(session *models.UserSession) error
	GetSessionByToken(token string) (*models.UserSession, error)
	UpdateSessionActivity(token string) error
	DeleteSession(token string) error
	DeleteExpiredSessions() error
	GetUserSessions(userID int64) ([]*models.UserSession, error)
	DeleteUserSessions(userID int64) error
	
	// Audit Logs
	CreateAuditLog(log *models.AuditLog) error
	GetAuditLogsByUser(userID int64, limit, offset int) ([]*models.AuditLog, error)
	GetAuditLogsByAction(action string, limit, offset int) ([]*models.AuditLog, error)
	GetAuditLogsByEntity(entityType string, entityID int64, limit, offset int) ([]*models.AuditLog, error)
	GetAuditLogsByDateRange(startDate, endDate time.Time, limit, offset int) ([]*models.AuditLog, error)
	
	// User Notifications
	CreateNotification(notification *models.UserNotification) error
	GetUserNotifications(userID int64, unreadOnly bool, limit, offset int) ([]*models.UserNotification, error)
	MarkNotificationRead(notificationID int64) error
	MarkAllNotificationsRead(userID int64) error
	DeleteNotification(notificationID int64) error
	GetUnreadNotificationCount(userID int64) (int64, error)
	
	// Notification Templates
	CreateNotificationTemplate(template *models.NotificationTemplate) error
	GetNotificationTemplateByName(name string) (*models.NotificationTemplate, error)
	UpdateNotificationTemplate(template *models.NotificationTemplate) error
	DeleteNotificationTemplate(name string) error
	ListNotificationTemplates() ([]*models.NotificationTemplate, error)
}

type PostgresSystemRepository struct {
	db *sqlx.DB
}

func NewPostgresSystemRepository(db *sqlx.DB) SystemRepository {
	return &PostgresSystemRepository{db: db}
}

// System Config methods
func (r *PostgresSystemRepository) GetConfigByKey(key string) (*models.SystemConfig, error) {
	config := &models.SystemConfig{}
	query := `SELECT * FROM system_config WHERE config_key = $1`
	
	err := r.db.Get(config, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get config by key: %w", err)
	}
	
	return config, nil
}

func (r *PostgresSystemRepository) SetConfig(key, value, configType string, userID int64) error {
	query := `
		INSERT INTO system_config (config_key, config_value, config_type, updated_by)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (config_key) 
		DO UPDATE SET config_value = $2, config_type = $3, updated_by = $4, updated_at = CURRENT_TIMESTAMP`
	
	_, err := r.db.Exec(query, key, value, configType, userID)
	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) GetConfigsByCategory(category string) ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig
	query := `SELECT * FROM system_config WHERE category = $1 ORDER BY config_key ASC`
	
	err := r.db.Select(&configs, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get configs by category: %w", err)
	}
	
	return configs, nil
}

func (r *PostgresSystemRepository) GetAllConfigs() ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig
	query := `SELECT * FROM system_config ORDER BY category ASC, config_key ASC`
	
	err := r.db.Select(&configs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all configs: %w", err)
	}
	
	return configs, nil
}

func (r *PostgresSystemRepository) DeleteConfig(key string) error {
	query := `DELETE FROM system_config WHERE config_key = $1 AND is_system = false`
	
	_, err := r.db.Exec(query, key)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}
	
	return nil
}

// User Sessions methods
func (r *PostgresSystemRepository) CreateSession(session *models.UserSession) error {
	query := `
		INSERT INTO user_sessions (user_id, session_token, ip_address, user_agent, expires_at)
		VALUES (:user_id, :session_token, :ip_address, :user_agent, :expires_at)
		RETURNING id, created_at`
	
	rows, err := r.db.NamedQuery(query, session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&session.ID, &session.CreatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created session")
}

func (r *PostgresSystemRepository) GetSessionByToken(token string) (*models.UserSession, error) {
	session := &models.UserSession{}
	query := `SELECT * FROM user_sessions WHERE session_token = $1 AND is_active = true AND expires_at > CURRENT_TIMESTAMP`
	
	err := r.db.Get(session, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session by token: %w", err)
	}
	
	return session, nil
}

func (r *PostgresSystemRepository) UpdateSessionActivity(token string) error {
	query := `UPDATE user_sessions SET last_activity_at = CURRENT_TIMESTAMP WHERE session_token = $1`
	
	_, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to update session activity: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) DeleteSession(token string) error {
	query := `DELETE FROM user_sessions WHERE session_token = $1`
	
	_, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) DeleteExpiredSessions() error {
	query := `DELETE FROM user_sessions WHERE expires_at <= CURRENT_TIMESTAMP`
	
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) GetUserSessions(userID int64) ([]*models.UserSession, error) {
	var sessions []*models.UserSession
	query := `
		SELECT * FROM user_sessions 
		WHERE user_id = $1 AND is_active = true AND expires_at > CURRENT_TIMESTAMP
		ORDER BY last_activity_at DESC`
	
	err := r.db.Select(&sessions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	
	return sessions, nil
}

func (r *PostgresSystemRepository) DeleteUserSessions(userID int64) error {
	query := `DELETE FROM user_sessions WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	
	return nil
}

// Audit Log methods
func (r *PostgresSystemRepository) CreateAuditLog(log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			user_id, action, entity_type, entity_id, old_values, new_values,
			ip_address, user_agent, success, error_message
		) VALUES (
			:user_id, :action, :entity_type, :entity_id, :old_values, :new_values,
			:ip_address, :user_agent, :success, :error_message
		) RETURNING id, created_at`
	
	rows, err := r.db.NamedQuery(query, log)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&log.ID, &log.CreatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created audit log")
}

func (r *PostgresSystemRepository) GetAuditLogsByUser(userID int64, limit, offset int) ([]*models.AuditLog, error) {
	var logs []*models.AuditLog
	query := `
		SELECT * FROM audit_logs 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`
	
	err := r.db.Select(&logs, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by user: %w", err)
	}
	
	return logs, nil
}

func (r *PostgresSystemRepository) GetAuditLogsByAction(action string, limit, offset int) ([]*models.AuditLog, error) {
	var logs []*models.AuditLog
	query := `
		SELECT * FROM audit_logs 
		WHERE action = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`
	
	err := r.db.Select(&logs, query, action, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by action: %w", err)
	}
	
	return logs, nil
}

func (r *PostgresSystemRepository) GetAuditLogsByEntity(entityType string, entityID int64, limit, offset int) ([]*models.AuditLog, error) {
	var logs []*models.AuditLog
	query := `
		SELECT * FROM audit_logs 
		WHERE entity_type = $1 AND entity_id = $2 
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4`
	
	err := r.db.Select(&logs, query, entityType, entityID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by entity: %w", err)
	}
	
	return logs, nil
}

func (r *PostgresSystemRepository) GetAuditLogsByDateRange(startDate, endDate time.Time, limit, offset int) ([]*models.AuditLog, error) {
	var logs []*models.AuditLog
	query := `
		SELECT * FROM audit_logs 
		WHERE created_at >= $1 AND created_at <= $2 
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4`
	
	err := r.db.Select(&logs, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by date range: %w", err)
	}
	
	return logs, nil
}

// User Notifications methods
func (r *PostgresSystemRepository) CreateNotification(notification *models.UserNotification) error {
	query := `
		INSERT INTO user_notifications (
			user_id, notification_type, title, message, priority, action_url, metadata
		) VALUES (
			:user_id, :notification_type, :title, :message, :priority, :action_url, :metadata
		) RETURNING id, created_at`
	
	rows, err := r.db.NamedQuery(query, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&notification.ID, &notification.CreatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created notification")
}

func (r *PostgresSystemRepository) GetUserNotifications(userID int64, unreadOnly bool, limit, offset int) ([]*models.UserNotification, error) {
	var notifications []*models.UserNotification
	var query string
	
	if unreadOnly {
		query = `
			SELECT * FROM user_notifications 
			WHERE user_id = $1 AND is_read = false 
			ORDER BY created_at DESC 
			LIMIT $2 OFFSET $3`
	} else {
		query = `
			SELECT * FROM user_notifications 
			WHERE user_id = $1 
			ORDER BY created_at DESC 
			LIMIT $2 OFFSET $3`
	}
	
	err := r.db.Select(&notifications, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}
	
	return notifications, nil
}

func (r *PostgresSystemRepository) MarkNotificationRead(notificationID int64) error {
	query := `UPDATE user_notifications SET is_read = true, read_at = CURRENT_TIMESTAMP WHERE id = $1`
	
	_, err := r.db.Exec(query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification read: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) MarkAllNotificationsRead(userID int64) error {
	query := `UPDATE user_notifications SET is_read = true, read_at = CURRENT_TIMESTAMP WHERE user_id = $1 AND is_read = false`
	
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications read: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) DeleteNotification(notificationID int64) error {
	query := `DELETE FROM user_notifications WHERE id = $1`
	
	_, err := r.db.Exec(query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) GetUnreadNotificationCount(userID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM user_notifications WHERE user_id = $1 AND is_read = false`
	
	err := r.db.Get(&count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread notification count: %w", err)
	}
	
	return count, nil
}

// Notification Template methods
func (r *PostgresSystemRepository) CreateNotificationTemplate(template *models.NotificationTemplate) error {
	query := `
		INSERT INTO notification_templates (
			template_name, template_type, subject_template, body_template, variables, is_active, created_by
		) VALUES (
			:template_name, :template_type, :subject_template, :body_template, :variables, :is_active, :created_by
		) RETURNING id, created_at, updated_at`
	
	rows, err := r.db.NamedQuery(query, template)
	if err != nil {
		return fmt.Errorf("failed to create notification template: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created notification template")
}

func (r *PostgresSystemRepository) GetNotificationTemplateByName(name string) (*models.NotificationTemplate, error) {
	template := &models.NotificationTemplate{}
	query := `SELECT * FROM notification_templates WHERE template_name = $1`
	
	err := r.db.Get(template, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notification template by name: %w", err)
	}
	
	return template, nil
}

func (r *PostgresSystemRepository) UpdateNotificationTemplate(template *models.NotificationTemplate) error {
	query := `
		UPDATE notification_templates SET
			template_name = :template_name, template_type = :template_type,
			subject_template = :subject_template, body_template = :body_template,
			variables = :variables, is_active = :is_active, updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`
	
	_, err := r.db.NamedExec(query, template)
	if err != nil {
		return fmt.Errorf("failed to update notification template: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) DeleteNotificationTemplate(name string) error {
	query := `DELETE FROM notification_templates WHERE template_name = $1`
	
	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete notification template: %w", err)
	}
	
	return nil
}

func (r *PostgresSystemRepository) ListNotificationTemplates() ([]*models.NotificationTemplate, error) {
	var templates []*models.NotificationTemplate
	query := `SELECT * FROM notification_templates ORDER BY template_name ASC`
	
	err := r.db.Select(&templates, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification templates: %w", err)
	}
	
	return templates, nil
}
