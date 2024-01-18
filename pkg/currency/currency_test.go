package currency_test

import (
	"encoding/json"
	"testing"

	"github.com/devdammit/shekel/pkg/currency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCode(t *testing.T) {
	rubCode, err := currency.NewCode("rub")
	require.NoError(t, err)
	usdCode, err := currency.NewCode("USD")
	require.NoError(t, err)

	rubCodeRaw, err := json.Marshal(rubCode)
	require.NoError(t, err)
	usdCodeRaw, err := json.Marshal(usdCode)
	require.NoError(t, err)

	assert.EqualValues(t, `"RUB"`, rubCodeRaw)
	assert.EqualValues(t, `"USD"`, usdCodeRaw)

	var newRUBCode, newUSDCode currency.Code
	err = json.Unmarshal(rubCodeRaw, &newRUBCode)
	require.NoError(t, err)
	err = json.Unmarshal(usdCodeRaw, &newUSDCode)
	require.NoError(t, err)

	require.EqualValues(t, rubCode, newRUBCode)
	require.EqualValues(t, usdCode, newUSDCode)

	t.Run("check invalid rates", func(t *testing.T) {
		var rates currency.Rates
		require.NoError(t, json.Unmarshal([]byte(`{"invalid": 10, "rub": 1}`), &rates))
		assert.Len(t, rates, 1)
		assert.EqualValues(t, rates[currency.RUB], 1)
	})
}

func TestFullRates_Convert(t *testing.T) {
	rates := currency.FullRates{
		currency.RUB: currency.Rates{
			currency.RUB: 1,
			currency.USD: 74.167795,
		},
		currency.EUR: currency.Rates{
			currency.RUB: 0.012572181767032444,
			currency.USD: 0.932451,
		},
	}

	amount, err := rates.Convert(&currency.Amount{CurrencyCode: currency.USD, Value: 1}, currency.RUB)
	require.NoError(t, err)
	assert.EqualValues(t, 74.167795, amount.Value)

	amount, err = rates.Convert(&currency.Amount{CurrencyCode: currency.RUB, Value: 100}, currency.RUB)
	require.NoError(t, err)
	assert.EqualValues(t, 100, amount.Value)

	amount, err = rates.Convert(&currency.Amount{CurrencyCode: currency.RUB, Value: 100}, currency.EUR)
	require.NoError(t, err)
	assert.EqualValues(t, 1.2572181767032444, amount.Value)

	amount, err = rates.Convert(&currency.Amount{CurrencyCode: currency.USD, Value: 100}, currency.EUR)
	require.NoError(t, err)
	assert.InEpsilon(t, 93.2451, amount.Value, 0.0001)
}
