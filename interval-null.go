package pgtypes

import (
	"fmt"
	"github.com/jackc/pgx"
)

// NullInterval represents an Interval that may be null.
// It implements the pgx.Scanner and pgx.Encoder interfaces so it may be used both as an argument to Query[Row] and a destination for Scan.
//
// If Valid is false then the value is NULL.
type NullInterval struct {
	Interval Interval
	Valid    bool
}

// Scan implements the pgx.Scanner interface.
func (n *NullInterval) Scan(vr *pgx.ValueReader) error {
	if vr.Type().DataType != IntervalOid {
		return pgx.SerializationError(fmt.Sprintf("NullInterval.Scan cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		n.Interval, n.Valid = Interval{0, 0, 0, PgPrecision}, false
		return nil
	}

	n.Valid = true
	return n.Interval.Scan(vr)
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
