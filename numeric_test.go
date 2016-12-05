package pgtypes

import (
	"github.com/apaxa-go/helper/strconvh"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestNumeric_FromToString(t *testing.T) {
	test := []string{
		"NaN",
		"0.12425345132423143452",
		"90.12425345132423143452",
		"890.12425345132423143452",
		"7890.12425345132423143452",
		"67890.12425345132423143452",
		"567890.12425345132423143452",
		"4567890.12425345132423143452",
		"34567890.12425345132423143452",
		"234567890.12425345132423143452",
		"1234567890.12425345132423143452",
		"1234567890.1242534513242314345",
		"1234567890.124253451324231434",
		"1234567890.12425345132423143",
		"1234567890.1242534513242314",
		"1234567890.124253451324231",
		"1234567890.12425345132423",
		"1234567890.1242534513242",
		"1234567890.124253451324",
		"1234567890.12425345132",
		"1234567890.1242534513",
		"1234567890.124253451",
		"1234567890.12425345",
		"1234567890.1242534",
		"1234567890.124253",
		"1234567890.12425",
		"1234567890.1242",
		"1234567890.124",
		"1234567890.12",
		"1234567890.1",
		"1234567890",
		"0.1",
		"0.12",
		"90.12",
		"90.124",
		"890.124",
		"890.1242",
		"7890.1242",
		"7890.12425",
		"67890.12425",
		"-123.456",
		"+123.456",
		"10000000000",
		"1000000000",
		"100000000",
		"10000000",
		"1000000",
		"100000",
		"10000",
		"1000",
		"100",
		"10",
		"1",
		"0.000000000001",
		"0.00000000001",
		"0.0000000001",
		"0.000000001",
		"0.00000001",
		"0.0000001",
		"0.000001",
		"0.00001",
		"0.0001",
		"0.001",
		"0.01",
		"0.1",
		"478997845379834578934789978543897534897978324897543789547856896548905649828940569346523457987.7734578789365789657894657895643789564786547865478657865798454697878956234789",
		"0",
		"-1",
		"-10",
		"-0.1",
		"-0.01",
		"123.",
	}
	for i, v := range test {
		var n Numeric
		if _, ok := n.SetString(v); !ok {
			t.Errorf("\nTestFromToString - %v. Unexpected error while parsing Numeric %v", i, v)
		}
		{ // Validate internal representation
			if n.sign == numericNaN { // NaN
				if n.weight != 0 || len(n.digits) != 0 {
					t.Errorf("\nTestFromToString - %v. %v - NaN should have weight=0 and no digits. Got: %#v", i, v, n)
				}
			} else if len(n.digits) == 0 { // Zero
				if n.sign != numericPositive || n.weight != 0 {
					t.Errorf("\nTestFromToString - %v. %v - zero value should have weight=0 and positive sign. Got: %#v", i, v, n)
				}
			} else { // Normal value
				if n.digits[0] == 0 || n.digits[len(n.digits)-1] == 0 {
					t.Errorf("\nTestFromToString - %v. %v - normal value should have non zero first and last elements in digits. Got: %#v", i, v, n)
				}
			}
		}

		if str := n.String(); strings.TrimSuffix(strings.TrimPrefix(v, "+"), ".") != str {
			t.Errorf("\nTestFromToString - %v. Expected: %v got: %v\nInternal representation: %#v", i, v, str, n)
		}
	}
}

func TestNumeric_SetString(t *testing.T) {
	badTests := []string{
		"",
		"a",
		"1a",
		"1a2",
		"a2",
		"a.",
		"1a.",
		"1a2.",
		"a2.",
		".a",
		".1a",
		".1a2",
		".a2",
		"1.2.3",
	}
	var n Numeric
	for _, v := range badTests {
		if _, ok := n.SetString(v); ok {
			t.Errorf("%v: expect not ok", v)
		}
	}
}

var numericTests = []int64{
	0,
	1,
	2,
	3,
	998,
	999,
	1000,
	1001,
	9999,
	12345678,
	123456789,
	math.MaxInt32 - 1,
	math.MaxInt32,
	-1,
	-2,
	-3,
	-998,
	-999,
	-1000,
	-1001,
	-9999,
	-12345678,
	-123456789,
	-math.MaxInt32 - 1,
	-math.MaxInt32,
}

func TestNumeric_Add(t *testing.T) {
	for _, v1 := range numericTests {
		for _, v2 := range numericTests {
			var a, b, c Numeric
			if _, ok := a.SetString(strconvh.FormatInt64(v1)); !ok {
				t.Errorf("TestAdd. Unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvh.FormatInt64(v2)); !ok {
				t.Errorf("TestAdd. Unable to parse int64 %v", v2)
			}
			c.Add(&a, &b)

			if s1, s2 := c.String(), strconvh.FormatInt64(v1+v2); s1 != s2 {
				t.Errorf("TestAdd. %v + %v expected %v, got %v", v1, v2, s2, s1)
			}
		}
	}
}

func TestNumeric_Cmp(t *testing.T) {
	for _, v1 := range numericTests {
		for _, v2 := range numericTests {
			var a, b Numeric
			if _, ok := a.SetString(strconvh.FormatInt64(v1)); !ok {
				t.Errorf("TestCmp. Unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvh.FormatInt64(v2)); !ok {
				t.Errorf("TestCmp. Unable to parse int64 %v", v2)
			}
			c := a.Cmp(&b)
			var r int
			switch {
			case v1 < v2:
				r = -1
			case v1 == v2:
				r = 0
			case v1 > v2:
				r = 1
			}
			if c != r {
				t.Errorf("TestCmp. %v ? %v expected %v, got %v", v1, v2, r, c)
			}
		}
	}
}

func TestNumeric_Cmp2(t *testing.T) {
	type testElement struct {
		s1, s2 string
		r      int
	}
	tests := []testElement{
		{"123.456", "0.0000789", 1},
		{"0.0000789", "123.456", -1},
		{"0.0000789", "0.0000789", 0},
		{"0.000078912345678", "0.0000789", 1},
		{"0.00007891", "0.000078912345678", -1},
		{"0.000078912345678", "0.00007891", 1},
		{"0.0000789", "0.000078912345678", -1},
		{"123.456", "123.457", -1},
		{"123.457", "123.456", 1},
		{"123.456", "1.2345678", 1},
		{"1.2345678", "123.456", -1},
		{"NaN", "1.2345678", 1},
		{"1.2345678", "NaN", -1},
		{"NaN", "NaN", 0},
	}
	for _, v := range tests {
		var n1, n2 Numeric
		if _, ok := n1.SetString(v.s1); !ok {
			t.Errorf("bad number '%v'", v.s1)
		}
		if _, ok := n2.SetString(v.s2); !ok {
			t.Errorf("bad number '%v'", v.s2)
		}
		if r := n1.Cmp(&n2); r != v.r {
			t.Errorf("%v,%v: expect %v, got %v", v.s1, v.s2, v.r, r)
		}
	}
}

func TestNumeric_Sub(t *testing.T) {
	for _, v1 := range numericTests {
		for _, v2 := range numericTests {
			var a, b, c Numeric
			if _, ok := a.SetString(strconvh.FormatInt64(v1)); !ok {
				t.Errorf("unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvh.FormatInt64(v2)); !ok {
				t.Errorf("unable to parse int64 %v", v2)
			}
			c.Sub(&a, &b)

			if s1, s2 := c.String(), strconvh.FormatInt64(v1-v2); s1 != s2 {
				t.Errorf("%v,%v: expected %v, got %v", v1, v2, s2, s1)
			}
		}
	}
}

func TestNumeric_Mul(t *testing.T) {
	for _, v1 := range numericTests {
		for _, v2 := range numericTests {
			var a, b, c Numeric
			if _, ok := a.SetString(strconvh.FormatInt64(v1)); !ok {
				t.Errorf("TestMul. Unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvh.FormatInt64(v2)); !ok {
				t.Errorf("TestMul. Unable to parse int64 %v", v2)
			}
			c.Mul(&a, &b)

			if s1, s2 := c.String(), strconvh.FormatInt64(v1*v2); s1 != s2 {
				t.Errorf("TestMul. %v * %v expected %v, got %v", v1, v2, s2, s1)
			}
		}
	}
}

func TestNumeric_SetIsNaN(t *testing.T) {
	var n Numeric
	if r := n.SetNaN(); r != &n || len(r.digits) != 0 || r.sign != numericNaN || r.weight != 0 || !r.IsNaN() {
		t.Error("error with NaN")
	}
}

// Just for coverage. No other way with current implementation to cover this.
func TestAddAbs(t *testing.T) {
	if d, w := addAbs(nil, 0, nil, 0); d != nil || w != 0 {
		t.Errorf("expect %v %v, got %v %v", nil, 0, d, w)
	}
}

func TestNumeric_Add2(t *testing.T) {
	var n1, n2, r1, r2, r3 Numeric
	n1.SetNaN()
	n2.SetZero()
	r1.Add(&n1, &n1)
	r2.Add(&n1, &n2)
	r3.Add(&n2, &n1)
	if !r1.IsNaN() || !r2.IsNaN() || !r3.IsNaN() {
		t.Error("error with NaN add")
	}
}

func TestNumeric_Neg(t *testing.T) {
	var n1, n2, n3, n4 Numeric
	var r1, r2, r3, r4 Numeric
	n1.SetZero()
	n2.SetNaN()
	if _, ok := n3.SetString("1"); !ok {
		t.Error("bad number")
	}
	if _, ok := n4.SetString("-1"); !ok {
		t.Error("bad number")
	}
	r1.Neg(&n1)
	r2.Neg(&n2)
	r3.Neg(&n3)
	r4.Neg(&n4)
	if !reflect.DeepEqual(r1, n1) || !reflect.DeepEqual(r2, n2) || !reflect.DeepEqual(r3, n4) || !reflect.DeepEqual(r4, n3) {
		t.Errorf("expect %v %v %v %v, got %v %v %v %v", &n1, &n2, &n4, &n3, &r1, &r2, &r3, &r4)
	}
}

func TestNumeric_Sub2(t *testing.T) {
	var n1, n2, r1, r2, r3 Numeric
	n1.SetNaN()
	n2.SetZero()
	r1.Sub(&n1, &n1)
	r2.Sub(&n1, &n2)
	r3.Sub(&n2, &n1)
	if !r1.IsNaN() || !r2.IsNaN() || !r3.IsNaN() {
		t.Error("error with NaN sub")
	}
}

// Just for coverage. No other way with current implementation to cover this.
func TestSubAbsOrdered(t *testing.T) {
	if d, w := subAbsOrdered(nil, 0, nil, 0); d != nil || w != 0 {
		t.Errorf("expect %v %v, got %v %v", nil, 0, d, w)
	}
}

// Just for coverage. No other way with current implementation to cover this.
func TestMulAbs(t *testing.T) {
	if d, w := mulAbs(nil, 0, nil, 0); d != nil || w != 0 {
		t.Errorf("expect %v %v, got %v %v", nil, 0, d, w)
	}
}

func TestNumeric_Mul2(t *testing.T) {
	var n1, n2, r1, r2, r3 Numeric
	n1.SetNaN()
	n2.SetZero()
	r1.Mul(&n1, &n1)
	r2.Mul(&n1, &n2)
	r3.Mul(&n2, &n1)
	if !r1.IsNaN() || !r2.IsNaN() || !r3.IsNaN() {
		t.Error("error with NaN mul")
	}
}

func TestNumeric_Sign(t *testing.T) {
	var n1, n2, n3, n4 Numeric
	n1.SetZero()
	n2.SetNaN()
	if _, ok := n3.SetString("1"); !ok {
		t.Error("bad number")
	}
	if _, ok := n4.SetString("-1"); !ok {
		t.Error("bad number")
	}
	if n1.Sign() != 0 || n2.Sign() != 2 || n3.Sign() != 1 || n4.Sign() != -1 {
		t.Error("error with Sign")
	}
}
