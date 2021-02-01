package models

import (
	"msig/types"
	"time"
)

type Contract struct {
	ID      uint64        `gorm:"column:ctr_id;primaryKey"`
	Address types.Address `gorm:"column:ctr_address"`
}

type RequestStatus string

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type Request struct {
	ID         uint64                   `gorm:"column:req_id;primaryKey" json:"-"`
	Hash       string                   `gorm:"column:req_hash" json:"operation_id"`
	ContractID uint64                   `gorm:"column:ctr_id" json:"-"`
	Counter    int64                    `gorm:"column:req_counter" json:"nonce"`
	Status     RequestStatus            `gorm:"column:req_status;default:pending" json:"status"`
	CreatedAt  time.Time                `gorm:"column:req_created_at" json:"created_at"`
	Info       ContractOperationRequest `gorm:"column:req_info" json:"operation_info"`
	NetworkID  string                   `gorm:"column:req_network_id" json:"network_id"`

	//TODO add request operation
	//OperationID string  `gorm:"column:req_operation_id" json:"operation_id"`
}

type RequestReport struct {
	Request
	Signatures Signatures `gorm:"column:signatures" json:"signatures,omitempty"`
}

type OperationToSignResp struct {
	OperationID string        `json:"operation_id"`
	Payload     types.Payload `json:"payload"`
}

type OperationSignatureResp struct {
	SigCount  int64 `json:"sig_count"`
	Threshold int64 `json:"threshold"`
}

type OperationParameter struct {
	Entrypoint string `json:"entrypoint"`
	Value      string `json:"value,omitempty"`
}
