//replacer:generated-file

package pgtypes

import (
	"github.com/apaxa-go/helper/mathh"
	"github.com/apaxa-go/helper/strconvh"
	"testing"
)

func TestNumeric_SetInt(t *testing.T) {
	tests := []int{
		mathh.MinInt,
		mathh.MinInt + 1,
		mathh.MinInt / 2,
		mathh.MinInt/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxInt,
		10000 % mathh.MaxInt,
		100000 % mathh.MaxInt,
		1000000 % mathh.MaxInt,
		10000000 % mathh.MaxInt,
		100000000 % mathh.MaxInt,
		1000000000 % mathh.MaxInt,
		10000000000 % mathh.MaxInt,
		100000000000 % mathh.MaxInt,
		mathh.MaxInt/2 - 1,
		mathh.MaxInt / 2,
		mathh.MaxInt/2 + 1,
		mathh.MaxInt - 1,
		mathh.MaxInt,
	}
	for _, v := range tests {
		var n Numeric
		n.SetInt(v)

		if str := strconvh.FormatInt(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Int(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Int(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Int(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Int to avoid "1000 overflow int8"
	n.SetInt(mathh.MaxInt).Mul(&tmp, &n).Int()
	n.SetInt(mathh.MinInt).Mul(&tmp, &n).Int()
}

func TestNumeric_SetInt8(t *testing.T) {
	tests := []int8{
		mathh.MinInt8,
		mathh.MinInt8 + 1,
		mathh.MinInt8 / 2,
		mathh.MinInt8/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxInt8,
		10000 % mathh.MaxInt8,
		100000 % mathh.MaxInt8,
		1000000 % mathh.MaxInt8,
		10000000 % mathh.MaxInt8,
		100000000 % mathh.MaxInt8,
		1000000000 % mathh.MaxInt8,
		10000000000 % mathh.MaxInt8,
		100000000000 % mathh.MaxInt8,
		mathh.MaxInt8/2 - 1,
		mathh.MaxInt8 / 2,
		mathh.MaxInt8/2 + 1,
		mathh.MaxInt8 - 1,
		mathh.MaxInt8,
	}
	for _, v := range tests {
		var n Numeric
		n.SetInt8(v)

		if str := strconvh.FormatInt8(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Int8(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Int8(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Int8(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Int8 to avoid "1000 overflow int8"
	n.SetInt8(mathh.MaxInt8).Mul(&tmp, &n).Int8()
	n.SetInt8(mathh.MinInt8).Mul(&tmp, &n).Int8()
}

func TestNumeric_SetInt16(t *testing.T) {
	tests := []int16{
		mathh.MinInt16,
		mathh.MinInt16 + 1,
		mathh.MinInt16 / 2,
		mathh.MinInt16/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxInt16,
		10000 % mathh.MaxInt16,
		100000 % mathh.MaxInt16,
		1000000 % mathh.MaxInt16,
		10000000 % mathh.MaxInt16,
		100000000 % mathh.MaxInt16,
		1000000000 % mathh.MaxInt16,
		10000000000 % mathh.MaxInt16,
		100000000000 % mathh.MaxInt16,
		mathh.MaxInt16/2 - 1,
		mathh.MaxInt16 / 2,
		mathh.MaxInt16/2 + 1,
		mathh.MaxInt16 - 1,
		mathh.MaxInt16,
	}
	for _, v := range tests {
		var n Numeric
		n.SetInt16(v)

		if str := strconvh.FormatInt16(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Int16(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Int16(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Int16(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Int16 to avoid "1000 overflow int8"
	n.SetInt16(mathh.MaxInt16).Mul(&tmp, &n).Int16()
	n.SetInt16(mathh.MinInt16).Mul(&tmp, &n).Int16()
}

func TestNumeric_SetInt32(t *testing.T) {
	tests := []int32{
		mathh.MinInt32,
		mathh.MinInt32 + 1,
		mathh.MinInt32 / 2,
		mathh.MinInt32/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxInt32,
		10000 % mathh.MaxInt32,
		100000 % mathh.MaxInt32,
		1000000 % mathh.MaxInt32,
		10000000 % mathh.MaxInt32,
		100000000 % mathh.MaxInt32,
		1000000000 % mathh.MaxInt32,
		10000000000 % mathh.MaxInt32,
		100000000000 % mathh.MaxInt32,
		mathh.MaxInt32/2 - 1,
		mathh.MaxInt32 / 2,
		mathh.MaxInt32/2 + 1,
		mathh.MaxInt32 - 1,
		mathh.MaxInt32,
	}
	for _, v := range tests {
		var n Numeric
		n.SetInt32(v)

		if str := strconvh.FormatInt32(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Int32(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Int32(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Int32(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Int32 to avoid "1000 overflow int8"
	n.SetInt32(mathh.MaxInt32).Mul(&tmp, &n).Int32()
	n.SetInt32(mathh.MinInt32).Mul(&tmp, &n).Int32()
}

func TestNumeric_SetUint(t *testing.T) {
	tests := []uint{
		mathh.MinUint,
		mathh.MinUint + 1,
		mathh.MinUint / 2,
		mathh.MinUint/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxUint,
		10000 % mathh.MaxUint,
		100000 % mathh.MaxUint,
		1000000 % mathh.MaxUint,
		10000000 % mathh.MaxUint,
		100000000 % mathh.MaxUint,
		1000000000 % mathh.MaxUint,
		10000000000 % mathh.MaxUint,
		100000000000 % mathh.MaxUint,
		mathh.MaxUint/2 - 1,
		mathh.MaxUint / 2,
		mathh.MaxUint/2 + 1,
		mathh.MaxUint - 1,
		mathh.MaxUint,
	}
	for _, v := range tests {
		var n Numeric
		n.SetUint(v)

		if str := strconvh.FormatUint(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Uint(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Uint(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Uint(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Uint to avoid "1000 overflow int8"
	n.SetUint(mathh.MaxUint).Mul(&tmp, &n).Uint()
	n.SetUint(mathh.MinUint).Mul(&tmp, &n).Uint()
}

func TestNumeric_SetUint8(t *testing.T) {
	tests := []uint8{
		mathh.MinUint8,
		mathh.MinUint8 + 1,
		mathh.MinUint8 / 2,
		mathh.MinUint8/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxUint8,
		10000 % mathh.MaxUint8,
		100000 % mathh.MaxUint8,
		1000000 % mathh.MaxUint8,
		10000000 % mathh.MaxUint8,
		100000000 % mathh.MaxUint8,
		1000000000 % mathh.MaxUint8,
		10000000000 % mathh.MaxUint8,
		100000000000 % mathh.MaxUint8,
		mathh.MaxUint8/2 - 1,
		mathh.MaxUint8 / 2,
		mathh.MaxUint8/2 + 1,
		mathh.MaxUint8 - 1,
		mathh.MaxUint8,
	}
	for _, v := range tests {
		var n Numeric
		n.SetUint8(v)

		if str := strconvh.FormatUint8(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Uint8(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Uint8(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Uint8(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Uint8 to avoid "1000 overflow int8"
	n.SetUint8(mathh.MaxUint8).Mul(&tmp, &n).Uint8()
	n.SetUint8(mathh.MinUint8).Mul(&tmp, &n).Uint8()
}

func TestNumeric_SetUint16(t *testing.T) {
	tests := []uint16{
		mathh.MinUint16,
		mathh.MinUint16 + 1,
		mathh.MinUint16 / 2,
		mathh.MinUint16/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxUint16,
		10000 % mathh.MaxUint16,
		100000 % mathh.MaxUint16,
		1000000 % mathh.MaxUint16,
		10000000 % mathh.MaxUint16,
		100000000 % mathh.MaxUint16,
		1000000000 % mathh.MaxUint16,
		10000000000 % mathh.MaxUint16,
		100000000000 % mathh.MaxUint16,
		mathh.MaxUint16/2 - 1,
		mathh.MaxUint16 / 2,
		mathh.MaxUint16/2 + 1,
		mathh.MaxUint16 - 1,
		mathh.MaxUint16,
	}
	for _, v := range tests {
		var n Numeric
		n.SetUint16(v)

		if str := strconvh.FormatUint16(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Uint16(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Uint16(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Uint16(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Uint16 to avoid "1000 overflow int8"
	n.SetUint16(mathh.MaxUint16).Mul(&tmp, &n).Uint16()
	n.SetUint16(mathh.MinUint16).Mul(&tmp, &n).Uint16()
}

func TestNumeric_SetUint32(t *testing.T) {
	tests := []uint32{
		mathh.MinUint32,
		mathh.MinUint32 + 1,
		mathh.MinUint32 / 2,
		mathh.MinUint32/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxUint32,
		10000 % mathh.MaxUint32,
		100000 % mathh.MaxUint32,
		1000000 % mathh.MaxUint32,
		10000000 % mathh.MaxUint32,
		100000000 % mathh.MaxUint32,
		1000000000 % mathh.MaxUint32,
		10000000000 % mathh.MaxUint32,
		100000000000 % mathh.MaxUint32,
		mathh.MaxUint32/2 - 1,
		mathh.MaxUint32 / 2,
		mathh.MaxUint32/2 + 1,
		mathh.MaxUint32 - 1,
		mathh.MaxUint32,
	}
	for _, v := range tests {
		var n Numeric
		n.SetUint32(v)

		if str := strconvh.FormatUint32(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Uint32(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Uint32(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Uint32(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Uint32 to avoid "1000 overflow int8"
	n.SetUint32(mathh.MaxUint32).Mul(&tmp, &n).Uint32()
	n.SetUint32(mathh.MinUint32).Mul(&tmp, &n).Uint32()
}

func TestNumeric_SetUint64(t *testing.T) {
	tests := []uint64{
		mathh.MinUint64,
		mathh.MinUint64 + 1,
		mathh.MinUint64 / 2,
		mathh.MinUint64/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxUint64,
		10000 % mathh.MaxUint64,
		100000 % mathh.MaxUint64,
		1000000 % mathh.MaxUint64,
		10000000 % mathh.MaxUint64,
		100000000 % mathh.MaxUint64,
		1000000000 % mathh.MaxUint64,
		10000000000 % mathh.MaxUint64,
		100000000000 % mathh.MaxUint64,
		mathh.MaxUint64/2 - 1,
		mathh.MaxUint64 / 2,
		mathh.MaxUint64/2 + 1,
		mathh.MaxUint64 - 1,
		mathh.MaxUint64,
	}
	for _, v := range tests {
		var n Numeric
		n.SetUint64(v)

		if str := strconvh.FormatUint64(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Uint64(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}
	}
}

func TestNumeric_Uint64(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Uint64(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Uint64 to avoid "1000 overflow int8"
	n.SetUint64(mathh.MaxUint64).Mul(&tmp, &n).Uint64()
	n.SetUint64(mathh.MinUint64).Mul(&tmp, &n).Uint64()
}
