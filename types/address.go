package types

import (
	"fmt"
	"github.com/anchorageoss/tezosprotocol/v2"
)

type Address tezosprotocol.ContractID

const (
	AddressLength  = 36
	accountPrefix  = "tz"
	contractPrefix = "KT"
)

func (a Address) Validate() (err error) {
	if len(a) != AddressLength {
		return fmt.Errorf("address len")
	}

	//Check that address
	if a[:2] != accountPrefix && contractPrefix != contractPrefix {
		return fmt.Errorf("address format")
	}

	//Check base58 format
	_, _, err = tezosprotocol.Base58CheckDecode(string(a))
	if err != nil {
		return fmt.Errorf("wrong base58 format")
	}

	return nil
}

func (a Address) String() string {
	return string(a)
}

func (a Address) MarshalBinary() ([]byte, error) {
	return tezosprotocol.ContractID(a).MarshalBinary()
}
