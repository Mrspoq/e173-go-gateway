package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/e173-gateway/e173_go_gateway/internal/service"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SIPAccountHandler handles SIP account-related HTTP requests
type SIPAccountHandler struct {
	sipService service.SIPAccountService
	logger     *logrus.Logger
}

// NewSIPAccountHandler creates a new SIP account handler
func NewSIPAccountHandler(sipService service.SIPAccountService, logger *logrus.Logger) *SIPAccountHandler {
	return &SIPAccountHandler{
		sipService: sipService,
		logger:     logger,
	}
}

// CreateSIPAccountRequest represents the request to create a SIP account
type CreateSIPAccountRequest struct {
	AccountName          string `json:"account_name" binding:"required"`
	Extension            string `json:"extension"`
	CallerID             string `json:"caller_id"`
	CallerIDName         string `json:"caller_id_name"`
	MaxConcurrentCalls   int    `json:"max_concurrent_calls"`
	NATSupport           bool   `json:"nat_support"`
	EncryptionEnabled    bool   `json:"encryption_enabled"`
	CodecsAllowed        string `json:"codecs_allowed"`
	Notes                string `json:"notes"`
}

// SIPAccountResponse represents a SIP account in API responses
type SIPAccountResponse struct {
	ID                 int64                          `json:"id"`
	CustomerID         int64                          `json:"customer_id"`
	AccountName        string                         `json:"account_name"`
	Username           string                         `json:"username"`
	Domain             string                         `json:"domain"`
	Extension          string                         `json:"extension"`
	CallerID           string                         `json:"caller_id"`
	CallerIDName       string                         `json:"caller_id_name"`
	Context            string                         `json:"context"`
	Transport          string                         `json:"transport"`
	NATSupport         bool                           `json:"nat_support"`
	DirectMediaSupport bool                           `json:"direct_media_support"`
	EncryptionEnabled  bool                           `json:"encryption_enabled"`
	CodecsAllowed      string                         `json:"codecs_allowed"`
	MaxConcurrentCalls int                            `json:"max_concurrent_calls"`
	CurrentActiveCalls int                            `json:"current_active_calls"`
	Status             string                         `json:"status"`
	IsRegistered       bool                           `json:"is_registered"`
	LastRegisteredIP   *string                        `json:"last_registered_ip"`
	LastRegisteredAt   *time.Time                     `json:"last_registered_at"`
	LastCallAt         *time.Time                     `json:"last_call_at"`
	TotalCalls         int64                          `json:"total_calls"`
	TotalMinutes       int64                          `json:"total_minutes"`
	Notes              *string                        `json:"notes"`
	Permissions        *models.SIPAccountPermission   `json:"permissions,omitempty"`
	MonthlyUsage       *models.SIPAccountUsage        `json:"monthly_usage,omitempty"`
	CreatedAt          time.Time                      `json:"created_at"`
	UpdatedAt          time.Time                      `json:"updated_at"`
}

// CreateSIPAccount handles POST /api/v1/customers/:customer_id/sip-accounts
func (h *SIPAccountHandler) CreateSIPAccount(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var req CreateSIPAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create SIP account
	account := &models.SIPAccount{
		AccountName:        req.AccountName,
		Extension:          req.Extension,
		CallerID:           req.CallerID,
		CallerIDName:       req.CallerIDName,
		MaxConcurrentCalls: req.MaxConcurrentCalls,
		NATSupport:         req.NATSupport,
		EncryptionEnabled:  req.EncryptionEnabled,
		CodecsAllowed:      req.CodecsAllowed,
	}

	if req.Notes != "" {
		account.Notes = &req.Notes
	}

	if account.MaxConcurrentCalls == 0 {
		account.MaxConcurrentCalls = 2
	}

	err = h.sipService.CreateSIPAccount(c.Request.Context(), customerID, account, userID.(int64))
	if err != nil {
		h.logger.WithError(err).Error("Failed to create SIP account")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SIP account"})
		return
	}

	// Get permissions for response
	permissions, _ := h.sipService.GetSIPAccountPermissions(c.Request.Context(), account.ID)

	response := h.toResponse(account)
	response.Permissions = permissions

	c.JSON(http.StatusCreated, response)
}

// GetSIPAccount handles GET /api/v1/sip-accounts/:id
func (h *SIPAccountHandler) GetSIPAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIP account ID"})
		return
	}

	account, err := h.sipService.GetSIPAccountByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SIP account not found"})
		return
	}

	// Get permissions
	permissions, _ := h.sipService.GetSIPAccountPermissions(c.Request.Context(), id)

	// Get current month usage
	monthlyUsage, _ := h.sipService.GetCurrentMonthUsage(c.Request.Context(), id)

	response := h.toResponse(account)
	response.Permissions = permissions
	response.MonthlyUsage = monthlyUsage

	c.JSON(http.StatusOK, response)
}

// GetCustomerSIPAccounts handles GET /api/v1/customers/:customer_id/sip-accounts
func (h *SIPAccountHandler) GetCustomerSIPAccounts(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	accounts, err := h.sipService.GetCustomerSIPAccounts(c.Request.Context(), customerID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get customer SIP accounts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get SIP accounts"})
		return
	}

	responses := make([]SIPAccountResponse, 0, len(accounts))
	for _, account := range accounts {
		responses = append(responses, h.toResponse(account))
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": responses,
		"count":    len(responses),
	})
}

// UpdateSIPAccount handles PUT /api/v1/sip-accounts/:id
func (h *SIPAccountHandler) UpdateSIPAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIP account ID"})
		return
	}

	var req struct {
		AccountName        string `json:"account_name"`
		Extension          string `json:"extension"`
		CallerID           string `json:"caller_id"`
		CallerIDName       string `json:"caller_id_name"`
		MaxConcurrentCalls int    `json:"max_concurrent_calls"`
		NATSupport         bool   `json:"nat_support"`
		EncryptionEnabled  bool   `json:"encryption_enabled"`
		CodecsAllowed      string `json:"codecs_allowed"`
		Status             string `json:"status"`
		Notes              string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing account
	account, err := h.sipService.GetSIPAccountByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SIP account not found"})
		return
	}

	// Update fields
	account.AccountName = req.AccountName
	account.Extension = req.Extension
	account.CallerID = req.CallerID
	account.CallerIDName = req.CallerIDName
	account.MaxConcurrentCalls = req.MaxConcurrentCalls
	account.NATSupport = req.NATSupport
	account.EncryptionEnabled = req.EncryptionEnabled
	account.CodecsAllowed = req.CodecsAllowed
	account.Status = req.Status
	
	if req.Notes != "" {
		account.Notes = &req.Notes
	}

	err = h.sipService.UpdateSIPAccount(c.Request.Context(), account)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update SIP account")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SIP account"})
		return
	}

	c.JSON(http.StatusOK, h.toResponse(account))
}

// DeleteSIPAccount handles DELETE /api/v1/sip-accounts/:id
func (h *SIPAccountHandler) DeleteSIPAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIP account ID"})
		return
	}

	err = h.sipService.DeleteSIPAccount(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete SIP account")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SIP account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SIP account deleted successfully"})
}

// UpdateSIPAccountPermissions handles PUT /api/v1/sip-accounts/:id/permissions
func (h *SIPAccountHandler) UpdateSIPAccountPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIP account ID"})
		return
	}

	var permissions models.SIPAccountPermission
	if err := c.ShouldBindJSON(&permissions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissions.SIPAccountID = id

	err = h.sipService.UpdateSIPAccountPermissions(c.Request.Context(), &permissions)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update SIP account permissions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions updated successfully"})
}

// GetSIPAccountUsage handles GET /api/v1/sip-accounts/:id/usage
func (h *SIPAccountHandler) GetSIPAccountUsage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIP account ID"})
		return
	}

	// Parse date range from query params
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	usage, err := h.sipService.GetUsageStats(c.Request.Context(), id, startDate, endDate)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get SIP account usage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get usage statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usage":      usage,
		"start_date": startDate,
		"end_date":   endDate,
	})
}

// GenerateCredentials handles POST /api/v1/sip-accounts/:id/generate-credentials
func (h *SIPAccountHandler) GenerateCredentials(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIP account ID"})
		return
	}

	// Get existing account
	account, err := h.sipService.GetSIPAccountByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SIP account not found"})
		return
	}

	// Generate new password
	newPassword, err := h.sipService.GenerateSecurePassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password"})
		return
	}

	account.Password = newPassword

	// Note: In a real implementation, you'd update only the password field
	// This is simplified for demonstration
	err = h.sipService.UpdateSIPAccount(c.Request.Context(), account)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update SIP account password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Credentials updated successfully",
		"username": account.Username,
		"password": newPassword,
		"domain":   account.Domain,
		"sip_uri":  account.GetSIPURI(),
	})
}

// Helper function to convert model to response
func (h *SIPAccountHandler) toResponse(account *models.SIPAccount) SIPAccountResponse {
	return SIPAccountResponse{
		ID:                 account.ID,
		CustomerID:         account.CustomerID,
		AccountName:        account.AccountName,
		Username:           account.Username,
		Domain:             account.Domain,
		Extension:          account.Extension,
		CallerID:           account.CallerID,
		CallerIDName:       account.CallerIDName,
		Context:            account.Context,
		Transport:          account.Transport,
		NATSupport:         account.NATSupport,
		DirectMediaSupport: account.DirectMediaSupport,
		EncryptionEnabled:  account.EncryptionEnabled,
		CodecsAllowed:      account.CodecsAllowed,
		MaxConcurrentCalls: account.MaxConcurrentCalls,
		CurrentActiveCalls: account.CurrentActiveCalls,
		Status:             account.Status,
		IsRegistered:       account.IsRegistered(),
		LastRegisteredIP:   account.LastRegisteredIP,
		LastRegisteredAt:   account.LastRegisteredAt,
		LastCallAt:         account.LastCallAt,
		TotalCalls:         account.TotalCalls,
		TotalMinutes:       account.TotalMinutes,
		Notes:              account.Notes,
		CreatedAt:          account.CreatedAt,
		UpdatedAt:          account.UpdatedAt,
	}
}