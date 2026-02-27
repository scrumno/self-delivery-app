package v1

import (
	"github.com/gorilla/mux"
	"log"
	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
	"github.com/scrumno/scrumno-api/internal/api/v1/collector"
)

// SetupRouter создаёт маршруты
func SetupRouter(cfg *config.Config, actions *action.Actions) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	api := router.PathPrefix("/api/v1").Subrouter()

	healthPrefix := "/health"

	health := api.PathPrefix(healthPrefix).Subrouter()

	collectorRoutes := collector.NewEndpointCollector()

	collectorRoutes.HandleFuncWithPostman(
        health,
		healthPrefix,
        actions.CheckStatusConnectDB.Action,
        actions.CheckStatusConnectDB.GetInputType(),
        "GET",
        "/check-status-connect-db",
    )

	collectorRoutes.HandleFuncWithPostman(
        health,
		healthPrefix,
        actions.CheckStatusConnectDB.Action,
        actions.CheckStatusConnectDB.GetInputType(),
        "POST",
        "/check-1221",
    )

	userPrefix := "/users"

	user := api.PathPrefix(userPrefix).Subrouter()

	collectorRoutes.HandleFuncWithPostman(
        user,
		userPrefix,
        actions.CheckStatusConnectDB.Action,
        actions.CheckStatusConnectDB.GetInputType(),
        "GET",
        "/users",
    )

	collectorRoutes.HandleFuncWithPostman(
        user,
		userPrefix,
        actions.CheckStatusConnectDB.Action,
        actions.CheckStatusConnectDB.GetInputType(),
        "GET",
        "/usersusersuser/us",
    )

	err := collectorRoutes.GeneratePostmanCollections()
	if err != nil {
        log.Printf("Ошибка генерации Postman: %v", err)
    }

	return router
}