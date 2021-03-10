package models

import (
	"tezosign/types"
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
