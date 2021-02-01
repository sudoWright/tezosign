package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"msig/types"
)

type SignatureReq struct {
	PubKey    types.PubKey    `json:"pub_key"`
	Payload   types.Payload   `json:"payload"`
	Signature types.Signature `json:"signature"`
}

func (r SignatureReq) Validate() (err error) {

	//TODO refactor
	if r.Payload != "" {
		err = r.Payload.Validate()
		if err != nil {
			return err
		}
	}

	err = r.PubKey.Validate()
	if err != nil {
		return err
	}

	err = r.Signature.Validate()
	if err != nil {
		return err
	}

	return nil
}

type Signature struct {
	ID        uint64          `gorm:"column:sig_id;primaryKey"  json:"-"`
	RequestID uint64          `gorm:"column:req_id"  json:"-"`
	Index     int64           `gorm:"column:sig_index" json:"index"`
	Signature types.Signature `gorm:"column:sig_data" json:"signature"`
	Type      PayloadType     `gorm:"column:sig_type" json:"type"`
}

type Signatures []Signature

func (s *Signatures) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type")
	}

	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}

	return nil
}

func (s Signatures) Value() (driver.Value, error) {

	bt, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return string(bt), nil
}
