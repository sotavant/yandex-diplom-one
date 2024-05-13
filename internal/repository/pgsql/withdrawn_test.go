package pgsql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithdrawnRepository_Store(t *testing.T) {
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
	withdrawnRepo, err := NewWithdrawnRepository(ctx, pool)
	assert.NoError(t, err)

	var withdrawnSum float64 = 100
	domainUser := domain.User{
		Login:    "123",
		Password: "123134",
	}

	domainWithdrawn := domain.Withdrawn{
		OrderNum: "123",
		Sum:      withdrawnSum,
	}

	userID, err := userRepo.Store(ctx, domainUser)
	assert.NoError(t, err)

	domainUser.ID = userID
	domainWithdrawn.UserID = userID

	_, err = pool.Exec(ctx, "update users set current = $1 where id = $2", withdrawnSum, userID)
	assert.NoError(t, err)

	err = withdrawnRepo.Store(ctx, domainWithdrawn)
	assert.NoError(t, err)

	var userCurrentSum, userWithdrawnSum *float64
	err = pool.QueryRow(ctx, `select current, withdrawn from users where id = $1`, userID).Scan(&userCurrentSum, &userWithdrawnSum)
	assert.NoError(t, err)

	assert.Equal(t, float64(0), *userCurrentSum)
	assert.Equal(t, withdrawnSum, *userWithdrawnSum)
}
