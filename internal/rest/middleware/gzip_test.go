package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzip(t *testing.T) {
	type Want struct {
		status int
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var user domain.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	tests := []struct {
		name string
		body string
		gzip bool
		want Want
	}{
		{
			name: "with gzip",
			body: `{"login": "123", "password": "123123"}`,
			want: Want{status: http.StatusOK},
			gzip: true,
		},
		{
			name: "without gzip",
			body: `{"login": "123", "password": "123123"}`,
			want: Want{status: http.StatusOK},
			gzip: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf *bytes.Buffer
			if tt.gzip {
				buf = bytes.NewBuffer(nil)
				zb := gzip.NewWriter(buf)
				_, err := zb.Write([]byte(tt.body))
				assert.NoError(t, err)
				err = zb.Close()
				assert.NoError(t, err)
			} else {
				buf = bytes.NewBufferString(tt.body)
			}

			r := chi.NewRouter()
			r.Use(Gzip)
			r.Post("/gzip", handler)

			w := httptest.NewRecorder()
			reqFunc := func() *http.Request {
				req := httptest.NewRequest("POST", "/gzip", buf)
				if tt.gzip {
					req.Header.Set("Content-Encoding", "gzip")
				}

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
