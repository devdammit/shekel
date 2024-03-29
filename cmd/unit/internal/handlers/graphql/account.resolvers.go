package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/cmd/unit/internal/handlers/graphql/model"
	"github.com/devdammit/shekel/pkg/gql"
	"github.com/devdammit/shekel/pkg/pointer"
)

// Type is the resolver for the Type field.
func (r *accountResolver) Type(ctx context.Context, obj *entities.Account) (model.AccountType, error) {
	return model.AccountType(obj.Type), nil
}

// Balance is the resolver for the Balance field.
func (r *accountResolver) Balance(ctx context.Context, obj *entities.Account) (model.Amount, error) {
	return model.Amount{
		Currency: gql.Currency{Code: obj.Balance.CurrencyCode},
		Amount:   obj.Balance.Value,
	}, nil
}

// DeletedAt is the resolver for the deletedAt field.
func (r *accountResolver) DeletedAt(ctx context.Context, obj *entities.Account) (*gql.DateTime, error) {
	if obj.DeletedAt == nil {
		return nil, nil
	}

	return pointer.Ptr(gql.FromDateTime(*obj.DeletedAt)), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *accountResolver) CreatedAt(ctx context.Context, obj *entities.Account) (gql.DateTime, error) {
	return gql.FromDateTime(obj.CreatedAt), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *accountResolver) UpdatedAt(ctx context.Context, obj *entities.Account) (gql.DateTime, error) {
	return gql.FromDateTime(obj.UpdateAt), nil
}

// Account returns AccountResolver implementation.
func (r *Resolver) Account() AccountResolver { return &accountResolver{r} }

type accountResolver struct{ *Resolver }
