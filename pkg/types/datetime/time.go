package datetime

import (
	"time"
)

var timeLayout = "15:04:05"

// Time type represents a time with seconds precision.
// Time values can be created from HH:MM:SS strings and
// are serialized using the same layout.
// swaggen:type=string
type Time struct {
	time.Time
}

func NewTime(t time.Time) Time {
	return Time{Time: time.Date(0, time.January, 1, t.Hour(), t.Minute(), t.Second(), 0, t.Location())}
}

// ParseAsTime creates a Time value from a string in HH:MM:SS format.
// It will return an error if the string is not properly formatted.
func ParseAsTime(str string) (Time, error) {
	res, err := time.Parse(timeLayout, str)
	return NewTime(res), err
}

// MustParseAsTime creates a Time object from a string in HH:MM:SS format and
// panics if the string is not properly formatted.
func MustParseAsTime(str string) Time {
	res, err := ParseAsTime(str)
	if err != nil {
		panic(err)
	}
	return res
}

// String returns a string representation of a Time in HH:MM:SS format.
func (t Time) String() string {
	return t.Format(timeLayout)
}

// DefaultTimeFormat returns a string representation of a Time in HH:MM format.
func (t Time) DefaultTimeFormat() string {
	return t.Format(DefaultTimeFormat)
}

// MarshalJSON implements the json.Marshaler interface and
// serializes Time value in HH:MM:SS format string.
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeLayout)+2)
	return t.AppendFormat(b, `"`+timeLayout+`"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface and
// parses a string in HH:MM:SS format to Time value.
func (t *Time) UnmarshalJSON(data []byte) error {
	pt, err := time.Parse(`"`+timeLayout+`"`, string(data))
	*t = NewTime(pt)
	return err
}
