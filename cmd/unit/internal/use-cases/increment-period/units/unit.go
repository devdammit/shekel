package units

import (
	"context"
	"errors"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
)

type Unit interface {
	GetName() string
	Handle(ctx context.Context, request *Request, payload *Payload) (*Payload, error)
}

type PeriodsRepository interface {
	Update(ctx context.Context, period *entities.Period) error
	GetLast(ctx context.Context) (*entities.Period, error)
	Create(ctx context.Context) (*entities.Period, error)
}

type Payload struct {
	ActivePeriod *entities.Period
	Transactions []entities.Transaction
	Accounts     []entities.Account
}

type Request struct {
}

var (
	ErrPayloadCheckFailed = errors.New("payload check failed")
)

func NewPayload() *Payload {
	return &Payload{}
}

type ErrorUnitPanic string

func (e ErrorUnitPanic) Error() string {
	return string(e)
}
