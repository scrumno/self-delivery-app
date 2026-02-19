package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	updateUser "github.com/scrumno/scrumno-api/internal/users/command/update-user"
)

type UpdateUserAction struct {
	handler *updateUser.Handler
}

func NewUpdateUserAction(handler *updateUser.Handler) *UpdateUserAction {
	return &UpdateUserAction{handler: handler}
}

func (a *UpdateUserAction) Action(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		utils.JSONResponse(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	var cmd updateUser.Command
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		utils.JSONResponse(w, map[string]string{"error": "invalid request body"}, http.StatusBadRequest)
		return
	}
	cmd.ID = id

	// проверяем что хоть одно поле передано
	if cmd.Phone == nil && cmd.FullName == nil && cmd.BirthDate == nil {
		utils.JSONResponse(w, map[string]string{"error": "no fields to update"}, http.StatusBadRequest)
		return
	}

	userDTO, err := a.handler.Handle(cmd)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, userDTO, http.StatusOK)
}
