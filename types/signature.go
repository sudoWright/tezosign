package types

import "github.com/anchorageoss/tezosprotocol/v2"

type Signature tezosprotocol.Signature

func (s Signature) Validate() (err error) {
	//Todo add validation
	return nil
}

func (s Signature) String() string {
	return string(s)
}

func (s Signature) IsEmpty() bool {
	return len(s) == 0
}

func (s Signature) MarshalBinary() (bt []byte, err error) {
	return tezosprotocol.Signature(s).MarshalBinary()
}
