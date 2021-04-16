package contract

import (
	"errors"
	"fmt"
	"math/big"
	"tezosign/models"
	"tezosign/types"

	"blockwatch.cc/tzindex/micheline"
)

func buildActionParams(operationParams models.ContractOperationRequest) (actionParams *micheline.Prim, err error) {

	switch operationParams.Type {
	case models.Transfer:
		encodedDestinationAddress, err := operationParams.To.MarshalBinary()
		if err != nil {
			return actionParams, err
		}

		//Pair (address int)
		actionParams = &micheline.Prim{
			Type:   micheline.PrimBinary,
			OpCode: micheline.D_PAIR,
			Args: []*micheline.Prim{
				//Destination
				{
					Type:   micheline.PrimBytes,
					OpCode: micheline.T_BYTES,
					Bytes:  encodedDestinationAddress,
				},
				//Amount
				{
					Type:   micheline.PrimInt,
					OpCode: micheline.T_INT,
					Int:    big.NewInt(int64(operationParams.Amount)),
				},
			},
		}

	case models.FATransfer, models.FA2Transfer:

		actionParams, err = buildFATransferParams(operationParams)
		if err != nil {
			return actionParams, err
		}

	case models.Delegation:
		actionParams, err = buildDelegationPrim(operationParams.To)
		if err != nil {
			return actionParams, err
		}

	case models.StorageUpdate:
		actionParams, err = buildStorageMichelsonArgs(int64(operationParams.Threshold), operationParams.Keys)
		if err != nil {
			return actionParams, err
		}
	case models.VestingVest, models.VestingSetDelegate:
		actionParams, err = buildVestingTxPrim(operationParams)
		if err != nil {
			return actionParams, err
		}
	case models.CustomPayload:
		actionParams = &micheline.Prim{}
		if len(operationParams.CustomPayload) == 0 {
			return actionParams, nil
		}

		//Hex payload
		if operationParams.CustomPayload.HasPrefix() {
			bt, err := operationParams.CustomPayload.MarshalBinary()
			if err != nil {
				return actionParams, err
			}
			if len(bt) == 0 {
				return actionParams, nil
			}
			//Remove watermark
			if bt[0] == TextWatermark {
				bt = bt[1:]
			}

			err = actionParams.UnmarshalBinary(bt)
		} else {
			err = actionParams.UnmarshalJSON([]byte(operationParams.CustomPayload.String()))
		}
		if err != nil {
			return actionParams, err
		}
	default:
		return actionParams, fmt.Errorf("unknown action")
	}

	return actionParams, nil
}

func buildFATransferParams(operationParams models.ContractOperationRequest) (actionParams *micheline.Prim, err error) {
	encodedAssetAddress, err := operationParams.AssetID.MarshalBinary()
	if err != nil {
		return actionParams, err
	}

	var transferParam *micheline.Prim
	var opCode micheline.OpCode
	switch operationParams.Type {
	//(pair  address (pair address nat))
	case models.FATransfer:
		opCode = micheline.D_LEFT

		//Check TransferList len
		if len(operationParams.TransferList) != 1 || len(operationParams.TransferList[0].Txs) != 1 {
			return actionParams, errors.New("wrong transfers num")
		}

		//Use contract self address as default from
		from := operationParams.ContractID
		if !operationParams.TransferList[0].From.IsEmpty() {
			from = operationParams.TransferList[0].From
		}

		encodedFromAddress, err := from.MarshalBinary()
		if err != nil {
			return actionParams, err
		}

		encodedDestinationAddress, err := operationParams.TransferList[0].Txs[0].To.MarshalBinary()
		if err != nil {
			return actionParams, err
		}

		transferParam = &micheline.Prim{
			Type:   micheline.PrimBinary,
			OpCode: micheline.D_PAIR,
			Args: []*micheline.Prim{
				//From address
				{
					Type:   micheline.PrimBytes,
					OpCode: micheline.T_BYTES,
					Bytes:  encodedFromAddress,
				},
				{
					Type:   micheline.PrimBinary,
					OpCode: micheline.D_PAIR,
					Args: []*micheline.Prim{
						//Destination address
						{
							Type:   micheline.PrimBytes,
							OpCode: micheline.T_BYTES,
							Bytes:  encodedDestinationAddress,
						},
						//Amount
						{
							Type:   micheline.PrimInt,
							OpCode: micheline.T_INT,
							Int:    big.NewInt(int64(operationParams.TransferList[0].Txs[0].Amount)),
						},
					},
				},
			},
		}

		//(list (pair address (list (pair address (pair nat nat)))))
	case models.FA2Transfer:
		opCode = micheline.D_RIGHT

		transfers := make([]*micheline.Prim, len(operationParams.TransferList))

		for i := range operationParams.TransferList {
			//(list (pair address (pair nat nat)))
			txsPrim, err := buildTransferTxsPrim(operationParams.TransferList[i].Txs)
			if err != nil {
				return actionParams, err
			}

			//Use contract self address as default from
			from := getFATransferAddressFrom(operationParams, int64(i))

			encodedFromAddress, err := from.MarshalBinary()
			if err != nil {
				return actionParams, err
			}

			//(pair address (list (pair address (pair nat nat))))
			transfers[i] = &micheline.Prim{
				Type:   micheline.PrimBinary,
				OpCode: micheline.D_PAIR,
				Args: []*micheline.Prim{
					{
						Type:   micheline.PrimBytes,
						OpCode: micheline.T_BYTES,
						Bytes:  encodedFromAddress,
					},
					txsPrim,
				},
			}
		}

		// list
		transferParam = &micheline.Prim{
			Type:   micheline.PrimSequence,
			OpCode: micheline.T_LIST,
			Args:   transfers,
		}

	default:
		return actionParams, errors.New("unknown FA format")
	}

	//(pair address (or (pair  address (pair address nat)) (list (pair address (list (pair address (pair nat nat))))) ) )
	actionParams = &micheline.Prim{
		Type:   micheline.PrimBinary,
		OpCode: micheline.D_PAIR,
		Args: []*micheline.Prim{
			//Asset
			{
				Type:   micheline.PrimBytes,
				OpCode: micheline.T_BYTES,
				Bytes:  encodedAssetAddress,
			},
			//Contract call pair
			{
				Type:   micheline.PrimUnary,
				OpCode: opCode,
				Args:   []*micheline.Prim{transferParam},
			},
		},
	}

	return actionParams, nil
}

func buildTransferTxsPrim(txs []models.Tx) (txsPrim *micheline.Prim, err error) {

	txsPrimArgs := make([]*micheline.Prim, len(txs))

	for j := range txs {
		encodedDestinationAddress, err := txs[j].To.MarshalBinary()
		if err != nil {
			return txsPrim, err
		}

		//(pair address (pair nat nat))
		txsPrimArgs[j] = &micheline.Prim{
			Type:   micheline.PrimBinary,
			OpCode: micheline.D_PAIR,
			Args: []*micheline.Prim{
				//Destination address
				{
					Type:   micheline.PrimBytes,
					OpCode: micheline.T_BYTES,
					Bytes:  encodedDestinationAddress,
				},
				{
					Type:   micheline.PrimBinary,
					OpCode: micheline.D_PAIR,
					Args: []*micheline.Prim{
						//Token_ID
						{
							Type:   micheline.PrimInt,
							OpCode: micheline.T_INT,
							Int:    big.NewInt(int64(txs[j].TokenID)),
						},
						//Amount
						{
							Type:   micheline.PrimInt,
							OpCode: micheline.T_INT,
							Int:    big.NewInt(int64(txs[j].Amount)),
						},
					},
				},
			},
		}
	}

	txsPrim = &micheline.Prim{
		Type:   micheline.PrimSequence,
		OpCode: micheline.T_LIST,
		Args:   txsPrimArgs,
	}

	return txsPrim, nil
}

func getFATransferAddressFrom(operationParams models.ContractOperationRequest, index int64) (from types.Address) {
	//Use contract self address as default from
	from = operationParams.ContractID
	if operationParams.TransferList[index].From != "" {
		from = operationParams.TransferList[index].From
	}

	return from
}

func buildVestingTxPrim(operationParams models.ContractOperationRequest) (prim *micheline.Prim, err error) {

	encodedVestingContract, err := operationParams.VestingID.MarshalBinary()
	if err != nil {
		return prim, err
	}

	var opCode micheline.OpCode
	var arg *micheline.Prim
	switch operationParams.Type {
	case models.VestingSetDelegate: //	setDelegate
		opCode = micheline.D_LEFT

		arg, err = buildDelegationPrim(operationParams.To)
		if err != nil {
			return prim, err
		}

	case models.VestingVest: // vest
		opCode = micheline.D_RIGHT

		//nat
		arg = &micheline.Prim{
			Type:   micheline.PrimInt,
			OpCode: micheline.T_INT,
			Int:    big.NewInt(int64(operationParams.Amount)),
		}
	default:
		return prim, errors.New("wrong type")
	}

	prim = &micheline.Prim{
		Type:   micheline.PrimBinary,
		OpCode: micheline.D_PAIR,
		Args: []*micheline.Prim{
			{
				Type:   micheline.PrimBytes,
				OpCode: micheline.T_BYTES,
				Bytes:  encodedVestingContract,
			},
			{
				Type:   micheline.PrimUnary,
				OpCode: opCode,
				Args:   []*micheline.Prim{arg},
			},
		},
	}

	return prim, err
}

func buildDelegationPrim(paramTo types.Address) (delegationPrim *micheline.Prim, err error) {

	if paramTo.IsEmpty() {
		//None
		delegationPrim = &micheline.Prim{
			Type:   micheline.PrimNullary,
			OpCode: micheline.D_NONE,
		}
		return delegationPrim, nil
	}

	encodedDestinationAddress, err := paramTo.MarshalBinary()
	if err != nil {
		return delegationPrim, err
	}

	//option(key_hash)
	delegationPrim = &micheline.Prim{
		Type:   micheline.PrimUnary,
		OpCode: micheline.D_SOME,
		Args: []*micheline.Prim{
			{
				Type:   micheline.PrimBytes,
				OpCode: micheline.T_BYTES,
				//Remove address byte
				Bytes: encodedDestinationAddress[1:],
			},
		},
	}

	return delegationPrim, nil
}
