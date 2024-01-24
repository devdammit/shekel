package gql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/devdammit/shekel/pkg/currency"
)

type Currency struct {
	currency.Code
}

// MarshalGQL implements the graphql.Marshaler interface.
func (c Currency) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(c.String()))) // nolint: errcheck
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (c *Currency) UnmarshalGQL(v interface{}) error {
	rawCode, ok := v.(string)
	if !ok {
		return fmt.Errorf("code must be a string")
	}

	code, err := currency.NewCode(rawCode)
	if err != nil {
		return err
	}

	*c = Currency{code}
	return nil
}
