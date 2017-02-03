package pgtypes

import (
	"fmt"
	"github.com/jackc/pgx"
)

const (
	// IntervalOid is an object identifier (OID) of Interval type used in PostgreSQL.
	IntervalOid = 1186
	intervalLen = 8 + 4 + 4
)

func init() {
	// Register interval type in pgx as binary-compatible.
	// This may cause error if type other when apaxa-go's Interval will be used with pgx for interval storing.
	pgx.DefaultTypeFormats["interval"] = pgx.BinaryFormatCode
}

// ScanPgx implements the pgx.PgxScanner interface.
func (i *Interval) ScanPgx(vr *pgx.ValueReader) error {
	if vr.Type().DataType != IntervalOid {
		return pgx.SerializationError(fmt.Sprintf("Interval.ScanPgx cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		return pgx.SerializationError("Interval.ScanPgx cannot parse NULL value")
	}

	switch vr.Type().FormatCode {
	case pgx.TextFormatCode:
		var err error
		if *i, err = ParseInterval(vr.ReadString(vr.Len()), IntervalPgPrecision); err != nil {
			return pgx.SerializationError(fmt.Sprintf("received invalid Interval string: %v", err.Error())) // It is hard cover this case with test
		}
	case pgx.BinaryFormatCode:
		if vr.Len() != intervalLen {
			return pgx.SerializationError(fmt.Sprintf("received Interval with invalid length: %d", vr.Len())) // It is hard cover this case with test
		}

		i.precision = IntervalPgPrecision
		i.SomeSeconds = vr.ReadInt64()
		i.Days = vr.ReadInt32()
		i.Months = vr.ReadInt32()
	default:
		return fmt.Errorf("unknown format %v", vr.Type().FormatCode) // It is hard cover this case with test
	}

	return vr.Err()
}

// FormatCode implements the pgx.Encoder interface.
func (i Interval) FormatCode() int16 { return pgx.BinaryFormatCode }

// Encode implements the pgx.Encoder interface.
func (i Interval) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != IntervalOid {
		return pgx.SerializationError(fmt.Sprintf("Interval.Encode cannot encode into OID %d", oid))
	}

	w.WriteInt32(intervalLen)
	w.WriteInt64(someSecondsChangePrecision(i.SomeSeconds, i.precision, IntervalPgPrecision))
	w.WriteInt32(i.Days)
	w.WriteInt32(i.Months)

	return nil
}
