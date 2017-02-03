package pgtypes

import (
	"github.com/jackc/pgx"
	"strings"
	"testing"
)

func testNullInterval_ScanPgx(t *testing.T) {
	type testElement struct {
		sql string
		i   NullInterval
		err bool
	}

	var tests = []testElement{
		{"SELECT '3 years 2 days 12:34'::INTERVAL", NullInterval{Interval{3 * 12, 2, (12*3600 + 34*60) * 1e6, IntervalPgPrecision}, true}, false},
		{"SELECT 'string'::TEXT", NullInterval{}, true},
		{"SELECT null::interval", NullInterval{Interval{0, 0, 0, IntervalPgPrecision}, false}, false},
	}

	for _, v := range tests {
		if rows, err := pgxConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
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

func TestNullInterval_ScanPgx(t *testing.T) {
	testNullInterval_ScanPgx(t)

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

	testNullInterval_ScanPgx(t)

	pgx.DefaultTypeFormats["interval"] = save

	// Reconnect with old FormatCode
	if err = pgxConn.Close(); err != nil {
		panic(err)
	}
	if pgxConn, err = pgx.Connect(pgxConf); err != nil {
		panic(err)
	}
}

func TestNullInterval_Encode(t *testing.T) {
	var tests = []NullInterval{
		NullInterval{Interval{3 * 12, 2, (12*3600 + 34*60) * 1e6, IntervalPgPrecision}, true},
		NullInterval{Interval{0, 0, 0, IntervalPgPrecision}, false},
	}

	for _, v := range tests {
		if rows, err := pgxConn.Query("SELECT $1::INTERVAL", v); err != nil {
			t.Errorf("%v: bad query", v)
		} else {
			func() {
				var r NullInterval
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

func TestNullInterval_Encode2(t *testing.T) {
	rightPrefix := "NullInterval.Encode cannot encode into OID "
	if rows, err := pgxConn.Query("SELECT $1::INTEGER", NullInterval{Interval{1, 2, 3, IntervalPgPrecision}, true}); err == nil || !strings.HasPrefix(err.Error(), rightPrefix) {
		t.Errorf("expect '%v', got %v", rightPrefix, err)
		rows.Close()
	}
}

func TestInterval_Nullable(t *testing.T) {
	i := Interval{1, 2, 3, IntervalGoPrecision}
	ni := i.Nullable()
	if !ni.Valid || ni.Interval != i {
		t.Errorf("expect %v %v, got %v %v", true, i, ni.Valid, ni.Interval)
	}
}

func TestNullInterval_Scan(t *testing.T) {
	type testElement struct {
		sql string
		i   NullInterval
		err bool
	}

	var tests = []testElement{
		{"SELECT '3 years 2 days 12:34'::INTERVAL", NullInterval{Interval{3 * 12, 2, (12*3600 + 34*60) * 1e6, IntervalPgPrecision}, true}, false},
		{"SELECT 'string'::TEXT", NullInterval{}, true},
		{"SELECT null::interval", NullInterval{Interval{0, 0, 0, IntervalPgPrecision}, false}, false},
	}

	for _, v := range tests {
		if rows, err := pqConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r NullInterval
				defer func() { _ = rows.Close() }()
				if !rows.Next() {
					t.Errorf("%v: no row", v.sql)
				}
				if err := rows.Scan(&r); (err != nil) != v.err || (!v.err && r != v.i) {
					t.Errorf("%v: expect %v %v, got %v %v", v.sql, v.i, v.err, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.sql)
				}
			}()
		}
	}
}

func TestNullInterval_Value(t *testing.T) {
	var tests = []NullInterval{
		NullInterval{Interval{3 * 12, 2, (12*3600 + 34*60) * 1e6, IntervalPgPrecision}, true},
		NullInterval{Interval{0, 0, 0, IntervalPgPrecision}, false},
	}

	for _, v := range tests {
		if rows, err := pqConn.Query("SELECT $1::INTERVAL", v); err != nil {
			t.Errorf("%v: bad query", v)
		} else {
			func() {
				var r NullInterval
				defer func() { _ = rows.Close() }()
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
