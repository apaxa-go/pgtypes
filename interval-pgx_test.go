package pgtypes

import (
	"github.com/jackc/pgx"
	"strings"
	"testing"
)

type intervalTestElement struct {
	sql string
	i   Interval
	err bool
}

var intervalTests = []intervalTestElement{
	{"SELECT '0'::INTERVAL", Interval{0, 0, 0, PgPrecision}, false},
	{"SELECT '1 seconds'::INTERVAL", Interval{0, 0, 1e6, PgPrecision}, false},
	{"SELECT '2 seconds'::INTERVAL", Interval{0, 0, 2e6, PgPrecision}, false},
	{"SELECT '1 days'::INTERVAL", Interval{0, 1, 0, PgPrecision}, false},
	{"SELECT '2 days'::INTERVAL", Interval{0, 2, 0, PgPrecision}, false},
	{"SELECT '1 mons'::INTERVAL", Interval{1, 0, 0, PgPrecision}, false},
	{"SELECT '2 mons'::INTERVAL", Interval{2, 0, 0, PgPrecision}, false},
	{"SELECT '1 years'::INTERVAL", Interval{12, 0, 0, PgPrecision}, false},
	{"SELECT '2 years'::INTERVAL", Interval{24, 0, 0, PgPrecision}, false},
	{"SELECT '12'::INTERVAL", Interval{0, 0, 12e6, PgPrecision}, false},
	{"SELECT '13:12'::INTERVAL", Interval{0, 0, (13*3600 + 12*60) * 1e6, PgPrecision}, false},
	{"SELECT '14:13:12'::INTERVAL", Interval{0, 0, (14*3600 + 13*60 + 12) * 1e6, PgPrecision}, false},
	{"SELECT '2 day 1 seconds'::INTERVAL", Interval{0, 2, 1e6, PgPrecision}, false},
	{"SELECT '-2 day 1 seconds'::INTERVAL", Interval{0, -2, 1e6, PgPrecision}, false},
	{"SELECT '2 day -1 seconds'::INTERVAL", Interval{0, 2, -1e6, PgPrecision}, false},
	{"SELECT '3 year 2 day -1 seconds'::INTERVAL", Interval{3 * 12, 2, -1e6, PgPrecision}, false},
	{"SELECT '-3 year 2 day -1 seconds'::INTERVAL", Interval{-3 * 12, 2, -1e6, PgPrecision}, false},
	{"SELECT '-3 year 2 day -1.23456 seconds'::INTERVAL", Interval{-3 * 12, 2, -1234560, PgPrecision}, false},
	{"SELECT '-3 year 2 day -1.234567 seconds'::INTERVAL", Interval{-3 * 12, 2, -1234567, PgPrecision}, false},
	{"SELECT 'string'::TEXT", Interval{}, true},
	{"SELECT null::interval", Interval{}, true},
}

func testInterval_Scan(t *testing.T) {
	for _, v := range intervalTests {
		if rows, err := conn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r Interval
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

func TestInterval_Scan(t *testing.T) {
	testInterval_Scan(t)

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

	testInterval_Scan(t)

	pgx.DefaultTypeFormats["interval"] = save

	// Reconnect with old FormatCode
	if err = conn.Close(); err != nil {
		panic(err)
	}
	if conn, err = pgx.Connect(conf); err != nil {
		panic(err)
	}
}

func TestInterval_Encode(t *testing.T) {
	for _, v := range intervalTests {
		if v.err {
			continue
		}

		if rows, err := conn.Query("SELECT $1::INTERVAL", v.i); err != nil {
			t.Errorf("%v: bad query", v.i)
		} else {
			func() {
				var r Interval
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

func TestInterval_Encode2(t *testing.T) {
	rightPrefix := "Interval.Encode cannot encode into OID "
	if rows, err := conn.Query("SELECT $1::INTEGER", Interval{1, 2, 3, PgPrecision}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}
