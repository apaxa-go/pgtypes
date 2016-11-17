package interval

import (
	"fmt"
	"github.com/jackc/pgx"
)

const (
	intervalOid = 1186
	intervalLen = 8 + 4 + 4
)

func init() {
	// Register interval type in pgx as binary-compatible.
	// This may cause error if type other when apaxa-io's Interval will be used with pgx for interval storing.
	pgx.DefaultTypeFormats["interval"] = pgx.BinaryFormatCode
}

// Scan implements the pgx.Scanner interface.
func (u *Interval) Scan(vr *pgx.ValueReader) error {
	if vr.Type().DataType != intervalOid {
		return pgx.SerializationError(fmt.Sprintf("Interval.Scan cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		return pgx.SerializationError("Interval.Scan cannot parse NULL value")
	}

	switch vr.Type().FormatCode {
	case pgx.TextFormatCode:
		var err error
		if *u, err = Parse(vr.ReadString(vr.Len()), PostgreSQLPrecision); err != nil {
			return pgx.SerializationError(fmt.Sprintf("Received invalid Interval string: %v", err.Error()))
		}
	case pgx.BinaryFormatCode:
		if vr.Len() != intervalLen {
			return pgx.SerializationError(fmt.Sprintf("Received Interval with invalid length: %d", vr.Len()))
		}

		u.precision = PostgreSQLPrecision
		u.SomeSeconds = vr.ReadInt64()
		u.Days = vr.ReadInt32()
		u.Months = vr.ReadInt32()
	default:
		return fmt.Errorf("unknown format %v", vr.Type().FormatCode)
	}

	return vr.Err()
}

// FormatCode implements the pgx.Encoder interface.
func (u Interval) FormatCode() int16 {
	return pgx.BinaryFormatCode
}

// Encode implements the pgx.Encoder interface.
func (u Interval) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != intervalOid {
		return pgx.SerializationError(fmt.Sprintf("Interval.Encode cannot encode into OID %d", oid))
	}

	w.WriteInt32(intervalLen)
	w.WriteInt64(someSecondsChangePrecision(u.SomeSeconds, u.precision, PostgreSQLPrecision))
	w.WriteInt32(u.Days)
	w.WriteInt32(u.Months)

	return nil
}
