package query

import (
	"net/http"

	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/users/user/entity/user"
	"github.com/scrumno/scrumno-api/shared/utils"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []user.User

	if err := config.DB.Find(&users).Error; err != nil {
		utils.JSONResponse(w, map[string]string{
			"error": err.Error(),
		}, 500)

		return
	}

	utils.JSONResponse(w, users, http.StatusOK)
}
