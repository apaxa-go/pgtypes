package pgtypes

import (
	"reflect"
	"testing"
)

func TestNumeric_Scan(t *testing.T) {
	type testElement struct {
		sql string
		n   *Numeric
		err bool
	}
	setString := func(s string) *Numeric {
		var n Numeric
		n.setString(s)
		return &n
	}
	tests := []testElement{
		{"SELECT 'NaN'::Numeric", (&Numeric{}).SetNaN(), false},
		{"SELECT '0'::Numeric", (&Numeric{}).SetZero(), false},
		{"SELECT '1'::Numeric", (&Numeric{}).SetInt64(1), false},
		{"SELECT '-1'::Numeric", (&Numeric{}).SetInt64(-1), false},
		{"SELECT '9999'::Numeric", (&Numeric{}).SetInt64(9999), false},
		{"SELECT '-9999'::Numeric", (&Numeric{}).SetInt64(-9999), false},
		{"SELECT '1239999'::Numeric", (&Numeric{}).SetInt64(1239999), false},
		{"SELECT '-1239999'::Numeric", (&Numeric{}).SetInt64(-1239999), false},
		{"SELECT '1239900'::Numeric", (&Numeric{}).SetInt64(1239900), false},
		{"SELECT '-1239900'::Numeric", (&Numeric{}).SetInt64(-1239900), false},
		{"SELECT '1230000'::Numeric", (&Numeric{}).SetInt64(1230000), false},
		{"SELECT '-1230000'::Numeric", (&Numeric{}).SetInt64(-1230000), false},
		{"SELECT '123.456'::Numeric", setString("123.456"), false},
		{"SELECT '-123.456'::Numeric", setString("-123.456"), false},
		{"SELECT '0.456'::Numeric", setString("0.456"), false},
		{"SELECT '-0.456'::Numeric", setString("-0.456"), false},
		{"SELECT '0.0000456'::Numeric", setString("0.0000456"), false},
		{"SELECT '-0.0000456'::Numeric", setString("-0.0000456"), false},
		{"SELECT 'string'::TEXT", &Numeric{}, true},
		{"SELECT null::Numeric", &Numeric{}, true},
	}

	for _, v := range tests {
		if rows, err := pqConn.Query(v.sql); err != nil { // Do not use QueryRow because it is harder to split error origin.
			t.Errorf("%v: bad query", v.sql)
		} else {
			func() {
				var r Numeric
				defer func() { _ = rows.Close() }()
				if !rows.Next() {
					t.Errorf("%v: no row", v.sql)
				}
				if err := rows.Scan(&r); (err != nil) != v.err || !reflect.DeepEqual(r, *v.n) {
					t.Errorf("%v: expect %v %v, got %#v %v", v.sql, v.n, v.err, &r, err)
				}
				if rows.Next() {
					t.Errorf("%v: multiple row", v.sql)
				}
			}()
		}
	}
}

func TestNumeric_Value(t *testing.T) {
	setString := func(s string) *Numeric {
		var n Numeric
		n.setString(s)
		return &n
	}
	tests := []*Numeric{
		(&Numeric{}).SetNaN(),
		(&Numeric{}).SetZero(),
		(&Numeric{}).SetInt64(1),
		(&Numeric{}).SetInt64(-1),
		(&Numeric{}).SetInt64(9999),
		(&Numeric{}).SetInt64(-9999),
		(&Numeric{}).SetInt64(1239999),
		(&Numeric{}).SetInt64(-1239999),
		(&Numeric{}).SetInt64(1239900),
		(&Numeric{}).SetInt64(-1239900),
		(&Numeric{}).SetInt64(1230000),
		(&Numeric{}).SetInt64(-1230000),
		setString("123.456"),
		setString("-123.456"),
		setString("0.456"),
		setString("-0.456"),
		setString("0.0000456"),
		setString("-0.0000456"),
	}
	for _, v := range tests {
		if rows, err := pqConn.Query("SELECT $1::Numeric", v); err != nil {
			t.Error("bad query")
		} else {
			func() {
				var r Numeric
				defer func() { _ = rows.Close() }()
				if !rows.Next() {
					t.Error("no row")
					return
				}
				if err := rows.Scan(&r); !reflect.DeepEqual(r, *v) || err != nil {
					t.Errorf("expect %v %v, got %v %v", v, nil, &r, err)
				}
				if rows.Next() {
					t.Error("multiple row")
				}
			}()
		}
	}
}
