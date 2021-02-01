package services

import (
	"context"
	"msig/types"
)

func (s *ServiceFacade) GetOperationsList(userAddress types.Address, contractID types.Address, params interface{}) (resp interface{}, err error) {

	storage, err := s.getContractStorage(contractID.String())
	if err != nil {
		return
	}

	pubKey, err := s.rpcClient.ManagerKey(context.Background(), userAddress.String())
	if err != nil {
		return
	}

	_, isOwner := storage.Contains(types.PubKey(pubKey))

	contractRepo := s.repoProvider.GetContract()

	contract, err := contractRepo.GetOrCreateContract(contractID)
	if err != nil {
		return
	}

	//Get pending operations
	if isOwner {
		//Extend to Report
		resp, err = contractRepo.GetPayloadsReportByContractID(contract.ID)
		if err != nil {
			return
		}
	}

	//TODO Add income txs

	return resp, nil
}
