package services

import (
	"tezosign/models"
	"tezosign/services/contract"
)

func (s *ServiceFacade) BuildVestingContractInitStorage(storageReq models.VestingContractStorageRequest) (resp []byte, err error) {

	resp, err = contract.BuildVestingContractStorage(storageReq.VestingAddress, storageReq.DelegateAdmin, storageReq.Timestamp, storageReq.SecondsPerTick, storageReq.TokensPerTick)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
