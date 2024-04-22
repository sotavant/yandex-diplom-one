package rest

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/internal/repository/pgsql"
	"github.com/sotavant/yandex-diplom-one/internal/test"
	"github.com/sotavant/yandex-diplom-one/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Register(t *testing.T) {
	ctx := context.Background()

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)

	if pool == nil {
		fmt.Println("no databases init")
		return
	}

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool)
		assert.NoError(t, err)
	}(ctx, pool)

	userRepo, err := pgsql.NewUserRepository(ctx, pool)
	assert.NoError(t, err)

	userService := user.NewService(userRepo)

	type Want struct {
		status int
		body   string
	}

	tests := []struct {
		name string
		body string
		want Want
	}{
		{
			name: "register empty body",
			body: "",
			want: Want{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "register bad params",
			body: `{"login": "134", "pass": 345345}`,
			want: Want{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "register ok params",
			body: `{"login": "134", "password": "345345"}`,
			want: Want{
				status: http.StatusOK,
				body:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				Service: userService,
			}

			buf := bytes.NewBufferString(tt.body)

			req := httptest.NewRequest(http.MethodPost, "/api/user/register", buf)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			u.Auth(w, req)
			result := w.Result()

			defer func() {
				err := result.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, tt.want.status, result.StatusCode)
		})
	}
}
