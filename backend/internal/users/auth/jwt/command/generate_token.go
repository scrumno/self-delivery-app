package command

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/entity/claims"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

// GenerateAccessToken создаёт короткоживущий JWT для запросов к API (в заголовок Authorization).
func GenerateAccessToken(secret []byte, userID, username string) (string, error) {
	return generateToken(secret, userID, username, claims.TokenTypeAccess, AccessTokenTTL)
}

// GenerateRefreshToken создаёт длинноживущий JWT только для обмена на новую пару (не для API-запросов).
func GenerateRefreshToken(secret []byte, userID, username string) (string, error) {
	return generateToken(secret, userID, username, claims.TokenTypeRefresh, RefreshTokenTTL)
}

func generateToken(secret []byte, userID, username, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	c := claims.Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}
