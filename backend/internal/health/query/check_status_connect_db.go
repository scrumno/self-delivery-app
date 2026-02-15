package query

import (
	"net/http"

	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/shared/utils"
)

func CheckStatusConnectDB(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Message string `json:"message,omitempty"`
	}

	sqlDB, err := config.DB.DB()
	if err != nil {
		utils.JSONResponse(w, response{
			Message: "Не удалось установить соединение с БД",
		}, http.StatusServiceUnavailable)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		utils.JSONResponse(w, response{Message: err.Error()}, http.StatusServiceUnavailable)
		return
	}

	utils.JSONResponse(w, response{Message: "Подключение к базе данных - окей"}, http.StatusOK)
}
