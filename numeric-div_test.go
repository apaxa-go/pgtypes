package pgtypes

import (
	"github.com/apaxa-go/helper/strconvh"
	"reflect"
	"testing"
)

func TestNumeric_QuoRem(t *testing.T) {
	type testElement struct {
		a, b     string
		quo, rem string
	}
	tests := []testElement{
		{"10", "3", "3", "1"},
		{"-10", "3", "-3", "-1"},
		{"10", "-3", "-3", "1"},
		{"-10", "-3", "3", "-1"},
		{"NaN", "3", "NaN", "NaN"},
		{"10", "NaN", "NaN", "NaN"},
		{"NaN", "NaN", "NaN", "NaN"},
		{"12345", "1", "12345", "0"},
		{"12345", "2", "6172", "1"},
		{"123", "12345", "0", "123"},
		{"123002460012300", "12300", "10000200001", "0"},
		{"123002460012300", "10000200001", "12300", "0"},
		{"123456789", "5351", "23071", "3868"},
	}
	for _, v := range tests {
		var a, b, quo, rem Numeric
		if _, ok := a.SetString(v.a); !ok {
			t.Errorf("%v: bad Numeric", v.a)
		}
		if _, ok := b.SetString(v.b); !ok {
			t.Errorf("%v: bad Numeric", v.b)
		}
		if _, ok := quo.SetString(v.quo); !ok {
			t.Errorf("%v: bad Numeric", v.quo)
		}
		if _, ok := rem.SetString(v.rem); !ok {
			t.Errorf("%v: bad Numeric", v.rem)
		}
		var r1, r2, r3 Numeric
		r1.QuoRem(&a, &b, &r2)
		r3.Rem(&a, &b)
		if !reflect.DeepEqual(r1, quo) || !reflect.DeepEqual(r2, rem) || !reflect.DeepEqual(r3, rem) {
			t.Errorf("%v,%v: expect %v %v %v, got %v %v %v", &a, &b, &quo, &rem, &rem, &r1, &r2, &r3)
		}
	}
}

func TestNumeric_QuoRem2(t *testing.T) {
	var a, b Numeric
	if _, ok := a.SetString("1"); !ok {
		t.Errorf("%v: bad Numeric", "1")
	}
	b.SetZero()
	defer func() {
		if recover() == nil {
			t.Error("panic expected")
		}
	}()
	a.QuoRem(&a, &b, &b)
}

func TestNumeric_Quo(t *testing.T) {
	delta := float64(1e-7)
	var deltaN Numeric
	deltaN.setString(strconvh.FormatFloat64(delta))

	for _, v1 := range numericTests {
		for _, v2 := range numericTests {
			if v2 == 0 {
				continue
			}
			var a, b, c Numeric
			if _, ok := a.SetString(strconvh.FormatInt64(v1)); !ok {
				t.Errorf("unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvh.FormatInt64(v2)); !ok {
				t.Errorf("unable to parse int64 %v", v2)
			}
			c.Quo(&a, &b)

			r := float64(v1) / float64(v2)
			if r == 0 { // Avoid float -0
				r = 0
			}

			var rN, rD Numeric
			rN.setString(strconvh.FormatFloat64(r))
			rD.Sub(&c, &rN).Abs(&rD)

			if s1, s2 := c.String(), rN.String(); rD.Cmp(&deltaN) > 0 {
				t.Errorf("%v,%v: expected %v, got %v", v1, v2, s2, s1)
			}
		}
	}
}

func TestNumeric_Quo2(t *testing.T) {
	var a, b, quo, r Numeric
	if _, ok := a.SetString("0.00000000000000000000000000000000000000002"); !ok {
		t.Error("bad Numeric")
	}
	if _, ok := b.SetString("0.00000000000000000000000000000000000000001"); !ok {
		t.Error("bad Numeric")
	}
	if _, ok := quo.SetString("2"); !ok {
		t.Error("bad Numeric")
	}
	r.Quo(&a, &b)
	if !reflect.DeepEqual(r, quo) {
		t.Errorf("expect %v, got %v", &quo, &r)
	}
}

func TestNumeric_QuoPrec(t *testing.T) {
	type testElement struct {
		a, b  string
		p     int16
		round bool
		quo   string
	}
	tests := []testElement{
		{"0", "3", 0, true, "0"},
		{"0", "3", 0, false, "0"},
		{"10", "3", 0, true, "3"},
		{"10", "3", 0, false, "3"},
		{"10", "3", 1, true, "3.3"},
		{"10", "3", 1, false, "3.3"},
		{"12345", "2", 0, true, "6173"},
		{"12345", "2", 0, false, "6172"},
		{"12345", "2", 1, true, "6172.5"},
		{"12345", "2", 1, false, "6172.5"},
		{"12345", "2", 2, true, "6172.5"},
		{"12345", "2", 2, false, "6172.5"},
		{"12345", "2", 10, true, "6172.5"},
		{"12345", "2", 10, false, "6172.5"},
		{"1.2345", "2", 0, true, "1"},
		{"1.2345", "2", 0, false, "0"},
		{"1.2345", "2", 1, true, "0.6"},
		{"1.2345", "2", 1, false, "0.6"},
		{"1.2345", "2", 2, true, "0.62"},
		{"1.2345", "2", 2, false, "0.61"},
		{"1.2345", "2", 3, true, "0.617"},
		{"1.2345", "2", 3, false, "0.617"},
		{"1.2345", "2", 4, true, "0.6173"},
		{"1.2345", "2", 4, false, "0.6172"},
		{"1.2345", "2", 5, true, "0.61725"},
		{"1.2345", "2", 5, false, "0.61725"},
		{"1.2345", "2", 6, true, "0.61725"},
		{"1.2345", "2", 6, false, "0.61725"},
		{"1.2345", "2", 10, true, "0.61725"},
		{"1.2345", "2", 10, false, "0.61725"},
	}
	for _, v := range tests {
		var a, b, quo Numeric
		if _, ok := a.SetString(v.a); !ok {
			t.Errorf("%v: bad Numeric", v.a)
		}
		if _, ok := b.SetString(v.b); !ok {
			t.Errorf("%v: bad Numeric", v.b)
		}
		if _, ok := quo.SetString(v.quo); !ok {
			t.Errorf("%v: bad Numeric", v.quo)
		}
		var r Numeric
		r.QuoPrec(&a, &b, v.p, v.round)
		if !reflect.DeepEqual(r, quo) {
			t.Errorf("%v,%v,%v,%v: expect %#v, got %#v", &a, &b, v.p, v.round, quo, r)
		}
	}
}

// Just for coverage. No other way with current implementation to cover this.
func TestTruncAbs(t *testing.T) {
	if d, w := truncAbs([]int16{}, 1, 10); d != nil || w != 0 {
		t.Errorf("expect %v %v, got %v %v", nil, 0, d, w)
	}
}

// Just for coverage. No other way with current implementation to cover this.
func TestRoundAbs(t *testing.T) {
	if d, w := roundAbs([]int16{}, 1, 10); d != nil || w != 0 {
		t.Errorf("expect %v %v, got %v %v", nil, 0, d, w)
	}
}
