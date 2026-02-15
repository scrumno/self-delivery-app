package claims

import "github.com/golang-jwt/jwt/v5"

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"` // "access" | "refresh"
	jwt.RegisteredClaims
}
