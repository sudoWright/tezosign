package services

import (
	"context"
	"msig/common/apperrors"
	"msig/models"
	"msig/services/contract"
)

const maxAddressesNum = 20

func (s *ServiceFacade) BuildContractInitStorage(req models.ContractStorageRequest) (resp []byte, err error) {

	if req.Threshold > uint(len(req.Addresses)) {
		return nil, apperrors.New(apperrors.ErrBadParam, "threshold")
	}

	if len(req.Addresses) > maxAddressesNum {
		return nil, apperrors.New(apperrors.ErrBadParam, "addresses num")
	}

	pubKeys, err := s.getPubKeysByAddresses(req.Threshold, req.Addresses)
	if err != nil {
		return nil, err
	}

	resp, err = contract.BuildContractStorage(req.Threshold, pubKeys)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ServiceFacade) BuildContractStorageUpdateOperation(req models.ContractStorageRequest) (resp string, err error) {
	if req.Threshold > uint(len(req.Addresses)) {
		return resp, apperrors.New(apperrors.ErrBadParam, "threshold")
	}

	if len(req.Addresses) > maxAddressesNum {
		return resp, apperrors.New(apperrors.ErrBadParam, "addresses num")
	}

	pubKeys, err := s.getPubKeysByAddresses(req.Threshold, req.Addresses)
	if err != nil {
		return resp, err
	}

	resp, err = s.BuildContractOperation(models.ContractOperationRequest{
		Threshold: req.Threshold,
		Keys:      pubKeys,
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *ServiceFacade) getPubKeysByAddresses(threshold uint, addresses []models.Address) (pubKeys []string, err error) {
	var pubKey string
	pubKeys = make([]string, len(addresses))

	for i := range addresses {
		//TODO probably use indexed db
		pubKey, err = s.rpcClient.ManagerKey(context.Background(), string(addresses[i]))
		if err != nil {
			return
		}

		if len(pubKey) == 0 {
			return nil, apperrors.New(apperrors.ErrBadParam, "address")
		}
		pubKeys[i] = pubKey
	}

	return pubKeys, err
}

func (s *ServiceFacade) BuildContractOperation(req models.ContractOperationRequest) (resp string, err error) {

	chainID, err := s.rpcClient.ChainID(context.Background())
	if err != nil {
		return resp, err
	}

	//TODO get counter from contract storage
	counter := int64(0)

	//TODO check another txs with same counter

	resp, err = contract.BuildContractSignPayload(chainID, counter, req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
