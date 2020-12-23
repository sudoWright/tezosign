package models

import script "blockwatch.cc/tzindex/micheline"

type BigMap struct {
	Code    *script.Prim `json:"code"`
	Storage *script.Prim `json:"storage"`
}
