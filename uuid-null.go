package pgtypes

import (
	"fmt"
	"github.com/jackc/pgx"
)

// NullUUID represents an UUID that may be NULL.
// NullUUID implements the pgx.Scanner and pgx.Encoder interfaces so it may be used both as an argument to Query[Row] and a destination for Scan.
//
// If Valid is false then the value is NULL.
type NullUUID struct {
	UUID  UUID
	Valid bool // Valid is true if UUID is not NULL
}

// Scan implements the pgx.Scanner interface.
func (u *NullUUID) Scan(vr *pgx.ValueReader) error {
	if vr.Type().DataType != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("NullUUID.Scan cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		u.UUID, u.Valid = zeroUUID, false
		return nil
	}

	u.Valid = true
	return u.UUID.Scan(vr)
}

// FormatCode implements the pgx.Encoder interface.
func (u NullUUID) FormatCode() int16 { return pgx.BinaryFormatCode }

// Encode implements the pgx.Encoder interface.
func (u NullUUID) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("NullUUID.Encode cannot encode into OID %d", oid))
	}

	if !u.Valid {
		w.WriteInt32(-1)
		return nil
	}

	return u.UUID.Encode(w, oid)
}

// Nullable returns valid NullUUID with UUID u.
func (u UUID) Nullable() NullUUID {
	return NullUUID{u, true}
}
