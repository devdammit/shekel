package delete_account

import "context"

type Repository interface {
	Delete(ctx context.Context, accountID uint64) error
}

type DeleteAccountUseCase struct {
	accounts Repository
}

func NewUseCase(accounts Repository) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{
		accounts: accounts,
	}
}

// Execute
// Прежде чем удалять счет, нужно убедиться, что он не используется в транзакциях.
// Если используется, то архивируем счет, а не удаляем.
func (u *DeleteAccountUseCase) Execute(ctx context.Context, accountID uint64) (bool, error) {
	err := u.accounts.Delete(ctx, accountID)
	if err != nil {
		return false, err
	}

	return true, nil
}
