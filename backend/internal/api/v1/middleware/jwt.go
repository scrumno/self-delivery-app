package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/entity/claims"
	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/query"
	"github.com/scrumno/scrumno-api/shared/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

// UserInfo — данные пользователя из JWT, доступные в контексте.
type UserInfo struct {
	UserID   string
	Username string
}

// JWT возвращает middleware, проверяющий Bearer JWT. Secret инжектируется (DIP).
func JWT(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.JSONResponse(w, map[string]string{"message": "Не указан заголовок авторизации"}, http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.JSONResponse(w, map[string]string{"message": "Неверный формат заголовка Authorization"}, http.StatusUnauthorized)
				return
			}
			c, err := query.ValidateToken(secret, parts[1])
			if err != nil {
				utils.JSONResponse(w, map[string]string{"message": "Невалидный токен"}, http.StatusForbidden)
				return
			}
			if c.TokenType != claims.TokenTypeAccess {
				utils.JSONResponse(w, map[string]string{"message": "Требуется access токен, не refresh"}, http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, &UserInfo{
				UserID:   c.UserID,
				Username: c.Username,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
