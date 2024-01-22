package bbolt

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
)

var AccountsBucket = []byte("accounts")

type AccountsRepository struct {
	db   *resources.Bolt
	data map[uint64]entities.Account
	*sync.RWMutex
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

			r.Lock()
			defer r.Unlock()

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

		r.Lock()
		defer r.Unlock()

		for _, a := range r.data {
			if a.Name == account.Name && a.DeletedAt == nil {
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
	r.RLock()
	defer r.RUnlock()

	account, ok := r.data[id]

	if !ok {
		return nil, entities.ErrorAccountNotFound
	}

	return &account, nil
}

func (r *AccountsRepository) Update(account *entities.Account) (*entities.Account, error) {
	r.Lock()
	defer r.Unlock()

	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(AccountsBucket)

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

func (r *AccountsRepository) Delete(id uint64) error {
	r.Lock()
	defer r.Unlock()

	account, ok := r.data[id]

	if !ok {
		return entities.ErrorAccountNotFound
	}

	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(AccountsBucket)

		hasTransactions, err := r.hasTransactionsByAccountID(tx, id)

		if err != nil {
			return fmt.Errorf("failed to check has transactions: %w", err)
		}

		key := resources.Itob(int(account.ID))

		if hasTransactions {
			deletedAt := datetime.Now()
			account.DeletedAt = &deletedAt

			data, err := json.Marshal(account)

			if err != nil {
				return fmt.Errorf("failed to marshal: %w", err)
			}

			if err := bucket.Put(key, data); err != nil {
				return fmt.Errorf("failed to put in bucket: %w", err)
			}

			r.data[account.ID] = account

			return nil
		}

		if err := bucket.Delete(key); err != nil {
			return fmt.Errorf("failed to delete from bucket: %w", err)
		}

		delete(r.data, account.ID)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *AccountsRepository) hasTransactionsByAccountID(tx *bbolt.Tx, id uint64) (bool, error) {
	bucket := tx.Bucket(resources.BoltRootBucket).Bucket(TransactionsBucket)
	cursor := bucket.Cursor()

	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		var transaction entities.Transaction

		err := json.Unmarshal(v, &transaction)

		if err != nil {
			return false, fmt.Errorf("failed to unmarshal transaction: %w", err)
		}

		if transaction.From.ID == id || transaction.To.ID == id {
			return true, nil
		}
	}

	return false, nil
}
