package models

import (
	"msig/types"
	"time"
)

type TokenType string

const (
	TypeAuth    TokenType = "auth"
	TypeRefresh TokenType = "refresh"
)

type AuthToken struct {
	ID        uint64        `gorm:"column:atn_id;primaryKey"`
	Address   types.Address `gorm:"column:atn_address"`
	Type      TokenType     `gorm:"column:atn_type"`
	Data      string        `gorm:"column:atn_data"`
	IsUsed    bool          `gorm:"column:atn_is_used"`
	ExpiresAt time.Time     `gorm:"column:atn_expires_at"`
}

type AuthTokenReq struct {
	Address types.Address `json:"address"`
}

func (r AuthTokenReq) Validate() (err error) {
	err = r.Address.Validate()
	if err != nil {
		return err
	}
	return nil
}

type AuthTokenResp struct {
	Token string `json:"token"`
}

func (r AuthToken) Expired() bool {
	if r.ExpiresAt.Before(time.Now()) {
		return true
	}

	return false
}
