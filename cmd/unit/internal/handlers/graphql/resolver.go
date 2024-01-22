package graphql

import (
	"context"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type InitializeUseCase interface {
	Execute(ctx context.Context, startDate datetime.Date) error
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, params port.CreateAccountParams) (bool, error)
}

type UseCases struct {
	Initialize    InitializeUseCase
	CreateAccount CreateAccountUseCase
}

type Resolver struct {
	UseCases UseCases
}
