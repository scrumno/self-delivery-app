package config

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
	getUserByID "github.com/scrumno/scrumno-api/internal/users/query/get-user-by-id"
)

func DI() *action.Actions {
	// repository
	userRepo := user.NewUserRepository(DB)

	// service

	// command

	// query
	getUserByIDFetcher := getUserByID.NewFetcher(userRepo)

	return &action.Actions{
		GetUserByID: userAction.NewGetUserByIDAction(getUserByIDFetcher),
	}
}
