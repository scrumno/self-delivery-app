package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	deleteUser "github.com/scrumno/scrumno-api/internal/users/command/delete-user"
)

type DeleteUserAction struct {
	handler *deleteUser.Handler
}

func NewDeleteUserAction(handler *deleteUser.Handler) *DeleteUserAction {
	return &DeleteUserAction{handler: handler}
}

func (a *DeleteUserAction) Action(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		utils.JSONResponse(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	err := a.handler.Handle(deleteUser.Command{ID: id})
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, map[string]string{"message": "User deleted"}, http.StatusOK)
}
