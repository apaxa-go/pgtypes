package pgtypes

import (
	"database/sql/driver"
	"fmt"
)

// sqlScan implements the sql.Scanner interface.
// Warning: Because of conflict between pgx.Scanner and sql.Scanner this method has prefixed with "sql".
// Warning: If you want to use UUID with sql (not pgx) just remove methods "Scan", "FormatCode" & "Encode" and rename "sqlScan" to "Scan" and "sqlValue" to "Value".
// Can parse:
//  bytes as raw UUID representations (as-is)
//  string/bytes as default UUID string representation
func (u *UUID) sqlScan(src interface{}) (err error) {
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

// sqlValue implements the driver.Valuer interface.
// Warning: see sqlScan warning.
func (u UUID) sqlValue() (driver.Value, error) {
	return u.String(), nil
}
