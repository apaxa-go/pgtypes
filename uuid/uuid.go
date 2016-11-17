package uuid

import (
	"encoding/hex"
	"errors"
)

// UUID representation compliant with specification described in RFC 4122.
type UUID [uuidLen]byte

var nullUUID = UUID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// ParseString extract UUID from its standard string representation (like "6ba7b814-9dad-11d1-80b4-00c04fd430c8").
func (u *UUID) ParseString(s string) (err error) {
	if len(s) != uuidStringLen {
		err = errors.New("Invalid UUID string length")
		return
	}

	// TODO after testing remove check for "l"
	if s[sDelim0At] != uuidDelim || s[sDelim1At] != uuidDelim || s[sDelim2At] != uuidDelim || s[sDelim3At] != uuidDelim {
		err = errors.New("Invalid UUID string delimiters")
		return
	}

	b := []byte(s)

	if l, e := hex.Decode(u[part0From:part1From], b[sPart0From:sPart0To]); l != part0Len || e != nil {
		err = errors.New("Invalid UUID part 1")
		return
	}

	if l, e := hex.Decode(u[part1From:part2From], b[sPart1From:sPart1To]); l != part1Len || e != nil {
		err = errors.New("Invalid UUID part 2")
		return
	}

	if l, e := hex.Decode(u[part2From:part3From], b[sPart2From:sPart2To]); l != part2Len || e != nil {
		err = errors.New("Invalid UUID part 3")
		return
	}

	if l, e := hex.Decode(u[part3From:part4From], b[sPart3From:sPart3To]); l != part3Len || e != nil {
		err = errors.New("Invalid UUID part 4")
		return
	}

	if l, e := hex.Decode(u[part4From:], b[sPart4From:sPart4To]); l != part4Len || e != nil {
		err = errors.New("Invalid UUID part 5")
		return
	}
	return
}

// ParseCleanString extract UUID from clean string UUID representation (like "6ba7b8149dad11d180b400c04fd430c8").
func (u *UUID) ParseCleanString(s string) (err error) {
	if len(s) != uuidCleanStringLen {
		err = errors.New("Invalid UUID clean string length")
		return
	}

	if l, e := hex.Decode(u[:], []byte(s)); l != uuidLen || e != nil {
		err = errors.New("Unable to parse UUID from clean string")
	}

	return
}

// ParseBytes extract UUID from byte slice (as-is).
func (u *UUID) ParseBytes(b []byte) (err error) {
	if len(b) != uuidLen {
		err = errors.New("Given slice is not valid UUID sequence")
		return
	}
	copy(u[:], b)
	return
}

// FromString creates a UUID object from given hex string
// representation. Function accepts UUID string in following
// formats:
//
//     uuid.FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//
func FromString(s string) (u UUID, err error) {
	err = u.ParseString(s)
	return
}

// FromCleanString creates a UUID object from given hex string
// in lower case without any additional chars.
//
//     uuid.FromCleanString("6ba7b8149dad11d180b400c04fd430c8")
//
func FromCleanString(s string) (u UUID, err error) {
	err = u.ParseCleanString(s)
	return
}

// IsNull return true if UUID is NULL (all bytes are zero).
func (u UUID) IsNull() bool {
	return u == nullUUID
}

// FromBytes creates a UUID object from given bytes slice.
func FromBytes(b []byte) (u UUID, err error) {
	err = u.ParseBytes(b)
	return
}

// String returns string representation of UUID.
func (u UUID) String() string {
	buf := make([]byte, uuidStringLen)

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

// Null returns NULL UUID (all bytes are zero).
func Null() UUID {
	return nullUUID
}
