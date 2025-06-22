package repository

import (
	"context"
	"errors"

	"github.com/e173-gateway/e173_go_gateway/pkg/models" // Ensuring this import is correct for models.Cdr, models.Modem, models.SIMCard
	"github.com/google/uuid"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("repository: record not found")

// CdrRepository defines the interface for interacting with CDR data.
type CdrRepository interface {
	CreateCdr(ctx context.Context, cdr *models.Cdr) error // Ensure models.Cdr is resolvable
	GetCdrByID(ctx context.Context, id uuid.UUID) (*models.Cdr, error) // Ensure models.Cdr is resolvable
	GetRecentCDRs(ctx context.Context, limit int) ([]*models.Cdr, error) // Get recent CDRs for dashboard
	// Potentially: GetCdrByAsteriskUniqueID(ctx context.Context, asteriskUniqueID string) (*models.Cdr, error)
	// ListCdrs(ctx context.Context, // filters, pagination) ([]*models.Cdr, error)
}

// GatewayRepository defines the interface for gateway-related database operations.
type GatewayRepository interface {
	CreateGateway(ctx context.Context, gateway *models.Gateway) error
	GetGatewayByID(ctx context.Context, id string) (*models.Gateway, error)
	ListGateways(ctx context.Context) ([]*models.Gateway, error)
	UpdateGateway(ctx context.Context, gateway *models.Gateway) error
	UpdateGatewayHeartbeat(ctx context.Context, id string) error
	DeleteGateway(ctx context.Context, id string) error
	GetGatewayStats(ctx context.Context) (total, online, offline int, err error)
}

// ModemRepository defines the interface for interacting with modem data.
type ModemRepository interface {
	CreateModem(ctx context.Context, modem *models.Modem) (int, error)
	GetAllModems(ctx context.Context) ([]models.Modem, error)
	GetModemByID(ctx context.Context, id int) (*models.Modem, error) // Kept, will be implemented
	// UpdateModem(ctx context.Context, modem *models.Modem) error
	// DeleteModem(ctx context.Context, id int) error
}

// SIMCardRepository defines the interface for interacting with SIM card data.
type SIMCardRepository interface {
	CreateSIMCard(ctx context.Context, sim *models.SIMCard) (int64, error)
	GetSIMCardByID(ctx context.Context, id int64) (*models.SIMCard, error) // Changed id to int64
	GetSIMCardByICCID(ctx context.Context, iccid string) (*models.SIMCard, error)
	GetAllSIMCards(ctx context.Context) ([]models.SIMCard, error)
	UpdateSIMCard(ctx context.Context, sim *models.SIMCard) error
	DeleteSIMCard(ctx context.Context, id int64) error
	// Potentially more specific methods like:
	// GetSIMCardsByStatus(ctx context.Context, status string) ([]models.SIMCard, error)
	// GetSIMCardsByModemID(ctx context.Context, modemID int64) ([]models.SIMCard, error)
}
