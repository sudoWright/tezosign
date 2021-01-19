package models

import "msig/types"

type Contract struct {
	ID      uint64        `gorm:"column:ctr_id;PRIMARY_KEY;DEFAULT"`
	Address types.Address `gorm:"column:ctr_address"`
}

type Request struct {
	ID         uint64 `gorm:"column:req_id;PRIMARY_KEY;DEFAULT"`
	Hash       string `gorm:"column:req_hash"`
	ContractID uint64 `gorm:"column:ctr_id"`
	Counter    int64  `gorm:"column:req_counter"`
	Status     string `gorm:"column:req_status"`
	//Hexed payload with watermark byte
	Data types.Payload `gorm:"column:req_data"`
}

type Sig struct {
	ID        uint64          `gorm:"column:sig_id;PRIMARY_KEY;DEFAULT"`
	RequestID uint64          `gorm:"column:req_id"`
	Index     int64           `gorm:"column:sig_index"`
	Signature types.Signature `gorm:"column:sig_data"`
}
