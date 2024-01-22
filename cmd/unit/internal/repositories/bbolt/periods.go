package bbolt

import (
	"encoding/json"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories/periods"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
	"sync"
)

var PeriodsBucket = []byte("periods")

type AppConfigService interface {
	GetStartYear() (datetime.DateTime, error)
}

type DateTimeProvider interface {
	Now() datetime.DateTime
}

type PeriodsRepository struct {
	db      *resources.Bolt
	periods map[uint64]entities.Period

	appConfig AppConfigService
	dateTime  DateTimeProvider
	sync.RWMutex
}

func NewPeriodsRepository(db *resources.Bolt, appConfig AppConfigService, provider DateTimeProvider) *PeriodsRepository {
	return &PeriodsRepository{
		db:        db,
		periods:   make(map[uint64]entities.Period),
		appConfig: appConfig,
		dateTime:  provider,
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

			r.periods[period.ID] = period

			return nil
		})
		if err != nil {
			return err
		}

		log.Info("periods loaded", log.Int("count", len(r.periods)))

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *PeriodsRepository) Create(period entities.Period) error {
	if len(r.periods) != 0 && r.periods[uint64(len(r.periods)-1)].ClosedAt == nil {
		return port.ErrHasOpenedPeriod
	}

	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(PeriodsBucket)
		ID, _ := bucket.NextSequence()

		period.ID = ID

		data, err := json.Marshal(period)
		if err != nil {
			return err
		}

		return bucket.Put(resources.Itob(int(period.ID)), data)
	})
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()

	r.periods[period.ID] = period

	return nil
}
