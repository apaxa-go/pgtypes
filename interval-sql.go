package pgtypes

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// Scan implements the sql.Scanner interface.
func (i *Interval) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		*i, err = ParseInterval(string(src), IntervalPgPrecision)
		if err != nil {
			err = errors.New("interval: " + err.Error())
		}
		return
	case string:
		*i, err = ParseInterval(src, IntervalPgPrecision)
		if err != nil {
			err = errors.New("interval: " + err.Error())
		}
		return
	}

	return fmt.Errorf("interval: cannot convert %T to Interval", src)
}

// Value implements the driver.Valuer interface.
func (i Interval) Value() (driver.Value, error) {
	return i.String(), nil
}
