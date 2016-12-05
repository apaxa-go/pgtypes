package pgtypes

import (
	"github.com/jackc/pgx"
	"strings"
	"testing"
)

func testUUID_Scan(t *testing.T) {
	type testElement struct {
		sql string
		u   UUID
		err bool
	}
	tests := []testElement{
		{"SELECT '6ba7b814-9dad-11d1-80b4-00c04fd430c8'::UUID", UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, false},
		{"SELECT 'string'::TEXT", UUID{}, true},
		{"SELECT null::uuid", UUID{}, true},
	}

	for _, v := range tests {
		if rows, err := conn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r UUID
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v.sql)
				}
				if err := rows.Scan(&r); (err != nil) != v.err || r != v.u {
					t.Errorf("%v: expect %v %v, got %v %v", v.sql, v.u, v.err, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.sql)
				}
			}()
		}
	}
}

func TestUUID_Scan(t *testing.T) {
	testUUID_Scan(t)

	save := pgx.DefaultTypeFormats["uuid"]
	switch save {
	case pgx.TextFormatCode:
		pgx.DefaultTypeFormats["uuid"] = pgx.BinaryFormatCode
	case pgx.BinaryFormatCode:
		pgx.DefaultTypeFormats["uuid"] = pgx.TextFormatCode
	}

	// Reconnect with new FormatCode
	var err error
	if err = conn.Close(); err != nil {
		panic(err)
	}
	if conn, err = pgx.Connect(conf); err != nil {
		panic(err)
	}

	testUUID_Scan(t)

	pgx.DefaultTypeFormats["uuid"] = save

	// Reconnect with old FormatCode
	if err = conn.Close(); err != nil {
		panic(err)
	}
	if conn, err = pgx.Connect(conf); err != nil {
		panic(err)
	}
}

func TestUUID_Encode(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	if rows, err := conn.Query("SELECT $1::UUID", u); err != nil {
		t.Error("bad query")
	} else {
		func() {
			var r UUID
			defer rows.Close()
			if !rows.Next() {
				t.Error("no row")
			}
			if err := rows.Scan(&r); r != u || err != nil {
				t.Errorf("expect %v %v, got %v %v", u, nil, r, err)
			}
			if rows.Next() {
				t.Error("multiple row")
			}
		}()
	}
}

func TestUUID_Encode2(t *testing.T) {
	rightPrefix := "UUID.Encode cannot encode into OID "
	if rows, err := conn.Query("SELECT $1::INTEGER", UUID{}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}
