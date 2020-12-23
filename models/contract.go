package models

import "fmt"

type ContractRequest struct {
	Threshold uint      `json:"threshold"`
	Addresses []Address `json:"addresses"`
}

func (r ContractRequest) Validate() (err error) {
	if r.Threshold == 0 {
		return fmt.Errorf("zero threshold")
	}

	if len(r.Addresses) == 0 {
		return fmt.Errorf("empty addresses")
	}

	for i := range r.Addresses {
		err = r.Addresses[i].Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
