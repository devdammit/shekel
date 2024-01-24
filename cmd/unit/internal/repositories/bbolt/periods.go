package bbolt

import (
	"context"
	"encoding/json"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
)

var PeriodsBucket = []byte("periods") //nolint:gochecknoglobals

type DateTimeProvider interface {
	Now() datetime.DateTime
}

type PeriodsRepository struct {
	db   *resources.Bolt
	data map[uint64]entities.Period
	keys []uint64

	dateTime DateTimeProvider
}

func NewPeriodsRepository(db *resources.Bolt, provider DateTimeProvider) *PeriodsRepository {
	return &PeriodsRepository{
		db:       db,
		data:     make(map[uint64]entities.Period),
		keys:     make([]uint64, 0),
		dateTime: provider,
	}
}

func (r *PeriodsRepository) GetName() string {
	return "periods"
}

func (r *PeriodsRepository) Start() error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket(resources.BoltRootBucket)

		_, err := root.CreateBucketIfNotExists(PeriodsBucket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(PeriodsBucket)

		err = bucket.ForEach(func(k, v []byte) error {
			var period entities.Period

			err = json.Unmarshal(v, &period)
			if err != nil {
				return err
			}

			r.data[period.ID] = period
			r.keys = append(r.keys, period.ID)

			return nil
		})
		if err != nil {
			return err
		}

		log.Info("periods loaded", log.Int("count", len(r.data)))

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *PeriodsRepository) Create(ctx context.Context, period entities.Period) (*entities.Period, error) {
	var entity entities.Period

	err := r.db.Update(func(tx *bbolt.Tx) error {
		p, err := r.CreateTx(ctx, tx, period)
		if err != nil {
			return err
		}

		entity = *p

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *PeriodsRepository) CreateTx(
	_ context.Context,
	tx *bbolt.Tx,
	period entities.Period,
) (*entities.Period, error) {
	if len(r.keys) != 0 && r.data[r.keys[len(r.keys)-1]].ClosedAt == nil {
		return nil, port.ErrHasOpenedPeriod
	}

	bucket := tx.Bucket(resources.BoltRootBucket).Bucket(PeriodsBucket)
	id, _ := bucket.NextSequence()

	period.ID = id

	data, err := json.Marshal(period)
	if err != nil {
		return nil, err
	}

	err = bucket.Put(resources.Itob(int(period.ID)), data)
	if err != nil {
		return nil, err
	}

	r.data[period.ID] = period
	r.keys = append(r.keys, period.ID)

	return &period, nil
}

func (r *PeriodsRepository) GetCount(_ context.Context) (uint64, error) {
	return uint64(len(r.data)), nil
}

func (r *PeriodsRepository) GetLast(_ context.Context) (*entities.Period, error) {
	var period entities.Period
	if len(r.data) == 0 {
		return nil, port.ErrNotFound
	}

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(resources.BoltRootBucket).Bucket(PeriodsBucket)

		c := b.Cursor()

		_, v := c.Last()

		err := json.Unmarshal(v, &period)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &period, nil
}

func (r *PeriodsRepository) GetAll(_ context.Context, limit *uint64, offset *uint64) ([]entities.Period, error) {
	var periods []entities.Period

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(resources.BoltRootBucket).Bucket(PeriodsBucket)

		c := b.Cursor()

		var i uint64
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if limit != nil && uint64(len(periods)) >= *limit {
				break
			}

			if offset != nil && i < *offset {
				i++
				continue
			}

			var period entities.Period

			err := json.Unmarshal(v, &period)
			if err != nil {
				return err
			}

			periods = append(periods, period)

			i++
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return periods, nil
}
