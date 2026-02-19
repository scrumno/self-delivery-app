package user

import (
	"net/http"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	getUserByPhone "github.com/scrumno/scrumno-api/internal/users/query/get-user-by-phone"
)

type GetUserByPhoneAction struct {
	fetcher *getUserByPhone.Fetcher
}

func NewGetUserByPhoneAction(fetcher *getUserByPhone.Fetcher) *GetUserByPhoneAction {
	return &GetUserByPhoneAction{fetcher: fetcher}
}

func (a *GetUserByPhoneAction) Action(w http.ResponseWriter, r *http.Request) {

	phone := r.URL.Query().Get("phone")

	if phone == "" {
		utils.JSONResponse(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	query := getUserByPhone.Query{Phone: phone}
	userDTO, err := a.fetcher.Fetch(query)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, userDTO, http.StatusOK)
}
