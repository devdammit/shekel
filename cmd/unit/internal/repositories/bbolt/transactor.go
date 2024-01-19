package bbolt

import (
	"context"
	"github.com/devdammit/shekel/internal/resources"
	"go.etcd.io/bbolt"
)

type Transactor struct {
	db *resources.Bolt
}

func NewTransactor(db *resources.Bolt) *Transactor {
	return &Transactor{
		db: db,
	}
}

func (d *Transactor) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}

	defer func(tx *bbolt.Tx) {
		err := tx.Rollback()
		if err != nil {
			panic(err)
		}
	}(tx)

	err = fn(ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
