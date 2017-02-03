package pgtypes

import (
	"github.com/jackc/pgx"
	"reflect"
	"strings"
	"testing"
)

func testNullNumeric_ScanPgx(t *testing.T) {
	type testElement struct {
		sql string
		n   NullNumeric
		err bool
	}

	var tests = []testElement{
		{"SELECT '0'::Numeric", NullNumeric{*((&Numeric{}).SetZero()), true}, false},
		{"SELECT 'string'::TEXT", NullNumeric{}, true},
		{"SELECT null::Numeric", NullNumeric{Numeric{}, false}, false},
	}

	for _, v := range tests {
		if rows, err := pgxConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r NullNumeric
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v.sql)
					return
				}
				if err := rows.Scan(&r); (err != nil) != v.err || !reflect.DeepEqual(r, v.n) {
					t.Errorf("%v: expect %v %v, got %v %v", v.sql, v.n, v.err, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.sql)
				}
			}()
		}
	}
}

func TestNullNumeric_ScanPgx(t *testing.T) {
	testNullNumeric_ScanPgx(t)

	save := pgx.DefaultTypeFormats["numeric"]
	switch save {
	case pgx.TextFormatCode:
		pgx.DefaultTypeFormats["numeric"] = pgx.BinaryFormatCode
	case pgx.BinaryFormatCode:
		pgx.DefaultTypeFormats["numeric"] = pgx.TextFormatCode
	}

	// Reconnect with new FormatCode
	var err error
	if err = pgxConn.Close(); err != nil {
		panic(err)
	}
	if pgxConn, err = pgx.Connect(pgxConf); err != nil {
		panic(err)
	}

	testNullNumeric_ScanPgx(t)

	pgx.DefaultTypeFormats["numeric"] = save

	// Reconnect with old FormatCode
	if err = pgxConn.Close(); err != nil {
		panic(err)
	}
	if pgxConn, err = pgx.Connect(pgxConf); err != nil {
		panic(err)
	}
}

func TestNullNumeric_Encode(t *testing.T) {
	var tests = []NullNumeric{
		NullNumeric{*((&Numeric{}).SetZero()), true},
		NullNumeric{Numeric{}, false},
	}

	for _, v := range tests {
		if rows, err := pgxConn.Query("SELECT $1::Numeric", v); err != nil {
			t.Errorf("%v: bad query", v)
		} else {
			func() {
				var r NullNumeric
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v)
				}
				if err := rows.Scan(&r); !reflect.DeepEqual(r, v) || err != nil {
					t.Errorf("%v: expect %v %v, got %v %v", v, v, nil, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v)
				}
			}()
		}
	}
}

func TestNullNumeric_Encode2(t *testing.T) {
	rightPrefix := "NullNumeric.Encode cannot encode into OID "
	if rows, err := pgxConn.Query("SELECT $1::INTEGER", NullNumeric{Numeric{}, true}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}

func TestNumeric_Nullable(t *testing.T) {
	var n Numeric
	n.SetInt64(123)
	nn := n.Nullable()
	if !nn.Valid || !reflect.DeepEqual(nn.Numeric, n) {
		t.Errorf("expect %v %v, got %v %v", true, n, nn.Valid, nn.Numeric)
	}
}

func TestNullNumeric_Scan(t *testing.T) {
	type testElement struct {
		sql string
		n   NullNumeric
		err bool
	}

	var tests = []testElement{
		{"SELECT '0'::Numeric", NullNumeric{*((&Numeric{}).SetZero()), true}, false},
		{"SELECT 'string'::TEXT", NullNumeric{}, true},
		{"SELECT null::Numeric", NullNumeric{Numeric{}, false}, false},
	}

	for _, v := range tests {
		if rows, err := pqConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r NullNumeric
				defer func() { _ = rows.Close() }()
				if !rows.Next() {
					t.Errorf("%v: no row", v.sql)
					return
				}
				if err := rows.Scan(&r); (err != nil) != v.err || !reflect.DeepEqual(r, v.n) {
					t.Errorf("%v: expect %v %v, got %v %v", v.sql, v.n, v.err, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.sql)
				}
			}()
		}
	}
}

func TestNullNumeric_Value(t *testing.T) {
	var tests = []NullNumeric{
		NullNumeric{*((&Numeric{}).SetZero()), true},
		NullNumeric{Numeric{}, false},
	}

	for _, v := range tests {
		if rows, err := pqConn.Query("SELECT $1::Numeric", v); err != nil {
			t.Errorf("%v: bad query", v)
		} else {
			func() {
				var r NullNumeric
				defer func() { _ = rows.Close() }()
				if !rows.Next() {
					t.Errorf("%v: no row", v)
				}
				if err := rows.Scan(&r); !reflect.DeepEqual(r, v) || err != nil {
					t.Errorf("expect %v %v, got %v %v", v, nil, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v)
				}
			}()
		}
	}
}
