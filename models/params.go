package models

import "fmt"

type CommonParams struct {
	Limit  int
	Offset int
}

const MaxLimitSize = 500

func (c CommonParams) Validate() error {
	if c.Limit == 0 {
		return fmt.Errorf("limit not present")
	}

	if c.Limit > MaxLimitSize {
		return fmt.Errorf("overlimit")
	}

	return nil
}
