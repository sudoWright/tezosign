package models

import "fmt"

type ContractStorageRequest struct {
	Threshold uint      `json:"threshold"`
	Addresses []Address `json:"addresses"`
}

func (r ContractStorageRequest) Validate() (err error) {
	if r.Threshold == 0 {
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
	ContractID Address    `json:"contract_id"`
	Type       ActionType `json:"type"`

	Amount uint64 `json:"amount"`
	//Transfer Delegation
	To      Address `json:"to"`
	From    Address `json:"from"`
	AssetID Address `json:"asset_id"`
	//Custom json michelson payload
	CustomPayload []byte `json:"custom_payload"`
	//Internal params
	//Update storage
	Threshold uint     `json:"-"`
	Keys      []string `json:"-"`
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

	//if r.Amount == 0 {
	//	return fmt.Errorf("amount")
	//}

	return nil
}

type OperationSignature struct {
	ContractID Address `json:"contract_id"`
	PubKey     string  `json:"pub_key"`
	Payload    string  `json:"payload"`
	Signature  string  `json:"signature"`
}

func (r OperationSignature) Validate() (err error) {

	return nil
}
