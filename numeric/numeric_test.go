package numeric

import (
	"github.com/apaxa-io/strconvhelper"
	"math"
	"testing"
)

func TestFromToString(t *testing.T) {
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
	}
	for i, v := range test {
		var n Numeric
		if _, ok := n.SetString(v); !ok {
			t.Errorf("\nTestFromToString - %v. Unexpected error while parsing Numeric %v", i, v)
		}
		{ // Validate internal representation
			if n.sign == NaN { // NaN
				if n.weight != 0 || len(n.digits) != 0 {
					t.Errorf("\nTestFromToString - %v. %v - NaN should have weight=0 and no digits. Got: %#v", i, v, n)
				}
			} else if len(n.digits) == 0 { // Zero
				if n.sign != Positive || n.weight != 0 {
					t.Errorf("\nTestFromToString - %v. %v - zero value should have weight=0 and positive sign. Got: %#v", i, v, n)
				}
			} else { // Normal value
				if n.digits[0] == 0 || n.digits[len(n.digits)-1] == 0 {
					t.Errorf("\nTestFromToString - %v. %v - normal value should have non zero first and last elements in digits. Got: %#v", i, v, n)
				}
			}
		}

		if str := n.String(); str != v {
			t.Errorf("\nTestFromToString - %v. Expected: %v got: %v\nInternal representation: %#v", i, v, str, n)
		}
	}
}

var tests = []int64{
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

func TestAdd(t *testing.T) {
	for _, v1 := range tests {
		for _, v2 := range tests {
			var a, b, c Numeric
			if _, ok := a.SetString(strconvhelper.FormatInt64(v1)); !ok {
				t.Errorf("TestAdd. Unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvhelper.FormatInt64(v2)); !ok {
				t.Errorf("TestAdd. Unable to parse int64 %v", v2)
			}
			c.Add(&a, &b)

			if s1, s2 := c.String(), strconvhelper.FormatInt64(v1+v2); s1 != s2 {
				t.Errorf("TestAdd. %v + %v expected %v, got %v", v1, v2, s2, s1)
			}
		}
	}
}

func TestCmp(t *testing.T) {
	for _, v1 := range tests {
		for _, v2 := range tests {
			var a, b Numeric
			if _, ok := a.SetString(strconvhelper.FormatInt64(v1)); !ok {
				t.Errorf("TestCmp. Unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvhelper.FormatInt64(v2)); !ok {
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

func TestSub(t *testing.T) {
	for _, v1 := range tests {
		for _, v2 := range tests {
			var a, b, c Numeric
			if _, ok := a.SetString(strconvhelper.FormatInt64(v1)); !ok {
				t.Errorf("TestSub. Unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvhelper.FormatInt64(v2)); !ok {
				t.Errorf("TestSub. Unable to parse int64 %v", v2)
			}
			c.Sub(&a, &b)

			if s1, s2 := c.String(), strconvhelper.FormatInt64(v1-v2); s1 != s2 {
				t.Errorf("TestSub. %v - %v expected %v, got %v", v1, v2, s2, s1)
			}
		}
	}
}

func TestMul(t *testing.T) {
	for _, v1 := range tests {
		for _, v2 := range tests {
			var a, b, c Numeric
			if _, ok := a.SetString(strconvhelper.FormatInt64(v1)); !ok {
				t.Errorf("TestMul. Unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvhelper.FormatInt64(v2)); !ok {
				t.Errorf("TestMul. Unable to parse int64 %v", v2)
			}
			c.Mul(&a, &b)

			if s1, s2 := c.String(), strconvhelper.FormatInt64(v1*v2); s1 != s2 {
				t.Errorf("TestMul. %v * %v expected %v, got %v", v1, v2, s2, s1)
			}
		}
	}
}

func TestDiv(t *testing.T) {
	delta := float64(1e-7)
	var deltaN Numeric
	deltaN.setString(strconvhelper.FormatFloat64(delta))

	for _, v1 := range tests {
		for _, v2 := range tests {
			if v2 == 0 {
				continue
			}
			var a, b, c Numeric
			if _, ok := a.SetString(strconvhelper.FormatInt64(v1)); !ok {
				t.Errorf("TestDiv. Unable to parse int64 %v", v1)
			}
			if _, ok := b.SetString(strconvhelper.FormatInt64(v2)); !ok {
				t.Errorf("TestDiv. Unable to parse int64 %v", v2)
			}
			c.Div(&a, &b)

			r := float64(v1) / float64(v2)
			if r == 0 { // Avoid float -0
				r = 0
			}

			var rN, rD Numeric
			rN.setString(strconvhelper.FormatFloat64(r))
			rD.Sub(&c, &rN).Abs(&rD)

			if s1, s2 := c.String(), rN.String(); rD.Cmp(&deltaN) > 0 {
				t.Errorf("TestDiv. %v / %v expected %v, got %v", v1, v2, s2, s1)
			}
		}
	}
}
