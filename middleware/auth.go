package middleware

import (
	"net/http"
	"strings"
)

func Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header ", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "invalid token format ", http.StatusUnauthorized)
			return
		}

		// if everything works well we pass on the next handler func
		next.ServeHTTP(w, r)
	})

}
