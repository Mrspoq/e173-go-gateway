package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
	List(limit, offset int) ([]*models.User, error)
	UpdateLastLogin(userID int64) error
	IncrementFailedAttempts(userID int64) error
	LockAccount(userID int64, lockDuration time.Duration) error
	UnlockAccount(userID int64) error
	GetByRole(role string) ([]*models.User, error)
	SetPassword(userID int64, passwordHash string) error
	Enable2FA(userID int64, secret string) error
	Disable2FA(userID int64) error
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role, is_active, is_2fa_enabled, two_fa_secret)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`
	
	err := r.db.QueryRow(context.Background(), query, 
		user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName, 
		user.Role, user.IsActive, user.Is2FAEnabled, user.TwoFASecret).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, first_name, last_name, role, is_active, 
	          is_2fa_enabled, two_fa_secret, last_login_at, failed_login_attempts, locked_until, 
	          created_at, updated_at FROM users WHERE id = $1`
	
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.IsActive, &user.Is2FAEnabled, &user.TwoFASecret, &user.LastLoginAt,
		&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return user, nil
}

func (r *PostgresUserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, first_name, last_name, role, is_active, 
	          is_2fa_enabled, two_fa_secret, last_login_at, failed_login_attempts, locked_until, 
	          created_at, updated_at FROM users WHERE username = $1`
	
	err := r.db.QueryRow(context.Background(), query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.IsActive, &user.Is2FAEnabled, &user.TwoFASecret, &user.LastLoginAt,
		&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, first_name, last_name, role, is_active, 
	          is_2fa_enabled, two_fa_secret, last_login_at, failed_login_attempts, locked_until, 
	          created_at, updated_at FROM users WHERE email = $1`
	
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.IsActive, &user.Is2FAEnabled, &user.TwoFASecret, &user.LastLoginAt,
		&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

func (r *PostgresUserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, first_name = $4, last_name = $5, 
		    role = $6, is_active = $7, is_2fa_enabled = $8, 
		    two_fa_secret = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`
	
	_, err := r.db.Exec(context.Background(), query, user.ID, user.Username, user.Email, 
		user.FirstName, user.LastName, user.Role, user.IsActive, user.Is2FAEnabled, user.TwoFASecret)
	
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) List(limit, offset int) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT id, username, email, password_hash, first_name, last_name, role, is_active, 
	          is_2fa_enabled, two_fa_secret, last_login_at, failed_login_attempts, locked_until, 
	          created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
			&user.Role, &user.IsActive, &user.Is2FAEnabled, &user.TwoFASecret, &user.LastLoginAt,
			&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (r *PostgresUserRepository) UpdateLastLogin(userID int64) error {
	query := `UPDATE users SET last_login_at = CURRENT_TIMESTAMP, failed_login_attempts = 0 WHERE id = $1`
	
	_, err := r.db.Exec(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) IncrementFailedAttempts(userID int64) error {
	query := `UPDATE users SET failed_login_attempts = failed_login_attempts + 1 WHERE id = $1`
	
	_, err := r.db.Exec(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("failed to increment failed attempts: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) LockAccount(userID int64, lockDuration time.Duration) error {
	lockUntil := time.Now().Add(lockDuration)
	query := `UPDATE users SET locked_until = $1 WHERE id = $2`
	
	_, err := r.db.Exec(context.Background(), query, lockUntil, userID)
	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) UnlockAccount(userID int64) error {
	query := `UPDATE users SET locked_until = NULL, failed_login_attempts = 0 WHERE id = $1`
	
	_, err := r.db.Exec(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) GetByRole(role string) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT id, username, email, password_hash, first_name, last_name, role, is_active, 
	          is_2fa_enabled, two_fa_secret, last_login_at, failed_login_attempts, locked_until, 
	          created_at, updated_at FROM users WHERE role = $1 AND is_active = true ORDER BY created_at DESC`
	
	rows, err := r.db.Query(context.Background(), query, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
			&user.Role, &user.IsActive, &user.Is2FAEnabled, &user.TwoFASecret, &user.LastLoginAt,
			&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (r *PostgresUserRepository) SetPassword(userID int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	
	_, err := r.db.Exec(context.Background(), query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) Enable2FA(userID int64, secret string) error {
	query := `UPDATE users SET is_2fa_enabled = true, two_fa_secret = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	
	_, err := r.db.Exec(context.Background(), query, secret, userID)
	if err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) Disable2FA(userID int64) error {
	query := `UPDATE users SET is_2fa_enabled = false, two_fa_secret = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	
	_, err := r.db.Exec(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}
	
	return nil
}
