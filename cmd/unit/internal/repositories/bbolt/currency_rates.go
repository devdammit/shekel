package bbolt

import (
	"context"
	"encoding/json"
	"fmt"

	port "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories/currency-rates"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
)

var CurrencyRatesBucket = []byte("currency_rates")

type CurrencyRatesRepository struct {
	db *resources.Bolt

	data map[datetime.Date]currency.Rates
}

func NewCurrencyRatesRepository(bolt *resources.Bolt) *CurrencyRatesRepository {
	return &CurrencyRatesRepository{
		db:   bolt,
		data: make(map[datetime.Date]currency.Rates),
	}
}

func (r *CurrencyRatesRepository) GetName() string {
	return "currency_rates"
}

func (r *CurrencyRatesRepository) Start() error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
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

	return r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(CurrencyRatesBucket)

		err = bucket.ForEach(func(k, v []byte) error {
			var data map[datetime.Date]currency.Rates

			err = json.Unmarshal(v, &data)
			if err != nil {
				return err
			}

			r.data = data

			return nil
		})
		if err != nil {
			return err
		}

		log.Info("currency rates loaded", log.Int("count", len(r.data)))

		return nil
	})
}

func (r *CurrencyRatesRepository) GetCurrencyRateByDate(_ context.Context, code currency.Code, date datetime.Date) (float64, error) {
	if r.data[date] == nil {
		return 0, port.ErrRateNotFound
	}

	return r.data[date][code], nil
}
