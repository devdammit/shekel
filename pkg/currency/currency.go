package currency

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/devdammit/shekel/pkg/encoding/json"
	xcurrency "golang.org/x/text/currency"
)

// swaggen: type=string
// swaggen: format=currency
type Code struct {
	xcurrency.Unit
}

// NewCode constructs Code from the provided string
func NewCode(code string) (Code, error) {
	var ret Code
	err := ret.FromString(code)
	return ret, err
}

func MustNewCode(code string) Code {
	c, err := NewCode(code)
	if err != nil {
		panic(err.Error())
	}
	return c
}

// FromString sets the Code from its string representation.
func (cc *Code) FromString(str string) error {
	// avoid err if str is empty
	if str == "" {
		*cc = Code{xcurrency.XXX}
		return nil
	}

	unit, err := xcurrency.ParseISO(str)
	if err != nil {
		return err
	}

	*cc = Code{unit}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface and
// serializes Code to its string representation.
// NOTE: zero currency code marshals as "XXX", need fix this?
func (cc Code) MarshalText() ([]byte, error) {
	return []byte(cc.String()), nil
}

func (cc Code) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(cc.String())), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface and
// parses Code from its string representation.
func (cc *Code) UnmarshalText(data []byte) error {
	err := cc.FromString(string(data))
	if err != nil {
		return fmt.Errorf("failed to parse '%s': %w", data, err)
	}
	return nil
}

func (cc *Code) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("failed to parse '%s': %w", data, err)
	}
	err = cc.FromString(s)
	if err != nil {
		return fmt.Errorf("failed to parse '%s': %w", data, err)
	}
	return nil
}

func (cc Code) EncodeValues(key string, v *url.Values) error {
	v.Set(key, cc.String())
	return nil
}

// Currency amount
type Amount struct {
	// Currency in which the price is presented
	CurrencyCode Code `json:"currency_code" swaggen:"required,type=string,format=currency"`

	// Float value of price
	Value float64 `json:"value" swaggen:"required"`
}

func (a Amount) Subtract(b Amount) Amount {
	if a.CurrencyCode != b.CurrencyCode {
		panic(fmt.Errorf("cannot substract %s from %s", b.CurrencyCode, a.CurrencyCode))
	}

	return Amount{
		CurrencyCode: a.CurrencyCode,
		Value:        a.Value - b.Value,
	}
}

func (a Amount) Add(b Amount) Amount {
	if a.CurrencyCode != b.CurrencyCode {
		panic(fmt.Errorf("cannot add %s to %s", b.CurrencyCode, a.CurrencyCode))
	}

	return Amount{
		CurrencyCode: a.CurrencyCode,
		Value:        a.Value + b.Value,
	}
}

func (a Amount) IsNegative() bool {
	return a.Value < 0
}

// FullRates to convert from one currency to another. Example: 1 USD = Rates[RUB][USD] * 1 = 77.374996 * 1 = 77.374996 RUB
type FullRates map[Code]Rates

func (rates FullRates) Convert(amount *Amount, to Code) (*Amount, error) {
	if amount.CurrencyCode == to {
		return amount, nil
	}

	if rates[to] == nil {
		return nil, fmt.Errorf("rates[%s][%s] is nil", to, amount.CurrencyCode)
	}

	if rates[to][amount.CurrencyCode] == 0 {
		return nil, fmt.Errorf("rates[%s][%s] is 0", to, amount.CurrencyCode)
	}

	return &Amount{
		CurrencyCode: to,
		Value:        amount.Value * rates[to][amount.CurrencyCode],
	}, nil
}

type Rates map[Code]float64

func (rates *Rates) UnmarshalJSON(data []byte) error {
	m := make(map[string]float64)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	res := make(Rates)
	for code, rate := range m {
		var c Code
		// ignore invalid codes
		if err := c.FromString(code); err == nil {
			res[c] = rate
		}
	}
	*rates = res
	return nil
}
