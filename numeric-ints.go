package pgtypes

import "github.com/apaxa-go/helper/mathh"

//replacer:ignore
//go:generate go run $GOPATH/src/github.com/apaxa-go/generator/replacer/main.go -- $GOFILE

// SetInt8 sets z to x and returns z.
func (z *Numeric) SetInt8(x int8) *Numeric {
	if x == 0 {
		return z.SetZero()
	}

	if x < 0 {
		z.sign = numericNegative
	} else {
		z.sign = numericPositive
	}
	z.weight = 0
	z.digits = []int16{mathh.AbsInt16(int16(x))} // First update type, second abs!

	return z
}

// SetUint8 sets z to x and returns z.
func (z *Numeric) SetUint8(x uint8) *Numeric {
	if x == 0 {
		return z.SetZero()
	}

	z.sign = numericPositive
	z.weight = 0
	z.digits = []int16{int16(x)}

	return z
}

// Uint8 returns the uint8 representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an uint8, the result is undefined.
func (x *Numeric) Uint8() uint8 {
	if x.sign != numericPositive || x.weight != 0 || len(x.digits) == 0 {
		return 0
	}

	return uint8(x.digits[0])
}

// Int8 returns the int8 representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an int8, the result is undefined.
func (x *Numeric) Int8() int8 {
	if x.sign == numericNaN || x.weight != 0 || len(x.digits) == 0 {
		return 0
	}
	if x.sign == numericNegative {
		return int8(-x.digits[0]) // First - negate, only after it type conversion!
	}

	return int8(x.digits[0])
}

//replacer:replace
//replacer:old int64	Int64
//replacer:new int	Int
//replacer:new int16	Int16
//replacer:new int32	Int32

// SetInt64 sets z to x and returns z.
func (z *Numeric) SetInt64(x int64) *Numeric {
	if x == 0 {
		return z.SetZero()
	}

	if x < 0 {
		z.sign = numericNegative
	} else {
		z.sign = numericPositive
	}

	z.weight = -1
	z.digits = make([]int16, 0, 1) // as x!=0 there is at least 1 1000-base digit
	for x != 0 {
		d := mathh.AbsInt16(int16(x % numericBase))
		x /= numericBase
		if d != 0 || len(z.digits) > 0 { // avoid tailing zero
			z.digits = append([]int16{d}, z.digits...)
		}
		z.weight++
	}

	return z
}

// SetUint64 sets z to x and returns z.
func (z *Numeric) SetUint64(x uint64) *Numeric {
	if x == 0 {
		return z.SetZero()
	}

	z.sign = numericPositive
	z.weight = -1
	z.digits = make([]int16, 0, 1) // as x!=0 there is at least 1 1000-base digit
	for x != 0 {
		d := int16(x % numericBase)
		x /= numericBase
		if d != 0 || len(z.digits) > 0 { // avoid tailing zero
			z.digits = append([]int16{d}, z.digits...)
		}
		z.weight++
	}

	return z
}

// Uint64 returns the uint64 representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an uint64, the result is undefined.
func (x *Numeric) Uint64() uint64 {
	const maxWeight = mathh.Uint64Bytes / 2 // Interesting, this should work at least for 1-8 bytes [unsigned] integers
	if x.sign != numericPositive || len(x.digits) == 0 {
		return 0
	}
	if x.weight > maxWeight {
		return mathh.MaxUint64
	}

	to := mathh.Min2Int(int(x.weight), len(x.digits)-1)
	var r uint64
	for i := 0; i <= to; i++ {
		r = r*numericBase + uint64(x.digits[i])
	}
	for i := to + 1; i <= int(x.weight); i++ {
		r *= numericBase
	}

	return r
}

// Int64 returns the int64 representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an int64, the result is undefined.
func (x *Numeric) Int64() int64 {
	const maxWeight = mathh.Int64Bytes / 2 // Interesting, this should work at least for 1-8 bytes [unsigned] integers
	if x.sign == numericNaN || len(x.digits) == 0 {
		return 0
	}

	var sign int64
	if x.sign == numericPositive {
		if x.weight > maxWeight {
			return mathh.MaxInt64
		}
		sign = 1
	} else {
		if x.weight > maxWeight {
			return mathh.MinInt64
		}
		sign = -1
	}

	to := mathh.Min2Int(int(x.weight), len(x.digits)-1)
	var r int64
	for i := 0; i <= to; i++ {
		r = r*numericBase + sign*int64(x.digits[i])
	}
	for i := to + 1; i <= int(x.weight); i++ {
		r *= numericBase
	}

	return r
}

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

// NewInt64 allocates and returns a new Numeric set to x.
func NewInt64(x int64) *Numeric {
	var r Numeric
	return r.SetInt64(x)
}
