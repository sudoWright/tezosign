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
	"golang.org/x/crypto/blake2b"
	"golang.org/x/xerrors"
)

type PubKey tezosprotocol.PublicKey

func (a PubKey) String() string {
	return string(a)
}

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
}

func (a PubKey) MarshalBinary() ([]byte, error) {
	return tezosprotocol.PublicKey(a).MarshalBinary()
}

func (a *PubKey) UnmarshalBinary(data []byte) (err error) {
	var pubKey tezosprotocol.PublicKey
	err = pubKey.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	*a = PubKey(pubKey)

	return nil
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

const PubKeyHashLen = 20

//TODO add tests
func (a PubKey) Address() (Address, error) {
	pubKeyPrefix, pubKeyBytes, err := tezosprotocol.Base58CheckDecode(string(a))
	if err != nil {
		return "", err
	}

	var addressPrefix tezosprotocol.Base58CheckPrefix
	switch pubKeyPrefix {
	case tezosprotocol.PrefixEd25519PublicKey:
		addressPrefix = tezosprotocol.PrefixEd25519PublicKeyHash
	case tezosprotocol.PrefixP256PublicKey:
		addressPrefix = tezosprotocol.PrefixP256PublicKeyHash
	case tezosprotocol.PrefixSecp256k1PublicKey:
		addressPrefix = tezosprotocol.PrefixSecp256k1PublicKeyHash
	default:
		return "", fmt.Errorf("unsupported public key type %s", a)
	}

	// pubkey hash
	pubKeyHash, err := blake2b.New(PubKeyHashLen, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create blake2b hash: %w", err))
	}
	_, err = pubKeyHash.Write(pubKeyBytes)
	if err != nil {
		panic(fmt.Errorf("failed to write pubkey to hash: %w", err))
	}
	pubKeyHashBytes := pubKeyHash.Sum([]byte{})

	// base58check
	addr, err := tezosprotocol.Base58CheckEncode(addressPrefix, pubKeyHashBytes)
	if err != nil {
		return "", xerrors.Errorf("failed to base58check encode hash: %w", err)
	}

	return Address(addr), nil
}
