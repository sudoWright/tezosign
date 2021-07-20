package models

import (
	"database/sql"
	"tezosign/types"
	"time"
)

type Storage struct {
	Level     uint64         `gorm:"column:Level"`
	Current   bool           `gorm:"column:Current"`
	RawValue  types.TZKTPrim `gorm:"column:RawValue"`
	JsonValue string         `gorm:"column:JsonValue"`
}

type Script struct {
	Current         bool           `gorm:"column:Current"`
	ParameterSchema types.TZKTPrim `gorm:"column:ParameterSchema"`
	StorageSchema   types.TZKTPrim `gorm:"column:StorageSchema"`
	CodeSchema      types.TZKTPrim `gorm:"column:CodeSchema"`
}

type Account struct {
	Id      uint64        `gorm:"column:Id"`
	Address types.Address `gorm:"column:Address"`
	Type    uint8         `gorm:"column:Type"`
	Balance uint64        `gorm:"column:Balance"`

	DelegateID sql.NullInt64 `gorm:"column:DelegateId"`
}

type Block struct {
	Id    uint64 `gorm:"column:Id"`
	Level uint64 `gorm:"column:Level"`

	Hash      string    `gorm:"column:Hash"`
	Timestamp time.Time `gorm:"column:Timestamp"`
}
