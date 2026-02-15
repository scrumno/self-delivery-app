package command

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/command"
	jwtEntity "github.com/scrumno/scrumno-api/internal/users/auth/jwt/entity"
	"github.com/scrumno/scrumno-api/internal/users/user/entity/user"
	"github.com/scrumno/scrumno-api/shared/utils"
)

// LoginRequest — тело запроса на логин.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler возвращает обработчик логина. Секрет JWT инжектируется (DIP).
func LoginHandler(jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Неверное тело запроса"}, http.StatusBadRequest)
			return
		}
		if req.Username == "" || req.Password == "" {
			utils.JSONResponse(w, map[string]string{"message": "Логин и пароль обязательны"}, http.StatusBadRequest)
			return
		}
		phone := utils.NormalizePhone(req.Username)
		if !utils.ValidatePhone(req.Username) {
			utils.JSONResponse(w, map[string]string{"message": "Номер телефона в формате 79996663355 (11 цифр, первая 7)"}, http.StatusBadRequest)
			return
		}

		u, err := findUserByUsername(phone)
		if err != nil || u == nil {
			utils.JSONResponse(w, map[string]string{"message": "Неверный логин или пароль"}, http.StatusBadRequest)
			return
		}
		if u.Password != req.Password {
			utils.JSONResponse(w, map[string]string{"message": "Неверный логин или пароль"}, http.StatusBadRequest)
			return
		}

		userIDStr := strconv.FormatUint(uint64(u.ID), 10)
		accessToken, err := command.GenerateAccessToken(jwtSecret, userIDStr, u.Username)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Не удалось создать токен"}, http.StatusInternalServerError)
			return
		}
		refreshToken, err := command.GenerateRefreshToken(jwtSecret, userIDStr, u.Username)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Не удалось создать токен"}, http.StatusInternalServerError)
			return
		}
		if err := saveRefreshToken(userIDStr, refreshToken); err != nil {
			slog.Error("save refresh token", "error", err)
			utils.JSONResponse(w, map[string]string{"message": "Не удалось сохранить сессию"}, http.StatusInternalServerError)
			return
		}
		utils.JSONResponse(w, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    "900", // секунды жизни access (15 мин)
		}, http.StatusOK)
	}
}

func findUserByUsername(username string) (*user.User, error) {
	var u user.User
	err := config.DB.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func saveRefreshToken(userID, refreshToken string) error {
	expiresAt := time.Now().Add(command.RefreshTokenTTL)
	var ut jwtEntity.UserToken
	err := config.DB.Where("user_id = ?", userID).First(&ut).Error
	if err != nil {
		ut = jwtEntity.UserToken{UserID: userID, RefreshToken: refreshToken, ExpiresAt: expiresAt, CreatedAt: time.Now()}
		return config.DB.Create(&ut).Error
	}
	ut.RefreshToken = refreshToken
	ut.ExpiresAt = expiresAt
	return config.DB.Save(&ut).Error
}
