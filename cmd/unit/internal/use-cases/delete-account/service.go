package delete_account

type Repository interface {
	Delete(accountID uint64) error
}

type DeleteAccountUseCase struct {
	accounts Repository
}

// Execute
// Прежде чем удалять счет, нужно убедиться, что он не используется в транзакциях.
// Если используется, то архивируем счет, а не удаляем.
func (u *DeleteAccountUseCase) Execute(accountID uint64) (bool, error) {
	err := u.accounts.Delete(accountID)
	if err != nil {
		return false, err
	}

	return true, nil
}
