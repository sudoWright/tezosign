package services

import (
	"database/sql"
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

func (s *ServiceFacade) VestingsList(userPubKey types.PubKey, contractAddress types.Address, params models.CommonParams) (vestings []models.Vesting, err error) {

	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return vestings, err
	}

	if !isFound {
		return vestings, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(userPubKey, contractAddress)
	if err != nil {
		return vestings, err
	}

	//For viewer return empty arr
	if !isOwner {
		return vestings, nil
	}

	vestings, err = s.repoProvider.GetVesting().GetVestingsList(contract.ID, params.Limit, params.Offset)
	if err != nil {
		return vestings, err
	}

	//TODO add balances

	return vestings, err
}

func (s *ServiceFacade) ContractVesting(userPubKey types.PubKey, contractAddress types.Address, reqVesting models.Vesting) (vesting models.Vesting, err error) {

	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return vesting, err
	}

	if !isFound {
		return vesting, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(userPubKey, contractAddress)
	if err != nil {
		return vesting, err
	}

	if !isOwner {
		return vesting, apperrors.New(apperrors.ErrNotAllowed)
	}

	//Сheck contract for vesting type
	isVestingContract, err := s.checkVestingContract(reqVesting.Address)
	if err != nil {
		return vesting, err
	}

	if !isVestingContract {
		return vesting, apperrors.New(apperrors.ErrBadParam, "not vesting contract")
	}

	vestingRepo := s.repoProvider.GetVesting()
	vesting, isFound, err = vestingRepo.GetVesting(contract.ID, reqVesting.Address)
	if err != nil {
		return vesting, err
	}

	//Already created
	if isFound {
		return vesting, apperrors.New(apperrors.ErrAlreadyExists, "asset")
	}

	reqVesting.ContractID = sql.NullInt64{
		Int64: int64(contract.ID),
		Valid: true,
	}

	err = vestingRepo.CreateVesting(reqVesting)
	if err != nil {
		return vesting, err
	}

	return reqVesting, nil
}

func (s *ServiceFacade) ContractVestingEdit(userPubKey types.PubKey, contractAddress types.Address, reqVesting models.Vesting) (vesting models.Vesting, err error) {

	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return vesting, err
	}

	if !isFound {
		return vesting, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(userPubKey, contractAddress)
	if err != nil {
		return vesting, err
	}

	if !isOwner {
		return vesting, apperrors.New(apperrors.ErrNotAllowed)
	}

	vestingRepo := s.repoProvider.GetVesting()
	vesting, isFound, err = vestingRepo.GetVesting(contract.ID, reqVesting.Address)
	if err != nil {
		return vesting, err
	}

	//Not created
	if !isFound {
		return vesting, apperrors.New(apperrors.ErrNotFound, "asset")
	}

	//Global asset cannot be edited
	if !vesting.ContractID.Valid {
		return vesting, apperrors.New(apperrors.ErrNotAllowed, "global asset")
	}

	vesting.Name = reqVesting.Name

	err = vestingRepo.UpdateVesting(vesting)
	if err != nil {
		return vesting, err
	}

	return vesting, nil
}

func (s *ServiceFacade) RemoveContractVesting(userPubKey types.PubKey, contractAddress types.Address, vesting models.Vesting) (err error) {

	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return err
	}

	if !isFound {
		return apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(userPubKey, contractAddress)
	if err != nil {
		return err
	}

	if !isOwner {
		return apperrors.New(apperrors.ErrNotAllowed)
	}

	vestingRepo := s.repoProvider.GetVesting()
	vesting, isFound, err = vestingRepo.GetVesting(contract.ID, vesting.Address)
	if err != nil {
		return err
	}

	if !isFound {
		return apperrors.New(apperrors.ErrNotFound, "asset")
	}

	//Global asset cannot be removed
	if !vesting.ContractID.Valid {
		return apperrors.New(apperrors.ErrNotAllowed, "global asset")
	}

	err = vestingRepo.DeleteContractVesting(vesting.ID)
	if err != nil {
		return err
	}

	return nil
}
