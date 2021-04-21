package contract

import (
	"errors"
	"math/big"
	"tezosign/models"
	"tezosign/types"

	"blockwatch.cc/tzindex/micheline"
)

//parameter (or (option %setDelegate key_hash) (nat %vest))

//storage (pair (pair %wrapped (address %target) (address %delegateAdmin))
//              (pair (nat %vested)
//              (pair %schedule (timestamp %epoch)
//              (pair (nat %secondsPerTick) (nat %tokensPerTick)))))

func BuildVestingContractStorage(vestingAddress, delegateAdmin types.Address, timestamp, secondsPerTick, tokensPerTick uint64) (resp []byte, err error) {

	encodedVestingAddress, err := vestingAddress.MarshalBinary()
	if err != nil {
		return nil, err
	}

	encodedDelegateAdmin, err := delegateAdmin.MarshalBinary()
	if err != nil {
		return nil, err
	}

	storage := &micheline.Prim{
		Type:   micheline.PrimBinary,
		OpCode: micheline.D_PAIR,
		Args: []*micheline.Prim{
			{
				//(pair %wrapped (address %target) (address %delegateAdmin))
				Type:   micheline.PrimBinary,
				OpCode: micheline.D_PAIR,
				Args: []*micheline.Prim{
					//Target address
					{
						Type:   micheline.PrimBytes,
						OpCode: micheline.T_BYTES,
						Bytes:  encodedVestingAddress,
					},
					//Delegate admin address
					{
						Type:   micheline.PrimBytes,
						OpCode: micheline.T_BYTES,
						Bytes:  encodedDelegateAdmin,
					},
				},
			},
			{
				//(pair (nat %vested) ...
				Type:   micheline.PrimBinary,
				OpCode: micheline.D_PAIR,
				Args: []*micheline.Prim{
					// nat %vested
					{
						Type:   micheline.PrimInt,
						OpCode: micheline.T_NAT,
						Int:    big.NewInt(0),
					},
					//(pair %schedule (timestamp %epoch) ...
					{
						//(pair (nat %vested) ...
						Type:   micheline.PrimBinary,
						OpCode: micheline.D_PAIR,
						Args: []*micheline.Prim{
							// timestamp %epoch
							{
								Type:   micheline.PrimInt,
								OpCode: micheline.T_TIMESTAMP,
								Int:    big.NewInt(int64(timestamp)),
							},
							//(pair (nat %secondsPerTick) (nat %tokensPerTick))
							{
								Type:   micheline.PrimBinary,
								OpCode: micheline.D_PAIR,
								Args: []*micheline.Prim{
									// nat %secondsPerTick
									{
										Type:   micheline.PrimInt,
										OpCode: micheline.T_NAT,
										Int:    big.NewInt(int64(secondsPerTick)),
									},
									// nat %tokensPerTick
									{
										Type:   micheline.PrimInt,
										OpCode: micheline.T_NAT,
										Int:    big.NewInt(int64(tokensPerTick)),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return storage.MarshalJSON()
}

const (
	setDelegateEntrypoint = "setDelegate"
	vestEntrypoint        = "vest"
)

func VestingContractParamAndEntrypoint(req models.VestingContractOperation) (arg []byte, entrypoint string, err error) {

	var prim *micheline.Prim
	switch req.Type {
	case models.VestingSetDelegate:

		prim, err = buildDelegationPrim(req.To)
		if err != nil {
			return nil, "", err
		}

		entrypoint = setDelegateEntrypoint
	case models.VestingVest:
		prim = &micheline.Prim{
			Type:   micheline.PrimInt,
			OpCode: micheline.T_INT,
			Int:    big.NewInt(int64(req.Ticks)),
		}
		entrypoint = vestEntrypoint
	default:
		return nil, "", errors.New("wrong request type")
	}

	arg, err = prim.MarshalJSON()
	if err != nil {
		return nil, "", err
	}

	return arg, entrypoint, nil
}

type VestingContractStorageContainer struct {
	VestingAddress types.Address
	DelegateAdmin  types.Address
	VestedTicks    uint64
	Timestamp      uint64
	SecondsPerTick uint64
	TokensPerTick  uint64
	storage        *micheline.Prim
}

const (
	targetEntrypoint         = "target"
	delegateAdminEntrypoint  = "delegateadmin"
	vestedEntrypoint         = "vested"
	epochEntrypoint          = "epoch"
	secondsPerTickEntrypoint = "secondspertick"
	tokensPerTickEntrypoint  = "tokenspertick"
)

var vestingContractStorageEntrypoints = map[string]micheline.OpCode{
	targetEntrypoint:         micheline.T_ADDRESS,
	delegateAdminEntrypoint:  micheline.T_ADDRESS,
	vestedEntrypoint:         micheline.T_NAT,
	epochEntrypoint:          micheline.T_TIMESTAMP,
	secondsPerTickEntrypoint: micheline.T_NAT,
	tokensPerTickEntrypoint:  micheline.T_NAT,
}

func NewVestingContractStorageContainer(script micheline.Script) (c VestingContractStorageContainer, err error) {

	e, err := InitAnnotsEntrypoints(script.Code.Storage)
	if err != nil {
		return c, err
	}

	err = checkFields(e, vestingContractStorageEntrypoints)
	if err != nil {
		return c, err
	}

	c.storage = script.Storage

	address, err := GetStorageValue(e[targetEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}

	err = c.VestingAddress.UnmarshalBinary(address.Bytes)
	if err != nil {
		return c, err
	}

	address, err = GetStorageValue(e[delegateAdminEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}

	err = c.DelegateAdmin.UnmarshalBinary(address.Bytes)
	if err != nil {
		return c, err
	}

	amount, err := GetStorageValue(e[vestedEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}

	c.VestedTicks = amount.Int.Uint64()

	amount, err = GetStorageValue(e[epochEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}
	c.Timestamp = amount.Int.Uint64()

	amount, err = GetStorageValue(e[secondsPerTickEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}
	c.SecondsPerTick = amount.Int.Uint64()

	amount, err = GetStorageValue(e[tokensPerTickEntrypoint], script.Storage)
	if err != nil {
		return c, err
	}
	c.TokensPerTick = amount.Int.Uint64()

	return
}
