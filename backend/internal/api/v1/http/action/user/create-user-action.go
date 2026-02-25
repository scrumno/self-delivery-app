package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	createUser "github.com/scrumno/scrumno-api/internal/users/command/create-user"
)

type CreateUserAction struct {
	Handler *createUser.CreateUserHandler
}

func NewCreateUserAction(handler *createUser.CreateUserHandler) *CreateUserAction {
	return &CreateUserAction{
		Handler: handler,
	}
}

type ErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}

type SuccessResponse struct {
	IsSuccess bool               `json:"isSuccess"`
	User      createUser.UserDTO `json:"user"`
}

func (a *CreateUserAction) Action(w http.ResponseWriter, r *http.Request) {
	phone := mux.Vars(r)["phone"]
	if phone == "" {
		utils.JSONResponse(w, ErrorResponse{false, "phone is empty"}, http.StatusBadRequest)
		return
	}

	cmd := createUser.CreateUserCommand{
		Phone: phone,
	}

	dto, err := a.Handler.Handle(r.Context(), cmd)
	if err != nil {
		utils.JSONResponse(w, ErrorResponse{false, err.Error()}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, SuccessResponse{true, dto}, http.StatusOK)
}
