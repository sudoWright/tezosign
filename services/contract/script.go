package contract

import (
	"math/big"
	"tezosign/models"

	"blockwatch.cc/tzindex/base58"
	"golang.org/x/crypto/blake2b"

	"blockwatch.cc/tzindex/micheline"
)

var (
	fa12TransferFields = map[string]micheline.OpCode{"from": micheline.T_ADDRESS, "to": micheline.T_ADDRESS, "value": micheline.T_NAT}
	fa2TransferFields  = map[string]micheline.OpCode{"from": micheline.T_ADDRESS, "txs": micheline.T_LIST, "to": micheline.T_ADDRESS, "tokenid": micheline.T_NAT, "amount": micheline.T_NAT}
	scriptHashPrefix   = []byte{13, 44, 64, 27}
)

func CheckFATransferMethod(script *micheline.Script, faType models.AssetType) (ok bool) {

	entrypoints, err := script.Entrypoints(true)
	if err != nil {
		return false
	}

	transferEntrypoint, ok := entrypoints["transfer"]
	if !ok {
		return false
	}

	//Init transfer method fields
	e, err := InitAnnotsEntrypoints(transferEntrypoint.Prim)
	if err != nil {
		return false
	}

	//Check FA1.2
	faFieldsTransfer := fa12TransferFields
	if faType == models.TypeFA2 {
		faFieldsTransfer = fa2TransferFields
	}

	err = checkFields(e, faFieldsTransfer)
	if err != nil {
		return false
	}

	return true
}

func GetBigMapKeyHash(tokenID int64) (hash string, err error) {
	p := micheline.Prim{
		Type:   micheline.PrimInt,
		OpCode: micheline.T_INT,
		Int:    big.NewInt(tokenID),
	}

	bt, err := p.MarshalBinary()
	if err != nil {
		return hash, err
	}

	bt = append([]byte{TextWatermark}, bt...)

	hh := blake2b.Sum256(bt)

	hash = base58.CheckEncode(hh[:], scriptHashPrefix)
	return
}
