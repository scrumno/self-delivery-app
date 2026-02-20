package user

import (
	"github.com/scrumno/scrumno-api/internal/api/utils"
	getAllUsers "github.com/scrumno/scrumno-api/internal/users/query/get-all-users"
	"net/http"
)

type GetAllUsersAction struct {
	fetcher *getAllUsers.Fetcher
}

func NewGetAllUsersAction(fetcher *getAllUsers.Fetcher) *GetAllUsersAction {
	return &GetAllUsersAction{fetcher: fetcher}
}

func (a *GetAllUsersAction) Action(w http.ResponseWriter, _ *http.Request) {
	usersDTO, err := a.fetcher.Fetch(getAllUsers.Query{})
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, usersDTO, http.StatusOK)
}
