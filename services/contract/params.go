package contract

import (
	"blockwatch.cc/tzindex/micheline"
	"fmt"
	"math/big"
	"tezosign/models"
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

	case models.Delegation:
		if operationParams.To == "" {
			//None
			actionParams = &micheline.Prim{
				Type:   micheline.PrimNullary,
				OpCode: micheline.D_NONE,
			}
			break
		}
		encodedDestinationAddress, err := operationParams.To.MarshalBinary()
		if err != nil {
			return actionParams, err
		}

		//option(key_hash)
		actionParams = &micheline.Prim{
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
	case models.StorageUpdate:
		actionParams, err = buildStorageMichelsonArgs(int64(operationParams.Threshold), operationParams.Keys)
		if err != nil {
			return actionParams, err
		}
	case models.FATransfer:
		encodedAssetAddress, err := operationParams.AssetID.MarshalBinary()
		if err != nil {
			return actionParams, err
		}

		//Use contract self address as default from
		from := operationParams.ContractID
		if operationParams.From != "" {
			from = operationParams.From
		}

		encodedFromAddress, err := from.MarshalBinary()
		if err != nil {
			return actionParams, err
		}

		encodedDestinationAddress, err := operationParams.To.MarshalBinary()
		if err != nil {
			return actionParams, err
		}

		//(pair address (pair  address (pair address nat)))
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
									Int:    big.NewInt(int64(operationParams.Amount)),
								},
							},
						},
					},
				},
			},
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
