package services

import (
	"log"
	"reflect"
	"testing"
	"tezosign/models"
	"tezosign/types"
)

func Test_GroupOperations(t *testing.T) {

	contracts := map[types.Address]models.Contract{
		"address1": {},
		"address2": {},
		"address3": {},
	}

	txs := []models.TransferUnit{
		{
			From: "address_from",
			Txs: []models.Tx{
				{
					To:      "address1",
					TokenID: 0,
					Amount:  111,
				},
				//Not our
				{
					To:      "address4",
					TokenID: 0,
					Amount:  111,
				},
				//NFT token
				{
					To:      "address1",
					TokenID: 2,
					Amount:  1,
				},
			},
		},
		{
			From: "address_from_2",
			Txs: []models.Tx{
				{
					To:      "address1",
					TokenID: 0,
					Amount:  222,
				},
			},
		},
	}

	result := map[types.Address][]models.TransferUnit{
		"address1": {
			{
				From: "address_from",
				Txs: []models.Tx{
					{
						To:      "address1",
						TokenID: 0,
						Amount:  111,
					},
					//NFT token
					{
						To:      "address1",
						TokenID: 2,
						Amount:  1,
					},
				},
			},
			{
				From: "address_from_2",
				Txs: []models.Tx{
					{
						To:      "address1",
						TokenID: 0,
						Amount:  222,
					},
				},
			},
		},
	}

	res := groupOperations(contracts, txs)

	if !reflect.DeepEqual(res, result) {
		log.Print(res)
		log.Print(result)
	}
}
