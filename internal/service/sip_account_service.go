package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	internalRepo "github.com/e173-gateway/e173_go_gateway/internal/repository"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/e173-gateway/e173_go_gateway/pkg/repository"
	"github.com/sirupsen/logrus"
)

// SIPAccountService handles business logic for SIP accounts
type SIPAccountService interface {
	CreateSIPAccount(ctx context.Context, customerID int64, account *models.SIPAccount, createdBy int64) error
	GetSIPAccountByID(ctx context.Context, id int64) (*models.SIPAccount, error)
	GetSIPAccountByUsername(ctx context.Context, username string) (*models.SIPAccount, error)
	GetCustomerSIPAccounts(ctx context.Context, customerID int64) ([]*models.SIPAccount, error)
	UpdateSIPAccount(ctx context.Context, account *models.SIPAccount) error
	DeleteSIPAccount(ctx context.Context, id int64) error
	ListSIPAccounts(ctx context.Context, limit, offset int) ([]*models.SIPAccount, error)
	SearchSIPAccounts(ctx context.Context, query string, limit, offset int) ([]*models.SIPAccount, error)
	
	// Permissions
	GetSIPAccountPermissions(ctx context.Context, accountID int64) (*models.SIPAccountPermission, error)
	UpdateSIPAccountPermissions(ctx context.Context, permissions *models.SIPAccountPermission) error
	
	// Registration
	RegisterSIPAccount(ctx context.Context, username, contactURI, sourceIP string, sourcePort, expires int, userAgent string) error
	UnregisterSIPAccount(ctx context.Context, username string) error
	GetActiveRegistrations(ctx context.Context, accountID int64) ([]*models.SIPRegistration, error)
	
	// Usage
	RecordCallUsage(ctx context.Context, accountID int64, duration int, isIncoming, isInternational, isSuccessful bool) error
	GetUsageStats(ctx context.Context, accountID int64, startDate, endDate time.Time) ([]*models.SIPAccountUsage, error)
	GetCurrentMonthUsage(ctx context.Context, accountID int64) (*models.SIPAccountUsage, error)
	
	// Validation
	ValidateCallPermission(ctx context.Context, accountID int64, destination string) (bool, string)
	GenerateUniqueUsername(ctx context.Context, customerCode string) (string, error)
	GenerateSecurePassword() (string, error)
}

type sipAccountService struct {
	sipRepo      repository.SIPAccountRepository
	customerRepo internalRepo.CustomerRepository
	logger       *logrus.Logger
}

// NewSIPAccountService creates a new SIP account service
func NewSIPAccountService(
	sipRepo repository.SIPAccountRepository,
	customerRepo internalRepo.CustomerRepository,
	logger *logrus.Logger,
) SIPAccountService {
	return &sipAccountService{
		sipRepo:      sipRepo,
		customerRepo: customerRepo,
		logger:       logger,
	}
}

func (s *sipAccountService) CreateSIPAccount(ctx context.Context, customerID int64, account *models.SIPAccount, createdBy int64) error {
	// Verify customer exists and is active
	customer, err := s.customerRepo.GetByID(customerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}
	
	if customer == nil {
		return fmt.Errorf("customer not found")
	}
	
	if !customer.IsActive() {
		return fmt.Errorf("customer account is not active")
	}
	
	// Set defaults
	account.CustomerID = customerID
	account.CreatedBy = &createdBy
	account.Status = models.SIPAccountStatusActive
	
	if account.Domain == "" {
		account.Domain = "sip.e173gateway.com"
	}
	
	if account.Context == "" {
		account.Context = "default"
	}
	
	if account.Transport == "" {
		account.Transport = models.TransportUDP
	}
	
	if account.CodecsAllowed == "" {
		account.CodecsAllowed = models.DefaultCodecs
	}
	
	if account.MaxConcurrentCalls == 0 {
		account.MaxConcurrentCalls = 2
	}
	
	// Generate username if not provided
	if account.Username == "" {
		username, err := s.GenerateUniqueUsername(ctx, customer.CustomerCode)
		if err != nil {
			return fmt.Errorf("failed to generate username: %w", err)
		}
		account.Username = username
	}
	
	// Generate password if not provided
	if account.Password == "" {
		password, err := s.GenerateSecurePassword()
		if err != nil {
			return fmt.Errorf("failed to generate password: %w", err)
		}
		account.Password = password
	}
	
	// Create the account
	err = s.sipRepo.CreateSIPAccount(ctx, account)
	if err != nil {
		return fmt.Errorf("failed to create SIP account: %w", err)
	}
	
	s.logger.WithFields(logrus.Fields{
		"sip_account_id": account.ID,
		"customer_id":    customerID,
		"username":       account.Username,
	}).Info("SIP account created")
	
	return nil
}

func (s *sipAccountService) GetSIPAccountByID(ctx context.Context, id int64) (*models.SIPAccount, error) {
	account, err := s.sipRepo.GetSIPAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Load customer information
	customer, err := s.customerRepo.GetByID(account.CustomerID)
	if err == nil {
		account.Customer = customer
	}
	
	return account, nil
}

func (s *sipAccountService) GetSIPAccountByUsername(ctx context.Context, username string) (*models.SIPAccount, error) {
	return s.sipRepo.GetSIPAccountByUsername(ctx, username)
}

func (s *sipAccountService) GetCustomerSIPAccounts(ctx context.Context, customerID int64) ([]*models.SIPAccount, error) {
	return s.sipRepo.GetSIPAccountsByCustomerID(ctx, customerID)
}

func (s *sipAccountService) UpdateSIPAccount(ctx context.Context, account *models.SIPAccount) error {
	// Get existing account to check permissions
	existing, err := s.sipRepo.GetSIPAccountByID(ctx, account.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing account: %w", err)
	}
	
	// Don't allow changing certain fields
	account.CustomerID = existing.CustomerID
	account.Username = existing.Username
	account.Password = existing.Password // Password changes should use a separate method
	account.Domain = existing.Domain
	account.CreatedBy = existing.CreatedBy
	account.CreatedAt = existing.CreatedAt
	
	err = s.sipRepo.UpdateSIPAccount(ctx, account)
	if err != nil {
		return fmt.Errorf("failed to update SIP account: %w", err)
	}
	
	s.logger.WithFields(logrus.Fields{
		"sip_account_id": account.ID,
		"status":         account.Status,
	}).Info("SIP account updated")
	
	return nil
}

func (s *sipAccountService) DeleteSIPAccount(ctx context.Context, id int64) error {
	err := s.sipRepo.DeleteSIPAccount(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete SIP account: %w", err)
	}
	
	s.logger.WithField("sip_account_id", id).Info("SIP account deleted")
	return nil
}

func (s *sipAccountService) ListSIPAccounts(ctx context.Context, limit, offset int) ([]*models.SIPAccount, error) {
	return s.sipRepo.ListSIPAccounts(ctx, limit, offset)
}

func (s *sipAccountService) SearchSIPAccounts(ctx context.Context, query string, limit, offset int) ([]*models.SIPAccount, error) {
	return s.sipRepo.SearchSIPAccounts(ctx, query, limit, offset)
}

func (s *sipAccountService) GetSIPAccountPermissions(ctx context.Context, accountID int64) (*models.SIPAccountPermission, error) {
	return s.sipRepo.GetSIPAccountPermissions(ctx, accountID)
}

func (s *sipAccountService) UpdateSIPAccountPermissions(ctx context.Context, permissions *models.SIPAccountPermission) error {
	err := s.sipRepo.UpdateSIPAccountPermissions(ctx, permissions)
	if err != nil {
		return fmt.Errorf("failed to update permissions: %w", err)
	}
	
	s.logger.WithField("sip_account_id", permissions.SIPAccountID).Info("SIP account permissions updated")
	return nil
}

func (s *sipAccountService) RegisterSIPAccount(ctx context.Context, username, contactURI, sourceIP string, sourcePort, expires int, userAgent string) error {
	// Get account by username
	account, err := s.sipRepo.GetSIPAccountByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	
	if !account.IsActive() {
		return fmt.Errorf("account is not active")
	}
	
	// Create registration record
	registration := &models.SIPRegistration{
		SIPAccountID:   account.ID,
		ContactURI:     contactURI,
		SourceIP:       sourceIP,
		SourcePort:     sourcePort,
		UserAgent:      &userAgent,
		ExpiresSeconds: expires,
		RegisteredAt:   time.Now(),
		ExpiredAt:      time.Now().Add(time.Duration(expires) * time.Second),
		IsActive:       true,
	}
	
	err = s.sipRepo.CreateRegistration(ctx, registration)
	if err != nil {
		return fmt.Errorf("failed to create registration: %w", err)
	}
	
	s.logger.WithFields(logrus.Fields{
		"username":   username,
		"source_ip":  sourceIP,
		"expires":    expires,
	}).Info("SIP account registered")
	
	return nil
}

func (s *sipAccountService) UnregisterSIPAccount(ctx context.Context, username string) error {
	// Get account by username
	account, err := s.sipRepo.GetSIPAccountByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	
	err = s.sipRepo.UpdateRegistrationStatus(ctx, account.ID, "", false)
	if err != nil {
		return fmt.Errorf("failed to update registration status: %w", err)
	}
	
	s.logger.WithField("username", username).Info("SIP account unregistered")
	return nil
}

func (s *sipAccountService) GetActiveRegistrations(ctx context.Context, accountID int64) ([]*models.SIPRegistration, error) {
	return s.sipRepo.GetActiveRegistrations(ctx, accountID)
}

func (s *sipAccountService) RecordCallUsage(ctx context.Context, accountID int64, duration int, isIncoming, isInternational, isSuccessful bool) error {
	today := time.Now().Truncate(24 * time.Hour)
	
	// Get or create today's usage record
	usage, err := s.sipRepo.GetUsageByDate(ctx, accountID, today)
	if err != nil {
		if err == repository.ErrNotFound {
			usage = &models.SIPAccountUsage{
				SIPAccountID: accountID,
				Date:         today,
			}
		} else {
			return fmt.Errorf("failed to get usage: %w", err)
		}
	}
	
	// Update usage statistics
	usage.TotalCalls++
	if isSuccessful {
		usage.SuccessfulCalls++
	} else {
		usage.FailedCalls++
	}
	
	if duration > 0 {
		usage.TotalMinutes += (duration + 59) / 60 // Round up to nearest minute
	}
	
	if isIncoming {
		usage.IncomingCalls++
	} else {
		usage.OutgoingCalls++
	}
	
	if isInternational {
		usage.InternationalCalls++
	}
	
	// Calculate average call duration
	if usage.TotalCalls > 0 {
		usage.AverageCallDuration = (usage.TotalMinutes * 60) / usage.TotalCalls
	}
	
	// Update the account's total statistics
	_, err = s.sipRepo.GetSIPAccountByID(ctx, accountID)
	if err == nil {
		// Update account totals (this would be done through a separate query)
		// For now, we just record the daily usage
	}
	
	err = s.sipRepo.RecordUsage(ctx, usage)
	if err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}
	
	return nil
}

func (s *sipAccountService) GetUsageStats(ctx context.Context, accountID int64, startDate, endDate time.Time) ([]*models.SIPAccountUsage, error) {
	return s.sipRepo.GetUsageStats(ctx, accountID, startDate, endDate)
}

func (s *sipAccountService) GetCurrentMonthUsage(ctx context.Context, accountID int64) (*models.SIPAccountUsage, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	
	stats, err := s.sipRepo.GetUsageStats(ctx, accountID, startOfMonth, endOfMonth)
	if err != nil {
		return nil, err
	}
	
	// Aggregate monthly stats
	monthlyUsage := &models.SIPAccountUsage{
		SIPAccountID: accountID,
		Date:         startOfMonth,
	}
	
	for _, daily := range stats {
		monthlyUsage.TotalCalls += daily.TotalCalls
		monthlyUsage.SuccessfulCalls += daily.SuccessfulCalls
		monthlyUsage.FailedCalls += daily.FailedCalls
		monthlyUsage.TotalMinutes += daily.TotalMinutes
		monthlyUsage.IncomingCalls += daily.IncomingCalls
		monthlyUsage.OutgoingCalls += daily.OutgoingCalls
		monthlyUsage.InternationalCalls += daily.InternationalCalls
		if daily.PeakConcurrentCalls > monthlyUsage.PeakConcurrentCalls {
			monthlyUsage.PeakConcurrentCalls = daily.PeakConcurrentCalls
		}
	}
	
	if monthlyUsage.TotalCalls > 0 {
		monthlyUsage.AverageCallDuration = (monthlyUsage.TotalMinutes * 60) / monthlyUsage.TotalCalls
	}
	
	return monthlyUsage, nil
}

func (s *sipAccountService) ValidateCallPermission(ctx context.Context, accountID int64, destination string) (bool, string) {
	// Get account
	account, err := s.sipRepo.GetSIPAccountByID(ctx, accountID)
	if err != nil {
		return false, "Account not found"
	}
	
	if !account.IsActive() {
		return false, "Account is not active"
	}
	
	if !account.CanMakeCall() {
		return false, "Maximum concurrent calls reached"
	}
	
	// Get permissions
	permissions, err := s.sipRepo.GetSIPAccountPermissions(ctx, accountID)
	if err != nil {
		return false, "Failed to get permissions"
	}
	
	// Check if it's an international call
	isInternational := strings.HasPrefix(destination, "+") && !strings.HasPrefix(destination, "+1")
	
	if isInternational && !permissions.AllowInternational {
		return false, "International calls not allowed"
	}
	
	// Check blocked prefixes
	if permissions.BlockedPrefixes != nil && *permissions.BlockedPrefixes != "" {
		blockedPrefixes := strings.Split(*permissions.BlockedPrefixes, ",")
		for _, prefix := range blockedPrefixes {
			if strings.HasPrefix(destination, strings.TrimSpace(prefix)) {
				return false, "Destination prefix is blocked"
			}
		}
	}
	
	// Check allowed prefixes
	if permissions.AllowedPrefixes != nil && *permissions.AllowedPrefixes != "" {
		allowedPrefixes := strings.Split(*permissions.AllowedPrefixes, ",")
		allowed := false
		for _, prefix := range allowedPrefixes {
			if strings.HasPrefix(destination, strings.TrimSpace(prefix)) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "Destination prefix not allowed"
		}
	}
	
	// Check daily limits
	if permissions.DailyCallLimit != nil {
		usage, err := s.sipRepo.GetUsageByDate(ctx, accountID, time.Now().Truncate(24*time.Hour))
		if err == nil && usage.TotalCalls >= *permissions.DailyCallLimit {
			return false, "Daily call limit reached"
		}
	}
	
	// TODO: Check time restrictions, country restrictions, etc.
	
	return true, ""
}

func (s *sipAccountService) GenerateUniqueUsername(ctx context.Context, customerCode string) (string, error) {
	// Clean customer code
	base := strings.ToLower(strings.ReplaceAll(customerCode, "-", ""))
	
	// Try to generate a unique username
	for i := 1; i <= 999; i++ {
		username := fmt.Sprintf("%s%03d", base, i)
		
		// Check if username exists
		_, err := s.sipRepo.GetSIPAccountByUsername(ctx, username)
		if err == repository.ErrNotFound {
			return username, nil
		} else if err != nil {
			return "", fmt.Errorf("failed to check username: %w", err)
		}
	}
	
	return "", fmt.Errorf("unable to generate unique username")
}

func (s *sipAccountService) GenerateSecurePassword() (string, error) {
	return models.GenerateSecurePassword(16)
}