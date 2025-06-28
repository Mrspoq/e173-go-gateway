package api

import (
	"context"
	"net/http"
	"time"

	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/e173-gateway/e173_go_gateway/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GatewayHandler handles gateway-related HTTP requests.
type GatewayHandler struct {
	gatewayRepo repository.GatewayRepository
	logger      *logrus.Logger
}

// NewGatewayHandler creates a new instance of GatewayHandler.
func NewGatewayHandler(gatewayRepo repository.GatewayRepository, logger *logrus.Logger) *GatewayHandler {
	return &GatewayHandler{
		gatewayRepo: gatewayRepo,
		logger:      logger,
	}
}

// CreateGateway handles POST /api/v1/gateways
func (h *GatewayHandler) CreateGateway(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Location    string `json:"location"`
		AMIHost     string `json:"ami_host" binding:"required"`
		AMIPort     string `json:"ami_port"`
		AMIUser     string `json:"ami_user" binding:"required"`
		AMIPass     string `json:"ami_pass" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default port if not provided
	if req.AMIPort == "" {
		req.AMIPort = "5038"
	}

	gateway := &models.Gateway{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		AMIHost:     req.AMIHost,
		AMIPort:     req.AMIPort,
		AMIUser:     req.AMIUser,
		AMIPass:     req.AMIPass,
		Status:      models.GatewayStatusOffline,
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.gatewayRepo.CreateGateway(ctx, gateway); err != nil {
		h.logger.WithError(err).Error("Failed to create gateway")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gateway"})
		return
	}

	c.JSON(http.StatusCreated, gateway)
}

// GetGatewayByID handles GET /api/v1/gateways/:id
func (h *GatewayHandler) GetGatewayByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gateway ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	gateway, err := h.gatewayRepo.GetGatewayByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gateway not found"})
			return
		}
		h.logger.WithError(err).Error("Failed to get gateway")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gateway"})
		return
	}

	c.JSON(http.StatusOK, gateway)
}

// ListGateways handles GET /api/v1/gateways
func (h *GatewayHandler) ListGateways(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	gateways, err := h.gatewayRepo.ListGateways(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list gateways")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list gateways"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"gateways": gateways})
}

// UpdateGateway handles PUT /api/v1/gateways/:id
func (h *GatewayHandler) UpdateGateway(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gateway ID is required"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
		AMIHost     string `json:"ami_host"`
		AMIPort     string `json:"ami_port"`
		AMIUser     string `json:"ami_user"`
		AMIPass     string `json:"ami_pass"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Get existing gateway
	gateway, err := h.gatewayRepo.GetGatewayByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gateway not found"})
			return
		}
		h.logger.WithError(err).Error("Failed to get gateway")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gateway"})
		return
	}

	// Update fields
	if req.Name != "" {
		gateway.Name = req.Name
	}
	if req.Description != "" {
		gateway.Description = req.Description
	}
	if req.Location != "" {
		gateway.Location = req.Location
	}
	if req.AMIHost != "" {
		gateway.AMIHost = req.AMIHost
	}
	if req.AMIPort != "" {
		gateway.AMIPort = req.AMIPort
	}
	if req.AMIUser != "" {
		gateway.AMIUser = req.AMIUser
	}
	if req.AMIPass != "" {
		gateway.AMIPass = req.AMIPass
	}

	// Save updates
	if err := h.gatewayRepo.UpdateGateway(ctx, gateway); err != nil {
		h.logger.WithError(err).Error("Failed to update gateway")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gateway"})
		return
	}

	c.JSON(http.StatusOK, gateway)
}

// DeleteGateway handles DELETE /api/v1/gateways/:id
func (h *GatewayHandler) DeleteGateway(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gateway ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.gatewayRepo.DeleteGateway(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gateway not found"})
			return
		}
		h.logger.WithError(err).Error("Failed to delete gateway")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete gateway"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gateway deleted successfully"})
}

// Heartbeat handles POST /api/v1/gateways/heartbeat
// This endpoint is called by gateways to report their status
func (h *GatewayHandler) Heartbeat(c *gin.Context) {
	var req struct {
		GatewayID string `json:"gateway_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Verify gateway exists
	gateway, err := h.gatewayRepo.GetGatewayByID(ctx, req.GatewayID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gateway not found"})
			return
		}
		h.logger.WithError(err).Error("Failed to get gateway")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gateway"})
		return
	}

	// Update heartbeat
	if err := h.gatewayRepo.UpdateGatewayHeartbeat(ctx, gateway.ID); err != nil {
		h.logger.WithError(err).Error("Failed to update gateway heartbeat")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update heartbeat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Heartbeat received",
		"gateway": gin.H{
			"id":   gateway.ID,
			"name": gateway.Name,
		},
	})
}

// GetGatewayStats handles GET /api/stats/gateways
// Returns an HTMX-compatible HTML fragment for the dashboard
func (h *GatewayHandler) GetGatewayStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	total, online, _, err := h.gatewayRepo.GetGatewayStats(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get gateway stats")
		// Return a default card on error
		total, online = 0, 0
	}

	// Return HTMX-compatible HTML fragment
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
        <div id="gateway-stats-card" hx-get="/api/stats/gateways" hx-trigger="every 10s" hx-swap="outerHTML" class="dashboard-card">
                <div class="flex items-center">
                        <div class="flex-shrink-0">
                                <div class="w-8 h-8 bg-blue-100 rounded-md flex items-center justify-center">
                                        <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"></path>
                                        </svg>
                                </div>
                        </div>
                        <div class="ml-5 w-0 flex-1">
                                <dl>
                                        <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Gateways</dt>
                                        <dd class="flex items-baseline">
                                                <div class="text-2xl font-semibold text-gray-900 dark:text-white">%d / %d</div>
                                                <div class="ml-2 text-sm text-gray-500 dark:text-gray-400">online</div>
                                        </dd>
                                </dl>
                        </div>
                </div>
        </div>`, online, total)
}

// Gateway UI Handlers

// GetGatewayListUI handles GET /gateways
// Returns the gateway management page
func (h *GatewayHandler) GetGatewayListUI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	gateways, err := h.gatewayRepo.ListGateways(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list gateways")
		gateways = []*models.Gateway{} // Empty list on error
	}

	c.HTML(http.StatusOK, "gateways/list.html", gin.H{
		"title":    "Gateway Management",
		"Gateways": gateways,
	})
}

// GetGatewayCreateUI handles GET /gateways/create
// Returns the gateway creation form
func (h *GatewayHandler) GetGatewayCreateUI(c *gin.Context) {
	c.HTML(http.StatusOK, "gateways/create.html", gin.H{
		"title": "Create Gateway",
	})
}

// GetGatewayEditUI handles GET /gateways/:id/edit
// Returns the gateway edit form
func (h *GatewayHandler) GetGatewayEditUI(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "Error",
			"Error": "Invalid gateway ID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	gateway, err := h.gatewayRepo.GetGatewayByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"Title": "Error",
				"Error": "Gateway not found",
			})
			return
		}
		h.logger.WithError(err).Error("Failed to get gateway")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "Error",
			"Error": "Failed to get gateway",
		})
		return
	}

	c.HTML(http.StatusOK, "gateways/edit.html", gin.H{
		"title":   "Edit Gateway",
		"Gateway": gateway,
	})
}
