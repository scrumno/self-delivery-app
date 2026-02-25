package config

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	"github.com/scrumno/scrumno-api/internal/health/entity/status"
	checkStatusConnectDB "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
	createUser "github.com/scrumno/scrumno-api/internal/users/command/create-user"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
	"github.com/scrumno/scrumno-api/shared/factory"
)

func DI() *action.Actions {
	// repository
	statusRepo := status.NewStatusRepository(DB)
	userRepo := factory.NewGormRepository[user.User](DB)

	// service
	checkStatusFetcher := checkStatusConnectDB.NewFetcher(statusRepo)

	// command
	createUserHandler := createUser.NewCreateUserHandler(userRepo)

	// query

	return &action.Actions{
		CheckStatusConnectDB: healthAction.NewCheckStatusConnectDBAction(checkStatusFetcher),

		// users
		CreateUser: userAction.NewCreateUserAction(createUserHandler),
	}
}
