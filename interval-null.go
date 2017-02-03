package pgtypes

import (
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx"
)

// NullInterval represents an Interval that may be null.
// It implements the pgx.PgxScanner and pgx.Encoder interfaces so it may be used both as an argument to Query[Row] and a destination for ScanPgx.
//
// If Valid is false then the value is NULL.
type NullInterval struct {
	Interval Interval
	Valid    bool
}

// ScanPgx implements the pgx.PgxScanner interface.
func (n *NullInterval) ScanPgx(vr *pgx.ValueReader) error {
	if vr.Type().DataType != IntervalOid {
		return pgx.SerializationError(fmt.Sprintf("NullInterval.ScanPgx cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		n.Interval, n.Valid = Interval{0, 0, 0, IntervalPgPrecision}, false
		return nil
	}

	n.Valid = true
	return n.Interval.ScanPgx(vr)
}

// FormatCode implements the pgx.Encoder interface.
func (n NullInterval) FormatCode() int16 { return pgx.BinaryFormatCode }

// Encode implements the pgx.Encoder interface.
func (n NullInterval) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != IntervalOid {
		return pgx.SerializationError(fmt.Sprintf("NullInterval.Encode cannot encode into OID %d", oid))
	}

	if !n.Valid {
		w.WriteInt32(-1)
		return nil
	}

	return n.Interval.Encode(w, oid)
}

// Nullable returns valid NullInterval with Interval i.
func (i Interval) Nullable() NullInterval {
	return NullInterval{i, true}
}

// Scan implements the sql.Scanner interface.
func (n *NullInterval) Scan(src interface{}) (err error) {
	if src == nil {
		*n = NullInterval{Interval: Interval{precision: IntervalPgPrecision}}
		return nil
	}
	err = n.Interval.Scan(src)
	n.Valid = err == nil
	return
}

// Value implements the driver.Valuer interface.
func (n NullInterval) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Interval.Value()
}
