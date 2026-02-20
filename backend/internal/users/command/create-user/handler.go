package create_user

import (
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

type Handler struct {
	repository user.UserRepositoryInterface
}

func NewHandler(repository user.UserRepositoryInterface) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Handle(cmd Command) (UserDTO, error) {
	u := user.User{
		Phone:    cmd.Phone,
		FullName: cmd.FullName,
	}

	created, err := h.repository.Create(u)
	if err != nil {
		return UserDTO{}, err
	}

	return UserDTO{
		ID:        created.ID,
		Phone:     created.Phone,
		FullName:  created.FullName,
		IsActive:  created.IsActive,
		CreatedAt: created.CreatedAt,
	}, nil
}
