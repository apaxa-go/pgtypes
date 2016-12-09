//replacer:generated-file

package pgtypes

import "github.com/apaxa-go/helper/mathh"

// SetInt sets z to x and returns z.
func (z *Numeric) SetInt(x int) *Numeric {
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

// SetUint sets z to x and returns z.
func (z *Numeric) SetUint(x uint) *Numeric {
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

// Uint returns the uint representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an uint, the result is undefined.
func (x *Numeric) Uint() uint {
	const maxWeight = mathh.UintBytes / 2 // Interesting, this should work at least for 1-8 bytes [unsigned] integers
	if x.sign != numericPositive || len(x.digits) == 0 {
		return 0
	}
	if x.weight > maxWeight {
		return mathh.MaxUint
	}

	to := mathh.Min2Int(int(x.weight), len(x.digits)-1)
	var r uint
	for i := 0; i <= to; i++ {
		r = r*numericBase + uint(x.digits[i])
	}
	for i := to + 1; i <= int(x.weight); i++ {
		r *= numericBase
	}

	return r
}

// Int returns the int representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an int, the result is undefined.
func (x *Numeric) Int() int {
	const maxWeight = mathh.IntBytes / 2 // Interesting, this should work at least for 1-8 bytes [unsigned] integers
	if x.sign == numericNaN || len(x.digits) == 0 {
		return 0
	}

	var sign int
	if x.sign == numericPositive {
		if x.weight > maxWeight {
			return mathh.MaxInt
		}
		sign = 1
	} else {
		if x.weight > maxWeight {
			return mathh.MinInt
		}
		sign = -1
	}

	to := mathh.Min2Int(int(x.weight), len(x.digits)-1)
	var r int
	for i := 0; i <= to; i++ {
		r = r*numericBase + sign*int(x.digits[i])
	}
	for i := to + 1; i <= int(x.weight); i++ {
		r *= numericBase
	}

	return r
}

// SetInt16 sets z to x and returns z.
func (z *Numeric) SetInt16(x int16) *Numeric {
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

// SetUint16 sets z to x and returns z.
func (z *Numeric) SetUint16(x uint16) *Numeric {
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

// Uint16 returns the uint16 representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an uint16, the result is undefined.
func (x *Numeric) Uint16() uint16 {
	const maxWeight = mathh.Uint16Bytes / 2 // Interesting, this should work at least for 1-8 bytes [unsigned] integers
	if x.sign != numericPositive || len(x.digits) == 0 {
		return 0
	}
	if x.weight > maxWeight {
		return mathh.MaxUint16
	}

	to := mathh.Min2Int(int(x.weight), len(x.digits)-1)
	var r uint16
	for i := 0; i <= to; i++ {
		r = r*numericBase + uint16(x.digits[i])
	}
	for i := to + 1; i <= int(x.weight); i++ {
		r *= numericBase
	}

	return r
}

// Int16 returns the int16 representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an int16, the result is undefined.
func (x *Numeric) Int16() int16 {
	const maxWeight = mathh.Int16Bytes / 2 // Interesting, this should work at least for 1-8 bytes [unsigned] integers
	if x.sign == numericNaN || len(x.digits) == 0 {
		return 0
	}

	var sign int16
	if x.sign == numericPositive {
		if x.weight > maxWeight {
			return mathh.MaxInt16
		}
		sign = 1
	} else {
		if x.weight > maxWeight {
			return mathh.MinInt16
		}
		sign = -1
	}

	to := mathh.Min2Int(int(x.weight), len(x.digits)-1)
	var r int16
	for i := 0; i <= to; i++ {
		r = r*numericBase + sign*int16(x.digits[i])
	}
	for i := to + 1; i <= int(x.weight); i++ {
		r *= numericBase
	}

	return r
}

// SetInt32 sets z to x and returns z.
func (z *Numeric) SetInt32(x int32) *Numeric {
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

// SetUint32 sets z to x and returns z.
func (z *Numeric) SetUint32(x uint32) *Numeric {
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

// Uint32 returns the uint32 representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an uint32, the result is undefined.
func (x *Numeric) Uint32() uint32 {
	const maxWeight = mathh.Uint32Bytes / 2 // Interesting, this should work at least for 1-8 bytes [unsigned] integers
	if x.sign != numericPositive || len(x.digits) == 0 {
		return 0
	}
	if x.weight > maxWeight {
		return mathh.MaxUint32
	}

	to := mathh.Min2Int(int(x.weight), len(x.digits)-1)
	var r uint32
	for i := 0; i <= to; i++ {
		r = r*numericBase + uint32(x.digits[i])
	}
	for i := to + 1; i <= int(x.weight); i++ {
		r *= numericBase
	}

	return r
}

// Int32 returns the int32 representation of x.
// If x is NaN, the result is 0.
// If x cannot be represented in an int32, the result is undefined.
func (x *Numeric) Int32() int32 {
	const maxWeight = mathh.Int32Bytes / 2 // Interesting, this should work at least for 1-8 bytes [unsigned] integers
	if x.sign == numericNaN || len(x.digits) == 0 {
		return 0
	}

	var sign int32
	if x.sign == numericPositive {
		if x.weight > maxWeight {
			return mathh.MaxInt32
		}
		sign = 1
	} else {
		if x.weight > maxWeight {
			return mathh.MinInt32
		}
		sign = -1
	}

	to := mathh.Min2Int(int(x.weight), len(x.digits)-1)
	var r int32
	for i := 0; i <= to; i++ {
		r = r*numericBase + sign*int32(x.digits[i])
	}
	for i := to + 1; i <= int(x.weight); i++ {
		r *= numericBase
	}

	return r
}

// NewInt allocates and returns a new Numeric set to x.
func NewInt(x int) *Numeric {
	var r Numeric
	return r.SetInt(x)
}

// NewInt8 allocates and returns a new Numeric set to x.
func NewInt8(x int8) *Numeric {
	var r Numeric
	return r.SetInt8(x)
}

// NewInt16 allocates and returns a new Numeric set to x.
func NewInt16(x int16) *Numeric {
	var r Numeric
	return r.SetInt16(x)
}

// NewInt32 allocates and returns a new Numeric set to x.
func NewInt32(x int32) *Numeric {
	var r Numeric
	return r.SetInt32(x)
}

// NewUint allocates and returns a new Numeric set to x.
func NewUint(x uint) *Numeric {
	var r Numeric
	return r.SetUint(x)
}

// NewUint8 allocates and returns a new Numeric set to x.
func NewUint8(x uint8) *Numeric {
	var r Numeric
	return r.SetUint8(x)
}

// NewUint16 allocates and returns a new Numeric set to x.
func NewUint16(x uint16) *Numeric {
	var r Numeric
	return r.SetUint16(x)
}

// NewUint32 allocates and returns a new Numeric set to x.
func NewUint32(x uint32) *Numeric {
	var r Numeric
	return r.SetUint32(x)
}

// NewUint64 allocates and returns a new Numeric set to x.
func NewUint64(x uint64) *Numeric {
	var r Numeric
	return r.SetUint64(x)
}
