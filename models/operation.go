package models

import (
	"tezosign/types"
)

type Contract struct {
	ID                      uint64        `gorm:"column:ctr_id;primaryKey"`
	Address                 types.Address `gorm:"column:ctr_address"`
	LastOperationBlockLevel uint64        `gorm:"column:ctr_last_block_level"`
}

type RequestStatus string

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
	//Status for incoming transfers
	StatusSuccess = "success"
)

type Request struct {
	ID         uint64                   `gorm:"column:req_id;primaryKey" json:"-"`
	Hash       string                   `gorm:"column:req_hash" json:"operation_id,omitempty"`
	ContractID uint64                   `gorm:"column:ctr_id" json:"-"`
	Counter    *int64                   `gorm:"column:req_counter" json:"nonce,omitempty"`
	Status     RequestStatus            `gorm:"column:req_status;default:pending" json:"status"`
	CreatedAt  types.JSONTimestamp      `gorm:"column:req_created_at" json:"created_at"`
	Info       ContractOperationRequest `gorm:"column:req_info" json:"operation_info"`
	NetworkID  string                   `gorm:"column:req_network_id" json:"network_id"`

	OperationID string `gorm:"column:req_operation_id" json:"tx_id,omitempty"`
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

//Indexer Tezos operation
type TezosOperation struct {
	Level              uint64              `gorm:"column:Level"`
	OpHash             string              `gorm:"column:OpHash"`
	Status             int                 `gorm:"column:Status"`
	Errors             string              `gorm:"column:Errors"`
	BakerFee           uint64              `gorm:"column:BakerFee"`
	StorageFee         uint64              `gorm:"column:StorageFee"`
	AllocationFee      uint64              `gorm:"column:AllocationFee"`
	Amount             uint64              `gorm:"column:Amount"`
	Parameters         string              `gorm:"column:Parameters"`
	InternalOperations uint64              `gorm:"column:InternalOperations"`
	Timestamp          types.JSONTimestamp `gorm:"column:Timestamp"`
}
