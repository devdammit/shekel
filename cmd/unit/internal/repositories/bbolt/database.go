package bbolt

import (
	"github.com/devdammit/shekel/internal/resources"
	"go.etcd.io/bbolt"
)

type Database struct {
	db *resources.Bolt
}

func NewDatabase(db *resources.Bolt) *Database {
	return &Database{db: db}
}

func (d *Database) Transaction(fn func() error) error {
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

	err = fn()
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
