package middleware

import (
	"context"
	"harmony/pkg/jwt"
	"net/http"
	"strings"
)

type AuthData struct {
	Id string
}

func writeUnauthed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func IsAuthed(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authedHeader := r.Header.Get("Authorization")
			if authedHeader == "" || !strings.HasPrefix(authedHeader, "Bearer ") {
				writeUnauthed(w)
				return
			}
			token := strings.TrimPrefix(authedHeader, "Bearer ")
			isValid, data := jwt.NewJWT(secret).Parse(token)
			if !isValid {
				writeUnauthed(w)
				return
			}
			ctx := context.WithValue(r.Context(), "authData", AuthData{
				Id: data.Id,
			})
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}
}
