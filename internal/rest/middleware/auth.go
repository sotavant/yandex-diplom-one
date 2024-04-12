package middleware

import "net/http"

const authorizationHeader = "Authorization"

func Auth(h http.Handler) http.Handler {
	authFn := func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(authorizationHeader)
		if auth == "" {
			authFailed(w)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(authFn)
}

func authFailed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}
