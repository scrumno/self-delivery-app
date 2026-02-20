package action

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
)

type Actions struct {
	// db
	CheckStatusConnectDB *health.CheckStatusConnectDBAction

	// users
	GetUserByID    *userAction.GetUserByIDAction
	GetUserByPhone *userAction.GetUserByPhoneAction
	GetAllUsers    *userAction.GetAllUsersAction
	CreateUser     *userAction.CreateUserAction
	UpdateUserById *userAction.UpdateUserByIdAction
	DeleteUser     *userAction.DeleteUserAction
}
