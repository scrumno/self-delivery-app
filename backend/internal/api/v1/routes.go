package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
	"github.com/scrumno/scrumno-api/internal/health/query"
	"github.com/scrumno/scrumno-api/internal/users/auth/command"
	userQuery "github.com/scrumno/scrumno-api/internal/users/user/query"
)

// SetupRouter создаёт маршруты. Конфиг передаётся для инъекции JWT secret (DIP).
func SetupRouter(cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	api := router.PathPrefix("/api/v1").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", command.LoginHandler(cfg.JWT.SecretKey)).Methods(http.MethodPost)
	auth.HandleFunc("/signup", command.SignUpHandler(cfg.JWT.SecretKey)).Methods(http.MethodPost)
	auth.HandleFunc("/refresh", command.RefreshHandler(cfg.JWT.SecretKey)).Methods(http.MethodPost)

	healthProtected := api.PathPrefix("/health").Subrouter()
	healthProtected.Use(middleware.JWT(cfg.JWT.SecretKey))
	healthProtected.HandleFunc("/check-status-db-connect", query.CheckStatusConnectDB)

	users := api.PathPrefix("/users").Subrouter()
	users.Use(middleware.JWT(cfg.JWT.SecretKey))
	users.HandleFunc("", userQuery.GetAllUsers).Methods(http.MethodGet)

	usersProtected := api.PathPrefix("/users").Subrouter()
	usersProtected.Use(middleware.JWT(cfg.JWT.SecretKey))
	usersProtected.HandleFunc("", userQuery.GetAllUsers).Methods(http.MethodGet)
	usersProtected.HandleFunc("/{id}", userQuery.GetUserByID).Methods(http.MethodGet)

	return router
}
