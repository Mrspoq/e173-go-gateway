package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/e173-gateway/e173_go_gateway/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RechargeHandler struct {
	rechargeRepo repository.RechargeRepository
	simCardRepo  repository.SIMCardRepository
	logger       *logrus.Entry
}

func NewRechargeHandler(rechargeRepo repository.RechargeRepository, simCardRepo repository.SIMCardRepository, logger *logrus.Entry) *RechargeHandler {
	return &RechargeHandler{
		rechargeRepo: rechargeRepo,
		simCardRepo:  simCardRepo,
		logger:       logger,
	}
}

// CreateRechargeCode handles POST /api/v1/recharge/codes
func (h *RechargeHandler) CreateRechargeCode(c *gin.Context) {
	var code models.RechargeCode
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	code.Status = "active"
	code.CreatedAt = time.Now()

	if err := h.rechargeRepo.CreateRechargeCode(ctx, &code); err != nil {
		h.logger.WithError(err).Error("Failed to create recharge code")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recharge code"})
		return
	}

	c.JSON(http.StatusCreated, code)
}

// RechargeSimCard handles POST /api/v1/sims/:id/recharge
func (h *RechargeHandler) RechargeSimCard(c *gin.Context) {
	simIDStr := c.Param("id")
	simID, err := strconv.ParseInt(simIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIM card ID"})
		return
	}

	var req struct {
		Code     string  `json:"code" binding:"required"`
		Amount   float64 `json:"amount,omitempty"`
		Operator string  `json:"operator,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Get SIM card
	simCard, err := h.simCardRepo.GetSIMCardByID(ctx, simID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "SIM card not found"})
			return
		}
		h.logger.WithError(err).Error("Failed to get SIM card")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get SIM card"})
		return
	}

	// If code is provided, validate it
	var operator string
	if simCard.OperatorName.Valid {
		operator = simCard.OperatorName.String
	}
	
	if req.Code != "" {
		rechargeCode, err := h.rechargeRepo.GetRechargeCodeByCode(ctx, req.Code, operator)
		if err != nil {
			if err == repository.ErrNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recharge code"})
				return
			}
			h.logger.WithError(err).Error("Failed to validate recharge code")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate recharge code"})
			return
		}

		if rechargeCode.Status != "active" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Recharge code is not active"})
			return
		}

		req.Amount = rechargeCode.Amount
	}

	// Update SIM card balance
	currentBalance := 0.0
	if simCard.Balance.Valid {
		currentBalance = simCard.Balance.Float64
	}
	newBalance := currentBalance + req.Amount
	simCard.Balance = sql.NullFloat64{Float64: newBalance, Valid: true}
	
	if err := h.simCardRepo.UpdateSIMCard(ctx, simCard); err != nil {
		h.logger.WithError(err).Error("Failed to update SIM card balance")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	// Get phone number from SIM card
	phoneNumber := ""
	if simCard.MSISDN.Valid {
		phoneNumber = simCard.MSISDN.String
	}
	
	// Create recharge history
	history := &models.RechargeHistory{
		SimCardID:     simCard.ID,
		PhoneNumber:   phoneNumber,
		Amount:        req.Amount,
		BalanceBefore: sql.NullFloat64{Float64: currentBalance, Valid: true},
		BalanceAfter:  sql.NullFloat64{Float64: newBalance, Valid: true},
		Method:        "api",
		Status:        "completed",
		Attempts:      1,
		ProcessedBy:   1, // TODO: Get from current user
		ProcessedAt:   time.Now(),
	}

	if err := h.rechargeRepo.CreateRechargeHistory(ctx, history); err != nil {
		h.logger.WithError(err).Error("Failed to create recharge history")
		// Don't fail the request, history is not critical
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Recharge successful",
		"balance": newBalance,
		"history": history,
	})
}

// GetRechargeHistory handles GET /api/v1/sims/:id/recharge/history
func (h *RechargeHandler) GetRechargeHistory(c *gin.Context) {
	simIDStr := c.Param("id")
	simID, err := strconv.ParseInt(simIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIM card ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	history, err := h.rechargeRepo.GetRechargeHistory(ctx, simID, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get recharge history")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recharge history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetRechargeUI handles GET /sims/:id/recharge
func (h *RechargeHandler) GetRechargeUI(c *gin.Context) {
	simIDStr := c.Param("id")
	simID, err := strconv.ParseInt(simIDStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "Invalid SIM card ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Get SIM card
	simCard, err := h.simCardRepo.GetSIMCardByID(ctx, simID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.HTML(http.StatusNotFound, "error.tmpl", gin.H{"error": "SIM card not found"})
			return
		}
		h.logger.WithError(err).Error("Failed to get SIM card")
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Failed to get SIM card"})
		return
	}

	// Get recharge history
	history, err := h.rechargeRepo.GetRechargeHistory(ctx, simID, 10)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get recharge history")
		history = []*models.RechargeHistory{} // Empty on error
	}

	templateData := gin.H{
		"title":   "Recharge SIM Card",
		"SIMCard": simCard,
		"History": history,
	}

	// Add current user if available
	if user, exists := c.Get("currentUser"); exists {
		templateData["CurrentUser"] = user
	}

	c.HTML(http.StatusOK, "sims/recharge.html", templateData)
}

// AutoRecharge handles automatic recharge for low balance SIMs
func (h *RechargeHandler) AutoRecharge(threshold float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get SIM cards with low balance
	sims, err := h.rechargeRepo.GetSimCardsForAutoRecharge(ctx, threshold)
	if err != nil {
		return err
	}

	h.logger.Infof("Found %d SIM cards for auto-recharge", len(sims))

	for _, sim := range sims {
		// TODO: Implement auto-recharge logic
		// This could involve:
		// 1. Checking if auto-recharge is enabled for the SIM
		// 2. Getting the configured recharge amount
		// 3. Processing the recharge
		// 4. Sending notifications
		
		h.logger.Infof("Auto-recharge needed for SIM %s (balance: %.2f)", sim.ICCID, sim.Balance)
	}

	return nil
}