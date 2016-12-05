package pgtypes

import (
	"github.com/apaxa-go/helper/mathh"
	"github.com/apaxa-go/helper/timeh"
	"math"
	"testing"
	"time"
)

// TODO - add more direct tests
// TODO - make tests for the following:
/*
When adding an interval value to (or subtracting an interval value from) a timestamp with time zone value, the days component advances or decrements the date of the timestamp with time zone by the indicated number of days. Across daylight saving time changes (when the session time zone is set to a time zone that recognizes DST), this means interval '1 day' does not necessarily equal interval '24 hours'. For example, with the session time zone set to CST7CDT, timestamp with time zone '2005-04-02 12:00-07' + interval '1 day' will produce timestamp with time zone '2005-04-03 12:00-06', while adding interval '24 hours' to the same initial timestamp with time zone produces timestamp with time zone '2005-04-03 13:00-06', as there is a change in daylight saving time at 2005-04-03 02:00 in time zone CST7CDT.

Note there can be ambiguity in the months field returned by age because different months have different numbers of days. PostgreSQL's approach uses the month from the earlier of the two dates when calculating partial months. For example, age('2004-06-01', '2004-04-30') uses April to yield 1 mon 1 day, while using May would yield 1 mon 2 days because May has 31 days, while April has only 30.

Subtraction of dates and timestamps can also be complex. One conceptually simple way to perform subtraction is to convert each value to a number of seconds using EXTRACT(EPOCH FROM ...), then subtract the results; this produces the number of seconds between the two values. This will adjust for the number of days in each month, timezone changes, and daylight saving time adjustments. Subtraction of date or timestamp values with the "-" operator returns the number of days (24-hours) and hours/minutes/seconds between the values, making the same adjustments. The age function returns years, months, days, and hours/minutes/seconds, performing field-by-field subtraction and then adjusting for negative field values. The following queries illustrate the differences in these approaches. The sample results were produced with timezone = 'US/Eastern'; there is a daylight saving time change between the two dates used:

SELECT EXTRACT(EPOCH FROM timestamptz '2013-07-01 12:00:00') -
       EXTRACT(EPOCH FROM timestamptz '2013-03-01 12:00:00');
Result: 10537200
SELECT (EXTRACT(EPOCH FROM timestamptz '2013-07-01 12:00:00') -
        EXTRACT(EPOCH FROM timestamptz '2013-03-01 12:00:00'))
        / 60 / 60 / 24;
Result: 121.958333333333
SELECT timestamptz '2013-07-01 12:00:00' - timestamptz '2013-03-01 12:00:00';
Result: 121 days 23:00:00
SELECT age(timestamptz '2013-07-01 12:00:00', timestamptz '2013-03-01 12:00:00');
Result: 4 mons
*/

func TestParse(t *testing.T) {
	type testElement struct {
		s    string
		prec uint8
		i    Interval
		err  bool
	}

	test := []testElement{
		{"-1 year -2 mons +3 days -04:05:06", IntervalMicrosecondPrecision, Interval{-14, 3, -14706 * 1e9, IntervalNanosecondPrecision}, false},
		{"-1 year 2 mons -3 days 04:05:06.789", IntervalMicrosecondPrecision, Interval{-10, -3, 14706789 * 1e6, IntervalNanosecondPrecision}, false},
		{"", IntervalMicrosecondPrecision, Interval{0, 0, 0, IntervalNanosecondPrecision}, false},
		{"00:00", IntervalMicrosecondPrecision, Interval{}, true},
		{"year mons days", IntervalMicrosecondPrecision, Interval{}, true},
		{"1.5 year", IntervalMicrosecondPrecision, Interval{}, true},
		{"1,5 year", IntervalMicrosecondPrecision, Interval{}, true},
		{"99999999999 year -2 mons +3 days -04:05:06", IntervalMicrosecondPrecision, Interval{}, true},
		{"9 year 9999999999 mons +3 days -04:05:06", IntervalMicrosecondPrecision, Interval{}, true},
		{"9 year -2 mons +99999999999 days -04:05:06", IntervalMicrosecondPrecision, Interval{}, true},
		{"9 year -2 mons +9 days 040506", IntervalMicrosecondPrecision, Interval{}, true},
		{"9 year -2 mons +9 days 04:06:99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", IntervalMicrosecondPrecision, Interval{}, true},
		//TODO waiting check overflow
		//{"2147483647 year 2147483647 mons 2147483647 days 00:00:00", MicrosecondPrecision, Interval{2147483647, 2147483647, 0, MicrosecondPrecision}, false},
		//-2147483648 to 2147483647
		//TODO waiting fix spaces
		//{"   ", MicrosecondPrecision, Interval{}, true},
		{"1 year", 15, Interval{12, 0, 0, IntervalPicosecondPrecision}, false},
		{"132428754353245897987092345345:00:00", IntervalMicrosecondPrecision, Interval{}, true},
		{"00:132428754353245897987092345345:00", IntervalMicrosecondPrecision, Interval{}, true},
		{"00:00:132428754353245897987092345345", IntervalMicrosecondPrecision, Interval{}, true},
		{"00:00:00.132428754353245897987092345345", IntervalMillisecondPrecision, Interval{0, 0, 132, IntervalMillisecondPrecision}, false},
		{"00:00:00.132528754353245897987092345345", IntervalMillisecondPrecision, Interval{0, 0, 133, IntervalMillisecondPrecision}, false},
		{"00:00:00.132628754353245897987092345345", IntervalMillisecondPrecision, Interval{0, 0, 133, IntervalMillisecondPrecision}, false},
	}

	for _, v := range test {
		i, err := ParseInterval(v.s, v.prec)
		if (err != nil) != v.err {
			t.Errorf("%v,%v: expect error %v, got %v", v.s, v.prec, v.err, err)
		}
		if !v.err && (err == nil) {
			if !i.Equal(v.i) {
				t.Errorf("%v,%v: expect %v, got %v", v.s, v.prec, v.i, i)
			}
		}
	}
}

func TestString(t *testing.T) {
	type testElement struct {
		s   string
		i   Interval
		err bool
	}

	test := []testElement{
		// 0
		{
			s:   "-1 year -2 mons 3 days -04:05:06",
			i:   Interval{-14, 3, -14706 * 1e9, IntervalNanosecondPrecision},
			err: false,
		},

		// 1
		{
			s:   "-10 mons -3 days 04:05:06.789",
			i:   Interval{-10, -3, 14706789 * 1e6, IntervalNanosecondPrecision},
			err: false,
		},

		// 2
		{
			s:   "1 mons",
			i:   Interval{1, 0, 0, IntervalNanosecondPrecision},
			err: false,
		},

		// 3
		{
			s:   "2 year -34:57:18",
			i:   Interval{24, 0, -125838 * 1e9, IntervalNanosecondPrecision},
			err: false,
		},

		// 4
		{
			s:   "00:00:00",
			i:   Interval{0, 0, 0, IntervalNanosecondPrecision},
			err: false,
		},

		// 5
		{
			s:   "83 year 4 mons",
			i:   Interval{1000, 0, 0, IntervalNanosecondPrecision},
			err: false,
		},

		// 6
		{
			s:   "1000 days",
			i:   Interval{0, 1000, 0, IntervalNanosecondPrecision},
			err: false,
		},

		// 7
		{
			s:   "-1 mons",
			i:   Interval{-1, 0, 0, IntervalNanosecondPrecision},
			err: false,
		},

		// 8
		{
			s:   "-1 mons",
			i:   Interval{-1, 0, 0, IntervalNanosecondPrecision},
			err: false,
		},

		//-2147483648 to 2147483647
		// 9
		{
			s:   "178956970 year 7 mons 2147483647 days",
			i:   Interval{2147483647, 2147483647, 0, IntervalNanosecondPrecision},
			err: false,
		},

		// 10
		{
			s:   "-178956970 year -8 mons -2147483648 days",
			i:   Interval{-2147483648, -2147483648, 0, IntervalNanosecondPrecision},
			err: false,
		},
	}

	for j, v := range test {
		s := v.i.String()
		if s != v.s {
			t.Errorf("Test-%v. Strings not equal.\nExpected:\n%s\ngot:\n%s", j, v.s, s)
		}
	}
}

func TestAdd(t *testing.T) {
	type testElement struct {
		i   Interval
		add Interval
		res Interval
	}

	test := []testElement{
		// 0
		{
			Interval{-14, 3, -14706 * 1e9, IntervalNanosecondPrecision},
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			Interval{-13, 5, -14703 * 1e9, IntervalNanosecondPrecision},
		},

		// 1
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{0, 0, 0, IntervalNanosecondPrecision},
		},

		// 2
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{-14, 3, -14706 * 1e9, IntervalNanosecondPrecision},
			Interval{-14, 3, -14706 * 1e9, IntervalNanosecondPrecision},
		},

		// 3
		{
			Interval{-14, -15, -16 * 1e9, IntervalNanosecondPrecision},
			Interval{-14, -15, -16 * 1e9, IntervalNanosecondPrecision},
			Interval{-28, -30, -32 * 1e9, IntervalNanosecondPrecision},
		},

		// 4
		{
			Interval{-14, -15, -16 * 1e9, IntervalNanosecondPrecision},
			Interval{14, 15, 16 * 1e9, IntervalNanosecondPrecision},
			Interval{0, 0, 0, IntervalNanosecondPrecision},
		},

		// 5
		{
			Interval{14, 15, 16 * 1e9, IntervalNanosecondPrecision},
			Interval{100, 200, 300 * 1e9, IntervalNanosecondPrecision},
			Interval{114, 215, 316 * 1e9, IntervalNanosecondPrecision},
		},
		// 6
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{14, 15, 16 * 1e9, IntervalNanosecondPrecision},
			Interval{14, 15, 16 * 1e9, IntervalNanosecondPrecision},
		},

		// 7
		{
			Interval{14, 15, 16 * 1e9, IntervalNanosecondPrecision},
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{14, 15, 16 * 1e9, IntervalNanosecondPrecision},
		},
	}

	for j, v := range test {
		i := v.i.Add(v.add)
		if !i.Equal(v.res) {
			t.Errorf("Test-%v. Intervals are not equal.\nExpected:\n%v\ngot:\n%v", j, v.res, i)
		}
	}
}

func TestDuration(t *testing.T) {

	type testElement struct {
		i            Interval
		daysInMonth  uint8
		minutesInDay uint32
		d            time.Duration
	}

	test := []testElement{
		// 0
		{
			Interval{0, 0, 86400 * 1e9, IntervalNanosecondPrecision},
			30,
			1440,
			86400 * time.Second,
		},

		// 1
		{
			Interval{0, 10, 1 * 1e9, IntervalNanosecondPrecision},
			30,
			1440,
			864001 * time.Second,
		},

		// 2
		{
			Interval{10, 10, 1 * 1e9, IntervalNanosecondPrecision},
			30,
			1440,
			26784001 * time.Second, //2562000
		},

		// 3
		{
			Interval{20, 10, 1 * 1e9, IntervalNanosecondPrecision},
			0,
			0,
			time.Second,
		},

		// 4
		{
			Interval{-10, -5, -1 * 1e9, IntervalNanosecondPrecision},
			30,
			1400,
			-25620001 * time.Second,
		},

		// 5
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			30,
			1400,
			0,
		},
	}

	for j, v := range test {
		d := v.i.Duration(v.daysInMonth, v.minutesInDay)
		if d != v.d {
			t.Errorf("Test-%v. Wrong duration. Expected: %v, got: %v", j, v.d, d)
		}
	}
}

func TestSub(t *testing.T) {
	type testElement struct {
		i   Interval
		sub Interval
		res Interval
	}

	test := []testElement{
		// 0
		{
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			Interval{2, 3, 4 * 1e9, IntervalNanosecondPrecision},
			Interval{-1, -1, -1 * 1e9, IntervalNanosecondPrecision},
		},

		// 1
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{2, 3, 4 * 1e9, IntervalNanosecondPrecision},
			Interval{-2, -3, -4 * 1e9, IntervalNanosecondPrecision},
		},

		// 2
		{
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
		},

		// 3
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			Interval{0, 0, 0, IntervalNanosecondPrecision},
		},

		// 4
		{
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			Interval{0, 0, 0, IntervalNanosecondPrecision},
		},

		// 5
		{
			Interval{-1, -2, -3 * 1e9, IntervalNanosecondPrecision},
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			Interval{-2, -4, -6 * 1e9, IntervalNanosecondPrecision},
		},

		// 6
		{
			Interval{-2147483648, -2147483648, -3 * 1e9, IntervalNanosecondPrecision},
			Interval{-1, -2, -3 * 1e9, IntervalNanosecondPrecision},
			Interval{-2147483647, -2147483646, 0, IntervalNanosecondPrecision},
		},

		// 7
		{
			Interval{2147483647, 2147483647, -3 * 1e9, IntervalNanosecondPrecision},
			Interval{1, 2, -3 * 1e9, IntervalNanosecondPrecision},
			Interval{2147483646, 2147483645, 0, IntervalNanosecondPrecision},
		},
	}

	for j, v := range test {
		s := v.i.Sub(v.sub)
		if !s.Equal(v.res) {
			t.Errorf("Test-%v. Wrong sub.\nExpected interval:%v\ngot:%v", j, v.res, s)
		}
	}
}

func TestMul(t *testing.T) {
	type testElement struct {
		i   Interval
		mul int64
		res Interval
	}

	test := []testElement{
		// 0
		{
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			2,
			Interval{2, 4, 6 * 1e9, IntervalNanosecondPrecision},
		},

		// 1
		{
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			-2,
			Interval{-2, -4, -6 * 1e9, IntervalNanosecondPrecision},
		},

		// 2
		// TODO no more valid since moved from float64 to int64
		//{
		//	Interval{1, 2, 3 * 1e9, NanosecondPrecision},
		//	1.05,
		//	Interval{1, 2, 3150 * 1e6, NanosecondPrecision},
		//},

		// 3
		{
			Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision},
			0,
			Interval{0, 0, 0, IntervalNanosecondPrecision},
		},

		// 4
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			-2,
			Interval{0, 0, 0, IntervalNanosecondPrecision},
		},
	}

	for j, v := range test {
		i := v.i.Mul(v.mul)
		//if (i.Months != v.res.Months) || (i.Days != v.res.Days) || (math.Abs(i.Seconds-v.res.Seconds) > inaccuracySeconds) {
		if (i.Months != v.res.Months) || (i.Days != v.res.Days) || (i.SomeSeconds != v.res.SomeSeconds) {
			t.Errorf("Test-%v. Wrong interval.\nExpected:%v\ngot:%v", j, v.res, i)
		}
	}
}

func TestDiv(t *testing.T) {
	type testElement struct {
		i   Interval
		div int64
		res Interval
	}

	test := []testElement{
		// 0
		{
			Interval{4, 6, 8 * 1e9, IntervalNanosecondPrecision},
			2,
			Interval{2, 3, 4 * 1e9, IntervalNanosecondPrecision},
		},

		// 1
		// TODO no more valid since moved from float64 to int64
		//{
		//	Interval{4, 6, 8 * 1e9, NanosecondPrecision},
		//	1.1,
		//	Interval{3, 5, 7272727272, NanosecondPrecision},
		//},
	}

	for j, v := range test {
		i := v.i.Div(v.div)
		//if (i.Months != v.res.Months) || (i.Days != v.res.Days) || (math.Abs(i.Seconds-v.res.Seconds) > inaccuracySeconds) {
		if (i.Months != v.res.Months) || (i.Days != v.res.Days) || (i.SomeSeconds != v.res.SomeSeconds) {
			t.Errorf("Test-%v. Wrong interval.\nExpected:%v\ngot:%v", j, v.res, i)
		}
	}
}

func TestInterval_Cmp(t *testing.T) {
	type testElement struct {
		i          Interval
		i2         Interval
		comparable bool
		sign       int
	}
	tests := []testElement{
		{Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision}, Interval{2, 3, 4 * 1e9, IntervalNanosecondPrecision}, true, -1},
		{Interval{10, 2, 3 * 1e9, IntervalNanosecondPrecision}, Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision}, true, 1},
		{Interval{1, 20, 3 * 1e9, IntervalNanosecondPrecision}, Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision}, true, 1},
		{Interval{10, 20, 30 * 1e9, IntervalNanosecondPrecision}, Interval{2, 3, 4 * 1e9, IntervalNanosecondPrecision}, true, 1},
		{Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision}, Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision}, true, 0},
		{Interval{1, 2, 3 * 1e9, IntervalNanosecondPrecision}, Interval{1, 2, 3 * 1e6, IntervalMicrosecondPrecision}, true, 0},
		{Interval{1, 0, 86400 * 1e9, IntervalNanosecondPrecision}, Interval{1, 1, 0, IntervalNanosecondPrecision}, false, 0},
		{Interval{1, 0, 86400 * 1e3, IntervalMillisecondPrecision}, Interval{1, 1, 0, IntervalNanosecondPrecision}, false, 0},
		{Interval{-2147483648, -2147483648, 0, IntervalNanosecondPrecision}, Interval{2147483647, -2147483647, 0, IntervalNanosecondPrecision}, true, -1},
	}

	// Extend tests: Cmp(a,b)=-Cmp(b,a)
	var tests2 []testElement
	for _, v := range tests {
		v.sign *= -1
		v.i, v.i2 = v.i2, v.i
		tests2 = append(tests2, v)
	}
	tests = append(tests, tests2...)

	// Run
	for _, v := range tests {
		var less, lessEq, greater, greaterEq, eq bool
		if v.comparable {
			less = v.sign == -1
			lessEq = v.sign <= 0
			greater = v.sign == 1
			greaterEq = v.sign >= 0
			eq = v.sign == 0
		}

		if s, ok := v.i.Cmp(v.i2); s != v.sign || ok != v.comparable {
			t.Errorf("%v,%v: expect %v %v, got %v %v", v.i, v.i2, v.sign, v.comparable, s, ok)
		}

		if r := v.i.Comparable(v.i2); r != v.comparable {
			t.Errorf("%v,%v: expect %v, got %v", v.i, v.i2, v.comparable, r)
		}
		if r := v.i.Less(v.i2); r != less {
			t.Errorf("%v,%v: expect %v, got %v", v.i, v.i2, less, r)
		}
		if r := v.i.LessOrEqual(v.i2); r != lessEq {
			t.Errorf("%v,%v: expect %v, got %v", v.i, v.i2, lessEq, r)
		}
		if r := v.i.Greater(v.i2); r != greater {
			t.Errorf("%v,%v: expect %v, got %v", v.i, v.i2, greater, r)
		}
		if r := v.i.GreaterOrEqual(v.i2); r != greaterEq {
			t.Errorf("%v,%v: expect %v, got %v", v.i, v.i2, greaterEq, r)
		}
		if r := v.i.Equal(v.i2); r != eq {
			t.Errorf("%v,%v: expect %v, got %v", v.i, v.i2, eq, r)
		}
	}
}

func TestAddToAndSubFrom(t *testing.T) {
	const inaccuracySeconds = 5

	type testElement struct {
		i   Interval
		t   time.Time
		res time.Time
	}

	test := []testElement{
		// 0
		{
			Interval{0, 0, 1 * 1e9, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(1, 0),
		},

		// 1
		{
			Interval{0, 0, 1 * 1e9, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(1, 0),
		},

		// 2
		{
			Interval{0, 0, 1 * 1e9, IntervalNanosecondPrecision},
			time.Unix(86400, 0),
			time.Unix(86401, 0),
		},

		// 3
		{
			Interval{0, 0, 1 * 1e9, IntervalNanosecondPrecision},
			time.Unix(0, 9223372035854775807),
			time.Unix(0, 9223372036854775807),
		},

		// 4
		{
			Interval{0, 0, 9223372036854775807, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(9223372036, 854775807),
		},

		// 5
		{
			Interval{0, 0, -4775808, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(0, -4775808),
		},
	}

	for j, v := range test {
		tA := v.i.AddTo(v.t)
		if time.Duration(math.Abs(float64(tA.Sub(v.res)))) > inaccuracySeconds*time.Second {
			//if t1 != v.res {
			t.Errorf("TestAddTo - %v. Wrong time\nExpected time:\n%v\ngot:\n%v", j, v.res.UTC(), tA.UTC())
		}
		tS := v.i.SubFrom(v.res)
		if time.Duration(math.Abs(float64(tS.Sub(v.res)))) > inaccuracySeconds*time.Second {
			t.Errorf("TestSubFrom - %v. Wrong time\nExpected time:\n%v\ngot:\n%v", j, v.t.UTC(), tS.UTC())
		}
	}

}

func TestNormal(t *testing.T) {
	type testElement struct {
		i    Interval
		year int32
		mon  int32
		day  int32
		hour int64
		min  int64
		sec  int64
		nsec int64
	}

	test := []testElement{
		// 0
		{
			Interval{1001, 101, 10013 * 1e8, IntervalNanosecondPrecision},
			83,
			5,
			101,
			0,
			16,
			41,
			3 * 1e8,
		},

		// 1
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		},

		// 2
		{
			Interval{-128, 97, 24001789 * 1e6, IntervalNanosecondPrecision},
			-10,
			-8,
			97,
			6,
			40,
			1,
			789000000,
		},
	}

	for j, v := range test {
		y := v.i.NormalYears()
		if y != v.year {
			t.Errorf("Test-%v. Ecpected normal year: %v, got: %v", j, v.year, y)
		}
		m := v.i.NormalMonths()
		if m != v.mon {
			t.Errorf("Test-%v. Ecpected normal month: %v, got: %v", j, v.mon, m)
		}
		d := v.i.NormalDays()
		if d != v.day {
			t.Errorf("Test-%v. Ecpected normal days: %v, got: %v", j, v.day, d)
		}
		h := v.i.NormalHours()
		if h != v.hour {
			t.Errorf("Test-%v. Ecpected normal hours: %v, got: %v", j, v.hour, h)
		}
		min := v.i.NormalMinutes()
		if min != v.min {
			t.Errorf("Test-%v. Ecpected normal minutes: %v, got: %v", j, v.min, min)
		}
		s := v.i.NormalSeconds()
		if s != v.sec {
			t.Errorf("Test-%v. Ecpected normal seconds: %v, got: %v", j, v.sec, s)
		}
		ns := v.i.NormalNanoseconds()
		if ns != v.nsec {
			t.Errorf("Test-%v. Ecpected normal nanoseconds: %v, got: %v\n%v", j, v.nsec, ns, s)
		}
	}
}

func TestAll(t *testing.T) {
	i := Picosecond()
	if i.SomeSeconds != 1 || i.precision != IntervalPicosecondPrecision {
		t.Error("Error")
	}

	i = Nanosecond()
	if i.SomeSeconds != 1 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}

	i = Microsecond()
	if i.SomeSeconds != 1e3 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}

	i = Millisecond()
	if i.SomeSeconds != 1e6 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}

	i = Second()
	if i.SomeSeconds != 1e9 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}

	i = Minute()
	if i.SomeSeconds != 60*1e9 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}

	i = Hour()
	if i.SomeSeconds != 3600*1e9 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}

	i = Day()
	if i.Days != 1 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}

	i = Month()
	if i.Months != 1 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}

	i = Year()
	if i.Months != 12 || i.precision != IntervalNanosecondPrecision {
		t.Error("Error")
	}
}

func TestFromDuration(t *testing.T) {
	type testElement struct {
		i Interval
		d time.Duration
	}
	test := []testElement{
		// 0
		{
			Interval{0, 0, 86400 * 1e9, IntervalNanosecondPrecision},
			86400 * time.Second,
		},

		// 1
		{
			Interval{0, 0, 8 * 1e9, IntervalNanosecondPrecision},
			8 * time.Second,
		},

		// 2
		{
			Interval{0, 0, -9223372036854775808, IntervalNanosecondPrecision},
			-9223372036854775808,
		},

		// 3
		{
			Interval{0, 0, 9223372036854775807, IntervalNanosecondPrecision},
			9223372036854775807,
		},

		// 4
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			0,
		},

		// 5
		{
			Interval{0, 0, -0000000001, IntervalNanosecondPrecision},
			-1,
		},
	}
	for j, v := range test {
		i := FromDuration(v.d)
		if i != v.i {
			t.Errorf("Test-%v. Wrong interval. Expected: %v, got: %v", j, v.i, i)
		}
	}
}

func TestDiff(t *testing.T) {
	type testElement struct {
		i    Interval
		from time.Time
		to   time.Time
	}
	test := []testElement{
		// 0
		{
			Interval{0, 0, 1 * 1e9, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(1, 0),
		},

		// 1
		{
			Interval{0, 0, -1 * 1e9, IntervalNanosecondPrecision},
			time.Unix(1, 0),
			time.Unix(0, 0),
		},

		// 2
		{
			Interval{0, 0, -0000000001, IntervalNanosecondPrecision},
			time.Unix(0, 1),
			time.Unix(0, 0),
		},

		// 3
		{
			Interval{0, 0, 0000000001, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(0, 1),
		},

		// 4
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(0, 0),
		},

		// 5
		{
			Interval{0, 0, 5854775807, IntervalNanosecondPrecision},
			time.Unix(1, 0),
			time.Unix(0, 6854775807),
		},

		// 6
		{
			Interval{0, 0, -9223372036854775807, IntervalNanosecondPrecision},
			time.Unix(0, 9223372036854775807),
			time.Unix(0, 0),
		},

		// 7
		{
			Interval{0, 0, 9223372036854775807, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(9223372036, 854775807),
		},

		// 8
		{
			Interval{0, 0, -9223372036854775808, IntervalNanosecondPrecision},
			time.Unix(0, 0),
			time.Unix(0, -9223372036854775808),
		},
	}

	for j, v := range test {
		i := Diff(v.from, v.to)
		if i != v.i {
			t.Errorf("Test-%v. Wrong interval\nExpected:\n%v\ngot:\n%v", j, v.i, i)
		}
	}
}

func TestDiffExtended(t *testing.T) {
	type testElement struct {
		i     Interval
		sFrom string
		sTo   string
	}
	test := []testElement{
		// 0
		{
			Interval{3507, 10, 85636854775807, IntervalNanosecondPrecision},
			"1970-01-01T00:00:00Z",
			"2262-04-11T23:47:16.854775807Z",
		},

		// 1
		{
			Interval{0, 0, 0, IntervalNanosecondPrecision},
			"1970-01-01T00:00:00Z",
			"1970-01-01T00:00:00Z",
		},

		// 2
		{
			Interval{0, 0, 1 * 1e9, IntervalNanosecondPrecision},
			"1970-01-01T00:00:58Z",
			"1970-01-01T00:00:59Z",
		},

		// 3
		{
			Interval{12, 0, 1260 * 1e9, IntervalNanosecondPrecision},
			"1970-01-01T00:11:00Z",
			"1971-01-01T00:32:00Z",
		},

		// 4
		{
			Interval{12, 0, 3600 * 1e9, IntervalNanosecondPrecision},
			"1970-01-01T22:00:00Z",
			"1971-01-01T23:00:00Z",
		},

		// 5
		{
			Interval{0, 11, 0, IntervalNanosecondPrecision},
			"1970-01-14T00:00:00Z",
			"1970-01-25T00:00:00Z",
		},

		// 6
		{
			Interval{8, 11, 0, IntervalNanosecondPrecision},
			"1970-03-14T00:00:00Z",
			"1970-11-25T00:00:00Z",
		},

		// 7
		{
			Interval{12, 0, 0, IntervalNanosecondPrecision},
			"1970-01-01T00:00:00Z",
			"1971-01-01T00:00:00Z",
		},

		// 8
		{
			Interval{852, 0, 0, IntervalNanosecondPrecision},
			"1900-01-01T00:00:00Z",
			"1971-01-01T00:00:00Z",
		},

		// 9
		{
			Interval{-3507, -10, -85636854775807, IntervalNanosecondPrecision},
			"2262-04-11T23:47:16.854775807Z",
			"1970-01-01T00:00:00Z",
		},

		// 10
		{
			Interval{-1192, 11, 0, IntervalNanosecondPrecision},
			"2000-03-01T00:00:00Z",
			"1900-11-12T00:00:00Z",
		},
	}

	for j, v := range test {
		// RFC3339Nano = "2006-01-02T15:04:05.999 999 999Z07:00"
		timeFrom, err := time.Parse(time.RFC3339Nano, v.sFrom)
		if err != nil {
			t.Errorf("Test-%v. Parsing string:%v\ngot err: %v", j, v.sFrom, err)
		}
		timeTo, err1 := time.Parse(time.RFC3339Nano, v.sTo)
		if err1 != nil {
			t.Errorf("Test-%v. Got err: %v, while parsing:%v", j, v.sTo, err1)
		}
		i := DiffExtended(timeFrom, timeTo)
		if i != v.i {
			t.Errorf("Test-%v. Wrong interval\nExpected:\n%v\ngot:\n%v", j, v.i, i)
		}
	}

}

func TestSince(t *testing.T) {
	const inaccuracySeconds = 1
	test := []time.Time{time.Unix(1, 0), time.Unix(1e9, 1e18), time.Unix(0, 0)}
	//TODO check whats wrong with big values
	// max time: time.Unix(1<<63-62135596801, 999999999)
	//time.Unix(- 9223372036854775808, -9223372036854775808)
	for j, v := range test {
		nsec := time.Since(v)
		i := Since(v)
		if (i.Months != 0) || (i.Days != 0) || mathh.AbsInt64(i.SomeSeconds-int64(nsec)) > inaccuracySeconds*timeh.NanosecsInSec {
			t.Errorf("Test-%v. Wrong time since: %v\nExpected (time.Since):\n%v\ngot (Since):\n%v", j, v, nsec, time.Duration(i.SomeSeconds))
		}
	}

}

func TestSinceExtended(t *testing.T) {
	const inaccuracySeconds = 5
	test := []time.Time{time.Unix(1, 0), time.Unix(1e9, 1e18), time.Unix(0, 0)}
	for j, v := range test {
		i := SinceExtended(v)
		v1 := v.AddDate(0, int(i.Months), int(i.Days))
		v1 = v1.Add(time.Duration(i.SomeSeconds) * time.Nanosecond)
		//if time.Since(v1) > inaccuracySeconds*time.Second || time.Since(v1) < -inaccuracySeconds*time.Second {
		if time.Since(v1) > inaccuracySeconds*time.Second || time.Since(v1) < -inaccuracySeconds*time.Second {
			t.Errorf("Test-%v\nWrong time since: %v\nGit interval:%v\ntime now(v1):%v\nexpected time since(ts):%v", j, v, i, v1, time.Since(v1))
		}

	}
}

func TestNewInterval(t *testing.T) {
	if i := NewInterval(IntervalMicrosecondPrecision); i.Days != 0 || i.Months != 0 || i.SomeSeconds != 0 || i.precision != IntervalMicrosecondPrecision {
		t.Errorf("expect %v %v, got %v", 0, IntervalMicrosecondPrecision, i)
	}
	if i := NewInterval(15); i.Days != 0 || i.Months != 0 || i.SomeSeconds != 0 || i.precision != IntervalMaxPrecision {
		t.Errorf("expect %v %v, got %v", 0, IntervalMaxPrecision, i)
	}
}

func TestNewGoInterval(t *testing.T) {
	if i := NewGoInterval(); i.Days != 0 || i.Months != 0 || i.SomeSeconds != 0 || i.precision != IntervalGoPrecision {
		t.Errorf("expect %v %v, got %v", 0, IntervalGoPrecision, i)
	}
}

func TestNewPgInterval(t *testing.T) {
	if i := NewPgInterval(); i.Days != 0 || i.Months != 0 || i.SomeSeconds != 0 || i.precision != IntervalPgPrecision {
		t.Errorf("expect %v %v, got %v", 0, IntervalPgPrecision, i)
	}
}

func TestInterval_SetPrecision(t *testing.T) {
	type testElement struct {
		i Interval
		p uint8
		r Interval
	}
	tests := []testElement{
		{Interval{1, 2, 1, IntervalSecondPrecision}, IntervalMillisecondPrecision, Interval{1, 2, 1000, IntervalMillisecondPrecision}},
		{Interval{1, 2, 1, IntervalMillisecondPrecision}, IntervalMillisecondPrecision, Interval{1, 2, 1, IntervalMillisecondPrecision}},
		{Interval{1, 2, 1499, IntervalMillisecondPrecision}, IntervalSecondPrecision, Interval{1, 2, 1, IntervalSecondPrecision}},
		{Interval{1, 2, 1500, IntervalMillisecondPrecision}, IntervalSecondPrecision, Interval{1, 2, 2, IntervalSecondPrecision}},
		{Interval{1, 2, 1, IntervalNanosecondPrecision}, 15, Interval{1, 2, 1000, IntervalPicosecondPrecision}},
	}
	for _, v := range tests {
		if r := v.i.SetPrecision(v.p); r != v.r {
			t.Errorf("%v,%v: expect %v, got %v", v.i, v.p, v.r, r)
		}
	}
}

func TestInterval_Precision(t *testing.T) {
	i := Interval{0, 0, 0, IntervalSecondPrecision}
	if i.Precision() != IntervalSecondPrecision {
		t.Error("error")
	}
	i = Interval{0, 0, 0, IntervalMillisecondPrecision}
	if i.Precision() != IntervalMillisecondPrecision {
		t.Error("error")
	}
}

func TestInterval_SafePrec(t *testing.T) {
	type testElement struct {
		i Interval
		p uint8
	}
	// 	10 seconds => 0 (second precision, can not be less)
	//	10 milliseconds => 2 (2 digit after decimal point precision)
	// 	123 microseconds => 6 (microsecond precision)
	tests := []testElement{
		{Interval{1, 2, 10, IntervalSecondPrecision}, IntervalSecondPrecision},
		{Interval{1, 2, 10, IntervalMillisecondPrecision}, 2},
		{Interval{1, 2, 123, IntervalMicrosecondPrecision}, IntervalMicrosecondPrecision},
	}
	for _, v := range tests {
		if r := v.i.SafePrec(); r != v.p {
			t.Errorf("%v: expect %v, got %v", v.i, v.p, r)
		}
	}
}

func TestInterval_In(t *testing.T) {
	type testElement struct {
		i1, i2 Interval
		r      int64
	}
	tests := []testElement{
		{Interval{0, 0, 1, IntervalGoPrecision}, Interval{0, 0, 10, IntervalGoPrecision}, 10},
		{Interval{0, 0, 1, IntervalMillisecondPrecision}, Interval{0, 0, 1, IntervalSecondPrecision}, 1000},
		{Interval{0, 1, 0, IntervalGoPrecision}, Interval{1, 0, 0, IntervalSecondPrecision}, 30},
		{Interval{0, 0, 1, IntervalSecondPrecision}, Interval{0, 1, 0, IntervalGoPrecision}, 3600 * 24},
	}
	for _, v := range tests {
		if r := v.i1.In(v.i2); r != v.r {
			t.Errorf("%v,%v: expect %v, got %v", v.i1, v.i2, v.r, r)
		}
	}
}
