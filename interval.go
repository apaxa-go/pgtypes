package pgtypes

import (
	"errors"
	"github.com/apaxa-go/helper/mathh"
	"github.com/apaxa-go/helper/strconvh"
	"github.com/apaxa-go/helper/stringsh"
	"github.com/apaxa-go/helper/timeh"
	"regexp"
	"strings"
	"time"
)

// Predefined precisions for Interval.
const (
	IntervalSecondPrecision      = 0
	IntervalMillisecondPrecision = 3
	IntervalMicrosecondPrecision = 6
	IntervalNanosecondPrecision  = 9
	IntervalPicosecondPrecision  = 12
	IntervalGoPrecision          = IntervalNanosecondPrecision
	IntervalPgPrecision          = IntervalMicrosecondPrecision
	IntervalMaxPrecision         = 12
)

// RE for parse interval in postgres style specification.
// http://www.postgresql.org/docs/9.4/interactive/datatype-datetime.html#DATATYPE-INTERVAL-OUTPUT
var re = regexp.MustCompile(`^(?:([+-]?[0-9]+) years?)? ?(?:([+-]?[0-9]+) mons?)? ?(?:([+-]?[0-9]+) days?)? ?(?:([+-])?([0-9]+):([0-9]+):([0-9]+)(?:,|.([0-9]+))?)?$`)

// Interval represent time interval in Postgres-compatible way.
// It consists of 3 public fields:
// 	Months - number months
// 	Days - number of days
// 	SomeSeconds - number of seconds or some smaller units (depends on precision).
// All fields are signed. Sign of one field is independent from sign of others.
// Interval internally stores precision. Precision is number of digits after comma in 10-based representation of seconds.
// Precision can be from [0; 12] where 0 means that SomeSeconds is seconds and 12 means that SomeSeconds is picoseconds.
// If Interval created without calling constructor when it has 0 precision (i.e. SomeSeconds is just seconds).
// If Interval created with calling constructor and its documentation does not say another when it has precision = 9 (i.e. SomeSeconds is nanoseconds). This is because default Go time type has nanosecond precision.
// If interval is used to store PostgreSQL Interval when recommended precision is 6 (microsecond) because PostgreSQL use microsecond precision.
// This type is similar to Postgres interval data type.
// Value from one field is never automatically translated to value of another field, so <60*60*24 seconds> != <1 days> and so on.
// This is because of compatibility with Postgres, moreover day may have different amount of seconds and month may have different amount of days.
type Interval struct {
	Months      int32
	Days        int32
	SomeSeconds int64
	precision   uint8
}

// Picosecond returns new Interval equal to 1 picosecond.
// This constructor return interval with precision = 12 (picosecond).
func Picosecond() Interval {
	return Interval{SomeSeconds: 1, precision: IntervalPicosecondPrecision}
}

// Nanosecond returns new Interval equal to 1 nanosecond.
func Nanosecond() Interval {
	return Interval{SomeSeconds: 1, precision: IntervalGoPrecision}
}

// Microsecond returns new Interval equal to 1 microsecond.
func Microsecond() Interval {
	return Interval{SomeSeconds: timeh.NanosecsInMicrosec, precision: IntervalGoPrecision}
}

// Millisecond returns new Interval equal to 1 millisecond.
func Millisecond() Interval {
	return Interval{SomeSeconds: timeh.NanosecsInMillisec, precision: IntervalGoPrecision}
}

// Second returns new Interval equal to 1 second.
func Second() Interval {
	return Interval{SomeSeconds: timeh.NanosecsInSec, precision: IntervalGoPrecision}
}

// Minute returns new Interval equal to 1 minute (60 seconds).
func Minute() Interval {
	return Interval{SomeSeconds: timeh.NanosecsInSec * timeh.SecsInMin, precision: IntervalGoPrecision}
}

// Hour returns new Interval equal to 1 hour (3600 seconds).
func Hour() Interval {
	return Interval{SomeSeconds: timeh.NanosecsInSec * timeh.SecsInHour, precision: IntervalGoPrecision}
}

// Day returns new Interval equal to 1 day.
func Day() Interval {
	return Interval{Days: 1, precision: IntervalGoPrecision}
}

// Month returns new Interval equal to 1 month.
func Month() Interval {
	return Interval{Months: 1, precision: IntervalGoPrecision}
}

// Year returns new Interval equal to 1 year (12 months).
func Year() Interval {
	return Interval{Months: timeh.MonthsInYear, precision: IntervalGoPrecision}
}

// ParseInterval parses incoming string and extract interval with requested precision p.
// Format is postgres style specification for interval output format.
// Examples:
// 	-1 year 2 mons -3 days 04:05:06.789
// 	1 mons
// 	2 year -34:56:78
// 	00:00:00
//
// BUG(unsacrificed): ParseInterval may overflow SomeSeconds if computed SomeSeconds should be MinInt64.
func ParseInterval(s string, p uint8) (i Interval, err error) {
	if p > IntervalMaxPrecision {
		i.precision = IntervalMaxPrecision
	} else {
		i.precision = p
	}

	parts := re.FindStringSubmatch(s)
	if parts == nil || len(parts) != 9 {
		err = errors.New("Unable to parse interval from string " + s)
		return
	}

	var ti32 int32

	// Store as months:

	// years
	if parts[1] != "" {
		ti32, err = strconvh.ParseInt32(parts[1])
		if err != nil {
			return
		}
		i.Months = ti32 * timeh.MonthsInYear
	}

	// months
	if parts[2] != "" {
		ti32, err = strconvh.ParseInt32(parts[2])
		if err != nil {
			return
		}
		i.Months += ti32
	}

	// Store as days:

	// days
	if parts[3] != "" {
		ti32, err = strconvh.ParseInt32(parts[3])
		if err != nil {
			return
		}
		i.Days = ti32
	}

	var ti64 int64

	// Store as seconds:

	negativeTime := parts[4] == "-"

	// hours
	if parts[5] != "" {
		ti64, err = strconvh.ParseInt64(parts[5])
		if err != nil {
			return
		}
		i.SomeSeconds = ti64 // Now SomeSeconds contains hours
	}
	i.SomeSeconds *= timeh.MinsInHour // Now SomeSeconds contains minutes

	// minutes
	if parts[6] != "" {
		ti64, err = strconvh.ParseInt64(parts[6])
		if err != nil {
			return
		}
		i.SomeSeconds += ti64
	}
	i.SomeSeconds *= timeh.SecsInMin // Now SomeSeconds contains seconds

	// seconds
	if parts[7] != "" {
		ti64, err = strconvh.ParseInt64(parts[7]) // Possible overflow
		if err != nil {
			return
		}
		i.SomeSeconds += ti64
	}

	i.SomeSeconds *= mathh.PowInt64(10, int64(p)) // Now SomeSeconds contains units with required precision

	if parts[8] != "" {
		if len(parts[8]) < int(p) {
			parts[8] = stringsh.PadRightWithByte(parts[8], '0', int(p))
		}
		ti64, err = strconvh.ParseInt64(parts[8][:p]) // Possible overflow
		if err != nil {
			return // It is impossible to cover this case because of RegExp and limits on p (p<Digits(MaxInt64))
		}
		i.SomeSeconds += ti64

		if len(parts[8]) > int(p) && parts[8][int(p)] >= '5' { // Round-to-upper if needed
			i.SomeSeconds++
		}
	}

	if negativeTime {
		i.SomeSeconds *= -1
	}

	return
}

// FromDuration returns new Interval equivalent for given time.Duration (convert time.Duration to Interval).
func FromDuration(d time.Duration) Interval {
	return Interval{SomeSeconds: d.Nanoseconds(), precision: IntervalGoPrecision}
}

// Diff calculates difference between given timestamps (time.Time) as nanoseconds and returns result as Interval (=to-from).
// Result always have months & days parts set to zero.
func Diff(from, to time.Time) Interval {
	return Interval{SomeSeconds: to.UnixNano() - from.UnixNano(), precision: IntervalGoPrecision}
}

// DiffExtended is similar to Diff but calculates difference in months, days & nanoseconds instead of just nanoseconds (=to-from).
// Result may have non-zero months & days parts.
// DiffExtended use Location of both passed times while calculation. Most of time it is better to pass times with the same Location (UTC or not).
func DiffExtended(from, to time.Time) (i Interval) {
	fromYear, fromMonth, fromDay := from.Date()
	toYear, toMonth, toDay := to.Date()

	i.Months = int32((toYear-fromYear)*timeh.MonthsInYear + int(toMonth-fromMonth))
	i.Days = int32(toDay - fromDay)

	i.SomeSeconds = to.UnixNano() - i.AddTo(from).UnixNano()
	i.precision = IntervalGoPrecision

	return
}

// Since returns elapsed time since given timestamp as Interval (=Diff(t, time.New()).
// Result always have months & days parts set to zero.
func Since(t time.Time) Interval {
	return Diff(t, time.Now())
}

// SinceExtended returns elapsed time since given timestamp as Interval (=DiffExtended(t, time.New()).
// Result may have non-zero months & days parts.
func SinceExtended(t time.Time) Interval {
	return DiffExtended(t, time.Now().In(t.Location()))
}

// NewInterval returns zero interval with specified precision p.
func NewInterval(p uint8) Interval {
	if p > IntervalMaxPrecision {
		p = IntervalMaxPrecision
	}
	return Interval{precision: p}
}

// NewGoInterval returns zero interval with GoLang precision (= nanosecond).
func NewGoInterval() Interval {
	return Interval{0, 0, 0, IntervalGoPrecision}
}

// NewPgInterval returns zero interval with PostgreSQL precision (= microsecond).
func NewPgInterval() Interval {
	return Interval{0, 0, 0, IntervalPgPrecision}
}

// SetPrecision returns new interval with changed precision (and do appropriate recalculation).
// Possible precision is 0..12 where 0 means second precision and 9 means nanosecond precision.
// If passed p>12 it will be silently replaced with p=12.
func (i Interval) SetPrecision(p uint8) Interval {
	if p > IntervalMaxPrecision {
		p = IntervalMaxPrecision
	}
	if p == i.precision {
		return i
	}
	return Interval{Months: i.Months, Days: i.Days, SomeSeconds: someSecondsChangePrecision(i.SomeSeconds, i.precision, p), precision: p}
}

// Precision returns internally stored precision.
func (i Interval) Precision() uint8 {
	return i.precision
}

// String returns string representation of interval.
// Output format is the same as for Parse.
//
// BUG(unsacrificed): String may overflow SomeSeconds if SomeSeconds is MinInt64.
func (i Interval) String() string {
	if i.Months == 0 && i.Days == 0 && i.SomeSeconds == 0 {
		return "00:00:00"
	}

	y := i.NormalYears()
	mon := i.NormalMonths()

	str := ""
	if y != 0 {
		str += strconvh.FormatInt32(y) + " year "
	}
	if mon != 0 {
		str += strconvh.FormatInt32(mon) + " mons "
	}
	if i.Days != 0 {
		str += strconvh.FormatInt32(i.Days) + " days "
	}

	if i.SomeSeconds != 0 {
		negativeTime := i.SomeSeconds < 0
		if negativeTime {
			i.SomeSeconds *= -1 // Possible overflow
		}

		tmp := mathh.PowInt64(10, int64(i.precision))
		h := i.NormalHours()
		m := i.NormalMinutes()
		f := i.SomeSeconds % (tmp * timeh.SecsInMin)
		s := f / tmp
		f -= s * tmp

		if negativeTime {
			str += "-"
		}

		str += stringsh.PadLeftWithByte(strconvh.FormatInt64(h), '0', 2) + ":" +
			stringsh.PadLeftWithByte(strconvh.FormatInt64(m), '0', 2) + ":" +
			stringsh.PadLeftWithByte(strconvh.FormatInt64(s), '0', 2)
		if f != 0 {
			str += "." + strings.TrimRight(
				stringsh.PadLeftWithByte(strconvh.FormatInt64(f), '0', int(i.precision)),
				"0",
			)
		}

		return str
	}
	// As all null interval filtered at the beginning of method there is a space at the end of string
	return str[:len(str)-1]
}

// Duration convert Interval to time.Duration.
// It is required to pass number of days in month (usually 30 or something near)
// and number of minutes in day (usually 1440) because of converting months and days parts of original Interval to time.Duration nanoseconds.
// Warning: this method is inaccuracy because in real life daysInMonth & minutesInDay vary and depends on relative timestamp.
func (i Interval) Duration(daysInMonth uint8, minutesInDay uint32) time.Duration {
	return time.Duration((int64(i.Months)*int64(daysInMonth)+int64(i.Days))*int64(minutesInDay)*mathh.PowInt64(10, int64(i.precision))*timeh.SecsInMin + i.SomeSeconds)
}

// someSecondsChangePrecision recalculates s (with precision from) to precision to and return result.
// Example: someSecondsChangePrecision(20 000 000, 6, 3)=20 000
//	it can be read as convert 1 000 000 microseconds (1e-6 seconds) to milliseconds (1e-3).
func someSecondsChangePrecision(s int64, from, to uint8) int64 {
	if to >= from {
		return s * mathh.PowInt64(10, int64(to-from))
	}
	return mathh.DivideRoundFixInt64(s, mathh.PowInt64(10, int64(from-to)))
}

// Add returns i+add.
func (i Interval) Add(add Interval) Interval {
	i.Months += add.Months
	i.Days += add.Days
	i.SomeSeconds += someSecondsChangePrecision(add.SomeSeconds, add.precision, i.precision)
	return i
}

// Sub returns i-sub.
func (i Interval) Sub(sub Interval) Interval {
	i.Months -= sub.Months
	i.Days -= sub.Days
	i.SomeSeconds -= someSecondsChangePrecision(sub.SomeSeconds, sub.precision, i.precision)
	return i
}

// Mul returns interval i multiplied by mul. Each part of Interval multiples independently.
func (i Interval) Mul(mul int64) Interval {
	i.Months, i.Days, i.SomeSeconds = int32(int64(i.Months)*mul), int32(int64(i.Days)*mul), int64(int64(i.SomeSeconds)*mul)
	return i
}

// Div divides interval by mul and returns result. Each part of Interval divides independently.
// Round rule: 0.4=>0 ; 0.5=>1 ; 0.6=>1 ; -0.4=>0 ; -0.5=>-1 ; -0.6=>-1
func (i Interval) Div(div int64) Interval {
	i.Months = int32(mathh.DivideRoundFixInt64(int64(i.Months), div))
	i.Days = int32(mathh.DivideRoundFixInt64(int64(i.Days), div))
	i.SomeSeconds = mathh.DivideRoundFixInt64(i.SomeSeconds, div)
	return i
}

// SafePrec returns minimal precision which can be used for current Interval without data loss.
// Examples:
// 	10 seconds => 0 (second precision, can not be less)
//	10 milliseconds => 2 (2 digit after decimal point precision)
// 	123 microseconds => 6 (microsecond precision)
func (i Interval) SafePrec() uint8 {
	ss := i.SomeSeconds
	p := i.precision
	for p > 0 && ss%10 == 0 {
		ss /= 10
		p--
	}
	return p
}

// In counts how many i contains in i2 (=i2/i).
// Round rule: 0.4=>0 ; 0.5=>1 ; 0.6=>1 ; -0.4=>0 ; -0.5=>-1 ; -0.6=>-1
func (i Interval) In(i2 Interval) int64 {
	prec := mathh.Max2Uint8(i.SafePrec(), i2.SafePrec())
	pow := mathh.PowInt64(10, int64(prec))

	iv := (int64(i.Months)*timeh.DaysInMonth + int64(i.Days)) * timeh.SecsInDay
	iv *= pow
	iv += someSecondsChangePrecision(i.SomeSeconds, i.precision, prec)

	i2v := (int64(i2.Months)*timeh.DaysInMonth + int64(i2.Days)) * timeh.SecsInDay
	i2v *= pow
	i2v += someSecondsChangePrecision(i2.SomeSeconds, i2.precision, prec)

	return mathh.DivideRoundFixInt64(i2v, iv)
}

// Cmp returns compare two Intervals. ok indicates is i and i2 comparable.
// sign means the following:
// 	sign<0 => i<i2
// 	sign=0 => i=i2
// 	sign>0 => i>i2
func (i Interval) Cmp(i2 Interval) (sign int, ok bool) {
	var mSign, dSign, sSign int

	switch {
	case i.Months < i2.Months:
		mSign = -1
	case i.Months > i2.Months:
		mSign = 1
	}

	switch {
	case i.Days < i2.Days:
		dSign = -1
	case i.Days > i2.Days:
		dSign = 1
	}

	// Compare SomeSecond. Requirements: smallP<=bigP
	ssCmp := func(v1 int64, smallP uint8, v2 int64, bigP uint8) int {
		tmp := mathh.PowInt64(10, int64(bigP-smallP))
		switch {
		case v1 < v2/tmp:
			return -1
		case v1 > v2/tmp:
			return 1
		default:
			return int(mathh.SignInt64(v2 % tmp))
		}
	}
	if i.precision <= i2.precision {
		sSign = ssCmp(i.SomeSeconds, i.precision, i2.SomeSeconds, i2.precision)
	} else {
		sSign = -ssCmp(i2.SomeSeconds, i2.precision, i.SomeSeconds, i.precision)
	}

	// Check if Intervals comparable
	if (mSign < 0 || dSign < 0 || sSign < 0) && (mSign > 0 || dSign > 0 || sSign > 0) {
		return 0, false
	}

	sign = mSign + dSign + sSign

	switch {
	case sign < -1:
		sign = -1
	case sign > 1:
		sign = 1
	}

	return sign, true
}

// Comparable returns true only if it is possible to compare Intervals.
// Intervals "A" and "B" can be compared only if:
//   1) all parts of "A" are less or equal to relative parts of "B"
//   or
//   2) all parts of "B" are less or equal to relative parts of "A".
// In the other words, it is impossible to compare "30 days"-Interval with "1 month"-Interval.
func (i Interval) Comparable(i2 Interval) bool {
	_, ok := i.Cmp(i2)
	return ok
}

// Equal compare original Interval with given for full equality part by part.
func (i Interval) Equal(i2 Interval) bool {
	sign, ok := i.Cmp(i2)
	return ok && sign == 0
}

// LessOrEqual returns true if all parts of original Interval are less or equal to relative parts of i2.
func (i Interval) LessOrEqual(i2 Interval) bool {
	sign, ok := i.Cmp(i2)
	return ok && sign <= 0
}

// Less returns true if at least one part of original Interval is less then relative part of i2 and all other parts of original Interval are less or equal to relative parts of i2.
func (i Interval) Less(i2 Interval) bool {
	sign, ok := i.Cmp(i2)
	return ok && sign < 0
}

// GreaterOrEqual returns true if all parts of original Interval are greater or equal to relative parts of i2.
func (i Interval) GreaterOrEqual(i2 Interval) bool {
	sign, ok := i.Cmp(i2)
	return ok && sign >= 0
}

// Greater returns true if at least one part of original Interval is greater then relative part of i2 and all other parts of original Interval are greater or equal to relative parts of i2.
func (i Interval) Greater(i2 Interval) bool {
	sign, ok := i.Cmp(i2)
	return ok && sign > 0
}

// NormalYears return number of years in month part (as i.Months / 12).
func (i Interval) NormalYears() int32 {
	return i.Months / timeh.MonthsInYear
}

// NormalMonths return number of months in month part after subtracting NormalYears*12 (as i.Months % 12).
// Examples: if .Months = 11 then NormalMonths = 11, but if .Months = 13 then NormalMonths = 1.
func (i Interval) NormalMonths() int32 {
	return i.Months % timeh.MonthsInYear
}

// NormalDays just returns Days part.
func (i Interval) NormalDays() int32 {
	return i.Days
}

// NormalHours returns number of hours in seconds part.
func (i Interval) NormalHours() int64 {
	pow := mathh.PowInt64(10, int64(i.precision))
	return i.SomeSeconds / (pow * timeh.SecsInHour)
}

// NormalMinutes returns number of hours in seconds part after subtracting NormalHours.
func (i Interval) NormalMinutes() int64 {
	pow := mathh.PowInt64(10, int64(i.precision))
	return (i.SomeSeconds - i.NormalHours()*pow*timeh.SecsInHour) / (pow * timeh.SecsInMin)
}

// NormalSeconds returns number of seconds in seconds part after subtracting NormalHours*3600 and NormalMinutes*60 (as i.Seconds % 60).
func (i Interval) NormalSeconds() int64 {
	pow := mathh.PowInt64(10, int64(i.precision))
	return (i.SomeSeconds / pow) % timeh.SecsInMin
}

// NormalNanoseconds returns number of nanoseconds in fraction part of seconds part.
func (i Interval) NormalNanoseconds() int64 {
	pow := mathh.PowInt64(10, int64(i.precision))
	return i.SomeSeconds % pow
}

// AddTo adds original Interval to given timestamp and return result.
func (i Interval) AddTo(t time.Time) time.Time {
	return t.AddDate(0, int(i.Months), int(i.Days)).Add(time.Duration(someSecondsChangePrecision(i.SomeSeconds, i.precision, IntervalNanosecondPrecision)))
}

// SubFrom subtract original Interval from given timestamp and return result.
func (i Interval) SubFrom(t time.Time) time.Time {
	return i.Mul(-1).AddTo(t)
}
