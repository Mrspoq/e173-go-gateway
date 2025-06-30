package repository

import (
	"database/sql"
	"fmt"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/jmoiron/sqlx"
)

type BlacklistRepository interface {
	Add(number string, reason string) error
	Remove(number string) error
	IsBlacklisted(number string) (bool, error)
	GetAll() ([]models.BlacklistEntry, error)
	GetByNumber(number string) (*models.BlacklistEntry, error)
}

type blacklistRepository struct {
	db *sqlx.DB
}

func NewBlacklistRepository(db *sqlx.DB) BlacklistRepository {
	return &blacklistRepository{db: db}
}

func (r *blacklistRepository) Add(number string, reason string) error {
	query := `
		INSERT INTO blacklist (phone_number, reason, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (phone_number) DO UPDATE
		SET reason = $2, updated_at = NOW()
	`
	_, err := r.db.Exec(query, number, reason)
	return err
}

func (r *blacklistRepository) Remove(number string) error {
	query := `DELETE FROM blacklist WHERE phone_number = $1`
	result, err := r.db.Exec(query, number)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("number not found in blacklist")
	}
	
	return nil
}

func (r *blacklistRepository) IsBlacklisted(number string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM blacklist WHERE phone_number = $1`
	err := r.db.Get(&count, query, number)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *blacklistRepository) GetAll() ([]models.BlacklistEntry, error) {
	var entries []models.BlacklistEntry
	query := `SELECT * FROM blacklist ORDER BY created_at DESC`
	err := r.db.Select(&entries, query)
	return entries, err
}

func (r *blacklistRepository) GetByNumber(number string) (*models.BlacklistEntry, error) {
	var entry models.BlacklistEntry
	query := `SELECT * FROM blacklist WHERE phone_number = $1`
	err := r.db.Get(&entry, query, number)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &entry, err
}