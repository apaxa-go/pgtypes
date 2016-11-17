package uuid

import (
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx"
)

func init() {
	// Register UUID type in pgx as binary-compatible.
	// This may cause error if type other when apaxa-io UUID will be used with pgx for uuid storing.
	pgx.DefaultTypeFormats["uuid"] = pgx.BinaryFormatCode
}

// Scan implements the pgx.Scanner interface.
func (u *UUID) Scan(vr *pgx.ValueReader) error {
	if vr.Type().DataType != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("UUID.Scan cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		return pgx.SerializationError("UUID.Scan cannot parse NULL value")
	}

	switch vr.Type().FormatCode {
	case pgx.TextFormatCode:
		if vr.Len() != uuidStringLen {
			return pgx.SerializationError(fmt.Sprintf("Received UUID string with invalid length: %d", vr.Len()))
		}
		if err := u.ParseString(vr.ReadString(vr.Len())); err != nil {
			return pgx.SerializationError(fmt.Sprintf("Received invalid UUID string: %v", err.Error()))
		}
	case pgx.BinaryFormatCode:
		if vr.Len() != uuidLen {
			return pgx.SerializationError(fmt.Sprintf("Received UUID with invalid length: %d", vr.Len()))
		}
		if err := u.ParseBytes(vr.ReadBytes(uuidLen)); err != nil {
			return pgx.SerializationError(fmt.Sprintf("Received invalid UUID: %v", err.Error()))
		}
	default:
		return fmt.Errorf("unknown format %v", vr.Type().FormatCode)
	}

	return vr.Err()
}

// FormatCode implements the pgx.Encoder interface.
func (u UUID) FormatCode() int16 {
	return pgx.BinaryFormatCode
}

// Encode implements the pgx.Encoder interface.
func (u UUID) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("UUID.Encode cannot encode into OID %d", oid))
	}

	w.WriteInt32(uuidLen)
	w.WriteBytes(u[:])

	return nil
}

// sqlScan implements the sql.Scanner interface.
// Warning: Because of conflict between pgx.Scanner and sql.Scanner this method has prefixed with "sql".
// Warning: If you want to use UUID with sql (not pgx) just remove methods "Scan", "FormatCode" & "Encode" and rename "sqlScan" to "Scan" and "sqlValue" to "Value".
// Can parse:
//  bytes as raw UUID representations (as-is)
//  string/bytes as default UUID string representation
func (u *UUID) sqlScan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		switch len(src) {
		case uuidLen:
			return u.ParseBytes(src)
		case uuidStringLen:
			return u.ParseString(string(src))
		}
	case string:
		return u.ParseString(src)
	}

	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

// sqlValue implements the driver.Valuer interface.
// Warning: see sqlScan warning.
func (u UUID) sqlValue() (driver.Value, error) {
	return u.String(), nil
}
