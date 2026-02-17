package v1

import (
	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/config"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
)

// SetupRouter создаёт маршруты. Конфиг передаётся для инъекции JWT secret (DIP).
func SetupRouter(cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	api := router.PathPrefix("/api/v1").Subrouter()

	health := api.PathPrefix("/health").Subrouter()
	health.HandleFunc("/check-status-db-connect", healthAction.CheckStatusConnectDBAction).Methods("GET")

	return router
}
