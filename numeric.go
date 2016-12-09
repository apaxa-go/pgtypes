package pgtypes

import (
	"github.com/apaxa-go/helper/mathh"
	"github.com/apaxa-go/helper/strconvh"
	"github.com/apaxa-go/helper/stringsh"
	"strings"
)

type numericSign uint16

const (
	numericPositive numericSign = 0x0000
	numericNegative             = 0x4000
	numericNaN                  = 0xC000
)

const (
	numericNanStr    = "NaN"
	numericDelimiter = '.'
	numericBase      = 10000
	numericGroupLen  = 4 // Number of 10-based digits stored together, =lg(base)
)

// Numeric is a PostgreSQL Numeric type implementation in GoLang.
// It is arbitrary precision numbers type which can store numbers with a very large number of digits.
// It is especially recommended for storing monetary amounts and other quantities where exactness is required.
// Calculations with Numeric values yield exact results where possible, e.g. addition, subtraction, multiplication.
// However, calculations on Numeric values are very slow compared to the integer types, or to the floating-point types.
// Internally Numeric type has the same structure (except dscale field) as a PostgreSQL numeric type so it perfect for using in DB communications.
type Numeric struct {
	sign   numericSign
	digits []int16
	weight int16
}

func parseInteger(s string, fracPos int) (digits []int16, weight int16) {
	// Pad string left & right (on the left and on the right side of fracPos should be integer number of groupLen digits)
	shift := ((fracPos % numericGroupLen) + numericGroupLen) % numericGroupLen
	leftAdd := 0
	if shift != 0 {
		leftAdd = numericGroupLen - shift
	}
	shift = (leftAdd + len(s)) % numericGroupLen
	rightAdd := 0
	if shift != 0 {
		rightAdd = numericGroupLen - shift
	}

	if leftAdd != 0 || rightAdd != 0 {
		s = strings.Repeat("0", leftAdd) + s + strings.Repeat("0", rightAdd)
	}

	digits = make([]int16, len(s)/numericGroupLen)

	for i := range digits {
		digits[i] = int16(s[i*numericGroupLen+0]-'0')*1000 +
			int16(s[i*numericGroupLen+1]-'0')*100 +
			int16(s[i*numericGroupLen+2]-'0')*10 +
			int16(s[i*numericGroupLen+3]-'0')*1
	}

	weight = int16(mathh.DivideCeilInt(fracPos, numericGroupLen) - 1)

	return
}

// Find delimiter position (if not exists return len(s)) and check each char for validity
func findDelim(s string) (delimPos int, valid bool) {
	valid = false
	l := len(s)
	delimPos = l // ("123"=="123.")

	for i := 0; i < l; i++ {
		if s[i] == numericDelimiter {
			if delimPos == l {
				delimPos = i
			} else { // Second delimiter found - error
				return
			}
		} else if s[i] < '0' || s[i] > '9' { // If char is not delimiter than it can be only 10-base digit
			return
		}
	}
	valid = true
	return
}

func (z *Numeric) parseUnsigned(s string) (valid bool) {
	var delimPos int
	if delimPos, valid = findDelim(s); !valid {
		return
	}

	switch delimPos {
	case len(s):
		s = stringsh.TrimRightBytes(s, '0')
	case len(s) - 1:
		s = stringsh.TrimRightBytes(s[:delimPos], '0')
	default:
		s = stringsh.TrimRightBytes(s[:delimPos]+s[delimPos+1:], '0')
	}

	if len(s) == 0 {
		z.SetZero()
		return true
	}

	{
		// For now delimPos means position of first fraction char in string (may be out of index)
		delimPos -= len(s)
		s = stringsh.TrimLeftBytes(s, '0')
		delimPos += len(s)
	}

	z.digits, z.weight = parseInteger(s, delimPos)
	return true
}

func (z *Numeric) setString(s string) bool {
	if len(s) == 0 {
		return false
	}

	if s == numericNanStr {
		z.sign = numericNaN
		z.digits = nil
		z.weight = 0

		return true
	}

	switch s[0] {
	case '-':
		z.sign = numericNegative
		s = s[1:]
	case '+':
		z.sign = numericPositive
		s = s[1:]
	default:
		z.sign = numericPositive
	}

	return z.parseUnsigned(s)
}

// SetString sets z to the value of s and returns z and a boolean indicating success.
// s must be a floating-point number of one of format
// 	[+-]?[0-9]*\.[0-9]*	// "123.456", "123.", ".456", ".", "-123.456"
// 	[+-]?[0-9]+		// "123", "-123"
func (z *Numeric) SetString(s string) (*Numeric, bool) {
	if z.setString(s) {
		return z, true
	}
	return nil, false
}

// String converts the Number x to a string representation (10-base).
func (x *Numeric) String() (r string) {
	if x.sign == numericNaN {
		return numericNanStr
	}

	if x.sign == numericNegative {
		r = "-"
	}

	// Print integer part
	if x.weight < 0 || len(x.digits) == 0 {
		r += "0"
	} else {
		// Print integer part (before delimiter) using x.digits.
		// x.digits may not be enough to print all digits before delimiter (example: "10000000").
		for i := 0; i <= int(x.weight) && i < len(x.digits); i++ {
			if i == 0 {
				r += strconvh.FormatInt16(x.digits[i])
			} else {
				r += stringsh.PadLeftWithByte(strconvh.FormatInt16(x.digits[i]), '0', numericGroupLen)
			}
		}

		// Append some groups of zero if x.digits is not enough to print all digits before delimiter.
		appendZero := int(x.weight) + 1 - len(x.digits)
		if appendZero > 0 {
			r += strings.Repeat("0", appendZero*numericGroupLen)
		}
	}

	// Print fraction part
	if len(x.digits) > int(x.weight)+1 {
		r += string(numericDelimiter)
		if x.weight < -1 {
			r += strings.Repeat("0", numericGroupLen*(-int(x.weight)-1))
		}
		for i := int(mathh.Max2Int16(x.weight+1, 0)); i < len(x.digits); i++ {
			if i < len(x.digits)-1 {
				r += stringsh.PadLeftWithByte(strconvh.FormatInt16(x.digits[i]), '0', numericGroupLen)
			} else {
				r += stringsh.TrimRightBytes(stringsh.PadLeftWithByte(strconvh.FormatInt16(x.digits[i]), '0', numericGroupLen), '0')
			}
		}
	}

	return
}

// SetZero sets Number z to zero and return z.
func (z *Numeric) SetZero() *Numeric {
	z.sign = numericPositive
	z.weight = 0
	z.digits = nil
	return z
}

// SetNaN sets Number z to NaN and return z.
func (z *Numeric) SetNaN() *Numeric {
	z.sign = numericNaN
	z.weight = 0
	z.digits = nil
	return z
}

// IsZero reports whether x is zero.
func (x *Numeric) IsZero() bool {
	return x.sign == numericPositive && len(x.digits) == 0
}

// IsNaN reports whether x is NaN.
func (x *Numeric) IsNaN() bool {
	return x.sign == numericNaN
}

func digitByWeightAbs(d []int16, w int16, reqW int) int16 {
	if reqW > int(w) {
		return 0
	}
	if reqW <= int(w)-len(d) {
		return 0
	}
	return d[int(w)-reqW]
}

// cmpAbs compare absolute value of two Numeric (i.e. without sign).
//
//   -1 if d1 <  d2
//    0 if d1 == d2
//   +1 if d1 >  d2
//
func cmpAbs(d1 []int16, w1 int16, d2 []int16, w2 int16) int {
	if len(d1) == 0 {
		if len(d2) == 0 {
			return 0
		}
		return -1
	} else if len(d2) == 0 {
		return 1
	}

	if w1 < w2 {
		return -1
	} else if w1 > w2 {
		return 1
	}

	for i := 0; i < mathh.Min2Int(len(d1), len(d2)); i++ {
		if d1[i] < d2[i] {
			return -1
		} else if d1[i] > d2[i] {
			return 1
		}
	}

	if len(d1) < len(d2) {
		return -1
	} else if len(d1) > len(d2) {
		return 1
	}

	return 0
}

// Cmp compare two Numeric.
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
// NaN treats as equal to other NaN and greater then any other number. This is as in PostgreSQL.
func (x *Numeric) Cmp(y *Numeric) (r int) {
	// NaN logic
	switch {
	case x.sign == numericNaN && y.sign == numericNaN:
		return 0
	case x.sign == numericNaN:
		return 1
	case y.sign == numericNaN:
		return -1
	}

	if x.sign == numericPositive && y.sign == numericNegative {
		return +1
	} else if x.sign == numericNegative && y.sign == numericPositive {
		return -1
	}

	r = cmpAbs(x.digits, x.weight, y.digits, y.weight)

	if x.sign == numericNegative {
		r *= -1
	}

	return
}

func addAbs(d1 []int16, w1 int16, d2 []int16, w2 int16) (d3 []int16, w3 int16) {
	weightFrom := mathh.Min2Int(int(w1)-len(d1)+1, int(w2)-len(d2)+1)
	weightTo := mathh.Max2Int(int(w1), int(w2))
	var overflow int16
	var index int
	for i := weightFrom; i <= weightTo; i++ {
		tmp := digitByWeightAbs(d1, w1, i) + digitByWeightAbs(d2, w2, i) + overflow
		overflow = tmp / numericBase
		tmp = tmp % numericBase
		if d3 == nil && tmp != 0 {
			requiredLen := weightTo - i + 1 + 1 // Reserve 1 for overall overflow
			d3 = make([]int16, requiredLen)
			index = requiredLen - 1
		}
		if d3 != nil {
			d3[index] = tmp
		}
		index--
	}

	if d3 == nil && overflow == 0 {
		w3 = 0
	} else if overflow == 0 {
		d3 = d3[1:]
		w3 = int16(weightTo)
	} else {
		if d3 == nil {
			d3 = make([]int16, 1)
		}
		d3[0] = overflow
		w3 = int16(weightTo + 1)
	}

	return
}

func add(d1 []int16, w1 int16, n1 bool, d2 []int16, w2 int16, n2 bool) (d3 []int16, w3 int16, n3 bool) {
	if n1 == n2 {
		d3, w3 = addAbs(d1, w1, d2, w2)
		n3 = n1
		return
	}
	return sub(d1, w1, n1, d2, w2, !n2)
}

// Copy sets z to x and returns z.
// x is not changed even if z and x are the same.
func (z *Numeric) Copy(x *Numeric) *Numeric {
	if x != z {
		z.weight, z.sign = x.weight, x.sign
		z.digits = make([]int16, len(x.digits))
		copy(z.digits, x.digits)
	}
	return z
}

// Add sets z to the sum x+y and returns z.
func (z *Numeric) Add(x, y *Numeric) *Numeric {
	if x.sign == numericNaN || y.sign == numericNaN {
		return z.SetNaN()
	}
	if x.IsZero() {
		return z.Copy(y)
	}
	if y.IsZero() {
		return z.Copy(x)
	}

	var negative bool
	z.digits, z.weight, negative = add(x.digits, x.weight, x.sign == numericNegative, y.digits, y.weight, y.sign == numericNegative)
	if negative {
		z.sign = numericNegative
	} else {
		z.sign = numericPositive
	}

	return z
}

func sub(d1 []int16, w1 int16, n1 bool, d2 []int16, w2 int16, n2 bool) (d3 []int16, w3 int16, n3 bool) {
	if n1 == n2 {
		d3, w3, n3 = subAbs(d1, w1, d2, w2)
		if len(d3) != 0 && n1 {
			n3 = !n3
		}
		return
	}
	return add(d1, w1, n1, d2, w2, !n2)
}

// Neg sets z to -x and returns z.
func (z *Numeric) Neg(x *Numeric) *Numeric {
	if x.IsNaN() {
		return z.SetNaN()
	}
	if x.IsZero() {
		return z.SetZero()
	}
	z.Copy(x)

	if x.sign == numericPositive {
		z.sign = numericNegative
	} else {
		z.sign = numericPositive
	}

	return z
}

// Sub sets z to the difference x-y and returns z.
func (z *Numeric) Sub(x, y *Numeric) *Numeric {
	if x.sign == numericNaN || y.sign == numericNaN {
		return z.SetNaN()
	}
	if x.IsZero() {
		return z.Copy(y).Neg(z)
	}
	if y.IsZero() {
		return z.Copy(x)
	}

	var negative bool
	z.digits, z.weight, negative = sub(x.digits, x.weight, x.sign == numericNegative, y.digits, y.weight, y.sign == numericNegative)

	if negative {
		z.sign = numericNegative
	} else {
		z.sign = numericPositive
	}

	return z
}

// This function is for subAbs only. Do not use this function directly.
// Number 1 must be > number 2
func subAbsOrdered(d1 []int16, w1 int16, d2 []int16, w2 int16) (d3 []int16, w3 int16) {
	weightFrom := mathh.Min2Int(int(w1)-len(d1)+1, int(w2)-len(d2)+1)
	weightTo := mathh.Max2Int(int(w1), int(w2))
	var underflow int16
	var index int
	for i := weightFrom; i <= weightTo; i++ {
		tmp := digitByWeightAbs(d1, w1, i) - digitByWeightAbs(d2, w2, i) - underflow
		if tmp < 0 {
			tmp += numericBase
			underflow = 1
		} else {
			underflow = 0
		}
		if d3 == nil && tmp != 0 {
			requiredLen := weightTo - i + 1
			d3 = make([]int16, requiredLen)
			index = requiredLen - 1
		}
		if d3 != nil {
			d3[index] = tmp
		}
		index--
	}

	if d3 == nil {
		w3 = 0
	} else {
		w3 = int16(weightTo)
	}

	// Trim leading zero
	leadingZero := 0
	for i := 0; i < len(d3) && d3[i] == 0; i++ {
		leadingZero++
	}
	w3 -= int16(leadingZero)
	d3 = d3[leadingZero:]

	return
}

func subAbs(d1 []int16, w1 int16, d2 []int16, w2 int16) (d3 []int16, w3 int16, negative bool) {
	switch cmpAbs(d1, w1, d2, w2) {
	case -1:
		d3, w3 = subAbsOrdered(d2, w2, d1, w1)
		negative = true
	case 0:
		negative = false
	case 1:
		d3, w3 = subAbsOrdered(d1, w1, d2, w2)
		negative = false
	}
	return
}

func mulAbs(d1 []int16, w1 int16, d2 []int16, w2 int16) (d3 []int16, w3 int16) {
	// (abc)*(xy) = SUM:
	//     a*y b*y c*y
	// a*x b*x c*x
	// So number of columns = ( len(d1) + len(d2) - 1 ), but it is possible to overflow.
	//
	// d1*d2 < base^len(d1) * base^len(d2) = base^(len(d1) + len(d2))
	// len( base^(len(d1) + len(d2)) ) = len(d1) + len(d2) + 1
	// len( base^(len(d1) + len(d2)) - 1 ) = len(d1) + len(d2)
	// len(d1*d2) <= len(d1) + len(d2)
	col := len(d1) + len(d2) - 1
	var overflow int64
	var index int
	for i := 0; i < col; i++ {
		// tmp contains sum for digit at position i (counted from the right) in the result
		// max(tmp) - maximum value os sum
		// max(tmp) = (digitBase-1)^2 * maxDigitsInNumber
		// maxDigitInNumber = maxInt16 (because in Postgres number of digits in number stored in int16)
		// max(tmp) = 9999 * 9999 * 32767 = 3276044692767
		// MAX(tmp) - maximum tmp plus overflow
		// MAX(tmp) <= max(tmp) + max(tmp)/digitBase
		// MaxInt64 > MAX(tmp) > MaxInt32 => use int64 (for overflow and for tmp)

		var tmp = overflow
		for j1 := mathh.Max2Int(0, i-len(d2)+1); j1 <= mathh.Min2Int(len(d1)-1, i); j1++ {
			j2 := i - j1
			tmp += int64(d1[len(d1)-1-j1]) * int64(d2[len(d2)-1-j2])
		}
		overflow = tmp / numericBase
		tmp = tmp % numericBase
		if d3 == nil && tmp != 0 {
			requiredLen := col - i + 1 // Reserve 1 for overall overflow
			d3 = make([]int16, requiredLen)
			index = requiredLen - 1
		}
		if d3 != nil {
			d3[index] = int16(tmp)
		}
		index--
	}

	if d3 == nil && overflow == 0 {
		w3 = 0
	} else if overflow == 0 {
		d3 = d3[1:]
		w3 = int16(w1 + w2)
	} else {
		if d3 == nil {
			d3 = make([]int16, 1)
		}
		d3[0] = int16(overflow)
		w3 = int16(w1 + w2 + 1)
	}

	return
}

// Mul sets z to the product x*y and returns z.
func (z *Numeric) Mul(x, y *Numeric) *Numeric {
	if x.sign == numericNaN || y.sign == numericNaN {
		z.SetNaN()
		return z
	}
	if x.IsZero() || y.IsZero() {
		z.SetZero()
		return z
	}

	z.digits, z.weight = mulAbs(x.digits, x.weight, y.digits, y.weight)
	if x.sign == y.sign {
		z.sign = numericPositive
	} else {
		z.sign = numericNegative
	}

	return z
}

// Abs sets z to |x| (the absolute value of x) and returns z.
func (z *Numeric) Abs(x *Numeric) *Numeric {
	z.Copy(x)
	if x.sign == numericNegative {
		z.sign = numericPositive
	}
	return z
}

// Sign returns first value as following:
//
//	-1 if x <  0
//	 0 if x is 0
//	+1 if x >  0
//	+2 if x is NaN
//
func (x *Numeric) Sign() int {
	switch {
	case x.sign == numericNegative:
		return -1
	case x.sign == numericNaN:
		return 2
	case x.IsZero():
		return 0
	default:
		return 1
	}
}

// NewNumeric allocates and returns a new Numeric set to 0.
func NewNumeric() *Numeric {
	var r Numeric
	return &r
}
