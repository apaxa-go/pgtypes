package pgtypes

import (
	"github.com/jackc/pgx"
	"strings"
	"testing"
)

func testNullUUID_ScanPgx(t *testing.T) {
	type testElement struct {
		sql string
		i   NullUUID
		err bool
	}

	tests := []testElement{
		{"SELECT '6ba7b814-9dad-11d1-80b4-00c04fd430c8'::UUID", NullUUID{UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, true}, false},
		{"SELECT 'string'::TEXT", NullUUID{}, true},
		{"SELECT null::UUID", NullUUID{UUID{}, false}, false},
	}

	for _, v := range tests {
		if rows, err := pgxConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query: %v", v.sql, err)
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

func TestNullUUID_ScanPgx(t *testing.T) {
	testNullUUID_ScanPgx(t)

	save := pgx.DefaultTypeFormats["uuid"]
	switch save {
	case pgx.TextFormatCode:
		pgx.DefaultTypeFormats["uuid"] = pgx.BinaryFormatCode
	case pgx.BinaryFormatCode:
		pgx.DefaultTypeFormats["uuid"] = pgx.TextFormatCode
	}

	// Reconnect with new FormatCode
	var err error
	if err = pgxConn.Close(); err != nil {
		panic(err)
	}
	if pgxConn, err = pgx.Connect(pgxConf); err != nil {
		panic(err)
	}

	testNullUUID_ScanPgx(t)

	pgx.DefaultTypeFormats["uuid"] = save

	// Reconnect with old FormatCode
	if err = pgxConn.Close(); err != nil {
		panic(err)
	}
	if pgxConn, err = pgx.Connect(pgxConf); err != nil {
		panic(err)
	}
}

func TestNullUUID_Encode(t *testing.T) {
	tests := []NullUUID{
		NullUUID{UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, true},
		NullUUID{UUID{}, false},
	}

	for _, v := range tests {
		if rows, err := pgxConn.Query("SELECT $1::UUID", v); err != nil {
			t.Errorf("%v: bad query: %v", v, err)
		} else {
			func() {
				var r NullUUID
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v)
					return
				}
				if err := rows.Scan(&r); r != v || err != nil {
					t.Errorf("expect %v %v, got %v %v", v, nil, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v)
				}
			}()
		}
	}
}

func TestNullUUID_Encode2(t *testing.T) {
	rightPrefix := "NullUUID.Encode cannot encode into OID "
	if rows, err := pgxConn.Query("SELECT $1::INTEGER", NullUUID{UUID{}, true}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}

func TestUUID_Nullable(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	nu := u.Nullable()
	if !nu.Valid || nu.UUID != u {
		t.Errorf("expect %v %v, got %v %v", true, u, nu.Valid, nu.UUID)
	}
}

func TestNullUUID_Scan(t *testing.T) {
	type testElement struct {
		sql string
		i   NullUUID
		err bool
	}

	tests := []testElement{
		{"SELECT '6ba7b814-9dad-11d1-80b4-00c04fd430c8'::UUID", NullUUID{UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, true}, false},
		{"SELECT 'string'::TEXT", NullUUID{}, true},
		{"SELECT null::UUID", NullUUID{UUID{}, false}, false},
	}

	for _, v := range tests {
		if rows, err := pqConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query: %v", v.sql, err)
		} else {
			func() {
				var r NullUUID
				defer func() { _ = rows.Close() }()
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

func TestNullUUID_Value(t *testing.T) {
	tests := []NullUUID{
		NullUUID{UUID{0x6b, 0xa7, 0xb8, 0x14 /**/, 0x9d, 0xad /**/, 0x11, 0xd1 /**/, 0x80, 0xb4 /**/, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, true},
		NullUUID{UUID{}, false},
	}

	for _, v := range tests {
		if rows, err := pqConn.Query("SELECT $1::UUID", v); err != nil {
			t.Errorf("%v: bad query: %v", v, err)
		} else {
			func() {
				var r NullUUID
				defer func() { _ = rows.Close() }()
				if !rows.Next() {
					t.Errorf("%v: no row", v)
					return
				}
				if err := rows.Scan(&r); r != v || err != nil {
					t.Errorf("expect %v %v, got %v %v", v, nil, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v)
				}
			}()
		}
	}
}
