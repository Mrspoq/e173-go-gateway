package api

import (
	"net/http"

	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/e173-gateway/e173_go_gateway/pkg/repository"
	"github.com/e173-gateway/e173_go_gateway/pkg/logging"
	"github.com/gin-gonic/gin"
	"strconv"
	"errors"
)

// ErrorResponse is a generic JSON error response structure.
type ErrorResponse struct {
	Error string `json:"error"`
}

type SIMCardHandler struct {
	repo repository.SIMCardRepository
}

func NewSIMCardHandler(repo repository.SIMCardRepository) *SIMCardHandler {
	return &SIMCardHandler{repo: repo}
}

// CreateSIMCard godoc
// @Summary Create a new SIM card
// @Description Add a new SIM card to the system
// @Tags sim_cards
// @Accept  json
// @Produce  json
// @Param   sim_card  body   models.SIMCard  true  "SIM Card Object"
// @Success 201 {object} models.SIMCard "Successfully created SIM card"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /simcards [post]
func (h *SIMCardHandler) CreateSIMCard(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.Logger.WithContext(ctx)
	var sim models.SIMCard

	if err := c.ShouldBindJSON(&sim); err != nil {
		logger.WithError(err).Error("Failed to bind SIM card JSON")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Basic validation (can be expanded)
	if sim.ICCID == "" {
		logger.Error("ICCID is required")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ICCID is required"})
		return
	}
	
	// Default status if not provided, or validate allowed statuses
	if sim.Status == "" {
		sim.Status = "unknown" 
	}


	id, err := h.repo.CreateSIMCard(ctx, &sim)
	if err != nil {
		logger.WithError(err).Error("Failed to create SIM card in repository")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create SIM card"})
		return
	}

	sim.ID = id // Set the ID on the response object
	// We might want to fetch the full sim object again to get created_at, updated_at
	// For now, returning the input sim with the new ID.
	
	logger.WithField("sim_id", id).Info("SIM card created successfully via API")
	c.JSON(http.StatusCreated, sim)
}

// GetSIMCardByID godoc
// @Summary Get a SIM card by ID
// @Description Retrieve details of a specific SIM card using its ID
// @Tags sim_cards
// @Produce  json
// @Param   id   path   int  true  "SIM Card ID"
// @Success 200 {object} models.SIMCard "Successfully retrieved SIM card"
// @Failure 400 {object} ErrorResponse "Invalid SIM Card ID format"
// @Failure 404 {object} ErrorResponse "SIM Card not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /simcards/{id} [get]
func (h *SIMCardHandler) GetSIMCardByID(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.Logger.WithContext(ctx)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.WithError(err).Error("Invalid SIM card ID format")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid SIM Card ID format"})
		return
	}

	sim, err := h.repo.GetSIMCardByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.WithField("sim_id", id).Warn("SIM card not found")
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "SIM Card not found"})
		} else {
			logger.WithError(err).WithField("sim_id", id).Error("Failed to get SIM card from repository")
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve SIM card"})
		}
		return
	}

	logger.WithField("sim_id", id).Info("SIM card retrieved successfully via API")
	c.JSON(http.StatusOK, sim)
}

// GetAllSIMCards godoc
// @Summary Get all SIM cards
// @Description Retrieve a list of all SIM cards in the system
// @Tags sim_cards
// @Produce  json
// @Success 200 {array} models.SIMCard "Successfully retrieved list of SIM cards"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /simcards [get]
func (h *SIMCardHandler) GetAllSIMCards(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.Logger.WithContext(ctx)

	sims, err := h.repo.GetAllSIMCards(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to get all SIM cards from repository")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve SIM cards"})
		return
	}

	if sims == nil { // Ensure we return an empty array, not null, if no SIMs are found
		sims = []models.SIMCard{}
	}

	logger.Info("All SIM cards retrieved successfully via API")
	c.JSON(http.StatusOK, sims)
}

// UpdateSIMCard godoc
// @Summary Update an existing SIM card
// @Description Update details of an existing SIM card by its ID
// @Tags sim_cards
// @Accept  json
// @Produce  json
// @Param   id        path   int             true  "SIM Card ID"
// @Param   sim_card  body   models.SIMCard  true  "SIM Card Object with updated fields"
// @Success 200 {object} models.SIMCard "Successfully updated SIM card"
// @Failure 400 {object} ErrorResponse "Invalid request payload or SIM Card ID format"
// @Failure 404 {object} ErrorResponse "SIM Card not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /simcards/{id} [put]
func (h *SIMCardHandler) UpdateSIMCard(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.Logger.WithContext(ctx)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.WithError(err).Error("Invalid SIM card ID format for update")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid SIM Card ID format"})
		return
	}

	var simUpdates models.SIMCard
	if err := c.ShouldBindJSON(&simUpdates); err != nil {
		logger.WithError(err).Error("Failed to bind SIM card JSON for update")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Ensure the ID in the path matches the ID in the body if provided, or set it
	// For a PUT, the ID in the path is authoritative.
	simUpdates.ID = id

	// Basic validation (can be expanded)
	if simUpdates.ICCID == "" { // ICCID might be part of an update, ensure it's not being blanked if it's a required field.
		// Depending on business logic, you might disallow changing ICCID or require it.
		// For now, we assume it could be updated but shouldn't be set to empty if it's a key identifier.
		// If ICCID is immutable, this check should be different or removed if not updatable.
	}


	err = h.repo.UpdateSIMCard(ctx, &simUpdates)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.WithField("sim_id", id).Warn("SIM card not found for update")
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "SIM Card not found"})
		} else {
			logger.WithError(err).WithField("sim_id", id).Error("Failed to update SIM card in repository")
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update SIM card"})
		}
		return
	}
	
	// Fetch the updated record to return the full object with all fields (including any server-set fields like updated_at)
	updatedSim, err := h.repo.GetSIMCardByID(ctx, id)
	if err != nil {
		// Log this error, but we can still return the simUpdates object as a fallback
		logger.WithError(err).WithField("sim_id", id).Error("Failed to retrieve updated SIM card after update, returning input data")
		c.JSON(http.StatusOK, simUpdates) // Fallback to returning the input data if retrieval fails
		return
	}

	logger.WithField("sim_id", id).Info("SIM card updated successfully via API")
	c.JSON(http.StatusOK, updatedSim)
}

// DeleteSIMCard godoc
// @Summary Delete a SIM card by ID
// @Description Remove a SIM card from the system using its ID
// @Tags sim_cards
// @Produce  json
// @Param   id   path   int  true  "SIM Card ID"
// @Success 204 "Successfully deleted SIM card (No Content)"
// @Failure 400 {object} ErrorResponse "Invalid SIM Card ID format"
// @Failure 404 {object} ErrorResponse "SIM Card not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /simcards/{id} [delete]
func (h *SIMCardHandler) DeleteSIMCard(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.Logger.WithContext(ctx)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.WithError(err).Error("Invalid SIM card ID format for delete")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid SIM Card ID format"})
		return
	}

	err = h.repo.DeleteSIMCard(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.WithField("sim_id", id).Warn("SIM card not found for delete")
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "SIM Card not found"})
		} else {
			logger.WithError(err).WithField("sim_id", id).Error("Failed to delete SIM card from repository")
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete SIM card"})
		}
		return
	}

	logger.WithField("sim_id", id).Info("SIM card deleted successfully via API")
	c.Status(http.StatusNoContent) // 204 No Content for successful DELETE
}
