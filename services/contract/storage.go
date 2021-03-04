package contract

import (
	"errors"
	"math/big"
	"tezosign/types"

	"blockwatch.cc/tzindex/micheline"
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

const (
	counterEntrypoint   = "counter"
	keysEntrypoint      = "keys"
	thresholdEntrypoint = "threshold"
)

var contractStorageEntrypoints = map[string]micheline.OpCode{counterEntrypoint: micheline.T_NAT, keysEntrypoint: micheline.T_LIST, thresholdEntrypoint: micheline.T_NAT}

func NewContractStorageContainer(script micheline.Script) (c ContractStorageContainer, err error) {

	e, err := InitStorageAnnotsEntrypoints(script.Code.Storage)
	if err != nil {
		return c, err
	}

	err = checkStorage(e, contractStorageEntrypoints)
	if err != nil {
		return c, err
	}

	c.storage = script.Storage

	counter, err := GetStorageValue(e[counterEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}

	c.counter = counter.Int.Int64()

	threshold, err := GetStorageValue(e[thresholdEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}

	c.threshold = threshold.Int.Int64()

	keys, err := GetStorageValue(e[keysEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}

	c.keys = make([]types.PubKey, len(keys.Args))

	for i := range keys.Args {
		err = c.keys[i].UnmarshalBinary(keys.Args[i].Bytes)
		if err != nil {
			return c, err
		}
	}

	return
}

func checkStorage(e Entrypoints, contractStorageEntrypoints map[string]micheline.OpCode) error {

	var entrypoint Entrypoint
	var ok bool
	for eName, opCode := range contractStorageEntrypoints {
		if entrypoint, ok = e[eName]; !ok {
			return errors.New("entrypoint not found")
		}

		if entrypoint.OpCode != opCode {
			return errors.New("wrong entrypoint opcode")
		}
	}

	return nil
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

func InitStorageAnnotsEntrypoints(codeStorage *micheline.Prim) (e Entrypoints, err error) {
	if len(codeStorage.Args) == 0 {
		return e, errors.New("wrong code storage")
	}

	e = make(Entrypoints)
	dfs(e, newVertex(codeStorage.Args[0]), []micheline.OpCode{})

	return
}

func GetStorageValue(e Entrypoint, storage *micheline.Prim) (*micheline.Prim, error) {
	return getParamsByPath(storage, e.Branch)
}
