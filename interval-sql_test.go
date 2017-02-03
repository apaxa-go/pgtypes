package pgtypes

import "testing"

func TestInterval_Scan(t *testing.T) {
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
		if rows, err := pqConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r Interval
				defer func() { _ = rows.Close() }()
				if !rows.Next() {
					t.Errorf("%v: no row", v.sql)
				}
				if err := rows.Scan(&r); (err != nil) != v.err || (!v.err && r != v.i) {
					t.Errorf("%v: expect %#v %v, got %#v %v", v.sql, v.i, v.err, r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.sql)
				}
			}()
		}
	}
}

func TestInterval_Value(t *testing.T) {
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
		if rows, err := pqConn.Query("SELECT $1::INTERVAL", v); err != nil {
			t.Errorf("%v: bad query", v)
		} else {
			func() {
				var r Interval
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
