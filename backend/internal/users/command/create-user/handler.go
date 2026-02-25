package create_user

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/users/entity/user"
	"github.com/scrumno/scrumno-api/shared/base"
)

type CreateUserHandler struct {
	repo base.BaseRepository[user.User]
}

func NewCreateUserHandler(repo base.BaseRepository[user.User]) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

func (handler *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (UserDTO, error) {
	newUser := user.NewUser(cmd.Phone)

	if err := handler.repo.Create(ctx, newUser); err != nil {
		return UserDTO{}, err
	}

	return UserDTO{
		ID:        newUser.ID,
		Phone:     newUser.Phone,
		CreatedAt: newUser.CreatedAt,
		BirthDate: newUser.BirthDate,
		IsActive:  newUser.IsActive,
		FullName:  newUser.FullName,
	}, nil
}
