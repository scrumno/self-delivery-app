package query

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/entity/claims"
)

// ValidateToken проверяет JWT и возвращает claims. Secret передаётся снаружи (DIP).
func ValidateToken(secret []byte, tokenString string) (*claims.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*claims.Claims); ok && token.Valid {
		return c, nil
	}
	return nil, fmt.Errorf("invalid token")
}
