package pgtypes

import (
	"github.com/jackc/pgx"
	"strings"
	"testing"
)

type nullIntervalTestElement struct {
	sql string
	i   NullInterval
	err bool
}

var nullIntervalTests = []nullIntervalTestElement{
	{"SELECT '3 years 2 days 12:34'::INTERVAL", NullInterval{Interval{3 * 12, 2, (12*3600 + 34*60) * 1e6, PgPrecision}, true}, false},
	{"SELECT 'string'::TEXT", NullInterval{}, true},
	{"SELECT null::interval", NullInterval{Interval{0, 0, 0, PgPrecision}, false}, false},
}

func testNullInterval_Scan(t *testing.T) {
	for _, v := range nullIntervalTests {
		if rows, err := conn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r NullInterval
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

func TestNullInterval_Scan(t *testing.T) {
	testNullInterval_Scan(t)

	save := pgx.DefaultTypeFormats["interval"]
	switch save {
	case pgx.TextFormatCode:
		pgx.DefaultTypeFormats["interval"] = pgx.BinaryFormatCode
	case pgx.BinaryFormatCode:
		pgx.DefaultTypeFormats["interval"] = pgx.TextFormatCode
	}

	// Reconnect with new FormatCode
	var err error
	if err = conn.Close(); err != nil {
		panic(err)
	}
	if conn, err = pgx.Connect(conf); err != nil {
		panic(err)
	}

	testNullInterval_Scan(t)

	pgx.DefaultTypeFormats["interval"] = save

	// Reconnect with old FormatCode
	if err = conn.Close(); err != nil {
		panic(err)
	}
	if conn, err = pgx.Connect(conf); err != nil {
		panic(err)
	}
}

func TestNullInterval_Encode(t *testing.T) {
	for _, v := range nullIntervalTests {
		if v.err {
			continue
		}

		if rows, err := conn.Query("SELECT $1::INTERVAL", v.i); err != nil {
			t.Errorf("%v: bad query", v.i)
		} else {
			func() {
				var r NullInterval
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v.i)
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

func TestNullInterval_Encode2(t *testing.T) {
	rightPrefix := "NullInterval.Encode cannot encode into OID "
	if rows, err := conn.Query("SELECT $1::INTEGER", NullInterval{Interval{1, 2, 3, PgPrecision}, true}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}
