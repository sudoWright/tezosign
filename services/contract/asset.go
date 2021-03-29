package contract

import (
	"tezosign/models"
	"tezosign/types"

	"blockwatch.cc/tzindex/micheline"
)

func AssetOperation(prim *micheline.Prim, assetType models.AssetType) (transfers []models.TransferUnit) {

	//TODO add validations
	if assetType == models.TypeFA2 {
		for i := range prim.Args {
			//List
			txs := make([]models.Tx, len(prim.Args[i].Args[1].Args))
			for j, arg := range prim.Args[i].Args[1].Args {
				txs[j] = models.Tx{
					To:      types.Address(arg.Args[0].String),
					TokenID: arg.Args[1].Args[0].Int.Uint64(),
					Amount:  arg.Args[1].Args[1].Int.Uint64(),
				}
			}
			//prim.Args[i].Args[1]
			transfers = append(transfers, models.TransferUnit{
				From: types.Address(prim.Args[0].String),
				Txs:  nil,
			})
		}

		return transfers
	}

	transfers = []models.TransferUnit{
		{
			From: types.Address(prim.Args[0].String),
			Txs: []models.Tx{
				{
					To:     types.Address(prim.Args[1].Args[0].String),
					Amount: prim.Args[1].Args[1].Int.Uint64(),
				},
			},
		},
	}

	return transfers
}
