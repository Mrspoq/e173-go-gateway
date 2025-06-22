package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/e173-gateway/e173_go_gateway/internal/repository"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type CustomerService interface {
	CreateCustomer(customer *models.Customer, createdBy int64) error
	GetCustomerByID(id int64) (*models.Customer, error)
	GetCustomerByCode(code string) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer, updatedBy int64) error
	DeleteCustomer(id int64, deletedBy int64) error
	ListCustomers(limit, offset int) ([]*models.Customer, error)
	SearchCustomers(query string, limit, offset int) ([]*models.Customer, error)
	GetCustomersByStatus(status string, limit, offset int) ([]*models.Customer, error)
	GetCustomersByManager(userID int64, limit, offset int) ([]*models.Customer, error)
	
	// Balance management
	GetCustomerBalance(customerID int64) (float64, error)
	UpdateCustomerBalance(customerID int64, amount float64, transactionType string, description string, processedBy int64) error
	GetLowBalanceCustomers(threshold float64) ([]*models.Customer, error)
	GetCustomersNeedingRecharge() ([]*models.Customer, error)
	
	// Account management
	SuspendCustomer(customerID int64, reason string, suspendedBy int64) error
	ReactivateCustomer(customerID int64, reactivatedBy int64) error
	TerminateCustomer(customerID int64, reason string, terminatedBy int64) error
	
	// Statistics
	GetTotalCustomerCount() (int64, error)
	GetActiveCustomerCount() (int64, error)
	GetCustomerStats() (*CustomerStats, error)
}

type CustomerStats struct {
	TotalCustomers      int64   `json:"total_customers"`
	ActiveCustomers     int64   `json:"active_customers"`
	SuspendedCustomers  int64   `json:"suspended_customers"`
	TerminatedCustomers int64   `json:"terminated_customers"`
	LowBalanceCustomers int64   `json:"low_balance_customers"`
	TotalBalance        float64 `json:"total_balance"`
	AverageBalance      float64 `json:"average_balance"`
}

type PostgresCustomerService struct {
	customerRepo repository.CustomerRepository
	paymentRepo  repository.PaymentRepository
	systemRepo   repository.SystemRepository
}

func NewPostgresCustomerService(
	customerRepo repository.CustomerRepository,
	paymentRepo repository.PaymentRepository,
	systemRepo repository.SystemRepository,
) CustomerService {
	return &PostgresCustomerService{
		customerRepo: customerRepo,
		paymentRepo:  paymentRepo,
		systemRepo:   systemRepo,
	}
}

func (s *PostgresCustomerService) CreateCustomer(customer *models.Customer, createdBy int64) error {
	// Generate unique customer code if not provided
	if customer.CustomerCode == "" {
		code, err := s.generateCustomerCode()
		if err != nil {
			return fmt.Errorf("failed to generate customer code: %w", err)
		}
		customer.CustomerCode = code
	}
	
	// Validate unique customer code
	existing, err := s.customerRepo.GetByCustomerCode(customer.CustomerCode)
	if err != nil {
		return fmt.Errorf("failed to check customer code uniqueness: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("customer code already exists")
	}
	
	// Set defaults
	if customer.AccountStatus == "" {
		customer.AccountStatus = models.CustomerStatusActive
	}
	if customer.Timezone == "" {
		customer.Timezone = "UTC"
	}
	if customer.PreferredCurrency == "" {
		customer.PreferredCurrency = "USD"
	}
	
	customer.CreatedBy = &createdBy
	
	// Create customer
	err = s.customerRepo.Create(customer)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &createdBy,
		Action:     "create_customer",
		EntityType: stringPtr("customer"),
		EntityID:   &customer.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	// Create notification for assigned manager
	if customer.AssignedTo != nil {
		notification := &models.UserNotification{
			UserID:           *customer.AssignedTo,
			NotificationType: "new_customer_assigned",
			Title:            "New Customer Assigned",
			Message:          fmt.Sprintf("New customer %s has been assigned to you", customer.DisplayName()),
			Priority:         models.PriorityNormal,
			ActionURL:        stringPtr(fmt.Sprintf("/customers/%d", customer.ID)),
		}
		s.systemRepo.CreateNotification(notification)
	}
	
	return nil
}

func (s *PostgresCustomerService) GetCustomerByID(id int64) (*models.Customer, error) {
	return s.customerRepo.GetByID(id)
}

func (s *PostgresCustomerService) GetCustomerByCode(code string) (*models.Customer, error) {
	return s.customerRepo.GetByCustomerCode(code)
}

func (s *PostgresCustomerService) UpdateCustomer(customer *models.Customer, updatedBy int64) error {
	// Get existing customer for audit log
	oldCustomer, err := s.customerRepo.GetByID(customer.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing customer: %w", err)
	}
	
	if oldCustomer == nil {
		return fmt.Errorf("customer not found")
	}
	
	// Update customer
	err = s.customerRepo.Update(customer)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &updatedBy,
		Action:     "update_customer",
		EntityType: stringPtr("customer"),
		EntityID:   &customer.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	// Notify if status changed
	if oldCustomer.AccountStatus != customer.AccountStatus {
		if customer.AssignedTo != nil {
			notification := &models.UserNotification{
				UserID:           *customer.AssignedTo,
				NotificationType: "customer_status_changed",
				Title:            "Customer Status Changed",
				Message:          fmt.Sprintf("Customer %s status changed from %s to %s", customer.DisplayName(), oldCustomer.AccountStatus, customer.AccountStatus),
				Priority:         models.PriorityNormal,
				ActionURL:        stringPtr(fmt.Sprintf("/customers/%d", customer.ID)),
			}
			s.systemRepo.CreateNotification(notification)
		}
	}
	
	return nil
}

func (s *PostgresCustomerService) DeleteCustomer(id int64, deletedBy int64) error {
	// Get customer for audit log
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}
	
	if customer == nil {
		return fmt.Errorf("customer not found")
	}
	
	// Delete customer
	err = s.customerRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &deletedBy,
		Action:     "delete_customer",
		EntityType: stringPtr("customer"),
		EntityID:   &id,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresCustomerService) ListCustomers(limit, offset int) ([]*models.Customer, error) {
	return s.customerRepo.List(limit, offset)
}

func (s *PostgresCustomerService) SearchCustomers(query string, limit, offset int) ([]*models.Customer, error) {
	return s.customerRepo.Search(query, limit, offset)
}

func (s *PostgresCustomerService) GetCustomersByStatus(status string, limit, offset int) ([]*models.Customer, error) {
	return s.customerRepo.ListByStatus(status, limit, offset)
}

func (s *PostgresCustomerService) GetCustomersByManager(userID int64, limit, offset int) ([]*models.Customer, error) {
	return s.customerRepo.ListByAssignedTo(userID, limit, offset)
}

func (s *PostgresCustomerService) GetCustomerBalance(customerID int64) (float64, error) {
	return s.paymentRepo.GetCustomerBalance(customerID)
}

func (s *PostgresCustomerService) UpdateCustomerBalance(customerID int64, amount float64, transactionType string, description string, processedBy int64) error {
	// Generate payment reference
	reference, err := s.generatePaymentReference()
	if err != nil {
		return fmt.Errorf("failed to generate payment reference: %w", err)
	}
	
	// Create payment record
	payment := &models.Payment{
		CustomerID:       customerID,
		PaymentReference: reference,
		PaymentType:      transactionType,
		Amount:           amount,
		Currency:         "USD", // TODO: get from customer preferred currency
		Description:      &description,
		Status:           models.PaymentStatusCompleted,
		ProcessedAt:      &time.Time{},
		ProcessedBy:      &processedBy,
	}
	*payment.ProcessedAt = time.Now()
	
	err = s.paymentRepo.Create(payment)
	if err != nil {
		return fmt.Errorf("failed to create payment record: %w", err)
	}
	
	// Update customer balance
	newBalance, err := s.paymentRepo.GetCustomerBalance(customerID)
	if err != nil {
		return fmt.Errorf("failed to calculate new balance: %w", err)
	}
	
	err = s.customerRepo.UpdateBalance(customerID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update customer balance: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &processedBy,
		Action:     "update_balance",
		EntityType: stringPtr("customer"),
		EntityID:   &customerID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresCustomerService) GetLowBalanceCustomers(threshold float64) ([]*models.Customer, error) {
	return s.customerRepo.GetLowBalanceCustomers(threshold)
}

func (s *PostgresCustomerService) GetCustomersNeedingRecharge() ([]*models.Customer, error) {
	return s.customerRepo.GetCustomersNeedingRecharge()
}

func (s *PostgresCustomerService) SuspendCustomer(customerID int64, reason string, suspendedBy int64) error {
	customer, err := s.customerRepo.GetByID(customerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}
	
	if customer == nil {
		return fmt.Errorf("customer not found")
	}
	
	customer.AccountStatus = models.CustomerStatusSuspended
	customer.Notes = stringPtr(fmt.Sprintf("Suspended: %s", reason))
	
	err = s.customerRepo.Update(customer)
	if err != nil {
		return fmt.Errorf("failed to suspend customer: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &suspendedBy,
		Action:     "suspend_customer",
		EntityType: stringPtr("customer"),
		EntityID:   &customerID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresCustomerService) ReactivateCustomer(customerID int64, reactivatedBy int64) error {
	customer, err := s.customerRepo.GetByID(customerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}
	
	if customer == nil {
		return fmt.Errorf("customer not found")
	}
	
	customer.AccountStatus = models.CustomerStatusActive
	
	err = s.customerRepo.Update(customer)
	if err != nil {
		return fmt.Errorf("failed to reactivate customer: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &reactivatedBy,
		Action:     "reactivate_customer",
		EntityType: stringPtr("customer"),
		EntityID:   &customerID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresCustomerService) TerminateCustomer(customerID int64, reason string, terminatedBy int64) error {
	customer, err := s.customerRepo.GetByID(customerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}
	
	if customer == nil {
		return fmt.Errorf("customer not found")
	}
	
	customer.AccountStatus = models.CustomerStatusTerminated
	customer.Notes = stringPtr(fmt.Sprintf("Terminated: %s", reason))
	
	err = s.customerRepo.Update(customer)
	if err != nil {
		return fmt.Errorf("failed to terminate customer: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &terminatedBy,
		Action:     "terminate_customer",
		EntityType: stringPtr("customer"),
		EntityID:   &customerID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresCustomerService) GetTotalCustomerCount() (int64, error) {
	return s.customerRepo.GetTotalCount()
}

func (s *PostgresCustomerService) GetActiveCustomerCount() (int64, error) {
	return s.customerRepo.GetActiveCount()
}

func (s *PostgresCustomerService) GetCustomerStats() (*CustomerStats, error) {
	totalCustomers, err := s.customerRepo.GetTotalCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get total customer count: %w", err)
	}
	
	activeCustomers, err := s.customerRepo.GetActiveCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get active customer count: %w", err)
	}
	
	// Get low balance customers
	lowBalanceCustomers, err := s.customerRepo.GetLowBalanceCustomers(10.0) // TODO: make configurable
	if err != nil {
		return nil, fmt.Errorf("failed to get low balance customers: %w", err)
	}
	
	stats := &CustomerStats{
		TotalCustomers:      totalCustomers,
		ActiveCustomers:     activeCustomers,
		LowBalanceCustomers: int64(len(lowBalanceCustomers)),
		// TODO: Calculate other stats
	}
	
	return stats, nil
}

func (s *PostgresCustomerService) generateCustomerCode() (string, error) {
	// Generate a random 6-character hex string
	bytes := make([]byte, 3)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	
	code := "CUST" + hex.EncodeToString(bytes)
	
	// Check if code already exists
	existing, err := s.customerRepo.GetByCustomerCode(code)
	if err != nil {
		return "", err
	}
	
	if existing != nil {
		// Recursively try again if code exists
		return s.generateCustomerCode()
	}
	
	return code, nil
}

func (s *PostgresCustomerService) generatePaymentReference() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	
	return "PAY" + hex.EncodeToString(bytes), nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
