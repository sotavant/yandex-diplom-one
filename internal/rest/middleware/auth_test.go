package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/sotavant/yandex-diplom-one/internal/auth"
	user2 "github.com/sotavant/yandex-diplom-one/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	type Want struct {
		status int
	}

	var userId int64 = 123

	handler := func(w http.ResponseWriter, r *http.Request) {
		userContextId := r.Context().Value(user2.ContextUserIDKey{}).(int64)
		if userContextId != userId {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	token, err := auth.BuildJWTString(userId)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		withAuth bool
		want     Want
	}{
		{
			name:     "with auth",
			want:     Want{status: http.StatusOK},
			withAuth: true,
		},
		{
			name:     "without auth",
			want:     Want{status: http.StatusUnauthorized},
			withAuth: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(Auth)
			r.Post("/auth", handler)

			w := httptest.NewRecorder()
			reqFunc := func() *http.Request {
				req := httptest.NewRequest("POST", "/auth", nil)
				if tt.withAuth {
					req.Header.Set("Authorization", "Bearer "+token)
				}

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
