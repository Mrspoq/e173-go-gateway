package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	GetByID(id int64) (*models.Customer, error)
	GetByCustomerCode(customerCode string) (*models.Customer, error)
	Update(customer *models.Customer) error
	Delete(id int64) error
	List(limit, offset int) ([]*models.Customer, error)
	ListByStatus(status string, limit, offset int) ([]*models.Customer, error)
	ListByAssignedTo(userID int64, limit, offset int) ([]*models.Customer, error)
	UpdateBalance(customerID int64, newBalance float64) error
	GetLowBalanceCustomers(threshold float64) ([]*models.Customer, error)
	GetCustomersNeedingRecharge() ([]*models.Customer, error)
	Search(query string, limit, offset int) ([]*models.Customer, error)
	GetTotalCount() (int64, error)
	GetActiveCount() (int64, error)
}

type PostgresCustomerRepository struct {
	db *sqlx.DB
}

func NewPostgresCustomerRepository(db *sqlx.DB) CustomerRepository {
	return &PostgresCustomerRepository{db: db}
}

func (r *PostgresCustomerRepository) Create(customer *models.Customer) error {
	query := `
		INSERT INTO customers (
			customer_code, company_name, contact_person, email, phone, address, city, state, country, postal_code,
			billing_address, billing_city, billing_state, billing_country, billing_postal_code,
			account_status, credit_limit, current_balance, monthly_limit, timezone, preferred_currency,
			auto_recharge_enabled, auto_recharge_threshold, auto_recharge_amount, notes, created_by, assigned_to
		) VALUES (
			:customer_code, :company_name, :contact_person, :email, :phone, :address, :city, :state, :country, :postal_code,
			:billing_address, :billing_city, :billing_state, :billing_country, :billing_postal_code,
			:account_status, :credit_limit, :current_balance, :monthly_limit, :timezone, :preferred_currency,
			:auto_recharge_enabled, :auto_recharge_threshold, :auto_recharge_amount, :notes, :created_by, :assigned_to
		) RETURNING id, created_at, updated_at`
	
	rows, err := r.db.NamedQuery(query, customer)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created customer")
}

func (r *PostgresCustomerRepository) GetByID(id int64) (*models.Customer, error) {
	customer := &models.Customer{}
	query := `SELECT * FROM customers WHERE id = $1`
	
	err := r.db.Get(customer, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer by ID: %w", err)
	}
	
	return customer, nil
}

func (r *PostgresCustomerRepository) GetByCustomerCode(customerCode string) (*models.Customer, error) {
	customer := &models.Customer{}
	query := `SELECT * FROM customers WHERE customer_code = $1`
	
	err := r.db.Get(customer, query, customerCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer by code: %w", err)
	}
	
	return customer, nil
}

func (r *PostgresCustomerRepository) Update(customer *models.Customer) error {
	query := `
		UPDATE customers SET
			customer_code = :customer_code, company_name = :company_name, contact_person = :contact_person,
			email = :email, phone = :phone, address = :address, city = :city, state = :state,
			country = :country, postal_code = :postal_code, billing_address = :billing_address,
			billing_city = :billing_city, billing_state = :billing_state, billing_country = :billing_country,
			billing_postal_code = :billing_postal_code, account_status = :account_status,
			credit_limit = :credit_limit, current_balance = :current_balance, monthly_limit = :monthly_limit,
			timezone = :timezone, preferred_currency = :preferred_currency, auto_recharge_enabled = :auto_recharge_enabled,
			auto_recharge_threshold = :auto_recharge_threshold, auto_recharge_amount = :auto_recharge_amount,
			notes = :notes, assigned_to = :assigned_to, updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`
	
	_, err := r.db.NamedExec(query, customer)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	
	return nil
}

func (r *PostgresCustomerRepository) Delete(id int64) error {
	query := `DELETE FROM customers WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	
	return nil
}

func (r *PostgresCustomerRepository) List(limit, offset int) ([]*models.Customer, error) {
	var customers []*models.Customer
	query := `SELECT * FROM customers ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	err := r.db.Select(&customers, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	
	return customers, nil
}

func (r *PostgresCustomerRepository) ListByStatus(status string, limit, offset int) ([]*models.Customer, error) {
	var customers []*models.Customer
	query := `SELECT * FROM customers WHERE account_status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	err := r.db.Select(&customers, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers by status: %w", err)
	}
	
	return customers, nil
}

func (r *PostgresCustomerRepository) ListByAssignedTo(userID int64, limit, offset int) ([]*models.Customer, error) {
	var customers []*models.Customer
	query := `SELECT * FROM customers WHERE assigned_to = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	err := r.db.Select(&customers, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers by assigned user: %w", err)
	}
	
	return customers, nil
}

func (r *PostgresCustomerRepository) UpdateBalance(customerID int64, newBalance float64) error {
	query := `UPDATE customers SET current_balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	
	_, err := r.db.Exec(query, newBalance, customerID)
	if err != nil {
		return fmt.Errorf("failed to update customer balance: %w", err)
	}
	
	return nil
}

func (r *PostgresCustomerRepository) GetLowBalanceCustomers(threshold float64) ([]*models.Customer, error) {
	var customers []*models.Customer
	query := `
		SELECT * FROM customers 
		WHERE account_status = 'active' AND current_balance <= $1 
		ORDER BY current_balance ASC`
	
	err := r.db.Select(&customers, query, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get low balance customers: %w", err)
	}
	
	return customers, nil
}

func (r *PostgresCustomerRepository) GetCustomersNeedingRecharge() ([]*models.Customer, error) {
	var customers []*models.Customer
	query := `
		SELECT * FROM customers 
		WHERE account_status = 'active' 
		  AND auto_recharge_enabled = true 
		  AND auto_recharge_threshold IS NOT NULL 
		  AND current_balance <= auto_recharge_threshold
		ORDER BY current_balance ASC`
	
	err := r.db.Select(&customers, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers needing recharge: %w", err)
	}
	
	return customers, nil
}

func (r *PostgresCustomerRepository) Search(query string, limit, offset int) ([]*models.Customer, error) {
	var customers []*models.Customer
	searchQuery := `
		SELECT * FROM customers 
		WHERE customer_code ILIKE $1 
		   OR company_name ILIKE $1 
		   OR contact_person ILIKE $1 
		   OR email ILIKE $1 
		   OR phone ILIKE $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`
	
	searchTerm := "%" + query + "%"
	err := r.db.Select(&customers, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	
	return customers, nil
}

func (r *PostgresCustomerRepository) GetTotalCount() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM customers`
	
	err := r.db.Get(&count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get total customer count: %w", err)
	}
	
	return count, nil
}

func (r *PostgresCustomerRepository) GetActiveCount() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM customers WHERE account_status = 'active'`
	
	err := r.db.Get(&count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get active customer count: %w", err)
	}
	
	return count, nil
}
