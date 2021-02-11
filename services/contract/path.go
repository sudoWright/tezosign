package contract

import (
	"errors"
	"fmt"
	"tezosign/models"

	"blockwatch.cc/tzindex/micheline"
)

var (
	delegationPath    = []micheline.OpCode{micheline.D_LEFT, micheline.D_LEFT, micheline.D_LEFT, micheline.D_RIGHT}
	transferPath      = []micheline.OpCode{micheline.D_LEFT, micheline.D_LEFT, micheline.D_LEFT, micheline.D_LEFT}
	transferFAPath    = []micheline.OpCode{micheline.D_LEFT, micheline.D_LEFT, micheline.D_RIGHT}
	customPayloadPath = []micheline.OpCode{micheline.D_LEFT, micheline.D_RIGHT}
	updateStoragePath = []micheline.OpCode{micheline.D_RIGHT}
)

func buildMichelsonPath(actionType models.ActionType, actionParams *micheline.Prim) (pathParam *micheline.Prim, err error) {
	path, err := getPathByType(actionType)
	if err != nil {
		return pathParam, nil
	}

	pathParam = actionParams

	for i := len(path) - 1; i >= 0; i-- {
		switch path[i] {
		case micheline.D_LEFT:
			pathParam = &micheline.Prim{
				Type:   micheline.PrimUnary,
				OpCode: micheline.D_LEFT,
				Args:   []*micheline.Prim{pathParam},
			}
		case micheline.D_RIGHT:
			pathParam = &micheline.Prim{
				Type:   micheline.PrimUnary,
				OpCode: micheline.D_RIGHT,
				Args:   []*micheline.Prim{pathParam},
			}
		default:
			return nil, fmt.Errorf("unknown OpCode")
		}
	}

	return pathParam, err
}

func getMichelsonParamsByActionType(actionType models.ActionType, actionParams *micheline.Prim) (pathParam *micheline.Prim, err error) {
	path, err := getPathByType(actionType)
	if err != nil {
		return pathParam, nil
	}

	pathParam = actionParams
	for i := range path {
		switch path[i] {
		case micheline.D_LEFT:
			if pathParam.OpCode != micheline.D_LEFT {
				return nil, errors.New("wrong path opcode D_LEFT")
			}

		case micheline.D_RIGHT:
			if pathParam.OpCode != micheline.D_RIGHT {
				return nil, errors.New("wrong path opcode ")
			}
		default:
			return nil, fmt.Errorf("unknown OpCode")
		}
		pathParam = pathParam.Args[0]
	}

	return pathParam, nil
}

func getPathByType(actionType models.ActionType) (path []micheline.OpCode, err error) {
	switch actionType {
	case models.Transfer:
		path = transferPath
	case models.Delegation:
		path = delegationPath
	case models.StorageUpdate:
		path = updateStoragePath
	case models.FATransfer:
		path = transferFAPath
	case models.CustomPayload:
		path = customPayloadPath

	default:
		return nil, fmt.Errorf("unknown action")
	}

	return path, nil
}
