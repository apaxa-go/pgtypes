package pgtypes

import "github.com/apaxa-go/helper/mathh"

const (
	// Limit on the precision (and hence scale) specifiable in a NUMERIC typmod.
	// Note that the implementation limit on the length of a numeric value is much larger --- beware of what you use this for!
	// pgNumericMaxPrecision is a copy of PostgreSQL NUMERIC_MAX_PRECISION defined at "src/include/utils/numeric.h".
	pgNumericMaxPrecision = 1000
	// Internal limits on the scales chosen for calculation results
	// pgNumericMaxDisplayScale and pgNumericMinDisplayScale are copy of PostgreSQL NUMERIC_MAX_DISPLAY_SCALE/NUMERIC_MIN_DISPLAY_SCALE defined at "src/include/utils/numeric.h".
	pgNumericMaxDisplayScale = pgNumericMaxPrecision
	pgNumericMinDisplayScale = 0
	// For inherently inexact calculations such as division and square root, we try to get at least this many significant digits;
	// the idea is to deliver a result no worse than float8 would.
	// pgNumericMinSigDigits is a copy of PostgreSQL NUMERIC_MIN_SIG_DIGITS defined at "src/include/utils/numeric.h".
	pgNumericMinSigDigits = 16
)

// divAbs divide number (d1,w1) to (d2,w2).
//
//	(d3,w3) = (d1,w1)/(d2,w2)
//
// If round is true results will be rounded otherwise result will be truncated.
// s3 is number of decimal digits to produce in result.
// divAbs is based on PostgreSQL div_var function defined at "src/backend/utils/adt/numeric.c".
func divAbs(d1 []int16, w1 int16, d2 []int16, w2 int16, s3 int16, round bool) (d3 []int16, w3 int16) {
	w3 = w1 - w2

	{
		d3Len := int(w3) + 1 + (int(s3)+numericGroupLen-1)/numericGroupLen // The number of accurate result digits we need to produce,
		d3Len = mathh.Max2Int(d3Len, 1)                                    // but at least 1,
		if round {                                                         // and if rounding needed, figure one more digit to ensure correct result
			d3Len++
		}

		d3 = make([]int16, d3Len)
	}

	var d1C, d2C []int16
	{
		var d1CLen int

		// The working dividend (d1C) normally requires len(d3) + len(d2) digits, but make it at least len(d1) so we can load all of d1 into it.
		// (There will be an additional digit d1C[0] in the dividend space, but for consistency with Knuth's notation we don't count that in d1CLen.)
		d1CLen = len(d3) + len(d2)
		d1CLen = mathh.Max2Int(d1CLen, len(d1))

		// We need a workspace with room for the working dividend (d1CLen+1 digits).
		d1C = make([]int16, d1CLen+1)
	}
	// Also we need a workspace with room for the possibly-normalized divisor (len(d2) digits).
	// It is convenient also to have a zero at divisor[0] with the actual divisor data in divisor[1 .. len(d2)].
	d2C = make([]int16, len(d2)+1)
	copy(d1C[1:], d1)
	copy(d2C[1:], d2)

	//
	// Main part
	//
	if len(d2) == 1 {
		// If there's only a single divisor (d2) digit, we can use a fast path (cf. Knuth section 4.3.1 exercise 16).
		var underflow int
		for i := 0; i < len(d3); i++ {
			underflow = underflow*numericBase + int(d1C[i+1])
			d3[i] = int16(underflow / int(d2C[1]))
			underflow = underflow % int(d2C[1])
		}
	} else {
		// The full multiple-place algorithm is taken from Knuth volume 2, Algorithm 4.3.1D.

		// We need the first divisor digit (d2C[1]) to be >= NBASE/2.
		// If it isn't, make it so by scaling up both the divisor and dividend by the factor "d".
		// (The reason for allocating d1C[0] above is to leave room for possible overflow here.)
		if d2C[1] < numericBase/2 {
			d := numericBase / (int(d2C[1]) + 1)

			var overflow int
			for i := len(d2); i > 0; i-- {
				overflow += int(d2C[i]) * d
				d2C[i] = int16(overflow % numericBase)
				overflow = overflow / numericBase
			}
			overflow = 0
			// At this point only len(d1) of dividend (d1C) can be nonzero
			for i := len(d1); i >= 0; i-- {
				overflow += int(d1C[i]) * d
				d1C[i] = int16(overflow % numericBase)
				overflow = overflow / numericBase
			}
		}

		// Begin the main loop.
		// Each iteration of this loop produces the j'th quotient digit by dividing d1C[j .. j + len(d2)] by the d2C;
		// This is essentially the same as the common manual procedure for long division.
		for j := 0; j < len(d3); j++ {
			// Estimate quotient digit from the first two dividend digits
			next2Digits := int(d1C[j])*numericBase + int(d1C[j+1])

			// If next2Digits is 0, then quotient digit must be 0 and there's no need to adjust the working dividend.
			// It's worth testing here to fall out ASAP when processing trailing zeroes in a dividend.
			if next2Digits == 0 {
				d3[j] = 0
				continue
			}

			// Estimated quotient digit
			estDigit := int(numericBase - 1)
			if d1C[j] != d2C[1] {
				estDigit = next2Digits / int(d2C[1])
			}

			// Adjust quotient digit if it's too large.
			// Knuth proves that after this step, the quotient digit will be either correct or just one too large.
			// (Note: it's OK to use dividend[j+2] here because we know the divisor length is at least 2.)
			for int(d2C[2])*estDigit > (next2Digits-estDigit*int(d2C[1]))*numericBase+int(d1C[j+2]) {
				estDigit--
			}

			// As above, need do nothing more when quotient digit is 0
			if estDigit > 0 {
				// Multiply the divisor (d2C) by estDigit, and subtract that from the working dividend (d1C).
				// "overflow" tracks the multiplication, "borrow" the subtraction (could we fold these together?)
				var overflow, borrow int
				for i := len(d2); i >= 0; i-- {
					overflow += int(d2C[i]) * estDigit
					borrow -= overflow % numericBase
					overflow = overflow / numericBase
					borrow += int(d1C[j+i])
					if borrow < 0 {
						d1C[j+i] = int16(borrow + numericBase)
						borrow = -1
					} else {
						d1C[j+i] = int16(borrow)
						borrow = 0
					}
				}

				// If we got a borrow out of the top dividend (d1C) digit, then indeed estDigit was one too large.
				// Fix it, and add back the divisor (d2C) to correct the working dividend (d1C).
				// (Knuth proves that this will occur only about 3/NBASE of the time;
				// hence, it's a good idea to test this code with small NBASE to be sure this section gets exercised.)
				if borrow != 0 {
					estDigit--
					var overflow int
					for i := len(d2); i >= 0; i-- {
						overflow += int(d1C[j+i]) + int(d2C[i])
						if overflow >= numericBase {
							d1C[j+i] = int16(overflow - numericBase)
							overflow = 1
						} else {
							d1C[j+i] = int16(overflow)
							overflow = 0
						}
					}
				}
			}

			// And we're done with this quotient digit
			d3[j] = int16(estDigit)
		}
	}

	// Round or truncate to target rscale
	if round {
		d3, w3 = roundAbs(d3, w3, s3)
	} else {
		d3, w3 = truncAbs(d3, w3, s3)
	}

	return trimAbs(d3, w3)
}

// trimAbs trim 0 from digits d and adjust weight w if needed.
// If len(d)==0 trimAbs return correct zero value.
func trimAbs(d []int16, w int16) ([]int16, int16) {
	var skipLeft int
	for skipLeft = 0; skipLeft < len(d) && d[skipLeft] == 0; skipLeft++ {
	}
	var skipRight int
	for skipRight = len(d); skipRight > skipLeft+1 && d[skipRight-1] == 0; skipRight-- {
	}

	d = d[skipLeft:skipRight]

	if len(d) == 0 {
		w = 0
		d = nil
	} else {
		w -= int16(skipLeft)
	}

	return d, w
}

// roundAbs rounds the value of d to no more than s decimal digits after the decimal point and adjust w if needed.
// s<0 means rounding before the decimal point.
// If trimAbs called on (d,w) it is better to do it after calling this function, not before.
// roundAbs is based on PostgreSQL round_var function defined at "src/backend/utils/adt/numeric.c".
func roundAbs(d []int16, w int16, s int16) ([]int16, int16) {
	if len(d) == 0 {
		return nil, 0
	}

	decimalDigits := (int(w)+1)*numericGroupLen + int(s)
	// If di < 0 the result must be 0, but if di = 0, the value loses all digits, but could round up to 1 if its first extra digit is >= 5.
	if decimalDigits < 0 {
		return nil, 0
	}
	baseDigits := (decimalDigits + numericGroupLen - 1) / numericGroupLen
	decimalDigits %= numericGroupLen // 0, or number of decimal digits to keep in last base digit
	if baseDigits > len(d) || (baseDigits == len(d) && decimalDigits == 0) {
		return d, w
	}

	var carry int16
	if decimalDigits == 0 {
		var extraBaseDigit int16
		if baseDigits < len(d) {
			extraBaseDigit = d[baseDigits]
			d = d[:baseDigits]
		}

		if extraBaseDigit >= numericBase/2 {
			carry = 1
		} else {
			carry = 0
		}
	} else {
		// Must round within last base digit
		d = d[:baseDigits]
		var roundPowers = [4]int16{0, 1000, 100, 10}
		pow10 := roundPowers[decimalDigits]
		baseDigits--
		extra := d[baseDigits] % pow10
		d[baseDigits] -= extra
		carry = 0
		if extra >= pow10/2 {
			pow10 += d[baseDigits]
			if pow10 >= numericBase {
				pow10 -= numericBase
				carry = 1
			}
			d[baseDigits] = pow10
		}
	}

	// Propagate carry if needed
	for baseDigits--; carry > 0 && baseDigits >= 0; baseDigits-- {
		carry += d[baseDigits]
		if carry >= numericBase {
			d[baseDigits] = carry - numericBase
			carry = 1
		} else {
			d[baseDigits] = carry
			carry = 0
		}
	}
	if carry > 0 {
		d = append([]int16{carry}, d...)
		w++
	}

	return d, w
}

// truncAbs truncates the value of d at s decimal digits after the decimal point.
// s<0 means truncation before the decimal point.
// truncAbs is based on PostgreSQL trunc_var function defined at "src/backend/utils/adt/numeric.c".
func truncAbs(d []int16, w int16, rscale int16) ([]int16, int16) {
	if len(d) == 0 {
		return nil, 0
	}

	decimalDigits := (int(w)+1)*numericGroupLen + int(rscale)
	// If di <= 0, the value loses all digits.
	if decimalDigits <= 0 {
		return nil, 0
	}
	baseDigits := (decimalDigits + numericGroupLen - 1) / numericGroupLen

	if baseDigits <= len(d) {
		d = d[:baseDigits]

		// 0, or number of decimal digits to keep in last base digit
		decimalDigits %= numericGroupLen

		if decimalDigits > 0 {
			// Must truncate within last base digit
			var roundPowers = [4]int16{0, 1000, 100, 10}
			pow10 := roundPowers[decimalDigits]
			// Warning, here was groupLen specific code
			baseDigits--
			extra := d[baseDigits] % pow10
			d[baseDigits] -= extra
		}
	}

	return d, w
}

// selectDivScaleAbs calculates default scale for division (as PostgreSQL do it).
// selectDivScaleAbs is based on PostgreSQL select_div_scale function defined at "src/backend/utils/adt/numeric.c".
func selectDivScaleAbs(d1 []int16, w1 int16, d2 []int16, w2 int16) int16 {
	// The result scale of a division isn't specified in any SQL standard.
	// For PostgreSQL we select a result scale that will give at least NUMERIC_MIN_SIG_DIGITS significant digits,
	// so that numeric gives a result no less accurate than float8; but use a scale not less than either input's display scale.

	// Get the actual (normalized) weight and first digit of each input
	firstDigit1 := int16(0)
	for i := 0; i < len(d1); i++ {
		if d1[i] != 0 {
			firstDigit1 = d1[i]
			w1 -= int16(i)
			break
		}
	}
	if firstDigit1 == 0 {
		w1 = 0
	}

	firstDigit2 := int16(0)
	for i := 0; i < len(d2); i++ {
		if d2[i] != 0 {
			firstDigit2 = d2[i]
			w2 -= int16(i)
			break
		}
	}
	if firstDigit2 == 0 {
		w2 = 0
	}

	// Estimate weight of quotient.  If the two first digits are equal, we can't be sure, but assume that var1 is less than var2.
	qweight := w1 - w2
	if firstDigit1 <= firstDigit2 {
		qweight--
	}

	// Select result scale
	rscale := pgNumericMinSigDigits - qweight*numericGroupLen
	rscale = mathh.Max2Int16(rscale, getScaleAbs(d1, w1)) // Here used emulated scale of operand
	rscale = mathh.Max2Int16(rscale, getScaleAbs(d2, w2)) // Here used emulated scale of operand
	rscale = mathh.Max2Int16(rscale, pgNumericMinDisplayScale)
	rscale = mathh.Min2Int16(rscale, pgNumericMaxDisplayScale)

	return rscale
}

// getScaleAbs simulates PostgreSQL Numeric field dscale.
// It returns number of decimal digit after decimal point.
// Result is not exactly, because this function treats each base digit as a base number of decimal digits ("0.9000"="0.9" but getScaleAbs returns 4).
// Result value as always >= 0 as required be PostgreSQL field.
// This function used then calculation result depends on operands scale (div, sqrt, ...) and in communication with DB.
// Most of arithmetic operations does not require scale because they are absolutely accurate.
func getScaleAbs(d []int16, w int16) int16 {
	s := (int16(len(d)) - w - 1) * numericGroupLen
	if s <= 0 {
		return 0
	}
	return s
}

// Quo is just a shorthand for QuoPrec with default scale and rounding enabled.
// Default scale calculates as in PostgreSQL (more or less).
func (z *Numeric) Quo(x, y *Numeric) *Numeric {
	return z.QuoPrec(x, y, selectDivScaleAbs(x.digits, x.weight, y.digits, y.weight), true)
}

// QuoPrec sets z to the quotient x/y for y != 0 and returns z.
// Result will truncated or rounded up to scale decimal digits after decimal point.
// Result will be rounded if round is true, otherwise result will be truncated.
// If y == 0, a division-by-zero run-time panic occurs.
// QuoPrec implements truncated division (like Go); see QuoRem for more details.
func (z *Numeric) QuoPrec(x, y *Numeric, scale int16, round bool) *Numeric {
	if x.IsNaN() || y.IsNaN() {
		return z.SetNaN()
	}
	if y.IsZero() {
		panic("division by zero")
	}
	if x.IsZero() {
		return z.SetZero()
	}

	if x.sign == y.sign {
		z.sign = numericPositive
	} else {
		z.sign = numericNegative
	}

	z.digits, z.weight = divAbs(x.digits, x.weight, y.digits, y.weight, scale, round)

	return z
}

// QuoRem sets z to the quotient x/y and r to the remainder x%y and returns the pair (z, r) for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
//
// QuoRem implements T-division and modulus (like Go):
//
//	q = x/y      with the result truncated to zero
//	r = x - y*q
//
// (See Daan Leijen, ``Division and Modulus for Computer Scientists''.)
// Euclidean division and modulus (unlike Go) do not currently implemented for Numeric.
func (z *Numeric) QuoRem(x, y, m *Numeric) (*Numeric, *Numeric) {
	z.QuoPrec(x, y, 0, false)
	m.Sub(x, m.Mul(z, y))
	return z, m
}

// Rem sets z to the remainder x%y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Rem implements truncated modulus (like Go); see QuoRem for more details.
func (z *Numeric) Rem(x, y *Numeric) *Numeric {
	z.QuoRem(x, y, z)
	return z
}
