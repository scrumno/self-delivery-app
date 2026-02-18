package action

import userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"

type Actions struct {
	// users
	GetUserByID *userAction.GetUserByIDAction
}
