package uuid

import (
	"fmt"
	"github.com/apaxa-io/uuid"
	"github.com/jackc/pgx"
)

// NullUUID represent UUID that may be NULL.
// NullUUID implements the pgx.Scanner and pgx.Encoder interfaces.
type NullUUID struct {
	UUID  uuid.UUID
	Valid bool // Valid is true if UUID is not NULL
}

// Scan implements the pgx.Scanner interface.
func (u *NullUUID) Scan(vr *pgx.ValueReader) error {
	if vr.Type().DataType != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("UUID.Scan cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		u.UUID, u.Valid = uuid.Null(), false
		return nil
	}

	u.Valid = true
	return u.UUID.Scan(vr)
}

// FormatCode implements the pgx.Encoder interface.
func (u NullUUID) FormatCode() int16 {
	return pgx.BinaryFormatCode
}

// Encode implements the pgx.Encoder interface.
func (u NullUUID) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("UUID.Encode cannot encode into OID %d", oid))
	}

	if !u.Valid {
		w.WriteInt32(-1)
		return nil
	}

	return u.UUID.Encode(w, oid)
}

// TODO rename function. May be convert to method of UUID?
// FromUUID construct nullable UUID from given UUID.
func FromUUID(u uuid.UUID) NullUUID {
	return NullUUID{UUID: u, Valid: true}
}

// Null construct invalid (set to null) NullUUID.
func Null() NullUUID {
	return NullUUID{Valid: false}
}

