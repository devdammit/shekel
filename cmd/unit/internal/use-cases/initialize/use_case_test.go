package initialize_test

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/use-cases/initialize"
	"github.com/devdammit/shekel/cmd/unit/internal/use-cases/initialize"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("should create new period if none exists", func(t *testing.T) {
		var (
			mockController   = gomock.NewController(t)
			appConfig        = mocks.NewMockAppConfig(mockController)
			repository       = mocks.NewMockPeriodsRepository(mockController)
			datetimeProvider = mocks.NewMockDateTimeProvider(mockController)
		)

		appConfig.EXPECT().SetStartDate(gomock.Any(), datetime.MustParseDate("2023-11-01")).Return(nil)
		datetimeProvider.EXPECT().Now().Times(7).Return(datetime.MustParseDateTime("2024-01-22 19:19"))

		repository.EXPECT().GetCount(gomock.Any()).Times(1).Return(uint64(0), nil)

		repository.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, period entities.Period) {
			assert.Equal(t, datetime.MustParseDateTime("2023-11-01 00:00"), period.CreatedAt)
			assert.Equal(t, datetime.MustParseDateTime("2023-12-01 00:00"), *period.ClosedAt)
		}).Return(nil, nil)

		repository.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, period entities.Period) {
			assert.Equal(t, datetime.MustParseDateTime("2023-12-01 00:00"), period.CreatedAt)
			assert.Equal(t, datetime.MustParseDateTime("2024-01-01 00:00"), *period.ClosedAt)
		}).Return(nil, nil)

		repository.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, period entities.Period) {
			assert.Equal(t, datetime.MustParseDateTime("2024-01-01 00:00"), period.CreatedAt)
			assert.Nil(t, period.ClosedAt)
		}).Return(nil, nil)

		ctx := context.Background()
		uc := initialize.NewUseCase(repository, appConfig, datetimeProvider)

		err := uc.Execute(ctx, datetime.MustParseDate("2023-11-01"))
		assert.NoError(t, err)
	})

	t.Run("should return error if periods already initialized", func(t *testing.T) {
		var (
			mockController   = gomock.NewController(t)
			appConfig        = mocks.NewMockAppConfig(mockController)
			repository       = mocks.NewMockPeriodsRepository(mockController)
			datetimeProvider = mocks.NewMockDateTimeProvider(mockController)
		)

		repository.EXPECT().GetCount(gomock.Any()).Return(uint64(5), nil)

		ctx := context.Background()
		uc := initialize.NewUseCase(repository, appConfig, datetimeProvider)

		err := uc.Execute(ctx, datetime.MustParseDate("2023-11-01"))
		assert.EqualError(t, err, "periods already initialized")

	})
}
