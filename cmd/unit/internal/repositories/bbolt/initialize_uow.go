package bbolt

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
)

type InitializeUow struct {
	startDate *datetime.Date
	periods   []entities.Period

	db *resources.Bolt
}

func NewInitializeUow() *InitializeUow {
	return &InitializeUow{}
}

func (u *InitializeUow) SetStartDate(date *datetime.Date) {
	u.startDate = date
}

func (u *InitializeUow) CreatePeriods(periods []entities.Period) error {
	u.periods = periods
	return nil
}

func (u *InitializeUow) Commit(ctx context.Context) error {
	if u.startDate == nil || u.periods == nil {
		return errors.New("missing data")
	}

	err := u.db.Update(func(tx *bbolt.Tx) error {

	})

	return nil
}
