package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	GetByID(id int64) (*models.Payment, error)
	GetByReference(reference string) (*models.Payment, error)
	Update(payment *models.Payment) error
	Delete(id int64) error
	ListByCustomer(customerID int64, limit, offset int) ([]*models.Payment, error)
	ListByStatus(status string, limit, offset int) ([]*models.Payment, error)
	ListByDateRange(startDate, endDate time.Time, limit, offset int) ([]*models.Payment, error)
	GetCustomerBalance(customerID int64) (float64, error)
	GetCustomerPaymentHistory(customerID int64, days int) ([]*models.Payment, error)
	GetDailyRevenue(date time.Time) (float64, error)
	GetMonthlyRevenue(year int, month int) (float64, error)
	GetTotalRevenue() (float64, error)
	GetPendingPayments() ([]*models.Payment, error)
	ProcessPayment(paymentID int64, status string, processedBy int64) error
}

type PostgresPaymentRepository struct {
	db *sqlx.DB
}

func NewPostgresPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

func (r *PostgresPaymentRepository) Create(payment *models.Payment) error {
	query := `
		INSERT INTO payments (
			customer_id, payment_reference, payment_type, amount, currency, description,
			payment_method, transaction_id, gateway_response, status, notes
		) VALUES (
			:customer_id, :payment_reference, :payment_type, :amount, :currency, :description,
			:payment_method, :transaction_id, :gateway_response, :status, :notes
		) RETURNING id, created_at, updated_at`
	
	rows, err := r.db.NamedQuery(query, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created payment")
}

func (r *PostgresPaymentRepository) GetByID(id int64) (*models.Payment, error) {
	payment := &models.Payment{}
	query := `SELECT * FROM payments WHERE id = $1`
	
	err := r.db.Get(payment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by ID: %w", err)
	}
	
	return payment, nil
}

func (r *PostgresPaymentRepository) GetByReference(reference string) (*models.Payment, error) {
	payment := &models.Payment{}
	query := `SELECT * FROM payments WHERE payment_reference = $1`
	
	err := r.db.Get(payment, query, reference)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by reference: %w", err)
	}
	
	return payment, nil
}

func (r *PostgresPaymentRepository) Update(payment *models.Payment) error {
	query := `
		UPDATE payments SET
			customer_id = :customer_id, payment_reference = :payment_reference, payment_type = :payment_type,
			amount = :amount, currency = :currency, description = :description, payment_method = :payment_method,
			transaction_id = :transaction_id, gateway_response = :gateway_response, status = :status,
			processed_at = :processed_at, processed_by = :processed_by, notes = :notes,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`
	
	_, err := r.db.NamedExec(query, payment)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	
	return nil
}

func (r *PostgresPaymentRepository) Delete(id int64) error {
	query := `DELETE FROM payments WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	
	return nil
}

func (r *PostgresPaymentRepository) ListByCustomer(customerID int64, limit, offset int) ([]*models.Payment, error) {
	var payments []*models.Payment
	query := `
		SELECT * FROM payments 
		WHERE customer_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`
	
	err := r.db.Select(&payments, query, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments by customer: %w", err)
	}
	
	return payments, nil
}

func (r *PostgresPaymentRepository) ListByStatus(status string, limit, offset int) ([]*models.Payment, error) {
	var payments []*models.Payment
	query := `
		SELECT * FROM payments 
		WHERE status = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`
	
	err := r.db.Select(&payments, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments by status: %w", err)
	}
	
	return payments, nil
}

func (r *PostgresPaymentRepository) ListByDateRange(startDate, endDate time.Time, limit, offset int) ([]*models.Payment, error) {
	var payments []*models.Payment
	query := `
		SELECT * FROM payments 
		WHERE created_at >= $1 AND created_at <= $2 
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4`
	
	err := r.db.Select(&payments, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments by date range: %w", err)
	}
	
	return payments, nil
}

func (r *PostgresPaymentRepository) GetCustomerBalance(customerID int64) (float64, error) {
	var balance sql.NullFloat64
	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN payment_type IN ('credit', 'auto_recharge') AND status = 'completed' THEN amount
				WHEN payment_type IN ('debit') AND status = 'completed' THEN -amount
				ELSE 0
			END
		), 0) as balance
		FROM payments 
		WHERE customer_id = $1`
	
	err := r.db.Get(&balance, query, customerID)
	if err != nil {
		return 0, fmt.Errorf("failed to get customer balance: %w", err)
	}
	
	if balance.Valid {
		return balance.Float64, nil
	}
	return 0, nil
}

func (r *PostgresPaymentRepository) GetCustomerPaymentHistory(customerID int64, days int) ([]*models.Payment, error) {
	var payments []*models.Payment
	query := `
		SELECT * FROM payments 
		WHERE customer_id = $1 
		  AND created_at >= CURRENT_DATE - INTERVAL '%d days'
		ORDER BY created_at DESC`
	
	formattedQuery := fmt.Sprintf(query, days)
	err := r.db.Select(&payments, formattedQuery, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer payment history: %w", err)
	}
	
	return payments, nil
}

func (r *PostgresPaymentRepository) GetDailyRevenue(date time.Time) (float64, error) {
	var revenue sql.NullFloat64
	query := `
		SELECT COALESCE(SUM(amount), 0) as revenue
		FROM payments 
		WHERE payment_type IN ('debit') 
		  AND status = 'completed'
		  AND DATE(created_at) = DATE($1)`
	
	err := r.db.Get(&revenue, query, date)
	if err != nil {
		return 0, fmt.Errorf("failed to get daily revenue: %w", err)
	}
	
	if revenue.Valid {
		return revenue.Float64, nil
	}
	return 0, nil
}

func (r *PostgresPaymentRepository) GetMonthlyRevenue(year int, month int) (float64, error) {
	var revenue sql.NullFloat64
	query := `
		SELECT COALESCE(SUM(amount), 0) as revenue
		FROM payments 
		WHERE payment_type IN ('debit') 
		  AND status = 'completed'
		  AND EXTRACT(YEAR FROM created_at) = $1
		  AND EXTRACT(MONTH FROM created_at) = $2`
	
	err := r.db.Get(&revenue, query, year, month)
	if err != nil {
		return 0, fmt.Errorf("failed to get monthly revenue: %w", err)
	}
	
	if revenue.Valid {
		return revenue.Float64, nil
	}
	return 0, nil
}

func (r *PostgresPaymentRepository) GetTotalRevenue() (float64, error) {
	var revenue sql.NullFloat64
	query := `
		SELECT COALESCE(SUM(amount), 0) as revenue
		FROM payments 
		WHERE payment_type IN ('debit') 
		  AND status = 'completed'`
	
	err := r.db.Get(&revenue, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get total revenue: %w", err)
	}
	
	if revenue.Valid {
		return revenue.Float64, nil
	}
	return 0, nil
}

func (r *PostgresPaymentRepository) GetPendingPayments() ([]*models.Payment, error) {
	var payments []*models.Payment
	query := `
		SELECT * FROM payments 
		WHERE status = 'pending' 
		ORDER BY created_at ASC`
	
	err := r.db.Select(&payments, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending payments: %w", err)
	}
	
	return payments, nil
}

func (r *PostgresPaymentRepository) ProcessPayment(paymentID int64, status string, processedBy int64) error {
	query := `
		UPDATE payments 
		SET status = $1, processed_at = CURRENT_TIMESTAMP, processed_by = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`
	
	_, err := r.db.Exec(query, status, processedBy, paymentID)
	if err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}
	
	return nil
}
