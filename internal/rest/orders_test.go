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
	"github.com/sotavant/yandex-diplom-one/order"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrdersHandler_AddOrder(t *testing.T) {
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

	userID, err := userRepo.Store(ctx, domain.User{
		Login:    "123",
		Password: "123134",
	})
	assert.NoError(t, err)

	token, err := auth.BuildJWTString(userID)
	assert.NoError(t, err)

	orderService := order.NewOrderService(orderRepo)

	type Want struct {
		status int
	}

	tests := []struct {
		name string
		body string
		want Want
	}{
		{
			name: "new order",
			body: "87932",
			want: Want{
				status: http.StatusAccepted,
			},
		},
		{
			name: "already uploaded",
			body: "87932",
			want: Want{
				status: http.StatusOK,
			},
		},
		{
			name: "bad order number",
			body: "123123",
			want: Want{
				status: http.StatusUnprocessableEntity,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.body)
			u := &OrdersHandler{
				Service: orderService,
			}

			r := chi.NewRouter()
			r.Use(middleware.Auth)
			r.Post("/api/user/orders", u.AddOrder)

			w := httptest.NewRecorder()
			reqFunc := func() *http.Request {
				req := httptest.NewRequest("POST", "/api/user/orders", buf)
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "text/plain")

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
