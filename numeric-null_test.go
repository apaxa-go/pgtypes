package pgtypes

import (
	"github.com/jackc/pgx"
	"reflect"
	"strings"
	"testing"
)

type nullNumericTestElement struct {
	sql string
	n   NullNumeric
	err bool
}

var nullNumericTests = []nullNumericTestElement{
	{"SELECT '0'::Numeric", NullNumeric{*((&Numeric{}).SetZero()), true}, false},
	{"SELECT 'string'::TEXT", NullNumeric{}, true},
	{"SELECT null::Numeric", NullNumeric{Numeric{}, false}, false},
}

func testNullNumeric_Scan(t *testing.T) {
	for _, v := range nullNumericTests {
		if rows, err := conn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
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

func TestNullNumeric_Scan(t *testing.T) {
	testNullNumeric_Scan(t)

	save := pgx.DefaultTypeFormats["numeric"]
	switch save {
	case pgx.TextFormatCode:
		pgx.DefaultTypeFormats["numeric"] = pgx.BinaryFormatCode
	case pgx.BinaryFormatCode:
		pgx.DefaultTypeFormats["numeric"] = pgx.TextFormatCode
	}

	// Reconnect with new FormatCode
	var err error
	if err = conn.Close(); err != nil {
		panic(err)
	}
	if conn, err = pgx.Connect(conf); err != nil {
		panic(err)
	}

	testNullNumeric_Scan(t)

	pgx.DefaultTypeFormats["numeric"] = save

	// Reconnect with old FormatCode
	if err = conn.Close(); err != nil {
		panic(err)
	}
	if conn, err = pgx.Connect(conf); err != nil {
		panic(err)
	}
}

func TestNullNumeric_Encode(t *testing.T) {
	for _, v := range nullNumericTests {
		if v.err {
			continue
		}

		if rows, err := conn.Query("SELECT $1::Numeric", v.n); err != nil {
			t.Errorf("%v: bad query", v.n)
		} else {
			func() {
				var r NullNumeric
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v.n)
				}
				if err := rows.Scan(&r); !reflect.DeepEqual(r, v.n) || err != nil {
					t.Errorf("expect %v %v, got %v %v", v.n, nil, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.n)
				}
			}()
		}
	}
}

func TestNullNumeric_Encode2(t *testing.T) {
	rightPrefix := "NullNumeric.Encode cannot encode into OID "
	if rows, err := conn.Query("SELECT $1::INTEGER", NullNumeric{Numeric{}, true}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}
