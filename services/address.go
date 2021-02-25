package services

import (
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
