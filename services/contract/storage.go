package contract

import (
	"blockwatch.cc/tzindex/micheline"
	"math/big"
)

//Initial contract storage zero counter
const defaultCounter = 0

func BuildContractStorage(threshold uint, pubKeys []string) (resp []byte, err error) {

	storageParams, err := buildStorageMichelsonArgs(int64(threshold), pubKeys)
	if err != nil {
		return nil, err
	}

	storage := &micheline.Prim{
		Type:   micheline.PrimBinary,
		OpCode: micheline.D_PAIR,
		Args: []*micheline.Prim{
			//Counter
			{
				Type:   micheline.PrimInt,
				OpCode: micheline.T_INT,
				Int:    big.NewInt(defaultCounter),
			},
			//Pair (nat * list(key))
			storageParams,
		},
	}

	return storage.MarshalJSON()
}

func buildStorageMichelsonArgs(threshold int64, pubKeys []string) (actionParams *micheline.Prim, err error) {
	var encodedPubKey []byte

	pubKeysPrim := make([]*micheline.Prim, len(pubKeys))
	for i := range pubKeys {
		encodedPubKey, err = encodeBase58ToPrimBytes(pubKeys[i])
		if err != nil {
			return actionParams, err
		}

		pubKeysPrim[i] = &micheline.Prim{
			Type:   micheline.PrimBytes,
			OpCode: micheline.T_BYTES,
			Bytes:  encodedPubKey,
		}
	}

	//Pair (int * list(key))
	actionParams = &micheline.Prim{
		Type:   micheline.PrimBinary,
		OpCode: micheline.D_PAIR,
		Args: []*micheline.Prim{
			// Threshold
			{
				Type:   micheline.PrimInt,
				OpCode: micheline.T_INT,
				Int:    big.NewInt(threshold),
			},
			// list
			{
				Type:   micheline.PrimSequence,
				OpCode: micheline.I_SLICE,
				Args:   pubKeysPrim,
			},
		},
	}

	return actionParams, nil
}
