package models

import "github.com/wedancedalot/decimal"

type Quote struct {
	Id    uint64          `gorm:"column:Id" json:"-"`
	Level uint64          `gorm:"column:Level" json:"-"`
	BTC   decimal.Decimal `gorm:"column:Btc" json:"btc"`
	Eur   decimal.Decimal `gorm:"column:Eur" json:"eur"`
	Usd   decimal.Decimal `gorm:"column:Usd" json:"usd"`
	Cny   decimal.Decimal `gorm:"column:Cny" json:"cny"`
	Jpy   decimal.Decimal `gorm:"column:Jpy" json:"jpy"`
	Krw   decimal.Decimal `gorm:"column:Krw" json:"krw"`
}
