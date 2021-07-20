package contract

import (
	"errors"
	"fmt"
	"tezosign/models"

	"blockwatch.cc/tzindex/micheline"
)

var (
	transferPath   = []micheline.OpCode{micheline.D_LEFT, micheline.D_LEFT, micheline.D_LEFT, micheline.D_LEFT}
	delegationPath = []micheline.OpCode{micheline.D_LEFT, micheline.D_LEFT, micheline.D_LEFT, micheline.D_RIGHT}
	//Transfer have branching inside
	transferFAPath = []micheline.OpCode{micheline.D_LEFT, micheline.D_LEFT, micheline.D_RIGHT, micheline.D_LEFT}
	//Vesting have branching inside
	vestingPath       = []micheline.OpCode{micheline.D_LEFT, micheline.D_LEFT, micheline.D_RIGHT, micheline.D_RIGHT}
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
		return pathParam, err
	}

	pathParam, err = getParamsByPath(actionParams, path)
	if err != nil {
		return pathParam, err
	}

	return pathParam, nil
}

func getParamsByPath(params *micheline.Prim, path []micheline.OpCode) (pathParam *micheline.Prim, err error) {

	pathParam = params

	var index int64
	for i := range path {

		//D_LEFT or D_RIGHT prim
		if pathParam.Type == micheline.PrimUnary || pathParam.Type == micheline.PrimUnaryAnno {
			pathParam = pathParam.Args[0]
			continue
		}

		//D_PAIR prim
		//EDO format of pair --> compress to normal format
		if pathParam.OpCode == micheline.D_PAIR && len(pathParam.Args) > 2 {

			pairsCount := (len(pathParam.Args) / 2) + 1

			for i := 0; i < pairsCount; i++ {

				pair, err := compressPair(pathParam.Args[len(pathParam.Args)-2:])
				if err != nil {
					return pathParam, err
				}

				//Remove last elem
				pathParam.Args = pathParam.Args[:len(pathParam.Args)-1]

				//Replace last element
				pathParam.Args[len(pathParam.Args)-1] = pair
			}

		}

		index = 0
		if path[i] == micheline.D_RIGHT {
			index = 1
		}

		pathParam = pathParam.Args[index]
	}

	return pathParam, nil
}

func compressPair(param []*micheline.Prim) (pair *micheline.Prim, err error) {
	if len(param) != 2 {
		return pair, errors.New("wrong args num")
	}

	args := make([]*micheline.Prim, len(param))
	copy(args, param)

	return &micheline.Prim{
		Type:   micheline.PrimBinary,
		OpCode: micheline.D_PAIR,
		Args:   args,
	}, nil
}

func getPathByType(actionType models.ActionType) (path []micheline.OpCode, err error) {
	switch actionType {
	case models.Transfer:
		path = transferPath
	case models.Delegation:
		path = delegationPath
	case models.StorageUpdate:
		path = updateStoragePath
	case models.FATransfer, models.FA2Transfer:
		path = transferFAPath
	case models.VestingVest, models.VestingSetDelegate:
		path = vestingPath
	case models.CustomPayload:
		path = customPayloadPath

	default:
		return nil, fmt.Errorf("unknown action")
	}

	return path, nil
}
