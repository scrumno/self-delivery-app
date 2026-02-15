package command

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/users/auth/jwt/command"
	"github.com/scrumno/scrumno-api/internal/users/user/entity/user"
	"github.com/scrumno/scrumno-api/shared/utils"
)

// SignUpRequest — тело запроса регистрации. Username = номер телефона.
type SignUpRequest struct {
	Username string `json:"username"` // номер телефона
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// SignUpHandler возвращает обработчик регистрации. Секрет JWT инжектируется (DIP).
func SignUpHandler(jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Неверное тело запроса"}, http.StatusBadRequest)
			return
		}

		phone := utils.NormalizePhone(req.Username)
		if !utils.ValidatePhone(req.Username) {
			utils.JSONResponse(w, map[string]string{"message": "Номер телефона в формате 79996663355 (11 цифр, первая 7)"}, http.StatusBadRequest)
			return
		}
		if req.Password == "" {
			utils.JSONResponse(w, map[string]string{"message": "Пароль обязателен"}, http.StatusBadRequest)
			return
		}

		existing, _ := findUserByUsername(phone)
		if existing != nil {
			utils.JSONResponse(w, map[string]string{"message": "Пользователь с таким номером телефона уже зарегистрирован"}, http.StatusConflict)
			return
		}

		u := user.User{
			Username: phone,
			Password: req.Password,
			FullName: req.FullName,
		}
		if err := config.DB.Create(&u).Error; err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Не удалось создать пользователя"}, http.StatusInternalServerError)
			return
		}

		userIDStr := strconv.FormatUint(uint64(u.ID), 10)
		accessToken, err := command.GenerateAccessToken(jwtSecret, userIDStr, u.Username)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Пользователь создан, но не удалось выдать токен"}, http.StatusInternalServerError)
			return
		}
		refreshToken, err := command.GenerateRefreshToken(jwtSecret, userIDStr, u.Username)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Пользователь создан, но не удалось выдать токен"}, http.StatusInternalServerError)
			return
		}
		if err := saveRefreshToken(userIDStr, refreshToken); err != nil {
			utils.JSONResponse(w, map[string]string{"message": "Пользователь создан, но не удалось сохранить сессию"}, http.StatusInternalServerError)
			return
		}
		utils.JSONResponse(w, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    "900",
		}, http.StatusCreated)
	}
}
