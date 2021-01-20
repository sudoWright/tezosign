package types

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"fmt"
	"github.com/anchorageoss/tezosprotocol/v2"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type PubKey tezosprotocol.PublicKey

func (a PubKey) Validate() (err error) {
	b58prefix, _, err := tezosprotocol.Base58CheckDecode(string(a))
	if err != nil {
		return fmt.Errorf("wrong pubKey format")
	}

	switch b58prefix {
	case tezosprotocol.PrefixEd25519PublicKey, tezosprotocol.PrefixSecp256k1PublicKey, tezosprotocol.PrefixP256PublicKey:
		return nil
	default:
		return fmt.Errorf("wrong pubKey prefix")
	}

	return nil
}

func (a PubKey) MarshalBinary() ([]byte, error) {
	return tezosprotocol.PublicKey(a).MarshalBinary()
}

func (a PubKey) CryptoPublicKey() (crypto.PublicKey, error) {
	b58prefix, b58decoded, err := tezosprotocol.Base58CheckDecode(string(a))
	if err != nil {
		return nil, err
	}
	switch b58prefix {
	case tezosprotocol.PrefixEd25519PublicKey:
		return ed25519.PublicKey(b58decoded), nil
	case tezosprotocol.PrefixSecp256k1PublicKey:
		return secp256k1.PubKey(b58decoded), nil
	case tezosprotocol.PrefixP256PublicKey:

		x, y := elliptic.UnmarshalCompressed(elliptic.P256(), b58decoded)

		return ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		}, nil
	default:
		return nil, errors.Errorf("unexpected base58check prefix: %s", b58prefix)
	}
}
