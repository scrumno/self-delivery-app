package command

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/command"
	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/entity/claims"
	jwtEntity "github.com/scrumno/scrumno-api/internal/users/auth/jwt/entity"
	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/query"
	"github.com/scrumno/scrumno-api/shared/utils"
)

// RefreshRequest — тело запроса обновления токенов.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshHandler возвращает обработчик обмена refresh на новую пару access + refresh.
func RefreshHandler(jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Неверное тело запроса"}, http.StatusBadRequest)
			return
		}
		if req.RefreshToken == "" {
			utils.JSONResponse(w, map[string]string{"message": "refresh_token обязателен"}, http.StatusBadRequest)
			return
		}

		c, err := query.ValidateToken(jwtSecret, req.RefreshToken)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Невалидный или истёкший refresh токен"}, http.StatusUnauthorized)
			return
		}
		if c.TokenType != claims.TokenTypeRefresh {
			utils.JSONResponse(w, map[string]string{"message": "Требуется refresh токен"}, http.StatusBadRequest)
			return
		}

		var ut jwtEntity.UserToken
		err = config.DB.Where("user_id = ? AND refresh_token = ? AND expires_at > ?", c.UserID, req.RefreshToken, time.Now()).First(&ut).Error
		if err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Refresh токен отозван или не найден"}, http.StatusUnauthorized)
			return
		}

		accessToken, err := command.GenerateAccessToken(jwtSecret, c.UserID, c.Username)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Не удалось создать токены"}, http.StatusInternalServerError)
			return
		}
		newRefreshToken, err := command.GenerateRefreshToken(jwtSecret, c.UserID, c.Username)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Не удалось создать токены"}, http.StatusInternalServerError)
			return
		}
		if err := saveRefreshToken(c.UserID, newRefreshToken); err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Не удалось обновить сессию"}, http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, map[string]string{
			"access_token":  accessToken,
			"refresh_token": newRefreshToken,
			"expires_in":    "900",
		}, http.StatusOK)
	}
}
