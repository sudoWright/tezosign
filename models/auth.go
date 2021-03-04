package models

import (
	"tezosign/types"
	"time"
)

type TokenType string

const (
	TypeAuth    TokenType = "auth"
	TypeRefresh TokenType = "refresh"
)

type AuthToken struct {
	ID        uint64       `gorm:"column:atn_id;primaryKey"`
	PubKey    types.PubKey `gorm:"column:atn_pubkey"`
	Type      TokenType    `gorm:"column:atn_type"`
	Data      string       `gorm:"column:atn_data"`
	IsUsed    bool         `gorm:"column:atn_is_used"`
	ExpiresAt time.Time    `gorm:"column:atn_expires_at"`
}

type AuthTokenReq struct {
	PubKey types.PubKey `json:"pub_key"`
}

func (r AuthTokenReq) Validate() (err error) {

	err = r.PubKey.Validate()
	if err != nil {
		return err
	}
	return nil
}

type AuthTokenResp struct {
	Token string `json:"token"`
}

func (r AuthToken) Expired() bool {
	return r.ExpiresAt.Before(time.Now())
}

type AuthSignature struct {
	Payload types.Payload `json:"payload"`
	SignatureReq
}

func (s AuthSignature) Validate() (err error) {
	err = s.Payload.Validate()
	if err != nil {
		return err
	}

	err = s.SignatureReq.Validate()
	if err != nil {
		return err
	}
	return nil
}
