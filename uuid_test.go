package pgtypes

import (
	"reflect"
	"testing"
)

func TestParseUUID(t *testing.T) {
	type testElement struct {
		s   string
		u   UUID
		err bool
	}
	tests := []testElement{
		{"6ba7b814-9dad-11d1-80b4-00c04fd430c8", UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, false},
		{"6BA7B814-9dad-11D1-80b4-00C04FD430C8", UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, false},
		{"00000000-0000-0000-0000-000000000000", UUID{}, false},
		{"0000000-0000-0000-0000-000000000000", UUID{}, true},
		{"00000000-000-0000-0000-000000000000", UUID{}, true},
		{"00000000-0000-000-0000-000000000000", UUID{}, true},
		{"00000000-0000-0000-000-000000000000", UUID{}, true},
		{"00000000-0000-0000-0000-00000000000", UUID{}, true},
		{"000000000000-0000-0000-000000000000", UUID{}, true},
		{"00000000-00000000-0000-000000000000", UUID{}, true},
		{"00000000-0000-00000000-000000000000", UUID{}, true},
		{"00000000-0000-0000-0000000000000000", UUID{}, true},
		{"0000x000-0000-0000-0000-000000000000", UUID{}, true},
		{"00000000-000x-0000-0000-000000000000", UUID{}, true},
		{"00000000-0000-00x0-0000-000000000000", UUID{}, true},
		{"00000000-0000-0000-x000-000000000000", UUID{}, true},
		{"00000000-0000-0000-0000-0000000x0000", UUID{}, true},
		{"00000000-0000-0000-0000-0000000000000", UUID{}, true},
		{"00000000+0000-0000-0000-000000000000", UUID{}, true},
		{"00000000-0000+0000-0000-000000000000", UUID{}, true},
		{"00000000-0000-0000+0000-000000000000", UUID{}, true},
		{"00000000-0000-0000-0000+000000000000", UUID{}, true},
	}
	for _, v := range tests {
		if u, err := ParseUUID(v.s); u != v.u || err != nil != v.err {
			t.Errorf("expect %v %v, got %v %v", v.u, v.err, u, err)
		}
	}
}

func TestParseUUIDClean(t *testing.T) {
	type testElement struct {
		s   string
		u   UUID
		err bool
	}
	tests := []testElement{
		{"6ba7b8149dad11d180b400c04fd430c8", UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, false},
		{"6BA7B8149dad11D180b400C04FD430C8", UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, false},
		{"00000000000000000000000000000000", UUID{}, false},
		{"0000000000000000000000000000000", UUID{}, true},
		{"0000x000000000000000000000000000", UUID{}, true},
		{"000000000000000000000000000000000", UUID{}, true},
	}
	for _, v := range tests {
		if u, err := ParseUUIDClean(v.s); u != v.u || err != nil != v.err {
			t.Errorf("expect %v %v, got %v %v", v.u, v.err, u, err)
		}
	}
}

func TestParseUUIDBytes(t *testing.T) {
	type testElement struct {
		b   []byte
		u   UUID
		err bool
	}
	tests := []testElement{
		{[]byte{0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, false},
		{make([]byte, 16), UUID{}, false},
		{make([]byte, 15), UUID{}, true},
		{make([]byte, 17), UUID{}, true},
	}
	for _, v := range tests {
		if u, err := ParseUUIDBytes(v.b); u != v.u || err != nil != v.err {
			t.Errorf("expect %v %v, got %v %v", v.u, v.err, u, err)
		}
	}
}

func TestUUID_IsZero(t *testing.T) {
	if !(UUID{}).IsZero() {
		t.Error("zero UUID expected")
	}
	if (UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}).IsZero() {
		t.Error("non ero UUID expected")
	}
}

func TestUUID_String(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	s := "6ba7b814-9dad-11d1-80b4-00c04fd430c8"
	if u.String() != s {
		t.Error("expect equality")
	}
}

func TestUUID_CleanString(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	s := "6ba7b8149dad11d180b400c04fd430c8"
	if u.CleanString() != s {
		t.Error("expect equality")
	}
}

func TestUUID_Bytes(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b := []byte{0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	if !reflect.DeepEqual(u.Bytes(), b) {
		t.Error("expect equality")
	}
}
