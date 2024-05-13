package rest

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/auth"
	"github.com/sotavant/yandex-diplom-one/internal/repository/pgsql"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
	"github.com/sotavant/yandex-diplom-one/internal/test"
	"github.com/sotavant/yandex-diplom-one/withdrawn"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithdrawnHandler_Add(t *testing.T) {
	ctx := context.Background()
	internal.InitLogger()
	var current float64 = 500

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool)
		assert.NoError(t, err)
	}(ctx, pool)

	userRepo, err := pgsql.NewUserRepository(ctx, pool)
	assert.NoError(t, err)
	withdrawnRepo, err := pgsql.NewWithdrawnRepository(ctx, pool)
	assert.NoError(t, err)

	userID, err := userRepo.Store(ctx, domain.User{
		Login:    "123",
		Password: "123134",
	})
	assert.NoError(t, err)

	token, err := auth.BuildJWTString(userID)
	assert.NoError(t, err)

	withdrawnService := withdrawn.NewService(withdrawnRepo, userRepo)
	addUserBalance(ctx, t, pool, userID, current)

	type Want struct {
		status        int
		userCurrent   float64
		userWithdrawn float64
	}

	tests := []struct {
		name string
		body string
		want Want
	}{
		{
			name: "new withdraw",
			body: `{"order": "22962", "sum": 100}`,
			want: Want{
				status:        http.StatusOK,
				userCurrent:   400,
				userWithdrawn: 100,
			},
		},
		{
			name: "big sum",
			body: `{"order": "37473", "sum": 1000}`,
			want: Want{
				status: http.StatusPaymentRequired,
			},
		},
		{
			name: "bad order num",
			body: `{"order": "345345345", "sum": 1000}`,
			want: Want{
				status: http.StatusUnprocessableEntity,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.body)
			u := &WithdrawnHandler{
				Service: withdrawnService,
			}

			r := chi.NewRouter()
			r.Use(middleware.Auth)
			r.Post("/api/user/balance/withdraw", u.Add)

			w := httptest.NewRecorder()
			reqFunc := func() *http.Request {
				req := httptest.NewRequest("POST", "/api/user/balance/withdraw", buf)
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")

				return req
			}

			r.ServeHTTP(w, reqFunc())
			result := w.Result()
			defer func() {
				err := result.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, tt.want.status, result.StatusCode)
		})
	}
}

func addUserBalance(ctx context.Context, t *testing.T, pool *pgxpool.Pool, userID int64, current float64) {
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	query := `update users set current = $1 where id = $2`

	_, err = pool.Exec(ctx, query, current, userID)
	assert.NoError(t, err)

}
