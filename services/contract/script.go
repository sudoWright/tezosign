package contract

import (
	"blockwatch.cc/tzindex/micheline"
)

func CheckTransferMetod(script *micheline.Script) (ok bool) {

	entrypoints, err := script.Entrypoints(true)
	if err != nil {
		return false
	}

	/*transferEntrypoint*/
	_, ok = entrypoints["transfer"]
	if !ok {
		return false
	}

	return true
}
