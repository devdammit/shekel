package bbolt

import (
	"encoding/json"
	"fmt"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/log"
	openexchange "github.com/devdammit/shekel/pkg/open-exchange"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
)

var BaseCurrency = currency.USD
var CurrencyRatesBucket = []byte("currency_rates")
var SupportedCurrencies = []currency.Code{
	currency.USD,
	currency.EUR,
	currency.THB,
	currency.RUB,
}

type OpenExchangeRatesAPI interface {
	GetByDate(base currency.Code, symbols []currency.Code, date datetime.Date) (*openexchange.HistoricalRates, error)
}

type CurrencyRatesRepository struct {
	api OpenExchangeRatesAPI
	db  *resources.Bolt

	data map[datetime.Date]currency.Rates
}

func NewCurrencyRatesRepository(bolt *resources.Bolt, api OpenExchangeRatesAPI) *CurrencyRatesRepository {
	return &CurrencyRatesRepository{
		api:  api,
		db:   bolt,
		data: make(map[datetime.Date]currency.Rates),
	}
}

func (c *CurrencyRatesRepository) GetName() string {
	return "currency_rates"
}

func (c *CurrencyRatesRepository) Start() error {
	err := c.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket(resources.BoltRootBucket)

		_, err := root.CreateBucketIfNotExists(CurrencyRatesBucket)
		if err != nil {
			return fmt.Errorf("could not create skus bucket: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(CurrencyRatesBucket)

		err = bucket.ForEach(func(k, v []byte) error {
			var data map[datetime.Date]currency.Rates

			err = json.Unmarshal(v, &data)
			if err != nil {
				return err
			}

			c.data = data

			return nil
		})
		if err != nil {
			return err
		}

		log.Info("currency rates loaded", log.Int("count", len(c.data)))

		return nil
	})
}

func (c *CurrencyRatesRepository) GetCurrencyRateByDate(source, target currency.Code, date datetime.Date) (float64, error) {
	if source == target {
		return 1, nil
	}

	if c.data[date] == nil {
		rates, err := c.api.GetByDate(BaseCurrency, SupportedCurrencies, date)
		if err != nil {
			return 0, err
		}

		err = c.db.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket(resources.BoltRootBucket).Bucket(CurrencyRatesBucket)

			c.data[date] = make(currency.Rates)

			for code, rate := range rates.Rates {
				c.data[date][code] = rate
			}

			data, err := json.Marshal(c.data)
			if err != nil {
				return err
			}

			return bucket.Put([]byte(date.String()), data)
		})

		if err != nil {
			return 0, err
		}
	}

	return c.data[date][target] / c.data[date][source], nil
}
