package numeric

import (
	"github.com/apaxa-io/mathhelper"
	"github.com/apaxa-io/strconvhelper"
	"github.com/apaxa-io/stringshelper"
	"strings"
)

// TODO move to parent package
// TODO implement NullNumeric

type Sign uint16

const (
	Positive Sign = 0x0000
	Negative      = 0x4000
	NaN           = 0xC000
)

const (
	nan   = "NaN"
	delim = '.'
)

const (
	base     = 10000
	groupLen = 4 // Number of 10-based digits stored together, =lg(base)
)

type Numeric struct {
	sign   Sign
	digits []int16
	weight int16
}

func parseInteger(s string, fracPos int) (digits []int16, weight int16) {
	// Pad string left & right (on the left and on the right side of fracPos should be integer number of groupLen digits)
	shift := ((fracPos % groupLen) + groupLen) % groupLen
	leftAdd := 0
	if shift != 0 {
		leftAdd = groupLen - shift
	}
	shift = (leftAdd + len(s)) % groupLen
	rightAdd := 0
	if shift != 0 {
		rightAdd = groupLen - shift
	}

	if leftAdd != 0 || rightAdd != 0 {
		s = strings.Repeat("0", leftAdd) + s + strings.Repeat("0", rightAdd)
	}

	digits = make([]int16, len(s)/groupLen)

	for i := range digits {
		//k:=int16(1)
		//for j:=groupLen-1; j>=0; j--{
		//	r[i]=(s[i*groupLen+j]-'0')*k
		//	k*=10
		//}
		digits[i] = int16(s[i*groupLen+0]-'0')*1000 +
			int16(s[i*groupLen+1]-'0')*100 +
			int16(s[i*groupLen+2]-'0')*10 +
			int16(s[i*groupLen+3]-'0')*1
	}

	weight = int16(mathhelper.DivideCeilInt(fracPos, groupLen) - 1)

	return
}

// Find delimiter position (if not exists return len(s)) and check each char for validity
func findDelim(s string) (delimPos int, valid bool) {
	valid = false
	l := len(s)
	delimPos = l // ("123"=="123.")

	for i := 0; i < l; i++ {
		if s[i] == delim {
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
		s = stringshelper.TrimRightBytes(s, '0')
	case len(s) - 1:
		s = stringshelper.TrimRightBytes(s[:delimPos], '0')
	default:
		s = stringshelper.TrimRightBytes(s[:delimPos]+s[delimPos+1:], '0')
	}

	if len(s) == 0 {
		z.SetZero()
		return true
	}

	{
		// For now delimPos means position of first fraction char in string (may be out of index)
		delimPos -= len(s)
		s = stringshelper.TrimLeftBytes(s, '0')
		delimPos += len(s)
	}

	z.digits, z.weight = parseInteger(s, delimPos)
	return true
}

func (z *Numeric) setString(s string) bool {
	if len(s) == 0 {
		return false
	}

	if s == nan {
		z.sign = NaN
		z.digits = nil
		z.weight = 0

		return true
	}

	switch s[0] {
	case '-':
		z.sign = Negative
		s = s[1:]
	case '+':
		z.sign = Positive
		s = s[1:]
	default:
		z.sign = Positive
	}

	return z.parseUnsigned(s)
}

func (z *Numeric) SetString(s string) (*Numeric, bool) {
	if z.setString(s) {
		return z, true
	}
	return nil, false
}

func (x *Numeric) String() (r string) {
	if x.sign == NaN {
		return nan
	}

	if x.sign == Negative {
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
				r += strconvhelper.FormatInt16(x.digits[i])
			} else {
				r += stringshelper.PadLeftWithByte(strconvhelper.FormatInt16(x.digits[i]), '0', groupLen)
			}
		}

		// Append some groups of zero if x.digits is not enough to print all digits before delimiter.
		appendZero := int(x.weight) + 1 - len(x.digits)
		if appendZero > 0 {
			r += strings.Repeat("0", appendZero*groupLen)
		}
	}

	// Print fraction part
	if len(x.digits) > int(x.weight)+1 {
		r += string(delim)
		if x.weight < -1 {
			r += strings.Repeat("0", groupLen*(-int(x.weight)-1))
		}
		for i := int(mathhelper.Max2Int16(x.weight+1, 0)); i < len(x.digits); i++ {
			if i < len(x.digits)-1 {
				r += stringshelper.PadLeftWithByte(strconvhelper.FormatInt16(x.digits[i]), '0', groupLen)
			} else {
				r += stringshelper.TrimRightBytes(stringshelper.PadLeftWithByte(strconvhelper.FormatInt16(x.digits[i]), '0', groupLen), '0')
			}
		}
	}

	return
}

func (z *Numeric) SetZero() *Numeric {
	z.sign = Positive
	z.weight = 0
	z.digits = []int16{}
	return z
}

func (z *Numeric) SetNaN() *Numeric {
	z.sign = NaN
	z.weight = 0
	z.digits = []int16{}
	return z
}

func (z *Numeric) IsZero() bool {
	return z.sign == Positive && len(z.digits) == 0
}

func (z *Numeric) IsNaN() bool {
	return z.sign == NaN
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

// cmp compare absolute value of two Numeric (i.e. without sign).
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

	for i := 0; i < mathhelper.Min2Int(len(d1), len(d2)); i++ {
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
func (x *Numeric) Cmp(y *Numeric) (r int) {
	if x.sign == NaN || y.sign == NaN {
		return -1 // TODO is it good
	}

	if x.sign == Positive && y.sign == Negative {
		return +1
	} else if x.sign == Negative && y.sign == Positive {
		return -1
	}

	r = cmpAbs(x.digits, x.weight, y.digits, y.weight)

	if x.sign == Negative {
		r *= -1
	}

	return
}

func addAbs(d1 []int16, w1 int16, d2 []int16, w2 int16) (d3 []int16, w3 int16) {
	weightFrom := mathhelper.Min2Int(int(w1)-len(d1)+1, int(w2)-len(d2)+1)
	weightTo := mathhelper.Max2Int(int(w1), int(w2))
	var overflow int16
	var index int
	for i := weightFrom; i <= weightTo; i++ {
		tmp := digitByWeightAbs(d1, w1, i) + digitByWeightAbs(d2, w2, i) + overflow
		overflow = tmp / base
		tmp = tmp % base
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

func (z *Numeric) Copy(x *Numeric) *Numeric {
	z.digits, z.weight, z.sign = x.digits, x.weight, x.sign
	return z
}

func (z *Numeric) Add(x, y *Numeric) *Numeric {
	if x.sign == NaN || y.sign == NaN {
		return z.SetNaN()
	}
	if x.IsZero() {
		return z.Copy(y)
	}
	if y.IsZero() {
		return z.Copy(x)
	}

	var negative bool
	z.digits, z.weight, negative = add(x.digits, x.weight, x.sign == Negative, y.digits, y.weight, y.sign == Negative)
	if negative {
		z.sign = Negative
	} else {
		z.sign = Positive
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

func (z *Numeric) Neg(x *Numeric) *Numeric {
	if x.IsNaN() {
		return z.SetNaN()
	}
	if x.IsZero() {
		return z.SetZero()
	}
	z.Copy(x)

	if x.sign == Positive {
		z.sign = Negative
	} else {
		x.sign = Positive
	}

	return z
}

func (z *Numeric) Sub(x, y *Numeric) *Numeric {
	if x.sign == NaN || y.sign == NaN {
		return z.SetNaN()
	}
	if x.IsZero() {
		return z.Copy(y).Neg(z)
	}
	if y.IsZero() {
		return z.Copy(x)
	}

	var negative bool
	if x.sign == y.sign {
		z.digits, z.weight, negative = sub(x.digits, x.weight, x.sign == Negative, y.digits, y.weight, y.sign == Negative)
	} else {
		z.digits, z.weight, negative = add(x.digits, x.weight, x.sign == Negative, y.digits, y.weight, !(y.sign == Negative))
	}

	if negative {
		z.sign = Negative
	} else {
		z.sign = Positive
	}

	return z
}

// This function is for subAbs only. Do not use this function directly.
// Number 1 must be > number 2
func subAbsOrdered(d1 []int16, w1 int16, d2 []int16, w2 int16) (d3 []int16, w3 int16) {
	weightFrom := mathhelper.Min2Int(int(w1)-len(d1)+1, int(w2)-len(d2)+1)
	weightTo := mathhelper.Max2Int(int(w1), int(w2))
	var underflow int16
	var index int
	for i := weightFrom; i <= weightTo; i++ {
		tmp := digitByWeightAbs(d1, w1, i) - digitByWeightAbs(d2, w2, i) - underflow
		if tmp < 0 {
			tmp += base
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
		for j1 := mathhelper.Max2Int(0, i-len(d2)+1); j1 <= mathhelper.Min2Int(len(d1)-1, i); j1++ {
			j2 := i - j1
			tmp += int64(d1[len(d1)-1-j1]) * int64(d2[len(d2)-1-j2])
		}
		overflow = tmp / base
		tmp = tmp % base
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

func (z *Numeric) Mul(x, y *Numeric) *Numeric {
	if x.sign == NaN || y.sign == NaN {
		z.SetNaN()
		return z
	}
	if x.IsZero() || y.IsZero() {
		z.SetZero()
		return z
	}

	z.digits, z.weight = mulAbs(x.digits, x.weight, y.digits, y.weight)
	if x.sign == y.sign {
		z.sign = Positive
	} else {
		z.sign = Negative
	}

	return z
}

func (z *Numeric) Abs(x *Numeric) *Numeric {
	z.Copy(x)
	if x.sign == Negative {
		z.sign = Positive
	}
	return z
}

// Sign returns first value as following:
//
//	-1 if x <   0
//	 0 if x is ±0 or NaN
//	+1 if x >   0
//
func (x *Numeric) Sign() int {
	switch {
	case x.IsZero():
		return 0
	case x.sign == Positive:
		return 1
	case x.sign == Negative:
		return -1
	default:
		return 0 // TODO what to return for NaN?
	}
}

