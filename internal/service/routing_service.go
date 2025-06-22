package service

import (
	"fmt"
	"time"

	"github.com/e173-gateway/e173_go_gateway/internal/repository"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type RoutingService interface {
	// Call Routing
	RouteCall(callerNumber, destinationNumber string, customerID *int64) (*models.CallRoutingResult, error)
	GetRoutingRules(limit, offset int) ([]*models.RoutingRule, error)
	CreateRoutingRule(rule *models.RoutingRule, createdBy int64) error
	UpdateRoutingRule(rule *models.RoutingRule, updatedBy int64) error
	DeleteRoutingRule(id int64, deletedBy int64) error
	GetRoutingRuleByID(id int64) (*models.RoutingRule, error)
	
	// Blacklist Management
	CheckNumberBlacklisted(number string, direction string) (*models.Blacklist, error)
	AddToBlacklist(entry *models.Blacklist, addedBy int64) error
	RemoveFromBlacklist(id int64, removedBy int64) error
	GetBlacklistEntries(limit, offset int) ([]*models.Blacklist, error)
	GetBlacklistEntryByID(id int64) (*models.Blacklist, error)
	UpdateBlacklistEntry(entry *models.Blacklist, updatedBy int64) error
	
	// Spam Detection
	DetectSpam(callerNumber, destinationNumber string, callDurationSeconds int) (bool, string, error)
	AutoBlacklistNumber(number string, reason string, detectionMethod string) error
	GetSpamStatistics() (*SpamStats, error)
	
	// SIM Pool Management
	GetAvailableSIMForCall(poolName string, customerID *int64) (*models.SIMPoolAssignment, error)
	CreateSIMPool(pool *models.SIMPool, createdBy int64) error
	AssignSIMToPool(simPoolID, simCardID int64, priority int, assignedBy int64) error
	RemoveSIMFromPool(simPoolID, simCardID int64, removedBy int64) error
	GetSIMPools(limit, offset int) ([]*models.SIMPool, error)
	GetSIMsInPool(poolName string) ([]*models.SIMPoolAssignment, error)
}

type SpamStats struct {
	TotalBlacklistedNumbers int64 `json:"total_blacklisted_numbers"`
	AutoBlacklistedToday    int64 `json:"auto_blacklisted_today"`
	BlockedCallsToday       int64 `json:"blocked_calls_today"`
	TopSpamPrefixes         []SpamPrefixStat `json:"top_spam_prefixes"`
}

type SpamPrefixStat struct {
	Prefix      string `json:"prefix"`
	BlockedCalls int64  `json:"blocked_calls"`
}

type PostgresRoutingService struct {
	routingRepo repository.RoutingRepository
	systemRepo  repository.SystemRepository
}

func NewPostgresRoutingService(
	routingRepo repository.RoutingRepository,
	systemRepo repository.SystemRepository,
) RoutingService {
	return &PostgresRoutingService{
		routingRepo: routingRepo,
		systemRepo:  systemRepo,
	}
}

func (s *PostgresRoutingService) RouteCall(callerNumber, destinationNumber string, customerID *int64) (*models.CallRoutingResult, error) {
	result := &models.CallRoutingResult{
		Success: false,
	}
	
	// Check if caller is blacklisted
	blacklistEntry, err := s.routingRepo.CheckNumberBlacklisted(callerNumber)
	if err != nil {
		result.ErrorMessage = stringPtr(fmt.Sprintf("Failed to check blacklist: %v", err))
		return result, nil
	}
	
	if blacklistEntry != nil && blacklistEntry.ShouldBlock("outbound") {
		result.IsBlocked = true
		result.BlockReason = stringPtr(fmt.Sprintf("Caller number %s is blacklisted: %s", callerNumber, *blacklistEntry.Reason))
		return result, nil
	}
	
	// Check if destination is blacklisted
	blacklistEntry, err = s.routingRepo.CheckNumberBlacklisted(destinationNumber)
	if err != nil {
		result.ErrorMessage = stringPtr(fmt.Sprintf("Failed to check destination blacklist: %v", err))
		return result, nil
	}
	
	if blacklistEntry != nil && blacklistEntry.ShouldBlock("inbound") {
		result.IsBlocked = true
		result.BlockReason = stringPtr(fmt.Sprintf("Destination number %s is blacklisted: %s", destinationNumber, *blacklistEntry.Reason))
		return result, nil
	}
	
	// Get routing rules for the destination number
	rules, err := s.routingRepo.GetRoutingRulesForNumber(destinationNumber)
	if err != nil {
		result.ErrorMessage = stringPtr(fmt.Sprintf("Failed to get routing rules: %v", err))
		return result, nil
	}
	
	// Find the best matching rule
	var selectedRule *models.RoutingRule
	for _, rule := range rules {
		// Check customer restrictions
		if len(rule.CustomerRestrictions) > 0 && customerID != nil {
			allowed := false
			for _, allowedCustomerID := range rule.CustomerRestrictions {
				if allowedCustomerID == *customerID {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}
		
		// Check time restrictions (simplified - would need proper JSON parsing in production)
		// For now, just use the first matching rule
		selectedRule = rule
		break
	}
	
	if selectedRule == nil {
		result.ErrorMessage = stringPtr("No routing rule found for destination number")
		return result, nil
	}
	
	result.Success = true
	result.RoutingRuleID = &selectedRule.ID
	result.CostMarkup = selectedRule.CostMarkupPercent
	
	// Route to specific modem or pool
	if selectedRule.RouteToModemID != nil {
		result.RouteToModemID = selectedRule.RouteToModemID
	} else if selectedRule.RouteToPool != nil {
		result.RouteToPool = selectedRule.RouteToPool
		
		// Get available SIM from pool
		simAssignment, err := s.GetAvailableSIMForCall(*selectedRule.RouteToPool, customerID)
		if err != nil {
			result.Success = false
			result.ErrorMessage = stringPtr(fmt.Sprintf("Failed to get SIM from pool: %v", err))
			return result, nil
		}
		
		if simAssignment != nil {
			result.SelectedSIMID = &simAssignment.SIMCardID
		}
	}
	
	return result, nil
}

func (s *PostgresRoutingService) GetRoutingRules(limit, offset int) ([]*models.RoutingRule, error) {
	return s.routingRepo.ListRoutingRules(limit, offset)
}

func (s *PostgresRoutingService) CreateRoutingRule(rule *models.RoutingRule, createdBy int64) error {
	rule.CreatedBy = &createdBy
	
	err := s.routingRepo.CreateRoutingRule(rule)
	if err != nil {
		return fmt.Errorf("failed to create routing rule: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &createdBy,
		Action:     "create_routing_rule",
		EntityType: stringPtr("routing_rule"),
		EntityID:   &rule.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) UpdateRoutingRule(rule *models.RoutingRule, updatedBy int64) error {
	err := s.routingRepo.UpdateRoutingRule(rule)
	if err != nil {
		return fmt.Errorf("failed to update routing rule: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &updatedBy,
		Action:     "update_routing_rule",
		EntityType: stringPtr("routing_rule"),
		EntityID:   &rule.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) DeleteRoutingRule(id int64, deletedBy int64) error {
	err := s.routingRepo.DeleteRoutingRule(id)
	if err != nil {
		return fmt.Errorf("failed to delete routing rule: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &deletedBy,
		Action:     "delete_routing_rule",
		EntityType: stringPtr("routing_rule"),
		EntityID:   &id,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) GetRoutingRuleByID(id int64) (*models.RoutingRule, error) {
	return s.routingRepo.GetRoutingRuleByID(id)
}

func (s *PostgresRoutingService) CheckNumberBlacklisted(number string, direction string) (*models.Blacklist, error) {
	entry, err := s.routingRepo.CheckNumberBlacklisted(number)
	if err != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", err)
	}
	
	if entry != nil && !entry.ShouldBlock(direction) {
		return nil, nil
	}
	
	return entry, nil
}

func (s *PostgresRoutingService) AddToBlacklist(entry *models.Blacklist, addedBy int64) error {
	entry.CreatedBy = &addedBy
	
	err := s.routingRepo.CreateBlacklistEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to add to blacklist: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &addedBy,
		Action:     "add_blacklist",
		EntityType: stringPtr("blacklist"),
		EntityID:   &entry.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) RemoveFromBlacklist(id int64, removedBy int64) error {
	err := s.routingRepo.DeleteBlacklistEntry(id)
	if err != nil {
		return fmt.Errorf("failed to remove from blacklist: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &removedBy,
		Action:     "remove_blacklist",
		EntityType: stringPtr("blacklist"),
		EntityID:   &id,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) GetBlacklistEntries(limit, offset int) ([]*models.Blacklist, error) {
	return s.routingRepo.ListBlacklistEntries(limit, offset)
}

func (s *PostgresRoutingService) GetBlacklistEntryByID(id int64) (*models.Blacklist, error) {
	return s.routingRepo.GetBlacklistEntryByID(id)
}

func (s *PostgresRoutingService) UpdateBlacklistEntry(entry *models.Blacklist, updatedBy int64) error {
	err := s.routingRepo.UpdateBlacklistEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to update blacklist entry: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &updatedBy,
		Action:     "update_blacklist",
		EntityType: stringPtr("blacklist"),
		EntityID:   &entry.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) DetectSpam(callerNumber, destinationNumber string, callDurationSeconds int) (bool, string, error) {
	// Get spam detection configuration
	shortCallThresholdConfig, err := s.systemRepo.GetConfigByKey("spam_short_call_threshold")
	if err != nil {
		return false, "", fmt.Errorf("failed to get spam config: %w", err)
	}
	
	shortCallThreshold := 10 // default
	if shortCallThresholdConfig != nil {
		shortCallThreshold = shortCallThresholdConfig.GetIntValue()
	}
	
	// Check for short call spam
	if callDurationSeconds < shortCallThreshold {
		return true, "short_call", nil
	}
	
	// TODO: Implement more sophisticated spam detection:
	// - High frequency calling patterns
	// - Known spam number patterns
	// - Call success rate analysis
	// - Carrier feedback integration
	
	return false, "", nil
}

func (s *PostgresRoutingService) AutoBlacklistNumber(number string, reason string, detectionMethod string) error {
	// Check if auto-blacklisting is enabled
	autoBlacklistConfig, err := s.systemRepo.GetConfigByKey("auto_blacklist_enabled")
	if err != nil {
		return fmt.Errorf("failed to get auto blacklist config: %w", err)
	}
	
	if autoBlacklistConfig == nil || !autoBlacklistConfig.GetBoolValue() {
		return nil // Auto-blacklisting disabled
	}
	
	// Check if number is already blacklisted
	existing, err := s.routingRepo.CheckNumberBlacklisted(number)
	if err != nil {
		return fmt.Errorf("failed to check existing blacklist: %w", err)
	}
	
	if existing != nil {
		// Update violation count
		existing.ViolationCount++
		existing.LastViolationAt = time.Now()
		return s.routingRepo.UpdateBlacklistEntry(existing)
	}
	
	// Create new blacklist entry
	entry := &models.Blacklist{
		NumberPattern:    number,
		BlacklistType:    models.BlacklistTypeNumber,
		Reason:           &reason,
		AutoAdded:        true,
		DetectionMethod:  &detectionMethod,
		BlockInbound:     true,
		BlockOutbound:    false,
		ViolationCount:   1,
		LastViolationAt:  time.Now(),
	}
	
	err = s.routingRepo.CreateBlacklistEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to auto-blacklist number: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		Action:     "auto_blacklist",
		EntityType: stringPtr("blacklist"),
		EntityID:   &entry.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) GetSpamStatistics() (*SpamStats, error) {
	// Get total blacklisted numbers
	allEntries, err := s.routingRepo.GetActiveBlacklistEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklist entries: %w", err)
	}
	
	// Count auto-blacklisted today
	autoBlacklistedToday := int64(0)
	today := time.Now().Truncate(24 * time.Hour)
	
	for _, entry := range allEntries {
		if entry.AutoAdded && entry.CreatedAt.After(today) {
			autoBlacklistedToday++
		}
	}
	
	stats := &SpamStats{
		TotalBlacklistedNumbers: int64(len(allEntries)),
		AutoBlacklistedToday:    autoBlacklistedToday,
		// TODO: Calculate blocked calls and top prefixes from CDR data
	}
	
	return stats, nil
}

func (s *PostgresRoutingService) GetAvailableSIMForCall(poolName string, customerID *int64) (*models.SIMPoolAssignment, error) {
	// Get SIMs in pool
	assignments, err := s.routingRepo.GetSIMsInPool(poolName)
	if err != nil {
		return nil, fmt.Errorf("failed to get SIMs in pool: %w", err)
	}
	
	if len(assignments) == 0 {
		return nil, fmt.Errorf("no SIMs available in pool %s", poolName)
	}
	
	// Simple round-robin selection for now
	// TODO: Implement proper load balancing based on pool configuration
	return assignments[0], nil
}

func (s *PostgresRoutingService) CreateSIMPool(pool *models.SIMPool, createdBy int64) error {
	pool.CreatedBy = &createdBy
	
	err := s.routingRepo.CreateSIMPool(pool)
	if err != nil {
		return fmt.Errorf("failed to create SIM pool: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &createdBy,
		Action:     "create_sim_pool",
		EntityType: stringPtr("sim_pool"),
		EntityID:   &pool.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) AssignSIMToPool(simPoolID, simCardID int64, priority int, assignedBy int64) error {
	assignment := &models.SIMPoolAssignment{
		SIMPoolID:  simPoolID,
		SIMCardID:  simCardID,
		Priority:   priority,
		IsActive:   true,
		AssignedBy: &assignedBy,
	}
	
	err := s.routingRepo.AssignSIMToPool(assignment)
	if err != nil {
		return fmt.Errorf("failed to assign SIM to pool: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &assignedBy,
		Action:     "assign_sim_to_pool",
		EntityType: stringPtr("sim_pool_assignment"),
		EntityID:   &assignment.ID,
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) RemoveSIMFromPool(simPoolID, simCardID int64, removedBy int64) error {
	err := s.routingRepo.RemoveSIMFromPool(simPoolID, simCardID)
	if err != nil {
		return fmt.Errorf("failed to remove SIM from pool: %w", err)
	}
	
	// Create audit log
	auditLog := &models.AuditLog{
		UserID:     &removedBy,
		Action:     "remove_sim_from_pool",
		EntityType: stringPtr("sim_pool_assignment"),
		Success:    true,
	}
	s.systemRepo.CreateAuditLog(auditLog)
	
	return nil
}

func (s *PostgresRoutingService) GetSIMPools(limit, offset int) ([]*models.SIMPool, error) {
	return s.routingRepo.ListSIMPools(limit, offset)
}

func (s *PostgresRoutingService) GetSIMsInPool(poolName string) ([]*models.SIMPoolAssignment, error) {
	return s.routingRepo.GetSIMsInPool(poolName)
}
