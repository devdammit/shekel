package bbolt

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/log"
	"go.etcd.io/bbolt"
)

var AccountsBucket = []byte("accounts")

type data struct {
	val map[uint64]entities.Account
	mu  sync.RWMutex
}

type AccountsRepository struct {
	db   *resources.Bolt
	data data
}

func NewAccountsRepository(bolt *resources.Bolt) *AccountsRepository {
	return &AccountsRepository{
		db: bolt,
		data: data{
			val: make(map[uint64]entities.Account),
			mu:  sync.RWMutex{},
		},
	}
}

func (r *AccountsRepository) Start() error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket(resources.BoltRootBucket)

		_, err := root.CreateBucketIfNotExists(AccountsBucket)

		if err != nil {
			return fmt.Errorf("failed to create accounts bucket: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(AccountsBucket)

		err := bucket.ForEach(func(_, v []byte) error {
			var account entities.Account

			err := json.Unmarshal(v, &account)

			if err != nil {
				return fmt.Errorf("failed to unmarshal account: %w", err)
			}

			r.data.mu.Lock()
			defer r.data.mu.Unlock()

			r.data.val[account.ID] = account

			return nil
		})

		if err != nil {
			return err
		}

		log.Info("accounts loaded", log.Int("count", len(r.data.val)))

		return nil
	})
}

func (r *AccountsRepository) Create(account *entities.Account) (*entities.Account, error) {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(AccountsBucket)

		r.data.mu.RLock()
		defer r.data.mu.RUnlock()

		for _, a := range r.data.val {
			if a.Name == account.Name {
				return entities.ErrorAccountExists
			}
		}
		
		newID, _ := bucket.NextSequence()
		account.ID = newID

		data, err := json.Marshal(account)

		if err != nil {
			return fmt.Errorf("failed to marshal: %w", err)
		}

		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, account.ID)

		err = bucket.Put(key, data)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return account, nil
}
