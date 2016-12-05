package pgtypes

import (
	"encoding/hex"
	"errors"
)

const (
	// UUIDLen is a length of UUID in bytes
	UUIDLen = 16
	// UUIDCleanStringLen is a length of raw UUID string representation length in chars ("6ba7b8149dad11d180b400c04fd430c8")
	UUIDCleanStringLen = UUIDLen * 2
	// UUIDStringLen is a length of default UUID string representation length in chars ("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
	UUIDStringLen = UUIDLen*2 + 4
)

// String representation details.
const (
	// UUID parts lengths in bytes
	part0Len = 4
	part1Len = 2
	part2Len = 2
	part3Len = 2
	part4Len = 6

	uuidDelim = '-'

	// UUID delimiters positions in standard string representation
	sDelim0At = part0Len * 2
	sDelim1At = sDelim0At + 1 + part1Len*2
	sDelim2At = sDelim1At + 1 + part2Len*2
	sDelim3At = sDelim2At + 1 + part3Len*2

	// UUID parts position in standard string representation
	sPart0From = 0
	sPart0To   = sPart0From + part0Len*2
	sPart1From = sPart0To + 1
	sPart1To   = sPart1From + part1Len*2
	sPart2From = sPart1To + 1
	sPart2To   = sPart2From + part2Len*2
	sPart3From = sPart2To + 1
	sPart3To   = sPart3From + part3Len*2
	sPart4From = sPart3To + 1
	sPart4To   = sPart4From + part4Len*2

	// UUID parts position in internal representation
	// No need to store partXTo as partXTo = part(X+1)From
	part0From = 0
	part1From = part0From + part0Len
	part2From = part1From + part1Len
	part3From = part2From + part2Len
	part4From = part3From + part3Len
)

// UUID representation compliant with specification described in RFC 4122.
type UUID [UUIDLen]byte

var zeroUUID = UUID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// ParseUUID parses an UUID standard string representation (like "6ba7b814-9dad-11d1-80b4-00c04fd430c8").
func ParseUUID(s string) (u UUID, err error) {
	if len(s) != UUIDStringLen {
		err = errors.New("invalid UUID string length")
		return
	}

	if s[sDelim0At] != uuidDelim || s[sDelim1At] != uuidDelim || s[sDelim2At] != uuidDelim || s[sDelim3At] != uuidDelim {
		err = errors.New("invalid UUID string delimiters")
		return
	}

	b := []byte(s)

	if l, e := hex.Decode(u[part0From:part1From], b[sPart0From:sPart0To]); l != part0Len || e != nil {
		err = errors.New("invalid UUID part 1")
		return
	}

	if l, e := hex.Decode(u[part1From:part2From], b[sPart1From:sPart1To]); l != part1Len || e != nil {
		err = errors.New("invalid UUID part 2")
		return
	}

	if l, e := hex.Decode(u[part2From:part3From], b[sPart2From:sPart2To]); l != part2Len || e != nil {
		err = errors.New("invalid UUID part 3")
		return
	}

	if l, e := hex.Decode(u[part3From:part4From], b[sPart3From:sPart3To]); l != part3Len || e != nil {
		err = errors.New("invalid UUID part 4")
		return
	}

	if l, e := hex.Decode(u[part4From:], b[sPart4From:sPart4To]); l != part4Len || e != nil {
		err = errors.New("invalid UUID part 5")
		return
	}
	return
}

// ParseUUIDClean parses an UUID clean string representation (like "6ba7b8149dad11d180b400c04fd430c8").
func ParseUUIDClean(s string) (u UUID, err error) {
	if len(s) != UUIDCleanStringLen {
		err = errors.New("Invalid UUID clean string length")
		return
	}

	if l, e := hex.Decode(u[:], []byte(s)); l != UUIDLen || e != nil {
		err = errors.New("Unable to parse UUID from clean string")
	}

	return
}

// ParseUUIDBytes parses byte slice (as-is) with UUID.
func ParseUUIDBytes(b []byte) (u UUID, err error) {
	if len(b) != UUIDLen {
		err = errors.New("Given slice is not valid UUID sequence")
		return
	}
	copy(u[:], b)
	return
}

// IsZero return true if UUID is zero (all bytes are zero).
func (u UUID) IsZero() bool {
	return u == zeroUUID
}

// String returns standard string representation of UUID.
func (u UUID) String() string {
	buf := make([]byte, UUIDStringLen)

	hex.Encode(buf[sPart0From:sPart0To], u[part0From:part1From])
	buf[sDelim0At] = uuidDelim
	hex.Encode(buf[sPart1From:sPart1To], u[part1From:part2From])
	buf[sDelim1At] = uuidDelim
	hex.Encode(buf[sPart2From:sPart2To], u[part2From:part3From])
	buf[sDelim2At] = uuidDelim
	hex.Encode(buf[sPart3From:sPart3To], u[part3From:part4From])
	buf[sDelim3At] = uuidDelim
	hex.Encode(buf[sPart4From:sPart4To], u[part4From:])

	return string(buf)
}

// CleanString returns clean string representation of UUID (without any additional chars and in lower case).
func (u UUID) CleanString() string {
	return hex.EncodeToString(u[:])
}

// Bytes returns byte slice containing UUID (as-is).
func (u UUID) Bytes() []byte {
	return u[:]
}
