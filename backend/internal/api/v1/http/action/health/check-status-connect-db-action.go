package health

import (
	"net/http"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	checkStatusConnectDB "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
)

type response struct {
	IsOk bool `json:"isOk"`
}

func CheckStatusConnectDBAction(responseWriter http.ResponseWriter, request *http.Request) {
	res := checkStatusConnectDB.Fetcher()

	if !res.Status {
		utils.JSONResponse(responseWriter, response{IsOk: false}, http.StatusInternalServerError)
	}

	utils.JSONResponse(responseWriter, response{IsOk: true}, http.StatusOK)
}
