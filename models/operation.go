package models

import "msig/types"

type Contract struct {
	ID      uint64        `gorm:"column:ctr_id;primaryKey"`
	Address types.Address `gorm:"column:ctr_address"`
}

type Request struct {
	ID         uint64 `gorm:"column:req_id;primaryKey"`
	Hash       string `gorm:"column:req_hash"`
	ContractID uint64 `gorm:"column:ctr_id"`
	Counter    int64  `gorm:"column:req_counter"`
	Status     string `gorm:"column:req_status;default:wait"`
	//Hexed payload with watermark byte
	Data types.Payload `gorm:"column:req_data"`
}

type Signature struct {
	ID        uint64          `gorm:"column:sig_id;primaryKey"`
	RequestID uint64          `gorm:"column:req_id"`
	Index     int64           `gorm:"column:sig_index"`
	Signature types.Signature `gorm:"column:sig_data"`
}

type OperationToSignResp struct {
	OperationID string `json:"operation_id"`
	Payload     string `json:"payload"`
}

type OperationSignatureResp struct {
	SigCount  int64 `json:"sig_count"`
	Threshold int64 `json:"threshold"`
}

type OperationParameter struct {
	Entrypoint string `json:"entrypoint"`
	Value      string `json:"value,omitempty"`
}
