package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"tezosign/types"
)

type ContractStorageRequest struct {
	Threshold uint            `json:"threshold"`
	Addresses []types.Address `json:"addresses"`
}

func (r ContractStorageRequest) Validate() (err error) {
	if r.Threshold <= 0 {
		return fmt.Errorf("zero threshold")
	}

	if len(r.Addresses) == 0 {
		return fmt.Errorf("empty addresses")
	}

	for i := range r.Addresses {
		err = r.Addresses[i].Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type ContractInfo struct {
	Address   types.Address `json:"address"`
	Balance   int64         `json:"balance"`
	Threshold int64         `json:"threshold"`
	Counter   int64         `json:"counter"`
	Owners    []Owner       `json:"owners"`
}

type Owner struct {
	PubKey  types.PubKey  `json:"pub_key"`
	Address types.Address `json:"address"`
}

type ContractOperationRequest struct {
	ContractID types.Address `json:"contract_id"`
	Type       ActionType    `json:"type"`

	Amount uint64 `json:"amount,omitempty"`

	//Transfer Delegation
	To types.Address `json:"to"`

	From types.Address `json:"from,omitempty"`

	//FA transfer
	AssetID types.Address `json:"asset_id,omitempty"`

	//Custom json michelson payload
	CustomPayload types.Payload `json:"custom_payload,omitempty"`
	//Internal params
	//Update storage
	Threshold uint           `json:"-"`
	Keys      []types.PubKey `json:"-"`
}

func (r ContractOperationRequest) Validate() (err error) {

	err = r.ContractID.Validate()
	if err != nil {
		return err
	}

	switch r.Type {
	case Delegation:
		//Delegation destination can be empty
		if r.To.String() != "" {
			err = r.To.Validate()
		}
		if err != nil {
			return err
		}
	case FATransfer:
		err = r.AssetID.Validate()
		if err != nil {
			return err
		}

		//If From field not presented send from current contract
		if r.From.String() != "" {
			err = r.To.Validate()
		}
		if err != nil {
			return err
		}

		fallthrough
	case Transfer:
		err = r.To.Validate()
		if err != nil {
			return err
		}

		if r.Amount == 0 {
			return fmt.Errorf("wrong amount")
		}
	case CustomPayload:
		if !json.Valid([]byte(r.CustomPayload)) {
			return fmt.Errorf("wrong custom payload")
		}

	default:
		return fmt.Errorf("wrong operation type")
	}

	return nil
}

func (r *ContractOperationRequest) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}
	data, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type")
	}

	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal([]byte(data), r)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}

	return nil
}

func (r ContractOperationRequest) Value() (driver.Value, error) {

	bt, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return string(bt), nil
}

type PayloadType string

const (
	TypeApprove = "approve"
	TypeReject  = "reject"
)

type OperationSignature struct {
	ContractID types.Address `json:"contract_id"`
	SignatureReq
	Type PayloadType `json:"type"`
}

func (r OperationSignature) Validate() (err error) {

	err = r.SignatureReq.Validate()
	if err != nil {
		return err
	}

	err = r.ContractID.Validate()
	if err != nil {
		return err
	}

	err = r.Type.Validate()
	if err != nil {
		return fmt.Errorf("wrong signature type")
	}

	return nil
}

func (p PayloadType) Validate() (err error) {

	if p != TypeApprove && p != TypeReject {
		return fmt.Errorf("wrong type")
	}

	return nil
}
