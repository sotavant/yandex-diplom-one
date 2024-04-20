package middleware

import (
	"context"
	"github.com/sotavant/yandex-diplom-one/internal/auth"
	"github.com/sotavant/yandex-diplom-one/user"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	tokenSubstr         = "Bearer"
)

func Auth(h http.Handler) http.Handler {
	authFn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorizationHeader)
		if authHeader == "" {
			authFailed(w)
			return
		}

		if !strings.Contains(authHeader, tokenSubstr) {
			authFailed(w)
			return
		}

		token := strings.TrimSpace(strings.Replace(authHeader, tokenSubstr, "", -1))
		userID, err := auth.GetUserID(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userID == -1 {
			authFailed(w)
			return
		}

		ctx := context.WithValue(r.Context(), user.ContextUserIDKey{}, userID)

		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(authFn)
}

func authFailed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}
