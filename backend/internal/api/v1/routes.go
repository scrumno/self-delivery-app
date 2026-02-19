package v1

import (
	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
)

// SetupRouter создаёт маршруты. Конфиг передаётся для инъекции JWT secret (DIP).
func SetupRouter(cfg *config.Config, actions *action.Actions) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	api := router.PathPrefix("/api/v1").Subrouter()

	health := api.PathPrefix("/health").Subrouter()
	health.HandleFunc("/check-status-connect-db", actions.CheckStatusConnectDB.Action).Methods("GET")

	userRouter := api.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", actions.GetAllUsers.Action).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9a-fA-F-]+}", actions.GetUserByID.Action).Methods("GET")
	userRouter.HandleFunc("/search", actions.GetUserByPhone.Action).Methods("GET")

	userRouter.HandleFunc("", actions.CreateUser.Action).Methods("POST")
	userRouter.HandleFunc("/{id:[0-9a-fA-F-]+}", actions.UpdateUser.Action).Methods("PATCH")
	userRouter.HandleFunc("/{id:[0-9a-fA-F-]+}", actions.DeleteUser.Action).Methods("DELETE")

	return router
}
