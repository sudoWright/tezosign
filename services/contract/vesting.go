package contract

import (
	"math/big"
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
