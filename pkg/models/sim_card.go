package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type SIMCard struct {
	ID                      int64            `json:"id"`
	ModemID                 sql.NullInt64    `json:"modem_id,omitempty"` // Foreign key, can be NULL
	ICCID                   string           `json:"iccid"`
	IMSI                    sql.NullString   `json:"imsi,omitempty"`
	MSISDN                  sql.NullString   `json:"msisdn,omitempty"`
	OperatorName            sql.NullString   `json:"operator_name,omitempty"`
	NetworkCountryCode      sql.NullString   `json:"network_country_code,omitempty"`
	Balance                 sql.NullFloat64  `json:"balance,omitempty"` // DECIMAL(10, 4)
	BalanceCurrency         sql.NullString   `json:"balance_currency,omitempty"`
	BalanceLastCheckedAt    sql.NullTime     `json:"balance_last_checked_at,omitempty"`
	DataAllowanceMB         sql.NullInt32    `json:"data_allowance_mb,omitempty"`
	DataUsedMB              sql.NullInt32    `json:"data_used_mb,omitempty"`
	Status                  string           `json:"status"` // NOT NULL DEFAULT 'unknown'
	PIN1                    sql.NullString   `json:"pin1,omitempty"`
	PUK1                    sql.NullString   `json:"puk1,omitempty"`
	PIN2                    sql.NullString   `json:"pin2,omitempty"`
	PUK2                    sql.NullString   `json:"puk2,omitempty"`
	ActivationDate          sql.NullTime     `json:"activation_date,omitempty"` // DATE
	ExpiryDate              sql.NullTime     `json:"expiry_date,omitempty"`   // DATE
	RechargeHistory         json.RawMessage  `json:"recharge_history,omitempty"` // JSONB
	Notes                   sql.NullString   `json:"notes,omitempty"`
	CellID                  sql.NullString   `json:"cell_id,omitempty"`
	LAC                     sql.NullString   `json:"lac,omitempty"`
	PSC                     sql.NullString   `json:"psc,omitempty"`
	RSCP                    sql.NullInt32    `json:"rscp,omitempty"`
	ECIO                    sql.NullInt32    `json:"ecio,omitempty"`
	BTSInfoHistory          json.RawMessage  `json:"bts_info_history,omitempty"` // JSONB
	CreatedAt               time.Time        `json:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at"`
}
