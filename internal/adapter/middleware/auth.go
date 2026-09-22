package middleware

import (
	"context"
	"net/http"
	"strings"
	"user-management-api/internal/adapter/client"
	"user-management-api/internal/constant"
	"user-management-api/internal/core/port"
)

type contextKey string

const UserIDKey contextKey = "userID"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Auth(tokenProvider port.TokenProvider) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			const prefix = "Bearer "
			if authHeader == "" || !strings.HasPrefix(authHeader, prefix) {
				client.WriteFail(w, http.StatusUnauthorized, constant.UNAUTHORIZED_CODE, "missing token")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, prefix)

			userID, err := tokenProvider.ValidateToken(tokenStr)
			if err != nil {
				client.WriteFail(w, http.StatusUnauthorized, constant.UNAUTHORIZED_CODE, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next(w, r.WithContext(ctx))
		}

	}
}
