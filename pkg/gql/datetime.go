package gql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/devdammit/shekel/pkg/types/datetime"
)

type Date struct {
	datetime.Date
}

// MarshalGQL implements the graphql.Marshaler interface.
func (d Date) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(d.String()))) // nolint: errcheck
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (d *Date) UnmarshalGQL(v interface{}) error {
	rawDate, ok := v.(string)
	if !ok {
		return fmt.Errorf("date must be a string")
	}

	date, err := datetime.ParseDate(rawDate)
	if err != nil {
		return err
	}

	*d = Date{date}
	return nil
}

type DateTime struct {
	datetime.DateTime
}

func FromDateTime(dt datetime.DateTime) DateTime {
	return DateTime{dt}
}
func FromDate(dt datetime.Date) Date {
	return Date{dt}
}

// MarshalGQL implements the graphql.Marshaler interface.
func (dt DateTime) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(dt.String()))) // nolint: errcheck
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (dt *DateTime) UnmarshalGQL(v interface{}) error {
	rawDate, ok := v.(string)
	if !ok {
		return fmt.Errorf("date must be a string")
	}

	dateTime, err := datetime.ParseDateTime(rawDate)
	if err != nil {
		return err
	}

	*dt = DateTime{dateTime}
	return nil
}
