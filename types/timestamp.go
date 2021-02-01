package types

import (
	"database/sql/driver"
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

func (t JSONTimestamp) Value() (driver.Value, error) {
	return t.Time(), nil
}
