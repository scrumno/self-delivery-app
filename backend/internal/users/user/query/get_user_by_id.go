package query

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/users/user/entity/user"
	"github.com/scrumno/scrumno-api/shared/utils"
)

// GetUserByID возвращает пользователя по ID из пути /users/{id}.
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		utils.JSONResponse(w, map[string]string{"message": "ID обязателен"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"message": "Некорректный ID"}, http.StatusBadRequest)
		return
	}

	var u user.User
	err = config.DB.First(&u, id).Error
	if err != nil {
		utils.JSONResponse(w, map[string]string{"message": "Пользователь не найден"}, http.StatusNotFound)
		return
	}
	utils.JSONResponse(w, u, http.StatusOK)
}
