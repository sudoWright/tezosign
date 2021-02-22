package services

import (
	"context"
	"encoding/json"
	"tezosign/models"
	"tezosign/services/contract"

	contractRepo "tezosign/repos/contract"
	"tezosign/types"
)

func (s *ServiceFacade) GetOperationsList(userAddress types.Address, contractID types.Address, params interface{}) (resp []models.RequestReport, err error) {

	storage, err := s.getContractStorage(contractID.String())
	if err != nil {
		return resp, err
	}

	pubKey, err := s.rpcClient.ManagerKey(context.Background(), userAddress.String())
	if err != nil {
		return resp, err
	}

	_, isOwner := storage.Contains(types.PubKey(pubKey))

	repo := s.repoProvider.GetContract()

	contract, err := repo.GetOrCreateContract(contractID)
	if err != nil {
		return resp, err
	}

	//Get pending operations
	if isOwner {
		//Extend to Report
		resp, err = repo.GetPayloadsReportByContractID(contract.ID)
		if err != nil {
			return
		}
	}

	//TODO Add income txs

	return resp, nil
}

func (s *ServiceFacade) CheckOperations() (counter int64, err error) {
	//Init transaction
	s.repoProvider.Start(context.Background())
	defer s.repoProvider.RollbackUnlessCommitted()

	// Get contracts
	repo := s.repoProvider.GetContract()

	indexerRepo := s.indexerRepoProvider.GetIndexer()

	//Todo init limit,offset params
	limit := 100

	contracts, err := repo.GetContractsList(limit, 0)
	if err != nil {
		return counter, err
	}

	networkID, err := s.rpcClient.ChainID(context.Background())
	if err != nil {
		return counter, err
	}

	//var parameter contract.Operation
	for i := range contracts {

		operations, err := indexerRepo.GetContractOperations(contracts[i].Address, contracts[i].LastOperationBlockLevel)
		if err != nil {
			return counter, err
		}

		if len(operations) == 0 {
			continue
		}

		lastOperationBlockLevel := operations[len(operations)-1].Level

		counter, err = s.processOperations(repo, contracts[i], networkID, operations)
		if err != nil {
			return counter, err
		}

		err = repo.UpdateContractLastOperationBlock(contracts[i].ID, lastOperationBlockLevel)
		if err != nil {
			return counter, err
		}

	}

	err = s.repoProvider.Commit()
	if err != nil {
		return counter, err
	}

	return counter, nil
}

func (s *ServiceFacade) processOperations(repo contractRepo.Repo, c models.Contract, networkID string, operations []models.TransactionOperation) (counter int64, err error) {
	var parameter contract.Operation

	for j := range operations {
		//Not success tx
		if operations[j].Status != 1 {
			continue
		}

		//TODO add assets income transfers
		//Default entrypoint
		if len(operations[j].Parameters) == 0 {
			err = repo.SavePayload(models.Request{
				Hash:       operationID(operations[j].OpHash),
				ContractID: c.ID,
				Counter:    nil,
				Status:     models.StatusSuccess,
				CreatedAt:  operations[j].Timestamp,
				Info: models.ContractOperationRequest{
					ContractID: c.Address,
					Type:       models.IncomeTransfer,
					Amount:     operations[j].Amount,
					//TODO probably add From To
				},
				NetworkID:   networkID,
				OperationID: &operations[j].OpHash,
			})
			if err != nil {
				return counter, err
			}
			//Increment updated operations
			counter++
			continue
		}

		err = json.Unmarshal([]byte(operations[j].Parameters), &parameter)
		if err != nil {
			return counter, err
		}

		//Parse value
		counter, isReject, err := contract.GetOperationCounter(parameter)
		if err != nil {
			return counter, err
		}

		//TODO check that operation payload equal to db operation payload
		payload, isFound, err := repo.GetPayloadByContractAndCounter(c.ID, counter)
		if err != nil {
			return counter, err
		}

		if !isFound {
			//Probably some manual operation
			continue
		}

		payload.Status = models.StatusApproved
		if isReject {
			payload.Status = models.StatusRejected
		}

		//TODO process update signers request

		payload.OperationID = &operations[j].OpHash

		err = repo.UpdatePayload(payload)
		if err != nil {
			return counter, err
		}

		//Increment updated operations
		counter++
	}

	return counter, nil
}
