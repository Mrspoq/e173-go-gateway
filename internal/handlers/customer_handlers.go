package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/e173-gateway/e173_go_gateway/internal/service"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type CustomerHandlers struct {
	customerService service.CustomerService
}

func NewCustomerHandlers(customerService service.CustomerService) *CustomerHandlers {
	return &CustomerHandlers{
		customerService: customerService,
	}
}

type CreateCustomerRequest struct {
	CustomerCode            string   `json:"customer_code,omitempty"`
	CompanyName             *string  `json:"company_name"`
	ContactPerson           *string  `json:"contact_person"`
	Email                   *string  `json:"email"`
	Phone                   *string  `json:"phone"`
	Address                 *string  `json:"address"`
	City                    *string  `json:"city"`
	State                   *string  `json:"state"`
	Country                 *string  `json:"country"`
	PostalCode              *string  `json:"postal_code"`
	CreditLimit             float64  `json:"credit_limit"`
	MonthlyLimit            *float64 `json:"monthly_limit"`
	Timezone                string   `json:"timezone"`
	PreferredCurrency       string   `json:"preferred_currency"`
	AutoRechargeEnabled     bool     `json:"auto_recharge_enabled"`
	AutoRechargeThreshold   *float64 `json:"auto_recharge_threshold"`
	AutoRechargeAmount      *float64 `json:"auto_recharge_amount"`
	Notes                   *string  `json:"notes"`
	AssignedTo              *int64   `json:"assigned_to"`
}

type UpdateBalanceRequest struct {
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	Description     string  `json:"description"`
}

func (h *CustomerHandlers) ListCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limit := 50
	offset := 0
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	var customers []*models.Customer
	var err error

	// Handle different query types
	if search != "" {
		customers, err = h.customerService.SearchCustomers(search, limit, offset)
	} else if status != "" {
		customers, err = h.customerService.GetCustomersByStatus(status, limit, offset)
	} else {
		customers, err = h.customerService.ListCustomers(limit, offset)
	}

	if err != nil {
		http.Error(w, "Failed to get customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"customers": customers,
		"count":     len(customers),
	})
}

func (h *CustomerHandlers) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create customer model
	customer := &models.Customer{
		CustomerCode:            req.CustomerCode,
		CompanyName:             req.CompanyName,
		ContactPerson:           req.ContactPerson,
		Email:                   req.Email,
		Phone:                   req.Phone,
		Address:                 req.Address,
		City:                    req.City,
		State:                   req.State,
		Country:                 req.Country,
		PostalCode:              req.PostalCode,
		AccountStatus:           models.CustomerStatusActive,
		CreditLimit:             req.CreditLimit,
		CurrentBalance:          0.0,
		MonthlyLimit:            req.MonthlyLimit,
		Timezone:                req.Timezone,
		PreferredCurrency:       req.PreferredCurrency,
		AutoRechargeEnabled:     req.AutoRechargeEnabled,
		AutoRechargeThreshold:   req.AutoRechargeThreshold,
		AutoRechargeAmount:      req.AutoRechargeAmount,
		Notes:                   req.Notes,
		AssignedTo:              req.AssignedTo,
	}

	err := h.customerService.CreateCustomer(customer, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandlers) GetCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract customer ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customer, err := h.customerService.GetCustomerByID(customerID)
	if err != nil {
		http.Error(w, "Failed to get customer", http.StatusInternalServerError)
		return
	}

	if customer == nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandlers) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract customer ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	customer.ID = customerID
	err = h.customerService.UpdateCustomer(&customer, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandlers) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract customer ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	err = h.customerService.DeleteCustomer(customerID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Customer deleted successfully",
	})
}

func (h *CustomerHandlers) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract customer ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var req UpdateBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate transaction type
	validTypes := []string{models.PaymentTypeCredit, models.PaymentTypeDebit, models.PaymentTypeAdjustment}
	isValid := false
	for _, validType := range validTypes {
		if req.TransactionType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		http.Error(w, "Invalid transaction type", http.StatusBadRequest)
		return
	}

	err = h.customerService.UpdateCustomerBalance(customerID, req.Amount, req.TransactionType, req.Description, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get updated balance
	newBalance, err := h.customerService.GetCustomerBalance(customerID)
	if err != nil {
		http.Error(w, "Failed to get updated balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"message":     "Balance updated successfully",
		"new_balance": newBalance,
	})
}

func (h *CustomerHandlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract customer ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	balance, err := h.customerService.GetCustomerBalance(customerID)
	if err != nil {
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"customer_id": customerID,
		"balance":     balance,
	})
}

func (h *CustomerHandlers) SuspendCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract customer ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.customerService.SuspendCustomer(customerID, req.Reason, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Customer suspended successfully",
	})
}

func (h *CustomerHandlers) ReactivateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract customer ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	err = h.customerService.ReactivateCustomer(customerID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Customer reactivated successfully",
	})
}

func (h *CustomerHandlers) GetCustomerStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.customerService.GetCustomerStats()
	if err != nil {
		http.Error(w, "Failed to get customer statistics", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return HTML for HTMX
		w.Header().Set("Content-Type", "text/html")
		// Return just the count formatted as HTML
		w.Write([]byte(strconv.FormatInt(stats.TotalCustomers, 10)))
		return
	}

	// Return JSON for API requests
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
