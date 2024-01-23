package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"
	"fmt"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/cmd/unit/internal/handlers/graphql/model"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/gql"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

// Initialize is the resolver for the initialize field.
func (r *mutationResolver) Initialize(ctx context.Context, startDate gql.Date) (bool, error) {
	err := r.UseCases.Initialize.Execute(ctx, datetime.NewDate(startDate.Time))
	if err != nil {
		return false, err
	}

	return true, nil
}

// AddAccount is the resolver for the AddAccount field.
func (r *mutationResolver) AddAccount(ctx context.Context, contact model.CreateAccountInput) (bool, error) {
	var contactType entities.AccountType

	switch contact.Type {
	case model.AccountTypeCash:
		contactType = entities.AccountTypeCash
	case model.AccountTypeCredit:
		contactType = entities.AccountTypeCredit
	case model.AccountTypeDebit:
		contactType = entities.AccountTypeDebit
	default:
		return false, fmt.Errorf("invalid account type: %s", contact.Type)
	}

	ok, err := r.UseCases.CreateAccount.Execute(ctx, port.CreateAccountParams{
		Name:        contact.Name,
		Description: contact.Description,
		Type:        contactType,
		Balance: currency.Amount{
			Value:        contact.Balance.Amount,
			CurrencyCode: contact.Balance.Currency.Code,
		},
	})
	if err != nil {
		return false, err
	}

	return ok, nil
}

// Contacts is the resolver for the contacts field.
func (r *queryResolver) Contacts(ctx context.Context) ([]*model.Contact, error) {
	panic(fmt.Errorf("not implemented: Contacts - contacts"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
