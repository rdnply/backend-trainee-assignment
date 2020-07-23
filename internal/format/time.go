package format

import (
	"database/sql/driver"
	"time"
)

type JSONTime struct {
	time.Time
}

func NewTime(t time.Time) JSONTime {
	return JSONTime{t}
}

func (jt JSONTime) MarshalJSON() ([]byte, error) {
	s := jt.Time.Format(time.RFC3339)
	return []byte(`"` + s + `"`), nil
}

func (jt *JSONTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	jt.Time = value.(time.Time)
	return nil
}

func (jt JSONTime) Value() (driver.Value, error) {
	return jt.Time.String(), nil
}
