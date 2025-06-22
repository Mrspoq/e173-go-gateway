package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type RoutingRepository interface {
	// Routing Rules
	CreateRoutingRule(rule *models.RoutingRule) error
	GetRoutingRuleByID(id int64) (*models.RoutingRule, error)
	UpdateRoutingRule(rule *models.RoutingRule) error
	DeleteRoutingRule(id int64) error
	ListRoutingRules(limit, offset int) ([]*models.RoutingRule, error)
	GetActiveRoutingRules() ([]*models.RoutingRule, error)
	GetRoutingRulesForNumber(number string) ([]*models.RoutingRule, error)
	
	// Blacklist
	CreateBlacklistEntry(entry *models.Blacklist) error
	GetBlacklistEntryByID(id int64) (*models.Blacklist, error)
	UpdateBlacklistEntry(entry *models.Blacklist) error
	DeleteBlacklistEntry(id int64) error
	ListBlacklistEntries(limit, offset int) ([]*models.Blacklist, error)
	GetActiveBlacklistEntries() ([]*models.Blacklist, error)
	CheckNumberBlacklisted(number string) (*models.Blacklist, error)
	GetAutoBlacklistedNumbers() ([]*models.Blacklist, error)
	
	// SIM Pools
	CreateSIMPool(pool *models.SIMPool) error
	GetSIMPoolByID(id int64) (*models.SIMPool, error)
	GetSIMPoolByName(name string) (*models.SIMPool, error)
	UpdateSIMPool(pool *models.SIMPool) error
	DeleteSIMPool(id int64) error
	ListSIMPools(limit, offset int) ([]*models.SIMPool, error)
	GetActiveSIMPools() ([]*models.SIMPool, error)
	
	// SIM Pool Assignments
	AssignSIMToPool(assignment *models.SIMPoolAssignment) error
	RemoveSIMFromPool(simPoolID, simCardID int64) error
	GetSIMPoolAssignments(simPoolID int64) ([]*models.SIMPoolAssignment, error)
	GetSIMsInPool(poolName string) ([]*models.SIMPoolAssignment, error)
}

type PostgresRoutingRepository struct {
	db *sqlx.DB
}

func NewPostgresRoutingRepository(db *sqlx.DB) RoutingRepository {
	return &PostgresRoutingRepository{db: db}
}

// Routing Rules methods
func (r *PostgresRoutingRepository) CreateRoutingRule(rule *models.RoutingRule) error {
	query := `
		INSERT INTO routing_rules (
			rule_name, rule_order, prefix_pattern, destination_pattern, caller_id_pattern,
			route_to_modem_id, route_to_pool, max_channels, time_restrictions,
			customer_restrictions, cost_markup_percent, is_active, notes, created_by
		) VALUES (
			:rule_name, :rule_order, :prefix_pattern, :destination_pattern, :caller_id_pattern,
			:route_to_modem_id, :route_to_pool, :max_channels, :time_restrictions,
			:customer_restrictions, :cost_markup_percent, :is_active, :notes, :created_by
		) RETURNING id, created_at, updated_at`
	
	rows, err := r.db.NamedQuery(query, rule)
	if err != nil {
		return fmt.Errorf("failed to create routing rule: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created routing rule")
}

func (r *PostgresRoutingRepository) GetRoutingRuleByID(id int64) (*models.RoutingRule, error) {
	rule := &models.RoutingRule{}
	query := `SELECT * FROM routing_rules WHERE id = $1`
	
	err := r.db.Get(rule, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get routing rule by ID: %w", err)
	}
	
	return rule, nil
}

func (r *PostgresRoutingRepository) UpdateRoutingRule(rule *models.RoutingRule) error {
	query := `
		UPDATE routing_rules SET
			rule_name = :rule_name, rule_order = :rule_order, prefix_pattern = :prefix_pattern,
			destination_pattern = :destination_pattern, caller_id_pattern = :caller_id_pattern,
			route_to_modem_id = :route_to_modem_id, route_to_pool = :route_to_pool,
			max_channels = :max_channels, time_restrictions = :time_restrictions,
			customer_restrictions = :customer_restrictions, cost_markup_percent = :cost_markup_percent,
			is_active = :is_active, notes = :notes, updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`
	
	_, err := r.db.NamedExec(query, rule)
	if err != nil {
		return fmt.Errorf("failed to update routing rule: %w", err)
	}
	
	return nil
}

func (r *PostgresRoutingRepository) DeleteRoutingRule(id int64) error {
	query := `DELETE FROM routing_rules WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete routing rule: %w", err)
	}
	
	return nil
}

func (r *PostgresRoutingRepository) ListRoutingRules(limit, offset int) ([]*models.RoutingRule, error) {
	var rules []*models.RoutingRule
	query := `SELECT * FROM routing_rules ORDER BY rule_order ASC, created_at DESC LIMIT $1 OFFSET $2`
	
	err := r.db.Select(&rules, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list routing rules: %w", err)
	}
	
	return rules, nil
}

func (r *PostgresRoutingRepository) GetActiveRoutingRules() ([]*models.RoutingRule, error) {
	var rules []*models.RoutingRule
	query := `SELECT * FROM routing_rules WHERE is_active = true ORDER BY rule_order ASC`
	
	err := r.db.Select(&rules, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active routing rules: %w", err)
	}
	
	return rules, nil
}

func (r *PostgresRoutingRepository) GetRoutingRulesForNumber(number string) ([]*models.RoutingRule, error) {
	var rules []*models.RoutingRule
	query := `
		SELECT * FROM routing_rules 
		WHERE is_active = true 
		  AND $1 LIKE prefix_pattern || '%'
		ORDER BY rule_order ASC, LENGTH(prefix_pattern) DESC`
	
	err := r.db.Select(&rules, query, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get routing rules for number: %w", err)
	}
	
	return rules, nil
}

// Blacklist methods
func (r *PostgresRoutingRepository) CreateBlacklistEntry(entry *models.Blacklist) error {
	query := `
		INSERT INTO blacklist (
			number_pattern, blacklist_type, reason, auto_added, detection_method,
			block_inbound, block_outbound, temporary_until, violation_count,
			last_violation_at, created_by
		) VALUES (
			:number_pattern, :blacklist_type, :reason, :auto_added, :detection_method,
			:block_inbound, :block_outbound, :temporary_until, :violation_count,
			:last_violation_at, :created_by
		) RETURNING id, created_at, updated_at`
	
	rows, err := r.db.NamedQuery(query, entry)
	if err != nil {
		return fmt.Errorf("failed to create blacklist entry: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created blacklist entry")
}

func (r *PostgresRoutingRepository) GetBlacklistEntryByID(id int64) (*models.Blacklist, error) {
	entry := &models.Blacklist{}
	query := `SELECT * FROM blacklist WHERE id = $1`
	
	err := r.db.Get(entry, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get blacklist entry by ID: %w", err)
	}
	
	return entry, nil
}

func (r *PostgresRoutingRepository) UpdateBlacklistEntry(entry *models.Blacklist) error {
	query := `
		UPDATE blacklist SET
			number_pattern = :number_pattern, blacklist_type = :blacklist_type, reason = :reason,
			auto_added = :auto_added, detection_method = :detection_method, block_inbound = :block_inbound,
			block_outbound = :block_outbound, temporary_until = :temporary_until,
			violation_count = :violation_count, last_violation_at = :last_violation_at,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`
	
	_, err := r.db.NamedExec(query, entry)
	if err != nil {
		return fmt.Errorf("failed to update blacklist entry: %w", err)
	}
	
	return nil
}

func (r *PostgresRoutingRepository) DeleteBlacklistEntry(id int64) error {
	query := `DELETE FROM blacklist WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete blacklist entry: %w", err)
	}
	
	return nil
}

func (r *PostgresRoutingRepository) ListBlacklistEntries(limit, offset int) ([]*models.Blacklist, error) {
	var entries []*models.Blacklist
	query := `SELECT * FROM blacklist ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	err := r.db.Select(&entries, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list blacklist entries: %w", err)
	}
	
	return entries, nil
}

func (r *PostgresRoutingRepository) GetActiveBlacklistEntries() ([]*models.Blacklist, error) {
	var entries []*models.Blacklist
	query := `
		SELECT * FROM blacklist 
		WHERE temporary_until IS NULL OR temporary_until > CURRENT_TIMESTAMP
		ORDER BY created_at DESC`
	
	err := r.db.Select(&entries, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active blacklist entries: %w", err)
	}
	
	return entries, nil
}

func (r *PostgresRoutingRepository) CheckNumberBlacklisted(number string) (*models.Blacklist, error) {
	entry := &models.Blacklist{}
	query := `
		SELECT * FROM blacklist 
		WHERE (temporary_until IS NULL OR temporary_until > CURRENT_TIMESTAMP)
		  AND (
		    (blacklist_type = 'number' AND number_pattern = $1) OR
		    (blacklist_type = 'prefix' AND $1 LIKE number_pattern || '%') OR
		    (blacklist_type = 'pattern' AND $1 LIKE REPLACE(number_pattern, '*', '%'))
		  )
		ORDER BY LENGTH(number_pattern) DESC
		LIMIT 1`
	
	err := r.db.Get(entry, query, number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check number blacklisted: %w", err)
	}
	
	return entry, nil
}

func (r *PostgresRoutingRepository) GetAutoBlacklistedNumbers() ([]*models.Blacklist, error) {
	var entries []*models.Blacklist
	query := `SELECT * FROM blacklist WHERE auto_added = true ORDER BY created_at DESC`
	
	err := r.db.Select(&entries, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get auto blacklisted numbers: %w", err)
	}
	
	return entries, nil
}

// SIM Pool methods
func (r *PostgresRoutingRepository) CreateSIMPool(pool *models.SIMPool) error {
	query := `
		INSERT INTO sim_pools (
			pool_name, description, load_balance_method, max_channels_per_sim, is_active, created_by
		) VALUES (
			:pool_name, :description, :load_balance_method, :max_channels_per_sim, :is_active, :created_by
		) RETURNING id, created_at, updated_at`
	
	rows, err := r.db.NamedQuery(query, pool)
	if err != nil {
		return fmt.Errorf("failed to create SIM pool: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&pool.ID, &pool.CreatedAt, &pool.UpdatedAt)
	}
	
	return fmt.Errorf("failed to retrieve created SIM pool")
}

func (r *PostgresRoutingRepository) GetSIMPoolByID(id int64) (*models.SIMPool, error) {
	pool := &models.SIMPool{}
	query := `SELECT * FROM sim_pools WHERE id = $1`
	
	err := r.db.Get(pool, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get SIM pool by ID: %w", err)
	}
	
	return pool, nil
}

func (r *PostgresRoutingRepository) GetSIMPoolByName(name string) (*models.SIMPool, error) {
	pool := &models.SIMPool{}
	query := `SELECT * FROM sim_pools WHERE pool_name = $1`
	
	err := r.db.Get(pool, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get SIM pool by name: %w", err)
	}
	
	return pool, nil
}

func (r *PostgresRoutingRepository) UpdateSIMPool(pool *models.SIMPool) error {
	query := `
		UPDATE sim_pools SET
			pool_name = :pool_name, description = :description, load_balance_method = :load_balance_method,
			max_channels_per_sim = :max_channels_per_sim, is_active = :is_active,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`
	
	_, err := r.db.NamedExec(query, pool)
	if err != nil {
		return fmt.Errorf("failed to update SIM pool: %w", err)
	}
	
	return nil
}

func (r *PostgresRoutingRepository) DeleteSIMPool(id int64) error {
	query := `DELETE FROM sim_pools WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete SIM pool: %w", err)
	}
	
	return nil
}

func (r *PostgresRoutingRepository) ListSIMPools(limit, offset int) ([]*models.SIMPool, error) {
	var pools []*models.SIMPool
	query := `SELECT * FROM sim_pools ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	err := r.db.Select(&pools, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list SIM pools: %w", err)
	}
	
	return pools, nil
}

func (r *PostgresRoutingRepository) GetActiveSIMPools() ([]*models.SIMPool, error) {
	var pools []*models.SIMPool
	query := `SELECT * FROM sim_pools WHERE is_active = true ORDER BY pool_name ASC`
	
	err := r.db.Select(&pools, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active SIM pools: %w", err)
	}
	
	return pools, nil
}

// SIM Pool Assignment methods
func (r *PostgresRoutingRepository) AssignSIMToPool(assignment *models.SIMPoolAssignment) error {
	query := `
		INSERT INTO sim_pool_assignments (sim_pool_id, sim_card_id, priority, is_active, assigned_by)
		VALUES (:sim_pool_id, :sim_card_id, :priority, :is_active, :assigned_by)
		RETURNING id, assigned_at`
	
	rows, err := r.db.NamedQuery(query, assignment)
	if err != nil {
		return fmt.Errorf("failed to assign SIM to pool: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		return rows.Scan(&assignment.ID, &assignment.AssignedAt)
	}
	
	return fmt.Errorf("failed to retrieve created SIM pool assignment")
}

func (r *PostgresRoutingRepository) RemoveSIMFromPool(simPoolID, simCardID int64) error {
	query := `DELETE FROM sim_pool_assignments WHERE sim_pool_id = $1 AND sim_card_id = $2`
	
	_, err := r.db.Exec(query, simPoolID, simCardID)
	if err != nil {
		return fmt.Errorf("failed to remove SIM from pool: %w", err)
	}
	
	return nil
}

func (r *PostgresRoutingRepository) GetSIMPoolAssignments(simPoolID int64) ([]*models.SIMPoolAssignment, error) {
	var assignments []*models.SIMPoolAssignment
	query := `
		SELECT * FROM sim_pool_assignments 
		WHERE sim_pool_id = $1 AND is_active = true 
		ORDER BY priority ASC, assigned_at ASC`
	
	err := r.db.Select(&assignments, query, simPoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get SIM pool assignments: %w", err)
	}
	
	return assignments, nil
}

func (r *PostgresRoutingRepository) GetSIMsInPool(poolName string) ([]*models.SIMPoolAssignment, error) {
	var assignments []*models.SIMPoolAssignment
	query := `
		SELECT spa.* FROM sim_pool_assignments spa
		JOIN sim_pools sp ON spa.sim_pool_id = sp.id
		WHERE sp.pool_name = $1 AND spa.is_active = true AND sp.is_active = true
		ORDER BY spa.priority ASC, spa.assigned_at ASC`
	
	err := r.db.Select(&assignments, query, poolName)
	if err != nil {
		return nil, fmt.Errorf("failed to get SIMs in pool: %w", err)
	}
	
	return assignments, nil
}
