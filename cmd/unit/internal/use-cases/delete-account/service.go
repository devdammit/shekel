package delete_account

type Repository interface {
	Delete(accountID uint64) error
	Archive(accountID uint64) error
}

type TransactionsRepository interface {
	GetCountByAccount(accountID uint64) (uint64, error)
}

type DeleteAccountUseCase struct {
	repo Repository

	transactionsRepo TransactionsRepository
}

// Execute
// Прежде чем удалять счет, нужно убедиться, что он не используется в транзакциях.
// Если используется, то архивируем счет, а не удаляем.
func (u *DeleteAccountUseCase) Execute(accountID uint64) (bool, error) {
	err := u.repo.Delete(accountID)
	if err != nil {
		return false, err
	}

	return true, nil
}
