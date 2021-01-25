package models

import "msig/types"

type SignatureReq struct {
	PubKey    types.PubKey    `json:"pub_key"`
	Payload   types.Payload   `json:"payload"`
	Signature types.Signature `json:"signature"`
}

func (r SignatureReq) Validate() (err error) {

	err = r.Payload.Validate()
	if err != nil {
		return err
	}

	err = r.PubKey.Validate()
	if err != nil {
		return err
	}

	err = r.Signature.Validate()
	if err != nil {
		return err
	}

	return nil
}
