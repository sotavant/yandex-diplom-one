package order

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/repository/pgsql"
	"github.com/sotavant/yandex-diplom-one/internal/test"
	"github.com/sotavant/yandex-diplom-one/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_List(t *testing.T) {
	ctx := context.Background()
	internal.InitLogger()

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool)
		assert.NoError(t, err)
	}(ctx, pool)

	userRepo, err := pgsql.NewUserRepository(ctx, pool)
	assert.NoError(t, err)
	orderRepo, err := pgsql.NewOrderRepository(ctx, pool)
	assert.NoError(t, err)
	service := NewOrderService(orderRepo)

	domainUser := domain.User{
		Login:    "123",
		Password: "123134",
	}

	domainOrder := domain.Order{
		Number: "123",
		Status: StatusNew,
	}

	userID, err := userRepo.Store(ctx, domainUser)
	assert.NoError(t, err)

	domainUser.ID = userID
	domainOrder.UserID = userID

	ctxWithValue := context.WithValue(ctx, user.ContextUserIDKey{}, userID)
	t.Run("no orders", func(t *testing.T) {
		orders, resp, err := service.List(ctxWithValue)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(orders))
		assert.Equal(t, domain.RespNoDataToResponse, resp)
	})

	t.Run("one order", func(t *testing.T) {
		_, err = orderRepo.Store(ctx, domainOrder)
		assert.NoError(t, err)

		orders, resp, err := service.List(ctxWithValue)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(orders))
		assert.Equal(t, "", resp)
	})
}
