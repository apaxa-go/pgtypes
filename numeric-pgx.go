package pgtypes

import (
	"fmt"
	"github.com/apaxa-go/helper/mathh"
	"github.com/jackc/pgx"
)

const (
	// NumericOid is an object identifier (OID) of Numeric type used in PostgreSQL.
	NumericOid       = 1700
	numericHeaderLen = 4 * 2
)

func init() {
	// Register Numeric type in pgx as binary-compatible.
	// This may cause error if type other when apaxa-go's Numeric will be used with pgx for uuid storing.
	pgx.DefaultTypeFormats["numeric"] = pgx.BinaryFormatCode
}

// Scan implements the pgx.Scanner interface.
func (n *Numeric) Scan(vr *pgx.ValueReader) error {
	if vr.Type().DataType != NumericOid {
		return pgx.SerializationError(fmt.Sprintf("Numeric.Scan cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		return pgx.SerializationError("Numeric.Scan cannot parse NULL value")
	}

	switch vr.Type().FormatCode {
	case pgx.TextFormatCode:
		if str := vr.ReadString(vr.Len()); !n.setString(str) {
			return pgx.SerializationError(fmt.Sprintf("received invalid Numeric string: '%v'", str)) // It is hard cover this case with test
		}
	case pgx.BinaryFormatCode:
		if vr.Len() < numericHeaderLen {
			return pgx.SerializationError(fmt.Sprintf("Received Numeric with invalid length: %d", vr.Len())) // It is hard cover this case with test
		}

		l := vr.ReadInt16()

		if l < 0 || vr.Len()+2 != numericHeaderLen+int32(l)*2 {
			return pgx.SerializationError(fmt.Sprintf("Received inconsistent Numeric: length %d, number of digits = %d", vr.Len(), l)) // It is hard cover this case with test
		}

		n.weight = vr.ReadInt16()
		n.sign = numericSign(vr.ReadInt16())

		if n.sign == numericNaN {
			if l > 0 {
				return pgx.SerializationError(fmt.Sprintf("Received inconsistent Numeric: NaN with number of digits = %d", n.sign)) // It is hard cover this case with test
			}
		} else if n.sign != numericPositive && n.sign != numericNegative {
			return pgx.SerializationError(fmt.Sprintf("Received Numeric with invalid sign: %d", n.sign)) // It is hard cover this case with test
		}

		vr.ReadInt16() // Here we read dscale. Currently we just ignore it.

		if l == 0 {
			n.weight = 0 // PostgreSQL can return not very expected combination (9.4 can return NaN with Weight=99). Here it will be normalized.
			n.digits = nil
		} else {
			n.digits = make([]int16, l)
			for i := 0; i < int(l); i++ {
				n.digits[i] = vr.ReadInt16()
				// We believe PostgreSQL :=) ... and do not like overhead.
				//if n.digits[i]<0 || n.digits[i]>=numericBase {
				//	return pgx.SerializationError(fmt.Sprintf("Received Numeric with invalid digit: %d", n.digits[i]))	// It is hard cover this case with test
				//}
			}
		}
	default:
		return fmt.Errorf("unknown format %v", vr.Type().FormatCode) // It is hard cover this case with test
	}

	return vr.Err()
}

// FormatCode implements the pgx.Encoder interface.
func (n Numeric) FormatCode() int16 { return pgx.BinaryFormatCode }

// Encode implements the pgx.Encoder interface.
func (n Numeric) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != NumericOid {
		return pgx.SerializationError(fmt.Sprintf("Numeric.Encode cannot encode into OID %d", oid))
	}

	l := len(n.digits)
	if l > mathh.MaxInt16 {
		return pgx.SerializationError(fmt.Sprintf("Numeric.Encode cannot encode so much digits: %d", l))
	}

	w.WriteInt32(numericHeaderLen + 2*int32(l))
	w.WriteInt16(int16(l))
	w.WriteInt16(n.weight)
	w.WriteInt16(int16(n.sign))
	w.WriteInt16(getScaleAbs(n.digits, n.weight))
	for _, v := range n.digits {
		w.WriteInt16(v)
	}

	return nil
}
