package models

import (
	"tezosign/types"

	"github.com/wedancedalot/decimal"
)

type AssetBalances struct {
	Total    uint64  `json:"total"`
	Balances []Token `json:"balances"`
}

type Token struct {
	Asset types.Address `json:"contract"`
	TokenBalance
}

type TokenBalance struct {
	Balance decimal.Decimal `json:"balance"`
	TokenId uint64          `json:"token_id"`
}
