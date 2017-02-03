package pgtypes

import (
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx"
)

// NullNumeric represents an Numeric that may be NULL.
// NullNumeric implements the pgx.PgxScanner and pgx.Encoder interfaces so it may be used both as an argument to Query[Row] and a destination for ScanPgx.
//
// If Valid is false then the value is NULL.
type NullNumeric struct {
	Numeric Numeric
	Valid   bool // Valid is true if Numeric is not NULL
}

// ScanPgx implements the pgx.PgxScanner interface.
func (n *NullNumeric) ScanPgx(vr *pgx.ValueReader) error {
	if vr.Type().DataType != NumericOid {
		return pgx.SerializationError(fmt.Sprintf("NullNumeric.ScanPgx cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		n.Valid = false
		n.Numeric.SetZero()
		return nil
	}

	n.Valid = true
	return n.Numeric.ScanPgx(vr)
}

// FormatCode implements the pgx.Encoder interface.
func (n NullNumeric) FormatCode() int16 { return pgx.BinaryFormatCode }

// Encode implements the pgx.Encoder interface.
func (n NullNumeric) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != NumericOid {
		return pgx.SerializationError(fmt.Sprintf("NullNumeric.Encode cannot encode into OID %d", oid))
	}

	if !n.Valid {
		w.WriteInt32(-1)
		return nil
	}

	return n.Numeric.Encode(w, oid)
}

// Nullable returns valid NullNumeric with Numeric n.
func (n Numeric) Nullable() NullNumeric {
	return NullNumeric{n, true}
}

// Scan implements the sql.Scanner interface.
func (n *NullNumeric) Scan(src interface{}) (err error) {
	if src == nil {
		*n = NullNumeric{}
		return nil
	}
	err = n.Numeric.Scan(src)
	n.Valid = err == nil
	return
}

// Value implements the driver.Valuer interface.
func (n NullNumeric) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Numeric.Value()
}
