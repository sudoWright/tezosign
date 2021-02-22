package contract

import (
	"errors"
	"math/big"

	"blockwatch.cc/tzindex/micheline"
)

const (
	tokenPoolAnno = "tokenpool"
	xtzPoolAnno   = "xtzpool"
)

func GetDexterContractTokenPool(e Entrypoints, storage *micheline.Prim) (*big.Int, error) {
	entrypoint, ok := e[tokenPoolAnno]
	if !ok {
		return nil, errors.New("entrypoint not found")
	}

	prim, err := GetStorageValue(entrypoint, storage)
	if err != nil {
		return nil, err
	}

	return prim.Int, nil
}

func GetDexterContractXTZPool(e Entrypoints, storage *micheline.Prim) (*big.Int, error) {
	entrypoint, ok := e[xtzPoolAnno]
	if !ok {
		return nil, errors.New("entrypoint not found")
	}

	prim, err := GetStorageValue(entrypoint, storage)
	if err != nil {
		return nil, err
	}

	return prim.Int, nil
}
