package datetime

import (
	"net/url"
	"time"
)

var dateLayout = "2006-01-02"

// Date type represents a date with day precision.
// Date values can be created from YYYY-MM-DD strings and
// are serialized using the same layout.
// swaggen:type=string
// swaggen:format=date
type Date struct {
	time.Time
}

func NewDate(t time.Time) Date {
	return Date{Time: t}
}

// ParseDate creates a Date value from a string in YYYY-MM-DD format.
// It will return an error if the string is not properly formatted.
func ParseDate(str string) (Date, error) {
	res, err := time.Parse(dateLayout, str)
	return NewDate(res), err
}

// MustParseDate creates a Date object from a string in YYYY-MM-DD format and
// panics if the string is not properly formatted.
func MustParseDate(str string) Date {
	res, err := ParseDate(str)
	if err != nil {
		panic(err)
	}
	return res
}

// String returns a string representation of a Date in YYYY-MM-DD format.
func (d Date) String() string {
	return d.Format(dateLayout)
}

// MarshalJSON implements the json.Marshaler interface and
// serializes Date value in YYYY-MM-DD format string.
func (d Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(dateLayout)+2)
	return d.AppendFormat(b, `"`+dateLayout+`"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface and
// parses a string in YYYY-MM-DD format to Date value.
func (d *Date) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(`"`+dateLayout+`"`, string(data))
	*d = NewDate(t)
	return err
}

func (d Date) EncodeValues(key string, v *url.Values) error {
	v.Set(key, d.String())
	return nil
}
