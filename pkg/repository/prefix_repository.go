package repository

import (
	"database/sql"
	"fmt"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/jmoiron/sqlx"
)

type PrefixRepository interface {
	Create(prefix *models.Prefix) error
	GetByID(id string) (*models.Prefix, error)
	GetByPrefix(prefix string) (*models.Prefix, error)
	GetAllActive() ([]models.Prefix, error)
	Update(prefix *models.Prefix) error
	Delete(id string) error
}

type prefixRepository struct {
	db *sqlx.DB
}

func NewPrefixRepository(db *sqlx.DB) PrefixRepository {
	return &prefixRepository{db: db}
}

func (r *prefixRepository) Create(prefix *models.Prefix) error {
	query := `
		INSERT INTO prefixes (id, prefix, country, operator, gateway_id, rate_per_minute, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	_, err := r.db.Exec(query, prefix.ID, prefix.Prefix, prefix.Country, prefix.Operator, 
		prefix.GatewayID, prefix.RatePerMinute, prefix.IsActive)
	return err
}

func (r *prefixRepository) GetByID(id string) (*models.Prefix, error) {
	var prefix models.Prefix
	query := `SELECT * FROM prefixes WHERE id = $1`
	err := r.db.Get(&prefix, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &prefix, err
}

func (r *prefixRepository) GetByPrefix(prefixStr string) (*models.Prefix, error) {
	var prefix models.Prefix
	query := `SELECT * FROM prefixes WHERE prefix = $1 AND is_active = true`
	err := r.db.Get(&prefix, query, prefixStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &prefix, err
}

func (r *prefixRepository) GetAllActive() ([]models.Prefix, error) {
	var prefixes []models.Prefix
	query := `SELECT * FROM prefixes WHERE is_active = true ORDER BY LENGTH(prefix) DESC`
	err := r.db.Select(&prefixes, query)
	return prefixes, err
}

func (r *prefixRepository) Update(prefix *models.Prefix) error {
	query := `
		UPDATE prefixes 
		SET prefix = $2, country = $3, operator = $4, gateway_id = $5, 
		    rate_per_minute = $6, is_active = $7, updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.Exec(query, prefix.ID, prefix.Prefix, prefix.Country, 
		prefix.Operator, prefix.GatewayID, prefix.RatePerMinute, prefix.IsActive)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("prefix not found")
	}
	
	return nil
}

func (r *prefixRepository) Delete(id string) error {
	query := `DELETE FROM prefixes WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("prefix not found")
	}
	
	return nil
}