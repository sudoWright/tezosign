package models

import "database/sql"

type Asset struct {
	ID            uint64  `gorm:"column:ast_id;primaryKey" json:"-"`
	Name          string  `gorm:"column:ast_name" json:"name"`
	ContractType  string  `gorm:"column:ast_contract_type" json:"contract_type"`
	Address       string  `gorm:"column:ast_address" json:"address"`
	DexterAddress *string `gorm:"column:ast_dexter_address" json:"dexter_address,omitempty"`
	Scale         uint8   `gorm:"column:ast_scale" json:"scale"`
	Ticker        string  `gorm:"column:ast_ticker" json:"ticker"`

	ContractID sql.NullInt64 `gorm:"column:ctr_id" json:"-"`
}
