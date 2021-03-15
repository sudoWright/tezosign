package models

import (
	"database/sql"
	"errors"
	"tezosign/types"
)

type AssetType string

const (
	TypeFA12 AssetType = "FA1.2"
	TypeFA2  AssetType = "FA2"
)

type Asset struct {
	ID            uint64        `gorm:"column:ast_id;primaryKey" json:"-"`
	Name          string        `gorm:"column:ast_name" json:"name"`
	ContractType  AssetType     `gorm:"column:ast_contract_type" json:"contract_type"`
	Address       types.Address `gorm:"column:ast_address" json:"address"`
	DexterAddress *string       `gorm:"column:ast_dexter_address" json:"-"`
	Scale         uint8         `gorm:"column:ast_scale" json:"scale"`
	Ticker        string        `gorm:"column:ast_ticker" json:"ticker"`

	ContractID sql.NullInt64 `gorm:"column:ctr_id" json:"-"`

	Balances []TokenBalance `gorm:"-" json:"balances"`
	IsGlobal bool           `gorm:"-" json:"is_global"`
}

func (a Asset) Validate() (err error) {

	if err = a.Address.Validate(); err != nil {
		return err
	}

	if len(a.Name) == 0 || len(a.Name) > 32 {
		return errors.New("name")
	}

	if len(a.Ticker) == 0 || len(a.Ticker) > 5 {
		return errors.New("ticker")
	}

	if a.ContractType != TypeFA12 && a.ContractType != TypeFA2 {
		return errors.New("contract_type")
	}

	return nil
}
