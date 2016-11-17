package numeric

import "github.com/apaxa-io/mathhelper"

func (z *Numeric) SetInt8(x int8) *Numeric{
	if x==0{
		return z.SetZero()
	}

	if x<0 {
		z.sign=Negative
	}else {
		z.sign=Positive
	}
	z.digits=[]int16{int16(mathhelper.AbsInt8(x))}

	return z
}