package format

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"
)

type NullTime struct {
	V sql.NullTime
}

func NewNullTime() *NullTime {
	return &NullTime{V: sql.NullTime{Time: time.Now().UTC(), Valid: true}}
}

const layout = "2006-01-02T15:04:05Z"

func (nt *NullTime) Scan(value interface{}) error {
	var t sql.NullTime
	if err := t.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*nt = NullTime{V: sql.NullTime{Time: t.Time, Valid: false}}
	} else {
		t := t.Time.UTC()
		s := t.Format(time.RFC3339)
		t, err := time.Parse(layout, s)
		if err != nil {
			nt.V.Valid = false
			return err
		}

		*nt = NullTime{V: sql.NullTime{Time: t, Valid: true}}
	}

	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.V.Valid {
		return nil, nil
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}

	nt.V.Valid = true

	return nt.V.Time.In(loc), nil
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.V.Valid {
		return nil, nil
	}

	t := nt.V.Time.UTC()
	s := t.Format(time.RFC3339)

	t, err := time.Parse(layout, s)
	if err != nil {
		nt.V.Valid = false

		return nil, err
	}

	s = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return []byte(`"` + s + `"`), nil
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	s := string(b)

	t, err := time.Parse(`"`+layout+`"`, s)
	if err != nil {
		nt.V.Valid = false
		return err
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}

	nt.V.Time = t.In(loc)
	nt.V.Valid = true

	return nil
}
