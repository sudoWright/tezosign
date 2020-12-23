package services

import (
	"context"
	"fmt"
	"log"
	"msig/models"
	"msig/services/contract"
)

const maxAddressesNum = 20

func (s *ServiceFacade) BuildContract(req models.ContractRequest) (resp []byte, err error) {

	if req.Threshold > uint(len(req.Addresses)) {
		return nil, fmt.Errorf("")
	}

	if len(req.Addresses) > maxAddressesNum {
		return nil, fmt.Errorf("")
	}

	var pubKey string
	pubKeys := make([]string, len(req.Addresses))

	for i := range req.Addresses {
		//TODO probably use indexed db
		pubKey, err = s.rpcClient.ManagerKey(context.Background(), string(req.Addresses[i]))
		if err != nil {
			return
		}

		pubKeys[i] = pubKey
	}

	resp, err = contract.BuildContractStorage(req.Threshold, pubKeys)
	if err != nil {
		return nil, err
	}

	log.Print(string(resp))
	return resp, nil
}
