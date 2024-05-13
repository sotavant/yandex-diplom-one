package pgsql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/test"
	"github.com/sotavant/yandex-diplom-one/order"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderRepository_SetAccrual(t *testing.T) {
	ctx := context.Background()
	internal.InitLogger()

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool)
		assert.NoError(t, err)
	}(ctx, pool)

	userRepo, err := NewUserRepository(ctx, pool)
	assert.NoError(t, err)
	orderRepo, err := NewOrderRepository(ctx, pool)
	assert.NoError(t, err)

	var accrualSum float64 = 100
	domainUser := domain.User{
		Login:    "123",
		Password: "123134",
	}

	domainOrder := domain.Order{
		Number:  "123",
		Status:  order.StatusNew,
		Accrual: &accrualSum,
	}

	userID, err := userRepo.Store(ctx, domainUser)
	assert.NoError(t, err)

	domainUser.ID = userID
	domainOrder.UserID = userID

	orderID, err := orderRepo.Store(ctx, domainOrder)
	assert.NoError(t, err)
	domainOrder.ID = orderID

	err = orderRepo.SetAccrual(ctx, domainOrder)
	assert.NoError(t, err)

	var userCurrentSum *float64
	err = pool.QueryRow(ctx, `select current from users where id = $1`, userID).Scan(&userCurrentSum)
	assert.NoError(t, err)

	assert.Equal(t, accrualSum, *userCurrentSum)
}
