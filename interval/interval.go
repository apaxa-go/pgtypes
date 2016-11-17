package interval

import (
	"errors"
	"github.com/apaxa-io/mathhelper"
	"github.com/apaxa-io/strconvhelper"
	"github.com/apaxa-io/stringshelper"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TODO move to parent package?
// TODO implement NullInterval

const (
	SecondPrecision      = 0
	MillisecondPrecision = 3
	MicrosecondPrecision = 6
	NanosecondPrecision  = 9
	PicosecondPrecision  = 12
	GoPrecision          = NanosecondPrecision
	PostgreSQLPrecision  = MicrosecondPrecision
	maxPrecision         = 12
	defaultPrecision     = GoPrecision
)

// RE for parse interval in postgres style specification.
// http://www.postgresql.org/docs/9.4/interactive/datatype-datetime.html#DATATYPE-INTERVAL-OUTPUT
var re = regexp.MustCompile(`^(?:([+-]?[0-9]+) year)? ?(?:([+-]?[0-9]+) mons)? ?(?:([+-]?[0-9]+) days)? ?(?:([+-])?([0-9]+):([0-9]+):([0-9]+)(?:,|.([0-9]+))?)?$`)

// Interval represent time interval in Postgres-compatible way.
// It consists of 3 public fields:
// 	Months - number months
// 	Days - number of days
// 	SomeSeconds - number of seconds or some smaller units (depends on precision).
// All fields are signed. Sign of one field is independent from sign of other field.
// Interval internally stores precision. Precision is number of digits after comma in 10-based representation of seconds.
// Precision can be from [0; 12] where 0 means that SomeSeconds is seconds and 12 means that SomeSeconds is picoseconds.
// If Interval created without calling constructor when it has 0 precision (i.e. SomeSeconds is just seconds).
// If Interval created with calling constructor and its documentation does not say another when it has precision = 9 (i.e. SomeSeconds is nanoseconds). This is because default Go time type has nanosecond precision.
// If interval is used to store PostgreSQL Interval when recommended precision is 6 (microsecond) because PostgreSQL use microsecond.
// This type is similar to Postgres interval data type.
// Value from one field is never automatically translated to value of another field, so <60*60*24 seconds> != <1 days> and so on.
// This is because of:
// 	1) compatibility with Postgres;
// 	2) day may have different amount of seconds and month may have different amount of days.
type Interval struct {
	Months      int32
	Days        int32
	SomeSeconds int64
	precision   uint8
}

// Nanosecond returns new Interval equal to 1 picosecond.
// This constructor return interval with precision = 12 (picosecond).
func Picosecond() Interval {
	return Interval{SomeSeconds: 1, precision: PicosecondPrecision}
}

// Nanosecond returns new Interval equal to 1 nanosecond
func Nanosecond() Interval {
	return Interval{SomeSeconds: 1, precision: GoPrecision}
}

// Microsecond returns new Interval equal to 1 microsecond
func Microsecond() Interval {
	return Interval{SomeSeconds: NanosecsInMicrosec, precision: GoPrecision}
}

// Millisecond returns new Interval equal to 1 millisecond
func Millisecond() Interval {
	return Interval{SomeSeconds: NanosecsInMillisec, precision: GoPrecision}
}

// Second returns new Interval equal to 1 second
func Second() Interval {
	return Interval{SomeSeconds: NanosecsInSec, precision: GoPrecision}
}

// Minute returns new Interval equal to 1 minute (60 seconds)
func Minute() Interval {
	return Interval{SomeSeconds: NanosecsInSec * SecsInMin, precision: GoPrecision}
}

// Hour returns new Interval equal to 1 hour (3600 seconds)
func Hour() Interval {
	return Interval{SomeSeconds: NanosecsInSec * SecsInHour, precision: GoPrecision}
}

// Day returns new Interval equal to 1 day
func Day() Interval {
	return Interval{Days: 1, precision: GoPrecision}
}

// Month returns new Interval equal to 1 month
func Month() Interval {
	return Interval{Months: 1, precision: GoPrecision}
}

// Year returns new Interval equal to 1 year (12 months)
func Year() Interval {
	return Interval{Months: MonthsInYear, precision: GoPrecision}
}

// Parse parses incoming string and extract interval with requested precision p.
// Format is postgres style specification for interval output format.
// Examples:
// 	-1 year 2 mons -3 days 04:05:06.789
// 	1 mons
// 	2 year -34:56:78
// 	00:00:00
func Parse(s string, p uint8) (i Interval, err error) {
	//TODO string of 1-3 spaces are parse ok
	//TODO add check for overflow

	if p > maxPrecision {
		i.precision = maxPrecision
	} else {
		i.precision = p
	}

	parts := re.FindStringSubmatch(s)
	if parts == nil || len(parts) != 9 {
		err = errors.New("Unable to parse interval from string " + s)
		return
	}

	var ti int64

	// Store as months:

	// years
	if parts[1] != "" {
		ti, err = strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return
		}
		i.Months = int32(ti) * MonthsInYear
	}

	// months
	if parts[2] != "" {
		ti, err = strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			return
		}
		i.Months += int32(ti)
	}

	// Store as days:

	// days
	if parts[3] != "" {
		ti, err = strconv.ParseInt(parts[3], 10, 32)
		if err != nil {
			return
		}
		i.Days = int32(ti)
	}

	// Store as seconds:

	negativeTime := parts[4] == "-" // TODO problem with MinInt64 because scanning as positive

	// hours
	if parts[5] != "" {
		ti, err = strconv.ParseInt(parts[5], 10, 64)
		if err != nil {
			return
		}
		i.SomeSeconds = ti // Now somesecs contains hours
	}
	i.SomeSeconds *= MinsInHour // Now somesecs contains minutes

	// minutes
	if parts[6] != "" {
		ti, err = strconv.ParseInt(parts[6], 10, 64)
		if err != nil {
			return
		}
		i.SomeSeconds += ti
	}
	i.SomeSeconds *= SecsInMin // Now somesecs contains seconds

	// seconds
	if parts[7] != "" {
		ti, err = strconv.ParseInt(parts[7], 10, 64)
		if err != nil {
			return
		}
		i.SomeSeconds += ti
	}

	i.SomeSeconds *= mathhelper.PowInt64(10, int64(p)) // Now someseconds contains required precision units

	if parts[8] != "" {
		if len(parts[8]) < int(p) {
			parts[8] = stringshelper.PadRightWithByte(parts[8], '0', int(p))
		}
		ti, err = strconv.ParseInt(parts[8][:p], 10, 64)
		if err != nil {
			return
		}
		i.SomeSeconds += ti

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
	return Interval{SomeSeconds: d.Nanoseconds(), precision: GoPrecision}
}

// Diff calculates difference between given timestamps (time.Time) as nanoseconds and returns result as Interval (=to-from).
// Result always have months & days parts set to zero.
func Diff(from, to time.Time) Interval {
	return Interval{SomeSeconds: to.UnixNano() - from.UnixNano(), precision: GoPrecision}
}

// DiffExtended is similar to Diff but calculates difference in months, days & nanoseconds instead of just nanoseconds (=to-from).
// Result may have non-zero months & days parts.
// DiffExtended use Location of both passed times while calculation. Most of time it is better to pass times with the same Location (UTC or not).
func DiffExtended(from, to time.Time) (i Interval) {
	fromYear, fromMonth, fromDay := from.Date()
	toYear, toMonth, toDay := to.Date()

	i.Months = int32((toYear-fromYear)*MonthsInYear + int(toMonth-fromMonth))
	i.Days = int32(toDay - fromDay)

	i.SomeSeconds = to.UnixNano() - i.AddTo(from).UnixNano()
	i.precision = GoPrecision

	return
}

// Since returns elapsed time since given timestamp as Interval (=Diff(t, time.New())
// Result always have months & days parts set to zero.
func Since(t time.Time) Interval {
	return Diff(t, time.Now())
}

// SinceExtended returns elapsed time since given timestamp as Interval (=DiffExtended(t, time.New())
// Result may have non-zero months & days parts.
func SinceExtended(t time.Time) Interval {
	return DiffExtended(t, time.Now().In(t.Location()))
}

// New returns zero interval with specified precision p
func NewInterval(p uint8) Interval {
	if p > maxPrecision {
		p = maxPrecision
	}
	return Interval{precision: p}
}

// New returns zero interval with GoLang precision (= nanosecond)
func NewGoInterval() Interval {
	return Interval{0, 0, 0, GoPrecision}
}

// New returns zero interval with PostgreSQL precision (= microsecond)
func NewPgInterval() Interval {
	return Interval{0, 0, 0, PostgreSQLPrecision}
}

// SetPrecision change interval precision and do appropriate stored value recalculation.
// Possible precision is 0..12 where 0 means second precision and 9 means nanosecond precision.
// If passed p>12 it will be silently replaced with p=12.
func (i Interval) SetPrecision(p uint8) Interval {
	if p > maxPrecision {
		p = maxPrecision
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
// Output format is the same as for Parse
func (i Interval) String() string {
	if i.Months == 0 && i.Days == 0 && i.SomeSeconds == 0 {
		return "00:00:00"
	}

	y := i.NormalYears()
	mon := i.NormalMonths()

	str := ""
	if y != 0 {
		str += strconvhelper.FormatInt32(y) + " year "
	}
	if mon != 0 {
		str += strconvhelper.FormatInt32(mon) + " mons "
	}
	if i.Days != 0 {
		str += strconvhelper.FormatInt32(i.Days) + " days "
	}

	if i.SomeSeconds != 0 {
		negativeTime := i.SomeSeconds < 0
		if negativeTime { // TODO possible overflow because of MinInt64*-1
			i.SomeSeconds *= -1
		}

		tmp := mathhelper.PowInt64(10, int64(i.precision))
		h := i.NormalHours()
		m := i.NormalMinutes()
		f := i.SomeSeconds % (tmp * SecsInMin)
		s := f / tmp
		f -= s * tmp

		if negativeTime {
			str += "-"
		}

		str += stringshelper.PadLeftWithByte(strconvhelper.FormatInt64(h), '0', 2) + ":" +
			stringshelper.PadLeftWithByte(strconvhelper.FormatInt8(m), '0', 2) + ":" +
			stringshelper.PadLeftWithByte(strconvhelper.FormatInt64(s), '0', 2)
		if f != 0 {
			str += "." + strings.TrimRight(
				stringshelper.PadLeftWithByte(strconvhelper.FormatInt64(f), '0', int(i.precision)),
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
	return time.Duration((int64(i.Months)*int64(daysInMonth)+int64(i.Days))*int64(minutesInDay)*mathhelper.PowInt64(10, int64(i.precision))*SecsInMin + i.SomeSeconds)
}

// someSecondsChangePrecision recalculates s (with precision from) to precision to and return result.
// Example: someSecondsChangePrecision(20 000 000, 6, 3)=20 000
//	it can be read as convert 1 000 000 microseconds (1e-6 seconds) to milliseconds (1e-3).
func someSecondsChangePrecision(s int64, from, to uint8) int64 {
	if to >= from {
		return s * mathhelper.PowInt64(10, int64(to-from))
	}
	return mathhelper.DivideRoundFixInt64(s, mathhelper.PowInt64(10, int64(from-to)))
}

// Add adds given Interval to original Interval.
// Original Interval will be changed.
// TODO 'will be changed'?
func (i Interval) Add(add Interval) Interval {
	i.Months += add.Months
	i.Days += add.Days
	i.SomeSeconds += someSecondsChangePrecision(add.SomeSeconds, add.precision, i.precision)
	return i
}

// Sub subtracts given Interval from original Interval.
// Original Interval will be changed.
// TODO 'will be changed'?
func (i Interval) Sub(sub Interval) Interval {
	i.Months -= sub.Months
	i.Days -= sub.Days
	i.SomeSeconds -= someSecondsChangePrecision(sub.SomeSeconds, sub.precision, i.precision)
	return i
}

// Mul multiples interval by mul. Each part of Interval multiples independently.
// Original Interval will be changed.
// TODO 'will be changed'?
func (i Interval) Mul(mul int64) Interval {
	i.Months, i.Days, i.SomeSeconds = int32(int64(i.Months)*mul), int32(int64(i.Days)*mul), int64(int64(i.SomeSeconds)*mul)
	return i
}

// Div divides interval by mul. Each part of Interval divides independently.
// Round rule: 0.4=>0 ; 0.5=>1 ; 0.6=>1 ; -0.4=>0 ; -0.5=>-1 ; -0.6=>-1
// Original Interval will be changed.
// TODO 'will be changed'?
func (i Interval) Div(div int64) Interval {
	i.Months = int32(mathhelper.DivideRoundFixInt64(int64(i.Months), div))
	i.Days = int32(mathhelper.DivideRoundFixInt64(int64(i.Days), div))
	i.SomeSeconds = mathhelper.DivideRoundFixInt64(i.SomeSeconds, div)
	return i
}

// In counts how many i contains in i2 (=i2/i).
// Round rule: 0.4=>0 ; 0.5=>1 ; 0.6=>1 ; -0.4=>0 ; -0.5=>-1 ; -0.6=>-1
// TODO A lot of overflows
func (i Interval) In(i2 Interval) int64 {
	iv := (int64(i.Months)*DaysInMonth+int64(i.Days))*SecsInDay*mathhelper.PowInt64(10, int64(i.precision)) + i.SomeSeconds
	i2v := (int64(i2.Months)*DaysInMonth+int64(i2.Days))*SecsInDay*mathhelper.PowInt64(10, int64(i2.precision)) + i2.SomeSeconds
	if i.precision > i2.precision {
		i2v = someSecondsChangePrecision(i2v, i2.precision, i.precision)
	} else {
		iv = someSecondsChangePrecision(iv, i.precision, i2.precision)
	}
	return mathhelper.DivideRoundFixInt64(i2v, iv)
}

// Comparable returns true only if it is possible to compare Intervals.
// Intervals "A" and "B" can be compared only if:
//   1) all parts of "A" are less or equal to relative parts of "B"
//   or
//   2) all parts of "B" are less or equal to relative parts of "A".
// In the other words, it is impossible to compare "30 days"-Interval with "1 month"-Interval.
func (i Interval) Comparable(i2 Interval) bool {
	return i.LessOrEqual(i2) || i.GreaterOrEqual(i2)
}

// Equal compare original Interval with given for full equality part by part.
func (i Interval) Equal(i2 Interval) bool {
	if i.Months != i2.Months || i.Days != i2.Days {
		return false
	}
	if (i.precision == i2.precision && i.SomeSeconds == i2.SomeSeconds) || (i.SomeSeconds == 0 && i2.SomeSeconds == 0) {
		return true
	}
	if i.precision > i2.precision {
		tmp := mathhelper.PowInt64(10, int64(i.precision-i2.precision))
		return i.SomeSeconds%tmp == 0 && i.SomeSeconds/tmp == i2.SomeSeconds
	} else {
		tmp := mathhelper.PowInt64(10, int64(i2.precision-i.precision))
		return i2.SomeSeconds%tmp == 0 && i2.SomeSeconds/tmp == i.SomeSeconds
	}
}

// LessOrEqual returns true if all parts of original Interval are less or equal to relative parts of i2.
func (i Interval) LessOrEqual(i2 Interval) bool {
	if i.Months > i2.Months || i.Days > i2.Days {
		return false
	}
	if (i.precision == i2.precision && i.SomeSeconds <= i2.SomeSeconds) || (i.SomeSeconds == 0 && i2.SomeSeconds == 0) {
		return true
	}
	if i.precision > i2.precision {
		tmp := mathhelper.PowInt64(10, int64(i.precision-i2.precision))
		switch {
		case i.SomeSeconds/tmp < i2.SomeSeconds:
			return true
		case i.SomeSeconds/tmp > i2.SomeSeconds:
			return false
		default:
			return i.SomeSeconds%tmp <= 0
		}
	} else {
		tmp := mathhelper.PowInt64(10, int64(i2.precision-i.precision))
		switch {
		case i.SomeSeconds < i2.SomeSeconds/tmp:
			return true
		case i.SomeSeconds > i2.SomeSeconds/tmp:
			return false
		default:
			return i2.SomeSeconds%tmp <= 0
		}
	}
}

// Less returns true if at least one part of original Interval is less then relative part of i2 and all other parts of original Interval are less or equal to relative parts of i2.
func (i Interval) Less(i2 Interval) bool {
	return !i.Equal(i2) && i.LessOrEqual(i2)
}

// GreaterOrEqual returns true if all parts of original Interval are greater or equal to relative parts of i2.
func (i Interval) GreaterOrEqual(i2 Interval) bool {
	//return i.Months >= i2.Months && i.Days >= i2.Days && i.SomeSeconds >= i2.SomeSeconds
	return i2.LessOrEqual(i)
}

// Greater returns true if at least one part of original Interval is greater then relative part of i2 and all other parts of original Interval are greater or equal to relative parts of i2.
func (i Interval) Greater(i2 Interval) bool {
	return !i.Equal(i2) && i.GreaterOrEqual(i2)
}

// NormalYears return number of years in month part (as i.Months / 12).
func (i Interval) NormalYears() int32 {
	// TODO what about sign?
	return i.Months / MonthsInYear
}

// NormalMonths return number of months in month part after subtracting NormalYears*12 (as i.Months % 12).
// Examples: if .Months = 11 then NormalMonths = 11, but if .Months = 13 then NormalMonths = 1.
func (i Interval) NormalMonths() int32 {
	// TODO what about sign?
	return i.Months % MonthsInYear
}

// NormalDays just returns Days part.
func (i Interval) NormalDays() int32 {
	return i.Days
}

// NormalHours returns number of hours in seconds part.
func (i Interval) NormalHours() int64 {
	// TODO what about sign?
	return int64(i.SomeSeconds / (mathhelper.PowInt64(10, int64(i.precision)) * SecsInHour))
}

// NormalMinutes returns number of hours in seconds part after subtracting NormalHours.
func (i Interval) NormalMinutes() int8 {
	// TODO what about sign?
	tmp := mathhelper.PowInt64(10, int64(i.precision))
	return int8((i.SomeSeconds - int64(i.NormalHours())*tmp*SecsInHour) / (tmp * SecsInMin))
}

// NormalSeconds returns number of seconds in seconds part after subtracting NormalHours*3600 and NormalMinutes*60 (as i.Seconds % 60).
func (i Interval) NormalSeconds() int8 {
	// TODO what about sign?
	return int8((i.SomeSeconds / mathhelper.PowInt64(10, int64(i.precision))) % SecsInMin)
}

// NormalNanoseconds returns number of nanoseconds in fraction part of seconds part.
func (i Interval) NormalNanoseconds() int32 {
	// TODO what about sign?
	return int32(i.SomeSeconds % mathhelper.PowInt64(10, int64(i.precision)))
}

// AddTo adds original Interval to given timestamp and return result.
func (i Interval) AddTo(t time.Time) time.Time {
	return t.AddDate(0, int(i.Months), int(i.Days)).Add(time.Duration(someSecondsChangePrecision(i.SomeSeconds, i.precision, NanosecondPrecision)))
}

// SubFrom subtract original Interval from given timestamp and return result.
func (i Interval) SubFrom(t time.Time) time.Time {
	return i.Mul(-1).AddTo(t) // TODO possible overflow (MinInt64)
}
