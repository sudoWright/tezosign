package contract

import (
	"strings"

	"blockwatch.cc/tzindex/micheline"
)

type Entrypoints map[string]Entrypoint

type Entrypoint struct {
	Id     int                `json:"id"`
	Branch []micheline.OpCode `json:"branch"`
	OpCode micheline.OpCode   `json:"type"`
	Prim   *micheline.Prim    `json:"prim,omitempty"`
}

type vertex struct {
	visited    bool
	value      *micheline.Prim
	neighbours []*micheline.Prim
}

func newVertex(prim *micheline.Prim) *vertex {
	if prim == nil {
		return nil
	}

	return &vertex{
		visited:    false,
		value:      prim,
		neighbours: prim.Args,
	}
}

func dfs(e Entrypoints, vertex *vertex, path []micheline.OpCode) {
	if vertex == nil || vertex.visited {
		return
	}

	vertex.visited = true

	pathInit(e, vertex.value, path)

	for i, v := range vertex.neighbours {
		var pathIndex micheline.OpCode
		switch i {
		case 0:
			pathIndex = micheline.D_LEFT
			if vertex.neighbours[0].OpCode == micheline.D_RIGHT {
				pathIndex = 1
			}
		case 1:
			pathIndex = micheline.D_RIGHT
		default:
		}

		dfs(e, newVertex(v), append(path, pathIndex))
	}
}

func pathInit(e Entrypoints, prim *micheline.Prim, path []micheline.OpCode) {
	//For now process only values with annotation
	if prim == nil || prim.GetAnno() == "" {
		return
	}

	//Init new list
	p := make([]micheline.OpCode, len(path))
	copy(p, path)
	anno := prim.GetAnno()

	//Dexter delphi and mainnet have diff anno
	anno = strings.ToLower(anno)
	anno = strings.ReplaceAll(anno, "_", "")

	e[anno] = Entrypoint{
		Id:     len(e),
		Branch: p,
		OpCode: prim.OpCode,
		Prim:   prim,
	}
}
