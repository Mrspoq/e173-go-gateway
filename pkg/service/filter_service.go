package service

import (
	"fmt"
	"strings"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/e173-gateway/e173_go_gateway/pkg/repository"
	"github.com/e173-gateway/e173_go_gateway/pkg/validation"
)

type FilterService interface {
	ProcessCall(call *models.Call) (*FilterResult, error)
}

type FilterResult struct {
	Action    string // route, reject, blackhole
	GatewayID string
	Reason    string
	Prefix    string
}

type filterService struct {
	blacklistRepo      repository.BlacklistRepository
	prefixRepo         repository.PrefixRepository
	whatsappValidator  validation.WhatsAppValidator
	phoneValidator     validation.PhoneNumberValidator
}

func NewFilterService(
	blacklistRepo repository.BlacklistRepository,
	prefixRepo repository.PrefixRepository,
	whatsappValidator validation.WhatsAppValidator,
	phoneValidator validation.PhoneNumberValidator,
) FilterService {
	return &filterService{
		blacklistRepo:     blacklistRepo,
		prefixRepo:        prefixRepo,
		whatsappValidator: whatsappValidator,
		phoneValidator:    phoneValidator,
	}
}

func (s *filterService) ProcessCall(call *models.Call) (*FilterResult, error) {
	// 1. Check if source number is blacklisted
	blacklisted, err := s.blacklistRepo.IsBlacklisted(call.SourceNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if blacklisted {
		return &FilterResult{
			Action: "blackhole",
			Reason: "Source number is blacklisted",
		}, nil
	}

	// 2. Source number validation is minimal - just check it's not empty
	// Source can be anonymous, private, international, etc.
	if call.SourceNumber == "" {
		return &FilterResult{
			Action: "reject",
			Reason: "Empty source number",
		}, nil
	}

	// 3. Validate DESTINATION number format using libphonenumber
	// This is critical because we're terminating calls to Morocco
	if !s.phoneValidator.IsValid(call.DestNumber) {
		return &FilterResult{
			Action: "reject",
			Reason: "Invalid destination number format (failed libphonenumber validation)",
		}, nil
	}

	// 3a. For Morocco DESTINATION numbers, ensure exact length (9 digits after country code)
	// We only validate Morocco numbers since we're a Morocco termination gateway
	destCleaned := strings.TrimPrefix(strings.TrimPrefix(call.DestNumber, "+"), "00")
	if !strings.HasPrefix(destCleaned, "212") {
		return &FilterResult{
			Action: "reject", 
			Reason: "Non-Morocco destination number (this gateway only terminates to Morocco)",
		}, nil
	}
	
	// Morocco number must be exactly 12 digits total
	if len(destCleaned) != 12 {
		return &FilterResult{
			Action: "reject",
			Reason: fmt.Sprintf("Invalid Morocco number length: expected 12 digits, got %d", len(destCleaned)),
		}, nil
	}

	// 4. Both numbers are valid per libphonenumber, now check WhatsApp for destination
	// This ensures we only call the paid WhatsApp API for valid numbers
	whatsappStatus, err := s.whatsappValidator.ValidateNumber(call.DestNumber)
	if err != nil {
		// Log error but continue - WhatsApp check might be temporarily unavailable
		fmt.Printf("WhatsApp validation error for %s: %v\n", call.DestNumber, err)
		// Optionally, you might want to reject or allow based on configuration
		// For now, we'll continue to route if WhatsApp check fails
	} else if whatsappStatus != nil && !whatsappStatus.HasWhatsApp {
		return &FilterResult{
			Action: "reject",
			Reason: "Destination number not active on WhatsApp",
		}, nil
	}

	// 5. Future: HLR check would go here for additional validation
	// if s.hlrValidator != nil {
	//     hlrStatus, err := s.hlrValidator.ValidateNumber(call.DestNumber)
	//     ...
	// }

	// 6. Find matching prefix and gateway
	prefix, gateway, err := s.findBestPrefixMatch(call.DestNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to find prefix match: %w", err)
	}

	if prefix == nil || gateway == nil {
		return &FilterResult{
			Action: "reject",
			Reason: "No route found for destination",
		}, nil
	}

	// 7. All checks passed - route the call
	return &FilterResult{
		Action:    "route",
		GatewayID: gateway.ID,
		Prefix:    prefix.Prefix,
	}, nil
}

func (s *filterService) findBestPrefixMatch(number string) (*models.Prefix, *models.Gateway, error) {
	// Clean the number (remove + if present)
	cleanNumber := strings.TrimPrefix(number, "+")

	// Get all active prefixes
	prefixes, err := s.prefixRepo.GetAllActive()
	if err != nil {
		return nil, nil, err
	}

	// Find the longest matching prefix
	var bestMatch *models.Prefix
	maxLength := 0

	for _, prefix := range prefixes {
		if strings.HasPrefix(cleanNumber, prefix.Prefix) && len(prefix.Prefix) > maxLength {
			bestMatch = &prefix
			maxLength = len(prefix.Prefix)
		}
	}

	if bestMatch == nil {
		return nil, nil, nil
	}

	// For now, return the prefix with a dummy gateway
	// In production, this should select the best gateway based on load, cost, etc.
	gateway := &models.Gateway{
		ID: bestMatch.GatewayID,
	}

	return bestMatch, gateway, nil
}