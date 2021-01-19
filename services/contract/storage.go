package contract

import (
	"blockwatch.cc/tzindex/micheline"
	"fmt"
	"math/big"
	"msig/types"
)

//Initial contract storage zero counter
const defaultCounter = 0

func BuildContractStorage(threshold uint, pubKeys []types.PubKey) (resp []byte, err error) {

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

func buildStorageMichelsonArgs(threshold int64, pubKeys []types.PubKey) (actionParams *micheline.Prim, err error) {
	var encodedPubKey []byte

	pubKeysPrim := make([]*micheline.Prim, len(pubKeys))
	for i := range pubKeys {
		encodedPubKey, err = pubKeys[i].MarshalBinary()
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
				OpCode: micheline.T_LIST,
				Args:   pubKeysPrim,
			},
		},
	}

	return actionParams, nil
}

type ContractStorageContainer struct {
	counter   int64
	threshold int64
	keys      []types.PubKey
	storage   *micheline.Prim
}

func NewContractStorageContainer(rawStorage string) (c ContractStorageContainer, err error) {

	storage := &micheline.Prim{}

	err = storage.UnmarshalJSON([]byte(rawStorage))
	if err != nil {
		return c, err
	}

	//Validate storage params
	//TODO make errors wrap
	if storage.OpCode != micheline.D_PAIR && len(storage.Args) != 2 {
		return c, fmt.Errorf("Wrong storage struct")
	}

	//Counter
	if storage.Args[0].OpCode != micheline.K_PARAMETER {
		return c, fmt.Errorf("Wrong counter type")
	}

	c.counter = storage.Args[0].Int.Int64()

	//pair(int * list(key)
	if storage.Args[1].OpCode != micheline.D_PAIR && len(storage.Args[1].Args) != 2 {
		return c, fmt.Errorf("Wrong storage pair struct")
	}

	//Threshold
	if storage.Args[1].Args[0].OpCode != micheline.K_PARAMETER {
		return c, fmt.Errorf("Wrong threshold type")
	}

	c.threshold = storage.Args[1].Args[0].Int.Int64()

	//Keys
	if storage.Args[1].Args[1].OpCode != micheline.T_LIST && len(storage.Args[1].Args[1].Args) < 2 {
		return c, fmt.Errorf("Wrong keys list")
	}

	c.keys = make([]types.PubKey, len(storage.Args[1].Args[1].Args))

	for i := range storage.Args[1].Args[1].Args {
		c.keys[i] = types.PubKey(storage.Args[1].Args[1].Args[i].String)
	}

	//Check storage input storage
	return c, nil
}

func (c ContractStorageContainer) Counter() int64 {
	return c.counter
}

func (c ContractStorageContainer) PubKeys() []types.PubKey {
	return c.keys
}

func (c ContractStorageContainer) Threshold() int64 {
	return c.threshold
}

func (c ContractStorageContainer) Contains(pubKey types.PubKey) (index int64, isFound bool) {
	for i := range c.keys {
		if c.keys[i] == pubKey {
			return int64(i), true
		}
	}
	return 0, false
}
