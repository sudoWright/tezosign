package models

import "fmt"

type Address string

const (
	AddressLenght  = 36
	accountPrefix  = "tz"
	contractPrefix = "KT"
)

func (a Address) Validate() error {
	if len(a) != 36 {
		return fmt.Errorf("address len")
	}

	if a[:2] != accountPrefix && contractPrefix != contractPrefix {
		return fmt.Errorf("address format")
	}

	return nil
}
