package types

import (
	"encoding/json"
	"time"
)

type JSONTimestamp time.Time

func (t JSONTimestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Unix())
}

func (t JSONTimestamp) Time() time.Time {
	return time.Time(t)
}
