package bbolt

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/internal/resources"
	"go.etcd.io/bbolt"
)

var AccountsBucket = []byte("accounts")

type AccountsRepository struct {
	db *resources.Bolt
}

func NewAccountsRepository(bolt *resources.Bolt) *AccountsRepository {
	return &AccountsRepository{
		db: bolt,
	}
}

func (r *AccountsRepository) Create(account *entities.Account) (*entities.Account, error) {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(resources.BoltRootBucket).Bucket(AccountsBucket)

		err := bucket.ForEach(func(k, v []byte) error {
			var a entities.Account

			if err := json.Unmarshal(v, &a); err != nil {
				return fmt.Errorf("failed to unmarshal: %w", err)
			}

			if a.Name == account.Name {
				return entities.ErrorAccountExists
			}

			return nil
		})

		if err != nil {
			return err
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
