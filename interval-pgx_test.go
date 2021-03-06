package pgtypes

import (
	"github.com/jackc/pgx"
	"strings"
	"testing"
)

func testInterval_ScanPgx(t *testing.T) {
	type testElement struct {
		sql string
		i   Interval
		err bool
	}

	var tests = []testElement{
		{"SELECT '0'::INTERVAL", Interval{0, 0, 0, IntervalPgPrecision}, false},
		{"SELECT '1 seconds'::INTERVAL", Interval{0, 0, 1e6, IntervalPgPrecision}, false},
		{"SELECT '2 seconds'::INTERVAL", Interval{0, 0, 2e6, IntervalPgPrecision}, false},
		{"SELECT '1 days'::INTERVAL", Interval{0, 1, 0, IntervalPgPrecision}, false},
		{"SELECT '2 days'::INTERVAL", Interval{0, 2, 0, IntervalPgPrecision}, false},
		{"SELECT '1 mons'::INTERVAL", Interval{1, 0, 0, IntervalPgPrecision}, false},
		{"SELECT '2 mons'::INTERVAL", Interval{2, 0, 0, IntervalPgPrecision}, false},
		{"SELECT '1 years'::INTERVAL", Interval{12, 0, 0, IntervalPgPrecision}, false},
		{"SELECT '2 years'::INTERVAL", Interval{24, 0, 0, IntervalPgPrecision}, false},
		{"SELECT '12'::INTERVAL", Interval{0, 0, 12e6, IntervalPgPrecision}, false},
		{"SELECT '13:12'::INTERVAL", Interval{0, 0, (13*3600 + 12*60) * 1e6, IntervalPgPrecision}, false},
		{"SELECT '14:13:12'::INTERVAL", Interval{0, 0, (14*3600 + 13*60 + 12) * 1e6, IntervalPgPrecision}, false},
		{"SELECT '2 day 1 seconds'::INTERVAL", Interval{0, 2, 1e6, IntervalPgPrecision}, false},
		{"SELECT '-2 day 1 seconds'::INTERVAL", Interval{0, -2, 1e6, IntervalPgPrecision}, false},
		{"SELECT '2 day -1 seconds'::INTERVAL", Interval{0, 2, -1e6, IntervalPgPrecision}, false},
		{"SELECT '3 year 2 day -1 seconds'::INTERVAL", Interval{3 * 12, 2, -1e6, IntervalPgPrecision}, false},
		{"SELECT '-3 year 2 day -1 seconds'::INTERVAL", Interval{-3 * 12, 2, -1e6, IntervalPgPrecision}, false},
		{"SELECT '-3 year 2 day -1.23456 seconds'::INTERVAL", Interval{-3 * 12, 2, -1234560, IntervalPgPrecision}, false},
		{"SELECT '-3 year 2 day -1.234567 seconds'::INTERVAL", Interval{-3 * 12, 2, -1234567, IntervalPgPrecision}, false},
		{"SELECT 'string'::TEXT", Interval{}, true},
		{"SELECT null::interval", Interval{}, true},
	}

	for _, v := range tests {
		if rows, err := pgxConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
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

func TestInterval_ScanPgx(t *testing.T) {
	testInterval_ScanPgx(t)

	save := pgx.DefaultTypeFormats["interval"]
	switch save {
	case pgx.TextFormatCode:
		pgx.DefaultTypeFormats["interval"] = pgx.BinaryFormatCode
	case pgx.BinaryFormatCode:
		pgx.DefaultTypeFormats["interval"] = pgx.TextFormatCode
	}

	// Reconnect with new FormatCode
	var err error
	if err = pgxConn.Close(); err != nil {
		panic(err)
	}
	if pgxConn, err = pgx.Connect(pgxConf); err != nil {
		panic(err)
	}

	testInterval_ScanPgx(t)

	pgx.DefaultTypeFormats["interval"] = save

	// Reconnect with old FormatCode
	if err = pgxConn.Close(); err != nil {
		panic(err)
	}
	if pgxConn, err = pgx.Connect(pgxConf); err != nil {
		panic(err)
	}
}

func TestInterval_Encode(t *testing.T) {
	var tests = []Interval{
		Interval{0, 0, 0, IntervalPgPrecision},
		Interval{0, 0, 1e6, IntervalPgPrecision},
		Interval{0, 0, 2e6, IntervalPgPrecision},
		Interval{0, 1, 0, IntervalPgPrecision},
		Interval{0, 2, 0, IntervalPgPrecision},
		Interval{1, 0, 0, IntervalPgPrecision},
		Interval{2, 0, 0, IntervalPgPrecision},
		Interval{12, 0, 0, IntervalPgPrecision},
		Interval{24, 0, 0, IntervalPgPrecision},
		Interval{0, 0, 12e6, IntervalPgPrecision},
		Interval{0, 0, (13*3600 + 12*60) * 1e6, IntervalPgPrecision},
		Interval{0, 0, (14*3600 + 13*60 + 12) * 1e6, IntervalPgPrecision},
		Interval{0, 2, 1e6, IntervalPgPrecision},
		Interval{0, -2, 1e6, IntervalPgPrecision},
		Interval{0, 2, -1e6, IntervalPgPrecision},
		Interval{3 * 12, 2, -1e6, IntervalPgPrecision},
		Interval{-3 * 12, 2, -1e6, IntervalPgPrecision},
		Interval{-3 * 12, 2, -1234560, IntervalPgPrecision},
		Interval{-3 * 12, 2, -1234567, IntervalPgPrecision},
	}

	for _, v := range tests {
		if rows, err := pgxConn.Query("SELECT $1::INTERVAL", v); err != nil {
			t.Errorf("%v: bad query", v)
		} else {
			func() {
				var r Interval
				defer rows.Close()
				if !rows.Next() {
					t.Errorf("%v: no row", v)
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

func TestInterval_Encode2(t *testing.T) {
	rightPrefix := "Interval.Encode cannot encode into OID "
	if rows, err := pgxConn.Query("SELECT $1::INTEGER", Interval{1, 2, 3, IntervalPgPrecision}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}
