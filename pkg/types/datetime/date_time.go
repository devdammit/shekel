package datetime

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// DateTime type represents a time with date with minute precision.
// swaggen: type=string
// swaggen: format=date-time
type DateTime struct {
	time.Time
}

func NewDateTime(t time.Time) DateTime {
	return DateTime{Time: t}
}

type DateTimeProvider struct{}

func (p DateTimeProvider) Now() DateTime {
	return NewDateTime(time.Now())
}

// ParseDateTime creates a DateTime value from a string in "YYYY-MM-DD HH:mm" format.
// It will return an error if the string is not properly formatted.
func ParseDateTime(str string) (DateTime, error) {
	res, err := time.Parse(DefaultDateTimeFormat, str)
	return NewDateTime(res), err
}

// MustParseDateTime creates a DateTime object from a string in "YYYY-MM-DD HH:mm" format and
// panics if the string is not properly formatted.
func MustParseDateTime(str string) DateTime {
	res, err := ParseDateTime(str)
	if err != nil {
		panic(err)
	}
	return res
}

// ParseTime creates a DateTime value from a string in "HH:mm" format.
// It will return DateTime with date in "0000-01-01" format.
// It will return an error if the string is not properly formatted.
func ParseTime(str string) (DateTime, error) {
	res, err := time.Parse(DefaultTimeFormat, str)
	return NewDateTime(res), err
}

// MustParseTime creates a DateTime object from a string in "HH:mm" format.
// It will return DateTime with date in "0000-01-01" format.
// It will panic if the string is not properly formatted.
func MustParseTime(str string) DateTime {
	res, err := ParseTime(str)
	if err != nil {
		panic(err)
	}
	return res
}

// UnmarshalJSON implements the json.Unmarshaler interface and
// parses a string in "YYYY-MM-DD HH:MM" format to DateTime value.
func (dt *DateTime) UnmarshalJSON(input []byte) error {
	stringDateTime := strings.Trim(string(input), `"`)
	t, err := time.Parse(DefaultDateTimeFormat, stringDateTime)
	*dt = DateTime{t}
	return err
}

// MarshalJSON implements the json.Marshaler interface and
// serializes DateTime value in "YYYY-MM-DD HH:mm" format string.
func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", dt.String())), nil
}

// String returns a string representation of a DateTime in "YYYY-MM-DD HH:mm" format.
func (dt DateTime) String() string {
	return dt.Format(DefaultDateTimeFormat)
}

// DefaultDateFormat returns a string representation of a DateTime in "YYYY-MM-DD" format.
func (dt DateTime) DefaultDateFormat() string {
	return dt.Format(DefaultDateFormat)
}

// DefaultTimeFormat returns a string representation of a DateTime in "HH:mm" format.
func (dt DateTime) DefaultTimeFormat() string {
	return dt.Format(DefaultTimeFormat)
}

// LocalTimestamp returns unix timestamp + time zone offset.
func (dt DateTime) LocalTimestamp() int64 {
	_, offset := dt.Zone()
	return dt.Unix() + int64(offset)
}

func (dt DateTime) EncodeValues(key string, v *url.Values) error {
	v.Set(key, dt.String())
	return nil
}

func (dt DateTime) RoundDown(d time.Duration) DateTime {
	return NewDateTime(dt.Truncate(d))
}

func (dt DateTime) MinutesSinceMidnight() int {
	return dt.Hour()*60 + dt.Minute()
}

func (dt DateTime) MinutesSinceDateMidnight(since DateTime) int {
	sinceDate := NewDate(time.Date(since.Year(), since.Month(), since.Day(), 0, 0, 0, 0, since.Location()))
	return int(dt.Sub(sinceDate.Time).Minutes())
}

func (dt DateTime) RoundUp(d time.Duration) DateTime {
	if dt.Truncate(d) == dt.Time {
		return dt
	}
	return NewDateTime(dt.Add(d).Truncate(d))
}

func (dt DateTime) EndOfDay() DateTime {
	return NewDateTime(time.Date(dt.Year(), dt.Month(), dt.Day(), 23, 59, 59, 0, dt.Location()))
}

func (dt DateTime) WithDate(year int, month time.Month, day int) DateTime {
	return DateTime{
		time.Date(year, month, day, dt.Hour(), dt.Minute(), 0, 0, dt.Location()),
	}
}

func (dt DateTime) AsTime() Time {
	return NewTime(dt.Time)
}

func (dt DateTime) TruncateDate() DateTime {
	return NewDateTime(dt.AddDate(-dt.Year(), int(-dt.Month()+1), -dt.Day()+1))
}

func Now() DateTime {
	return NewDateTime(time.Now())
}
