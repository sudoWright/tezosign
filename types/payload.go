package types

import (
	"encoding/hex"
	"fmt"
)

const hexPrefix = "0x"

//Hexed payload
type Payload string

func (p Payload) Validate() error {
	//Hex payload should be at least len 2
	if len(p) < 2 {
		return fmt.Errorf("empty payload")
	}

	_, err := hex.DecodeString(p.WithoutPrefix())
	if err != nil {
		return fmt.Errorf("wrong hex format")
	}

	return nil
}

func (p Payload) String() string {
	return string(p)
}

func (p Payload) HasPrefix() bool {
	if p[0:2] == hexPrefix {
		return true
	}
	return false
}

func (p Payload) WithoutPrefix() string {
	if p.HasPrefix() {
		return string(p[2:])
	}
	return string(p)
}

func (p Payload) MarshalBinary() ([]byte, error) {
	return hex.DecodeString(p.WithoutPrefix())
}
