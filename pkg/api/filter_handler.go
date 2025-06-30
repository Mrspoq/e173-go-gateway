package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/e173-gateway/e173_go_gateway/pkg/service"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type FilterHandler struct {
	filterService service.FilterService
}

func NewFilterHandler(filterService service.FilterService) *FilterHandler {
	return &FilterHandler{
		filterService: filterService,
	}
}

type FilterRequest struct {
	Source      string `json:"source" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	GatewayID   string `json:"gateway_id"`
}

type FilterResponse struct {
	Action    string `json:"action"` // route, reject, blackhole
	GatewayID string `json:"gateway_id,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
}

// CheckCall handles POST /api/v1/filter/check
func (h *FilterHandler) CheckCall(c *gin.Context) {
	var req FilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create call object for filtering
	call := &models.Call{
		SourceNumber: req.Source,
		DestNumber:   req.Destination,
		GatewayID:    req.GatewayID,
	}

	// Run through filter engine
	result, err := h.filterService.ProcessCall(call)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare response
	response := FilterResponse{
		Action:    result.Action,
		GatewayID: result.GatewayID,
		Reason:    result.Reason,
		Prefix:    result.Prefix,
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers all filter-related routes
func (h *FilterHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/filter/check", h.CheckCall)
}