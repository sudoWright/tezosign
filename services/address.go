package services

import (
	"tezosign/common/apperrors"
	"tezosign/types"
)

func (s *ServiceFacade) AddressRevealed(address types.Address) (isRevealed bool, err error) {
	index := s.indexerRepoProvider.GetIndexer()

	_, isRevealed, err = index.GetContractRevealOperation(address)
	if err != nil {
		return isRevealed, err
	}

	return isRevealed, nil
}

func (s *ServiceFacade) AddressBalance(address types.Address) (balance uint64, err error) {

	acc, isFound, err := s.indexerRepoProvider.GetIndexer().GetAccount(address)
	if err != nil {
		return balance, err
	}

	if !isFound {
		return balance, apperrors.New(apperrors.ErrNotFound, "address")
	}

	return acc.Balance, nil
}
