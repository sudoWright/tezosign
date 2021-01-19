package models

import (
	"fmt"
	"msig/types"
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

type ContractOperationRequest struct {
	ContractID types.Address `json:"contract_id"`
	Type       ActionType    `json:"type"`

	Amount uint64 `json:"amount"`
	//Transfer Delegation
	To      types.Address `json:"to"`
	From    types.Address `json:"from"`
	AssetID types.Address `json:"asset_id"`
	//Custom json michelson payload
	CustomPayload types.Payload `json:"custom_payload"`
	//Internal params
	//Update storage
	Threshold uint           `json:"-"`
	Keys      []types.PubKey `json:"-"`
}

func (r ContractOperationRequest) Validate() (err error) {
	err = r.To.Validate()
	if err != nil {
		return err
	}

	err = r.ContractID.Validate()
	if err != nil {
		return err
	}

	return nil
}

type OperationSignature struct {
	ContractID types.Address   `json:"contract_id"`
	PubKey     types.PubKey    `json:"pub_key"`
	Payload    types.Payload   `json:"payload"`
	Signature  types.Signature `json:"signature"`
}

func (r OperationSignature) Validate() (err error) {
	//TODO add validation
	return nil
}
