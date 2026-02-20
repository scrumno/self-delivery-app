package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	getUserByID "github.com/scrumno/scrumno-api/internal/users/query/get-user-by-id"
)

type GetUserByIDAction struct {
	fetcher *getUserByID.Fetcher
}

func NewGetUserByIDAction(fetcher *getUserByID.Fetcher) *GetUserByIDAction {
	return &GetUserByIDAction{fetcher: fetcher}
}

func (a *GetUserByIDAction) Action(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.JSONResponse(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	query := getUserByID.Query{ID: id}
	userDTO, err := a.fetcher.Fetch(query)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, userDTO, http.StatusOK)
}
