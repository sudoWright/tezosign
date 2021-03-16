package contract

import (
	"tezosign/models"

	"blockwatch.cc/tzindex/micheline"
)

var fa12TransferFields = map[string]micheline.OpCode{"from": micheline.T_ADDRESS, "to": micheline.T_ADDRESS, "value": micheline.T_NAT}
var fa2TransferFields = map[string]micheline.OpCode{"from": micheline.T_ADDRESS, "txs": micheline.T_LIST, "to": micheline.T_ADDRESS, "tokenid": micheline.T_NAT, "amount": micheline.T_NAT}

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
