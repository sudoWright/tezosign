package contract

import (
	"blockwatch.cc/tzindex/micheline"
)

func CheckTransferMethod(script *micheline.Script) (ok bool) {

	entrypoints, err := script.Entrypoints(true)
	if err != nil {
		return false
	}

	transferEntrypoint, ok := entrypoints["transfer"]
	if !ok {
		return false
	}

	//FA Transfer method checks
	if transferEntrypoint.Prim.OpCode != micheline.T_PAIR {
		return false
	}

	//FA transfer(from: to: address value: nat)

	//From Address
	if transferEntrypoint.Prim.Args[0].OpCode != micheline.T_ADDRESS {
		return false
	}

	// Pair to: address value: nat
	if transferEntrypoint.Prim.Args[1].OpCode != micheline.T_PAIR {
		return false
	}

	//to : address
	if transferEntrypoint.Prim.Args[1].Args[0].OpCode != micheline.T_ADDRESS {
		return false
	}

	//value: nat
	if transferEntrypoint.Prim.Args[1].Args[1].OpCode != micheline.T_NAT {
		return false
	}

	return true
}
