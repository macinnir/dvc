package time

import (
	"strconv"
	"strings"
	"time"
)

type ITime interface {
	NowMilliseconds() int64
}

// https://www.golangprograms.com/add-n-number-of-year-month-day-hour-minute-second-millisecond-microsecond-and-nanosecond-to-current-date-time.html

// NowMilliseconds returns a unix timestamp (int64) in milliseconds
func NowMilliseconds() int64 {
	return ToMilliseconds(time.Now())
}

// ToMilliseconds converts a time to milliseconds
func ToMilliseconds(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// NowAddDateMilliseconds adds year, month, day to a date and returns the milliseconds
func NowAddDateMilliseconds(year int, month int, day int) int64 {
	return ToMilliseconds(time.Now().AddDate(year, month, day))
}

// ParseRFC3339Offset parses an RFC3339 and returns its offset in seconds
func ParseRFC3339Offset(dateString string) (offsetSeconds int64, e error) {

	offsetString := dateString[len(dateString)-6:]
	offsetParts := strings.Split(offsetString, ":")

	var offsetHour int64 = 0
	var offsetMinute int64 = 0

	offsetPlusMinus := offsetParts[0][0]

	if offsetHour, e = strconv.ParseInt(offsetParts[0][1:], 10, 64); e != nil {
		return
	}

	if offsetMinute, e = strconv.ParseInt(offsetParts[1], 10, 64); e != nil {
		return
	}

	offsetSeconds = (offsetHour * 3600) + (offsetMinute * 60)

	if offsetPlusMinus == '-' {
		offsetSeconds = 0 - offsetSeconds
	}

	return
}

// type Timer struct {
// 	t *time.Time
// }

// func NewTimer() *Timer {

// 	now := time.Now()

// 	return &Timer{
// 		t: &now,
// 	}
// }

// func (t *Timer) AddYear(year int) *Timer {
// 	t.t.AddDate(year, 0, 0)
// 	return t
// }

// func (t *Timer) AddMonth(month int) *Timer {
// 	t.t.AddDate(0, month, 0)
// 	return t
// }

// func (t *Timer) AddDay(day int) *Timer {
// 	t.t.AddDate(0, 0, day)
// 	return t
// }

// func (t *Timer) AddHour(hour int) *Timer {
// 	t.t.Add(time.Duration(hour) * time.Hour)
// 	return t
// }

// func (t *Timer) AddMinute(minute int) *Timer {
// 	t.t.Add(time.Duration(minute) * time.Minute)
// 	return t
// }

// func (t *Timer) AddSecond(second int) *Timer {
// 	t.t.Add(time.Duration(second) * time.Second)
// 	return t
// }
