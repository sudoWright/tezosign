package types

import (
	"fmt"
	"github.com/anchorageoss/tezosprotocol/v2"
)

type Signature tezosprotocol.Signature

func (s Signature) Validate() (err error) {
	b58prefix, _, err := tezosprotocol.Base58CheckDecode(string(s))
	if err != nil {
		return fmt.Errorf("wrong signature format")
	}

	switch b58prefix {
	case tezosprotocol.PrefixEd25519Signature, tezosprotocol.PrefixSecp256k1Signature, tezosprotocol.PrefixP256Signature, tezosprotocol.PrefixGenericSignature:
		return nil
	default:
		return fmt.Errorf("wrong signature prefix")
	}

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
