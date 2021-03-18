package models

import (
	"errors"
	"tezosign/types"
)

//storage (pair (pair %wrapped (address %target) (address %delegateAdmin))
//              (pair (nat %vested)
//              (pair %schedule (timestamp %epoch)
//              (pair (nat %secondsPerTick) (nat %tokensPerTick)))))

type VestingContractStorageRequest struct {
	VestingAddress types.Address `json:"vesting_address"`
	DelegateAdmin  types.Address `json:"delegate_admin"`
	Timestamp      uint64        `json:"timestamp"`
	SecondsPerTick uint64        `json:"seconds_per_tick"`
	TokensPerTick  uint64        `json:"tokens_per_tick"`
}

func (v VestingContractStorageRequest) Validate() (err error) {

	if err = v.VestingAddress.Validate(); err != nil {
		return err
	}

	if err = v.DelegateAdmin.Validate(); err != nil {
		return err
	}

	if v.Timestamp == 0 {
		return errors.New("timestamp")
	}

	if v.SecondsPerTick == 0 {
		return errors.New("seconds per tick")
	}

	if v.TokensPerTick == 0 {
		return errors.New("tokens per tick")
	}

	return nil
}

type VestingContractOperation struct {
	Type   ActionType    `json:"type"`
	Amount uint64        `json:"amount,omitempty"`
	To     types.Address `json:"to,omitempty"`
}

func (v VestingContractOperation) Validate() (err error) {

	switch v.Type {
	case VestingVest:
		if v.Amount == 0 {
			return errors.New("amount")
		}
	case VestingSetDelegate:
		if err = v.To.Validate(); err != nil {
			return err
		}
	default:
		return errors.New("wrong type")
	}

	return nil
}
