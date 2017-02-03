package pgtypes

import (
	"fmt"
	"github.com/jackc/pgx"
)

func init() {
	// Register UUID type in pgx as binary-compatible.
	// This may cause error if type other when apaxa-go's UUID will be used with pgx for uuid storing.
	pgx.DefaultTypeFormats["uuid"] = pgx.BinaryFormatCode
}

// ScanPgx implements the pgx.PgxScanner interface.
func (u *UUID) ScanPgx(vr *pgx.ValueReader) error {
	if vr.Type().DataType != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("UUID.ScanPgx cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		return pgx.SerializationError("UUID.ScanPgx cannot parse NULL value")
	}

	switch vr.Type().FormatCode {
	case pgx.TextFormatCode:
		if vr.Len() != UUIDStringLen {
			return pgx.SerializationError(fmt.Sprintf("Received UUID string with invalid length: %d", vr.Len())) // It is hard cover this case with test
		}
		var err error
		if *u, err = ParseUUID(vr.ReadString(vr.Len())); err != nil {
			return pgx.SerializationError(fmt.Sprintf("Received invalid UUID string: %v", err.Error())) // It is hard cover this case with test
		}
	case pgx.BinaryFormatCode:
		if vr.Len() != UUIDLen {
			return pgx.SerializationError(fmt.Sprintf("Received UUID with invalid length: %d", vr.Len())) // It is hard cover this case with test
		}
		var err error
		if *u, err = ParseUUIDBytes(vr.ReadBytes(UUIDLen)); err != nil {
			return pgx.SerializationError(fmt.Sprintf("Received invalid UUID: %v", err.Error())) // It is hard cover this case with test
		}
	default:
		return fmt.Errorf("unknown format %v", vr.Type().FormatCode) // It is hard cover this case with test
	}

	return vr.Err()
}

// FormatCode implements the pgx.Encoder interface.
func (u UUID) FormatCode() int16 { return pgx.BinaryFormatCode }

// Encode implements the pgx.Encoder interface.
func (u UUID) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("UUID.Encode cannot encode into OID %d", oid))
	}

	w.WriteInt32(UUIDLen)
	w.WriteBytes(u[:])

	return nil
}
