package delete_contact

import "context"

type ContractsRepository interface {
	Remove(ctx context.Context, id uint64) error
}

type UseCase struct {
	contracts ContractsRepository
}

func NewUseCase(contracts ContractsRepository) *UseCase {
	return &UseCase{
		contracts: contracts,
	}
}

func (uc *UseCase) Execute(ctx context.Context, ID uint64) error {
	return uc.contracts.Remove(ctx, ID)
}
