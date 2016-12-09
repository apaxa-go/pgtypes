package pgtypes

import (
	"github.com/apaxa-go/helper/mathh"
	"github.com/apaxa-go/helper/strconvh"
	"testing"
)

//replacer:ignore
//go:generate go run $GOPATH/src/github.com/apaxa-go/generator/replacer/main.go -- $GOFILE
//replacer:replace
//replacer:old int64	Int64
//replacer:new int	Int
//replacer:new int8	Int8
//replacer:new int16	Int16
//replacer:new int32	Int32
//replacer:new uint	Uint
//replacer:new uint8	Uint8
//replacer:new uint16	Uint16
//replacer:new uint32	Uint32
//replacer:new uint64	Uint64

func TestNumeric_SetInt64(t *testing.T) {
	tests := []int64{
		mathh.MinInt64,
		mathh.MinInt64 + 1,
		mathh.MinInt64 / 2,
		mathh.MinInt64/2 + 1,
		0,
		1,
		10,
		100,
		1000 % mathh.MaxInt64,
		10000 % mathh.MaxInt64,
		100000 % mathh.MaxInt64,
		1000000 % mathh.MaxInt64,
		10000000 % mathh.MaxInt64,
		100000000 % mathh.MaxInt64,
		1000000000 % mathh.MaxInt64,
		10000000000 % mathh.MaxInt64,
		100000000000 % mathh.MaxInt64,
		mathh.MaxInt64/2 - 1,
		mathh.MaxInt64 / 2,
		mathh.MaxInt64/2 + 1,
		mathh.MaxInt64 - 1,
		mathh.MaxInt64,
	}
	for _, v := range tests {
		var n Numeric
		n.SetInt64(v)

		if str := strconvh.FormatInt64(v); n.String() != str {
			t.Errorf("expect %v, got %v", str, n.String())
		}
		if i := n.Int64(); i != v {
			t.Errorf("expect %v, got %v", v, i)
		}

		n2 := NewInt64(v)
		if str := strconvh.FormatInt64(v); n2.String() != str {
			t.Errorf("expect %v, got %v", str, n2.String())
		}
	}
}

func TestNumeric_Int64(t *testing.T) {
	// NaN
	var n Numeric
	if i := n.SetNaN().Int64(); i != 0 {
		t.Errorf("expect %v, got %v", 0, i)
	}

	// Just for coverage
	var tmp Numeric
	tmp.SetInt32(numericBase) // Use Int32 instead of Int64 to avoid "1000 overflow int8"
	n.SetInt64(mathh.MaxInt64).Mul(&tmp, &n).Int64()
	n.SetInt64(mathh.MinInt64).Mul(&tmp, &n).Int64()
}
