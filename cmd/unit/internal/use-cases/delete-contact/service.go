package delete_contact

import "context"

type ContractsRepository interface {
	DeleteContact(ctx context.Context, id uint64) error
}

type DeleteContactUseCase struct {
	contracts ContractsRepository
}

func NewDeleteContactUseCase(contracts ContractsRepository) *DeleteContactUseCase {
	return &DeleteContactUseCase{
		contracts: contracts,
	}
}

func (uc *DeleteContactUseCase) Handle(ctx context.Context, ID uint64) error {
	return uc.contracts.DeleteContact(ctx, ID)
}
