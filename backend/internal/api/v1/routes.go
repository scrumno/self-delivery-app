package v1

import (
	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
)

// SetupRouter создаёт маршруты
func SetupRouter(cfg *config.Config, actions *action.Actions) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	api := router.PathPrefix("/api/v1").Subrouter()

	health := api.PathPrefix("/health").Subrouter()
	health.HandleFunc("/check-status-connect-db", actions.CheckStatusConnectDB.Action).Methods("GET")

	return router
}
