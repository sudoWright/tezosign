package contract

import (
	"blockwatch.cc/tzindex/micheline"
	"math/big"
)

//Initial contract storage with 0 counter
const initialContractStorage = `{"prim": "Pair","args": [{"int": "0"},{"prim": "Pair","args": []}]}`

func BuildContractStorage(threshold uint, pubKeys []string) (resp []byte, err error) {
	storage := micheline.Prim{}
	err = storage.UnmarshalJSON([]byte(initialContractStorage))
	if err != nil {
		return nil, err
	}

	pubKeysPrim := make([]*micheline.Prim, len(pubKeys))
	for i := range pubKeys {
		pubKeysPrim[i] = &micheline.Prim{
			Type:   micheline.PrimString,
			OpCode: micheline.T_STRING,
			String: pubKeys[i],
		}
	}

	storage.Args[1].Args = []*micheline.Prim{
		{
			Type:   micheline.PrimInt,
			OpCode: micheline.T_INT,
			Int:    big.NewInt(int64(threshold)),
		},
		{
			Type:   micheline.PrimSequence,
			OpCode: micheline.I_SLICE,
			Args:   pubKeysPrim,
		},
	}

	return storage.MarshalJSON()
}
