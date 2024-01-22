package bbolt

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/configs"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
)

var AppConfigBucket = []byte("app_config")

type AppConfigRepository struct {
	config configs.AppConfig

	db *resources.Bolt
}

func NewAppConfigRepository(bolt *resources.Bolt) *AppConfigRepository {
	return &AppConfigRepository{
		db: bolt,
	}
}

func (r *AppConfigRepository) GetName() string {
	return "app_config"
}

func (r *AppConfigRepository) Start() error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket(resources.BoltRootBucket)

		_, err := root.CreateBucketIfNotExists(AppConfigBucket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(AppConfigBucket)

		err = bucket.ForEach(func(k, v []byte) error {
			var data configs.AppConfig

			err = json.Unmarshal(v, &data)
			if err != nil {
				return err
			}

			r.config = data

			return nil
		})
		if err != nil {
			return err
		}

		if r.config.DateStart != nil {
			log.Info("app config loaded", log.String("start_date", r.config.DateStart.String()), log.Bool("initialized", r.config.Initialized))
		} else {
			log.Info("app config loaded", log.Bool("initialized", r.config.Initialized))
		}

		return nil
	})
}

func (r *AppConfigRepository) Get() configs.AppConfig {
	return r.config
}

func (c *AppConfigRepository) GetStartDate() (*datetime.Date, error) {
	if c.config.DateStart == nil {
		return nil, errors.New("start date is not set")
	}

	return c.config.DateStart, nil
}

func (c *AppConfigRepository) SetStartDate(ctx context.Context, date datetime.Date) error {
	c.config.DateStart = &date

	return c.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(AppConfigBucket)

		data, err := json.Marshal(c.config)
		if err != nil {
			return err
		}

		return bucket.Put([]byte("config"), data)
	})
}
