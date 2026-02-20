package config

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	"github.com/scrumno/scrumno-api/internal/health/entity/status"
	checkStatusConnectDB "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
	createUser "github.com/scrumno/scrumno-api/internal/users/command/create-user"
	deleteUser "github.com/scrumno/scrumno-api/internal/users/command/delete-user"
	updateUserById "github.com/scrumno/scrumno-api/internal/users/command/update-user-by-id"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
	getAllUsers "github.com/scrumno/scrumno-api/internal/users/query/get-all-users"
	getUserByID "github.com/scrumno/scrumno-api/internal/users/query/get-user-by-id"
	getUserByPhone "github.com/scrumno/scrumno-api/internal/users/query/get-user-by-phone"
)

func DI() *action.Actions {
	// repository
	userRepo := user.NewUserRepository(DB)
	statusRepo := status.NewStatusRepository(DB) // добавил
	// service

	// command
	updateUserHandler := updateUserById.NewHandler(userRepo)
	deleteUserHandler := deleteUser.NewHandler(userRepo) // добавил
	createUserHandler := createUser.NewHandler(userRepo) // добавил

	// query
	getUserByIDFetcher := getUserByID.NewFetcher(userRepo)
	getUserByPhoneFetcher := getUserByPhone.NewFetcher(userRepo)
	getAllUsersFetcher := getAllUsers.NewFetcher(userRepo)            //добавил
	checkStatusFetcher := checkStatusConnectDB.NewFetcher(statusRepo) //добавил

	return &action.Actions{
		CheckStatusConnectDB: healthAction.NewCheckStatusConnectDBAction(checkStatusFetcher), //добавил
		GetUserByID:          userAction.NewGetUserByIDAction(getUserByIDFetcher),
		GetUserByPhone:       userAction.NewGetUserByPhoneAction(getUserByPhoneFetcher),
		GetAllUsers:          userAction.NewGetAllUsersAction(getAllUsersFetcher), // добавил
		CreateUser:           userAction.NewCreateUserAction(createUserHandler),   // добавил
		UpdateUserById:       userAction.NewUpdateUserByIdAction(updateUserHandler),
		DeleteUser:           userAction.NewDeleteUserAction(deleteUserHandler), //добавил
	}
}
