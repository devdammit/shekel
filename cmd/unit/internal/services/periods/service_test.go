package periods_test

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/services/periods"
	"github.com/devdammit/shekel/cmd/unit/internal/services/periods"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestService_InitPeriods(t *testing.T) {
	t.Run("should create new period if none exists", func(t *testing.T) {
		var (
			mockController   = gomock.NewController(t)
			appConfig        = mocks.NewMockAppConfig(mockController)
			repository       = mocks.NewMockRepository(mockController)
			datetimeProvider = mocks.NewMockDateTimeProvider(mockController)
		)

		appConfig.EXPECT().GetStartYear().Return(datetime.MustParseDateTime("2023-11-01 14:20"), nil)
		datetimeProvider.EXPECT().Now().Times(7).Return(datetime.MustParseDateTime("2024-01-22 19:19"))

		repository.EXPECT().GetCount(gomock.Any()).Return(uint64(0), nil)

		repository.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, period entities.Period) {
			assert.Equal(t, datetime.MustParseDateTime("2023-11-01 14:20"), period.CreatedAt)
			assert.Equal(t, datetime.MustParseDateTime("2023-12-01 14:20"), *period.ClosedAt)
		}).Return(nil, nil)

		repository.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, period entities.Period) {
			assert.Equal(t, datetime.MustParseDateTime("2023-12-01 14:20"), period.CreatedAt)
			assert.Equal(t, datetime.MustParseDateTime("2024-01-01 14:20"), *period.ClosedAt)
		}).Return(nil, nil)

		repository.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, period entities.Period) {
			assert.Equal(t, datetime.MustParseDateTime("2024-01-01 14:20"), period.CreatedAt)
			assert.Nil(t, period.ClosedAt)
		}).Return(nil, nil)

		ctx := context.Background()
		service := periods.NewService(repository, appConfig, datetimeProvider)

		err := service.InitPeriods(ctx)
		assert.NoError(t, err)
	})

	t.Run("should return error if periods already initialized", func(t *testing.T) {
		var (
			mockController   = gomock.NewController(t)
			appConfig        = mocks.NewMockAppConfig(mockController)
			repository       = mocks.NewMockRepository(mockController)
			datetimeProvider = mocks.NewMockDateTimeProvider(mockController)
		)

		repository.EXPECT().GetCount(gomock.Any()).Return(uint64(5), nil)

		ctx := context.Background()
		service := periods.NewService(repository, appConfig, datetimeProvider)

		err := service.InitPeriods(ctx)
		assert.EqualError(t, err, "periods already initialized")

	})
}
