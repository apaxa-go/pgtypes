package pgtypes

import (
	"github.com/jackc/pgx"
	"strings"
	"testing"
)

type nullUUIDTestElement struct {
	sql string
	i   NullUUID
	err bool
}

var nullUUIDTests = []nullUUIDTestElement{
	{"SELECT '6ba7b814-9dad-11d1-80b4-00c04fd430c8'::UUID", NullUUID{UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, true}, false},
	{"SELECT 'string'::TEXT", NullUUID{}, true},
	{"SELECT null::UUID", NullUUID{UUID{}, false}, false},
}

func testNullUUID_Scan(t *testing.T) {
	for _, v := range nullUUIDTests {
		if rows, err := conn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r NullUUID
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v.sql)
				}
				if err := rows.Scan(&r); (err != nil) != v.err || r != v.i {
					t.Errorf("%v: expect %v %v, got %v %v", v.sql, v.i, v.err, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.sql)
				}
			}()
		}
	}
}

func TestNullUUID_Scan(t *testing.T) {
	testNullUUID_Scan(t)

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

	testNullUUID_Scan(t)

	pgx.DefaultTypeFormats["uuid"] = save

	// Reconnect with old FormatCode
	if err = conn.Close(); err != nil {
		panic(err)
	}
	if conn, err = pgx.Connect(conf); err != nil {
		panic(err)
	}
}

func TestNullUUID_Encode(t *testing.T) {
	for _, v := range nullUUIDTests {
		if v.err {
			continue
		}

		if rows, err := conn.Query("SELECT $1::UUID", v.i); err != nil {
			t.Errorf("%v: bad query", v.i)
		} else {
			func() {
				var r NullUUID
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v.i)
					return
				}
				if err := rows.Scan(&r); r != v.i || err != nil {
					t.Errorf("expect %v %v, got %v %v", v.i, nil, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.i)
				}
			}()
		}
	}
}

func TestNullUUID_Encode2(t *testing.T) {
	rightPrefix := "NullUUID.Encode cannot encode into OID "
	if rows, err := conn.Query("SELECT $1::INTEGER", NullUUID{UUID{}, true}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}
