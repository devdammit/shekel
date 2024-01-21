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
	data map[uint64]entities.Account
	mu   sync.RWMutex
}

func NewAccountsRepository(bolt *resources.Bolt) *AccountsRepository {
	return &AccountsRepository{
		db:   bolt,
		data: make(map[uint64]entities.Account),
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

			r.mu.Lock()
			defer r.mu.Unlock()

			r.data[account.ID] = account

			return nil
		})

		if err != nil {
			return err
		}

		log.Info("accounts loaded", log.Int("count", len(r.data)))

		return nil
	})
}

func (r *AccountsRepository) Create(account *entities.Account) (*entities.Account, error) {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(AccountsBucket)

		r.mu.Lock()
		defer r.mu.Unlock()

		for _, a := range r.data {
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

		err = bucket.Put(resources.Itob(int(account.ID)), data)

		if err != nil {
			return fmt.Errorf("failed to put in bucket: %w", err)
		}

		r.data[account.ID] = *account

		return nil
	})

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountsRepository) GetByID(id uint64) (*entities.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, ok := r.data[id]

	if !ok {
		return nil, entities.ErrorAccountNotFound
	}

	return &account, nil
}
