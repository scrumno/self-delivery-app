package user

import (
	"encoding/json"
	"net/http"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	createUser "github.com/scrumno/scrumno-api/internal/users/command/create-user"
)

type CreateUserAction struct {
	handler *createUser.Handler
}

func NewCreateUserAction(handler *createUser.Handler) *CreateUserAction {
	return &CreateUserAction{handler: handler}
}

func (a *CreateUserAction) Action(w http.ResponseWriter, r *http.Request) {
	var cmd createUser.Command
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		utils.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusUnprocessableEntity)
		return
	}

	if cmd.Phone == "" {
		utils.JSONResponse(w, map[string]string{"error": "phone is required"}, http.StatusUnprocessableEntity)
		return
	}

	userDTO, err := a.handler.Handle(cmd)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, userDTO, http.StatusCreated)
}
