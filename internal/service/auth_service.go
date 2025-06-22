package service

import (
	"fmt"
	"time"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"

	"github.com/e173-gateway/e173_go_gateway/internal/repository"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type AuthService interface {
	Login(username, password, ipAddress, userAgent string) (*models.User, *models.UserSession, error)
	Logout(sessionToken string) error
	ValidateSession(sessionToken string) (*models.User, error)
	ChangePassword(userID int64, oldPassword, newPassword string) error
	ResetPassword(userID int64, newPassword string) error
	Enable2FA(userID int64) (string, error)
	Disable2FA(userID int64) error
	Verify2FA(userID int64, token string) error
	RefreshSession(sessionToken string) error
	GetUserSessions(userID int64) ([]*models.UserSession, error)
	RevokeUserSessions(userID int64) error
	CleanupExpiredSessions() error
}

type PostgresAuthService struct {
	userRepo   repository.UserRepository
	systemRepo repository.SystemRepository
}

func NewPostgresAuthService(userRepo repository.UserRepository, systemRepo repository.SystemRepository) AuthService {
	return &PostgresAuthService{
		userRepo:   userRepo,
		systemRepo: systemRepo,
	}
}

func (s *PostgresAuthService) Login(username, password, ipAddress, userAgent string) (*models.User, *models.UserSession, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	if user == nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}
	
	// Check if user is active
	if !user.IsActive {
		return nil, nil, fmt.Errorf("account is disabled")
	}
	
	// Check if user is locked
	if user.IsLocked() {
		return nil, nil, fmt.Errorf("account is temporarily locked")
	}
	
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// Increment failed attempts
		s.userRepo.IncrementFailedAttempts(user.ID)
		
		// Lock account if too many failed attempts
		if user.FailedLoginAttempts+1 >= 5 { // TODO: make configurable
			lockDuration := 30 * time.Minute // TODO: make configurable
			s.userRepo.LockAccount(user.ID, lockDuration)
		}
		
		return nil, nil, fmt.Errorf("invalid credentials")
	}
	
	// Generate session token
	sessionToken, err := s.generateSessionToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate session token: %w", err)
	}
	
	// Create session
	session := &models.UserSession{
		UserID:         user.ID,
		SessionToken:   sessionToken,
		IPAddress:      &ipAddress,
		UserAgent:      &userAgent,
		ExpiresAt:      time.Now().Add(24 * time.Hour), // TODO: make configurable
		IsActive:       true,
		LastActivityAt: time.Now(),
	}
	
	err = s.systemRepo.CreateSession(session)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}
	
	// Update last login
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Warning: failed to update last login for user %d: %v\n", user.ID, err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &user.ID,
		Action:     "login",
		IPAddress:  &ipAddress,
		UserAgent:  &userAgent,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return user, session, nil
}

func (s *PostgresAuthService) Logout(sessionToken string) error {
	// Get session to get user ID for audit log
	session, err := s.systemRepo.GetSessionByToken(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	
	// Delete session
	err = s.systemRepo.DeleteSession(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	
	// Create audit log
	if session != nil {
		auditLog := &models.AuditLog{
			UserID:     &session.UserID,
			Action:     "logout",
			IPAddress:  session.IPAddress,
			UserAgent:  session.UserAgent,
			Success:    true,
		}
		s.systemRepo.CreateAuditLog(auditLog)
	}
	
	return nil
}

func (s *PostgresAuthService) ValidateSession(sessionToken string) (*models.User, error) {
	// Get session
	session, err := s.systemRepo.GetSessionByToken(sessionToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	
	if session == nil {
		return nil, fmt.Errorf("invalid session")
	}
	
	if session.IsExpired() {
		s.systemRepo.DeleteSession(sessionToken)
		return nil, fmt.Errorf("session expired")
	}
	
	// Get user
	user, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	if user == nil {
		s.systemRepo.DeleteSession(sessionToken)
		return nil, fmt.Errorf("user not found")
	}
	
	if !user.IsActive {
		s.systemRepo.DeleteSession(sessionToken)
		return nil, fmt.Errorf("user account disabled")
	}
	
	// Update session activity
	s.systemRepo.UpdateSessionActivity(sessionToken)
	
	return user, nil
}

func (s *PostgresAuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	
	if user == nil {
		return fmt.Errorf("user not found")
	}
	
	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return fmt.Errorf("invalid current password")
	}
	
	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Update password
	err = s.userRepo.SetPassword(userID, string(passwordHash))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &userID,
		Action:     "change_password",
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresAuthService) ResetPassword(userID int64, newPassword string) error {
	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Update password
	err = s.userRepo.SetPassword(userID, string(passwordHash))
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}
	
	// Unlock account if locked
	err = s.userRepo.UnlockAccount(userID)
	if err != nil {
		// Log error but don't fail reset
		fmt.Printf("Warning: failed to unlock account for user %d: %v\n", userID, err)
	}
	
	// Revoke all user sessions
	s.systemRepo.DeleteUserSessions(userID)
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &userID,
		Action:     "reset_password",
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresAuthService) Enable2FA(userID int64) (string, error) {
	// Generate 2FA secret
	secret, err := s.generate2FASecret()
	if err != nil {
		return "", fmt.Errorf("failed to generate 2FA secret: %w", err)
	}
	
	// Enable 2FA for user
	err = s.userRepo.Enable2FA(userID, secret)
	if err != nil {
		return "", fmt.Errorf("failed to enable 2FA: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &userID,
		Action:     "enable_2fa",
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return secret, nil
}

func (s *PostgresAuthService) Disable2FA(userID int64) error {
	// Disable 2FA for user
	err := s.userRepo.Disable2FA(userID)
	if err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &userID,
		Action:     "disable_2fa",
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresAuthService) Verify2FA(userID int64, token string) error {
	// TODO: Implement TOTP verification
	// For now, just return success
	return nil
}

func (s *PostgresAuthService) RefreshSession(sessionToken string) error {
	session, err := s.systemRepo.GetSessionByToken(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	
	if session == nil {
		return fmt.Errorf("session not found")
	}
	
	// Extend session expiry
	session.ExpiresAt = time.Now().Add(24 * time.Hour) // TODO: make configurable
	
	// Update session activity
	return s.systemRepo.UpdateSessionActivity(sessionToken)
}

func (s *PostgresAuthService) GetUserSessions(userID int64) ([]*models.UserSession, error) {
	return s.systemRepo.GetUserSessions(userID)
}

func (s *PostgresAuthService) RevokeUserSessions(userID int64) error {
	err := s.systemRepo.DeleteUserSessions(userID)
	if err != nil {
		return fmt.Errorf("failed to revoke user sessions: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &userID,
		Action:     "revoke_sessions",
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresAuthService) CleanupExpiredSessions() error {
	return s.systemRepo.DeleteExpiredSessions()
}

func (s *PostgresAuthService) generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *PostgresAuthService) generate2FASecret() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
