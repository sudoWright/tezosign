package services

import (
	"tezosign/common/apperrors"
	"tezosign/models"
	"tezosign/services/contract"
	"tezosign/types"
	"time"

	"blockwatch.cc/tzindex/micheline"
)

func (s *ServiceFacade) BuildVestingContractInitStorage(storageReq models.VestingContractStorageRequest) (resp []byte, err error) {

	resp, err = contract.BuildVestingContractStorage(storageReq.VestingAddress, storageReq.DelegateAdmin, storageReq.Timestamp, storageReq.SecondsPerTick, storageReq.TokensPerTick)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *ServiceFacade) VestingContractInfo(contractID types.Address) (info models.VestingContractInfo, err error) {

	indexerRepo := s.indexerRepoProvider.GetIndexer()

	account, isFound, err := indexerRepo.GetAccount(contractID)
	if err != nil {
		return info, err
	}

	if !isFound {
		return info, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	var delegate types.Address
	if account.DelegateID.Valid {
		delegateAccount, isFound, err := indexerRepo.GetAccountByID(account.Id)
		if err != nil {
			return info, err
		}

		if !isFound {
			return info, apperrors.New(apperrors.ErrNotFound, "delegate")
		}

		delegate = delegateAccount.Address
	}

	script, isFound, err := indexerRepo.GetContractScript(contractID)
	if err != nil {
		return info, err
	}

	if !isFound {
		return info, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	storage, isFound, err := indexerRepo.GetContractStorage(contractID)
	if err != nil {
		return info, err
	}

	if !isFound {
		return info, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	storageContainer, err := contract.NewVestingContractStorageContainer(micheline.Script{
		Code: &micheline.Code{
			Param:   script.ParameterSchema.MichelinePrim(),
			Storage: script.StorageSchema.MichelinePrim(),
			Code:    script.CodeSchema.MichelinePrim(),
		},
		Storage: storage.RawValue.MichelinePrim(),
	})
	if err != nil {
		return info, err
	}

	//Calc already opened amount
	openedAmount := (uint64(time.Now().Unix()) - storageContainer.Timestamp) / storageContainer.SecondsPerTick * storageContainer.TokensPerTick

	return models.VestingContractInfo{
		Balance:       account.Balance,
		OpenedBalance: openedAmount,
		Delegate:      delegate,
		Storage: models.VestingContractStorageRequest{
			VestingAddress: storageContainer.VestingAddress,
			DelegateAdmin:  storageContainer.DelegateAdmin,
			Timestamp:      storageContainer.Timestamp,
			SecondsPerTick: storageContainer.SecondsPerTick,
			TokensPerTick:  storageContainer.TokensPerTick,
		},
	}, nil
}

func (s *ServiceFacade) VestingContractOperation(req models.VestingContractOperation) (param models.OperationParameter, err error) {

	value, entrypoint, err := contract.VestingContractParamAndEntrypoint(req)
	if err != nil {
		return param, err
	}

	return models.OperationParameter{
		Entrypoint: entrypoint,
		Value:      string(value),
	}, nil
}
