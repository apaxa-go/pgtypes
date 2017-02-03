package pgtypes

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements the sql.Scanner interface.
// Can parse:
//  bytes as raw UUID representations (as-is)
//  string/bytes as default UUID string representation
func (u *UUID) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		switch len(src) {
		case UUIDLen:
			*u, err = ParseUUIDBytes(src)
			return
		case UUIDStringLen:
			*u, err = ParseUUID(string(src))
			return
		}
	case string:
		*u, err = ParseUUID(src)
		return
	}

	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

// Value implements the driver.Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}
