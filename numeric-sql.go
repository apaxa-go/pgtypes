package pgtypes

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements the sql.Numeric interface.
func (n *Numeric) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		if !n.setString(string(src)) {
			err = fmt.Errorf(`numeric: cannot convert []byte %v to Numeric`, src)
		}
		return
	case string:
		if !n.setString(src) {
			err = fmt.Errorf(`numeric: cannot convert string "%v" to Numeric`, src)
		}
		return
	}

	return fmt.Errorf("numeric: cannot convert %T to Numeric", src)
}

// Value implements the driver.Valuer interface.
func (n Numeric) Value() (driver.Value, error) {
	return n.String(), nil
}
