package contract

import (
	"github.com/anchorageoss/tezosprotocol/v2"
)

const (
	publicKeyHashPrefix = 0x00
	contractHashPrefix  = 0x01

	ed25519Prefix   = 0x00
	secp256k1Prefix = 0x01
	p256Prefix      = 0x02
)

func encodeBase58ToPrimBytes(base58 string) (encodedBytes []byte, err error) {
	var prefix tezosprotocol.Base58CheckPrefix
	prefix, encodedBytes, err = tezosprotocol.Base58CheckDecode(base58)
	if err != nil {
		return encodedBytes, err
	}

	switch prefix {
	case tezosprotocol.PrefixEd25519PublicKey:
		encodedBytes = append([]byte{ed25519Prefix}, encodedBytes...)
	case tezosprotocol.PrefixEd25519PublicKeyHash:
		encodedBytes = append([]byte{publicKeyHashPrefix, ed25519Prefix}, encodedBytes...)
	case tezosprotocol.PrefixSecp256k1PublicKey:
		encodedBytes = append([]byte{secp256k1Prefix}, encodedBytes...)
	case tezosprotocol.PrefixSecp256k1PublicKeyHash:
		encodedBytes = append([]byte{publicKeyHashPrefix, secp256k1Prefix}, encodedBytes...)
	case tezosprotocol.PrefixP256PublicKey:
		encodedBytes = append([]byte{p256Prefix}, encodedBytes...)
	case tezosprotocol.PrefixP256PublicKeyHash:
		encodedBytes = append([]byte{publicKeyHashPrefix, p256Prefix}, encodedBytes...)
	case tezosprotocol.PrefixContractHash:
		encodedBytes = append([]byte{contractHashPrefix}, encodedBytes...)
		encodedBytes = append(encodedBytes, 0x00)
	}

	return encodedBytes, nil
}
