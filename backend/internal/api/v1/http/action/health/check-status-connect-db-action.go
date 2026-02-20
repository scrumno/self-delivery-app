package health

import (
	"net/http"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	checkStatus "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
)

type CheckStatusConnectDBAction struct {
	fetcher *checkStatus.Fetcher
}

func NewCheckStatusConnectDBAction(fetcher *checkStatus.Fetcher) *CheckStatusConnectDBAction {
	return &CheckStatusConnectDBAction{fetcher: fetcher}
}

func (a *CheckStatusConnectDBAction) Action(w http.ResponseWriter, _ *http.Request) {
	dto := a.fetcher.Fetch(checkStatus.Query{})

	if !dto.IsConnected {
		utils.JSONResponse(w, map[string]bool{"isOk": false}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]bool{"isOk": true}, http.StatusOK)
}
